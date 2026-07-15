package domain

import "time"

// ProfitType represents the category of profit.
type ProfitType string

const (
	ProfitTypeShipping ProfitType = "shipping"
	ProfitTypeService  ProfitType = "service"
	ProfitTypeOther    ProfitType = "other"
)

// OrderProfit represents profit calculation for a single order.
type OrderProfit struct {
	ID               int64      `json:"id"`
	TenantID         int64      `json:"tenant_id"`
	OrderID          int64      `json:"order_id"`
	OrderNo          string     `json:"order_no"`
	ClientID         int64      `json:"client_id"`
	RouteID          int64      `json:"route_id"`
	Revenue          float64    `json:"revenue"`
	Cost             float64    `json:"cost"`
	GrossProfit      float64    `json:"gross_profit"`
	ProfitMargin     float64    `json:"profit_margin"`
	ShippingCost     float64    `json:"shipping_cost"`
	ServiceCost      float64    `json:"service_cost"`
	CreatedAt        time.Time  `json:"created_at"`
}

// ServiceProfit represents profit from a value-added service.
type ServiceProfit struct {
	ID           int64      `json:"id"`
	TenantID     int64      `json:"tenant_id"`
	ServiceOrderID int64    `json:"service_order_id"`
	OrderID      int64      `json:"order_id"`
	ClientID     int64      `json:"client_id"`
	ServiceType  string     `json:"service_type"`
	Revenue      float64    `json:"revenue"`
	Cost         float64    `json:"cost"`
	GrossProfit  float64    `json:"gross_profit"`
	CreatedAt    time.Time  `json:"created_at"`
}

// ClientProfit represents aggregate profit for a client over a period.
type ClientProfit struct {
	ID              int64     `json:"id"`
	TenantID        int64     `json:"tenant_id"`
	ClientID        int64     `json:"client_id"`
	Period          string    `json:"period"`
	TotalRevenue    float64   `json:"total_revenue"`
	TotalCost       float64   `json:"total_cost"`
	TotalProfit     float64   `json:"total_profit"`
	OrderCount      int       `json:"order_count"`
	ParcelCount     int       `json:"parcel_count"`
	AvgProfitPerOrder float64 `json:"avg_profit_per_order"`
	CreatedAt       time.Time `json:"created_at"`
}

// RouteProfit represents aggregate profit for a route over a period.
type RouteProfit struct {
	ID              int64     `json:"id"`
	TenantID        int64     `json:"tenant_id"`
	RouteID         int64     `json:"route_id"`
	Period          string    `json:"period"`
	TotalRevenue    float64   `json:"total_revenue"`
	TotalCost       float64   `json:"total_cost"`
	TotalProfit     float64   `json:"total_profit"`
	OrderCount      int       `json:"order_count"`
	TotalWeight     float64   `json:"total_weight"`
	AvgProfitPerKg  float64   `json:"avg_profit_per_kg"`
	CreatedAt       time.Time `json:"created_at"`
}
