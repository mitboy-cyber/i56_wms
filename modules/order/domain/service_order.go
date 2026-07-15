package domain

import "time"

// ServiceOrderType represents the type of value-added service.
type ServiceOrderType string

const (
	ServiceTypePhotos      ServiceOrderType = "photos"
	ServiceTypeInspection  ServiceOrderType = "inspection"
	ServiceTypeRepack      ServiceOrderType = "repack"
	ServiceTypeRemoveBox   ServiceOrderType = "remove_box"
	ServiceTypeReinforce   ServiceOrderType = "reinforce"
	ServiceTypeInsurance   ServiceOrderType = "insurance"
	ServiceTypeCustomDecl  ServiceOrderType = "customs_declaration"
)

// ServiceOrderStatus represents the status of a value-added service order.
type ServiceOrderStatus string

const (
	ServStatusPending    ServiceOrderStatus = "pending"
	ServStatusProcessing ServiceOrderStatus = "processing"
	ServStatusCompleted  ServiceOrderStatus = "completed"
	ServStatusCancelled  ServiceOrderStatus = "cancelled"
)

// ServiceOrder represents a value-added service request (拍照, 验货, 加固, etc.).
type ServiceOrder struct {
	ID            int64              `json:"id"`
	TenantID      int64              `json:"tenant_id"`
	ClientID      int64              `json:"client_id"`
	ParcelID      int64              `json:"parcel_id"`
	OrderID       int64              `json:"order_id"`
	ServiceType   ServiceOrderType   `json:"service_type"`
	ServiceName   string             `json:"service_name"`
	Status        ServiceOrderStatus `json:"status"`
	Price         float64            `json:"price"`
	OperatorID    int64              `json:"operator_id"`
	ResultNote    string             `json:"result_note"`
	ResultImages  []string           `json:"result_images"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
	CompletedAt   *time.Time         `json:"completed_at"`
}

func ServiceValidTransitions() map[ServiceOrderStatus][]ServiceOrderStatus {
	return map[ServiceOrderStatus][]ServiceOrderStatus{
		ServStatusPending:    {ServStatusProcessing, ServStatusCancelled},
		ServStatusProcessing: {ServStatusCompleted},
		ServStatusCompleted:  {},
		ServStatusCancelled:  {},
	}
}

func (o *ServiceOrder) CanTransitionTo(target ServiceOrderStatus) bool {
	allowed, ok := ServiceValidTransitions()[o.Status]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == target {
			return true
		}
	}
	return false
}
