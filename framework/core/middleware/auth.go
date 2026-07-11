package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/i56/framework/core/auth"
	"github.com/i56/framework/core/response"
)

type contextKey string

const (
	ClaimsKey  contextKey = "claims"
	TenantKey  contextKey = "tenant_id"
	UserIDKey  contextKey = "user_id"
)

// AuthRequired creates middleware that validates JWT access tokens.
func AuthRequired(tm *auth.TokenManager) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractBearerToken(r)
			if token == "" {
				response.Error(w, nil) // 401
				return
			}

			claims, err := tm.ValidateAccessToken(token)
			if err != nil {
				response.Error(w, err)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, ClaimsKey, claims)
			ctx = context.WithValue(ctx, TenantKey, claims.TenantID)
			ctx = context.WithValue(ctx, UserIDKey, claims.Subject)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractBearerToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return ""
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return ""
	}
	return parts[1]
}
