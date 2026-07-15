package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/modules/finance/domain"
)

// MemOrderProfitRepo is an in-memory implementation for order profits.
type MemOrderProfitRepo struct {
	mu      sync.RWMutex
	profits map[int64]*domain.OrderProfit
	nextID  int64
}

func NewMemOrderProfitRepo() *MemOrderProfitRepo {
	r := &MemOrderProfitRepo{profits: make(map[int64]*domain.OrderProfit), nextID: 1}
	r.seed()
	return r
}

func (r *MemOrderProfitRepo) seed() {
	now := time.Now()
	seeds := []struct {
		orderID    int64
		orderNo    string
		clientID   int64
		routeID    int64
		revenue    float64
		cost       float64
		grossProfit float64
		shipCost   float64
		servCost   float64
	}{
		{1, "ORD202401001", 1, 1, 320.00, 200.00, 120.00, 180.00, 20.00},
		{2, "ORD202401002", 2, 1, 550.00, 380.00, 170.00, 350.00, 30.00},
		{3, "ORD202401003", 1, 2, 480.00, 320.00, 160.00, 300.00, 20.00},
		{4, "ORD202401004", 3, 3, 280.00, 180.00, 100.00, 160.00, 20.00},
	}
	for _, s := range seeds {
		id := atomic.AddInt64(&r.nextID, 1) - 1
		margin := 0.0
		if s.revenue > 0 {
			margin = s.grossProfit / s.revenue * 100
		}
		r.profits[id] = &domain.OrderProfit{
			ID:             id,
			TenantID:       1,
			OrderID:        s.orderID,
			OrderNo:        s.orderNo,
			ClientID:       s.clientID,
			RouteID:        s.routeID,
			Revenue:        s.revenue,
			Cost:           s.cost,
			GrossProfit:    s.grossProfit,
			ProfitMargin:   margin,
			ShippingCost:   s.shipCost,
			ServiceCost:    s.servCost,
			CreatedAt:      now,
		}
	}
}

func (r *MemOrderProfitRepo) Create(ctx context.Context, p *domain.OrderProfit) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p.ID = atomic.AddInt64(&r.nextID, 1) - 1
	p.CreatedAt = time.Now()
	r.profits[p.ID] = p
	return nil
}

func (r *MemOrderProfitRepo) GetByOrderID(ctx context.Context, orderID int64) (*domain.OrderProfit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.profits {
		if p.OrderID == orderID {
			return p, nil
		}
	}
	return nil, nil
}

func (r *MemOrderProfitRepo) ListByClient(ctx context.Context, tenantID, clientID int64) ([]domain.OrderProfit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.OrderProfit
	for _, p := range r.profits {
		if p.TenantID == tenantID && p.ClientID == clientID {
			result = append(result, *p)
		}
	}
	return result, nil
}

func (r *MemOrderProfitRepo) ListByRoute(ctx context.Context, routeID int64) ([]domain.OrderProfit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.OrderProfit
	for _, p := range r.profits {
		if p.RouteID == routeID {
			result = append(result, *p)
		}
	}
	return result, nil
}

var _ OrderProfitRepository = (*MemOrderProfitRepo)(nil)
