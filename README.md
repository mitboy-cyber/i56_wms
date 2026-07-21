# I56 WMS — 跨境物流仓储管理系统

> **版本**: v3.0 | **语言**: Go 1.24+ / TypeScript | **架构**: Modular Monolith  
> **生产环境**: https://wms.mikaplay.com | **竞品对标**: 八方云仓 (BFT56)  
> **最新提交**: d224c77 | **GitHub**: github.com/mitboy-cyber/i56_wms

---

## 🚀 项目简介

I56 WMS 是一套面向跨境物流的全功能 SaaS 平台，覆盖订单管理、仓储作业、包裹追踪、财务报表、客户余额管理、物流运输调度、Webhook 投递等全链路业务。

系统采用 Go 模块化单体架构（可演进至微服务），前端基于 React + TypeScript + Vite + Tailwind CSS，支持 PDA 扫码作业、多语言切换、全局搜索、CSV 导出等企业级特性。

### 对标竞品 (BFT56)

BFT56（八方云仓）是同一代码基的生产级部署，运行于真实跨境物流业务。I56 WMS 已实现完全对标，并在以下维度具备独特优势：

| 维度 | BFT56 | I56 WMS v3.0 |
|------|:--:|:--:|
| 功能覆盖 | 56子页面 | **56子页面全对标** ✅ |
| Dashboard 可视化 | Chart.js | **SVG 饼图+柱状图** ✅ |
| 全局搜索 | 有 | **Enter 一键导航** ✅ |
| 多语言 | ZH/EN | **🇨🇳/🇺🇸 实时切换** ✅ |
| CSV 导出 | 有 | **BOM UTF-8 全列表** ✅ |
| 余额充值 | 有 | **充值弹窗+实时余额卡片** ✅ |
| Webhook | 有 | **投递监控+重试** ✅ |
| API 校验 | 有 | **字段级中文 422 错误** ✅ |
| 文档 | 无 | **PRD.md + README.md** ✅ |

---

## 📋 技术架构

```
Presentation       React SPA (Vite + TypeScript + Tailwind)  +  PDA (HTML5)
     │
Application        Admin API / Client API / PDA API (Go net/http)
     │
Domain             modules/order, modules/parcel, modules/customer,
                   modules/wms, modules/transport, modules/finance
     │
Infrastructure     MySQL, Redis, RabbitMQ, MinIO (已规划)
     │
Framework Core     router, eventbus, scheduler, logger, auth, tenant
```

### 认证方案
- **Admin**: Cookie (`admin_session`) + Bearer Token 双通道，Bearer 优先
- **Client**: Cookie (`client_session`)  
- **PDA**: Cookie (`pda_session`)

### RBAC 角色
| 角色 | 权限范围 |
|------|---------|
| 系统管理员 | 全部功能 |
| 仓库管理员 | 入库/出库/上架/包裹管理 |
| 客服人员 | 查看订单与包裹，处理咨询 |
| 财务人员 | 财务报表与对账 |
| 操作员 | PDA 扫码操作 |

---

## 📦 模块详情 (56 子页面)

### 1. 首页

| 页面 | 路由 | 功能 |
|------|------|------|
| 仪表盘 | `/admin` | KPI 卡片 + SVG 订单状态饼图 + 线路营收柱状图 |
| 仓库看板 | `/admin/warehouse-board` | 实时仓库状态（待拣货/待打包/装箱中） |

### 2. 订单管理

| 页面 | 路由 | 关键字段 |
|------|------|---------|
| 集运订单 | `/admin/orders` | 订单编号, 收件人, 路线, 包裹数, 总价, 实重, 计费重, 报关单号, 快递跟踪号, 状态, 备注, 创建时间 |
| 附加服务订单 | `/admin/service-orders` | 订单号, 客户, 服务模板, 状态, 负责人, 价格, 时间 |

#### 订单状态机
```
待拣货 → 拣货中 → 待打包 → 待装车 → 已装车 → 运输中
    → 清关中 → 已发货 → 已送达 → 已完成
任意状态 → 已取消
```

### 3. 仓库管理

| 页面 | 路由 | 关键字段 |
|------|------|---------|
| 包裹列表 | `/admin/parcels` | 快递单号, 品名, 包裹名, 状态, 货物类型, 实重, 快递公司, 入库时间 |
| 附加服务工单 | `/admin/service-orders` | 工单号, 订单, 客户, 服务类型, 状态, 价格 |
| 附加服务模板 | `/admin/service-templates` | 名称, 描述, 价格, 包含服务类型 |
| 附加服务类型 | `/admin/service-types` | 名称, 描述, 单价 |
| PDA 在线会话 | `/admin/pda-sessions` | 仓库, 操作员, 设备编号, 登录时间, 心跳, 当前页面/区域/货位, 在线状态 |
| 集装柜管理 | `/admin/containers` | 柜号, 仓库, 线路, 状态, 限重, 创建时间 |
| 仓库列表 | `/admin/warehouses` | 名称, 编码, 地址, 联系人, 电话, 包裹数 |
| 入库看板 | `/admin/inbound-board` | 快递单号, 品名, 重量, 快递, 状态(预申报/签收/称重/上架), 时间 |
| 仓库作业台 | `/admin/warehouse-console` | 待收货, 待称重, 待上架, 待拣货, 待打包, 待装车, 完成, 异常 |
| 员工任务监控 | `/admin/work-orders` | 入库任务, 拣货任务, 打包任务, 装车任务, 异常处理 |
| PDA 工单模板 | `/admin/pda-templates` | 模板名称, 步骤, 适用仓库 |
| 工单流程管理 | `/admin/workflow-management` | 流程名称, 步骤数, 所属模块, 状态, 更新时间 |
| 工单列表 | `/admin/orders` (管理员视图) | 所有待处理工单 |
| 异常记录 | `/admin/exception-reports` | 异常类型, 包裹, 描述, 处理状态, 时间 |

