// Package service provides business logic services for the WMS backend.
// These services delegate to the github.com/i56/modules repositories.
package service

import (
	"context"

	orderRepo "github.com/i56/modules/order/repository"
	orderSvc "github.com/i56/modules/order/service"
)

// OrderService wraps the order service with tenant-aware business logic.
type OrderService struct {
	svc *orderSvc.OrderService
}

// NewOrderService creates a new OrderService.
func NewOrderService(or *orderRepo.MemOrderRepo) *OrderService {
	return &OrderService{svc: orderSvc.NewOrderService(or)}
}

// List returns all orders for a tenant.
func (s *OrderService) List(ctx context.Context, tenantID int64) ([]interface{}, int64, error) {
	orders, total, err := s.svc.List(ctx, tenantID, 0, 200)
	if err != nil {
		return nil, 0, err
	}
	result := make([]interface{}, len(orders))
	for i, o := range orders {
		result[i] = o
	}
	return result, total, nil
}
