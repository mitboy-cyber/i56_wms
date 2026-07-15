package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/modules/tms/domain"
)

// MemCarrierRepo is an in-memory implementation for carriers.
type MemCarrierRepo struct {
	mu       sync.RWMutex
	carriers map[int64]*domain.Carrier
	nextID   int64
}

func NewMemCarrierRepo() *MemCarrierRepo {
	r := &MemCarrierRepo{carriers: make(map[int64]*domain.Carrier), nextID: 1}
	r.seed()
	return r
}

func (r *MemCarrierRepo) seed() {
	seeds := []struct {
		name         string
		code         string
		phone        string
		email        string
		website      string
		accountNo    string
	}{
		{"顺丰速运", "SF", "95338", "service@sf-express.com", "https://www.sf-express.com", "ACC-SF-001"},
		{"圆通速递", "YTO", "95554", "service@yto.net.cn", "https://www.yto.net.cn", "ACC-YTO-001"},
		{"中通快递", "ZTO", "95311", "service@zto.com", "https://www.zto.com", "ACC-ZTO-001"},
		{"韵达快递", "YUNDA", "95546", "service@yundaex.com", "https://www.yundaex.com", "ACC-YD-001"},
	}
	now := time.Now()
	for _, s := range seeds {
		id := atomic.AddInt64(&r.nextID, 1) - 1
		r.carriers[id] = &domain.Carrier{
			ID:           id,
			TenantID:     1,
			Name:         s.name,
			Code:         s.code,
			ContactPhone: s.phone,
			ContactEmail: s.email,
			Website:      s.website,
			AccountNo:    s.accountNo,
			IsActive:     true,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
	}
}

func (r *MemCarrierRepo) Create(ctx context.Context, c *domain.Carrier) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c.ID = atomic.AddInt64(&r.nextID, 1) - 1
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	r.carriers[c.ID] = c
	return nil
}

func (r *MemCarrierRepo) GetByID(ctx context.Context, id int64) (*domain.Carrier, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.carriers[id]
	if !ok {
		return nil, nil
	}
	return c, nil
}

func (r *MemCarrierRepo) GetByCode(ctx context.Context, tenantID int64, code string) (*domain.Carrier, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, c := range r.carriers {
		if c.TenantID == tenantID && c.Code == code {
			return c, nil
		}
	}
	return nil, nil
}

func (r *MemCarrierRepo) List(ctx context.Context, tenantID int64) ([]domain.Carrier, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.Carrier
	for _, c := range r.carriers {
		if c.TenantID == tenantID {
			result = append(result, *c)
		}
	}
	return result, nil
}

func (r *MemCarrierRepo) Update(ctx context.Context, c *domain.Carrier) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.carriers[c.ID]; !ok {
		return nil
	}
	c.UpdatedAt = time.Now()
	r.carriers[c.ID] = c
	return nil
}

func (r *MemCarrierRepo) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.carriers, id)
	return nil
}

var _ CarrierRepository = (*MemCarrierRepo)(nil)
