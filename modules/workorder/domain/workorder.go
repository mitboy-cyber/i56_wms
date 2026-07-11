package domain
import "time"

type WorkOrder struct {
	ID           int64     `json:"id"`
	TenantID     int64     `json:"tenant_id"`
	WarehouseID  int64     `json:"warehouse_id"`
	TemplateID   int64     `json:"template_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Status       string    `json:"status"`
	Priority     int       `json:"priority"`
	AssignedTo   *int64    `json:"assigned_to"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
}

type WorkOrderTemplate struct {
	ID          int64  `json:"id"`
	TenantID    int64  `json:"tenant_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Steps       []string `json:"steps"`
}

type WorkflowProcess struct {
	ID          int64     `json:"id"`
	TenantID    int64     `json:"tenant_id"`
	Name        string    `json:"name"`
	States      []string  `json:"states"`
	Transitions map[string][]string `json:"transitions"`
}
