package domain

import "time"

type LedgerEntry struct {
	ID          int64     `json:"id"`
	TenantID    int64     `json:"tenant_id"`
	ClientID    int64     `json:"client_id"`
	Amount      float64   `json:"amount"`
	BalanceAfter float64  `json:"balance_after"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
