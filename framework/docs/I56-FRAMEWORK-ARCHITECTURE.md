# I56 Framework 1.0 LTS — 架构设计文档

> 版本：1.0 LTS  
> 语言：Go 1.24+  
> 架构模式：Modular Monolith（可演进为微服务）  
> 对标系统：BFT56 八方云仓（Laravel Filament）  
> 定位：企业级应用开发平台  

---

## 一、愿景与目标

### 1.1 核心愿景

构建一套可支撑未来 10 年企业应用开发的统一平台：

```
一套 Framework → 多种业务模块 → 多个行业产品
                                    ├── 多租户 SaaS
                                    ├── 私有化部署
                                    └── 云原生部署
```

### 1.2 Framework 边界

Framework **不包含行业逻辑**，只提供基础能力。基于 Framework 可快速构建：

| 产品 | 说明 | 优先级 |
|------|------|--------|
| I56 WMS | 仓库管理系统 | P0 |
| I56 OMS | 订单管理系统 | P0 |
| I56 TMS | 运输管理系统 | P0 |
| I56 CRM | 客户关系管理 | P1 |
| I56 Finance | 财务系统 | P1 |
| I56 ERP | 企业资源计划 | P1 |
| I56 BI | 商业智能 | P2 |
| I56 OA | 办公自动化 | P2 |
| I56 MES | 制造执行系统 | P2 |
| I56 SRM | 供应商关系管理 | P2 |

---

## 二、八项设计原则

| 原则 | 说明 | 实践 |
|------|------|------|
| **Simple** | 简单优先 | 优先使用标准库，避免过度抽象 |
| **Stable** | 稳定优先 | API 向后兼容，LTS 3年支持周期 |
| **Modular** | 模块化 | 每个模块独立目录，统一接口 |
| **Convention** | 约定优于配置 | 目录结构/命名/错误码全 Framework 统一 |
| **DDD Lite** | 领域驱动设计（轻量） | Entity/ValueObject/Aggregate/DomainEvent，不过度 |
| **Event Driven** | 事件驱动 | 模块间通过 EventBus 通信，零直接依赖 |
| **Cloud Native Ready** | 云原生就绪 | 12-Factor App，Docker/K8s/Helm 部署支持 |
| **Long Term Support** | 长期支持 | 1.0 LTS 3年，安全补丁 + 数据库迁移工具 |

---

## 三、三仓库策略

```
i56-framework/     ← 企业级开发框架（Core）
i56-admin/         ← 通用管理后台（认证/RBAC/多租户/工作流/通知）
i56-apps/          ← 业务应用
  ├── i56-wms/
  ├── i56-oms/
  ├── i56-tms/
  ├── i56-crm/
  └── i56-erp/
```

### 仓库职责

| 仓库 | 职责 | 依赖 | 独立发布 |
|------|------|------|----------|
| `i56-framework` | Core 库，被所有其他仓库引用 | 无 | ✅ Go Module |
| `i56-admin` | 通用后台 UI + 认证/RBAC/多租户入口 | i56-framework | ✅ Docker Image |
| `i56-apps/i56-wms` | 仓库管理业务模块 | i56-framework + i56-admin | ✅ Docker Image |
| `i56-apps/i56-oms` | 订单管理业务模块 | i56-framework + i56-admin | ✅ Docker Image |

---

## 四、总体系统架构

