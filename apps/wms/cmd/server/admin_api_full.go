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
	psRepo "github.com/i56/modules/parcel_service/repository"
	pricingRepo "github.com/i56/modules/pricing/repository"
	rbacRepo "github.com/i56/modules/rbac/repository"
	tmsRepo "github.com/i56/modules/transport/repository"
	whSvc "github.com/i56/modules/warehouse/service"
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
	rbac *rbacRepo.MemRBACRepo,
) {
	const t int64 = 1

	// ── WAREHOUSES ──
	r.GET("/admin/api/warehouses", a(func(w http.ResponseWriter, req *http.Request) {
		wh, _, _ := ws.List(req.Context(), t, 0, 200)
		apiJSON(w, 200, wh)
	}))
	r.POST("/admin/api/warehouses", a(func(w http.ResponseWriter, req *http.Request) {
		var body struct{ Name, Code, Address, Contact, Phone string }
		json.NewDecoder(req.Body).Decode(&body)
		wh, _ := ws.Create(req.Context(), t, body.Name, body.Code, body.Address, body.Contact, body.Phone)
		apiJSON(w, 201, wh)
	}))

	// ── PARCELS (list + predeclare) ──
	r.GET("/admin/api/parcels", a(func(w http.ResponseWriter, req *http.Request) {
		parcels, _, _ := ps.List(req.Context(), t, 0, 200)
		apiJSON(w, 200, parcels)
	}))
	r.POST("/admin/api/parcels", a(func(w http.ResponseWriter, req *http.Request) {
		var body struct{ TrackingNumber, ProductName, CourierCode string; WarehouseID int64 }
		json.NewDecoder(req.Body).Decode(&body)
		if body.WarehouseID == 0 { body.WarehouseID = 1 }
		p, err := ps.PreDeclare(req.Context(), &parcelDomain.Parcel{
			TrackingNumber: body.TrackingNumber, ProductName: body.ProductName,
			TenantID: t, WarehouseID: body.WarehouseID, CourierCode: body.CourierCode,
		})
		if err != nil { apiJSON(w, 400, map[string]string{"error": err.Error()}); return }
		apiJSON(w, 201, p)
	}))

	// ── ORDERS ──
	r.GET("/admin/api/orders", a(func(w http.ResponseWriter, req *http.Request) {
		orders, _, _ := osvc.List(req.Context(), t, 0, 200)
		apiJSON(w, 200, orders)
	}))
	r.POST("/admin/api/orders", a(func(w http.ResponseWriter, req *http.Request) {
		var body struct{ OrderNo, RecipientName string; ParcelCount int; TotalPrice float64; RouteID int64 }
		json.NewDecoder(req.Body).Decode(&body)
		o := &orderDomain.Order{
			OrderNo: body.OrderNo, RecipientName: body.RecipientName,
			ParcelCount: body.ParcelCount, TotalPrice: body.TotalPrice,
			RouteID: body.RouteID, TenantID: t,
		}
		created, err := osvc.Create(req.Context(), o)
		if err != nil { apiJSON(w, 400, map[string]string{"error": err.Error()}); return }
		apiJSON(w, 201, created)
	}))
	r.GET("/admin/api/orders/{id}", a(func(w http.ResponseWriter, req *http.Request) {
		idStr := req.PathValue("id")
		id, _ := strconv.ParseInt(idStr, 10, 64)
		o, _ := osvc.GetByOrderNo(req.Context(), t, idStr)
		if o == nil { o, _ = osvc.GetByID(req.Context(), t, id) }
		if o == nil { apiJSON(w, 404, map[string]string{"error": "not found"}); return }
		apiJSON(w, 200, o)
	}))

	// ── CLIENTS ──
	r.GET("/admin/api/clients", a(func(w http.ResponseWriter, req *http.Request) {
		cl, _, _ := cr.List(req.Context(), t, 0, 200)
		apiJSON(w, 200, cl)
	}))
	r.POST("/admin/api/clients", a(func(w http.ResponseWriter, req *http.Request) {
		var body struct{ Name, Code string }
		json.NewDecoder(req.Body).Decode(&body)
		c := &custDomain.Client{Name: body.Name, Code: body.Code}
		cr.Create(req.Context(), t, c)
		apiJSON(w, 201, c)
	}))

	// ── EMPLOYEES & ROLES — already registered in admin_api.go ──

	// ── CARRIERS (routes) ──
	r.GET("/admin/api/carriers", a(func(w http.ResponseWriter, req *http.Request) {
		routes, _, _ := rr.List(req.Context(), t, 0, 200)
		apiJSON(w, 200, routes)
	}))

	// ── DECLARANTS ──
	r.GET("/admin/api/declarants", a(func(w http.ResponseWriter, req *http.Request) {
		d, _, _ := dr.List(req.Context(), 1, 0, 200)
		apiJSON(w, 200, d)
	}))

	// ── MEMBERS ──
	r.GET("/admin/api/members", a(func(w http.ResponseWriter, req *http.Request) {
		m, _, _ := mr.List(req.Context(), 1, 0, 200)
		apiJSON(w, 200, m)
	}))

	// ── ADDRESSES ──
	r.GET("/admin/api/addresses", a(func(w http.ResponseWriter, req *http.Request) {
		addr, _ := ar.List(req.Context(), 1)
		apiJSON(w, 200, addr)
	}))

	// ── LEDGER ──
	r.GET("/admin/api/ledger", a(func(w http.ResponseWriter, req *http.Request) {
		entries := lr.GetByClient(req.Context(), 1, 1)
		apiJSON(w, 200, entries)
	}))

	// ── PRICING ──
	r.GET("/admin/api/pricing/routes", a(func(w http.ResponseWriter, req *http.Request) {
		apiJSON(w, 200, rpr.List())
	}))
	r.GET("/admin/api/pricing/delivery", a(func(w http.ResponseWriter, req *http.Request) {
		apiJSON(w, 200, dfr.List())
	}))
	r.GET("/admin/api/pricing/surcharges", a(func(w http.ResponseWriter, req *http.Request) {
		apiJSON(w, 200, scr.List())
	}))

	// ── SERVICE ORDERS ──
	r.GET("/admin/api/service-orders", a(func(w http.ResponseWriter, req *http.Request) {
		so, _, _ := sr.List(req.Context(), t, 0, 200)
		apiJSON(w, 200, so)
	}))

	// ── WORK ORDERS ──
	r.GET("/admin/api/work-orders", a(func(w http.ResponseWriter, req *http.Request) {
		wo, _, _ := wor.List(req.Context(), t, 0, 200)
		apiJSON(w, 200, wo)
	}))

	// ── CREDENTIALS ──
	r.GET("/admin/api/credentials", a(func(w http.ResponseWriter, req *http.Request) {
		apiJSON(w, 200, acr.List())
	}))

	_ = orderDomain.Order{} // suppress unused import
}
