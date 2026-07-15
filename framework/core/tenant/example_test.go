package tenant_test

import (
	"fmt"
	"net/http"

	"github.com/i56/framework/core/tenant"
)

// ExampleHeaderResolver demonstrates resolving tenant from HTTP headers.
func ExampleHeaderResolver() {
	store := tenant.NewInMemTenantStore()
	store.AddTenant(&tenant.TenantInfo{
		ID:       "acme-corp",
		Name:     "ACME Corporation",
		Strategy: tenant.StrategyShared,
	})
	store.AddTenant(&tenant.TenantInfo{
		ID:           "big-client",
		Name:         "Big Client Inc",
		Strategy:     tenant.StrategyDatabase,
		DatabaseName: "tenant_big_client",
	})

	resolver := tenant.NewHeaderResolver("X-Tenant-ID", store)

	// Simulate an incoming HTTP request
	req, _ := http.NewRequest("GET", "/api/orders", nil)
	req.Header.Set("X-Tenant-ID", "big-client")

	info, err := resolver.Resolve(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Tenant:", info.Name)
	fmt.Println("Strategy:", info.Strategy)
	fmt.Println("Database:", info.DatabaseName)
	// Output:
	// Tenant: Big Client Inc
	// Strategy: database
	// Database: tenant_big_client
}

// ExampleDSNBuilder demonstrates constructing tenant-specific DSNs.
func ExampleDSNBuilder() {
	builder := tenant.NewDSNBuilder(
		"postgres://user:***@localhost:5432/%s?sslmode=disable",
		tenant.StrategyDatabase,
	)

	info := &tenant.TenantInfo{
		ID:           "t1",
		DatabaseName: "db_tenant_t1",
	}

	dsn := builder.BuildDSN(info)
	fmt.Println(dsn)
	// Output:
	// postgres://user:***@localhost:5432/db_tenant_t1?sslmode=disable
}

// ExampleMultiResolver demonstrates trying multiple resolution strategies.
func ExampleMultiResolver() {
	store := tenant.NewInMemTenantStore()
	store.AddTenant(&tenant.TenantInfo{ID: "acme", Name: "ACME", Strategy: tenant.StrategyShared})

	headerResolver := tenant.NewHeaderResolver("X-Tenant-ID", store)
	subdomainResolver := tenant.NewSubdomainResolver(store)
	multi := tenant.NewMultiResolver(headerResolver, subdomainResolver)

	// Try with subdomain
	req, _ := http.NewRequest("GET", "/api/orders", nil)
	req.Host = "acme.example.com"

	info, err := multi.Resolve(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(info.Name)
	// Output:
	// ACME
}
