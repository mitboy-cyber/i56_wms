# I56 WMS — API 规范 (OpenAPI 3.0)

> 基于当前 v114 实现，记录所有可用 API 端点

## 认证

所有管理 API 需要 `admin_session` Cookie：
```
POST /admin/login → Set-Cookie: admin_session=xxx; Path=/admin
```

所有客户端 API 需要 `client_token` Cookie：
```
POST /client/login → Set-Cookie: client_token=xxx; Path=/client
```

## 统一响应格式

```json
// 列表
GET /admin/api/{resource} → Array<Object>

// 创建
POST /admin/api/{resource} → { "id": 1 }

// 更新
PUT /admin/api/{resource} → { "id": 1 }

// 删除
DELETE /admin/api/{resource} → 200
```

## Admin API 端点 (48)

### 仪表盘
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/admin/api/dashboard/stats` | 仪表盘统计 |
| GET | `/admin/api/dashboard/orders` | 最近订单 |
| GET | `/admin/api/dashboard/parcels` | 包裹状态分布 |

### 订单管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/admin/api/orders` | 订单列表 |
| POST | `/admin/api/orders` | 创建订单 |
| PUT | `/admin/api/orders` | 更新订单 |
| DELETE | `/admin/api/orders` | 删除订单 |

### 仓库管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/admin/api/parcels` | 包裹列表 |
| GET | `/admin/api/warehouses` | 仓库列表 |
| GET | `/admin/api/containers` | 集装柜 |
| GET | `/admin/api/container-loadings` | 装箱记录 |
| GET | `/admin/api/service-templates` | 服务模板 |
| GET | `/admin/api/service-types` | 服务类型 |
| GET | `/admin/api/service-order-records` | 服务订单(Store) |
| GET | `/admin/api/service-workorders` | 工单列表 |
| GET | `/admin/api/exceptions` | 异常记录 |
| GET | `/admin/api/exception-reports` | 异常报告 |
| GET | `/admin/api/ai-exceptions` | AI异常 |
| GET | `/admin/api/pda-sessions` | PDA会话 |
| GET | `/admin/api/pda-workorder-templates` | PDA模板 |
| GET | `/admin/api/workflow` | 工作流 |
| GET | `/admin/api/warehouse-board` | 仓库看板 |
| GET | `/admin/api/inbound-board` | 入库看板 |
| GET | `/admin/api/warehouse-console` | 仓库作业台 |

### 财务报表
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/admin/api/report/order-profit` | 订单盈利 |
| GET | `/admin/api/report/service-profit` | 服务盈利 |
| GET | `/admin/api/report/client-profit` | 客户盈利 |
| GET | `/admin/api/report/route-profit` | 路线盈利 |
| GET | `/admin/api/monthly-statements` | 对账单 |

### 物流管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/admin/api/couriers` | 快递公司 |
| GET | `/admin/api/shipping-providers` | 承运商 |
| GET | `/admin/api/area-groups` | 区域组 |
| GET | `/admin/api/route-templates` | 路线模板 |
| GET | `/admin/api/cargo-types` | 货物类型 |
| GET | `/admin/api/transport-modes` | 运输方式 |
| GET | `/admin/api/customs-brokers` | 清关公司 |
| GET | `/admin/api/customs-points` | 清关点 |
| GET | `/admin/api/logistics-tracking` | 物流追踪 |

### 客户管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/admin/api/client-accounts` | 客户列表 |
| GET | `/admin/api/client-members` | 会员 |
| GET | `/admin/api/client-ledgers` | 账本 |
| GET | `/admin/api/client-recharges` | 充值记录 |
| GET | `/admin/api/client-pricing` | 客户报价 |
| GET | `/admin/api/client-permissions` | 客户权限(自动) |
| GET | `/admin/api/customer-declarants` | 申报人 |
| GET | `/admin/api/customer-addresses` | 地址簿 |

### 系统管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/admin/api/roles` | 角色 |
| GET | `/admin/api/notifications` | 通知 |
| GET | `/admin/api/printers` | 打印机 |
| GET | `/admin/api/print-templates` | 打印模板 |
| GET | `/admin/api/storage` | 存储配置 |
| GET | `/admin/api/system/params` | 系统参数 |
| GET | `/admin/api/system/brand` | 品牌设置 |
| GET | `/admin/api/system/scheduler` | 定时任务 |
| GET | `/admin/api/system/audit-logs` | 审计日志 |
| GET | `/admin/api/system/reports` | 报表 |
| GET | `/admin/api/system/notification-channels` | 通知渠道 |

### 计费管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/admin/api/pricing-services` | 计价服务 |
| GET | `/admin/api/pricing-routes` | 路线计费 |
| GET | `/admin/api/pricing-delivery` | 派送计费 |
| GET | `/admin/api/pricing-surcharges` | 加收计费 |

## Client API 端点 (13)

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/client/api/me` | 认证检查 |
| GET | `/client/api/dashboard` | 客户端仪表盘 |
| GET | `/client/api/parcels` | 我的包裹 |
| GET | `/client/api/orders` | 我的订单 |
| GET | `/client/api/declarants` | 申报人 |
| GET | `/client/api/members` | 会员 |
| GET | `/client/api/addresses` | 地址 |
| GET | `/client/api/ledger` | 账本 |
| GET | `/client/api/warehouses` | 仓库 |
| GET | `/client/api/route-prices` | 路线价格 |
| GET | `/client/api/couriers` | 快递 |
| GET | `/client/api/service-orders` | 服务订单 |
| GET | `/client/api/delivery-fees` | 派送费 |
