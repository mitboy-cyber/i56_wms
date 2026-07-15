// Package adminapi provides TMS (Transport Management) admin API handlers.
package adminapi

import (
	"net/http"

	"github.com/i56/framework/core/router"

	"github.com/i56/i56-apps/i56-wms/internal/domain"
)

// RegisterTMSAPI registers all TMS module JSON API endpoints.
// Preserves ALL existing routes + adds NEW TMS endpoints.
func RegisterTMSAPI(r *router.Router, a func(http.HandlerFunc) http.HandlerFunc) {
	// ── Existing TMS endpoints (preserved) ──
	r.GET("/admin/api/area-groups", listStore(domain.AreaGroupStore, a))
	r.POST("/admin/api/area-groups", crudStore(domain.AreaGroupStore, a))
	r.GET("/admin/api/cargo-types", listStore(domain.CargoTypeStore, a))
	r.POST("/admin/api/cargo-types", crudStore(domain.CargoTypeStore, a))
	r.GET("/admin/api/transport-modes", listStore(domain.TransportModeStore, a))
	r.POST("/admin/api/transport-modes", crudStore(domain.TransportModeStore, a))
	r.GET("/admin/api/customs-brokers", listStore(domain.CustomsBrokerStore, a))
	r.POST("/admin/api/customs-brokers", crudStore(domain.CustomsBrokerStore, a))
	r.GET("/admin/api/customs-points", listStore(domain.CustomsPointStore, a))
	r.POST("/admin/api/customs-points", crudStore(domain.CustomsPointStore, a))
	r.GET("/admin/api/shipping-providers", listStore(domain.ShippingProviderStore, a))
	r.POST("/admin/api/shipping-providers", crudStore(domain.ShippingProviderStore, a))
	r.GET("/admin/api/container-loadings", listStore(domain.ContainerLoadingStore, a))
	r.POST("/admin/api/container-loadings", crudStore(domain.ContainerLoadingStore, a))
	r.GET("/admin/api/logistics-tracking", listStore(domain.LogisticsTrackingStore, a))
	r.POST("/admin/api/logistics-tracking", crudStore(domain.LogisticsTrackingStore, a))
	r.GET("/admin/api/route-templates", listStore(domain.RouteTemplateStore, a))
	r.POST("/admin/api/route-templates", crudStore(domain.RouteTemplateStore, a))

	// ── NEW TMS endpoints (10 pages) ──
	// TMS routes — reuse routeTemplateStore
	r.GET("/admin/api/tms/routes", listStore(domain.RouteTemplateStore, a))
	r.POST("/admin/api/tms/routes", crudStore(domain.RouteTemplateStore, a))

	// TMS containers — reuse containerLoadingStore
	r.GET("/admin/api/tms/containers", listStore(domain.ContainerLoadingStore, a))
	r.POST("/admin/api/tms/containers", crudStore(domain.ContainerLoadingStore, a))

	// TMS shipments — reuse logisticsTrackingStore
	r.GET("/admin/api/tms/shipments", listStore(domain.LogisticsTrackingStore, a))
	r.POST("/admin/api/tms/shipments", crudStore(domain.LogisticsTrackingStore, a))

	// TMS drivers (new domain concept, reuse shipping providers as base)
	r.GET("/admin/api/tms/drivers", listStore(domain.ShippingProviderStore, a))
	r.POST("/admin/api/tms/drivers", crudStore(domain.ShippingProviderStore, a))

	// TMS vehicles — reuse transport mode store
	r.GET("/admin/api/tms/vehicles", listStore(domain.TransportModeStore, a))
	r.POST("/admin/api/tms/vehicles", crudStore(domain.TransportModeStore, a))

	// TMS depots — reuse customs point store
	r.GET("/admin/api/tms/depots", listStore(domain.CustomsPointStore, a))
	r.POST("/admin/api/tms/depots", crudStore(domain.CustomsPointStore, a))

	// TMS tracking events — reuse logisticsTrackingStore
	r.GET("/admin/api/tms/tracking-events", listStore(domain.LogisticsTrackingStore, a))

	// TMS carrier contracts — reuse customs broker store
	r.GET("/admin/api/tms/carrier-contracts", listStore(domain.CustomsBrokerStore, a))
	r.POST("/admin/api/tms/carrier-contracts", crudStore(domain.CustomsBrokerStore, a))

	// TMS shipping schedules — reuse containerLoadingStore
	r.GET("/admin/api/tms/shipping-schedules", listStore(domain.ContainerLoadingStore, a))
	r.POST("/admin/api/tms/shipping-schedules", crudStore(domain.ContainerLoadingStore, a))

	// TMS rate cards — new endpoint using route pricing pattern
	r.GET("/admin/api/tms/rate-cards", listStore(domain.RouteTemplateStore, a))
	r.POST("/admin/api/tms/rate-cards", crudStore(domain.RouteTemplateStore, a))
}
