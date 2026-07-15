# I56 Platform — Ubiquitous Language

## 仓库 (Warehouse)

| Term | Definition | Aliases to avoid |
|------|-----------|-----------------|
| **Warehouse** | 物理仓储地点，是包裹收发、操作的基本单元 | 仓库, 库房, 仓 |
| **Warehouse Zone** | 仓库内物理分区，用于库位管理 | 区域, Zone |
| **Location** | 仓库内最小存储单元，包裹上架后占据的物理位置 | 库位, Storage Location |
| **Container** | 集装柜，用于将多个包裹合并运输的物理容器 | 柜, 集装箱 |
| **Inbound** | 包裹进入仓库的流程 | 入库, 收货 |
| **Outbound** | 包裹从仓库发出的流程 | 出库, 发货 |
| **Picking** | 根据订单从库位拣选包裹的操作 | 拣货, 捡货 |
| **Stock** | 包裹在仓库中的库存状态 | 库存, 在库 |

## 包裹 (Parcel)

| Term | Definition | Aliases to avoid |
|------|-----------|-----------------|
| **Parcel** | 客户寄送到仓库的单个邮包，是最小操作单元 | 包裹, Package, 快件 |
| **Parcel No** | 快递公司分配给包裹的唯一运单号 | 快递单号, Tracking No |
| **Tracking No** | 同 Parcel No，建议统一使用 Parcel No | 运单号 |
| **Parcel Status** | 包裹在系统中的生命周期状态 | 包裹状态 |
| **Predeclare** | 客户在包裹到达仓库前预先登记的操作 | 预报, 预录 |
| **Received** | 包裹已被仓库签收 | 已签收, 已收货 |
| **On Shelf** | 包裹已完成入库上架，放入库位 | 已上架, 入库完成 |
| **Actual Weight** | 仓库实际称重的包裹重量(kg) | 实重 |
| **Volumetric Weight** | 根据包裹尺寸计算的体积重量(kg) | 体积重, 计费重 |
| **Dimensions** | 包裹的长宽高(cm) | 尺寸 |

## 订单 (Order)

| Term | Definition | Aliases to avoid |
|------|-----------|-----------------|
| **Consolidation Order** | 客户发起的"将多个已上架包裹合并运输"的订单 | 集运订单, 合单, 打包订单 |
| **Order No** | 系统为集运订单生成的唯一编号 | 订单号 |
| **Order Status** | 订单在生命周期中的状态 | 订单状态 |
| **Additional Service Order** | 客户为包裹请求附加服务的订单（如拍照、加固等） | 附加服务订单, 增值服务订单 |
| **Service Type** | 附加服务的类型定义（如拍照、拆包、加固） | 服务类型 |
| **Service Template** | 附加服务的预设模板，包含定价和操作流程 | 服务模板 |
| **Declaration** | 向海关申报的货物清单信息 | 申报, 报关 |

## 运输 (Transport)

| Term | Definition | Aliases to avoid |
|------|-----------|-----------------|
| **Route** | 从发货仓库到目的地的完整运输线路模板 | 线路, 路线, 线路模板 |
| **Transport Type** | 运输方式：海运、空运、陆运、铁路等 | 运输方式 |
| **Carrier** | 负责干线运输的承运商（如船公司、航空公司） | 承运商 |
| **Courier** | 负责末端派送的快递公司 | 快递公司 |
| **Customs Broker** | 负责清关服务的代理公司 | 清关公司 |
| **Customs Clearance Point** | 清关口岸/地点 | 清关点 |
| **Cargo Type** | 货物分类：普货、敏感货、特货等 | 货物类型, 货类 |
| **Shipping Provider** | 运输公司，泛指各类运输服务商 | 运输公司 |

## 客户 (Client)

