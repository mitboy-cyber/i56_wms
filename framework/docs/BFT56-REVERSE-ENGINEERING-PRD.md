# BFT56 八方云仓 — 完整逆向工程 PRD

> 分析日期：2026-07-10  
> 目标系统：https://bft56.com  
> 技术栈：Laravel + Filament + Spatie Shield + Tailwind CSS + MySQL  
> 市场定位：跨境集运物流 SaaS（台湾/香港市场）  
> 语言：繁体中文（zh-TW）

---

## 一、产品概述

**八方云仓** 是一套面向跨境集运行业的 SaaS 平台，核心业务为：中国大陆仓库收货 → 集运合单 → 报关清关 → 台湾末端派送。系统由 **管理后台** 和 **客户端门户** 两大子系统组成，支持多租户（多客户/多仓库）。

### 核心业务流

```
用户网购 → 包裹寄往大陆仓库 → 仓库签收入库 → 客户创建集运订单
→ 仓库拣货打包 → 装柜 → 报关清关 → 承运商运输 → 台湾末端派送 → 签收
```

---

## 二、系统架构概览

### 2.1 URL 模块划分

| 前缀 | 系统 | 说明 |
|------|------|------|
| `/admin` | 管理后台 | Laravel Filament Admin Panel |
| `/admin/o-m-s/` | Order Management | 订单、包裹、申报人、客户核心数据 |
| `/admin/w-m-s/` | Warehouse Management | 仓库、工单、流程、库位 |
| `/admin/t-m-s/` | Transport Management | 承运商、线路、清关、物流 |
| `/admin/system/` | System | 用户、通知、打印、PDA、日志 |
| `/admin/shield/` | Shield RBAC | Spatie Permission 角色权限 |
| `/client` | 客户端门户 | 客户自助操作平台 |

### 2.2 技术检测

- **框架**：Laravel (PHP)
- **后台面板**：Filament (Livewire + Tailwind)
- **权限系统**：Spatie Shield (`/admin/shield/roles`)
- **前端渲染**：SSR + Livewire SPA
- **搜索**：全局搜索框（Filament Global Search）
- **国际化**：zh-TW（繁体中文）

---

## 三、管理后台完整菜单树

### 3.1 顶级导航

```
├── 首页 (/admin)
├── 仓库看板 (/admin/warehouse-board)
├── 订单管理
├── 仓库管理
├── 财务报表
├── 物流管理
├── 客户管理
└── 系统
```

### 3.2 订单管理 (OMS)

| 菜单项 | URL | 说明 |
|--------|-----|------|
| 集运订单 | `/admin/o-m-s/orders` | 核心订单管理，18 列数据 |
| 附加服务订单 | `/admin/o-m-s/parcel-service-orders` | 增值服务订单 |

**集运订单 表列（18列）**：
编号、仓库、订单号、客户名称、客户会员、收件人、快递单号、运输线路、包裹数、订单状态、所属柜、总实重(kg)、总计费量(kg)、总价、清关单号、承运商单号、备注、下单时间

**筛选器**：批量订单号、客户、仓库、线路、所属柜、订单状态、起止时间、客户会员编号

**订单操作**：打印、设置清关公司、设置承运商单号、流水、导出申报单

**订单状态**（观测到）：
- `待拣货` — 已生成，等待仓库拣选
- `待装柜` — 已拣货打包，等待装入集装箱
- `已取消` — 已取消

### 3.3 仓库管理 (WMS)

| 菜单项 | URL | 说明 |
|--------|-----|------|
| 包裹列表 | `/admin/o-m-s/parcels` | 20 列，核心包裹管理 |
| 附加服务工单 | `/admin/o-m-s/parcel-service-order-items` | 增值服务执行 |
| 附加服务模板 | `/admin/o-m-s/parcel-service-templates` | 服务模板定义 |
| 附加服务类型 | `/admin/parcel-services-type` | 服务类型分类 |
| PDA 在线会话 | `/admin/system/operator-sessions` | PDA 设备在线状态 |
| 仓库列表 | `/admin/w-m-s/warehouses` | 仓库配置 |
| 入库看板 | `/admin/inbound-board` | 入库可视化 |
| 仓库作业台 | `/admin/warehouse-console` | 仓库操作面板 |
| 员工任务监控 | `/admin/system/work-orders` | 任务分配与监控 |
| PDA 工单模板 | `/admin/w-m-s/work-order-templates` | PDA 任务模板 |
| 工单流程管理 | `/admin/w-m-s/workflow-processes` | 工作流设计 |
| 工单列表 | `/admin/w-m-s/work-order-lists` | 工单执行列表 |
| 异常记录 | `/admin/o-m-s/exception-reports` | 包裹异常管理 |

