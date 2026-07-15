# I56 Framework 1.0 LTS — Product Requirements Document

> **版本**: 1.0.0 · **语言**: Go 1.24+ · **模式**: Modular Monolith · **对标**: BFT56 八方云仓
> **作者**: Peter · **最后更新**: 2026-07-15

---

## 一、竞品分析摘要 (BFT56 八方云仓)

### 1.1 系统规模
- **管理后台**: 48 页，7 大模块组 + 首页/仓库看板
- **客户端**: 10+ 页面，含钱包、订单、包裹、会计、定价、API
- **PDA 端**: 收货、上架、拣货、打包、异常上报
- **技术栈**: Filament PHP 3 + Livewire + MySQL + Redis

### 1.2 BFT56 模块全景

| 模块组 | 子页面数 | 对标 BFT56 路由 |
|--------|---------|-----------------|
| 首页/看板 | 2 | `/admin`, `/admin/warehouse-board` |
| 订单管理 | 2 | `/admin/o-m-s/orders`, `/admin/o-m-s/parcel-service-orders` |
| 仓库管理 | 14 | parcels, containers, warehouses, inbound-board, warehouse-console, work-orders, work-order-templates, workflow-processes, work-order-lists, exception-reports, operator-sessions, parcel-service-templates, parcel-service-type, parcel-service-order-items |
| 财务报表 | 4 | order-profit, service-profit, client-profit, route-profit |
| 物流管理 | 10 | area-groups, cargo-types, carriers, couriers, customs-brokers, customs-clearance-points, routes, shipping-providers, transport-types, logistics-trackings |
| 客户管理 | 11 | client-member-addresses, declarants, clients, client-users, client-members, client-recharges, client-ledgers, client-recharge-logs, client-route-prices, client-statements, client-panel-permissions |
| 系统管理 | 5 | notifications, print-templates, roles, users, system-settings |

### 1.3 BFT56 关键特性 (I56 当前缺失)

| 特性 | BFT56 实现 | I56 现状 |
|------|-----------|---------|
| 客户会员体系 | Member + Address + Declarant 三层 | ❌ 缺失 |
| 充值/余额/授信 | 充值→余额日志→授信额度→月结 | ❌ 仅存根 |
| 财务报表 | 4 维盈利分析 (订单/服务/客户/路线) | ❌ 缺失 |
| 线路模板 | Route + TransportType + CargoType 组合 | ❌ 缺失 |
| 承运商/快递/清关 | 三级物流服务商管理 | ❌ 缺失 |
| 工单系统 | Workflow + Template + List | ❌ 缺失 |
| 集装柜 | Container 管理 | ❌ 缺失 |
| 列管理 | 动态显示/隐藏表格列 | ❌ 缺失 |
| 打印模板 | 面单打印 | ❌ 缺失 |
| 搜索/筛选 | 全局搜索 + 高级筛选 | ❌ 简易搜索 |
| 多标签页 | 浏览器式 Tab 导航 | ❌ 缺失 |
| 导出功能 | 导出申报单 CSV | ❌ 缺失 |

### 1.4 客户端功能对比

| 特性 | BFT56 客户端 | I56 现状 |
|------|------------|---------|
| 仪表盘 | KPI 卡片 (余额/额度/包裹统计) | ✅ 基础 |
| 余额/授信 | 实时余额 + 授信额度 | ❌ |
| 包裹管理 | 列表 + 预报 + 状态流转 | ⚠️ 基础 |
| 订单管理 | 集运订单 + 附加服务 | ⚠️ 基础 |
| 线路价格 | 客户专属价格表 | ❌ |
| 承运商价格 | 派送价 + 加收价 | ❌ |
| Webhook | Webhook 投递配置 | ❌ |
| API 凭证 | OpenAPI Key 管理 | ❌ |

---

## 二、I56 Framework 1.0 LTS 产品定义

### 2.1 三层架构

