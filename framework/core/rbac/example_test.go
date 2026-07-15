package rbac_test

import (
	"context"
	"fmt"

	"github.com/i56/framework/core/rbac"
)

// ExampleEnforcer demonstrates RBAC permission checking.
func ExampleEnforcer() {
	store := rbac.NewInMemPermissionStore()
	store.AddRole("warehouse-operator", map[string][]string{
		rbac.ResourceOrder:  {rbac.ActionView, rbac.ActionCreate, rbac.ActionUpdate},
		rbac.ResourceParcel: {rbac.ActionView, rbac.ActionUpdate},
	})
	store.SetDataScope(rbac.ResourceOrder, rbac.ScopeWarehouse)

	enforcer := rbac.NewEnforcer(store)
	subject := rbac.Subject{
		UserID:       "user-1",
		TenantID:     "tenant-1",
		RoleIDs:      []string{"warehouse-operator"},
		WarehouseIDs: []string{"wh-shanghai", "wh-beijing"},
	}

	ctx := context.Background()

	// Check a specific permission
	canCreate := enforcer.Enforce(ctx, subject, rbac.ResourceOrder, rbac.ActionCreate)
	fmt.Println("Can create order:", canCreate)

	// Get data scope for a resource
	scope := enforcer.DataScope(ctx, subject, rbac.ResourceOrder)
	fmt.Println("Order data scope:", scope)

	// Apply data scope filter
	filter := rbac.ApplyDataScope(subject, scope)
	fmt.Println("Needs warehouse filter:", filter.NeedsWarehouseFilter())
	fmt.Println("Warehouse IDs:", filter.WarehouseIDs)
	// Output:
	// Can create order: true
	// Order data scope: warehouse
	// Needs warehouse filter: true
	// Warehouse IDs: [wh-shanghai wh-beijing]
}

// ExampleDataScopeFilter demonstrates SQL query filtering.
func ExampleDataScopeFilter() {
	subject := rbac.Subject{
		UserID:       "user-1",
		TenantID:     "t1",
		WarehouseIDs: []string{"wh-1", "wh-2"},
	}

	filter := rbac.ApplyDataScope(subject, rbac.ScopeWarehouse)

	// Use filter.TenantID and filter.WarehouseIDs in SQL WHERE clauses
	fmt.Println("WHERE tenant_id =", filter.TenantID)
	fmt.Println("AND warehouse_id IN", filter.WarehouseIDs)
	// Output:
	// WHERE tenant_id = t1
	// AND warehouse_id IN [wh-1 wh-2]
}
