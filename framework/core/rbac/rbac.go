// Package rbac provides Role-Based Access Control with DataScope filtering.
// DataScope controls which records are visible: All, Tenant, Warehouse, Department, or Self.
// When a subject has ScopeWarehouse, their WarehouseIDs are used to filter queries.
package rbac

import (
	"context"
	"sync"
)

// ---------------------------------------------------------------------------
// DataScope
// ---------------------------------------------------------------------------

// DataScope defines the data visibility level.
type DataScope int

const (
	ScopeAll        DataScope = iota // All data (super admin)
	ScopeTenant                      // Current tenant
	ScopeWarehouse                   // Assigned warehouses (use WarehouseIDs)
	ScopeDepartment                  // Current department
	ScopeSelf                        // Own data only
)

func (s DataScope) String() string {
	switch s {
	case ScopeAll:
		return "all"
	case ScopeTenant:
		return "tenant"
	case ScopeWarehouse:
		return "warehouse"
	case ScopeDepartment:
		return "department"
	case ScopeSelf:
		return "self"
	default:
		return "unknown"
	}
}

// ---------------------------------------------------------------------------
// Subject
// ---------------------------------------------------------------------------

// Subject is the security principal making the request.
type Subject struct {
	UserID       string   `json:"user_id"`
	TenantID     string   `json:"tenant_id"`
	DeptID       string   `json:"dept_id,omitempty"`
	RoleIDs      []string `json:"role_ids"`
	Permissions  []string `json:"permissions"`
	WarehouseIDs []string `json:"warehouse_ids,omitempty"`
}

// HasWarehouseAccess checks if a warehouse ID is in the subject's assigned list.
func (s Subject) HasWarehouseAccess(warehouseID string) bool {
	for _, id := range s.WarehouseIDs {
		if id == warehouseID {
			return true
		}
	}
	return false
}

// HasRole checks whether the subject has a specific role.
func (s Subject) HasRole(roleID string) bool {
	for _, id := range s.RoleIDs {
		if id == roleID {
			return true
		}
	}
	return false
}

