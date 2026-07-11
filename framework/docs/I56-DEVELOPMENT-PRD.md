# I56 Framework 1.0 LTS — 开发级产品需求说明书

> **版本**: 5.0 — 合并版
> **日期**: 2026-07-10
> **来源**: BFT56 八方云仓 完整逆向工程 + 2026-07-08 PRD v4.0 + 2026-07-10 浏览器实时验证
> **目标**: I56 Framework 1.0 LTS (Go 1.24+ Modular Monolith)
> **验收标准**: 关闭 BFT56 原系统，开发团队仅凭此文档可完整重建

---

## 第一部分：产品定位与商业模式

### 1.1 产品定位

| 维度 | 分析结论 |
|------|---------|
| **行业** | 跨境电商集运物流 (Cross-border Consolidation Logistics) |
| **目标市场** | 台湾消费者在中国大陆电商平台购物后的跨境集运 |
| **核心场景** | 台湾买家在淘宝/京东/拼多多购物→寄到大陆仓库→合包→跨境运输→清关→台湾末端派送 |
| **目标客户** | 集运物流商（如 "EZ集运通"），非终端消费者 |
| **SaaS 模式** | 多租户平台，每个租户（物流公司）独立运营自己的客户群 |

### 1.2 商业模式推导

```
收费模式（基于系统数据推导）:
├── SaaS 订阅费: 租户按月/年支付平台使用费
├── 交易佣金: 按订单流水抽成（订单盈利报表体现）
├── 附加服务: 13 种增值服务明码标价
└── 客户预充值: 余额模式，平台沉淀资金

收入来源层级:
  租户（嗨购邦集团）→ 平台费
  终端客户（EZ集运通）→ 运费 + 服务费
  平台运营方 → SaaS 费 + 佣金
```

### 1.3 系统解决的问题

| 痛点 | 解决方案 |
|------|---------|
| 买家多平台购物包裹分散 | 统一寄到仓库，集运合包 |
| 跨境物流链路长不可见 | 全链路状态追踪 + Webhook 推送 |
| 清关资料复杂 | 申报人管理 + 清关公司/清关点/清关号四层自动化 |
| 费用计算不透明 | 多维计价引擎（线路×货类×税档×重量） |
| 仓库操作效率低 | PDA 扫码全流程 + 库位精确管理 |
| 财务核算困难 | 自动扣费 + 四维盈利分析 + 月结对账 |

---

## 第二部分：用户角色全画像

### 2.1 角色矩阵

| 角色 | 权限数 | 工作目标 | 核心职责 |
|------|--------|---------|---------|
| **超级管理员** | 379 项 | 平台整体运营管控 | 系统配置、全部数据查看操作、客户管理、价格配置、报表分析 |
| **仓库管理员** | 259 项 | 仓库运营效率最大化 | 包裹入库、PDA 工单、库区维护、任务分配、入库看板、异常处理 |
| **财务人员** | 361 项 | 财务数据准确、资金安全 | 充值审核、余额审查、月结对账、四维盈利分析、成本录入 |
| **客服人员** | 36 项 | 客户问题快速响应 | 订单/包裹查询、客户余额查询、基础问题解答 |
| **PDA 操作员** | 27 项 | 执行仓库操作任务 | PDA 扫码入库、拣货确认、打包确认、出库扫描 |
| **平台客户** | — | 管理自己的包裹和订单 | 预报包裹、创建订单、管理申报人、查看物流、充值、API 对接 |

### 2.2 典型工作日

**超级管理员**:
- 08:00 查看仓库看板（入库量/订单量）
- 09:00 审核客户充值申请
- 10:00 调整线路价格
- 11:00 处理异常订单
- 14:00 查看盈利报表
- 16:00 配置新客户账号

**仓库管理员**:
- 08:00 查看入库看板 → 08:30 分配拣货任务 → 09:00 处理异常包裹 → 10:00 监控仓库作业台 → 14:00 库位优化 → 17:00 确认装柜计划

---

## 第三部分：完整功能地图

### 3.1 系统总览

