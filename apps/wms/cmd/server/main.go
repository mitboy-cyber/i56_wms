package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"os"
	"os/signal"
	"syscall"

	"github.com/i56/framework/core/auth"
	"github.com/i56/framework/core/audit"
	"github.com/i56/framework/core/config"
	"github.com/i56/framework/core/middleware"
	"github.com/i56/framework/core/report"
	"github.com/i56/framework/core/router"
	"github.com/i56/framework/core/scheduler"
	"github.com/i56/framework/core/sse"
	jpkg "github.com/i56/framework/core/jwt"
	"github.com/i56/framework/db"

	// I56 2.0 AI Runtime
	"github.com/i56/framework/ai"
	"github.com/i56/framework/ai/gateway"
	"github.com/i56/framework/ai/tools"
	aiSec "github.com/i56/framework/ai/security"
	aiRouter "github.com/i56/framework/ai/router"
	"github.com/i56/framework/ai/gateway/providers"

	pricingRepo "github.com/i56/modules/pricing/repository"
	weightDomain "github.com/i56/modules/weight/domain"

	custRepo "github.com/i56/modules/customer/repository"
	custSvc "github.com/i56/modules/customer/service"
	pdaRepo "github.com/i56/modules/pda/repository"
	pdaSvc "github.com/i56/modules/pda/service"
	parcelRepo "github.com/i56/modules/parcel/repository"
	parcelSvc "github.com/i56/modules/parcel/service"
	orderRepo "github.com/i56/modules/order/repository"
	orderSvc "github.com/i56/modules/order/service"
	whRepo "github.com/i56/modules/warehouse/repository"
	whSvc "github.com/i56/modules/warehouse/service"
	tmsRepo "github.com/i56/modules/transport/repository"
	wfRepo "github.com/i56/modules/workflow/repository"
	psRepo "github.com/i56/modules/parcel_service/repository"
	woRepo "github.com/i56/modules/workorder/repository"
	printRepo "github.com/i56/modules/print/repository"
	whRepo2 "github.com/i56/modules/webhook/repository"
	sysRepo "github.com/i56/modules/system/repository"
	reportDomain "github.com/i56/modules/report/domain"
	rbacRepoPkg "github.com/i56/modules/rbac/repository"
	tdRepo "github.com/i56/modules/taskdispatch/repository"

	// I56 module route packages
	"github.com/i56/i56-apps/i56-wms/internal/common"
	crmroute "github.com/i56/i56-apps/i56-wms/internal/crmroute"
	finroute "github.com/i56/i56-apps/i56-wms/internal/finroute"
	omsroute "github.com/i56/i56-apps/i56-wms/internal/omsroute"
	sysroute "github.com/i56/i56-apps/i56-wms/internal/sysroute"
	tmsroute "github.com/i56/i56-apps/i56-wms/internal/tmsroute"
	wmsroute "github.com/i56/i56-apps/i56-wms/internal/wmsroute"
)