// HasPermission checks whether the subject has a specific permission.
func (s Subject) HasPermission(perm string) bool {
	for _, p := range s.Permissions {
		if p == perm {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// PermissionStore
// ---------------------------------------------------------------------------

// PermissionStore is the interface for checking permissions.
type PermissionStore interface {
	HasPermission(ctx context.Context, subject Subject, resource, action string) (bool, error)
	GetDataScope(ctx context.Context, subject Subject, resource string) (DataScope, error)
}

// ---------------------------------------------------------------------------
// Enforcer
// ---------------------------------------------------------------------------

// Enforcer is the RBAC enforcement point.
type Enforcer struct {
	store PermissionStore
}

// NewEnforcer creates a new Enforcer.
func NewEnforcer(store PermissionStore) *Enforcer {
	return &Enforcer{store: store}
}

// Enforce checks if a subject has permission for a resource+action.
func (e *Enforcer) Enforce(ctx context.Context, subject Subject, resource, action string) bool {
	ok, err := e.store.HasPermission(ctx, subject, resource, action)
	if err != nil {
		return false
	}
	return ok
}

// DataScope returns the data scope for a subject on a resource.
func (e *Enforcer) DataScope(ctx context.Context, subject Subject, resource string) DataScope {
	scope, err := e.store.GetDataScope(ctx, subject, resource)
	if err != nil {
		return ScopeSelf // Most restrictive by default
	}
	return scope
}

// HasAnyPermission checks if subject has any of the given permissions.
func (e *Enforcer) HasAnyPermission(ctx context.Context, subject Subject, permissions ...[2]string) bool {
	for _, p := range permissions {
		if e.Enforce(ctx, subject, p[0], p[1]) {
			return true
		}
	}
	return false
}

// HasAllPermissions checks if subject has all of the given permissions.
func (e *Enforcer) HasAllPermissions(ctx context.Context, subject Subject, permissions ...[2]string) bool {
	for _, p := range permissions {
		if !e.Enforce(ctx, subject, p[0], p[1]) {
			return false
		}
	}
	return true
}

// ---------------------------------------------------------------------------
// DataScopeFilter
// ---------------------------------------------------------------------------

// DataScopeFilter computes SQL WHERE conditions and query parameters based on
// the subject's DataScope and WarehouseIDs. Use this to build secure queries.
type DataScopeFilter struct {
	WarehouseIDs []string `json:"warehouse_ids,omitempty"`
	TenantID     string   `json:"tenant_id,omitempty"`
	DeptID       string   `json:"dept_id,omitempty"`
	UserID       string   `json:"user_id,omitempty"`
	Scope        DataScope `json:"scope"`
}

// ApplyDataScope builds a DataScopeFilter from a subject and a resource's scope.
// Callers use the filter to add WHERE clauses and query parameters.
func ApplyDataScope(subject Subject, scope DataScope) DataScopeFilter {
	f := DataScopeFilter{Scope: scope}
	switch scope {
	case ScopeAll:
		// No filtering needed
	case ScopeTenant:
		f.TenantID = subject.TenantID
	case ScopeWarehouse:
		f.TenantID = subject.TenantID
		f.WarehouseIDs = subject.WarehouseIDs
	case ScopeDepartment:
		f.TenantID = subject.TenantID
		f.DeptID = subject.DeptID
	case ScopeSelf:
		f.TenantID = subject.TenantID
		f.UserID = subject.UserID
	}
	return f
}

// NeedsWarehouseFilter returns true when the scope requires warehouse-level filtering.
func (f DataScopeFilter) NeedsWarehouseFilter() bool {
	return f.Scope == ScopeWarehouse && len(f.WarehouseIDs) > 0
}

// NeedsNoFilter returns true when all data is visible (ScopeAll).
func (f DataScopeFilter) NeedsNoFilter() bool {
	return f.Scope == ScopeAll
}

// ---------------------------------------------------------------------------
// In-memory Permission Store
// ---------------------------------------------------------------------------

// InMemPermissionStore is a simple permission store for testing and prototyping.
type InMemPermissionStore struct {
	mu          sync.RWMutex
	roles       map[string]roleDef              // roleID → permissions
	userRoles   map[string][]string             // userID → roleIDs
	scopes      map[string]DataScope            // resource → scope mapping
}

type roleDef struct {
	permissions map[string][]string // resource → list of actions
}

// NewInMemPermissionStore creates an in-memory permission store.
func NewInMemPermissionStore() *InMemPermissionStore {
	return &InMemPermissionStore{
		roles:     make(map[string]roleDef),
		userRoles: make(map[string][]string),
		scopes:    make(map[string]DataScope),
	}
}

// AddRole adds a role with its resource-action permissions.
func (s *InMemPermissionStore) AddRole(roleID string, permissions map[string][]string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.roles[roleID] = roleDef{permissions: permissions}
}

// AssignRole assigns a role to a user.
func (s *InMemPermissionStore) AssignRole(userID, roleID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.userRoles[userID] = append(s.userRoles[userID], roleID)
}

// SetDataScope sets the data scope for a resource.
func (s *InMemPermissionStore) SetDataScope(resource string, scope DataScope) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.scopes[resource] = scope
}

func (s *InMemPermissionStore) HasPermission(ctx context.Context, subject Subject, resource, action string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, roleID := range subject.RoleIDs {
		if rd, ok := s.roles[roleID]; ok {
			if actions, ok := rd.permissions[resource]; ok {
				for _, a := range actions {
					if a == action || a == "*" {
						return true, nil
					}
				}
			}
		}
	}
	return false, nil
}

func (s *InMemPermissionStore) GetDataScope(ctx context.Context, subject Subject, resource string) (DataScope, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if scope, ok := s.scopes[resource]; ok {
		return scope, nil
	}
	return ScopeTenant, nil // default
}

// ---------------------------------------------------------------------------
// Predefined resource actions
// ---------------------------------------------------------------------------

const (
	ActionList    = "list"
	ActionView    = "view"
	ActionCreate  = "create"
	ActionUpdate  = "update"
	ActionDelete  = "delete"
	ActionExport  = "export"
	ActionImport  = "import"
	ActionApprove = "approve"
	ActionCancel  = "cancel"
)

// Predefined resources.
const (
	ResourceOrder     = "order"
	ResourceParcel    = "parcel"
	ResourceClient    = "client"
	ResourceWarehouse = "warehouse"
	ResourceRoute     = "route"
	ResourceFinance   = "finance"
	ResourceUser      = "user"
	ResourceRole      = "role"
	ResourceReport    = "report"
)
