package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/modules/tms/domain"
)

// MemCargoTypeRepo is an in-memory implementation for cargo types.
type MemCargoTypeRepo struct {
	mu         sync.RWMutex
	cargoTypes map[int64]*domain.CargoType
	nextID     int64
}

func NewMemCargoTypeRepo() *MemCargoTypeRepo {
	r := &MemCargoTypeRepo{cargoTypes: make(map[int64]*domain.CargoType), nextID: 1}
	r.seed()
	return r
}

func (r *MemCargoTypeRepo) seed() {
	now := time.Now()
	for _, ct := range domain.DefaultCargoTypes() {
		ct.TenantID = 1
		ct.CreatedAt = now
		ct.UpdatedAt = now
		r.cargoTypes[ct.ID] = &ct
		r.nextID = ct.ID + 1
	}
}

func (r *MemCargoTypeRepo) List(ctx context.Context, tenantID int64) ([]domain.CargoType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.CargoType
	for _, c := range r.cargoTypes {
		if c.TenantID == tenantID {
			result = append(result, *c)
		}
	}
	return result, nil
}

func (r *MemCargoTypeRepo) GetByID(ctx context.Context, id int64) (*domain.CargoType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.cargoTypes[id]
	if !ok {
		return nil, nil
	}
	return c, nil
}

func (r *MemCargoTypeRepo) GetByCode(ctx context.Context, tenantID int64, code string) (*domain.CargoType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, c := range r.cargoTypes {
		if c.TenantID == tenantID && c.Code == code {
			return c, nil
		}
	}
	return nil, nil
}

func (r *MemCargoTypeRepo) Create(ctx context.Context, c *domain.CargoType) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c.ID = atomic.AddInt64(&r.nextID, 1) - 1
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	r.cargoTypes[c.ID] = c
	return nil
}

func (r *MemCargoTypeRepo) Update(ctx context.Context, c *domain.CargoType) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.cargoTypes[c.ID]; !ok {
		return nil
	}
	c.UpdatedAt = time.Now()
	r.cargoTypes[c.ID] = c
	return nil
}

var _ CargoTypeRepository = (*MemCargoTypeRepo)(nil)
