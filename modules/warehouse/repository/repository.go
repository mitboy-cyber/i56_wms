package repository

import (
	"context"
	"github.com/i56/modules/warehouse/domain"
)

type WarehouseRepository interface {
	Create(ctx context.Context, tenantID int64, w *domain.Warehouse) error
	GetByID(ctx context.Context, tenantID, id int64) (*domain.Warehouse, error)
	List(ctx context.Context, tenantID int64, offset, limit int) ([]domain.Warehouse, int64, error)
	Update(ctx context.Context, tenantID, id int64, w *domain.Warehouse) error
}

// ContainerRepository defines persistence operations for containers.
type ContainerRepository interface {
	Create(ctx context.Context, c *domain.Container) error
	GetByID(ctx context.Context, id int64) (*domain.Container, error)
	GetByContainerNo(ctx context.Context, containerNo string) (*domain.Container, error)
	ListByWarehouse(ctx context.Context, warehouseID int64) ([]domain.Container, error)
	ListByStatus(ctx context.Context, warehouseID int64, status domain.ContainerStatus) ([]domain.Container, error)
	ListByRoute(ctx context.Context, routeID int64) ([]domain.Container, error)
	Update(ctx context.Context, c *domain.Container) error
	Delete(ctx context.Context, id int64) error
}
