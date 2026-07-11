package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/i56/framework/core/auth"
	"github.com/i56/framework/core/config"
	"github.com/i56/framework/core/middleware"
	"github.com/i56/framework/core/response"
	"github.com/i56/framework/core/router"

	custDomain "github.com/i56/framework/internal/modules/customer/domain"
	custHandler "github.com/i56/framework/internal/modules/customer/handler"
	custRepo "github.com/i56/framework/internal/modules/customer/repository"
	custSvc "github.com/i56/framework/internal/modules/customer/service"

	whDomain "github.com/i56/framework/internal/modules/warehouse/domain"
	whRepo "github.com/i56/framework/internal/modules/warehouse/repository"
	whSvc "github.com/i56/framework/internal/modules/warehouse/service"

	parcelDomain "github.com/i56/framework/internal/modules/parcel/domain"
	parcelRepo "github.com/i56/framework/internal/modules/parcel/repository"
	parcelSvc "github.com/i56/framework/internal/modules/parcel/service"

	orderDomain "github.com/i56/framework/internal/modules/order/domain"
	orderRepo "github.com/i56/framework/internal/modules/order/repository"
	orderSvc "github.com/i56/framework/internal/modules/order/service"

	tmsDomain "github.com/i56/framework/internal/modules/transport/domain"
	tmsRepo "github.com/i56/framework/internal/modules/transport/repository"

	psRepo "github.com/i56/framework/internal/modules/parcel_service/repository"
	psDomain "github.com/i56/framework/internal/modules/parcel_service/domain"
	woRepo "github.com/i56/framework/internal/modules/workorder/repository"
	woDomain "github.com/i56/framework/internal/modules/workorder/domain"
	printRepo "github.com/i56/framework/internal/modules/print/repository"
	reportDomain "github.com/i56/framework/internal/modules/report/domain"
	whRepo2 "github.com/i56/framework/internal/modules/webhook/repository"
	whDomain2 "github.com/i56/framework/internal/modules/webhook/domain"
)

