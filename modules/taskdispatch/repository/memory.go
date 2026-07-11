package repository

import (
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/modules/taskdispatch/domain"
)

// MemTaskDispatchRepo is the in-memory task dispatch engine.
// It manages a task pool with capability-based operator matching,
// timeout detection, and automatic reassignment.
type MemTaskDispatchRepo struct {
	mu         sync.RWMutex
	tasks      map[int64]*domain.WarehouseTask
	operators  map[int64]*domain.OperatorCapability
	nextTaskID int64
	nextOpID   int64
}

// NewMemTaskDispatchRepo creates a new repository with seed data.
func NewMemTaskDispatchRepo() *MemTaskDispatchRepo {
	r := &MemTaskDispatchRepo{
		tasks:     make(map[int64]*domain.WarehouseTask),
		operators: make(map[int64]*domain.OperatorCapability),
		nextTaskID: 1,
		nextOpID:   1,
	}
	r.seedOperators()
	r.seedTasks()
	return r
}

// ─── Seeding ──────────────────────────────────────────────────────────

func (r *MemTaskDispatchRepo) seedOperators() {
	r.operators[1] = &domain.OperatorCapability{
		OperatorID: 1, OperatorName: "张操作员",
		Capabilities: []string{domain.CapForklift, domain.CapHeavyLift},
		IsOnline: true, CurrentTaskCount: 0, WarehouseID: 1,
	}
	r.operators[2] = &domain.OperatorCapability{
		OperatorID: 2, OperatorName: "李操作员",
		Capabilities: []string{domain.CapFragile, domain.CapEcommerce},
		IsOnline: true, CurrentTaskCount: 0, WarehouseID: 1,
	}
	r.operators[3] = &domain.OperatorCapability{
		OperatorID: 3, OperatorName: "王操作员",
		Capabilities: []string{domain.CapHazmat, domain.CapColdChain, domain.CapForklift},
		IsOnline: true, CurrentTaskCount: 0, WarehouseID: 1,
	}
	r.operators[4] = &domain.OperatorCapability{
		OperatorID: 4, OperatorName: "赵操作员",
		Capabilities: []string{domain.CapHeavyLift, domain.CapFragile, domain.CapEcommerce},
		IsOnline: true, CurrentTaskCount: 0, WarehouseID: 1,
	}
	r.nextOpID = 5
}

func (r *MemTaskDispatchRepo) seedTasks() {
	now := time.Now()
	parcel1 := int64(1)
	parcel2 := int64(2)
	order1 := int64(1)

	tasks := []*domain.WarehouseTask{
		{
			TaskCode: "TASK-001", TaskType: domain.TaskTypeReceive,
			ParcelID: &parcel1, ParcelTrackingNumber: "SF1234567890",
			WarehouseID: 1, LocationCode: "A-01-01",
			Status: domain.StatusPending,
			RequiredCapabilities: []string{},
			TimeoutMinutes: 30,
		},
		{
			TaskCode: "TASK-002", TaskType: domain.TaskTypePutaway,
			ParcelID: &parcel1, ParcelTrackingNumber: "SF1234567890",
			WarehouseID: 1, LocationCode: "A-02-03",
			Status: domain.StatusPending,
			RequiredCapabilities: []string{domain.CapForklift},
			TimeoutMinutes: 30,
		},
		{
			TaskCode: "TASK-003", TaskType: domain.TaskTypeWeigh,
			ParcelID: &parcel2, ParcelTrackingNumber: "ZTO9876543210",
			WarehouseID: 1, LocationCode: "B-01-02",
			Status: domain.StatusPending,
			RequiredCapabilities: []string{},
			TimeoutMinutes: 30,
		},
		{
			TaskCode: "TASK-004", TaskType: domain.TaskTypePick,
			OrderID: &order1,
			WarehouseID: 1, LocationCode: "C-03-01",
			Status: domain.StatusPending,
			RequiredCapabilities: []string{domain.CapFragile},
			TimeoutMinutes: 30,
		},
		{
			TaskCode: "TASK-005", TaskType: domain.TaskTypePack,
			OrderID: &order1,
			WarehouseID: 1, LocationCode: "PACK-01",
			Status: domain.StatusPending,
			RequiredCapabilities: []string{domain.CapEcommerce},
			TimeoutMinutes: 30,
		},
		{
			TaskCode: "TASK-006", TaskType: domain.TaskTypeWeightCheck,
			ParcelID: &parcel2, ParcelTrackingNumber: "ZTO9876543210",
			WarehouseID: 1, LocationCode: "WEIGH-02",
			Status: domain.StatusPending,
			RequiredCapabilities: []string{},
			TimeoutMinutes: 15,
		},
		{
			TaskCode: "TASK-007", TaskType: domain.TaskTypeException,
			ParcelID: &parcel1, ParcelTrackingNumber: "SF1234567890",
			WarehouseID: 1, LocationCode: "EXCEPTION-ZONE",
			Status: domain.StatusPending,
			RequiredCapabilities: []string{domain.CapHazmat},
			TimeoutMinutes: 45,
		},
		{
			TaskCode: "TASK-008", TaskType: domain.TaskTypeLoad,
			OrderID: &order1,
			WarehouseID: 1, LocationCode: "DOCK-03",
			Status: domain.StatusPending,
			RequiredCapabilities: []string{domain.CapHeavyLift, domain.CapForklift},
			TimeoutMinutes: 60,
		},
	}

	for _, t := range tasks {
		id := atomic.AddInt64(&r.nextTaskID, 1) - 1
		t.ID = id
		t.CreatedAt = now
		t.UpdatedAt = now
		r.tasks[id] = t
	}
}

