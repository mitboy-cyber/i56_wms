// Package demodata provides reusable demo data helpers for admin API endpoints.
package demodata

import "time"

// ─── Transport / Logistics ──────────────────────────────────────────

var Routes = []map[string]any{
	{"id": 1, "name": "厦门→台湾(空运)", "transport_type": "air", "from": "厦门", "to": "桃园", "base_price": 20.0, "min_days": 1, "max_days": 3, "is_active": true},
	{"id": 2, "name": "厦门→台湾(海快)", "transport_type": "sea_express", "from": "厦门", "to": "基隆", "base_price": 15.0, "min_days": 3, "max_days": 7, "is_active": true},
	{"id": 3, "name": "厦门→台湾(海运)", "transport_type": "sea", "from": "厦门", "to": "台北", "base_price": 12.0, "min_days": 5, "max_days": 14, "is_active": true},
	{"id": 4, "name": "厦门→台湾(经济)", "transport_type": "economy", "from": "厦门", "to": "台中", "base_price": 8.0, "min_days": 7, "max_days": 21, "is_active": false},
}

var Carriers = []map[string]any{
	{"id": 1, "name": "新竹物流", "code": "HCT", "contact": "林经理", "phone": "02-12345678"},
	{"id": 2, "name": "黑猫宅急便", "code": "YAMATO", "contact": "张经理", "phone": "02-87654321"},
	{"id": 3, "name": "全家超商取货", "code": "FAMILY", "contact": "", "phone": ""},
	{"id": 4, "name": "7-11超商取货", "code": "SEVEN", "contact": "", "phone": ""},
}

var Couriers = []map[string]any{
	{"id": 1, "name": "顺丰速运", "code": "SF"},
	{"id": 2, "name": "圆通快递", "code": "YT"},
	{"id": 3, "name": "中通快递", "code": "ZTO"},
	{"id": 4, "name": "韵达快递", "code": "YTO"},
	{"id": 5, "name": "DHL国际", "code": "DHL"},
}

const AreaGroups = `[{"id":1,"name":"北台湾","cities":["台北市","新北市","基隆市","桃园市"],"fee":10},{"id":2,"name":"中台湾","cities":["台中市","彰化县","南投县"],"fee":20},{"id":3,"name":"南台湾","cities":["高雄市","台南市","屏东县"],"fee":30}]`

const CargoTypes = `[{"id":1,"name":"普货","code":"general"},{"id":2,"name":"易碎品","code":"fragile"},{"id":3,"name":"液体","code":"liquid"},{"id":4,"name":"电子产品","code":"electronics"},{"id":5,"name":"食品","code":"food"},{"id":6,"name":"化妆品","code":"cosmetics"}]`

const TransportModes = `[{"id":1,"name":"空运","code":"air"},{"id":2,"name":"海运","code":"sea"},{"id":3,"name":"海快","code":"sea_express"},{"id":4,"name":"陆运","code":"land"},{"id":5,"name":"经济","code":"economy"}]`

// ─── Customer Management ───────────────────────────────────────────

var PricingRoutes = []map[string]any{
	{"id": 1, "name": "厦门→台湾(空运)首重", "type": "首重", "base_price": 20.0, "over_price": 30.0, "valid_from": "2026-01-01", "valid_to": "2026-12-31"},
	{"id": 2, "name": "厦门→台湾(海快)首重", "type": "首重", "base_price": 15.0, "over_price": 22.0, "valid_from": "2026-01-01", "valid_to": "2026-12-31"},
	{"id": 3, "name": "深圳→台湾(空运)首重", "type": "首重", "base_price": 22.0, "over_price": 32.0, "valid_from": "2026-01-01", "valid_to": "2026-12-31"},
}

var PricingDelivery = []map[string]any{
	{"id": 1, "carrier": "新竹物流", "region": "北台湾", "first_kg": 60, "add_kg": 20, "valid_from": "2026-01-01", "valid_to": "2026-12-31"},
	{"id": 2, "carrier": "新竹物流", "region": "中台湾", "first_kg": 80, "add_kg": 25, "valid_from": "2026-01-01", "valid_to": "2026-12-31"},
	{"id": 3, "carrier": "黑猫宅急便", "region": "全台湾", "first_kg": 90, "add_kg": 30, "valid_from": "2026-01-01", "valid_to": "2026-12-31"},
}

var Surcharges = []map[string]any{
	{"id": 1, "name": "偏远地区附加费", "type": "delivery", "amount": 50.0},
	{"id": 2, "name": "超大型件处理费", "type": "size", "amount": 100.0},
	{"id": 3, "name": "危险品处理费", "type": "dangerous", "amount": 200.0},
	{"id": 4, "name": "保价费(3%)", "type": "insurance", "amount": 0},
}

var ClientCredentials = []map[string]any{
	{"id": 1, "client": "EZ集運通", "api_key": "ez-api-key-demo-2026", "secret": "sk_live_ezjyt_demo_key_2026", "status": "active", "created_at": "2026-01-01T00:00:00Z"},
}

