package tenant

import (
	"context"
	"net/http"
	"testing"
)

func TestHeaderResolver_Resolve(t *testing.T) {
	store := NewInMemTenantStore()
	store.AddTenant(&TenantInfo{ID: "t1", Name: "Tenant One", Strategy: StrategyShared})
	store.AddTenant(&TenantInfo{ID: "t2", Name: "Tenant Two", Strategy: StrategyDatabase, DatabaseName: "db_t2"})

	resolver := NewHeaderResolver("X-Tenant-ID", store)

	req, _ := http.NewRequest("GET", "/api/orders", nil)
	req.Header.Set("X-Tenant-ID", "t2")

	info, err := resolver.Resolve(req)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if info.ID != "t2" {
		t.Errorf("expected t2, got %q", info.ID)
	}
	if info.Strategy != StrategyDatabase {
		t.Errorf("expected database strategy, got %q", info.Strategy)
	}
	if info.DatabaseName != "db_t2" {
		t.Errorf("expected db_t2, got %q", info.DatabaseName)
	}
}

func TestHeaderResolver_MissingHeader(t *testing.T) {
	resolver := NewHeaderResolver("X-Tenant-ID", nil)

	req, _ := http.NewRequest("GET", "/api/orders", nil)
	_, err := resolver.Resolve(req)
	if err == nil {
		t.Error("expected error for missing header")
	}
}

func TestSubdomainResolver_Resolve(t *testing.T) {
	store := NewInMemTenantStore()
	store.AddTenant(&TenantInfo{ID: "acme", Name: "ACME Corp", Strategy: StrategyShared})

	resolver := NewSubdomainResolver(store)
	req, _ := http.NewRequest("GET", "/api/orders", nil)
	req.Host = "acme.example.com"

	info, err := resolver.Resolve(req)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if info.ID != "acme" {
		t.Errorf("expected acme, got %q", info.ID)
	}
}

func TestSubdomainResolver_HostNoSubdomain(t *testing.T) {
	resolver := NewSubdomainResolver(nil)
	req, _ := http.NewRequest("GET", "/", nil)
	req.Host = "localhost" // no dots = no subdomain

	_, err := resolver.Resolve(req)
	if err == nil {
		t.Error("expected error for host without subdomain")
	}
}

func TestPathResolver_Resolve(t *testing.T) {
	store := NewInMemTenantStore()
	store.AddTenant(&TenantInfo{ID: "wh1", Name: "Warehouse 1", Strategy: StrategyShared})

	resolver := NewPathResolver(store, 1) // second segment
	req, _ := http.NewRequest("GET", "/api/wh1/orders", nil)

	info, err := resolver.Resolve(req)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if info.ID != "wh1" {
		t.Errorf("expected wh1, got %q", info.ID)
	}
}

func TestMultiResolver(t *testing.T) {
	store := NewInMemTenantStore()
	store.AddTenant(&TenantInfo{ID: "t1", Name: "T1", Strategy: StrategyShared})

	headerRes := NewHeaderResolver("X-Tenant-ID", store)
	subdomainRes := NewSubdomainResolver(store)

	multi := NewMultiResolver(headerRes, subdomainRes)

	// Test with header
	req1, _ := http.NewRequest("GET", "/", nil)
	req1.Header.Set("X-Tenant-ID", "t1")
	req1.Host = "example.com"

	info, err := multi.Resolve(req1)
	if err != nil {
		t.Fatalf("MultiResolver with header: %v", err)
	}
	if info.ID != "t1" {
		t.Errorf("expected t1, got %q", info.ID)
	}

	// Test fallback to subdomain when no header
	req2, _ := http.NewRequest("GET", "/", nil)
	req2.Host = "t1.example.com"

	info, err = multi.Resolve(req2)
	if err != nil {
		t.Fatalf("MultiResolver with subdomain: %v", err)
	}
	if info.ID != "t1" {
		t.Errorf("expected t1, got %q", info.ID)
	}
}

func TestMultiResolver_NoMatch(t *testing.T) {
	multi := NewMultiResolver()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Host = "example.com"

	_, err := multi.Resolve(req)
	if err == nil {
		t.Error("expected error when no resolver matches")
	}
}