```
                    Users
        ┌─────────────┼─────────────┐
    PC Browser    Tablet/PDA    Mobile    OpenAPI
        └─────────────┼─────────────┘
                      │
        ┌─────────────▼─────────────┐
        │   HTML5 Presentation      │
        │   Bootstrap 5 + HTMX      │
        │   Alpine.js + Chart.js    │
        └─────────────┬─────────────┘
                      │ HTTPS / WebSocket / SSE
        ┌─────────────▼─────────────┐
        │    I56 HTTP Gateway       │
        │    (Go 1.24+ net/http)    │
        └─────────────┬─────────────┘
                      │
        ┌─────────────▼─────────────────────────┐
        │         I56 Framework Core             │
        │  Auth │ RBAC │ Tenant │ Workflow       │
        │  EventBus │ Scheduler │ Cache          │
        │  Logger │ Config │ Validator           │
        │  Storage │ Notification │ Audit        │
        │  Report Engine │ OpenAPI               │
        └─────────────┬─────────────────────────┘
                      │
        ┌─────────────▼─────────────────────────┐
        │         Business Modules               │
        │  Customer │ Warehouse │ Parcel          │
        │  Inventory │ Order │ Purchase           │
        │  Sales │ Finance │ Report               │
        │  Dashboard │ CRM │ Work Order           │
        └─────────────┬─────────────────────────┘
                      │
        ┌─────────────▼─────────────────────────┐
        │         Infrastructure                  │
        │  MySQL │ Redis │ RabbitMQ               │
        │  MinIO │ Elasticsearch                  │
        │  SMTP │ SMS │ LINE │ Telegram           │
        └───────────────────────────────────────┘
                      │
        ┌─────────────▼─────────────────────────┐
        │  Docker │ Linux │ Kubernetes            │
        └───────────────────────────────────────┘
```

---

## 五、Framework 七层分层

```
Presentation   ← HTML, HTMX, API, WebSocket
     ↓
Application    ← 业务流程编排 (CreateOrder → CreateInvoice → PublishEvent)
     ↓
Domain         ← 纯业务 (Entity, ValueObject, Aggregate, Domain Event)
     ↓
Infrastructure ← 数据库, Redis, MQ, 文件, 第三方接口
     ↓
Framework Core ← Config, Logger, Cache, Queue, Security, RBAC, Workflow...
     ↓
Storage        ← MySQL, PostgreSQL, SQLite
     ↓
Operating System ← Linux, Darwin, Windows
```

### 各层职责

| 层 | 职责 | 依赖方向 |
|----|------|----------|
| Presentation | HTTP handler，参数绑定，响应格式化 | → Application |
| Application | 用例编排，事务管理，权限检查 | → Domain |
| Domain | 业务规则，状态机，不依赖任何外部 | 无 |
| Infrastructure | Repository 实现，外部服务适配 | ← Domain (实现接口) |
| Core | 横切关注点：日志/缓存/配置/事件 | 所有层可引用 |

---

## 六、Framework Core — 17 个核心模块

### 6.1 模块依赖图

```
config ←── (无依赖，所有模块依赖它)
logger ←── config
errors ←── config
response ←── config
validator ←── config
middleware ←── logger, errors, response
router ←── middleware
tenant ←── config, errors
auth ←── config, tenant, errors
rbac ←── auth, tenant
eventbus ←── logger, config
scheduler ←── logger, config
cache ←── config, logger
storage ←── config
notification ←── config, logger, eventbus
audit ←── logger, config, eventbus
workflow ←── config, logger, eventbus
```

### 6.2 模块规格

#### 6.2.1 config — 配置管理

```go
// 多源配置：环境变量 > 配置文件 > 默认值
type Config struct {
    App      AppConfig
    Database DatabaseConfig
    Redis    RedisConfig
    Storage  StorageConfig
    Auth     AuthConfig
    Tenant   TenantConfig
}

// 接口
func Load(opts ...Option) (*Config, error)
func (c *Config) Get(key string) any
```

#### 6.2.2 logger — 结构化日志

```go
// 基于 Go 1.21+ log/slog
type Logger interface {
    Debug(msg string, args ...any)
    Info(msg string, args ...any)
    Warn(msg string, args ...any)
    Error(msg string, args ...any)
    With(args ...any) Logger
    WithGroup(name string) Logger
}
```

#### 6.2.3 errors — 统一错误码

