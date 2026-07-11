package repository

import (
	"context"
	"github.com/i56/modules/order/domain"
)

type OrderRepository interface {
	Create(ctx context.Context, o *domain.Order) error
	GetByID(ctx context.Context, tenantID, id int64) (*domain.Order, error)
	GetByOrderNo(ctx context.Context, tenantID int64, orderNo string) (*domain.Order, error)
	List(ctx context.Context, tenantID int64, offset, limit int) ([]domain.Order, int64, error)
	ListByClient(ctx context.Context, tenantID, clientID int64, offset, limit int) ([]domain.Order, int64, error)
	Update(ctx context.Context, o *domain.Order) error
}