func setupRouter() http.Handler {
	cfg, _ := config.Load()
	tm, _ := auth.NewTokenManager(cfg.Auth)

	clientRepo := custRepo.NewMemClientRepo()
	warehouseRepo := whRepo.NewMemWarehouseRepo()
	pRepo := parcelRepo.NewMemParcelRepo()
	oRepo := orderRepo.NewMemOrderRepo()
	routeRepo := tmsRepo.NewMemRouteRepo()
	courierRepo := tmsRepo.NewMemCourierRepo()
	ledgerRepo := custRepo.NewMemLedgerRepo()
	svcRepo := psRepo.NewMemServiceRepo()
	woR := woRepo.NewMemWorkOrderRepo()
	prRepo := printRepo.NewMemPrintRepo()
	whRepo := whRepo2.NewMemWebhookRepo()
	reportSvc := reportDomain.NewReportService()

	clientSvc := custSvc.NewClientService(clientRepo)
	warehouseSvc := whSvc.NewWarehouseService(warehouseRepo)
	parcelSvc := parcelSvc.NewParcelService(pRepo)
	orderSvc := orderSvc.NewOrderService(oRepo)
	clientH := custHandler.NewClientHandler(clientSvc)

	// Seed
	seedIntegration(clientRepo, warehouseRepo, routeRepo, pRepo, oRepo, courierRepo, ledgerRepo, svcRepo, woR, prRepo, whRepo, reportSvc)

	r := router.New()
	r.Use(middleware.RequestID())

	r.GET("/api/health", func(w http.ResponseWriter, req *http.Request) {
		response.JSON(w, 200, map[string]string{"status": "ok", "name": cfg.App.Name, "version": cfg.App.Version})
	})

	r.POST("/api/v1/auth/login", func(w http.ResponseWriter, req *http.Request) {
		var b struct{ Username, Password string }
		json.NewDecoder(req.Body).Decode(&b)
		if b.Username == "" || b.Password == "" {
			response.Error(w, nil)
			return
		}
		at, _ := tm.IssueAccessToken(b.Username, "tenant-1", []string{"admin"}, []string{"*"})
		rt, _ := tm.IssueRefreshToken(b.Username, "tenant-1")
		response.JSON(w, 200, map[string]any{"access_token": at, "refresh_token": rt, "token_type": "Bearer", "expires_in": int(tm.AccessTTL().Seconds())})
	})

	api := router.New().WithPrefix("/api/v1")
	api.Use(middleware.AuthRequired(tm))

	api.GET("/me", func(w http.ResponseWriter, req *http.Request) {
		c := req.Context().Value(middleware.ClaimsKey).(*auth.Claims)
		response.JSON(w, 200, map[string]any{"user_id": c.Subject, "tenant_id": c.TenantID, "roles": c.Roles})
	})
	api.GET("/clients", clientH.List)
	api.GET("/clients/{id}", clientH.GetByID)
	api.POST("/clients", clientH.Create)
	api.PATCH("/clients/{id}", clientH.Update)
	api.DELETE("/clients/{id}", clientH.Delete)

	api.GET("/warehouses", paginatedHandler(func(ctx context.Context) ([]any, int64, error) {
		ws, t, e := warehouseSvc.List(ctx, 1, 0, 50)
		d := make([]any, len(ws))
		for i, w := range ws { d[i] = w }
		return d, t, e
	}))
	api.GET("/parcels", paginatedHandler(func(ctx context.Context) ([]any, int64, error) {
		ps, t, e := parcelSvc.List(ctx, 1, 0, 100)
		d := make([]any, len(ps))
		for i, p := range ps { d[i] = p }
		return d, t, e
	}))
	api.GET("/orders", paginatedHandler(func(ctx context.Context) ([]any, int64, error) {
		os, t, e := orderSvc.List(ctx, 1, 0, 100)
		d := make([]any, len(os))
		for i, o := range os { d[i] = o }
		return d, t, e
	}))
	api.GET("/routes", paginatedHandler(func(ctx context.Context) ([]any, int64, error) {
		rs, t, e := routeRepo.List(ctx, 1, 0, 50)
		d := make([]any, len(rs))
		for i, r := range rs { d[i] = r }
		return d, t, e
	}))
	api.GET("/couriers", func(w http.ResponseWriter, req *http.Request) {
		cs, _ := courierRepo.List(req.Context())
		d := make([]any, len(cs))
		for i, c := range cs { d[i] = c }
		response.PaginatedJSON(w, d, int64(len(cs)), 1, len(cs))
	})
	api.GET("/dashboard", func(w http.ResponseWriter, req *http.Request) {
		ps, pt, _ := parcelSvc.List(req.Context(), 1, 0, 1000)
		os, ot, _ := orderSvc.List(req.Context(), 1, 0, 1000)
		sc := map[string]int{}
		for _, p := range ps { sc[string(p.Status)]++ }
		response.JSON(w, 200, map[string]any{"total_parcels": pt, "total_orders": ot, "parcel_status": sc, "active_orders": len(os)})
	})

	// P1 endpoints
	api.GET("/services/types", func(w http.ResponseWriter, req *http.Request) { response.JSON(w, 200, svcRepo.ListTypes()) })
	api.GET("/services/orders", paginatedHandler(func(ctx context.Context) ([]any, int64, error) {
		os, t, e := svcRepo.List(ctx, 1, 0, 50)
		d := make([]any, len(os))
		for i, o := range os { d[i] = o }
		return d, t, e
	}))
	api.GET("/workorders", paginatedHandler(func(ctx context.Context) ([]any, int64, error) {
		wos, t, e := woR.List(ctx, 1, 0, 50)
		d := make([]any, len(wos))
		for i, wo := range wos { d[i] = wo }
		return d, t, e
	}))
	api.GET("/prints", func(w http.ResponseWriter, req *http.Request) {
		pts, _ := prRepo.List(req.Context(), 1)
		response.JSON(w, 200, pts)
	})
	api.GET("/reports/orders", func(w http.ResponseWriter, req *http.Request) { response.JSON(w, 200, reportSvc.OrderProfit()) })
	api.GET("/reports/clients", func(w http.ResponseWriter, req *http.Request) { response.JSON(w, 200, reportSvc.ClientProfit()) })
	api.GET("/reports/routes", func(w http.ResponseWriter, req *http.Request) { response.JSON(w, 200, reportSvc.RouteProfit()) })
	api.GET("/webhooks", func(w http.ResponseWriter, req *http.Request) {
		subs, _ := whRepo.ListSubs(req.Context(), 1)
		response.JSON(w, 200, subs)
	})
	api.POST("/webhooks", func(w http.ResponseWriter, req *http.Request) {
		var s whDomain2.WebhookSubscription
		json.NewDecoder(req.Body).Decode(&s)
		s.TenantID = 1
		whRepo.CreateSub(req.Context(), &s)
		response.Created(w, s)
	})
	api.GET("/webhooks/logs", func(w http.ResponseWriter, req *http.Request) {
		logs := whRepo.ListLogs(req.Context(), 50)
		response.JSON(w, 200, logs)
	})

	r.Handle("/api/v1/", api)

	return r
}

