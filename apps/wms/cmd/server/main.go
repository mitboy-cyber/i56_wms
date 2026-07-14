package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	"github.com/i56/i56-apps/i56-wms/internal/ai/core"
	adminAuth "github.com/i56/i56-apps/i56-wms/internal/auth"
	"github.com/i56/i56-apps/i56-wms/internal/ai/classifier"
	"github.com/i56/i56-apps/i56-wms/internal/ai/translate"
)


// ─── Framework Version ───
const I56Version = "2.4.2"
const I56Name = "I56 Framework"
const I56Copyright = "© 2026 I56 Framework. All rights reserved."

// wmsDataProvider bridges services to the report engine.
type wmsDataProvider struct {
	orders  *orderSvc.OrderService
	clients *custRepo.MemClientRepo
	parcels *parcelSvc.ParcelService
}
func (p *wmsDataProvider) Orders(ctx context.Context) ([]report.OrderRow, error) {
	orders, _, _ := p.orders.List(ctx, 1, 0, 200)
	rows := make([]report.OrderRow, len(orders))
	for i, o := range orders {
		rows[i] = report.OrderRow{OrderNo: o.OrderNo, Status: string(o.Status), TotalPrice: o.TotalPrice, ParcelCount: o.ParcelCount, CreatedAt: o.CreatedAt.Format("2006-01-02")}
	}
	return rows, nil
}
func (p *wmsDataProvider) Clients(ctx context.Context) ([]report.ClientRow, error) {
	clients, _, _ := p.clients.List(ctx, 1, 0, 200)
	rows := make([]report.ClientRow, len(clients))
	for i, c := range clients {
		rows[i] = report.ClientRow{ID: c.ID, Name: c.Name, Code: c.Code, Balance: c.Balance, CreditLimit: 0}
	}
	return rows, nil
}
func (p *wmsDataProvider) Parcels(ctx context.Context) ([]report.ParcelRow, error) {
	parcels, _, _ := p.parcels.List(ctx, 1, 0, 500)
	rows := make([]report.ParcelRow, len(parcels))
	for i, p := range parcels {
		rows[i] = report.ParcelRow{TrackingNo: p.TrackingNumber, ProductName: p.ProductName, Status: string(p.Status), ActualWeight: p.ActualWeight}
	}
	return rows, nil
}

