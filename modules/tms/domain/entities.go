package domain

import "time"

// CargoType represents a type of cargo for classification and pricing.
type CargoType struct {
	ID        int64     `json:"id"`
	TenantID  int64     `json:"tenant_id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	IsDangerous bool    `json:"is_dangerous"`
	IsFragile   bool    `json:"is_fragile"`
	Remark      string  `json:"remark"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Carrier represents a logistics carrier company (SF, YTO, ZTO, etc.).
type Carrier struct {
	ID        int64     `json:"id"`
	TenantID  int64     `json:"tenant_id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	ContactPhone string `json:"contact_phone"`
	ContactEmail string `json:"contact_email"`
	Website   string    `json:"website"`
	AccountNo string    `json:"account_no"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Courier represents an individual courier/driver.
type Courier struct {
	ID          int64     `json:"id"`
	CarrierID   int64     `json:"carrier_id"`
	TenantID    int64     `json:"tenant_id"`
	Name        string    `json:"name"`
	Phone       string    `json:"phone"`
	VehiclePlate string   `json:"vehicle_plate"`
	IDNumber    string    `json:"id_number"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RouteType represents the type of shipping route.
type RouteType string

const (
	RouteTypeAir        RouteType = "air"
	RouteTypeSea        RouteType = "sea"
	RouteTypeSeaExpress RouteType = "sea_express"
	RouteTypeLand       RouteType = "land"
	RouteTypeRail       RouteType = "rail"
)

// Route represents a shipping route from origin to destination.
type Route struct {
	ID               int64     `json:"id"`
	TenantID         int64     `json:"tenant_id"`
	Name             string    `json:"name"`
	Code             string    `json:"code"`
	RouteType        RouteType `json:"route_type"`
	OriginCountry    string    `json:"origin_country"`
	OriginCity       string    `json:"origin_city"`
	DestCountry      string    `json:"dest_country"`
	DestCity         string    `json:"dest_city"`
	ShippingProviderID int64   `json:"shipping_provider_id"`
	CarrierID        int64     `json:"carrier_id"`
	CustomsBrokerID  int64     `json:"customs_broker_id"`
	CustomsPointID   int64     `json:"customs_point_id"`
	EstimatedDays    int       `json:"estimated_days"`
	BasePrice        float64   `json:"base_price"`
	PricePerKg       float64   `json:"price_per_kg"`
	MinWeight        float64   `json:"min_weight"`
	MaxWeight        float64   `json:"max_weight"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func DefaultCargoTypes() []CargoType {
	return []CargoType{
		{ID: 1, Name: "普通货物", Code: "general", IsActive: true},
		{ID: 2, Name: "易碎品", Code: "fragile", IsFragile: true, IsActive: true},
		{ID: 3, Name: "液体", Code: "liquid", IsActive: true},
		{ID: 4, Name: "电子产品", Code: "electronics", IsActive: true},
		{ID: 5, Name: "危险品", Code: "dangerous", IsDangerous: true, IsActive: true},
	}
}
