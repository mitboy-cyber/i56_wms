# I56 WMS — 跨境物流仓储管理系统

> **版本**: 2.4.2 | **语言**: Go 1.24+ / TypeScript | **架构**: Modular Monolith  
> **生产环境**: https://wms.mikaplay.com | **竞品对标**: 八方云仓 (BFT56)

---

## 项目简介

I56 WMS 是一套面向跨境物流的服务平台，覆盖订单管理、仓储管理、包裹追踪、财务报表、客户管理、物流运输调度等全链路业务。系统采用 Go 模块化单体架构，前端基于 React + TypeScript + Vite + Tailwind CSS，支持 PDA 扫码作业、客户端多角色登录。

### 对标竞品 (BFT56)

BFT56（八方云仓）是同代码基的生产级部署，运行于真实跨境物流业务。BFT56 目前已覆盖：
- 订单管理层 → 集运订单 + 附加服务订单
- 仓库作业层 → 入库看板 / 仓库作业台 / PDA 协同
- 物流管理层 → 区域组 / 线路模板 / 承运商 / 运输方式 / 物流追踪
- 客户管理层 → 客户账号 / 会员 / 充值 / 余额 / 价格 / 月结对账
- 财务层 → 集运订单盈利 / 附加服务盈利 / 客户盈利 / 路线盈利
- 系统层 → 通知 / 打印模板 / RBAC / 员工 / 系统参数

---

## 技术架构

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

## 模块详情

### 1. 首页

| 页面 | 路由 | 功能 |
|------|------|------|
| 仪表盘 | `/admin/dashboard` | KPI 卡片（订单/包裹/营收统计） |
| 仓库看板 | `/admin/warehouse-board` | 实时仓库状态（待拣货/待打包/装箱中） |

### 2. 订单管理

| 页面 | 路由 | 关键字段 |
|------|------|---------|
| 集运订单 | `/admin/orders` | 订单编号, 收件人, 路线, 包裹数, 总价, 实重, 计费重, 报关单号, 快递跟踪号, 状态, 备注, 创建时间 |
| 附加服务订单 | `/admin/service-orders` | 订单号, 客户, 服务模板, 状态, 负责人, 价格, 时间 |

