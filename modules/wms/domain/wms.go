package domain
import "time"

type Zone struct {
	ID          int64     `json:"id"`
	WarehouseID int64     `json:"warehouse_id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	ZoneTypeID  int64     `json:"zone_type_id"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

type ZoneType struct {
	ID   int64  `json:"id"`
	Name string `json:"name"` // 收货区/存储区/拣货区/发货区/退货区
	Code string `json:"code"`
}

type Location struct {
	ID           int64     `json:"id"`
	ZoneID       int64     `json:"zone_id"`
	Code         string    `json:"code"`           // A-01-02
	LocationTypeID int64   `json:"location_type_id"`
	Barcode      string    `json:"barcode"`        // QR code for PDA scanning
	MaxWeight    float64   `json:"max_weight_kg"`
	IsOccupied   bool      `json:"is_occupied"`
	CurrentParcelID *int64 `json:"current_parcel_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type LocationType struct {
	ID   int64  `json:"id"`
	Name string `json:"name"` // 货架位/托盘位/地面位/流水线
	Code string `json:"code"`
}

type Container struct {
	ID          int64     `json:"id"`
	WarehouseID int64     `json:"warehouse_id"`
	Name        string    `json:"name"`
	ContainerNo string    `json:"container_no"`    // 柜号 TCLU1234567
	SealNo      string    `json:"seal_no"`         // 封条号
	Type        string    `json:"type"`            // 20GP/40GP/40HQ/45HQ
	MaxWeight   float64   `json:"max_weight_kg"`
	Status      string    `json:"status"`          // empty/loading/loaded/sealed/shipped
	CreatedAt   time.Time `json:"created_at"`
}

type InboundMachine struct {
	ID          int64     `json:"id"`
	WarehouseID int64     `json:"warehouse_id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

func DefaultZoneTypes() []ZoneType {
	return []ZoneType{
		{ID:1,Name:"收货区",Code:"RECEIVING"},
		{ID:2,Name:"存储区",Code:"STORAGE"},
		{ID:3,Name:"拣货区",Code:"PICKING"},
		{ID:4,Name:"发货区",Code:"SHIPPING"},
		{ID:5,Name:"退货区",Code:"RETURNS"},
	}
}

func DefaultLocationTypes() []LocationType {
	return []LocationType{
		{ID:1,Name:"货架位",Code:"SHELF"},
		{ID:2,Name:"托盘位",Code:"PALLET"},
		{ID:3,Name:"地面位",Code:"FLOOR"},
	}
}
