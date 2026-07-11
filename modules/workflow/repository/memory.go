package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/modules/workflow/domain"
)

// ===================================================================
// MemWorkflowRepo — in-memory workflow engine with real seed pipelines
// ===================================================================

type MemWorkflowRepo struct {
	mu        sync.RWMutex
	processes map[int64]*domain.WorkflowProcess
	orders    map[int64]*domain.WorkOrder
	nextPID   int64
	nextWOID  int64
}

func NewMemWorkflowRepo() *MemWorkflowRepo {
	r := &MemWorkflowRepo{
		processes: make(map[int64]*domain.WorkflowProcess),
		orders:    make(map[int64]*domain.WorkOrder),
		nextPID:   1,
		nextWOID:  1,
	}
	r.seed()
	return r
}

// ===================================================================
// Seed 2 real pipelines with full step definitions
// ===================================================================

func (r *MemWorkflowRepo) seed() {
	now := time.Now()
	tenant := int64(1)

	// ─── Pipeline 1: 标准入库流程 (Inbound) ──────────────────────────
	// Steps: 收货确认→称重测量→上架入库→完成
	inbound := &domain.WorkflowProcess{
		ID:           1,
		TenantID:     tenant,
		Name:         "标准入库流程",
		Code:         "inbound",
		TriggerEvent: domain.TriggerParcelReceived,
		IsActive:     true,
		Steps: []domain.WorkflowStep{
			{ID: 1, ProcessID: 1, Name: domain.StepReceiveConfirm, DisplayName: "收货确认", OrderIndex: 1, Required: true, Assignable: true, TimeoutMinutes: 60},
			{ID: 2, ProcessID: 1, Name: domain.StepWeighMeasure, DisplayName: "称重测量", OrderIndex: 2, Required: true, Assignable: true, TimeoutMinutes: 30},
			{ID: 3, ProcessID: 1, Name: domain.StepPutaway, DisplayName: "上架入库", OrderIndex: 3, Required: true, Assignable: true, TimeoutMinutes: 120},
			{ID: 4, ProcessID: 1, Name: domain.StepComplete, DisplayName: "完成", OrderIndex: 4, Required: true, Assignable: false, TimeoutMinutes: 0},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	r.processes[1] = inbound
	r.nextPID = 2

	// ─── Pipeline 2: 标准出库流程 (Outbound) ─────────────────────────
	// Steps: 拣货→送打包→打包→核重→送出库→送装柜→装柜→完成
	outbound := &domain.WorkflowProcess{
		ID:           2,
		TenantID:     tenant,
		Name:         "标准出库流程",
		Code:         "outbound",
		TriggerEvent: domain.TriggerOrderCreated,
		IsActive:     true,
		Steps: []domain.WorkflowStep{
			{ID: 5, ProcessID: 2, Name: domain.StepPick, DisplayName: "拣货", OrderIndex: 1, Required: true, Assignable: true, TimeoutMinutes: 120},
			{ID: 6, ProcessID: 2, Name: domain.StepSendToPack, DisplayName: "送打包", OrderIndex: 2, Required: true, Assignable: true, TimeoutMinutes: 30},
			{ID: 7, ProcessID: 2, Name: domain.StepPack, DisplayName: "打包", OrderIndex: 3, Required: true, Assignable: true, TimeoutMinutes: 60},
			{ID: 8, ProcessID: 2, Name: domain.StepWeightCheck, DisplayName: "核重", OrderIndex: 4, Required: true, Assignable: true, TimeoutMinutes: 30},
			{ID: 9, ProcessID: 2, Name: domain.StepSendOut, DisplayName: "送出库", OrderIndex: 5, Required: true, Assignable: true, TimeoutMinutes: 30},
			{ID: 10, ProcessID: 2, Name: domain.StepSendToContainer, DisplayName: "送装柜", OrderIndex: 6, Required: true, Assignable: true, TimeoutMinutes: 30},
			{ID: 11, ProcessID: 2, Name: domain.StepLoadContainer, DisplayName: "装柜", OrderIndex: 7, Required: true, Assignable: true, TimeoutMinutes: 120},
			{ID: 12, ProcessID: 2, Name: domain.StepComplete, DisplayName: "完成", OrderIndex: 8, Required: true, Assignable: false, TimeoutMinutes: 0},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	r.processes[2] = outbound
	r.nextPID = 3

	// ─── Seed some sample work orders ─────────────────────────────────
	user1 := int64(1) // 大宝
	user2 := int64(2) // 安冉
	parcel1 := int64(1)
	parcel2 := int64(2)
	order1 := int64(1)
	order2 := int64(2)

	// WO-1: 入库 — 收货确认→称重测量 (at step 2)
	r.orders[1] = &domain.WorkOrder{
		ID: 1, TenantID: tenant, WarehouseID: 1,
		ProcessID: 1, ProcessName: "标准入库流程",
		ParcelID: &parcel1, Status: domain.WOStatusInProgress,
		CurrentStep: 2, AssignedTo: &user1, AssignedName: "大宝",
		Title: "SF1234567890 入库", Description: "顺丰包裹入库作业",
		Priority: 0,
		CreatedAt: now.Add(-2 * time.Hour), UpdatedAt: now.Add(-1 * time.Hour),
	}

	// WO-2: 入库 — 已完成
	completedTime := now.Add(-30 * time.Minute)
	r.orders[2] = &domain.WorkOrder{
		ID: 2, TenantID: tenant, WarehouseID: 1,
		ProcessID: 1, ProcessName: "标准入库流程",
		ParcelID: &parcel2, Status: domain.WOStatusCompleted,
		CurrentStep: 4, AssignedTo: &user1, AssignedName: "大宝",
		Title: "ZTO9876543210 入库", Description: "中通包裹入库作业",
		Priority: 0,
		CreatedAt: now.Add(-4 * time.Hour), UpdatedAt: now.Add(-30 * time.Minute),
		CompletedAt: &completedTime,
	}

	// WO-3: 出库 — 拣货→送打包→打包 (at step 3)
	r.orders[3] = &domain.WorkOrder{
		ID: 3, TenantID: tenant, WarehouseID: 1,
		ProcessID: 2, ProcessName: "标准出库流程",
		OrderID: &order1, Status: domain.WOStatusInProgress,
		CurrentStep: 3, AssignedTo: &user2, AssignedName: "安冉",
		Title: "订单80020737681100020001 出库", Description: "王仁照的集运出库",
		Priority: 0,
		CreatedAt: now.Add(-3 * time.Hour), UpdatedAt: now.Add(-45 * time.Minute),
	}

	// WO-4: 出库 — 待处理
	r.orders[4] = &domain.WorkOrder{
		ID: 4, TenantID: tenant, WarehouseID: 1,
		ProcessID: 2, ProcessName: "标准出库流程",
		OrderID: &order2, Status: domain.WOStatusPending,
		CurrentStep: 1, AssignedTo: nil, AssignedName: "",
		Title: "订单YT7631606603205 出库", Description: "琦立工作室的集运出库",
		Priority: 1,
		CreatedAt: now.Add(-1 * time.Hour), UpdatedAt: now.Add(-1 * time.Hour),
	}

	// WO-5: 入库 — 待处理 (urgent)
	r.orders[5] = &domain.WorkOrder{
		ID: 5, TenantID: tenant, WarehouseID: 1,
		ProcessID: 1, ProcessName: "标准入库流程",
		ParcelID: &parcel1, Status: domain.WOStatusPending,
		CurrentStep: 1, AssignedTo: nil, AssignedName: "",
		Title: "YTO1111222233 入库(急)", Description: "圆通急件入库",
		Priority: 2,
		CreatedAt: now.Add(-30 * time.Minute), UpdatedAt: now.Add(-30 * time.Minute),
	}

	r.nextWOID = 6
}

// ===================================================================
// WorkflowProcess CRUD
// ===================================================================

// ListProcesses returns all processes for a tenant
func (r *MemWorkflowRepo) ListProcesses(ctx context.Context, tenantID int64) ([]domain.WorkflowProcess, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.WorkflowProcess
	for _, p := range r.processes {
		if p.TenantID == tenantID {
			result = append(result, *p)
		}
	}
	return result, nil
}

// GetProcessByID returns a single process by ID
func (r *MemWorkflowRepo) GetProcessByID(ctx context.Context, tenantID, id int64) (*domain.WorkflowProcess, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.processes[id]
	if !ok || p.TenantID != tenantID {
		return nil, nil
	}
	return p, nil
}

// CreateProcess inserts a new workflow process
func (r *MemWorkflowRepo) CreateProcess(ctx context.Context, p *domain.WorkflowProcess) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p.ID = atomic.AddInt64(&r.nextPID, 1) - 1
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	r.processes[p.ID] = p
	return nil
}

// ===================================================================
// WorkOrder CRUD
// ===================================================================

// ListWorkOrders returns work orders for a tenant with pagination
func (r *MemWorkflowRepo) ListWorkOrders(ctx context.Context, tenantID int64, offset, limit int) ([]domain.WorkOrder, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.WorkOrder
	for _, wo := range r.orders {
		if wo.TenantID == tenantID {
			result = append(result, *wo)
		}
	}
	total := int64(len(result))
	if offset >= int(total) {
		return nil, total, nil
	}
	end := offset + limit
	if end > int(total) {
		end = int(total)
	}
	return result[offset:end], total, nil
}

// GetWorkOrderByID returns a single work order
func (r *MemWorkflowRepo) GetWorkOrderByID(ctx context.Context, tenantID, id int64) (*domain.WorkOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	wo, ok := r.orders[id]
	if !ok || wo.TenantID != tenantID {
		return nil, nil
	}
	return wo, nil
}

// CreateWorkOrder inserts a new work order
func (r *MemWorkflowRepo) CreateWorkOrder(ctx context.Context, wo *domain.WorkOrder) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	wo.ID = atomic.AddInt64(&r.nextWOID, 1) - 1
	now := time.Now()
	wo.CreatedAt = now
	wo.UpdatedAt = now
	r.orders[wo.ID] = wo
	return nil
}

// UpdateWorkOrder modifies an existing work order
func (r *MemWorkflowRepo) UpdateWorkOrder(ctx context.Context, wo *domain.WorkOrder) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	wo.UpdatedAt = time.Now()
	r.orders[wo.ID] = wo
	return nil
}

// DeleteWorkOrder removes a work order
func (r *MemWorkflowRepo) DeleteWorkOrder(ctx context.Context, tenantID, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if wo, ok := r.orders[id]; ok && wo.TenantID == tenantID {
		delete(r.orders, id)
	}
	return nil
}

// GetProcessForWorkOrder returns the process associated with a work order
func (r *MemWorkflowRepo) GetProcessForWorkOrder(ctx context.Context, tenantID, processID int64) (*domain.WorkflowProcess, error) {
	return r.GetProcessByID(ctx, tenantID, processID)
}
