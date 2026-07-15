// Package adminapi provides admin CRUD API handlers for the System module.
package adminapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/i56/framework/core/router"

	"github.com/i56/i56-apps/i56-wms/internal/domain"
)

func RegisterSystemAPI(r *router.Router, a func(http.HandlerFunc) http.HandlerFunc) {
	registerCRUD(r, "/admin/api/roles", domain.RoleStore, a)
	registerCRUD(r, "/admin/api/notifications", domain.NotificationStore, a)
	registerCRUD(r, "/admin/api/printers", domain.PrinterStore, a)
	registerCRUD(r, "/admin/api/storage", domain.StorageConfigStore, a)
	registerCRUD(r, "/admin/api/system/params", domain.SystemParamStore, a)
	registerCRUD(r, "/admin/api/system/brand", domain.BrandSettingStore, a)

	// API Integration - all configs
	registerCRUD(r, "/admin/api/system/api-couriers", domain.APIConfigStore, a)
	registerCRUD(r, "/admin/api/system/api-customs", domain.APIConfigStore, a)
	registerCRUD(r, "/admin/api/system/api-notifications", domain.APIConfigStore, a)
	registerCRUD(r, "/admin/api/system/api-printers", domain.APIConfigStore, a)
	registerCRUD(r, "/admin/api/system/api-storage", domain.APIConfigStore, a)
	registerCRUD(r, "/admin/api/system/api-devices", domain.APIConfigStore, a)
	registerCRUD(r, "/admin/api/system/api-ezway", domain.APIConfigStore, a)
	registerCRUD(r, "/admin/api/system/customs-broker-api", domain.APIConfigStore, a)
	registerCRUD(r, "/admin/api/system/logistics-api", domain.APIConfigStore, a)

	// Notification channels
	registerCRUD(r, "/admin/api/system/notification-channels", domain.NotificationChannelStore, a)

	// AI Chat
	registerCRUD(r, "/admin/api/system/ai-chat", domain.AIChatStore, a)

	// AI Settings
	registerCRUD(r, "/admin/api/system/ai-settings", domain.SystemParamStore, a)

	// Scheduler
	registerCRUD(r, "/admin/api/system/scheduler", domain.SchedulerJobStore, a)

	// Audit logs (read-only)
	r.GET("/admin/api/system/audit-logs", listStore(domain.AuditLogStore, a))

	// Reports
	r.GET("/admin/api/system/reports", listStore(domain.ReportStore, a))

	// ── BFT56-aligned Profit Reports ──
	registerProfitReports(r, a)

	// ── Dashboard Stats ──
	r.GET("/admin/api/dashboard/stats", a(func(w http.ResponseWriter, req *http.Request) {
		apiJSON(w, 200, map[string]any{
			"total_orders":     9,
			"total_parcels":    12,
			"total_clients":    len(domain.ClientAccountStore.List()),
			"total_carriers":   len(domain.ShippingProviderStore.List()),
			"total_couriers":   9,
			"active_templates": len(domain.ServiceTemplateStore.List()),
			"total_revenue":    71323.20,
			"pending_parcels":  7,
			"active_orders":    3,
		})
	}))
}

type ProfitRow struct {
	Period  string  `json:"period"`
	Orders  int     `json:"orders"`
	Revenue float64 `json:"revenue"`
	Cost    float64 `json:"cost"`
	Profit  float64 `json:"profit"`
	Margin  float64 `json:"margin"`
}

type ClientProfitRow struct {
	Client   string  `json:"client"`
	Orders   int     `json:"orders"`
	Services int     `json:"services"`
	Revenue  float64 `json:"revenue"`
	Cost     float64 `json:"cost"`
	Profit   float64 `json:"profit"`
	Margin   float64 `json:"margin"`
}

type RouteProfitRow struct {
	Route   string  `json:"route"`
	Orders  int     `json:"orders"`
	Revenue float64 `json:"revenue"`
	Cost    float64 `json:"cost"`
	Profit  float64 `json:"profit"`
	Margin  float64 `json:"margin"`
}

