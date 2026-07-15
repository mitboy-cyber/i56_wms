package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/modules/finance/domain"
)

// MemServiceProfitRepo is an in-memory implementation for service profits.
type MemServiceProfitRepo struct {
	mu      sync.RWMutex
	profits map[int64]*domain.ServiceProfit
	nextID  int64
}

func NewMemServiceProfitRepo() *MemServiceProfitRepo {
	r := &MemServiceProfitRepo{profits: make(map[int64]*domain.ServiceProfit), nextID: 1}
	r.seed()
	return r
}

func (r *MemServiceProfitRepo) seed() {
	now := time.Now()
	seeds := []struct {
		serviceOrderID int64
		orderID        int64
		clientID       int64
		serviceType    string
		revenue        float64
		cost           float64
		grossProfit    float64
	}{
		{1, 1, 1, "photos", 15.00, 5.00, 10.00},
		{2, 2, 2, "reinforce", 20.00, 8.00, 12.00},
		{4, 2, 2, "inspection", 10.00, 3.00, 7.00},
		{5, 3, 3, "insurance", 25.00, 18.00, 7.00},
	}
	for _, s := range seeds {
		id := atomic.AddInt64(&r.nextID, 1) - 1
		r.profits[id] = &domain.ServiceProfit{
			ID:             id,
			TenantID:       1,
			ServiceOrderID: s.serviceOrderID,
			OrderID:        s.orderID,
			ClientID:       s.clientID,
			ServiceType:    s.serviceType,
			Revenue:        s.revenue,
			Cost:           s.cost,
			GrossProfit:    s.grossProfit,
			CreatedAt:      now,
		}
	}
}

func (r *MemServiceProfitRepo) Create(ctx context.Context, p *domain.ServiceProfit) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p.ID = atomic.AddInt64(&r.nextID, 1) - 1
	p.CreatedAt = time.Now()
	r.profits[p.ID] = p
	return nil
}

func (r *MemServiceProfitRepo) GetByServiceOrderID(ctx context.Context, serviceOrderID int64) (*domain.ServiceProfit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.profits {
		if p.ServiceOrderID == serviceOrderID {
			return p, nil
		}
	}
	return nil, nil
}

func (r *MemServiceProfitRepo) ListByOrder(ctx context.Context, orderID int64) ([]domain.ServiceProfit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.ServiceProfit
	for _, p := range r.profits {
		if p.OrderID == orderID {
			result = append(result, *p)
		}
	}
	return result, nil
}

var _ ServiceProfitRepository = (*MemServiceProfitRepo)(nil)
