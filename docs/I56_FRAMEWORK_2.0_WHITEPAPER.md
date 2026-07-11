# I56 Framework 2.0 LTS — AI Native Architecture

```
版本：2.0 LTS
语言：Go 1.24+
架构：Modular Monolith × AI Runtime
定位：AI-Native Enterprise Application Platform
```

---

## 一、2.0 范式跃迁

在 1.x 时代，I56 是一个**确定性**的企业应用平台：输入→处理→输出。

在 2.0 时代，I56 进化为**非确定性协同引擎**：

```
Context + Intention → Planning → Tool Execution → Human-in-the-Loop
```

### 架构对比

| 维度 | 1.x (Business Platform) | 2.0 (AI Infrastructure) |
|:--|:--|:--|
| 核心抽象 | HTTP Router → Controller → Service | AI Runtime → Agent → Tool |
| 数据流 | 请求/响应 | 流式 SSE + 事件驱动 |
| 权限模型 | RBAC | RBAC × AI Guardrail × PII Mask |
| 路由决策 | 固定路由表 | Model Router (成本/速度/复杂度) |
| 前端入口 | URL → Page | Command Bar (Ctrl+K) + URL |
| 工作流 | 硬编码状态机 | AI Planner → 动态 DAG |
| 集成协议 | REST/OpenAPI | OpenAPI + MCP (Model Context Protocol) |
| 用户界面 | 表单+表格 | Workspace (人+Agent协作空间) |

---

## 二、新六大支柱 (Six Pillars)

```
┌──────────────────────────────────────────────────────┐
│                  Human User & Agent                   │
├──────────────────────────────────────────────────────┤
│  BDL 2.0 (AI Workspace / Command Bar / Chat Panel)   │
├──────────────────────────────────────────────────────┤
│              framework/ai (独立层)                     │
│        12 Modules: Gateway Router Context Memory       │
│        Knowledge Tools MCP Workflow Agent              │
│        Security Prompt RAG                             │
├──────────────────────────────────────────────────────┤
│              Framework Core (RBAC/Audit/Event)         │
├──────────────────────────────────────────────────────┤
│            Business Tools (WMS/OMS/CRM DDD)           │
└──────────────────────────────────────────────────────┘
```

### Pillar 1: AI First

系统启动时，第一行初始化代码是 AI Runtime，而非业务路由。

```go
func main() {
    aiRuntime := ai.NewAIService(...)  // ① AI Runtime 先于一切
    router := core.NewRouter()          // ② 传统路由
    aiRuntime.Mount(router)             // ③ AI 感知所有路由
}
```

所有 HTTP 请求、数据库变更、事件总线消息默认暴露给 AI Runtime 监听。

### Pillar 2: Human + Agent Collaboration

```
订单详情页：
  ┌─────────────────────────────────────┐
  │ 修改人：张三 | 协作Agent：库存优化器  │
  │ Agent建议：该SKU建议补货200件        │
  │ [人工确认] [拒绝] [修改参数]         │
  └─────────────────────────────────────┘
```

Agent 可独立触发工作流，关键节点通过 Human-in-the-Loop 确认。

### Pillar 3: Model Agnostic

```go
// 上层业务完全感知不到底层模型切换
type Gateway interface {
    Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
}

// Provider 热插拔
var providers = map[string]Provider{
    "openai":    &OpenAIProvider{},
    "claude":    &ClaudeProvider{},
    "deepseek":  &DeepSeekProvider{},
    "ollama":    &OllamaProvider{},
}
```

模型的倾销、降价、升级对上层零感知。

### Pillar 4: Tool Native

```go
// 通过 struct tag 自动暴露为 AI Tool
type OrderService struct{}

// @i56-tool(name="create_order", desc="创建补货订单")
func (s *OrderService) CreateOrder(
    ctx context.Context,
    tenantID string,  // 自动从 Session 注入
    sku      string,  // 必填, AI 提取
    qty      int,     // 必填, AI 提取
) (orderNo string, err error) {
    // 业务代码 —— 零 AI 适配成本
}
```

框架自动将 `CreateOrder` 序列化为符合 OpenAI/Claude Tool Schema 的 JSON。

### Pillar 5: Context Aware

用户在 UI 点击某 SKU 后唤醒 AI：

```
Context Manager 自动收集：
  TenantID:     TW001
  Warehouse:    Taipei-No2
  ActiveRow:    SKU-996
  Role:         warehouse_manager
  RecentOps:    [received, weighed]

→ 组装为 System Prompt 垫底
→ AI 第一个字即针对该仓库该 SKU 的精确回答
```