func TestMiddleware(t *testing.T) {
	store := NewInMemTenantStore()
	store.AddTenant(&TenantInfo{ID: "acme", Name: "ACME", Strategy: StrategyShared})

	resolver := NewHeaderResolver("X-Tenant-ID", store)
	mw := Middleware(resolver)

	called := false
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		info := FromContext(r.Context())
		if info == nil {
			t.Error("expected tenant in context")
		}
		if info.ID != "acme" {
			t.Errorf("expected 'acme', got %q", info.ID)
		}
		called = true
		w.WriteHeader(200)
	}))

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-Tenant-ID", "acme")
	rr := &testResponseWriter{}

	handler.ServeHTTP(rr, req)

	if !called {
		t.Error("handler was not called")
	}
	if rr.status != 200 {
		t.Errorf("expected status 200, got %d", rr.status)
	}
}

func TestMiddleware_InvalidTenant(t *testing.T) {
	resolver := NewHeaderResolver("X-Tenant-ID", nil)
	mw := Middleware(resolver)

	req, _ := http.NewRequest("GET", "/", nil)
	rr := &testResponseWriter{}

	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called for missing tenant")
	})).ServeHTTP(rr, req)

	if rr.status != 401 {
		t.Errorf("expected status 401, got %d", rr.status)
	}
}

func TestDSNBuilder_Shared(t *testing.T) {
	builder := NewDSNBuilder("postgres://user:pass@host/db?sslmode=disable", StrategyShared)
	info := &TenantInfo{ID: "t1"}
	dsn := builder.BuildDSN(info)
	if dsn != "postgres://user:pass@host/db?sslmode=disable" {
		t.Errorf("unexpected DSN: %s", dsn)
	}
}

func TestDSNBuilder_Database(t *testing.T) {
	builder := NewDSNBuilder("postgres://user:pass@host/%s?sslmode=disable", StrategyDatabase)
	info := &TenantInfo{ID: "t1", DatabaseName: "tenant_t1"}
	dsn := builder.BuildDSN(info)
	if dsn != "postgres://user:pass@host/tenant_t1?sslmode=disable" {
		t.Errorf("unexpected DSN: %s", dsn)
	}
}

func TestDSNBuilder_DatabaseDefaultName(t *testing.T) {
	builder := NewDSNBuilder("postgres://user:pass@host/%s?sslmode=disable", StrategyDatabase)
	info := &TenantInfo{ID: "t1"} // No DatabaseName set
	dsn := builder.BuildDSN(info)
	if dsn != "postgres://user:pass@host/tenant_t1?sslmode=disable" {
		t.Errorf("unexpected DSN: %s", dsn)
	}
}

func TestContextHelpers(t *testing.T) {
	info := &TenantInfo{ID: "test", Strategy: StrategyShared}
	ctx := WithContext(context.Background(), info)

	got := FromContext(ctx)
	if got == nil || got.ID != "test" {
		t.Errorf("FromContext: expected 'test', got %v", got)
	}

	if MustFromContext(ctx).ID != "test" {
		t.Error("MustFromContext mismatch")
	}
}

func TestMustFromContext_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	MustFromContext(context.Background())
}

func TestInMemTenantStore(t *testing.T) {
	store := NewInMemTenantStore()
	store.AddTenant(&TenantInfo{ID: "a", Name: "A"})
	store.AddTenant(&TenantInfo{ID: "b", Name: "B"})

	ctx := context.Background()

	info, err := store.GetTenant(ctx, "a")
	if err != nil {
		t.Fatalf("GetTenant: %v", err)
	}
	if info.Name != "A" {
		t.Errorf("expected 'A', got %q", info.Name)
	}

	_, err = store.GetTenant(ctx, "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent tenant")
	}

	all, err := store.ListTenants(ctx)
	if err != nil {
		t.Fatalf("ListTenants: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 tenants, got %d", len(all))
	}
}

// testResponseWriter for middleware tests
type testResponseWriter struct {
	status int
	body   []byte
}

func (rw *testResponseWriter) Header() http.Header         { return http.Header{} }
func (rw *testResponseWriter) Write(b []byte) (int, error)  { rw.body = append(rw.body, b...); return len(b), nil }
func (rw *testResponseWriter) WriteHeader(code int)          { rw.status = code }
