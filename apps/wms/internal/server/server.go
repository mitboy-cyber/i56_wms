// Package server provides HTTP server setup, route wiring, and seed data.
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"text/template"
	"time"

	"github.com/i56/framework/core/auth"
	"github.com/i56/framework/core/config"
	"github.com/i56/framework/core/middleware"
	"github.com/i56/framework/core/router"
	"github.com/i56/framework/core/scheduler"
	jpkg "github.com/i56/framework/core/jwt"
	"github.com/i56/framework/core/sse"
	"github.com/i56/framework/core/eventbus"
	"github.com/i56/framework/core/logger"
	"github.com/i56/framework/core/tenant"
	"github.com/i56/framework/core/rbac"
	"github.com/i56/framework/core/storage"
	"github.com/i56/framework/core/workflow"
	"github.com/i56/framework/core/notification"
	"github.com/i56/framework/core/plugin"
	"github.com/i56/framework/db"

	adminAuth "github.com/i56/i56-apps/i56-wms/internal/auth"
	"github.com/i56/i56-apps/i56-wms/internal/adminapi"
	"github.com/i56/i56-apps/i56-wms/internal/clientapi"
	"github.com/i56/i56-apps/i56-wms/internal/domain"
	"github.com/i56/i56-apps/i56-wms/internal/pdaapi"
	wmsMiddleware "github.com/i56/i56-apps/i56-wms/internal/middleware"

	custRepo "github.com/i56/modules/customer/repository"
	custDomain "github.com/i56/modules/customer/domain"
	custSvc "github.com/i56/modules/customer/service"
	orderRepo "github.com/i56/modules/order/repository"
	orderDomain "github.com/i56/modules/order/domain"
	orderSvc "github.com/i56/modules/order/service"
	parcelRepo "github.com/i56/modules/parcel/repository"
	parcelDomain "github.com/i56/modules/parcel/domain"
	parcelSvc "github.com/i56/modules/parcel/service"
	pdaRepo "github.com/i56/modules/pda/repository"
	pdaSvc "github.com/i56/modules/pda/service"
	pricingRepo "github.com/i56/modules/pricing/repository"
	psRepo "github.com/i56/modules/parcel_service/repository"
	printRepo "github.com/i56/modules/print/repository"
	rbacRepoPkg "github.com/i56/modules/rbac/repository"
	sysRepo "github.com/i56/modules/system/repository"
	taskdispatch "github.com/i56/modules/taskdispatch/repository"
	tmsRepo "github.com/i56/modules/transport/repository"
	tmsDomain "github.com/i56/modules/transport/domain"
	whRepo "github.com/i56/modules/warehouse/repository"
	whDomain "github.com/i56/modules/warehouse/domain"
	whSvc "github.com/i56/modules/warehouse/service"
	whRepo2 "github.com/i56/modules/webhook/repository"
	wfRepo "github.com/i56/modules/workflow/repository"
	twoRepo "github.com/i56/modules/workorder/repository"
	wmsRepo "github.com/i56/modules/wms/repository"
	wmsDomain "github.com/i56/modules/wms/domain"
	wdDomain "github.com/i56/modules/workorder/domain"
)

const I56Version = "2.4.2"

// Server encapsulates the HTTP server and all its dependencies.
type Server struct {
	cfg         Config
	router      *router.Router
	httpSrv     *http.Server
	dbAvailable bool
	dbConn     interface{} // *sql.DB handle for graceful shutdown

	SessionMgr *adminAuth.SessionManager
	TokenMgr   *auth.TokenManager
	EventBus   *eventbus.EventBus
	EventLog   []map[string]interface{} // in-memory event log
	evtMu      sync.RWMutex             // guards EventLog
	TenantMgr  *tenant.InMemTenantStore
	PermStore  *rbac.InMemPermissionStore
	RBAC       *rbac.Enforcer
	Storage    storage.StorageProvider
	Workflow   *workflow.Engine
	NotifySvc  *notification.Service
	NotifChan  *memNotifChannel
	Registry   *plugin.Registry // unified service locator

	// Repos (singletons)
	ClientRepo        *custRepo.MemClientRepo
	WarehouseRepo     *whRepo.MemWarehouseRepo
	ParcelRepo        *parcelRepo.MemParcelRepo
	OrderRepo         *orderRepo.MemOrderRepo
	CourierRepo       *tmsRepo.MemCourierRepo
	RouteRepo         *tmsRepo.MemRouteRepo
	LedgerRepo        *custRepo.MemLedgerRepo
	ServiceRepo       *psRepo.MemServiceRepo
	WorkOrderRepo     *twoRepo.MemWorkOrderRepo
	PrintRepo         *printRepo.MemPrintRepo
	WebhookRepo       *whRepo2.MemWebhookRepo
	RBACRepo          *rbacRepoPkg.MemRBACRepo
	WorkflowRepo      *wfRepo.MemWorkflowRepo
	TaskDispatchRepo  *taskdispatch.MemTaskDispatchRepo
	DeclarantRepo     *custRepo.MemDeclarantRepo
	MemberRepo        *custRepo.MemMemberRepo
	AddressRepo       *custRepo.MemAddressRepo
	RoutePriceRepo    *pricingRepo.MemRoutePriceRepo
	DeliveryFeeRepo   *pricingRepo.MemDeliveryFeeRepo
	SurchargeRepo     *pricingRepo.MemSurchargeRepo
	APICredentialRepo *pricingRepo.MemApiCredentialRepo
	PdaRepo           *pdaRepo.MemPDARepo
	WMSRepo           *wmsRepo.MemWMSRepo

	ParcelSvc    *parcelSvc.ParcelService
	OrderSvc     *orderSvc.OrderService
	ClientSvc    *custSvc.ClientService
	WarehouseSvc *whSvc.WarehouseService
	PDAOps       *pdaSvc.PDAOperations

	Scheduler    *scheduler.Scheduler
	AuditLogger  interface{}
	SSEHub       *sse.Hub
	TemplateMap  map[string]*template.Template
	ClientTplMap map[string]*template.Template
}

// Config holds server configuration.
type Config struct {
	Port      int
	DBDSN     string
	StaticDir string
}