```go
type AppError struct {
    Code     string // "ORDER_NOT_FOUND"
    Message  string // human-readable
    HTTPStatus int  // 404
    Details  []ErrorDetail
}

// 预定义错误码
const (
    ErrNotFound          = "NOT_FOUND"
    ErrValidation        = "VALIDATION_ERROR"
    ErrUnauthorized      = "UNAUTHORIZED"
    ErrForbidden         = "FORBIDDEN"
    ErrConflict          = "CONFLICT"
    ErrInternal          = "INTERNAL_ERROR"
    ErrTenantRequired    = "TENANT_REQUIRED"
    ErrInvalidTransition = "INVALID_STATE_TRANSITION"
)
```

#### 6.2.4 response — 统一响应格式

```go
// 成功响应
type Envelope struct {
    Data  any        `json:"data,omitempty"`
    Meta  *Meta      `json:"meta,omitempty"`
    Error *APIError  `json:"error,omitempty"`
}

type Meta struct {
    Total      int64  `json:"total,omitempty"`
    PageSize   int    `json:"page_size,omitempty"`
    NextCursor string `json:"next_cursor,omitempty"`
    RequestID  string `json:"request_id"`
}

// 分页响应
type PaginatedResponse struct {
    Data       []any  `json:"data"`
    Total      int64  `json:"total"`
    Page       int    `json:"page"`
    PageSize   int    `json:"page_size"`
    TotalPages int    `json:"total_pages"`
}
```

#### 6.2.5 validator — 链式校验

```go
type Validator struct {
    errors []ErrorDetail
}

func New() *Validator
func (v *Validator) Required(field, value string) *Validator
func (v *Validator) MaxLength(field, value string, max int) *Validator
func (v *Validator) Email(field, value string) *Validator
func (v *Validator) In(field, value string, allowed []string) *Validator
func (v *Validator) Custom(fn func() error) *Validator
func (v *Validator) Valid() bool
func (v *Validator) Errors() []ErrorDetail
```

#### 6.2.6 middleware — 中间件链

```go
type Middleware func(http.Handler) http.Handler

// 内置中间件
func Recovery(logger Logger) Middleware
func RequestID() Middleware
func Logger(logger Logger) Middleware
func CORS(opts *CORSOptions) Middleware
func RateLimit(store RateLimitStore, limit int, window time.Duration) Middleware
func TenantResolver(resolver TenantResolver) Middleware
func AuthRequired(auth Authenticator) Middleware
func RBAC(enforcer Enforcer, resource, action string) Middleware

// 链式组合
func Chain(handler http.Handler, middlewares ...Middleware) http.Handler
```

#### 6.2.7 router — 路由注册（Go 1.22+ 方法路由 + 前缀支持）

```go
type Router struct {
    mux    *http.ServeMux
    prefix string
    mws    []Middleware
}

func New() *Router
func (r *Router) WithPrefix(prefix string) *Router
func (r *Router) Use(mws ...Middleware)
func (r *Router) GET(pattern string, handler http.HandlerFunc)
func (r *Router) POST(pattern string, handler http.HandlerFunc)
func (r *Router) PUT(pattern string, handler http.HandlerFunc)
func (r *Router) PATCH(pattern string, handler http.HandlerFunc)
func (r *Router) DELETE(pattern string, handler http.HandlerFunc)
func (r *Router) Handle(pattern string, handler http.Handler) // 挂载子路由
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request)

// CRITICAL: 前缀处理使用 strings.Cut 而非 strings.SplitN
// 确保 "GET /path" → "GET /prefix/path" 而非 "/prefixGET /path"
```

#### 6.2.8 tenant — 多租户解析

```go
type TenantInfo struct {
    ID     string
    Name   string
    Schema string // 用于 Schema Per Tenant
}

type TenantResolver interface {
    Resolve(r *http.Request) (*TenantInfo, error)
}

// 内置解析器
func HeaderResolver(headerName string) TenantResolver    // X-Tenant-ID
func SubdomainResolver() TenantResolver                   // tenant.example.com
func PathResolver() TenantResolver                        // /t/{tenant_id}/
func JWTResolver() TenantResolver                         // JWT claims
```

