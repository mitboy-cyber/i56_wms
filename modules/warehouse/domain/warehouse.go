package domain

import "time"

type Warehouse struct {
	ID        int64     `json:"id"`
	TenantID  int64     `json:"tenant_id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Address   string    `json:"address"`
	Contact   string    `json:"contact"`
	Phone     string    `json:"phone"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Zone struct {
	ID          int64     `json:"id"`
	WarehouseID int64     `json:"warehouse_id"`
	Name        string    `json:"name"`
	ZoneType    string    `json:"zone_type"`
	Code        string    `json:"code"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Location struct {
	ID          int64     `json:"id"`
	ZoneID      int64     `json:"zone_id"`
	Code        string    `json:"code"`
	LocationType string   `json:"location_type"`
	IsOccupied  bool      `json:"is_occupied"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
