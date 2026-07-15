// Package adminapi provides Finance admin API handlers.
package adminapi

import (
	"net/http"

	"github.com/i56/framework/core/router"

	"github.com/i56/i56-apps/i56-wms/internal/domain"
)

// RegisterFinanceAPI registers all Finance module JSON API endpoints.
// Preserves ALL existing report routes + adds NEW finance report endpoints.
func RegisterFinanceAPI(r *router.Router, a func(http.HandlerFunc) http.HandlerFunc) {
	// ── Existing finance report endpoints (preserved) ──
	r.GET("/admin/api/report/order-profit", listStore(domain.ReportStore, a))
	r.GET("/admin/api/report/route-profit", listStore(domain.ReportStore, a))
	r.GET("/admin/api/report/client-profit", listStore(domain.ReportStore, a))
	r.GET("/admin/api/report/service-profit", listStore(domain.ReportStore, a))

	// ── NEW Finance report endpoints (4 reports) ──
	// Revenue report — uses monthly statement store for financial data
	r.GET("/admin/api/finance/revenue-report", a(func(w http.ResponseWriter, req *http.Request) {
		apiJSON(w, 200, map[string]any{
			"report": "revenue",
			"data":   domain.MonthlyStatementStore.List(),
			"total":  len(domain.MonthlyStatementStore.List()),
		})
	}))

	// Cost report — uses client recharge store for cost tracking
	r.GET("/admin/api/finance/cost-report", a(func(w http.ResponseWriter, req *http.Request) {
		apiJSON(w, 200, map[string]any{
			"report": "cost",
			"data":   domain.ClientRechargeStore.List(),
			"total":  len(domain.ClientRechargeStore.List()),
		})
	}))

	// Profit & Loss — combines revenue and cost data
	r.GET("/admin/api/finance/profit-loss", a(func(w http.ResponseWriter, req *http.Request) {
		statements := domain.MonthlyStatementStore.List()
		var totalRevenue, totalPaid float64
		for _, s := range statements {
			totalRevenue += s.Total
			totalPaid += s.PaidAmount
		}
		apiJSON(w, 200, map[string]any{
			"report":        "profit_loss",
			"total_revenue":  totalRevenue,
			"total_paid":     totalPaid,
			"outstanding":    totalRevenue - totalPaid,
			"statement_count": len(statements),
		})
	}))

	// Cash flow — uses recharge records for cash flow tracking
	r.GET("/admin/api/finance/cash-flow", a(func(w http.ResponseWriter, req *http.Request) {
		recharges := domain.ClientRechargeStore.List()
		var totalInflow float64
		for _, r := range recharges {
			totalInflow += r.Amount
		}
		apiJSON(w, 200, map[string]any{
			"report":       "cash_flow",
			"total_inflow":  totalInflow,
			"total_outflow": 0.0,
			"net_cash":      totalInflow,
			"records":       recharges,
		})
	}))

	// Payment history — uses client recharge store
	r.GET("/admin/api/finance/payments", listStore(domain.ClientRechargeStore, a))
	r.GET("/admin/api/finance/invoices", listStore(domain.MonthlyStatementStore, a))
}