**包裹列表 表列（20列）**：
编号、仓库、包裹编号、快递公司、所属客户、客户会员、货物类型、集运订单、品名、运输方式、包裹名、实重(kg)、尺寸(cm)、状态、异常、附加服务、库位、最后操作、操作

**包裹操作**：手动认领、预报、打印、流转事件、标记异常、拒收

**包裹状态**（从权限推断完整流转）：
预报 → 已入库 → 已称重 → 已上架 → 已拣货 → 已打包 → 已出库

### 3.4 财务报表

| 菜单项 | URL | 说明 |
|--------|-----|------|
| 集运订单盈利 | `/admin/order-profit-report-page` | 订单维度利润分析 |
| 附加服务盈利 | `/admin/service-profit-report-page` | 增值服务利润 |
| 客户盈利 | `/admin/client-profit-report-page` | 客户维度利润汇总 |
| 路线盈利 | `/admin/route-profit-report-page` | 线路维度利润分析 |

### 3.5 物流管理 (TMS)

| 菜单项 | URL | 说明 |
|--------|-----|------|
| 区域组管理 | `/admin/t-m-s/area-groups` | 地理区域分组 |
| 货物类型 | `/admin/t-m-s/cargo-types` | 货物分类（普货/特货/敏感货） |
| 承运商列表 | `/admin/t-m-s/carriers` | 末端派送承运商 |
| 装柜记录 | `/admin/t-m-s/container-loading-records` | 集装箱装载记录 |
| 快递公司 | `/admin/t-m-s/couriers` | 国内快递公司 |
| 清关公司 | `/admin/t-m-s/customs-brokers` | 报关行管理 |
| 清关点管理 | `/admin/t-m-s/customs-clearance-points` | 清关口岸 |
| 线路模板 | `/admin/t-m-s/routes` | 运输线路与定价 |
| 运输公司 | `/admin/t-m-s/shipping-providers` | 干线运输商 |
| 运输方式 | `/admin/t-m-s/transport-types` | 空运/海运/海快等 |
| 物流追踪 | `/admin/t-m-s/logistics-trackings` | 轨迹追踪 |

**线路模板 表列**：编号、啟用(开关)、發貨倉庫、範本名稱、運輸方式、支援貨類

### 3.6 客户管理 (CRM)

| 菜单项 | URL | 说明 |
|--------|-----|------|
| 客户收件地址 | `/admin/o-m-s/client-member-addresses` | 会员收货地址 |
| 客户申报人 | `/admin/o-m-s/declarants` | 报关人信息 |
| 客户管理 | `/admin/o-m-s/clients` | 企业客户 |
| 客户账号 | `/admin/o-m-s/client-users` | 登录账号 |
| 客户会员 | `/admin/o-m-s/client-members` | 终端会员(收件人) |
| 客户充值 | `/admin/o-m-s/client-recharges` | 充值审核 |
| 余额日志 | `/admin/o-m-s/client-ledgers` | 账户流水 |
| 充值记录 | `/admin/o-m-s/client-recharge-logs` | 充值历史 |
| 客户价格 | `/admin/o-m-s/client-route-prices` | 客户专属报价 |
| 月结对账单 | `/admin/o-m-s/client-statements` | 月度结算单 |
| 客户端权限 | `/admin/client-panel-permissions` | 客户门户功能开关 |

### 3.7 系统设置

| 菜单项 | URL | 说明 |
|--------|-----|------|
| 通知管理 | `/admin/system/notifications` | 系统通知推送 |
| 打印模板 | `/admin/system/print-templates` | 面单/标签模板 |
| 角色管理 | `/admin/shield/roles` | RBAC 角色配置 |
| 员工管理 | `/admin/system/users` | 后台操作员 |
| 任务派发参数 | `/admin/task-dispatch-settings-page` | 自动派单规则 |

---

## 四、隐藏资源注册表（来自权限页面）

以下模块不在主导航菜单中直接展示，但通过权限系统暴露，说明系统已实现但入口在其他地方或仅超管可见：

