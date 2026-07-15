package domain

import "time"

// WorkflowProcessStatus represents the status of a workflow process step.
type WorkflowProcessStatus string

const (
	WFPStatusPending   WorkflowProcessStatus = "pending"
	WFPStatusRunning   WorkflowProcessStatus = "running"
	WFPStatusCompleted WorkflowProcessStatus = "completed"
	WFPStatusFailed    WorkflowProcessStatus = "failed"
	WFPStatusSkipped   WorkflowProcessStatus = "skipped"
)

// WorkflowProcess represents a step within a work order's execution lifecycle.
type WorkflowProcess struct {
	ID           int64                 `json:"id"`
	WorkOrderID  int64                 `json:"work_order_id"`
	StepName     string                `json:"step_name"`
	StepOrder    int                   `json:"step_order"`
	Status       WorkflowProcessStatus `json:"status"`
	OperatorID   int64                 `json:"operator_id"`
	StartedAt    *time.Time            `json:"started_at"`
	CompletedAt  *time.Time            `json:"completed_at"`
	ResultNote   string                `json:"result_note"`
	CreatedAt    time.Time             `json:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at"`
}

func WorkflowProcessValidTransitions() map[WorkflowProcessStatus][]WorkflowProcessStatus {
	return map[WorkflowProcessStatus][]WorkflowProcessStatus{
		WFPStatusPending:   {WFPStatusRunning, WFPStatusSkipped},
		WFPStatusRunning:   {WFPStatusCompleted, WFPStatusFailed},
		WFPStatusCompleted: {},
		WFPStatusFailed:    {WFPStatusRunning},
		WFPStatusSkipped:   {},
	}
}