```
┌──────────────────────────────────────────────┐
│  Applications (应用产品)                       │
│  i56-wms │ i56-oms │ i56-tms │ i56-crm        │
│  i56-erp │ i56-finance │ i56-pda              │
├──────────────────────────────────────────────┤
│  Business Modules (业务模块)                    │
│  WMS │ OMS │ TMS │ CRM │ Finance │ Report      │
│  System │ OpenAPI │ Notification │ PDA         │
├──────────────────────────────────────────────┤
│  Framework Core (核心框架)                       │
│  Auth │ RBAC │ Tenant │ Workflow │ EventBus     │
│  Cache │ Logger │ Config │ Storage │ Scheduler  │
│  Validator │ Response │ Exception │ Audit       │
└──────────────────────────────────────────────┘
```

### 2.2 三个独立仓库

| 仓库 | 用途 | 依赖 |
|------|------|------|
| **i56-framework** | Core 框架 + 共享工具包 | 无 |
| **i56-admin** | 通用后台 (auth/RBAC/tenant/workflow) | i56-framework |
| **i56-apps** | 业务应用 (wms/oms/tms/crm/finance/pda) | i56-framework + i56-admin |

---

## 三、Framework Core 1.0 (P0 — 必须交付)

### 3.1 Core 模块清单

```
core/
├── app/          # 应用生命周期管理
├── auth/         # 认证 (Session/JWT/OAuth2)
├── rbac/         # 角色权限 + DataScope
├── tenant/       # 多租户 (Shared DB / Schema / DB Per Tenant)
├── workflow/     # 审批工作流引擎
├── eventbus/     # 事件总线 (内存 + Redis/RabbitMQ)
├── cache/        # 缓存抽象 (Memory/Redis)
├── logger/       # 结构化日志 (zap/slog)
├── config/       # 配置管理 (YAML + 环境变量)
├── storage/      # 文件存储抽象 (Local/MinIO/S3/OSS)
├── scheduler/    # 定时任务 (无需 Linux Cron)
├── validator/    # 参数校验 (i18n)
├── response/     # 统一 API 响应格式
├── exception/    # 统一错误处理
├── middleware/    # HTTP 中间件集合
├── router/       # 路由注册抽象
├── database/     # 数据库抽象 + 迁移
├── queue/        # 消息队列抽象
└── notification/ # 通知 (Email/SMS/Webhook)
```

### 3.2 技术选型

| 组件 | 选型 | 理由 |
|------|------|------|
| HTTP 框架 | Gin | 性能最高，生态成熟 |
| ORM | GORM v2 | Go 生态最完整 |
| 日志 | zap | 零分配，结构化 |
| 配置 | viper | 多格式，热重载 |
| 缓存 | go-redis + sync.Map | 两级缓存 |
| 事件总线 | watermill | 统一消息抽象 |
| 迁移 | golang-migrate | 纯 SQL，支持多 DB |
| 校验 | go-playground/validator | i18n 支持 |
| JWT | golang-jwt v5 | 标准实现 |

---

## 四、业务模块详细规划 (P1 — 对标 BFT56)

### 4.1 系统管理模块 (System)

| 页面 | 路由 | 功能 | BFT56 对标 |
|------|------|------|-----------|
| 角色管理 | `/admin/system/roles` | CRUD + 权限树 + 数据范围 | ✅ `/admin/shield/roles` |
| 员工管理 | `/admin/system/users` | CRUD + 仓库绑定 + 角色分配 + 状态 | ✅ `/admin/system/users` |
| 系统参数 | `/admin/system/settings` | 全局配置 KV | ✅ `/admin/system-settings-page` |
| 通知管理 | `/admin/system/notifications` | 通知模板 + 发送记录 | ✅ `/admin/system/notifications` |
| 打印模板 | `/admin/system/print-templates` | 面单打印模板 | ✅ 新增对标 |
| 审计日志 | `/admin/system/audit-logs` | 操作审计 | I56 已有 |

### 4.2 仓库管理模块 (WMS)

