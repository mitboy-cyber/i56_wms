package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/modules/order/domain"
)

type MemOrderRepo struct {
	mu     sync.RWMutex
	orders map[int64]*domain.Order
	nextID int64
}

func NewMemOrderRepo() *MemOrderRepo {
	return &MemOrderRepo{orders: make(map[int64]*domain.Order), nextID: 1}
}

func (r *MemOrderRepo) Create(ctx context.Context, o *domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	o.ID = atomic.AddInt64(&r.nextID, 1) - 1
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()
	r.orders[o.ID] = o
	return nil
}

func (r *MemOrderRepo) GetByID(ctx context.Context, tenantID, id int64) (*domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	o, ok := r.orders[id]
	if !ok || o.TenantID != tenantID { return nil, nil }
	return o, nil
}

func (r *MemOrderRepo) GetByOrderNo(ctx context.Context, tenantID int64, orderNo string) (*domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, o := range r.orders {
		if o.TenantID == tenantID && o.OrderNo == orderNo { return o, nil }
	}
	return nil, nil
}

func (r *MemOrderRepo) List(ctx context.Context, tenantID int64, offset, limit int) ([]domain.Order, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.Order
	for _, o := range r.orders {
		if o.TenantID == tenantID { result = append(result, *o) }
	}
	total := int64(len(result))
	if offset >= int(total) { return nil, total, nil }
	end := offset + limit
	if end > int(total) { end = int(total) }
	return result[offset:end], total, nil
}

func (r *MemOrderRepo) ListByClient(ctx context.Context, tenantID, clientID int64, offset, limit int) ([]domain.Order, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.Order
	for _, o := range r.orders {
		if o.TenantID == tenantID && o.ClientID == clientID { result = append(result, *o) }
	}
	total := int64(len(result))
	if offset >= int(total) { return nil, total, nil }
	end := offset + limit
	if end > int(total) { end = int(total) }
	return result[offset:end], total, nil
}

func (r *MemOrderRepo) Update(ctx context.Context, o *domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	o.UpdatedAt = time.Now()
	r.orders[o.ID] = o
	return nil
}

var _ OrderRepository = (*MemOrderRepo)(nil)
