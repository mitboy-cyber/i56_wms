package domain
import "time"

type ServiceTemplate struct {
	ID          int64            `json:"id"`
	TenantID    int64            `json:"tenant_id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Services    []TemplateService `json:"services"` // combo of service types
	TotalPrice  float64          `json:"total_price"`
	IsActive    bool             `json:"is_active"`
	CreatedAt   time.Time        `json:"created_at"`
}

type TemplateService struct {
	ServiceCode string `json:"service_code"`
	Quantity    int    `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
}
