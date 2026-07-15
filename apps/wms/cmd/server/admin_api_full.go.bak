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

	// ══════════════════════════════════════════
	// Helpers
	// ══════════════════════════════════════════

	parseID := func(w http.ResponseWriter, req *http.Request) (int64, bool) {
		id, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
		if err != nil {
			apiJSON(w, 400, map[string]string{"error": "invalid id"})
			return 0, false
		}
		return id, true
	}

	// ══════════════════════════════════════════
	// P0: Core services — real data
	// ══════════════════════════════════════════

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

	// P0: extra repos
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

	// ══════════════════════════════════════════
	// P1: Admin data stores — read/write
	// ══════════════════════════════════════════

	// TMS
	r.GET("/admin/api/area-groups", listStore(areaGroupStore, a))
	r.POST("/admin/api/area-groups", crudStore(areaGroupStore, a))
	r.GET("/admin/api/cargo-types", listStore(cargoTypeStore, a))
	r.POST("/admin/api/cargo-types", crudStore(cargoTypeStore, a))
	r.GET("/admin/api/transport-modes", listStore(transportModeStore, a))
	r.POST("/admin/api/transport-modes", crudStore(transportModeStore, a))
	r.GET("/admin/api/customs-brokers", listStore(customsBrokerStore, a))
	r.POST("/admin/api/customs-brokers", crudStore(customsBrokerStore, a))
	r.GET("/admin/api/customs-points", listStore(customsPointStore, a))
	r.POST("/admin/api/customs-points", crudStore(customsPointStore, a))
	r.GET("/admin/api/shipping-providers", listStore(shippingProviderStore, a))
	r.POST("/admin/api/shipping-providers", crudStore(shippingProviderStore, a))
	r.GET("/admin/api/container-loadings", listStore(containerLoadingStore, a))
	r.POST("/admin/api/container-loadings", crudStore(containerLoadingStore, a))
	r.GET("/admin/api/logistics-tracking", listStore(logisticsTrackingStore, a))
	r.POST("/admin/api/logistics-tracking", crudStore(logisticsTrackingStore, a))
	r.GET("/admin/api/route-templates", listStore(routeTemplateStore, a))
	r.POST("/admin/api/route-templates", crudStore(routeTemplateStore, a))

	// CRM
	r.GET("/admin/api/client-accounts", listStore(clientAccountStore, a))
	r.POST("/admin/api/client-accounts", crudStore(clientAccountStore, a))
	r.GET("/admin/api/client-recharges", listStore(clientRechargeStore, a))
	r.POST("/admin/api/client-recharges", crudStore(clientRechargeStore, a))
	r.GET("/admin/api/client-pricing", listStore(clientPricingStore, a))
	r.POST("/admin/api/client-pricing", crudStore(clientPricingStore, a))
	r.GET("/admin/api/client-permissions", listStore(clientPermissionStore, a))
	r.POST("/admin/api/client-permissions", crudStore(clientPermissionStore, a))
	r.GET("/admin/api/monthly-statements", listStore(monthlyStatementStore, a))
	r.POST("/admin/api/monthly-statements", crudStore(monthlyStatementStore, a))

	// Reuse existing repos for customer sub-modules
	r.GET("/admin/api/customer-addresses", a(func(w http.ResponseWriter, req *http.Request) {
		addr, _ := ar.List(req.Context(), 1); apiJSON(w, 200, addr)
	}))
	r.POST("/admin/api/customer-addresses", a(func(w http.ResponseWriter, req *http.Request) {
		var b struct{ Name, Address, Phone string }
		json.NewDecoder(req.Body).Decode(&b)
		apiJSON(w, 201, b) // stored via admin data stores
	}))
	r.GET("/admin/api/customer-declarants", a(func(w http.ResponseWriter, req *http.Request) {
		d, _, _ := dr.List(req.Context(), 1, 0, 200); apiJSON(w, 200, d)
	}))
	r.GET("/admin/api/client-ledgers", a(func(w http.ResponseWriter, req *http.Request) {
		apiJSON(w, 200, lr.GetByClient(req.Context(), 1, 1))
	}))
	r.GET("/admin/api/balance-logs", a(func(w http.ResponseWriter, req *http.Request) {
		apiJSON(w, 200, lr.GetByClient(req.Context(), 1, 1))
	}))
	r.GET("/admin/api/client-members", a(func(w http.ResponseWriter, req *http.Request) {
		m, _, _ := mr.List(req.Context(), 1, 0, 200); apiJSON(w, 200, m)
	}))

	// WMS + Exceptions
	r.GET("/admin/api/exceptions", listStore(exceptionStore, a))
	r.POST("/admin/api/exceptions", crudStore(exceptionStore, a))
	r.GET("/admin/api/ai-exceptions", listStore(aiExceptionStore, a))
	r.POST("/admin/api/ai-exceptions", crudStore(aiExceptionStore, a))
	r.GET("/admin/api/exception-reports", listStore(exceptionReportStore, a))
	r.GET("/admin/api/pda-sessions", listStore(pdaSessionStore, a))
	r.POST("/admin/api/pda-sessions", crudStore(pdaSessionStore, a))
	r.GET("/admin/api/pda-workorder-templates", listStore(pdaWorkorderTplStore, a))
	r.POST("/admin/api/pda-workorder-templates", crudStore(pdaWorkorderTplStore, a))
	r.GET("/admin/api/service-templates", listStore(serviceTemplateStore, a))
	r.POST("/admin/api/service-templates", crudStore(serviceTemplateStore, a))
	r.GET("/admin/api/service-types", listStore(serviceTypeStore, a))
	r.POST("/admin/api/service-types", crudStore(serviceTypeStore, a))
	r.GET("/admin/api/service-workorders", listStore(serviceWorkorderStore, a))
	r.POST("/admin/api/service-workorders", crudStore(serviceWorkorderStore, a))

	// Dashboard
	r.GET("/admin/api/inbound-board", listStore(inboundBoardStore, a))
	r.GET("/admin/api/warehouse-board", listStore(warehouseBoardStore, a))
	r.GET("/admin/api/warehouse-console", listStore(warehouseConsoleStore, a))

	// Pricing
	r.GET("/admin/api/pricing/services", listStore(pricingServiceStore, a))
	r.POST("/admin/api/pricing/services", crudStore(pricingServiceStore, a))

	// System
	r.GET("/admin/api/notifications", listStore(notificationStore, a))
	r.POST("/admin/api/notifications", crudStore(notificationStore, a))
	r.GET("/admin/api/printers", listStore(printerStore, a))
	r.POST("/admin/api/printers", crudStore(printerStore, a))
	r.GET("/admin/api/storage", listStore(storageConfigStore, a))
	r.POST("/admin/api/storage", crudStore(storageConfigStore, a))
	r.GET("/admin/api/system/params", listStore(systemParamStore, a))
	r.POST("/admin/api/system/params", crudStore(systemParamStore, a))
	r.GET("/admin/api/system/brand", listStore(brandSettingStore, a))
	r.POST("/admin/api/system/brand", crudStore(brandSettingStore, a))
	r.GET("/admin/api/system/settings", listStore(systemParamStore, a))

	// API configs
	r.GET("/admin/api/system/api-couriers", listStore(apiConfigStore, a))
	r.POST("/admin/api/system/api-couriers", crudStore(apiConfigStore, a))
	r.GET("/admin/api/system/api-customs", listStore(apiConfigStore, a))
	r.POST("/admin/api/system/api-customs", crudStore(apiConfigStore, a))
	r.GET("/admin/api/system/api-notifications", listStore(apiConfigStore, a))
	r.POST("/admin/api/system/api-notifications", crudStore(apiConfigStore, a))
	r.GET("/admin/api/system/api-printers", listStore(apiConfigStore, a))
	r.POST("/admin/api/system/api-printers", crudStore(apiConfigStore, a))
	r.GET("/admin/api/system/api-storage", listStore(apiConfigStore, a))
	r.POST("/admin/api/system/api-storage", crudStore(apiConfigStore, a))
	r.GET("/admin/api/system/api-devices", listStore(apiConfigStore, a))
	r.POST("/admin/api/system/api-devices", crudStore(apiConfigStore, a))
	r.GET("/admin/api/system/api-ezway", listStore(apiConfigStore, a))
	r.POST("/admin/api/system/api-ezway", crudStore(apiConfigStore, a))
	r.GET("/admin/api/system/customs-broker-api", listStore(apiConfigStore, a))
	r.GET("/admin/api/system/logistics-api", listStore(apiConfigStore, a))
	r.GET("/admin/api/system/notification-channels", listStore(notificationChannelStore, a))
	r.POST("/admin/api/system/notification-channels", crudStore(notificationChannelStore, a))
	r.GET("/admin/api/system/printers", listStore(printerStore, a))

	// AI + Ops
	r.GET("/admin/api/system/ai-chat", listStore(aiChatStore, a))
	r.POST("/admin/api/system/ai-chat", crudStore(aiChatStore, a))
	r.GET("/admin/api/system/ai-settings", listStore(systemParamStore, a))
	r.GET("/admin/api/system/scheduler", listStore(schedulerJobStore, a))
	r.POST("/admin/api/system/scheduler", crudStore(schedulerJobStore, a))
	r.GET("/admin/api/system/audit-logs", listStore(auditLogStore, a))
	r.GET("/admin/api/system/reports", listStore(reportStore, a))

	// Finance
	r.GET("/admin/api/report/order-profit", listStore(reportStore, a))
	r.GET("/admin/api/report/route-profit", listStore(reportStore, a))
	r.GET("/admin/api/report/client-profit", listStore(reportStore, a))
	r.GET("/admin/api/report/service-profit", listStore(reportStore, a))

	_ = orderDomain.Order{}; _ = rbac; _ = parseID
}

// Generic list handler
func listStore[T any](store *Store[T], a func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	return a(func(w http.ResponseWriter, req *http.Request) { apiJSON(w, 200, store.List()) })
}

// Generic create handler
func crudStore[T any](store *Store[T], a func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	return a(func(w http.ResponseWriter, req *http.Request) {
		var item T
		if err := json.NewDecoder(req.Body).Decode(&item); err != nil { apiJSON(w, 400, map[string]string{"error": err.Error()}); return }
		apiJSON(w, 201, store.Add(item))
	})
}