// ─── Task Pool ────────────────────────────────────────────────────────

// TaskPool returns all pending (unclaimed) tasks — the 抢单池.
func (r *MemTaskDispatchRepo) TaskPool() []*domain.WarehouseTask {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*domain.WarehouseTask
	for _, t := range r.tasks {
		if t.Status == domain.StatusPending {
			result = append(result, t)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})
	return result
}

// GetTaskByID returns a single task by ID.
func (r *MemTaskDispatchRepo) GetTaskByID(taskID int64) (*domain.WarehouseTask, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tasks[taskID]
	if !ok {
		return nil, fmt.Errorf("task %d not found", taskID)
	}
	return t, nil
}

// ─── Task Lifecycle ───────────────────────────────────────────────────

// ClaimTask moves a task from pending → claimed and assigns it to an operator.
func (r *MemTaskDispatchRepo) ClaimTask(taskID, operatorID int64) (*domain.WarehouseTask, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	t, ok := r.tasks[taskID]
	if !ok {
		return nil, fmt.Errorf("任务不存在")
	}
	if t.Status != domain.StatusPending {
		return nil, fmt.Errorf("任务状态为 %s，无法认领（仅pending可认领）", domain.StatusDisplay(t.Status))
	}

	op, ok := r.operators[operatorID]
	if !ok || !op.IsOnline {
		return nil, fmt.Errorf("操作员不在线或不存在")
	}

	now := time.Now()
	t.Status = domain.StatusClaimed
	t.AssignedOperatorID = &operatorID
	t.AssignedAt = &now
	t.UpdatedAt = now
	op.CurrentTaskCount++

	return t, nil
}

// StartTask moves a task from claimed → in_progress.
func (r *MemTaskDispatchRepo) StartTask(taskID int64) (*domain.WarehouseTask, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	t, ok := r.tasks[taskID]
	if !ok {
		return nil, fmt.Errorf("任务不存在")
	}
	if t.Status != domain.StatusClaimed {
		return nil, fmt.Errorf("任务状态为 %s，无法开始（需先认领）", domain.StatusDisplay(t.Status))
	}

	now := time.Now()
	t.Status = domain.StatusInProgress
	t.UpdatedAt = now
	return t, nil
}

// CompleteTask moves a task to completed.
func (r *MemTaskDispatchRepo) CompleteTask(taskID int64) (*domain.WarehouseTask, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	t, ok := r.tasks[taskID]
	if !ok {
		return nil, fmt.Errorf("任务不存在")
	}
	if t.Status != domain.StatusInProgress && t.Status != domain.StatusClaimed {
		return nil, fmt.Errorf("任务状态为 %s，无法完成", domain.StatusDisplay(t.Status))
	}

	now := time.Now()
	t.Status = domain.StatusCompleted
	t.UpdatedAt = now

	// Decrement operator task count
	if t.AssignedOperatorID != nil {
		if op, ok := r.operators[*t.AssignedOperatorID]; ok && op.CurrentTaskCount > 0 {
			op.CurrentTaskCount--
		}
	}

	return t, nil
}