### Pillar 6: Workflow Driven

AI 不再只是"对话工具"，而是系统的**路由器和规划师**：

```
用户: "把今天所有超时订单通知客户并生成报表"
  ↓
Planner: 拆解为 DAG
  ├── Step 1: FetchTimeoutOrders(today)
  ├── Step 2: ForEach → RenderTemplate → SendNotification
  └── Step 3: Aggregate → GenerateReport
  ↓
Workflow Engine: 顺序执行，Human-in-Loop 节点暂停等待
```

---

## 三、AI Layer 十二大模块

### 模块架构总图

```
┌─────────────────────────────────────────────────────┐
│                   ai.Facade                          │
│            AIService.Chat() / Execute()              │
├──────────┬──────────┬──────────┬────────────────────┤
│ Gateway  │ Router   │ Context  │ Memory             │
│(统一网关) │(智能路由) │(上下文)  │(长短期记忆)         │
├──────────┼──────────┼──────────┼────────────────────┤
│Knowledge │ Tools    │ MCP      │ Workflow           │
│(RAG引擎) │(工具调用) │(MCP协议) │(自主编排)           │
├──────────┼──────────┼──────────┼────────────────────┤
│ Security │ Prompt   │ Agent    │ RAG                │
│(安全护栏) │(编排引擎) │(领域Agent)│(检索增强)         │
└──────────┴──────────┴──────────┴────────────────────┘
```

### 1. Gateway — 统一模型网关

```go
type Gateway struct {
    providers map[string]LLMProvider
    fallback  []string  // 降级链
}

type LLMProvider interface {
    Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
    ChatStream(ctx context.Context, req *ChatRequest) (<-chan StreamChunk, error)
}
```

### 2. Router — 智能路由引擎

| 任务类型 | 模型选择 | 考量 |
|:--|:--|:--|
| 轻量 (格式化/分类) | Qwen-2.5-7B (本地) | 速度+成本 |
| 中等 (报表/分析) | DeepSeek-V3 | 性价比 |
| 重 (合同审查/合规) | Claude 3.5 Sonnet / GPT-5 | 准确率 |
| 敏感数据 | 私有化本地模型 | 100%安全 |

### 3. Context Manager

运行时自动收集：
- Session: TenantID, UserID, Role, Permissions
- UI: ActivePage, ActiveEntity, SelectedRows
- Business: WarehouseID, ActiveOrder, RecentOperations
- Environment: Timezone, Locale, FeatureFlags

### 4. Memory

```
User Space Memory:     快捷键偏好、常用表单模板、最近操作
Enterprise Memory:     组织惯例、常配物流商、历史决策
```

存储: PostgreSQL JSONB，会话前压缩为 Memory Summary。

### 5. Knowledge Base (RAG)

```
PDF/SOP/FAQ → Chunking → Embedding → pgvector Index → Semantic Search
```

零配置：导入文档后自动触发 Micro-pipeline。

### 6. Tool Calling

```go
type ToolRegistry struct {
    tools map[string]*ToolDef
}

type ToolDef struct {
    Name        string          `json:"name"`
    Description string          `json:"description"`
    Parameters  json.RawMessage `json:"parameters"`  // JSON Schema
    Handler     func(ctx context.Context, args map[string]any) (any, error)
}
```

通过 Go reflect 自动扫描 struct tag `@i56-tool`，生成 JSON Schema。

### 7. MCP Runtime

```
I56 Framework
    ↓ 作为 MCP Server
Cursor / Claude Desktop / LangChain Agent
    ↓ 通过 JSON-RPC 2.0 over stdio/SSE
直接调用 WMS/OMS/CRM 业务工具
```

```go
type MCPRuntime struct {
    Registry  *Registry     // tools + resources + prompts
    Transport Transport     // stdio | SSE | HTTP
}

type Transport interface {
    Listen(ctx context.Context, handler JSONRPCHandler) error
}
```

### 8. Workflow Planner

```
宏命令 → ReAct Loop → DAG 任务图 → 提交 Workflow Engine
         ↓ (每步)
    Thought → Action → Observation → (循环)
```

### 9. AI Security

| 方向 | 防护 |
|:--|:--|
| 入方向 (Ingress) | Prompt 注入检测、敏感关键词过滤 |
| 出方向 (Egress) | PII 动态掩码替换 |
| 执行层 | AI RBAC 权限检查、DeleteDatabase 熔断 |
| 审计 | 所有 AI 调用全量审计日志 |

