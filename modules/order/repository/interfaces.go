package repository

import (
	"context"
	"github.com/i56/modules/order/domain"
)

// ConsolidationOrderRepository defines persistence operations for consolidation orders.
type ConsolidationOrderRepository interface {
	Create(ctx context.Context, o *domain.ConsolidationOrder) error
	GetByID(ctx context.Context, tenantID, id int64) (*domain.ConsolidationOrder, error)
	GetByOrderNo(ctx context.Context, tenantID int64, orderNo string) (*domain.ConsolidationOrder, error)
	List(ctx context.Context, tenantID int64, offset, limit int) ([]domain.ConsolidationOrder, int64, error)
	ListByClient(ctx context.Context, tenantID, clientID int64) ([]domain.ConsolidationOrder, error)
	ListByMember(ctx context.Context, tenantID, memberID int64) ([]domain.ConsolidationOrder, error)
	Update(ctx context.Context, o *domain.ConsolidationOrder) error
	Delete(ctx context.Context, tenantID, id int64) error
}

// ServiceOrderRepository defines persistence operations for service orders.
type ServiceOrderRepository interface {
	Create(ctx context.Context, o *domain.ServiceOrder) error
	GetByID(ctx context.Context, tenantID, id int64) (*domain.ServiceOrder, error)
	ListByParcel(ctx context.Context, tenantID, parcelID int64) ([]domain.ServiceOrder, error)
	ListByOrder(ctx context.Context, tenantID, orderID int64) ([]domain.ServiceOrder, error)
	ListByClient(ctx context.Context, tenantID, clientID int64) ([]domain.ServiceOrder, error)
	Update(ctx context.Context, o *domain.ServiceOrder) error
}

// OrderParcelRepository defines persistence operations for the OrderParcel join entity.
type OrderParcelRepository interface {
	Create(ctx context.Context, op *domain.OrderParcel) error
	GetByID(ctx context.Context, id int64) (*domain.OrderParcel, error)
	ListByOrder(ctx context.Context, orderID int64) ([]domain.OrderParcel, error)
	ListByParcel(ctx context.Context, parcelID int64) ([]domain.OrderParcel, error)
	Update(ctx context.Context, op *domain.OrderParcel) error
	Delete(ctx context.Context, id int64) error
}
