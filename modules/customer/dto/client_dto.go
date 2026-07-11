// Package dto defines data transfer objects for the customer module.
package dto

// CreateClientRequest is the input for creating a new client.
type CreateClientRequest struct {
	Name         string `json:"name"`
	ClientType   string `json:"client_type"`
	ContactName  string `json:"contact_name"`
	ContactPhone string `json:"contact_phone"`
	ContactEmail string `json:"contact_email"`
	Remark       string `json:"remark"`
}

// UpdateClientRequest is the input for updating a client.
type UpdateClientRequest struct {
	Name         *string `json:"name,omitempty"`
	ContactName  *string `json:"contact_name,omitempty"`
	ContactPhone *string `json:"contact_phone,omitempty"`
	ContactEmail *string `json:"contact_email,omitempty"`
	Remark       *string `json:"remark,omitempty"`
	IsActive     *bool   `json:"is_active,omitempty"`
}

// ClientResponse is the standard client output.
type ClientResponse struct {
	ID           int64   `json:"id"`
	TenantID     int64   `json:"tenant_id"`
	Name         string  `json:"name"`
	Code         string  `json:"code"`
	ClientType   string  `json:"client_type"`
	ContactName  string  `json:"contact_name"`
	ContactPhone string  `json:"contact_phone"`
	ContactEmail string  `json:"contact_email"`
	Balance      float64 `json:"balance"`
	IsActive     bool    `json:"is_active"`
	Remark       string  `json:"remark"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

// CreateDeclarantRequest is the input for creating a declarant.
type CreateDeclarantRequest struct {
	MemberID int64  `json:"member_id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	IDNumber string `json:"id_number"`
	Phone    string `json:"phone"`
}

// UpdateDeclarantRequest is the input for updating a declarant.
type UpdateDeclarantRequest struct {
	Name     *string `json:"name,omitempty"`
	IDNumber *string `json:"id_number,omitempty"`
	Phone    *string `json:"phone,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}