#### 订单状态机
```
待拣货(pending_picking) → 拣货中(picking) → 待打包(pending_packing)
→ 待装车(pending_loading) → 已装车(loaded) → 运输中(in_transit)
→ 清关中(customs_clearance) → 已发货(shipped) → 已送达(delivered) → 已完成(completed)
任意状态 → 已取消(cancelled)
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

#### 包裹状态机
```
预申报(pre_declared) → 已签收(received) → 已称重(weighed)
→ 已上架(stored) → 已打包(packed) → 已发货(shipped) → 已送达(delivered)
→ 已签收(signed)
```

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
| 物流追踪 | `/admin/logistics-tracking` | 单号, 线路, 状态, 段落, 承运商, 详细信息, 失败次数, 关联订单 |

### 6. 客户管理

| 页面 | 路由 | 关键字段 |
|------|------|---------|
| 客户收件地址 | `/admin/customer-addresses` | 收件人, 电话, 城市, 区域, 详细地址, 是否默认 |
| 客户申报人 | `/admin/customer-declarants` | 类型(个人/公司), 姓名, 身份证, 公司统编, 电话, 归属会员, 认证状态, 启用状态 |
| 客户管理 | `/admin/clients` | 名称, 编码, 类型(平台/直客/代理), 联系人, 电话 |
| 客户账号 | `/admin/client-accounts` | 用户名, 所属客户, 角色, 状态 |
| 客户会员 | `/admin/client-members` | 姓名, 手机, 邮箱, 会员编码, 所属客户 |
| 客户充值 | `/admin/client-recharges` | 客户, 金额, 方式, 参考号, 状态, 时间 |
| 余额日志 | `/admin/client-ledgers` | 客户, 类型(充值/扣款/退款), 描述, 金额, 余额, 时间 |
| 充值记录 | `/admin/recharge-records` | (同客户充值) |
| 客户价格 | `/admin/pricing/routes` | 线路, 类型, 首重价, 续重价, 有效期 |
| 月结对账单 | `/admin/monthly-statements` | 客户, 年月, 周期, 订单数, 总金额, 总成本, 利润, 状态 |
| 客户端权限 | `/admin/client-panel-perms` | 客户端菜单权限配置 |

### 7. 系统

| 页面 | 路由 | 关键字段 |
|------|------|---------|
| 通知管理 | `/admin/notifications` | 标题, 类型(系统通知/公告/任务通知), 优先级, 发送范围, 内容, 渠道(站内信/邮件/短信/微信), 发送人, 已发送, 发送时间 |
| 打印模板 | `/admin/print-templates` | 模板名称, 类型(标签/面单/发票), 内容 |
| 角色管理 | `/admin/roles` | 角色名称, 描述, 启用状态 |
| 员工管理 | `/admin/employees` | 仓库归属, 姓名, 账号, 角色, 电话, 邮箱, 创建时间 |
| 系统参数 | `/admin/system/params` | 键名, 值, 分组(system/finance/logistics/order/warehouse/customs), 标签 |

---

## API 架构

所有 API 使用 RESTful 风格，Content-Type: application/json。

### Admin API (Bearer Token 认证)
```
GET    /admin/api/dashboard/stats       — 仪表盘 KPI
GET    /admin/api/orders                — 订单列表
POST   /admin/api/orders                — 创建订单 (含必填校验)
PUT    /admin/api/orders/{id}/status    — 订单状态流转
GET    /admin/api/parcels               — 包裹列表
POST   /admin/api/parcels               — 包裹预申报
GET    /admin/api/area-groups           — 区域组
GET    /admin/api/cargo-types           — 货物类型
GET    /admin/api/transport-modes       — 运输方式
...
```

### API 响应格式 (v188+)
```json
// 成功
{"success": true, "data": [...]}
// 校验失败 (422)
{"success": false, "error": "数据校验失败", "fields": {"recipient_name": "为必填项"}}
// 内部错误 (500) — 不暴露 DB 细节
{"success": false, "error": "服务器内部错误，请稍后重试"}
```

---

## 竞品对标差距

| 差距项 | BFT56 | I56 WMS | 优先级 |
|--------|-------|---------|:--:|
| 仪表盘图表 | Chart.js 可视化 | KPI 数字卡片 | P1 |
| 订单导出 Excel | 内置 | 无 | P1 |
| 全局搜索 | 顶部搜索栏 | 无 | P2 |
| 多语言 (ZH/EN) | 语言切换 | 仅中文 | P2 |
| 主题密度 (紧凑/标准/大号) | 三档密度 | 无 | P2 |
| 打印模板管理 | 完整 | 基础 | P2 |
| 物流追踪详情 | 物流段落+失败次数 | 占位 | P1 |
| 客户充值流程 | 充值→余额日志→记录 | 种子数据 | P1 |
| 月结对账单 | 完整 | 占位 | P1 |
| PDA 实时协同 | WebSocket | 轮询 | P1 |
| Webhook 投递 | 客户端可配 | 无 | P2 |

---

## 路线图

### v2.5 (当前迭代)
- [x] Admin 全页面白屏修复
- [x] 后端 API 验证框架 (httputil + validate)
- [x] 64列表页 GenericListPage → MinimalListPage 替换
- [ ] 仪表盘 Chart.js 可视化
- [ ] 订单导出 Excel/CSV
- [ ] 物流追踪详情页补全

### v2.6
- [ ] 客户端充值与余额流水（端到端）
- [ ] 月结对账单生成逻辑
- [ ] PDA WebSocket 实时推送
- [ ] 全局搜索实现

### v3.0
- [ ] 多租户隔离 (Schema/DB per Tenant)
- [ ] Webhook 投递通知
- [ ] 多语言国际化
- [ ] 移动端适配

---

## 贡献

```bash
git clone https://github.com/mitboy-cyber/i56_wms.git
cd i56_wms/apps/wms

# 后端
go build -o i56-server ./cmd/server/ && ./i56-server

# 前端
cd frontend && npm install && npm run dev
```

### 分支策略
- `main` — 生产分支
- `develop` — 开发分支
- `feature/*` — 功能分支
- 每次提交格式: `v<版本号>: <描述>`

---

**维护方**: I56 Team | **许可**: MIT | **最后更新**: 2026-07-21
