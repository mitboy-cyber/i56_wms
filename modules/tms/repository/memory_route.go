package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/modules/tms/domain"
)

// MemRouteRepo is an in-memory implementation for routes.
type MemRouteRepo struct {
	mu     sync.RWMutex
	routes map[int64]*domain.Route
	nextID int64
}

func NewMemRouteRepo() *MemRouteRepo {
	r := &MemRouteRepo{routes: make(map[int64]*domain.Route), nextID: 1}
	r.seed()
	return r
}

func (r *MemRouteRepo) seed() {
	seeds := []struct {
		name              string
		code              string
		routeType         domain.RouteType
		originCountry     string
		originCity        string
		destCountry       string
		destCity          string
		shippingProviderID int64
		carrierID         int64
		customsBrokerID   int64
		customsPointID    int64
		estimatedDays     int
		basePrice         float64
		pricePerKg        float64
		minWeight         float64
		maxWeight         float64
	}{
		{"深圳→台北(空运)", "SZ-TPE-AIR", domain.RouteTypeAir, "CN", "深圳", "TW", "台北", 1, 1, 1, 1, 3, 50.00, 25.00, 0.5, 1000.0},
		{"广州→台中(海运)", "GZ-TXG-SEA", domain.RouteTypeSea, "CN", "广州", "TW", "台中", 2, 2, 2, 2, 15, 200.00, 8.00, 10.0, 28000.0},
		{"厦门→高雄(海快)", "XM-KHH-SE", domain.RouteTypeSeaExpress, "CN", "厦门", "TW", "高雄", 3, 3, 1, 3, 7, 120.00, 12.00, 5.0, 20000.0},
		{"上海→台北(空运)", "SH-TPE-AIR", domain.RouteTypeAir, "CN", "上海", "TW", "台北", 1, 1, 1, 1, 2, 60.00, 28.00, 0.5, 800.0},
		{"义乌→基隆(海运)", "YW-KL-SEA", domain.RouteTypeSea, "CN", "义乌", "TW", "基隆", 2, 4, 2, 2, 18, 180.00, 7.00, 10.0, 28000.0},
	}
	now := time.Now()
	for _, s := range seeds {
		id := atomic.AddInt64(&r.nextID, 1) - 1
		r.routes[id] = &domain.Route{
			ID:                 id,
			TenantID:           1,
			Name:               s.name,
			Code:               s.code,
			RouteType:          s.routeType,
			OriginCountry:      s.originCountry,
			OriginCity:         s.originCity,
			DestCountry:        s.destCountry,
			DestCity:           s.destCity,
			ShippingProviderID: s.shippingProviderID,
			CarrierID:          s.carrierID,
			CustomsBrokerID:    s.customsBrokerID,
			CustomsPointID:     s.customsPointID,
			EstimatedDays:      s.estimatedDays,
			BasePrice:          s.basePrice,
			PricePerKg:         s.pricePerKg,
			MinWeight:          s.minWeight,
			MaxWeight:          s.maxWeight,
			IsActive:           true,
			CreatedAt:          now,
			UpdatedAt:          now,
		}
	}
}

func (r *MemRouteRepo) Create(ctx context.Context, route *domain.Route) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	route.ID = atomic.AddInt64(&r.nextID, 1) - 1
	now := time.Now()
	route.CreatedAt = now
	route.UpdatedAt = now
	r.routes[route.ID] = route
	return nil
}

func (r *MemRouteRepo) GetByID(ctx context.Context, id int64) (*domain.Route, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	route, ok := r.routes[id]
	if !ok {
		return nil, nil
	}
	return route, nil
}

func (r *MemRouteRepo) GetByCode(ctx context.Context, tenantID int64, code string) (*domain.Route, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, route := range r.routes {
		if route.TenantID == tenantID && route.Code == code {
			return route, nil
		}
	}
	return nil, nil
}

func (r *MemRouteRepo) List(ctx context.Context, tenantID int64) ([]domain.Route, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.Route
	for _, route := range r.routes {
		if route.TenantID == tenantID {
			result = append(result, *route)
		}
	}
	return result, nil
}

func (r *MemRouteRepo) ListByRouteType(ctx context.Context, tenantID int64, routeType domain.RouteType) ([]domain.Route, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.Route
	for _, route := range r.routes {
		if route.TenantID == tenantID && route.RouteType == routeType {
			result = append(result, *route)
		}
	}
	return result, nil
}

func (r *MemRouteRepo) Update(ctx context.Context, route *domain.Route) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.routes[route.ID]; !ok {
		return nil
	}
	route.UpdatedAt = time.Now()
	r.routes[route.ID] = route
	return nil
}

func (r *MemRouteRepo) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.routes, id)
	return nil
}

var _ RouteRepository = (*MemRouteRepo)(nil)
