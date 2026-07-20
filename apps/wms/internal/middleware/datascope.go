// Package middleware provides DataScope enrichment for query filtering.
package middleware

import (
	"context"
	"net/http"

	"github.com/i56/framework/core/rbac"
)

type scopeKey struct{}

// ScopeFromContext extracts DataScope from context (set by middleware).
func ScopeFromContext(ctx context.Context) (rbac.DataScope, []string) {
	if v, ok := ctx.Value(scopeKey{}).(*scopeCtx); ok {
		return v.Scope, v.WarehouseIDs
	}
	return rbac.ScopeAll, nil
}

type scopeCtx struct {
	Scope        rbac.DataScope
	WarehouseIDs []string
}

// DataScopeMiddleware injects DataScope into request context.
// Uses the Enforcer's DataScope method per-request based on the authenticated subject.
func DataScopeMiddleware(enforcer *rbac.Enforcer, subjectBuilder func(*http.Request) rbac.Subject) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			subject := subjectBuilder(r)
			scope := enforcer.DataScope(r.Context(), subject, "parcel")

			sc := &scopeCtx{Scope: scope, WarehouseIDs: subject.WarehouseIDs}
			ctx := context.WithValue(r.Context(), scopeKey{}, sc)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