// New creates a new Server with all dependencies wired.
func New(cfg Config) (*Server, error) {
	s := &Server{cfg: cfg}

	// JWT
	_, err := jpkg.NewService("i56-framework")
	if err != nil {
		return nil, fmt.Errorf("jwt: %w", err)
	}

	sysCfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	tm, err := auth.NewTokenManager(sysCfg.Auth)
	if err != nil {
		return nil, fmt.Errorf("token manager: %w", err)
	}
	s.TokenMgr = tm
	s.SessionMgr = adminAuth.NewSessionManager()

	// Event Bus
	s.EventBus = eventbus.New(nil)
	s.EventLog = make([]map[string]interface{}, 0, 200)
	s.registerEventHandlers()

	// Tenant
	s.TenantMgr = tenant.NewInMemTenantStore()
	s.TenantMgr.AddTenant(&tenant.TenantInfo{ID: "default", Name: "默认租户", Strategy: tenant.StrategyShared})
	s.TenantMgr.AddTenant(&tenant.TenantInfo{ID: "t2", Name: "测试租户2", Strategy: tenant.StrategyShared})

	// RBAC
	s.PermStore = rbac.NewInMemPermissionStore()
	s.PermStore.AddRole("admin", map[string][]string{"*": {"*"}})
	s.PermStore.AddRole("operator", map[string][]string{
		"parcel": {"read", "create", "update"}, "order": {"read", "create"},
		"warehouse": {"read"}, "report": {"read"},
	})
	s.PermStore.AddRole("viewer", map[string][]string{
		"parcel": {"read"}, "order": {"read"}, "report": {"read"},
	})
	s.PermStore.AssignRole("admin", "admin")
	s.PermStore.AssignRole("OP001", "operator")
	s.PermStore.SetDataScope("parcel", rbac.ScopeWarehouse)
	s.PermStore.SetDataScope("order", rbac.ScopeTenant)
	s.RBAC = rbac.NewEnforcer(s.PermStore)

	// Storage
	st, err := storage.NewLocalStorage("/opt/i56/storage")
	if err != nil {
		st, err = storage.NewLocalStorage("/tmp/i56-storage")
	}
	if err != nil || st == nil {
		return nil, fmt.Errorf("storage init failed: %w", err)
	}
	s.Storage = st

	// Workflow Engine
	wfStore := &memWorkflowStore{instances: make(map[string]*workflow.ProcessInstance)}
	// Use a noop logger to prevent nil-pointer SEGV during process transitions
	s.Workflow = workflow.NewEngine(wfStore, &noopLogger{})
	// Register purchase approval workflow
	s.Workflow.RegisterDefinition(&workflow.ProcessDefinition{
		ID:   "purchase-approval",
		Name: "采购审批流程",
		States: []workflow.State{
			{ID: "start", Name: "开始", Type: workflow.StateStart},
			{ID: "dept_approve", Name: "部门审批", Type: workflow.StateTask},
			{ID: "finance_approve", Name: "财务审批", Type: workflow.StateTask},
			{ID: "gm_approve", Name: "总经理审批", Type: workflow.StateTask},
			{ID: "end", Name: "完成", Type: workflow.StateEnd},
		},
		Transitions: []workflow.Transition{
			{From: "start", To: "dept_approve"},
			{From: "dept_approve", To: "finance_approve", Condition: "amount>5000"},
			{From: "dept_approve", To: "end", Condition: "amount<=5000"},
			{From: "finance_approve", To: "gm_approve", Condition: "amount>50000"},
			{From: "finance_approve", To: "end", Condition: "amount<=50000"},
			{From: "gm_approve", To: "end"},
		},
	})

	// Notification Center
	memChan := &memNotifChannel{}
	s.NotifChan = memChan
	s.NotifySvc = notification.NewService(&noopLogger{})
	s.NotifySvc.Register(memChan)

	// Plugin Registry — unified service locator for all framework modules
	s.Registry = plugin.NewRegistry(&noopLogger{})
	s.Registry.Provide("eventbus", s.EventBus)
	s.Registry.Provide("tenant", s.TenantMgr)
	s.Registry.Provide("rbac", s.RBAC)
	s.Registry.Provide("storage", s.Storage)
	s.Registry.Provide("workflow", s.Workflow)
	s.Registry.Provide("notification", s.NotifySvc)

	// PostgreSQL
	if err := db.Connect(cfg.DBDSN); err == nil {
		s.dbAvailable = true
		s.dbConn = db.Pool // keep handle for later close
		log.Println("[DB] Using PostgreSQL")
	}

	// Init repos (singletons)
	s.ClientRepo = custRepo.NewMemClientRepo()
	s.WarehouseRepo = whRepo.NewMemWarehouseRepo()
	s.ParcelRepo = parcelRepo.NewMemParcelRepo()
	s.OrderRepo = orderRepo.NewMemOrderRepo()
	s.CourierRepo = tmsRepo.NewMemCourierRepo()
	s.RouteRepo = tmsRepo.NewMemRouteRepo()
	s.LedgerRepo = custRepo.NewMemLedgerRepo()
	s.ServiceRepo = psRepo.NewMemServiceRepo()
	s.WorkOrderRepo = twoRepo.NewMemWorkOrderRepo()
	s.PrintRepo = printRepo.NewMemPrintRepo()
	s.WebhookRepo = whRepo2.NewMemWebhookRepo()
	s.RBACRepo = rbacRepoPkg.NewMemRBACRepo()
	s.WorkflowRepo = wfRepo.NewMemWorkflowRepo()
	s.TaskDispatchRepo = taskdispatch.NewMemTaskDispatchRepo()
	s.DeclarantRepo = custRepo.NewMemDeclarantRepo()
	s.MemberRepo = custRepo.NewMemMemberRepo()
	s.AddressRepo = custRepo.NewMemAddressRepo()
	s.RoutePriceRepo = pricingRepo.NewMemRoutePriceRepo()
	s.DeliveryFeeRepo = pricingRepo.NewMemDeliveryFeeRepo()
	s.SurchargeRepo = pricingRepo.NewMemSurchargeRepo()
	s.APICredentialRepo = pricingRepo.NewMemApiCredentialRepo()
	s.PdaRepo = pdaRepo.NewMemPDARepo()
	_ = sysRepo.NewMemSystemConfigRepo()

	domain.SeedAll()
	seedRealData(s)

	// Services
	s.ParcelSvc = parcelSvc.NewParcelService(s.ParcelRepo)
	s.OrderSvc = orderSvc.NewOrderService(s.OrderRepo)
	s.ClientSvc = custSvc.NewClientService(s.ClientRepo)
	s.WarehouseSvc = whSvc.NewWarehouseService(s.WarehouseRepo)

	// WMS repo for PDA operations (must init before PDAOps)
	s.WMSRepo = wmsRepo.NewMemWMSRepo()
	s.WMSRepo.CreateLocation(context.Background(), &wmsDomain.Location{
		ZoneID: 1, Code: "A-01-01", Barcode: "A-01-01",
	})

	s.PDAOps = pdaSvc.NewPDAOperations(s.ParcelRepo, s.OrderRepo, s.WMSRepo, nil)

	// Framework
	s.SSEHub = sse.NewHub()
	s.Scheduler = scheduler.New()
	scheduler.DemoJobs(s.Scheduler)
	s.Scheduler.Start()

	// Framework services
	s.TemplateMap = initAdminTemplates()
	s.ClientTplMap = initClientTemplates()

	// Router
	s.router = router.New()
	// Tenant resolver: try header, fall back to "default"
	tr := tenant.NewMultiResolver(
		tenant.NewHeaderResolver("X-Tenant-ID", s.TenantMgr),
	)
	tenMw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			info, err := tr.Resolve(r)
			if err != nil || info == nil {
				info = &tenant.TenantInfo{ID: "default", Name: "默认租户", Strategy: tenant.StrategyShared}
			}
			ctx := tenant.WithContext(r.Context(), info)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
	s.router.Use(tenMw, middleware.Recovery(nil), middleware.RequestID(), middleware.CORS(nil))

	return s, nil
}

