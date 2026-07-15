// Package adminapi provides admin CRUD API handlers for the System module.
package adminapi

import (
	"encoding/json"
	"net/http"

	"github.com/i56/framework/core/router"

	"github.com/i56/i56-apps/i56-wms/internal/domain"
)

// RegisterSystemAPI registers all System module JSON API endpoints.
// Preserves ALL existing routes from admin_api_full.go.
func RegisterSystemAPI(r *router.Router, a func(http.HandlerFunc) http.HandlerFunc) {
	// Notifications
	r.GET("/admin/api/notifications", listStore(domain.NotificationStore, a))
	r.POST("/admin/api/notifications", crudStore(domain.NotificationStore, a))

	// Printers
	r.GET("/admin/api/printers", listStore(domain.PrinterStore, a))
	r.POST("/admin/api/printers", crudStore(domain.PrinterStore, a))
	r.GET("/admin/api/system/printers", listStore(domain.PrinterStore, a))

	// Storage
	r.GET("/admin/api/storage", listStore(domain.StorageConfigStore, a))
	r.POST("/admin/api/storage", crudStore(domain.StorageConfigStore, a))

	// System params
	r.GET("/admin/api/system/params", listStore(domain.SystemParamStore, a))
	r.POST("/admin/api/system/params", crudStore(domain.SystemParamStore, a))

	// Brand settings
	r.GET("/admin/api/system/brand", listStore(domain.BrandSettingStore, a))
	r.POST("/admin/api/system/brand", crudStore(domain.BrandSettingStore, a))

	// System settings (alias for params)
	r.GET("/admin/api/system/settings", listStore(domain.SystemParamStore, a))

	// API configs — multiple views on the same store
	r.GET("/admin/api/system/api-couriers", listStore(domain.APIConfigStore, a))
	r.POST("/admin/api/system/api-couriers", crudStore(domain.APIConfigStore, a))
	r.GET("/admin/api/system/api-customs", listStore(domain.APIConfigStore, a))
	r.POST("/admin/api/system/api-customs", crudStore(domain.APIConfigStore, a))
	r.GET("/admin/api/system/api-notifications", listStore(domain.APIConfigStore, a))
	r.POST("/admin/api/system/api-notifications", crudStore(domain.APIConfigStore, a))
	r.GET("/admin/api/system/api-printers", listStore(domain.APIConfigStore, a))
	r.POST("/admin/api/system/api-printers", crudStore(domain.APIConfigStore, a))
	r.GET("/admin/api/system/api-storage", listStore(domain.APIConfigStore, a))
	r.POST("/admin/api/system/api-storage", crudStore(domain.APIConfigStore, a))
	r.GET("/admin/api/system/api-devices", listStore(domain.APIConfigStore, a))
	r.POST("/admin/api/system/api-devices", crudStore(domain.APIConfigStore, a))
	r.GET("/admin/api/system/api-ezway", listStore(domain.APIConfigStore, a))
	r.POST("/admin/api/system/api-ezway", crudStore(domain.APIConfigStore, a))
	r.GET("/admin/api/system/customs-broker-api", listStore(domain.APIConfigStore, a))
	r.GET("/admin/api/system/logistics-api", listStore(domain.APIConfigStore, a))

	// Notification channels
	r.GET("/admin/api/system/notification-channels", listStore(domain.NotificationChannelStore, a))
	r.POST("/admin/api/system/notification-channels", crudStore(domain.NotificationChannelStore, a))

	// AI Chat
	r.GET("/admin/api/system/ai-chat", listStore(domain.AIChatStore, a))
	r.POST("/admin/api/system/ai-chat", crudStore(domain.AIChatStore, a))

	// AI Settings (reuse system params)
	r.GET("/admin/api/system/ai-settings", listStore(domain.SystemParamStore, a))

	// Scheduler
	r.GET("/admin/api/system/scheduler", listStore(domain.SchedulerJobStore, a))
	r.POST("/admin/api/system/scheduler", crudStore(domain.SchedulerJobStore, a))

	// Audit logs
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

// apiJSON writes a JSON response with the given status code.
func apiJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
