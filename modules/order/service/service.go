package service

import (
	"context"
	"fmt"
	"time"

	"github.com/i56/framework/core/errors"
	"github.com/i56/modules/order/domain"
	"github.com/i56/modules/order/repository"
)

type OrderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) Create(ctx context.Context, o *domain.Order) (*domain.Order, error) {
	o.OrderNo = fmt.Sprintf("%s%08d", time.Now().Format("20060102150405"), time.Now().UnixNano()%100000000)
	o.Status = domain.StatusPendingPicking
	o.ParcelCount = 1 // simplified
	if err := s.repo.Create(ctx, o); err != nil { return nil, err }
	return o, nil
}

func (s *OrderService) Cancel(ctx context.Context, tenantID, id int64) error {
	o, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil || o == nil { return errors.NewNotFound("Order") }
	if !o.IsCancellable() {
		return errors.NewValidation("only pending_picking orders can be cancelled")
	}
	o.Status = domain.StatusCancelled
	return s.repo.Update(ctx, o)
}

func (s *OrderService) Transition(ctx context.Context, tenantID, id int64, target domain.OrderStatus) error {
	o, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil || o == nil { return errors.NewNotFound("Order") }
	if !o.CanTransitionTo(target) {
		return errors.NewInvalidTransition(string(o.Status), string(target))
	}
	o.Status = target
	return s.repo.Update(ctx, o)
}

func (s *OrderService) GetByID(ctx context.Context, tenantID, id int64) (*domain.Order, error) {
	o, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil { return nil, err }
	if o == nil { return nil, errors.NewNotFound("Order") }
	return o, nil
}

func (s *OrderService) List(ctx context.Context, tenantID int64, offset, limit int) ([]domain.Order, int64, error) {
	return s.repo.List(ctx, tenantID, offset, limit)
}

func (s *OrderService) ListByClient(ctx context.Context, tenantID, clientID int64, offset, limit int) ([]domain.Order, int64, error) {
	return s.repo.ListByClient(ctx, tenantID, clientID, offset, limit)
}
