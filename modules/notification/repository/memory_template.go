package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/modules/notification/domain"
)

// MemNotificationTemplateRepo is an in-memory implementation for notification templates.
type MemNotificationTemplateRepo struct {
	mu        sync.RWMutex
	templates map[int64]*domain.NotificationTemplate
	nextID    int64
}

func NewMemNotificationTemplateRepo() *MemNotificationTemplateRepo {
	r := &MemNotificationTemplateRepo{templates: make(map[int64]*domain.NotificationTemplate), nextID: 1}
	r.seed()
	return r
}

func (r *MemNotificationTemplateRepo) seed() {
	now := time.Now()
	for _, t := range domain.DefaultNotificationTemplates() {
		id := atomic.AddInt64(&r.nextID, 1) - 1
		t.ID = id
		t.TenantID = 1
		t.IsActive = true
		t.CreatedAt = now
		t.UpdatedAt = now
		r.templates[id] = &t
	}
}

func (r *MemNotificationTemplateRepo) Create(ctx context.Context, t *domain.NotificationTemplate) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	t.ID = atomic.AddInt64(&r.nextID, 1) - 1
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now
	r.templates[t.ID] = t
	return nil
}

func (r *MemNotificationTemplateRepo) GetByID(ctx context.Context, id int64) (*domain.NotificationTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.templates[id]
	if !ok {
		return nil, nil
	}
	return t, nil
}

func (r *MemNotificationTemplateRepo) GetByCode(ctx context.Context, tenantID int64, code string) (*domain.NotificationTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, t := range r.templates {
		if t.TenantID == tenantID && t.Code == code {
			return t, nil
		}
	}
	return nil, nil
}

func (r *MemNotificationTemplateRepo) List(ctx context.Context, tenantID int64) ([]domain.NotificationTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.NotificationTemplate
	for _, t := range r.templates {
		if t.TenantID == tenantID {
			result = append(result, *t)
		}
	}
	return result, nil
}

func (r *MemNotificationTemplateRepo) ListByType(ctx context.Context, tenantID int64, ntype domain.NotificationType) ([]domain.NotificationTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.NotificationTemplate
	for _, t := range r.templates {
		if t.TenantID == tenantID && t.Type == ntype {
			result = append(result, *t)
		}
	}
	return result, nil
}

func (r *MemNotificationTemplateRepo) Update(ctx context.Context, t *domain.NotificationTemplate) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.templates[t.ID]; !ok {
		return nil
	}
	t.UpdatedAt = time.Now()
	r.templates[t.ID] = t
	return nil
}

func (r *MemNotificationTemplateRepo) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.templates, id)
	return nil
}

var _ NotificationTemplateRepository = (*MemNotificationTemplateRepo)(nil)
