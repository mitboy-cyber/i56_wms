package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/modules/order/domain"
)

// MemConsolidationOrderRepo is an in-memory implementation for consolidation orders.
type MemConsolidationOrderRepo struct {
	mu     sync.RWMutex
	orders map[int64]*domain.ConsolidationOrder
	nextID int64
}

func NewMemConsolidationOrderRepo() *MemConsolidationOrderRepo {
	r := &MemConsolidationOrderRepo{orders: make(map[int64]*domain.ConsolidationOrder), nextID: 1}
	r.seed()
	return r
}

func (r *MemConsolidationOrderRepo) seed() {
	seeds := []struct {
		tenantID       int64
		orderNo        string
		clientID       int64
		memberID       int64
		addressID      int64
		warehouseID    int64
		routeID        int64
		parcelIDs      []int64
		parcelCount    int
		totalWeight    float64
		totalChargeable float64
		shippingFee    float64
		serviceFee     float64
		totalPrice     float64
		status         domain.ConsolidationOrderStatus
		declarantID    int64
		customsNumber  string
		carrierTN      string
		remark         string
	}{
		{1, "CON202401001", 1, 1, 1, 1, 1, []int64{1, 2}, 2, 0.23, 0.25, 45.00, 10.00, 55.00, domain.ConsolStatusPacked, 1, "CN202401001", "SF1234567890", "合并发货"},
		{1, "CON202401002", 2, 2, 3, 1, 1, []int64{3, 4}, 2, 0.60, 0.70, 80.00, 15.00, 95.00, domain.ConsolStatusMerged, 2, "", "", "待打包"},
		{1, "CON202401003", 1, 1, 2, 1, 2, []int64{5}, 1, 0.55, 0.55, 50.00, 5.00, 55.00, domain.ConsolStatusShipped, 1, "CN202401003", "ZTO9876543210", "已发货"},
		{1, "CON202401004", 3, 5, 5, 1, 3, []int64{}, 0, 0, 0, 0, 0, 0, domain.ConsolStatusDraft, 3, "", "", "新建合并单"},
	}
	now := time.Now()
	for _, s := range seeds {
		id := atomic.AddInt64(&r.nextID, 1) - 1
		r.orders[id] = &domain.ConsolidationOrder{
			ID:                id,
			TenantID:          s.tenantID,
			OrderNo:           s.orderNo,
			ClientID:          s.clientID,
			MemberID:          s.memberID,
			MemberAddressID:   s.addressID,
			WarehouseID:       s.warehouseID,
			RouteID:           s.routeID,
			ParcelIDs:         s.parcelIDs,
			ParcelCount:       s.parcelCount,
			TotalWeight:       s.totalWeight,
			TotalChargeable:   s.totalChargeable,
			ShippingFee:       s.shippingFee,
			ServiceFee:        s.serviceFee,
			TotalPrice:        s.totalPrice,
			Status:            s.status,
			DeclarantID:       s.declarantID,
			CustomsNumber:     s.customsNumber,
			CarrierTrackingNo: s.carrierTN,
			Remark:            s.remark,
			CreatedAt:         now,
			UpdatedAt:         now,
		}
	}
}

func (r *MemConsolidationOrderRepo) Create(ctx context.Context, o *domain.ConsolidationOrder) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	o.ID = atomic.AddInt64(&r.nextID, 1) - 1
	now := time.Now()
	o.CreatedAt = now
	o.UpdatedAt = now
	r.orders[o.ID] = o
	return nil
}

func (r *MemConsolidationOrderRepo) GetByID(ctx context.Context, tenantID, id int64) (*domain.ConsolidationOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	o, ok := r.orders[id]
	if !ok || o.TenantID != tenantID {
		return nil, nil
	}
	return o, nil
}

func (r *MemConsolidationOrderRepo) GetByOrderNo(ctx context.Context, tenantID int64, orderNo string) (*domain.ConsolidationOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, o := range r.orders {
		if o.TenantID == tenantID && o.OrderNo == orderNo {
			return o, nil
		}
	}
	return nil, nil
}

func (r *MemConsolidationOrderRepo) List(ctx context.Context, tenantID int64, offset, limit int) ([]domain.ConsolidationOrder, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.ConsolidationOrder
	for _, o := range r.orders {
		if o.TenantID == tenantID {
			result = append(result, *o)
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

func (r *MemConsolidationOrderRepo) ListByClient(ctx context.Context, tenantID, clientID int64) ([]domain.ConsolidationOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.ConsolidationOrder
	for _, o := range r.orders {
		if o.TenantID == tenantID && o.ClientID == clientID {
			result = append(result, *o)
		}
	}
	return result, nil
}

func (r *MemConsolidationOrderRepo) ListByMember(ctx context.Context, tenantID, memberID int64) ([]domain.ConsolidationOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.ConsolidationOrder
	for _, o := range r.orders {
		if o.TenantID == tenantID && o.MemberID == memberID {
			result = append(result, *o)
		}
	}
	return result, nil
}

func (r *MemConsolidationOrderRepo) Update(ctx context.Context, o *domain.ConsolidationOrder) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.orders[o.ID]; !ok {
		return nil
	}
	o.UpdatedAt = time.Now()
	r.orders[o.ID] = o
	return nil
}

func (r *MemConsolidationOrderRepo) Delete(ctx context.Context, tenantID, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.orders[id]
	if !ok || existing.TenantID != tenantID {
		return nil
	}
	delete(r.orders, id)
	return nil
}

var _ ConsolidationOrderRepository = (*MemConsolidationOrderRepo)(nil)