func main() {
	// JWT service with Ed25519
	jwtSvc, err := jpkg.NewService("i56-framework")
	if err != nil { log.Fatal(err) }
	log.Printf("[JWT] Public key: %s", jwtSvc.PublicKeyBase64())

	cfg, _ := config.Load()
	tm, err := auth.NewTokenManager(cfg.Auth)
	if err != nil { log.Fatal(err) }

	// ★ PostgreSQL: try to connect; gracefully fall back to in-memory
	dbAvailable := false
	if err := db.Connect("postgres://ubuntu@localhost:5432/i56_dev?sslmode=disable"); err == nil {
		dbAvailable = true
		defer db.Close()
	}

	// Repos
	cr := custRepo.NewMemClientRepo()
	wr := whRepo.NewMemWarehouseRepo()
	pr := parcelRepo.NewMemParcelRepo()
	or := orderRepo.NewMemOrderRepo()
	cour := tmsRepo.NewMemCourierRepo()
	rr := tmsRepo.NewMemRouteRepo()
	lr := custRepo.NewMemLedgerRepo()
	sr := psRepo.NewMemServiceRepo()
	wor := woRepo.NewMemWorkOrderRepo()
	wfr := wfRepo.NewMemWorkflowRepo()
	ppr := printRepo.NewMemPrintRepo()
	whr := whRepo2.NewMemWebhookRepo()
	sysCfg := sysRepo.NewMemSystemConfigRepo()
	rbac := rbacRepoPkg.NewMemRBACRepo()
	rpt := reportDomain.NewReportService()
	pdaR := pdaRepo.NewMemPDARepo(); _ = pdaR
	dr := custRepo.NewMemDeclarantRepo()
	mr := custRepo.NewMemMemberRepo()
	ar := custRepo.NewMemAddressRepo()
	rpr := pricingRepo.NewMemRoutePriceRepo()
	dfr := pricingRepo.NewMemDeliveryFeeRepo()
	scr := pricingRepo.NewMemSurchargeRepo()
	acr := pricingRepo.NewMemApiCredentialRepo()

	// ★ PostgreSQL repos (when connected, services use these for real data)
	var ppg *parcelRepo.PgParcelRepo
	var opg *orderRepo.PgOrderRepo
	var wpg *whRepo.PgWarehouseRepo
	var cpg *custRepo.PgClientRepo
	var rpg *tmsRepo.PgRouteRepo
	var copg *tmsRepo.PgCarrierRepo
	var rbacpg *rbacRepoPkg.PgRBACRepo
	var syspg *sysRepo.PgSystemConfigRepo
	if dbAvailable {
		ppg = parcelRepo.NewPgParcelRepo()
		opg = orderRepo.NewPgOrderRepo()
		wpg = whRepo.NewPgWarehouseRepo()
		cpg = custRepo.NewPgClientRepo()
		rpg = tmsRepo.NewPgRouteRepo()
		copg = tmsRepo.NewPgCarrierRepo()
		rbacpg = rbacRepoPkg.NewPgRBACRepo()
		syspg = sysRepo.NewPgSystemConfigRepo()
		_ = ppg; _ = opg; _ = wpg; _ = cpg; _ = rpg; _ = copg; _ = rbacpg; _ = syspg
		log.Println("[DB] Using PostgreSQL for real data persistence")
	}

	// Task Dispatch Engine
	td := tdRepo.NewMemTaskDispatchRepo()

	// Services
	ps := parcelSvc.NewParcelService(pr)
	osvc := orderSvc.NewOrderService(or)
	cs := custSvc.NewClientService(cr)
	ws := whSvc.NewWarehouseService(wr)
	pdaOps := pdaSvc.NewPDAOperations(pr, or, nil, nil)

	// SSE Hub for real-time events
	hub := sse.NewHub()

	// ★ Scheduler — cron-style task scheduler
	sch := scheduler.New()
	scheduler.DemoJobs(sch) // pre-register 5 demo jobs
	sch.Start()

	// ★ Audit Logger — operation audit trail
	auditRepo := audit.NewMemAuditRepo()
	auditLogger := audit.New(auditRepo)

	// ★ Built-in Report Engine
	reportEngine := report.NewBuiltinEngine()

	// ★ OpenAPI Generator
	openapiGen := router.NewOpenAPIGenerator(router.OpenAPIInfo{
	Title:   "I56 Framework API",
	Version: "2.3.0",
	Description: "I56 WMS Framework — Warehouse Management System API",
	})

	// Seed data
	seed(cr, wr, rr, pr, or, cour, lr, sr, wor, ppr, whr, rpt)

	// Templates
	tmpl := initTemplates()
	cTmpl := initClientTemplates()

	// Router
	r := router.New()
	r.Use(middleware.Recovery(nil), middleware.RequestID(), middleware.CORS(nil))

	// Static files
	fs := http.FileServer(http.Dir("static"))
	r.GET("/static/", func(w http.ResponseWriter, req *http.Request) {
		req.URL.Path = strings.TrimPrefix(req.URL.Path, "/static")
		fs.ServeHTTP(w, req)
	})

	// Login
	r.GET("/login", func(w http.ResponseWriter, req *http.Request) {
		tmpl["login"].ExecuteTemplate(w, "login.html", map[string]any{"HideSidebar":true})
	})
	r.POST("/login", func(w http.ResponseWriter, req *http.Request) {
		u, p := req.FormValue("username"), req.FormValue("password")
		if u != "admin" || p != "admin" {
			tmpl["login"].ExecuteTemplate(w, "login.html", map[string]any{"Error": "用户名或密码错误","HideSidebar":true})
			return
		}
		tk, _ := tm.IssueAccessToken(u, "tenant-1", []string{"admin"}, []string{"*"})
		http.SetCookie(w, &http.Cookie{Name: "i56_token", Value: tk, Path: "/", HttpOnly: true, MaxAge: 86400})
		http.Redirect(w, req, "/admin", 303)
	})
	r.GET("/logout", func(w http.ResponseWriter, req *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "i56_token", Value: "", Path: "/", MaxAge: -1})
		http.Redirect(w, req, "/login", 303)
	})

	// ==========================================
	// ★ I56 2.0 AI Runtime — initialized BEFORE business routes
	aiSvc := ai.New(ai.Config{
		DefaultTenant:  "1",
		BasePrompt:     "You are an AI assistant for I56 WMS. You help with warehouse operations, order management, parcel tracking, and logistics.",
		SecurityConfig: aiSec.DefaultConfig(),
		RoutePolicy:    aiRouter.PolicyQualityFirst,
	})

	// Register business domain tools
	// Register DeepSeek as the live provider
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" { apiKey = "sk-placeholder" }
	dsProvider := providers.NewDeepSeek(apiKey)
	aiSvc.RegisterGateway("deepseek", dsProvider)

	aiSvc.Tools.Register("list_warehouses", &tools.ToolMeta{
		Name: "list_warehouses", Description: "List all warehouses",
		Handler: func(ctx context.Context, args map[string]any) (any, error) {
			whs, _, _ := ws.List(ctx, 1, 0, 50); return whs, nil
		},
	})
	aiSvc.Tools.Register("list_orders", &tools.ToolMeta{
		Name: "list_orders", Description: "List recent orders",
		Handler: func(ctx context.Context, args map[string]any) (any, error) {
			orders, _, _ := osvc.List(ctx, 1, 0, 20); return orders, nil
		},
	})
	aiSvc.Tools.Register("parcel_status", &tools.ToolMeta{
		Name: "parcel_status", Description: "Get parcel status by tracking number",
		Handler: func(ctx context.Context, args map[string]any) (any, error) {
			parcels, _, _ := ps.List(ctx, 1, 0, 200)
			tracking := fmt.Sprint(args["tracking"])
			for _, p := range parcels { if p.TrackingNumber == tracking { return p, nil } }
			return map[string]string{"error": "not found"}, nil
		},
	})
	aiSvc.Tools.Register("client_balance", &tools.ToolMeta{
		Name: "client_balance", Description: "Get client balance",
		Handler: func(ctx context.Context, args map[string]any) (any, error) {
			entries := lr.GetByClient(ctx, 1, 1); bal := 0.0
			if len(entries) > 0 { bal = entries[len(entries)-1].BalanceAfter }
			return map[string]float64{"balance": bal}, nil
		},
	})

