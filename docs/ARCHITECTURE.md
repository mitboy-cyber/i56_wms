# I56 Framework 1.0 LTS — 核心架构设计

## Core 模块规范

### 1. `core/auth` — 认证
```
接口:
  Authenticator.Authenticate(credential) → Token
  TokenManager.Validate(token) → Claims
  TokenManager.Refresh(token) → Token

实现:
  - JWT (默认)
  - Session (可选)
  - OAuth2 (插件)
```

### 2. `core/rbac` — 权限控制
```
模型: Tenant → Role → Permission → DataScope

DataScope:
  ALL(1)        — 全部数据
  ENTERPRISE(2) — 企业级
  WAREHOUSE(3)  — 仓库级
  DEPARTMENT(4) — 部门级
  SELF(5)       — 本人

中间件:
  RequirePermission(perm) — 检查权限
  ScopeQuery(ctx, query)  — 注入数据范围
```

### 3. `core/tenant` — 多租户
```
接口:
  TenantProvider.CurrentID(ctx) → int64
  TenantProvider.Isolate(ctx, tenantID) → context

策略:
  - SharedTable (1.0 默认)
  - SchemaPerTenant (2.0)
  - DatabasePerTenant (3.0)
```

### 4. `core/events` — 事件总线
```
接口:
  Publisher.Publish(event) → error
  Subscriber.Subscribe(eventType, handler)
  EventBus.Dispatch(event) → []error

事件格式:
  type Event struct {
      ID        string
      Type      string
      Source    string
      Timestamp time.Time
      TenantID  int64
      Payload   json.RawMessage
  }
```

### 5. `core/workflow` — 工作流引擎
```
状态机:
  StateMachine.Define(from, to, action)
  StateMachine.Transition(ctx, entity, to) → error
  StateMachine.GetTransitions(from) → []State

审批流:
  ApprovalChain.Define(steps)
  ApprovalChain.Submit(ctx, document) → error
  ApprovalChain.Approve(ctx, step, comment) → error
```

### 6. `core/scheduler` — 定时任务
```
接口:
  Job.Name() → string
  Job.Cron() → string
  Job.Execute(ctx) → error

调度器:
  Scheduler.Register(job)
  Scheduler.Trigger(name) → error
  Scheduler.Status(name) → JobStatus
```

### 7. `core/storage` — 对象存储
```
接口:
  Storage.Put(bucket, key, reader) → error
  Storage.Get(bucket, key) → (reader, error)
  Storage.Delete(bucket, key) → error
  Storage.PresignedURL(bucket, key, ttl) → url

驱动: MinIO / S3 / OSS / COS
```

### 8. `core/notification` — 通知中心
```
接口:
  Channel.Send(message) → error
  Notifier.Register(channel)
  Notifier.Notify(event) → error

渠道: Email / SMS / WeChat / Telegram / Slack / Webhook
```

## 模块依赖关系

```
                    ┌──────────┐
                    │   App    │ (业务产品)
                    └────┬─────┘
         ┌───────────────┼───────────────┐
    ┌────┴────┐    ┌────┴────┐    ┌────┴────┐
    │   WMS   │    │   OMS   │    │   CRM   │
    └────┬────┘    └────┬────┘    └────┬────┘
         │              │              │
    ┌────┴──────────────┴──────────────┴────┐
    │            Event Bus                  │
    └────────────────┬──────────────────────┘
                     │
    ┌────────────────┴──────────────────────┐
    │              CORE                     │
    │  auth  rbac  tenant  events  cache   │
    │  logger  config  validator  errors   │
    │  storage  notification  scheduler    │
    │  workflow  response  middleware      │
    └────────────────┬──────────────────────┘
                     │
    ┌────────────────┴──────────────────────┐
    │         INFRASTRUCTURE                │
    │   MySQL  Redis  RabbitMQ  MinIO       │
    │   SMTP   SMS    Elasticsearch         │
    └───────────────────────────────────────┘
```

## 向后兼容策略

### 1.0 LTS 承诺
- **API 兼容**: `/api/v1/*` 端点不向后不兼容变更
- **数据库迁移**: 提供自动迁移工具和手动SQL脚本
- **配置兼容**: 新增配置项提供默认值，不修改现有键名
- **SDK 兼容**: Go SDK 语义版本控制

### 废弃标记
```go
// Deprecated: Use core.NewEventBus() instead.
// Will be removed in I56 Framework 2.0.
func NewEvents() *Events { ... }
```