// registerRoutes wires all routes into the router.
func (s *Server) registerRoutes() {
	r := s.router
	aAPI := wmsMiddleware.AdminOnlyAPI(s.SessionMgr)

	// Static files
	fs := http.FileServer(http.Dir(s.cfg.StaticDir))
	r.GET("/static/", func(w http.ResponseWriter, req *http.Request) {
		req.URL.Path = strings.TrimPrefix(req.URL.Path, "/static")
		fs.ServeHTTP(w, req)
	})

	// Login page
	r.GET("/login", func(w http.ResponseWriter, req *http.Request) {
		s.TemplateMap["login"].ExecuteTemplate(w, "login.html", map[string]any{"HideSidebar": true})
	})
	r.POST("/login", func(w http.ResponseWriter, req *http.Request) {
		u, p := req.FormValue("username"), req.FormValue("password")
		if u != "admin" || p != "admin" {
			s.TemplateMap["login"].ExecuteTemplate(w, "login.html", map[string]any{"Error": "用户名或密码错误", "HideSidebar": true})
			return
		}
		tk, err := s.TokenMgr.IssueAccessToken(u, "tenant-1", []string{"admin"}, []string{"*"})
		if err != nil {
			s.TemplateMap["login"].ExecuteTemplate(w, "login.html", map[string]any{"Error": "令牌生成失败", "HideSidebar": true})
			return
		}
		http.SetCookie(w, &http.Cookie{Name: "i56_token", Value: tk, Path: "/", HttpOnly: true, MaxAge: 86400})
		http.Redirect(w, req, "/admin", 303)
	})
	r.GET("/logout", func(w http.ResponseWriter, req *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "i56_token", Value: "", Path: "/", MaxAge: -1})
		http.Redirect(w, req, "/login", 303)
	})

	// Login page — redirect to React SPA (SPA handles its own /admin/login route)
	r.GET("/admin/login", func(w http.ResponseWriter, req *http.Request) {
		// If already authenticated, redirect to admin dashboard
		if ck, err := req.Cookie("admin_session"); err == nil && s.SessionMgr.ValidateSession(ck.Value) != nil {
			http.Redirect(w, req, "/admin", 303)
			return
		}
		// Serve React SPA index.html from frontend build directory
		http.ServeFile(w, req, "/opt/i56/frontend/index.html")
	})
	r.POST("/admin/login", func(w http.ResponseWriter, req *http.Request) {
		u, p := req.FormValue("username"), req.FormValue("password")
		if !s.SessionMgr.Authenticate(u, p) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(401)
			w.Write([]byte(`{"error":"用户名或密码错误"}`))
			return
		}
		cookieValue := s.SessionMgr.CreateSession(u)
		http.SetCookie(w, &http.Cookie{Name: "admin_session", Value: cookieValue, Path: "/admin", HttpOnly: true, MaxAge: int(adminAuth.SessionTTL.Seconds())})
		// Return JSON for React SPA; also support form-based redirect
		accept := req.Header.Get("Accept")
		isAjax := req.Header.Get("X-Requested-With") == "XMLHttpRequest"
		if isAjax || strings.Contains(accept, "application/json") {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"success":true,"username":"` + u + `"}`))
			return
		}
		http.Redirect(w, req, "/admin", 303)
	})
	r.GET("/admin/logout", func(w http.ResponseWriter, req *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "admin_session", Value: "", Path: "/admin", MaxAge: -1})
		http.Redirect(w, req, "/admin/login", 303)
	})
	// Client login
	r.GET("/client/login", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, "/opt/i56/frontend/index.html")
	})
	r.POST("/client/login", func(w http.ResponseWriter, req *http.Request) {
		u, p := req.FormValue("username"), req.FormValue("password")
		if u == "" {
			u = req.FormValue("email")
		}
		if u == "" || p == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(401)
			w.Write([]byte(`{"error":"请输入账号和密码"}`))
			return
		}
		tk, err := s.TokenMgr.IssueAccessToken(u, "tenant-1", []string{"client"}, []string{"read:parcels"})
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			w.Write([]byte(`{"error":"token_generation_failed"}`))
			return
		}
		http.SetCookie(w, &http.Cookie{Name: "client_token", Value: tk, Path: "/client", HttpOnly: true, MaxAge: 86400})
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success":true,"username":"` + u + `"}`))
	})
	r.GET("/client/logout", func(w http.ResponseWriter, req *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "client_token", Value: "", Path: "/client", MaxAge: -1})
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success":true}`))
	})

	// ══════════════════════════════════════════
	// JSON API routes
	// ══════════════════════════════════════════
	// Device & Shelf API (explicit handlers to avoid route conflicts)
	r.GET("/admin/api/devices", func(w http.ResponseWriter, req *http.Request) {
		all := domain.DeviceStore.List()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(all)
	})
	r.GET("/admin/api/shelves", func(w http.ResponseWriter, req *http.Request) {
		all := domain.ShelfStore.List()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(all)
	})

	// Employee API — list users from RBAC repo
	r.GET("/admin/api/employees", func(w http.ResponseWriter, req *http.Request) {
		users, _, err := s.RBACRepo.ListUsers(req.Context(), 1, 0, 50)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			w.Write([]byte(`{"error":"` + err.Error() + `"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	})

	adminapi.RegisterSystemAPI(r, aAPI)
	adminapi.RegisterOMSAPI(r, aAPI, s.ParcelSvc, s.OrderSvc, s.WarehouseSvc,
		s.ClientRepo, s.RouteRepo, s.CourierRepo,
		s.ServiceRepo, s.LedgerRepo, s.DeclarantRepo, s.MemberRepo, s.AddressRepo,
		s.RoutePriceRepo, s.DeliveryFeeRepo, s.SurchargeRepo, s.APICredentialRepo,
		s.EventBus)
	adminapi.RegisterWMSAPI(r, aAPI, s.PrintRepo, s.WorkflowRepo, s.TaskDispatchRepo, s.WebhookRepo, s.WorkOrderRepo)
	adminapi.RegisterTMSAPI(r, aAPI)
	adminapi.RegisterCRMAPI(r, aAPI, s.ClientSvc, s.ClientRepo, s.LedgerRepo, s.DeclarantRepo, s.MemberRepo, s.AddressRepo, s.PdaRepo)
	adminapi.RegisterFinanceAPI(r, aAPI, s.OrderRepo, s.LedgerRepo, s.RouteRepo)
	adminapi.RegisterDashboardAPI(r, aAPI, s.OrderRepo, s.WarehouseRepo, s.ParcelRepo, s.PdaRepo)
	s.registerEventAPI(r, aAPI)
	s.registerTenantAPI(r, aAPI)
	s.registerRBACAPI(r, aAPI)
	s.registerStorageAPI(r, aAPI)
	s.registerWorkflowAPI(r, aAPI)
	s.registerNotificationAPI(r, aAPI)
	s.registerPluginAPI(r, aAPI)

	// ── Session check endpoint (used by React SPA login flow) ──
	r.GET("/admin/api/me", aAPI(func(w http.ResponseWriter, req *http.Request) {
		ck, err := req.Cookie("admin_session")
		if err != nil || ck.Value == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(401)
			w.Write([]byte(`{"error":"unauthorized"}`))
			return
		}
		sess := s.SessionMgr.ValidateSession(ck.Value)
		if sess == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(401)
			w.Write([]byte(`{"error":"invalid_session"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":1,"username":"` + sess.Username + `","real_name":"` + sess.Username + `","role_id":1,"role_name":"系统管理员"}`))
	}))

	clientapi.RegisterClientAPI(r, s.TokenMgr, s.ParcelSvc, s.OrderSvc,
		s.RouteRepo, s.CourierRepo, s.WarehouseSvc, s.LedgerRepo,
		s.DeclarantRepo, s.MemberRepo, s.ServiceRepo, s.WebhookRepo,
		s.AddressRepo, s.RoutePriceRepo, s.DeliveryFeeRepo, s.SurchargeRepo, s.APICredentialRepo)

	pdaapi.RegisterPDAAPI(r, s.PdaRepo, s.PDAOps, s.ParcelRepo, s.OrderRepo)

	// Generic Admin Delete
	r.POST("/admin/delete", aAPI(func(w http.ResponseWriter, req *http.Request) {
		page, idStr := req.URL.Query().Get("page"), req.URL.Query().Get("id")
		if page == "" || idStr == "" {
			http.Error(w, "missing page or id", 400)
			return
		}
		id, _ := strconv.ParseInt(idStr, 10, 64)
		switch page {
		case "warehouses":
			s.WarehouseRepo.Delete(req.Context(), 1, id)
		case "parcels":
			s.ParcelRepo.Delete(req.Context(), 1, id)
		case "route-templates":
			s.RouteRepo.Delete(req.Context(), 1, id)
		case "couriers":
			s.CourierRepo.Delete(req.Context(), idStr)
		case "clients":
			s.ClientRepo.Delete(req.Context(), 1, id)
		case "employees":
			s.RBACRepo.DeleteUser(req.Context(), 1, id)
		case "roles":
			s.RBACRepo.DeleteRole(req.Context(), 1, id)
		default:
			http.Error(w, "unsupported page: "+page, 400)
			return
		}
		w.Header().Set("HX-Refresh", "true")
		w.WriteHeader(200)
	}))

	// Health
	r.GET("/api/v1/health", func(w http.ResponseWriter, req *http.Request) {
		deps := "in-memory"
		if s.dbAvailable {
			deps = "postgres"
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":{"name":"I56 Framework","version":"` + I56Version + `","status":"ok","deps":"` + deps + `"}}`))
	})

	// Catch-all
	r.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			w.WriteHeader(404)
			return
		}
		http.Redirect(w, req, "/login", 303)
	}))
}

