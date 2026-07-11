package repository

import (
	"context"
	"github.com/i56/modules/transport/domain"
)

type RouteRepository interface {
	Create(ctx context.Context, r *domain.Route) error
	GetByID(ctx context.Context, tenantID, id int64) (*domain.Route, error)
	List(ctx context.Context, tenantID int64, offset, limit int) ([]domain.Route, int64, error)
	Update(ctx context.Context, r *domain.Route) error
}

type CourierRepository interface {
	Create(ctx context.Context, c *domain.Courier) error
	List(ctx context.Context) ([]domain.Courier, error)
	DetectByTrackingNo(trackingNo string) *domain.Courier
}
