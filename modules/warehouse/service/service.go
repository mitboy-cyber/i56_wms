package service

import (
	"context"
	"time"

	"github.com/i56/framework/core/errors"
	"github.com/i56/modules/warehouse/domain"
	"github.com/i56/modules/warehouse/repository"
)

type WarehouseService struct {
	repo repository.WarehouseRepository
}

func NewWarehouseService(repo repository.WarehouseRepository) *WarehouseService {
	return &WarehouseService{repo: repo}
}

func (s *WarehouseService) Create(ctx context.Context, tenantID int64, name, code, address, contact, phone string) (*domain.Warehouse, error) {
	w := &domain.Warehouse{
		TenantID: tenantID, Name: name, Code: code,
		Address: address, Contact: contact, Phone: phone,
		IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	if err := s.repo.Create(ctx, tenantID, w); err != nil {
		return nil, err
	}
	return w, nil
}

func (s *WarehouseService) GetByID(ctx context.Context, tenantID, id int64) (*domain.Warehouse, error) {
	w, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil { return nil, err }
	if w == nil { return nil, errors.NewNotFound("Warehouse") }
	return w, nil
}

func (s *WarehouseService) List(ctx context.Context, tenantID int64, offset, limit int) ([]domain.Warehouse, int64, error) {
	return s.repo.List(ctx, tenantID, offset, limit)
}
