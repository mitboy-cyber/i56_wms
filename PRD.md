# I56 WMS 物流 SaaS 平台 — 产品需求文档 (PRD)

> 版本: 2.4.2 | 最后更新: 2026-07-21 | 文档状态: 正式版

---

## 目录

1. [产品概述与愿景](#1-产品概述与愿景)
2. [完整模块细分](#2-完整模块细分)
3. [关键字段定义](#3-关键字段定义)
4. [业务逻辑流程](#4-业务逻辑流程)
5. [状态机](#5-状态机)
6. [RBAC 角色与权限](#6-rbac-角色与权限)
7. [数据模型概述](#7-数据模型概述)
8. [API 架构](#8-api-架构)
9. [与竞品差距对比](#9-与竞品差距对比)
10. [路线图](#10-路线图)

---

## 1. 产品概述与愿景

### 1.1 产品定位

I56 WMS 是一套面向跨境電商物流企业的仓库管理系统（WMS）SaaS 平台。核心业务场景为：**中国大陆仓库收件 → 集运合单 → 跨境运输 → 台湾清关 → 末端派送**。

平台以"集运"（Consolidation）为业务核心，一个客户将多个来自不同电商平台的快递包裹集中入仓，合并为一个集运订单，通过空运/海运/海快等方式发往台湾，在台湾完成清关后由末端快递派送给最终收件人。

### 1.2 产品愿景

构建业界领先的跨境物流 SaaS 平台，实现端到端的供应链全链条数字化：从包裹预申报、入库称重、拣货打包、装柜出运、清关追踪到末端签收的全生命周期管理，并通过多维度定价引擎和实时盈利分析，帮助物流企业实现精细化运营。

### 1.3 技术架构概览

| 层级 | 技术选型 |
|------|---------|
| 后端语言 | Go 1.24+ |
| 前端框架 | React 18 + TypeScript + Vite + Tailwind CSS |
| 数据库 | PostgreSQL (主) + 内存存储 (开发/测试) |
| 模块管理 | Go Modules Monorepo (github.com/mitboy-cyber/i56_wms) |
| 认证 | JWT Bearer Token + Session Cookie 双模式 |
| 权限模型 | RBAC (基于角色的访问控制) + 多租户数据隔离 |
| 事件系统 | EventBus (Pub/Sub 模式) |
| 工作流 | 内置工作流引擎 |
| 通知 | 多通道通知 (Web/SMS/Email/微信) |
| 部署 | Docker / Docker Compose / Kubernetes / Helm |

### 1.4 应用入口

| 入口 | 路径 | 用户 |
|------|------|------|
| 管理后台 (Admin SPA) | `/admin` | 平台运营人员 |
| 客户端门户 (Client Portal) | `/client` | 客户/代理商 |
| PDA 接口 | `/pda` | 仓库操作员手持设备 |
| OpenAPI | `/api/v1/*` | 第三方系统集成 |
| 健康检查 | `/api/v1/health` | 运维监控 |

---

## 2. 完整模块细分

平台按侧边栏导航结构分为七大功能模块，共 **56 个功能页面/子模块**：

### 2.1 首页 (Dashboard)

| 页面 | 路由 | 功能说明 |
|------|------|---------|
| 仪表盘 | `/admin/dashboard` | 核心运营指标总览：订单数、包裹数、营收、待处理任务、状态分布 |
| 仓库看板 | `/admin/warehouse-board` | 各仓库实时状态：待收货、在库、拣货中、出库中 |

### 2.2 订单管理 (OMS)

| 页面 | 路由 | 功能说明 |
|------|------|---------|
| 集运订单 | `/admin/orders` | 集运订单 CRUD、状态流转、承运商单号设置、导出申报单 |
| 附加服务订单 | `/admin/service-orders` | 附加服务订单管理（拍照、验货、加固、保险等） |

### 2.3 仓库管理 (WMS)

| 页面 | 路由 | 功能说明 |
|------|------|---------|
| 包裹列表 | `/admin/parcels` | 包裹预申报、入库、称重、上架、查询、导出 |
| 附加服务工单 | `/admin/service-workorders` | 附加服务作业工单分配与追踪 |
| 附加服务模板 | `/admin/service-templates` | 附加服务组合模板（打包类、加固类、开箱类） |
| 附加服务类型 | `/admin/service-types` | 服务类型定义（退货寄回、打木箱、内容物拍照等 12 种） |
| PDA 在线会话 | `/admin/pda-sessions` | 实时查看 PDA 在线操作员、设备、当前页面/区域 |
| 集装柜管理 | `/admin/containers` | 集装柜（Container）创建、装柜/封柜/发运状态管理 |
| 仓库列表 | `/admin/warehouses` | 仓库 CRUD、库位管理、区域管理 |
| 入库看板 | `/admin/inbound-board` | 预期到货包裹看板 |
| 仓库作业台 | `/admin/warehouse-console` | 仓库现场大屏作业界面 |
| 员工任务监控 | `/admin/task-monitor` | 操作员任务池分配与进度 |
| PDA 工单模板 | `/admin/pda-workorder-templates` | PDA 端工单模板（收货/拣货/打包/装柜/异常处理） |
| 工单流程管理 | `/admin/workflow-management` | 工作流定义与流程实例管理 |
| 工单列表 | `/admin/work-orders` | 工单全列表（draft → assigned → in_progress → completed） |
| 异常记录 | `/admin/exceptions` | 包裹异常登记、处理、AI 异常检测 |

### 2.4 财务报表 (Finance)

| 页面 | 路由 | 功能说明 |
|------|------|---------|
| 集运订单盈利 | `/admin/reports/order-profit` | 按日期维度的集运订单收入/成本/毛利/毛利率 |
| 附加服务盈利 | `/admin/reports/service-profit` | 附加服务的收入/成本/毛利报表 |
| 客户盈利 | `/admin/reports/client-profit` | 按客户维度的收入/成本/利润/毛利率排行 |
| 路线盈利 | `/admin/reports/route-profit` | 按运输路线的收入/成本/利润分析 |

### 2.5 物流管理 (TMS)

| 页面 | 路由 | 功能说明 |
|------|------|---------|
| 区域组管理 | `/admin/area-groups` | 配送区域组定义（北部/中部/南部/东部） |
| 货物类型 | `/admin/cargo-types` | 货物分类（普货/家具类/一类~六类/易碎品/特货） |
| 承运商列表 | `/admin/carriers` | 运输承运商管理（新竹物流/黑猫宅急便/顺丰等） |
| 快递公司 | `/admin/couriers` | 国内快递公司管理（韵达/中通/极兔等） |
| 清关公司 | `/admin/customs-brokers` | 清关公司（报关行）管理 |
| 清关点管理 | `/admin/customs-points` | 清关点定义（台北/台中/高雄） |
| 线路模板 | `/admin/route-templates` | 运输线路模板 + 定价配置 |
| 运输公司 | `/admin/shipping-providers` | 运输公司管理 |
| 运输方式 | `/admin/transport-modes` | 运输方式定义（空运/海运/海快/陆运/铁路） |
| 物流追踪 | `/admin/logistics-tracking` | 运输轨迹追踪记录 |

### 2.6 客户管理 (CRM)

| 页面 | 路由 | 功能说明 |
|------|------|---------|
| 客户收件地址 | `/admin/customer-addresses` | 台湾收件地址管理（省/市/区/详细地址） |
| 客户申报人 | `/admin/customer-declarants` | 海关申报人管理（个人/公司，实名认证状态） |
| 客户管理 | `/admin/clients` | 客户 CRUD（平台客户/虾皮商家/大客户/同行/普通客户） |
| 客户账号 | `/admin/client-accounts` | 客户端登录账号管理 |
| 客户会员 | `/admin/client-members` | 客户下的终端收件人（会员） |
| 客户充值 | `/admin/client-recharges` | 充值申请审核（银行转账/微信/支付宝/线下） |
| 余额日志 | `/admin/balance-logs` | 客户余额变动流水 |
| 充值记录 | `/admin/recharge-records` | 充值申请记录 |
| 客户价格 | `/admin/pricing` | 客户五维定价体系 |
| 月结对账单 | `/admin/monthly-statements` | 按月生成的客户账单 |
| 客户端权限 | `/admin/client-panel-permissions` | 控制客户在门户可使用的功能模块 |

### 2.7 系统设置 (System)

| 页面 | 路由 | 功能说明 |
|------|------|---------|
| 通知管理 | `/admin/notifications` | 系统通知的查看与发送 |
| 打印模板 | `/admin/print-templates` | 面单/清关单/承运商单/标签/发票模板 |
| 角色管理 | `/admin/roles` | RBAC 角色定义与权限分配 |
| 员工管理 | `/admin/employees` | 系统操作员工（用户）管理 |
| 系统参数 | `/admin/system/params` | 全局系统参数配置 |
| 设备管理 | `/admin/devices` | 打印机/扫码枪/称重机等硬件设备 |
| 存储配置 | `/admin/storage` | 对象存储配置（MinIO/S3/OSS） |
| 品牌设置 | `/admin/system/brand` | 平台品牌信息配置 |
| API 集成配置 | `/admin/system/api-*` | 快递/清关/通知/打印/存储/设备/EZWay 等 API 对接配置 |
| 定时任务 | `/admin/system/scheduler` | 调度任务管理 |
| 审计日志 | `/admin/system/audit-logs` | 操作审计追踪 |
| AI 设置 | `/admin/system/ai-settings` | AI 智能分类/异常检测/翻译设置 |

---

## 3. 关键字段定义

### 3.1 集运订单 (Consolidation Order)

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 主键 |
| `tenant_id` | int64 | 租户 ID |
| `order_no` | string | 订单编号（自动生成） |
| `client_id` | int64 | 客户 ID |
| `member_id` | int64 | 会员（收件人）ID |
| `member_address_id` | int64 | 收件地址 ID |
| `warehouse_id` | int64 | 发货仓库 ID |
| `route_id` | int64 | 运输线路 ID |
| `parcel_ids` | []int64 | 包含的包裹 ID 列表 |
| `parcel_count` | int | 包裹数量 |
| `total_weight` | float64 | 总实际重量 (kg) |
| `total_chargeable_weight` | float64 | 总计费重量 (kg) |
| `shipping_fee` | float64 | 运输费用 |
| `service_fee` | float64 | 附加服务费 |
| `total_price` | float64 | 订单总价 |
| `status` | ConsolidationOrderStatus | 订单状态 |
| `declarant_id` | int64 | 申报人 ID |
| `customs_number` | string | 海关申报单号 |
| `carrier_tracking_no` | string | 承运商追踪号 |
| `remark` | string | 备注 |

### 3.2 包裹 (Parcel)

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 主键 |
| `tenant_id` | int64 | 租户 ID |
| `warehouse_id` | int64 | 入库仓库 ID |
| `client_id` | int64 | 所属客户 ID |
| `tracking_number` | string | 快递追踪号 |
| `courier_code` | string | 快递公司代码 |
| `cargo_type` | string | 货物类型 |
| `product_name` | string | 商品名称 |
| `parcel_name` | string | 包裹名称/标签 |
| `actual_weight` | float64 | 实际重量 (kg) |
| `length` | float64 | 长 (cm) |
| `width` | float64 | 宽 (cm) |
| `height` | float64 | 高 (cm) |
| `status` | ParcelStatus | 包裹状态 |
| `is_abnormal` | bool | 是否异常 |
| `location_code` | string | 库位编号 |
| `image_urls` | []string | 包裹图片 |
| `order_id` | *int64 | 关联订单 ID |

**体积重计算公式**：`(长 × 宽 × 高) / 6000`  
**计费重**：`max(实际重, 体积重)`

### 3.3 客户 (Client)

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 主键 |
| `tenant_id` | int64 | 租户 ID |
| `name` | string | 客户名称 |
| `code` | string | 客户编码 |
| `client_type` | ClientType | 客户类型：platform / shopee / major / peer / normal |
| `contact_name` | string | 联系人 |
| `contact_phone` | string | 联系电话 |
| `contact_email` | string | 联系邮箱 |
| `balance` | float64 | 账户余额 |
| `is_active` | bool | 是否启用 |

### 3.4 仓库 (Warehouse)

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 主键 |
| `tenant_id` | int64 | 租户 ID |
| `name` | string | 仓库名称 |
| `code` | string | 仓库编码 |
| `address` | string | 仓库地址 |
| `contact` | string | 联系人 |
| `phone` | string | 联系电话 |
| `is_active` | bool | 是否启用 |

### 3.5 集装柜 (Container)

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 主键 |
| `warehouse_id` | int64 | 所属仓库 |
| `container_no` | string | 柜号 |
| `container_type` | ContainerType | 柜型：20GP / 40GP / 40HQ / 45HQ |
| `seal_no` | string | 封条号 |
| `route_id` | int64 | 运输线路 |
| `status` | ContainerStatus | 状态 |
| `max_capacity` | float64 | 最大装载量 (kg) |
| `current_weight` | float64 | 当前装载重量 (kg) |
| `parcel_count` | int | 已装载包裹数 |

### 3.6 运输线路模板 (Route Template)

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 主键 |
| `name` | string | 线路名称（如"深圳→台北(空运)"） |
| `code` | string | 线路编码（如"WH001-AIR"） |
| `transport_type` | string | 运输方式：air / sea / sea_express / air_special / land / rail |
| `cargo_type` | string | 货物类型 |
| `tax_type` | string | 包税类型：不包税 / 频税 / 全包税 |
| `weight_price` | float64 | 重量单价 (元/kg) |
| `volume_price` | float64 | 体积单价 (元/才) |
| `min_charge` | float64 | 最低收费 |
| `first_weight` | float64 | 首重 (kg) |
| `first_weight_price` | float64 | 首重价格 |
| `cont_weight_price` | float64 | 续重单价 |
| `volume_coeff` | int | 体积系数（默认 6000） |
| `estimated_days` | int | 预计运输天数 |

### 3.7 附加服务订单 (Service Order)

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 主键 |
| `parcel_id` | int64 | 关联包裹 |
| `order_id` | int64 | 关联订单 |
| `service_type` | ServiceOrderType | 服务类型：photos / inspection / repack / remove_box / reinforce / insurance / customs_declaration |
| `service_name` | string | 服务名称 |
| `status` | ServiceOrderStatus | 状态 |
| `price` | float64 | 服务费 |
| `operator_id` | int64 | 操作员 |
| `result_note` | string | 结果备注 |
| `result_images` | []string | 结果照片 |

### 3.8 财务流水 (Ledger Entry)

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 主键 |
| `client_id` | int64 | 客户 ID |
| `amount` | float64 | 金额（正=收入/充值，负=支出） |
| `balance_after` | float64 | 变动后余额 |
| `type` | LedgerType | 类型：recharge / deduction / refund / adjustment |
| `reference_no` | string | 关联单号 |
| `description` | string | 描述 |

### 3.9 申报人 (Declarant)

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 主键 |
| `client_id` | int64 | 客户 ID |
| `member_id` | int64 | 会员 ID |
| `type` | DeclarantType | 类型：individual / company |
| `name` | string | 姓名/公司名 |
| `id_number` | string | 身份证号（个人） |
| `company_tax_id` | string | 统一编号（公司） |
| `phone` | string | 联系电话 |
| `auth_status` | DeclarantAuthStatus | 认证状态：pending / verifying / verified / failed |

### 3.10 客户定价五维体系

**维度一：客户 × 线路价 (RoutePriceModel)**

| 字段 | 类型 | 说明 |
|------|------|------|
| `client_id` | int64 | 客户 |
| `route_name` | string | 线路名 |
| `transport_type` | string | 运输方式 |
| `cargo_type` | string | 货物类型 |
| `tax_mode` | string | 包税模式 |
| `weight_price` | float64 | 重量单价 |
| `first_weight_price` | float64 | 首重价格 |
| `cont_weight_price` | float64 | 续重单价 |

**维度二：客户 × 仓储价 (StoragePriceModel)**

| 字段 | 类型 | 说明 |
|------|------|------|
| `client_id` | int64 | 客户 |
| `warehouse_id` | int64 | 仓库 |
| `free_days` | int | 免费存放天数 |
| `daily_rate` | float64 | 超期单价 (元/kg/天) |
| `max_storage_days` | int | 最大存放天数 |

**维度三：客户 × 派送费 (DeliveryFeeModel)**

| 字段 | 类型 | 说明 |
|------|------|------|
| `client_id` | int64 | 客户 |
| `carrier_id` | int64 | 承运商 |
| `customs_point` | string | 清关点：台北/台中/高雄 |
| `area` | string | 配送区域：北部/中部/南部/东部 |
| `delivery_method` | string | 派送方式：宅配/专车/自取 |
| `fee` | float64 | 基础费用 |
| `free_threshold` | float64 | 免运门槛（重量阈值） |

**维度四：客户 × 加收费 (SurchargeModel)**

| 字段 | 类型 | 说明 |
|------|------|------|
| `client_id` | int64 | 客户 |
| `carrier_id` | int64 | 承运商 |
| `charge_type` | string | 超长费/超材费/栈板费/偏远费/上楼费 |
| `tier` | string | 小板/大板/— |
| `price` | float64 | 加收金额 |

**维度五：客户 × 附加服务 (ServicePriceModel)**

| 字段 | 类型 | 说明 |
|------|------|------|
| `client_id` | int64 | 客户 |
| `service_type` | string | 服务类型（木箱/验货/拍照/换标） |
| `unit_price` | float64 | 单价 |
| `price_mode` | string | fixed / per_item / per_kg / per_order |

---

## 4. 业务逻辑流程

### 4.1 集运订单生命周期

```
客户创建预报 → 包裹入仓 → 客户发起集运 → 仓库拣货 → 打包复核 → 装柜 → 发运 → 
清关 → 末端派送 → 签收 → 订单完成
```

**详细步骤：**

1. **包裹预申报 (Pre-declare)**：客户在门户/API 预报快递单号、商品名称、货物类型
2. **包裹入库**：快递到达仓库 → 操作员 PDA 扫描签收 → 称重测量长宽高 → 拍照 → 分配库位上架
3. **创建集运订单**：客户选择多个已上架包裹 → 选择运输线路 → 选择申报人 → 选择收件地址 → 提交合并
4. **拣货 (Picking)**：系统生成拣货任务 → 操作员 PDA 按单拣取包裹 → 扫描确认
5. **打包复核 (Packing)**：包裹移至打包区 → 合并打包 → 称重复核 → 打印面单
6. **装柜 (Loading)**：包裹送至装柜区 → 扫码装柜 → 记录柜号
7. **发运 (Shipping)**：封柜 → 生成承运商追踪号 → 发运
8. **清关 (Customs)**：到达台湾 → 海关申报 → 查验 → 放行
9. **派送 (Delivery)**：末端快递取件 → 配送 → 收件人签收
10. **完成**：订单完结，生成盈利数据

### 4.2 包裹生命周期

```
预申报 → 已入库 → 已核重 → 已上架 → 已拣货 → 已打包 → 已出库 → 
已送装柜区 → 已装柜 → 运输中 → 清关中 → 派送中 → 已签收
```

**异常分支**：任何中间状态可转入"异常"状态；异常包裹可转回"已上架"或"已退件"。

### 4.3 财务流程

**充值流程**：
```
客户提交充值申请（上传凭证） → 财务审核 → 确认/拒绝 → 更新余额 → 生成流水记录
```

**自动扣费流程**：
```
订单创建 → 冻结费用 → 订单发运 → 正式扣除 → 生成扣费流水 → 更新盈利报表
```

**对账流程**：
```
月末 → 系统自动生成月度对账单 → 客户确认 → 开票 → 收款核销
```

**盈利分析流程**：
```
订单完成 → 计算收入（运费+服务费） → 计算成本（运输成本+操作成本） → 
计算毛利 → 按订单/客户/路线/附加服务四个维度汇总
```

### 4.4 附加服务工单流程

```
客户下单附加服务 → 系统生成服务工单 → 分配操作员 → 
操作员 PDA 接单 → 执行操作（拍照/验货/加固等） → 
提交结果（照片+备注） → 工单完成 → 计入费用
```

### 4.5 PDA 仓库作业流程

```
操作员登录 PDA（PIN码） → 选择功能菜单：
  ├── 📦 包裹收货 → 扫描快递单号 → 称重 → 测量尺寸 → 拍照 → 分配库位 → 上架
  ├── ⚖️ 称重核重 → 扫描包裹 → 复称 → 确认
  ├── 📍 上架入库 → 扫描包裹 → 扫描库位 → 确认上架
  ├── 🛒 订单拣货 → 扫描订单 → 逐个扫描包裹 → 确认拣齐
  ├── 📋 打包复核 → 扫描包裹 → 合并打包 → 称重核对 → 打印面单
  └── 🔍 快件查询 → 扫描/输入单号 → 显示包裹状态与位置
```

---

## 5. 状态机

### 5.1 集运订单状态机 (ConsolidationOrderStatus)

```
draft ───────────────────────→ cancelled
  │
  ▼
pending_merge
  │
  ▼
merged
  │
  ▼
weighing
  │
  ▼
packing
  │
  ▼
packed
  │
  ▼
outbound
  │
  ▼
shipped
  │
  ▼
completed
```

| 状态 | 中文 | 可转入状态 |
|------|------|-----------|
| `draft` | 草稿 | pending_merge, cancelled |
| `pending_merge` | 待合单 | merged |
| `merged` | 已合单 | weighing |
| `weighing` | 称重中 | packing |
| `packing` | 打包中 | packed |
| `packed` | 已打包 | outbound |
| `outbound` | 已出库 | shipped |
| `shipped` | 已发运 | completed |
| `completed` | 已完成 | — (终态) |
| `cancelled` | 已取消 | — (终态) |

### 5.2 包裹状态机 (ParcelStatus)

```
pre_declared → received → weighed → stored ──→ picked → packed → outbound
                                           │
                                           ├──→ container_area → loaded → shipped → customs → delivering → delivered
                                           │
                                           └──→ abnormal → received / stored / returned
                                                                              returned (终态)
```

| 状态 | 中文 | 可转入状态 |
|------|------|-----------|
| `pre_declared` | 已预报 | received, returned |
| `received` | 已入库 | weighed, returned |
| `weighed` | 已核重 | stored, returned |
| `stored` | 已上架 | picked, container_area, returned, abnormal |
| `picked` | 已拣货 | packed, returned |
| `packed` | 已打包 | shipped, loaded, container_area |
| `outbound` | 已出库 | container_area |
| `container_area` | 已送装柜区 | loaded, shipped |
| `loaded` | 已装柜 | shipped, customs |
| `shipped` | 运输中 | delivering, delivered, customs |
| `customs` | 清关中 | delivering, delivered, returned, abnormal |
| `delivering` | 派送中 | delivered, abnormal |
| `delivered` | 已签收 | — (终态) |
| `abnormal` | 异常 | stored, returned |
| `returned` | 已退件 | — (终态) |

**关键状态转换所需字段校验**：
- `pre_declared → received`：需要 weight, length, width, height
- `received → stored` / `weighed → stored`：需要 location_barcode
- `stored → picked`：需要 order_id

### 5.3 集装柜状态机 (ContainerStatus)

```
available → loading → loaded → sealed → shipped
```

| 状态 | 中文 | 可转入状态 |
|------|------|-----------|
| `available` | 空闲 | loading |
| `loading` | 装柜中 | loaded |
| `loaded` | 已装载 | sealed |
| `sealed` | 已封柜 | shipped |
| `shipped` | 已发运 | — (终态) |

### 5.4 附加服务订单状态机 (ServiceOrderStatus)

```
pending → processing → completed
    ↘
     cancelled
```

| 状态 | 中文 | 可转入状态 |
|------|------|-----------|
| `pending` | 待处理 | processing, cancelled |
| `processing` | 处理中 | completed |
| `completed` | 已完成 | — (终态) |
| `cancelled` | 已取消 | — (终态) |

### 5.5 工单状态机 (WorkOrderStatus)

```
draft → assigned → in_progress → completed
  ↘                    ↘
   cancelled         cancelled
```

| 状态 | 中文 | 可转入状态 |
|------|------|-----------|
| `draft` | 草稿 | assigned, cancelled |
| `assigned` | 已分配 | in_progress, cancelled |
| `in_progress` | 进行中 | completed, cancelled |
| `completed` | 已完成 | — (终态) |
| `cancelled` | 已取消 | — (终态) |

### 5.6 充值状态机 (RechargeStatus)

```
pending → confirmed
   ↘
 rejected
```

---

## 6. RBAC 角色与权限

### 6.1 角色定义

系统内置五大角色：

| 角色 | Slug | 说明 | 权限范围 |
|------|------|------|---------|
| 超级管理员 | `super_admin` | 拥有系统全部权限 | 所有模块全部操作 |
| 仓库管理员 | `warehouse_manager` | 负责仓库日常运营 | 包裹管理、仓库管理、订单查看、入库看板、报表查看 |
| 财务专员 | `finance_specialist` | 负责财务管理 | 财务流水、客户充值、对账结算、财务报表面临 |
| 客服专员 | `customer_service` | 处理客户咨询 | 包裹查看、订单查看/创建/编辑、客户查看、线路查看 |
| 操作员 | `operator` | PDA 手持终端操作 | 包裹收货、称重、上架、拣货、打包（通过 PDA 授权） |

### 6.2 权限模块矩阵

| 模块 | 权限 Slug | 超管 | 仓管 | 财务 | 客服 | 操作员 |
|------|----------|:--:|:--:|:--:|:--:|:--:|
| 仪表板 | `dashboard:view` | ✅ | ✅ | ✅ | ✅ | — |
| 包裹 | `parcel:list/create/update/delete/inbound/export` | ✅ | ✅ | ✅ | ✅ | ✅ |
| 订单 | `order:list/create/update/delete/approve/cancel/export` | ✅ | ✅ | — | ✅ | — |
| 客户 | `client:list/create/update/delete/finance/permissions` | ✅ | ✅ | ✅ | ✅ | — |
| 仓库 | `warehouse:list/create/update/delete/locations/inbound-board` | ✅ | ✅ | — | ✅ | — |
| 线路 | `route:list/create/update/delete/pricing` | ✅ | ✅ | — | ✅ | — |
| 财务 | `finance:list/recharge/reports/export/settlement` | ✅ | — | ✅ | — | — |
| 报表 | `report:list/export` | ✅ | ✅ | ✅ | — | — |
| 附加服务 | `service:list/manage` | ✅ | ✅ | — | ✅ | — |
| 系统 | `user:list/create/update/delete` | ✅ | — | — | — | — |
| 系统 | `role:list/create/update/delete` | ✅ | — | — | — | — |
| 系统 | `system:settings/audit` | ✅ | — | — | — | — |

### 6.3 数据权限范围 (DataScope)

| 级别 | 范围 | 适用场景 |
|------|------|---------|
| `ALL (1)` | 全部数据 | 超级管理员 |
| `TENANT (2)` | 本租户全部数据 | 仓库管理员/财务 |
| `WAREHOUSE (3)` | 指定仓库数据 | 仓库操作员 |
| `SELF (5)` | 仅本人数据 | PDA 操作员个人任务 |

### 6.4 客户端权限 (Client Panel Permission)

客户门户支持独立的权限控制体系，每个客户可配置对不同功能模块的访问权限（查看/创建/编辑/删除/导出）。权限级别支持分级（Level），如基础版/标准版/专业版，与客户套餐绑定。

---

## 7. 数据模型概述

### 7.1 核心实体关系图

```
┌─────────────┐     ┌──────────────────┐     ┌────────────────┐
│   Tenant    │1───*│    Warehouse     │1───*│   Container    │
└─────────────┘     └──────────────────┘     └────────────────┘
       │                      │
       │                   1──*│
       │              ┌──────────────────┐
       │              │     Location     │
       │              └──────────────────┘
       │
       ├──────────*──┐
       │     ┌──────────────┐     ┌──────────────────┐
       │     │    Client    │1───*│  Consolidation   │
       │     └──────────────┘     │     Order        │
       │            │             └──────┬───────────┘
       │            │                    │
       │     1─────*│              *─────│─────*
       │     ┌──────────────┐     ┌──────────────────┐
       │     │ ClientMember │     │   OrderParcel    │
       │     └──────┬───────┘     └────────┬─────────┘
       │            │                      │
       │     1─────*│                *─────│─────*
       │     ┌──────────────┐     ┌──────────────────┐
       │     │MemberAddress │     │     Parcel       │
       │     └──────────────┘     └──────────────────┘
       │            │                      │
       │     1─────*│                *─────│─────*
       │     ┌──────────────┐     ┌──────────────────┐
       │     │  Declarant   │     │  ServiceOrder    │
       │     └──────────────┘     └──────────────────┘
       │
       ├──────────*──┐
       │     ┌──────────────┐     ┌──────────────────┐
       │     │   Route      │1───*│ RouteTemplate    │
       │     └──────────────┘     └──────────────────┘
       │            │
       │     1─────*│
       │     ┌──────────────┐
       │     │ Courier/     │
       │     │ Carrier      │
       │     └──────────────┘
       │
       └──────────*──┐
              ┌──────────────┐
              │  LedgerEntry │
              └──────────────┘
```

### 7.2 数据库表概览（预估 ~68 张表）

| 领域 | 表数量 | 核心表 |
|------|--------|--------|
| Framework Core | 8 | users, roles, permissions, role_permission, audit_logs, notifications, scheduler_jobs, tenant_configs |
| Customer (客户) | 14 | clients, client_users, client_members, member_addresses, declarants, client_ledgers, client_recharges, client_route_prices, client_delivery_fees, client_surcharges, client_storage_prices, client_service_overrides, statements, warehouse_authorizations |
| Warehouse (仓库) | 9 | warehouses, zones, zone_types, locations, location_types, containers, container_loading_records, inbound_machines, warehouse_configs |
| Parcel (包裹) | 5 | parcels, parcel_events, exception_reports, parcel_photos, parcel_dimensions |
| Order (订单) | 6 | consolidation_orders, order_parcels, order_events, service_orders, service_templates, service_types |
| Transport (运输) | 12 | area_groups, cargo_types, carriers, carrier_numbers, couriers, customs_brokers, customs_points, customs_numbers, routes, shipping_providers, transport_types, trackings |
| Finance (财务) | 4 | invoices, payments, recharge_logs, profit_reports |
| WorkOrder (工单) | 5 | work_orders, work_order_templates, workflow_processes, process_instances, task_assignments |
| System (系统) | 5 | notifications, print_templates, pda_versions, operator_sessions, api_call_logs |

---

## 8. API 架构

### 8.1 认证机制

系统支持双认证模式：
- **Cookie Session**：管理后台 SPA 使用 `admin_session` Cookie（HttpOnly, Secure, SameSite=Lax）
- **Bearer Token**：API 调用和客户端门户使用 JWT Bearer Token（Header: `Authorization: Bearer <token>`）

### 8.2 多租户

通过 `X-Tenant-ID` 请求头指定租户上下文，中间件自动注入 Context。默认租户为 `default`。

### 8.3 API 端点结构

```
/admin/api/*          → 管理后台 JSON API (需要 admin session / bearer token)
/client/api/*         → 客户端门户 API (需要 JWT client token)
/pda/api/*            → PDA 手持终端 API
/admin/api/v1/*       → OpenAPI v1 端点
/admin/login          → 管理后台登录
/client/login         → 客户端登录
/admin/logout         → 管理后台登出
```

### 8.4 管理后台 API 汇总

**订单管理 (OMS)：**
```
GET    /admin/api/orders                    → 集运订单列表
POST   /admin/api/orders                    → 创建集运订单
GET    /admin/api/orders/{id}               → 订单详情
PUT    /admin/api/orders/{id}/status        → 订单状态流转
GET    /admin/api/service-orders            → 附加服务订单列表
```

**仓库管理 (WMS)：**
```
GET    /admin/api/parcels                   → 包裹列表
POST   /admin/api/parcels                   → 预申报包裹
GET    /admin/api/warehouses                → 仓库列表
POST   /admin/api/warehouses                → 创建仓库
GET    /admin/api/containers                → 集装柜列表
GET    /admin/api/pda-sessions              → PDA 在线会话
GET    /admin/api/work-orders               → 工单列表
GET    /admin/api/print-templates           → 打印模板
GET    /admin/api/exceptions                → 异常记录
GET    /admin/api/inbound-board             → 入库看板
GET    /admin/api/warehouse-board           → 仓库看板
GET    /admin/api/warehouse-console         → 仓库作业台
GET    /admin/api/task-monitor              → 员工任务监控
```

**物流管理 (TMS)：**
```
GET    /admin/api/area-groups               → 区域组列表
GET    /admin/api/cargo-types               → 货物类型
GET    /admin/api/carriers                  → 承运商列表
GET    /admin/api/couriers                  → 快递公司
GET    /admin/api/customs-brokers           → 清关公司
GET    /admin/api/customs-points            → 清关点
GET    /admin/api/route-templates           → 线路模板
GET    /admin/api/transport-modes           → 运输方式
GET    /admin/api/logistics-tracking        → 物流追踪
```

**客户管理 (CRM)：**
```
GET    /admin/api/clients                   → 客户列表
POST   /admin/api/clients                   → 创建客户
GET    /admin/api/declarants                → 申报人列表
GET    /admin/api/members                   → 会员列表
GET    /admin/api/addresses                 → 收件地址
GET    /admin/api/ledger                    → 余额流水
GET    /admin/api/client-accounts           → 客户账号
GET    /admin/api/client-recharges          → 充值记录
GET    /admin/api/client-pricing            → 客户价格
GET    /admin/api/monthly-statements        → 月结对账单
GET    /admin/api/client-panel-perms        → 客户端权限
POST   /admin/api/client-panel-perms/batch  → 批量保存客户端权限
```

**财务管理 (Finance)：**
```
GET    /admin/api/finance/order-profit      → 订单盈利报表
GET    /admin/api/finance/revenue-report    → 收入报表
GET    /admin/api/finance/cost-report       → 成本报表
GET    /admin/api/finance/profit-loss       → 损益报表
GET    /admin/api/finance/cash-flow         → 现金流报表
GET    /admin/api/report/order-profit       → 订单盈利（日期维度）
GET    /admin/api/report/service-profit     → 附加服务盈利
GET    /admin/api/report/client-profit      → 客户盈利
GET    /admin/api/report/route-profit       → 路线盈利
```

**系统管理 (System)：**
```
GET    /admin/api/roles                     → 角色列表
GET    /admin/api/employees                 → 员工列表
GET    /admin/api/notifications             → 通知列表
GET    /admin/api/system/params             → 系统参数
GET    /admin/api/system/audit-logs         → 审计日志
GET    /admin/api/system/scheduler          → 定时任务
```

**仪表盘 (Dashboard)：**
```
GET    /admin/api/dashboard/stats           → 核心指标统计
GET    /admin/api/dashboard                 → 仪表盘详细数据
GET    /admin/api/dashboard/order-status    → 订单状态分布
GET    /admin/api/dashboard/revenue-by-route → 按路线营收
GET    /admin/api/pda/online-sessions       → PDA 在线统计
GET    /admin/api/pda/scan-logs             → PDA 扫描日志
```

### 8.5 响应格式

所有 API 统一 JSON 响应，成功时直接返回数据对象或数组，错误时返回：

```json
{
  "error": "错误信息描述"
}
```

创建成功返回 HTTP 201 并包含创建的资源对象。

---

## 9. 与竞品差距对比

> 以 BFT56（同行业成熟的集运 WMS 系统）为主要对标竞品。

### 9.1 功能覆盖总览

| 模块 | BFT56 功能数 | I56 已覆盖 | 覆盖率 |
|------|:---------:|:--------:|:-----:|
| CRM/客户管理 | 15 | 8 | 53% |
| WMS/仓库管理 | 10 | 3 | 30% |
| OMS/订单管理 | 1 | 1 | 100% |
| 附加服务 | 4 | 3 | 75% |
| TMS/物流管理 | 13 | 4 | 31% |
| 工单系统 | 5 | 3 | 60% |
| 系统管理 | 9 | 6 | 67% |
| **总计** | **56** | **28** | **50%** |

### 9.2 差距详细分析

#### 🔴 P0 — 客户定价体系（8 项缺失）

BFT56 的定价体系覆盖 5 个维度，I56 仅在各 domain 模型中定义了结构但前端页面和管理功能尚未完整实现：

| 差距项 | BFT56 实现 | I56 现状 |
|--------|-----------|---------|
| 客户线路价 | 按客户+线路+货物类型+包税模式的多维定价表 | Domain 模型已定义，前端 CRUD 待完善 |
| 客户派送费 | 按客户+承运商+清关点+区域+派送方式的费用表 | Domain 模型已定义，前端 CRUD 待完善 |
| 客户加收费 | 超长费/超材费/栈板费/偏远费/上楼费 | Domain 模型已定义，前端 CRUD 待完善 |
| 客户仓储价 | 免仓期+超期日费率+最大存放天数 | Domain 模型已定义，前端 CRUD 待完善 |
| 客户附加服务覆盖 | 按客户覆盖默认服务单价 | Domain 模型已定义，前端 CRUD 待完善 |
| 客户线路价 | — | 同上 |
| 承运商授权 | 控制客户可使用的承运商 | 缺失 |
| 仓库授权 | 控制客户可使用的仓库 | 缺失 |

#### 🟡 P1 — 仓库维度（7 项缺失）

| 差距项 | BFT56 实现 | I56 现状 |
|--------|-----------|---------|
| 库位管理 | 分级库位：货架/托盘位/地面位 | 部分定义（Zone/Location domain），前端缺失 |
| 库位类型 | 库位类型分类管理 | 缺失 |
| 区域管理 | 收货区/存储区/拣货区/发货区 | 部分定义，前端缺失 |
| 区域类型 | 区域分类管理 | 缺失 |
| 集装柜管理 | 柜号/封条号/装载管理 | Domain 已定义，前端有 Container 页面 |
| 入库机 | 自动化入库设备管理 | 缺失 |
| 员工任务 | PDA 任务分配与监控 | Domain 已定义，前端 TaskMonitor 页面已有 |

#### 🟡 P1 — 物流维度（9 项缺失）

| 差距项 | BFT56 实现 | I56 现状 |
|--------|-----------|---------|
| 区域组 | 配送区域组（北部/中部/南部/东部） | Domain 模型已定义，前端页面已有 |
| 承运商单号池 | 预先导入承运商追踪号池 | 缺失 |
| 装柜记录 | 记录每个柜装载哪些订单 | 缺失 |
| 清关公司 | 报关行管理 | Domain 模型已定义，前端页面已有 |
| 清关点 | 清关地点（台北/台中/高雄） | Domain 模型已定义，前端页面已有 |
| 清关单号池 | 海关申报单号池 | 缺失 |
| 物流追踪 | 轨迹追踪记录 | 缺失 |
| 运输公司 | 物流供应商管理 | Domain 模型已定义，前端页面已有 |
| 运输方式 | 空运/海运/海快/陆运/铁路 | Domain 模型已定义，前端页面已有 |

#### 🟢 P2 — 系统增强

| 差距项 | BFT56 实现 | I56 现状 |
|--------|-----------|---------|
| 审计日志 | API 调用记录 + 操作审计 | 后端 skeleton 已有 |
| API 日志 | 详细请求/响应日志 | 缺失 |
| PDA 版本管理 | 手持终端 APK 版本控制 | 缺失 |
| 对账单 | 月度客户对账 | Domain 模型已定义，前端页面已有 |
| BI 报表 | 可视化分析报表 | 4 维盈利报表已实现 |
| Webhook | 事件推送 | 框架支持，Webhook repo 已定义 |

### 9.3 I56 架构优势（超越 BFT56）

| 维度 | BFT56 | I56 Framework |
|------|-------|---------------|
| 架构模式 | 单体 Laravel (Filament) | 模块化 Monolith → 可演进至微服务 |
| 事件驱动 | 无 | EventBus (Pub/Sub) |
| 工作流引擎 | 无 | 内置 Workflow Engine |
| 多租户 | 基础支持 | 支持 SharedTable/Schema/Database 三种策略 |
| API/SDK | REST API | REST + SDK (Go/Python/JS/Java) |
| 多语言 | 繁体中文 | i18n 完整支持 |
| 部署 | 单服务器 | Docker Compose / K8s / Helm |
| 插件系统 | 无 | Plugin Registry + Shopify/FedEx/Stripe 插件 |
| AI 能力 | 无 | AI 包裹分类、异常检测、智能翻译 |

---

## 10. 路线图

### Phase 1: Foundation Core ✅ (已完成)

- [x] Framework Core: auth, rbac, tenant, events, workflow, scheduler, storage, notification, logger
- [x] Domain models: Order, Parcel, Client, Warehouse, Route, Finance, Service, WorkOrder
- [x] Admin SPA + React Router + API integration
- [x] In-memory repositories (开发阶段)
- [x] PostgreSQL 数据库连接支持
- [x] PDA 基础 API (收货/称重/上架/拣货/打包/查询)
- [x] Plugin System (Shopify, FedEx, Stripe)
- [x] AI 模块：包裹分类器、异常检测器、翻译引擎、成本优化器

### Phase 2: Complete Business Closure (Q3 2026)

- [ ] 客户多维定价体系完整实现（5 维定价前端 CRUD）
- [ ] 库位/区域管理体系完善（Zone/Location CRUD）
- [ ] 集装柜管理完善（装柜记录、单号池）
- [ ] 承运商/清关单号池
- [ ] 物流追踪 API 集成
- [ ] 月结对账单自动生成
- [ ] 承运商授权 + 仓库授权
- [ ] PDA 版本管理
- [ ] Webhook 事件推送完善

### Phase 3: Production Hardening (Q4 2026)

- [ ] PostgreSQL 迁移（从内存存储切换为 PostgreSQL 生产级）
- [ ] 数据迁移工具
- [ ] 性能优化（数据库索引、查询优化、缓存层 Redis）
- [ ] 全文搜索（Elasticsearch）
- [ ] 审计日志完整实现
- [ ] API 限流与安全加固
- [ ] 负载测试与容量规划
- [ ] Kubernetes Helm Chart 生产级部署

### Phase 4: Ecosystem & Growth (2027)

- [ ] OpenAPI 规范生成（Swagger/OpenAPI 3.0）
- [ ] 多语言 SDK 发布（Go / Python / JavaScript / Java）
- [ ] 可视化 BI 报表引擎（拖拽式自定义报表）
- [ ] 打印模板可视化设计器
- [ ] 工作流可视化设计器
- [ ] 插件市场
- [ ] Flutter 跨平台 PDA 应用
- [ ] 多仓库协同调度
- [ ] 实时库存大盘
- [ ] Fulfillment (一件代发) 业务线

---

## 附录 A: 术语对照表 (Ubiquitous Language)

| 中文 | English | 定义 |
|------|---------|------|
| 租户 | Tenant | 平台上的独立企业/组织 |
| 仓库 | Warehouse | 物理仓储设施 |
| 客户 | Client | 使用平台的集运客户 |
| 会员 | Member | 客户的子账户（收件人） |
| 申报人 | Declarant | 海关申报人，需实名认证 |
| 承运商 | Carrier | 运输服务提供商 |
| 快递公司 | Courier | 国内快递公司 |
| 包裹 | Parcel | 单个快递包裹 |
| 集运订单 | Consolidation Order | 合并多个包裹的统一发运订单 |
| 集装柜 | Container | 海运/空运集装箱 |
| 线路 | Route | 运输路线 |
| 计费重 | Chargeable Weight | max(实际重, 体积重) |
| 体积重 | Volumetric Weight | (长×宽×高)/6000 |
| 库位 | Location | 仓库内存储位置 |
| 区域 | Zone | 仓库功能分区 |

## 附录 B: 附加服务类型清单

| 服务 | Code | 分类 | 价格模式 | 单价 |
|------|------|------|---------|------|
| 退货寄回 | RETURN_GOODS | 退货类 | fixed | ¥2.00 |
| 易碎品贴纸 | FRAGILE_STICKER | 打包类 | fixed | ¥0.10 |
| 外箱标识拍照 | OUTERBOX_PHOTO | 加固类 | fixed | ¥0.10 |
| 打木箱 | WOODEN_CRATE | 加固类 | fixed | ¥80.00 |
| 打木架 | WOODEN_FRAME | 加固类 | fixed | ¥50.00 |
| 包装气柱袋 | WRAP_AIRBAG | 加固类 | fixed | ¥5.00 |
| 包装气泡棉 | WRAP_BUBBLE | 加固类 | fixed | ¥2.00 |
| 包裹拆分 | SPLIT_PARCEL | 开箱类 | per_qty | ¥2.00 |
| 填充气柱袋 | FILL_AIRBAG | 开箱类 | fixed | ¥5.00 |
| 内容物拍照 | CONTENT_PHOTO | 开箱类 | fixed | ¥1.00 |
| 清点数量 | COUNT_QTY | 开箱类 | per_qty | ¥0.10 |
| 确认型号 | CONFIRM_MODEL | 开箱类 | fixed | ¥0.10 |
| 开箱验货 | OPEN_INSPECT | 开箱类 | fixed | ¥0.00 |

## 附录 C: 运输方式清单

| 名称 | Code |
|------|------|
| 空运 | air |
| 海运 | sea |
| 海快 | sea_express |
| 陆运 | land |
| 铁路 | rail |

## 附录 D: 客户类型

| 名称 | Code |
|------|------|
| 平台客户 | platform |
| 虾皮商家 | shopee |
| 大客户 | major |
| 同行 | peer |
| 普通客户 | normal |

---

> 文档基于 I56 Framework v2.4.2 代码库分析生成。  
> 技术栈：Go 1.24+ / React 18 + TypeScript + Vite + Tailwind CSS / PostgreSQL  
> 仓库：github.com/mitboy-cyber/i56_wms