| 资源 | 权限操作 | 业务含义 |
|------|---------|----------|
| 客户派送费 | CRUD + 从默认同步 | Customer Delivery Fee |
| 承运商授权 | CRUD + 批量授权 + 一键绑定 | Carrier Authorization |
| 客户加收费 | CRUD + 从默认同步 | Customer Surcharge |
| 流水 | CRUD | Transaction Ledger |
| 仓库授权 | CRUD + 批量授权 | Warehouse Access Authorization |
| 客户仓储价 | CRUD + 从默认同步 | Customer Storage Pricing |
| 日志 | CRUD | System Operation Logs |
| Api 调用日志 | CRUD | API Call Audit Logs |
| Pda 版本 | CRUD + 下载APK | PDA App Version Management |
| 承运商单号 | CRUD + 作废 + 批量导入 | Carrier Tracking Number Pool |
| 清关单号 | CRUD + 作废 + 批量生成 + 批量导入 | Customs Declaration Number Pool |
| 集装柜 | CRUD | Container Management |
| 入库机 | CRUD + 重置Token | Inbound Machine Registration |
| 库位 | CRUD + 打印 | Storage Location |
| 库位类型 | CRUD | Location Type Dictionary |
| 区域 | CRUD + 打印 | Warehouse Zone |
| 区域类型 | CRUD | Zone Type Dictionary |
| 权限 | CRUD | Shield Permission Management |
| 工单 | CRUD | Work Order Instances |
| 工单模板 | CRUD + 管理 + 复制 | Work Order Templates |
| 流程实例 | CRUD + 取消 + 重启 | Process Instances |
| 工单流程 | CRUD + 管理 + 复制 | Process Definitions |

### 特殊权限（功能开关型）

| 权限名 | 含义 |
|--------|------|
| 客户端权限 | 控制客户门户功能可见性 |
| 仓库看板 | 访问仓库数据看板 |
| 客户盈利汇总 | 访问客户维度的利润报表 |
| 集运订单盈利报表 | 访问订单利润报表 |
| 路线盈利汇总 | 访问路线利润报表 |
| 附加服务盈利报表 | 访问增值服务利润报表 |
| 入库看板 | 访问入库可视化面板 |
| 任务派发参数 | 配置自动任务派发规则 |
| 仓库作业台 | 访问仓库操作面板 |
| Global Stats Widget | 首页全局统计组件 |
| Parcel Status Stats Widget | 包裹状态统计图 |
| Online Employees Board Widget | 在线员工面板 |
| Parcel Order Flow Widget | 包裹订单流程图 |
| Open Pool Widget | 开放单号池组件 |
| Alerts Widget | 预警通知组件 |
| Customs Number Pool Widget | 清关单号池组件 |
| Location Stats Widget | 库位统计组件 |
| Carrier Tracking Number Pool Widget | 承运商单号池组件 |
| 驾驶舱管理 | Cockpit/Dashboard Admin |
| PDA 打印(面单/标签) | PDA 打印功能 |
| Webhook 投递日志 重发 | Webhook 重发（超管专用） |
| 查看全部仓库 | 跨仓数据查看（否则锁本仓） |
| 修改/删除已封柜或已发运的集装柜 | 高危操作权限 |
| 字典管理 增删改 | 字典表（库位类型/区域类型）管理 |
| 通知 查看本公司全部(跨仓)+ 发全公司广播 | 通知管理高级权限 |
| 员工管理 查看本公司全部账号 | 员工列表跨仓查看 |
| 查看订单财务(成本/盈利/盈利率) | 订单财务敏感数据 |

---

## 五、客户端门户 (Client Portal)

### 5.1 完整菜单

```
客户端门户 (/client)
├── 主控台 (Dashboard)
├── 收件地址 (Addresses)
├── 客戶會員 (Members)
├── 申報人 (Declarants)
├── 我的訂單 (My Orders)
├── 我的包裹 (My Parcels)
├── 附加服務訂單 (Service Orders)
├── 帳務 (Billing)
│   ├── 餘額明細 (Balance Ledger)
│   └── 月結對帳單 (Monthly Statement)
├── 價格 / 承運商 (Pricing)
│   ├── 倉庫資訊 (Warehouse Info)
│   ├── 線路價格 (Route Pricing)
│   ├── 承運商派送價 (Carrier Delivery Fee)
│   ├── 承運商加收價 (Carrier Surcharge)
│   └── Webhook 投遞 (Webhook Delivery Logs)
└── API 凭证 (API Credentials)
```

### 5.2 客户端订单列表

**表列（13列）**：訂單號、倉庫、客戶會員、收件人、快遞單號、包裹數、線路、承運商、總實重(kg)、總計費量(kg)、金額、狀態、打包時間

**操作**：檢視、下附加服務單、取消

### 5.3 申报人列表

**表列（9列）**：類型(个人/公司)、姓名/公司名稱、身分證、公司統編、電話、歸屬客戶會員、狀態(认证状态)、啟用狀態