### 4. 财务报表

| 页面 | 路由 | 关键字段 |
|------|------|---------|
| 集运订单盈利 | `/admin/report/order-profit` | 订单数/营收/均单价, 明细表(订单号/客户/收入/成本/毛利) |
| 附加服务盈利 | `/admin/report/service-profit` | 服务单数/营收, 按服务类型聚合 |
| 客户盈利 | `/admin/report/client-profit` | 活跃客户数, 销售总额, 成本总额, 毛利, 毛利率 |
| 路线盈利 | `/admin/report/route-profit` | 按路线聚合的订单数/营收/利润 |

### 5. 物流管理

| 页面 | 路由 | 关键字段 |
|------|------|---------|
| 区域组管理 | `/admin/area-groups` | 名称, 编码, 描述 |
| 货物类型 | `/admin/cargo-types` | 名称, 编码 (普货/特货/海快普货/空运特货/食品/易碎品/液体/化妆品...) |
| 承运商列表 | `/admin/shipping-providers` | 名称, 编码, 联系人 |
| 快递公司 | `/admin/couriers` | 名称, 编码 |
| 清关公司 | `/admin/customs-brokers` | 名称, 编码, 联系人, 电话 |
| 清关点管理 | `/admin/customs-points` | 名称, 海关代码, 城市, 国家 |
| 线路模板 | `/admin/route-templates` | 名称, 出发地, 目的地, 承运商, 预计天数 |
| 运输公司 | `/admin/shipping-providers` | (同承运商列表) |
| 运输方式 | `/admin/transport-modes` | 名称, 编码 (海运/空运/海快/空运特货/商业海快) |
| **物流追踪** | `/admin/logistics-tracking` | **双栏面板 + 时间线 + 失败次数标记 + 承运商/出发地/目的地详情 + 刷新/复制** |

### 6. 客户管理

| 页面 | 路由 | 关键字段 |
|------|------|---------|
| 客户收件地址 | `/admin/customer-addresses` | 收件人, 电话, 城市, 区域, 详细地址, 是否默认 |
| 客户申报人 | `/admin/customer-declarants` | 类型(个人/公司), 姓名, 身份证, 公司统编, 电话, 归属会员, 认证状态, 启用状态 |
| 客户管理 | `/admin/clients` | 名称, 编码, 类型(平台/直客/代理), 联系人, 电话 |
| 客户账号 | `/admin/client-accounts` | 用户名, 所属客户, 角色, 状态 |
| 客户会员 | `/admin/client-members` | 姓名, 手机, 邮箱, 会员编码, 所属客户 |
| 客户充值 | `/admin/client-recharges` | 客户, 金额, 方式, 参考号, 状态, 时间 |
| **余额日志** | `/admin/client-ledgers` | **渐变余额卡片 + 充值弹窗 + 实时余额 + 交易明细表(充值/扣费标签)** |
| 充值记录 | `/admin/recharge-records` | (同客户充值) |
| 客户价格 | `/admin/pricing/routes` | 线路, 类型, 首重价, 续重价, 有效期 |
| **月结对账单** | `/admin/monthly-statements` | 客户, 年月, 周期, 订单数, 总金额, 总成本, 利润, 状态(待结算/已结算) |
| 客户端权限 | `/admin/client-panel-perms` | 客户端菜单权限配置 |

### 7. 系统

| 页面 | 路由 | 关键字段 |
|------|------|---------|
| **通知管理** | `/admin/notifications` | KPI 卡片 + 类型/优先级/范围/渠道/发送人/已发送/发送时间 完整表格 |
| 打印模板 | `/admin/print-templates` | 模板名称, 类型(标签/面单/发票), 内容 |
| 角色管理 | `/admin/roles` | 角色名称, 描述, 启用状态 |
| 员工管理 | `/admin/employees` | 仓库归属, 姓名, 账号, 角色, 电话, 邮箱, 创建时间 |
| 系统参数 | `/admin/system/params` | 键名, 值, 分组(system/finance/logistics/order/warehouse/customs), 标签 |
| **Webhook 投递** | `/admin/system/webhooks` | **配置 Webhook 端点 + 投递日志(客户/事件/URL/状态/响应码/时间)** |

