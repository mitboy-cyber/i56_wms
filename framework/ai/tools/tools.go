// Package tools provides automatic tool registration from Go struct types.
// Struct tags are parsed to generate JSON Schema definitions compatible with
// OpenAI function calling, Anthropic tool use, and other LLM tool APIs.
package tools

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// ToolMeta captures the tool's identity and schema.
type ToolMeta struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
	Handler     any            `json:"-"` // The tool implementation (func or struct method)
}

// ToolCall represents an invocation of a registered tool.
type ToolCall struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Arguments string `json:"arguments"` // JSON-encoded arguments
}

// ToolResult is the outcome of a tool execution.
type ToolResult struct {
	CallID string `json:"call_id"`
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

// Registry manages tool registration, schema generation, and invocation.
type Registry struct {
	mu    sync.RWMutex
	tools map[string]*ToolMeta
}

// NewRegistry creates an empty tool registry.
func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]*ToolMeta),
	}
}

// JSONSchemaType maps Go types to JSON Schema type strings.
var JSONSchemaType = map[reflect.Kind]string{
	reflect.String:  "string",
	reflect.Int:     "integer",
	reflect.Int8:    "integer",
	reflect.Int16:   "integer",
	reflect.Int32:   "integer",
	reflect.Int64:   "integer",
	reflect.Float32: "number",
	reflect.Float64: "number",
	reflect.Bool:    "boolean",
	reflect.Slice:   "array",
	reflect.Map:     "object",
	reflect.Struct:  "object",
}

// toolTag is the struct tag key for tool metadata.
const toolTag = "tool"

// RegisterFromStruct parses a Go struct and registers it as a tool.
// It inspects the struct's fields and `tool` tags to generate a JSON Schema.
//
// Fields tagged with `tool:"description=..."` contribute to the schema.
// The struct's doc comment should serve as the tool description.
func (r *Registry) RegisterFromStruct(v any) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("tools: RegisterFromStruct requires a struct, got %s", t.Kind())
	}

	name := toSnakeCase(t.Name())
	meta := &ToolMeta{
		Name:        name,
		Description: name + " tool", // TODO: extract from struct doc comment
		Parameters:  buildJSONSchema(t),
		Handler:     v,
	}
	r.tools[name] = meta
	return nil
}

// Register adds a named tool with a manual schema.
func (r *Registry) Register(name string, meta *ToolMeta) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tools[name] = meta
}

// Get retrieves a tool by name.
func (r *Registry) Get(name string) (*ToolMeta, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tools[name]
	return t, ok
}

// List returns all registered tool names.
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.tools))
	for n := range r.tools {
		names = append(names, n)
	}
	return names
}

// ToOpenAISchemas converts all registered tools to OpenAI-compatible function schemas.
func (r *Registry) ToOpenAISchemas() []map[string]any {
	r.mu.RLock()
	defer r.mu.RUnlock()

	schemas := make([]map[string]any, 0, len(r.tools))
	for _, meta := range r.tools {
		schemas = append(schemas, map[string]any{
			"type": "function",
			"function": map[string]any{
				"name":        meta.Name,
				"description": meta.Description,
				"parameters":  meta.Parameters,
			},
		})
	}
	return schemas
}

// Call invokes a registered tool by name with JSON-encoded arguments.
func (r *Registry) Call(name string, argsJSON string) (*ToolResult, error) {
	r.mu.RLock()
	meta, ok := r.tools[name]
	r.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("tools: unknown tool %q", name)
	}

	// TODO: Implement reflection-based invocation using meta.Handler + args
	_ = meta
	_ = argsJSON

	return &ToolResult{
		CallID: name,
		Output: fmt.Sprintf("stub: tool %q called with args: %s", name, argsJSON),
	}, nil
}

// Count returns the number of registered tools.
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.tools)
}

// buildJSONSchema builds a JSON Schema object from a struct type.
func buildJSONSchema(t reflect.Type) map[string]any {
	props := make(map[string]any)
	required := []string{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		name := fieldName(field)
		fieldSchema := buildFieldSchema(field)

		props[name] = fieldSchema

		// Fields without `omitempty` or with `required` tag are required.
		tag := field.Tag.Get(toolTag)
		if !strings.Contains(tag, "omitempty") {
			required = append(required, name)
		}
	}

	schema := map[string]any{
		"type":       "object",
		"properties": props,
	}
	if len(required) > 0 {
		schema["required"] = required
	}
	return schema
}

func buildFieldSchema(field reflect.StructField) map[string]any {
	kind := field.Type.Kind()
	schema := map[string]any{}

	// Map Go kind to JSON Schema type
	if jsType, ok := JSONSchemaType[kind]; ok {
		schema["type"] = jsType
	} else {
		schema["type"] = "string"
	}

	// Parse tool tag for description
	tag := field.Tag.Get(toolTag)
	if tag != "" {
		parts := parseTag(tag)
		if desc, ok := parts["description"]; ok {
			schema["description"] = desc
		}
	}

	// Use json tag for name
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		jsonTag = toSnakeCase(field.Name)
	} else {
		// Strip options like ",omitempty"
		if idx := strings.Index(jsonTag, ","); idx >= 0 {
			jsonTag = jsonTag[:idx]
		}
	}

	// Handle nested structs
	if kind == reflect.Struct && field.Type != reflect.TypeOf(json.RawMessage{}) {
		nested := buildJSONSchema(field.Type)
		schema = nested
	}

	// Handle slices of structs
	if kind == reflect.Slice {
		elemKind := field.Type.Elem().Kind()
		if elemKind == reflect.Struct {
			schema["items"] = buildJSONSchema(field.Type.Elem())
		} else if jsType, ok := JSONSchemaType[elemKind]; ok {
			schema["items"] = map[string]any{"type": jsType}
		}
	}

	return schema
}

func fieldName(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return toSnakeCase(field.Name)
	}
	if idx := strings.Index(jsonTag, ","); idx >= 0 {
		return jsonTag[:idx]
	}
	return jsonTag
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

func parseTag(tag string) map[string]string {
	result := make(map[string]string)
	for _, pair := range strings.Split(tag, ",") {
		kv := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(kv) == 2 {
			result[kv[0]] = kv[1]
		}
	}
	return result
}
