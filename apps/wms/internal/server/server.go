// Package server provides HTTP server setup, route wiring, and seed data.
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/i56/framework/core/auth"
	"github.com/i56/framework/core/config"
	"github.com/i56/framework/core/middleware"
	"github.com/i56/framework/core/router"
	"github.com/i56/framework/core/scheduler"
	jpkg "github.com/i56/framework/core/jwt"
	"github.com/i56/framework/core/sse"
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

	SessionMgr *adminAuth.SessionManager
	TokenMgr   *auth.TokenManager

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

	sysCfg, _ := config.Load()
	tm, err := auth.NewTokenManager(sysCfg.Auth)
	if err != nil {
		return nil, fmt.Errorf("token manager: %w", err)
	}
	s.TokenMgr = tm
	s.SessionMgr = adminAuth.NewSessionManager()

	// PostgreSQL
	if err := db.Connect(cfg.DBDSN); err == nil {
		s.dbAvailable = true
		defer db.Close()
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
	s.router.Use(middleware.Recovery(nil), middleware.RequestID(), middleware.CORS(nil))

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
		tk, _ := s.TokenMgr.IssueAccessToken(u, "tenant-1", []string{"admin"}, []string{"*"})
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
		tk, _ := s.TokenMgr.IssueAccessToken(u, "tenant-1", []string{"client"}, []string{"read:parcels"})
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
		s.ClientRepo, s.RouteRepo, s.CourierRepo, s.ServiceRepo,
		s.LedgerRepo, s.DeclarantRepo, s.MemberRepo, s.AddressRepo,
		s.RoutePriceRepo, s.DeliveryFeeRepo, s.SurchargeRepo, s.APICredentialRepo)
	adminapi.RegisterWMSAPI(r, aAPI, s.PrintRepo, s.WorkflowRepo, s.TaskDispatchRepo, s.WebhookRepo, s.WorkOrderRepo)
	adminapi.RegisterTMSAPI(r, aAPI)
	adminapi.RegisterCRMAPI(r, aAPI, s.ClientSvc, s.ClientRepo, s.LedgerRepo, s.DeclarantRepo, s.MemberRepo, s.AddressRepo, s.PdaRepo)
	adminapi.RegisterFinanceAPI(r, aAPI)

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
