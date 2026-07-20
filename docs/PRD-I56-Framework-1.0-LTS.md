# I56 Framework 1.0 LTS — 产品需求文档
## Enterprise Application Development Platform

---

### 文档信息

| 项 | 值 |
|---|-----|
| 版本 | 1.0 LTS |
| 日期 | 2026-07-20 |
| 状态 | v165 已部署 8/8 Framework 模块 |
| 对标竞品 | BFT56 (八方云仓) — 全栈 SaaS OMS+WMS+TMS+CRM+Finance |
| 仓库 | github.com/mitboy-cyber/i56_wms |

---

## 一、产品愿景

构建一套可支撑未来 10 年企业应用开发的统一平台：

```
一套 Framework → 多种业务模块 → 多个行业产品 → SaaS 租户 → 私有化/云原生部署
```

---

## 二、完成状态

### 2.1 Framework 核心模块 (8/8 ✅)

| 模块 | 版本 | 能力 |
|------|:--:|------|
| **EventBus** | v157 | Pub/Sub + 通配符 + 域事件驱动 |
| **Multi-Tenant** | v158 | 3 策略 (Shared/Schema/Database) |
| **RBAC** | v159 | 3 角色 + 数据范围 (all/warehouse/tenant/self) |
| **Go 1.26** | v160 | 统一构建链 |
| **Storage** | v161 | Local + MinIO/S3 抽象 |
| **Workflow** | v162 | 采购审批流 (5 状态) + 条件路由 |
| **Notification** | v163 | 多通道 (in_app/email/sms/slack) |
| **Plugin** | v164 | 统一服务定位器 + 生命周期管理 |

### 2.2 安全审计修复 (v165)

| 等级 | 修复项 |
|:--:|------|
| 🔴 P0 | DB 关闭移到 shutdown / config.Load 错误检查 / Storage fail-fast |
| 🟠 P1 | EventLog 读写锁 / Token 生成错误 / 未认证端点 |
| 🟡 P2 | JSON 注入 (保留低风险) |

### 2.3 已知架构 Gap (BFT56 对比)

| 领域 | 当前 | 目标 | 差距 |
|------|:--:|:--:|:--:|
| Multi-Tenant | 20% (中间件存在,数据层无) | 100% | **80%** |
| DataScope | 20% (API 存在,查询层无) | 100% | **80%** |
| EventBus | 20% (域事件未自动发布) | 100% | **80%** |
| Finance | 0% | 盈利报表/充值/对账 | **100%** |
| TMS | 10% | 承运商/线路/清关/追踪 | **90%** |

---

## 三、技术架构

### 3.1 技术栈

```
语言: Go 1.26+
前端: React 18 + TypeScript + Vite + Tailwind CSS v4 + TanStack Query
测试: go test + Playwright (待)
部署: Docker Compose / Kubernetes (待)
存储: PostgreSQL + Redis + MinIO + Elasticsearch (待)
```

### 3.2 目录结构

```
i56-framework/
├── cmd/server/    → 入口
├── core/          → 8 个 Framework 模块 (独立可版本化)
│   ├── eventbus/
│   ├── tenant/
│   ├── rbac/
│   ├── storage/
│   ├── workflow/
│   ├── notification/
│   ├── plugin/
│   └── router/middleware/cache/logger/config/validator/
├── apps/wms/      → 业务应用
│   ├── internal/server/       → Server.go (1048行)
│   ├── internal/adminapi/     → 6 个模块 API
│   ├── internal/domain/       → 领域模型
│   └── frontend/              → React SPA
├── sdk/           → 客户端 SDK
└── deployments/   → Docker + Helm
```

### 3.3 架构模式

```
Modular Monolith (可演进为微服务)
7 层: Presentation → Application → Domain → Infrastructure → Core → Storage → OS

所有 Framework 模块通过 plugin.Registry 统一管理
业务模块通过 EventBus 解耦，无直接依赖
```

---

## 四、业务模块规划

### 4.1 当前 (I56 WMS v2.4.2)

```
OMS   → 集运订单、订单 CRUD
WMS   → 包裹列表、仓库、货架、入库
TMS   → 承运商、路线、路由价格
CRM   → 客户、钱包、账本、申报人
System → 角色、用户、通知、插件
```

### 4.2 对标 BFT56 缺失模块

| 模块 | BFT56 功能 | I56 优先级 |
|------|----------|:--:|
| **Finance** | 集运盈利/服务盈利/客户盈利/路线盈利/充值/月结对账 | 🔴 P0 |
| **PDA** | 在线会话/工单模板/任务监控/操作台 | 🟠 P1 |
| **BI** | 仓库看板/入库看板/盈利分析 | 🟡 P2 |
| **Logistics** | 清关公司/清关点/区域组/货物类型/运输方式 | 🟡 P2 |
| **CRM++** | 客户价格/客户端权限/会员体系 | 🟡 P2 |

### 4.3 计划产品

```
I56 WMS  (已完成)
I56 OMS  (已集成)
I56 TMS  (部分)
I56 CRM  (部分)
I56 Finance ← 下一阶段
I56 BI    ← 后续
I56 PDA   ← 后续
```

---

## 五、API 端点总表 (v165 生产)

