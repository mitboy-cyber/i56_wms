package repository

import (
	"context"
	"github.com/i56/modules/finance/domain"
)

// ClientLedgerRepository defines persistence for client ledgers.
type ClientLedgerRepository interface {
	GetByClient(ctx context.Context, tenantID, clientID int64) (*domain.ClientLedger, error)
	AddEntry(ctx context.Context, entry *domain.LedgerEntry) error
	ListEntries(ctx context.Context, tenantID, clientID int64, offset, limit int) ([]domain.LedgerEntry, int64, error)
	Update(ctx context.Context, ledger *domain.ClientLedger) error
}

// RechargeRepository defines persistence for recharges.
type RechargeRepository interface {
	Create(ctx context.Context, recharge *domain.Recharge) error
	GetByID(ctx context.Context, id int64) (*domain.Recharge, error)
	GetByRechargeNo(ctx context.Context, tenantID int64, rechNo string) (*domain.Recharge, error)
	ListByClient(ctx context.Context, tenantID, clientID int64) ([]domain.Recharge, error)
	ListByStatus(ctx context.Context, tenantID int64, status domain.RechargeStatus) ([]domain.Recharge, error)
	Update(ctx context.Context, recharge *domain.Recharge) error
}

// RechargeLogRepository defines persistence for recharge audit logs.
type RechargeLogRepository interface {
	Create(ctx context.Context, log *domain.RechargeLog) error
	ListByRecharge(ctx context.Context, rechargeID int64) ([]domain.RechargeLog, error)
}

// StatementRepository defines persistence for statements.
type StatementRepository interface {
	Create(ctx context.Context, s *domain.Statement) error
	GetByID(ctx context.Context, id int64) (*domain.Statement, error)
	ListByClient(ctx context.Context, tenantID, clientID int64) ([]domain.Statement, error)
	GetByStatementNo(ctx context.Context, tenantID int64, stmtNo string) (*domain.Statement, error)
	Update(ctx context.Context, s *domain.Statement) error
}

// OrderProfitRepository defines persistence for order profits.
type OrderProfitRepository interface {
	Create(ctx context.Context, p *domain.OrderProfit) error
	GetByOrderID(ctx context.Context, orderID int64) (*domain.OrderProfit, error)
	ListByClient(ctx context.Context, tenantID, clientID int64) ([]domain.OrderProfit, error)
	ListByRoute(ctx context.Context, routeID int64) ([]domain.OrderProfit, error)
}

// ServiceProfitRepository defines persistence for service profits.
type ServiceProfitRepository interface {
	Create(ctx context.Context, p *domain.ServiceProfit) error
	GetByServiceOrderID(ctx context.Context, serviceOrderID int64) (*domain.ServiceProfit, error)
	ListByOrder(ctx context.Context, orderID int64) ([]domain.ServiceProfit, error)
}

// ClientProfitRepository defines persistence for client profit aggregation.
type ClientProfitRepository interface {
	Create(ctx context.Context, p *domain.ClientProfit) error
	GetByClientAndPeriod(ctx context.Context, tenantID, clientID int64, period string) (*domain.ClientProfit, error)
	ListByPeriod(ctx context.Context, tenantID int64, period string) ([]domain.ClientProfit, error)
}

// RouteProfitRepository defines persistence for route profit aggregation.
type RouteProfitRepository interface {
	Create(ctx context.Context, p *domain.RouteProfit) error
	GetByRouteAndPeriod(ctx context.Context, routeID int64, period string) (*domain.RouteProfit, error)
	ListByPeriod(ctx context.Context, tenantID int64, period string) ([]domain.RouteProfit, error)
}