---

## 🔧 API 架构

所有 API 使用 RESTful 风格，Content-Type: application/json。

### 关键端点

```
GET    /admin/api/dashboard/stats          — 仪表盘 KPI
GET    /admin/api/dashboard/order-status   — 订单状态分布 (v193+)
GET    /admin/api/dashboard/revenue-by-route — 线路营收 (v193+)
GET    /admin/api/orders                   — 订单列表
POST   /admin/api/orders                   — 创建订单 (含必填校验)
PUT    /admin/api/orders/{id}/status       — 订单状态流转
GET    /admin/api/client-ledgers?client_id={id} — 余额流水 (v195+)
POST   /admin/api/ledger-recharge          — 客户充值 (v195+)
GET    /admin/api/webhooks                 — Webhook 投递日志 (v197+)
GET    /admin/api/notifications            — 通知列表
GET    /admin/api/logistics-tracking       — 物流追踪
GET    /admin/api/area-groups              — 区域组列表
...
```

### API 响应格式 (v188+)
```json
// 成功
{"data": [...]}

// 校验失败 (422)
{"fields": {"recipient_name": "为必填项"}}

// 内部错误 (500) — 不暴露 DB 细节
{"error": "服务器内部错误，请稍后重试"}
```

---

## 🏗️ 架构亮点

### MinimalListPage — 零依赖列表组件
全部 64 个列表页面使用同一个纯 TypeScript + 内联样式的组件，**零外部图标库依赖**：
- ✅ 搜索 + 分页 + CRUD 弹窗  
- ✅ CSV 导出 (BOM UTF-8)  
- ✅ 自定义 render 列  
- ✅ 加载/错误/空状态全处理

### SVG 仪表盘图表 (v193+)
纯 SVG 实现，无 Chart.js 依赖：
- Donut Chart — 订单状态分布  
- Bar Chart — 线路营收 TOP

### 移动端适配 (v199+)
- 📱 汉堡菜单（侧边栏滑入/滑出）  
- 📊 表格横向滚动  
- 📈 图表竖排堆叠  
- 👆 触摸友好输入框（16px 防 iOS 缩放）  
- 🎨 CSS 媒体查询断点：移动端 ≤768px / 平板 ≤1024px

### 全局搜索 (v194+)
顶部搜索栏 → 输入关键词 → Enter 一键导航到对应页面

### i18n 国际化 (v197+)
`🇨🇳 中文 / 🇺🇸 EN` 切换按钮，支持 localStorage 持久化

---

## 📈 路线图

### ✅ v3.0（已完成）
- [x] Dashboard SVG 图表
- [x] 全局搜索
- [x] 语言切换 (zh/en)
- [x] CSV 导出
- [x] 客户充值 + 余额卡片
- [x] 物流追踪时间线
- [x] Webhook 投递监控
- [x] 月结对账单
- [x] 通知管理增强
- [x] API 字段级校验 (httputil+validate)
- [x] 安全错误处理
- [x] 64页白屏修复
- [x] PRD.md + README.md
- [x] 移动端响应式适配 🆕

### 🔜 v3.1（规划中）
- [ ] PDA WebSocket 实时推送
- [ ] 多租户隔离 (Schema/DB per Tenant)
- [ ] PDF 打印面单
- [ ] 订单导入 Excel

---

## 📊 版本链

```
57044fc → v199: mobile responsive (hamburger + CSS media queries)
6d49ce1 → v198: README.md v3.0
d224c77 → v197: v3.0-final (webhook + i18n + recharge UI)
827c580 → v195: recharge API + seed ledgers
bf6acc9 → v194: global search + tracking timeline
679fb9c → v193: SVG charts + CSV export
fc24e52 → v192: PRD.md (1023 lines)
c18611f → v191: README.md
1c2e582 → v190: queryFn fix
bc79f3b → v189: 64-page MinimalListPage
3de0292 → v188: validation framework
```
### 分支策略

## 📝 贡献

```bash
git clone https://github.com/mitboy-cyber/i56_wms.git
cd i56_wms/apps/wms

# 后端
go build -o i56-server ./cmd/server/ && ./i56-server

# 前端
cd frontend && npm install && npm run dev
```
- `main` — 生产分支
- `develop` — 开发分支
- `feature/*` — 功能分支
- 每次提交格式: `v<版本号>: <描述>`

---

## 📊 版本链

```
d224c77 → v197: v3.0-final (webhook + i18n + enhanced client ledgers + notifications)
827c580 → v195: client recharge API + 4-seed ledgers
bf6acc9 → v194: global search + logistics tracking timeline
679fb9c → v193: dashboard SVG charts + CSV export
fc24e52 → v192: PRD.md (1023 lines)
c18611f → v191: README.md (230 lines)
1c2e582 → v190: MinimalListPage queryFn fix
bc79f3b → v189: 64-page GenericListPage → MinimalListPage
3de0292 → v188: backend validation framework
```

---

**维护方**: I56 Team | **许可**: MIT | **最后更新**: 2026-07-21
