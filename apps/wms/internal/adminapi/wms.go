// Package adminapi provides WMS (Warehouse Management) admin API handlers.
package adminapi

import (
	"net/http"

	"github.com/i56/framework/core/router"

	"github.com/i56/i56-apps/i56-wms/internal/domain"
	printRepo "github.com/i56/modules/print/repository"
	taskdispatch "github.com/i56/modules/taskdispatch/repository"
	whRepo2 "github.com/i56/modules/webhook/repository"
	wfRepo "github.com/i56/modules/workflow/repository"
	twoRepo "github.com/i56/modules/workorder/repository"
)

// RegisterWMSAPI registers all WMS module JSON API endpoints.
func RegisterWMSAPI(
	r *router.Router, a func(http.HandlerFunc) http.HandlerFunc,
	ppr *printRepo.MemPrintRepo, wfr *wfRepo.MemWorkflowRepo,
	td *taskdispatch.MemTaskDispatchRepo, whr *whRepo2.MemWebhookRepo,
	wor *twoRepo.MemWorkOrderRepo,
) {
	const t int64 = 1

	// Exceptions
	r.GET("/admin/api/exceptions", listStore(domain.ExceptionStore, a))
	r.POST("/admin/api/exceptions", crudStore(domain.ExceptionStore, a))

	// AI Exceptions
	r.GET("/admin/api/ai-exceptions", listStore(domain.AIExceptionStore, a))
	r.POST("/admin/api/ai-exceptions", crudStore(domain.AIExceptionStore, a))

	// Exception Reports
	r.GET("/admin/api/exception-reports", listStore(domain.ExceptionReportStore, a))

	// PDA Sessions
	r.GET("/admin/api/pda-sessions", listStore(domain.PDASessionStore, a))
	r.POST("/admin/api/pda-sessions", crudStore(domain.PDASessionStore, a))

	// PDA Workorder Templates
	r.GET("/admin/api/pda-workorder-templates", listStore(domain.PDAWorkorderTplStore, a))
	r.POST("/admin/api/pda-workorder-templates", crudStore(domain.PDAWorkorderTplStore, a))

	// Service Templates
	r.GET("/admin/api/service-templates", listStore(domain.ServiceTemplateStore, a))
	r.POST("/admin/api/service-templates", crudStore(domain.ServiceTemplateStore, a))

	// Service Types
	r.GET("/admin/api/service-types", listStore(domain.ServiceTypeStore, a))
	r.POST("/admin/api/service-types", crudStore(domain.ServiceTypeStore, a))

	// Service Workorders
	r.GET("/admin/api/service-workorders", listStore(domain.ServiceWorkorderStore, a))
	r.POST("/admin/api/service-workorders", crudStore(domain.ServiceWorkorderStore, a))

	// Dashboard boards
	r.GET("/admin/api/inbound-board", listStore(domain.InboundBoardStore, a))
	r.GET("/admin/api/warehouse-board", listStore(domain.WarehouseBoardStore, a))
	r.GET("/admin/api/warehouse-console", listStore(domain.WarehouseConsoleStore, a))

	// Pricing services
	r.GET("/admin/api/pricing/services", listStore(domain.PricingServiceStore, a))
	r.POST("/admin/api/pricing/services", crudStore(domain.PricingServiceStore, a))

	// Work orders (real repo)
	r.GET("/admin/api/work-orders", a(func(w http.ResponseWriter, req *http.Request) {
		wo, _, _ := wor.List(req.Context(), t, 0, 200)
		apiJSON(w, 200, wo)
	}))

	// Print templates (real repo)
	r.GET("/admin/api/print-templates", a(func(w http.ResponseWriter, req *http.Request) {
		items, _ := ppr.List(req.Context(), t)
		apiJSON(w, 200, items)
	}))

	// Webhooks (real repo)
	r.GET("/admin/api/webhooks", a(func(w http.ResponseWriter, req *http.Request) {
		items, _ := whr.ListSubs(req.Context(), t)
		apiJSON(w, 200, items)
	}))

	// Workflow management (real repo)
	r.GET("/admin/api/workflow-management", a(func(w http.ResponseWriter, req *http.Request) {
		items := domain.WorkflowProcessStore.List()
		apiJSON(w, 200, items)
	}))

	// Task monitor (real repo)
	r.GET("/admin/api/task-monitor", a(func(w http.ResponseWriter, req *http.Request) {
		apiJSON(w, 200, td.TaskPool())
	}))
}
