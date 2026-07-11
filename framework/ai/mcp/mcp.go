// Package mcp implements the Model Context Protocol (MCP) runtime.
// It provides both server and client roles using JSON-RPC 2.0 over
// stdio and SSE transports, enabling AI agents to discover and
// interact with external tools, resources, and prompts.
package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

// ProtocolVersion is the MCP version implemented by this runtime.
const ProtocolVersion = "2024-11-05"

// Transport is the wire protocol for MCP communication.
type Transport string

const (
	TransportStdio Transport = "stdio"
	TransportSSE   Transport = "sse"
)

// JSONRPCMessage is the base JSON-RPC 2.0 envelope.
type JSONRPCMessage struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
}

// JSONRPCError is a JSON-RPC 2.0 error object.
type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// ToolDescriptor describes a tool that an MCP server exposes.
type ToolDescriptor struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"inputSchema"`
}

// ResourceDescriptor describes a resource exposed by an MCP server.
type ResourceDescriptor struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	MIMEType    string `json:"mimeType,omitempty"`
}

// PromptDescriptor describes a prompt template exposed by an MCP server.
type PromptDescriptor struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Arguments   []PromptArg    `json:"arguments,omitempty"`
}

// PromptArg describes an argument for a prompt template.
type PromptArg struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
}

// ToolHandler is the function signature for tool implementations.
type ToolHandler func(ctx context.Context, args json.RawMessage) (string, error)

// ServerCapabilities describes the server's feature set.
type ServerCapabilities struct {
	Tools     bool `json:"tools,omitempty"`
	Resources bool `json:"resources,omitempty"`
	Prompts   bool `json:"prompts,omitempty"`
}

// ServerInfo carries the server's identity metadata.
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Registry stores registered tools, resources, and prompts.
type Registry struct {
	mu        sync.RWMutex
	tools     map[string]*ToolDescriptor
	resources map[string]*ResourceDescriptor
	prompts   map[string]*PromptDescriptor
	handlers  map[string]ToolHandler
}

// NewRegistry creates an empty MCP registry.
func NewRegistry() *Registry {
	return &Registry{
		tools:     make(map[string]*ToolDescriptor),
		resources: make(map[string]*ResourceDescriptor),
		prompts:   make(map[string]*PromptDescriptor),
		handlers:  make(map[string]ToolHandler),
	}
}

// RegisterTool registers a tool with its handler.
func (r *Registry) RegisterTool(t ToolDescriptor, handler ToolHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tools[t.Name] = &t
	r.handlers[t.Name] = handler
}

// RegisterResource registers a resource.
func (r *Registry) RegisterResource(res ResourceDescriptor) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.resources[res.URI] = &res
}

// RegisterPrompt registers a prompt template.
func (r *Registry) RegisterPrompt(p PromptDescriptor) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.prompts[p.Name] = &p
}

// CallTool invokes a registered tool by name.
func (r *Registry) CallTool(ctx context.Context, name string, args json.RawMessage) (string, error) {
	r.mu.RLock()
	handler, ok := r.handlers[name]
	r.mu.RUnlock()
	if !ok {
		return "", fmt.Errorf("mcp: unknown tool %q", name)
	}
	return handler(ctx, args)
}

// ListTools returns all registered tool descriptors.
func (r *Registry) ListTools() []ToolDescriptor {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]ToolDescriptor, 0, len(r.tools))
	for _, t := range r.tools {
		out = append(out, *t)
	}
	return out
}

// ListResources returns all registered resource descriptors.
func (r *Registry) ListResources() []ResourceDescriptor {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]ResourceDescriptor, 0, len(r.resources))
	for _, res := range r.resources {
		out = append(out, *res)
	}
	return out
}

// ListPrompts returns all registered prompt descriptors.
func (r *Registry) ListPrompts() []PromptDescriptor {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]PromptDescriptor, 0, len(r.prompts))
	for _, p := range r.prompts {
		out = append(out, *p)
	}
	return out
}

// Server represents an MCP server instance.
type Server struct {
	Info         ServerInfo
	Capabilities ServerCapabilities
	Registry     *Registry
	Transport    Transport
}

// NewServer creates an MCP server.
func NewServer(info ServerInfo) *Server {
	return &Server{
		Info:     info,
		Registry: NewRegistry(),
		Capabilities: ServerCapabilities{
			Tools:     true,
			Resources: true,
			Prompts:   true,
		},
		Transport: TransportStdio,
	}
}

// Start begins listening on the configured transport.
func (s *Server) Start(ctx context.Context) error {
	// TODO: Implement transport-specific startup (stdio pipe or SSE listener)
	return nil
}

// Stop gracefully shuts down the server.
func (s *Server) Stop(ctx context.Context) error {
	// TODO: Implement graceful shutdown
	return nil
}

// Client represents an MCP client that connects to an MCP server.
type Client struct {
	ServerInfo ServerInfo
	Transport  Transport
	Endpoint   string // For SSE transport
}

// NewClient creates an MCP client.
func NewClient(transport Transport, endpoint string) *Client {
	return &Client{
		Transport: transport,
		Endpoint:  endpoint,
	}
}

// Connect establishes a connection to the MCP server.
func (c *Client) Connect(ctx context.Context) error {
	// TODO: Implement connection handshake (initialize + initialized)
	return nil
}

// Disconnect closes the connection.
func (c *Client) Disconnect(ctx context.Context) error {
	// TODO: Implement disconnection
	return nil
}

// ListTools retrieves available tools from the server.
func (c *Client) ListTools(ctx context.Context) ([]ToolDescriptor, error) {
	// TODO: Implement tools/list RPC
	return nil, nil
}

// CallTool invokes a tool on the server.
func (c *Client) CallTool(ctx context.Context, name string, args json.RawMessage) (string, error) {
	// TODO: Implement tools/call RPC
	return "", nil
}

// ListResources retrieves available resources from the server.
func (c *Client) ListResources(ctx context.Context) ([]ResourceDescriptor, error) {
	// TODO: Implement resources/list RPC
	return nil, nil
}

// ListPrompts retrieves available prompts from the server.
func (c *Client) ListPrompts(ctx context.Context) ([]PromptDescriptor, error) {
	// TODO: Implement prompts/list RPC
	return nil, nil
}
