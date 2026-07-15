package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/modules/pda/domain"
)

// MemWorkOrderRepo is an in-memory implementation for work orders.
type MemWorkOrderRepo struct {
	mu         sync.RWMutex
	workOrders map[int64]*domain.WorkOrder
	nextID     int64
}

func NewMemWorkOrderRepo() *MemWorkOrderRepo {
	r := &MemWorkOrderRepo{workOrders: make(map[int64]*domain.WorkOrder), nextID: 1}
	r.seed()
	return r
}

func (r *MemWorkOrderRepo) seed() {
	now := time.Now()
	seeds := []struct {
		warehouseID  int64
		workOrderNo  string
		templateID   int64
		wtype        string
		status       domain.WorkOrderStatus
		priority     domain.WorkOrderPriority
		assignedTo   int64
		orderID      int64
		parcelIDs    []int64
		locationCode string
		targetLoc    string
		instructions string
		createdBy    int64
	}{
		{1, "WO202401001", 1, "receive", domain.WOStatusInProgress, domain.WOPriorityNormal, 1, 0, []int64{6}, "A-01-01", "A-01-10", "收货并称重", 2},
		{1, "WO202401002", 2, "pick", domain.WOStatusAssigned, domain.WOPriorityHigh, 2, 1, []int64{1, 2}, "A-01-01", "PACK-01", "按订单1拣货", 2},
		{1, "WO202401003", 3, "pack", domain.WOStatusCompleted, domain.WOPriorityNormal, 3, 2, []int64{3, 4}, "PACK-01", "OUTBOUND-01", "打包并发货", 2},
		{1, "WO202401004", 4, "load", domain.WOStatusDraft, domain.WOPriorityUrgent, 0, 3, []int64{5}, "OUTBOUND-01", "CNTR202401003", "紧急装柜任务", 2},
	}
	for _, s := range seeds {
		id := atomic.AddInt64(&r.nextID, 1) - 1
		wo := &domain.WorkOrder{
			ID:             id,
			TenantID:       1,
			WarehouseID:    s.warehouseID,
			WorkOrderNo:    s.workOrderNo,
			TemplateID:     s.templateID,
			Type:           s.wtype,
			Status:         s.status,
			Priority:       s.priority,
			AssignedTo:     s.assignedTo,
			OrderID:        s.orderID,
			ParcelIDs:      s.parcelIDs,
			LocationCode:   s.locationCode,
			TargetLocation: s.targetLoc,
			Instructions:   s.instructions,
			CreatedBy:      s.createdBy,
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		if s.status == domain.WOStatusInProgress {
			wo.StartedAt = &now
		}
		if s.status == domain.WOStatusCompleted {
			wo.CompletedAt = &now
		}
		r.workOrders[id] = wo
	}
}

func (r *MemWorkOrderRepo) Create(ctx context.Context, wo *domain.WorkOrder) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	wo.ID = atomic.AddInt64(&r.nextID, 1) - 1
	now := time.Now()
	wo.CreatedAt = now
	wo.UpdatedAt = now
	r.workOrders[wo.ID] = wo
	return nil
}

func (r *MemWorkOrderRepo) GetByID(ctx context.Context, id int64) (*domain.WorkOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	wo, ok := r.workOrders[id]
	if !ok {
		return nil, nil
	}
	return wo, nil
}

func (r *MemWorkOrderRepo) GetByWorkOrderNo(ctx context.Context, tenantID int64, woNo string) (*domain.WorkOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, wo := range r.workOrders {
		if wo.TenantID == tenantID && wo.WorkOrderNo == woNo {
			return wo, nil
		}
	}
	return nil, nil
}

func (r *MemWorkOrderRepo) ListByOperator(ctx context.Context, operatorID int64) ([]domain.WorkOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.WorkOrder
	for _, wo := range r.workOrders {
		if wo.AssignedTo == operatorID {
			result = append(result, *wo)
		}
	}
	return result, nil
}

func (r *MemWorkOrderRepo) ListByWarehouse(ctx context.Context, tenantID, warehouseID int64, status domain.WorkOrderStatus) ([]domain.WorkOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.WorkOrder
	for _, wo := range r.workOrders {
		if wo.TenantID == tenantID && wo.WarehouseID == warehouseID && (status == "" || wo.Status == status) {
			result = append(result, *wo)
		}
	}
	return result, nil
}

func (r *MemWorkOrderRepo) ListByOrder(ctx context.Context, orderID int64) ([]domain.WorkOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.WorkOrder
	for _, wo := range r.workOrders {
		if wo.OrderID == orderID {
			result = append(result, *wo)
		}
	}
	return result, nil
}

func (r *MemWorkOrderRepo) Update(ctx context.Context, wo *domain.WorkOrder) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.workOrders[wo.ID]; !ok {
		return nil
	}
	wo.UpdatedAt = time.Now()
	r.workOrders[wo.ID] = wo
	return nil
}

func (r *MemWorkOrderRepo) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.workOrders, id)
	return nil
}

var _ WorkOrderRepository = (*MemWorkOrderRepo)(nil)
