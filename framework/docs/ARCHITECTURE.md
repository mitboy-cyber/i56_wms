# I56 Framework 1.0 LTS — Architecture

## Stack

```
┌──────────────────────────────────────────────────────────────────┐
│  Presentation   │ HTML5 + ES Modules + Web Components            │
│  (BDL 1.0)      │ I56 Design Language — no Vue/React            │
├──────────────────────────────────────────────────────────────────┤
│  Interaction    │ Command Palette (Ctrl+K) + AI Command Bar      │
│                 │ Right-click menus + Keyboard shortcuts         │
├──────────────────────────────────────────────────────────────────┤
│  Real-time      │ SSE — AI streaming / dashboard live feed       │
│                 │ WebSocket — logistics tracking / PDA sync      │
├──────────────────────────────────────────────────────────────────┤
│  API Gateway    │ Gin (primary) or Fiber (alternative)           │
│                 │ JWT / RateLimit / CORS / Logging               │
├──────────────────────────────────────────────────────────────────┤
│  Application    │ Go Modular Monolith → evolvable to microservices│
│                 │ 7 layers: Presentation → App → Domain → Infra  │
├──────────────────────────────────────────────────────────────────┤
│  Persistence    │ PostgreSQL 16 — core business data             │
│                 │ Redis 7 — cache / sessions / message queue     │
│                 │ MinIO — file / document storage                │
├──────────────────────────────────────────────────────────────────┤
│  Deployment     │ Docker + Compose (dev) / Kubernetes (prod)     │
│                 │ CDN for static assets                          │
└──────────────────────────────────────────────────────────────────┘
```

## Design Language — I56 BDL 1.0

```
Principles:
  Minimal    — 极简布局，减少视觉噪音
  Consistent — 全局统一间距、字体、图标
  AI Native  — AI 助手是一级能力
  KB First   — 全局快捷键与 Command Palette
  Fast       — 首屏快、动画轻、响应快
  Modular    — 可组合的卡片和模块
  Responsive — HTML5 适配 PC/平板/移动端
  Accessible — 高对比度、键盘可操作、语义化

Visual:
  Dark (#0a0a0f) / Light (#ffffff)
  Brand: #6366f1 (indigo-500)
  Font: system-ui | Monospace: JetBrains Mono
  Radius: 6/8/12px | Spacing: 4-64 scale
  Animation: 150ms ease all transitions
```

## Component Library

```
<i56-button>     <i56-card>       <i56-table>
<i56-input>      <i56-select>     <i56-form-group>
<i56-modal>      <i56-toast>      <i56-badge>
<i56-tabs>       <i56-avatar>     <i56-spinner>
<i56-timeline>   <i56-command-palette>
```

## Modules

```
core/gateway/    — Gin engine, middleware, WS Hub, SSE Hub
core/auth/       — JWT Ed25519
core/config/     — env / file config
core/logger/     — structured logging
core/cache/      — Redis cache
core/queue/      — Redis streams (message bus)
core/storage/    — MinIO / S3
core/scheduler/  — cron jobs
core/workflow/   — approval engine
core/report/     — report engine

internal/modules/
  customer/      — clients, members, addresses
  order/         — orders, service orders
  parcel/        — parcels, state machine
  warehouse/     — inventory, bin locations
  transport/     — routes, carriers, pricing
  finance/       — billing, ledger, statements
  rbac/          — roles, permissions, users
  weight/        — weight records
  pda/           — PDA operations
  webhook/       — webhook delivery

static/
  css/i56-bdl.css       — Design tokens (BDL 1.0)
  js/i56-theme.js       — Dark/Light toggle + brand
  js/i56-components.js  — Web Component library
  js/i56-command.js     — Command Palette (Ctrl+K)

templates/admin/
  base_new.html         — New BDL layout
  permissions_new.html  — RBAC permissions (sample)
```

## 7 Layers

```
Presentation  → HTML5 + Web Components + Go templates
Application   → Service layer (CreateOrder / PreDeclare / etc.)
Domain        → Entities, Value Objects, Aggregates, Domain Events
Infrastructure→ Repositories, External APIs, File Storage
Framework Core→ Auth, Config, Cache, Queue, Logger, Scheduler
Storage       → PostgreSQL, Redis, MinIO
OS            → Linux, Docker, Kubernetes
```

## API Design

```
/api/v1/health            — Health check (public)
/api/v1/weight-records    — Weight CRUD
/api/v1/admin/crud        — Admin CRUD (RBAC entities)
/ws                        — WebSocket (logistics push)
/sse?channel=xxx           — SSE (AI streaming)
```

## Deployment

```bash
# Dev (Docker Compose)
docker compose -f deployments/compose/docker-compose.yml up -d

# Prod (Kubernetes)
kubectl apply -f deployments/kubernetes/deployment.yaml
kubectl scale deployment i56-server --replicas=4
```

## File Structure

```
i56-framework/
├── cmd/server/           — Entry points (main.go, main_gin.go)
├── core/                 — Framework core
│   ├── gateway/          — Gin Gateway, middleware, WS, SSE
│   ├── auth/             — JWT, permissions
│   ├── config/           — Configuration
│   ├── cache/            — Redis cache
│   ├── queue/            — Message queue
│   ├── storage/          — File storage
│   ├── scheduler/        — Cron scheduler
│   ├── workflow/         — Approval engine
│   ├── report/           — Report engine
│   └── router/           — Custom router (legacy)
├── internal/modules/     — Business modules
├── static/               — Frontend assets
│   ├── css/i56-bdl.css
│   └── js/
├── templates/            — Go HTML templates
├── deployments/          — Docker, Compose, K8s
│   ├── docker/Dockerfile
│   ├── compose/docker-compose.yml
│   └── kubernetes/deployment.yaml
├── docs/ARCHITECTURE.md
├── main.go
└── go.mod
```

## Version Lifecycle

| Version | Type | Support |
|---------|------|---------|
| 1.0 LTS | Long Term Support | 3 years |
| 1.5 | Feature release | 12 months |
| 2.0 LTS | Next-gen architecture | 3-5 years |
