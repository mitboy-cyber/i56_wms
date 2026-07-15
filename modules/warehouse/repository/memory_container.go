package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/modules/warehouse/domain"
)

// MemContainerRepo is an in-memory implementation for containers.
type MemContainerRepo struct {
	mu         sync.RWMutex
	containers map[int64]*domain.Container
	nextID     int64
}

func NewMemContainerRepo() *MemContainerRepo {
	r := &MemContainerRepo{containers: make(map[int64]*domain.Container), nextID: 1}
	r.seed()
	return r
}

func (r *MemContainerRepo) seed() {
	seeds := []struct {
		warehouseID int64
		containerNo string
		ctype       domain.ContainerType
		sealNo      string
		routeID     int64
		status      domain.ContainerStatus
		maxCap      float64
		curWeight   float64
		parcelCount int
		remark      string
	}{
		{1, "CNTR202401001", domain.Container20GP, "", 1, domain.ContainerStatusAvailable, 28000, 0, 0, "台北线-20尺柜"},
		{1, "CNTR202401002", domain.Container40GP, "SEAL001", 1, domain.ContainerStatusSealed, 30000, 15200.5, 340, "台北线-40尺柜-已封签"},
		{1, "CNTR202401003", domain.Container40HQ, "SEAL002", 2, domain.ContainerStatusLoaded, 30500, 13200.0, 280, "台中线-40高柜"},
		{1, "CNTR202401004", domain.Container20GP, "", 3, domain.ContainerStatusLoading, 28000, 5200.3, 120, "高雄线-正在装柜"},
	}
	now := time.Now()
	for _, s := range seeds {
		id := atomic.AddInt64(&r.nextID, 1) - 1
		c := &domain.Container{
			ID:            id,
			WarehouseID:   s.warehouseID,
			ContainerNo:   s.containerNo,
			ContainerType: s.ctype,
			SealNo:        s.sealNo,
			RouteID:       s.routeID,
			Status:        s.status,
			MaxCapacity:   s.maxCap,
			CurrentWeight: s.curWeight,
			ParcelCount:   s.parcelCount,
			Remark:        s.remark,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		if s.status == domain.ContainerStatusSealed {
			c.SealedAt = &now
		}
		r.containers[id] = c
	}
}

func (r *MemContainerRepo) Create(ctx context.Context, c *domain.Container) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c.ID = atomic.AddInt64(&r.nextID, 1) - 1
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	r.containers[c.ID] = c
	return nil
}

func (r *MemContainerRepo) GetByID(ctx context.Context, id int64) (*domain.Container, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.containers[id]
	if !ok {
		return nil, nil
	}
	return c, nil
}

func (r *MemContainerRepo) GetByContainerNo(ctx context.Context, containerNo string) (*domain.Container, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, c := range r.containers {
		if c.ContainerNo == containerNo {
			return c, nil
		}
	}
	return nil, nil
}

func (r *MemContainerRepo) ListByWarehouse(ctx context.Context, warehouseID int64) ([]domain.Container, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.Container
	for _, c := range r.containers {
		if c.WarehouseID == warehouseID {
			result = append(result, *c)
		}
	}
	return result, nil
}

func (r *MemContainerRepo) ListByStatus(ctx context.Context, warehouseID int64, status domain.ContainerStatus) ([]domain.Container, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.Container
	for _, c := range r.containers {
		if c.WarehouseID == warehouseID && c.Status == status {
			result = append(result, *c)
		}
	}
	return result, nil
}

func (r *MemContainerRepo) ListByRoute(ctx context.Context, routeID int64) ([]domain.Container, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.Container
	for _, c := range r.containers {
		if c.RouteID == routeID {
			result = append(result, *c)
		}
	}
	return result, nil
}

func (r *MemContainerRepo) Update(ctx context.Context, c *domain.Container) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.containers[c.ID]; !ok {
		return nil
	}
	c.UpdatedAt = time.Now()
	r.containers[c.ID] = c
	return nil
}

func (r *MemContainerRepo) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.containers, id)
	return nil
}

var _ ContainerRepository = (*MemContainerRepo)(nil)
