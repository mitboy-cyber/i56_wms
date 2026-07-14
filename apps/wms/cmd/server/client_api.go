package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	fwAuth "github.com/i56/framework/core/auth"
	"github.com/i56/framework/core/router"

	custRepo "github.com/i56/modules/customer/repository"
	orderSvc "github.com/i56/modules/order/service"
	parcelDomain "github.com/i56/modules/parcel/domain"
	parcelSvc "github.com/i56/modules/parcel/service"
	psRepo "github.com/i56/modules/parcel_service/repository"
	pricingRepo "github.com/i56/modules/pricing/repository"
	tmsRepo "github.com/i56/modules/transport/repository"
	whSvc "github.com/i56/modules/warehouse/service"
	whRepo2 "github.com/i56/modules/webhook/repository"
)

func registerClientJSONAPI(
	r *router.Router, tm *fwAuth.TokenManager,
	ps *parcelSvc.ParcelService, osvc *orderSvc.OrderService,
	rr *tmsRepo.MemRouteRepo, cour *tmsRepo.MemCourierRepo,
	ws *whSvc.WarehouseService, lr *custRepo.MemLedgerRepo,
	dr *custRepo.MemDeclarantRepo, mr *custRepo.MemMemberRepo,
	sr *psRepo.MemServiceRepo, whr *whRepo2.MemWebhookRepo,
	ar *custRepo.MemAddressRepo, rpr *pricingRepo.MemRoutePriceRepo,
	dfr *pricingRepo.MemDeliveryFeeRepo, scr *pricingRepo.MemSurchargeRepo,
	acr *pricingRepo.MemApiCredentialRepo,
) {
	const t int64 = 1

	ca := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			ck, err := req.Cookie("client_token")
			if err != nil {
				apiJSON(w, 401, map[string]string{"error": "unauthorized"})
				return
			}
			if _, err := tm.ValidateAccessToken(ck.Value); err != nil {
				apiJSON(w, 401, map[string]string{"error": "invalid_token"})
				return
			}
			next(w, req)
		}
	}

	r.GET("/client/api/me", ca(func(w http.ResponseWriter, req *http.Request) {
		apiJSON(w, 200, map[string]string{"client": "EZ集運通", "tenant": "1"})
	}))

	r.GET("/client/api/dashboard", ca(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		_, pt, _ := ps.List(ctx, t, 0, 8)
		_, rc, _ := rr.List(ctx, t, 0, 50)
		_, oc, _ := osvc.List(ctx, t, 0, 100)
		balance := 0.0
		if entries := lr.GetByClient(ctx, 1, 1); len(entries) > 0 {
			balance = entries[len(entries)-1].BalanceAfter
		}
		apiJSON(w, 200, map[string]interface{}{
			"balance": balance, "total_parcels": pt,
			"order_count": oc, "route_count": rc,
		})
	}))

	r.GET("/client/api/parcels", ca(func(w http.ResponseWriter, req *http.Request) {
		parcels, _, _ := ps.List(req.Context(), t, 0, 100)
		apiJSON(w, 200, parcels)
	}))

	r.POST("/client/api/parcels/predeclare", ca(func(w http.ResponseWriter, req *http.Request) {
		var body struct {
			TrackingNumber string `json:"tracking_number"`
			ProductName    string `json:"product_name"`
			WarehouseID    int64  `json:"warehouse_id"`
			CourierCode    string `json:"courier_code"`
		}
		json.NewDecoder(req.Body).Decode(&body)
		if body.TrackingNumber == "" {
			apiJSON(w, 400, map[string]string{"error": "tracking_number required"})
			return
		}
		whID := body.WarehouseID
		if whID == 0 { whID = 1 }
		p, err := ps.PreDeclare(req.Context(), &parcelDomain.Parcel{
			TrackingNumber: body.TrackingNumber, ProductName: body.ProductName,
			TenantID: t, WarehouseID: whID, CourierCode: body.CourierCode,
		})
		if err != nil {
			apiJSON(w, 500, map[string]string{"error": err.Error()})
			return
		}
		apiJSON(w, 201, p)
	}))

	r.GET("/client/api/orders", ca(func(w http.ResponseWriter, req *http.Request) {
		orders, _, _ := osvc.List(req.Context(), t, 0, 50)
		apiJSON(w, 200, orders)
	}))

	r.GET("/client/api/orders/{id}", ca(func(w http.ResponseWriter, req *http.Request) {
		idStr := req.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			o, oErr := osvc.GetByOrderNo(req.Context(), t, idStr)
			if oErr != nil { apiJSON(w, 404, map[string]string{"error": "not found"}); return }
			apiJSON(w, 200, o)
			return
		}
		o, err := osvc.GetByID(req.Context(), t, id)
		if err != nil || o == nil { apiJSON(w, 404, map[string]string{"error": "not found"}); return }
		apiJSON(w, 200, o)
	}))

	r.GET("/client/api/ledger", ca(func(w http.ResponseWriter, req *http.Request) {
		entries := lr.GetByClient(req.Context(), 1, 1)
		apiJSON(w, 200, entries)
	}))

	r.GET("/client/api/declarants", ca(func(w http.ResponseWriter, req *http.Request) {
		d, _, _ := dr.List(req.Context(), 1, 0, 50)
		apiJSON(w, 200, d)
	}))

	r.GET("/client/api/members", ca(func(w http.ResponseWriter, req *http.Request) {
		m, _, _ := mr.List(req.Context(), 1, 0, 50)
		apiJSON(w, 200, m)
	}))

	r.GET("/client/api/addresses", ca(func(w http.ResponseWriter, req *http.Request) {
		a, _ := ar.List(req.Context(), 1)
		apiJSON(w, 200, a)
	}))

	r.GET("/client/api/warehouses", ca(func(w http.ResponseWriter, req *http.Request) {
		wh, _, _ := ws.List(req.Context(), t, 0, 50)
		apiJSON(w, 200, wh)
	}))

	r.GET("/client/api/route-prices", ca(func(w http.ResponseWriter, req *http.Request) {
		rp := rpr.List()
		apiJSON(w, 200, rp)
	}))

	r.GET("/client/api/service-orders", ca(func(w http.ResponseWriter, req *http.Request) {
		so, _, _ := sr.List(req.Context(), t, 0, 50)
		apiJSON(w, 200, so)
	}))

	r.GET("/client/api/webhooks", ca(func(w http.ResponseWriter, req *http.Request) {
		wh, _ := whr.ListSubs(req.Context(), t)
		apiJSON(w, 200, wh)
	}))

	r.GET("/client/api/credentials", ca(func(w http.ResponseWriter, req *http.Request) {
		creds := acr.List()
		apiJSON(w, 200, creds)
	}))

	r.GET("/client/api/delivery-fees", ca(func(w http.ResponseWriter, req *http.Request) {
		df := dfr.List()
		apiJSON(w, 200, df)
	}))

	r.GET("/client/api/surcharges", ca(func(w http.ResponseWriter, req *http.Request) {
		sc := scr.List()
		apiJSON(w, 200, sc)
	}))

	r.GET("/client/api/couriers", ca(func(w http.ResponseWriter, req *http.Request) {
		c, _ := cour.List(req.Context())
		apiJSON(w, 200, c)
	}))
}
