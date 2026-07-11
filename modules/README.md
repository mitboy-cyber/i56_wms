# I56 Modules — Official Business Modules

> 可复用的业务模块库。

## 定位

`i56-modules` 提供官方业务模块，每个模块独立、可插拔、可组合。

```
i56-framework/  (平台能力)
        ↓ 依赖
i56-modules/    ← 这个仓库 (业务模块)
        ↓ 被装配
i56-apps/       (产品交付)
```

## 模块清单

```
i56-modules/
├── customer/       — 客户管理 (clients, members, declarants, addresses, ledger)
├── order/          — 订单管理 (consolidation orders, service orders)
├── parcel/         — 包裹管理 (state machine: pre_declared→received→stored→picked→packed→shipped)
├── warehouse/      — 仓库管理 (warehouses, bin locations, inventory)
├── transport/      — 物流管理 (routes, carriers, couriers, cargo types)
├── finance/        — 财务管理 (billing, invoices, statements, profit reports)
├── rbac/           — 权限管理 (permissions, roles, users, client permissions)
├── weight/         — 称重管理 (weight records, dashboard)
├── pda/            — PDA 操作 (receive, putaway, pick, pack, load)
├── webhook/        — Webhook 管理
├── report/         — 报表引擎
├── notification/   — 通知中心 (email, SMS, webhook)
└── openapi/        — OpenAPI / SDK
```

## 模块规范

每个模块统一结构：

```
module/
├── controller/     — HTTP handlers
├── service/        — business logic
├── repository/     — data access
├── domain/         — entities, value objects, aggregates
├── dto/            — data transfer objects
├── validator/      — request validation
├── migration/      — database migrations
├── menu/           — admin menu registration
├── permission/     — permission registration
├── routes/         — route registration
└── module.go       — module definition & registration
```

## 使用

```go
import (
    "github.com/i56/modules/order"
    "github.com/i56/modules/parcel"
    "github.com/i56/modules/rbac"
)

app := framework.New()
app.RegisterModule(order.Module)
app.RegisterModule(parcel.Module)
app.RegisterModule(rbac.Module)
app.Run()
```

## 许可证

MIT