**操作**：檢視、同步認證、停用

---

## 六、核心领域模型

### 6.1 实体关系图（ER 概要）

```
Tenant (租户/客户公司)
 ├── Client (客户)
 │    ├── ClientAccount (客户账号 - 登录用)
 │    │    └── WarehouseAuthorization (仓库授权)
 │    ├── ClientMember (客户会员 - 终端收件人)
 │    │    ├── MemberAddress (收件地址)
 │    │    └── Declarant (申报人/报关人)
 │    ├── ClientLedger (余额流水)
 │    ├── ClientRecharge (充值申请)
 │    ├── ClientRoutePrice (客户线路价格)
 │    ├── ClientServiceOverride (附加服务覆盖)
 │    ├── ClientStoragePrice (仓储价格)
 │    ├── ClientDeliveryFee (派送费)
 │    └── ClientSurcharge (加收费)
 │
 ├── Warehouse (仓库)
 │    ├── Zone (区域)
 │    ├── Location (库位)
 │    ├── InboundMachine (入库机/PDA)
 │    └── Container (集装柜)
 │
 ├── Order (集运订单)
 │    ├── Parcel (包裹) — many-to-many via order_parcels
 │    ├── Declarant (申报人)
 │    └── MemberAddress (收件地址)
 │
 ├── Parcel (包裹)
 │    ├── Courier (快递公司)
 │    ├── CargoType (货物类型)
 │    └── ExceptionReport (异常记录)
 │
 ├── ParcelServiceOrder (附加服务订单)
 │    └── ParcelServiceTemplate (服务模板)
 │
 ├── Route (线路模板)
 │    ├── TransportType (运输方式)
 │    ├── CargoType (货物类型)
 │    ├── Carrier (承运商)
 │    └── RouteTransportType (线路-运输方式关联)
 │
 ├── Carrier (承运商)
 │    ├── CarrierTrackingNumber (承运商单号池)
 │    └── CarrierAuthorization (承运商授权)
 │
 ├── CustomsBroker (清关公司)
 │    └── CustomsClearancePoint (清关点)
 │
 ├── WorkOrder (工单)
 │    ├── WorkOrderTemplate (工单模板)
 │    └── WorkflowProcess (流程定义)
 │         └── ProcessInstance (流程实例)
 │
 └── System
      ├── User (员工/操作员)
      ├── Role (角色) — Shield RBAC
      ├── Permission (权限)
      ├── Notification (通知)
      ├── PrintTemplate (打印模板)
      ├── PdaVersion (PDA版本)
      └── ApiCallLog (API日志)
```

### 6.2 关键状态机

**包裹状态流转**:
```
预报 → 已入库 → 已称重 → 已上架 → 已拣货 → 已打包 → 已出库
  ↓                                              ↓
拒收                                           异常(可标记)
```

**订单状态流转**:
```
待拣货 → 待装柜 → 待发运 → 运输中 → 已签收
  ↓         ↓
已取消    待装柜(异常挂起)
```

**申报人认证状态**:
```
待认证 → 认证中 → 认证成功
                 → 认证失败
启用 ⇄ 停用
```

**充值申请状态**:
```
待确认 → 已确认(到账)
       → 已驳回
```

**工单状态**:
```
待开始 → 进行中 → 已完成
       → 已取消
```

### 6.3 定价模型

BFT56 的定价采用**多层次叠加模型**：

```
最终价格 = 线路基础价(按运输方式/货物类型/重量阶梯)
         + 客户线路价覆盖(ClientRoutePrice - 从默认同步或自定义)
         + 附加服务费(按服务类型 × 数量)
         + 客户附加服务覆盖(ClientServiceOverride)
         + 承运商派送费(CarrierDeliveryFee - 末端派送)
         + 承运商加收费(CarrierSurcharge - 偏远/超重等)
         + 仓储费(ClientStoragePrice - 按天/按体积)
         - 客户折扣覆盖
```

价格优先级：客户专属价 > 线路默认价

---

## 七、RBAC 权限模型

### 7.1 角色体系

| 角色 | 权限数 | 说明 |
|------|--------|------|
| 公司管理员 | 362 | 全部功能 + 跨仓数据 |
| 仓库管理员 | 256 | 仓库操作 + 限本仓数据 |
| 客服 | — | 工单/客户查询 |
| 财务 | — | 财务报表 + 充值审核 |

### 7.2 数据权限范围

```
全部 (All) → 企业 (Tenant) → 仓库 (Warehouse) → 部门 (Dept) → 本人 (Self)
```

