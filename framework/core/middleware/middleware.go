// Package middleware provides HTTP middleware components.
package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/i56/framework/core/logger"
	"github.com/i56/framework/core/response"
)

// Middleware is a function that wraps an http.Handler.
type Middleware func(http.Handler) http.Handler

// Chain applies middlewares in order — first middleware in list is outermost
// (runs first on request, last on response).
func Chain(handler http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		handler = mws[i](handler)
	}
	return handler
}

// Recovery catches panics and returns 500.
func Recovery(log logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.Error("panic recovered",
						"error", rec,
						"stack", string(debug.Stack()),
					)
					response.Error(w, nil) // nil → 500 internal
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// RequestID injects a unique request ID into the context and response header.
func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Use UUID in production; for now use timestamp+hash
			reqID := r.Header.Get("X-Request-ID")
			if reqID == "" {
				reqID = generateRequestID()
			}
			w.Header().Set("X-Request-ID", reqID)
			ctx := context.WithValue(r.Context(), "request_id", reqID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Logger logs each request with method, path, status, and duration.
func Logger(log logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			// Wrap ResponseWriter to capture status code
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(wrapped, r)
			log.Info("request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", wrapped.statusCode),
				slog.Duration("duration", time.Since(start)),
			)
		})
	}
}

// CORS adds Cross-Origin Resource Sharing headers.
func CORS(opts *CORSOptions) Middleware {
	if opts == nil {
		opts = &CORSOptions{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			allowed := false
			for _, o := range opts.AllowedOrigins {
				if o == "*" || o == origin {
					allowed = true
					break
				}
			}
			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", join(opts.AllowedMethods))
				w.Header().Set("Access-Control-Allow-Headers", join(opts.AllowedHeaders))
				if opts.AllowCredentials {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
			}
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// CORSOptions configures CORS behavior.
type CORSOptions struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + randomHex(8)
}

func randomHex(n int) string {
	const letters = "0123456789abcdef"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}

func join(ss []string) string {
	result := ""
	for i, s := range ss {
		if i > 0 {
			result += ", "
		}
		result += s
	}
	return result
}
