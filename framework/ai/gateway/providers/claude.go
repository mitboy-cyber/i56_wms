package providers

import (
	"context"

	"github.com/i56/framework/ai/gateway"
)

// ClaudeGateway implements gateway.Gateway for Anthropic Claude.
type ClaudeGateway struct {
	apiKey  string
	baseURL string
	info    gateway.ProviderInfo
}

// NewClaude creates a new Claude gateway instance.
func NewClaude(apiKey string) *ClaudeGateway {
	return &ClaudeGateway{
		apiKey:  apiKey,
		baseURL: "https://api.anthropic.com/v1",
		info: gateway.ProviderInfo{
			Name:              "claude",
			Models:            []string{"claude-sonnet-4-20250514", "claude-opus-4-20250514", "claude-3-5-haiku-20241022"},
			DefaultModel:      "claude-sonnet-4-20250514",
			SupportsStreaming: true,
			SupportsTools:     true,
		},
	}
}

// Chat implements the Gateway interface for Claude.
func (g *ClaudeGateway) Chat(ctx context.Context, req *gateway.ChatRequest) (*gateway.ChatResponse, error) {
	// TODO: Implement Claude messages API call
	return &gateway.ChatResponse{
		Content: "Claude stub: implement me",
	}, nil
}

// ChatStream implements streaming for Claude.
func (g *ClaudeGateway) ChatStream(ctx context.Context, req *gateway.ChatRequest) (<-chan gateway.StreamEvent, error) {
	// TODO: Implement Claude streaming
	ch := make(chan gateway.StreamEvent)
	close(ch)
	return ch, nil
}

// Info returns provider metadata.
func (g *ClaudeGateway) Info() gateway.ProviderInfo {
	return g.info
}

// Health checks API connectivity.
func (g *ClaudeGateway) Health(ctx context.Context) error {
	return nil
}
