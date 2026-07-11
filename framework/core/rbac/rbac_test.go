package rbac

import (
	"context"
	"testing"
)

func TestEnforcer_BasicPermission(t *testing.T) {
	store := NewInMemPermissionStore()
	store.AddRole("admin", map[string][]string{
		ResourceOrder: {ActionCreate, ActionView, ActionUpdate, ActionDelete},
	})
	store.SetDataScope(ResourceOrder, ScopeAll)

	enforcer := NewEnforcer(store)
	subject := Subject{
		UserID:  "user-1",
		RoleIDs: []string{"admin"},
	}

	if !enforcer.Enforce(context.Background(), subject, ResourceOrder, ActionCreate) {
		t.Error("expected admin to have create on order")
	}
	if !enforcer.Enforce(context.Background(), subject, ResourceOrder, ActionView) {
		t.Error("expected admin to have view on order")
	}
	if enforcer.Enforce(context.Background(), subject, ResourceOrder, ActionCancel) {
		t.Error("expected admin to NOT have cancel on order")
	}
}

func TestEnforcer_HasAnyPermission(t *testing.T) {
	store := NewInMemPermissionStore()
	store.AddRole("viewer", map[string][]string{
		ResourceOrder: {ActionView},
	})
	enforcer := NewEnforcer(store)
	subject := Subject{UserID: "u1", RoleIDs: []string{"viewer"}}

	if !enforcer.HasAnyPermission(context.Background(), subject,
		[2]string{ResourceOrder, ActionView},
		[2]string{ResourceOrder, ActionDelete},
	) {
		t.Error("expected HasAnyPermission to return true")
	}

	if enforcer.HasAnyPermission(context.Background(), subject,
		[2]string{ResourceOrder, ActionDelete},
		[2]string{ResourceOrder, ActionCancel},
	) {
		t.Error("expected HasAnyPermission to return false")
	}
}

func TestEnforcer_HasAllPermissions(t *testing.T) {
	store := NewInMemPermissionStore()
	store.AddRole("operator", map[string][]string{
		ResourceOrder:  {ActionView, ActionCreate, ActionUpdate},
		ResourceParcel: {ActionView},
	})
	enforcer := NewEnforcer(store)
	subject := Subject{UserID: "u1", RoleIDs: []string{"operator"}}

	if !enforcer.HasAllPermissions(context.Background(), subject,
		[2]string{ResourceOrder, ActionView},
		[2]string{ResourceParcel, ActionView},
	) {
		t.Error("expected HasAllPermissions to return true")
	}

	if enforcer.HasAllPermissions(context.Background(), subject,
		[2]string{ResourceOrder, ActionView},
		[2]string{ResourceOrder, ActionDelete},
	) {
		t.Error("expected HasAllPermissions to return false")
	}
}

func TestDataScope_String(t *testing.T) {
	tests := []struct {
		scope    DataScope
		expected string
	}{
		{ScopeAll, "all"},
		{ScopeTenant, "tenant"},
		{ScopeWarehouse, "warehouse"},
		{ScopeDepartment, "department"},
		{ScopeSelf, "self"},
	}
	for _, tt := range tests {
		if got := tt.scope.String(); got != tt.expected {
			t.Errorf("DataScope(%d).String() = %q, want %q", tt.scope, got, tt.expected)
		}
	}
}

func TestDataScopeFilter_All(t *testing.T) {
	subject := Subject{UserID: "u1", TenantID: "t1"}
	f := ApplyDataScope(subject, ScopeAll)

	if !f.NeedsNoFilter() {
		t.Error("ScopeAll should need no filter")
	}
	if f.NeedsWarehouseFilter() {
		t.Error("ScopeAll should not need warehouse filter")
	}
}

func TestDataScopeFilter_Warehouse(t *testing.T) {
	subject := Subject{
		UserID:       "u1",
		TenantID:     "t1",
		WarehouseIDs: []string{"wh-1", "wh-2"},
	}
	f := ApplyDataScope(subject, ScopeWarehouse)

	if f.Scope != ScopeWarehouse {
		t.Errorf("expected ScopeWarehouse, got %v", f.Scope)
	}
	if !f.NeedsWarehouseFilter() {
		t.Error("ScopeWarehouse should need warehouse filter")
	}
	if len(f.WarehouseIDs) != 2 {
		t.Errorf("expected 2 warehouse IDs, got %d", len(f.WarehouseIDs))
	}
	if f.WarehouseIDs[0] != "wh-1" {
		t.Errorf("expected wh-1, got %q", f.WarehouseIDs[0])
	}
}

func TestDataScopeFilter_Self(t *testing.T) {
	subject := Subject{UserID: "u1", TenantID: "t1"}
	f := ApplyDataScope(subject, ScopeSelf)

	if f.UserID != "u1" {
		t.Errorf("expected user ID 'u1', got %q", f.UserID)
	}
	if f.Scope != ScopeSelf {
		t.Errorf("expected ScopeSelf, got %v", f.Scope)
	}
}

func TestSubject_HasWarehouseAccess(t *testing.T) {
	subject := Subject{WarehouseIDs: []string{"wh-a", "wh-b"}}

	if !subject.HasWarehouseAccess("wh-a") {
		t.Error("expected access to wh-a")
	}
	if subject.HasWarehouseAccess("wh-c") {
		t.Error("expected no access to wh-c")
	}
}

func TestSubject_HasRole(t *testing.T) {
	subject := Subject{RoleIDs: []string{"admin", "operator"}}

	if !subject.HasRole("admin") {
		t.Error("expected HasRole('admin') true")
	}
	if subject.HasRole("viewer") {
		t.Error("expected HasRole('viewer') false")
	}
}

func TestSubject_HasPermission(t *testing.T) {
	subject := Subject{Permissions: []string{"order:view", "order:create"}}

	if !subject.HasPermission("order:view") {
		t.Error("expected HasPermission('order:view') true")
	}
	if subject.HasPermission("order:delete") {
		t.Error("expected HasPermission('order:delete') false")
	}
}

func TestInMemPermissionStore_RolePermissions(t *testing.T) {
	store := NewInMemPermissionStore()
	store.AddRole("ops", map[string][]string{
		ResourceOrder: {ActionView, ActionCreate},
	})
	store.AssignRole("user-1", "ops")

	subject := Subject{UserID: "user-1", RoleIDs: []string{"ops"}}

	ok, err := store.HasPermission(context.Background(), subject, ResourceOrder, ActionView)
	if err != nil {
		t.Fatalf("HasPermission: %v", err)
	}
	if !ok {
		t.Error("expected ops to view order")
	}

	ok, err = store.HasPermission(context.Background(), subject, ResourceOrder, ActionDelete)
	if err != nil {
		t.Fatalf("HasPermission: %v", err)
	}
	if ok {
		t.Error("expected ops NOT to delete order")
	}
}

func TestInMemPermissionStore_Wildcard(t *testing.T) {
	store := NewInMemPermissionStore()
	store.AddRole("superadmin", map[string][]string{
		ResourceOrder: {"*"},
	})

	subject := Subject{UserID: "u1", RoleIDs: []string{"superadmin"}}

	if ok, _ := store.HasPermission(context.Background(), subject, ResourceOrder, ActionCreate); !ok {
		t.Error("wildcard should allow create")
	}
	if ok, _ := store.HasPermission(context.Background(), subject, ResourceOrder, ActionDelete); !ok {
		t.Error("wildcard should allow delete")
	}
	if ok, _ := store.HasPermission(context.Background(), subject, ResourceOrder, "some_custom_action"); !ok {
		t.Error("wildcard should allow any action")
	}
}
