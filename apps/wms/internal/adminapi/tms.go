// Package adminapi provides TMS admin API handlers.
package adminapi

import (
	"net/http"

	"github.com/i56/framework/core/router"

	"github.com/i56/i56-apps/i56-wms/internal/domain"
)

func RegisterTMSAPI(r *router.Router, a func(http.HandlerFunc) http.HandlerFunc) {
	registerCRUD(r, "/admin/api/area-groups", domain.AreaGroupStore, a)
	registerCRUD(r, "/admin/api/cargo-types", domain.CargoTypeStore, a)
	registerCRUD(r, "/admin/api/transport-modes", domain.TransportModeStore, a)
	registerCRUD(r, "/admin/api/customs-brokers", domain.CustomsBrokerStore, a)
	registerCRUD(r, "/admin/api/customs-points", domain.CustomsPointStore, a)
	registerCRUD(r, "/admin/api/shipping-providers", domain.ShippingProviderStore, a)
	registerCRUD(r, "/admin/api/container-loadings", domain.ContainerLoadingStore, a)
	registerCRUD(r, "/admin/api/logistics-tracking", domain.LogisticsTrackingStore, a)
	registerCRUD(r, "/admin/api/route-templates", domain.RouteTemplateStore, a)

	// TMS module sub-routes
	registerCRUD(r, "/admin/api/tms/routes", domain.RouteTemplateStore, a)
	registerCRUD(r, "/admin/api/tms/containers", domain.ContainerLoadingStore, a)
	registerCRUD(r, "/admin/api/tms/shipments", domain.LogisticsTrackingStore, a)
	registerCRUD(r, "/admin/api/tms/drivers", domain.ShippingProviderStore, a)
	registerCRUD(r, "/admin/api/tms/vehicles", domain.TransportModeStore, a)
	registerCRUD(r, "/admin/api/tms/depots", domain.CustomsPointStore, a)
	registerCRUD(r, "/admin/api/tms/carrier-contracts", domain.CustomsBrokerStore, a)
	registerCRUD(r, "/admin/api/tms/shipping-schedules", domain.ContainerLoadingStore, a)
	registerCRUD(r, "/admin/api/tms/rate-cards", domain.RouteTemplateStore, a)
}
