package domain

import "time"

// ─── Task type constants ─────────────────────────────────────────────
const (
	TaskTypeReceive     = "receive"
	TaskTypeWeigh       = "weigh"
	TaskTypePutaway     = "putaway"
	TaskTypePick        = "pick"
	TaskTypePack        = "pack"
	TaskTypeWeightCheck = "weight_check"
	TaskTypeLoad        = "load"
	TaskTypeException   = "exception"
)

// ─── Task status constants ───────────────────────────────────────────
const (
	StatusPending    = "pending"     // 抢单池中
	StatusClaimed    = "claimed"     // 已认领
	StatusInProgress = "in_progress" // 执行中
	StatusCompleted  = "completed"   // 已完成
	StatusCancelled  = "cancelled"   // 已取消
	StatusTimeout    = "timeout"     // 超时
)

// DefaultTimeoutMinutes is the default SLA timeout for tasks.
const DefaultTimeoutMinutes = 30

// ─── Capacity constants ──────────────────────────────────────────────
const (
	CapForklift   = "forklift"
	CapHeavyLift  = "heavy_lift"
	CapHazmat     = "hazmat"
	CapColdChain  = "cold_chain"
	CapFragile    = "fragile"
	CapEcommerce  = "ecommerce"
)

// ─── WarehouseTask ───────────────────────────────────────────────────
// A warehouse task in the dispatch pool, claimed by operators via PDA.
type WarehouseTask struct {
	ID                   int64      `json:"id"`
	TaskCode             string     `json:"task_code"`             // e.g. "TASK-001"
	TaskType             string     `json:"task_type"`             // receive/weigh/putaway/pick/pack/weight_check/load/exception
	ParcelID             *int64     `json:"parcel_id,omitempty"`
	ParcelTrackingNumber string     `json:"parcel_tracking_number,omitempty"`
	OrderID              *int64     `json:"order_id,omitempty"`
	WarehouseID          int64      `json:"warehouse_id"`
	LocationCode         string     `json:"location_code,omitempty"`
	Status               string     `json:"status"` // pending/claimed/in_progress/completed/cancelled/timeout
	RequiredCapabilities []string   `json:"required_capabilities,omitempty"` // e.g. ["forklift","heavy_lift"]
	AssignedOperatorID   *int64     `json:"assigned_operator_id,omitempty"`
	AssignedAt           *time.Time `json:"assigned_at,omitempty"`
	TimeoutMinutes       int        `json:"timeout_minutes"` // default 30
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

// IsPending returns true if the task is in the claim pool.
func (t *WarehouseTask) IsPending() bool { return t.Status == StatusPending }

// IsClaimed returns true if the task has been claimed.
func (t *WarehouseTask) IsClaimed() bool { return t.Status == StatusClaimed }

// IsTimedOut checks whether a claimed/in-progress task has exceeded its SLA.
func (t *WarehouseTask) IsTimedOut() bool {
	if t.AssignedAt == nil {
		return false
	}
	if t.Status != StatusClaimed && t.Status != StatusInProgress {
		return false
	}
	timeout := time.Duration(t.TimeoutMinutes) * time.Minute
	if timeout <= 0 {
		timeout = DefaultTimeoutMinutes * time.Minute
	}
	return time.Since(*t.AssignedAt) > timeout
}

// StatusDisplay returns the Chinese status name.
func StatusDisplay(status string) string {
	switch status {
	case StatusPending:
		return "抢单池中"
	case StatusClaimed:
		return "已认领"
	case StatusInProgress:
		return "执行中"
	case StatusCompleted:
		return "已完成"
	case StatusCancelled:
		return "已取消"
	case StatusTimeout:
		return "超时"
	default:
		return status
	}
}

// TaskTypeDisplay returns the Chinese task type name.
func TaskTypeDisplay(taskType string) string {
	switch taskType {
	case TaskTypeReceive:
		return "收货"
	case TaskTypeWeigh:
		return "称重"
	case TaskTypePutaway:
		return "上架"
	case TaskTypePick:
		return "拣货"
	case TaskTypePack:
		return "打包"
	case TaskTypeWeightCheck:
		return "核重"
	case TaskTypeLoad:
		return "装柜"
	case TaskTypeException:
		return "异常处理"
	default:
		return taskType
	}
}

// ─── OperatorCapability ──────────────────────────────────────────────
// Operator skills and online status for capability-based matching.
type OperatorCapability struct {
	OperatorID       int64    `json:"operator_id"`
	OperatorName     string   `json:"operator_name"`
	Capabilities     []string `json:"capabilities"`       // forklift/heavy_lift/hazmat/cold_chain/fragile/ecommerce
	IsOnline         bool     `json:"is_online"`
	CurrentTaskCount int      `json:"current_task_count"`
	WarehouseID      int64    `json:"warehouse_id"`
}

// HasCapability checks if the operator has the given capability.
func (op *OperatorCapability) HasCapability(cap string) bool {
	for _, c := range op.Capabilities {
		if c == cap {
			return true
		}
	}
	return false
}

// MatchScore computes how well an operator matches a task's required capabilities.
// Returns the number of matching capabilities (higher = better match).
func (op *OperatorCapability) MatchScore(required []string) int {
	if len(required) == 0 {
		return 1 // no special requirements = any operator can do it
	}
	score := 0
	for _, rc := range required {
		if op.HasCapability(rc) {
			score++
		}
	}
	return score
}
