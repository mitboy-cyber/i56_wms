package repository

import (
	"context"
	"github.com/i56/modules/pda/domain"
)

// WorkOrderRepository defines persistence for work orders.
type WorkOrderRepository interface {
	Create(ctx context.Context, wo *domain.WorkOrder) error
	GetByID(ctx context.Context, id int64) (*domain.WorkOrder, error)
	GetByWorkOrderNo(ctx context.Context, tenantID int64, woNo string) (*domain.WorkOrder, error)
	ListByOperator(ctx context.Context, operatorID int64) ([]domain.WorkOrder, error)
	ListByWarehouse(ctx context.Context, tenantID, warehouseID int64, status domain.WorkOrderStatus) ([]domain.WorkOrder, error)
	ListByOrder(ctx context.Context, orderID int64) ([]domain.WorkOrder, error)
	Update(ctx context.Context, wo *domain.WorkOrder) error
	Delete(ctx context.Context, id int64) error
}

// WorkOrderTemplateRepository defines persistence for work order templates.
type WorkOrderTemplateRepository interface {
	Create(ctx context.Context, t *domain.WorkOrderTemplate) error
	GetByID(ctx context.Context, id int64) (*domain.WorkOrderTemplate, error)
	GetByCode(ctx context.Context, tenantID int64, code string) (*domain.WorkOrderTemplate, error)
	List(ctx context.Context, tenantID int64) ([]domain.WorkOrderTemplate, error)
	ListByType(ctx context.Context, tenantID int64, wtype string) ([]domain.WorkOrderTemplate, error)
	Update(ctx context.Context, t *domain.WorkOrderTemplate) error
	Delete(ctx context.Context, id int64) error
}

// WorkflowProcessRepository defines persistence for workflow processes.
type WorkflowProcessRepository interface {
	Create(ctx context.Context, wp *domain.WorkflowProcess) error
	GetByID(ctx context.Context, id int64) (*domain.WorkflowProcess, error)
	ListByWorkOrder(ctx context.Context, workOrderID int64) ([]domain.WorkflowProcess, error)
	ListByOperator(ctx context.Context, operatorID int64) ([]domain.WorkflowProcess, error)
	Update(ctx context.Context, wp *domain.WorkflowProcess) error
}
