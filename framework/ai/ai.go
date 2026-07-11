// Package ai is the main entry point for the I56 Framework AI subsystem.
// The AIService provides a unified facade over all AI capabilities:
// model routing, context management, security guardrails, tool orchestration,
// prompt composition, and RAG-based knowledge retrieval.
package ai

import (
	"context"

	aicontext "github.com/i56/framework/ai/context"
	"github.com/i56/framework/ai/gateway"
	"github.com/i56/framework/ai/memory"
	"github.com/i56/framework/ai/prompt"
	"github.com/i56/framework/ai/router"
	"github.com/i56/framework/ai/security"
	"github.com/i56/framework/ai/tools"
)

// AIService is the unified facade for all AI operations.
// It orchestrates routing, security, context injection, tool invocation,
// and prompt composition into a single, coherent API.
type AIService struct {
	Gateway  gateway.Gateway
	Router   *router.Router
	Context  *aicontext.Manager
	Security *security.Guardrail
	Tools    *tools.Registry
	Prompt   *prompt.Engine

	// UserMemory provides per-user conversation context.
	UserMemory *memory.UserMemory
}

// Config holds the configuration for initializing AIService.
type Config struct {
	DefaultTenant  string
	BasePrompt     string
	SecurityConfig security.Config
	RoutePolicy    router.RoutePolicy
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		DefaultTenant:  "default",
		BasePrompt:     "You are a helpful AI assistant for the I56 warehouse management platform.",
		SecurityConfig: security.DefaultConfig(),
		RoutePolicy:    router.PolicyQualityFirst,
	}
}

// New creates an AIService with the given configuration.
func New(cfg Config) *AIService {
	return &AIService{
		Router:   router.New(cfg.RoutePolicy),
		Context:  aicontext.NewManager(cfg.DefaultTenant),
		Security: security.New(cfg.SecurityConfig),
		Tools:    tools.NewRegistry(),
		Prompt:   prompt.NewEngine(cfg.BasePrompt),
	}
}

// RegisterGateway adds an LLM provider to the router.
func (s *AIService) RegisterGateway(name string, gw gateway.Gateway) {
	s.Router.Register(name, gw)
	if s.Gateway == nil {
		s.Gateway = gw // Set default gateway
	}
}

// RegisterTool registers a Go struct as an AI-callable tool.
func (s *AIService) RegisterTool(v any) error {
	return s.Tools.RegisterFromStruct(v)
}

// Chat sends a chat request through the full AI pipeline:
//  1. Pre-flight security checks (PII, injection, RBAC, circuit breaker)
//  2. Prompt composition (base + overrides + context + memory + tools)
//  3. Model routing (tier-based gateway selection)
//  4. Gateway execution
//  5. Post-flight processing (PII masking on response)
func (s *AIService) Chat(ctx context.Context, sessionID string, tier router.TaskTier, req *gateway.ChatRequest) (*gateway.ChatResponse, error) {
	// 1. Security pre-flight
	ac := s.Context.Get(sessionID)
	if s.Security != nil {
		role := ""
		if ac != nil {
			role = ac.Role
		}
		// Check the user message for injection
		for _, msg := range req.Messages {
			if msg.Role == gateway.RoleUser {
				if err := s.Security.PreFlight(ctx, msg.Content, role, "ai:chat"); err != nil {
					return nil, err
				}
				// Mask PII
				if masked, changed := s.Security.MaskPII(msg.Content); changed {
					msg.Content = masked
				}
			}
		}
	}

	// 2. Build tool definitions string
	toolDefs := ""
	if s.Tools != nil && s.Tools.Count() > 0 {
		toolDefs = "Available tools: " // TODO: JSON-serialize tool schemas
	}

	// 3. Compose the full message list with system prompt
	composedMessages := s.Prompt.BuildMessages(req.Messages, ac, s.UserMemory, toolDefs)
	req.Messages = composedMessages

	// 4. Route to the best gateway
	resp, err := s.Router.Chat(ctx, tier, req)
	if err != nil {
		if s.Security != nil {
			s.Security.RecordFailure()
		}
		return nil, err
	}
	if resp == nil {
		// Fallback: use default gateway
		if s.Gateway != nil {
			resp, err = s.Gateway.Chat(ctx, req)
		}
	}
	if err != nil {
		if s.Security != nil {
			s.Security.RecordFailure()
		}
		return nil, err
	}

	// 5. Post-flight: mask PII in response
	if s.Security != nil && resp != nil {
		if masked, _ := s.Security.MaskPII(resp.Content); masked != resp.Content {
			resp.Content = masked
		}
		s.Security.RecordSuccess()
	}

	return resp, nil
}

// ChatStream sends a streaming chat request through the full pipeline.
func (s *AIService) ChatStream(ctx context.Context, sessionID string, tier router.TaskTier, req *gateway.ChatRequest) (<-chan gateway.StreamEvent, error) {
	// Security pre-flight
	ac := s.Context.Get(sessionID)
	if s.Security != nil {
		role := ""
		if ac != nil {
			role = ac.Role
		}
		for _, msg := range req.Messages {
			if msg.Role == gateway.RoleUser {
				if err := s.Security.PreFlight(ctx, msg.Content, role, "ai:chat"); err != nil {
					return nil, err
				}
			}
		}
	}

	// Compose messages with system prompt
	toolDefs := ""
	if s.Tools != nil && s.Tools.Count() > 0 {
		toolDefs = "Available tools: "
	}
	req.Messages = s.Prompt.BuildMessages(req.Messages, ac, s.UserMemory, toolDefs)

	// Route and stream
	ch, err := s.Router.ChatStream(ctx, tier, req)
	if err != nil {
		if s.Security != nil {
			s.Security.RecordFailure()
		}
		return nil, err
	}
	if ch == nil && s.Gateway != nil {
		ch, err = s.Gateway.ChatStream(ctx, req)
	}
	return ch, err
}

// WithUserMemory attaches user-specific memory to the service.
func (s *AIService) WithUserMemory(um *memory.UserMemory) *AIService {
	s.UserMemory = um
	return s
}

// AddPromptOverride adds a tenant-specific prompt override layer.
func (s *AIService) AddPromptOverride(layer prompt.Layer) {
	s.Prompt.AddOverride(layer)
}