特殊数据权限控制：
- `查看全部仓库` — 超级管理员可跨仓查看所有数据
- `查看订单财务` — 控制是否可见成本/利润数据
- `修改/删除已封柜或已发运的集装柜` — 高危操作需单独授权

### 7.3 客户端权限模型

`客户端权限 管理(类型默认 + 账号覆盖)` — 两层结构：
1. **客户类型默认权限**：按客户类型（平台客户/企业客户/个人客户）设置默认可见功能
2. **账号覆盖权限**：针对特定客户账号单独调整

---

## 八、关键业务功能细节

### 8.1 包裹管理

- **预报包裹**：客户提前通知即将到达的包裹（快递单号+品名+货类）
- **手动认领**：仓库收到无主包裹后手动匹配客户
- **包裹入库**：扫描快递单号 → 称重 → 测量尺寸 → 分配库位
- **异常处理**：标记异常(破损/少件/禁运品) + 拒收
- **流转事件**：完整记录包裹生命周期每一步

### 8.2 集运订单

- **创建订单**：会员选择包裹 → 选择线路 → 选择申报人 → 生成订单
- **打印面单**：支持多种打印模板（设为默认/复制/套用示例）
- **清关处理**：设置清关公司 + 分配清关单号（从单号池批量生成）
- **承运商处理**：设置承运商单号（从单号池分配）
- **装柜记录**：订单关联集装柜
- **导出申报单**：支持当前筛选和全部导出

### 8.3 单号池管理

- **承运商单号池**：预生成/批量导入 → 按需分配 → 作废回收
- **清关单号池**：预生成/批量生成/批量导入 → 分配 → 作废

### 8.4 附加服务系统

- **服务类型**：定义可用服务种类（如：合箱/分箱/加固/拍照/退货等）
- **服务模板**：预设服务组合
- **服务订单**：客户下单 → 开始执行 → 完成/取消
- **服务工单**：仓库端执行跟踪

### 8.5 打印模板系统

- 支持面单/标签等多种类型
- 操作：设为默认、复制、套用示例
- PDA 端直接打印

### 8.6 Webhook 系统

- 事件投递至客户配置的 Webhook URL
- 投递日志可查看（超管可见）
- 支持重发失败的投递

### 8.7 PDA 系统

- PDA 版本管理 + APK 下载
- 在线会话监控
- 工单模板 + 任务派发
- 扫码入库/上架/拣货/出库

---

## 九、API 分析

### 9.1 API 凭证管理

客户端提供 API 凭证（API Key/Token）管理，供开发者对接：
- 创建/删除凭证
- 重置 Token

### 9.2 API 调用日志

`/admin/system/...` 路径下存在 `Api 调用日志` 资源，完整记录：
- 调用时间
- 调用方
- 接口路径
- 请求/响应
- 耗时

---

## 十、技术架构总结

| 维度 | BFT56 实现 |
|------|-----------|
| 后端框架 | Laravel (PHP) |
| 管理面板 | Filament (Livewire + Tailwind) |
| 权限系统 | Spatie Shield (RBAC) |
| 客户端 | 独立 Laravel 应用 |
| 数据库 | MySQL |
| 文件存储 | 本地/云存储 |
| 前端渲染 | SSR + Livewire SPA |
| 国际化 | zh-TW |
| 部署 | 单服务器 Docker/Nginx |
| 移动端 | PDA (Android APK) |

---

## 十一、与 I56 Framework 对标分析

### BFT56 的优势

1. **业务完整度高**：OMS + WMS + TMS + CRM + Finance 五合一
2. **权限细粒度**：362 个权限项，数据范围（仓库级/公司级/全局）
3. **定价灵活**：多层叠加定价模型，支持客户定制
4. **单号池管理**：承运商单号/清关单号的预生成和分配机制
5. **PDA 集成**：完整的移动端仓库作业方案

### BFT56 的局限（I56 可超越的方向）

1. **单体架构**：Laravel 单体，难以拆分为独立产品
2. **PHP 性能上限**：高并发场景受限
3. **无事件驱动**：模块间强耦合
4. **无工作流引擎**：工单流程较简单
5. **无多语言 SDK**：仅提供 API 凭证，无 SDK
6. **无 Marketplace**：插件/扩展无标准化机制
7. **无多租户隔离能力**：Shared DB 模式，无 Schema/Database 级别隔离
8. **报表能力弱**：仅 4 个盈利报表，无 BI 分析能力
9. **无消息队列/事件总线**：模块间直接调用，拆微服务困难

---

*本文档基于 2026-07-10 对 https://bft56.com 的完整浏览器逆向分析生成。*