| 分类 | 端点 | 方法 |
|------|------|:--:|
| Health | `/api/v1/health` | GET |
| Auth | `/admin/login`, `/client/login` | POST |
| Session | `/admin/api/me` | GET |
| EventBus | `/admin/api/events`, `/admin/api/events/publish` | GET/POST |
| Tenant | `/admin/api/tenants`, `/admin/api/tenant` | GET |
| RBAC | `/admin/api/rbac/subject`, `/check`, `/datascope` | GET |
| Storage | `/admin/api/storage/upload`, `/list` | POST/GET |
| Workflow | `/admin/api/workflow/start`, `/transition`, `/definitions` | POST/GET |
| Notification | `/admin/api/notify/send`, `/notify/inbox` | POST/GET |
| Plugin | `/admin/api/plugins`, `/plugins/resolve` | GET |
| Business | `/admin/api/orders`, `/parcels`, `/warehouses`, `/clients`, `/couriers`, `/employees` | CRUD |

---

## 六、LTS 路线图

```
v1.0 LTS (当前)     → 2026 Q3. 8 Framework + WMS
v1.5                → 2026 Q4. Finance 模块 + 域事件 + DataScope 渗透
v2.0 LTS            → 2027 Q1. Microservice 拆分 + K8s + PostgreSQL
v3.0 LTS            → 2028.    Marketplace + Plugin 生态
```

### v1.0 LTS 支撑承诺

- 数据库迁移工具
- API 兼容策略
- 安全补丁 (3 年)
- Go 1.26+ 长期兼容

---

## 七、竞品对标总结

| 维度 | BFT56 | I56 WMS | I56 PRD 目标 | 结果 |
|------|:--:|:--:|:--:|:--:|
| Framework 抽象 | ❌ 单体 | ✅ 8模块可复用 | ✅ | **赢** |
| 多租户 | ❌ 单例 | ✅ 2 Tenant | ✅ 3策略 | **赢** |
| RBAC | ⚠️ 固定角色 | ✅ Enforcer+通配符 | ✅ DataScope | **持平** |
| Workflow | ✅ 工单流程 | ✅ 采购审批 | ⚠️ 需条件引擎 | **持平** |
| Finance | ✅ 4维盈利 | ❌ | 🔴 缺失 | **输** |
| TMS | ✅ 完整物流 | ⚠️ 基础 | 🟠 待扩展 | **输** |
| PDA | ✅ 在线会话 | ❌ | 🟠 待开发 | **输** |
| 技术栈 | PHP Filament | Go+React | Go+React | **赢** |

---

## 八、BFT56 客户端端分析

### 8.1 客户端模块树 (EZ集运通 — plat_ezjyt)

```
主控台 (Dashboard)
├── 收件地址 (Delivery Addresses)
├── 客戶會員 (Customer Members)
├── 申報人 (Declarants)              ← 用户核心需求
├── 我的訂單 (My Orders)
├── 我的包裹 (My Parcels)
├── 附加服務訂單 (Additional Service)
├── 帳務 (Billing)
│   ├── 餘額明細 (Balance Details)
│   └── 月結對帳單 (Monthly Statements)
├── 價格 / 承運商 (Pricing/Carriers)
│   ├── 倉庫資訊 (Warehouse Info)
│   ├── 線路價格 (Route Pricing)
│   ├── 承運商派送價 (Carrier Delivery)
│   ├── 承運商加收價 (Carrier Surcharge)
│   └── Webhook 投遞 (Webhook Delivery)
└── API 凭证 (API Credentials)
```

### 8.2 订单字段分析 (真实数据)

| 字段 | 示例值 | 说明 |
|------|--------|------|
| 訂單號 | 20260720181715687146 | 时间戳+随机数字 |
| 倉庫 | 厦门仓 | 单一仓库 |
| 客戶會員 | 127154 | 会员ID |
| 收件人 | 蔡于萱 | 收件人姓名 |
| 快遞單號 | 773431715690557 | 快递追踪号 |
| 包裹數 | 1-2 | 多包裹支持 |
| 線路 | 海快 | 运输线路 |
| 承運商 | 新竹物流 | 台湾本地物流 |
| 總實重(kg) | 1.22 / 23.07 | 实际重量 |
| 總計費量(kg) | 1.22 / 23.07 | 计费量(体积重折算) |
| 金額 | ¥ 28.30 / ¥191.48 | 人民币 |
| 狀態 | 待揀貨 | 订单状态 (待拣货) |
| 打包時間 | — | 未打包 |
| 裝櫃時間 | — | 未装柜 |

### 8.3 I56 WMS 客户端对比

| 功能 | BFT56 | I56 WMS | 差距 |
|------|:--:|:--:|:--:|
| 订单列表 | ✅ 含状态/重量/价格 | ✅ | 持平 |
| 包裹管理 | ✅ | ✅ | 持平 |
| 会员管理 | ✅ | ✅ | 持平 |
| 申报人管理 | ✅ | ✅ | 持平 |
| 余额/对账 | ✅ 余额明细+月结 | ⚠️ 仅账本 | **输** |
| 线路价格 | ✅ 客户可见 | ❌ | **输** |
| Webhook | ✅ 客户端配置 | ❌ | **输** |
| API 凭证 | ✅ 自助管理 | ❌ | **输** |
| 订单状态流 | ✅ 待拣货→打包→装柜 | ⚠️ 静态 | **输** |

### 8.4 BFT56 订单状态流转 (推断)

```
待揀貨 (Pending Picking)
  → 已揀貨 (Picked)
  → 已打包 (Packed)        ← 打包時間
  → 已裝櫃 (Containerized) ← 裝櫃時間
  → 已出貨 (Shipped)
  → 已簽收 (Delivered)
```

---

*最后更新: 2026-07-20 — v165 生产正式版*
