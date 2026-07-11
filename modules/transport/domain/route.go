package domain

import "time"

type Route struct {
	ID              int64     `json:"id"`
	TenantID        int64     `json:"tenant_id"`
	WarehouseID     int64     `json:"warehouse_id"`
	Name            string    `json:"name"`
	TransportType   string    `json:"transport_type"`
	AreaGroupID     int64     `json:"area_group_id"`
	CargoTypes      []string  `json:"cargo_types"`
	MinWeight       float64   `json:"min_weight"`
	MaxWeight       float64   `json:"max_weight"`
	VolumeCoeff     int       `json:"volume_coeff"`
	WeightRounding  float64   `json:"weight_rounding"`
	MinAmount       float64   `json:"min_amount"`
	MinDays         int       `json:"min_days"`
	MaxDays         int       `json:"max_days"`
	BaseWeightPrice float64   `json:"base_weight_price"`
	BaseVolumePrice float64   `json:"base_volume_price"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Courier struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	CountryRegion string `json:"country_region"`
}

type CargoType struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type TransportType struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}
