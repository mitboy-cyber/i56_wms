package repository

import (
	"context"
	"github.com/i56/modules/parcel/domain"
)

type ParcelRepository interface {
	Create(ctx context.Context, p *domain.Parcel) error
	GetByID(ctx context.Context, tenantID, id int64) (*domain.Parcel, error)
	GetByTrackingNo(ctx context.Context, tenantID int64, trackingNo string) (*domain.Parcel, error)
	List(ctx context.Context, tenantID int64, offset, limit int) ([]domain.Parcel, int64, error)
	ListByClient(ctx context.Context, tenantID, clientID int64, offset, limit int) ([]domain.Parcel, int64, error)
	Update(ctx context.Context, p *domain.Parcel) error
}
