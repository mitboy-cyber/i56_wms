// Package tenant provides multi-tenant resolution with pluggable strategies:
// shared schema (single DB, tenant_id column), database-per-tenant, and schema-per-tenant.
package tenant

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/i56/framework/core/errors"
)

// ---------------------------------------------------------------------------
// Strategy
// ---------------------------------------------------------------------------

// Strategy defines how tenant isolation is implemented at the data layer.
type Strategy string

const (
	// StrategyShared uses a single database/schema with a tenant_id column.
	StrategyShared Strategy = "shared"
	// StrategyDatabase uses a separate database per tenant.
	StrategyDatabase Strategy = "database"
	// StrategySchema uses a separate schema per tenant in the same database.
	StrategySchema Strategy = "schema"
)

// ---------------------------------------------------------------------------
// TenantInfo
// ---------------------------------------------------------------------------

// TenantInfo holds resolved tenant information.
type TenantInfo struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Strategy Strategy `json:"strategy"`
	// DSN template: e.g. "postgres://user:pass@host/%s?sslmode=disable"
	//    where %s is replaced with the database or schema name.
	DSNTemplate string `json:"dsn_template,omitempty"`
	// DatabaseName used when Strategy is database-per-tenant.
	DatabaseName string `json:"database_name,omitempty"`
	// SchemaName used when Strategy is schema-per-tenant.
	SchemaName string `json:"schema_name,omitempty"`
}

// ---------------------------------------------------------------------------
// Context helpers
// ---------------------------------------------------------------------------

type contextKey string

const tenantKey contextKey = "tenant"

// WithContext stores tenant info in a context.
func WithContext(ctx context.Context, info *TenantInfo) context.Context {
	return context.WithValue(ctx, tenantKey, info)
}

// FromContext extracts tenant info from context.
func FromContext(ctx context.Context) *TenantInfo {
	info, _ := ctx.Value(tenantKey).(*TenantInfo)
	return info
}

// MustFromContext extracts tenant info or panics.
func MustFromContext(ctx context.Context) *TenantInfo {
	info := FromContext(ctx)
	if info == nil {
		panic("tenant not found in context")
	}
	return info
}

// ---------------------------------------------------------------------------
// Resolver interface
// ---------------------------------------------------------------------------

// Resolver resolves tenant from an HTTP request.
type Resolver interface {
	Resolve(r *http.Request) (*TenantInfo, error)
}

// ---------------------------------------------------------------------------
// TenantStore
// ---------------------------------------------------------------------------

// TenantStore is the interface for looking up tenant configuration.
type TenantStore interface {
	// GetTenant retrieves tenant configuration by ID.
	GetTenant(ctx context.Context, tenantID string) (*TenantInfo, error)
	// ListTenants returns all tenants (for admin use).
	ListTenants(ctx context.Context) ([]*TenantInfo, error)
}

// ---------------------------------------------------------------------------
// Resolver implementations
// ---------------------------------------------------------------------------

// HeaderResolver resolves tenant from a header (e.g., X-Tenant-ID).
type HeaderResolver struct {
	headerName string
	store      TenantStore
}

// NewHeaderResolver creates a header-based resolver.
func NewHeaderResolver(headerName string, store TenantStore) *HeaderResolver {
	if headerName == "" {
		headerName = "X-Tenant-ID"
	}
	return &HeaderResolver{headerName: headerName, store: store}
}

func (r *HeaderResolver) Resolve(req *http.Request) (*TenantInfo, error) {
	id := req.Header.Get(r.headerName)
	if id == "" {
		return nil, errors.NewUnauthorized("missing tenant header: " + r.headerName)
	}
	return r.resolveTenant(req.Context(), id)
}

// SubdomainResolver resolves tenant from subdomain (tenant.example.com).
type SubdomainResolver struct {
	store TenantStore
}

// NewSubdomainResolver creates a subdomain-based resolver.
func NewSubdomainResolver(store TenantStore) *SubdomainResolver {
	return &SubdomainResolver{store: store}
}

func (r *SubdomainResolver) Resolve(req *http.Request) (*TenantInfo, error) {
	host := req.Host
	// Strip port if present
	if i := strings.Index(host, ":"); i != -1 {
		host = host[:i]
	}

	// Extract first subdomain component
	for i := 0; i < len(host); i++ {
		if host[i] == '.' {
			id := host[:i]
			return r.resolveTenant(req.Context(), id)
		}
	}
	return nil, errors.NewUnauthorized("cannot resolve tenant from host: " + host)
}

// JWTAuthResolver resolves tenant from a JWT claims object (interface for DI).
type JWTAuthResolver struct {
	store TenantStore
}

// NewJWTAuthResolver creates a JWT-auth resolver. The caller injects the
// tenant ID extraction logic via middleware.
func NewJWTAuthResolver(store TenantStore) *JWTAuthResolver {
	return &JWTAuthResolver{store: store}
}

func (r *JWTAuthResolver) ResolveFromID(ctx context.Context, tenantID string) (*TenantInfo, error) {
	return r.resolveTenant(ctx, tenantID)
}

// PathResolver resolves tenant from a URL path segment (/api/{tenant}/...).
type PathResolver struct {
	store     TenantStore
	pathIndex int // 0-based index in the path
}