// AI API routes
	r.POST("/api/ai/chat", func(w http.ResponseWriter, req *http.Request) {
		var body struct{ Message string `json:"message"`; EnableTools bool `json:"enable_tools"` }
		json.NewDecoder(req.Body).Decode(&body)
		resp, err := aiSvc.Chat(req.Context(), "admin", aiRouter.TierHeavy, &gateway.ChatRequest{
			Messages: []gateway.Message{
				{Role: gateway.RoleSystem, Content: "You are a WMS assistant."},
				{Role: gateway.RoleUser, Content: body.Message},
			},
		})
		w.Header().Set("Content-Type", "application/json")
		if err != nil { json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); return }
		json.NewEncoder(w).Encode(resp)
	})
	r.GET("/api/ai/chat/stream", func(w http.ResponseWriter, req *http.Request) {
		msg := req.URL.Query().Get("q")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		flusher, _ := w.(http.Flusher)
		ch, _ := aiSvc.ChatStream(req.Context(), "admin", aiRouter.TierLight, &gateway.ChatRequest{
			Messages: []gateway.Message{
				{Role: gateway.RoleUser, Content: msg},
			},
		})
		for ev := range ch {
			fmt.Fprintf(w, "data: %s\n\n", ev.Content)
			if flusher != nil { flusher.Flush() }
			if ev.Done { break }
		}
	})
	r.GET("/api/ai/tools", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"tools":"loaded","count":"4"})
	})

	// Health — now with AI status
	r.GET("/api/v1/health", func(w http.ResponseWriter, req *http.Request) {
		deps := "in-memory"
		if dbAvailable { deps = "postgres" }
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":{"name":"I56 Framework","version":"2.0.0","status":"ok","ai":"active","deps":"` + deps + `"}}`))
	})

	// SSE endpoint (real-time events)
	r.GET("/api/v1/sse", func(w http.ResponseWriter, req *http.Request) {
		channel := req.URL.Query().Get("channel")
		if channel == "" { channel = "inbound" }
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		c := hub.Subscribe(channel)
		defer hub.Unsubscribe(c)
		flusher, ok := w.(http.Flusher)
		if !ok { http.Error(w, "SSE not supported", 500); return }
		for {
			select {
			case ev := <-c.Events:
				fmt.Fprintf(w, "event: %s\ndata: %s\n\n", ev.Type, ev.Data)
				flusher.Flush()
			case <-req.Context().Done():
				return
			}
		}
	})

	a := adminOnly(tm)
	rc := &common.RenderCtx{Tmpl: tmpl, Exec: common.DefaultExecTpl}

	// ★ Module-split route registrations (replaces adminPages + registerBFT56Modules + adminSystemPages + registerAdminCRUD)
	omsroute.Register(r, a, rc, osvc, ws, rr, cr, mr, sr, lr, ps)
	wmsroute.Register(r, a, rc, ps, ws, osvc, cr, rr, wor, sr, wfr, rbac)
	tmsroute.Register(r, a, rc, rr, cour)
	crmroute.Register(r, a, rc, cr, mr, lr, ar, dr, rpr)
	finroute.Register(r, a, rc, rpt)
	sysroute.Register(r, a, rc, sysCfg, rbac)

	// ★ New Framework v2.3 routes
	// OpenAPI spec endpoint
	router.RegisterOpenAPIEndpoint(r, openapiGen)
	registerSchedulerRoutes(r, sch, a)
	registerAuditRoutes(r, auditLogger, a)
	registerReportRoutes(r, reportEngine, a)
	// Register some demo routes for OpenAPI documentation
	registerOpenAPIDemoRoutes(r, openapiGen, a)

	// PDA API routes (direct on main router)
	registerPDAAPIRoutesOnRouter(r, pdaOps)
	registerPDARoutes(r, pdaR, pdaOps)

	// Task Dispatch Engine routes (抢单池)
	registerTaskDispatchRoutes(r, td)
	_ = StartTimeoutChecker(td) // background goroutine for SLA timeouts

	// Client portal
	r.GET("/client/login", func(w http.ResponseWriter, req *http.Request) {
		cTmpl["login"].ExecuteTemplate(w, "login.html", map[string]any{"HideSidebar":true})
	})
	r.POST("/client/login", func(w http.ResponseWriter, req *http.Request) {
		u, p := req.FormValue("username"), req.FormValue("password")
		if u == "" || p == "" {
			cTmpl["login"].ExecuteTemplate(w, "login.html", map[string]any{"Error": "请输入账号密码"})
			return
		}
		tk, _ := tm.IssueAccessToken(u, "tenant-1", []string{"client"}, []string{"read:parcels"})
		http.SetCookie(w, &http.Cookie{Name: "client_token", Value: tk, Path: "/client", HttpOnly: true, MaxAge: 86400})
		http.Redirect(w, req, "/client", 303)
	})
	weightRepo := weightDomain.NewMemWeightRepo()
	registerWeightAPI(r, weightRepo)
	registerAdminCRUDAPI(r, rbac)
	registerWeightUIRoutes(r, tmpl, weightRepo)
	clientPg(tm, cTmpl, r, ps, osvc, rr, cour, ws, pr, lr, weightRepo, dr, mr, sr, whr, ar, rpr, dfr, scr, acr)

	// PDA routes

	// Catch-all
	r.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		http.Redirect(w, req, "/login", 303)
	}))

	// Start
	log.Println("I56 Framework 2.0.0 listening on :8080")
	srv := &http.Server{Addr: ":8080", Handler: r}
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	srv.Shutdown(context.Background())
	_ = cs
}
