package main

import (
	"sync"
	"sync/atomic"
	"time"
)

// ═══════════════════════════════════════════════════════════════
// Admin Data Layer — in-memory stores with seed data
// ═══════════════════════════════════════════════════════════════

type Store[T any] struct {
	mu     sync.RWMutex
	items  []T
	nextID int64
}

func NewStore[T any]() *Store[T] { return &Store[T]{items: []T{}} }

func (s *Store[T]) NextID() int64 { return atomic.AddInt64(&s.nextID, 1) }

func (s *Store[T]) List() []T {
	s.mu.RLock(); defer s.mu.RUnlock()
	result := make([]T, len(s.items)); copy(result, s.items); return result
}

func (s *Store[T]) Add(item T) T {
	s.mu.Lock(); defer s.mu.Unlock()
	s.items = append(s.items, item); return item
}

func (s *Store[T]) Seed(items ...T) {
	s.mu.Lock(); defer s.mu.Unlock()
	s.items = append(s.items, items...)
	s.nextID = int64(len(s.items)) + 1
}

// ── Domain models ──

type AdminAreaGroup struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
	Desc string `json:"description"`
}

type AdminCargoType struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type AdminTransportMode struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type AdminCustomsBroker struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	License string `json:"license"`
	Contact string `json:"contact"`
	Phone   string `json:"phone"`
}

type AdminCustomsPoint struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Code    string `json:"code"`
	Port    string `json:"port"`
	Country string `json:"country"`
}

type AdminShippingProvider struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Code    string `json:"code"`
	Contact string `json:"contact"`
}

type AdminContainerLoading struct {
	ID          int64     `json:"id"`
	ContainerNo string    `json:"container_no"`
	Vessel      string    `json:"vessel"`
	PortFrom    string    `json:"port_from"`
	PortTo      string    `json:"port_to"`
	ParcelCount int       `json:"parcel_count"`
	LoadedAt    time.Time `json:"loaded_at"`
}

type AdminLogisticsTracking struct {
	ID         int64     `json:"id"`
	TrackingNo string    `json:"tracking_no"`
	Location   string    `json:"location"`
	Status     string    `json:"status"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type AdminClientAccount struct {
	ID       int64   `json:"id"`
	Username string  `json:"username"`
	RealName string  `json:"real_name"`
	Email    string  `json:"email"`
	Balance  float64 `json:"balance"`
	Status   string  `json:"status"`
}

type AdminClientRecharge struct {
	ID       int64     `json:"id"`
	ClientID int64     `json:"client_id"`
	Amount   float64   `json:"amount"`
	Method   string    `json:"method"`
	Remark   string    `json:"remark"`
	Time     time.Time `json:"time"`
}

type AdminClientPricing struct {
	ID       int64   `json:"id"`
	ClientID int64   `json:"client_id"`
	RouteID  int64   `json:"route_id"`
	Price    float64 `json:"price"`
	Discount float64 `json:"discount"`
}

type AdminClientPermission struct {
	ID       int64  `json:"id"`
	ClientID int64  `json:"client_id"`
	Module   string `json:"module"`
	CanRead  bool   `json:"can_read"`
	CanWrite bool   `json:"can_write"`
}

type AdminMonthlyStatement struct {
	ID         int64     `json:"id"`
	ClientID   int64     `json:"client_id"`
	Period     string    `json:"period"`
	Total      float64   `json:"total"`
	PaidAmount float64   `json:"paid_amount"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type AdminNotification struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Channel   string    `json:"channel"`
	Recipient string    `json:"recipient"`
	Sent      bool      `json:"sent"`
	CreatedAt time.Time `json:"created_at"`
}

type AdminPDAWorkorderTemplate struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	ProcessType string `json:"process_type"`
	Steps       int    `json:"steps"`
}

type AdminServiceTemplate struct {
	ID   int64   `json:"id"`
	Name string  `json:"name"`
	Type string  `json:"type"`
	Desc string  `json:"description"`
	Fee  float64 `json:"fee"`
}

type AdminServiceType struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type AdminServiceWorkorder struct {
	ID         int64     `json:"id"`
	OrderNo    string    `json:"order_no"`
	Type       string    `json:"type"`
	Status     string    `json:"status"`
	AssignedTo string    `json:"assigned_to"`
	CreatedAt  time.Time `json:"created_at"`
}

type AdminRouteTemplate struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	From      string `json:"from"`
	To        string `json:"to"`
	CarrierID int64  `json:"carrier_id"`
	EstDays   int    `json:"est_days"`
}