// NewPathResolver creates a path-based resolver. pathIndex is the segment
// position after splitting by "/". Default 0 = first segment after "/".
func NewPathResolver(store TenantStore, pathIndex int) *PathResolver {
	return &PathResolver{store: store, pathIndex: pathIndex}
}

func (r *PathResolver) Resolve(req *http.Request) (*TenantInfo, error) {
	path := strings.Trim(req.URL.Path, "/")
	parts := strings.Split(path, "/")
	if r.pathIndex >= len(parts) {
		return nil, errors.NewUnauthorized("cannot resolve tenant from path: " + req.URL.Path)
	}
	id := parts[r.pathIndex]
	return r.resolveTenant(req.Context(), id)
}

// MultiResolver tries multiple resolvers in order, returning the first success.
type MultiResolver struct {
	resolvers []Resolver
}

// NewMultiResolver creates a multi-resolver that chains multiple resolvers.
func NewMultiResolver(resolvers ...Resolver) *MultiResolver {
	return &MultiResolver{resolvers: resolvers}
}

func (r *MultiResolver) Resolve(req *http.Request) (*TenantInfo, error) {
	for _, resolver := range r.resolvers {
		info, err := resolver.Resolve(req)
		if err == nil {
			return info, nil
		}
	}
	return nil, errors.NewUnauthorized("unable to resolve tenant from request")
}

// resolveTenant looks up the tenant in the store (if available) or returns a minimal info.
func (r *HeaderResolver) resolveTenant(ctx context.Context, id string) (*TenantInfo, error) {
	if r.store != nil {
		return r.store.GetTenant(ctx, id)
	}
	return &TenantInfo{ID: id, Name: id, Strategy: StrategyShared}, nil
}

func (r *SubdomainResolver) resolveTenant(ctx context.Context, id string) (*TenantInfo, error) {
	if r.store != nil {
		return r.store.GetTenant(ctx, id)
	}
	return &TenantInfo{ID: id, Name: id, Strategy: StrategyShared}, nil
}

func (r *JWTAuthResolver) resolveTenant(ctx context.Context, id string) (*TenantInfo, error) {
	if r.store != nil {
		return r.store.GetTenant(ctx, id)
	}
	return &TenantInfo{ID: id, Name: id, Strategy: StrategyShared}, nil
}

func (r *PathResolver) resolveTenant(ctx context.Context, id string) (*TenantInfo, error) {
	if r.store != nil {
		return r.store.GetTenant(ctx, id)
	}
	return &TenantInfo{ID: id, Name: id, Strategy: StrategyShared}, nil
}

// ---------------------------------------------------------------------------
// Middleware
// ---------------------------------------------------------------------------

// Middleware creates an HTTP middleware that resolves and stores tenant in context.
func Middleware(resolver Resolver) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			info, err := resolver.Resolve(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			ctx := WithContext(r.Context(), info)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ---------------------------------------------------------------------------
// Tenant-aware DB
// ---------------------------------------------------------------------------

// DSNBuilder constructs a database DSN based on tenant strategy.
type DSNBuilder struct {
	baseDSN string // e.g. "postgres://user:pass@host:5432/%s?sslmode=disable"
	strategy Strategy
}

// NewDSNBuilder creates a new DSN builder.
func NewDSNBuilder(baseDSN string, strategy Strategy) *DSNBuilder {
	return &DSNBuilder{baseDSN: baseDSN, strategy: strategy}
}

// BuildDSN constructs the tenant-specific DSN.
func (b *DSNBuilder) BuildDSN(info *TenantInfo) string {
	switch b.strategy {
	case StrategyShared:
		return b.baseDSN
	case StrategyDatabase:
		dbName := info.DatabaseName
		if dbName == "" {
			dbName = "tenant_" + info.ID
		}
		return fmt.Sprintf(b.baseDSN, dbName)
	case StrategySchema:
		schemaName := info.SchemaName
		if schemaName == "" {
			schemaName = "tenant_" + info.ID
		}
		return fmt.Sprintf(b.baseDSN, schemaName) + "&search_path=" + schemaName
	default:
		return b.baseDSN
	}
}

// ---------------------------------------------------------------------------
// In-memory Tenant Store
// ---------------------------------------------------------------------------

// InMemTenantStore is an in-memory implementation of TenantStore.
type InMemTenantStore struct {
	tenants map[string]*TenantInfo
}

// NewInMemTenantStore creates an in-memory tenant store.
func NewInMemTenantStore() *InMemTenantStore {
	return &InMemTenantStore{tenants: make(map[string]*TenantInfo)}
}

// AddTenant registers a tenant.
func (s *InMemTenantStore) AddTenant(info *TenantInfo) {
	s.tenants[info.ID] = info
}

func (s *InMemTenantStore) GetTenant(ctx context.Context, tenantID string) (*TenantInfo, error) {
	info, ok := s.tenants[tenantID]
	if !ok {
		return nil, fmt.Errorf("tenant %q not found", tenantID)
	}
	return info, nil
}

func (s *InMemTenantStore) ListTenants(ctx context.Context) ([]*TenantInfo, error) {
	result := make([]*TenantInfo, 0, len(s.tenants))
	for _, t := range s.tenants {
		result = append(result, t)
	}
	return result, nil
}
