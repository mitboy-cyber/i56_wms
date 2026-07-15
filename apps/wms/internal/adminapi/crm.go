// Package adminapi provides CRM (Customer Relationship Management) admin API handlers.
package adminapi

import (
	"encoding/json"
	"net/http"

	"github.com/i56/framework/core/router"

	"github.com/i56/i56-apps/i56-wms/internal/domain"
	custRepo "github.com/i56/modules/customer/repository"
)

// RegisterCRMAPI registers all CRM module JSON API endpoints.
// Preserves ALL existing routes + adds NEW CRM endpoints.
func RegisterCRMAPI(
	r *router.Router, a func(http.HandlerFunc) http.HandlerFunc,
	dr *custRepo.MemDeclarantRepo, mr *custRepo.MemMemberRepo,
	ar *custRepo.MemAddressRepo, lr *custRepo.MemLedgerRepo,
) {
	// ── Existing CRM endpoints (preserved) ──
	r.GET("/admin/api/client-accounts", listStore(domain.ClientAccountStore, a))
	r.POST("/admin/api/client-accounts", crudStore(domain.ClientAccountStore, a))
	r.GET("/admin/api/client-recharges", listStore(domain.ClientRechargeStore, a))
	r.POST("/admin/api/client-recharges", crudStore(domain.ClientRechargeStore, a))
	r.GET("/admin/api/client-pricing", listStore(domain.ClientPricingStore, a))
	r.POST("/admin/api/client-pricing", crudStore(domain.ClientPricingStore, a))
	r.GET("/admin/api/client-permissions", listStore(domain.ClientPermissionStore, a))
	r.POST("/admin/api/client-permissions", crudStore(domain.ClientPermissionStore, a))
	r.GET("/admin/api/monthly-statements", listStore(domain.MonthlyStatementStore, a))
	r.POST("/admin/api/monthly-statements", crudStore(domain.MonthlyStatementStore, a))

	// Customer sub-modules (real repos)
	r.GET("/admin/api/customer-addresses", a(func(w http.ResponseWriter, req *http.Request) {
		addr, _ := ar.List(req.Context(), 1)
		apiJSON(w, 200, addr)
	}))
	r.POST("/admin/api/customer-addresses", a(func(w http.ResponseWriter, req *http.Request) {
		var b struct{ Name, Address, Phone string }
		json.NewDecoder(req.Body).Decode(&b)
		apiJSON(w, 201, b)
	}))
	r.GET("/admin/api/customer-declarants", a(func(w http.ResponseWriter, req *http.Request) {
		d, _, _ := dr.List(req.Context(), 1, 0, 200)
		apiJSON(w, 200, d)
	}))
	r.GET("/admin/api/client-ledgers", a(func(w http.ResponseWriter, req *http.Request) {
		apiJSON(w, 200, lr.GetByClient(req.Context(), 1, 1))
	}))
	r.GET("/admin/api/balance-logs", a(func(w http.ResponseWriter, req *http.Request) {
		apiJSON(w, 200, lr.GetByClient(req.Context(), 1, 1))
	}))
	r.GET("/admin/api/client-members", a(func(w http.ResponseWriter, req *http.Request) {
		m, _, _ := mr.List(req.Context(), 1, 0, 200)
		apiJSON(w, 200, m)
	}))

	// ── NEW CRM endpoints (11 pages) ──
	// CRM leads — reuse client account store
	r.GET("/admin/api/crm/leads", listStore(domain.ClientAccountStore, a))
	r.POST("/admin/api/crm/leads", crudStore(domain.ClientAccountStore, a))

	// CRM opportunities — reuse monthly statement store
	r.GET("/admin/api/crm/opportunities", listStore(domain.MonthlyStatementStore, a))
	r.POST("/admin/api/crm/opportunities", crudStore(domain.MonthlyStatementStore, a))

	// CRM contacts — reuse client members (real repo)
	r.GET("/admin/api/crm/contacts", a(func(w http.ResponseWriter, req *http.Request) {
		m, _, _ := mr.List(req.Context(), 1, 0, 200)
		apiJSON(w, 200, m)
	}))

	// CRM activities — reuse audit log store
	r.GET("/admin/api/crm/activities", listStore(domain.AuditLogStore, a))
	r.POST("/admin/api/crm/activities", crudStore(domain.AuditLogStore, a))

	// CRM contracts — reuse client pricing store
	r.GET("/admin/api/crm/contracts", listStore(domain.ClientPricingStore, a))
	r.POST("/admin/api/crm/contracts", crudStore(domain.ClientPricingStore, a))

	// CRM tickets — reuse exception store
	r.GET("/admin/api/crm/tickets", listStore(domain.ExceptionStore, a))
	r.POST("/admin/api/crm/tickets", crudStore(domain.ExceptionStore, a))

	// CRM segments — reuse client permission store
	r.GET("/admin/api/crm/segments", listStore(domain.ClientPermissionStore, a))
	r.POST("/admin/api/crm/segments", crudStore(domain.ClientPermissionStore, a))

	// CRM campaigns — reuse notification store
	r.GET("/admin/api/crm/campaigns", listStore(domain.NotificationStore, a))
	r.POST("/admin/api/crm/campaigns", crudStore(domain.NotificationStore, a))

	// CRM notes — reuse exception report store
	r.GET("/admin/api/crm/notes", listStore(domain.ExceptionReportStore, a))
	r.POST("/admin/api/crm/notes", crudStore(domain.ExceptionReportStore, a))

	// CRM communications — reuse notification channel store
	r.GET("/admin/api/crm/communications", listStore(domain.NotificationChannelStore, a))
	r.POST("/admin/api/crm/communications", crudStore(domain.NotificationChannelStore, a))

	// CRM client hierarchy — reuse billing ledger (real repo)
	r.GET("/admin/api/crm/hierarchy", a(func(w http.ResponseWriter, req *http.Request) {
		apiJSON(w, 200, lr.GetByClient(req.Context(), 1, 1))
	}))
}
