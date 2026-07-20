// Package adminapi provides Finance admin API handlers.
package adminapi

import (
	"context"
	"net/http"

	"github.com/i56/framework/core/router"

	custRepo "github.com/i56/modules/customer/repository"
	orderRepo "github.com/i56/modules/order/repository"
	tmsRepo "github.com/i56/modules/transport/repository"

	"github.com/i56/i56-apps/i56-wms/internal/domain"
)

// RegisterFinanceAPI registers all Finance module JSON API endpoints.
func RegisterFinanceAPI(r *router.Router, a func(http.HandlerFunc) http.HandlerFunc,
	or *orderRepo.MemOrderRepo, lr *custRepo.MemLedgerRepo, rr *tmsRepo.MemRouteRepo) {

	// ── Existing reports ──
	r.GET("/admin/api/finance/revenue-report", a(func(w http.ResponseWriter, req *http.Request) {
		apiJSON(w, 200, map[string]any{"report": "revenue", "data": domain.MonthlyStatementStore.List()})
	}))
	r.GET("/admin/api/finance/cost-report", a(func(w http.ResponseWriter, req *http.Request) {
		apiJSON(w, 200, map[string]any{"report": "cost", "data": domain.ClientRechargeStore.List()})
	}))
	r.GET("/admin/api/finance/profit-loss", a(func(w http.ResponseWriter, req *http.Request) {
		stmts := domain.MonthlyStatementStore.List()
		var rev, paid float64
		for _, s := range stmts { rev += s.Total; paid += s.PaidAmount }
		apiJSON(w, 200, map[string]any{"report": "profit_loss", "total_revenue": rev, "total_paid": paid, "outstanding": rev - paid})
	}))
	r.GET("/admin/api/finance/cash-flow", a(func(w http.ResponseWriter, req *http.Request) {
		recharges := domain.ClientRechargeStore.List()
		var inflow float64
		for _, r := range recharges { inflow += r.Amount }
		apiJSON(w, 200, map[string]any{"report": "cash_flow", "total_inflow": inflow, "net_cash": inflow})
	}))

	// ── NEW: Real aggregation from repos (BFT56-aligned 4 dimensions) ──

	// Order profitability
	r.GET("/admin/api/finance/order-profit", a(func(w http.ResponseWriter, req *http.Request) {
		orders, _, _ := or.List(req.Context(), 1, 0, 200)
		var revenue float64
		for _, o := range orders { revenue += o.TotalPrice }
		apiJSON(w, 200, map[string]any{"report": "order_profit", "orders": len(orders), "total_revenue": revenue,
			"avg_order": func() float64 { if len(orders) == 0 { return 0 }; return revenue / float64(len(orders)) }()})
	}))

	// Customer balance
	r.GET("/admin/api/finance/customer-balance", a(func(w http.ResponseWriter, req *http.Request) {
		entries, _, _ := lr.List(context.Background(), 1, 0, 0, 200)
		balances := map[int64]float64{}
		for _, e := range entries { balances[e.ClientID] += e.Amount }
		type cb struct { ClientID int64 `json:"client_id"`; Balance float64 `json:"balance"` }
		var out []cb
		for id, bal := range balances { out = append(out, cb{id, bal}) }
		apiJSON(w, 200, map[string]any{"report": "customer_balance", "customers": out})
	}))

	// Aliases for frontend compatibility
	// /client-profit → customer-balance alias
	r.GET("/admin/api/finance/client-profit", a(func(w http.ResponseWriter, req *http.Request) {
		routes, _, _ := rr.List(req.Context(), 1, 0, 100)
		orders, _, _ := or.List(req.Context(), 1, 0, 200)
		var rev float64
		for _, o := range orders { rev += o.TotalPrice }
		apiJSON(w, 200, map[string]any{"report": "client_profit", "total_revenue": rev, "clients": len(routes), "avg_per_client": func() float64 {
			if len(routes) == 0 { return 0 }; return rev / float64(len(routes))
		}()})
	}))

	// /income-statement → profit-loss alias
	r.GET("/admin/api/finance/income-statement", a(func(w http.ResponseWriter, req *http.Request) {
		stmts := domain.MonthlyStatementStore.List()
		var rev, paid float64
		for _, s := range stmts { rev += s.Total; paid += s.PaidAmount }
		apiJSON(w, 200, map[string]any{"report": "income_statement", "total_revenue": rev, "total_paid": paid, "outstanding": rev - paid})
	}))

	// Fix route-profit: add total_revenue
	r.GET("/admin/api/finance/route-profit", a(func(w http.ResponseWriter, req *http.Request) {
		routes, _, _ := rr.List(req.Context(), 1, 0, 100)
		orders, _, _ := or.List(req.Context(), 1, 0, 200)
		var rev float64
		routeRev := map[int64]float64{}
		for _, o := range orders {
			rev += o.TotalPrice
			routeRev[o.RouteID] += o.TotalPrice
		}
		type rt struct { Name string `json:"name"`; Price float64 `json:"base_price"`; Type string `json:"type"`; Revenue float64 `json:"revenue"` }
		var out []rt
		for _, r := range routes {
			out = append(out, rt{r.Name, r.BaseWeightPrice, r.TransportType, routeRev[r.ID]})
		}
		apiJSON(w, 200, map[string]any{"report": "route_profit", "total_revenue": rev, "routes": out})
	}))
}
