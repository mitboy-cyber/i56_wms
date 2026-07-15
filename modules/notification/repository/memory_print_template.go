package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/modules/notification/domain"
)

// MemPrintTemplateRepo is an in-memory implementation for print templates.
type MemPrintTemplateRepo struct {
	mu        sync.RWMutex
	templates map[int64]*domain.PrintTemplate
	nextID    int64
}

func NewMemPrintTemplateRepo() *MemPrintTemplateRepo {
	r := &MemPrintTemplateRepo{templates: make(map[int64]*domain.PrintTemplate), nextID: 1}
	r.seed()
	return r
}

func (r *MemPrintTemplateRepo) seed() {
	now := time.Now()
	for _, t := range domain.DefaultPrintTemplates() {
		id := atomic.AddInt64(&r.nextID, 1) - 1
		t.ID = id
		t.TenantID = 1
		t.IsActive = true
		t.CreatedAt = now
		t.UpdatedAt = now
		r.templates[id] = &t
	}
}

func (r *MemPrintTemplateRepo) Create(ctx context.Context, t *domain.PrintTemplate) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	t.ID = atomic.AddInt64(&r.nextID, 1) - 1
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now
	r.templates[t.ID] = t
	return nil
}

func (r *MemPrintTemplateRepo) GetByID(ctx context.Context, id int64) (*domain.PrintTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.templates[id]
	if !ok {
		return nil, nil
	}
	return t, nil
}

func (r *MemPrintTemplateRepo) GetByCode(ctx context.Context, tenantID int64, code string) (*domain.PrintTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, t := range r.templates {
		if t.TenantID == tenantID && t.Code == code {
			return t, nil
		}
	}
	return nil, nil
}

func (r *MemPrintTemplateRepo) List(ctx context.Context, tenantID int64) ([]domain.PrintTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.PrintTemplate
	for _, t := range r.templates {
		if t.TenantID == tenantID {
			result = append(result, *t)
		}
	}
	return result, nil
}

func (r *MemPrintTemplateRepo) ListByType(ctx context.Context, tenantID int64, ptype domain.PrintTemplateType) ([]domain.PrintTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.PrintTemplate
	for _, t := range r.templates {
		if t.TenantID == tenantID && t.Type == ptype {
			result = append(result, *t)
		}
	}
	return result, nil
}

func (r *MemPrintTemplateRepo) SetDefault(ctx context.Context, tenantID, id int64, ptype domain.PrintTemplateType) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, t := range r.templates {
		if t.TenantID == tenantID && t.Type == ptype {
			t.IsDefault = false
		}
	}
	if t, ok := r.templates[id]; ok && t.TenantID == tenantID {
		t.IsDefault = true
	}
	return nil
}

func (r *MemPrintTemplateRepo) Update(ctx context.Context, t *domain.PrintTemplate) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.templates[t.ID]; !ok {
		return nil
	}
	t.UpdatedAt = time.Now()
	r.templates[t.ID] = t
	return nil
}

func (r *MemPrintTemplateRepo) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.templates, id)
	return nil
}

var _ PrintTemplateRepository = (*MemPrintTemplateRepo)(nil)
