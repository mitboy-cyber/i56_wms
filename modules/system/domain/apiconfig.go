package domain

import "time"

// ===================================================================
// CourierAPIConfig — 快递公司 API 集成配置
// Used for: 顺丰、圆通、中通、EMS 等快递公司 API 对接
// ===================================================================
type CourierAPIConfig struct {
	ID              int64
	TenantID        int64
	CourierID       int64
	Name            string // e.g. "顺丰速运API", "圆通快递API"
	APIEndpoint     string // e.g. "https://open.sf-express.com/std/service"
	APIKey          string
	APISecret       string
	TrackingPattern string    // e.g. "^SF\\d{12}$" 运单号正则
	AuthType        string    // api_key, hmac, oauth2
	ExtraHeaders    string    // JSON string: {"X-Custom":"value"}
	IsActive        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ===================================================================
// CustomsBrokerConfig — 报关行 API 集成配置
// Used for: 电子报关、海关申报系统对接
// ===================================================================
type CustomsBrokerConfig struct {
	ID                 int64
	TenantID           int64
	BrokerID           int64
	Name               string // e.g. "厦门清关代理"
	DeclarationAPIURL  string // e.g. "https://customs.xm-port.gov.cn/api/v2"
	APIKey             string
	APISecret          string
	CustomsPointID     string // 海关口岸编号
	NumberPrefix       string // 报关单号前缀 e.g. "776XM"
	SupportedDocuments string // JSON array of supported doc types
	IsActive           bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// ===================================================================
// NotificationChannel — 通知渠道配置 (富字段版)
// Used for: 邮件、短信、Line、Telegram、Webhook 等通知发送
// ===================================================================
type NotificationChannel struct {
	ID          int64
	TenantID    int64
	Name        string    // e.g. "系统邮件", "阿里云短信"
	ChannelType string    // email, sms, line, telegram, webhook
	Provider    string    // smtp, sendgrid, aliyun_sms, twilio, line_api, telegram_bot
	ConfigJSON  string    // provider-specific config as JSON
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ===================================================================
// PrintTemplate — 打印模板配置
// Used for: 标签、发票、装箱单、运单等打印模板
// ===================================================================
type PrintTemplate struct {
	ID              int64
	TenantID        int64
	Name            string // e.g. "顺丰面单模板", "标准发票模板"
	Type            string // label, invoice, packing_list, waybill
	PaperSize       string // e.g. "100x150mm", "A4", "4x6inch"
	TemplateContent string // ZPL/PDF/HTML template body
	PrinterType     string // thermal, laser, inkjet
	IsActive        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ===================================================================
// StorageConfig — 对象存储配置
// Used for: MinIO, S3, OSS, COS 等对象存储对接
// ===================================================================
type StorageConfig struct {
	ID        int64
	TenantID  int64
	Name      string // e.g. "厦门仓MinIO", "阿里云OSS"
	Provider  string // minio, s3, oss, cos
	Bucket    string
	Endpoint  string
	AccessKey string
	SecretKey string
	Region    string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
