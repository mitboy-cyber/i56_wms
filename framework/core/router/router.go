// Package router provides Go 1.22+ method-based routing with prefix support.
package router

import (
	"net/http"
	"strings"

	"github.com/i56/framework/core/middleware"
)

// Router wraps http.ServeMux with prefix and middleware support.
type Router struct {
	mux    *http.ServeMux
	prefix string
	mws    []middleware.Middleware
}

// New creates a new Router.
func New() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

// WithPrefix sets a path prefix for all registered routes.
func (r *Router) WithPrefix(prefix string) *Router {
	r.prefix = strings.TrimRight(prefix, "/")
	return r
}

// Use adds middleware to the router.
func (r *Router) Use(mws ...middleware.Middleware) {
	r.mws = append(r.mws, mws...)
}

// GET registers a GET route.
func (r *Router) GET(pattern string, handler http.HandlerFunc) {
	r.register("GET", pattern, handler)
}

// POST registers a POST route.
func (r *Router) POST(pattern string, handler http.HandlerFunc) {
	r.register("POST", pattern, handler)
}

// PUT registers a PUT route.
func (r *Router) PUT(pattern string, handler http.HandlerFunc) {
	r.register("PUT", pattern, handler)
}

// PATCH registers a PATCH route.
func (r *Router) PATCH(pattern string, handler http.HandlerFunc) {
	r.register("PATCH", pattern, handler)
}

// DELETE registers a DELETE route.
func (r *Router) DELETE(pattern string, handler http.HandlerFunc) {
	r.register("DELETE", pattern, handler)
}

// Handle mounts a sub-handler (typically a sub-router) at the given pattern.
func (r *Router) Handle(pattern string, handler http.Handler) {
	r.mux.Handle(pattern, handler)
}

// ServeHTTP implements the http.Handler interface.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// register applies prefix and middleware, then registers with ServeMux.
func (r *Router) register(method, pattern string, handler http.HandlerFunc) {
	fullPattern := r.buildPattern(method, pattern)
	h := middleware.Chain(http.HandlerFunc(handler), r.mws...)
	r.mux.HandleFunc(fullPattern, h.ServeHTTP)
}

// buildPattern constructs the full method+path pattern with prefix.
// Uses strings.Cut for reliable method/path separation.
// Input: method="GET", pattern="/orders"
// With prefix="/api/v1": returns "GET /api/v1/orders"
// Without prefix: returns "GET /orders"
func (r *Router) buildPattern(method, pattern string) string {
	if r.prefix == "" {
		return method + " " + pattern
	}
	fullPath := r.prefix + pattern
	return method + " " + fullPath
}