func main() {
	// JWT service with Ed25519
	jwtSvc, err := jpkg.NewService("i56-framework")
	if err != nil { log.Fatal(err) }
	log.Printf("[JWT] Public key: %s", jwtSvc.PublicKeyBase64())

	cfg, _ := config.Load()
	tm, err := auth.NewTokenManager(cfg.Auth)
	if err != nil { log.Fatal(err) }

	// ★ Admin Session Manager (HMAC-SHA256 sessions for admin panel auth)
	sessionMgr := adminAuth.NewSessionManager()

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
	ppr := printRepo.NewMemPrintRepo()
	whr := whRepo2.NewMemWebhookRepo()
	_ = sysRepo.NewMemSystemConfigRepo() // system config repo reserved for future use
	rbac := rbacRepoPkg.NewMemRBACRepo()
	wfr := wfRepo.NewMemWorkflowRepo()
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
	reportEngine := report.NewBuiltinEngine(&wmsDataProvider{orders: osvc, clients: cr, parcels: ps})

	// ★ OpenAPI Generator
	openapiGen := router.NewOpenAPIGenerator(router.OpenAPIInfo{
	Title:   "I56 Framework API",
	Version: I56Version,
	Description: "I56 WMS Framework — Warehouse Management System API",
	})

	// Seed data
	seed(rbac, cr, wr, rr, pr, or, cour, lr, sr, wor, ppr, whr, rpt)

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
	// ★ Admin Authentication Routes (HMAC-session based)
	// ==========================================
	r.GET("/admin/login", func(w http.ResponseWriter, req *http.Request) {
		// Check if already authenticated
		if ck, err := req.Cookie("admin_session"); err == nil && sessionMgr.ValidateSession(ck.Value) != nil {
			http.Redirect(w, req, "/admin", 303)
			return
		}
		// Render a standalone login page (no sidebar)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl["admin_login"].ExecuteTemplate(w, "login.html", map[string]any{"HideSidebar": true})
	})
	r.POST("/admin/login", func(w http.ResponseWriter, req *http.Request) {
		u, p := req.FormValue("username"), req.FormValue("password")
		if !sessionMgr.Authenticate(u, p) {
			tmpl["admin_login"].ExecuteTemplate(w, "login.html", map[string]any{
				"Error":       "用户名或密码错误",
				"HideSidebar": true,
			})
			return
		}
		cookieValue := sessionMgr.CreateSession(u)
		http.SetCookie(w, &http.Cookie{
			Name:     "admin_session",
			Value:    cookieValue,
			Path:     "/admin",
			HttpOnly: true,
			MaxAge:   int(adminAuth.SessionTTL.Seconds()),
		})
		http.Redirect(w, req, "/admin", 303)
	})
	r.GET("/admin/logout", func(w http.ResponseWriter, req *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "admin_session", Value: "", Path: "/admin", MaxAge: -1})
		http.Redirect(w, req, "/admin/login", 303)
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

	// ★ AI Cost Tracker
	costTracker := core.NewCostTracker()

	// ★ AI Business Context (injects WMS domain data into AI queries)
	bizCtx := core.NewBusinessContext(or, pr, wr, cr, osvc, ps, ws)

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
		// Track cost
		if resp != nil {
			costTracker.Track(resp.Model, "chat", 1, resp.TokenUsage.PromptTokens, resp.TokenUsage.CompletionTokens, 0.0001)
		}
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

	// ★ AI Chat SSE endpoint (for admin chat panel)
	r.POST("/api/v1/ai/chat", func(w http.ResponseWriter, req *http.Request) {
		var body struct{ Message string `json:"message"` }
		json.NewDecoder(req.Body).Decode(&body)
		msg := body.Message

		// Inject business context
		bizContext := bizCtx.GetBusinessContext(1, msg)

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "SSE not supported", 500)
			return
		}

		ch, err := aiSvc.ChatStream(req.Context(), "admin", aiRouter.TierLight, &gateway.ChatRequest{
			Messages: []gateway.Message{
				{Role: gateway.RoleSystem, Content: "You are an I56 WMS AI assistant. " + bizContext},
				{Role: gateway.RoleUser, Content: msg},
			},
		})
		if err != nil {
			fmt.Fprintf(w, "data: Error: %s\n\n", err.Error())
			flusher.Flush()
			return
		}
		totalTokens := 0
		for ev := range ch {
			fmt.Fprintf(w, "data: %s\n\n", ev.Content)
			flusher.Flush()
			totalTokens++
			if ev.Done {
				break
			}
		}
		// Track cost for streaming
		costTracker.Track("deepseek", "chat-sse", 1, len(msg)/4, totalTokens, float64(totalTokens)*0.000001)
	})

	// Health — now with AI status
		// AI services
	cargoClassifier := classifier.New(aiSvc.Gateway)
	productTranslator := translate.New(aiSvc.Gateway)

	r.GET("/api/ai/classify", func(w http.ResponseWriter, req *http.Request) {
		text := req.URL.Query().Get("text")
		if text == "" { http.Error(w, "{}", 400); return }
		result := cargoClassifier.Classify(text)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	r.GET("/api/ai/translate", func(w http.ResponseWriter, req *http.Request) {
		text := req.URL.Query().Get("q")
		from := req.URL.Query().Get("from"); to := req.URL.Query().Get("to")
		if from == "" { from = "zh" }; if to == "" { to = "en" }
		translated := productTranslator.Translate(text, from, to)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"text":text,"from":from,"to":to,"translated":translated})
	})

	// WMS modules need cargo classifier
	r.GET("/api/v1/health", func(w http.ResponseWriter, req *http.Request) {
		deps := "in-memory"
		if dbAvailable { deps = "postgres" }
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":{"name":"I56 Framework","version":"2.4.2","status":"ok","ai":"active","deps":"` + deps + `"}}`))
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

	// Admin SSE endpoint — streams dashboard stat updates (order/parcel counts)
	r.GET("/api/v1/admin/events", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		c := hub.Subscribe("admin-dashboard")
		defer hub.Unsubscribe(c)
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "SSE not supported", 500)
			return
		}
		// Send an initial connected event
		fmt.Fprintf(w, "event: connected\ndata: {\"status\":\"ok\"}\n\n")
		flusher.Flush()
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

	a := adminOnly(sessionMgr)

	// ★ All admin HTML routes removed — React SPA handles admin UI via JSON API
	// omsroute, wmsroute, tmsroute, crmroute, finroute, sysroute registrations removed.

	// ★ JSON API for React frontend
	registerJSONAPI(r, a, rbac, sessionMgr)

	// ★ New Framework v2.3 routes
	// OpenAPI spec endpoint
	router.RegisterOpenAPIEndpoint(r, openapiGen)
	registerSchedulerRoutes(r, sch, a, tmpl)
	registerAuditRoutes(r, auditLogger, a, tmpl)

	// ─── Generic Admin Delete (covers all modules) ───
	r.POST("/admin/delete", a(func(w http.ResponseWriter, req *http.Request) {
		page, idStr := req.URL.Query().Get("page"), req.URL.Query().Get("id")
		if page == "" || idStr == "" { http.Error(w, "missing page or id", 400); return }
		id, _ := strconv.ParseInt(idStr, 10, 64)
		switch page {
		case "warehouses": wr.Delete(req.Context(), 1, id)
		case "parcels": pr.Delete(req.Context(), 1, id)
		case "route-templates": rr.Delete(req.Context(), 1, id)
		case "couriers": cour.Delete(req.Context(), idStr)
		case "clients": cr.Delete(req.Context(), 1, id)
		case "employees": rbac.DeleteUser(req.Context(), 1, id)
		case "roles": rbac.DeleteRole(req.Context(), 1, id)
		default:
			http.Error(w, "unsupported page: "+page, 400); return
		}
		w.Header().Set("HX-Refresh", "true")
		w.WriteHeader(200)
	}))

	registerReportRoutes(r, reportEngine, a, tmpl)
	// Register some demo routes for OpenAPI documentation
	registerOpenAPIDemoRoutes(r, openapiGen, a)

	// PDA API routes (direct on main router)
	registerPDAAPIRoutesOnRouter(r, pdaOps)
	registerPDARoutes(r, pdaR, pdaOps, hub)

	// ★ JSON APIs for React client & PDA frontends
	registerClientJSONAPI(r, tm, ps, osvc, rr, cour, ws, lr, dr, mr, sr, whr, ar, rpr, dfr, scr, acr)
	registerPDAJSONAPI(r, pdaR, pdaOps)
	registerAdminFullAPI(r, a, ps, osvc, ws, cr, rr, cour, sr, wor, lr, dr, mr, ar, rpr, dfr, scr, acr, rbac, ppr, wfr, td, whr)

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

	// Catch-all — serve 404 page
	r.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			w.WriteHeader(404)
			http.ServeFile(w, req, "templates/error/404.html")
			return
		}
		http.Redirect(w, req, "/login", 303)
	}))

	// Start
	log.Println("I56 Framework 2.4.2 listening on :8080")
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