var LedgerEntries = []map[string]any{
	{"id": 1, "client_id": 1, "type": "recharge", "description": "在线充值-银行转账", "amount": 5000, "balance_after": 5000, "created_at": "2026-07-01T10:00:00+08:00"},
	{"id": 2, "client_id": 1, "type": "order_deduct", "description": "ORD-20260706-001 扣除", "amount": -98, "balance_after": 4902, "created_at": "2026-07-16T14:30:00+08:00"},
	{"id": 3, "client_id": 1, "type": "refund", "description": "ORD-20260705-001 退款", "amount": 9.5, "balance_after": 4911.5, "created_at": "2026-07-18T09:15:00+08:00"},
}

var RechargeRecords = []map[string]any{
	{"id": 1, "client": "EZ集運通", "amount": 5000, "method": "银行转账", "ref_no": "BANK-20260701001", "status": "已到账", "created_at": "2026-07-01T10:00:00+08:00"},
}

var MonthlyStatements = []map[string]any{
	{"id": 1, "client": "EZ集運通", "year": 2026, "month": 6, "period": "2026-06-01~2026-06-30", "order_count": 15, "total_amount": 480.0, "total_cost": 360.0, "profit": 120.0, "status": "已结清"},
}

// ─── Warehouse / WMS ───────────────────────────────────────────────

var Containers = []map[string]any{
	{"id": 1, "no": "CNTR-20260715-001", "warehouse": "厦门仓", "line": "海快-新竹物流", "status": "装货中", "max_weight": 5000, "created_at": time.Now().AddDate(0, 0, -1).Format("2006-01-02 15:04")},
	{"id": 2, "no": "CNTR-20260714-002", "warehouse": "厦门仓", "line": "海快-新竹物流", "status": "已发运", "max_weight": 4800, "created_at": time.Now().AddDate(0, 0, -2).Format("2006-01-02 15:04")},
	{"id": 3, "no": "CNTR-20260713-003", "warehouse": "厦门仓", "line": "空运-桃园", "status": "已完成", "max_weight": 3500, "created_at": time.Now().AddDate(0, 0, -7).Format("2006-01-02 15:04")},
	{"id": 4, "no": "CNTR-20260716-004", "warehouse": "厦门仓", "line": "海运-基隆", "status": "待装货", "max_weight": 8000, "created_at": time.Now().Format("2006-01-02 15:04")},
}

var PDASessions = []map[string]any{
	{"id": 1, "warehouse": "厦门仓", "operator": "出库员-小蓝", "device_sn": "PDA-A01", "login_at": "2026-07-05T21:17:30+08:00", "last_ping": time.Now().Format("2006-01-02T15:04:05+08:00"), "current_page": "称重", "current_zone": "出库区", "current_location": "A-01-01", "is_online": true, "logout_at": nil},
	{"id": 2, "warehouse": "厦门仓", "operator": "入库员-小王", "device_sn": "PDA-B03", "login_at": "2026-07-16T02:17:30+08:00", "last_ping": time.Now().Add(-1*time.Hour).Format("2006-01-02T15:04:05+08:00"), "current_page": "收货", "current_zone": "入库区", "current_location": "B-03-02", "is_online": true, "logout_at": nil},
	{"id": 3, "warehouse": "厦门仓", "operator": "拣货员-阿杰", "device_sn": "PDA-C05", "login_at": "2026-07-17T18:17:30+08:00", "last_ping": time.Now().Add(-2*time.Hour).Format("2006-01-02T15:04:05+08:00"), "current_page": "拣货", "current_zone": "拣货区", "current_location": "C-05-01", "is_online": true, "logout_at": nil},
	{"id": 6, "warehouse": "厦门仓", "operator": "上架员-小陈", "device_sn": "PDA-F03", "login_at": "2026-07-20T02:17:30+08:00", "last_ping": time.Now().Add(-8*time.Hour).Format("2006-01-02T15:04:05+08:00"), "current_page": "上架", "current_zone": "仓库A区", "current_location": "A-02-05", "is_online": false, "logout_at": time.Now().Format("2006-01-02T15:04:05+08:00")},
}

var Workflows = []map[string]any{
	{"id": 1, "name": "入库流程", "steps": "收货→核重→上架→入库完成", "module": "仓库管理", "status": "已启用", "updated_at": "2026-06-21T02:17:00+08:00"},
	{"id": 2, "name": "出库流程(大仓)", "steps": "拣货→打包→核重→送出库→装柜", "module": "仓库管理", "status": "已启用", "updated_at": "2026-06-26T02:17:00+08:00"},
	{"id": 3, "name": "退件处理流程", "steps": "签收→核验→入库→上架", "module": "仓库管理", "status": "已启用", "updated_at": "2026-07-06T02:17:00+08:00"},
}

