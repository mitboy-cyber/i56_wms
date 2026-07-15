package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/modules/finance/domain"
)

// MemStatementRepo is an in-memory implementation for statements.
type MemStatementRepo struct {
	mu         sync.RWMutex
	statements map[int64]*domain.Statement
	nextID     int64
}

func NewMemStatementRepo() *MemStatementRepo {
	r := &MemStatementRepo{statements: make(map[int64]*domain.Statement), nextID: 1}
	r.seed()
	return r
}

func (r *MemStatementRepo) seed() {
	now := time.Now()
	seeds := []struct {
		clientID       int64
		statementNo    string
		startDate      string
		endDate        string
		openingBalance float64
		totalRech      float64
		totalShipFee   float64
		totalServFee   float64
		totalAdjust    float64
		closingBalance float64
		status         string
	}{
		{1, "STM202401001", "2024-01-01", "2024-01-15", 2000.00, 5000.00, 1200.00, 150.00, 0, 5650.00, "closed"},
		{2, "STM202401002", "2024-01-01", "2024-01-15", 1500.00, 8000.00, 3500.00, 200.00, 0, 5800.00, "closed"},
		{3, "STM202401003", "2024-01-01", "2024-01-15", 500.00, 2000.00, 800.00, 50.00, 0, 1650.00, "closed"},
		{1, "STM202401004", "2024-01-16", "2024-01-31", 5650.00, 0, 450.00, 100.00, -50.00, 5050.00, "open"},
	}
	for _, s := range seeds {
		id := atomic.AddInt64(&r.nextID, 1) - 1
		r.statements[id] = &domain.Statement{
			ID:         id,
			TenantID:   1,
			ClientID:   s.clientID,
			StatementNo: s.statementNo,
			Period: domain.StatementPeriod{
				StartDate: s.startDate,
				EndDate:   s.endDate,
			},
			OpeningBalance:   s.openingBalance,
			TotalRecharged:   s.totalRech,
			TotalShippingFee: s.totalShipFee,
			TotalServiceFee:  s.totalServFee,
			TotalAdjustments: s.totalAdjust,
			ClosingBalance:   s.closingBalance,
			Status:           s.status,
			GeneratedAt:      now,
			CreatedAt:        now,
		}
	}
}

func (r *MemStatementRepo) Create(ctx context.Context, s *domain.Statement) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	s.ID = atomic.AddInt64(&r.nextID, 1) - 1
	s.CreatedAt = time.Now()
	r.statements[s.ID] = s
	return nil
}

func (r *MemStatementRepo) GetByID(ctx context.Context, id int64) (*domain.Statement, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.statements[id]
	if !ok {
		return nil, nil
	}
	return s, nil
}

func (r *MemStatementRepo) ListByClient(ctx context.Context, tenantID, clientID int64) ([]domain.Statement, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.Statement
	for _, s := range r.statements {
		if s.TenantID == tenantID && s.ClientID == clientID {
			result = append(result, *s)
		}
	}
	return result, nil
}

func (r *MemStatementRepo) GetByStatementNo(ctx context.Context, tenantID int64, stmtNo string) (*domain.Statement, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, s := range r.statements {
		if s.TenantID == tenantID && s.StatementNo == stmtNo {
			return s, nil
		}
	}
	return nil, nil
}

func (r *MemStatementRepo) Update(ctx context.Context, s *domain.Statement) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.statements[s.ID]; !ok {
		return nil
	}
	r.statements[s.ID] = s
	return nil
}

var _ StatementRepository = (*MemStatementRepo)(nil)
