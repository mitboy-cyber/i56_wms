package domain

import "time"

// WorkOrderStatus represents the lifecycle of a work order.
type WorkOrderStatus string

const (
	WOStatusDraft      WorkOrderStatus = "draft"
	WOStatusAssigned   WorkOrderStatus = "assigned"
	WOStatusInProgress WorkOrderStatus = "in_progress"
	WOStatusCompleted  WorkOrderStatus = "completed"
	WOStatusCancelled  WorkOrderStatus = "cancelled"
)

// WorkOrderPriority represents the priority level.
type WorkOrderPriority string

const (
	WOPriorityLow    WorkOrderPriority = "low"
	WOPriorityNormal WorkOrderPriority = "normal"
	WOPriorityHigh   WorkOrderPriority = "high"
	WOPriorityUrgent WorkOrderPriority = "urgent"
)

// WorkOrder represents a task assigned to a warehouse operator via PDA.
type WorkOrder struct {
	ID               int64             `json:"id"`
	TenantID         int64             `json:"tenant_id"`
	WarehouseID      int64             `json:"warehouse_id"`
	WorkOrderNo      string            `json:"work_order_no"`
	TemplateID       int64             `json:"template_id"`
	Type             string            `json:"type"`
	Status           WorkOrderStatus   `json:"status"`
	Priority         WorkOrderPriority `json:"priority"`
	AssignedTo       int64             `json:"assigned_to"`
	OrderID          int64             `json:"order_id"`
	ParcelIDs        []int64           `json:"parcel_ids"`
	LocationCode     string            `json:"location_code"`
	TargetLocation   string            `json:"target_location"`
	Instructions     string            `json:"instructions"`
	ResultNote       string            `json:"result_note"`
	StartedAt        *time.Time        `json:"started_at"`
	CompletedAt      *time.Time        `json:"completed_at"`
	CreatedBy        int64             `json:"created_by"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
}

func WorkOrderValidTransitions() map[WorkOrderStatus][]WorkOrderStatus {
	return map[WorkOrderStatus][]WorkOrderStatus{
		WOStatusDraft:      {WOStatusAssigned, WOStatusCancelled},
		WOStatusAssigned:   {WOStatusInProgress, WOStatusCancelled},
		WOStatusInProgress: {WOStatusCompleted, WOStatusCancelled},
		WOStatusCompleted:  {},
		WOStatusCancelled:  {},
	}
}