#### 6.2.9 auth — JWT 认证

```go
type TokenManager struct {
    issuer     string
    accessTTL  time.Duration
    refreshTTL time.Duration
    signKey    ed25519.PrivateKey
}

func NewTokenManager(cfg AuthConfig) *TokenManager
func (tm *TokenManager) IssueAccessToken(subject string, claims map[string]any) (string, error)
func (tm *TokenManager) IssueRefreshToken(subject string) (string, error)
func (tm *TokenManager) ValidateAccessToken(tokenStr string) (*Claims, error)
func (tm *TokenManager) RefreshAccessToken(refreshToken string) (string, string, error)

type Claims struct {
    Subject  string
    TenantID string
    Roles    []string
    Scopes   []string
}
```

#### 6.2.10 rbac — RBAC + DataScope

```go
type Enforcer struct {
    store PermissionStore
}

// 权限检查
func (e *Enforcer) Enforce(subject Subject, resource string, action string) bool

// DataScope 过滤
func (e *Enforcer) DataScope(subject Subject, resource string) DataScope

type DataScope int
const (
    ScopeAll       DataScope = iota // 全部数据
    ScopeTenant                     // 本企业
    ScopeWarehouse                  // 指定仓库
    ScopeDepartment                 // 本部门
    ScopeSelf                       // 本人
)

type Subject struct {
    UserID     string
    TenantID   string
    DeptID     string
    RoleIDs    []string
    WarehouseIDs []string
}
```

#### 6.2.11 eventbus — 事件总线

```go
type Event interface {
    EventName() string
    OccurredAt() time.Time
}

type EventBus struct {
    syncHandlers  map[string][]EventHandler
    asyncHandlers map[string][]EventHandler
}

func NewEventBus() *EventBus
func (eb *EventBus) Subscribe(eventName string, handler EventHandler, async bool)
func (eb *EventBus) Publish(ctx context.Context, event Event) error // async: fire-and-forget
func (eb *EventBus) PublishSync(ctx context.Context, event Event) error // sync: wait all

type EventHandler func(ctx context.Context, event Event) error
```

#### 6.2.12 scheduler — 定时任务

```go
type Scheduler struct {
    jobs map[string]*Job
}

type Job struct {
    Name     string
    Schedule string // cron expression
    Handler  func(ctx context.Context) error
}

func NewScheduler() *Scheduler
func (s *Scheduler) AddJob(job *Job) error
func (s *Scheduler) RemoveJob(name string)
func (s *Scheduler) Start(ctx context.Context) error
func (s *Scheduler) Stop() error

// 预置任务示例
// @every 1h  → 生成账单
// @daily 2am → 数据备份
// @every 30m → 统计数据刷新
```

#### 6.2.13 cache — 多级缓存

```go
type Cache interface {
    Get(ctx context.Context, key string) (string, error)
    Set(ctx context.Context, key string, value string, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
    Exists(ctx context.Context, key string) (bool, error)
}

// 实现
func NewMemoryCache() Cache
func NewRedisCache(client *redis.Client) Cache
func NewMultiLevelCache(l1, l2 Cache) Cache // L1: memory, L2: redis
```

#### 6.2.14 storage — 统一存储

```go
type Storage interface {
    Put(ctx context.Context, path string, data io.Reader, opts *PutOptions) error
    Get(ctx context.Context, path string) (io.ReadCloser, error)
    Delete(ctx context.Context, path string) error
    URL(ctx context.Context, path string, expiry time.Duration) (string, error)
}

// 实现
func NewMinIOStorage(cfg MinIOConfig) Storage
func NewS3Storage(cfg S3Config) Storage
func NewLocalStorage(basePath string) Storage
func NewAliyunOSSStorage(cfg OSSConfig) Storage
```