| Term | Definition | Aliases to avoid |
|------|-----------|-----------------|
| **Client** | 使用平台的法人/组织，一个 Client 下有多个 Client User 和 Member | 客户, 企业客户 |
| **Client Account** | Client 在平台上的登录账号（即 Client User） | 客户账号 |
| **Client Member** | Client 下的最终收件人/会员，一个 Client 有多个 Member | 客户会员, 会员 |
| **Member Address** | Client Member 的收件地址 | 收件地址 |
| **Declarant** | 负责向海关申报包裹内容的自然人 | 申报人 |
| **Client Settlement** | 客户的结算方式：预充值/月结 | 结算方式 |
| **Monthly Statement** | 按月生成的对账单 | 月结对账单 |
| **Credit Limit** | 平台授予客户的授信额度 | 授信额度 |
| **Balance** | 客户在平台的账户余额 | 余额 |
| **Recharge** | 客户向账户充值 | 充值 |

## 财务 (Finance)

| Term | Definition | Aliases to avoid |
|------|-----------|-----------------|
| **Profit Report** | 按不同维度（订单/服务/客户/线路）的盈利分析报表 | 盈利报表 |
| **Ledger** | 客户账户的收支明细记录 | 余额日志, 账本 |
| **Recharge Log** | 充值操作的完整记录 | 充值记录 |
| **Route Price** | 针对特定客户、特定线路的定价 | 客户价格 |
| **Total Price** | 订单的总费用，含运费+附加服务费 | 总价 |

## 作业 (Operation)

| Term | Definition | Aliases to avoid |
|------|-----------|-----------------|
| **PDA** | 仓库操作员使用的手持终端设备 | 手持终端 |
| **PDA Session** | PDA 操作员的一次在线会话 | 在线会话 |
| **Work Order** | 系统或人工创建的仓库作业任务单 | 工单, 任务单 |
| **Work Order Template** | PDA 工单的预设模板 | 工单模板 |
| **Workflow Process** | 工单的处理流程定义 | 工单流程 |
| **Employee Task Monitor** | 员工任务分配与监控面板 | 任务监控 |
| **Exception Report** | 包裹/订单/工单中出现的异常记录 | 异常记录 |

## 系统 (System)

| Term | Definition | Aliases to avoid |
|------|-----------|-----------------|
| **Employee** | 平台方（非客户方）的员工账号 | 员工, 管理员 |
| **Role** | 权限角色的集合，关联到多个 Permission | 角色 |
| **Permission** | 单个操作权限点 | 权限 |
| **System Parameter** | 全局系统配置参数 | 系统参数 |
| **Notification** | 系统发送的通知消息 | 通知 |
| **Print Template** | 打印面单/标签的模板 | 打印模板 |

## 关系 (Relationships)

- A **Warehouse** has many **Parcels**, **Orders**, and **Employees**
- A **Parcel** belongs to one **Client** and optionally one **Client Member**
- A **Consolidation Order** contains multiple **Parcels**
- A **Parcel** may have multiple **Additional Service Orders**
- A **Client** has multiple **Client Accounts**, **Client Members**, and **Recharges**
- A **Client Member** has multiple **Member Addresses**
- A **Route** defines one **Transport Type**, supports multiple **Cargo Types**, originates from one **Warehouse**
- A **Carrier** may serve multiple **Routes**
- A **Role** contains multiple **Permissions**; an **Employee** has one or more **Roles**
- A **Work Order** follows one **Workflow Process** and may use one **Work Order Template**

## 状态模型 (Status Models)

### Parcel Status
predeclare → received → on_shelf → picking → packed → outbound

### Order Status
draft → submitted → processing → shipping → customs → delivered → completed
any → cancelled

### Settlement Method
prepaid (预充值) | monthly (月结)

## 租户模型 (Tenant Model)

- **Platform** manages multiple **Clients**
- Each **Client** has its own **Client Users**, **Members**, **Parcels**, **Orders**, **Prices**
- Client data is fully isolated — no cross-Client visibility
- Platform-level **Employees** can view all **Client** data with appropriate **Roles**
