package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/modules/pda/domain"
)

// MemWorkflowProcessRepo is an in-memory implementation for workflow processes.
type MemWorkflowProcessRepo struct {
	mu       sync.RWMutex
	processes map[int64]*domain.WorkflowProcess
	nextID   int64
}

func NewMemWorkflowProcessRepo() *MemWorkflowProcessRepo {
	r := &MemWorkflowProcessRepo{processes: make(map[int64]*domain.WorkflowProcess), nextID: 1}
	r.seed()
	return r
}

func (r *MemWorkflowProcessRepo) seed() {
	now := time.Now()
	seeds := []struct {
		workOrderID int64
		stepName    string
		stepOrder   int
		status      domain.WorkflowProcessStatus
		operatorID  int64
		resultNote  string
	}{
		{1, "scan_tracking", 1, domain.WFPStatusCompleted, 1, "已扫描运单"},
		{1, "weigh", 2, domain.WFPStatusCompleted, 1, "包裹重量0.48kg"},
		{1, "verify_weight", 3, domain.WFPStatusRunning, 1, ""},
		{1, "putaway", 4, domain.WFPStatusPending, 0, ""},
		{3, "scan_parcel", 1, domain.WFPStatusCompleted, 3, "已完成扫描"},
		{3, "pack", 2, domain.WFPStatusCompleted, 3, "打包完成"},
		{3, "confirm_pack", 3, domain.WFPStatusCompleted, 3, "复核无误"},
		{3, "outbound", 4, domain.WFPStatusCompleted, 3, "出库完成"},
	}
	for _, s := range seeds {
		id := atomic.AddInt64(&r.nextID, 1) - 1
		wp := &domain.WorkflowProcess{
			ID:          id,
			WorkOrderID: s.workOrderID,
			StepName:    s.stepName,
			StepOrder:   s.stepOrder,
			Status:      s.status,
			OperatorID:  s.operatorID,
			ResultNote:  s.resultNote,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		if s.status == domain.WFPStatusCompleted {
			wp.CompletedAt = &now
		}
		if s.status == domain.WFPStatusRunning {
			wp.StartedAt = &now
		}
		r.processes[id] = wp
	}
}

func (r *MemWorkflowProcessRepo) Create(ctx context.Context, wp *domain.WorkflowProcess) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	wp.ID = atomic.AddInt64(&r.nextID, 1) - 1
	now := time.Now()
	wp.CreatedAt = now
	wp.UpdatedAt = now
	r.processes[wp.ID] = wp
	return nil
}

func (r *MemWorkflowProcessRepo) GetByID(ctx context.Context, id int64) (*domain.WorkflowProcess, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	wp, ok := r.processes[id]
	if !ok {
		return nil, nil
	}
	return wp, nil
}

func (r *MemWorkflowProcessRepo) ListByWorkOrder(ctx context.Context, workOrderID int64) ([]domain.WorkflowProcess, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.WorkflowProcess
	for _, wp := range r.processes {
		if wp.WorkOrderID == workOrderID {
			result = append(result, *wp)
		}
	}
	return result, nil
}

func (r *MemWorkflowProcessRepo) ListByOperator(ctx context.Context, operatorID int64) ([]domain.WorkflowProcess, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.WorkflowProcess
	for _, wp := range r.processes {
		if wp.OperatorID == operatorID {
			result = append(result, *wp)
		}
	}
	return result, nil
}

func (r *MemWorkflowProcessRepo) Update(ctx context.Context, wp *domain.WorkflowProcess) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.processes[wp.ID]; !ok {
		return nil
	}
	wp.UpdatedAt = time.Now()
	r.processes[wp.ID] = wp
	return nil
}

var _ WorkflowProcessRepository = (*MemWorkflowProcessRepo)(nil)
