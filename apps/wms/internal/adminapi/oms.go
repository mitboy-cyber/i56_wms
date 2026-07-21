// Package adminapi provides OMS (Order Management) admin API handlers.
package adminapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/i56/framework/core/router"
	"github.com/i56/framework/core/eventbus"
	"github.com/i56/framework/core/tenant"

	"github.com/i56/i56-apps/i56-wms/internal/httputil"
	"github.com/i56/i56-apps/i56-wms/internal/types"
	"github.com/i56/i56-apps/i56-wms/internal/validate"
	custDomain "github.com/i56/modules/customer/domain"
	custRepo "github.com/i56/modules/customer/repository"
	orderDomain "github.com/i56/modules/order/domain"
	orderSvc "github.com/i56/modules/order/service"
	parcelDomain "github.com/i56/modules/parcel/domain"
	parcelSvc "github.com/i56/modules/parcel/service"
	psRepo "github.com/i56/modules/parcel_service/repository"
	pricingRepo "github.com/i56/modules/pricing/repository"
	tmsRepo "github.com/i56/modules/transport/repository"
	whSvc "github.com/i56/modules/warehouse/service"
)

// RegisterOMSAPI registers all OMS module JSON API endpoints.
// Handles orders, parcels, warehouses, services — real data from module repos.
func RegisterOMSAPI(
	r *router.Router, a func(http.HandlerFunc) http.HandlerFunc,
	ps *parcelSvc.ParcelService, osvc *orderSvc.OrderService,
	ws *whSvc.WarehouseService, cr *custRepo.MemClientRepo,
	rr *tmsRepo.MemRouteRepo, cour *tmsRepo.MemCourierRepo,
	sr *psRepo.MemServiceRepo, lr *custRepo.MemLedgerRepo,
	dr *custRepo.MemDeclarantRepo, mr *custRepo.MemMemberRepo,
	ar *custRepo.MemAddressRepo, rpr *pricingRepo.MemRoutePriceRepo,
	dfr *pricingRepo.MemDeliveryFeeRepo, scr *pricingRepo.MemSurchargeRepo,
	acr *pricingRepo.MemApiCredentialRepo,
	eb *eventbus.EventBus,
) {
	// Extract tenant from context (fallback to 1 for backward compatibility)
	tenantID := func(req *http.Request) int64 {
		if ti := tenant.FromContext(req.Context()); ti != nil {
			switch ti.ID {
			case "default": return 1
			case "t2": return 2
			default: return 1
			}
		}
		return 1
	}

	// Warehouses
	r.GET("/admin/api/warehouses", a(func(w http.ResponseWriter, req *http.Request) {
		wh, _, _ := ws.List(req.Context(), tenantID(req), 0, 200)
		apiJSON(w, 200, wh)
	}))
	r.POST("/admin/api/warehouses", a(func(w http.ResponseWriter, req *http.Request) {
		var b types.CreateWarehouseRequest
		if err := json.NewDecoder(req.Body).Decode(&b); err != nil {
			httputil.BadRequest(w, "请求数据格式错误")
			return
		}
		if errs := validate.Struct(&b); errs != nil {
			httputil.ValidationError(w, errs)
			return
		}
		wh, err := ws.Create(req.Context(), tenantID(req), b.Name, b.Code, b.Address, b.Contact, b.Phone)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Created(w, wh)
	}))

	// Parcels
	r.GET("/admin/api/parcels", a(func(w http.ResponseWriter, req *http.Request) {
		px, _, _ := ps.List(req.Context(), tenantID(req), 0, 200)
		apiJSON(w, 200, px)
	}))
	r.POST("/admin/api/parcels", a(func(w http.ResponseWriter, req *http.Request) {
		var b types.CreateParcelRequest
		if err := json.NewDecoder(req.Body).Decode(&b); err != nil {
			httputil.BadRequest(w, "请求数据格式错误")
			return
		}
		if errs := validate.Struct(&b); errs != nil {
			httputil.ValidationError(w, errs)
			return
		}
		if b.WarehouseID == 0 {
			b.WarehouseID = 1
		}
		p, err := ps.PreDeclare(req.Context(), &parcelDomain.Parcel{
			TrackingNumber: b.TrackingNumber, ProductName: b.ProductName,
			TenantID: tenantID(req), WarehouseID: b.WarehouseID, CourierCode: b.CourierCode,
			ActualWeight: b.ActualWeight, CargoType: b.CargoType,
		})
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Created(w, p)
	}))

	// Orders
	r.GET("/admin/api/orders", a(func(w http.ResponseWriter, req *http.Request) {
		ox, _, _ := osvc.List(req.Context(), tenantID(req), 0, 200)
		apiJSON(w, 200, ox)
	}))
	r.POST("/admin/api/orders", a(func(w http.ResponseWriter, req *http.Request) {
		var b types.CreateOrderRequest
		if err := json.NewDecoder(req.Body).Decode(&b); err != nil {
			httputil.BadRequest(w, "请求数据格式错误")
			return
		}
		if errs := validate.Struct(&b); errs != nil {
			httputil.ValidationError(w, errs)
			return
		}
		o := &orderDomain.Order{
			RecipientName: b.RecipientName, ParcelCount: b.ParcelCount,
			TotalPrice: b.TotalPrice, RouteID: b.RouteID, TenantID: tenantID(req),
			TrackingNumbers: b.TrackingNumbers, Remark: b.Remark,
		}
		created, err := osvc.Create(req.Context(), o)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		eb.Publish(req.Context(), &DataEvent{
			BaseEvent: eventbus.NewEvent("order.created"),
			Data: map[string]interface{}{
				"order_id": created.ID, "order_no": created.OrderNo,
				"recipient": created.RecipientName, "amount": created.TotalPrice,
				"parcel_count": created.ParcelCount,
			},
		})
		httputil.Created(w, created)
	}))
	r.GET("/admin/api/orders/{id}", a(func(w http.ResponseWriter, req *http.Request) {
		idStr := req.PathValue("id")
		id, _ := strconv.ParseInt(idStr, 10, 64)
		o, _ := osvc.GetByOrderNo(req.Context(), tenantID(req), idStr)
		if o == nil {
			o, _ = osvc.GetByID(req.Context(), tenantID(req), id)
		}
		if o == nil {
			apiJSON(w, 404, map[string]string{"error": "not found"})
			return
		}
		apiJSON(w, 200, o)
	}))
	// Order status transition
	r.PUT("/admin/api/orders/{id}/status", a(func(w http.ResponseWriter, req *http.Request) {
		id, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
		if err != nil {
			httputil.BadRequest(w, "无效的订单ID")
			return
		}
		var b struct{ Status string `json:"status"` }
		if err := json.NewDecoder(req.Body).Decode(&b); err != nil {
			httputil.BadRequest(w, "请求数据格式错误")
			return
		}
		if b.Status == "" {
			httputil.BadRequest(w, "状态不能为空")
			return
		}
		if err := osvc.Transition(req.Context(), tenantID(req), id, orderDomain.OrderStatus(b.Status)); err != nil {
			httputil.InternalError(w, err)
			return
		}
		o, _ := osvc.GetByID(req.Context(), tenantID(req), id)
		httputil.OK(w, o)
	}))

	// Clients (real repos)
	r.GET("/admin/api/clients", a(func(w http.ResponseWriter, req *http.Request) {
		cl, _, _ := cr.List(req.Context(), tenantID(req), 0, 200)
		apiJSON(w, 200, cl)
	}))
	r.POST("/admin/api/clients", a(func(w http.ResponseWriter, req *http.Request) {
		var b types.CreateClientRequest
		if err := json.NewDecoder(req.Body).Decode(&b); err != nil {
			httputil.BadRequest(w, "请求数据格式错误")
			return
		}
		if errs := validate.Struct(&b); errs != nil {
			httputil.ValidationError(w, errs)
			return
		}
		c := &custDomain.Client{Name: b.Name, Code: b.Code, ContactName: b.Contact, ContactPhone: b.Phone}
		if err := cr.Create(req.Context(), tenantID(req), c); err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Created(w, c)
	}))

	// Sub-resources
	r.GET("/admin/api/declarants", a(func(w http.ResponseWriter, req *http.Request) {
		d, _, _ := dr.List(req.Context(), 1, 0, 200)
		apiJSON(w, 200, d)
	}))
	r.GET("/admin/api/members", a(func(w http.ResponseWriter, req *http.Request) {
		m, _, _ := mr.List(req.Context(), 1, 0, 200)
		apiJSON(w, 200, m)
	}))
	r.GET("/admin/api/addresses", a(func(w http.ResponseWriter, req *http.Request) {
		addr, _ := ar.List(req.Context(), 1)
		apiJSON(w, 200, addr)
	}))
	r.GET("/admin/api/ledger", a(func(w http.ResponseWriter, req *http.Request) {
		apiJSON(w, 200, lr.GetByClient(req.Context(), 1, 1))
	}))
	r.GET("/admin/api/service-orders", a(func(w http.ResponseWriter, req *http.Request) {
		so, _, _ := sr.List(req.Context(), tenantID(req), 0, 200)
		apiJSON(w, 200, so)
	}))

	// Transport
	r.GET("/admin/api/carriers", a(func(w http.ResponseWriter, req *http.Request) {
		routes, _, _ := rr.List(req.Context(), tenantID(req), 0, 200)
		apiJSON(w, 200, routes)
	}))
	r.GET("/admin/api/couriers", a(func(w http.ResponseWriter, req *http.Request) {
		c, _ := cour.List(req.Context())
		apiJSON(w, 200, c)
	}))

	// Pricing
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
}

// dataEvent wraps BaseEvent with payload for domain events.
type DataEvent struct {
	eventbus.BaseEvent
	Data map[string]interface{}
}
