package domain

import "time"

// SeedAll populates all in-memory stores with BFT56-matching seed data.
func SeedAll() {
	now := time.Now()

	// NOTE: Employees & Warehouses are seeded via real repos in server.go

	// ── 角色 ──
	RoleStore.Seed(
		Role{1, "系统管理员", "拥有全部系统权限", true},
		Role{2, "仓库管理员", "管理仓库入库、出库、上架等操作", true},
		Role{3, "客服人员", "查看订单与包裹，处理客户咨询", true},
		Role{4, "财务人员", "查看财务报表与对账", true},
		Role{5, "操作员", "PDA扫码操作权限", true},
	)

	// ── 仓库看板/控制台 ──
	WarehouseBoardStore.Seed(
		WarehouseBoardEntry{1, 125, 3420, 15, 98},
		WarehouseBoardEntry{2, 98, 2800, 12, 76},
	)
	WarehouseConsoleStore.Seed(
		WarehouseConsoleEntry{1, 1, "入库组", "receive", "进行中"},
		WarehouseConsoleEntry{2, 1, "拣货组", "pick", "等待中"},
		WarehouseConsoleEntry{3, 1, "打包组", "pack", "进行中"},
		WarehouseConsoleEntry{4, 1, "出库组", "ship", "进行中"},
	)
	InboundBoardStore.Seed(
		InboundBoardEntry{1, "YT7625763166053", "厦门仓", "待上架", now},
		InboundBoardEntry{2, "9822467437512", "厦门仓", "已签收", now.Add(-1 * time.Hour)},
		InboundBoardEntry{3, "SF1234567890", "厦门仓", "已上架", now.Add(-3 * time.Hour)},
		InboundBoardEntry{4, "ZTO9876543210", "厦门仓", "已称重", now.Add(-4 * time.Hour)},
		InboundBoardEntry{5, "JD9999000011", "厦门仓", "待称重", now},
	)

	// ── PDA 在线会话 ──
	PDASessionStore.Seed(
		PDASession{1, "厦门仓", "出库员-小蓝", "PDA-A01", now.Add(-200*time.Hour), now, "200h", "packing", "出库区", "A-01-01", true, nil},
		PDASession{2, "厦门仓", "入库员-小王", "PDA-B03", now.Add(-120*time.Hour), now.Add(-1*time.Hour), "119h", "receive", "入库区", "B-03-02", true, nil},
		PDASession{3, "厦门仓", "拣货员-阿杰", "PDA-C05", now.Add(-80*time.Hour), now.Add(-2*time.Hour), "78h", "pick", "拣货区", "C-05-01", true, nil},
	)

	// ── PDA 工单模板 ──
	PDAWorkorderTplStore.Seed(
		PDAWorkorderTemplate{1, "厦门仓", "WT001", "入库上架", "入库", "WP-RECEIVE", 5, true, now.Add(-30*24*time.Hour)},
		PDAWorkorderTemplate{2, "厦门仓", "WT002", "拣货-按集运单", "拣货", "WP-BIG-DELIVER", 3, true, now.Add(-20*24*time.Hour)},
		PDAWorkorderTemplate{3, "厦门仓", "WT003", "打包装箱", "打包", "WP-BIG-DELIVER", 4, true, now.Add(-15*24*time.Hour)},
	)

	// ── 工单流程 ──
	WorkflowProcessStore.Seed(
		WorkflowProcess{1, "厦门仓", "WP-RECEIVE", "入库流程", "收货→核重→上架→入库完成", "包裹签收", true, now.Add(-30*24*time.Hour)},
		WorkflowProcess{2, "厦门仓", "WP-BIG-DELIVER", "出库流程(大仓)", "拣货→打包→核重→送出库→装柜", "订单创建", true, now.Add(-25*24*time.Hour)},
		WorkflowProcess{3, "厦门仓", "WP-RETURN", "退件处理流程", "签收→核验→入库→上架", "退件签收", true, now.Add(-15*24*time.Hour)},
	)

	// ── 集装柜 ──
	ContainerStore.Seed(
		Container{1, "厦门仓", "CNTR-20260715-001", "海快-新竹物流", "装货中", 5000, now.Add(-1*24*time.Hour)},
		Container{2, "厦门仓", "CNTR-20260714-002", "海快-新竹物流", "已发运", 4800, now.Add(-2*24*time.Hour)},
		Container{3, "厦门仓", "CNTR-20260713-003", "空运-桃园", "已完成", 3500, now.Add(-7*24*time.Hour)},
		Container{4, "厦门仓", "CNTR-20260716-004", "海运-基隆", "待装货", 8000, now},
	)

	// ── 客户 (BFT56: 4 clients, EZ集运通 main) ──
	ClientAccountStore.Seed(
		ClientAccount{1, "plat_ezjyt", "EZ集运通", "vip@ezjyt.com", 7670.40, "active"},
		ClientAccount{2, "plat_i56", "i56", "admin@i56.com", 0, "active"},
		ClientAccount{3, "plat_fb", "付呗", "fb@test.com", 0, "active"},
		ClientAccount{4, "plat_hgez", "嗨购EZ", "ez@test.com", 0, "active"},
	)
	ClientMemberStore.Seed(
		ClientMember{1, 1, "蕭惠昱", "0912345678", "A123456789", "EZ集运通", now.Add(-30*24*time.Hour)},
		ClientMember{2, 1, "王小明", "0923456789", "B234567890", "EZ集运通", now.Add(-25*24*time.Hour)},
		ClientMember{3, 1, "林大華", "0933334444", "F123456789", "EZ集运通", now.Add(-20*24*time.Hour)},
		ClientMember{4, 1, "陳怡君", "0944445555", "D223456789", "EZ集运通", now.Add(-15*24*time.Hour)},
		ClientMember{5, 1, "張志強", "0955556666", "E123456789", "EZ集运通", now.Add(-10*24*time.Hour)},
	)

	// ── 客户充值 ──
	ClientRechargeStore.Seed(
		ClientRecharge{1, 1, 5000, "银行转账", "银行转账充值", now.Add(-30*24*time.Hour)},
		ClientRecharge{2, 1, 3000, "微信支付", "微信充值", now.Add(-20*24*time.Hour)},
		ClientRecharge{3, 1, 2000, "银行转账", "7月预充值", now.Add(-5*24*time.Hour)},
		ClientRecharge{4, 2, 5000, "银行转账", "开户充值", now.Add(-25*24*time.Hour)},
	)

	// ── 余额日志 ──
	BalanceLogStore.Seed(
		BalanceLog{1, 1, "充值", 5000, 5000, "银行转账充值", now.Add(-30*24*time.Hour)},
		BalanceLog{2, 1, "消费", -329.60, 4670.40, "订单20260715120525777938", now.Add(-1*time.Hour)},
		BalanceLog{3, 1, "消费", -410.58, 4259.82, "订单20260708001", now.Add(-2*324*time.Hour)},
		BalanceLog{4, 1, "充值", 3000, 7259.82, "微信充值", now.Add(-20*24*time.Hour)},
		BalanceLog{5, 1, "消费", -189.42, 7070.40, "附加服务扣费", now.Add(-12*time.Hour)},
		BalanceLog{6, 1, "消费", -600.00, 6470.40, "集运订单批量扣费", now.Add(-6*time.Hour)},
		BalanceLog{7, 1, "充值", 2000, 8470.40, "7月预充值", now.Add(-5*24*time.Hour)},
		BalanceLog{8, 1, "消费", -800.00, 7670.40, "7月账单结算", now.Add(-3*24*time.Hour)},
	)

	// ── 充值记录 ──
	RechargeRecordStore.Seed(
		RechargeRecord{1, 1, 5000, "银行转账", "已完成", now.Add(-30*24*time.Hour)},
		RechargeRecord{2, 1, 3000, "微信支付", "已完成", now.Add(-20*24*time.Hour)},
	)

	// ── 客户端权限 (BFT56: 根据客户等级控制模块 + 增删改查导出权限) ──
	ClientPanelPermStore.Seed(
		// EZ集运通 — enterprise 全功能
		ClientPanelPerm{1, 1, "EZ集运通", "订单管理", "我的订单", true, true, true, true, true, "enterprise", "active", now.Add(-180*24*time.Hour), now.Add(365*24*time.Hour), "企业级全功能"},
		ClientPanelPerm{2, 1, "EZ集运通", "订单管理", "包裹列表", true, false, true, false, true, "enterprise", "active", now.Add(-180*24*time.Hour), now.Add(365*24*time.Hour), ""},
		ClientPanelPerm{3, 1, "EZ集运通", "订单管理", "集运下单", true, true, true, true, false, "enterprise", "active", now.Add(-180*24*time.Hour), now.Add(365*24*time.Hour), ""},
		ClientPanelPerm{4, 1, "EZ集运通", "申报管理", "申报人管理", true, true, true, true, false, "enterprise", "active", now.Add(-180*24*time.Hour), now.Add(365*24*time.Hour), ""},
		ClientPanelPerm{5, 1, "EZ集运通", "财务管理", "余额明细", true, false, false, false, true, "enterprise", "active", now.Add(-180*24*time.Hour), now.Add(365*24*time.Hour), ""},
		ClientPanelPerm{6, 1, "EZ集运通", "地址管理", "地址簿", true, true, true, true, false, "enterprise", "active", now.Add(-180*24*time.Hour), now.Add(365*24*time.Hour), ""},
		ClientPanelPerm{7, 1, "EZ集运通", "会员管理", "子账户", true, true, true, true, false, "enterprise", "active", now.Add(-180*24*time.Hour), now.Add(365*24*time.Hour), ""},
		ClientPanelPerm{8, 1, "EZ集运通", "系统设置", "API密钥", true, false, true, false, false, "enterprise", "active", now.Add(-180*24*time.Hour), now.Add(365*24*time.Hour), "只读+编辑API密钥"},
		// i56平台 — pro 中等级
		ClientPanelPerm{9, 2, "i56平台", "订单管理", "我的订单", true, true, true, false, true, "pro", "active", now.Add(-90*24*time.Hour), now.Add(180*24*time.Hour), "专业版"},
		ClientPanelPerm{10, 2, "i56平台", "申报管理", "申报人管理", true, true, true, false, false, "pro", "active", now.Add(-90*24*time.Hour), now.Add(180*24*time.Hour), ""},
		ClientPanelPerm{11, 2, "i56平台", "财务管理", "余额明细", true, false, false, false, false, "pro", "active", now.Add(-90*24*time.Hour), now.Add(180*24*time.Hour), "只读余额"},
		// 测试客户A — basic (已过期)
		ClientPanelPerm{12, 3, "测试客户A", "订单管理", "我的订单", true, false, false, false, false, "basic", "expired", now.Add(-60*24*time.Hour), now.Add(-1*24*time.Hour), "已过期需续费"},
		ClientPanelPerm{13, 3, "测试客户A", "财务管理", "余额明细", true, false, false, false, false, "basic", "expired", now.Add(-60*24*time.Hour), now.Add(-1*24*time.Hour), "已过期"},
	)

	// ── 客户定价 ──
	ClientPricingStore.Seed(
		ClientPricing{1, 1, 1, 15.0, 0.9},
		ClientPricing{2, 1, 2, 8.30, 0.95},
		ClientPricing{3, 1, 3, 3.20, 1.0},
		ClientPricing{4, 2, 1, 18.0, 0.85},
	)

	// ── 月度对账单 ──
	MonthlyStatementStore.Seed(
		MonthlyStatement{1, 1, "2026-06", 8560.30, 8560.30, "已结算", now},
		MonthlyStatement{2, 1, "2026-07", 12345.67, 0, "待结算", now},
		MonthlyStatement{3, 2, "2026-07", 4520.00, 0, "待结算", now},
	)

	// ── TMS: 区域组 (BFT56: 6 regions) ──
	AreaGroupStore.Seed(
		AreaGroup{1, "台湾线", "TW", "台湾地区全境"},
		AreaGroup{2, "离岛", "TW-OFF", "澎湖/金门/马祖"},
		AreaGroup{3, "东部", "TW-E", "宜兰/花莲/台东"},
		AreaGroup{4, "南部", "TW-S", "高雄/台南/屏东/嘉义"},
		AreaGroup{5, "中部", "TW-C", "台中/彰化/南投/云林/苗栗"},
		AreaGroup{6, "北部", "TW-N", "台北/新北/基隆/桃园/新竹"},
	)
	// ── 线路模板 (BFT56: 5 templates) ──
	RouteTemplateStore.Seed(
		RouteTemplate{1, "空运-台湾", "厦门", "台北", 1, 2},
		RouteTemplate{2, "海快-台湾", "厦门", "台北", 1, 3},
		RouteTemplate{3, "海运-台湾", "厦门", "台中", 1, 10},
		RouteTemplate{4, "空运特货-台湾", "厦门", "桃园", 1, 3},
		RouteTemplate{5, "商业海快-台湾", "厦门", "台北", 1, 5},
	)
	// ── 货物类型 (BFT56: key types) ──
	CargoTypeStore.Seed(
		CargoType{1, "普货", "GENERAL"},
		CargoType{2, "特货", "SPECIAL"},
		CargoType{3, "海快普货", "SEA_FAST"},
		CargoType{4, "空运特货", "AIR_SPECIAL"},
		CargoType{5, "食品", "FOOD"},
		CargoType{6, "生活用品", "DAILY"},
		CargoType{7, "文创用品", "STATIONERY"},
		CargoType{8, "电子商品", "ELECTRONICS"},
		CargoType{9, "化妆品", "COSMETICS"},
		CargoType{10, "易碎品", "FRAGILE"},
		CargoType{11, "液体", "LIQUID"},
	)
	// ── 运输方式 ──
	TransportModeStore.Seed(
		TransportMode{1, "海运", "SEA"},
		TransportMode{2, "空运", "AIR"},
		TransportMode{3, "海快", "SEA_FAST"},
		TransportMode{4, "空运特货", "AIR_SPECIAL"},
		TransportMode{5, "商业海快", "BIZ_FAST"},
	)
	// ── 清关公司 ──
	CustomsBrokerStore.Seed(
		CustomsBroker{1, "顺达报关行", "XM-BK-001", "张经理", "0592-2222222"},
		CustomsBroker{2, "德通实业", "XM-DT-001", "李经理", "0592-3333333"},
		CustomsBroker{3, "台通报关", "TW-TG-001", "王经理", "02-23456789"},
	)
	// ── 清关点 ──
	CustomsPointStore.Seed(
		CustomsPoint{1, "厦门海关", "CNXMN", "厦门", "中国"},
		CustomsPoint{2, "台北海关", "TWTPE", "台北", "台湾"},
		CustomsPoint{3, "台中海关", "TWTCG", "台中", "台湾"},
	)
	// ── 运输公司 ──
	ShippingProviderStore.Seed(
		ShippingProvider{1, "新竹物流", "HCT", "0920000001"},
		ShippingProvider{2, "中远海运", "COSCO", "021-12345678"},
		ShippingProvider{3, "万海航运", "WHL", "02-25679888"},
	)

	// ── 物流追踪 ──
	LogisticsTrackingStore.Seed(
		LogisticsTracking{1, "YT7625763166053", "厦门→台北", "运输中", now},
		LogisticsTracking{2, "435212825957725", "厦门→桃园", "已签收", now.Add(-1*24*time.Hour)},
	)
	ContainerLoadingStore.Seed(
		ContainerLoading{1, "CNTR-001", "厦门轮", "厦门", "台北", 450, now},
		ContainerLoading{2, "CNTR-002", "EVA AIR Cargo", "深圳", "桃园", 320, now.Add(-2 * 24 * time.Hour)},
		ContainerLoading{3, "CNTR-003", "WAN HAI 235", "厦门", "基隆", 580, now.Add(-1 * 24 * time.Hour)},
	)

	// ── 通知 ──
	NotificationStore.Seed(
		Notification{1, "系统上线通知", "I56 WMS 系统正式上线运行", "system", "all", true, now},
		Notification{2, "订单异常提醒", "订单20260715120525777938需人工审核", "alert", "admin", false, now},
	)

	// ── 定时任务 ──
	SchedulerJobStore.Seed(
		SchedulerJob{1, "每日账单生成", "0 8 * * *", true, now.Format("2006-01-02")},
		SchedulerJob{2, "数据库备份", "0 2 * * *", true, now.Format("2006-01-02")},
	)

	// ── 审计日志 ──
	AuditLogStore.Seed(
		AuditLog{1, 1, "login", "system", "管理员登录系统", now.Add(-2*time.Hour)},
		AuditLog{2, 1, "create", "order", "创建订单20260715120525777938", now.Add(-1*time.Hour)},
	)

	// ── AI Chat ──
	AIChatStore.Seed(
		AIChatMessage{1, "assistant", "你好，我是 I56 智能助手，有什么可以帮你？", time.Now()},
		AIChatMessage{2, "user", "查询今天入库包裹数量", time.Now()},
	)

	// ── 系统参数 ──
	SystemParamStore.Seed(
		SystemParam{1, "site_name", "I56 WMS", "system", "站点名称"},
		SystemParam{2, "default_warehouse", "厦门仓", "system", "默认仓库"},
		SystemParam{3, "currency", "TWD", "finance", "结算货币"},
		SystemParam{4, "volume_factor", "6000", "logistics", "体积重系数"},
		SystemParam{5, "tax_rate", "0.05", "finance", "税率"},
		SystemParam{6, "max_package_weight", "70", "logistics", "单件限重(kg)"},
		SystemParam{7, "declaration_limit", "50000", "customs", "申报限额(NTD)"},
		SystemParam{8, "auto_cancel_hours", "48", "order", "自动取消(小时)"},
		SystemParam{9, "free_storage_days", "30", "warehouse", "免费仓储(天)"},
		SystemParam{10, "storage_fee_daily", "5", "warehouse", "超期仓储费/天(NTD)"},
	)

	// ── 品牌 ──
	BrandSettingStore.Seed(
		BrandSetting{1, "logo", "/static/logo.png", "brand"},
		BrandSetting{2, "theme_color", "#1a73e8", "brand"},
	)

	// ── 存储配置 ──
	StorageConfigStore.Seed(
		StorageConfig{1, "minio", "MinIO", "i56-bucket", "us-east-1"},
		StorageConfig{2, "s3-backup", "AWS S3", "i56-backup", "ap-southeast-1"},
	)

	// ── 通知渠道 ──
	NotificationChannelStore.Seed(
		NotificationChannel{1, "邮件通知", "email", "smtp://smtp.example.com"},
		NotificationChannel{2, "短信通知", "sms", "aliyun-sms"},
		NotificationChannel{3, "企业微信", "wechat_work", "webhook://wx.example.com"},
	)
	PrinterStore.Seed(
		Printer{1, "主打印机", "USB", "192.168.1.100"},
		Printer{2, "标签打印机", "Network", "192.168.1.101"},
	)

	// ── 设备 (扫码枪/地磅/入库机/PDA/打印机) ──
	DeviceStore.Seed(
		Device{1, "扫码枪-01", "barcode_scanner", "SCAN-WH-001", "192.168.1.50", "online", 1},
		Device{2, "扫码枪-02", "barcode_scanner", "SCAN-WH-002", "192.168.1.51", "online", 1},
		Device{3, "地磅-01", "weighbridge", "WEIGH-WH-001", "192.168.1.10", "online", 1},
		Device{4, "入库机-01", "inbound_machine", "INB-WH-001", "192.168.1.20", "online", 1},
		Device{5, "PDA-01", "pda", "PDA-WH-001", "192.168.1.70", "online", 1},
		Device{6, "PDA-02", "pda", "PDA-WH-002", "192.168.1.71", "offline", 1},
		Device{7, "打印机-01", "printer", "PRT-WH-001", "192.168.1.100", "online", 1},
		Device{8, "打印机-02", "printer", "PRT-WH-002", "192.168.1.101", "online", 1},
	)

	// ── 仓位/货架 ──
	ShelfStore.Seed(
		Shelf{1, 1, "A-01-01", "A区", "01", 1, "empty"},
		Shelf{2, 1, "A-01-02", "A区", "01", 2, "occupied"},
		Shelf{3, 1, "A-02-01", "A区", "02", 1, "empty"},
		Shelf{4, 1, "B-01-01", "B区", "01", 1, "occupied"},
		Shelf{5, 1, "B-02-01", "B区", "02", 1, "empty"},
	)

	// ── API 配置 ──
	APIConfigStore.Seed(
		APIConfig{1, "顺丰快递", "SF", "https://api.sf.com", "sk-sf-test", "active"},
		APIConfig{2, "圆通快递", "YTO", "https://api.yto.com", "sk-yto-test", "active"},
	)

	// ── 附加服务 (BFT56: 10 templates, 5 types) ──
	ServiceTemplateStore.Seed(
		ServiceTemplate{1, "标准加固", "PACKAGING", "气泡膜+纸箱加固", 15.00},
		ServiceTemplate{2, "合并打包", "PACKAGING", "多件合并到一个包裹", 20.00},
		ServiceTemplate{3, "打木箱", "PACKAGING", "打木箱加固(需重新入库)", 80.00},
		ServiceTemplate{4, "打木架", "PACKAGING", "打木架加固(需重新入库)", 50.00},
		ServiceTemplate{5, "包装气柱袋", "PACKAGING", "气泡柱包装防护", 10.00},
		ServiceTemplate{6, "退货寄回", "RETURN", "退件寄回原地址", 2.00},
		ServiceTemplate{7, "易碎品贴纸", "LABEL", "易碎品标识贴纸", 0.10},
		ServiceTemplate{8, "内容物拍照", "INSPECTION", "包裹内容物拍照", 1.00},
		ServiceTemplate{9, "清点数量", "INSPECTION", "清点包裹内物品数量", 0.10},
		ServiceTemplate{10, "外箱标识拍照", "INSPECTION", "外箱标识与标签拍照", 1.00},
	)
	ServiceTypeStore.Seed(
		ServiceType{1, "加固包装", "PACKAGING"},
		ServiceType{2, "验货拍照", "INSPECTION"},
		ServiceType{3, "退货处理", "RETURN"},
		ServiceType{4, "标签贴纸", "LABEL"},
		ServiceType{5, "合并打包", "CONSOLIDATE"},
	)

	// ── 异常记录 ──
	ExceptionStore.Seed(
		Exception{1, 1, "破损", "外箱有明显压痕", "待处理", now.Add(-24*time.Hour)},
		Exception{2, 2, "短少", "应收3件实收2件", "处理中", now.Add(-12*time.Hour)},
		Exception{3, 3, "面单脱落", "快递面单脱落无法识别", "待处理", now.Add(-6*time.Hour)},
	)
	AIExceptionStore.Seed(
		AIException{1, 1, "包裹重量异常：申报1kg实际3.5kg", 0.85, false, now.Add(-12*time.Hour)},
	)
	ExceptionReportStore.Seed(
		ExceptionReport{1, "破损", 3, "2026-07", now},
	)

	// ── 报表 (BFT56: 4 profit reports) ──
	ReportStore.Seed(
		Report{1, "集运订单盈利-2026-07", "order_profit", "completed", now},
		Report{2, "附加服务盈利-2026-07", "service_profit", "completed", now},
		Report{3, "客户盈利汇总-2026-07", "client_profit", "completed", now},
		Report{4, "路线盈利汇总-2026-07", "route_profit", "completed", now},
	)

	// ── 计费 ──
	PricingServiceStore.Seed(
		PricingService{1, "加固包装", 15.00, "per_item"},
		PricingService{2, "海快基础费", 20.00, "per_kg"},
		PricingService{3, "空运基础费", 35.00, "per_kg"},
	)

	// ── 服务工单 ──
	ServiceWorkorderStore.Seed(
		ServiceWorkorder{1, "WO-20260715-001", "加固", "pending", "操作员01", now.Add(-2 * time.Hour)},
		ServiceWorkorder{2, "WO-20260715-002", "拍照", "processing", "操作员02", now.Add(-1 * time.Hour)},
		ServiceWorkorder{3, "WO-20260713-001", "打包", "completed", "操作员01", now.Add(-3 * 24 * time.Hour)},
		ServiceWorkorder{4, "WO-20260714-001", "退货", "pending", "", now.Add(-1 * 24 * time.Hour)},
	)

	// ── 客户权限 (补充) ──
	// ── 计费 ──
	RechargeRecordStore.Seed(
		RechargeRecord{1, 1, 5000, "银行转账", "completed", now.Add(-30 * 24 * time.Hour)},
		RechargeRecord{2, 1, 3000, "微信支付", "completed", now.Add(-20 * 24 * time.Hour)},
		RechargeRecord{3, 1, 2000, "银行转账", "completed", now.Add(-5 * 24 * time.Hour)},
	)
}