#### 6.2.15 notification — 通知中心

```go
type NotificationService struct {
    channels map[string]Channel
}

type Channel interface {
    Send(ctx context.Context, msg Message) error
}

type Message struct {
    Title    string
    Body     string
    To       []string
    Template string
    Data     map[string]any
}

// 内置渠道
func NewEmailChannel(cfg SMTPConfig) Channel
func NewSMSChannel(cfg SMSConfig) Channel
func NewLINEChannel(cfg LINEConfig) Channel
func NewTelegramChannel(cfg TelegramConfig) Channel
func NewWebhookChannel(cfg WebhookConfig) Channel
```

#### 6.2.16 audit — 审计日志

```go
type AuditService struct {
    logger  Logger
    storage AuditStorage
}

type AuditEntry struct {
    ID         string
    TenantID   string
    UserID     string
    Action     string // "order.create", "parcel.delete"
    Resource   string // "Order", "Parcel"
    ResourceID string
    Details    map[string]any
    IP         string
    UserAgent  string
    Timestamp  time.Time
}

func (a *AuditService) Log(ctx context.Context, entry AuditEntry) error
func (a *AuditService) Query(ctx context.Context, filter AuditFilter) ([]AuditEntry, error)
```

#### 6.2.17 workflow — 工作流引擎

```go
type WorkflowEngine struct {
    definitions map[string]*ProcessDefinition
    instances   ProcessStore
}

type ProcessDefinition struct {
    ID          string
    Name        string
    States      []State
    Transitions []Transition
}

type State struct {
    ID   string
    Name string
    Type StateType // Start, Task, Gateway, End
}

type Transition struct {
    From     string
    To       string
    Condition string // 条件表达式
    Action   string // 自动执行的动作
}

type ProcessInstance struct {
    ID           string
    DefinitionID string
    CurrentState string
    Variables    map[string]any
    Status       string // Running, Completed, Cancelled
}

func (e *WorkflowEngine) StartProcess(defID string, vars map[string]any) (*ProcessInstance, error)
func (e *WorkflowEngine) CompleteTask(instanceID, stateID string, vars map[string]any) error
func (e *WorkflowEngine) CancelProcess(instanceID string) error
```

---

## 七、模块体系（Plugin Module Pattern）

### 7.1 模块目录结构

```
internal/modules/
└── customer/
    ├── handler/        # HTTP handlers (presentation)
    │   ├── customer_handler.go
    │   └── customer_dto.go
    ├── service/        # Application services
    │   └── customer_service.go
    ├── repository/     # Data access
    │   └── customer_repository.go
    ├── domain/         # Domain model
    │   ├── customer.go
    │   ├── customer_member.go
    │   ├── address.go
    │   └── declarant.go
    ├── migration/      # DB migrations
    │   └── 001_create_customers.sql
    ├── menu/           # Admin menu registration
    │   └── menu.go
    ├── permission/     # Permission definitions
    │   └── permissions.go
    ├── routes/         # Route registration
    │   └── routes.go
    └── module.go       # Module interface implementation
```

### 7.2 模块接口

```go
type Module interface {
    Name() string
    Version() string
    RegisterRoutes(r *router.Router)
    RegisterPermissions() []Permission
    RegisterMenu() []MenuItem
    RegisterEventHandlers(bus *eventbus.EventBus)
    Migrate(db *sql.DB) error
}
```

### 7.3 模块注册

```go
// cmd/server/main.go
func main() {
    // ... init core ...
    
    modules := []core.Module{
        customer.New(),
        warehouse.New(),
        parcel.New(),
        order.New(),
        finance.New(),
        // ... 按需添加
    }
    
    for _, m := range modules {
        m.RegisterRoutes(appRouter)
        m.RegisterPermissions()
        m.RegisterMenu()
        m.RegisterEventHandlers(eventBus)
    }
}
```

---

## 八、事件驱动架构

### 8.1 核心事件

