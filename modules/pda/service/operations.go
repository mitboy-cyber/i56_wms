package service
import (
	"context"; "fmt"; "sync"; "time"
	"github.com/i56/modules/pda/domain"
	parcelDomain "github.com/i56/modules/parcel/domain"
	parcelRepo "github.com/i56/modules/parcel/repository"
	orderDomain "github.com/i56/modules/order/domain"
	orderRepo "github.com/i56/modules/order/repository"
	wmsRepo "github.com/i56/modules/wms/repository"
	tmsRepo "github.com/i56/modules/tms/repository"
)

// WarehouseOperator carries full context for PDA operations
type WarehouseOperator struct {
	OperatorID  int64
	WarehouseID int64
	TenantID    int64
}

type PDAOperations struct {
	mu         sync.RWMutex
	parcels    *parcelRepo.MemParcelRepo
	orders     *orderRepo.MemOrderRepo
	wms        *wmsRepo.MemWMSRepo
	tms        *tmsRepo.MemTMSRepo
	scanLogs   []domain.ScanLog
	operatorStore map[int64]*WarehouseOperator
}

func NewPDAOperations(pr *parcelRepo.MemParcelRepo, or *orderRepo.MemOrderRepo, wms *wmsRepo.MemWMSRepo, tms *tmsRepo.MemTMSRepo) *PDAOperations {
	return &PDAOperations{
		parcels: pr, orders: or, wms: wms, tms: tms,
		operatorStore: map[int64]*WarehouseOperator{
			1: {OperatorID: 1, WarehouseID: 1, TenantID: 1},
			2: {OperatorID: 2, WarehouseID: 1, TenantID: 1},
			3: {OperatorID: 3, WarehouseID: 1, TenantID: 1},
		},
	}
}

func (p *PDAOperations) GetOperator(opID int64) *WarehouseOperator {
	p.mu.RLock(); defer p.mu.RUnlock()
	return p.operatorStore[opID]
}

func (p *PDAOperations) log(action, barcode, trackingNo, msg string, opID int64, ok bool) {
	p.mu.Lock(); defer p.mu.Unlock()
	p.scanLogs = append(p.scanLogs, domain.ScanLog{
		TenantID: 1, WarehouseID: 1, OperatorID: opID,
		Action: action, Barcode: barcode, TrackingNumber: trackingNo,
		Success: ok, Message: msg, ScannedAt: time.Now(),
	})
}

func (p *PDAOperations) RecentLogs(limit int) []domain.ScanLog {
	p.mu.RLock(); defer p.mu.RUnlock()
	if limit > len(p.scanLogs) { limit = len(p.scanLogs) }
	if limit <= 0 { return nil }
	return p.scanLogs[len(p.scanLogs)-limit:]
}

// ============================
// 1. RECEIVE (入库) — scan→weigh→measure→receive
// ============================
func (p *PDAOperations) Receive(ctx context.Context, opID int64, trackingNo string, weight, length, width, height float64, locationBarcode string) (*parcelDomain.Parcel, error) {
	pr, err := p.parcels.GetByTrackingNo(ctx, 1, trackingNo)
	if err != nil || pr == nil { return nil, fmt.Errorf("包裹 %s 未找到", trackingNo) }
	if pr.Status != parcelDomain.StatusPreDeclared { return nil, fmt.Errorf("包裹状态为 %s，不能入库", pr.Status) }

	pr.ActualWeight = weight; pr.Length = length; pr.Width = width; pr.Height = height
	pr.Status = parcelDomain.StatusReceived
	if locationBarcode != "" {
		loc := p.wms.GetLocationByBarcode(ctx, locationBarcode)
		if loc != nil { pr.LocationCode = loc.Code }
	}
	pr.UpdatedAt = time.Now()
	if err := p.parcels.Update(ctx, pr); err != nil { return nil, err }
	p.log("receive", trackingNo, trackingNo, fmt.Sprintf("入库 %.2fkg", weight), opID, true)
	return pr, nil
}

// ============================
// 2. WEIGH (核重) — scan→show actual weight→confirm
// ============================
func (p *PDAOperations) Weigh(ctx context.Context, opID int64, trackingNo string, weight float64) (*parcelDomain.Parcel, float64, error) {
	pr, err := p.parcels.GetByTrackingNo(ctx, 1, trackingNo)
	if err != nil || pr == nil { return nil, 0, fmt.Errorf("包裹 %s 未找到", trackingNo) }
	if pr.Status != parcelDomain.StatusReceived && pr.Status != parcelDomain.StatusPreDeclared {
		return nil, 0, fmt.Errorf("包裹状态为 %s，不能核重", pr.Status)
	}
	oldW := pr.ActualWeight
	pr.ActualWeight = weight; pr.Status = parcelDomain.StatusWeighed; pr.UpdatedAt = time.Now()
	p.parcels.Update(ctx, pr)
	p.log("weigh", trackingNo, trackingNo, fmt.Sprintf("核重 %.2f→%.2fkg", oldW, weight), opID, true)
	return pr, oldW, nil
}

