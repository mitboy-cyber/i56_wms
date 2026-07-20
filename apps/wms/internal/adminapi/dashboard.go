// Package adminapi provides Dashboard API — warehouse & operational metrics.
package adminapi

import (
	"net/http"

	"github.com/i56/framework/core/router"

	orderRepo "github.com/i56/modules/order/repository"
	parcelRepo "github.com/i56/modules/parcel/repository"
	pdaRepo "github.com/i56/modules/pda/repository"
	whRepo "github.com/i56/modules/warehouse/repository"
)

// RegisterDashboardAPI registers warehouse dashboard endpoints.
func RegisterDashboardAPI(r *router.Router,
	a func(http.HandlerFunc) http.HandlerFunc,
	or *orderRepo.MemOrderRepo,
	wr *whRepo.MemWarehouseRepo,
	pr *parcelRepo.MemParcelRepo,
	pdaR *pdaRepo.MemPDARepo) {

	r.GET("/admin/api/dashboard", a(func(w http.ResponseWriter, req *http.Request) {
		orders, _, _ := or.List(req.Context(), 1, 0, 500)
		warehouses, _, _ := wr.List(req.Context(), 1, 0, 100)
		parcels, _, _ := pr.List(req.Context(), 1, 0, 500)

		var revenue float64
		statusMap := map[string]int{}
		for _, o := range orders {
			revenue += o.TotalPrice
			statusMap[string(o.Status)]++
		}
		statusDist := []map[string]any{}
		for s, c := range statusMap {
			statusDist = append(statusDist, map[string]any{"status": s, "count": c})
		}
		avgOrder := 0.0
		if len(orders) > 0 {
			avgOrder = revenue / float64(len(orders))
		}
		apiJSON(w, 200, map[string]any{
			"total_orders":        len(orders),
			"total_parcels":       len(parcels),
			"total_revenue":       revenue,
			"warehouse_count":     len(warehouses),
			"avg_order_value":     avgOrder,
			"status_distribution": statusDist,
			"pending_tasks":       statusMap["pending_picking"] + statusMap["pending_packing"],
		})
	}))

	r.GET("/admin/api/dashboard/order-status", a(func(w http.ResponseWriter, req *http.Request) {
		orders, _, _ := or.List(req.Context(), 1, 0, 500)
		statusMap := map[string]int{}
		for _, o := range orders {
			statusMap[string(o.Status)]++
		}
		apiJSON(w, 200, statusMap)
	}))

	r.GET("/admin/api/dashboard/revenue-by-route", a(func(w http.ResponseWriter, req *http.Request) {
		orders, _, _ := or.List(req.Context(), 1, 0, 500)
		routeRevenue := map[int64]float64{}
		for _, o := range orders {
			routeRevenue[o.RouteID] += o.TotalPrice
		}
		apiJSON(w, 200, routeRevenue)
	}))

	// PDA monitoring
	r.GET("/admin/api/pda/online-sessions", a(func(w http.ResponseWriter, req *http.Request) {
		count := pdaR.ActiveSessionCount()
		type opInfo struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
			Code string `json:"code"`
		}
		operators := pdaR.ListOperators()
		var result []opInfo
		for _, op := range operators {
			result = append(result, opInfo{ID: op.ID, Name: op.Name, Code: op.Code})
		}
		apiJSON(w, 200, map[string]any{"active_sessions": count, "operators": result})
	}))

	r.GET("/admin/api/pda/scan-logs", a(func(w http.ResponseWriter, req *http.Request) {
		logs := pdaR.RecentScans(50)
		type logInfo struct {
			Action  string `json:"action"`
			Barcode string `json:"barcode"`
			Success bool   `json:"success"`
		}
		var result []logInfo
		for _, l := range logs {
			result = append(result, logInfo{Action: l.Action, Barcode: l.Barcode, Success: l.Success})
		}
		apiJSON(w, 200, map[string]any{"total": len(result), "logs": result})
	}))
}