```go
// 包裹事件
ParcelArrived{ParcelID, WarehouseID, CourierCode, Weight, Dimensions}
ParcelWeighed{ParcelID, ActualWeight, VolumetricWeight}
ParcelStored{ParcelID, LocationID}
ParcelMarkedAbnormal{ParcelID, Reason, Description}

// 订单事件
OrderCreated{OrderID, ClientID, MemberID, RouteID, TotalWeight}
OrderPicked{OrderID, OperatorID}
OrderPacked{OrderID, ContainerID}
OrderShipped{OrderID, CarrierTrackingNo}
OrderDelivered{OrderID, SignedBy}

// 财务事件
InvoiceGenerated{InvoiceID, OrderID, Amount}
PaymentReceived{PaymentID, ClientID, Amount}
BalanceLow{ClientID, CurrentBalance, Threshold}

// 系统事件
WebhookDelivered{SubscriptionID, Event, StatusCode, Latency}
WebhookFailed{SubscriptionID, Event, Error, RetryCount}
```

### 8.2 事件消费

```go
// 订单创建后 → 财务模块生成账单
eventBus.Subscribe("order.created", func(ctx context.Context, e Event) error {
    order := e.(OrderCreated)
    return financeService.GenerateInvoice(ctx, order.OrderID)
}, true) // async

// 包裹入库后 → 通知客户
eventBus.Subscribe("parcel.arrived", func(ctx context.Context, e Event) error {
    parcel := e.(ParcelArrived)
    return notificationService.NotifyClient(ctx, parcel.ClientID, "包裹已入库")
}, true)

// Webhook 投递
eventBus.Subscribe("*", webhookDispatcher, true)
```

---

## 九、多租户模型

### 9.1 三种隔离级别

| 模式 | 隔离度 | 成本 | 适用场景 |
|------|--------|------|----------|
| Shared DB + Shared Schema | 低 | 最低 | 小微客户 SaaS |
| Schema Per Tenant | 中 | 中 | 中型企业 |
| Database Per Tenant | 高 | 高 | 大型企业/合规要求 |

### 9.2 Tenant Provider 统一抽象

```go
type TenantProvider interface {
    // 返回当前租户的 DB 连接
    DB(ctx context.Context) *sql.DB
    // 返回当前租户的 Schema（用于 Schema Per Tenant）
    Schema(ctx context.Context) string
}

// Framework 统一调用
func (s *OrderService) CreateOrder(ctx context.Context, input CreateOrderInput) error {
    db := tenant.FromContext(ctx).DB(ctx)
    // 业务代码不感知租户隔离策略
}
```

---

## 十、RBAC 完整模型

```
Tenant (企业)
  └── Department (部门)
       └── Role (角色)
            └── Permission (权限)
                 └── DataScope (数据范围)
                      ├── all       (全部)
                      ├── tenant    (本企业)
                      ├── warehouse (指定仓库)
                      ├── dept      (本部门)
                      └── self      (本人)
```

### 预置角色

| 角色 | 权限范围 | DataScope |
|------|----------|-----------|
| super_admin | 全部资源 + 全部操作 | all |
| tenant_admin | 本企业全部资源 | tenant |
| warehouse_admin | 仓库相关资源 | warehouse |
| warehouse_operator | 仓库操作资源 | warehouse |
| cs_agent | 客服资源 | tenant |
| finance | 财务资源 | tenant |
| client | 客户端门户资源 | self |

---

## 十一、开发路线图

### Phase 0：Framework Core（8 周）

| 周 | 任务 | 产出 |
|----|------|------|
| W1-2 | config + logger + errors + response | `/core/{config,logger,errors,response}/` |
| W3 | validator + middleware + router | `/core/{validator,middleware,router}/` |
| W4 | tenant + auth | `/core/{tenant,auth}/` |
| W5 | rbac + eventbus | `/core/{rbac,eventbus}/` |
| W6 | scheduler + cache + storage | `/core/{scheduler,cache,storage}/` |
| W7 | notification + audit | `/core/{notification,audit}/` |
| W8 | workflow | `/core/workflow/` + 全量测试 |

