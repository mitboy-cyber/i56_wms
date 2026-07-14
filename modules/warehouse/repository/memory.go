package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/framework/core/errors"
	"github.com/i56/modules/warehouse/domain"
)

type MemWarehouseRepo struct {
	mu         sync.RWMutex
	warehouses map[int64]*domain.Warehouse
	nextID     int64
}

func NewMemWarehouseRepo() *MemWarehouseRepo {
	return &MemWarehouseRepo{warehouses: make(map[int64]*domain.Warehouse), nextID: 1}
}

func (r *MemWarehouseRepo) Create(ctx context.Context, tenantID int64, w *domain.Warehouse) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	w.ID = atomic.AddInt64(&r.nextID, 1) - 1
	w.CreatedAt = time.Now()
	w.UpdatedAt = time.Now()
	r.warehouses[w.ID] = w
	return nil
}

func (r *MemWarehouseRepo) GetByID(ctx context.Context, tenantID, id int64) (*domain.Warehouse, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	w, ok := r.warehouses[id]
	if !ok || w.TenantID != tenantID {
		return nil, nil
	}
	return w, nil
}

func (r *MemWarehouseRepo) List(ctx context.Context, tenantID int64, offset, limit int) ([]domain.Warehouse, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.Warehouse
	for _, w := range r.warehouses {
		if w.TenantID == tenantID {
			result = append(result, *w)
		}
	}
	total := int64(len(result))
	if offset >= int(total) { return nil, total, nil }
	end := offset + limit
	if end > int(total) { end = int(total) }
	return result[offset:end], total, nil
}

func (r *MemWarehouseRepo) Update(ctx context.Context, tenantID, id int64, w *domain.Warehouse) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.warehouses[id]
	if !ok || existing.TenantID != tenantID {
		return errors.NewNotFound("Warehouse")
	}
	w.UpdatedAt = time.Now()
	r.warehouses[id] = w
	return nil
}

var _ WarehouseRepository = (*MemWarehouseRepo)(nil)

func (r *MemWarehouseRepo) Delete(ctx context.Context, tenantID, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.warehouses, id)
	return nil
}
