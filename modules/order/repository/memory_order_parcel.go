package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/modules/order/domain"
)

// MemOrderParcelRepo is an in-memory implementation for the OrderParcel join entity.
type MemOrderParcelRepo struct {
	mu      sync.RWMutex
	entries map[int64]*domain.OrderParcel
	nextID  int64
}

func NewMemOrderParcelRepo() *MemOrderParcelRepo {
	r := &MemOrderParcelRepo{entries: make(map[int64]*domain.OrderParcel), nextID: 1}
	r.seed()
	return r
}

func (r *MemOrderParcelRepo) seed() {
	seeds := []struct {
		orderID   int64
		parcelID  int64
		status    domain.OrderParcelStatus
		sortOrder int
	}{
		{1, 1, domain.OPStatusPacked, 1},
		{1, 2, domain.OPStatusPacked, 2},
		{2, 3, domain.OPStatusPicked, 1},
		{2, 4, domain.OPStatusPicked, 2},
		{3, 5, domain.OPStatusShipped, 1},
	}
	now := time.Now()
	for _, s := range seeds {
		id := atomic.AddInt64(&r.nextID, 1) - 1
		r.entries[id] = &domain.OrderParcel{
			ID:        id,
			OrderID:   s.orderID,
			ParcelID:  s.parcelID,
			Status:    s.status,
			SortOrder: s.sortOrder,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
}

func (r *MemOrderParcelRepo) Create(ctx context.Context, op *domain.OrderParcel) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	op.ID = atomic.AddInt64(&r.nextID, 1) - 1
	now := time.Now()
	op.CreatedAt = now
	op.UpdatedAt = now
	r.entries[op.ID] = op
	return nil
}

func (r *MemOrderParcelRepo) GetByID(ctx context.Context, id int64) (*domain.OrderParcel, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	op, ok := r.entries[id]
	if !ok {
		return nil, nil
	}
	return op, nil
}

func (r *MemOrderParcelRepo) ListByOrder(ctx context.Context, orderID int64) ([]domain.OrderParcel, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.OrderParcel
	for _, op := range r.entries {
		if op.OrderID == orderID {
			result = append(result, *op)
		}
	}
	return result, nil
}

func (r *MemOrderParcelRepo) ListByParcel(ctx context.Context, parcelID int64) ([]domain.OrderParcel, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.OrderParcel
	for _, op := range r.entries {
		if op.ParcelID == parcelID {
			result = append(result, *op)
		}
	}
	return result, nil
}

func (r *MemOrderParcelRepo) Update(ctx context.Context, op *domain.OrderParcel) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.entries[op.ID]; !ok {
		return nil
	}
	op.UpdatedAt = time.Now()
	r.entries[op.ID] = op
	return nil
}

func (r *MemOrderParcelRepo) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, id)
	return nil
}

var _ OrderParcelRepository = (*MemOrderParcelRepo)(nil)
