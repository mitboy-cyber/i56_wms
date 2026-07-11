// Package context manages ambient execution context that is automatically
// injected into AI requests. It carries tenant, warehouse, and user identity
// through the AI pipeline without requiring manual propagation.
package context

import (
	"context"
	"sync"
	"time"
)

// AmbientContext carries business-level metadata that the AI layer needs
// for authorization, scoping, and personalization.
type AmbientContext struct {
	// TenantID identifies the multi-tenant account.
	TenantID string `json:"tenant_id"`

	// WarehouseID is the active warehouse within a tenant.
	WarehouseID string `json:"warehouse_id"`

	// UserID identifies the end user making the request.
	UserID string `json:"user_id"`

	// Role is the user's role for RBAC enforcement.
	Role string `json:"role"`

	// SessionID ties requests to a conversation session.
	SessionID string `json:"session_id"`

	// TraceID for distributed tracing across services.
	TraceID string `json:"trace_id"`

	// Permissions is the set of granted permission codes.
	Permissions []string `json:"permissions"`

	// Metadata carries arbitrary key-value extensions.
	Metadata map[string]string `json:"metadata"`

	// CreatedAt records when this context was initialized.
	CreatedAt time.Time `json:"created_at"`
}

// contextKey is the unexported key type for storing AmbientContext in context.Context.
type contextKey struct{}

// NewAmbientContext creates a new AmbientContext with defaults.
func NewAmbientContext(tenantID string) *AmbientContext {
	return &AmbientContext{
		TenantID:  tenantID,
		Metadata:  make(map[string]string),
		CreatedAt: time.Now(),
	}
}

// WithTenant sets the tenant.
func (a *AmbientContext) WithTenant(id string) *AmbientContext {
	a.TenantID = id
	return a
}

// WithWarehouse sets the active warehouse.
func (a *AmbientContext) WithWarehouse(id string) *AmbientContext {
	a.WarehouseID = id
	return a
}

// WithUser sets the user identity.
func (a *AmbientContext) WithUser(id, role string) *AmbientContext {
	a.UserID = id
	a.Role = role
	return a
}

// WithSession sets the session identifier.
func (a *AmbientContext) WithSession(id string) *AmbientContext {
	a.SessionID = id
	return a
}

// WithTrace sets the trace ID.
func (a *AmbientContext) WithTrace(id string) *AmbientContext {
	a.TraceID = id
	return a
}

// WithPermissions sets the granted permissions.
func (a *AmbientContext) WithPermissions(perms ...string) *AmbientContext {
	a.Permissions = perms
	return a
}

// HasPermission checks if a specific permission is granted.
func (a *AmbientContext) HasPermission(perm string) bool {
	for _, p := range a.Permissions {
		if p == perm {
			return true
		}
	}
	return false
}

// ToContext stores the AmbientContext inside a Go context.Context.
func (a *AmbientContext) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKey{}, a)
}

// FromContext retrieves the AmbientContext from a Go context.Context.
func FromContext(ctx context.Context) *AmbientContext {
	if a, ok := ctx.Value(contextKey{}).(*AmbientContext); ok {
		return a
	}
	return nil
}

// Clone returns a deep copy.
func (a *AmbientContext) Clone() *AmbientContext {
	clone := *a
	clone.Permissions = make([]string, len(a.Permissions))
	copy(clone.Permissions, a.Permissions)
	clone.Metadata = make(map[string]string, len(a.Metadata))
	for k, v := range a.Metadata {
		clone.Metadata[k] = v
	}
	return &clone
}

// Manager provides a store for AmbientContext instances keyed by session or request.
type Manager struct {
	mu      sync.RWMutex
	store   map[string]*AmbientContext // sessionID → AmbientContext
	defaults *AmbientContext            // fallback when no session found
}

// NewManager creates a Context manager with optional defaults.
func NewManager(defaultTenant string) *Manager {
	return &Manager{
		store: make(map[string]*AmbientContext),
		defaults: &AmbientContext{
			TenantID:  defaultTenant,
			Metadata:  make(map[string]string),
		},
	}
}

// Set stores an AmbientContext for a session.
func (m *Manager) Set(sessionID string, ac *AmbientContext) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.store[sessionID] = ac
}

// Get retrieves the AmbientContext for a session, falling back to defaults.
func (m *Manager) Get(sessionID string) *AmbientContext {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if ac, ok := m.store[sessionID]; ok {
		return ac
	}
	return m.defaults
}

// Delete removes a session's AmbientContext.
func (m *Manager) Delete(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.store, sessionID)
}

// AutoInject wraps a standard context with the AmbientContext for the given session.
func (m *Manager) AutoInject(ctx context.Context, sessionID string) context.Context {
	ac := m.Get(sessionID)
	return ac.ToContext(ctx)
}
