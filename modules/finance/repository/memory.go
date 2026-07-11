package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/modules/finance/domain"
)

type MemLedgerRepo struct {
	mu      sync.RWMutex
	entries map[int64]*domain.LedgerEntry
	nextID  int64
}

func NewMemLedgerRepo() *MemLedgerRepo {
	return &MemLedgerRepo{entries: make(map[int64]*domain.LedgerEntry), nextID: 1}
}

func (r *MemLedgerRepo) Add(ctx context.Context, entry *domain.LedgerEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	entry.ID = atomic.AddInt64(&r.nextID, 1) - 1
	entry.CreatedAt = time.Now()
	r.entries[entry.ID] = entry
	return nil
}

func (r *MemLedgerRepo) List(ctx context.Context, tenantID, clientID int64, offset, limit int) ([]domain.LedgerEntry, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.LedgerEntry
	for _, e := range r.entries {
		if e.TenantID == tenantID && e.ClientID == clientID {
			result = append(result, *e)
		}
	}
	total := int64(len(result))
	if offset >= int(total) { return nil, total, nil }
	end := offset + limit
	if end > int(total) { end = int(total) }
	return result[offset:end], total, nil
}

var _ LedgerRepository = (*MemLedgerRepo)(nil)
