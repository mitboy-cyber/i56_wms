package domain
import "time"

// Operator represents a warehouse worker using a PDA device.
type Operator struct {
	ID          int64     `json:"id"`
	TenantID    int64     `json:"tenant_id"`
	WarehouseID int64     `json:"warehouse_id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Pin         string    `json:"pin"` // 4-digit login PIN
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

// Session tracks an active PDA login session.
type Session struct {
	ID         int64     `json:"id"`
	OperatorID int64     `json:"operator_id"`
	Token      string    `json:"token"`
	DeviceID   string    `json:"device_id"`
	IsActive   bool      `json:"is_active"`
	LoginAt    time.Time `json:"login_at"`
	LastSeen   time.Time `json:"last_seen"`
}

// ScanLog records every barcode scan in the warehouse.
type ScanLog struct {
	ID             int64     `json:"id"`
	TenantID       int64     `json:"tenant_id"`
	WarehouseID    int64     `json:"warehouse_id"`
	OperatorID     int64     `json:"operator_id"`
	Action         string    `json:"action"` // receive | pick | pack | load | query
	Barcode        string    `json:"barcode"`
	TrackingNumber string    `json:"tracking_number"`
	OrderNo        string    `json:"order_no,omitempty"`
	LocationCode   string    `json:"location_code,omitempty"`
	ContainerCode  string    `json:"container_code,omitempty"`
	Weight         float64   `json:"weight,omitempty"`
	Success        bool      `json:"success"`
	Message        string    `json:"message,omitempty"`
	ScannedAt      time.Time `json:"scanned_at"`
}

// PDAMenu item for the main screen
type PDAMenu struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	Icon  string `json:"icon"`
	Color string `json:"color"`
}

func DefaultMenus() []PDAMenu {
	return []PDAMenu{
		{Code: "receive", Name: "📦 包裹收货", Icon: "box-arrow-in-down", Color: "#0ea5e9"},
		{Code: "weigh", Name: "⚖️ 称重核重", Icon: "speedometer", Color: "#f59e0b"},
		{Code: "putaway", Name: "📍 上架入库", Icon: "pin-map", Color: "#10b981"},
		{Code: "pick", Name: "🛒 订单拣货", Icon: "basket", Color: "#8b5cf6"},
		{Code: "pack", Name: "📋 打包复核", Icon: "box-seam", Color: "#f97316"},
		{Code: "query", Name: "🔍 快件查询", Icon: "search", Color: "#64748b"},
	}
}
