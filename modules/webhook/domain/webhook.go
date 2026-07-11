package domain
import "time"

type WebhookSubscription struct {
	ID        int64     `json:"id"`
	TenantID  int64     `json:"tenant_id"`
	ClientID  int64     `json:"client_id"`
	Event     string    `json:"event"`
	URL       string    `json:"url"`
	Secret    string    `json:"secret"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type WebhookDeliveryLog struct {
	ID             int64     `json:"id"`
	SubscriptionID int64     `json:"subscription_id"`
	Event          string    `json:"event"`
	Payload        string    `json:"payload"`
	StatusCode     int       `json:"status_code"`
	Error          string    `json:"error"`
	RetryCount     int       `json:"retry_count"`
	DeliveredAt    time.Time `json:"delivered_at"`
}

var SupportedEvents = []string{
	"parcel.arrived", "parcel.weighed", "parcel.stored", "parcel.shipped",
	"order.created", "order.shipped", "order.delivered", "order.cancelled",
}
