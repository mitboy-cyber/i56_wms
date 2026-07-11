package service

import (
	"context"

	"github.com/i56/framework/core/errors"
	"github.com/i56/modules/parcel/domain"
	"github.com/i56/modules/parcel/repository"
)

type ParcelService struct {
	repo repository.ParcelRepository
}

func NewParcelService(repo repository.ParcelRepository) *ParcelService {
	return &ParcelService{repo: repo}
}

func (s *ParcelService) PreDeclare(ctx context.Context, p *domain.Parcel) (*domain.Parcel, error) {
	existing, _ := s.repo.GetByTrackingNo(ctx, p.TenantID, p.TrackingNumber)
	if existing != nil {
		return nil, errors.NewConflict("tracking number already pre-declared")
	}
	p.Status = domain.StatusPreDeclared
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *ParcelService) Receive(ctx context.Context, tenantID, id int64, weight, length, width, height float64) (*domain.Parcel, error) {
	p, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil || p == nil { return nil, errors.NewNotFound("Parcel") }
	if !p.CanTransitionTo(domain.StatusReceived) {
		return nil, errors.NewInvalidTransition(string(p.Status), string(domain.StatusReceived))
	}
	p.ActualWeight = weight
	p.Length = length
	p.Width = width
	p.Height = height
	p.Status = domain.StatusReceived
	if err := s.repo.Update(ctx, p); err != nil { return nil, err }
	return p, nil
}

func (s *ParcelService) GetByID(ctx context.Context, tenantID, id int64) (*domain.Parcel, error) {
	p, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil { return nil, err }
	if p == nil { return nil, errors.NewNotFound("Parcel") }
	return p, nil
}

func (s *ParcelService) List(ctx context.Context, tenantID int64, offset, limit int) ([]domain.Parcel, int64, error) {
	return s.repo.List(ctx, tenantID, offset, limit)
}

func (s *ParcelService) ListByClient(ctx context.Context, tenantID, clientID int64, offset, limit int) ([]domain.Parcel, int64, error) {
	return s.repo.ListByClient(ctx, tenantID, clientID, offset, limit)
}
