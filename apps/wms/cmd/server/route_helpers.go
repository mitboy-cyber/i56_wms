package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/i56/framework/core/auth"
	"github.com/i56/framework/core/router"

	custRepo "github.com/i56/modules/customer/repository"
	orderDomain "github.com/i56/modules/order/domain"
	orderSvc "github.com/i56/modules/order/service"
	parcelDomain "github.com/i56/modules/parcel/domain"
	parcelRepo "github.com/i56/modules/parcel/repository"
	parcelSvc "github.com/i56/modules/parcel/service"
	psRepo "github.com/i56/modules/parcel_service/repository"
	pricingRepo "github.com/i56/modules/pricing/repository"
	tmsRepo "github.com/i56/modules/transport/repository"
	whSvc "github.com/i56/modules/warehouse/service"
	whRepo2 "github.com/i56/modules/webhook/repository"
	weightDomain "github.com/i56/modules/weight/domain"
)

// adminOnly returns an auth middleware that currently passes through all requests.
// TODO: implement proper token validation.
func adminOnly(tm *auth.TokenManager) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) { next(w, r) }
	}
}

// clientPg registers all client portal routes (old-style inline rendering).
// This also calls registerClientP01Routes for the P0/P1 migrated routes.
func clientPg(
	tm *auth.TokenManager,
	cTmpl map[string]*template.Template,
	r *router.Router,
	ps *parcelSvc.ParcelService,
	osvc *orderSvc.OrderService,
	rr *tmsRepo.MemRouteRepo,
	cour *tmsRepo.MemCourierRepo,
	ws *whSvc.WarehouseService,
	pr *parcelRepo.MemParcelRepo,
	lr *custRepo.MemLedgerRepo,
	weightRepo *weightDomain.MemWeightRepo,
	dr *custRepo.MemDeclarantRepo,
	mr *custRepo.MemMemberRepo,
	sr *psRepo.MemServiceRepo,
	whr *whRepo2.MemWebhookRepo,
	ar *custRepo.MemAddressRepo,
	rpr *pricingRepo.MemRoutePriceRepo,
	dfr *pricingRepo.MemDeliveryFeeRepo,
	scr *pricingRepo.MemSurchargeRepo,
	acr *pricingRepo.MemApiCredentialRepo,
) {
	ca := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			ck, err := req.Cookie("client_token")
			if err != nil {
				http.Redirect(w, req, "/client/login", 303)
				return
			}
			if _, err := tm.ValidateAccessToken(ck.Value); err != nil {
				http.Redirect(w, req, "/client/login?error=expired", 303)
				return
			}
			next(w, req)
		}
	}

	r.GET("/client", ca(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		parcels, pt, _ := ps.List(ctx, 1, 0, 8)
		_, rc, _ := rr.List(ctx, 1, 0, 50)
		orders, oc, _ := osvc.List(ctx, 1, 0, 100)
		balance := 0.0
		creditLimit := 5000.0
		if entries := lr.GetByClient(ctx, 1, 1); len(entries) > 0 {
			balance = entries[len(entries)-1].BalanceAfter
		}

		// Build orders for dashboard display
		orderMaps := make([]map[string]any, 0, len(orders))
		activeCount := 0
		statusCN := map[orderDomain.OrderStatus]string{
			orderDomain.StatusPendingPicking:  "待拣货",
			orderDomain.StatusPicking:         "拣货中",
			orderDomain.StatusPendingPacking:  "待打包",
			orderDomain.StatusPendingLoading:  "待装柜",
			orderDomain.StatusLoaded:          "已装柜",
			orderDomain.StatusInTransit:       "运输中",
			orderDomain.StatusCustomsClearance: "清关中",
			orderDomain.StatusOutForDelivery:  "派送中",
			orderDomain.StatusCompleted:       "已完成",
			orderDomain.StatusCancelled:       "已取消",
			orderDomain.StatusShipped:         "已发货",
		}
		for _, o := range orders {
			s := statusCN[o.Status]
			if s == "" {
				s = string(o.Status)
			}
			orderMaps = append(orderMaps, map[string]any{
				"OrderNo":       o.OrderNo,
				"ReceiverName":  o.RecipientName,
				"ParcelCount":   o.ParcelCount,
				"Weight":        fmt.Sprintf("%.2f", o.TotalActualWeight),
				"Amount":        fmt.Sprintf("%.2f", o.TotalPrice),
				"Status":        s,
			})
			if o.Status != orderDomain.StatusCompleted && o.Status != orderDomain.StatusCancelled {
				activeCount++
			}
		}
		// Build parcels for dashboard display
		type parcelDash struct {
			TrackingNumber, ProductName, StatusLabel, StatusColor string
			Weight                                                  float64
		}
		var pdList []parcelDash
		statusLabel := func(s parcelDomain.ParcelStatus) (string, string) {
			switch s {
			case parcelDomain.StatusPreDeclared: return "预报", "secondary"
			case parcelDomain.StatusReceived:    return "已入仓", "info"
			case parcelDomain.StatusWeighed:     return "已称重", "primary"
			case parcelDomain.StatusStored:      return "已上架", "success"
			case parcelDomain.StatusPicked:      return "已拣货", "warning"
			case parcelDomain.StatusPacked:      return "已打包", "success"
			case parcelDomain.StatusLoaded:      return "已装柜", "primary"
			case parcelDomain.StatusOutbound:    return "已出货", "dark"
			default:                             return string(s), "secondary"
			}
		}
		for _, p := range parcels {
			lb, sc := statusLabel(p.Status)
			pdList = append(pdList, parcelDash{p.TrackingNumber, p.ProductName, lb, sc, p.ActualWeight})
		}
		// Parcel status counts
		var preDec, recvd, weighed, stored, picked, packed, shipped int
		allParcels, _, _ := ps.List(ctx, 1, 0, 200)
		for _, p := range allParcels {
			switch p.Status {
			case parcelDomain.StatusPreDeclared: preDec++
			case parcelDomain.StatusReceived:    recvd++
			case parcelDomain.StatusWeighed:     weighed++
			case parcelDomain.StatusStored:      stored++
			case parcelDomain.StatusPicked:      picked++
			case parcelDomain.StatusPacked:      packed++
			case parcelDomain.StatusShipped, parcelDomain.StatusLoaded, parcelDomain.StatusOutbound: shipped++
			}
		}
		// Today's shipment count vs yesterday
		today := time.Now().Truncate(24 * time.Hour)
		todayStr := today.Format("2006-01-02")
		yesterdayStr := today.Add(-24 * time.Hour).Format("2006-01-02")
		todayShipments := 0
		yesterdayShipments := 0
		for _, o := range orders {
			if o.Status == orderDomain.StatusShipped || o.Status == orderDomain.StatusInTransit || o.Status == orderDomain.StatusCompleted {
				od := o.CreatedAt.Format("2006-01-02")
				if od == todayStr {
					todayShipments++
				} else if od == yesterdayStr {
					yesterdayShipments++
				}
			}
		}
		shipmentDelta := todayShipments - yesterdayShipments
		shipmentDeltaPct := 0
		if yesterdayShipments > 0 {
			shipmentDeltaPct = (shipmentDelta * 100) / yesterdayShipments
		}

		// Pre-calculate parcel circle percentages
		totalPc := int(pt)
		type pctData struct {
			Label string
			Count int
			Pct   float64
			Color string
		}
		pctList := []pctData{
			{"预报", preDec, 0, "#6366f1"},
			{"入仓", recvd, 0, "#0ea5e9"},
			{"上架", stored, 0, "#10b981"},
			{"出货", shipped, 0, "#f59e0b"},
		}
		otherCount := totalPc - preDec - recvd - stored - shipped
		if totalPc > 0 {
			for i := range pctList {
				pctList[i].Pct = float64(pctList[i].Count) * 100.0 / float64(totalPc)
			}
		}
		if otherCount > 0 {
			pctList = append(pctList, pctData{"其他", otherCount, float64(otherCount) * 100.0 / float64(totalPc), "#8b5cf6"})
		}

		// Notifications (last 5 from orders)
		type notif struct {
			Time    string
			Icon    string
			Message string
			Type    string
		}
		notifications := []notif{}
		for i, o := range orders {
			if i >= 5 {
				break
			}
			s := statusCN[o.Status]
			if s == "" {
				s = string(o.Status)
			}
			msg := fmt.Sprintf("订单 %s %s", o.OrderNo, s)
			notifications = append(notifications, notif{
				Time:    o.CreatedAt.Format("01-02 15:04"),
				Icon:    "📋",
				Message: msg,
				Type:    "order",
			})
		}

		execTpl(cTmpl, "dashboard", w, "dashboard.html", map[string]any{
			"Title":              "主控台",
			"Balance":            balance,
			"CreditLimit":        creditLimit,
			"AvailableCredit":    creditLimit - balance,
			"TotalParcels":       pt,
			"ParcelCount":        pt,
			"OrderCount":         oc,
			"ActiveOrderCount":   activeCount,
			"Parcels":            pdList,
			"Orders":             orderMaps,
			"RouteCount":         rc,
			"PreDeclaredCount":   preDec,
			"ReceivedCount":      recvd,
			"WeighedCount":       weighed,
			"StoredCount":        stored,
			"PickedCount":        picked,
			"PackingCount":       0,
			"PackedCount":        packed,
			"ShippedCount":       shipped,
			"TodayShipments":     todayShipments,
			"YesterdayShipments": yesterdayShipments,
			"ShipmentDelta":      shipmentDelta,
			"ShipmentDeltaPct":   shipmentDeltaPct,
			"ShipmentDeltaUp":    shipmentDelta >= 0,
			"Notifications":      notifications,
			"PctList":            pctList,
		})
	}))

	r.GET("/client/predeclare", ca(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		warehouses, _, _ := ws.List(ctx, 1, 0, 50)
		allParcels, pt, _ := ps.List(ctx, 1, 0, 200)
		var preDeclared, received, weighed, stored, picked, packed int64
		for _, p := range allParcels {
			switch p.Status {
			case parcelDomain.StatusPreDeclared: preDeclared++
			case parcelDomain.StatusReceived:    received++
			case parcelDomain.StatusWeighed:     weighed++
			case parcelDomain.StatusStored:      stored++
			case parcelDomain.StatusPicked:      picked++
			case parcelDomain.StatusPacked:      packed++
			}
		}
		type member struct{ ID int; Name string }
		type stat struct{ Label string; Count int64 }
		type recent struct{ TN, Name, Status string }
		recentList := []recent{}
		for i, p := range allParcels {
			if i >= 5 {
				break
			}
			recentList = append(recentList, recent{p.TrackingNumber, p.ProductName, string(p.Status)})
		}
		execTpl(cTmpl, "predeclare", w, "predeclare.html", map[string]any{
			"Warehouses": warehouses,
			"Members":    []member{{1, "王仁照"}, {2, "吴欣如"}, {3, "张致廷"}},
			"Stats": []stat{
				{"全部包裹", pt}, {"预报", preDeclared}, {"已入仓", received},
				{"已上架", stored}, {"待打包", picked}, {"打包中", 0}, {"已打包", packed},
			},
			"Recent": recentList,
		})
	}))

	r.GET("/client/parcels", ca(func(w http.ResponseWriter, req *http.Request) {
		parcels, _, _ := ps.List(req.Context(), 1, 0, 100)
		execTpl(cTmpl, "parcels", w, "parcels.html", map[string]any{"Parcels": parcels, "Total": len(parcels)})
	}))

	r.POST("/client/predeclare", ca(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		tn := req.FormValue("tracking_number")
		if tn == "" {
			execTpl(cTmpl, "predeclare", w, "predeclare.html", map[string]any{"Error": "快递单号必填"})
			return
		}
		whID := int64(1)
		if v := req.FormValue("warehouse_id"); v != "" {
			if id, err := strconv.ParseInt(v, 10, 64); err == nil {
				whID = id
			}
		}
		ps.PreDeclare(req.Context(), &parcelDomain.Parcel{
			TrackingNumber: tn,
			ProductName:    req.FormValue("product_name"),
			TenantID:       1,
			WarehouseID:    whID,
			CourierCode:    req.FormValue("courier_code"),
		})
		w.Header().Set("HX-Redirect", "/client/parcels")
		w.WriteHeader(200)
	}))

	r.GET("/client/logout", func(w http.ResponseWriter, req *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "client_token", Value: "", Path: "/client", MaxAge: -1})
		http.Redirect(w, req, "/client/login", 303)
	})

	registerClientP01Routes(r, cTmpl, ps, rr, ca, weightRepo, osvc, lr, ws, cour, dr, mr, sr, whr, ar, rpr, dfr, scr, acr)
}
