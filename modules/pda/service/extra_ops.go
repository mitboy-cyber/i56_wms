package service
import (
	"context"; "fmt"; "time"
	parcelDomain "github.com/i56/modules/parcel/domain"
	orderDomain "github.com/i56/modules/order/domain"
)

// SendToOutbound moves packed order to outbound zone
func (p *PDAOperations) SendToOutbound(ctx context.Context, opID int64, orderNo string) error {
	orders := p.orders.ListAll(ctx, 1)
	for _, o := range orders {
		if o.OrderNo == orderNo || fmt.Sprintf("ORD-%d", o.ID) == orderNo {
			if o.Status != orderDomain.StatusPendingPacking && o.Status != orderDomain.StatusPicking {
				// Allow any post-pick status
			}
			o.Status = orderDomain.StatusShipped; o.UpdatedAt = time.Now()
			p.orders.Update(ctx, &o)
			p.log("outbound", orderNo, orderNo, "已送出库区", opID, true)
			return nil
		}
	}
	return fmt.Errorf("订单 %s 未找到", orderNo)
}

// SendToContainerZone moves order to container loading zone
func (p *PDAOperations) SendToContainerZone(ctx context.Context, opID int64, orderNo string) error {
	orders := p.orders.ListAll(ctx, 1)
	for _, o := range orders {
		if o.OrderNo == orderNo || fmt.Sprintf("ORD-%d", o.ID) == orderNo {
			o.Status = orderDomain.StatusPendingLoading; o.UpdatedAt = time.Now()
			p.orders.Update(ctx, &o)
			p.log("container_zone", orderNo, orderNo, "已送装柜区", opID, true)
			return nil
		}
	}
	return fmt.Errorf("订单 %s 未找到", orderNo)
}

// ForceAssign assigns a task to an operator
type AssignedTask struct {
	ID           int64  `json:"id"`
	TaskType     string `json:"task_type"` // receive/pick/pack/load
	TargetID     int64  `json:"target_id"` // parcel_id or order_id
	AssignedTo   int64  `json:"assigned_to"`
	AssignedName string `json:"assigned_name"`
	AssignedAt   string `json:"assigned_at"`
}

var assignedTasks []AssignedTask

func (p *PDAOperations) ForceAssign(ctx context.Context, opID int64, taskType string, targetID int64, assignTo int64) (*AssignedTask, error) {
	at := &AssignedTask{
		ID: int64(len(assignedTasks) + 1),
		TaskType: taskType, TargetID: targetID,
		AssignedTo: assignTo, AssignedName: fmt.Sprintf("操作员%d", assignTo),
		AssignedAt: time.Now().Format("15:04:05"),
	}
	assignedTasks = append(assignedTasks, *at)
	p.log("force_assign", taskType, fmt.Sprintf("#%d→op%d", targetID, assignTo), "强制指派", opID, true)
	return at, nil
}

func (p *PDAOperations) GetAssignedTasks(opID int64) []AssignedTask {
	var result []AssignedTask
	for _, t := range assignedTasks {
		if t.AssignedTo == opID { result = append(result, t) }
	}
	return result
}

func (p *PDAOperations) PendingLoad(ctx context.Context) []orderDomain.Order {
	return p.orders.ListByStatus(ctx, 1, orderDomain.StatusPendingLoading)
}

// Warehouse console stats
type WarehouseStats struct {
	PendingReceive int64 `json:"pending_receive"`
	PendingPutAway int64 `json:"pending_putaway"`
	PendingPick    int64 `json:"pending_pick"`
	PendingPack    int64 `json:"pending_pack"`
	PendingLoad    int64 `json:"pending_load"`
	PendingShip    int64 `json:"pending_ship"`
	AbnormalCount  int64 `json:"abnormal_count"`
}

func (p *PDAOperations) WarehouseStats(ctx context.Context) WarehouseStats {
	var ws WarehouseStats
	for _, p := range p.parcels.ListByStatus(ctx, 1, parcelDomain.StatusPreDeclared) { _ = p; ws.PendingReceive++ }
	for _, p := range p.parcels.ListByStatus(ctx, 1, parcelDomain.StatusReceived) { _ = p; ws.PendingPutAway++ }
	ws.PendingPick = int64(len(p.orders.ListByStatus(ctx, 1, orderDomain.StatusPendingPicking)))
	ws.PendingPack = int64(len(p.orders.ListByStatus(ctx, 1, orderDomain.StatusPicking)))
	ws.PendingLoad = int64(len(p.orders.ListByStatus(ctx, 1, orderDomain.StatusPendingLoading)))
	ws.AbnormalCount = int64(len(p.parcels.ListByStatus(ctx, 1, "abnormal")))
	return ws
}
