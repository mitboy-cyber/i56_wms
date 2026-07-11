package domain
import "time"

type AreaGroup struct {
	ID      int64  `json:"id"`
	TenantID int64 `json:"tenant_id"`
	Name    string `json:"name"`
	Areas   []string `json:"areas"` // ["台北市","新北市","基隆市"]
}

type CarrierNumber struct {
	ID        int64     `json:"id"`
	CarrierID int64     `json:"carrier_id"`
	Prefix    string    `json:"prefix"`
	StartNo   int64     `json:"start_no"`
	EndNo     int64     `json:"end_no"`
	CurrentNo int64     `json:"current_no"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type CustomsBroker struct {
	ID        int64  `json:"id"`
	TenantID  int64  `json:"tenant_id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	Contact   string `json:"contact"`
	Phone     string `json:"phone"`
	IsActive  bool   `json:"is_active"`
}

type CustomsPoint struct {
	ID       int64  `json:"id"`
	TenantID int64  `json:"tenant_id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	IsActive bool   `json:"is_active"`
}

type CustomsNumber struct {
	ID            int64     `json:"id"`
	CustomsPointID int64    `json:"customs_point_id"`
	Prefix        string    `json:"prefix"`
	StartNo       int64     `json:"start_no"`
	EndNo         int64     `json:"end_no"`
	CurrentNo     int64     `json:"current_no"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
}

type ContainerLoading struct {
	ID          int64     `json:"id"`
	ContainerID int64     `json:"container_id"`
	OrderID     int64     `json:"order_id"`
	ParcelIDs   []int64   `json:"parcel_ids"`
	LoadedBy   string    `json:"loaded_by"`
	LoadedAt   time.Time `json:"loaded_at"`
}

type ShippingProvider struct {
	ID       int64  `json:"id"`
	TenantID int64  `json:"tenant_id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	IsActive bool   `json:"is_active"`
}

type TransportType struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`  // 空运/海运/海快/陆运/铁路
	Code string `json:"code"`
}

type Tracking struct {
	ID              int64     `json:"id"`
	OrderID         int64     `json:"order_id"`
	TrackingNumber  string    `json:"tracking_number"`
	CarrierCode     string    `json:"carrier_code"`
	Status          string    `json:"status"`
	StatusDetail    string    `json:"status_detail"`
	Location        string    `json:"location"`
	EventTime       time.Time `json:"event_time"`
	CreatedAt       time.Time `json:"created_at"`
}

func DefaultTransportTypes() []TransportType {
	return []TransportType{
		{ID:1,Name:"空运",Code:"air"},{ID:2,Name:"海运",Code:"sea"},
		{ID:3,Name:"海快",Code:"sea_express"},{ID:4,Name:"陆运",Code:"land"},{ID:5,Name:"铁路",Code:"rail"},
	}
}
