// Package router implements intelligent model routing based on task
// complexity, cost constraints, and latency requirements.
package router

import (
	"context"

	"github.com/i56/framework/ai/gateway"
)

// TaskTier classifies the task complexity for routing decisions.
type TaskTier string

const (
	TierLight     TaskTier = "light"     // Simple Q&A, classification, summarization
	TierHeavy     TaskTier = "heavy"     // Complex reasoning, long generation, multi-step
	TierSensitive TaskTier = "sensitive" // PII, financial, legal — requires on-prem/private models
)

// RoutePolicy defines how tasks are distributed across providers.
type RoutePolicy string

const (
	PolicyCostFirst    RoutePolicy = "cost_first"
	PolicyLatencyFirst RoutePolicy = "latency_first"
	PolicyQualityFirst RoutePolicy = "quality_first"
	PolicyFallback     RoutePolicy = "fallback"
)

// Route contains the resolved target for a model request.
type Route struct {
	Provider  string          `json:"provider"`
	Model     string          `json:"model"`
	Gateway   gateway.Gateway `json:"-"`
	Reasoning string          `json:"reasoning"`
}

// TierConfig maps a task tier to a preferred provider chain.
type TierConfig struct {
	Tier     TaskTier
	Primary  string   // Primary provider (e.g., "openai")
	Fallback []string // Fallback providers in priority order
	Models   []string // Preferred models for this tier
}

// Router intelligently routes AI requests to the optimal provider and model.
type Router struct {
	gateways map[string]gateway.Gateway // provider name → gateway
	tiers    map[TaskTier]*TierConfig   // tier → routing config
	policy   RoutePolicy               // default routing strategy
}

// New creates a Router with the given policy.
func New(policy RoutePolicy) *Router {
	return &Router{
		gateways: make(map[string]gateway.Gateway),
		tiers:    make(map[TaskTier]*TierConfig),
		policy:   policy,
	}
}

// Register adds a gateway under a provider name.
func (r *Router) Register(name string, gw gateway.Gateway) {
	r.gateways[name] = gw
}

// SetTier configures routing for a task tier.
func (r *Router) SetTier(cfg *TierConfig) {
	r.tiers[cfg.Tier] = cfg
}

// Route resolves the best gateway for a given request and tier.
func (r *Router) Route(ctx context.Context, tier TaskTier, req *gateway.ChatRequest) (*Route, error) {
	tierCfg, ok := r.tiers[tier]
	if !ok {
		// Fall back to any registered gateway
		for name, gw := range r.gateways {
			return &Route{
				Provider:  name,
				Model:     req.Model,
				Gateway:   gw,
				Reasoning: "default fallback (no tier config)",
			}, nil
		}
		return nil, nil // no gateways available; caller should handle
	}

	// Try primary
	if gw, ok := r.gateways[tierCfg.Primary]; ok {
		model := selectModel(tierCfg.Models, req.Model)
		return &Route{
			Provider:  tierCfg.Primary,
			Model:     model,
			Gateway:   gw,
			Reasoning: "primary tier match",
		}, nil
	}

	// Try fallbacks
	for _, fb := range tierCfg.Fallback {
		if gw, ok := r.gateways[fb]; ok {
			model := selectModel(tierCfg.Models, req.Model)
			return &Route{
				Provider:  fb,
				Model:     model,
				Gateway:   gw,
				Reasoning: "fallback tier match",
			}, nil
		}
	}

	return nil, nil
}

// Chat routes and executes a chat request through the appropriate gateway.
func (r *Router) Chat(ctx context.Context, tier TaskTier, req *gateway.ChatRequest) (*gateway.ChatResponse, error) {
	route, err := r.Route(ctx, tier, req)
	if err != nil {
		return nil, err
	}
	if route == nil {
		return nil, nil
	}
	return route.Gateway.Chat(ctx, req)
}

// CatStream routes and executes a streaming chat request.
func (r *Router) ChatStream(ctx context.Context, tier TaskTier, req *gateway.ChatRequest) (<-chan gateway.StreamEvent, error) {
	route, err := r.Route(ctx, tier, req)
	if err != nil {
		return nil, err
	}
	if route == nil {
		return nil, nil
	}
	return route.Gateway.ChatStream(ctx, req)
}

// Count returns the number of registered gateways.
func (r *Router) Count() int {
	return len(r.gateways)
}

// selectModel picks the best model: request override > tier preferred > empty (let gateway use default).
func selectModel(tierModels []string, requestedModel string) string {
	if requestedModel != "" {
		return requestedModel
	}
	if len(tierModels) > 0 {
		return tierModels[0]
	}
	return ""
}
