package domain
import "time"

type ProcessInstance struct {
	ID          int64              `json:"id"`
	TenantID    int64              `json:"tenant_id"`
	ProcessName string             `json:"process_name"` // e.g. "采购审批"
	BusinessKey string             `json:"business_key"` // e.g. "order:123"
	CurrentStep string             `json:"current_step"`
	Status      string             `json:"status"`       // running | approved | rejected | cancelled
	Steps       []ProcessStep      `json:"steps"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

type ProcessStep struct {
	Name      string     `json:"name"`
	Assignee  string     `json:"assignee"`
	Status    string     `json:"status"`  // pending | approved | rejected
	Comment   string     `json:"comment,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
}
