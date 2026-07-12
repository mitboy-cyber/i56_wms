// Package router provides Go 1.22+ method-based routing with OpenAPI support.
package router

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"sync"
)

// RouteDoc describes a registered route for OpenAPI generation.
type RouteDoc struct {
	Method  string     `json:"method"`
	Path    string     `json:"path"`
	Summary string     `json:"summary,omitempty"`
	Tags    []string   `json:"tags,omitempty"`
	Params  []ParamDoc `json:"params,omitempty"`
}

// ParamDoc describes a route parameter.
type ParamDoc struct {
	Name        string `json:"name"`
	In          string `json:"in"` // "path", "query", "header"
	Required    bool   `json:"required"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"` // "string", "integer", "boolean"
}

// OpenAPIGenerator produces OpenAPI 3.0 specs from registered routes.
type OpenAPIGenerator struct {
	mu     sync.RWMutex
	routes []RouteDoc
	info   OpenAPIInfo
}

// OpenAPIInfo holds the OpenAPI document info section.
type OpenAPIInfo struct {
	Title       string `json:"title"`
	Version     string `json:"version"`
	Description string `json:"description,omitempty"`
}

// NewOpenAPIGenerator creates a new OpenAPI generator.
func NewOpenAPIGenerator(info OpenAPIInfo) *OpenAPIGenerator {
	return &OpenAPIGenerator{
		routes: make([]RouteDoc, 0),
		info:   info,
	}
}

// RegisterRoute adds a route to the OpenAPI documentation.
func (g *OpenAPIGenerator) RegisterRoute(doc RouteDoc) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.routes = append(g.routes, doc)
}

// Routes returns a copy of all registered routes.
func (g *OpenAPIGenerator) Routes() []RouteDoc {
	g.mu.RLock()
	defer g.mu.RUnlock()
	result := make([]RouteDoc, len(g.routes))
	copy(result, g.routes)
	return result
}

// GenerateOpenAPI produces the OpenAPI 3.0 JSON spec.
func (g *OpenAPIGenerator) GenerateOpenAPI(host string) map[string]interface{} {
	g.mu.RLock()
	routes := make([]RouteDoc, len(g.routes))
	copy(routes, g.routes)
	g.mu.RUnlock()

	// Sort routes by path for readability
	sort.Slice(routes, func(i, j int) bool {
		if routes[i].Path == routes[j].Path {
			return routes[i].Method < routes[j].Method
		}
		return routes[i].Path < routes[j].Path
	})

	// Build paths map
	paths := make(map[string]map[string]interface{})
	for _, route := range routes {
		pathKey := normalizeOpenAPIPath(route.Path)
		if _, ok := paths[pathKey]; !ok {
			paths[pathKey] = make(map[string]interface{})
		}

		params := make([]map[string]interface{}, 0)
		for _, p := range route.Params {
			params = append(params, map[string]interface{}{
				"name":        p.Name,
				"in":          p.In,
				"required":    p.Required,
				"description": p.Description,
				"schema": map[string]interface{}{
					"type": p.Type,
				},
			})
		}

		operation := map[string]interface{}{
			"summary":  route.Summary,
			"tags":     route.Tags,
			"responses": map[string]interface{}{
				"200": map[string]interface{}{
					"description": "Successful response",
				},
			},
		}
		if len(params) > 0 {
			operation["parameters"] = params
		}

		paths[pathKey][strings.ToLower(route.Method)] = operation
	}

	spec := map[string]interface{}{
		"openapi": "3.0.3",
		"info": map[string]interface{}{
			"title":       g.info.Title,
			"version":     g.info.Version,
			"description": g.info.Description,
		},
		"servers": []map[string]interface{}{
			{
				"url":         host,
				"description": "I56 Framework API Server",
			},
		},
		"paths": paths,
	}

	return spec
}

// RegisterOpenAPIEndpoint registers the /openapi.json endpoint on the router.
func RegisterOpenAPIEndpoint(r *Router, gen *OpenAPIGenerator) {
	r.GET("/openapi.json", func(w http.ResponseWriter, req *http.Request) {
		host := req.Host
		if host == "" {
			host = "localhost:8080"
		}
		scheme := "http"
		if req.TLS != nil {
			scheme = "https"
		}
		fullHost := scheme + "://" + host

		spec := gen.GenerateOpenAPI(fullHost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(spec)
	})
}

// normalizeOpenAPIPath converts Go 1.22 path patterns to OpenAPI path templates.
// "{id}" → "{id}", stays the same since OpenAPI uses the same syntax.
func normalizeOpenAPIPath(path string) string {
	return path
}
