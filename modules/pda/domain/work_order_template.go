package domain

import "time"

// WorkOrderTemplate represents a predefined template for creating work orders.
type WorkOrderTemplate struct {
	ID          int64     `json:"id"`
	TenantID    int64     `json:"tenant_id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Steps       []string  `json:"steps"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func DefaultWorkOrderTemplates() []WorkOrderTemplate {
	return []WorkOrderTemplate{
		{Name: "标准收货流程", Code: "receive", Type: "receive", Description: "包裹收货→称重→核重→上架", Steps: []string{"scan_tracking", "weigh", "verify_weight", "putaway"}},
		{Name: "标准拣货流程", Code: "pick", Type: "pick", Description: "按单拣货→扫描确认→移至打包区", Steps: []string{"scan_order", "scan_parcel", "confirm_pick", "move_to_pack"}},
		{Name: "标准打包流程", Code: "pack", Type: "pack", Description: "扫描→打包→扫描确认→出库", Steps: []string{"scan_parcel", "pack", "confirm_pack", "outbound"}},
		{Name: "装柜流程", Code: "load", Type: "load", Description: "扫描容器→扫描包裹→确认装柜", Steps: []string{"scan_container", "scan_parcel", "confirm_load"}},
		{Name: "异常处理流程", Code: "abnormal", Type: "exception", Description: "登记异常→拍照→提交处理", Steps: []string{"register_exception", "take_photo", "submit"}},
	}
}
