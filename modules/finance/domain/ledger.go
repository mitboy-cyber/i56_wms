package domain

import "time"

// LedgerType represents the type of ledger entry.
type LedgerType string

const (
	LedgerRecharge  LedgerType = "recharge"
	LedgerDeduction LedgerType = "deduction"
	LedgerRefund    LedgerType = "refund"
	LedgerAdjust    LedgerType = "adjustment"
)

// LedgerEntry represents a single entry in a client's financial ledger.
type LedgerEntry struct {
	ID           int64      `json:"id"`
	TenantID     int64      `json:"tenant_id"`
	ClientID     int64      `json:"client_id"`
	Amount       float64    `json:"amount"`
	BalanceAfter float64    `json:"balance_after"`
	Type         LedgerType `json:"type"`
	ReferenceNo  string     `json:"reference_no"`
	Description  string     `json:"description"`
	CreatedAt    time.Time  `json:"created_at"`
}

// ClientLedger represents the current balance and aggregate information for a client.
type ClientLedger struct {
	ID              int64     `json:"id"`
	TenantID        int64     `json:"tenant_id"`
	ClientID        int64     `json:"client_id"`
	Balance         float64   `json:"balance"`
	TotalRecharged  float64   `json:"total_recharged"`
	TotalSpent      float64   `json:"total_spent"`
	CreditLimit     float64   `json:"credit_limit"`
	FrozenAmount    float64   `json:"frozen_amount"`
	IsActive        bool      `json:"is_active"`
	LastRechargeAt  *time.Time `json:"last_recharge_at"`
	LastDeductionAt *time.Time `json:"last_deduction_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// RechargeMethod represents how the client recharged their account.
type RechargeMethod string

const (
	RechargeMethodBankTransfer RechargeMethod = "bank_transfer"
	RechargeMethodWechat        RechargeMethod = "wechat"
	RechargeMethodAlipay       RechargeMethod = "alipay"
	RechargeMethodOffline       RechargeMethod = "offline"
)

// RechargeStatus represents the status of a recharge request.
type RechargeStatus string

const (
	RechargeStatusPending   RechargeStatus = "pending"
	RechargeStatusConfirmed RechargeStatus = "confirmed"
	RechargeStatusRejected  RechargeStatus = "rejected"
)

// Recharge represents a client account recharge request.
type Recharge struct {
	ID             int64          `json:"id"`
	TenantID       int64          `json:"tenant_id"`
	ClientID       int64          `json:"client_id"`
	RechargeNo     string         `json:"recharge_no"`
	Amount         float64        `json:"amount"`
	Method         RechargeMethod `json:"method"`
	Status         RechargeStatus `json:"status"`
	ProofImageURL  string         `json:"proof_image_url"`
	ReviewedBy     int64          `json:"reviewed_by"`
	ReviewNote     string         `json:"review_note"`
	ReviewedAt     *time.Time     `json:"reviewed_at"`
	Remark         string         `json:"remark"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

// RechargeLog represents an audit log entry for a recharge.
type RechargeLog struct {
	ID         int64     `json:"id"`
	RechargeID int64     `json:"recharge_id"`
	Action     string    `json:"action"`
	OperatorID int64     `json:"operator_id"`
	Note       string    `json:"note"`
	CreatedAt  time.Time `json:"created_at"`
}