// Start begins listening and blocks until SIGINT/SIGTERM.
func (s *Server) Start() error {
	s.registerRoutes()

	log.Printf("I56 Framework %s listening on :%d", I56Version, s.cfg.Port)
	s.httpSrv = &http.Server{Addr: fmt.Sprintf(":%d", s.cfg.Port), Handler: s.router}

	go func() {
		if err := s.httpSrv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	return s.httpSrv.Shutdown(context.Background())
}

// seedRealData populates real module repos with demo data.
func seedRealData(s *Server) {
	ctx := context.Background()
	now := time.Now()

	s.WarehouseRepo.Create(ctx, 1, &whDomain.Warehouse{
		Name: "厦门仓", Code: "XM", Address: "福建省厦门市集美区",
		Contact: "仓库管理员", Phone: "0592-1234567", IsActive: true, TenantID: 1,
	})

	c := &custDomain.Client{
		Name: "EZ集運通", Code: "EZ001", ClientType: custDomain.ClientTypePlatform,
		ContactName: "运营经理", ContactPhone: "13800001111", ContactEmail: "ez@example.com",
		Balance: 10000, IsActive: true, TenantID: 1,
	}
	s.ClientRepo.Create(ctx, 1, c)

	s.RouteRepo.Create(ctx, &tmsDomain.Route{TenantID: 1, WarehouseID: 1, Name: "厦门→台湾(空运)", TransportType: "air", MinWeight: 0.5, VolumeCoeff: 6000, BaseWeightPrice: 20.0, BaseVolumePrice: 20.0, MinAmount: 50, MinDays: 1, MaxDays: 3, IsActive: true})
	s.RouteRepo.Create(ctx, &tmsDomain.Route{TenantID: 1, WarehouseID: 1, Name: "厦门→台湾(海快)", TransportType: "sea_express", MinWeight: 1.0, VolumeCoeff: 6000, BaseWeightPrice: 8.30, BaseVolumePrice: 15.0, MinAmount: 50, MinDays: 3, MaxDays: 7, IsActive: true})
	s.RouteRepo.Create(ctx, &tmsDomain.Route{TenantID: 1, WarehouseID: 1, Name: "厦门→台湾(海运)", TransportType: "sea", MinWeight: 10.0, VolumeCoeff: 6000, BaseWeightPrice: 3.20, BaseVolumePrice: 10.0, MinAmount: 50, MinDays: 5, MaxDays: 14, IsActive: true})

	type pdSeed struct{ tn, pn string; s parcelDomain.ParcelStatus; w float64 }
	for i, pd := range []pdSeed{
		{"SF1234567890", "手机壳", "pre_declared", 0.12},
		{"ZTO9876543210", "运动鞋", "received", 0.80},
		{"YTO1111222233", "T恤", "weighed", 0.25},
		{"STO4444555566", "蓝牙耳机", "stored", 0.15},
		{"HTKY7777888899", "数据线", "stored", 0.08},
		{"JD9999000011", "充电宝", "stored", 0.30},
		{"EMS1213141516", "化妆品套装", "stored", 1.20},
	} {
		tn := now.Add(-time.Duration(10-i) * 24 * time.Hour)
		s.ParcelRepo.Create(ctx, &parcelDomain.Parcel{
			TenantID: 1, WarehouseID: 1, ClientID: c.ID,
			TrackingNumber: pd.tn, ProductName: pd.pn, ParcelName: pd.pn,
			Status: parcelDomain.ParcelStatus(pd.s), CourierCode: "SF",
			CargoType: "general", ActualWeight: pd.w,
			CreatedAt: tn, UpdatedAt: tn,
		})
	}

	type orderSeed struct {
		orderNo, recipient, tracking, carrierTrack, customsNo, remark string
		memberID, routeID, daysAgo, parcelCount                       int
		status                                                        orderDomain.OrderStatus
		weight, chgWeight, price                                       float64
	}
	today := now
	for _, od := range []orderSeed{
		{"ORD-20260711-001", "王仁照", "80020737681100020001", "CT-8837291", "CN-20260711001", "空运急件", 1, 2, 0, 1, orderDomain.StatusInTransit, 0.56, 0.60, 8.00},
		{"ORD-20260711-002", "琦立工作室", "YT7631606603205", "", "", "", 2, 2, 0, 2, orderDomain.StatusPendingLoading, 1.05, 1.50, 18.00},
		{"ORD-20260710-001", "张致廷", "HTKY7777888899,JD9999000011", "", "", "", 1, 1, 1, 2, orderDomain.StatusPendingPacking, 0.33, 0.50, 11.50},
		{"ORD-20260709-001", "吴欣如", "ZTO20250601001,SF120011223344", "CT-8837292", "CN-20260709001", "已签收", 2, 3, 2, 3, orderDomain.StatusCompleted, 12.80, 15.00, 56.20},
		{"ORD-20260708-001", "王仁照", "YTO8822110011", "CT-8837293", "CN-20260708001", "", 1, 2, 3, 1, orderDomain.StatusCustomsClearance, 2.30, 2.50, 22.00},
		{"ORD-20260707-001", "琦立工作室", "STO5555666677", "", "", "大件运输", 2, 1, 4, 2, orderDomain.StatusLoaded, 4.50, 5.00, 45.00},
		{"ORD-20260706-001", "张致廷", "EMS9988776655,EMS1122334455", "CT-8837294", "CN-20260706001", "", 1, 3, 5, 4, orderDomain.StatusShipped, 28.50, 30.00, 98.00},
		{"ORD-20260705-001", "吴欣如", "SF5566778899,YTO4433221100", "", "", "待拣货", 2, 2, 6, 2, orderDomain.StatusPendingPicking, 0.78, 1.00, 9.50},
		// Fresh demo order for PDA workflow
		{"ORD-DEMO-001", "陈小美", "HTKY7777888899,JD9999000011", "", "", "PDA演示订单-待拣货", 1, 1, 0, 2, orderDomain.StatusPendingPicking, 0.45, 0.50, 12.00},
	} {
		s.OrderRepo.Create(ctx, &orderDomain.Order{
			TenantID: 1, WarehouseID: 1, ClientID: c.ID,
			OrderNo: od.orderNo, MemberID: int64(od.memberID), RouteID: int64(od.routeID),
			RecipientName: od.recipient, TrackingNumbers: od.tracking,
			Status: od.status, ParcelCount: od.parcelCount,
			TotalActualWeight: od.weight, TotalChargeableWeight: od.chgWeight,
			TotalPrice: od.price, CarrierTrackingNo: od.carrierTrack,
			CustomsNumber: od.customsNo, Remark: od.remark,
		})
	}

	// Patch dates
	if o1, _ := s.OrderRepo.GetByOrderNo(ctx, 1, "ORD-20260711-001"); o1 != nil {
		o1.CreatedAt = today; o1.UpdatedAt = today; s.OrderRepo.Update(ctx, o1)
	}
	if o2, _ := s.OrderRepo.GetByOrderNo(ctx, 1, "ORD-20260711-002"); o2 != nil {
		o2.CreatedAt = today; o2.UpdatedAt = today; s.OrderRepo.Update(ctx, o2)
	}
	if o3, _ := s.OrderRepo.GetByOrderNo(ctx, 1, "ORD-20260710-001"); o3 != nil {
		o3.CreatedAt = today.Add(-24 * time.Hour); o3.UpdatedAt = today.Add(-24 * time.Hour); s.OrderRepo.Update(ctx, o3)
	}
	if o4, _ := s.OrderRepo.GetByOrderNo(ctx, 1, "ORD-20260709-001"); o4 != nil {
		o4.CreatedAt = today.Add(-2 * 24 * time.Hour); o4.UpdatedAt = today.Add(-2 * 24 * time.Hour); s.OrderRepo.Update(ctx, o4)
	}
	if o5, _ := s.OrderRepo.GetByOrderNo(ctx, 1, "ORD-20260708-001"); o5 != nil {
		o5.CreatedAt = today.Add(-3 * 24 * time.Hour); o5.UpdatedAt = today.Add(-3 * 24 * time.Hour); s.OrderRepo.Update(ctx, o5)
	}
	if o6, _ := s.OrderRepo.GetByOrderNo(ctx, 1, "ORD-20260707-001"); o6 != nil {
		o6.CreatedAt = today.Add(-4 * 24 * time.Hour); o6.UpdatedAt = today.Add(-4 * 24 * time.Hour); s.OrderRepo.Update(ctx, o6)
	}
	if o7, _ := s.OrderRepo.GetByOrderNo(ctx, 1, "ORD-20260706-001"); o7 != nil {
		o7.CreatedAt = today.Add(-5 * 24 * time.Hour); o7.UpdatedAt = today.Add(-5 * 24 * time.Hour); s.OrderRepo.Update(ctx, o7)
	}
	if o8, _ := s.OrderRepo.GetByOrderNo(ctx, 1, "ORD-20260705-001"); o8 != nil {
		o8.CreatedAt = today.Add(-6 * 24 * time.Hour); o8.UpdatedAt = today.Add(-6 * 24 * time.Hour); s.OrderRepo.Update(ctx, o8)
	}

	s.LedgerRepo.Add(ctx, &custRepo.LedgerEntry{TenantID: 1, ClientID: c.ID, Amount: 5000, BalanceAfter: 5000, Type: "recharge", Description: ""})

	// Seed Members
	s.MemberRepo.Create(ctx, 1, &custDomain.ClientMember{
		ClientID: c.ID, Name: "王仁照", Phone: "886912345678", Email: "wang@example.com", MemberCode: "M001",
	})
	s.MemberRepo.Create(ctx, 1, &custDomain.ClientMember{
		ClientID: c.ID, Name: "琦立工作室", Phone: "886923456789", Email: "qili@example.com", MemberCode: "M002",
	})
	s.MemberRepo.Create(ctx, 1, &custDomain.ClientMember{
		ClientID: c.ID, Name: "张致廷", Phone: "886934567890", Email: "zhang@example.com", MemberCode: "M003",
	})
	s.MemberRepo.Create(ctx, 1, &custDomain.ClientMember{
		ClientID: c.ID, Name: "吴欣如", Phone: "886945678901", Email: "wu@example.com", MemberCode: "M004",
	})

	// Seed Declarants (BFT56-aligned with TW ID cards + verification)
	s.DeclarantRepo.Create(ctx, 1, &custDomain.Declarant{Name: "江明威", IDNumber: "T121272432", Phone: "0938150360", AuthStatus: "pending"})
	s.DeclarantRepo.Create(ctx, 1, &custDomain.Declarant{Name: "施慶堂", IDNumber: "F124136343", Phone: "0978285977", AuthStatus: "approved"})
	s.DeclarantRepo.Create(ctx, 1, &custDomain.Declarant{Name: "李咏靜", IDNumber: "T221353452", Phone: "0925822280", AuthStatus: "approved"})
	s.DeclarantRepo.Create(ctx, 1, &custDomain.Declarant{Name: "李安宜", IDNumber: "I200562760", Phone: "0908908031", AuthStatus: "approved"})
	s.DeclarantRepo.Create(ctx, 1, &custDomain.Declarant{Name: "王仁照", IDNumber: "A123456789", Phone: "886912345678", AuthStatus: "approved"})
	s.DeclarantRepo.Create(ctx, 1, &custDomain.Declarant{Name: "张三", IDNumber: "G091220000", Phone: "0912200008", AuthStatus: "rejected"})

	// Seed Addresses (BFT56-aligned Taiwan addresses)
	s.AddressRepo.Create(ctx, 1, &custDomain.MemberAddress{MemberID: 1, RecipientName: "李咏靜", Phone: "0925822280", City: "高雄市", District: "楠梓區", Address: "大學東路200號", IsDefault: true})
	s.AddressRepo.Create(ctx, 1, &custDomain.MemberAddress{MemberID: 1, RecipientName: "宋宜玲", Phone: "0975033593", City: "台北市", District: "內湖區", Address: "康寧路三段189巷2號", IsDefault: false})
	s.AddressRepo.Create(ctx, 1, &custDomain.MemberAddress{MemberID: 2, RecipientName: "彭如嫣", Phone: "0935106089", City: "高雄市", District: "燕巢區", Address: "角宿路398號", IsDefault: true})
	s.AddressRepo.Create(ctx, 1, &custDomain.MemberAddress{MemberID: 2, RecipientName: "林琝翰", Phone: "0909079333", City: "新北市", District: "永和區", Address: "忠孝街34巷9弄5號", IsDefault: false})

	// Seed Work Orders
	s.WorkOrderRepo.Create(ctx, &wdDomain.WorkOrder{TenantID: 1, WarehouseID: 1, Title: "入库上架-包裹SF1234567890", Description: "手机壳上架到A-01-01", Status: "pending", Priority: 1})
	s.WorkOrderRepo.Create(ctx, &wdDomain.WorkOrder{TenantID: 1, WarehouseID: 1, Title: "拣货打包-ORD-20260711-001", Description: "空运急件拣货", Status: "in_progress", Priority: 2})
	s.WorkOrderRepo.Create(ctx, &wdDomain.WorkOrder{TenantID: 1, WarehouseID: 1, Title: "异常处理-包裹破损", Description: "ZTO9876543210 外包装破损需重新包装", Status: "pending", Priority: 1})
}

// Template initialization — removed unused Go HTML login templates (React SPA handles /admin/login)
func initAdminTemplates() map[string]*template.Template {
	return map[string]*template.Template{}
}

func initClientTemplates() map[string]*template.Template {
	return map[string]*template.Template{}
}

// ─── Event Bus ────────────────────────────────────────────────────────────

// registerEventHandlers wires up domain event handlers.
func (s *Server) registerEventHandlers() {
	// Audit: log all events (wildcard)
	s.EventBus.Subscribe("*", func(ctx context.Context, e eventbus.Event) error {
		entry := map[string]interface{}{
			"name": e.EventName(), "time": e.OccurredAt().Format(time.RFC3339),
		}
		domain.AuditLogStore.Add(domain.AuditLog{Action: e.EventName(), Detail: fmt.Sprintf("%v", entry)})
		s.logEvent(entry)
		return nil
	}, true) // async

	// Business events — create notifications on order events
	s.EventBus.Subscribe("order.created", func(ctx context.Context, e eventbus.Event) error {
		s.logEvent(map[string]interface{}{"name": e.EventName(), "data": e.(*adminapi.DataEvent).Data})
		// Send in-app notification
		if s.NotifySvc != nil {
			data := e.(*adminapi.DataEvent).Data
			title := "新订单创建"
			body := fmt.Sprintf("订单 %v 已创建", data["order_no"])
			s.NotifySvc.Send(ctx, "in_app", notification.Message{
				Title: title, Body: body,
				To: []string{"admin"},
			})
		}
		// Fire webhooks
		if s.WebhookRepo != nil {
			subs, _ := s.WebhookRepo.ListSubs(ctx, 1)
			for _, sub := range subs {
				if sub.IsActive && sub.Event == "order.created" {
					go func(url string) {
						payload, _ := json.Marshal(e.(*adminapi.DataEvent).Data)
						resp, err := http.Post(url, "application/json", strings.NewReader(string(payload)))
						if err == nil {
							resp.Body.Close()
						}
					}(sub.URL)
				}
			}
		}
		return nil
	}, true)
	s.EventBus.Subscribe("parcel.received", func(ctx context.Context, e eventbus.Event) error {
		s.logEvent(map[string]interface{}{"name": e.EventName()})
		return nil
	}, true)
}

// logEvent records event in buffer (ring buffer, max 200).
func (s *Server) logEvent(entry map[string]interface{}) {
	s.evtMu.Lock()
	defer s.evtMu.Unlock()
	s.EventLog = append(s.EventLog, entry)
	if len(s.EventLog) > 200 {
		s.EventLog = s.EventLog[1:]
	}
}

// registerEventAPI exposes the event stream.
func (s *Server) registerEventAPI(r *router.Router, a func(h http.HandlerFunc) http.HandlerFunc) {
	r.GET("/admin/api/events", a(func(w http.ResponseWriter, req *http.Request) {
		s.evtMu.RLock()
		defer s.evtMu.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s.EventLog)
	}))
	r.POST("/admin/api/events/publish", a(func(w http.ResponseWriter, req *http.Request) {
		var payload struct {
			Name string                 `json:"name"`
			Data map[string]interface{} `json:"data,omitempty"`
		}
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			w.Write([]byte(`{"error":"invalid event payload"}`))
			return
		}
		ev := eventbus.NewEvent(payload.Name)
		s.EventBus.Publish(req.Context(), ev)
		s.logEvent(map[string]interface{}{"name": payload.Name, "data": payload.Data})
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":true}`))
	}))
}

