# I56 AI Framework — Capability Layer Design Document

> **版本**: 1.0 Draft  
> **日期**: 2026-07-12  
> **定位**: I56 Framework 1.0 LTS 的 AI 能力层 — 七模块设计、架构、API 规范与实施路线图  
> **依赖**: `framework/ai/` (Gateway, Router, Security, Tools, Context, Prompt, Memory, MCP, RAG, Knowledge, Agent, Workflow)  
> **状态**: Planning → Construction

---

## 目录

1. [设计原则](#1-设计原则)
2. [架构总览](#2-架构总览)
3. [模块一：AI Chat Assistant (chatbot/)](#3-模块一ai-chat-assistant-chatbot)
4. [模块二：Smart Classification (classifier/)](#4-模块二smart-classification-classifier)
5. [模块三：Anomaly Detection (anomaly/)](#5-模块三anomaly-detection-anomaly)
6. [模块四：Route Optimization (optimizer/)](#6-模块四route-optimization-optimizer)
7. [模块五：Document Intelligence (ocr/)](#7-模块五document-intelligence-ocr)
8. [模块六：Predictive Analytics (forecast/)](#8-模块六predictive-analytics-forecast)
9. [模块七：Translation Engine (translate/)](#9-模块七translation-engine-translate)
10. [AI Core (模型路由 & 成本追踪)](#10-ai-core-模型路由--成本追踪)
11. [CLI 命令层](#11-cli-命令层)
12. [配置层](#12-配置层)
13. [API 规范](#13-api-规范)
14. [数据模型](#14-数据模型)
15. [实施路线图](#15-实施路线图)
16. [技术决策记录](#16-技术决策记录)

---

## 1. 设计原则

| 原则 | 说明 | 实践 |
|------|------|------|
| **Modular** | 每个 AI 能力是一个独立 Go 包，通过清晰接口暴露功能 | `internal/ai/{module}/` 独立目录 + `Service` 接口 |
| **Optional** | AI 模块注入失败不影响核心系统启动 | Feature flag + nil-check guard 模式 |
| **Multi-Model** | 统一 Gateway 抽象，一行配置切换 OpenAI/Claude/DeepSeek/Ollama | 复用 `framework/ai/gateway.Gateway` 接口 |
| **Observable** | 每次 AI 调用：日志、trace、token 用量、成本、延迟 | 复用 `framework/core/audit` + `CostTracker` |
| **Domain-Aware** | AI 理解 WMS 领域上下文（仓库、客户、包裹、航线） | 复用 `framework/ai/context.Manager` (AmbientContext) |
| **Fail-Safe** | AI 不可用时降级到确定性规则引擎 | `classifier` / `anomaly` 各有 rule-based fallback |
| **Tenant-Isolated** | 每个租户独立 AI 数据、记忆、配置 | `TenantID → ContextManager` + 租户级 prompt override |

---

## 2. 架构总览

### 2.1 分层架构

```
┌─────────────────────────────────────────────────────────────────────┐
│                        I56 WMS Application                          │
│  apps/wms/cmd/server/main.go — AI capability wiring                │
├─────────────────────────────────────────────────────────────────────┤
│  apps/wms/internal/ai/             ← AI Capability Layer (NEW)      │
│  ┌──────────┬───────────┬──────────┬──────────┬───────┬──────────┐ │
│  │ chatbot  │ classifier│ anomaly  │ optimizer│  ocr  │ forecast │ │
│  │          │           │          │          │       │ translate│ │
│  └──────────┴───────────┴──────────┴──────────┴───────┴──────────┘ │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │ core/              ← AI Core (Cost Tracker, Model Router,     │   │
│  │                       Prompt Templates, Feature Flags)        │   │
│  └──────────────────────────────────────────────────────────────┘   │
├─────────────────────────────────────────────────────────────────────┤
│  framework/ai/                     ← AI Runtime (ALREADY BUILT)     │
│  gateway/ router/ security/ tools/ context/ prompt/                 │
│  memory/ mcp/ agent/ workflow/ rag/ knowledge/                      │
├─────────────────────────────────────────────────────────────────────┤
│  framework/core/                   ← Core Framework                 │
│  router/ middleware/ audit/ scheduler/ db/ cache/ events/           │
├─────────────────────────────────────────────────────────────────────┤
│  modules/                          ← Domain Modules                 │
│  parcel/ order/ client/ warehouse/ carrier/ route/ rbac/           │
└─────────────────────────────────────────────────────────────────────┘
```

### 2.2 模块职责边界

```
                    ┌────────────────────────────────┐
                    │         AI Core (core/)         │
                    │  • ModelRouter (Tier routing)   │
                    │  • CostTracker (token/cost/op)  │
                    │  • PromptCatalog (templates)     │
                    │  • FeatureFlags (on/off/fallback)│
                    │  • AIAuditWriter (trace log)    │
                    └──────────┬─────────────────────┘
                               │ depends on
         ┌─────────────────────┼─────────────────────┐
         │                     │                     │
  ┌──────▼──────┐    ┌─────────▼──────┐    ┌───────▼──────┐
  │  chatbot/   │    │  classifier/   │    │  forecast/   │
  │ • Chat API  │    │ • Cargo type   │    │ • Volume     │
  │ • Co-pilot  │    │ • HS code      │    │ • Delay      │
  │ • Auto-reply│    │ • Risk score   │    │ • Seasonal   │
  └─────────────┘    └────────────────┘    └──────────────┘
         │                     │                     │
         └─────────────────────┼─────────────────────┘
                               │
                    ┌──────────▼──────────────┐
                    │  framework/ai/          │
                    │  • Gateway (4 providers)│
                    │  • Security (PII/inject)│
                    │  • Context (Ambient)    │
                    │  • Prompt Engine        │
                    └─────────────────────────┘
```

### 2.3 依赖方向 (Strict)

```
apps/wms/internal/ai/* → framework/ai/* → framework/core/* → modules/*
                         (One-way dependency, never reverse)
```

**规则**：
- `internal/ai/` 模块引用 `framework/ai/` (Gateway, Router, Security, Prompt)
- `internal/ai/` 模块引用 `modules/` (从 repos 读取业务数据作为 AI 上下文)
- `modules/` 不引用 `internal/ai/` (module 层保持 AI-free)
- `framework/ai/` 不引用 `internal/ai/` (framework 不依赖 application 层)

### 2.4 AI 请求生命周期

```
1. HTTP Request → Module Handler
2. Handler → AI Core.ModelRouter.Select(module, task)
3. Core → SecurityGuardrail.PreFlight(userInput, role, resource)
4. Core → ContextManager.AmbientInject(sessionID, tenantID, warehouseID)
5. Core → PromptCatalog.Load(module, templateName)
6. Core → Gateway.Chat(ctx, composedRequest)
7. Core → CostTracker.Record(callID, provider, model, tokenUsage, cost)
8. Core → AIAuditWriter.Log(action, module, input_hash, output_hash, latency)
9. Response → Handler → SSE stream (if chat) / JSON (if structured)
```

---

## 3. 模块一：AI Chat Assistant (chatbot/)

```
目录: apps/wms/internal/ai/chatbot/
接口: ChatService
触发词: "本月订单量多少？", "帮我查一下 VN230712-001 的状态"
```

### 3.1 能力概述

| 场景 | 入口 | 用户 | 说明 |
|------|------|------|------|
| **Admin Co-pilot** | Admin 页面底部 Chat Panel | 运营/仓库管理 | 自然语言查询 + 操作建议 |
| **Client Chatbot** | Client Portal 右下角 Widget | 外部客户 | 包裹查询、价格咨询、问题 FAQ |
| **PDA Assistant** | PDA 页面语音/文本输入 | 仓库操作员 | "扫描了 VN123, 下一步做什么?" |

### 3.2 子模块

#### 3.2.1 Admin Co-pilot (`admin_copilot.go`)

```go
// AdminCopilot 提供管理后台的 AI 辅助
type AdminCopilot struct {
    aiSvc      *ai.AIService
    core       *core.AICore
    toolRepo   *tools.Registry
    parcelRepo modules.ParcelRepository
    orderRepo  modules.OrderRepository
    // ...
}

// HandleChat 接收自然语言查询，返回 SSE 流式响应
func (c *AdminCopilot) HandleChat(ctx context.Context, req *ChatRequest) (<-chan StreamEvent, error)

// SuggestActions 分析当前页面上下文，推荐操作
func (c *AdminCopilot) SuggestActions(ctx context.Context, page string, entityID int64) ([]ActionSuggestion, error)

// ExplainData 用自然语言解释数据/图表
func (c *AdminCopilot) ExplainData(ctx context.Context, query string) (string, error)
```

**核心 Tool 注册** (函数调用工具，AI 可自动调用):

```go
// AI 可调用的工具函数 (通过 framework/ai/tools 注册)
type QueryOrdersInput struct {
    DateRange string `json:"date_range" tool:"description=日期范围 (today/week/month)"`
    Status    string `json:"status,omitempty" tool:"description=订单状态"`
    ClientID  int64  `json:"client_id,omitempty" tool:"description=客户ID"`
}
// → 返回: "今天共 156 笔订单, 已完成 89 笔, 待处理 67 笔"

type QueryParcelInput struct {
    TrackingNo string `json:"tracking_no" tool:"description=跟踪号"`
}
// → 返回: 包裹详情 + 当前状态 + 位置

type QueryInventoryInput struct {
    WarehouseID int64  `json:"warehouse_id" tool:"description=仓库ID"`
    SKU         string `json:"sku,omitempty" tool:"description=SKU"`
}
// → 返回: 库存数据
```

**SSE 流式聊天端点** (复用框架 SSE Hub):

```
GET /api/v1/ai/chat/stream?session={id}&message={query}
 → SSE event stream: text deltas + tool_call events + done
```

#### 3.2.2 Client Chatbot (`client_chatbot.go`)

```go
type ClientChatbot struct {
    aiSvc    *ai.AIService
    core     *core.AICore
    faqRepo  modules.FAQRepository // 知识库 FAQ
}

// AutoReply 自动回复客户常见问题
func (b *ClientChatbot) AutoReply(ctx context.Context, clientID int64, message string) (string, error)

// SearchParcels 自然语言包裹查询
func (b *ClientChatbot) SearchParcels(ctx context.Context, clientID int64, query string) ([]ParcelSummary, error)
```

**FAQ 匹配策略**: 先用本地向量匹配 (语义相似度 > 0.85 → 直接返回), 匹配失败再调用 LLM。

#### 3.2.3 PDA Assistant (`pda_assistant.go`)

```go
type PDAAssistant struct {
    aiSvc     *ai.AIService
    core      *core.AICore
    pdaOps    *pdaSvc.PDAOperations
}

// GuideNextStep 根据当前操作状态引导下一步
// "你刚扫描了 VN123, 包裹重量 2.3kg, 下一步请扫描库位条码进行上架"
func (a *PDAAssistant) GuideNextStep(ctx context.Context, opID int64, currentStep string) (string, error)
```

### 3.3 System Prompt 模板

```yaml
# config/ai/prompts/admin_copilot.yaml
system: |
  你是 I56 仓储管理平台的 AI 助手。
  你可以查询订单、包裹、库存、客户数据，并给出操作建议。
  当前上下文: 租户=${tenant_name}, 仓库=${warehouse_name}, 用户=${user_name}(角色=${role})
  
  规则:
  1. 回答前先调用查询工具获取最新数据
  2. 数据用中文呈现，数字加千位分隔符
  3. 涉及金额时标注币种 (USD/RMB)
  4. 不确定时明确说"需要进一步确认"
  5. 操作建议需标注风险等级 (🟢低/🟡中/🔴高)
```

### 3.4 API 端点

| Method | Path | Auth | 说明 |
|--------|------|------|------|
| `GET` | `/api/v1/ai/chat/stream` | JWT + Tenant | SSE 流式聊天 |
| `POST` | `/api/v1/ai/chat/sync` | JWT + Tenant | 同步聊天 (工具调用) |
| `GET` | `/api/v1/ai/chat/suggestions` | JWT + Tenant | 页面操作建议 |
| `GET` | `/api/v1/ai/chat/history` | JWT + Tenant | 对话历史 |
| `DELETE` | `/api/v1/ai/chat/history/{sessionID}` | JWT | 清除会话 |
| `GET` | `/api/v1/client/ai/chat` | Client Token | Client chatbot 聊天 |
| `GET` | `/api/v1/pda/ai/guide` | PDA Token | PDA 操作引导 |

---

## 4. 模块二：Smart Classification (classifier/)

```
目录: apps/wms/internal/ai/classifier/
接口: ClassifierService
触发词: 货物入库时自动判断类型、HS Code 推荐、风险分级
```

### 4.1 能力概述

| 功能 | 输入 | 输出 | 延迟要求 |
|------|------|------|----------|
| **货物类型检测** | 产品名称 + 描述 | `cargo_type` (general/electronics/food/liquid/fragile/dangerous) | < 500ms |
| **HS Code 推荐** | 货物信息 + 目的国 | HS Code Top-3 候选 + 置信度 | < 1s |
| **风险分级** | 包裹来源/内容/目的地 | `risk_level` (low/medium/high) + 原因 | < 300ms |

### 4.2 核心接口

```go
type ClassifierService struct {
    aiSvc       *ai.AIService
    core        *core.AICore
    ruleEngine  *RuleEngine     // 规则引擎 fallback (AI 不可用时)
    cache       core.CacheStore // 分类结果缓存 (命中率 > 80%)
}

// ClassifyCargo 根据产品描述分类货物类型
func (s *ClassifierService) ClassifyCargo(ctx context.Context, input *CargoClassifyInput) (*CargoClassifyOutput, error)

// SuggestHSCode 推荐海关 HS 编码
func (s *ClassifierService) SuggestHSCode(ctx context.Context, input *HSCodeInput) (*HSCodeOutput, error)

// AssessRisk 评估包裹风险等级
func (s *ClassifierService) AssessRisk(ctx context.Context, input *RiskInput) (*RiskOutput, error)

// BatchClassify 批量分类 (异步，用于大批量入库)
func (s *ClassifierService) BatchClassify(ctx context.Context, items []CargoClassifyInput) ([]CargoClassifyOutput, error)
```

### 4.3 Rule Engine Fallback (AI 不可用时)

```go
type RuleEngine struct {
    keywordRules map[string]cargoType    // "电池" → dangerous
    forbiddenWords []string              // 违禁品词库
    hsCountryMap map[string][]HSCode     // 常用国家 → HS 码映射
}

// 规则优先级: keyword match → regex pattern → AI classification → 返回 "general"
func (r *RuleEngine) Classify(productName, description string) cargoType
```

### 4.4 分类 Prompt 模板

```
你是 I56 跨境物流货物分类专家。
根据商品名称和描述，判断货物类型。
分类: general(普货) | electronics(电子) | food(食品) | liquid(液体)
       | fragile(易碎) | dangerous(危险品) | textile(纺织品) | machinery(机械)

商品名: "iPhone 15 Pro Max 256GB 黑色"
描述: "全新未拆封智能手机"
出库国: CN
目的国: US

请返回 JSON:
{"cargo_type": "...", "confidence": 0.95, "reasoning": "...", "hs_code_suggestions": ["8517.12.00", "..."]}

注意:
- 含电池的电子设备必须标注 battery_required=true
- 食品类需要标注是否需要 FDA 认证
- 危险品必须立即标记 risk_level=high
```

### 4.5 API 端点

| Method | Path | Auth | 说明 |
|--------|------|------|------|
| `POST` | `/api/v1/ai/classify/cargo` | JWT | 单件货物分类 |
| `POST` | `/api/v1/ai/classify/cargo/batch` | JWT | 批量分类 |
| `POST` | `/api/v1/ai/classify/hs-code` | JWT | HS Code 推荐 |
| `POST` | `/api/v1/ai/classify/risk` | JWT | 风险评级 |

### 4.6 缓存策略

```go
// 货物分类结果按 product_name hash 缓存 24h
cacheKey := fmt.Sprintf("classify:cargo:%s", md5Hash(input.ProductName))
if cached, ok := s.cache.Get(ctx, cacheKey); ok {
    return cached, nil
}
result, _ := s.callAI(ctx, input)
s.cache.Set(ctx, cacheKey, result, 24*time.Hour)
```

---

## 5. 模块三：Anomaly Detection (anomaly/)

```
目录: apps/wms/internal/ai/anomaly/
接口: AnomalyService
触发词: 重量差异 > 阈值、异常发货模式、可疑地址
```

### 5.1 检测维度

| 异常类型 | 检测方法 | 严重度 | 触发条件 |
|----------|----------|--------|----------|
| **重量差异** | 预报重量 vs 实际称重差值 | 🟡 中 | 差值 > 15% |
| **体积差异** | 预报体积 vs 实际测量体积 | 🟡 中 | 差值 > 20% |
| **高频发货** | 同一地址 24h 内 > N 包裹 | 🟠 高 | > 10 包裹/地址/天 |
| **异常时段** | 深夜创建大量预报 | 🟡 中 | 凌晨 2-5 点 > 5 预报 |
| **可疑地址** | 地址模式匹配欺诈库 | 🔴 严重 | 已知欺诈地址 |
| **品类突变** | 客户突然发全新品类 | 🟡 中 | 新品类占比 > 50% |
| **价格异常** | 申报价值显著偏离均值 | 🟡 中 | 偏离均值 > 3σ |

### 5.2 核心接口

```go
type AnomalyService struct {
    aiSvc         *ai.AIService
    core          *core.AICore
    ruleEngine    *AnomalyRuleEngine      // 统计规则 fallback
    parcelRepo    modules.ParcelRepository
    orderRepo     modules.OrderRepository
    eventBus      events.Bus              // 发布 AnomalyDetected 事件
}

// DetectWeightAnomaly 检测重量差异 (称重时触发)
func (s *AnomalyService) DetectWeightAnomaly(ctx context.Context, parcel *Parcel, actualWeight float64) (*AnomalyResult, error)

// DetectPatternAnomaly 检测发货模式异常 (定时任务触发)
func (s *AnomalyService) DetectPatternAnomaly(ctx context.Context, tenantID int64, window time.Duration) ([]AnomalyResult, error)

// DetectAddressAnomaly 检测可疑地址
func (s *AnomalyService) DetectAddressAnomaly(ctx context.Context, address *Address) (*AnomalyResult, error)

// GetAnomalyHistory 查询异常历史
func (s *AnomalyService) GetAnomalyHistory(ctx context.Context, filter AnomalyFilter) ([]AnomalyResult, int64, error)
```

### 5.3 事件驱动流

```
ParcelWeighed Event
    ↓
AnomalyService.DetectWeightAnomaly()
    ↓
AnomalyDetected Event → [NotificationCenter] → SMS/Line 通知运营
                      → [AuditLog] → 记录异常日志
                      → [ParcelRepo] → 更新 parcel.status = "flagged"
```

### 5.4 混合检测策略

```go
// 双层检测: 快速规则过滤 → AI 深度分析
func (s *AnomalyService) DetectWeightAnomaly(ctx context.Context, p *Parcel, actual float64) (*AnomalyResult, error) {
    diff := math.Abs(actual - p.DeclaredWeight) / p.DeclaredWeight

    // Layer 1: 规则引擎快速判断 (零延迟)
    if diff < 0.10 {
        return &AnomalyResult{Anomaly: false}, nil
    }
    if diff > 0.50 {
        return &AnomalyResult{Anomaly: true, Severity: "critical", Method: "rule"}, nil
    }

    // Layer 2: AI 深度分析 (15%-50% 灰色区域)
    return s.aiDeepCheck(ctx, p, actual)
}

func (s *AnomalyService) aiDeepCheck(ctx context.Context, p *Parcel, actual float64) (*AnomalyResult, error) {
    // 结合历史数据: 该客户历史重量偏差、该品类典型重量、最近是否有批量差异
    histData := s.gatherHistoricalContext(ctx, p.TenantID, p.ClientID, p.CargoType)
    prompt := s.core.Prompts.Load("anomaly/weight_check",
        "parcel", p, "history", histData, "actual_weight", actual)
    resp, _ := s.aiSvc.Chat(ctx, sessionID, router.TierLight, prompt)
    return parseAnomalyResult(resp.Content), nil
}
```

### 5.5 API 端点

| Method | Path | Auth | 说明 |
|--------|------|------|------|
| `POST` | `/api/v1/ai/anomaly/check` | JWT | 单件异常检测 |
| `GET` | `/api/v1/ai/anomaly/list` | JWT | 异常列表 (分页) |
| `GET` | `/api/v1/ai/anomaly/stats` | JWT | 异常统计数据 |
| `POST` | `/api/v1/ai/anomaly/resolve` | JWT + admin | 标记异常已处理 |

---

## 6. 模块四：Route Optimization (optimizer/)

```
目录: apps/wms/internal/ai/optimizer/
接口: OptimizerService
触发词: "从深圳到纽约, 选哪条线路最省钱?", "这批货和早上那批可以拼柜"
```

### 6.1 能力概述

| 功能 | 说明 | 输出 |
|------|------|------|
| **最优路线推荐** | 基于 cost/speed/reliability 三维评分 | 路线 Top-3 + 评分 |
| **拼货建议** | 相同目的地、相近时间窗的包裹合并 | 拼货组 + 预计节省 |
| **配送时间预测** | 结合历史数据和实时状态 | ETA 区间 (best/expected/worst) |
| **承运商评分** | 基于历史表现给承运商打分 | 时效准确率、破损率、成本性价比 |

### 6.2 核心接口

```go
type OptimizerService struct {
    aiSvc       *ai.AIService
    core        *core.AICore
    routeRepo   modules.RouteRepository
    carrierRepo modules.CarrierRepository
    parcelRepo  modules.ParcelRepository
    orderRepo   modules.OrderRepository
}

// SuggestRoute 推荐最优运输路线
func (s *OptimizerService) SuggestRoute(ctx context.Context, input *RouteInput) ([]RouteSuggestion, error)

// SuggestConsolidation 建议拼货组合
func (s *OptimizerService) SuggestConsolidation(ctx context.Context, tenantID int64) ([]ConsolidationGroup, error)

// PredictDeliveryTime 预测配送到达时间
func (s *OptimizerService) PredictDeliveryTime(ctx context.Context, parcelID int64) (*DeliveryPrediction, error)

// ScoreCarrier 承运商综合评分
func (s *OptimizerService) ScoreCarrier(ctx context.Context, carrierID int64) (*CarrierScore, error)
```

### 6.3 路线评分算法

```
RouteScore = w1 * CostScore + w2 * SpeedScore + w3 * ReliabilityScore

CostScore:   Normalize(estimated_cost, min_cost, max_cost)  [0,1]
SpeedScore:  Normalize(estimated_days, min_days, max_days)  [0,1] (反向)
ReliabilityScore: on_time_rate * 0.6 + damage_rate * 0.4     [0,1]

权重可配置 (默认: w1=0.4, w2=0.3, w3=0.3)
```

### 6.4 拼货逻辑

```go
// ConsolidationGroup 拼货建议
type ConsolidationGroup struct {
    Parcels      []Parcel   // 可合并的包裹
    Destination  string     // 共同目的地
    Savings      float64    // 预估节省金额
    SavingsPct   float64    // 节省百分比
    Reason       string     // AI 生成的理由
}

// 匹配条件:
// 1. 同一目的国/城市
// 2. 发货时间窗口重合 (24h 内)
// 3. 同一客户或客户允许拼货
// 4. 货物类型兼容 (危险品不能和普货拼)
```

### 6.5 API 端点

| Method | Path | Auth | 说明 |
|--------|------|------|------|
| `POST` | `/api/v1/ai/optimize/route` | JWT | 路线推荐 |
| `GET` | `/api/v1/ai/optimize/consolidate` | JWT | 拼货建议 |
| `GET` | `/api/v1/ai/optimize/delivery/{parcelID}` | JWT | 配送时间预测 |
| `GET` | `/api/v1/ai/optimize/carrier/{carrierID}/score` | JWT | 承运商评分 |

---

## 7. 模块五：Document Intelligence (ocr/)

```
目录: apps/wms/internal/ai/ocr/
接口: OCRService
触发词: 拍照上传运单 → 自动填充表单; 上传发票 → 解析金额/抬头
```

### 7.1 能力概述

| 功能 | 输入 | 输出 | 模型 |
|------|------|------|------|
| **运单识别** | 运单照片 (jpg/png) | 跟踪号、发件人、收件人、重量 | Vision LLM (GPT-4o / Claude) |
| **发票解析** | 发票 PDF/图片 | 金额、日期、抬头、税号、明细 | Vision LLM + 结构化输出 |
| **表单自动填充** | 预填模板 + 上传图片 | 填充后的表单字段 | OCR → 字段映射 |
| **条码识别** | 条码照片 | 条码值 + 类型 (Code128/EAN13) | zbar / local decoder |

### 7.2 核心接口

```go
type OCRService struct {
    aiSvc        *ai.AIService
    core         *core.AICore
    storageSvc   storage.StorageService  // 文件存储 (Minio/Local)
    barcodeDecoder *BarcodeDecoder       // 本地条码解码器
}

// ExtractShippingLabel 从运单图片提取结构化数据
func (s *OCRService) ExtractShippingLabel(ctx context.Context, imageURL string) (*ShippingLabelData, error)

// ParseInvoice 解析发票/账单
func (s *OCRService) ParseInvoice(ctx context.Context, fileURL string) (*InvoiceData, error)

// AutoFillForm 自动填充表单 (预填已知字段 + AI 补全)
func (s *OCRService) AutoFillForm(ctx context.Context, templateID string, imageURLs []string) (map[string]string, error)

// DecodeBarcode 解码条码 (本地 zbar, 不调用 AI)
func (s *OCRService) DecodeBarcode(ctx context.Context, imageData []byte) ([]BarcodeResult, error)
```

### 7.3 运单提取 Schema

```go
type ShippingLabelData struct {
    TrackingNumber  string  `json:"tracking_number"`
    SenderName      string  `json:"sender_name"`
    SenderAddress   string  `json:"sender_address"`
    SenderPhone     string  `json:"sender_phone"`
    ReceiverName    string  `json:"receiver_name"`
    ReceiverAddress string  `json:"receiver_address"`
    ReceiverPhone   string  `json:"receiver_phone"`
    Weight          float64 `json:"weight"`
    Pieces          int     `json:"pieces"`
    DeclaredValue   float64 `json:"declared_value"`
    CarrierName     string  `json:"carrier_name"`
    Confidence      float64 `json:"confidence"` // 整体置信度
    FieldConfidence map[string]float64 `json:"field_confidence"` // 逐字段置信度
}
```

### 7.4 处理流水线

```
Image Upload → Preprocessing (resize/normalize)
    ├──→ [BarcodeDetector] → 条码值
    └──→ [VisionLLM] → 结构化 JSON
            ├── Confidence < 0.7 → 人工审核标记
            └── Confidence ≥ 0.7 → 自动填充表单
```

### 7.5 API 端点

| Method | Path | Auth | 说明 |
|--------|------|------|------|
| `POST` | `/api/v1/ai/ocr/label` | JWT | 运单提取 |
| `POST` | `/api/v1/ai/ocr/invoice` | JWT | 发票解析 |
| `POST` | `/api/v1/ai/ocr/form/auto-fill` | JWT | 表单自动填充 |
| `POST` | `/api/v1/ai/ocr/barcode` | JWT | 条码解码 |

---

## 8. 模块六：Predictive Analytics (forecast/)

```
目录: apps/wms/internal/ai/forecast/
接口: ForecastService
触发词: 下月订单量预测、圣诞节高峰期资源规划
```

### 8.1 预测维度

| 指标 | 时间粒度 | 用途 | 数据源 |
|------|----------|------|--------|
| **订单量预测** | 日/周/月 | 人员排班、资源规划 | 历史订单 + 季节性因子 |
| **配送延迟预测** | 实时/每日 | 提前预警客户 | 当前在途 + 历史时效 + 天气 |
| **季节趋势分析** | 月/季度 | 促销计划、仓储资源 | 3 年历史数据 + 节假日历 |
| **客户流失风险** | 周/月 | 客户留存 | 发货频率变化 + 投诉记录 |
| **仓库容量预测** | 周/月 | 仓储扩容决策 | 入库量 + 在库时长 |

### 8.2 核心接口

```go
type ForecastService struct {
    aiSvc      *ai.AIService
    core       *core.AICore
    orderRepo  modules.OrderRepository
    parcelRepo modules.ParcelRepository
    scheduler  *scheduler.Scheduler   // 定时重新训练
}

// ForecastOrderVolume 预测订单量
func (s *ForecastService) ForecastOrderVolume(ctx context.Context, input *ForecastInput) (*VolumeForecast, error)

// PredictDeliveryDelay 预测配送延迟
func (s *ForecastService) PredictDeliveryDelay(ctx context.Context, parcelIDs []int64) ([]DelayPrediction, error)

// AnalyzeSeasonalTrends 分析季节趋势
func (s *ForecastService) AnalyzeSeasonalTrends(ctx context.Context, tenantID int64, months int) (*SeasonalReport, error)

// PredictChurnRisk 预测客户流失风险
func (s *ForecastService) PredictChurnRisk(ctx context.Context, clientID int64) (*ChurnPrediction, error)
```

### 8.3 预测方法

```go
// 混合预测: 统计模型 + AI 推理
type VolumeForecast struct {
    Period        string          // "2026-08"
    EstimatedVolume int64
    LowerBound    int64           // 80% 置信区间下界
    UpperBound    int64           // 80% 置信区间上界
    TrendFactor   float64         // 趋势因子 (YoY growth)
    SeasonFactor  float64         // 季节性因子
    AIReasoning   string          // AI 的预测解释
    DataPoints    int             // 用于预测的数据点数量
}

func (s *ForecastService) ForecastOrderVolume(ctx context.Context, input *ForecastInput) (*VolumeForecast, error) {
    // Step 1: 从仓库获取历史数据
    histData := s.loadHistoricalOrders(ctx, input.TenantID, input.LookbackMonths)

    // Step 2: 统计基线预测 (移动平均 + 季节性分解)
    statForecast := s.statisticalForecast(histData, input.TargetPeriod)

    // Step 3: AI 增强 (考虑促销、节假日、行业趋势等不可量化因素)
    context := map[string]any{
        "stat_baseline": statForecast,
        "upcoming_holidays": s.getUpcomingHolidays(input.TargetPeriod),
        "client_growth": s.getClientGrowthRate(input.TenantID),
        "market_news": s.getMarketNews(),        // 可选: 行业新闻
    }
    prompt := s.core.Prompts.Load("forecast/volume", context)
    resp, _ := s.aiSvc.Chat(ctx, sessionID, router.TierHeavy, prompt)

    // Step 4: 融合统计预测和 AI 预测 (加权平均)
    return s.fuse(statForecast, resp), nil
}
```

### 8.4 定时任务集成

```go
// 每日定时执行预测
scheduler.Register("ai-forecast-daily", "0 6 * * *", func(ctx context.Context) {
    for _, tenant := range activeTenants() {
        forecast, _ := forecastSvc.ForecastOrderVolume(ctx, &ForecastInput{
            TenantID: tenant.ID, LookbackMonths: 12, TargetPeriod: "next_month",
        })
        storeForecast(tenant.ID, forecast)
    }
})
```

### 8.5 API 端点

| Method | Path | Auth | 说明 |
|--------|------|------|------|
| `GET` | `/api/v1/ai/forecast/volume` | JWT | 订单量预测 |
| `GET` | `/api/v1/ai/forecast/delivery-delay` | JWT | 配送延迟预测 |
| `GET` | `/api/v1/ai/forecast/seasonal` | JWT | 季节趋势分析 |
| `GET` | `/api/v1/ai/forecast/churn` | JWT | 客户流失预测 |
| `GET` | `/api/v1/ai/forecast/capacity` | JWT | 仓库容量预测 |

---

## 9. 模块七：Translation Engine (translate/)

```
目录: apps/wms/internal/ai/translate/
接口: TranslateService
触发词: 商品名 CN→EN→TW 自动翻译, 客户通知多语言
```

### 9.1 能力概述

| 功能 | 方向 | 说明 |
|------|------|------|
| **商品名翻译** | CN ↔ EN ↔ TW | 用于海关申报、国际物流 |
| **通知翻译** | CN → EN/TW/JP/KO | 包裹状态通知多语言 |
| **地址翻译** | CN → EN | 国际地址格式转换 |
| **批量翻译** | 任意方向 | 表格批量翻译 |

### 9.2 核心接口

```go
type TranslateService struct {
    aiSvc      *ai.AIService
    core       *core.AICore
    termCache  core.CacheStore   // "手机壳" → "phone case" 永久缓存
}

// TranslateProduct 翻译商品名称 (用于海关申报)
func (s *TranslateService) TranslateProduct(ctx context.Context, input *TranslationInput) (*TranslationOutput, error)

// TranslateNotification 翻译客户通知模板
func (s *TranslateService) TranslateNotification(ctx context.Context, templateCode string, lang string, params map[string]any) (string, error)

// TranslateAddress 翻译并格式化地址
func (s *TranslateService) TranslateAddress(ctx context.Context, address string, fromLang, toLang string) (string, error)

// BatchTranslate 批量翻译
func (s *TranslateService) BatchTranslate(ctx context.Context, items []TranslationInput) ([]TranslationOutput, error)
```

### 9.3 翻译质量保证

```go
type TranslationOutput struct {
    Original     string  `json:"original"`
    Translated   string  `json:"translated"`
    FromLang     string  `json:"from_lang"`
    ToLang       string  `json:"to_lang"`
    Confidence   float64 `json:"confidence"`    // AI 自评置信度
    Alternatives []string `json:"alternatives,omitempty"` // 备选翻译
    Warning      string  `json:"warning,omitempty"`        // 低置信度警告
}

// 关键术语一致性检查
var productTerms = map[string]map[string]string{
    "包裹":     {"en": "parcel", "tw": "包裹"},
    "仓储":     {"en": "warehousing", "tw": "倉儲"},
    "报关":     {"en": "customs declaration", "tw": "報關"},
    "运费到付": {"en": "freight collect", "tw": "運費到付"},
    // ... 200+ 物流术语
}
```

### 9.4 通知模板翻译

```yaml
# 通知模板 (中文源)
order_created:
  zh-CN: "您的订单 {order_no} 已创建，共 {count} 件，预计 {eta} 发出"
  en:    "Your order {order_no} has been created with {count} items, estimated shipment {eta}"
  tw:    "您的訂單 {order_no} 已建立，共 {count} 件，預計 {eta} 發出"

parcel_delivered:
  zh-CN: "您的包裹 {tracking_no} 已签收"
  en:    "Your parcel {tracking_no} has been delivered"
  tw:    "您的包裹 {tracking_no} 已簽收"
```

### 9.5 API 端点

| Method | Path | Auth | 说明 |
|--------|------|------|------|
| `POST` | `/api/v1/ai/translate/product` | JWT | 商品名翻译 |
| `POST` | `/api/v1/ai/translate/notification` | JWT | 通知模板翻译 |
| `POST` | `/api/v1/ai/translate/address` | JWT | 地址翻译 |
| `POST` | `/api/v1/ai/translate/batch` | JWT | 批量翻译 |

---

## 10. AI Core (模型路由 & 成本追踪)

```
目录: apps/wms/internal/ai/core/
包:    github.com/i56/i56-apps/i56-wms/internal/ai/core
```

### 10.1 核心组件

```go
// AICore 是 AI 能力层的中枢，封装框架 AI Runtime 并为各模块提供统一服务
type AICore struct {
    Router      *ModelRouter      // 模型智能路由
    Tracker     *CostTracker      // 成本追踪
    Prompts     *PromptCatalog    // Prompt 模板目录
    Flags       *FeatureFlags     // 功能开关
    AuditWriter *AIAuditWriter    // AI 调用审计
}

func NewAICore(cfg AICoreConfig, aiSvc *ai.AIService) *AICore { ... }
```

### 10.2 ModelRouter — 模型智能路由

```go
type ModelRouter struct {
    inner      *router.Router             // 复用 framework/ai/router
    tierConfig map[string]TierMapping     // 模块 → Tier 映射
}

// TierMapping: 每个 AI 模块绑定的模型层级
// chatbot → TierHeavy (需要推理能力)
// classifier → TierLight (需要快速响应)
// ocr → TierSensitive (含 PII 数据)

func (r *ModelRouter) Select(ctx context.Context, module string, task string) (*Route, error) {
    tier := r.tierConfig[module].Tier
    return r.inner.Route(ctx, tier, &gateway.ChatRequest{...})
}
```

**默认 Tier 配置**:

| 模块 | Tier | 首选 Provider | Fallback | 原因 |
|------|------|-------------|----------|------|
| chatbot/admin | heavy | claude | deepseek | 复杂推理 + 工具调用 |
| chatbot/client | light | deepseek | ollama | 高并发 + 低成本 |
| classifier | light | openai | deepseek | 结构化输出 (JSON mode) |
| anomaly | light | openai | deepseek | 规则为主, AI 为辅 |
| optimizer | heavy | claude | openai | 多因素权衡推理 |
| ocr | heavy | openai | claude | Vision model 能力 |
| forecast | heavy | claude | openai | 统计融合 + 推理 |
| translate | light | deepseek | openai | 高并发翻译 |

### 10.3 CostTracker — 成本追踪

```go
type CostTracker struct {
    prices map[string]PriceTable   // provider → model → per-token cost
    repo   CostRepository          // 持久化到 PostgreSQL
    mu     sync.RWMutex
}

// PriceTable 模型定价 (USD per 1K tokens)
// openai / gpt-4o:        input=$0.0025, output=$0.01
// claude / claude-sonnet: input=$0.003,  output=$0.015
// deepseek / deepseek-chat: input=$0.00027, output=$0.0011
// ollama / llama3.1:       input=$0 (local), output=$0

type CostRecord struct {
    ID           int64     `json:"id"`
    CallID       string    `json:"call_id"`
    Module       string    `json:"module"`        // chatbot/classifier/...
    Provider     string    `json:"provider"`
    Model        string    `json:"model"`
    InputTokens  int       `json:"input_tokens"`
    OutputTokens int       `json:"output_tokens"`
    Cost         float64   `json:"cost"`          // USD
    LatencyMs    int       `json:"latency_ms"`
    TenantID     int64     `json:"tenant_id"`
    UserID       int64     `json:"user_id"`
    CreatedAt    time.Time `json:"created_at"`
}

func (t *CostTracker) Record(ctx context.Context, record *CostRecord) error

// 查询接口
func (t *CostTracker) GetDailyCost(ctx context.Context, tenantID int64, date time.Time) (float64, error)
func (t *CostTracker) GetModuleCost(ctx context.Context, tenantID int64, module string, from, to time.Time) (float64, error)
func (t *CostTracker) GetUsageStats(ctx context.Context, tenantID int64) (*UsageStats, error)
```

**成本限制**:

```go
type FeatureFlags struct {
    DailyLimit   float64  // 每日最大 AI 成本 (default: $50)
    MonthlyLimit float64  // 每月最大 AI 成本 (default: $1000)
    ModuleLimits map[string]float64 // 每个模块每月限额
    FallbackRule string   // 超限后行为: block / fallback_to_cheaper / alert
}
```

### 10.4 PromptCatalog — Prompt 模板管理

```go
type PromptCatalog struct {
    templates map[string]*PromptTemplate  // module/templateName → template
}

type PromptTemplate struct {
    Name        string
    Module      string
    System      string                 // System prompt
    Variables   []string               // Required variables
    Tools       []string               // Associated tool names
    Version     int                    // Template version
}

// Load 加载模板并用参数渲染
func (c *PromptCatalog) Load(module, templateName string, params map[string]any) (*gateway.ChatRequest, error)

// 模板示例: catalog.Load("classifier", "cargo_type", {"product_name": "...", "description": "..."})
```

**模板文件组织**:

```
apps/wms/config/ai/prompts/
├── chatbot/
│   ├── admin_copilot.yaml
│   ├── client_chatbot.yaml
│   └── pda_assistant.yaml
├── classifier/
│   ├── cargo_type.yaml
│   ├── hs_code.yaml
│   └── risk_assessment.yaml
├── anomaly/
│   └── weight_check.yaml
├── optimizer/
│   ├── route_suggestion.yaml
│   └── consolidation.yaml
├── ocr/
│   └── label_extraction.yaml
├── forecast/
│   ├── volume.yaml
│   └── churn.yaml
└── translate/
    └── product_name.yaml
```

### 10.5 FeatureFlags — 功能开关

```go
type FeatureFlags struct {
    // 全局
    AIEnabled          bool     // 全局 AI 开关
    Env                string   // dev/stg/prod

    // 模块开关
    ChatbotEnabled     bool     // 聊天助手
    ClassifierEnabled  bool     // 智能分类
    AnomalyEnabled     bool     // 异常检测
    OptimizerEnabled   bool     // 路线优化
    OCREnabled         bool     // 文档识别
    ForecastEnabled    bool     // 预测分析
    TranslateEnabled   bool     // 翻译引擎

    // Fallback 配置
    ClassifierFallback string   // "rule" | "ai" | "hybrid"
    AnomalyFallback    string   // "rule" | "ai" | "hybrid"
}

// 配置来源: 环境变量 → config.yaml → DB 动态配置
func LoadFeatureFlags(cfg *config.Config) *FeatureFlags
```

### 10.6 AIAuditWriter — AI 调用审计

```go
type AIAuditWriter struct {
    auditSvc *audit.AuditLogger  // 复用框架 audit
}

func (w *AIAuditWriter) LogCall(ctx context.Context, call *AICall) error {
    // 不记录原始 prompt/response 内容 (隐私保护)
    // 记录: 模块、操作、输入 MD5、输出 MD5、模型、延迟、成本、成功/失败
    data := map[string]any{
        "module":        call.Module,
        "action":        call.Action,
        "input_hash":    md5Hash(call.Input),
        "output_hash":   md5Hash(call.Output),
        "model":         call.Model,
        "provider":      call.Provider,
        "latency_ms":    call.LatencyMs,
        "cost":          call.Cost,
        "input_tokens":  call.InputTokens,
        "output_tokens": call.OutputTokens,
        "success":       call.Error == nil,
        "error":         truncate(call.Error, 200),
    }
    return w.auditSvc.Log(ctx, audit.Action("AI_CALL"), call.Module, call.ID, data)
}
```

---

## 11. CLI 命令层

```bash
# 模块状态查看
i56 ai status                          # 各模块健康状态 + 成本统计
i56 ai status --module=chatbot         # 单模块详情

# 模型管理
i56 ai model list                      # 列出所有注册的模型
i56 ai model test --provider=openai    # 测试模型连通性
i56 ai model switch --module=chatbot --provider=deepseek  # 切换模块模型

# 成本管理
i56 ai cost today                      # 今日成本
i56 ai cost month                      # 本月成本 (分模块)
i56 ai cost limit --daily=100          # 设置每日限额

# 功能开关
i56 ai enable chatbot                  # 启用聊天助手
i56 ai disable anomaly                 # 禁用异常检测
i56 ai feature-flags                   # 查看所有开关状态

# Prompt 管理
i56 ai prompt list                     # 列出所有模板
i56 ai prompt show --module=classifier --name=cargo_type
i56 ai prompt test --module=classifier --input='{"product_name":"iPhone"}'

# 审计日志
i56 ai audit today                     # 今日调用日志
i56 ai audit --module=chatbot --limit=50

# 批量操作
i56 ai batch classify --file=products.csv  # 批量分类
i56 ai batch translate --from=zh --to=en --file=names.csv
```

---

## 12. 配置层

### 12.1 配置结构

```yaml
# apps/wms/config/ai/config.yaml
ai:
  enabled: true
  env: production

  models:
    default: openai
    providers:
      openai:
        api_key: ${OPENAI_API_KEY}
        default_model: gpt-4o
        base_url: https://api.openai.com/v1
      claude:
        api_key: ${ANTHROPIC_API_KEY}
        default_model: claude-sonnet-4-20250514
        base_url: https://api.anthropic.com/v1
      deepseek:
        api_key: ${DEEPSEEK_API_KEY}
        default_model: deepseek-chat
        base_url: https://api.deepseek.com
      ollama:
        base_url: http://localhost:11434
        default_model: llama3.1

  router:
    policy: quality_first  # cost_first | latency_first | quality_first
    tiers:
      light:
        primary: deepseek
        fallback: [openai, ollama]
      heavy:
        primary: claude
        fallback: [openai]
      sensitive:
        primary: ollama
        fallback: []

  cost_limits:
    daily: 50.0           # USD
    monthly: 1000.0       # USD
    per_call: 5.0         # USD
    soft_limit_pct: 0.8   # 80% 时发出警告
    over_limit_action: fallback_to_cheaper

  features:
    chatbot: true
    classifier: true
    anomaly: true
    optimizer: true
    ocr: true
    forecast: true
    translate: true

  fallback:
    classifier: hybrid    # rule | ai | hybrid
    anomaly: hybrid       # rule | ai | hybrid
    translate: ai         # ai | cache_only

  cache:
    classifier_ttl: 24h
    translate_ttl: 720h   # 30 days
    anomaly_ttl: 1h

  retry:
    max_attempts: 3
    backoff_base: 1s
    backoff_max: 30s

  timeout:
    chat_stream: 300s
    classifier: 10s
    anomaly: 5s
    optimizer: 30s
    ocr: 60s
    forecast: 45s
    translate: 15s
```

### 12.2 环境变量 (12-Factor App)

```bash
# 全局
I56_AI_ENABLED=true

# 模型密钥
OPENAI_API_KEY=sk-xxx
ANTHROPIC_API_KEY=sk-ant-xxx
DEEPSEEK_API_KEY=sk-xxx

# 成本限制
I56_AI_DAILY_LIMIT=50.0
I56_AI_MONTHLY_LIMIT=1000.0

# 功能开关
I56_AI_CHATBOT_ENABLED=true
I56_AI_CLASSIFIER_ENABLED=true
I56_AI_OCR_ENABLED=false   # e.g. 测试环境禁用 OCR 以节约成本
```

---

## 13. API 规范

### 13.1 统一响应格式

```go
// 成功
{
    "status": "ok",
    "data": { ... },
    "meta": {
        "request_id": "req_abc123",
        "ai_cost": 0.0025,
        "ai_latency_ms": 342,
        "model": "gpt-4o",
        "tier": "light",
        "cached": false
    }
}

// 错误
{
    "status": "error",
    "error": {
        "code": "AI_COST_LIMIT_EXCEEDED",
        "message": "当日 AI 调用成本已达上限 ($50.00)",
        "detail": "请在明日重试或联系管理员调整限额"
    },
    "meta": {
        "request_id": "req_abc456",
        "fallback": "rule_engine"
    }
}
```

### 13.2 错误码体系

| 错误码 | HTTP | 说明 |
|--------|------|------|
| `AI_COST_LIMIT_EXCEEDED` | 429 | 成本限制已超 |
| `AI_RATE_LIMITED` | 429 | 调用频率限制 |
| `AI_MODEL_UNAVAILABLE` | 503 | 模型暂时不可用 |
| `AI_GATEWAY_TIMEOUT` | 504 | 模型响应超时 |
| `AI_SECURITY_BLOCKED` | 403 | 安全拦截 (PII/注入) |
| `AI_MODULE_DISABLED` | 503 | AI 模块已禁用 |
| `AI_CONTENT_UNSAFE` | 422 | 输入内容不安全 |
| `AI_FALLBACK_ACTIVE` | 200 | 使用规则引擎 fallback (meta 标注) |

### 13.3 调用审计响应头

每个 AI API 端点的响应包含自定义 HTTP 头:

```
X-AI-Model: gpt-4o
X-AI-Provider: openai
X-AI-Cost: 0.0025
X-AI-Latency-Ms: 342
X-AI-Tier: light
X-AI-Cached: false
X-AI-Request-Id: req_abc123
```

---

## 14. 数据模型

### 14.1 PostgreSQL 表 (新增)

```sql
-- AI 调用记录
CREATE TABLE ai_call_logs (
    id BIGSERIAL PRIMARY KEY,
    call_id VARCHAR(64) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    user_id BIGINT,
    module VARCHAR(50) NOT NULL,        -- chatbot/classifier/anomaly/...
    action VARCHAR(100) NOT NULL,        -- classify_cargo/chat/suggest_route
    provider VARCHAR(30) NOT NULL,       -- openai/claude/deepseek/ollama
    model VARCHAR(100) NOT NULL,         -- gpt-4o/claude-sonnet-4-20250514
    input_tokens INT DEFAULT 0,
    output_tokens INT DEFAULT 0,
    cost NUMERIC(12,6) DEFAULT 0,
    latency_ms INT DEFAULT 0,
    success BOOLEAN DEFAULT true,
    error_message TEXT,
    input_hash VARCHAR(64),             -- MD5 of input (for dedup detection)
    output_hash VARCHAR(64),            -- MD5 of output
    cached BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_ai_logs_tenant_date ON ai_call_logs(tenant_id, created_at);
CREATE INDEX idx_ai_logs_module ON ai_call_logs(module, created_at);

-- AI 每日成本汇总 (物化视图或定时任务维护)
CREATE TABLE ai_daily_costs (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    date DATE NOT NULL,
    module VARCHAR(50) NOT NULL,
    call_count INT DEFAULT 0,
    total_tokens INT DEFAULT 0,
    total_cost NUMERIC(12,6) DEFAULT 0,
    avg_latency_ms INT DEFAULT 0,
    UNIQUE(tenant_id, date, module)
);

-- Prompt 模板版本
CREATE TABLE ai_prompt_templates (
    id BIGSERIAL PRIMARY KEY,
    module VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    version INT DEFAULT 1,
    system_prompt TEXT NOT NULL,
    variables JSONB DEFAULT '[]',
    tool_names JSONB DEFAULT '[]',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(module, name, version)
);

-- AI 配置 (tenant-level overrides)
CREATE TABLE ai_tenant_configs (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT UNIQUE NOT NULL,
    features JSONB DEFAULT '{}',         -- { "chatbot": true, "classifier": false, ... }
    cost_limit_daily NUMERIC(12,6),
    cost_limit_monthly NUMERIC(12,6),
    preferred_models JSONB DEFAULT '{}', -- { "light": "deepseek", "heavy": "claude" }
    prompt_overrides JSONB DEFAULT '{}', -- { "classifier/cargo_type": "custom prompt text" }
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

---

## 15. 实施路线图

### Phase 0: Foundation (Week 1-2) — **P0**

```
目标: AI Core 就绪，所有模块可建立在统一基础上
```

| # | 任务 | 产出 | 依赖 |
|---|------|------|------|
| 0.1 | 创建 `internal/ai/core/` 包 | `ModelRouter`, `CostTracker`, `PromptCatalog`, `FeatureFlags`, `AIAuditWriter` | `framework/ai/*` |
| 0.2 | 配置系统集成 | `config/ai/config.yaml` + 环境变量加载 | `framework/core/config` |
| 0.3 | 数据库迁移 | `ai_call_logs`, `ai_daily_costs`, `ai_prompt_templates` 表 | PostgreSQL |
| 0.4 | `main.go` AI wiring | `aiCore := core.NewAICore(...)` 装配 + health check 端点 | 0.1, 0.2 |
| 0.5 | CLI 命令脚手架 | `i56 ai status` 基础命令 | Cobra CLI |
| 0.6 | 管理页面基架 | AI 模块状态面板 ( `admin/system/ai-status` ) | 0.3 |
| 0.7 | Admin 菜单集成 | sidebar 增加 "🤖 AI 智能" 菜单组 | `base.html` |

**验证标准**: `i56 ai status` 输出各模块状态、可用模型、今日成本

### Phase 1: Chat Assistant (Week 3-4) — **P1**

```
目标: Admin Co-pilot + Client Chatbot 可用，SSE 流式聊天
```

| # | 任务 | 产出 | 依赖 |
|---|------|------|------|
| 1.1 | `chatbot/admin_copilot.go` | Tool 函数注册 + SSE 聊天端点 | Phase 0 |
| 1.2 | `chatbot/client_chatbot.go` | Client FAQ + 包裹查询 | Phase 0 |
| 1.3 | `chatbot/pda_assistant.go` | PDA 操作引导 | Phase 0 + PDA module |
| 1.4 | Prompt 模板 | `admin_copilot.yaml`, `client_chatbot.yaml`, `pda_assistant.yaml` | Phase 0 |
| 1.5 | Admin 页面 Chat Panel | 底部聊天面板 (HTML/CSS/JS + SSE) | `base.html` |
| 1.6 | Client Portal Chat Widget | 右下角聊天组件 | `client/base.html` |

**验证标准**: Admin 页面输入 "今天有多少订单？" → AI 返回真实数据

### Phase 2: Classification + Translation (Week 5-6) — **P2**

```
目标: 货物分类自动化 + 多语言翻译可用
```

| # | 任务 | 产出 | 依赖 |
|---|------|------|------|
| 2.1 | `classifier/service.go` | 完整 CargoClassify + HSCode + Risk 实现 | Phase 0 |
| 2.2 | `classifier/rule_engine.go` | 规则引擎 fallback | 2.1 |
| 2.3 | `classifier/cache.go` | 分类结果缓存 | Phase 0 |
| 2.4 | `translate/service.go` | 翻译 + 术语表 + 缓存 | Phase 0 |
| 2.5 | `translate/templates.go` | 通知模板多语言 | 2.4 |
| 2.6 | 预报页面集成 | 创建预报时自动分类货物类型 | `parcel/create.html` |

**验证标准**: 创建预报 "iPhone 15" → 自动填充 cargo_type=electronics, hs_code=8517.12

### Phase 3: Anomaly + Optimization (Week 7-9) — **P3**

```
目标: 异常检测自动告警 + 路线优化推荐
```

| # | 任务 | 产出 | 依赖 |
|---|------|------|------|
| 3.1 | `anomaly/service.go` | 5 种异常检测 + 规则引擎 | Phase 0 |
| 3.2 | `anomaly/events.go` | AnomalyDetected 事件发布 | Phase 0 + Event Bus |
| 3.3 | `optimizer/service.go` | 路线推荐 + 拼货建议 | Phase 0 |
| 3.4 | `optimizer/scoring.go` | 承运商评分算法 | 3.3 |
| 3.5 | 称重页面集成 | 称重时自动检测重量异常 | `pda/weigh.html` |
| 3.6 | 订单页面集成 | 创建订单时推荐最优路线 | `orders/create.html` |

**验证标准**: 称重 2.3kg (预报 1.0kg) → 弹出异常警告

### Phase 4: OCR + Forecast (Week 10-12) — **P4**

```
目标: 文档智能提取 + 预测分析上线
```

| # | 任务 | 产出 | 依赖 |
|---|------|------|------|
| 4.1 | `ocr/service.go` | 运单提取 (Vision LLM) | Phase 0 |
| 4.2 | `ocr/barcode.go` | 本地条码解码器 (zbar) | 4.1 |
| 4.3 | `forecast/service.go` | 订单量预测 + 配送延迟预测 | Phase 0 |
| 4.4 | `forecast/scheduler.go` | 定时预测任务注册 | Phase 0 + Scheduler |
| 4.5 | PDA 拍照上传 | 扫描包裹 → 拍照运单 → 自动填充 | 4.1 + PDA |
| 4.6 | Dashboard 集成 | AI 预测卡片 (预计订单量/延迟风险) | `dashboard.html` |

**验证标准**: Dashboard 显示 "预计下月订单量 4,200 ± 300，高峰周 9/15-9/21"

---

## 16. 技术决策记录

### Decision 1: AI 模块放在 `apps/wms/internal/ai/` 而非 `framework/ai/`

**原因**: `framework/ai/` 是通用 AI Runtime (Gateway/Router/Security/Tools 等)，与业务无关。WMS 特定的 AI 能力 (货物分类、运单识别、路线优化) 属于应用层逻辑，放在 `apps/wms/internal/ai/`。

### Decision 2: Go 原生实现，不引入 Python 微服务

**原因**: I56 Framework 是纯 Go 生态。LLM 调用通过 HTTP API (OpenAI 兼容) 完成，不需要 Python 桥接。对于需要本地模型推理的场景，通过 Ollama (Go SDK) 或 vLLM 服务。

### Decision 3: 双层检测策略 (Rule + AI) 用于 classifier 和 anomaly

**原因**: 
- 规则引擎: 零延迟，低成本，覆盖 80% 常规场景
- AI 引擎: 处理灰色地带和复杂案例，覆盖剩余 20%
- 这种混合模式在 AI 不可用时仍能工作 (graceful degradation)

### Decision 4: SSE 流式聊天，复用框架 SSE Hub

**原因**: 
- `framework/core/sse` 已经提供 Hub/Client/Event 模型
- SSE 比 WebSocket 更简单，浏览器原生支持
- HTMX + SSE extension 无需前端框架依赖

### Decision 5: 分类结果和翻译结果缓存

**原因**: 
- 货物分类: 同一产品名返回相同结果，缓存命中率 > 80%
- 翻译: 术语库固化，不需要每次调用 AI
- 显著降低 AI 调用成本和延迟

### Decision 6: 成本限制实行软硬结合

**原因**:
- 80% 软限制: 发送告警但不阻断 (管理员知晓)
- 100% 硬限制: 阻断调用或降级到更便宜的模型
- 避免 "AI 突然停掉导致业务中断" 的尴尬

---

## 附录

### A. 文件清单 (新增)

```
apps/wms/
├── internal/ai/
│   ├── core/
│   │   ├── core.go              # AICore 总装配
│   │   ├── router.go            # ModelRouter
│   │   ├── tracker.go           # CostTracker
│   │   ├── prompts.go           # PromptCatalog
│   │   ├── flags.go             # FeatureFlags
│   │   ├── audit.go             # AIAuditWriter
│   │   └── core_test.go
│   ├── chatbot/
│   │   ├── service.go           # ChatService 接口
│   │   ├── admin_copilot.go     # Admin Co-pilot 实现
│   │   ├── client_chatbot.go    # Client Chatbot 实现
│   │   ├── pda_assistant.go     # PDA Assistant 实现
│   │   ├── tools.go             # AI 可调用工具注册
│   │   └── chatbot_test.go
│   ├── classifier/
│   │   ├── service.go           # ClassifierService
│   │   ├── rule_engine.go       # RuleEngine fallback
│   │   ├── cache.go             # 分类缓存
│   │   └── classifier_test.go
│   ├── anomaly/
│   │   ├── service.go           # AnomalyService
│   │   ├── detectors.go         # 各检测器实现
│   │   ├── rule_engine.go       # 统计规则引擎
│   │   └── anomaly_test.go
│   ├── optimizer/
│   │   ├── service.go           # OptimizerService
│   │   ├── scoring.go           # 评分算法
│   │   └── optimizer_test.go
│   ├── ocr/
│   │   ├── service.go           # OCRService
│   │   ├── barcode.go           # 条码解码器
│   │   └── ocr_test.go
│   ├── forecast/
│   │   ├── service.go           # ForecastService
│   │   ├── statistical.go       # 统计基线预测
│   │   └── forecast_test.go
│   └── translate/
│       ├── service.go           # TranslateService
│       ├── glossary.go          # 术语表管理
│       └── translate_test.go
├── config/ai/
│   ├── config.yaml              # AI 配置
│   └── prompts/                 # Prompt 模板目录
│       ├── chatbot/
│       ├── classifier/
│       ├── anomaly/
│       ├── optimizer/
│       ├── ocr/
│       ├── forecast/
│       └── translate/
├── cmd/server/
│   ├── ai_routes.go             # AI API 路由注册
│   └── ai_wiring.go             # AI 模块 DI 装配
└── migrations/
    └── 009_create_ai_tables.up.sql  # AI 数据表
```

### B. 测试策略

| 测试层 | 范围 | 工具 |
|--------|------|------|
| **Unit** | 规则引擎、评分算法、缓存逻辑 | Go testing + testify |
| **Integration** | Gateway mock → 模块接口 → 数据库 | Go testing + mock HTTP server |
| **E2E** | 端到端场景 (预报 → 分类 → 入库) | Playwright / API tests |
| **AI Quality** | 分类准确率、翻译质量、异常召回率 | 标注数据集 + metrics |

### C. 性能指标 (SLO)

| 模块 | 目标 P50 | 目标 P95 | 目标 P99 |
|------|----------|----------|----------|
| chatbot (非流式) | < 2s | < 5s | < 10s |
| chatbot (流式 TTFT) | < 500ms | < 1s | < 2s |
| classifier (cached) | < 5ms | < 20ms | < 50ms |
| classifier (AI) | < 500ms | < 1.5s | < 3s |
| anomaly (rule) | < 1ms | < 5ms | < 10ms |
| anomaly (AI) | < 300ms | < 1s | < 2s |
| optimizer | < 3s | < 8s | < 15s |
| ocr | < 5s | < 15s | < 30s |
| forecast | < 10s | < 30s | < 60s |
| translate (cached) | < 5ms | < 20ms | < 50ms |
| translate (AI) | < 500ms | < 1.5s | < 3s |

---

> **文档维护**: 本文档由 Hermes Agent 生成于 2026-07-12，随实施进展持续更新。  
> **关联文档**: [I56 Framework Architecture](framework/docs/I56-FRAMEWORK-ARCHITECTURE.md) | [AI API Surface](framework/ai/) | [v2.1 Roadmap](docs/V2.1_ROADMAP.md)
