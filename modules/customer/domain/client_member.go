package domain

import "time"

// ClientMember represents an end consumer (recipient in Taiwan).
type ClientMember struct {
	ID          int64     `json:"id"`
	ClientID    int64     `json:"client_id"`
	MemberCode  string    `json:"member_code"`
	Name        string    `json:"name"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
