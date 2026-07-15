package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/modules/finance/domain"
)

// MemClientLedgerRepo is an in-memory implementation for client ledgers.
type MemClientLedgerRepo struct {
	mu       sync.RWMutex
	ledgers  map[int64]*domain.ClientLedger
	entries  map[int64]*domain.LedgerEntry
	nextID   int64
}

func NewMemClientLedgerRepo() *MemClientLedgerRepo {
	r := &MemClientLedgerRepo{
		ledgers: make(map[int64]*domain.ClientLedger),
		entries: make(map[int64]*domain.LedgerEntry),
		nextID:  1,
	}
	r.seed()
	return r
}

func (r *MemClientLedgerRepo) seed() {
	now := time.Now()
	seeds := []struct {
		clientID    int64
		balance     float64
		totalRech   float64
		totalSpent  float64
		creditLimit float64
	}{
		{1, 5000.00, 10000.00, 5000.00, 20000.00},
		{2, 3200.00, 8000.00, 4800.00, 15000.00},
		{3, 1500.00, 3000.00, 1500.00, 5000.00},
	}
	for _, s := range seeds {
		id := atomic.AddInt64(&r.nextID, 1) - 1
		r.ledgers[id] = &domain.ClientLedger{
			ID:             id,
			TenantID:       1,
			ClientID:       s.clientID,
			Balance:        s.balance,
			TotalRecharged: s.totalRech,
			TotalSpent:     s.totalSpent,
			CreditLimit:    s.creditLimit,
			IsActive:       true,
			CreatedAt:      now,
			UpdatedAt:      now,
		}
	}
}

func (r *MemClientLedgerRepo) GetByClient(ctx context.Context, tenantID, clientID int64) (*domain.ClientLedger, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, l := range r.ledgers {
		if l.TenantID == tenantID && l.ClientID == clientID {
			return l, nil
		}
	}
	return nil, nil
}

func (r *MemClientLedgerRepo) AddEntry(ctx context.Context, entry *domain.LedgerEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	entry.ID = atomic.AddInt64(&r.nextID, 1) - 1
	entry.CreatedAt = time.Now()
	r.entries[entry.ID] = entry
	// Update balance
	for _, l := range r.ledgers {
		if l.TenantID == entry.TenantID && l.ClientID == entry.ClientID {
			l.Balance = entry.BalanceAfter
			l.TotalSpent += entry.Amount
			now := time.Now()
			l.LastDeductionAt = &now
			l.UpdatedAt = now
		}
	}
	return nil
}

func (r *MemClientLedgerRepo) ListEntries(ctx context.Context, tenantID, clientID int64, offset, limit int) ([]domain.LedgerEntry, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.LedgerEntry
	for _, e := range r.entries {
		if e.TenantID == tenantID && e.ClientID == clientID {
			result = append(result, *e)
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

func (r *MemClientLedgerRepo) Update(ctx context.Context, ledger *domain.ClientLedger) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.ledgers[ledger.ID]; !ok {
		return nil
	}
	ledger.UpdatedAt = time.Now()
	r.ledgers[ledger.ID] = ledger
	return nil
}

var _ ClientLedgerRepository = (*MemClientLedgerRepo)(nil)