type AdminPricingService struct {
	ID   int64   `json:"id"`
	Name string  `json:"name"`
	Fee  float64 `json:"fee"`
	Unit string  `json:"unit"`
}

type AdminException struct {
	ID          int64     `json:"id"`
	ParcelID    int64     `json:"parcel_id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type AdminPDASession struct {
	ID         int64     `json:"id"`
	OperatorID int64     `json:"operator_id"`
	Device     string    `json:"device"`
	LoginAt    time.Time `json:"login_at"`
}

type AdminPrinter struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	IP   string `json:"ip"`
}

type AdminStorageConfig struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Provider string `json:"provider"`
	Bucket   string `json:"bucket"`
	Region   string `json:"region"`
}

type AdminSystemParam struct {
	ID    int64  `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
	Group string `json:"group"`
	Label string `json:"label"`
}

type AdminBrandSetting struct {
	ID    int64  `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
	Group string `json:"group"`
}

type AdminAPIConfig struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Provider string `json:"provider"`
	Endpoint string `json:"endpoint"`
	APIKey   string `json:"api_key"`
	Status   string `json:"status"`
}

type AdminAIChatMessage struct {
	ID      int64     `json:"id"`
	Role    string    `json:"role"`
	Content string    `json:"content"`
	Time    time.Time `json:"time"`
}

type AdminSchedulerJob struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Cron    string `json:"cron"`
	Enabled bool   `json:"enabled"`
	LastRun string `json:"last_run"`
}

type AdminAuditLog struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	Detail    string    `json:"detail"`
	CreatedAt time.Time `json:"created_at"`
}

type AdminReport struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type AdminNotificationChannel struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Config string `json:"config"`
}

type AdminInboundBoardEntry struct {
	ID         int64     `json:"id"`
	ParcelNo   string    `json:"parcel_no"`
	Warehouse  string    `json:"warehouse"`
	Status     string    `json:"status"`
	ExpectedAt time.Time `json:"expected_at"`
}

type AdminWarehouseBoardEntry struct {
	ID             int64 `json:"id"`
	PendingReceive int   `json:"pending_receive"`
	InStock        int   `json:"in_stock"`
	Picking        int   `json:"picking"`
	Outbound       int   `json:"outbound"`
}

type AdminWarehouseConsoleEntry struct {
	ID          int64  `json:"id"`
	WarehouseID int64  `json:"warehouse_id"`
	Name        string `json:"name"`
	Machine     string `json:"machine"`
	Status      string `json:"status"`
}

type AdminAIException struct {
	ID         int64     `json:"id"`
	ParcelID   int64     `json:"parcel_id"`
	Reason     string    `json:"reason"`
	Confidence float64   `json:"confidence"`
	Reviewed   bool      `json:"reviewed"`
	CreatedAt  time.Time `json:"created_at"`
}

type AdminExceptionReport struct {
	ID        int64     `json:"id"`
	Type      string    `json:"type"`
	Count     int       `json:"count"`
	Period    string    `json:"period"`
	CreatedAt time.Time `json:"created_at"`
}

// ── Global data stores ──

var (
	areaGroupStore           = NewStore[AdminAreaGroup]()
	cargoTypeStore           = NewStore[AdminCargoType]()
	transportModeStore       = NewStore[AdminTransportMode]()
	customsBrokerStore       = NewStore[AdminCustomsBroker]()
	customsPointStore        = NewStore[AdminCustomsPoint]()
	shippingProviderStore    = NewStore[AdminShippingProvider]()
	containerLoadingStore    = NewStore[AdminContainerLoading]()
	logisticsTrackingStore   = NewStore[AdminLogisticsTracking]()
	clientAccountStore       = NewStore[AdminClientAccount]()
	clientRechargeStore      = NewStore[AdminClientRecharge]()
	clientPricingStore       = NewStore[AdminClientPricing]()
	clientPermissionStore    = NewStore[AdminClientPermission]()
	monthlyStatementStore    = NewStore[AdminMonthlyStatement]()
	notificationStore        = NewStore[AdminNotification]()
	pdaWorkorderTplStore     = NewStore[AdminPDAWorkorderTemplate]()
	serviceTemplateStore     = NewStore[AdminServiceTemplate]()
	serviceTypeStore         = NewStore[AdminServiceType]()
	serviceWorkorderStore    = NewStore[AdminServiceWorkorder]()
	routeTemplateStore       = NewStore[AdminRouteTemplate]()
	pricingServiceStore      = NewStore[AdminPricingService]()
	exceptionStore           = NewStore[AdminException]()
	pdaSessionStore          = NewStore[AdminPDASession]()
	printerStore             = NewStore[AdminPrinter]()
	storageConfigStore       = NewStore[AdminStorageConfig]()
	systemParamStore         = NewStore[AdminSystemParam]()
	brandSettingStore        = NewStore[AdminBrandSetting]()
	apiConfigStore           = NewStore[AdminAPIConfig]()
	aiChatStore              = NewStore[AdminAIChatMessage]()
	schedulerJobStore        = NewStore[AdminSchedulerJob]()
	auditLogStore            = NewStore[AdminAuditLog]()
	reportStore              = NewStore[AdminReport]()
	notificationChannelStore = NewStore[AdminNotificationChannel]()
	inboundBoardStore        = NewStore[AdminInboundBoardEntry]()
	warehouseBoardStore      = NewStore[AdminWarehouseBoardEntry]()
	warehouseConsoleStore    = NewStore[AdminWarehouseConsoleEntry]()
	aiExceptionStore         = NewStore[AdminAIException]()
	exceptionReportStore     = NewStore[AdminExceptionReport]()
)

// ═══════════════════════════════════════════════════════════════
// Seed data — injected at startup
// ═══════════════════════════════════════════════════════════════

func seedAll() {
	now := time.Now()
	areaGroupStore.Seed(
		AdminAreaGroup{1, "华南区", "HN", "广东、广西、海南"},
		AdminAreaGroup{2, "华东区", "HD", "上海、江苏、浙江"},
		AdminAreaGroup{3, "华北区", "HB", "北京、天津、河北"},
	)
	cargoTypeStore.Seed(
		AdminCargoType{1, "普货", "GENERAL"},
		AdminCargoType{2, "易碎品", "FRAGILE"},
		AdminCargoType{3, "液体", "LIQUID"},
		AdminCargoType{4, "电子产品", "ELECTRONICS"},
		AdminCargoType{5, "食品", "FOOD"},
	)
	transportModeStore.Seed(
		AdminTransportMode{1, "海运", "SEA"},
		AdminTransportMode{2, "空运", "AIR"},
		AdminTransportMode{3, "陆运", "LAND"},
		AdminTransportMode{4, "铁路", "RAIL"},
	)
	customsBrokerStore.Seed(
		AdminCustomsBroker{1, "深圳报关行", "CB20230001", "张经理", "13800138001"},
		AdminCustomsBroker{2, "上海报关行", "CB20230002", "李经理", "13900139002"},
		AdminCustomsBroker{3, "广州报关行", "CB20230003", "王经理", "13700137003"},
	)
	customsPointStore.Seed(
		AdminCustomsPoint{1, "深圳蛇口", "SZKOU", "蛇口港", "中国"},
		AdminCustomsPoint{2, "上海外高桥", "SHWGQ", "外高桥港", "中国"},
		AdminCustomsPoint{3, "广州南沙", "GZNS", "南沙港", "中国"},
	)
	shippingProviderStore.Seed(
		AdminShippingProvider{1, "中远海运", "COSCO", "400-810-8888"},
		AdminShippingProvider{2, "马士基", "MAERSK", "400-820-8888"},
		AdminShippingProvider{3, "地中海航运", "MSC", "400-830-8888"},
	)
	containerLoadingStore.Seed(
		AdminContainerLoading{1, "COSU1234567", "东方号", "深圳蛇口", "洛杉矶", 120, now.Add(-72 * time.Hour)},
		AdminContainerLoading{2, "MAEU2345678", "海洋号", "上海外高桥", "鹿特丹", 85, now.Add(-48 * time.Hour)},
	)
	logisticsTrackingStore.Seed(
		AdminLogisticsTracking{1, "TRK20240001", "深圳集散中心", "已揽收", now.Add(-24 * time.Hour)},
		AdminLogisticsTracking{2, "TRK20240002", "广州转运中心", "运输中", now.Add(-12 * time.Hour)},
		AdminLogisticsTracking{3, "TRK20240003", "上海分拨中心", "派送中", now.Add(-3 * time.Hour)},
	)
	routeTemplateStore.Seed(
		AdminRouteTemplate{1, "深圳→洛杉矶", "深圳", "洛杉矶", 1, 18},
		AdminRouteTemplate{2, "上海→鹿特丹", "上海", "鹿特丹", 2, 25},
		AdminRouteTemplate{3, "广州→悉尼", "广州", "悉尼", 1, 15},
	)
	clientAccountStore.Seed(
		AdminClientAccount{1, "plat_ezjyt", "易捷物流", "admin@ezjyt.com", 50000.00, "active"},
		AdminClientAccount{2, "plat_szhy", "深圳华远", "admin@szhy.com", 20000.00, "active"},
		AdminClientAccount{3, "plat_shjl", "上海捷联", "admin@shjl.com", 35000.00, "active"},
	)
	clientRechargeStore.Seed(
		AdminClientRecharge{1, 1, 10000.00, "银行转账", "预充值", now.Add(-168 * time.Hour)},
		AdminClientRecharge{2, 2, 5000.00, "微信支付", "首充", now.Add(-72 * time.Hour)},
	)
	clientPricingStore.Seed(
		AdminClientPricing{1, 1, 1, 25.50, 0.90},
		AdminClientPricing{2, 2, 2, 35.00, 0.85},
	)
	clientPermissionStore.Seed(
		AdminClientPermission{1, 1, "parcels", true, true},
		AdminClientPermission{2, 1, "orders", true, true},
		AdminClientPermission{3, 2, "parcels", true, false},
	)
	monthlyStatementStore.Seed(
		AdminMonthlyStatement{1, 1, "2026-06", 28500.00, 28500.00, "已结清", now.Add(-720 * time.Hour)},
		AdminMonthlyStatement{2, 2, "2026-06", 15200.00, 10000.00, "部分付款", now.Add(-720 * time.Hour)},
	)
	exceptionStore.Seed(
		AdminException{1, 1, "破损", "外箱有明显压痕", "待处理", now.Add(-24 * time.Hour)},
		AdminException{2, 2, "短少", "应收3件实收2件", "处理中", now.Add(-12 * time.Hour)},
	)
	pdaSessionStore.Seed(
		AdminPDASession{1, 1, "PDA-001", now.Add(-8 * time.Hour)},
		AdminPDASession{2, 2, "PDA-002", now.Add(-4 * time.Hour)},
	)
	pdaWorkorderTplStore.Seed(
		AdminPDAWorkorderTemplate{1, "标准收货流程", "RECEIVE", 4},
		AdminPDAWorkorderTemplate{2, "标准出库流程", "OUTBOUND", 5},
		AdminPDAWorkorderTemplate{3, "退件处理流程", "RETURN", 3},
	)
	serviceTemplateStore.Seed(
		AdminServiceTemplate{1, "标准加固", "PACKAGING", "气泡膜+纸箱加固", 15.00},
		AdminServiceTemplate{2, "合并打包", "PACKAGING", "多件合并到一个包裹", 20.00},
		AdminServiceTemplate{3, "拍照验货", "INSPECTION", "1-3张照片", 5.00},
	)
	serviceTypeStore.Seed(
		AdminServiceType{1, "加固包装", "PACKAGING"},
		AdminServiceType{2, "验货拍照", "INSPECTION"},
		AdminServiceType{3, "分箱服务", "SPLIT"},
	)
	serviceWorkorderStore.Seed(
		AdminServiceWorkorder{1, "SW20240001", "PACKAGING", "pending", "OP001", now.Add(-24 * time.Hour)},
		AdminServiceWorkorder{2, "SW20240002", "INSPECTION", "completed", "OP002", now.Add(-48 * time.Hour)},
	)
	pricingServiceStore.Seed(
		AdminPricingService{1, "拆包服务", 10.00, "次"},
		AdminPricingService{2, "转寄服务", 25.00, "次"},
		AdminPricingService{3, "退件服务", 30.00, "次"},
	)
	notificationStore.Seed(
		AdminNotification{1, "系统升级通知", "系统将于本周六凌晨2点升级", "email", "all", false, now.Add(-48 * time.Hour)},
		AdminNotification{2, "新功能上线", "包裹追踪功能已上线", "sms", "all", true, now.Add(-72 * time.Hour)},
	)
	printerStore.Seed(
		AdminPrinter{1, "仓库A打印机", "热敏标签", "192.168.1.100"},
		AdminPrinter{2, "办公室打印机", "激光", "192.168.1.101"},
	)
	storageConfigStore.Seed(
		AdminStorageConfig{1, "包裹图片", "minio", "parcel-images", "cn-east-1"},
		AdminStorageConfig{2, "文档存储", "oss", "i56-documents", "cn-hangzhou"},
	)
	systemParamStore.Seed(
		AdminSystemParam{1, "site.name", "I56 WMS", "system", "站点名称"},
		AdminSystemParam{2, "site.logo", "/assets/logo.png", "system", "站点Logo"},
		AdminSystemParam{3, "order.auto_confirm", "false", "order", "自动确认订单"},
		AdminSystemParam{4, "parcel.max_weight", "30.0", "parcel", "包裹最大重量(kg)"},
	)
	brandSettingStore.Seed(
		AdminBrandSetting{1, "brand.primary_color", "#2563EB", "theme"},
		AdminBrandSetting{2, "brand.company_name", "I56 Framework", "company"},
	)
	apiConfigStore.Seed(
		AdminAPIConfig{1, "顺丰快递API", "SF", "https://sfapi.sf-express.com", "sf_api_key_xxx", "active"},
		AdminAPIConfig{2, "中国海关API", "CUSTOMS", "https://api.customs.gov.cn", "customs_key_xxx", "active"},
		AdminAPIConfig{3, "阿里云短信API", "ALIYUN_SMS", "https://dysmsapi.aliyuncs.com", "aliyun_key_xxx", "active"},
		AdminAPIConfig{4, "易联云打印API", "YLY", "https://open-api.10ss.net", "yly_key_xxx", "active"},
		AdminAPIConfig{5, "七牛云存储API", "QINIU", "https://up.qiniu.com", "qiniu_key_xxx", "active"},
	)
	aiChatStore.Seed(
		AdminAIChatMessage{1, "user", "帮我查询最近一周的订单量", now.Add(-1 * time.Hour)},
		AdminAIChatMessage{2, "assistant", "最近一周共有 125 个订单，其中已发货 98 个，待处理 27 个。", now.Add(-1*time.Hour + time.Second)},
	)
	schedulerJobStore.Seed(
		AdminSchedulerJob{1, "每日账单生成", "0 2 * * *", true, now.Add(-24*time.Hour).Format("2006-01-02 15:04")},
		AdminSchedulerJob{2, "库存同步", "0 */4 * * *", true, now.Add(-4*time.Hour).Format("2006-01-02 15:04")},
		AdminSchedulerJob{3, "数据备份", "0 3 * * *", false, "从未执行"},
	)
	auditLogStore.Seed(
		AdminAuditLog{1, 1, "login", "system", "管理员登录", now.Add(-1 * time.Hour)},
		AdminAuditLog{2, 1, "create_order", "orders", "创建订单 ORD20240001", now.Add(-30 * time.Minute)},
		AdminAuditLog{3, 2, "update_parcel", "parcels", "更新包裹状态为已入库", now.Add(-15 * time.Minute)},
	)
	reportStore.Seed(
		AdminReport{1, "月度运营报告", "monthly", "completed", now.Add(-720 * time.Hour)},
		AdminReport{2, "季度财务报表", "quarterly", "generating", now.Add(-48 * time.Hour)},
	)
	notificationChannelStore.Seed(
		AdminNotificationChannel{1, "系统邮件", "email", `{"smtp":"smtp.i56.com","port":587}`},
		AdminNotificationChannel{2, "短信通道", "sms", `{"provider":"aliyun","sign":"I56"}`},
	)
	inboundBoardStore.Seed(
		AdminInboundBoardEntry{1, "TRK001", "深圳仓", "已到港", now.Add(-24 * time.Hour)},
		AdminInboundBoardEntry{2, "TRK002", "上海仓", "清关中", now.Add(-12 * time.Hour)},
		AdminInboundBoardEntry{3, "TRK003", "广州仓", "运输中", now.Add(-36 * time.Hour)},
	)
	warehouseBoardStore.Seed(
		AdminWarehouseBoardEntry{1, 15, 342, 28, 12},
	)
	warehouseConsoleStore.Seed(
		AdminWarehouseConsoleEntry{1, 1, "深圳仓-1号机", "DWS-1000", "运行中"},
		AdminWarehouseConsoleEntry{2, 1, "深圳仓-2号机", "DWS-2000", "待机"},
	)
	aiExceptionStore.Seed(
		AdminAIException{1, 1, "包裹外箱损坏概率78%", 0.78, false, now.Add(-2 * time.Hour)},
		AdminAIException{2, 2, "申报品名与实际不符", 0.92, true, now.Add(-6 * time.Hour)},
	)
	exceptionReportStore.Seed(
		AdminExceptionReport{1, "包裹破损", 12, "2026-06", now.Add(-720 * time.Hour)},
		AdminExceptionReport{2, "申报异常", 5, "2026-06", now.Add(-720 * time.Hour)},
	)
}
