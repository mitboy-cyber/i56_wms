package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/framework/core/errors"
	"github.com/i56/modules/rbac/domain"
)

// MemRBACRepo is an in-memory implementation of the RBAC repository.
// It manages permissions, roles, client permissions, and users with tenant-aware filtering.
type MemRBACRepo struct {
	mu                sync.RWMutex
	permissions       map[int64]*domain.Permission
	roles             map[int64]*domain.Role
	clientPermissions map[int64]*domain.ClientPermission
	users             map[int64]*domain.User
	nextPermID        int64
	nextRoleID        int64
	nextClientPermID  int64
	nextUserID        int64
}

// NewMemRBACRepo creates a new MemRBACRepo pre-seeded with default data.
func NewMemRBACRepo() *MemRBACRepo {
	r := &MemRBACRepo{
		permissions:       make(map[int64]*domain.Permission),
		roles:             make(map[int64]*domain.Role),
		clientPermissions: make(map[int64]*domain.ClientPermission),
		users:             make(map[int64]*domain.User),
		nextPermID:        1,
		nextRoleID:        1,
		nextClientPermID:  1,
		nextUserID:        1,
	}
	r.seed()
	return r
}

// seed populates the repository with default data.
func (r *MemRBACRepo) seed() {
	now := time.Now()

	// Seed permissions
	for _, p := range domain.DefaultPermissions() {
		cp := p
		cp.CreatedAt = now
		cp.UpdatedAt = now
		r.permissions[cp.ID] = &cp
		if cp.ID >= r.nextPermID {
			r.nextPermID = cp.ID + 1
		}
	}

	// Seed roles
	for _, role := range domain.DefaultRoles() {
		cr := role
		cr.CreatedAt = now
		cr.UpdatedAt = now
		r.roles[cr.ID] = &cr
		if cr.ID >= r.nextRoleID {
			r.nextRoleID = cr.ID + 1
		}
	}

	// Seed users
	for _, u := range domain.DefaultUsers() {
		cu := u
		cu.CreatedAt = now
		cu.UpdatedAt = now
		r.users[cu.ID] = &cu
		if cu.ID >= r.nextUserID {
			r.nextUserID = cu.ID + 1
		}
	}

	// Seed client permissions
	for _, cp := range domain.DefaultClientPermissions() {
		ccp := cp
		ccp.CreatedAt = now
		ccp.UpdatedAt = now
		r.clientPermissions[ccp.ID] = &ccp
		if ccp.ID >= r.nextClientPermID {
			r.nextClientPermID = ccp.ID + 1
		}
	}
}

// =============================================================================
// Permission CRUD
// =============================================================================

// CreatePermission adds a new permission.
func (r *MemRBACRepo) CreatePermission(ctx context.Context, p *domain.Permission) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p.ID = atomic.AddInt64(&r.nextPermID, 1) - 1
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	r.permissions[p.ID] = p
	return nil
}

// GetPermissionByID returns a permission by ID.
func (r *MemRBACRepo) GetPermissionByID(ctx context.Context, id int64) (*domain.Permission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.permissions[id]
	if !ok {
		return nil, nil
	}
	return p, nil
}

// ListPermissions returns all permissions with offset/limit pagination.
func (r *MemRBACRepo) ListPermissions(ctx context.Context, offset, limit int) ([]domain.Permission, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.Permission
	for _, p := range r.permissions {
		result = append(result, *p)
	}
	total := int64(len(result))
	if offset >= int(total) {
		return nil, total, nil
	}
	end := offset + limit
	if end > int(total) {
		end = int(total)
	}
	return result[offset:end], total, nil
}

// UpdatePermission updates an existing permission.
func (r *MemRBACRepo) UpdatePermission(ctx context.Context, id int64, p *domain.Permission) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.permissions[id]; !ok {
		return errors.NewNotFound("Permission")
	}
	p.ID = id
	p.UpdatedAt = time.Now()
	r.permissions[id] = p
	return nil
}