### 10. Prompt Engine

```
Final Prompt = SystemBase
             + TenantOverride
             + ContextTokens
             + MemorySnippet
             + ToolSchemas
             + UserMessage
```

管理后台支持 Prompt A/B 测试与版本回滚。

### 11. Enterprise Agent

```go
type EnterpriseAgent struct {
    Name        string
    Trigger     []EventType       // 监听的事件总线事件
    State       AgentState        // idle / running / paused / error
    Interval    time.Duration     // 轮询间隔
    Handler     func(ctx context.Context, event Event) error
}
```

长驻 Goroutine，事件驱动苏醒，自主执行业务操作。

### 12. AI Workspace

前端 BDL 2.0 顶层呈现：

```
┌──────────────────────────────┬────────────────────────┐
│                              │  <i56-ai-workspace>    │
│  业务页面 (订单/仓库/报表)     │  ┌──────────────────┐ │
│  <i56-wms-order-table>       │  │ <i56-ai-bar>      │ │
│                              │  │ Ctrl+K 唤醒...    │ │
│                              │  └──────────────────┘ │
│                              │  ┌──────────────────┐ │
│                              │  │ <i56-ai-chat>    │ │
│                              │  │ SSE逐字流式       │ │
│                              │  └──────────────────┘ │
│                              │  ┌──────────────────┐ │
│                              │  │Agent运行看板      │ │
│                              │  └──────────────────┘ │
└──────────────────────────────┴────────────────────────┘
```

---

## 四、终态工程目录

```
i56-framework/
├── core/            # 传统基础设施
│   ├── auth/        # JWT 认证
│   ├── rbac/        # 权限控制
│   ├── tenant/      # 多租户
│   ├── eventbus/    # 事件总线
│   ├── workflow/    # 工作流引擎
│   └── audit/       # 审计日志
│
├── ai/              # ★ 2.0 战略核心
│   ├── ai.go        # Facade 门面
│   ├── gateway/     # 统一模型网关
│   ├── router/      # 智能路由引擎
│   ├── context/     # 上下文自动注入
│   ├── memory/      # 企业级记忆管理
│   ├── prompt/      # Prompt 编排引擎
│   ├── agent/       # 领域 Agent 容器
│   ├── workflow/    # AI 规划器 + DAG
│   ├── rag/         # RAG 流水线
│   ├── knowledge/   # 知识库管理
│   ├── tools/       # Tool 注册与调用
│   ├── mcp/         # MCP Runtime
│   └── security/    # AI 安全护栏
│
├── modules/         # 业务模块 (DDD)
│   ├── wms/
│   ├── oms/
│   ├── crm/
│   └── finance/
│
├── apps/            # 产品装配
│   ├── i56-wms/
│   └── i56-erp/
│
└── plugins/         # 第三方 + AI Agent 市场
    ├── shopify/
    ├── fedex/
    └── agents/      # AI Agent 市场
```

---

## 五、1.1 → 2.0 迁移路径

| 阶段 | 内容 | 兼容性 |
|:--|:--|:--|
| Phase 1 | `ai/` 包独立编译, 零侵入 core/ | 100% 向后兼容 |
| Phase 2 | `core/` 增加事件钩子供 AI 监听 | 99% 兼容 |
| Phase 3 | `modules/` 增加 `@i56-tool` struct tag | 向后兼容 |
| Phase 4 | BDL 2.0 AI 组件可选引入 | 渐进增强 |
| Phase 5 | 默认启用 AI Runtime | 可配置关闭 |

### 核心原则: 绝不断崖式升级

```
I56 Framework 1.1 代码
  ↓
import "github.com/i56/i56-framework/ai"  // 仅增加 import
  ↓
编译通过, 1.1 功能全部正常运行, AI 能力可选激活
```

---

## 六、商业价值

1. **模型变迁风险解耦** — 更换模型如同拔插 U 盘
2. **MCP 生态无限扩展** — Cursor/Claude Desktop 原生挂载操控 WMS/ERP
3. **Agent 插件市场** — 按需下载"AI 虚拟员工"(智能催收/库存优化/自动对账)
4. **Human-in-the-Loop 治理** — Agent 自治同时保留人类最终决策权

---

*I56 Framework 2.0 — 智能化工业软件的基础设施平台*
