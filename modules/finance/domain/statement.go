package domain

import "time"

// StatementPeriod represents a billing cycle.
type StatementPeriod struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// Statement represents a periodic billing statement for a client.
type Statement struct {
	ID               int64           `json:"id"`
	TenantID         int64           `json:"tenant_id"`
	ClientID         int64           `json:"client_id"`
	StatementNo      string          `json:"statement_no"`
	Period           StatementPeriod `json:"period"`
	OpeningBalance   float64         `json:"opening_balance"`
	TotalRecharged   float64         `json:"total_recharged"`
	TotalShippingFee float64         `json:"total_shipping_fee"`
	TotalServiceFee  float64         `json:"total_service_fee"`
	TotalAdjustments float64         `json:"total_adjustments"`
	ClosingBalance   float64         `json:"closing_balance"`
	Status           string          `json:"status"`
	GeneratedAt      time.Time       `json:"generated_at"`
	CreatedAt        time.Time       `json:"created_at"`
}
