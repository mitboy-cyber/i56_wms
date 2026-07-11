package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/modules/transport/domain"
)

type MemRouteRepo struct {
	mu     sync.RWMutex
	routes map[int64]*domain.Route
	nextID int64
}

func NewMemRouteRepo() *MemRouteRepo {
	return &MemRouteRepo{routes: make(map[int64]*domain.Route), nextID: 1}
}

func (r *MemRouteRepo) Create(ctx context.Context, route *domain.Route) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	route.ID = atomic.AddInt64(&r.nextID, 1) - 1
	route.CreatedAt = time.Now()
	route.UpdatedAt = time.Now()
	r.routes[route.ID] = route
	return nil
}

func (r *MemRouteRepo) GetByID(ctx context.Context, tenantID, id int64) (*domain.Route, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rt, ok := r.routes[id]
	if !ok || rt.TenantID != tenantID { return nil, nil }
	return rt, nil
}

func (r *MemRouteRepo) List(ctx context.Context, tenantID int64, offset, limit int) ([]domain.Route, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.Route
	for _, rt := range r.routes {
		if rt.TenantID == tenantID { result = append(result, *rt) }
	}
	total := int64(len(result))
	if offset >= int(total) { return nil, total, nil }
	end := offset + limit
	if end > int(total) { end = int(total) }
	return result[offset:end], total, nil
}

func (r *MemRouteRepo) Update(ctx context.Context, route *domain.Route) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	route.UpdatedAt = time.Now()
	r.routes[route.ID] = route
	return nil
}

func (r *MemRouteRepo) Delete(ctx context.Context, tenantID, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.routes, id)
	return nil
}

var _ RouteRepository = (*MemRouteRepo)(nil)

type MemCourierRepo struct {
	mu       sync.RWMutex
	couriers map[string]*domain.Courier
}

func NewMemCourierRepo() *MemCourierRepo {
	r := &MemCourierRepo{couriers: make(map[string]*domain.Courier)}
	seeds := []struct{ name, code string }{
		{"顺丰速运", "SF"}, {"中通快递", "ZTO"}, {"圆通速递", "YTO"},
		{"申通快递", "STO"}, {"韵达快递", "YD"}, {"百世快递", "HTKY"},
		{"极兔速递", "JTSD"}, {"邮政EMS", "EMS"}, {"京东物流", "JD"},
	}
	for _, s := range seeds {
		r.couriers[s.code] = &domain.Courier{Name: s.name, Code: s.code, CountryRegion: "中国大陆"}
	}
	return r
}

func (r *MemCourierRepo) Create(ctx context.Context, c *domain.Courier) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.couriers[c.Code] = c
	return nil
}

func (r *MemCourierRepo) List(ctx context.Context) ([]domain.Courier, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]domain.Courier, 0, len(r.couriers))
	for _, c := range r.couriers { result = append(result, *c) }
	return result, nil
}

func (r *MemCourierRepo) GetByCode(ctx context.Context, code string) (*domain.Courier, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.couriers[code]
	if !ok { return nil, nil }
	return c, nil
}

func (r *MemCourierRepo) Update(ctx context.Context, code string, c *domain.Courier) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.couriers[code]; !ok { return nil }
	r.couriers[code] = c
	return nil
}

func (r *MemCourierRepo) Delete(ctx context.Context, code string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.couriers, code)
	return nil
}

func (r *MemCourierRepo) DetectByTrackingNo(trackingNo string) *domain.Courier {
	if len(trackingNo) < 2 { return nil }
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, c := range r.couriers {
		if len(trackingNo) >= 2 && trackingNo[:2] == c.Code[:2] { return c }
	}
	return nil
}

var _ CourierRepository = (*MemCourierRepo)(nil)