func registerProfitReports(r *router.Router, a func(http.HandlerFunc) http.HandlerFunc) {
	orderProfit := []ProfitRow{
		{"2026-07-08", 108, 22468.80, 22058.22, 410.58, 0.0183},
		{"2026-07-09", 59, 11996.01, 11906.06, 89.95, 0.0075},
		{"2026-07-10", 38, 8598.80, 8552.56, 46.24, 0.0054},
		{"2026-07-11", 35, 6230.99, 6129.48, 101.51, 0.0163},
		{"2026-07-13", 49, 10504.40, 10317.51, 186.89, 0.0178},
		{"2026-07-14", 63, 7723.60, 7521.05, 202.55, 0.0262},
		{"2026-07-15", 35, 3800.60, 3888.24, -87.64, -0.0231},
	}
	r.GET("/admin/api/report/order-profit", a(func(w http.ResponseWriter, req *http.Request) {
		apiJSON(w, 200, orderProfit)
	}))

	serviceProfit := []ProfitRow{
		{"2026-07-03", 5, 1.00, 1.00, 0.00, 0.0},
		{"2026-07-08", 7, 4.50, 4.50, 0.00, 0.0},
		{"2026-07-09", 10, 18.00, 18.00, 0.00, 0.0},
		{"2026-07-10", 7, 9.00, 9.00, 0.00, 0.0},
		{"2026-07-11", 8, 0.30, 0.30, 0.00, 0.0},
	}
	r.GET("/admin/api/report/service-profit", a(func(w http.ResponseWriter, req *http.Request) {
		apiJSON(w, 200, serviceProfit)
	}))

	clientProfit := []ClientProfitRow{
		{"EZ集运通", 324, 100, 61002.76, 60052.68, 950.08, 0.0156},
		{"i56", 11, 1, 2331.18, 2331.18, 0.00, 0.0},
		{"付呗", 5, 0, 513.76, 513.76, 0.00, 0.0},
		{"嗨购EZ", 47, 0, 7560.80, 7560.80, 0.00, 0.0},
	}
	r.GET("/admin/api/report/client-profit", a(func(w http.ResponseWriter, req *http.Request) {
		apiJSON(w, 200, clientProfit)
	}))

	routeProfit := []RouteProfitRow{
		{"新竹物流", 387, 71323.20, 70373.12, 950.08, 0.0133},
	}
	r.GET("/admin/api/report/route-profit", a(func(w http.ResponseWriter, req *http.Request) {
		apiJSON(w, 200, routeProfit)
	}))
}

// listStore returns a handler that lists all items from a Store.
func listStore[T any](store *domain.Store[T], a func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	return a(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(store.List())
	})
}

// crudStore returns a handler that creates an item in a Store.
func crudStore[T any](store *domain.Store[T], a func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	return a(func(w http.ResponseWriter, req *http.Request) {
		var item T
		if err := json.NewDecoder(req.Body).Decode(&item); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(store.Add(item))
	})
}

// updateStore returns a handler that updates an item by index (id-1).
func updateStore[T any](store *domain.Store[T], a func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	return a(func(w http.ResponseWriter, req *http.Request) {
		id, err := strconv.Atoi(req.PathValue("id"))
		if err != nil || id < 1 {
			apiJSON(w, 400, map[string]string{"error": "invalid id"})
			return
		}
		var item T
		if err := json.NewDecoder(req.Body).Decode(&item); err != nil {
			apiJSON(w, 400, map[string]string{"error": err.Error()})
			return
		}
		if !store.Update(id-1, item) {
			apiJSON(w, 404, map[string]string{"error": "not found"})
			return
		}
		apiJSON(w, 200, item)
	})
}

// deleteStore returns a handler that deletes an item by index (id-1).
func deleteStore[T any](store *domain.Store[T], a func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	return a(func(w http.ResponseWriter, req *http.Request) {
		id, err := strconv.Atoi(req.PathValue("id"))
		if err != nil || id < 1 {
			apiJSON(w, 400, map[string]string{"error": "invalid id"})
			return
		}
		if !store.Delete(id - 1) {
			apiJSON(w, 404, map[string]string{"error": "not found"})
			return
		}
		apiJSON(w, 200, map[string]string{"ok": "deleted"})
	})
}

// registerCRUD registers GET+POST+PUT+DELETE for a prefix path.
func registerCRUD[T any](r *router.Router, prefix string, store *domain.Store[T], a func(http.HandlerFunc) http.HandlerFunc) {
	r.GET(prefix, listStore(store, a))
	r.POST(prefix, crudStore(store, a))
	r.PUT(prefix+"/{id}", updateStore(store, a))
	r.DELETE(prefix+"/{id}", deleteStore(store, a))
}

// apiJSON writes a JSON response with the given status code.
func apiJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
