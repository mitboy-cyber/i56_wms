// Package gateway provides a unified LLM abstraction layer.
// It defines the core types and interface that all AI providers implement,
// enabling protocol-agnostic access to OpenAI, Claude, DeepSeek, Ollama, and more.
package gateway

import (
	"context"
	"time"
)

// Role represents the message originator in a conversation.
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

// Message is a single turn in a conversation.
type Message struct {
	Role       Role        `json:"role"`
	Content    string      `json:"content,omitempty"`
	Name       string      `json:"name,omitempty"`
	ToolCalls  []ToolCall  `json:"tool_calls,omitempty"`
	ToolCallID string      `json:"tool_call_id,omitempty"`
}

// ToolCall represents a request from the model to invoke a tool.
type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

// FunctionCall is the function name and arguments for a tool call.
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// ToolSchema describes a tool that the model can call.
type ToolSchema struct {
	Type     string      `json:"type"`
	Function FunctionDef `json:"function"`
}

// FunctionDef defines a function's name, description, and JSON Schema parameters.
type FunctionDef struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

// TokenUsage tracks token consumption for billing and rate limiting.
type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatRequest is the unified request structure for all LLM providers.
type ChatRequest struct {
	Messages    []Message    `json:"messages"`
	Model       string       `json:"model"`
	Temperature float64      `json:"temperature,omitempty"`
	MaxTokens   int          `json:"max_tokens,omitempty"`
	Tools       []ToolSchema `json:"tools,omitempty"`
	Stream      bool         `json:"stream,omitempty"`
	// Metadata carries provider-specific extensions (e.g., top_p, stop sequences).
	Metadata map[string]any `json:"metadata,omitempty"`
}

// ChatResponse is the unified response from any LLM provider.
type ChatResponse struct {
	ID           string     `json:"id"`
	Model        string     `json:"model"`
	Content      string     `json:"content"`
	ToolCalls    []ToolCall `json:"tool_calls,omitempty"`
	TokenUsage   TokenUsage `json:"usage"`
	FinishReason string     `json:"finish_reason,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

// StreamEvent represents a single chunk in a streaming response.
type StreamEvent struct {
	Content       string    `json:"content,omitempty"`
	ToolCallDelta *ToolCall `json:"tool_call_delta,omitempty"`
	FinishReason  string    `json:"finish_reason,omitempty"`
	Done          bool      `json:"done"`
	Error         error     `json:"error,omitempty"`
}

// ProviderInfo describes a provider's capabilities and defaults.
type ProviderInfo struct {
	Name              string   `json:"name"`
	Models            []string `json:"models"`
	DefaultModel      string   `json:"default_model"`
	SupportsStreaming bool     `json:"supports_streaming"`
	SupportsTools     bool     `json:"supports_tools"`
}

// Gateway is the unified interface that every LLM provider must implement.
// It abstracts away provider-specific SDKs, authentication, and wire protocols.
type Gateway interface {
	// Chat sends a request and blocks until a complete response is available.
	Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)

	// ChatStream sends a request and returns a channel of streaming events.
	ChatStream(ctx context.Context, req *ChatRequest) (<-chan StreamEvent, error)

	// Info returns the provider's metadata and supported models.
	Info() ProviderInfo

	// Health checks connectivity to the provider's API.
	Health(ctx context.Context) error
}
