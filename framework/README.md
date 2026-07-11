# I56 Framework — Enterprise Application Development Platform

> Go 1.24+ Modular Monolith — 可演进为微服务。

## 定位

`i56-framework` 是平台核心。提供认证、权限、多租户、工作流、事件总线、缓存、日志等基础能力。不包含任何行业逻辑。

```
i56-framework/        ← 这个仓库 (平台能力)
    ↓ 被依赖
i56-modules/          (业务模块)
    ↓ 被装配
i56-apps/             (产品交付)
```

## 版本

| 版本 | 状态 | 支持周期 |
|------|------|---------|
| 1.0 LTS | ✅ 当前 | 3 年 |

## 架构

```
┌──────────────────────────────────────────────────┐
│  API Gateway  │  Gin — JWT / RateLimit / CORS    │
│  Real-time    │  SSE + WebSocket                 │
├──────────────────────────────────────────────────┤
│  Application  │  Go Modular Monolith (7 layers) │
├──────────────────────────────────────────────────┤
│  Persistence  │  PostgreSQL 16 + Redis 7 + MinIO │
├──────────────────────────────────────────────────┤
│  Deployment   │  Docker + Compose + Kubernetes   │
└──────────────────────────────────────────────────┘
```

## 7 层架构

```
Presentation   → HTML5 + Web Components + Go templates
Application    → Service layer (CreateOrder / PreDeclare)
Domain         → Entities, ValueObjects, Aggregates, DomainEvents
Infrastructure → Repositories, External APIs, File Storage
Framework Core → Auth, Config, Cache, Queue, Logger, Scheduler
Storage        → PostgreSQL, Redis, MinIO
OS             → Linux, Docker, Kubernetes
```

## Core 模块

```
core/gateway/      — Gin 引擎、中间件、WS Hub、SSE Hub
core/auth/         — JWT Ed25519
core/config/       — 环境 / 文件配置
core/cache/        — Redis 缓存抽象
core/queue/        — Redis Streams 消息队列
core/storage/      — MinIO / S3 文件存储
core/scheduler/    — 定时任务
core/workflow/     — 审批引擎
core/report/       — 报表引擎
core/router/       — 自定义路由（legacy）
core/middleware/   — 中间件
core/logger/       — 结构化日志
core/eventbus/     — 事件总线
core/plugin/       — 插件系统
```

## API

```
/api/v1/health             — 健康检查
/api/v1/admin/crud         — 管理后台 CRUD
/api/v1/weight-records     — 称重记录
/ws                         — WebSocket
/sse?channel=xxx            — SSE 流
```

## 部署

```bash
# 开发
docker compose -f deployments/compose/docker-compose.yml up -d

# 生产
kubectl apply -f deployments/kubernetes/deployment.yaml
```

## 许可证

MIT