// registerTenantAPI exposes tenant management endpoints.
func (s *Server) registerTenantAPI(r *router.Router, a func(h http.HandlerFunc) http.HandlerFunc) {
	r.GET("/admin/api/tenants", a(func(w http.ResponseWriter, req *http.Request) {
		tenants, err := s.TenantMgr.ListTenants(req.Context())
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			w.Write([]byte(`{"error":"` + err.Error() + `"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tenants)
	}))
	// Current tenant info from context
	r.GET("/admin/api/tenant", a(func(w http.ResponseWriter, req *http.Request) {
		info := tenant.FromContext(req.Context())
		if info == nil {
			info = &tenant.TenantInfo{ID: "default", Name: "默认租户", Strategy: tenant.StrategyShared}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(info)
	}))
}

// registerRBACAPI exposes RBAC permission checks and subject info.
func (s *Server) registerRBACAPI(r *router.Router, a func(h http.HandlerFunc) http.HandlerFunc) {
	// Helper: build Subject from request session
	buildSubject := func(req *http.Request) rbac.Subject {
		ck, _ := req.Cookie("admin_session")
		username := "anonymous"
		roleName := "viewer"
		if ck != nil {
			if sess := s.SessionMgr.ValidateSession(ck.Value); sess != nil {
				username = sess.Username
				// Map username to role
				if username == "admin" {
					roleName = "admin"
				} else {
					roleName = "operator"
				}
			}
		}
		ti := tenant.FromContext(req.Context())
		tid := "default"
		if ti != nil { tid = ti.ID }
		return rbac.Subject{UserID: username, TenantID: tid, RoleIDs: []string{roleName}}
	}

	// DataScope middleware — enforces warehouse/tenant data visibility
	s.router.Use(wmsMiddleware.DataScopeMiddleware(s.RBAC, buildSubject))

	// Current subject
	r.GET("/admin/api/rbac/subject", a(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(buildSubject(req))
	}))
	// Permission check: GET /admin/api/rbac/check?resource=parcel&action=create
	r.GET("/admin/api/rbac/check", a(func(w http.ResponseWriter, req *http.Request) {
		resource := req.URL.Query().Get("resource")
		action := req.URL.Query().Get("action")
		subj := buildSubject(req)
		ok := s.RBAC.Enforce(req.Context(), subj, resource, action)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"resource": resource, "action": action, "allowed": ok})
	}))
	// Data scope
	r.GET("/admin/api/rbac/datascope", a(func(w http.ResponseWriter, req *http.Request) {
		resource := req.URL.Query().Get("resource")
		subj := buildSubject(req)
		scope := s.RBAC.DataScope(req.Context(), subj, resource)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"resource": resource, "scope": scope.String()})
	}))
}

// registerStorageAPI exposes file upload/download endpoints.
func (s *Server) registerStorageAPI(r *router.Router, a func(h http.HandlerFunc) http.HandlerFunc) {
	r.POST("/admin/api/storage/upload", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseMultipartForm(32 << 20) // 32 MB
		file, header, err := req.FormFile("file")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			w.Write([]byte(`{"error":"no file uploaded"}`))
			return
		}
		defer file.Close()
		bucket := req.FormValue("bucket")
		if bucket == "" { bucket = "default" }
		// Use tenant-aware prefix if available
		if ti := tenant.FromContext(req.Context()); ti != nil && ti.ID != "" {
			bucket = ti.ID + "/" + bucket
		}
		url, err := s.Storage.Upload(req.Context(), bucket, header.Filename, file, header.Size, header.Header.Get("Content-Type"))
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			w.Write([]byte(`{"error":"upload failed: ` + err.Error() + `"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": true, "url": url, "filename": header.Filename, "size": header.Size})
	}))
	// List files in a bucket
	r.GET("/admin/api/storage/list", a(func(w http.ResponseWriter, req *http.Request) {
		bucket := req.URL.Query().Get("bucket")
		if bucket == "" { bucket = "default" }
		// For local storage, just list directory
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"bucket":"` + bucket + `","files":[],"note":"use upload to add files"}`))
	}))
}

// ─── In-memory Workflow Store ─────────────────────────────────────────

type memWorkflowStore struct {
	mu        sync.RWMutex
	instances map[string]*workflow.ProcessInstance
}

func (s *memWorkflowStore) GetDefinition(ctx context.Context, id string) (*workflow.ProcessDefinition, error) { return nil, nil }
func (s *memWorkflowStore) SaveInstance(ctx context.Context, inst *workflow.ProcessInstance) error {
	s.mu.Lock(); defer s.mu.Unlock()
	s.instances[inst.ID] = inst
	return nil
}
func (s *memWorkflowStore) GetInstance(ctx context.Context, id string) (*workflow.ProcessInstance, error) {
	s.mu.RLock(); defer s.mu.RUnlock()
	if inst, ok := s.instances[id]; ok { return inst, nil }
	return nil, fmt.Errorf("instance not found: %s", id)
}

// ─── Workflow API ────────────────────────────────────────────────────

func (s *Server) registerWorkflowAPI(r *router.Router, a func(h http.HandlerFunc) http.HandlerFunc) {
	// Start new approval process: POST /admin/api/workflow/start {definition_id, amount}
	r.POST("/admin/api/workflow/start", a(func(w http.ResponseWriter, req *http.Request) {
		var body struct {
			DefinitionID string  `json:"definition_id"`
			Amount       float64 `json:"amount"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400); w.Write([]byte(`{"error":"invalid json"}`))
			return
		}
		inst, err := s.Workflow.StartProcess(req.Context(), body.DefinitionID, map[string]any{"amount": body.Amount})
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500); w.Write([]byte(`{"error":"` + err.Error() + `"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(inst)
	}))
	// Approve/reject step: POST /admin/api/workflow/transition {instance_id, to_state}
	r.POST("/admin/api/workflow/transition", a(func(w http.ResponseWriter, req *http.Request) {
		var body struct {
			InstanceID string `json:"instance_id"`
			ToState    string `json:"to_state"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400); w.Write([]byte(`{"error":"invalid json"}`))
			return
		}
		if err := s.Workflow.Transition(req.Context(), body.InstanceID, body.ToState, nil); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400); w.Write([]byte(`{"error":"` + err.Error() + `"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":true}`))
	}))
	// List definitions
	r.GET("/admin/api/workflow/definitions", a(func(w http.ResponseWriter, req *http.Request) {
		defs := []map[string]interface{}{{
			"id": "purchase-approval", "name": "采购审批流程",
			"states": []string{"开始", "部门审批", "财务审批", "总经理审批", "完成"},
		}}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(defs)
	}))
}

