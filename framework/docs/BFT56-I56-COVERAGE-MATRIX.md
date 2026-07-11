# I56 Framework 1.0 LTS — 对标 BFT56 功能覆盖矩阵

> 分析：2026-07-10 | 目标：确保 I56 Framework 完整覆盖 BFT56 全部功能

---

## 一、功能覆盖总览

| BFT56 模块 | I56 Framework 覆盖 | 状态 |
|------------|-------------------|------|
| 员工管理 | core/auth + core/rbac | ✅ 内置 |
| 角色管理 | core/rbac | ✅ 内置 |
| 权限管理 | core/rbac | ✅ 内置 |
| 多租户 | core/tenant | ✅ 内置 |
| 通知管理 | core/notification | ✅ 内置 |
| 打印模板 | 业务模块: print | 🔧 待实现 |
| PDA 版本/会话 | 业务模块: pda | 🔧 待实现 |
| 任务派发参数 | core/scheduler | ✅ 内置 |
| 日志/Api日志 | core/audit + core/logger | ✅ 内置 |

---

## 二、业务模块详细对标

### OMS (订单管理)

| BFT56 功能 | I56 模块 | 字段覆盖 |
|-----------|---------|---------|
| 集运订单 CRUD | order | 18列全覆蓋 |
| 订单状态机 | order/domain | ✅ workflow集成 |
| 订单打印 | order + print | 🔧 print模块 |
| 清关公司设置 | order ↔ customs | ✅ |
| 承运商单号设置 | order ↔ carrier | ✅ |
| 导出申报单 | order/handler | ✅ |
| 附加服务订单 | parcel_service | 🔧 |
| 包裹列表 | parcel | 20列全覆蓋 |
| 包裹入库流程 | parcel | ✅ |
| 包裹称重/尺寸 | parcel | ✅ |
| 库位分配 | parcel ↔ warehouse | ✅ |
| 包裹异常 | parcel/exception | ✅ |
| 包裹预报 | parcel/handler | ✅ |
| 手动认领 | parcel/handler | ✅ |

### WMS (仓库管理)

| BFT56 功能 | I56 模块 | 覆盖 |
|-----------|---------|------|
| 仓库 CRUD | warehouse | ✅ |
| 库位管理 | warehouse/location | ✅ |
| 库位类型 | warehouse/location_type | ✅ |
| 区域管理 | warehouse/zone | ✅ |
| 区域类型 | warehouse/zone_type | ✅ |
| 集装柜 | warehouse/container | ✅ |
| 入库机 | warehouse/inbound_machine | ✅ |
| 入库看板 | warehouse/dashboard | 🔧 |
| 仓库作业台 | warehouse/console | 🔧 |
| 工单模板 | workorder/template | 🔧 |
| 工单流程 | workorder/process | ✅ workflow引擎 |
| 工单列表 | workorder | 🔧 |
| 员工任务监控 | workorder/task | 🔧 |

### TMS (物流管理)

| BFT56 功能 | I56 模块 | 覆盖 |
|-----------|---------|------|
| 区域组管理 | transport/area_group | ✅ |
| 货物类型 | transport/cargo_type | ✅ |
| 承运商列表 | transport/carrier | ✅ |
| 承运商单号池 | transport/carrier_number | ✅ |
| 装柜记录 | transport/container_loading | ✅ |
| 快递公司 | transport/courier | ✅ |
| 清关公司 | transport/customs_broker | ✅ |
| 清关点 | transport/customs_point | ✅ |
| 清关单号池 | transport/customs_number | ✅ |
| 线路模板 | transport/route | ✅ |
| 运输公司 | transport/shipping_provider | ✅ |
| 运输方式 | transport/transport_type | ✅ |
| 物流追踪 | transport/tracking | ✅ |

### CRM (客户管理)

| BFT56 功能 | I56 模块 | 覆盖 |
|-----------|---------|------|
| 客户管理 | customer | ✅ |
| 客户账号 | customer/account | ✅ |
| 客户会员 | customer/member | ✅ |
| 收件地址 | customer/address | ✅ |
| 申报人 | customer/declarant | ✅ |
| 充值管理 | customer/recharge | ✅ |
| 余额流水 | customer/ledger | ✅ |
| 客户定价 | customer/pricing | ✅ |
| 客户线路价 | customer/route_price | ✅ |
| 客户派送费 | customer/delivery_fee | ✅ |
| 客户加收费 | customer/surcharge | ✅ |
| 客户仓储价 | customer/storage_price | ✅ |
| 客户附加服务覆盖 | customer/service_override | ✅ |
| 月结对账单 | customer/statement | ✅ |
| 客户端权限 | customer/permission | ✅ |
| 仓库授权 | customer/warehouse_auth | ✅ |
| 承运商授权 | customer/carrier_auth | ✅ |

