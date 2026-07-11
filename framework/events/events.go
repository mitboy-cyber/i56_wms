package events

import (
	"encoding/json"
	"time"
)

// ─── Domain Event Types ────────────────────────────────────────────────────

// OrderCreated is published when a new order is created.
type OrderCreated struct {
	OrderID    int64   `json:"order_id"`
	OrderNo    string  `json:"order_no"`
	ClientID   int64   `json:"client_id"`
	TotalPrice float64 `json:"total_price"`
	Timestamp  string  `json:"timestamp"`
}

func (e OrderCreated) EventName() string     { return "order.created" }
func (e OrderCreated) OccurredAt() time.Time { t, _ := time.Parse(time.RFC3339, e.Timestamp); return t }

// ParcelReceived is published when a parcel is received at a warehouse.
type ParcelReceived struct {
	ParcelID     int64  `json:"parcel_id"`
	TrackingNo   string `json:"tracking_no"`
	WarehouseID  int64  `json:"warehouse_id"`
	ProductName  string `json:"product_name"`
	Timestamp    string `json:"timestamp"`
}

func (e ParcelReceived) EventName() string     { return "parcel.received" }
func (e ParcelReceived) OccurredAt() time.Time { t, _ := time.Parse(time.RFC3339, e.Timestamp); return t }

// ParcelCreated is published when a new parcel is pre-declared.
type ParcelCreated struct {
	ParcelID     int64  `json:"parcel_id"`
	TrackingNo   string `json:"tracking_no"`
	WarehouseID  int64  `json:"warehouse_id"`
	ClientID     int64  `json:"client_id"`
	ProductName  string `json:"product_name"`
	Timestamp    string `json:"timestamp"`
}

func (e ParcelCreated) EventName() string     { return "parcel.created" }
func (e ParcelCreated) OccurredAt() time.Time { t, _ := time.Parse(time.RFC3339, e.Timestamp); return t }

// StatusChanged is published when any entity's status changes.
type StatusChanged struct {
	EntityType string `json:"entity_type"` // "order" or "parcel"
	EntityID   int64  `json:"entity_id"`
	OldStatus  string `json:"old_status"`
	NewStatus  string `json:"new_status"`
	Timestamp  string `json:"timestamp"`
}

func (e StatusChanged) EventName() string     { return "status.changed" }
func (e StatusChanged) OccurredAt() time.Time { t, _ := time.Parse(time.RFC3339, e.Timestamp); return t }

// PaymentReceived is published when a client payment is recorded.
type PaymentReceived struct {
	ClientID  int64   `json:"client_id"`
	Amount    float64 `json:"amount"`
	TxID      string  `json:"tx_id"`
	Timestamp string  `json:"timestamp"`
}

func (e PaymentReceived) EventName() string     { return "payment.received" }
func (e PaymentReceived) OccurredAt() time.Time { t, _ := time.Parse(time.RFC3339, e.Timestamp); return t }

// ContainerClosed is published when a container is closed/sealed.
type ContainerClosed struct {
	ContainerID int64 `json:"container_id"`
	ParcelCount int   `json:"parcel_count"`
	Timestamp   string `json:"timestamp"`
}

func (e ContainerClosed) EventName() string     { return "container.closed" }
func (e ContainerClosed) OccurredAt() time.Time { t, _ := time.Parse(time.RFC3339, e.Timestamp); return t }

// ─── JSON helpers ──────────────────────────────────────────────────────────

// ToJSON marshals an event to its JSON representation.
func ToJSON(e interface{ EventName() string; OccurredAt() time.Time }) ([]byte, error) {
	return json.Marshal(e)
}
