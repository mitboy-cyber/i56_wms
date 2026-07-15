package domain

import "time"

// NotificationType represents the category of notification.
type NotificationType string

const (
	NotifTypeSystem    NotificationType = "system"
	NotifTypeOrder     NotificationType = "order"
	NotifTypeParcel    NotificationType = "parcel"
	NotifTypeFinance   NotificationType = "finance"
	NotifTypeAlarm     NotificationType = "alarm"
)

// NotificationChannel represents how the notification is delivered.
type NotificationChannel string

const (
	ChannelWeb      NotificationChannel = "web"
	ChannelSMS      NotificationChannel = "sms"
	ChannelEmail    NotificationChannel = "email"
	ChannelWechat   NotificationChannel = "wechat"
	ChannelPush     NotificationChannel = "push"
)

// NotificationStatus represents delivery status.
type NotificationStatus string

const (
	NotifStatusPending   NotificationStatus = "pending"
	NotifStatusSent      NotificationStatus = "sent"
	NotifStatusRead      NotificationStatus = "read"
	NotifStatusFailed    NotificationStatus = "failed"
)

// Notification represents a notification message sent to users.
type Notification struct {
	ID              int64               `json:"id"`
	TenantID        int64               `json:"tenant_id"`
	TemplateID      int64               `json:"template_id"`
	RecipientID     int64               `json:"recipient_id"`
	RecipientType   string              `json:"recipient_type"`
	Type            NotificationType    `json:"type"`
	Channel         NotificationChannel `json:"channel"`
	Title           string              `json:"title"`
	Content         string              `json:"content"`
	Status          NotificationStatus  `json:"status"`
	ReferenceType   string              `json:"reference_type"`
	ReferenceID     int64               `json:"reference_id"`
	ReadAt          *time.Time          `json:"read_at"`
	SentAt          *time.Time          `json:"sent_at"`
	ErrorMsg        string              `json:"error_msg"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
}