// ─── Timeout Handling ─────────────────────────────────────────────────

// CheckTimeouts scans claimed/in-progress tasks for SLA breaches.
// Timed-out tasks are reset to pending and their operator assignment is cleared.
// Returns the list of tasks that were reassigned.
func (r *MemTaskDispatchRepo) CheckTimeouts() []*domain.WarehouseTask {
	r.mu.Lock()
	defer r.mu.Unlock()

	var timedOut []*domain.WarehouseTask
	now := time.Now()

	for _, t := range r.tasks {
		if t.Status != domain.StatusClaimed && t.Status != domain.StatusInProgress {
			continue
		}
		if !t.IsTimedOut() {
			continue
		}

		// Decrement operator task count
		if t.AssignedOperatorID != nil {
			if op, ok := r.operators[*t.AssignedOperatorID]; ok && op.CurrentTaskCount > 0 {
				op.CurrentTaskCount--
			}
		}

		// Reset to pending — back in the pool
		t.Status = domain.StatusTimeout // mark timeout first for audit trail
		t.AssignedOperatorID = nil
		t.AssignedAt = nil
		t.UpdatedAt = now

		// Reassign: try to match a new operator automatically
		bestOp := r.matchOperatorLocked(t)
		if bestOp != nil {
			t.Status = domain.StatusPending // make claimable again
			// Auto-claim with best match
			t.Status = domain.StatusClaimed
			t.AssignedOperatorID = &bestOp.OperatorID
			t.AssignedAt = &now
			bestOp.CurrentTaskCount++
		} else {
			// No operator available, return to pool
			t.Status = domain.StatusPending
		}
		t.UpdatedAt = now
		timedOut = append(timedOut, t)
	}

	return timedOut
}

// ─── Operator Matching ────────────────────────────────────────────────

// MatchOperator finds the best operator for a task based on capabilities.
// Returns the best match or nil if no online operator qualifies.
func (r *MemTaskDispatchRepo) MatchOperator(task *domain.WarehouseTask) *domain.OperatorCapability {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.matchOperatorLocked(task)
}

// matchOperatorLocked is the internal matching logic (caller must hold lock).
func (r *MemTaskDispatchRepo) matchOperatorLocked(task *domain.WarehouseTask) *domain.OperatorCapability {
	if task.RequiredCapabilities == nil {
		task.RequiredCapabilities = []string{}
	}

	type candidate struct {
		op    *domain.OperatorCapability
		score int
	}
	var candidates []candidate

	for _, op := range r.operators {
		if !op.IsOnline {
			continue
		}
		if op.WarehouseID != task.WarehouseID {
			continue
		}
		score := op.MatchScore(task.RequiredCapabilities)
		// Must match all required capabilities
		if len(task.RequiredCapabilities) > 0 && score < len(task.RequiredCapabilities) {
			continue
		}
		candidates = append(candidates, candidate{op: op, score: score})
	}

	if len(candidates) == 0 {
		return nil
	}

	// Sort by: highest score, then lowest task count (load balance)
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].score != candidates[j].score {
			return candidates[i].score > candidates[j].score
		}
		return candidates[i].op.CurrentTaskCount < candidates[j].op.CurrentTaskCount
	})

	return candidates[0].op
}

// ─── Operator Queries ─────────────────────────────────────────────────

// GetOperators returns all operators.
func (r *MemTaskDispatchRepo) GetOperators() []*domain.OperatorCapability {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*domain.OperatorCapability
	for _, op := range r.operators {
		result = append(result, op)
	}
	return result
}

// GetOperatorByID returns an operator by ID.
func (r *MemTaskDispatchRepo) GetOperatorByID(operatorID int64) (*domain.OperatorCapability, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	op, ok := r.operators[operatorID]
	if !ok {
		return nil, fmt.Errorf("操作员 %d 不存在", operatorID)
	}
	return op, nil
}

// GetOperatorTasks returns all current tasks for a specific operator
// (claimed or in-progress, not completed/cancelled).
func (r *MemTaskDispatchRepo) GetOperatorTasks(operatorID int64) []*domain.WarehouseTask {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*domain.WarehouseTask
	for _, t := range r.tasks {
		if t.AssignedOperatorID != nil && *t.AssignedOperatorID == operatorID {
			if t.Status == domain.StatusClaimed || t.Status == domain.StatusInProgress {
				result = append(result, t)
			}
		}
	}
	return result
}