// ============================
// 3. PUT-AWAY (上架) — received→stored. Scan parcel + location
// ============================
func (p *PDAOperations) PutAway(ctx context.Context, opID int64, trackingNo, locationBarcode string) (*parcelDomain.Parcel, error) {
	pr, err := p.parcels.GetByTrackingNo(ctx, 1, trackingNo)
	if err != nil || pr == nil { return nil, fmt.Errorf("包裹 %s 未找到", trackingNo) }
	if pr.Status != parcelDomain.StatusReceived && pr.Status != parcelDomain.StatusWeighed {
		return nil, fmt.Errorf("包裹状态为 %s，不能上架（需先入库或核重）", pr.Status)
	}
	loc := p.wms.GetLocationByBarcode(ctx, locationBarcode)
	if loc == nil { return nil, fmt.Errorf("库位 %s 不存在", locationBarcode) }
	if loc.IsOccupied { return nil, fmt.Errorf("库位 %s 已被占用", locationBarcode) }

	loc.IsOccupied = true; loc.CurrentParcelID = &pr.ID
	pr.LocationCode = loc.Code; pr.Status = parcelDomain.StatusStored; pr.UpdatedAt = time.Now()
	p.parcels.Update(ctx, pr)
	p.log("putaway", locationBarcode, trackingNo, fmt.Sprintf("→ 库位 %s", loc.Code), opID, true)
	return pr, nil
}

// ============================
// 4. PICK (拣货) — scan order→show parcels→confirm
// ============================
func (p *PDAOperations) Pick(ctx context.Context, opID int64, orderNo string) (*orderDomain.Order, []parcelDomain.Parcel, error) {
	orders := p.orders.ListAll(ctx, 1)
	for _, o := range orders {
		if o.OrderNo == orderNo || fmt.Sprintf("ORD-%d", o.ID) == orderNo {
			if o.Status != orderDomain.StatusPendingPicking {
				return &o, nil, fmt.Errorf("订单状态为 %s，不能拣货", o.Status)
			}
			o.Status = orderDomain.StatusPicking; o.UpdatedAt = time.Now()
			p.orders.Update(ctx, &o)
			// Get parcels by tracking numbers (split by 、)
			var parcels []parcelDomain.Parcel
			tns := splitTrackingNos(o.TrackingNumbers)
			for _, tn := range tns {
				pr, _ := p.parcels.GetByTrackingNo(ctx, 1, tn)
				if pr != nil {
					pr.Status = parcelDomain.StatusPicked; pr.UpdatedAt = time.Now()
					p.parcels.Update(ctx, pr)
					parcels = append(parcels, *pr)
				}
			}
			p.log("pick", orderNo, orderNo, fmt.Sprintf("拣货 %d 件", len(parcels)), opID, true)
			return &o, parcels, nil
		}
	}
	return nil, nil, fmt.Errorf("订单 %s 未找到", orderNo)
}

func splitTrackingNos(s string) []string {
	var result []string
	cur := ""
	for _, c := range s {
		if c == '、' || c == '，' || c == ',' {
			if cur != "" { result = append(result, cur); cur = "" }
		} else { cur += string(c) }
	}
	if cur != "" { result = append(result, cur) }
	return result
}

// ============================
// 5. PACK (打包) — scan parcels→verify→seal
// ============================
func (p *PDAOperations) Pack(ctx context.Context, opID int64, orderNo string) (*orderDomain.Order, error) {
	orders := p.orders.ListAll(ctx, 1)
	for _, o := range orders {
		if o.OrderNo == orderNo || fmt.Sprintf("ORD-%d", o.ID) == orderNo {
			if o.Status != orderDomain.StatusPicking && o.Status != orderDomain.StatusPendingPacking {
				return &o, fmt.Errorf("订单状态为 %s，不能打包", o.Status)
			}
			o.Status = orderDomain.StatusPendingPacking; o.UpdatedAt = time.Now()
			p.orders.Update(ctx, &o)
			tns := splitTrackingNos(o.TrackingNumbers)
			for _, tn := range tns {
				pr, _ := p.parcels.GetByTrackingNo(ctx, 1, tn)
				if pr != nil { pr.Status = parcelDomain.StatusPacked; pr.UpdatedAt = time.Now(); p.parcels.Update(ctx, pr) }
			}
			p.log("pack", orderNo, orderNo, "打包完成", opID, true)
			return &o, nil
		}
	}
	return nil, fmt.Errorf("订单 %s 未找到", orderNo)
}

// ============================
// 6. LOAD (装柜) — scan container→scan orders→confirm
// ============================
func (p *PDAOperations) LoadContainer(ctx context.Context, opID int64, containerNo, orderNo string) error {
	orders := p.orders.ListAll(ctx, 1)
	for _, o := range orders {
		if o.OrderNo == orderNo || fmt.Sprintf("ORD-%d", o.ID) == orderNo {
			o.Status = orderDomain.StatusLoaded; dummy := int64(1); o.ContainerLoadingID = &dummy; o.UpdatedAt = time.Now()
			p.orders.Update(ctx, &o)
			p.log("load", containerNo, orderNo, "装柜完成", opID, true)
			return nil
		}
	}
	return fmt.Errorf("订单 %s 未找到", orderNo)
}

