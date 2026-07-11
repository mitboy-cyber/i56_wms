package domain
import "time"

type TaskDispatchConfig struct {
	ID            int64     `json:"id"`
	TenantID      int64     `json:"tenant_id"`
	WarehouseID   int64     `json:"warehouse_id"`
	TaskType      string    `json:"task_type"`      // receive | pick | pack | load
	Strategy     string     `json:"strategy"`        // round_robin | least_busy | nearest
	MaxTasksPerOp int       `json:"max_tasks_per_op"`
	AutoAssign    bool      `json:"auto_assign"`
	Priority      int       `json:"priority"`
	CreatedAt     time.Time `json:"created_at"`
}