// ─── Noop Logger (for framework modules that require non-nil logger) ───

type noopLogger struct{}

func (l noopLogger) Debug(msg string, args ...any)      {}
func (l noopLogger) Info(msg string, args ...any)       {}
func (l noopLogger) Warn(msg string, args ...any)       {}
func (l noopLogger) Error(msg string, args ...any)      {}
func (l noopLogger) With(args ...any) logger.Logger        { return l }
func (l noopLogger) WithGroup(name string) logger.Logger   { return l }

// ─── Mem Notification Channel ────────────────────────────────────────

type memNotifChannel struct {
	mu       sync.RWMutex
	messages []map[string]interface{}
}

func (c *memNotifChannel) Name() string { return "in_app" }
func (c *memNotifChannel) Send(ctx context.Context, msg notification.Message) error {
	c.mu.Lock(); defer c.mu.Unlock()
	c.messages = append(c.messages, map[string]interface{}{
		"title": msg.Title, "body": msg.Body, "to": msg.To, "data": msg.Data,
		"sent_at": time.Now().Format(time.RFC3339),
	})
	return nil
}
func (c *memNotifChannel) List() []map[string]interface{} {
	c.mu.RLock(); defer c.mu.RUnlock()
	return c.messages
}

// ─── Notification API ────────────────────────────────────────────────