```
I56 Framework 企业级应用平台
│
├── 🏠 系统中心 (Framework Core)
│   ├── 员工管理 ─── CRUD + 重置密码
│   ├── 角色管理 ─── 6 角色（379/259/361/36/27 + 客户端）
│   ├── 权限管理 ─── 379 权限（含特殊数据范围权限）
│   ├── 通知管理 ─── 系统公告 + 多通道推送
│   ├── 打印模板 ─── 面单/清关单/承运商面单（复制/套用示例/设为默认）
│   ├── 任务派发参数 ─── 自动派发规则配置
│   ├── 操作日志 ─── 全量操作审计
│   ├── API 调用日志 ─── 外部调用审计
│   └── PDA 版本管理 ─── APK 下载 + 版本控制
│
├── 🏭 仓库中心 (WMS)
│   ├── 仓库列表 ─── CRUD + 仓储费配置
│   ├── 库区管理 ─── Zone CRUD + 打印 + 区域类型字典
│   ├── 库位管理 ─── Location CRUD + 打印 + 库位类型字典
│   ├── 入库机管理 ─── 设备注册 + Token 重置
│   ├── 集装柜管理 ─── Container CRUD
│   ├── 入库看板 ─── 可视化入库统计
│   ├── 仓库作业台 ─── 操作员任务面板
│   ├── PDA 在线会话 ─── 设备连接状态监控
│   ├── 员工任务监控 ─── 任务分配/取消/强制重派/优先级调整
│   ├── PDA 工单模板 ─── CRUD + 管理 + 复制
│   ├── 工单流程管理 ─── CRUD + 管理 + 复制
│   ├── 流程实例 ─── 运行中实例 + 取消 + 重启
│   ├── 工单列表 ─── 工单执行记录
│   └── 异常记录 ─── 包裹异常管理
│
├── 📦 包裹中心 (Parcel)
│   ├── 包裹列表 ─── 22 字段列表（详见 §4.2）
│   ├── 预报包裹 ─── 客户/管理员预报
│   ├── 手动认领 ─── 无人认领包裹归属匹配
│   ├── 包裹操作 ─── 打印/流转事件/标记异常/拒收
│   └── Excel 批量导入 ─── 批量预报
│
├── 📋 订单中心 (OMS)
│   ├── 集运订单 ─── 18 字段列表（详见 §4.1）
│   │   ├── 状态机: 待拣货→待装柜→已装柜→运输中→清关中→派送中→已完成
│   │   ├── 打印: 面单/承运商面单/清关单/导出申报单
│   │   ├── 操作: 设置清关公司/设置承运商单号/流水查看
│   │   └── 取消: 仅待拣货可取消
│   └── 附加服务订单 ─── 13 种服务
│       ├── 状态: 待开始→进行中→已完成/已取消
│       └── 打印
│
├── 🚚 运输中心 (TMS)
│   ├── 线路模板 ─── 35 字段含计价矩阵（详见 §4.3）
│   ├── 区域组管理 ─── 台湾北部/中部/南部等
│   ├── 货物类型 ─── 普货/特货/敏感货
│   ├── 运输方式 ─── 空运/海运/海快
│   ├── 承运商列表 ─── 末端派送商
│   ├── 承运商单号池 ─── CRUD + 作废 + 批量导入
│   ├── 运输公司 ─── 干线运输商
│   ├── 快递公司 ─── 718 家 + 单号自动识别
│   ├── 装柜记录 ─── 柜号 + 封条号
│   ├── 清关公司 ─── 清关代理
│   ├── 清关点管理 ─── 清关口岸
│   ├── 清关单号池 ─── CRUD + 作废 + 批量生成 + 批量导入
│   └── 物流追踪 ─── 节点状态更新
│
├── 👥 客户中心 (CRM)
│   ├── 客户管理 ─── 5 种客户类型 + 余额调账
│   ├── 客户账号 ─── CRUD + 重置密码
│   ├── 客户会员 ─── 127,569+ 会员 ID
│   ├── 收件地址 ─── 台湾地址管理
│   ├── 申报人 ─── 认证状态机（待认证→认证成功/失败 + 启用/停用）
│   ├── 客户充值 ─── 审核流程（确认/驳回）
│   ├── 余额日志 ─── 每笔变动记录
│   ├── 充值记录 ─── 操作日志
│   ├── 客户线路价 ─── CRUD + 从默认同步
│   ├── 客户派送费 ─── CRUD + 从默认同步
│   ├── 客户加收费 ─── CRUD + 从默认同步
│   ├── 客户仓储价 ─── CRUD + 从默认同步
│   ├── 客户附加服务覆盖 ─── CRUD + 从默认同步
│   ├── 仓库授权 ─── CRUD + 批量授权
│   ├── 承运商授权 ─── CRUD + 批量授权 + 一键绑定
│   ├── 月结对账单 ─── CRUD + 生成
│   └── 客户端权限 ─── 类型默认 + 账号覆盖
│
├── 💰 财务中心 (Finance)
│   ├── 流水 ─── CRUD
│   ├── 集运订单盈利 ─── 收入 - 成本
│   ├── 附加服务盈利 ─── 服务维度
│   ├── 客户盈利 ─── 客户维度
│   ├── 路线盈利 ─── 线路维度
│   ├── 仓库看板 ─── Global Stats / Parcel Status Stats / Online Employees / Parcel Order Flow
│   └── 驾驶舱管理 ─── 自定义仪表板
│
└── 🔌 开放平台 (API)
    ├── API 凭证 ─── Token CRUD + 重置
    ├── Webhook 投递 ─── 8 种事件 + 投递日志 + 重发
    ├── API 请求日志 ─── 审计
    └── 应用发布 ─── 版本管理
```

---

## 第四部分：字段级完整规格

### 4.1 集运订单（18 字段）

| # | 字段 | 类型 | 来源 | 规则 |
|---|------|------|------|------|
| 1 | id | INT PK | AUTO_INCREMENT | 系统生成 |
| 2 | warehouse_id | FK→warehouses | 创建时选定 | 不可修改 |
| 3 | order_no | VARCHAR(20) | 系统生成 | `yyyyMMddHHmmss` + 随机 8 位 |
| 4 | client_id | FK→clients | client 表 | 只读关联 |
| 5 | client_member_id | FK→client_members | member 表 | 只读关联，显示 member_code |
| 6 | recipient_name | VARCHAR | member_addresses | 收件人姓名 |
| 7 | tracking_numbers | TEXT | parcels.tracking_number | 多单号用 `、` 分隔 |
| 8 | route_id | FK→routes | routes 表 | 创建时选定 |
| 9 | parcel_count | INT | COUNT(parcels) | 自动计算 |
| 10 | status | ENUM | 状态机 | 10 种状态（见 §5.1） |
| 11 | container_loading_id | FK→container_loading_records | 装柜后关联 | 可空 |
| 12 | total_actual_weight | DECIMAL(10,2) | SUM(包裹实重) | 打包称重后更新 |
| 13 | total_chargeable_weight | DECIMAL(10,2) | MAX(实重, 材积) | 打包后计算 |
| 14 | total_price | DECIMAL(10,2) | 计价引擎 | 创建时计算 |
| 15 | customs_number | VARCHAR | 手动输入/分配 | 清关单号，可空 |
| 16 | carrier_tracking_no | VARCHAR | 手动/API | 承运商单号，可空 |
| 17 | remark | TEXT | 手动输入 | 可空 |
| 18 | created_at | DATETIME | NOW() | 创建时间 |

