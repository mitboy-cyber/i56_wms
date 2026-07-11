package domain

import "time"

// MemberAddress is a Taiwan delivery address for a client member.
type MemberAddress struct {
	ID           int64     `json:"id"`
	MemberID     int64     `json:"member_id"`
	RecipientName string   `json:"recipient_name"`
	Phone        string    `json:"phone"`
	PostalCode   string    `json:"postal_code"`
	City         string    `json:"city"`
	District     string    `json:"district"`
	Address      string    `json:"address"`
	IsDefault    bool      `json:"is_default"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
