package repository

import (
	"context"
	"github.com/i56/modules/tms/domain"
)

// CargoTypeRepository defines persistence operations for cargo types.
type CargoTypeRepository interface {
	Create(ctx context.Context, c *domain.CargoType) error
	GetByID(ctx context.Context, id int64) (*domain.CargoType, error)
	GetByCode(ctx context.Context, tenantID int64, code string) (*domain.CargoType, error)
	List(ctx context.Context, tenantID int64) ([]domain.CargoType, error)
	Update(ctx context.Context, c *domain.CargoType) error
}

// CarrierRepository defines persistence operations for carriers.
type CarrierRepository interface {
	Create(ctx context.Context, c *domain.Carrier) error
	GetByID(ctx context.Context, id int64) (*domain.Carrier, error)
	GetByCode(ctx context.Context, tenantID int64, code string) (*domain.Carrier, error)
	List(ctx context.Context, tenantID int64) ([]domain.Carrier, error)
	Update(ctx context.Context, c *domain.Carrier) error
	Delete(ctx context.Context, id int64) error
}

// CourierRepository defines persistence operations for couriers.
type CourierRepository interface {
	Create(ctx context.Context, c *domain.Courier) error
	GetByID(ctx context.Context, id int64) (*domain.Courier, error)
	ListByCarrier(ctx context.Context, carrierID int64) ([]domain.Courier, error)
	List(ctx context.Context, tenantID int64) ([]domain.Courier, error)
	Update(ctx context.Context, c *domain.Courier) error
	Delete(ctx context.Context, id int64) error
}

// RouteRepository defines persistence operations for routes.
type RouteRepository interface {
	Create(ctx context.Context, route *domain.Route) error
	GetByID(ctx context.Context, id int64) (*domain.Route, error)
	GetByCode(ctx context.Context, tenantID int64, code string) (*domain.Route, error)
	List(ctx context.Context, tenantID int64) ([]domain.Route, error)
	ListByRouteType(ctx context.Context, tenantID int64, routeType domain.RouteType) ([]domain.Route, error)
	Update(ctx context.Context, route *domain.Route) error
	Delete(ctx context.Context, id int64) error
}
