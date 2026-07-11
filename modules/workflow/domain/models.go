package domain

import "time"

// ===================================================================
// I56 WMS Workflow Engine — Real domain models
//
// Pipelines:
//   1. 入库流程 (Inbound): 收货确认→称重测量→上架入库→完成  (4 steps)
//   2. 出库流程 (Outbound): 拣货→送打包→打包→核重→送出库→送装柜→装柜→完成 (8 steps)
// Triggers:
//   - order_created     →  generate outbound pipeline per order
//   - parcel_received   →  generate inbound pipeline per parcel
// ===================================================================

// ─── Step type constants ─────────────────────────────────────────────
const (
	StepReceiveConfirm   = "receive_confirm"   // 收货确认
	StepWeighMeasure     = "weigh_measure"     // 称重测量
	StepPutaway          = "putaway"           // 上架入库
	StepPick             = "pick"              // 拣货
	StepSendToPack       = "send_to_pack"      // 送打包
	StepPack             = "pack"              // 打包
	StepWeightCheck      = "weight_check"      // 核重
	StepSendOut          = "send_out"          // 送出库
	StepSendToContainer  = "send_to_container" // 送装柜
	StepLoadContainer    = "load_container"    // 装柜
	StepComplete         = "complete"          // 完成 (terminal step)

	// Triggers
	TriggerOrderCreated   = "order_created"
	TriggerParcelReceived = "parcel_received"

	// WorkOrder statuses
	WOStatusPending    = "pending"
	WOStatusInProgress = "in_progress"
	WOStatusCompleted  = "completed"
	WOStatusCancelled  = "cancelled"
)

// StepChineseNames maps step codes to Chinese display names
var StepChineseNames = map[string]string{
	StepReceiveConfirm:  "收货确认",
	StepWeighMeasure:    "称重测量",
	StepPutaway:         "上架入库",
	StepPick:            "拣货",
	StepSendToPack:      "送打包",
	StepPack:            "打包",
	StepWeightCheck:     "核重",
	StepSendOut:         "送出库",
	StepSendToContainer: "送装柜",
	StepLoadContainer:   "装柜",
	StepComplete:        "完成",
}

// StepNameCN returns the Chinese display name for a step code
func StepNameCN(code string) string {
	if name, ok := StepChineseNames[code]; ok {
		return name
	}
	return code
}

// ─── WorkflowProcess ─────────────────────────────────────────────────
// A workflow pipeline definition: ordered steps with metadata.
type WorkflowProcess struct {
	ID           int64           `json:"id"`
	TenantID     int64           `json:"tenant_id"`
	Name         string          `json:"name"`          // e.g. "标准入库流程"
	Code         string          `json:"code"`          // e.g. "inbound"
	TriggerEvent string          `json:"trigger_event"` // order_created / parcel_received
	Steps        []WorkflowStep  `json:"steps"`
	IsActive     bool            `json:"is_active"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// ─── WorkflowStep ────────────────────────────────────────────────────
// A single step within a workflow pipeline.
type WorkflowStep struct {
	ID             int64  `json:"id"`
	ProcessID      int64  `json:"process_id"`
	Name           string `json:"name"`            // e.g. "receive_confirm"
	DisplayName    string `json:"display_name"`    // e.g. "收货确认"
	OrderIndex     int    `json:"order_index"`     // 1-based position in pipeline
	Required       bool   `json:"required"`        // must be completed to advance
	Assignable     bool   `json:"assignable"`      // can be assigned to a user
	TimeoutMinutes int    `json:"timeout_minutes"` // SLA timeout (0 = none)
}

// StepsDisplay returns a string like "收货确认→称重测量→上架入库→完成"
func (p *WorkflowProcess) StepsDisplay() string {
	s := ""
	for i, step := range p.Steps {
		if i > 0 {
			s += "→"
		}
		s += StepNameCN(step.Name)
	}
	return s
}

// ─── WorkOrder ───────────────────────────────────────────────────────
// A running instance of a workflow process, tied to a parcel or order.
type WorkOrder struct {
	ID          int64      `json:"id"`
	TenantID    int64      `json:"tenant_id"`
	WarehouseID int64      `json:"warehouse_id"`
	ProcessID   int64      `json:"process_id"`   // FK → WorkflowProcess
	ProcessName string     `json:"process_name"` // denormalized for display
	ParcelID    *int64     `json:"parcel_id"`    // for inbound workflows
	OrderID     *int64     `json:"order_id"`     // for outbound workflows
	Status      string     `json:"status"`       // pending / in_progress / completed / cancelled
	CurrentStep int        `json:"current_step"` // 1-based index into process steps
	AssignedTo  *int64     `json:"assigned_to"`  // operator user ID
	AssignedName string    `json:"assigned_name"` // denormalized
	Title       string     `json:"title"`        // human-readable title
	Description string     `json:"description"`
	Priority    int        `json:"priority"`     // 0=normal, 1=urgent, 2=critical
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// StatusDisplay returns the Chinese status name
func StatusDisplay(status string) string {
	switch status {
	case WOStatusPending:
		return "待处理"
	case WOStatusInProgress:
		return "进行中"
	case WOStatusCompleted:
		return "已完成"
	case WOStatusCancelled:
		return "已取消"
	default:
		return status
	}
}

// ParcelOrOrderRef returns a human-readable parcel/order reference
func (wo *WorkOrder) ParcelOrOrderRef() string {
	if wo.ParcelID != nil {
		return "包裹"
	}
	if wo.OrderID != nil {
		return "订单"
	}
	return "—"
}
