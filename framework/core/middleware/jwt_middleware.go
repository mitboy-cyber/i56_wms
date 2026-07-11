package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/i56/framework/core/jwt"
)

// JWTAuth creates middleware that validates JWT and injects claims into context
func JWTAuth(jwtSvc *jwt.Service) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			token := extractToken(r)
			if token == "" {
				http.Error(w, `{"error":"missing token"}`, 401)
				return
			}
			claims, err := jwtSvc.Verify(token)
			if err != nil {
				http.Error(w, `{"error":"invalid token"}`, 401)
				return
			}
			// Inject into context
			ctx := r.Context()
			ctx = context.WithValue(ctx, TenantKey, claims.TenantID)
			ctx = context.WithValue(ctx, UserIDKey, claims.Sub)
			if len(claims.WarehouseIDs) > 0 {
			}
			next(w, r.WithContext(ctx))
		}
	}
}

// JWTAuthOptional validates JWT if present, but doesn't block if missing
func JWTAuthOptional(jwtSvc *jwt.Service) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			token := extractToken(r)
			if token != "" {
				claims, err := jwtSvc.Verify(token)
				if err == nil {
					ctx := r.Context()
					ctx = context.WithValue(ctx, TenantKey, claims.TenantID)
					ctx = context.WithValue(ctx, UserIDKey, claims.Sub)
					r = r.WithContext(ctx)
				}
			}
			next(w, r)
		}
	}
}

// TenantFromContext extracts tenant_id from context
func TenantFromContext(ctx context.Context) int64 {
	if v, ok := ctx.Value(TenantKey).(int64); ok {
		return v
	}
	return 1 // default tenant
}

func extractToken(r *http.Request) string {
	// Authorization header
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	// Cookie fallback
	if ck, err := r.Cookie("i56_token"); err == nil {
		return ck.Value
	}
	// Query param
	return r.URL.Query().Get("token")
}
