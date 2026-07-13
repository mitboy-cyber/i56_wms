package main

import (
	"context"
	"time"

	custDomain "github.com/i56/modules/customer/domain"
	custRepo "github.com/i56/modules/customer/repository"
	orderDomain "github.com/i56/modules/order/domain"
	orderRepo "github.com/i56/modules/order/repository"
	rbacRepo "github.com/i56/modules/rbac/repository"
	parcelDomain "github.com/i56/modules/parcel/domain"
	parcelRepo "github.com/i56/modules/parcel/repository"
	psRepo "github.com/i56/modules/parcel_service/repository"
	printRepo "github.com/i56/modules/print/repository"
	tmsDomain "github.com/i56/modules/transport/domain"
	tmsRepo "github.com/i56/modules/transport/repository"
	whDomain "github.com/i56/modules/warehouse/domain"
	whRepo "github.com/i56/modules/warehouse/repository"
	whRepo2 "github.com/i56/modules/webhook/repository"
	twoRepo "github.com/i56/modules/workorder/repository"
	reportDomain "github.com/i56/modules/report/domain"
)

// seed populates in-memory repositories with demo data for development.
func seed(
	rbac *rbacRepo.MemRBACRepo,
	cr *custRepo.MemClientRepo,
	wr *whRepo.MemWarehouseRepo,
	rr *tmsRepo.MemRouteRepo,
	pr *parcelRepo.MemParcelRepo,
	or *orderRepo.MemOrderRepo,
	_ *tmsRepo.MemCourierRepo,
	lr *custRepo.MemLedgerRepo,
	sr *psRepo.MemServiceRepo,
	wor *twoRepo.MemWorkOrderRepo,
	_ *printRepo.MemPrintRepo,
	whr *whRepo2.MemWebhookRepo,
	rpt *reportDomain.ReportService,
) {
	ctx := context.Background()
	now := time.Now()

	wr.Create(ctx, 1, &whDomain.Warehouse{
		Name: "厦门仓", Code: "XM", Address: "福建省厦门市集美区",
		Contact: "仓库管理员", Phone: "0592-1234567", IsActive: true, TenantID: 1,
	})

	c := &custDomain.Client{
		Name: "EZ集运通", Code: "EZ001", ClientType: custDomain.ClientTypePlatform,
		ContactName: "运营经理", ContactPhone: "13800001111", ContactEmail: "ez@example.com",
		Balance: 10000, IsActive: true, TenantID: 1,
	}
	cr.Create(ctx, 1, c)

	// Seed routes for real route data
	rr.Create(ctx, &tmsDomain.Route{TenantID: 1, WarehouseID: 1, Name: "厦门→台湾(空运)", TransportType: "air", MinWeight: 0.5, VolumeCoeff: 6000, BaseWeightPrice: 20.0, BaseVolumePrice: 20.0, MinAmount: 50, MinDays: 1, MaxDays: 3, IsActive: true})
	rr.Create(ctx, &tmsDomain.Route{TenantID: 1, WarehouseID: 1, Name: "厦门→台湾(海快)", TransportType: "sea_express", MinWeight: 1.0, VolumeCoeff: 6000, BaseWeightPrice: 8.30, BaseVolumePrice: 15.0, MinAmount: 50, MinDays: 3, MaxDays: 7, IsActive: true})
	rr.Create(ctx, &tmsDomain.Route{TenantID: 1, WarehouseID: 1, Name: "厦门→台湾(海运)", TransportType: "sea", MinWeight: 10.0, VolumeCoeff: 6000, BaseWeightPrice: 3.20, BaseVolumePrice: 10.0, MinAmount: 50, MinDays: 5, MaxDays: 14, IsActive: true})

	for i, pd := range []struct {
		tn, pn string
		s      parcelDomain.ParcelStatus
		w      float64
	}{
		{"SF1234567890", "手机壳", "pre_declared", 0.12},
		{"ZTO9876543210", "运动鞋", "received", 0.80},
		{"YTO1111222233", "T恤", "weighed", 0.25},
		{"STO4444555566", "蓝牙耳机", "stored", 0.15},
		{"HTKY7777888899", "数据线", "stored", 0.08},
		{"JD9999000011", "充电宝", "stored", 0.30},
		{"EMS1213141516", "化妆品套装", "stored", 1.20},
	} {
		tn := now.Add(-time.Duration(10-i) * 24 * time.Hour)
		pr.Create(ctx, &parcelDomain.Parcel{
			TenantID: 1, WarehouseID: 1, ClientID: c.ID,
			TrackingNumber: pd.tn, ProductName: pd.pn, ParcelName: pd.pn,
			Status: parcelDomain.ParcelStatus(pd.s), CourierCode: "SF",
			CargoType: "general", ActualWeight: pd.w,
			CreatedAt: tn, UpdatedAt: tn,
		})
	}

	// Seed 8 orders spread across last 7 days
	type orderSeed struct {
		orderNo, recipient, tracking, carrierTrack, customsNo, remark string
		memberID, routeID, daysAgo, parcelCount                       int
		status                                                        orderDomain.OrderStatus
		weight, chgWeight, price                                       float64
	}
	today := now
	orders := []orderSeed{
		{"ORD-20260711-001", "王仁照", "80020737681100020001", "CT-8837291", "CN-20260711001", "空运急件", 1, 2, 0, 1, orderDomain.StatusInTransit, 0.56, 0.60, 8.00},
		{"ORD-20260711-002", "琦立工作室", "YT7631606603205", "", "", "", 2, 2, 0, 2, orderDomain.StatusPendingLoading, 1.05, 1.50, 18.00},
		{"ORD-20260710-001", "张致廷", "HTKY7777888899,JD9999000011", "", "", "", 1, 1, 1, 2, orderDomain.StatusPendingPacking, 0.33, 0.50, 11.50},
		{"ORD-20260709-001", "吴欣如", "ZTO20250601001,SF120011223344", "CT-8837292", "CN-20260709001", "已签收", 2, 3, 2, 3, orderDomain.StatusCompleted, 12.80, 15.00, 56.20},
		{"ORD-20260708-001", "王仁照", "YTO8822110011", "CT-8837293", "CN-20260708001", "", 1, 2, 3, 1, orderDomain.StatusCustomsClearance, 2.30, 2.50, 22.00},
		{"ORD-20260707-001", "琦立工作室", "STO5555666677", "", "", "大件运输", 2, 1, 4, 2, orderDomain.StatusLoaded, 4.50, 5.00, 45.00},
		{"ORD-20260706-001", "张致廷", "EMS9988776655,EMS1122334455", "CT-8837294", "CN-20260706001", "", 1, 3, 5, 4, orderDomain.StatusShipped, 28.50, 30.00, 98.00},
		{"ORD-20260705-001", "吴欣如", "SF5566778899,YTO4433221100", "", "", "待拣货", 2, 2, 6, 2, orderDomain.StatusPendingPicking, 0.78, 1.00, 9.50},
	}
	for _, od := range orders {
		or.Create(ctx, &orderDomain.Order{
			TenantID: 1, WarehouseID: 1, ClientID: c.ID,
			OrderNo: od.orderNo, MemberID: int64(od.memberID), RouteID: int64(od.routeID),
			RecipientName: od.recipient, TrackingNumbers: od.tracking,
			Status: od.status, ParcelCount: od.parcelCount,
			TotalActualWeight: od.weight, TotalChargeableWeight: od.chgWeight,
			TotalPrice: od.price, CarrierTrackingNo: od.carrierTrack,
			CustomsNumber: od.customsNo, Remark: od.remark,
		})
	}

	// Patch seed order dates to match their order_no dates
	if o1, _ := or.GetByOrderNo(ctx, 1, "ORD-20260711-001"); o1 != nil { o1.CreatedAt = today; o1.UpdatedAt = today; or.Update(ctx, o1) }
	if o2, _ := or.GetByOrderNo(ctx, 1, "ORD-20260711-002"); o2 != nil { o2.CreatedAt = today; o2.UpdatedAt = today; or.Update(ctx, o2) }
	if o3, _ := or.GetByOrderNo(ctx, 1, "ORD-20260710-001"); o3 != nil { o3.CreatedAt = today.Add(-24 * time.Hour); o3.UpdatedAt = today.Add(-24 * time.Hour); or.Update(ctx, o3) }
	if o4, _ := or.GetByOrderNo(ctx, 1, "ORD-20260709-001"); o4 != nil { o4.CreatedAt = today.Add(-2 * 24 * time.Hour); o4.UpdatedAt = today.Add(-2 * 24 * time.Hour); or.Update(ctx, o4) }
	if o5, _ := or.GetByOrderNo(ctx, 1, "ORD-20260708-001"); o5 != nil { o5.CreatedAt = today.Add(-3 * 24 * time.Hour); o5.UpdatedAt = today.Add(-3 * 24 * time.Hour); or.Update(ctx, o5) }
	if o6, _ := or.GetByOrderNo(ctx, 1, "ORD-20260707-001"); o6 != nil { o6.CreatedAt = today.Add(-4 * 24 * time.Hour); o6.UpdatedAt = today.Add(-4 * 24 * time.Hour); or.Update(ctx, o6) }
	if o7, _ := or.GetByOrderNo(ctx, 1, "ORD-20260706-001"); o7 != nil { o7.CreatedAt = today.Add(-5 * 24 * time.Hour); o7.UpdatedAt = today.Add(-5 * 24 * time.Hour); or.Update(ctx, o7) }
	if o8, _ := or.GetByOrderNo(ctx, 1, "ORD-20260705-001"); o8 != nil { o8.CreatedAt = today.Add(-6 * 24 * time.Hour); o8.UpdatedAt = today.Add(-6 * 24 * time.Hour); or.Update(ctx, o8) }

	lr.Add(ctx, &custRepo.LedgerEntry{TenantID: 1, ClientID: c.ID, Amount: 5000, BalanceAfter: 5000, Type: "recharge", Description: ""})
	_ = sr
	_ = wor
	_ = whr
	_ = rpt
}
