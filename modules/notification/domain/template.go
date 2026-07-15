package domain

import "time"

// NotificationTemplate represents a reusable template for notifications.
type NotificationTemplate struct {
	ID              int64               `json:"id"`
	TenantID        int64               `json:"tenant_id"`
	Name            string              `json:"name"`
	Code            string              `json:"code"`
	Type            NotificationType    `json:"type"`
	Channel         NotificationChannel `json:"channel"`
	TitleTemplate   string              `json:"title_template"`
	ContentTemplate string              `json:"content_template"`
	Variables       []string            `json:"variables"`
	IsActive        bool                `json:"is_active"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
}

func DefaultNotificationTemplates() []NotificationTemplate {
	return []NotificationTemplate{
		{
			Name: "包裹入库通知", Code: "parcel_received", Type: NotifTypeParcel, Channel: ChannelWeb,
			TitleTemplate: "包裹已入库 - {{.TrackingNumber}}",
			ContentTemplate: "您的包裹 {{.TrackingNumber}}（{{.ProductName}}）已成功入库，重量 {{.Weight}}kg，存放位置 {{.Location}}。",
			Variables: []string{"TrackingNumber", "ProductName", "Weight", "Location"},
		},
		{
			Name: "订单发货通知", Code: "order_shipped", Type: NotifTypeOrder, Channel: ChannelWeb,
			TitleTemplate: "订单已发货 - {{.OrderNo}}",
			ContentTemplate: "您的订单 {{.OrderNo}} 已于 {{.ShippedAt}} 发货，承运单号 {{.CarrierTrackingNo}}，预计 {{.EstimatedDays}} 天到达。",
			Variables: []string{"OrderNo", "ShippedAt", "CarrierTrackingNo", "EstimatedDays"},
		},
		{
			Name: "余额不足提醒", Code: "low_balance", Type: NotifTypeFinance, Channel: ChannelWeb,
			TitleTemplate: "账户余额不足提醒",
			ContentTemplate: "您的账户当前余额为 ¥{{.Balance}}，低于预警值，请及时充值以确保服务不中断。",
			Variables: []string{"Balance"},
		},
		{
			Name: "异常包裹告警", Code: "abnormal_parcel", Type: NotifTypeAlarm, Channel: ChannelSMS,
			TitleTemplate: "异常包裹告警 - {{.TrackingNumber}}",
			ContentTemplate: "包裹 {{.TrackingNumber}} 出现异常：{{.Reason}}，请及时处理。",
			Variables: []string{"TrackingNumber", "Reason"},
		},
		{
			Name: "包裹签收通知", Code: "parcel_delivered", Type: NotifTypeParcel, Channel: ChannelWeb,
			TitleTemplate: "包裹已签收 - {{.TrackingNumber}}",
			ContentTemplate: "您的包裹 {{.TrackingNumber}} 已于 {{.DeliveredAt}} 签收。感谢您的使用！",
			Variables: []string{"TrackingNumber", "DeliveredAt"},
		},
	}
}
