package domain

import "time"

// ClientUserRole represents the role of a user within a client organization.
type ClientUserRole string

const (
	ClientUserRoleAdmin    ClientUserRole = "admin"
	ClientUserRoleOperator ClientUserRole = "operator"
	ClientUserRoleViewer   ClientUserRole = "viewer"
)

// ClientUser represents a login account for a client (customer) organization.
// A client may have multiple users who can log into the client portal.
type ClientUser struct {
	ID           int64          `json:"id"`
	ClientID     int64          `json:"client_id"`
	Username     string         `json:"username"`
	PasswordHash string         `json:"password_hash"`
	DisplayName  string         `json:"display_name"`
	Email        string         `json:"email"`
	Phone        string         `json:"phone"`
	Role         ClientUserRole `json:"role"`
	IsActive     bool           `json:"is_active"`
	LastLoginAt  *time.Time     `json:"last_login_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}
