// Package service provides business logic services for the WMS backend.
package service

import (
	"context"

	parcelRepo "github.com/i56/modules/parcel/repository"
	parcelSvc "github.com/i56/modules/parcel/service"
	parcelDomain "github.com/i56/modules/parcel/domain"
)

// ParcelService wraps the parcel service with tenant-aware business logic.
type ParcelService struct {
	svc *parcelSvc.ParcelService
}

// NewParcelService creates a new ParcelService.
func NewParcelService(pr *parcelRepo.MemParcelRepo) *ParcelService {
	return &ParcelService{svc: parcelSvc.NewParcelService(pr)}
}

// List returns all parcels for a tenant.
func (s *ParcelService) List(ctx context.Context, tenantID int64) ([]interface{}, int64, error) {
	parcels, total, err := s.svc.List(ctx, tenantID, 0, 200)
	if err != nil {
		return nil, 0, err
	}
	result := make([]interface{}, len(parcels))
	for i, p := range parcels {
		result[i] = p
	}
	return result, total, nil
}

// PreDeclare pre-declares a parcel.
func (s *ParcelService) PreDeclare(ctx context.Context, p *parcelDomain.Parcel) (interface{}, error) {
	return s.svc.PreDeclare(ctx, p)
}
