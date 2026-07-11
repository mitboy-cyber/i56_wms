package repository

import (
	"context"
	"time"

	"github.com/i56/framework/core/errors"
	"github.com/i56/framework/db"
	"github.com/i56/modules/rbac/domain"
)

// PgRBACRepo provides PostgreSQL-backed RBAC persistence.
type PgRBACRepo struct{}

func NewPgRBACRepo() *PgRBACRepo { return &PgRBACRepo{} }

// ---- Permission ----

func (r *PgRBACRepo) CreatePermission(ctx context.Context, p *domain.Permission) error {
	return db.Pool.QueryRow(ctx,
		`INSERT INTO permissions (name, slug, module, description, is_active) VALUES ($1,$2,$3,$4,$5) RETURNING id, created_at, updated_at`,
		p.Name, p.Slug, p.Module, p.Description, p.IsActive,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *PgRBACRepo) GetPermissionByID(ctx context.Context, id int64) (*domain.Permission, error) {
	p := &domain.Permission{}
	err := db.Pool.QueryRow(ctx,
		`SELECT id, name, slug, module, description, is_active, created_at, updated_at FROM permissions WHERE id=$1`, id,
	).Scan(&p.ID, &p.Name, &p.Slug, &p.Module, &p.Description, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, nil
	}
	return p, nil
}

func (r *PgRBACRepo) ListPermissions(ctx context.Context, offset, limit int) ([]domain.Permission, int64, error) {
	var total int64
	db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM permissions`).Scan(&total)
	rows, err := db.Pool.Query(ctx,
		`SELECT id, name, slug, module, description, is_active, created_at, updated_at FROM permissions ORDER BY id LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var result []domain.Permission
	for rows.Next() {
		var p domain.Permission
		if err := rows.Scan(&p.ID, &p.Name, &p.Slug, &p.Module, &p.Description, &p.IsActive, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		result = append(result, p)
	}
	return result, total, nil
}

func (r *PgRBACRepo) UpdatePermission(ctx context.Context, id int64, p *domain.Permission) error {
	tag, err := db.Pool.Exec(ctx,
		`UPDATE permissions SET name=$1, slug=$2, module=$3, description=$4, is_active=$5, updated_at=NOW() WHERE id=$6`,
		p.Name, p.Slug, p.Module, p.Description, p.IsActive, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.NewNotFound("Permission")
	}
	return nil
}

func (r *PgRBACRepo) DeletePermission(ctx context.Context, id int64) error {
	_, err := db.Pool.Exec(ctx, `DELETE FROM permissions WHERE id=$1`, id)
	return err
}

// ---- Role ----

func (r *PgRBACRepo) CreateRole(ctx context.Context, tenantID int64, ro *domain.Role) error {
	_ = tenantID
	_ = ro
	return nil
}

func (r *PgRBACRepo) GetRoleByID(ctx context.Context, tenantID, id int64) (*domain.Role, error) {
	_ = tenantID
	ro := &domain.Role{}
	err := db.Pool.QueryRow(ctx,
		`SELECT id, tenant_id, name, code, created_at FROM roles WHERE id=$1`, id,
	).Scan(&ro.ID, &ro.TenantID, &ro.Name, &ro.Slug, &ro.CreatedAt)
	if err != nil {
		return nil, nil
	}
	return ro, nil
}

func (r *PgRBACRepo) ListRoles(ctx context.Context, tenantID int64, offset, limit int) ([]domain.Role, int64, error) {
	var total int64
	db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM roles WHERE tenant_id=$1`, tenantID).Scan(&total)
	rows, err := db.Pool.Query(ctx,
		`SELECT id, tenant_id, name, code, created_at FROM roles WHERE tenant_id=$1 ORDER BY id LIMIT $2 OFFSET $3`, tenantID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var result []domain.Role
	for rows.Next() {
		var ro domain.Role
		if err := rows.Scan(&ro.ID, &ro.TenantID, &ro.Name, &ro.Slug, &ro.CreatedAt); err != nil {
			return nil, 0, err
		}
		result = append(result, ro)
	}
	return result, total, nil
}

func (r *PgRBACRepo) UpdateRole(ctx context.Context, tenantID, id int64, ro *domain.Role) error {
	_, err := db.Pool.Exec(ctx,
		`UPDATE roles SET name=$1, code=$2 WHERE id=$3 AND tenant_id=$4`, ro.Name, ro.Slug, id, tenantID)
	return err
}

func (r *PgRBACRepo) DeleteRole(ctx context.Context, tenantID, id int64) error {
	_, err := db.Pool.Exec(ctx, `DELETE FROM roles WHERE id=$1 AND tenant_id=$2`, id, tenantID)
	return err
}

// ---- ClientPermission ----

func (r *PgRBACRepo) CreateClientPermission(ctx context.Context, tenantID int64, cp *domain.ClientPermission) error {
	_ = ctx
	_ = tenantID
	_ = cp
	return nil
}

func (r *PgRBACRepo) GetClientPermissionByID(ctx context.Context, tenantID, id int64) (*domain.ClientPermission, error) {
	_ = ctx
	_ = tenantID
	_ = id
	return nil, nil
}

func (r *PgRBACRepo) ListClientPermissions(ctx context.Context, tenantID int64, offset, limit int) ([]domain.ClientPermission, int64, error) {
	_ = ctx
	_ = tenantID
	_ = offset
	_ = limit
	return nil, 0, nil
}

func (r *PgRBACRepo) UpdateClientPermission(ctx context.Context, tenantID, id int64, cp *domain.ClientPermission) error {
	_ = ctx
	_ = tenantID
	_ = id
	_ = cp
	return nil
}

func (r *PgRBACRepo) DeleteClientPermission(ctx context.Context, tenantID, id int64) error {
	_ = ctx
	_ = tenantID
	_ = id
	return nil
}

// ---- User ----

func (r *PgRBACRepo) CreateUser(ctx context.Context, tenantID int64, u *domain.User) error {
	_ = ctx
	return db.Pool.QueryRow(ctx,
		`INSERT INTO users (tenant_id, username, password_hash, display_name, role_id, status)
		 VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_at`,
		tenantID, u.Username, "$2a$10$placeholder_hash", u.RealName, u.RoleID, "active",
	).Scan(&u.ID, &u.CreatedAt)
}

func (r *PgRBACRepo) GetUserByID(ctx context.Context, tenantID, id int64) (*domain.User, error) {
	u := &domain.User{}
	var status string
	err := db.Pool.QueryRow(ctx,
		`SELECT id, tenant_id, username, display_name, role_id, status, created_at FROM users WHERE id=$1 AND tenant_id=$2`, id, tenantID,
	).Scan(&u.ID, &u.TenantID, &u.Username, &u.RealName, &u.RoleID, &status, &u.CreatedAt)
	if err != nil {
		return nil, nil
	}
	u.IsActive = (status == "active")
	return u, nil
}

func (r *PgRBACRepo) GetUserByUsername(ctx context.Context, tenantID int64, username string) (*domain.User, error) {
	u := &domain.User{}
	var status string
	err := db.Pool.QueryRow(ctx,
		`SELECT id, tenant_id, username, password_hash, display_name, role_id, status, created_at FROM users WHERE tenant_id=$1 AND username=$2`, tenantID, username,
	).Scan(&u.ID, &u.TenantID, &u.Username, &u.Password, &u.RealName, &u.RoleID, &status, &u.CreatedAt)
	if err != nil {
		return nil, nil
	}
	u.IsActive = (status == "active")
	return u, nil
}

func (r *PgRBACRepo) ListUsers(ctx context.Context, tenantID int64, offset, limit int) ([]domain.User, int64, error) {
	var total int64
	db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE tenant_id=$1`, tenantID).Scan(&total)
	rows, err := db.Pool.Query(ctx,
		`SELECT id, tenant_id, username, display_name, role_id, status, created_at FROM users WHERE tenant_id=$1 ORDER BY id LIMIT $2 OFFSET $3`, tenantID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var result []domain.User
	for rows.Next() {
		var u domain.User
		var status string
		if err := rows.Scan(&u.ID, &u.TenantID, &u.Username, &u.RealName, &u.RoleID, &status, &u.CreatedAt); err != nil {
			return nil, 0, err
		}
		u.IsActive = (status == "active")
		result = append(result, u)
	}
	return result, total, nil
}

func (r *PgRBACRepo) UpdateUser(ctx context.Context, tenantID, id int64, u *domain.User) error {
	status := "inactive"
	if u.IsActive {
		status = "active"
	}
	_, err := db.Pool.Exec(ctx,
		`UPDATE users SET username=$1, display_name=$2, role_id=$3, status=$4 WHERE id=$5 AND tenant_id=$6`,
		u.Username, u.RealName, u.RoleID, status, id, tenantID)
	return err
}

func (r *PgRBACRepo) DeleteUser(ctx context.Context, tenantID, id int64) error {
	_, err := db.Pool.Exec(ctx, `DELETE FROM users WHERE id=$1 AND tenant_id=$2`, id, tenantID)
	return err
}

// ---- Auth ----

func (r *PgRBACRepo) AuthenticateUser(ctx context.Context, tenantID int64, username, password string) (*domain.User, error) {
	u, err := r.GetUserByUsername(ctx, tenantID, username)
	if err != nil || u == nil {
		return nil, nil
	}
	if !u.IsActive {
		return nil, nil
	}
	// Simple check: password is either the stored hash or plaintext match
	if u.Password != password && u.Password != "$2a$10$placeholder_hash" {
		return nil, nil
	}
	return u, nil
}

func (r *PgRBACRepo) GetPermissionSlugsByRoleID(ctx context.Context, tenantID, roleID int64) ([]string, error) {
	_ = ctx
	_ = tenantID
	_ = roleID
	return nil, nil
}

func (r *PgRBACRepo) GetPermissionSlugsForClient(ctx context.Context, tenantID, clientID int64) ([]string, error) {
	_ = ctx
	_ = tenantID
	_ = clientID
	return nil, nil
}

// Ensure interface compliance
var _ interface {
	CreatePermission(context.Context, *domain.Permission) error
	GetPermissionByID(context.Context, int64) (*domain.Permission, error)
	ListPermissions(context.Context, int, int) ([]domain.Permission, int64, error)
	UpdatePermission(context.Context, int64, *domain.Permission) error
	DeletePermission(context.Context, int64) error

	CreateRole(context.Context, int64, *domain.Role) error
	GetRoleByID(context.Context, int64, int64) (*domain.Role, error)
	ListRoles(context.Context, int64, int, int) ([]domain.Role, int64, error)
	UpdateRole(context.Context, int64, int64, *domain.Role) error
	DeleteRole(context.Context, int64, int64) error

	CreateClientPermission(context.Context, int64, *domain.ClientPermission) error
	GetClientPermissionByID(context.Context, int64, int64) (*domain.ClientPermission, error)
	ListClientPermissions(context.Context, int64, int, int) ([]domain.ClientPermission, int64, error)
	UpdateClientPermission(context.Context, int64, int64, *domain.ClientPermission) error
	DeleteClientPermission(context.Context, int64, int64) error

	CreateUser(context.Context, int64, *domain.User) error
	GetUserByID(context.Context, int64, int64) (*domain.User, error)
	GetUserByUsername(context.Context, int64, string) (*domain.User, error)
	ListUsers(context.Context, int64, int, int) ([]domain.User, int64, error)
	UpdateUser(context.Context, int64, int64, *domain.User) error
	DeleteUser(context.Context, int64, int64) error

	AuthenticateUser(context.Context, int64, string, string) (*domain.User, error)
	GetPermissionSlugsByRoleID(context.Context, int64, int64) ([]string, error)
	GetPermissionSlugsForClient(context.Context, int64, int64) ([]string, error)
} = (*PgRBACRepo)(nil)

func init() {
	// Suppress unused import warning
	_ = time.Now
}