// DeletePermission removes a permission by ID.
func (r *MemRBACRepo) DeletePermission(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.permissions[id]; !ok {
		return errors.NewNotFound("Permission")
	}
	delete(r.permissions, id)
	return nil
}

// =============================================================================
// Role CRUD
// =============================================================================

// CreateRole adds a new role.
func (r *MemRBACRepo) CreateRole(ctx context.Context, tenantID int64, ro *domain.Role) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	ro.ID = atomic.AddInt64(&r.nextRoleID, 1) - 1
	ro.TenantID = tenantID
	ro.CreatedAt = time.Now()
	ro.UpdatedAt = time.Now()
	r.roles[ro.ID] = ro
	return nil
}

// GetRoleByID returns a role by ID with tenant filter.
func (r *MemRBACRepo) GetRoleByID(ctx context.Context, tenantID, id int64) (*domain.Role, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ro, ok := r.roles[id]
	if !ok || ro.TenantID != tenantID {
		return nil, nil
	}
	return ro, nil
}

// ListRoles returns all roles for a tenant with offset/limit pagination.
func (r *MemRBACRepo) ListRoles(ctx context.Context, tenantID int64, offset, limit int) ([]domain.Role, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.Role
	for _, ro := range r.roles {
		if ro.TenantID == tenantID {
			result = append(result, *ro)
		}
	}
	total := int64(len(result))
	if offset >= int(total) {
		return nil, total, nil
	}
	end := offset + limit
	if end > int(total) {
		end = int(total)
	}
	return result[offset:end], total, nil
}

// UpdateRole updates an existing role.
func (r *MemRBACRepo) UpdateRole(ctx context.Context, tenantID, id int64, ro *domain.Role) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.roles[id]
	if !ok || existing.TenantID != tenantID {
		return errors.NewNotFound("Role")
	}
	ro.ID = id
	ro.TenantID = tenantID
	ro.UpdatedAt = time.Now()
	r.roles[id] = ro
	return nil
}

// DeleteRole removes a role by ID with tenant check.
func (r *MemRBACRepo) DeleteRole(ctx context.Context, tenantID, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.roles[id]
	if !ok || existing.TenantID != tenantID {
		return errors.NewNotFound("Role")
	}
	delete(r.roles, id)
	return nil
}

// =============================================================================
// ClientPermission CRUD
// =============================================================================

// CreateClientPermission adds a new client permission.
func (r *MemRBACRepo) CreateClientPermission(ctx context.Context, tenantID int64, cp *domain.ClientPermission) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp.ID = atomic.AddInt64(&r.nextClientPermID, 1) - 1
	cp.TenantID = tenantID
	cp.CreatedAt = time.Now()
	cp.UpdatedAt = time.Now()
	r.clientPermissions[cp.ID] = cp
	return nil
}

// GetClientPermissionByID returns a client permission by ID with tenant filter.
func (r *MemRBACRepo) GetClientPermissionByID(ctx context.Context, tenantID, id int64) (*domain.ClientPermission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cp, ok := r.clientPermissions[id]
	if !ok || cp.TenantID != tenantID {
		return nil, nil
	}
	return cp, nil
}

// ListClientPermissions returns all client permissions for a tenant.
func (r *MemRBACRepo) ListClientPermissions(ctx context.Context, tenantID int64, offset, limit int) ([]domain.ClientPermission, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.ClientPermission
	for _, cp := range r.clientPermissions {
		if cp.TenantID == tenantID {
			result = append(result, *cp)
		}
	}
	total := int64(len(result))
	if offset >= int(total) {
		return nil, total, nil
	}
	end := offset + limit
	if end > int(total) {
		end = int(total)
	}
	return result[offset:end], total, nil
}

// UpdateClientPermission updates an existing client permission.
func (r *MemRBACRepo) UpdateClientPermission(ctx context.Context, tenantID, id int64, cp *domain.ClientPermission) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.clientPermissions[id]
	if !ok || existing.TenantID != tenantID {
		return errors.NewNotFound("ClientPermission")
	}
	cp.ID = id
	cp.TenantID = tenantID
	cp.UpdatedAt = time.Now()
	r.clientPermissions[id] = cp
	return nil
}

