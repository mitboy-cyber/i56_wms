package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/modules/finance/domain"
)

// MemRouteProfitRepo is an in-memory implementation for route profit aggregation.
type MemRouteProfitRepo struct {
	mu      sync.RWMutex
	profits map[int64]*domain.RouteProfit
	nextID  int64
}

func NewMemRouteProfitRepo() *MemRouteProfitRepo {
	r := &MemRouteProfitRepo{profits: make(map[int64]*domain.RouteProfit), nextID: 1}
	r.seed()
	return r
}

func (r *MemRouteProfitRepo) seed() {
	now := time.Now()
	seeds := []struct {
		routeID      int64
		period       string
		totalRevenue float64
		totalCost    float64
		totalProfit  float64
		orderCount   int
		totalWeight  float64
	}{
		{1, "2024-01", 870.00, 580.00, 290.00, 2, 0.95},
		{2, "2024-01", 480.00, 320.00, 160.00, 1, 0.55},
		{3, "2024-01", 280.00, 180.00, 100.00, 1, 0.55},
	}
	for _, s := range seeds {
		id := atomic.AddInt64(&r.nextID, 1) - 1
		avgPerKg := 0.0
		if s.totalWeight > 0 {
			avgPerKg = s.totalProfit / s.totalWeight
		}
		r.profits[id] = &domain.RouteProfit{
			ID:             id,
			TenantID:       1,
			RouteID:        s.routeID,
			Period:         s.period,
			TotalRevenue:   s.totalRevenue,
			TotalCost:      s.totalCost,
			TotalProfit:    s.totalProfit,
			OrderCount:     s.orderCount,
			TotalWeight:    s.totalWeight,
			AvgProfitPerKg: avgPerKg,
			CreatedAt:      now,
		}
	}
}

func (r *MemRouteProfitRepo) Create(ctx context.Context, p *domain.RouteProfit) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p.ID = atomic.AddInt64(&r.nextID, 1) - 1
	p.CreatedAt = time.Now()
	r.profits[p.ID] = p
	return nil
}

func (r *MemRouteProfitRepo) GetByRouteAndPeriod(ctx context.Context, routeID int64, period string) (*domain.RouteProfit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.profits {
		if p.RouteID == routeID && p.Period == period {
			return p, nil
		}
	}
	return nil, nil
}

func (r *MemRouteProfitRepo) ListByPeriod(ctx context.Context, tenantID int64, period string) ([]domain.RouteProfit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.RouteProfit
	for _, p := range r.profits {
		if p.TenantID == tenantID && p.Period == period {
			result = append(result, *p)
		}
	}
	return result, nil
}

var _ RouteProfitRepository = (*MemRouteProfitRepo)(nil)
