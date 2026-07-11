package repository

import (
	"context"
	"github.com/i56/modules/finance/domain"
)

type LedgerRepository interface {
	Add(ctx context.Context, entry *domain.LedgerEntry) error
	List(ctx context.Context, tenantID, clientID int64, offset, limit int) ([]domain.LedgerEntry, int64, error)
}
