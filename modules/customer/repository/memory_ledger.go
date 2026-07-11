package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type LedgerEntry struct {
	ID           int64     `json:"id"`
	TenantID     int64     `json:"tenant_id"`
	ClientID     int64     `json:"client_id"`
	Amount       float64   `json:"amount"`
	BalanceAfter float64   `json:"balance_after"`
	Type         string    `json:"type"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
}

type MemLedgerRepo struct {
	mu      sync.RWMutex
	entries map[int64]*LedgerEntry
	nextID  int64
}

func NewMemLedgerRepo() *MemLedgerRepo {
	return &MemLedgerRepo{entries: make(map[int64]*LedgerEntry), nextID: 1}
}

func (r *MemLedgerRepo) Add(ctx context.Context, entry *LedgerEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	entry.ID = atomic.AddInt64(&r.nextID, 1) - 1
	entry.CreatedAt = time.Now()
	r.entries[entry.ID] = entry
	return nil
}

func (r *MemLedgerRepo) List(ctx context.Context, tenantID, clientID int64, offset, limit int) ([]LedgerEntry, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []LedgerEntry
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
func (r *MemLedgerRepo) GetByEntryID(ctx context.Context, id int64) (*LedgerEntry, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	e, ok := r.entries[id]
	if !ok { return nil, nil }
	return e, nil
}

func (r *MemLedgerRepo) Update(ctx context.Context, id int64, entry *LedgerEntry) error {
	r.mu.Lock(); defer r.mu.Unlock()
	if _, ok := r.entries[id]; !ok { return nil }
	entry.ID = id
	r.entries[id] = entry
	return nil
}
func (r *MemLedgerRepo) GetByClient(ctx context.Context, tenantID, clientID int64) []LedgerEntry {
	r.mu.RLock(); defer r.mu.RUnlock()
	var result []LedgerEntry
	for _, e := range r.entries { if e.TenantID==tenantID && e.ClientID==clientID { result=append(result, *e) } }
	return result
}
