package domain
import "time"

// ClientRoutePrice — per-client per-route pricing
type ClientRoutePrice struct {
	ID               int64     `json:"id"`
	TenantID         int64     `json:"tenant_id"`
	ClientID         int64     `json:"client_id"`
	RouteID          int64     `json:"route_id"`
	WeightPrice      float64   `json:"weight_price"`       // per kg
	VolumePrice      float64   `json:"volume_price"`       // per volumetric kg
	MinCharge        float64   `json:"min_charge"`         // minimum charge per order
	FirstWeightPrice float64   `json:"first_weight_price"` // first kg price
	FirstWeight      float64   `json:"first_weight"`       // first weight in kg
	AdditionalPrice  float64   `json:"additional_price"`   // per additional kg
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// ClientDeliveryFee — per-client delivery/pickup fees
type ClientDeliveryFee struct {
	ID             int64     `json:"id"`
	TenantID       int64     `json:"tenant_id"`
	ClientID       int64     `json:"client_id"`
	AreaGroupID    int64     `json:"area_group_id"`
	BaseFee        float64   `json:"base_fee"`
	PerKgFee       float64   `json:"per_kg_fee"`
	FreeThreshold  float64   `json:"free_threshold"` // free delivery above this amount
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
}

// ClientSurcharge — extra charges per client
type ClientSurcharge struct {
	ID          int64     `json:"id"`
	TenantID    int64     `json:"tenant_id"`
	ClientID    int64     `json:"client_id"`
	Name        string    `json:"name"`       // e.g. 报关费, 操作费
	ChargeType  string    `json:"charge_type"` // fixed | per_order | per_kg | percentage
	Amount      float64   `json:"amount"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

// ClientStoragePrice — per-client storage pricing
type ClientStoragePrice struct {
	ID             int64     `json:"id"`
	TenantID       int64     `json:"tenant_id"`
	ClientID       int64     `json:"client_id"`
	FreeDays       int       `json:"free_days"`        // free storage days
	DailyRate      float64   `json:"daily_rate"`       // per kg per day after free days
	MaxStorageDays int       `json:"max_storage_days"` // auto-return after this
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
}

// ClientServiceOverride — per-client service price override
type ClientServiceOverride struct {
	ID          int64     `json:"id"`
	TenantID    int64     `json:"tenant_id"`
	ClientID    int64     `json:"client_id"`
	ServiceCode string    `json:"service_code"` // WOODEN_CRATE, CONTENT_PHOTO, etc.
	UnitPrice   float64   `json:"unit_price"`   // overridden price
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

// MonthlyStatement — monthly billing statement
type MonthlyStatement struct {
	ID            int64     `json:"id"`
	TenantID      int64     `json:"tenant_id"`
	ClientID      int64     `json:"client_id"`
	Period        string    `json:"period"`       // "2026-07"
	OpeningBalance float64  `json:"opening_balance"`
	TotalCharges  float64   `json:"total_charges"`
	TotalPayments float64   `json:"total_payments"`
	ClosingBalance float64  `json:"closing_balance"`
	Status        string    `json:"status"`       // draft | sent | paid | overdue
	GeneratedAt   time.Time `json:"generated_at"`
}
