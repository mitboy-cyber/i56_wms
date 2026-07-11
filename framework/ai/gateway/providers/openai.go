package providers

import (
	"context"

	"github.com/i56/framework/ai/gateway"
)

// OpenAIGateway implements gateway.Gateway for OpenAI-compatible APIs.
type OpenAIGateway struct {
	apiKey  string
	baseURL string
	info    gateway.ProviderInfo
}

// NewOpenAI creates a new OpenAI gateway instance.
func NewOpenAI(apiKey string) *OpenAIGateway {
	return &OpenAIGateway{
		apiKey:  apiKey,
		baseURL: "https://api.openai.com/v1",
		info: gateway.ProviderInfo{
			Name:              "openai",
			Models:            []string{"gpt-4o", "gpt-4o-mini", "gpt-4-turbo", "gpt-3.5-turbo"},
			DefaultModel:      "gpt-4o",
			SupportsStreaming: true,
			SupportsTools:     true,
		},
	}
}

// Chat implements the Gateway interface for OpenAI.
func (g *OpenAIGateway) Chat(ctx context.Context, req *gateway.ChatRequest) (*gateway.ChatResponse, error) {
	// TODO: Implement OpenAI chat completions API call
	return &gateway.ChatResponse{
		Content: "OpenAI stub: implement me",
	}, nil
}

// ChatStream implements streaming for OpenAI.
func (g *OpenAIGateway) ChatStream(ctx context.Context, req *gateway.ChatRequest) (<-chan gateway.StreamEvent, error) {
	// TODO: Implement OpenAI streaming
	ch := make(chan gateway.StreamEvent)
	close(ch)
	return ch, nil
}

// Info returns provider metadata.
func (g *OpenAIGateway) Info() gateway.ProviderInfo {
	return g.info
}

// Health checks API connectivity.
func (g *OpenAIGateway) Health(ctx context.Context) error {
	// TODO: Implement health check ping
	return nil
}
