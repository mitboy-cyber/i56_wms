package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/modules/finance/domain"
)

// MemClientProfitRepo is an in-memory implementation for client profit aggregation.
type MemClientProfitRepo struct {
	mu      sync.RWMutex
	profits map[int64]*domain.ClientProfit
	nextID  int64
}

func NewMemClientProfitRepo() *MemClientProfitRepo {
	r := &MemClientProfitRepo{profits: make(map[int64]*domain.ClientProfit), nextID: 1}
	r.seed()
	return r
}

func (r *MemClientProfitRepo) seed() {
	now := time.Now()
	seeds := []struct {
		clientID        int64
		period          string
		totalRevenue    float64
		totalCost       float64
		totalProfit     float64
		orderCount      int
		parcelCount     int
	}{
		{1, "2024-01", 800.00, 520.00, 280.00, 2, 3},
		{2, "2024-01", 550.00, 380.00, 170.00, 1, 2},
		{3, "2024-01", 280.00, 180.00, 100.00, 1, 1},
	}
	for _, s := range seeds {
		id := atomic.AddInt64(&r.nextID, 1) - 1
		avgProfit := 0.0
		if s.orderCount > 0 {
			avgProfit = s.totalProfit / float64(s.orderCount)
		}
		r.profits[id] = &domain.ClientProfit{
			ID:               id,
			TenantID:         1,
			ClientID:         s.clientID,
			Period:           s.period,
			TotalRevenue:     s.totalRevenue,
			TotalCost:        s.totalCost,
			TotalProfit:      s.totalProfit,
			OrderCount:       s.orderCount,
			ParcelCount:      s.parcelCount,
			AvgProfitPerOrder: avgProfit,
			CreatedAt:        now,
		}
	}
}

func (r *MemClientProfitRepo) Create(ctx context.Context, p *domain.ClientProfit) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p.ID = atomic.AddInt64(&r.nextID, 1) - 1
	p.CreatedAt = time.Now()
	r.profits[p.ID] = p
	return nil
}

func (r *MemClientProfitRepo) GetByClientAndPeriod(ctx context.Context, tenantID, clientID int64, period string) (*domain.ClientProfit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.profits {
		if p.TenantID == tenantID && p.ClientID == clientID && p.Period == period {
			return p, nil
		}
	}
	return nil, nil
}

func (r *MemClientProfitRepo) ListByPeriod(ctx context.Context, tenantID int64, period string) ([]domain.ClientProfit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.ClientProfit
	for _, p := range r.profits {
		if p.TenantID == tenantID && p.Period == period {
			result = append(result, *p)
		}
	}
	return result, nil
}

var _ ClientProfitRepository = (*MemClientProfitRepo)(nil)