| 页面 | 路由 | 功能 | BFT56 对标 |
|------|------|------|-----------|
| 仓库列表 | `/admin/wms/warehouses` | CRUD | ✅ |
| 仓库看板 | `/admin/warehouse-board` | KPI 总览 | ✅ |
| 入库看板 | `/admin/inbound-board` | 入库跟踪 | ✅ |
| 仓库作业台 | `/admin/warehouse-console` | 操作员作业界面 | ✅ |
| 包裹列表 | `/admin/wms/parcels` | 全量包裹 + 筛选 + 列管理 | ✅ |
| 集装柜管理 | `/admin/wms/containers` | CRUD + 装柜操作 | ✅ 新增 |
| 员工任务监控 | `/admin/wms/task-monitor` | 任务分配 + 进度 | ✅ |
| 异常记录 | `/admin/wms/exceptions` | 异常列表 | ✅ |

### 4.3 订单管理模块 (OMS)

| 页面 | 路由 | 功能 | BFT56 对标 |
|------|------|------|-----------|
| 集运订单 | `/admin/oms/orders` | 全量 + 搜索 + 筛选 + 导出 | ✅ |
| 附加服务工单 | `/admin/oms/service-orders` | 服务工单列表 | ✅ |
| 附加服务模板 | `/admin/oms/service-templates` | 服务模板 CRUD | ✅ |
| 附加服务类型 | `/admin/oms/service-types` | 服务类型 CRUD | ✅ |

### 4.4 物流管理模块 (TMS)

| 页面 | 路由 | 功能 | BFT56 对标 |
|------|------|------|-----------|
| 区域组管理 | `/admin/tms/area-groups` | CRUD | ✅ |
| 货物类型 | `/admin/tms/cargo-types` | CRUD | ✅ |
| 承运商列表 | `/admin/tms/carriers` | CRUD | ✅ 新增 |
| 快递公司 | `/admin/tms/couriers` | CRUD | ✅ |
| 清关公司 | `/admin/tms/customs-brokers` | CRUD | ✅ |
| 清关点管理 | `/admin/tms/customs-points` | CRUD | ✅ 新增 |
| 线路模板 | `/admin/tms/routes` | CRUD + 运输方式 + 货类 | ✅ 新增 |
| 运输公司 | `/admin/tms/shipping-providers` | CRUD | ✅ 新增 |
| 运输方式 | `/admin/tms/transport-types` | CRUD | ✅ 新增 |
| 物流追踪 | `/admin/tms/tracking` | 物流状态查询 | ✅ 新增 |

### 4.5 客户管理模块 (CRM)

| 页面 | 路由 | 功能 | BFT56 对标 |
|------|------|------|-----------|
| 客户管理 | `/admin/crm/clients` | CRUD + 类型 + 结算方式 | ✅ |
| 客户账号 | `/admin/crm/client-accounts` | 关联账号管理 | ✅ |
| 客户会员 | `/admin/crm/client-members` | 会员 CRUD | ✅ 新增 |
| 客户收件地址 | `/admin/crm/member-addresses` | 地址管理 | ✅ 新增 |
| 客户申报人 | `/admin/crm/declarants` | 申报人 CRUD | ✅ 新增 |
| 客户充值 | `/admin/crm/recharges` | 充值操作 | ✅ 新增 |
| 余额日志 | `/admin/crm/ledgers` | 账目明细 | ✅ 新增 |
| 充值记录 | `/admin/crm/recharge-logs` | 充值记录 | ✅ 新增 |
| 客户价格 | `/admin/crm/route-prices` | 客户专属定价 | ✅ 新增 |
| 月结对账单 | `/admin/crm/statements` | 月度对账 | ✅ |
| 客户端权限 | `/admin/crm/client-permissions` | 客户端功能权限 | ✅ 新增 |

### 4.6 财务报表模块 (Finance)

