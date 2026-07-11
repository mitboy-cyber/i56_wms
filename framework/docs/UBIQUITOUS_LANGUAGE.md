# I56 Framework — Ubiquitous Language（领域术语表）

> 此术语表与 `I56-DEVELOPMENT-PRD.md` 配套使用。所有开发、讨论、API 命名必须统一使用此表中的规范术语。

---

## 核心业务实体

| 术语（中文） | 术语（英文） | 定义 | 禁止别名 |
|-------------|-------------|------|---------|
| **租户** | Tenant | 平台上的一个企业客户（如 "嗨购邦集团"、"EZ集运通"） | 公司、企业、organization |
| **客户** | Client | 租户下的业务客户（如 "EZ集运通-运营"） | 买家、用户 |
| **客户账号** | ClientUser | 客户登录系统的账号 | 账号、account、login |
| **客户会员** | ClientMember | 终端消费者（台湾收件人），属于客户 | 会员、收件人 |
| **申报人** | Declarant | 海关报关身份信息（个人身份证/公司统编） | 报关人、importer |
| **收件地址** | MemberAddress | 客户会员的台湾收货地址 | 地址、delivery address |
| **员工** | User | 后台操作人员（超管/仓管/财务/客服/PDA操作员） | 账号、admin |
| **角色** | Role | 权限的集合（超管/仓管/财务/客服/员工） | 权限组 |

## 仓库域

| 术语（中文） | 术语（英文） | 定义 |
|-------------|-------------|------|
| **仓库** | Warehouse | 物理仓储地点（如 "厦门仓"） |
| **库区** | Zone | 仓库内的功能分区（收货区/存储区/拣货区/打包区/发货区/异常区） |
| **库位** | Location | 库区内的精确存储位置（格式 `B1-00001`） |
| **集装柜** | Container | 用于跨境运输的集装箱 |
| **入库机** | InboundMachine | PDA 扫描设备 |
| **装柜记录** | ContainerLoadingRecord | 订单装入集装箱的操作记录 |

## 包裹域

| 术语（中文） | 术语（英文） | 定义 |
|-------------|-------------|------|
| **包裹** | Parcel | 一个快递包裹，有唯一快递单号 |
| **预报** | Pre-declare | 客户提前通知即将到达的包裹信息 |
| **认领** | Claim | 将无主包裹匹配到客户 |
| **实重** | ActualWeight | 包裹的实际物理重量（kg） |
| **材积重** | VolumetricWeight | `长×宽×高÷材积系数` 的计算重量 |
| **计费重** | ChargeableWeight | `MAX(实重, 材积重)`，用于计费 |
| **货物类型** | CargoType | 普货/特货/敏感货 |
| **快递公司** | Courier | 国内快递承运方（顺丰/中通等，718 家） |
| **快递单号** | TrackingNumber | 快递包裹的唯一追踪号 |

## 订单域

| 术语（中文） | 术语（英文） | 定义 |
|-------------|-------------|------|
| **集运订单** | Order / ConsolidationOrder | 客户创建的多包裹合并运输订单 |
| **线路** | Route | 运输线路模板（含计价矩阵） |
| **运输方式** | TransportType | 空运/海运/海快 |
| **附加服务** | ParcelService | 增值服务（加固/拍照/拆分/退货等 13 种） |
| **附加服务订单** | ParcelServiceOrder | 客户下的增值服务请求 |
| **面单** | Waybill / ShippingLabel | 贴在包裹上的运输标签 |

## 运输域

| 术语（中文） | 术语（英文） | 定义 |
|-------------|-------------|------|
| **承运商** | Carrier | 末端派送承运商（如 "新竹物流"） |
| **运输公司** | ShippingProvider | 干线运输商（船公司/航司） |
| **清关公司** | CustomsBroker | 报关代理 |
| **清关点** | CustomsClearancePoint | 清关口岸 |
| **承运商单号** | CarrierTrackingNumber | 承运商分配的追踪号 |
| **清关单号** | CustomsNumber | 海关申报单号 |

## 财务域

| 术语（中文） | 术语（英文） | 定义 |
|-------------|-------------|------|
| **充值** | Recharge | 客户向平台账户充值 |
| **余额** | Balance | 客户账户可用金额 |
| **流水** | Ledger | 客户账户的每笔收支记录 |
| **月结对账单** | MonthlyStatement | 每月自动生成的收支汇总 |
| **盈利报表** | ProfitReport | 按订单/服务/客户/路线维度分析的收入-成本 |

