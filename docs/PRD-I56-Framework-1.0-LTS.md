# I56 Framework 1.0 LTS — 产品需求文档 (PRD)

## 元信息
- **版本**: 1.0.0-draft
- **作者**: Peter (系统架构师)
- **状态**: 起草中
- **对标竞品**: BFT56 (https://bft56.com)
- **分析深度**: 完整管理后台76页 + 客户端19页

---

## 一、竞品分析摘要 (BFT56)

### 1.1 系统概况
BFT56 是一个集运物流 SaaS 平台，支持中国大陆→台湾的跨境集运业务。

| 维度 | 分析 |
|------|------|
| 架构 | PHP (推断), 单体应用, 前后端不分离 |
| 多租户 | Single DB + Shared Tables (tenant_id 隔离) |
| 用户角色 | 超管 / 平台客户 / 仓库操作员 / 申报人 |
| 技术栈 | Bootstrap + jQuery + 服务端渲染 |
| 部署 | 单机 Docker |

### 1.2 功能模块矩阵 (BFT56 完整审计)

#### 管理层 (admin) — 76 页, 7 组 49 菜单项

| 分组 | 菜单项 | 核心功能 |
|------|--------|---------|
| 首页 | 仪表盘 | KPI 卡片 (订单/包裹/营收/客户) |
| 订单管理 | 集运订单/服务订单 | 订单列表/详情/操作 |
| 仓库管理 | 包裹列表/入库看板/仓库作业台/集装柜/附加服务/PDA | 入库→上架→拣货→打包→出库 全流程 |
| 财务报表 | 集运订单盈利/服务盈利/客户盈利/路线盈利/对账单 | KPI 报表 + 图表 |
| 物流管理 | 快递公司/承运商/路线模板/运输方式/清关/物流追踪 | TMS 全链路 |
| 客户管理 | 客户列表/会员/账本/充值/报价/权限/申报人/地址 | CRM |
| 系统 | 角色/通知/打印机/存储/参数/定时任务/审计/AI | 系统管理 |

#### 客户端 (client) — 19 页

| 页面 | 功能 |
|------|------|
| 我的订单 | 创建/查看/追踪集运订单 |
| 我的包裹 | 入库包裹列表 |
| 申报人管理 | 台湾实名认证 (身份证+电话号码) |
| 余额明细 | 消费/充值记录 |
| 地址管理 | 台湾地址簿 |
| 会员管理 | 客户子账户 |
| 快递/路线/计费查询 | 参考信息 |

### 1.3 BFT56 核心业务流

```
客户下单 → 大陆卖家发货 → 包裹入库(厦门仓) → 称重测量
→ 上架存储 → 客户提交集运 → 拣货 → 打包 → 称重计费
→ 海关申报 → 装柜 → 海运/空运 → 台湾清关 → 派送 → 签收
```

### 1.4 BFT56 关键数据模型 (审计提取)

```
Tenant (租户)
  ├── Client (客户) ─ 余额/信用/等级
  │   ├── ClientMember (会员) ─ 电话/证件号
  │   ├── Declarant (申报人) ─ 台湾身份证/认证状态
  │   ├── Address (地址) ─ 台湾六都+县市
  │   ├── Ledger / Recharge / BalanceLog (财务)
  │   └── ClientPricing (客户专属报价)
  ├── Warehouse (仓库) ─ 厦门仓
  │   ├── Parcel (包裹) ─ 快递单号/状态流转
  │   ├── Shelf / Location (仓位)
  │   └── Container (集装柜) ─ 柜号/船名/起运港
  ├── Order (集运订单)
  │   ├── Status Flow: pending_picking → picking → packing → weighing → customs → loaded → shipped → completed
  │   ├── TrackingNumbers (多快递单号)
  │   └── Route (路线: 空运/海快/海运)
  ├── Route (线路) ─ 运输方式/价格
  ├── Courier (快递公司) ─ 顺丰/新竹/大荣...
  ├── ServiceOrder (附加服务: 加固/拍照/贴纸...)
  └── Finance
      ├── Statement (月结账单)
      └── Profit Report (盈利报表)
```

---

## 二、I56 Framework 1.0 LTS 定位

### 2.1 愿景

> **一套 Framework，多种业务模块，多个行业产品**

```
I56 Framework 1.0 LTS
  ├── I56 WMS    (仓储管理)
  ├── I56 OMS    (订单管理)
  ├── I56 TMS    (运输管理)
  ├── I56 ERP    (企业资源)
  ├── I56 CRM    (客户管理)
  └── I56 Finance (财务管理)
```

### 2.2 三层架构

```
Applications (产品)
  ├── 物流 SaaS
  ├── 海外仓系统
  └── 跨境电商 ERP

Business Modules (业务模块)
  ├── WMS  ─ Warehouse Management
  ├── OMS  ─ Order Management
  ├── TMS  ─ Transport Management
  ├── CRM  ─ Customer Management
  ├── FIN  ─ Financial Management
  └── RPT  ─ Report Engine

Framework (框架核心)
  ├── Auth / RBAC / Tenant
  ├── EventBus / Workflow / Scheduler
  ├── Cache / Logger / Config / Validator
  ├── Storage / Notification / OpenAPI
  └── HTTP Gateway / Router / Middleware
```

### 2.3 技术规格

| 项目 | 规格 |
|------|------|
| 语言 | Go 1.24+ |
| 架构 | Modular Monolith → 可演进微服务 |
| 数据库 | MySQL 8.0+ (多租户隔离) |
| 缓存 | Redis 7.0+ |
| 消息 | RabbitMQ / NATS |
| 存储 | MinIO (S3 兼容) |
| 搜索 | Elasticsearch 8.x |
| 前端 | Bootstrap 5 + HTMX + Alpine.js |
| 部署 | Docker / Kubernetes / Helm |

---

## 三、Framework 核心模块设计

### 3.1 Core 包结构

```
core/
├── app/          # 应用生命周期
├── auth/         # 认证 (JWT/Session)
├── cache/        # 缓存抽象
├── config/       # 配置管理
├── database/     # DB 抽象 + 迁移
├── errors/       # 错误处理
├── events/       # 事件总线
├── logger/       # 结构化日志
├── middleware/    # HTTP 中间件
├── queue/        # 消息队列
├── response/     # 统一响应格式
├── router/       # HTTP 路由
├── scheduler/    # 定时任务
├── security/     # RBAC + 数据权限
├── storage/      # 对象存储
├── tenant/       # 多租户 Provider
├── validator/    # 参数校验
└── workflow/     # 审批流引擎
```

### 3.2 认证鉴权 (Auth + RBAC)

```
Tenant → Department → Role → Permission → DataScope

DataScope 数据权限:
  - ALL         全部数据
  - ENTERPRISE  企业级
  - WAREHOUSE   仓库级
  - DEPARTMENT  部门级
  - SELF        本人
```

### 3.3 事件总线 (EventBus)

所有模块间通信必须通过事件：

```
OrderCreated → Publish → Finance (记账)
                       → Inventory (库存)
                       → Notification (通知)
                       → Webhook (回调)
                       → Audit (审计)
```

### 3.4 多租户 (Tenant)

支持三种隔离策略，通过 `TenantProvider` 接口抽象：
- **Shared Tables**: 所有租户共享表，`tenant_id` 字段隔离
- **Schema Per Tenant**: PostgreSQL 独立 Schema
- **Database Per Tenant**: 独立数据库

---

## 四、业务模块规划

### 4.1 WMS (仓库管理)

| 聚合 | 实体 | 说明 |
|------|------|------|
| Warehouse | 仓库/仓位/库区 | 厦门仓 → 台北仓 |
| Parcel | 包裹/状态流转 | 入库→称重→上架→拣货→打包→出库 |
| Container | 集装柜/装箱单 | 装柜→发运→到港 |
| Service | 附加服务/工单模板 | 加固/拍照/退货 |
| PDA | 设备/会话/任务 | 扫码操作 |

### 4.2 OMS (订单管理)

| 聚合 | 实体 | 说明 |
|------|------|------|
| Order | 集运订单/行项目 | 多包裹合并 |
| Route | 路线/运输方式 | 空运/海快/海运 |
| Tracking | 物流追踪 | 多快递单号 |

### 4.3 CRM (客户管理)

| 聚合 | 实体 | 说明 |
|------|------|------|
| Client | 平台客户/余额 | 4类客户 |
| Member | 会员子账户 | 实名认证 |
| Declarant | 申报人 | 台湾身份证 |
| Address | 地址簿 | 台湾六都 |

### 4.4 Finance (财务管理)

| 聚合 | 实体 | 说明 |
|------|------|------|
| Ledger | 账本/流水 | 充值+消费 |
| Statement | 月结对账单 | 月度汇总 |
| Recharge | 充值记录 | 支付渠道 |
| Profit | 盈利报表 | 订单/服务/客户/路线 |

---

## 五、数据库设计 (核心表)

### 5.1 多租户基础字段

```sql
-- 所有业务表必须包含
tenant_id    BIGINT NOT NULL,  -- 租户隔离
created_at   TIMESTAMP DEFAULT NOW(),
updated_at   TIMESTAMP DEFAULT NOW(),
deleted_at   TIMESTAMP NULL,   -- 软删除
```

### 5.2 核心表清单

```
tenants              # 租户
warehouses           # 仓库
clients              # 平台客户
client_members       # 会员
declarants           # 申报人
addresses            # 地址簿
parcels              # 包裹
parcel_logs          # 包裹日志
orders               # 集运订单
order_parcels        # 订单-包裹关联
routes               # 运输路线
couriers             # 快递公司
containers           # 集装柜
container_parcels    # 柜-包裹关联
service_templates    # 附加服务模板
service_orders       # 附加服务工单
ledger_entries       # 账本流水
monthly_statements   # 月结账单
client_recharges     # 充值记录
notifications        # 通知
audit_logs           # 审计日志
```

---

## 六、项目仓库规划

```
i56-framework/       # Core 框架 (独立仓库)
  ├── core/          # 框架核心
  ├── pkg/           # 公共工具
  └── sdk/           # Go SDK

i56-admin/           # 通用后台 (独立仓库)
  ├── templates/     # 管理后台模板
  └── static/        # 静态资源

i56-apps/            # 业务应用 (monorepo)
  ├── i56-wms/
  ├── i56-oms/
  ├── i56-tms/
  ├── i56-crm/
  └── i56-finance/
```

---

## 七、LTS 支持策略

| 版本 | 支持周期 | 内容 |
|------|---------|------|
| 1.0 LTS | 3 年 | 首个长期支持版 |
| 1.5 | 12 月 | 功能增强版 |
| 2.0 LTS | 3-5 年 | 下一代架构 |

每个 LTS 版本保证：
- 数据库迁移工具向前兼容
- API 兼容性策略
- 安全补丁定期推送
- 企业客户升级手册

---

## 八、当前 WMS 实现现状

### 8.1 已完成 (v1-v104)

| 模块 | 状态 | 页面数 |
|------|------|-------|
| 管理后台 | ✅ 完成 | 76 页全中文 |
| 认证 | ✅ | JWT + 员工编号 |
| RBAC | ✅ | 角色 CRUD 5 角色 |
| 订单管理 | ✅ | 9 订单, 状态流 |
| 包裹管理 | ✅ | 12 包裹, 状态 |
| 报表 | ✅ | 4 KPI 报表 |
| 仪表盘 | ✅ | 8 KPI 卡片 |
| 运营看板 | ✅ | 仓库/入库/作业台 |
| CRM | ✅ | 客户/会员/申报人/地址 |
| 物流 | ✅ | 承运商/快递/路线 |
| 系统 | ✅ | 角色/通知/定时/审计 |
| 种子数据 | ✅ | 35+ Store |

### 8.2 待完善

| 模块 | 优先级 |
|------|--------|
| 客户端门户 | P0 |
| PDA 功能 | P1 |
| 工作流引擎 | P1 |
| 消息队列 | P2 |
| 多租户隔离 | P2 |
| 国际化 (i18n) | P3 |

---

*文档版本: 1.0.0-draft | 基于 BFT56 完整审计 | 最后更新: 2026-07-15*
