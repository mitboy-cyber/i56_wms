package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/i56/framework/core/router"

	custRepo "github.com/i56/modules/customer/repository"
	custDomain "github.com/i56/modules/customer/domain"
	orderSvc "github.com/i56/modules/order/service"
	orderDomain "github.com/i56/modules/order/domain"
	parcelDomain "github.com/i56/modules/parcel/domain"
	parcelSvc "github.com/i56/modules/parcel/service"
	printRepo "github.com/i56/modules/print/repository"
	psRepo "github.com/i56/modules/parcel_service/repository"
	pricingRepo "github.com/i56/modules/pricing/repository"
	rbacRepo "github.com/i56/modules/rbac/repository"
	tdRepo "github.com/i56/modules/taskdispatch/repository"
	tmsRepo "github.com/i56/modules/transport/repository"
	whSvc "github.com/i56/modules/warehouse/service"
	whRepo2 "github.com/i56/modules/webhook/repository"
	wfRepo "github.com/i56/modules/workflow/repository"
	twoRepo "github.com/i56/modules/workorder/repository"
)

func registerAdminFullAPI(
	r *router.Router, a func(http.HandlerFunc) http.HandlerFunc,
	ps *parcelSvc.ParcelService, osvc *orderSvc.OrderService,
	ws *whSvc.WarehouseService, cr *custRepo.MemClientRepo,
	rr *tmsRepo.MemRouteRepo, cour *tmsRepo.MemCourierRepo,
	sr *psRepo.MemServiceRepo, wor *twoRepo.MemWorkOrderRepo,
	lr *custRepo.MemLedgerRepo, dr *custRepo.MemDeclarantRepo,
	mr *custRepo.MemMemberRepo, ar *custRepo.MemAddressRepo,
	rpr *pricingRepo.MemRoutePriceRepo, dfr *pricingRepo.MemDeliveryFeeRepo,
	scr *pricingRepo.MemSurchargeRepo, acr *pricingRepo.MemApiCredentialRepo,
	rbac *rbacRepo.MemRBACRepo, ppr *printRepo.MemPrintRepo,
	wfr *wfRepo.MemWorkflowRepo, td *tdRepo.MemTaskDispatchRepo,
	whr *whRepo2.MemWebhookRepo,
) {
	var t int64 = 1

	empty := func(label string) http.HandlerFunc {
		return a(func(w http.ResponseWriter, req *http.Request) {
			apiJSON(w, 200, map[string]any{"items": []any{}, "total": 0, "label": label})
		})
	}

	// ═══ P0: Real data from services ═══
	r.GET("/admin/api/warehouses", a(func(w http.ResponseWriter, req *http.Request) {
		wh, _, _ := ws.List(req.Context(), t, 0, 200); apiJSON(w, 200, wh)
	}))
	r.POST("/admin/api/warehouses", a(func(w http.ResponseWriter, req *http.Request) {
		var b struct{ Name, Code, Address, Contact, Phone string }; json.NewDecoder(req.Body).Decode(&b)
		wh, _ := ws.Create(req.Context(), t, b.Name, b.Code, b.Address, b.Contact, b.Phone); apiJSON(w, 201, wh)
	}))

	r.GET("/admin/api/parcels", a(func(w http.ResponseWriter, req *http.Request) {
		px, _, _ := ps.List(req.Context(), t, 0, 200); apiJSON(w, 200, px)
	}))
	r.POST("/admin/api/parcels", a(func(w http.ResponseWriter, req *http.Request) {
		var b struct{ TrackingNumber, ProductName, CourierCode string; WarehouseID int64 }; json.NewDecoder(req.Body).Decode(&b)
		if b.WarehouseID == 0 { b.WarehouseID = 1 }
		p, err := ps.PreDeclare(req.Context(), &parcelDomain.Parcel{TrackingNumber: b.TrackingNumber, ProductName: b.ProductName, TenantID: t, WarehouseID: b.WarehouseID, CourierCode: b.CourierCode})
		if err != nil { apiJSON(w, 400, map[string]string{"error": err.Error()}); return }
		apiJSON(w, 201, p)
	}))

	r.GET("/admin/api/orders", a(func(w http.ResponseWriter, req *http.Request) {
		ox, _, _ := osvc.List(req.Context(), t, 0, 200); apiJSON(w, 200, ox)
	}))
	r.POST("/admin/api/orders", a(func(w http.ResponseWriter, req *http.Request) {
		var b struct{ OrderNo, RecipientName string; ParcelCount int; TotalPrice float64; RouteID int64 }; json.NewDecoder(req.Body).Decode(&b)
		o := &orderDomain.Order{OrderNo: b.OrderNo, RecipientName: b.RecipientName, ParcelCount: b.ParcelCount, TotalPrice: b.TotalPrice, RouteID: b.RouteID, TenantID: t}
		created, err := osvc.Create(req.Context(), o)
		if err != nil { apiJSON(w, 400, map[string]string{"error": err.Error()}); return }
		apiJSON(w, 201, created)
	}))
	r.GET("/admin/api/orders/{id}", a(func(w http.ResponseWriter, req *http.Request) {
		idStr := req.PathValue("id"); id, _ := strconv.ParseInt(idStr, 10, 64)
		o, _ := osvc.GetByOrderNo(req.Context(), t, idStr)
		if o == nil { o, _ = osvc.GetByID(req.Context(), t, id) }
		if o == nil { apiJSON(w, 404, map[string]string{"error": "not found"}); return }
		apiJSON(w, 200, o)
	}))

	r.GET("/admin/api/clients", a(func(w http.ResponseWriter, req *http.Request) {
		cl, _, _ := cr.List(req.Context(), t, 0, 200); apiJSON(w, 200, cl)
	}))
	r.POST("/admin/api/clients", a(func(w http.ResponseWriter, req *http.Request) {
		var b struct{ Name, Code string }; json.NewDecoder(req.Body).Decode(&b)
		c := &custDomain.Client{Name: b.Name, Code: b.Code}; cr.Create(req.Context(), t, c); apiJSON(w, 201, c)
	}))

	r.GET("/admin/api/declarants", a(func(w http.ResponseWriter, req *http.Request) {
		d, _, _ := dr.List(req.Context(), 1, 0, 200); apiJSON(w, 200, d)
	}))
	r.GET("/admin/api/members", a(func(w http.ResponseWriter, req *http.Request) {
		m, _, _ := mr.List(req.Context(), 1, 0, 200); apiJSON(w, 200, m)
	}))
	r.GET("/admin/api/addresses", a(func(w http.ResponseWriter, req *http.Request) {
		addr, _ := ar.List(req.Context(), 1); apiJSON(w, 200, addr)
	}))
	r.GET("/admin/api/ledger", a(func(w http.ResponseWriter, req *http.Request) {
		apiJSON(w, 200, lr.GetByClient(req.Context(), 1, 1))
	}))
	r.GET("/admin/api/service-orders", a(func(w http.ResponseWriter, req *http.Request) {
		so, _, _ := sr.List(req.Context(), t, 0, 200); apiJSON(w, 200, so)
	}))
	r.GET("/admin/api/work-orders", a(func(w http.ResponseWriter, req *http.Request) {
		wo, _, _ := wor.List(req.Context(), t, 0, 200); apiJSON(w, 200, wo)
	}))
	r.GET("/admin/api/carriers", a(func(w http.ResponseWriter, req *http.Request) {
		routes, _, _ := rr.List(req.Context(), t, 0, 200); apiJSON(w, 200, routes)
	}))
	r.GET("/admin/api/couriers", a(func(w http.ResponseWriter, req *http.Request) {
		c, _ := cour.List(req.Context()); apiJSON(w, 200, c)
	}))
	r.GET("/admin/api/credentials", a(func(w http.ResponseWriter, req *http.Request) {
		apiJSON(w, 200, acr.List())
	}))
	r.GET("/admin/api/pricing/routes", a(func(w http.ResponseWriter, req *http.Request) {
		apiJSON(w, 200, rpr.List())
	}))
	r.GET("/admin/api/pricing/delivery", a(func(w http.ResponseWriter, req *http.Request) {
		apiJSON(w, 200, dfr.List())
	}))
	r.GET("/admin/api/pricing/surcharges", a(func(w http.ResponseWriter, req *http.Request) {
		apiJSON(w, 200, scr.List())
	}))

	// ═══ P1: Real data from support repos ═══
	r.GET("/admin/api/print-templates", a(func(w http.ResponseWriter, req *http.Request) {
		items, _ := ppr.List(req.Context(), t); apiJSON(w, 200, items)
	}))
	r.GET("/admin/api/webhooks", a(func(w http.ResponseWriter, req *http.Request) {
		items, _ := whr.ListSubs(req.Context(), t); apiJSON(w, 200, items)
	}))
	r.GET("/admin/api/workflow-management", a(func(w http.ResponseWriter, req *http.Request) {
		items, _ := wfr.ListProcesses(req.Context(), t); apiJSON(w, 200, items)
	}))
	r.GET("/admin/api/task-monitor", a(func(w http.ResponseWriter, req *http.Request) {
		apiJSON(w, 200, td.TaskPool())
	}))

	// ═══ P2: Stub endpoints (ready for real data wiring) ═══
	r.GET("/admin/api/ai-exceptions", empty("AI异常"))
	r.GET("/admin/api/area-groups", empty("区域分组"))
	r.GET("/admin/api/balance-logs", empty("收支明细"))
	r.GET("/admin/api/cargo-types", empty("货物类型"))
	r.GET("/admin/api/client-accounts", empty("客户账号"))
	r.GET("/admin/api/client-ledgers", empty("客户账本"))
	r.GET("/admin/api/client-members", empty("客户成员"))
	r.GET("/admin/api/client-permissions", empty("客户权限"))
	r.GET("/admin/api/client-pricing", empty("客户定价"))
	r.GET("/admin/api/client-recharges", empty("客户充值"))
	r.GET("/admin/api/container-loadings", empty("集装箱装货"))
	r.GET("/admin/api/customer-addresses", empty("收件地址"))
	r.GET("/admin/api/customer-declarants", empty("客户申报人"))
	r.GET("/admin/api/customs-brokers", empty("报关行"))
	r.GET("/admin/api/customs-points", empty("海关口岸"))
	r.GET("/admin/api/exception-reports", empty("异常报告"))
	r.GET("/admin/api/exceptions", empty("异常列表"))
	r.GET("/admin/api/inbound-board", empty("入库看板"))
	r.GET("/admin/api/logistics-tracking", empty("物流追踪"))
	r.GET("/admin/api/monthly-statements", empty("月度对账单"))
	r.GET("/admin/api/notifications", empty("通知管理"))
	r.GET("/admin/api/pda-sessions", empty("PDA会话"))
	r.GET("/admin/api/pda-workorder-templates", empty("PDA工单模板"))
	r.GET("/admin/api/pricing/services", empty("服务计费"))
	r.GET("/admin/api/printers", empty("打印机管理"))
	r.GET("/admin/api/report/client-profit", empty("客户利润报表"))
	r.GET("/admin/api/report/order-profit", empty("订单利润报表"))
	r.GET("/admin/api/report/route-profit", empty("线路利润报表"))
	r.GET("/admin/api/report/service-profit", empty("服务利润报表"))
	r.GET("/admin/api/route-templates", empty("线路模板"))
	r.GET("/admin/api/service-templates", empty("服务模板"))
	r.GET("/admin/api/service-types", empty("服务类型"))
	r.GET("/admin/api/service-workorders", empty("服务工单"))
	r.GET("/admin/api/shipping-providers", empty("承运商管理"))
	r.GET("/admin/api/storage", empty("存储配置"))
	r.GET("/admin/api/system/ai-chat", empty("AI聊天"))
	r.GET("/admin/api/system/ai-settings", empty("AI设置"))
	r.GET("/admin/api/system/api-couriers", empty("快递API"))
	r.GET("/admin/api/system/api-customs", empty("报关API"))
	r.GET("/admin/api/system/api-devices", empty("设备网关"))
	r.GET("/admin/api/system/api-ezway", empty("EZ Way"))
	r.GET("/admin/api/system/api-notifications", empty("通知API"))
	r.GET("/admin/api/system/api-printers", empty("打印API"))
	r.GET("/admin/api/system/api-storage", empty("存储API"))
	r.GET("/admin/api/system/audit-logs", empty("审计日志"))
	r.GET("/admin/api/system/brand", empty("品牌设置"))
	r.GET("/admin/api/system/customs-broker-api", empty("报关经纪API"))
	r.GET("/admin/api/system/logistics-api", empty("物流API"))
	r.GET("/admin/api/system/notification-channels", empty("通知渠道"))
	r.GET("/admin/api/system/params", empty("系统参数"))
	r.GET("/admin/api/system/printers", empty("系统打印机"))
	r.GET("/admin/api/system/reports", empty("系统报表"))
	r.GET("/admin/api/system/scheduler", empty("定时任务"))
	r.GET("/admin/api/system/settings", empty("系统设置"))
	r.GET("/admin/api/transport-modes", empty("运输方式"))
	r.GET("/admin/api/warehouse-board", empty("仓库看板"))
	r.GET("/admin/api/warehouse-console", empty("仓库控制台"))

	_ = orderDomain.Order{}; _ = rbac
}
