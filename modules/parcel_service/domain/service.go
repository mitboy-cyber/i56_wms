package domain
import "time"

type ServiceType struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Code      string  `json:"code"`
	Category  string  `json:"category"`
	UnitPrice float64 `json:"unit_price"`
	PriceMode string  `json:"price_mode"` // fixed | per_qty | per_kg
	Priority  int     `json:"priority"`
}

type ServiceOrder struct {
	ID          int64     `json:"id"`
	TenantID    int64     `json:"tenant_id"`
	ClientID    int64     `json:"client_id"`
	OrderID     *int64    `json:"order_id"`
	ParcelID    *int64    `json:"parcel_id"`
	ServiceType string    `json:"service_type"`
	Quantity    int       `json:"quantity"`
	TotalPrice  float64   `json:"total_price"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func DefaultServiceTypes() []ServiceType {
	return []ServiceType{
		{Name:"退货寄回",Code:"RETURN_GOODS",Category:"退货类",UnitPrice:2.00,PriceMode:"fixed"},
		{Name:"易碎品贴纸",Code:"FRAGILE_STICKER",Category:"打包类",UnitPrice:0.10,PriceMode:"fixed"},
		{Name:"外箱标识拍照",Code:"OUTERBOX_PHOTO",Category:"加固类",UnitPrice:0.10,PriceMode:"fixed"},
		{Name:"打木箱",Code:"WOODEN_CRATE",Category:"加固类",UnitPrice:80.00,PriceMode:"fixed"},
		{Name:"打木架",Code:"WOODEN_FRAME",Category:"加固类",UnitPrice:50.00,PriceMode:"fixed"},
		{Name:"包装气柱袋",Code:"WRAP_AIRBAG",Category:"加固类",UnitPrice:5.00,PriceMode:"fixed"},
		{Name:"包装气泡棉",Code:"WRAP_BUBBLE",Category:"加固类",UnitPrice:2.00,PriceMode:"fixed"},
		{Name:"包裹拆分",Code:"SPLIT_PARCEL",Category:"开箱类",UnitPrice:2.00,PriceMode:"per_qty"},
		{Name:"填充气柱袋",Code:"FILL_AIRBAG",Category:"开箱类",UnitPrice:5.00,PriceMode:"fixed"},
		{Name:"内容物拍照",Code:"CONTENT_PHOTO",Category:"开箱类",UnitPrice:1.00,PriceMode:"fixed"},
		{Name:"清点数量",Code:"COUNT_QTY",Category:"开箱类",UnitPrice:0.10,PriceMode:"per_qty"},
		{Name:"确认型号",Code:"CONFIRM_MODEL",Category:"开箱类",UnitPrice:0.10,PriceMode:"fixed"},
		{Name:"开箱验货",Code:"OPEN_INSPECT",Category:"开箱类",UnitPrice:0.00,PriceMode:"fixed"},
	}
}