**筛选器**: 批量订单号、客户、仓库、线路、所属柜、订单状态、起止时间、客户会员编号

### 4.2 包裹（22 字段）

| # | 字段 | 类型 | 来源 | 规则 |
|---|------|------|------|------|
| 1 | id | INT PK | AUTO_INCREMENT | — |
| 2 | warehouse_id | FK→warehouses | 入库扫描 | — |
| 3 | tracking_number | VARCHAR(50) | 快递单号 | 唯一 |
| 4 | courier_id | FK→couriers | 预报/识别 | 718 家自动匹配 |
| 5 | client_id | FK→clients | 预报匹配 | 无人认领可空 |
| 6 | client_member_id | FK→client_members | 预报匹配 | 可空 |
| 7 | cargo_type_id | FK→cargo_types | 预报/手动 | 普货/特货/敏感货 |
| 8 | order_id | FK→orders | 创建订单关联 | 可空 |
| 9 | product_name | VARCHAR(200) | 预报填写 | 货品名 |
| 10 | transport_type | VARCHAR(50) | 关联线路 | 空运/海快 |
| 11 | parcel_name | VARCHAR(100) | 预报填写 | 包裹名 |
| 12 | actual_weight | DECIMAL(8,2) | PDA 称重 | 单位 kg |
| 13 | dimensions | VARCHAR(30) | PDA 测量 | `L×W×H` 格式，cm |
| 14 | status | ENUM | 状态机 | 预报→已入库→已称重→已上架→已拣货→已打包→已出库 |
| 15 | is_abnormal | BOOLEAN | 操作员标记 | — |
| 16 | has_additional_service | BOOLEAN | 服务订单关联 | — |
| 17 | location_id | FK→locations | 入库分配 | `B1-00001` 格式 |
| 18 | last_operation_at | DATETIME | 自动更新 | — |
| 19 | inbound_photos | JSON | PDA 拍照 | 可多张 |
| 20 | inbound_at | DATETIME | NOW() | 入库时间 |
| 21 | operator_id | FK→users | 当前登录 PDA 用户 | — |
| 22 | created_at | DATETIME | 预报时间 | — |

### 4.3 线路模板（35 字段）

| 分类 | 字段 | 类型 | 规则 |
|------|------|------|------|
| **基本** | is_active | BOOLEAN | 启用/禁用开关 |
| | tenant_id | FK→tenants | 租户隔离 |
| | warehouse_id | FK→warehouses | 发货仓库 |
| | transport_type_id | FK→transport_types | 空运/海运/海快 |
| | name | VARCHAR(100) | 模板名称（空運/海快/空运特货） |
| | area_group_id | FK→area_groups | 服务地区 |
| | route_code | VARCHAR(50) | 内部编码 |
| | require_declarant | BOOLEAN | 是否必须填写申报人 |
| **限制** | min_weight | DECIMAL(8,2) | 最小重量 kg |
| | max_weight | DECIMAL(8,2) | 最大重量 kg |
| | max_length | DECIMAL(8,2) | 最大长度 cm |
| | max_width | DECIMAL(8,2) | 最大宽度 cm |
| | max_height | DECIMAL(8,2) | 最大高度 cm |
| **计价** | volume_coefficient | INT | 材积系数（默认 6000） |
| | weight_rounding | DECIMAL(4,2) | 实重进位（默认 0.5） |
| | volume_rounding | DECIMAL(4,2) | 材积进位（默认 0.5） |
| | min_amount | DECIMAL(8,2) | 最低收费金额 |
| **时效** | min_delivery_days | INT | 最短天数 |
| | max_delivery_days | INT | 最长天数 |
| **装柜** | loading_time | TIME | HH:MM |
| | loading_mon | BOOLEAN | 周一 |
| | loading_tue | BOOLEAN | 周二 |
| | loading_wed | BOOLEAN | 周三 |
| | loading_thu | BOOLEAN | 周四 |
| | loading_fri | BOOLEAN | 周五 |
| | loading_sat | BOOLEAN | 周六 |
| | loading_sun | BOOLEAN | 周日 |
| **价目矩阵** | cargo_type_id | FK→cargo_types | 货物类型 |
| | tax_type | ENUM | 全包税/不包税/部分包税 |
| | default_weight_price | DECIMAL(8,2) | 默认重价（元/kg） |
| | default_volume_price | DECIMAL(8,2) | 默认材价（元/材积kg） |

### 4.4 附加服务类型（13 种）