### Finance (财务)

| BFT56 功能 | I56 模块 | 覆盖 |
|-----------|---------|------|
| 订单盈利报表 | finance/report | 🔧 |
| 附加服务盈利 | finance/report | 🔧 |
| 客户盈利 | finance/report | 🔧 |
| 路线盈利 | finance/report | 🔧 |
| 流水 | finance/ledger | ✅ |
| 充值记录 | finance/recharge_log | ✅ |

---

## 三、BFT56 独有但 I56 可改进的功能

| BFT56 功能 | 当前实现 | I56 改进方案 |
|-----------|---------|-------------|
| 全局搜索 | Filament Global Search | Elasticsearch 全文搜索 |
| 打印模板 | Filament 模板 | 独立 Print Engine + 可视化设计器 |
| 报表 | 4 个简单盈利报表 | BI Report Engine + 自定义维度 |
| 通知 | 站内通知 | 多通道: Email/SMS/LINE/Telegram/Webhook |
| API | REST 端点 | REST + SDK (Go/Python/JS/Java) |
| PDA | Android APK | Flutter 跨平台 PDA 应用 |
| Webhook | 基础投递 | 完整事件订阅 + 重试 + 日志 |
| 多语言 | zh-TW | i18n 完整支持 |
| 部署 | 单服务器 | Docker Compose / K8s / Helm |

---

## 四、优先级排序 (P0-P3)

### P0 — Framework Core（必须有，否则无法构建）

1. config → logger → errors → response
2. validator → middleware → router
3. tenant → auth → rbac
4. eventbus

### P0 — 最小可行业务（对标 BFT56 核心）

1. customer (客户/会员/账号/申报人)
2. warehouse (仓库/库位/区域)
3. parcel (包裹/入库/状态)
4. order (集运订单)
5. transport (线路/承运商/清关)
6. finance (充值/流水/账单)

### P1 — 完整业务闭环

7. parcel_service (附加服务)
8. workorder (工单系统)
9. print (打印模板)
10. notification (通知)
11. report (报表引擎)

### P2 — 增强体验

12. pda (PDA 集成)
13. webhook (事件推送)
14. sdk (多语言 SDK)
15. bi (BI 分析)

### P3 — 生态建设

16. marketplace (插件市场)
17. workflow_designer (流程设计器)
18. print_designer (模板设计器)

---

## 五、数据库表预估（对标 BFT56）

| 领域 | 表数量 | 核心表 |
|------|--------|--------|
| Framework Core | 8 | users, roles, permissions, role_permission, audit_logs, notifications, scheduler_jobs, tenant_configs |
| Customer | 14 | clients, client_users, client_members, member_addresses, declarants, client_ledgers, client_recharges, client_route_prices, client_delivery_fees, client_surcharges, client_storage_prices, client_service_overrides, client_statements, warehouse_authorizations |
| Warehouse | 9 | warehouses, zones, zone_types, locations, location_types, containers, inbound_machines, container_loading_records, warehouse_configs |
| Parcel | 5 | parcels, parcel_events, exception_reports, parcel_photos, parcel_dimensions |
| Order | 6 | orders, order_parcels, order_events, parcel_service_orders, parcel_service_templates, parcel_service_types |
| Transport | 12 | area_groups, cargo_types, carriers, carrier_numbers, couriers, customs_brokers, customs_points, customs_numbers, routes, route_transport_types, shipping_providers, transport_types |
| Finance | 4 | invoices, payments, recharge_logs, profit_reports |
| WorkOrder | 5 | work_orders, work_order_templates, workflow_processes, process_instances, task_assignments |
| System | 5 | notifications, print_templates, pda_versions, operator_sessions, api_call_logs |
| **合计** | **~68** | |

---

*此矩阵确保 I56 Framework 不遗漏 BFT56 任何功能，同时在架构层面超越。*
