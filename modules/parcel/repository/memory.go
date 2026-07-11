package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/modules/parcel/domain"
)

type MemParcelRepo struct {
	mu      sync.RWMutex
	parcels map[int64]*domain.Parcel
	nextID  int64
}

func NewMemParcelRepo() *MemParcelRepo {
	return &MemParcelRepo{parcels: make(map[int64]*domain.Parcel), nextID: 1}
}

func (r *MemParcelRepo) Create(ctx context.Context, p *domain.Parcel) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p.ID = atomic.AddInt64(&r.nextID, 1) - 1
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	r.parcels[p.ID] = p
	return nil
}

func (r *MemParcelRepo) GetByID(ctx context.Context, tenantID, id int64) (*domain.Parcel, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.parcels[id]
	if !ok || p.TenantID != tenantID { return nil, nil }
	return p, nil
}

func (r *MemParcelRepo) GetByTrackingNo(ctx context.Context, tenantID int64, trackingNo string) (*domain.Parcel, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.parcels {
		if p.TenantID == tenantID && p.TrackingNumber == trackingNo { return p, nil }
	}
	return nil, nil
}

func (r *MemParcelRepo) List(ctx context.Context, tenantID int64, offset, limit int) ([]domain.Parcel, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return filterAndPaginate(r.parcels, func(p *domain.Parcel) bool { return p.TenantID == tenantID }, offset, limit)
}

func (r *MemParcelRepo) ListByClient(ctx context.Context, tenantID, clientID int64, offset, limit int) ([]domain.Parcel, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return filterAndPaginate(r.parcels, func(p *domain.Parcel) bool {
		return p.TenantID == tenantID && p.ClientID == clientID
	}, offset, limit)
}

func (r *MemParcelRepo) Update(ctx context.Context, p *domain.Parcel) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p.UpdatedAt = time.Now()
	r.parcels[p.ID] = p
	return nil
}

func filterAndPaginate(m map[int64]*domain.Parcel, fn func(*domain.Parcel) bool, offset, limit int) ([]domain.Parcel, int64, error) {
	var result []domain.Parcel
	for _, p := range m {
		if fn(p) { result = append(result, *p) }
	}
	total := int64(len(result))
	if offset >= int(total) { return nil, total, nil }
	end := offset + limit
	if end > int(total) { end = int(total) }
	return result[offset:end], total, nil
}

var _ ParcelRepository = (*MemParcelRepo)(nil)