| 页面 | 路由 | 功能 | BFT56 对标 |
|------|------|------|-----------|
| 集运订单盈利 | `/admin/finance/order-profit` | 订单维度盈利 | ✅ |
| 附加服务盈利 | `/admin/finance/service-profit` | 服务维度盈利 | ✅ |
| 客户盈利 | `/admin/finance/client-profit` | 客户维度盈利 | ✅ |
| 路线盈利 | `/admin/finance/route-profit` | 路线维度盈利 | ✅ |

### 4.7 PDA 端

| 页面 | 路由 | 功能 | BFT56 对标 |
|------|------|------|-----------|
| PDA 仪表盘 | `/pda` | 待处理任务 | ✅ |
| 收货 | `/pda/receive` | 扫描入库 | ✅ |
| 上架 | `/pda/shelve` | 库位上架 | ✅ |
| 拣货 | `/pda/pick` | 订单拣货 | ✅ |
| 打包 | `/pda/pack` | 包裹打包 | ✅ |
| 异常上报 | `/pda/exception` | 拍照 + 异常类型 | ✅ |
| PDA 会话 | `/admin/pda/sessions` | 操作员在线管理 | ✅ |
| PDA 工单模板 | `/admin/pda/work-order-templates` | 模板管理 | ✅ |

### 4.8 客户端 (Client Portal)

| 页面 | 路由 | 功能 | BFT56 对标 |
|------|------|------|-----------|
| 主控台 | `/client` | KPI 仪表盘 | ✅ |
| 收件地址 | `/client/addresses` | CRUD | ✅ |
| 客户会员 | `/client/members` | CRUD | ✅ |
| 申报人 | `/client/declarants` | CRUD | ✅ |
| 我的订单 | `/client/orders` | 订单列表 + 详情 | ✅ |
| 我的包裹 | `/client/parcels` | 包裹列表 + 预报 | ✅ |
| 附加服务订单 | `/client/service-orders` | 服务订单 | ✅ |
| 余额明细 | `/client/ledger` | 账目流水 | ✅ |
| 月结对账单 | `/client/statements` | 月结账单 | ✅ |
| 仓库信息 | `/client/warehouse-info` | 仓库地址 | ✅ |
| 线路价格 | `/client/route-prices` | 价格查询 | ✅ |
| 承运商派送价 | `/client/carrier-delivery` | 末端价格 | ✅ |
| 承运商加收价 | `/client/carrier-surcharge` | 附加费 | ✅ |
| Webhook 投递 | `/client/webhooks` | 回调配置 | ✅ |
| API 凭证 | `/client/api-keys` | Key 管理 | ✅ |

---

## 五、数据模型核心实体

### 5.1 实体关系图 (核心)

```
Tenant
  ├── Warehouse (1:N)
  ├── Employee (1:N)
  ├── Role → Permission (N:M)
  ├── Client (1:N)
  │   ├── ClientUser (1:N)
  │   ├── ClientMember (1:N)
  │   │   └── MemberAddress (1:N)
  │   ├── Declarant (1:N)
  │   ├── ClientRecharge (1:N)
  │   ├── ClientLedger (1:N)
  │   ├── RoutePrice (1:N)
  │   ├── MonthlyStatement (1:N)
  │   └── Parcel (1:N)
  │       └── ServiceOrder (1:N)
  ├── ConsolidationOrder (1:N)
  │   └── Parcel (N:M via OrderParcel)
  ├── Route (1:N)
  │   ├── TransportType
  │   └── CargoType (N:M)
  ├── Carrier (1:N)
  ├── Courier (1:N)
  ├── CustomsBroker (1:N)
  ├── CustomsClearancePoint (1:N)
  ├── Container (1:N)
  │   └── Parcel (N:M)
  ├── WorkOrder (1:N)
  │   ├── WorkOrderTemplate
  │   └── WorkflowProcess
  ├── ExceptionReport (1:N)
  └── Notification (1:N)
```

### 5.2 关键数据字段 (对标 BFT56)

