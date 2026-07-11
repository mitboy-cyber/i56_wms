# I56 Framework vs BFT56 — 全量差距分析

> 分析日期: 2026-07-10 | 基于 BFT56 Shield 权限页面提取的 56 个 Resource

## 一、BFT56 完整资源注册表（按模块分组）

### CRM/客户管理 (15 resources)
| Resource | 中文名 | I56 状态 | API |
|----------|--------|:--:|-----|
| Client | 客户 | ✅ | GET/POST /api/v1/clients |
| ClientUser | 客户账号 | ✅ | 内嵌 |
| ClientMember | 客户会员 | ✅ | 内嵌 |
| RecipientAddress | 收件地址 | ✅ | member_address domain |
| Declarant | 申报人 | ✅ | domain exists |
| Ledger | 流水 | ✅ | client_ledger |
| Recharge | 充值记录 | ✅ | domain exists |
| ClientRoutePrice | 客户线路价 | ❌ | — |
| ClientDeliveryFee | 客户派送费 | ❌ | — |
| ClientSurcharge | 客户加收费 | ❌ | — |
| ClientStoragePrice | 客户仓储价 | ❌ | — |
| ClientServiceOverride | 附加服务覆盖 | ❌ | — |
| Statement | 对账单 | ❌ | — |
| WarehouseAuth | 仓库授权 | ❌ | — |
| CarrierAuth | 承运商授权 | ❌ | — |

### WMS/仓库管理 (10 resources)
| Resource | 中文名 | I56 状态 |
|----------|--------|:--:|
| Warehouse | 仓库 | ✅ |
| Location | 库位 | ❌ |
| LocationType | 库位类型 | ❌ |
| Zone | 区域 | ❌ |
| ZoneType | 区域类型 | ❌ |
| Container | 集装柜 | ❌ |
| InboundMachine | 入库机 | ❌ |
| Parcel | 包裹 | ✅ |
| Exception | 异常记录 | ✅ |
| EmployeeTask | 员工任务 | ❌ |

### OMS/订单 (1 resource)
| Resource | I56 状态 |
|----------|:--:|
| ConsolidationOrder | ✅ |

### 附加服务 (4 resources)
| Resource | I56 状态 |
|----------|:--:|
| ServiceOrder | ✅ |
| ServiceWorkOrder | ✅ |
| ServiceTemplate | ❌ |
| ServiceType | ✅ |

### TMS/物流 (13 resources)
| Resource | 中文名 | I56 状态 |
|----------|--------|:--:|
| AreaGroup | 区域组 | ❌ |
| CargoType | 货物类型 | ✅ |
| Carrier | 承运商 | ✅ |
| CarrierNumber | 承运商单号池 | ❌ |
| ContainerLoading | 装柜记录 | ❌ |
| Courier | 快递公司 | ✅ |
| CustomsBroker | 清关公司 | ❌ |
| CustomsPoint | 清关点 | ❌ |
| CustomsNumber | 清关单号池 | ❌ |
| Tracking | 物流追踪单 | ❌ |
| Route | 线路模板 | ✅ |
| ShippingProvider | 运输公司 | ❌ |
| TransportType | 运输方式 | ❌ |

### 工单系统 (5 resources)
| Resource | I56 状态 |
|----------|:--:|
| WorkOrder | ✅ |
| WorkOrderTemplate | ✅ |
| ProcessInstance | ❌ |
| WorkflowProcess | ✅ |
| EmployeeTask | ❌ (同WMS) |

### 系统 (7 resources)
| Resource | I56 状态 |
|----------|:--:|
| Role | ✅ |
| Employee | ✅ |
| Permission | ✅ |
| Notification | ✅ |
| PrintTemplate | ✅ |
| AuditLog | ❌ |
| ApiLog | ❌ |
| PdaVersion | ❌ |
| PdaSession | ✅ |

---

## 二、TOP 差距（按商业价值排序）

### 🔴 P0: 客户定价体系 (8资源缺失)
整个BFT56的定价分为5个维度，I56完全缺失：
- 客戶線路價：按客戶+線路設定不同的重量/體積單價
- 客戶派送費：按客戶+區域設定派送上門費用
- 客戶加收費：按客戶設定額外收費項目（報關費、操作費等）
- 客戶倉儲價：按客戶設定倉儲天數+超期費率
- 客戶附加服務覆蓋：按客戶覆蓋默認服務價格

### 🟡 P1: WMS倉庫維度 (7資源缺失)
- 庫位/庫位類型：貨架、托盤位、地面位的分級管理
- 區域/區域類型：收貨區、存儲區、揀貨區、發貨區
- 集裝櫃：海運櫃管理（櫃號、封條號、最大裝載）
- 入庫機：自動化入庫設備管理
- 員工任務：PDA派發給特定員工的任務

### 🟡 P1: TMS物流維度 (7資源缺失)
- 區域組：定義配送區域組（如北部、中部、南部）
- 承運商單號池：預先導入承運商提供的追蹤號碼池
- 裝櫃記錄：記錄每個集裝櫃裝了哪些訂單
- 清關公司+清關點+清關單號池：完整的清關管理體系
- 運輸公司+運輸方式：物流供應商管理

### 🟢 P2: 其他
- 對賬單（月結客戶的月度對帳）
- 客戶端權限（控制客戶在門戶可見的模塊）
- 倉庫作業台（倉庫現場的大屏作業界面）
- API調用日誌
- PDA版本管理
