// Package domain provides seed data for the WMS backend.
package domain

import "time"

// SeedAll populates all in-memory stores with demo data.
func SeedAll() {
	now := time.Now()

	AreaGroupStore.Seed(
		AreaGroup{1, "华南区", "HN", "广东、广西、海南"},
		AreaGroup{2, "华东区", "HD", "上海、江苏、浙江"},
		AreaGroup{3, "华北区", "HB", "北京、天津、河北"},
	)
	CargoTypeStore.Seed(
		CargoType{1, "普货", "GENERAL"},
		CargoType{2, "易碎品", "FRAGILE"},
		CargoType{3, "液体", "LIQUID"},
		CargoType{4, "电子产品", "ELECTRONICS"},
		CargoType{5, "食品", "FOOD"},
	)
	TransportModeStore.Seed(
		TransportMode{1, "海运", "SEA"},
		TransportMode{2, "空运", "AIR"},
		TransportMode{3, "陆运", "LAND"},
		TransportMode{4, "铁路", "RAIL"},
	)
	CustomsBrokerStore.Seed(
		CustomsBroker{1, "深圳报关行", "CB20230001", "张经理", "13800138001"},
		CustomsBroker{2, "上海报关行", "CB20230002", "李经理", "13900139002"},
		CustomsBroker{3, "广州报关行", "CB20230003", "王经理", "13700137003"},
	)
	CustomsPointStore.Seed(
		CustomsPoint{1, "深圳蛇口", "SZKOU", "蛇口港", "中国"},
		CustomsPoint{2, "上海外高桥", "SHWGQ", "外高桥港", "中国"},
		CustomsPoint{3, "广州南沙", "GZNS", "南沙港", "中国"},
	)
	ShippingProviderStore.Seed(
		ShippingProvider{1, "中远海运", "COSCO", "400-810-8888"},
		ShippingProvider{2, "马士基", "MAERSK", "400-820-8888"},
		ShippingProvider{3, "地中海航运", "MSC", "400-830-8888"},
	)
	ContainerLoadingStore.Seed(
		ContainerLoading{1, "COSU1234567", "东方号", "深圳蛇口", "洛杉矶", 120, now.Add(-72 * time.Hour)},
		ContainerLoading{2, "MAEU2345678", "海洋号", "上海外高桥", "鹿特丹", 85, now.Add(-48 * time.Hour)},
	)
	LogisticsTrackingStore.Seed(
		LogisticsTracking{1, "TRK20240001", "深圳集散中心", "已揽收", now.Add(-24 * time.Hour)},
		LogisticsTracking{2, "TRK20240002", "广州转运中心", "运输中", now.Add(-12 * time.Hour)},
		LogisticsTracking{3, "TRK20240003", "上海分拨中心", "派送中", now.Add(-3 * time.Hour)},
	)
	RouteTemplateStore.Seed(
		RouteTemplate{1, "深圳→洛杉矶", "深圳", "洛杉矶", 1, 18},
		RouteTemplate{2, "上海→鹿特丹", "上海", "鹿特丹", 2, 25},
		RouteTemplate{3, "广州→悉尼", "广州", "悉尼", 1, 15},
	)
	ClientAccountStore.Seed(
		ClientAccount{1, "plat_ezjyt", "易捷物流", "admin@ezjyt.com", 50000.00, "active"},
		ClientAccount{2, "plat_szhy", "深圳华远", "admin@szhy.com", 20000.00, "active"},
		ClientAccount{3, "plat_shjl", "上海捷联", "admin@shjl.com", 35000.00, "active"},
	)
	ClientRechargeStore.Seed(
		ClientRecharge{1, 1, 10000.00, "银行转账", "预充值", now.Add(-168 * time.Hour)},
		ClientRecharge{2, 2, 5000.00, "微信支付", "首充", now.Add(-72 * time.Hour)},
	)
	ClientPricingStore.Seed(
		ClientPricing{1, 1, 1, 25.50, 0.90},
		ClientPricing{2, 2, 2, 35.00, 0.85},
	)
	ClientPermissionStore.Seed(
		ClientPermission{1, 1, "parcels", true, true},
		ClientPermission{2, 1, "orders", true, true},
		ClientPermission{3, 2, "parcels", true, false},
	)
	MonthlyStatementStore.Seed(
		MonthlyStatement{1, 1, "2026-06", 28500.00, 28500.00, "已结清", now.Add(-720 * time.Hour)},
		MonthlyStatement{2, 2, "2026-06", 15200.00, 10000.00, "部分付款", now.Add(-720 * time.Hour)},
	)
	ExceptionStore.Seed(
		Exception{1, 1, "破损", "外箱有明显压痕", "待处理", now.Add(-24 * time.Hour)},
		Exception{2, 2, "短少", "应收3件实收2件", "处理中", now.Add(-12 * time.Hour)},
	)
	PDASessionStore.Seed(
		PDASession{1, 1, "PDA-001", now.Add(-8 * time.Hour)},
		PDASession{2, 2, "PDA-002", now.Add(-4 * time.Hour)},
	)
	PDAWorkorderTplStore.Seed(
		PDAWorkorderTemplate{1, "标准收货流程", "RECEIVE", 4},
		PDAWorkorderTemplate{2, "标准出库流程", "OUTBOUND", 5},
		PDAWorkorderTemplate{3, "退件处理流程", "RETURN", 3},
	)
	ServiceTemplateStore.Seed(
		ServiceTemplate{1, "标准加固", "PACKAGING", "气泡膜+纸箱加固", 15.00},
		ServiceTemplate{2, "合并打包", "PACKAGING", "多件合并到一个包裹", 20.00},
		ServiceTemplate{3, "拍照验货", "INSPECTION", "1-3张照片", 5.00},
	)
	ServiceTypeStore.Seed(
		ServiceType{1, "加固包装", "PACKAGING"},
		ServiceType{2, "验货拍照", "INSPECTION"},
		ServiceType{3, "分箱服务", "SPLIT"},
	)
	ServiceWorkorderStore.Seed(
		ServiceWorkorder{1, "SW20240001", "PACKAGING", "pending", "OP001", now.Add(-24 * time.Hour)},
		ServiceWorkorder{2, "SW20240002", "INSPECTION", "completed", "OP002", now.Add(-48 * time.Hour)},
	)
	PricingServiceStore.Seed(
		PricingService{1, "拆包服务", 10.00, "次"},
		PricingService{2, "转寄服务", 25.00, "次"},
		PricingService{3, "退件服务", 30.00, "次"},
	)
	NotificationStore.Seed(
		Notification{1, "系统升级通知", "系统将于本周六凌晨2点升级", "email", "all", false, now.Add(-48 * time.Hour)},
		Notification{2, "新功能上线", "包裹追踪功能已上线", "sms", "all", true, now.Add(-72 * time.Hour)},
	)
	PrinterStore.Seed(
		Printer{1, "仓库A打印机", "热敏标签", "192.168.1.100"},
		Printer{2, "办公室打印机", "激光", "192.168.1.101"},
	)
	StorageConfigStore.Seed(
		StorageConfig{1, "包裹图片", "minio", "parcel-images", "cn-east-1"},
		StorageConfig{2, "文档存储", "oss", "i56-documents", "cn-hangzhou"},
	)
	SystemParamStore.Seed(
		SystemParam{1, "site.name", "I56 WMS", "system", "站点名称"},
		SystemParam{2, "site.logo", "/assets/logo.png", "system", "站点Logo"},
		SystemParam{3, "order.auto_confirm", "false", "order", "自动确认订单"},
		SystemParam{4, "parcel.max_weight", "30.0", "parcel", "包裹最大重量(kg)"},
	)
	BrandSettingStore.Seed(
		BrandSetting{1, "brand.primary_color", "#2563EB", "theme"},
		BrandSetting{2, "brand.company_name", "I56 Framework", "company"},
	)
	APIConfigStore.Seed(
		APIConfig{1, "顺丰快递API", "SF", "https://sfapi.sf-express.com", "sf_api_key_xxx", "active"},
		APIConfig{2, "中国海关API", "CUSTOMS", "https://api.customs.gov.cn", "customs_key_xxx", "active"},
		APIConfig{3, "阿里云短信API", "ALIYUN_SMS", "https://dysmsapi.aliyuncs.com", "aliyun_key_xxx", "active"},
		APIConfig{4, "易联云打印API", "YLY", "https://open-api.10ss.net", "yly_key_xxx", "active"},
		APIConfig{5, "七牛云存储API", "QINIU", "https://up.qiniu.com", "qiniu_key_xxx", "active"},
	)
	AIChatStore.Seed(
		AIChatMessage{1, "user", "帮我查询最近一周的订单量", now.Add(-1 * time.Hour)},
		AIChatMessage{2, "assistant", "最近一周共有 125 个订单，其中已发货 98 个，待处理 27 个。", now.Add(-1*time.Hour + time.Second)},
	)
	SchedulerJobStore.Seed(
		SchedulerJob{1, "每日账单生成", "0 2 * * *", true, now.Add(-24 * time.Hour).Format("2006-01-02 15:04")},
		SchedulerJob{2, "库存同步", "0 */4 * * *", true, now.Add(-4 * time.Hour).Format("2006-01-02 15:04")},
		SchedulerJob{3, "数据备份", "0 3 * * *", false, "从未执行"},
	)
	AuditLogStore.Seed(
		AuditLog{1, 1, "login", "system", "管理员登录", now.Add(-1 * time.Hour)},
		AuditLog{2, 1, "create_order", "orders", "创建订单 ORD20240001", now.Add(-30 * time.Minute)},
		AuditLog{3, 2, "update_parcel", "parcels", "更新包裹状态为已入库", now.Add(-15 * time.Minute)},
	)
	ReportStore.Seed(
		Report{1, "月度运营报告", "monthly", "completed", now.Add(-720 * time.Hour)},
		Report{2, "季度财务报表", "quarterly", "generating", now.Add(-48 * time.Hour)},
	)
	NotificationChannelStore.Seed(
		NotificationChannel{1, "系统邮件", "email", `{"smtp":"smtp.i56.com","port":587}`},
		NotificationChannel{2, "短信通道", "sms", `{"provider":"aliyun","sign":"I56"}`},
	)
	InboundBoardStore.Seed(
		InboundBoardEntry{1, "TRK001", "深圳仓", "已到港", now.Add(-24 * time.Hour)},
		InboundBoardEntry{2, "TRK002", "上海仓", "清关中", now.Add(-12 * time.Hour)},
		InboundBoardEntry{3, "TRK003", "广州仓", "运输中", now.Add(-36 * time.Hour)},
	)
	WarehouseBoardStore.Seed(
		WarehouseBoardEntry{1, 15, 342, 28, 12},
	)
	WarehouseConsoleStore.Seed(
		WarehouseConsoleEntry{1, 1, "深圳仓-1号机", "DWS-1000", "运行中"},
		WarehouseConsoleEntry{2, 1, "深圳仓-2号机", "DWS-2000", "待机"},
	)
	AIExceptionStore.Seed(
		AIException{1, 1, "包裹外箱损坏概率78%", 0.78, false, now.Add(-2 * time.Hour)},
		AIException{2, 2, "申报品名与实际不符", 0.92, true, now.Add(-6 * time.Hour)},
	)
	ExceptionReportStore.Seed(
		ExceptionReport{1, "包裹破损", 12, "2026-06", now.Add(-720 * time.Hour)},
		ExceptionReport{2, "申报异常", 5, "2026-06", now.Add(-720 * time.Hour)},
	)
}
