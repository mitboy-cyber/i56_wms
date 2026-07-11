package domain

import "time"

// ===================================================================
// BFT56 5-Tab Customer Pricing System — multi-dimensional domain models
//
// Tab 1: 客户×线路价 — client × route × cargo_type × tax_mode
// Tab 2: 客户×仓储价 — client × warehouse
// Tab 3: 客户×派送费 — client × carrier × customs_point × area × delivery_method
// Tab 4: 客户×加收费 — client × carrier × charge_type × tier × customs_point × area
// Tab 5: 客户×附加服务 — client × service_type
// ===================================================================

// ─── Tab 1: Route Pricing ───────────────────────────────────────────
// 客户×线路价: client × route × cargo_type × tax_mode → weight_price + volume_price + min_charge

type RoutePriceModel struct {
	ID               int64     `json:"id"`
	TenantID         int64     `json:"tenant_id"`
	ClientID         int64     `json:"client_id"`
	ClientName       string    `json:"client_name"`
	RouteName        string    `json:"route_name"`        // 线路名 e.g. "深圳→台湾(空运)"
	TransportType    string    `json:"transport_type"`    // air / sea / sea_express / air_special
	CargoType        string    `json:"cargo_type"`        // 普货 / 家具类 / 一类~六类 / 易碎品 / 特货
	TaxMode          string    `json:"tax_mode"`          // 全包税 / 频税 / 不包税
	WeightPrice      float64   `json:"weight_price"`      // ¥/kg — weight unit price
	VolumePrice      float64   `json:"volume_price"`      // ¥/才 — volume unit price (1才 ≈ 27872 cm³)
	MinCharge        float64   `json:"min_charge"`         // minimum charge per order
	FirstWeight      float64   `json:"first_weight"`       // first tier weight (kg)
	FirstWeightPrice float64   `json:"first_weight_price"` // first tier price
	ContWeightPrice  float64   `json:"cont_weight_price"`  // continuing weight price (¥/kg after first tier)
	FirstVolume      float64   `json:"first_volume"`       // first tier volume (才)
	FirstVolumePrice float64   `json:"first_volume_price"` // first tier volume price
	ContVolumePrice  float64   `json:"cont_volume_price"`  // continuing volume price
	VolumeCoeff      int       `json:"volume_coeff"`       // divisor for dim weight (6000 = 1m³=167kg)
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// ─── Tab 2: Storage Pricing ─────────────────────────────────────────
// 客户×仓储价: client × warehouse → storage_rate_per_day, free_days

type StoragePriceModel struct {
	ID             int64     `json:"id"`
	TenantID       int64     `json:"tenant_id"`
	ClientID       int64     `json:"client_id"`
	ClientName     string    `json:"client_name"`
	WarehouseID    int64     `json:"warehouse_id"`
	WarehouseName  string    `json:"warehouse_name"`
	FreeDays       int       `json:"free_days"`         // free storage days
	DailyRate      float64   `json:"daily_rate"`        // ¥/per-kg-per-day after free days
	MaxStorageDays int       `json:"max_storage_days"`  // auto-return after this many days
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ─── Tab 3: Delivery Fees ───────────────────────────────────────────
// 客户×派送费: client × carrier × customs_point × area × delivery_method → fee + free_threshold

type DeliveryFeeModel struct {
	ID             int64     `json:"id"`
	TenantID       int64     `json:"tenant_id"`
	ClientID       int64     `json:"client_id"`
	ClientName     string    `json:"client_name"`
	CarrierID      int64     `json:"carrier_id"`
	CarrierName    string    `json:"carrier_name"`    // 新竹物流 / 黑猫宅急便 / 顺丰速运
	CustomsPoint   string    `json:"customs_point"`   // 清關點: 台北 / 台中 / 高雄
	Area           string    `json:"area"`             // 區域: 預設 / 北部 / 中部 / 南部 / 东部
	DeliveryMethod string    `json:"delivery_method"` // 宅配 / 專車 / 自取
	Condition      string    `json:"condition"`        // 条件: 重量>39.8 / 单边长>=600或重量>=500
	Fee            float64   `json:"fee"`               // base fee in ¥
	FreeThreshold  float64   `json:"free_threshold"`   // free if weight >= this kg (0 = never free)
	FreeThresholdTxt string  `json:"free_threshold_txt"` // human-readable "≥10kg"
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ─── Tab 4: Surcharges ──────────────────────────────────────────────
// 客户×加收费: client × carrier × charge_type × tier × customs_point × area → price

type SurchargeModel struct {
	ID            int64     `json:"id"`
	TenantID      int64     `json:"tenant_id"`
	ClientID      int64     `json:"client_id"`
	ClientName    string    `json:"client_name"`
	CarrierID     int64     `json:"carrier_id"`
	CarrierName   string    `json:"carrier_name"`  // 新竹物流 / 黑猫宅急便
	ChargeType    string    `json:"charge_type"`   // 超長費 / 超材費 / 棧板費 / 偏遠費 / 上樓費
	Tier          string    `json:"tier"`           // 小板 / 大板 / — (for non-pallet charges)
	CustomsPoint  string    `json:"customs_point"` // 清關點: 台北 / 台中 / 高雄
	Area          string    `json:"area"`           // 區域
	Condition     string    `json:"condition"`      // trigger: 單邊>150cm / 體積重>實重2倍
	Price         float64   `json:"price"`          // 加收金额 (¥)
	PriceDesc     string    `json:"price_desc"`     // human-readable price description
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ─── Tab 5: Additional Services ─────────────────────────────────────
// 客户×附加服务: client × service_type → price

type ServicePriceModel struct {
	ID          int64     `json:"id"`
	TenantID    int64     `json:"tenant_id"`
	ClientID    int64     `json:"client_id"`
	ClientName  string    `json:"client_name"`
	ServiceType string    `json:"service_type"` // 木箱包装 / 开箱验货 / 拍照存证 / 换标 / 合箱
	ServiceCode string    `json:"service_code"` // WOODEN_CRATE / OPEN_INSPECT / PHOTO / RELABEL / MERGE
	UnitPrice   float64   `json:"unit_price"`   // unit price in ¥
	PriceMode   string    `json:"price_mode"`   // fixed / per_item / per_kg / per_order
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