func paginatedHandler(fn func(context.Context) ([]any, int64, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, t, err := fn(r.Context())
		if err != nil { response.Error(w, err); return }
		response.PaginatedJSON(w, d, t, 1, 50)
	}
}

func seedIntegration(
	cr *custRepo.MemClientRepo, wr *whRepo.MemWarehouseRepo, rr *tmsRepo.MemRouteRepo,
	pr *parcelRepo.MemParcelRepo, or *orderRepo.MemOrderRepo, _ *tmsRepo.MemCourierRepo, lr *custRepo.MemLedgerRepo,
	sr *psRepo.MemServiceRepo, woR *woRepo.MemWorkOrderRepo, prR *printRepo.MemPrintRepo,
	whR *whRepo2.MemWebhookRepo, rpt *reportDomain.ReportService,
) {
	ctx := context.Background()
	now := time.Now()
	wr.Create(ctx, 1, &whDomain.Warehouse{Name: "Test仓", Code: "TST", Address: "测试地址", Contact: "测试", Phone: "000", IsActive: true, TenantID: 1})
	for _, rt := range []tmsDomain.Route{
		{Name: "空运", TransportType: "air", BaseWeightPrice: 25, BaseVolumePrice: 12, MinAmount: 50, MinDays: 3, MaxDays: 7, TenantID: 1, WarehouseID: 1, VolumeCoeff: 6000, WeightRounding: 0.5, IsActive: true, CargoTypes: []string{"general"}},
		{Name: "海快", TransportType: "sea_express", BaseWeightPrice: 8, BaseVolumePrice: 6, MinAmount: 30, MinDays: 7, MaxDays: 15, TenantID: 1, WarehouseID: 1, VolumeCoeff: 6000, WeightRounding: 0.5, IsActive: true, CargoTypes: []string{"general"}},
	} {
		rr.Create(ctx, &rt)
	}
	c := &custDomain.Client{Name: "Test客户", Code: "TST001", ClientType: custDomain.ClientTypePlatform, ContactName: "测试", IsActive: true, TenantID: 1, Balance: 1000}
	cr.Create(ctx, 1, c)
	type pd struct{ tn, pn string; s parcelDomain.ParcelStatus; w, l, h, d float64 }
	for i, pd := range []pd{
		{"SF1001", "手机壳", parcelDomain.StatusPreDeclared, 0, 0, 0, 0},
		{"ZTO2001", "运动鞋", parcelDomain.StatusReceived, 0.85, 30, 20, 12},
		{"YTO3001", "T恤", parcelDomain.StatusWeighed, 0.35, 25, 18, 5},
		{"STO4001", "耳机", parcelDomain.StatusStored, 0.12, 12, 8, 4},
		{"HTKY5001", "数据线", parcelDomain.StatusStored, 0.08, 8, 5, 3},
		{"JD6001", "充电宝", parcelDomain.StatusStored, 0.25, 15, 10, 5},
		{"EMS7001", "化妆品", parcelDomain.StatusStored, 1.2, 22, 16, 10},
		{"YT8001", "键盘", parcelDomain.StatusPicked, 0.9, 45, 16, 5},
		{"YT9001", "鼠标", parcelDomain.StatusPacked, 0.15, 12, 7, 4},
		{"YT0001", "手表", parcelDomain.StatusShipped, 0.56, 8, 8, 6},
	} {
		tn := now.Add(-time.Duration(10-i) * 24 * time.Hour)
		pr.Create(ctx, &parcelDomain.Parcel{TenantID: 1, WarehouseID: 1, ClientID: c.ID, TrackingNumber: pd.tn, ProductName: pd.pn, ParcelName: pd.pn, Status: pd.s, ActualWeight: pd.w, Length: pd.l, Width: pd.h, Height: pd.d, CourierCode: "SF", CargoType: "general", CreatedAt: tn, UpdatedAt: tn})
	}
	or.Create(ctx, &orderDomain.Order{TenantID: 1, WarehouseID: 1, ClientID: c.ID, MemberID: 1, RouteID: 1, RecipientName: "用户A", TrackingNumbers: "YT0001", Status: orderDomain.StatusInTransit, TotalActualWeight: 0.56, TotalChargeableWeight: 0.56, TotalPrice: 25})
	or.Create(ctx, &orderDomain.Order{TenantID: 1, WarehouseID: 1, ClientID: c.ID, MemberID: 1, RouteID: 2, RecipientName: "用户B", TrackingNumbers: "YT8001,YT9001", Status: orderDomain.StatusPendingLoading, TotalActualWeight: 1.05, TotalChargeableWeight: 1.05, TotalPrice: 12.6})
	or.Create(ctx, &orderDomain.Order{TenantID: 1, WarehouseID: 1, ClientID: c.ID, MemberID: 1, RouteID: 1, RecipientName: "用户C", TrackingNumbers: "HTKY5001,JD6001", Status: orderDomain.StatusPendingPicking, TotalActualWeight: 0.33, TotalChargeableWeight: 0.33, TotalPrice: 50})
	sr.Create(ctx, &psDomain.ServiceOrder{TenantID: 1, ClientID: c.ID, ServiceType: "WOODEN_CRATE", Quantity: 1, TotalPrice: 80, Status: "pending"})
	sr.Create(ctx, &psDomain.ServiceOrder{TenantID: 1, ClientID: c.ID, ServiceType: "CONTENT_PHOTO", Quantity: 3, TotalPrice: 3, Status: "completed"})
	woR.Create(ctx, &woDomain.WorkOrder{TenantID: 1, WarehouseID: 1, Title: "入库作业", Status: "pending", Priority: 1})
	woR.Create(ctx, &woDomain.WorkOrder{TenantID: 1, WarehouseID: 1, Title: "拣货作业", Status: "in_progress", Priority: 2})
	whR.CreateSub(ctx, &whDomain2.WebhookSubscription{TenantID: 1, Event: "order.created", URL: "https://example.com/hook", Secret: "whsec_test", IsActive: true})
	whR.LogDelivery(ctx, &whDomain2.WebhookDeliveryLog{SubscriptionID: 1, Event: "order.created", Payload: `{"order_no":"test"}`, StatusCode: 200})
	rpt.RecordOrder("TST001", "Test客户", "海快", 12.6, 5)
	rpt.RecordOrder("TST002", "Test客户", "空运", 50, 25)
}