// DeleteClientPermission removes a client permission by ID.
func (r *MemRBACRepo) DeleteClientPermission(ctx context.Context, tenantID, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.clientPermissions[id]
	if !ok || existing.TenantID != tenantID {
		return errors.NewNotFound("ClientPermission")
	}
	delete(r.clientPermissions, id)
	return nil
}

// =============================================================================
// User CRUD
// =============================================================================

// CreateUser adds a new user.
func (r *MemRBACRepo) CreateUser(ctx context.Context, tenantID int64, u *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	u.ID = atomic.AddInt64(&r.nextUserID, 1) - 1
	u.TenantID = tenantID
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	r.users[u.ID] = u
	return nil
}

// GetUserByID returns a user by ID with tenant filter.
func (r *MemRBACRepo) GetUserByID(ctx context.Context, tenantID, id int64) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[id]
	if !ok || u.TenantID != tenantID {
		return nil, nil
	}
	return u, nil
}

// GetUserByUsername returns a user by username with tenant filter.
func (r *MemRBACRepo) GetUserByUsername(ctx context.Context, tenantID int64, username string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.users {
		if u.TenantID == tenantID && u.Username == username {
			return u, nil
		}
	}
	return nil, nil
}

// ListUsers returns all users for a tenant with offset/limit pagination.
func (r *MemRBACRepo) ListUsers(ctx context.Context, tenantID int64, offset, limit int) ([]domain.User, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.User
	for _, u := range r.users {
		if u.TenantID == tenantID {
			result = append(result, *u)
		}
	}
	total := int64(len(result))
	if offset >= int(total) {
		return nil, total, nil
	}
	end := offset + limit
	if end > int(total) {
		end = int(total)
	}
	return result[offset:end], total, nil
}

// UpdateUser updates an existing user.
func (r *MemRBACRepo) UpdateUser(ctx context.Context, tenantID, id int64, u *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.users[id]
	if !ok || existing.TenantID != tenantID {
		return errors.NewNotFound("User")
	}
	u.ID = id
	u.TenantID = tenantID
	u.UpdatedAt = time.Now()
	r.users[id] = u
	return nil
}

// DeleteUser removes a user by ID with tenant check.
func (r *MemRBACRepo) DeleteUser(ctx context.Context, tenantID, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.users[id]
	if !ok || existing.TenantID != tenantID {
		return errors.NewNotFound("User")
	}
	delete(r.users, id)
	return nil
}

// =============================================================================
// Convenience: authenticate user
// =============================================================================

// AuthenticateUser checks username/password and returns the user if valid.
func (r *MemRBACRepo) AuthenticateUser(ctx context.Context, tenantID int64, username, password string) (*domain.User, error) {
	u, err := r.GetUserByUsername(ctx, tenantID, username)
	if err != nil || u == nil {
		return nil, nil
	}
	if !u.IsActive || u.Password != password {
		return nil, nil
	}
	return u, nil
}

// =============================================================================
// Convenience: get role permissions
// =============================================================================

// GetPermissionSlugsByRoleID returns all permission slugs for a given role.
func (r *MemRBACRepo) GetPermissionSlugsByRoleID(ctx context.Context, tenantID, roleID int64) ([]string, error) {
	role, err := r.GetRoleByID(ctx, tenantID, roleID)
	if err != nil || role == nil {
		return nil, nil
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	slugs := make([]string, 0, len(role.PermissionIDs))
	for _, pid := range role.PermissionIDs {
		if p, ok := r.permissions[pid]; ok && p.IsActive {
			slugs = append(slugs, p.Slug)
		}
	}
	return slugs, nil
}

// GetPermissionSlugsForClient returns all permission slugs for a given client.
func (r *MemRBACRepo) GetPermissionSlugsForClient(ctx context.Context, tenantID, clientID int64) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, cp := range r.clientPermissions {
		if cp.TenantID == tenantID && cp.ClientID == clientID && cp.IsActive {
			return cp.PermissionSlugs, nil
		}
	}
	return nil, nil
}
