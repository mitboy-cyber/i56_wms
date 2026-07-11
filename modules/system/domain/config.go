package domain
import "time"

type LogisticsAPIConfig struct {
	ID int64;TenantID int64;CarrierID int64;Name string;BaseURL string
	AppKey string;AppSecret string;Token string;WebhookURL string
	TimeoutMs int;MaxRetry int;IsActive bool
	CreatedAt time.Time;UpdatedAt time.Time
}
type CustomsBrokerAPIConfig struct {
	ID int64;TenantID int64;BrokerID int64;Name string;Code string
	Prefix string;APIEndpoint string;APIKey string;APISecret string
	SupportedDocs string;IsActive bool
}
type PrinterConfig struct {
	ID int64;TenantID int64;Name string;PrinterType string
	PaperWidth int;PaperHeight int;DPI int;ConnectionType string
	IPAddress string;Port int;IsDefault bool;IsActive bool
}
type SystemSetting struct {
	ID int64;TenantID int64;Key string;Value string;Type string
	Group string;Label string;UpdatedAt time.Time
}
