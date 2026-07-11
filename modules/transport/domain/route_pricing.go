package domain

import "time"

// RouteTemplate represents a complete route configuration with pricing
type RouteTemplate struct {
	ID              int64     `json:"id"`
	TenantID        int64     `json:"tenant_id"`
	WarehouseID     int64     `json:"warehouse_id"`
	WarehouseName   string    `json:"warehouse_name"`
	Name            string    `json:"name"`            // 线路名
	Code            string    `json:"code"`            // 线路编码 e.g., "WH001-AIR"
	TransportType   string    `json:"transport_type"`  // 空运/海运/海快/空运特货/商业海快
	CargoType       string    `json:"cargo_type"`      // 普货/家具类/一类~六类/易碎品
	TaxType         string    `json:"tax_type"`        // 不包税/频税/全包税
	MaxLength       float64   `json:"max_length"`      // 单件长度限制(cm)
	IsActive        bool      `json:"is_active"`
	// Pricing
	WeightPrice       float64 `json:"weight_price"`       // 重量单价(元/kg)
	VolumePrice       float64 `json:"volume_price"`       // 体积单价(元/才)
	MinCharge         float64 `json:"min_charge"`         // 最低收费
	FirstWeight       float64 `json:"first_weight"`       // 首重(kg)
	FirstWeightPrice  float64 `json:"first_weight_price"` // 首重价格
	ContWeightPrice   float64 `json:"cont_weight_price"`  // 续重单价
	FirstVolume       float64 `json:"first_volume"`       // 首体积(才)
	FirstVolumePrice  float64 `json:"first_volume_price"` // 首体积价格
	ContVolumePrice   float64 `json:"cont_volume_price"`  // 续体积单价
	// Schedule
	ScheduleType    string    `json:"schedule_type"`   // 每周/每月
	ShipDays        string    `json:"ship_days"`       // 周一,周三,周五 或 1,15
	CutoffTime      string    `json:"cutoff_time"`     // 截止时间
	EstimatedDays   int       `json:"estimated_days"`  // 预计天数
	DeliveryMethod  string    `json:"delivery_method"` // 派送方式
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// CarrierDeliveryFee represents carrier delivery pricing
type CarrierDeliveryFee struct {
	ID             int64   `json:"id"`
	TenantID       int64   `json:"tenant_id"`
	CarrierID      int64   `json:"carrier_id"`
	CarrierName    string  `json:"carrier_name"`
	CustomsPoint   string  `json:"customs_point"`   // 清关点: 台北/台中/高雄
	Area           string  `json:"area"`            // 区域: 预设/台湾北部/台湾中部/台湾南部/台湾东部
	DeliveryMethod string  `json:"delivery_method"` // 宅配/专车/自取
	Condition      string  `json:"condition"`       // 条件: 重量>39.8/单边长≥600
	PriceMode      string  `json:"price_mode"`      // 固定费用/按重量/按体积
	Price          float64 `json:"price"`           // 价格
	FreeThreshold  string  `json:"free_threshold"`  // 免运门槛: ≥10kg
}
