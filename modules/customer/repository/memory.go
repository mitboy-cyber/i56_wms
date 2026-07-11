package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/framework/core/errors"
	"github.com/i56/modules/customer/domain"
)

// MemClientRepo is an in-memory implementation of ClientRepository.
type MemClientRepo struct {
	mu      sync.RWMutex
	clients map[int64]*domain.Client
	nextID  int64
}

func NewMemClientRepo() *MemClientRepo {
	return &MemClientRepo{clients: make(map[int64]*domain.Client), nextID: 1}
}

func (r *MemClientRepo) Create(ctx context.Context, tenantID int64, c *domain.Client) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c.ID = atomic.AddInt64(&r.nextID, 1) - 1
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	r.clients[c.ID] = c
	return nil
}

func (r *MemClientRepo) GetByID(ctx context.Context, tenantID, id int64) (*domain.Client, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.clients[id]
	if !ok || c.TenantID != tenantID {
		return nil, nil
	}
	return c, nil
}

func (r *MemClientRepo) GetByCode(ctx context.Context, tenantID int64, code string) (*domain.Client, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, c := range r.clients {
		if c.TenantID == tenantID && c.Code == code {
			return c, nil
		}
	}
	return nil, nil
}

func (r *MemClientRepo) List(ctx context.Context, tenantID int64, offset, limit int) ([]domain.Client, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.Client
	for _, c := range r.clients {
		if c.TenantID == tenantID {
			result = append(result, *c)
		}
	}
	total := int64(len(result))
	if offset >= int(total) {
		return nil, total, nil
	}
	end := offset + limit
	if end > int(total) {
		end = int(total)
	}
	return result[offset:end], total, nil
}

func (r *MemClientRepo) Update(ctx context.Context, tenantID, id int64, c *domain.Client) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.clients[id]
	if !ok || existing.TenantID != tenantID {
		return errors.NewNotFound("Client")
	}
	c.UpdatedAt = time.Now()
	r.clients[id] = c
	return nil
}

func (r *MemClientRepo) Delete(ctx context.Context, tenantID, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.clients[id]
	if !ok || existing.TenantID != tenantID {
		return errors.NewNotFound("Client")
	}
	delete(r.clients, id)
	return nil
}

// Ensure MemClientRepo implements ClientRepository.
var _ ClientRepository = (*MemClientRepo)(nil)
