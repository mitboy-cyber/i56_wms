package providers

import (
	"context"

	"github.com/i56/framework/ai/gateway"
)

// OllamaGateway implements gateway.Gateway for local Ollama models.
type OllamaGateway struct {
	baseURL string
	info    gateway.ProviderInfo
}

// NewOllama creates a new Ollama gateway instance.
func NewOllama(baseURL string) *OllamaGateway {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	return &OllamaGateway{
		baseURL: baseURL,
		info: gateway.ProviderInfo{
			Name:              "ollama",
			Models:            []string{"llama3.1", "mistral", "codellama"},
			DefaultModel:      "llama3.1",
			SupportsStreaming: true,
			SupportsTools:     true,
		},
	}
}

// Chat implements the Gateway interface for Ollama.
func (g *OllamaGateway) Chat(ctx context.Context, req *gateway.ChatRequest) (*gateway.ChatResponse, error) {
	// TODO: Implement Ollama chat API call
	return &gateway.ChatResponse{
		Content: "Ollama stub: implement me",
	}, nil
}

// ChatStream implements streaming for Ollama.
func (g *OllamaGateway) ChatStream(ctx context.Context, req *gateway.ChatRequest) (<-chan gateway.StreamEvent, error) {
	ch := make(chan gateway.StreamEvent)
	close(ch)
	return ch, nil
}

// Info returns provider metadata.
func (g *OllamaGateway) Info() gateway.ProviderInfo {
	return g.info
}

// Health checks API connectivity.
func (g *OllamaGateway) Health(ctx context.Context) error {
	return nil
}
