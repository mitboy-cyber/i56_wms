package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/modules/pda/domain"
)

// MemWorkOrderTemplateRepo is an in-memory implementation for work order templates.
type MemWorkOrderTemplateRepo struct {
	mu        sync.RWMutex
	templates map[int64]*domain.WorkOrderTemplate
	nextID    int64
}

func NewMemWorkOrderTemplateRepo() *MemWorkOrderTemplateRepo {
	r := &MemWorkOrderTemplateRepo{templates: make(map[int64]*domain.WorkOrderTemplate), nextID: 1}
	r.seed()
	return r
}

func (r *MemWorkOrderTemplateRepo) seed() {
	now := time.Now()
	for _, t := range domain.DefaultWorkOrderTemplates() {
		id := atomic.AddInt64(&r.nextID, 1) - 1
		t.ID = id
		t.TenantID = 1
		t.IsActive = true
		t.CreatedAt = now
		t.UpdatedAt = now
		r.templates[id] = &t
	}
}

func (r *MemWorkOrderTemplateRepo) Create(ctx context.Context, t *domain.WorkOrderTemplate) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	t.ID = atomic.AddInt64(&r.nextID, 1) - 1
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now
	r.templates[t.ID] = t
	return nil
}

func (r *MemWorkOrderTemplateRepo) GetByID(ctx context.Context, id int64) (*domain.WorkOrderTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.templates[id]
	if !ok {
		return nil, nil
	}
	return t, nil
}

func (r *MemWorkOrderTemplateRepo) GetByCode(ctx context.Context, tenantID int64, code string) (*domain.WorkOrderTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, t := range r.templates {
		if t.TenantID == tenantID && t.Code == code {
			return t, nil
		}
	}
	return nil, nil
}

func (r *MemWorkOrderTemplateRepo) List(ctx context.Context, tenantID int64) ([]domain.WorkOrderTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.WorkOrderTemplate
	for _, t := range r.templates {
		if t.TenantID == tenantID {
			result = append(result, *t)
		}
	}
	return result, nil
}

func (r *MemWorkOrderTemplateRepo) ListByType(ctx context.Context, tenantID int64, wtype string) ([]domain.WorkOrderTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.WorkOrderTemplate
	for _, t := range r.templates {
		if t.TenantID == tenantID && t.Type == wtype {
			result = append(result, *t)
		}
	}
	return result, nil
}

func (r *MemWorkOrderTemplateRepo) Update(ctx context.Context, t *domain.WorkOrderTemplate) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.templates[t.ID]; !ok {
		return nil
	}
	t.UpdatedAt = time.Now()
	r.templates[t.ID] = t
	return nil
}

func (r *MemWorkOrderTemplateRepo) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.templates, id)
	return nil
}

var _ WorkOrderTemplateRepository = (*MemWorkOrderTemplateRepo)(nil)