func (s *Server) registerNotificationAPI(r *router.Router, a func(h http.HandlerFunc) http.HandlerFunc) {
	// Send notification: POST /admin/api/notify/send {channel, title, body, to}
	r.POST("/admin/api/notify/send", a(func(w http.ResponseWriter, req *http.Request) {
		var body struct {
			Channel string   `json:"channel"`
			Title   string   `json:"title"`
			Body    string   `json:"body"`
			To      []string `json:"to,omitempty"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400); w.Write([]byte(`{"error":"invalid json"}`))
			return
		}
		msg := notification.Message{Title: body.Title, Body: body.Body, To: body.To}
		if body.Channel != "" {
			if err := s.NotifySvc.Send(req.Context(), body.Channel, msg); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(500); w.Write([]byte(`{"error":"`+err.Error()+`"}`))
				return
			}
		} else {
			s.NotifySvc.SendAll(req.Context(), msg)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":true}`))
	}))
	// List notifications: GET /admin/api/notify/inbox
	r.GET("/admin/api/notify/inbox", a(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if s.NotifChan != nil {
			json.NewEncoder(w).Encode(s.NotifChan.List())
		} else {
			w.Write([]byte(`[]`))
		}
	}))
}

// ─── Plugin API ──────────────────────────────────────────────────────

func (s *Server) registerPluginAPI(r *router.Router, a func(h http.HandlerFunc) http.HandlerFunc) {
	r.GET("/admin/api/plugins", a(func(w http.ResponseWriter, req *http.Request) {
		svcs := []map[string]string{
			{"name": "eventbus", "type": "events.Bus", "module": "framework"},
			{"name": "tenant", "type": "tenant.TenantStore", "module": "framework"},
			{"name": "rbac", "type": "rbac.Enforcer", "module": "framework"},
			{"name": "storage", "type": "storage.StorageProvider", "module": "framework"},
			{"name": "workflow", "type": "workflow.Engine", "module": "framework"},
			{"name": "notification", "type": "notification.Service", "module": "framework"},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(svcs)
	}))
	// Resolve a service: GET /admin/api/plugins/resolve?name=storage
	r.GET("/admin/api/plugins/resolve", a(func(w http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("name")
		svc := s.Registry.Resolve(name)
		w.Header().Set("Content-Type", "application/json")
		if svc == nil {
			w.Write([]byte(`{"name":"` + name + `","found":false}`))
		} else {
			json.NewEncoder(w).Encode(map[string]interface{}{"name": name, "found": true, "type": fmt.Sprintf("%T", svc)})
		}
	}))
}
