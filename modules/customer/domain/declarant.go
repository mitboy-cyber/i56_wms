package domain

import "time"

// DeclarantType represents whether the declarant is an individual or company.
type DeclarantType string

const (
	DeclarantIndividual DeclarantType = "individual" // 个人
	DeclarantCompany    DeclarantType = "company"    // 公司
)

// DeclarantAuthStatus represents the certification status.
type DeclarantAuthStatus string

const (
	AuthPending   DeclarantAuthStatus = "pending"
	AuthVerifying DeclarantAuthStatus = "verifying"
	AuthVerified  DeclarantAuthStatus = "verified"
	AuthFailed    DeclarantAuthStatus = "failed"
)

// Declarant represents customs declaration identity.
type Declarant struct {
	ID             int64               `json:"id"`
	ClientID       int64               `json:"client_id"`
	MemberID       int64               `json:"member_id"`
	Type           DeclarantType       `json:"type"`
	Name           string              `json:"name"`
	IDNumber       string              `json:"id_number"`
	CompanyTaxID   string              `json:"company_tax_id"`
	Phone          string              `json:"phone"`
	AuthStatus     DeclarantAuthStatus `json:"auth_status"`
	IsActive       bool                `json:"is_active"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
}