#### 集运订单 (ConsolidationOrder)
```
id, order_no, warehouse_id, client_id, client_member_id,
recipient_name, courier_tracking_no, route_id, tax_type,
parcel_count, order_status, container_id, total_actual_weight,
total_chargeable_weight, total_price, customs_declaration_no,
carrier_order_no, remark, manual_remark, created_at
```

#### 包裹 (Parcel)
```
id, parcel_no, warehouse_id, courier_id, client_id, client_member_id,
cargo_type_id, consolidation_order_id, product_name, transport_type_id,
parcel_name, actual_weight, length, width, height, status,
has_exception, has_service, location_code, last_operation, created_at
```

#### 客户 (Client)
```
id, name, client_type (大客户/普通), member_count, parcel_count,
settlement_method (预付/月结), balance, credit_limit, status, created_at
```

#### 线路模板 (Route)
```
id, enabled, warehouse_id, route_name, transport_type_id,
cargo_types (多选), created_at
```

---

## 六、实施路线图 (Roadmap)

### Phase 1: Framework Core (2 周)
- [x] 项目骨架 + Go Module 初始化
- [x] Auth 模块 (Session + JWT)
- [x] RBAC + DataScope
- [x] Tenant Provider
- [x] Logger + Config + Validator
- [ ] EventBus (watermill)
- [ ] Cache 抽象 (go-redis)
- [ ] Storage 抽象
- [ ] Scheduler 引擎
- [ ] Workflow 引擎
- [ ] Notification 中心

### Phase 2: 业务模块完善 (3 周)
- [ ] WMS: 仓库+看板+入库+作业台+集装柜+工单
- [ ] OMS: 订单+服务模板+导出
- [ ] TMS: 区域+货物+承运商+快递+清关+线路+运输+物流追踪
- [ ] CRM: 客户+会员+地址+申报人+充值+余额+价格+账单+权限
- [ ] Finance: 4 维盈利分析
- [ ] System: 角色+员工+参数+通知+打印+审计

### Phase 3: 前端升级 (2 周)
- [ ] 统一组件库: Table (列管理+排序+筛选), Form, Modal, Tab
- [ ] 搜索: 全局搜索 + 页面级搜索 + 高级筛选
- [ ] 列管理: 动态显示/隐藏/排序表格列
- [ ] 导出: CSV/XLSX 导出
- [ ] 多标签页导航
- [ ] 响应式适配 Mobile/PDA

### Phase 4: 客户端 + PDA + OpenAPI (1 周)
- [ ] 客户端: 钱包/余额/账单/价格/Webhook/API 凭证
- [ ] PDA: 收货/上架/拣货/打包/异常
- [ ] OpenAPI: RESTful API + 鉴权 + 限流 + 文档

---

## 七、质量指标

| 指标 | 目标 |
|------|------|
| Admin 页面 | 50+ 页 (对标 BFT56 的 48 页) |
| Client 页面 | 15 页 |
| PDA 页面 | 8 页 |
| 后端 API | 150+ 端点 |
| 单元测试覆盖 | 70%+ |
| API 响应时间 | P99 < 200ms |
| 页面加载时间 | FCP < 1.5s |
| 浏览器兼容 | Chrome/Firefox/Safari/Edge |

---

## 八、差异化优势 (vs BFT56)

| 维度 | BFT56 | I56 Framework 1.0 |
|------|-------|-------------------|
| 技术栈 | PHP/Filament/Livewire | Go 1.24+ / React / TypeScript |
| 架构 | Monolith | Modular Monolith → 可拆微服务 |
| 多租户 | 隐式 | 显式 Tenant Provider |
| 扩展性 | Filament Plugin | Framework Module 系统 |
| 开源 | 闭源 SaaS | 开源 Framework + 商业应用 |
| 部署 | SaaS Only | SaaS + 私有化 + K8s |
| 性能 | PHP-FPM 限制 | Go 原生并发 |
| 工作流 | 无 | 内置审批引擎 |
| 事件驱动 | 无 | EventBus 事件总线 |
