// Package prompt implements a dynamic prompt composition engine.
// It assembles the final prompt from multiple layers: base system prompt,
// tenant overrides, ambient context, user memory, and tool definitions.
package prompt

import (
	"strings"

	"github.com/i56/framework/ai/context"
	"github.com/i56/framework/ai/gateway"
	"github.com/i56/framework/ai/memory"
)

// Layer represents a composable prompt layer with a name and content.
type Layer struct {
	Name     string `json:"name"`
	Content  string `json:"content"`
	Priority int    `json:"priority"` // Lower number = earlier in prompt
}

// Engine composes prompts from multiple layers.
type Engine struct {
	basePrompt   string
	overrides    []Layer
	toolSchemas  string
	systemPrompt string // Cached fully composed system prompt
}

// NewEngine creates a prompt engine with a base system prompt.
func NewEngine(basePrompt string) *Engine {
	return &Engine{
		basePrompt: basePrompt,
	}
}

// SetBase replaces the base system prompt.
func (e *Engine) SetBase(prompt string) {
	e.basePrompt = prompt
	e.systemPrompt = "" // Invalidate cache
}

// AddOverride adds a tenant-specific or custom override layer.
func (e *Engine) AddOverride(layer Layer) {
	e.overrides = append(e.overrides, layer)
	e.systemPrompt = "" // Invalidate cache
}

// ClearOverrides removes all override layers.
func (e *Engine) ClearOverrides() {
	e.overrides = nil
	e.systemPrompt = ""
}

// ComposeSystemPrompt builds the full system prompt from all layers.
// Layers are assembled in this order:
//
//	Base → Overrides (by priority) → Context → Memory → Tools
func (e *Engine) ComposeSystemPrompt(ac *context.AmbientContext, um *memory.UserMemory, toolDefs string) string {
	var sb strings.Builder

	// 1. Base system prompt
	sb.WriteString(e.basePrompt)
	sb.WriteString("\n\n")

	// 2. Override layers (sorted by priority)
	for _, layer := range e.overrides {
		if layer.Content != "" {
			sb.WriteString("## ")
			sb.WriteString(layer.Name)
			sb.WriteString("\n")
			sb.WriteString(layer.Content)
			sb.WriteString("\n\n")
		}
	}

	// 3. Ambient context injection
	if ac != nil {
		sb.WriteString("## Current Context\n")
		sb.WriteString("- Tenant: " + ac.TenantID + "\n")
		if ac.WarehouseID != "" {
			sb.WriteString("- Warehouse: " + ac.WarehouseID + "\n")
		}
		if ac.Role != "" {
			sb.WriteString("- Role: " + ac.Role + "\n")
		}
		sb.WriteString("\n")
	}

	// 4. User memory
	if um != nil && len(um.Facts) > 0 {
		sb.WriteString("## User Context\n")
		for _, fact := range um.Facts {
			sb.WriteString("- " + fact + "\n")
		}
		sb.WriteString("\n")
	}

	// 5. Tool definitions
	if toolDefs != "" {
		sb.WriteString("## Available Tools\n")
		sb.WriteString(toolDefs)
		sb.WriteString("\n")
	}

	e.systemPrompt = sb.String()
	return e.systemPrompt
}

// BuildMessages constructs the full message array for a chat request.
// It prepends the system prompt (if set) to the user-provided messages.
func (e *Engine) BuildMessages(userMessages []gateway.Message, ac *context.AmbientContext, um *memory.UserMemory, toolDefs string) []gateway.Message {
	sysPrompt := e.ComposeSystemPrompt(ac, um, toolDefs)
	messages := make([]gateway.Message, 0, len(userMessages)+1)
	messages = append(messages, gateway.Message{
		Role:    gateway.RoleSystem,
		Content: sysPrompt,
	})
	messages = append(messages, userMessages...)
	return messages
}

// LastSystemPrompt returns the most recently composed system prompt.
func (e *Engine) LastSystemPrompt() string {
	return e.systemPrompt
}
