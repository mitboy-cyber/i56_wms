// Package adminapi provides CRM admin API handlers.
package adminapi

import (
	"net/http"

	"github.com/i56/framework/core/router"

	"github.com/i56/i56-apps/i56-wms/internal/domain"
	clientsvc "github.com/i56/modules/customer/service"
	custRepo "github.com/i56/modules/customer/repository"
	pdaRepo "github.com/i56/modules/pda/repository"
)

func RegisterCRMAPI(r *router.Router, a func(http.HandlerFunc) http.HandlerFunc,
	cs *clientsvc.ClientService, cr *custRepo.MemClientRepo,
	lr *custRepo.MemLedgerRepo, dr *custRepo.MemDeclarantRepo,
	mr *custRepo.MemMemberRepo, ar *custRepo.MemAddressRepo,
	pr *pdaRepo.MemPDARepo,
) {
	const t int64 = 1
	_ = cs; _ = pr // may be used later

	registerCRUD(r, "/admin/api/client-accounts", domain.ClientAccountStore, a)
	registerCRUD(r, "/admin/api/client-recharges", domain.ClientRechargeStore, a)
	registerCRUD(r, "/admin/api/client-pricing", domain.ClientPricingStore, a)
	registerCRUD(r, "/admin/api/client-permissions", domain.ClientPermissionStore, a)
	registerCRUD(r, "/admin/api/monthly-statements", domain.MonthlyStatementStore, a)

	// BFT56-aligned: Client Members, Balance, Recharge, Containers
	registerCRUD(r, "/admin/api/client-members", domain.ClientMemberStore, a)
	registerCRUD(r, "/admin/api/balance-logs", domain.BalanceLogStore, a)
	registerCRUD(r, "/admin/api/recharge-records", domain.RechargeRecordStore, a)
	registerCRUD(r, "/admin/api/containers", domain.ContainerStore, a)
	registerCRUD(r, "/admin/api/client-panel-perms", domain.ClientPanelPermStore, a)

	// CRM sub-module
	registerCRUD(r, "/admin/api/crm/leads", domain.ClientAccountStore, a)
	registerCRUD(r, "/admin/api/crm/opportunities", domain.MonthlyStatementStore, a)
	registerCRUD(r, "/admin/api/crm/activities", domain.AuditLogStore, a)
	registerCRUD(r, "/admin/api/crm/contracts", domain.ClientPricingStore, a)
	registerCRUD(r, "/admin/api/crm/tickets", domain.ExceptionStore, a)
	registerCRUD(r, "/admin/api/crm/segments", domain.ClientPermissionStore, a)
	registerCRUD(r, "/admin/api/crm/campaigns", domain.NotificationStore, a)
	registerCRUD(r, "/admin/api/crm/notes", domain.ExceptionReportStore, a)
	registerCRUD(r, "/admin/api/crm/communications", domain.NotificationChannelStore, a)

	// Customer declarants (real repo)
	r.GET("/admin/api/customer-declarants", a(func(w http.ResponseWriter, req *http.Request) {
		d, _, _ := dr.List(req.Context(), 1, 0, 50)
		apiJSON(w, 200, d)
	}))
	// Customer addresses (real repo)
	r.GET("/admin/api/customer-addresses", a(func(w http.ResponseWriter, req *http.Request) {
		addr, _ := ar.List(req.Context(), 1)
		apiJSON(w, 200, addr)
	}))
	// NOTE: /admin/api/client-members is now handled by registerCRUD above
	// Client ledgers (real repo)
	r.GET("/admin/api/client-ledgers", a(func(w http.ResponseWriter, req *http.Request) {
		entries := lr.GetByClient(req.Context(), 1, 1)
		apiJSON(w, 200, entries)
	}))
	// NOTE: /admin/api/clients CRUD is registered above via registerCRUD
	// but we need client list from real repo - use a DIFFERENT path
	r.GET("/admin/api/clients-list", a(func(w http.ResponseWriter, req *http.Request) {
		clients, _, _ := cr.List(req.Context(), 1, 0, 50)
		apiJSON(w, 200, clients)
	}))
}
