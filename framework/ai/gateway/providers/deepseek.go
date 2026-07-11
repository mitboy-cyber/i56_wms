// Package deepseek implements the Gateway interface for DeepSeek API.
// This is the first real provider implementation for I56 2.0 AI Runtime.
package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/i56/framework/ai/gateway"
)

// DeepSeekProvider implements gateway.Provider for DeepSeek API.
type DeepSeekProvider struct {
	apiKey string
	client *http.Client
}

// NewDeepSeek creates a DeepSeek provider.
func NewDeepSeek(apiKey string) *DeepSeekProvider {
	return &DeepSeekProvider{
		apiKey: apiKey,
		client: &http.Client{Timeout: 60 * time.Second},
	}
}

type dsRequest struct {
	Model    string            `json:"model"`
	Messages []dsMessage       `json:"messages"`
	Stream   bool              `json:"stream"`
}

type dsMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type dsResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message dsMessage `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type dsStreamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

// Chat implements gateway.Provider.
func (p *DeepSeekProvider) Chat(ctx context.Context, req *gateway.ChatRequest) (*gateway.ChatResponse, error) {
	dsReq := dsRequest{
		Model:  "deepseek-chat",
		Stream: false,
	}
	for _, m := range req.Messages {
		dsReq.Messages = append(dsReq.Messages, dsMessage{Role: string(m.Role), Content: m.Content})
	}

	body, _ := json.Marshal(dsReq)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", "https://api.deepseek.com/chat/completions", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return &gateway.ChatResponse{Content: fmt.Sprintf("AI 服务暂时不可用: %v", err), TokenUsage: gateway.TokenUsage{}}, nil
	}
	defer resp.Body.Close()

	var dsResp dsResponse
	if err := json.NewDecoder(resp.Body).Decode(&dsResp); err != nil {
		return &gateway.ChatResponse{Content: "API 响应解析失败，请稍后重试", TokenUsage: gateway.TokenUsage{}}, nil
	}

	if len(dsResp.Choices) == 0 {
		return &gateway.ChatResponse{Content: "未能理解您的问题，请换种方式提问", TokenUsage: gateway.TokenUsage{}}, nil
	}

	return &gateway.ChatResponse{
		Content: dsResp.Choices[0].Message.Content,
		TokenUsage: gateway.TokenUsage{
			PromptTokens:     dsResp.Usage.PromptTokens,
			CompletionTokens: dsResp.Usage.CompletionTokens,
			TotalTokens:      dsResp.Usage.TotalTokens,
		},
	}, nil
}

// ChatStream implements gateway.Provider with SSE streaming.
func (p *DeepSeekProvider) ChatStream(ctx context.Context, req *gateway.ChatRequest) (<-chan gateway.StreamEvent, error) {
	dsReq := dsRequest{
		Model:  "deepseek-chat",
		Stream: true,
	}
	for _, m := range req.Messages {
		dsReq.Messages = append(dsReq.Messages, dsMessage{Role: string(m.Role), Content: m.Content})
	}

	body, _ := json.Marshal(dsReq)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", "https://api.deepseek.com/chat/completions", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	ch := make(chan gateway.StreamEvent, 100)
	go func() {
		defer close(ch)
		resp, err := p.client.Do(httpReq)
		if err != nil {
			ch <- gateway.StreamEvent{Content: fmt.Sprintf("AI 服务暂时不可用: %v", err), Done: true}
			return
		}
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		for {
			var chunk dsStreamChunk
			if err := decoder.Decode(&chunk); err != nil {
				ch <- gateway.StreamEvent{Done: true}
				return
			}
			if len(chunk.Choices) > 0 {
				if chunk.Choices[0].FinishReason == "stop" {
					ch <- gateway.StreamEvent{Done: true}
					return
				}
				ch <- gateway.StreamEvent{Content: chunk.Choices[0].Delta.Content}
			}
		}
	}()

	return ch, nil
}

// Info implements gateway.Gateway.
func (p *DeepSeekProvider) Info() gateway.ProviderInfo {
	return gateway.ProviderInfo{
		Name:              "deepseek",
		Models:            []string{"deepseek-chat", "deepseek-reasoner"},
		DefaultModel:      "deepseek-chat",
		SupportsStreaming: true,
	}
}

// Health implements gateway.Gateway.
func (p *DeepSeekProvider) Health(ctx context.Context) error {
	return nil
}