### Phase 1：Admin 基础（4 周）

| 周 | 任务 |
|----|------|
| W9-10 | 认证登录 + 角色管理 + 员工管理 UI |
| W11-12 | 通用 CRUD 组件 + 菜单系统 + 仪表板 |

### Phase 2：核心业务模块（对标 BFT56）（12 周）

| 周 | 模块 |
|----|------|
| W13-14 | Customer（客户/会员/申报人/地址） |
| W15-16 | Warehouse（仓库/库位/区域/集装柜） |
| W17-18 | Parcel（包裹/入库/称重/上架） |
| W19-20 | Order（集运订单/状态机/打印） |
| W21-22 | TMS（线路/承运商/清关/物流追踪） |
| W23-24 | Finance（充值/流水/账单/盈利报表） |

### Phase 3：增强特性（8 周）

- PDA 移动端（Flutter/React Native）
- Webhook + OpenAPI
- 通知中心全渠道
- 打印模板引擎
- BI 报表引擎

---

## 十二、LTS 生命周期

| 版本 | 定位 | 支持周期 | Go 版本 |
|------|------|----------|---------|
| 1.0 LTS | 首个长期支持版 | 3 年 (2026-2029) | Go 1.24+ |
| 1.5 | 功能增强版 | 12 个月 | Go 1.26+ |
| 2.0 LTS | 下一代架构 | 3-5 年 | Go 1.30+ |

**LTS 承诺**：
- 数据库迁移工具（向前兼容）
- API 兼容策略（废弃字段保留 2 个大版本）
- 安全补丁（关键 CVE 48h 内修复）
- 不强制升级 Core 版本（业务模块可独立版本）

---

## 十三、最终目录结构

```
i56-framework/
├── cmd/
│   ├── server/main.go         # HTTP server
│   ├── migrate/main.go        # DB migration
│   ├── worker/main.go         # Background job worker
│   └── cli/main.go            # CLI tools
├── configs/                    # 配置文件模板
│   ├── config.yaml
│   └── config.prod.yaml
├── core/                       # Framework Core (17 modules)
│   ├── config/
│   ├── logger/
│   ├── errors/
│   ├── response/
│   ├── validator/
│   ├── middleware/
│   ├── router/
│   ├── tenant/
│   ├── auth/
│   ├── rbac/
│   ├── eventbus/
│   ├── scheduler/
│   ├── cache/
│   ├── storage/
│   ├── notification/
│   ├── audit/
│   └── workflow/
├── internal/
│   ├── modules/                # 业务模块
│   │   ├── customer/
│   │   ├── warehouse/
│   │   ├── parcel/
│   │   ├── order/
│   │   ├── finance/
│   │   ├── transport/
│   │   ├── workorder/
│   │   └── report/
│   ├── plugins/                # 第三方集成
│   └── services/               # 跨模块服务
├── pkg/                        # 可导出公共库
│   └── i56sdk/                 # Go SDK
├── sdk/                        # 多语言 SDK
│   ├── python/
│   ├── javascript/
│   └── java/
├── templates/                  # HTML 模板
├── static/                     # 静态资源
├── themes/                     # UI 主题
│   ├── default/
│   ├── dark/
│   └── enterprise/
├── migrations/                 # DB 迁移脚本
├── scripts/                    # 运维脚本
├── deployments/
│   ├── docker/Dockerfile
│   ├── compose/docker-compose.yml
│   ├── helm/
│   └── kubernetes/
├── docs/
│   ├── architecture.md
│   ├── api/
│   └── guides/
├── tests/
│   ├── unit/
│   ├── integration/
│   └── e2e/
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

*文档版本：1.0 | 日期：2026-07-10 | 作者：I56 Framework Team*