// ============================
// 7. MARK EXCEPTION (标异常)
// ============================
func (p *PDAOperations) MarkException(ctx context.Context, opID int64, trackingNo, reason string) error {
	pr, err := p.parcels.GetByTrackingNo(ctx, 1, trackingNo)
	if err != nil || pr == nil { return fmt.Errorf("包裹 %s 未找到", trackingNo) }
	pr.Status = "abnormal"; pr.IsAbnormal = true; pr.UpdatedAt = time.Now()
	p.parcels.Update(ctx, pr)
	p.log("exception", trackingNo, trackingNo, reason, opID, true)
	return nil
}

// ============================
// Support: list pending parcels for each operation
// ============================
func (p *PDAOperations) PendingReceive(ctx context.Context) []parcelDomain.Parcel {
	return p.parcels.ListByStatus(ctx, 1, parcelDomain.StatusPreDeclared)
}
func (p *PDAOperations) PendingPutAway(ctx context.Context) []parcelDomain.Parcel {
	return append(p.parcels.ListByStatus(ctx, 1, parcelDomain.StatusReceived),
		p.parcels.ListByStatus(ctx, 1, parcelDomain.StatusWeighed)...)
}
func (p *PDAOperations) PendingPick(ctx context.Context) []orderDomain.Order {
	return p.orders.ListByStatus(ctx, 1, orderDomain.StatusPendingPicking)
}

// GetOrderParcels returns all parcels associated with an order by its tracking numbers
func (p *PDAOperations) GetOrderParcels(ctx context.Context, order *orderDomain.Order) []parcelDomain.Parcel {
	var parcels []parcelDomain.Parcel
	tns := splitTrackingNos(order.TrackingNumbers)
	for _, tn := range tns {
		pr, _ := p.parcels.GetByTrackingNo(ctx, 1, tn)
		if pr != nil {
			parcels = append(parcels, *pr)
		}
	}
	return parcels
}

// ============================
// Pending lists for each operation
// ============================
func (p *PDAOperations) PendingWeigh(ctx context.Context) []parcelDomain.Parcel {
	return p.parcels.ListByStatus(ctx, 1, parcelDomain.StatusReceived)
}

func (p *PDAOperations) PendingPack(ctx context.Context) []orderDomain.Order {
	return p.orders.ListByStatus(ctx, 1, orderDomain.StatusPicking)
}

// ============================
// Lookup helpers for scan-first workflows
// ============================
func (p *PDAOperations) GetParcelForReceive(ctx context.Context, trackingNo string) (*parcelDomain.Parcel, error) {
	pr, err := p.parcels.GetByTrackingNo(ctx, 1, trackingNo)
	if err != nil || pr == nil { return nil, fmt.Errorf("包裹 %s 未找到", trackingNo) }
	return pr, nil
}

func (p *PDAOperations) GetParcelForWeigh(ctx context.Context, trackingNo string) (*parcelDomain.Parcel, float64, error) {
	pr, err := p.parcels.GetByTrackingNo(ctx, 1, trackingNo)
	if err != nil || pr == nil { return nil, 0, fmt.Errorf("包裹 %s 未找到", trackingNo) }
	return pr, pr.ActualWeight, nil
}

func (p *PDAOperations) GetParcelForPutAway(ctx context.Context, trackingNo string) (*parcelDomain.Parcel, float64, error) {
	pr, err := p.parcels.GetByTrackingNo(ctx, 1, trackingNo)
	if err != nil || pr == nil { return nil, 0, fmt.Errorf("包裹 %s 未找到", trackingNo) }
	return pr, pr.ActualWeight, nil
}

func (p *PDAOperations) LookupOrder(ctx context.Context, orderNo string) (*orderDomain.Order, []parcelDomain.Parcel, error) {
	orders := p.orders.ListAll(ctx, 1)
	for _, o := range orders {
		if o.OrderNo == orderNo || fmt.Sprintf("ORD-%d", o.ID) == orderNo {
			parcels := p.GetOrderParcels(ctx, &o)
			return &o, parcels, nil
		}
	}
	return nil, nil, fmt.Errorf("订单 %s 未找到", orderNo)
}

func (p *PDAOperations) LookupOrderForPack(ctx context.Context, orderNo string) (*orderDomain.Order, []parcelDomain.Parcel, error) {
	orders := p.orders.ListAll(ctx, 1)
	for _, o := range orders {
		if o.OrderNo == orderNo || fmt.Sprintf("ORD-%d", o.ID) == orderNo {
			if o.Status != orderDomain.StatusPicking {
				return &o, nil, fmt.Errorf("订单状态为 %s，不能打包（需先拣货）", o.Status)
			}
			parcels := p.GetOrderParcels(ctx, &o)
			return &o, parcels, nil
		}
	}
	return nil, nil, fmt.Errorf("订单 %s 未找到", orderNo)
}