| # | 服务名称 | 代码 | 大类 | 单价 | 价格模式 |
|---|---------|------|------|------|---------|
| 1 | 退货寄回 | RETURN_GOODS | 退货类 | ¥2.00 | fixed |
| 2 | 易碎品贴纸 | FRAGILE_STICKER | 打包类 | ¥0.10 | fixed |
| 3 | 外箱标识拍照 | OUTERBOX_PHOTO | 加固类 | ¥0.10 | fixed |
| 4 | 打木箱 | WOODEN_CRATE | 加固类 | ¥80.00 | fixed |
| 5 | 打木架 | WOODEN_FRAME | 加固类 | ¥50.00 | fixed |
| 6 | 包装气柱袋 | WRAP_AIRBAG | 加固类 | ¥5.00 | fixed |
| 7 | 包装气泡棉 | WRAP_BUBBLE | 加固类 | ¥2.00 | fixed |
| 8 | 包裹拆分 | SPLIT_PARCEL | 开箱类 | ¥2.00 | per_qty |
| 9 | 填充气柱袋 | FILL_AIRBAG | 开箱类 | ¥5.00 | fixed |
| 10 | 内容物拍照 | CONTENT_PHOTO | 开箱类 | ¥1.00 | fixed |
| 11 | 清点数量 | COUNT_QTY | 开箱类 | ¥0.10 | per_qty |
| 12 | 确认型号 | CONFIRM_MODEL | 开箱类 | ¥0.10 | fixed |
| 13 | 开箱验货 | OPEN_INSPECT | 开箱类 | ¥0.00 | fixed |

**价格模式枚举**: `fixed`（固定收费）| `per_qty`（按数量）| `per_kg`（按重量）

### 4.5 快递公司

| # | 字段 | 类型 | 规则 |
|---|------|------|------|
| 1 | id | INT PK | — |
| 2 | country_region | VARCHAR(50) | 中国大陆/台湾等 |
| 3 | name | VARCHAR(100) | 顺丰速递/中通/圆通等 |
| 4 | code | VARCHAR(20) | SF/STO/YTO 等 |

**特色功能**: 单号识别 — 输入快递单号自动识别快递公司（718 家公司编码库）

### 4.6 客户申报人（9 字段列表页）

| # | 字段 | 类型 | 规则 |
|---|------|------|------|
| 1 | type | ENUM | 个人/公司 |
| 2 | name | VARCHAR(100) | 姓名或公司名称 |
| 3 | id_number | VARCHAR(20) | 身份证号（台湾格式） |
| 4 | company_tax_id | VARCHAR(20) | 公司统一编号 |
| 5 | phone | VARCHAR(20) | 联系电话 |
| 6 | member_code | VARCHAR(20) | 归属客户会员编号 |
| 7 | auth_status | ENUM | 待认证/认证中/认证成功/认证失败 |
| 8 | is_active | BOOLEAN | 启用/停用 |
| 9 | operations | — | 查看/同步认证/停用 |

---

## 第五部分：状态机

### 5.1 订单状态机（10 态）

```
待拣货 → 待打包 → 待装柜 → 已装柜 → 运输中 → 清关中 → 派送中 → 已完成
  ↓                  ↓
已取消            异常挂起
```

| 状态 | 可执行操作 | 下一状态 |
|------|-----------|---------|
| 待拣货 | 取消 | 已取消 |
| | 拣货确认 | 待打包 |
| 待打包 | 打包确认 | 待装柜 |
| 待装柜 | 装柜确认 | 已装柜 |
| 已装柜 | 发运确认 | 运输中 |
| 运输中 | 清关开始 | 清关中 |
| 清关中 | 清关完成 | 派送中 |
| 派送中 | 签收确认 | 已完成 |
| 任意非终态 | 标记异常 | 异常挂起 |

### 5.2 包裹状态机（7 态）

```
预报 → 已入库 → 已称重 → 已上架 → 已拣货 → 已打包 → 已出库
  ↓
拒收
  ↓
异常（可标记）
```

### 5.3 申报人认证状态机

```
待认证 → 认证中 → 认证成功
                → 认证失败
启用 ⇄ 停用
```

### 5.4 充值审核状态机

```
待确认 → 已确认（到账）
       → 已驳回
```

### 5.5 附加服务订单状态机

```
待开始 → 进行中 → 已完成
       → 已取消
```

---

## 第六部分：BPMN 业务流程

### 6.1 端到端主流程

```
客户在淘宝下单 → 卖家发货 → 快递单号
                              │
┌─────────────────────────────▼──────────────────────────────┐
│ 客户操作 (Client Portal)                                    │
│  1. 登录客户端 → 2. 预报包裹(快递单号+快递公司+货品名)        │
│  3. 等待包裹到仓 → 4. 包裹已上架后创建集运订单               │
│     ├── 选择包裹(勾选) → 选择线路模板 → 选择收件地址         │
│     ├── 选择申报人(如线路要求) → 确认扣费                    │
│  5. 查看物流追踪 → 6. 签收确认                              │
└──────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────▼──────────────────────────────┐
│ 仓库操作 (PDA + Admin)                                      │
│  快递到仓→PDA扫码入库→拍照→称重→测尺寸→分配库位→已上架       │
│  订单创建→生成拣货工单→操作员PDA拣货→确认→待打包             │
│  打包操作→合箱→称重→打印面单→确认→待装柜                     │
│  装柜操作→扫描订单→输入柜号+封条号→确认→已装柜               │
└──────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────▼──────────────────────────────┐
│ 运输 (TMS)                                                  │
│  运输商提柜→物流追踪更新→运输中                              │
│  到达目的港→海关审核→清关中                                  │
│  海关放行→末端派送→派送中                                    │
│  收件人签收→已完成                                          │
└──────────────────────────────────────────────────────────┘
```

### 6.2 充值审核子流程

```
客户银行转账 → 系统提交充值申请(金额+凭证截图)
    │
    ▼
财务收到通知 → 核对银行到账记录 → 金额一致?
    ├── 是 → 审核通过 → 余额增加 → Ledger 记录 → 通知客户
    └── 否 → 审核驳回 → 注明原因 → 通知客户
```

### 6.3 清关子流程

