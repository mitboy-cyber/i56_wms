package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/modules/notification/domain"
)

// MemNotificationRepo is an in-memory implementation for notifications.
type MemNotificationRepo struct {
	mu            sync.RWMutex
	notifications map[int64]*domain.Notification
	nextID        int64
}

func NewMemNotificationRepo() *MemNotificationRepo {
	r := &MemNotificationRepo{notifications: make(map[int64]*domain.Notification), nextID: 1}
	r.seed()
	return r
}

func (r *MemNotificationRepo) seed() {
	now := time.Now()
	seeds := []struct {
		templateID    int64
		recipientID   int64
		recipientType string
		ntype         domain.NotificationType
		channel       domain.NotificationChannel
		title         string
		content       string
		status        domain.NotificationStatus
		refType       string
		refID         int64
	}{
		{1, 1, "client_user", domain.NotifTypeParcel, domain.ChannelWeb, "包裹已入库 - TN20240101001", "您的包裹 TN20240101001（手机壳-黑色）已成功入库，重量 0.15kg，存放位置 A-01-01。", domain.NotifStatusRead, "parcel", 1},
		{2, 1, "client_user", domain.NotifTypeOrder, domain.ChannelWeb, "订单已发货 - ORD202401001", "您的订单 ORD202401001 已发货，承运单号 SF1234567890，预计 3 天到达。", domain.NotifStatusSent, "order", 1},
		{3, 2, "client_user", domain.NotifTypeFinance, domain.ChannelWeb, "账户余额不足提醒", "您的账户当前余额为 ¥500.00，低于预警值，请及时充值。", domain.NotifStatusPending, "ledger", 1},
		{4, 1, "client_user", domain.NotifTypeAlarm, domain.ChannelSMS, "异常包裹告警 - TN20240101003", "包裹 TN20240101003 出现异常：陶瓷杯-白色发现破损，请及时处理。", domain.NotifStatusSent, "parcel", 3},
		{5, 1, "operator", domain.NotifTypeSystem, domain.ChannelPush, "新任务分配", "您有一个新的拣货任务 WO202401002，请及时处理。", domain.NotifStatusPending, "work_order", 2},
	}
	for _, s := range seeds {
		id := atomic.AddInt64(&r.nextID, 1) - 1
		n := &domain.Notification{
			ID:            id,
			TenantID:      1,
			TemplateID:    s.templateID,
			RecipientID:   s.recipientID,
			RecipientType: s.recipientType,
			Type:          s.ntype,
			Channel:       s.channel,
			Title:         s.title,
			Content:       s.content,
			Status:        s.status,
			ReferenceType: s.refType,
			ReferenceID:   s.refID,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		if s.status == domain.NotifStatusSent {
			n.SentAt = &now
		}
		if s.status == domain.NotifStatusRead {
			n.SentAt = &now
			n.ReadAt = &now
		}
		r.notifications[id] = n
	}
}

func (r *MemNotificationRepo) Create(ctx context.Context, n *domain.Notification) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	n.ID = atomic.AddInt64(&r.nextID, 1) - 1
	now := time.Now()
	n.CreatedAt = now
	n.UpdatedAt = now
	r.notifications[n.ID] = n
	return nil
}

func (r *MemNotificationRepo) GetByID(ctx context.Context, id int64) (*domain.Notification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	n, ok := r.notifications[id]
	if !ok {
		return nil, nil
	}
	return n, nil
}

func (r *MemNotificationRepo) ListByRecipient(ctx context.Context, tenantID, recipientID int64, recipientType string) ([]domain.Notification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.Notification
	for _, n := range r.notifications {
		if n.TenantID == tenantID && n.RecipientID == recipientID && n.RecipientType == recipientType {
			result = append(result, *n)
		}
	}
	return result, nil
}

func (r *MemNotificationRepo) ListUnreadByRecipient(ctx context.Context, tenantID, recipientID int64, recipientType string) ([]domain.Notification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.Notification
	for _, n := range r.notifications {
		if n.TenantID == tenantID && n.RecipientID == recipientID && n.RecipientType == recipientType && n.Status != domain.NotifStatusRead {
			result = append(result, *n)
		}
	}
	return result, nil
}

func (r *MemNotificationRepo) MarkAsRead(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	n, ok := r.notifications[id]
	if !ok {
		return nil
	}
	now := time.Now()
	n.Status = domain.NotifStatusRead
	n.ReadAt = &now
	n.UpdatedAt = now
	return nil
}

func (r *MemNotificationRepo) MarkAsSent(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	n, ok := r.notifications[id]
	if !ok {
		return nil
	}
	now := time.Now()
	n.Status = domain.NotifStatusSent
	n.SentAt = &now
	n.UpdatedAt = now
	return nil
}

func (r *MemNotificationRepo) Update(ctx context.Context, n *domain.Notification) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.notifications[n.ID]; !ok {
		return nil
	}
	n.UpdatedAt = time.Now()
	r.notifications[n.ID] = n
	return nil
}

var _ NotificationRepository = (*MemNotificationRepo)(nil)
