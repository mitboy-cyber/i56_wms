package repository

import (
	"context"
	"github.com/i56/modules/notification/domain"
)

// NotificationRepository defines persistence for notifications.
type NotificationRepository interface {
	Create(ctx context.Context, n *domain.Notification) error
	GetByID(ctx context.Context, id int64) (*domain.Notification, error)
	ListByRecipient(ctx context.Context, tenantID, recipientID int64, recipientType string) ([]domain.Notification, error)
	ListUnreadByRecipient(ctx context.Context, tenantID, recipientID int64, recipientType string) ([]domain.Notification, error)
	MarkAsRead(ctx context.Context, id int64) error
	MarkAsSent(ctx context.Context, id int64) error
	Update(ctx context.Context, n *domain.Notification) error
}

// NotificationTemplateRepository defines persistence for notification templates.
type NotificationTemplateRepository interface {
	Create(ctx context.Context, t *domain.NotificationTemplate) error
	GetByID(ctx context.Context, id int64) (*domain.NotificationTemplate, error)
	GetByCode(ctx context.Context, tenantID int64, code string) (*domain.NotificationTemplate, error)
	List(ctx context.Context, tenantID int64) ([]domain.NotificationTemplate, error)
	ListByType(ctx context.Context, tenantID int64, ntype domain.NotificationType) ([]domain.NotificationTemplate, error)
	Update(ctx context.Context, t *domain.NotificationTemplate) error
	Delete(ctx context.Context, id int64) error
}

// PrintTemplateRepository defines persistence for print templates.
type PrintTemplateRepository interface {
	Create(ctx context.Context, t *domain.PrintTemplate) error
	GetByID(ctx context.Context, id int64) (*domain.PrintTemplate, error)
	GetByCode(ctx context.Context, tenantID int64, code string) (*domain.PrintTemplate, error)
	List(ctx context.Context, tenantID int64) ([]domain.PrintTemplate, error)
	ListByType(ctx context.Context, tenantID int64, ptype domain.PrintTemplateType) ([]domain.PrintTemplate, error)
	SetDefault(ctx context.Context, tenantID, id int64, ptype domain.PrintTemplateType) error
	Update(ctx context.Context, t *domain.PrintTemplate) error
	Delete(ctx context.Context, id int64) error
}