```
订单已装柜
    │
    ▼
分配清关公司 → 指定清关点 → 生成清关号 → 生成分提单号
→ 提交申报资料(Declarant 信息 + 包裹清单 + 发票/装箱单)
    │
    ▼
海关审核 → 查验? ──是→ 补充资料或缴税 → 放行
    │
    否 → 放行 → 末端派送
```

### 6.4 定价计算流程

```
最终价格 = 线路基础价(按运输方式/货物类型/重量阶梯)
         + 客户线路价覆盖(ClientRoutePrice — 优先级: 客户专属价 > 默认价)
         + 附加服务费(按服务类型 × 数量)
         + 客户附加服务覆盖(ClientServiceOverride)
         + 承运商派送费(CarrierDeliveryFee)
         + 承运商加收费(CarrierSurcharge — 偏远/超重)
         + 仓储费(ClientStoragePrice — 按天/体积)
         - 客户折扣覆盖
```

---

## 第七部分：异常场景全覆盖

### 7.1 包裹异常

| 异常场景 | 触发条件 | 处理方式 | 恢复路径 |
|---------|---------|---------|---------|
| 快递单号重复预报 | 同一客户重复提交 | 前端拦截 + 后端校验 | 返回错误 "该单号已预报" |
| 快递未到仓 | 预报超过 N 天未入库 | 超期预警通知 | 客户确认是否取消 |
| 无法识别快递公司 | 单号不在 718 家库中 | 标记为 "未知快递" | 管理员手动指定 |
| 包裹破损 | PDA 入库拍照发现 | 标记异常 + 拍照存档 | 通知客户协商处理 |
| 包裹丢失 | 库存盘点差异 | 标记丢失 + 记录 | 赔偿流程 |
| 库位冲突 | 同一库位分配给两个包裹 | 系统禁止 | 提示 "库位已被占用" |
| 无人认领包裹 | 无预报直接到仓 | 进入认领池 | 客户通过 "认领包裹" 认领 |
| 超期未取 | 包裹超过 90 天未出库 | 仓储费累计/通知 | 客户付费或放弃 |

### 7.2 订单异常

| 异常场景 | 处理方式 |
|---------|---------|
| 创建订单时余额不足 | 前端校验 + 后端二次确认，返回 "余额不足，当前 ¥X，需 ¥Y" |
| 包裹已被其他订单占用 | 创建时校验，已占用包裹不可选 |
| 订单取消（非待拣货状态） | 返回 "只有待拣货状态的订单可取消" |
| 线路已禁用但订单使用中 | 已有订单不受影响，新订单不可选该线路 |
| 装柜超重 | 系统提示，需拆分装柜 |
| 清关被扣 | 标记订单异常，通知客户补充资料 |
| 客户价格过期 | 自动切换到默认价格，通知客户 |
| 订单并发操作 | 乐观锁 (version 字段)，版本号冲突回滚 |

### 7.3 财务异常

| 异常场景 | 处理方式 |
|---------|---------|
| 重复充值申请 | 同一转账凭证不可重复提交 |
| 充值金额与转账不符 | 财务驳回 + 注明原因 |
| 扣费时余额不足（并发） | 数据库行锁 + 事务，确保原子性 |
| 退款时余额已为负 | 记录为应收款项 |
| 月结对账单生成失败 | 自动重试 3 次 + 告警 |

### 7.4 系统异常

| 异常场景 | 处理方式 |
|---------|---------|
| PDA 离线 | 本地缓存操作，上线后同步 |
| 打印模板损坏 | 使用默认模板降级 |
| API 超时 | 3 次重试 + 指数退避 |
| Webhook 投递失败 | 3 次重试（1min/5min/15min 间隔）后标记失败 |
| 文件上传失败 | 重试 + 降级提示 |
| 数据库连接耗尽 | 连接池监控 + 自动扩容 |

---

## 第八部分：权限体系统一视图

### 8.1 RBAC 模型

```
Tenant（企业）
  └── Department（部门）
       └── Role（角色）
            └── Permission（权限）
                 └── DataScope（数据范围）
                      ├── all（全部）
                      ├── tenant（本企业）
                      ├── warehouse（指定仓库）
                      ├── dept（本部门）
                      └── self（本人）
```

### 8.2 按钮级权限矩阵

| 模块 | 操作 | 超管 | 仓管 | 财务 | 客服 | PDA员工 |
|------|------|:--:|:--:|:--:|:--:|:--:|
| 员工管理 | 查看 | ✅ | ✅ | ✅ | ✅ | ❌ |
| | 创建/编辑/删除 | ✅ | ❌ | ❌ | ❌ | ❌ |
| 集运订单 | 查看 | ✅ | ✅ | ✅ | ✅ | ❌ |
| | 取消 | ✅ | ✅ | ❌ | ❌ | ❌ |
| | 打印面单/导出申报 | ✅ | ✅ | ✅ | ❌ | ❌ |
| 包裹管理 | 查看 | ✅ | ✅ | ❌ | ✅ | ✅ |
| | 预报/入库 | ✅ | ✅ | ❌ | ❌ | ✅ |
| 客户管理 | 查看 | ✅ | ❌ | ✅ | ✅ | ❌ |
| | 创建/定价 | ✅ | ❌ | ❌ | ❌ | ❌ |
| 充值管理 | 查看 | ✅ | ❌ | ✅ | ❌ | ❌ |
| | 审核 | ✅ | ❌ | ✅ | ❌ | ❌ |
| 线路模板 | 查看 | ✅ | ✅ | ✅ | ❌ | ❌ |
| | 创建/编辑 | ✅ | ❌ | ❌ | ❌ | ❌ |
| 财务报表 | 查看 | ✅ | ❌ | ✅ | ❌ | ❌ |
| PDA 操作 | 入库/拣货/打包 | ❌ | ✅ | ❌ | ❌ | ✅ |

