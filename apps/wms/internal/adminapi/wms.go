// Package adminapi provides WMS admin API handlers.
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

func RegisterWMSAPI(
	r *router.Router, a func(http.HandlerFunc) http.HandlerFunc,
	ppr *printRepo.MemPrintRepo, wfr *wfRepo.MemWorkflowRepo,
	td *taskdispatch.MemTaskDispatchRepo, whr *whRepo2.MemWebhookRepo,
	wor *twoRepo.MemWorkOrderRepo,
) {
	const t int64 = 1

	registerCRUD(r, "/admin/api/exceptions", domain.ExceptionStore, a)
	registerCRUD(r, "/admin/api/ai-exceptions", domain.AIExceptionStore, a)
	registerCRUD(r, "/admin/api/pda-sessions", domain.PDASessionStore, a)
	registerCRUD(r, "/admin/api/pda-workorder-templates", domain.PDAWorkorderTplStore, a)
	registerCRUD(r, "/admin/api/service-templates", domain.ServiceTemplateStore, a)
	registerCRUD(r, "/admin/api/service-types", domain.ServiceTypeStore, a)
	registerCRUD(r, "/admin/api/service-workorders", domain.ServiceWorkorderStore, a)
	registerCRUD(r, "/admin/api/pricing/services", domain.PricingServiceStore, a)

	// Read-only dashboard boards
	r.GET("/admin/api/exception-reports", listStore(domain.ExceptionReportStore, a))
	r.GET("/admin/api/inbound-board", listStore(domain.InboundBoardStore, a))
	r.GET("/admin/api/warehouse-board", listStore(domain.WarehouseBoardStore, a))
	r.GET("/admin/api/warehouse-console", listStore(domain.WarehouseConsoleStore, a))

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
	// Workflow management
	r.GET("/admin/api/workflow-management", a(func(w http.ResponseWriter, req *http.Request) {
		items := domain.WorkflowProcessStore.List()
		apiJSON(w, 200, items)
	}))
	// Task monitor
	r.GET("/admin/api/task-monitor", a(func(w http.ResponseWriter, req *http.Request) {
		apiJSON(w, 200, td.TaskPool())
	}))
}
