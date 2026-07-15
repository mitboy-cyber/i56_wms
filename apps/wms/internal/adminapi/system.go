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