### 8.3 特殊数据权限

| 权限名 | 含义 |
|--------|------|
| 查看全部仓库 | 超级管理员跨仓查看所有数据（否则锁本仓） |
| 查看订单财务 | 是否可见成本/利润/盈利率 |
| 修改/删除已封柜的集装柜 | 高危操作单独授权 |
| 字典管理 增删改 | 库位类型/区域类型字典管理 |
| 通知 查看本公司全部 + 发全公司广播 | 通知管理高级权限 |
| Webhook 投递日志 重发 | 超管专用 |
| PDA 打印（面单/标签） | PDA 端打印功能 |

### 8.4 客户端权限模型

两级结构：
1. **客户类型默认权限** — 按客户类型设置默认可见功能
2. **账号覆盖权限** — 针对特定客户账号单独调整

---

## 第九部分：数据库完整设计

### 9.1 ER 总图

```
tenants (1) ──< (N) warehouses
    │                  ├── zones ──< locations
    │                  ├── inbound_machines
    │                  └── containers
    │
    ├── users (员工) ──< role_user ──> roles ──< permission_role ──> permissions
    │
    ├── clients (客户)
    │       ├── client_users (账号)
    │       ├── client_members (会员)
    │       │       ├── client_member_addresses
    │       │       └── declarants (申报人)
    │       ├── client_recharges
    │       ├── client_ledgers
    │       ├── client_route_prices
    │       └── client_panel_permissions
    │
    ├── routes (线路)
    │       ├── route_pricing_matrix
    │       ├── carrier_delivery_fees
    │       └── carrier_surcharge_fees
    │
    ├── orders (集运订单)
    │       ├── parcels (包裹)
    │       ├── parcel_service_orders
    │       ├── carrier_tracking_numbers
    │       ├── customs_numbers
    │       └── logistics_trackings
    │
    ├── container_loading_records
    │
    ├── work_orders ── work_order_templates ── workflow_processes
    │
    ├── cargo_types, transport_types, couriers (718)
    ├── customs_brokers, customs_clearance_points
    ├── carriers, shipping_providers, area_groups
    │
    ├── audit_logs, api_call_logs
    ├── notifications, print_templates
    └── pda_versions, operator_sessions
```

### 9.2 预估表清单（68 张表）

| 领域 | 表数 | 核心表 |
|------|------|--------|
| Framework Core | 8 | tenants, users, roles, permissions, role_permission, user_role, audit_logs, notifications |
| Customer | 14 | clients, client_users, client_members, member_addresses, declarants, client_ledgers, client_recharges, client_route_prices, client_delivery_fees, client_surcharges, client_storage_prices, client_service_overrides, client_statements, warehouse_authorizations |
| Warehouse | 9 | warehouses, zones, zone_types, locations, location_types, containers, inbound_machines, container_loading_records, warehouse_configs |
| Parcel | 5 | parcels, parcel_events, exception_reports, parcel_photos, parcel_dimensions |
| Order | 6 | orders, order_parcels, order_events, parcel_service_orders, parcel_service_templates, parcel_service_types |
| Transport | 12 | area_groups, cargo_types, carriers, carrier_numbers, couriers, customs_brokers, customs_points, customs_numbers, routes, route_pricing_matrix, shipping_providers, transport_types |
| Finance | 4 | invoices, payments, recharge_logs, profit_reports |
| WorkOrder | 5 | work_orders, work_order_templates, workflow_processes, process_instances, task_assignments |
| System | 5 | notifications, print_templates, pda_versions, operator_sessions, api_call_logs |

---

## 第十部分：API 契约

### 10.1 统一响应格式

```json
// 成功
{
  "data": { "id": "123", "name": "..." },
  "meta": { "total": 1234, "page_size": 20, "request_id": "req_abc" }
}

// 错误
{
  "error": {
    "code": "ORDER_NOT_FOUND",
    "message": "订单不存在",
    "details": [{"field": "order_id", "message": "找不到指定订单"}],
    "request_id": "req_abc"
  }
}
```

### 10.2 RESTful 端点规范

| 操作 | 方法 | 路径 | 响应码 |
|------|------|------|--------|
| 列表 | GET | `/api/v1/orders` | 200 + pagination meta |
| 详情 | GET | `/api/v1/orders/{id}` | 200 |
| 创建 | POST | `/api/v1/orders` | 201 |
| 更新 | PATCH | `/api/v1/orders/{id}` | 200 |
| 删除 | DELETE | `/api/v1/orders/{id}` | 204 |
| 动作 | POST | `/api/v1/orders/{id}/cancel` | 200 |

### 10.3 核心端点清单

```
# 认证
POST /api/v1/auth/login
POST /api/v1/auth/refresh
GET  /api/v1/me

# 包裹
GET    /api/v1/parcels
GET    /api/v1/parcels/{id}
POST   /api/v1/parcels
POST   /api/v1/parcels/{id}/claim    # 认领
POST   /api/v1/parcels/{id}/mark-abnormal

# 订单
GET    /api/v1/orders
GET    /api/v1/orders/{id}
POST   /api/v1/orders
POST   /api/v1/orders/{id}/cancel
POST   /api/v1/orders/{id}/print
GET    /api/v1/orders/{id}/ledger     # 流水

# 客户
GET    /api/v1/clients
POST   /api/v1/clients
GET    /api/v1/clients/{id}/members
GET    /api/v1/clients/{id}/ledgers

# 仓库
GET    /api/v1/warehouses
GET    /api/v1/warehouses/{id}/zones
GET    /api/v1/warehouses/{id}/locations

# 线路
GET    /api/v1/routes
POST   /api/v1/routes
GET    /api/v1/routes/{id}/pricing

# 财务
GET    /api/v1/finance/client-recharges
POST   /api/v1/finance/client-recharges/{id}/approve
POST   /api/v1/finance/client-recharges/{id}/reject
GET    /api/v1/finance/profit-reports

# Webhook
GET    /api/v1/webhooks/subscriptions
POST   /api/v1/webhooks/subscriptions
GET    /api/v1/webhooks/delivery-logs
POST   /api/v1/webhooks/delivery-logs/{id}/retry
```

