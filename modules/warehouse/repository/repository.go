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
