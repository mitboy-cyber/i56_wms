// Package types defines shared API request/response types.
package types

// ─── Order ────────────────────────────────────────────────────────────

// CreateOrderRequest is the validated request for creating an order.
type CreateOrderRequest struct {
	RecipientName    string  `json:"recipient_name"    validate:"required,min=1,max=128"`
	RouteID          int64   `json:"route_id"          validate:"required,gt=0"`
	ParcelCount      int     `json:"parcel_count"      validate:"required,gte=1,lte=9999"`
	TotalPrice       float64 `json:"total_price"       validate:"omitempty,gte=0"`
	TrackingNumbers  string  `json:"tracking_numbers"  validate:"omitempty,max=2048"`
	Remark           string  `json:"remark"            validate:"omitempty,max=512"`
}

// ─── Parcel ───────────────────────────────────────────────────────────

type CreateParcelRequest struct {
	TrackingNumber    string  `json:"tracking_number"    validate:"required,min=1,max=64"`
	ProductName       string  `json:"product_name"       validate:"required,min=1,max=256"`
	ActualWeight      float64 `json:"actual_weight"      validate:"omitempty,gte=0"`
	ChargeableWeight  float64 `json:"chargeable_weight"  validate:"omitempty,gte=0"`
	DeclaredValue     float64 `json:"declared_value"     validate:"omitempty,gte=0"`
	CourierCode       string  `json:"courier_code"       validate:"omitempty,max=32"`
	CargoType         string  `json:"cargo_type"         validate:"omitempty,max=32"`
	WarehouseID       int64   `json:"warehouse_id"       validate:"omitempty,gt=0"`
	Remark            string  `json:"remark"             validate:"omitempty,max=256"`
}

// ─── Warehouse ────────────────────────────────────────────────────────

type CreateWarehouseRequest struct {
	Name    string `json:"name"    validate:"required,min=1,max=128"`
	Code    string `json:"code"    validate:"required,min=1,max=32"`
	Address string `json:"address" validate:"required,min=1,max=512"`
	Contact string `json:"contact" validate:"omitempty,max=64"`
	Phone   string `json:"phone"   validate:"omitempty,max=32"`
}

// ─── Client ───────────────────────────────────────────────────────────

type CreateClientRequest struct {
	Name    string `json:"name"    validate:"required,min=1,max=128"`
	Code    string `json:"code"    validate:"required,min=1,max=64"`
	Type    string `json:"type"    validate:"omitempty,max=32"`
	Contact string `json:"contact" validate:"omitempty,max=64"`
	Phone   string `json:"phone"   validate:"omitempty,max=32"`
}