---

## 第十一部分：客户端门户（Client Portal）

### 11.1 完整页面清单

| 页面 | URL | 功能 |
|------|-----|------|
| 主控台 | `/client` | 包裹/订单概览仪表板 |
| 收件地址 | `/client/client-member-addresses` | 台湾收件地址 CRUD |
| 客户会员 | `/client/client-members` | 会员管理 |
| 申报人 | `/client/declarants` | 海关申报人管理 |
| 我的订单 | `/client/orders` | 13 列订单列表 + 取消/下服务单 |
| 我的包裹 | `/client/parcels` | 包裹列表 + 预报 |
| 附加服务订单 | `/client/parcel-service-orders` | 服务订单管理 |
| 余额明细 | `/client/client-ledgers` | 账户流水 |
| 月结对账单 | `/client/client-statements` | 月度账单 |
| 仓库信息 | `/client/warehouses` | 仓库地址 + 联系方式 |
| 线路价格 | `/client/routes-pricings` | 线路报价查询 |
| 承运商派送价 | `/client/carrier-delivery-fees` | 末端价格 |
| 承运商加收价 | `/client/carrier-surcharge-fees` | 附加费 |
| Webhook 投递 | `/client/webhook-logs` | 事件推送日志 |
| API 凭证 | `/client/credentials` | API Token 管理 |

---

## 第十二部分：开发实施计划（I56 Framework 1.0 LTS）

### 12.1 Phase 路线图

| Phase | 内容 | 关键交付 | 工期 |
|-------|------|---------|------|
| **P0** | Go 项目骨架 + 17 Core 模块 | 可编译运行 Framework | ✅ 已完成 |
| **P1** | config/logger/errors 完整实现 + 单元测试 | 核心基础设施 | 1 周 |
| **P2** | auth JWT (ed25519) + RBAC + tenant | 认证权限系统 | 1 周 |
| **P3** | Admin Shell (Bootstrap 5 + HTMX) + 登录/角色/用户 UI | 通用管理后台 | 2 周 |
| **P4** | 仓库/库区/库位/集装柜 CRUD | WMS 基础 | 1 周 |
| **P5** | 快递公司(718 家) + 单号识别 + 货物类型 + 运输方式 + 清关 | TMS 基础 | 1 周 |
| **P6** | 线路模板 + 计价矩阵 + 客户价格 + 承运商费用 | 计价引擎 | 1.5 周 |
| **P7** | 客户 + 会员 + 申报人 + 地址 + 5 客户类型 + 面板权限 | CRM | 1.5 周 |
| **P8** | 包裹预报 + 入库 + 认领 + Excel 导入 + 状态机 | 包裹管理 | 1.5 周 |
| **P9** | 集运订单 CRUD + 状态机 + 扣费 + 打印 + 导出 | 核心订单 | 2 周 |
| **P10** | 装柜 + 承运商单号池 + 清关单号池 + 物流追踪 | 运输管理 | 1.5 周 |
| **P11** | 充值审核 + 余额 + 月结 + 四维盈利报表 | 财务系统 | 1.5 周 |
| **P12** | 附加服务(13 种) + 工单模板 + 流程 + 仓库作业台 | 增值服务 | 1.5 周 |
| **P13** | 客户端门户 15 页面 (React/HTMX) | 客户端 | 2 周 |
| **P14** | API + Webhook(8 事件) + 凭证管理 + 日志 | 开放平台 | 1 周 |
| **P15** | 打印模板 + 面单渲染 | 文档输出 | 1 周 |
| **P16** | 通知中心全渠道 + BI 报表引擎 | 增强特性 | 1.5 周 |
| **P17** | PDA 移动端 (Flutter) | 移动端 | 2 周 |
| **P18** | 测试（单元 + 集成 + E2E）+ 安全审计 | 质量保障 | 2 周 |
| **P19** | Docker Compose + K8s + CI/CD + 监控 | 生产部署 | 1 周 |

**总工期: ~27 周（4 人团队）**

---

## 第十三部分：I56 Framework 技术架构

### 13.1 技术栈

| 层 | 技术选型 |
|----|---------|
| 语言 | Go 1.24+ |
| Web 框架 | net/http (标准库) |
| 数据库 | PostgreSQL 16 (生产) / SQLite (开发) |
| 缓存 | Redis 7 |
| 消息队列 | RabbitMQ / Redis Pub/Sub |
| 对象存储 | MinIO / S3 / OSS |
| 搜索 | Elasticsearch |
| 前端 | Bootstrap 5 + HTMX + Alpine.js + Chart.js |
| 移动端 | Flutter (PDA) |
| 部署 | Docker + Docker Compose + Kubernetes + Helm |

### 13.2 Framework Core 17 模块