var ServiceTypes = []map[string]any{
	{"id": 1, "name": "打包", "description": "标准包装", "price": 3.0},
	{"id": 2, "name": "加固", "description": "防震加固", "price": 8.0},
	{"id": 3, "name": "拍照", "description": "验货拍照", "price": 5.0},
	{"id": 4, "name": "合并", "description": "多包裹合并", "price": 10.0},
	{"id": 5, "name": "退货", "description": "退货处理", "price": 15.0},
	{"id": 6, "name": "销毁", "description": "商品销毁", "price": 20.0},
}

var ServiceTemplates = []map[string]any{
	{"id": 1, "name": "标准打包", "description": "纸箱+气泡膜标准打包", "price": 0.0, "type_ids": []int{1}},
	{"id": 2, "name": "易碎品加固", "description": "双层包装+防震+易碎标签", "price": 8.0, "type_ids": []int{2}},
	{"id": 3, "name": "验货拍照", "description": "开箱检查并拍照存档", "price": 5.0, "type_ids": []int{3}},
}

var ServiceWorkorders = []map[string]any{
	{"id": 1, "order_no": "ORD-20260706-001", "client": "EZ集運通", "template": "易碎品加固", "status": "completed", "assignee": "出库员-小蓝", "price": 8.0, "created_at": "2026-07-16T10:30:00+08:00"},
	{"id": 2, "order_no": "ORD-20260709-001", "client": "EZ集運通", "template": "标准打包", "status": "pending", "assignee": "", "price": 0.0, "created_at": "2026-07-20T14:00:00+08:00"},
}

var Notifications = []map[string]any{
	{"id": 1, "title": "系统维护通知", "type": "系统通知", "priority": "普通", "scope": "全员(跨所有公司)", "content": "系统将于7月20日维护升级", "channel": "站内信", "sender": "系统管理员", "sent": true, "send_time": "2026-07-20T02:17:30+08:00"},
	{"id": 2, "title": "新路线上线", "type": "公告", "priority": "紧急", "scope": "本公司全员", "content": "厦门→基隆 海运路线已开通", "channel": "邮件", "sender": "运营经理", "sent": true, "send_time": "2026-07-18T02:17:30+08:00"},
	{"id": 3, "title": "账户余额不足", "type": "任务通知", "priority": "普通", "scope": "指定用户", "content": "账户余额低于100元请及时充值", "channel": "短信", "sender": "财务系统", "sent": true, "send_time": "2026-07-19T02:17:30+08:00"},
	{"id": 4, "title": "包裹已签收", "type": "任务通知", "priority": "普通", "scope": "指定用户", "content": "包裹YT7625763166053已被签收", "channel": "微信", "sender": "物流系统", "sent": true, "send_time": "2026-07-20T14:17:30+08:00"},
	{"id": 5, "title": "假期运营安排", "type": "公告", "priority": "普通", "scope": "全员(跨所有公司)", "content": "端午假期正常运营部分路线时效延长", "channel": "站内信", "sender": "系统管理员", "sent": true, "send_time": "2026-07-21T02:17:30+08:00"},
}

var SystemParams = []map[string]any{
	{"id": 1, "key": "site_name", "value": "I56 WMS", "group": "system", "label": "站点名称"},
	{"id": 2, "key": "default_warehouse", "value": "厦门仓", "group": "system", "label": "默认仓库"},
	{"id": 3, "key": "currency", "value": "TWD", "group": "finance", "label": "结算货币"},
	{"id": 4, "key": "volume_factor", "value": "6000", "group": "logistics", "label": "体积重系数"},
	{"id": 5, "key": "tax_rate", "value": "0.05", "group": "finance", "label": "税率"},
	{"id": 6, "key": "max_package_weight", "value": "70", "group": "logistics", "label": "单件限重(kg)"},
	{"id": 7, "key": "declaration_limit", "value": "50000", "group": "customs", "label": "申报限额(NTD)"},
	{"id": 8, "key": "auto_cancel_hours", "value": "48", "group": "order", "label": "自动取消(小时)"},
	{"id": 9, "key": "free_storage_days", "value": "30", "group": "warehouse", "label": "免费仓储(天)"},
	{"id": 10, "key": "storage_fee_daily", "value": "5", "group": "warehouse", "label": "超期仓储费/天(NTD)"},
}

var Printers = []map[string]any{
	{"id": 1, "name": "入库标签打印机", "type": "thermal", "ip": "192.168.1.100", "status": "在线"},
	{"id": 2, "name": "出库面单打印机", "type": "thermal", "ip": "192.168.1.101", "status": "在线"},
	{"id": 3, "name": "发票打印机", "type": "laser", "ip": "192.168.1.102", "status": "离线"},
}

// Webhook logs
var Webhooks = []map[string]any{
	{"id": 1, "client": "EZ集運通", "event": "order.created", "url": "https://ezjyt.com/api/webhook/order", "status": "success", "response_code": 200, "created_at": time.Now().Add(-2*time.Hour).Format("2006-01-02T15:04:05+08:00")},
	{"id": 2, "client": "EZ集運通", "event": "parcel.received", "url": "https://ezjyt.com/api/webhook/order", "status": "success", "response_code": 200, "created_at": time.Now().Add(-3*time.Hour).Format("2006-01-02T15:04:05+08:00")},
}
