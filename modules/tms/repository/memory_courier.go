package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/modules/tms/domain"
)

// MemCourierRepo is an in-memory implementation for couriers.
type MemCourierRepo struct {
	mu       sync.RWMutex
	couriers map[int64]*domain.Courier
	nextID   int64
}

func NewMemCourierRepo() *MemCourierRepo {
	r := &MemCourierRepo{couriers: make(map[int64]*domain.Courier), nextID: 1}
	r.seed()
	return r
}

func (r *MemCourierRepo) seed() {
	seeds := []struct {
		carrierID    int64
		name         string
		phone        string
		vehiclePlate string
		idNumber     string
	}{
		{1, "张师傅", "886911111111", "ABC-1234", "A123456789"},
		{1, "李师傅", "886922222222", "DEF-5678", "B123456789"},
		{2, "王师傅", "886933333333", "GHI-9012", "C123456789"},
		{3, "赵师傅", "886944444444", "JKL-3456", "D123456789"},
		{4, "钱师傅", "886955555555", "MNO-7890", "E123456789"},
	}
	now := time.Now()
	for _, s := range seeds {
		id := atomic.AddInt64(&r.nextID, 1) - 1
		r.couriers[id] = &domain.Courier{
			ID:           id,
			CarrierID:    s.carrierID,
			TenantID:     1,
			Name:         s.name,
			Phone:        s.phone,
			VehiclePlate: s.vehiclePlate,
			IDNumber:     s.idNumber,
			IsActive:     true,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
	}
}

func (r *MemCourierRepo) Create(ctx context.Context, c *domain.Courier) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c.ID = atomic.AddInt64(&r.nextID, 1) - 1
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	r.couriers[c.ID] = c
	return nil
}

func (r *MemCourierRepo) GetByID(ctx context.Context, id int64) (*domain.Courier, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.couriers[id]
	if !ok {
		return nil, nil
	}
	return c, nil
}

func (r *MemCourierRepo) ListByCarrier(ctx context.Context, carrierID int64) ([]domain.Courier, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.Courier
	for _, c := range r.couriers {
		if c.CarrierID == carrierID {
			result = append(result, *c)
		}
	}
	return result, nil
}

func (r *MemCourierRepo) List(ctx context.Context, tenantID int64) ([]domain.Courier, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.Courier
	for _, c := range r.couriers {
		if c.TenantID == tenantID {
			result = append(result, *c)
		}
	}
	return result, nil
}

func (r *MemCourierRepo) Update(ctx context.Context, c *domain.Courier) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.couriers[c.ID]; !ok {
		return nil
	}
	c.UpdatedAt = time.Now()
	r.couriers[c.ID] = c
	return nil
}

func (r *MemCourierRepo) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.couriers, id)
	return nil
}

var _ CourierRepository = (*MemCourierRepo)(nil)