| # | 模块 | 文件 | 说明 |
|---|------|------|------|
| 1 | config | core/config/config.go | 多源配置 (env/file/etcd) |
| 2 | logger | core/logger/logger.go | 结构化日志 (slog) |
| 3 | errors | core/errors/errors.go | 统一错误码 + HTTP 映射 |
| 4 | response | core/response/response.go | 统一 JSON 信封 + 分页 |
| 5 | validator | core/validator/validator.go | 链式请求校验 |
| 6 | middleware | core/middleware/middleware.go | Recovery/RequestID/Logger/CORS/RateLimit |
| 7 | router | core/router/router.go | Go 1.22+ 方法路由 + 前缀 |
| 8 | tenant | core/tenant/tenant.go | 多租户解析 (Header/Subdomain) |
| 9 | auth | core/auth/auth.go | JWT 令牌管理 |
| 10 | rbac | core/rbac/rbac.go | RBAC + DataScope |
| 11 | eventbus | core/eventbus/eventbus.go | 进程内发布/订阅 |
| 12 | scheduler | core/scheduler/scheduler.go | Cron 任务调度 |
| 13 | cache | core/cache/cache.go | 多级缓存 (Memory/Redis) |
| 14 | storage | core/storage/storage.go | 统一存储 (MinIO/S3/本地) |
| 15 | notification | core/notification/notification.go | 多通道通知中心 |
| 16 | audit | core/audit/audit.go | 操作审计日志 |
| 17 | workflow | core/workflow/workflow.go | 状态机 + 流程引擎 |

### 13.3 Module 插件模式

每个业务模块遵循统一目录结构：

```
internal/modules/customer/
├── handler/           # HTTP handlers
├── service/           # Application services
├── repository/        # Data access
├── domain/            # Domain model (Entity/ValueObject/Aggregate)
├── dto/               # Request/Response DTO
├── migration/         # DB migration
├── menu/              # Admin menu registration
├── permission/        # Permission definitions
├── routes/            # Route registration
└── module.go          # Module interface (RegisterRoutes/RegisterPermissions/...)
```

---

## 第十四部分：QA 测试方案（精简）

### 关键测试用例

| TC-ID | 模块 | 场景 | 预期结果 |
|-------|------|------|---------|
| PKG-001 | 包裹 | 正常预报 | 包裹创建，状态=预报 |
| PKG-002 | 包裹 | 重复预报（同一客户+单号） | 提示 "该单号已预报" |
| PKG-003 | 包裹 | 入库扫描 | PDA 扫码→称重→分配库位→已上架 |
| PKG-004 | 包裹 | 无人认领 | 进入认领池，等待客户认领 |
| PKG-005 | 包裹 | Excel 批量导入 50 行 | 全部导入，重复跳过 |
| ORD-001 | 订单 | 正常创建（选包裹+线路+地址） | 订单创建，扣费，包裹锁定 |
| ORD-002 | 订单 | 余额不足 | 返回 "余额不足，当前¥X，需¥Y" |
| ORD-003 | 订单 | 包裹已被占用 | 包裹不可选或提示已占用 |
| ORD-004 | 订单 | 取消待拣货订单 | 包裹释放，费用退回 |
| ORD-005 | 订单 | 取消非待拣货订单 | 返回 "只有待拣货状态可取消" |
| ORD-006 | 订单 | 状态流转 | 待拣货→待打包→待装柜→已装柜 |
| FIN-001 | 财务 | 充值审核通过 | 余额增加，Ledger 记录，通知客户 |
| FIN-002 | 财务 | 重复充值（同一凭证） | 拒绝提交 |
| FIN-003 | 财务 | 月结对账自动生成 | 完整收支明细 |
| AUTH-001 | 权限 | 客服无权创建员工 | 403/按钮隐藏 |
| AUTH-002 | 权限 | 仓管无权审核充值 | 403/菜单不可见 |
| CON-001 | 并发 | 同一包裹同时被 2 个订单选择 | 只有一个成功 |

---

## 第十五部分：验收确认清单

- [x] 关闭原系统，开发团队仅凭此文档可设计数据库
- [x] 关闭原系统，开发团队仅凭此文档可开发后端
- [x] 关闭原系统，开发团队仅凭此文档可开发前端管理后台
- [x] 关闭原系统，开发团队仅凭此文档可开发客户端门户
- [x] 关闭原系统，开发团队仅凭此文档可实现全部业务流程
- [x] 关闭原系统，开发团队仅凭此文档可编写 OpenAPI 契约
- [x] 关闭原系统，开发团队仅凭此文档可编写测试用例
- [x] 关闭原系统，开发团队仅凭此文档可规划部署方案
- [x] Framework Core 17 模块代码骨架已就绪（Go 编译通过）
- [x] 定价引擎 6 层叠加模型完整定义
- [x] 5 种状态机完整定义
- [x] 68 张数据库表完整规划
- [x] BPML 核心流程完整定义
- [x] 异常场景 20+ 全覆盖
- [x] 按钮级权限矩阵完整定义

---

> **版本历史**
> - v1.0: 2026-07-08 初始逆向工程
> - v4.0: 2026-07-08 添加字段级分析、BPMN、状态机、QA、权限矩阵
> - v5.0: 2026-07-10 合并浏览器实时验证数据 + 24 隐藏资源 + I56 Framework 架构 + Go 实现路径

> **文档索引**
> - `docs/BFT56-I56-COVERAGE-MATRIX.md` — BFT56 vs I56 功能对标
> - `docs/I56-FRAMEWORK-ARCHITECTURE.md` — Framework 17 模块完整接口规格
> - `docs/DEPLOYMENT-ROADMAP.md` — 部署与路线图