## 状态枚举

### 订单状态 (OrderStatus)

| 状态值 | 中文 | 说明 |
|--------|------|------|
| `pending_picking` | 待拣货 | 已生成，等待仓库拣选 |
| `picking` | 拣货中 | 操作员正在拣货 |
| `pending_packing` | 待打包 | 已拣货，等待打包 |
| `pending_loading` | 待装柜 | 已打包，等待装入集装箱 |
| `loaded` | 已装柜 | 已装入集装箱 |
| `in_transit` | 运输中 | 跨境运输中 |
| `customs_clearance` | 清关中 | 海关审核中 |
| `out_for_delivery` | 派送中 | 末端派送中 |
| `completed` | 已完成 | 收件人签收 |
| `cancelled` | 已取消 | 订单已取消 |
| `exception` | 异常挂起 | 订单异常暂停 |

### 包裹状态 (ParcelStatus)

| 状态值 | 中文 | 说明 |
|--------|------|------|
| `pre_declared` | 预报 | 客户已预报 |
| `received` | 已入库 | 仓库已签收 |
| `weighed` | 已称重 | 已完成称重和测量 |
| `stored` | 已上架 | 已分配库位 |
| `picked` | 已拣货 | 已被订单拣选 |
| `packed` | 已打包 | 已合箱打包 |
| `shipped` | 已出库 | 已装柜发运 |
| `rejected` | 拒收 | 仓库拒收 |
| `abnormal` | 异常 | 标记为异常 |

### 申报人认证状态 (AuthStatus)

| 状态值 | 中文 |
|--------|------|
| `pending` | 待认证 |
| `verifying` | 认证中 |
| `verified` | 认证成功 |
| `failed` | 认证失败 |

### 充值状态 (RechargeStatus)

| 状态值 | 中文 |
|--------|------|
| `pending` | 待确认 |
| `confirmed` | 已确认 |
| `rejected` | 已驳回 |

---

## 关键关系

- 一个 **Tenant** 拥有多个 **Warehouse**
- 一个 **Warehouse** 包含多个 **Zone**，每个 **Zone** 包含多个 **Location**
- 一个 **Client** 拥有多个 **ClientMember**，每个 **ClientMember** 有多个 **MemberAddress** 和 **Declarant**
- 一个 **Order** 包含多个 **Parcel**（many-to-many 通过 `order_parcels`）
- 一个 **Order** 可选关联一个 **ContainerLoadingRecord**
- 一个 **Order** 通过 **Route** 确定运输方式和计价规则
- 一个 **Parcel** 属于一个 **Courier**，可能属于一个 **Order**
- 一个 **ParcelServiceOrder** 可应用于一个 **Parcel** 或一个 **Order**

---

## 示例对话

> **Dev:** "用户在创建 **Order** 时，如果 **Balance** 不足，系统应该怎么做？"

> **Domain Expert:** "前端和后端都要校验。后端在扣费时使用数据库行锁确保原子性——如果 **Balance** 不够覆盖运费加上所有 **ParcelService** 费用，返回特定错误码让前端提示用户去 **Recharge**。"

> **Dev:** "一个 **Parcel** 可以同时属于两个 **Order** 吗？"

> **Domain Expert:** "绝对不行。**Parcel** 被 **Order** 关联后，状态变为 `picked`，其他 **Order** 不应该能看到它。创建 **Order** 时需要加乐观锁或行锁防止并发抢占。"

> **Dev:** "**ContainerLoadingRecord** 和 **Order** 的关系是？"

> **Domain Expert:** "一对多——一个 **Container** 可以装多个 **Order**。装柜后 **Order** 状态从 `pending_loading` 变为 `loaded`。"

> **Dev:** "**Route** 的计价什么时候生效？"

> **Domain Expert:** "创建 **Order** 时快照当前 **Route** 的计价矩阵。如果之后 **Route** 价格变动，已创建的 **Order** 不受影响。但是如果客户有专属的 **ClientRoutePrice**，那优先级高于 **Route** 的默认价。"

---

*此术语表与 `I56-DEVELOPMENT-PRD.md` v5.0 配套使用*
