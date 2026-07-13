package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/i56/framework/core/router"
	"github.com/i56/framework/core/sse"
	pdaRepo "github.com/i56/modules/pda/repository"
	pdaSvc "github.com/i56/modules/pda/service"
	parcelDomain "github.com/i56/modules/parcel/domain"
)

func registerPDARoutes(r *router.Router, pdaR *pdaRepo.MemPDARepo, ops *pdaSvc.PDAOperations, hub *sse.Hub) {
	svc := pdaSvc.NewPDAService(pdaR, nil, nil)
	tmpl := initPDATemplates()

	// ---------- helper functions ----------
	parseScanParam := func(req *http.Request) string {
		return strings.TrimSpace(req.URL.Query().Get("scan"))
	}
	getOperatorID := func(req *http.Request) int64 {
		ck, _ := req.Cookie("pda_token")
		if ck == nil {
			return 1
		}
		sess := pdaR.ValidateSession(ck.Value)
		if sess == nil {
			return 1
		}
		return sess.OperatorID
	}

	// ---------- PDA auth middleware ----------
	pdaAuth := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ck, err := r.Cookie("pda_token")
			if err != nil || pdaR.ValidateSession(ck.Value) == nil {
				http.Redirect(w, r, "/pda/login", 303)
				return
			}
			next(w, r)
		}
	}

	// ==========================================
	// 0. LOGIN
	// ==========================================
	r.GET("/pda/login", func(w http.ResponseWriter, req *http.Request) {
		tmpl["login"].ExecuteTemplate(w, "login.html", map[string]any{"Active": "login"})
	})
	r.POST("/pda/login", func(w http.ResponseWriter, req *http.Request) {
		code := req.FormValue("code")
		pin := req.FormValue("pin")
		sess, err := svc.Login(code, pin, req.RemoteAddr)
		if err != nil {
			tmpl["login"].ExecuteTemplate(w, "login.html", map[string]any{"Error": err.Error(), "Code": code, "Active": "login"})
			return
		}
		http.SetCookie(w, &http.Cookie{Name: "pda_token", Value: sess.Token, Path: "/pda", HttpOnly: true, MaxAge: 43200})
		http.Redirect(w, req, "/pda", 303)
	})

	r.GET("/pda/logout", func(w http.ResponseWriter, req *http.Request) {
		if ck, err := req.Cookie("pda_token"); err == nil {
			pdaR.Logout(ck.Value)
		}
		http.SetCookie(w, &http.Cookie{Name: "pda_token", Value: "", Path: "/pda", MaxAge: -1})
		http.Redirect(w, req, "/pda/login", 303)
	})

	// ==========================================
	// 1. DASHBOARD (首页)
	// ==========================================
	r.GET("/pda", pdaAuth(func(w http.ResponseWriter, req *http.Request) {
		stats := ops.WarehouseStats(req.Context())
		logs := ops.RecentLogs(10)
		data := map[string]any{
			"Active":    "dashboard",
			"Stats":     stats,
			"OpID":      getOperatorID(req),
			"RecentLogs": logs,
			"Now":       time.Now().Format("15:04"),
		}
		tmpl["dashboard"].ExecuteTemplate(w, "dashboard.html", data)
	}))
	r.GET("/pda/menu", pdaAuth(func(w http.ResponseWriter, req *http.Request) {
		stats := ops.WarehouseStats(req.Context())
		logs := ops.RecentLogs(10)
		data := map[string]any{
			"Active":     "dashboard",
			"Stats":      stats,
			"OpID":       getOperatorID(req),
			"RecentLogs": logs,
			"Now":        time.Now().Format("15:04"),
		}
		tmpl["dashboard"].ExecuteTemplate(w, "dashboard.html", data)
	}))

	// ==========================================
	// 1b. TASK POOL (抢单池) — SSE + Page
	// ==========================================
	// SSE stream for real-time task pool updates
	r.GET("/pda/sse/task-pool", func(w http.ResponseWriter, req *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming unsupported", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		client := hub.Subscribe("pda-task-pool")
		defer hub.Unsubscribe(client)

		// Send initial state
		sendTaskPoolSSE(w, flusher, ops, req.Context())

		for {
			select {
			case ev := <-client.Events:
				fmt.Fprintf(w, "event: %s\ndata: %s\n\n", ev.Type, ev.Data)
				flusher.Flush()
			case <-req.Context().Done():
				return
			}
		}
	})

	// Full task pool page
	r.GET("/pda/task-pool", pdaAuth(func(w http.ResponseWriter, req *http.Request) {
		stats := ops.WarehouseStats(req.Context())
		receiveList := ops.PendingReceive(req.Context())
		pickList := ops.PendingPick(req.Context())
		packList := ops.PendingPack(req.Context())
		loadList := ops.PendingLoad(req.Context())
		putawayList := ops.PendingPutAway(req.Context())

		// Build task pool entries
		type TaskEntry struct {
			ID           string
			TaskType     string
			TaskLabel    string
			Icon         string
			TargetID     string
			TargetDesc   string
			ParcelCount  int
			Priority     string
			CreatedAt    string
		}
		var tasks []TaskEntry
		for i, p := range receiveList {
			tasks = append(tasks, TaskEntry{
				ID: fmt.Sprintf("recv-%d", i+1), TaskType: "receive", TaskLabel: "收货入库",
				Icon: "📦", TargetID: p.TrackingNumber, TargetDesc: p.ProductName,
				Priority: "normal", CreatedAt: time.Now().Add(-time.Duration(i)*time.Minute).Format("15:04"),
			})
		}
		for i, o := range pickList {
			tasks = append(tasks, TaskEntry{
				ID: fmt.Sprintf("pick-%d", i+1), TaskType: "pick", TaskLabel: "订单拣货",
				Icon: "🛒", TargetID: o.OrderNo, TargetDesc: o.RecipientName,
				ParcelCount: o.ParcelCount, Priority: "normal", CreatedAt: time.Now().Add(-time.Duration(i)*2*time.Minute).Format("15:04"),
			})
		}
		for i, o := range packList {
			tasks = append(tasks, TaskEntry{
				ID: fmt.Sprintf("pack-%d", i+1), TaskType: "pack", TaskLabel: "打包复核",
				Icon: "📋", TargetID: o.OrderNo, TargetDesc: o.RecipientName,
				ParcelCount: o.ParcelCount, Priority: "normal", CreatedAt: time.Now().Add(-time.Duration(i)*3*time.Minute).Format("15:04"),
			})
		}
		for i, o := range loadList {
			tasks = append(tasks, TaskEntry{
				ID: fmt.Sprintf("load-%d", i+1), TaskType: "load", TaskLabel: "装柜发货",
				Icon: "🚛", TargetID: o.OrderNo, TargetDesc: o.RecipientName,
				ParcelCount: o.ParcelCount, Priority: "normal", CreatedAt: time.Now().Add(-time.Duration(i)*5*time.Minute).Format("15:04"),
			})
		}
		for i, p := range putawayList {
			tasks = append(tasks, TaskEntry{
				ID: fmt.Sprintf("put-%d", i+1), TaskType: "putaway", TaskLabel: "上架入库",
				Icon: "📍", TargetID: p.TrackingNumber, TargetDesc: p.ProductName,
				Priority: "normal", CreatedAt: time.Now().Add(-time.Duration(i)*time.Minute).Format("15:04"),
			})
		}

		data := map[string]any{
			"Active":    "task_pool",
			"Tasks":     tasks,
			"Stats":     stats,
			"PoolCount": len(tasks),
		}
		tmpl["task_pool"].ExecuteTemplate(w, "task_pool.html", data)
	}))

	// POST /pda/task-pool - claim a task
	r.POST("/pda/task-pool", pdaAuth(func(w http.ResponseWriter, req *http.Request) {
		taskType := req.FormValue("task_type")
		targetID := req.FormValue("target_id")
		opID := getOperatorID(req)

		// Assign task to this operator (using the force-assign mechanism)
		var targetID64 int64
		fmt.Sscanf(targetID, "%d", &targetID64)
		if targetID64 == 0 {
			targetID64 = 1
		}
		at, err := ops.ForceAssign(req.Context(), opID, taskType, targetID64, opID)
		if err != nil {
			http.Redirect(w, req, "/pda/task-pool?error="+err.Error(), 303)
			return
		}
		_ = at

		// Broadcast updated task pool to all SSE clients
		broadcastTaskPoolUpdate(hub, ops, req.Context())

		// Redirect to the appropriate operation screen
		switch taskType {
		case "receive":
			http.Redirect(w, req, "/pda/receive?scan="+targetID, 303)
		case "putaway":
			http.Redirect(w, req, "/pda/putaway?scan="+targetID, 303)
		case "pick":
			http.Redirect(w, req, "/pda/pick?scan="+targetID, 303)
		case "pack":
			http.Redirect(w, req, "/pda/pack?scan="+targetID, 303)
		case "load":
			http.Redirect(w, req, "/pda/load?scan="+targetID, 303)
		default:
			http.Redirect(w, req, "/pda/my-tasks", 303)
		}
	}))

	// ==========================================
	// 1c. MY TASKS (我的任务)
	// ==========================================
	r.GET("/pda/my-tasks", pdaAuth(func(w http.ResponseWriter, req *http.Request) {
		opID := getOperatorID(req)
		assigned := ops.GetAssignedTasks(opID)
		stats := ops.WarehouseStats(req.Context())
		logs := ops.RecentLogs(20)

		data := map[string]any{
			"Active":       "my_tasks",
			"AssignedTasks": assigned,
			"Stats":        stats,
			"OpID":         opID,
			"RecentLogs":   logs,
		}
		tmpl["my_tasks"].ExecuteTemplate(w, "my_tasks.html", data)
	}))

	// ==========================================
	// 3. RECEIVE (收货)
	// ==========================================
	r.GET("/pda/receive", pdaAuth(func(w http.ResponseWriter, req *http.Request) {
		data := map[string]any{
			"Active":  "receive",
			"Pending": ops.PendingReceive(req.Context()),
		}
		if scan := parseScanParam(req); scan != "" {
			if p, err := ops.GetParcelForReceive(req.Context(), scan); err == nil && p != nil {
				data["ScannedParcel"] = p
			}
		}
		tmpl["receive"].ExecuteTemplate(w, "receive.html", data)
	}))
	r.POST("/pda/receive", pdaAuth(func(w http.ResponseWriter, req *http.Request) {
		tn := req.FormValue("tracking_no")
		wVal, _ := strconv.ParseFloat(req.FormValue("weight"), 64)
		lVal, _ := strconv.ParseFloat(req.FormValue("length"), 64)
		wiVal, _ := strconv.ParseFloat(req.FormValue("width"), 64)
		hVal, _ := strconv.ParseFloat(req.FormValue("height"), 64)
		loc := req.FormValue("location_barcode")
		if wVal <= 0 {
			wVal = 0.5
		}
		if lVal <= 0 {
			lVal = 20
		}
		if wiVal <= 0 {
			wiVal = 15
		}
		if hVal <= 0 {
			hVal = 10
		}
		opID := getOperatorID(req)
		p, err := ops.Receive(context.Background(), opID, tn, wVal, lVal, wiVal, hVal, loc)
		if err != nil {
			tmpl["result"].ExecuteTemplate(w, "result.html", map[string]any{
				"Active": "receive", "Error": err.Error(), "BackUrl": "/pda/receive",
			})
			return
		}
		tmpl["result"].ExecuteTemplate(w, "result.html", map[string]any{
			"Active":        "receive",
			"Message":       "收货成功",
			"ProductName":   p.ProductName,
			"TrackingNo":    p.TrackingNumber,
			"Weight":        p.ActualWeight,
			"Length":        p.Length, "Width": p.Width, "Height": p.Height,
			"Status":        string(p.Status),
			"BackUrl":       "/pda/receive",
			"NextStepUrl":   "/pda/weigh?scan=" + p.TrackingNumber,
			"NextStepLabel": "去核重",
		})
	}))

	// ==========================================
	// 4. WEIGH (称重核重)
	// ==========================================
	r.GET("/pda/weigh", pdaAuth(func(w http.ResponseWriter, req *http.Request) {
		data := map[string]any{
			"Active":       "weigh",
			"PendingWeigh": ops.PendingWeigh(req.Context()),
		}
		if scan := parseScanParam(req); scan != "" {
			if p, _, _ := ops.GetParcelForWeigh(req.Context(), scan); p != nil {
				data["ScannedParcel"] = p
			}
		}
		tmpl["weigh"].ExecuteTemplate(w, "weigh.html", data)
	}))
	r.POST("/pda/weigh", pdaAuth(func(w http.ResponseWriter, req *http.Request) {
		tn := req.FormValue("tracking_number")
		wVal, _ := strconv.ParseFloat(req.FormValue("weight_kg"), 64)
		if wVal <= 0 {
			wVal = 0.1
		}
		opID := getOperatorID(req)
		p, oldW, err := ops.Weigh(context.Background(), opID, tn, wVal)
		if err != nil {
			tmpl["result"].ExecuteTemplate(w, "result.html", map[string]any{
				"Active": "weigh", "Error": err.Error(), "BackUrl": "/pda/weigh",
			})
			return
		}
		tmpl["result"].ExecuteTemplate(w, "result.html", map[string]any{
			"Active":        "weigh",
			"Message":       fmt.Sprintf("核重完成 %.2f→%.2f kg", oldW, wVal),
			"ProductName":   p.ProductName,
			"TrackingNo":    p.TrackingNumber,
			"Weight":        p.ActualWeight,
			"Length":        p.Length, "Width": p.Width, "Height": p.Height,
			"Status":        string(p.Status),
			"BackUrl":       "/pda/weigh",
			"NextStepUrl":   "/pda/putaway?scan=" + p.TrackingNumber,
			"NextStepLabel": "去上架",
		})
	}))

	// ==========================================
	// 5. PUTAWAY (上架)
	// ==========================================
	r.GET("/pda/putaway", pdaAuth(func(w http.ResponseWriter, req *http.Request) {
		data := map[string]any{
			"Active":         "putaway",
			"PendingPutAway": ops.PendingPutAway(req.Context()),
		}
		if scan := parseScanParam(req); scan != "" {
			if p, _, _ := ops.GetParcelForPutAway(req.Context(), scan); p != nil {
				data["ScannedParcel"] = p
			}
		}
		tmpl["putaway"].ExecuteTemplate(w, "putaway.html", data)
	}))
	r.POST("/pda/putaway", pdaAuth(func(w http.ResponseWriter, req *http.Request) {
		tn := req.FormValue("tracking_no")
		loc := req.FormValue("location_barcode")
		opID := getOperatorID(req)
		p, err := ops.PutAway(context.Background(), opID, tn, loc)
		if err != nil {
			tmpl["result"].ExecuteTemplate(w, "result.html", map[string]any{
				"Active": "putaway", "Error": err.Error(), "BackUrl": "/pda/putaway",
			})
			return
		}
		tmpl["result"].ExecuteTemplate(w, "result.html", map[string]any{
			"Active":        "putaway",
			"Message":       "上架成功",
			"ProductName":   p.ProductName,
			"TrackingNo":    p.TrackingNumber,
			"Weight":        p.ActualWeight,
			"Length":        p.Length, "Width": p.Width, "Height": p.Height,
			"Status":        string(p.Status),
			"BackUrl":       "/pda/putaway",
			"NextStepUrl":   "/pda",
			"NextStepLabel": "返回首页",
		})
	}))

	// ==========================================
	// 6. PICK (拣货)
	// ==========================================
	r.GET("/pda/pick", pdaAuth(func(w http.ResponseWriter, req *http.Request) {
		data := map[string]any{
			"Active":            "pick",
			"PendingPickOrders": ops.PendingPick(req.Context()),
		}
		if scan := parseScanParam(req); scan != "" {
			if order, parcels, err := ops.LookupOrder(req.Context(), scan); err == nil && order != nil {
				data["ScannedOrder"] = order
				data["ScannedOrderParcels"] = parcels
			}
		}
		tmpl["pick"].ExecuteTemplate(w, "pick.html", data)
	}))
	r.POST("/pda/pick", pdaAuth(func(w http.ResponseWriter, req *http.Request) {
		orderNo := req.FormValue("order_no")
		opID := getOperatorID(req)
		order, parcels, err := ops.Pick(context.Background(), opID, orderNo)
		if err != nil {
			tmpl["result"].ExecuteTemplate(w, "result.html", map[string]any{
				"Active": "pick", "Error": err.Error(), "BackUrl": "/pda/pick",
			})
			return
		}
		tmpl["result"].ExecuteTemplate(w, "result.html", map[string]any{
			"Active":        "pick",
			"Message":       "拣货完成",
			"OrderNo":       order.OrderNo,
			"Parcels":       parcels,
			"Count":         len(parcels),
			"Status":        "picked",
			"BackUrl":       "/pda/pick",
			"NextStepUrl":   "/pda/pack?scan=" + order.OrderNo,
			"NextStepLabel": "去打包",
		})
	}))

	// ==========================================
	// 7. PACK (打包)
	// ==========================================
	r.GET("/pda/pack", pdaAuth(func(w http.ResponseWriter, req *http.Request) {
		data := map[string]any{
			"Active":           "pack",
			"PendingPackOrders": ops.PendingPack(req.Context()),
		}
		if scan := parseScanParam(req); scan != "" {
			if order, parcels, err := ops.LookupOrderForPack(req.Context(), scan); err == nil && order != nil {
				data["ScannedOrder"] = order
				data["ScannedOrderParcels"] = parcels
			}
		}
		tmpl["pack"].ExecuteTemplate(w, "pack.html", data)
	}))
	r.POST("/pda/pack", pdaAuth(func(w http.ResponseWriter, req *http.Request) {
		orderNo := req.FormValue("order_no")
		opID := getOperatorID(req)
		order, err := ops.Pack(context.Background(), opID, orderNo)
		if err != nil {
			tmpl["result"].ExecuteTemplate(w, "result.html", map[string]any{
				"Active": "pack", "Error": err.Error(), "BackUrl": "/pda/pack",
			})
			return
		}
		parcels := ops.GetOrderParcels(context.Background(), order)
		tmpl["result"].ExecuteTemplate(w, "result.html", map[string]any{
			"Active":        "pack",
			"Message":       "打包完成",
			"OrderNo":       order.OrderNo,
			"Parcels":       parcels,
			"Count":         len(parcels),
			"Status":        "packed",
			"BackUrl":       "/pda/pack",
			"NextStepUrl":   "/pda/load?scan=" + order.OrderNo,
			"NextStepLabel": "去装柜",
		})
	}))

	// ==========================================
	// 8. LOAD (装柜)
	// ==========================================
	r.GET("/pda/load", pdaAuth(func(w http.ResponseWriter, req *http.Request) {
		data := map[string]any{
			"Active":      "load",
			"PendingLoad": ops.PendingLoad(req.Context()),
		}
		if scan := parseScanParam(req); scan != "" {
			if order, parcels, err := ops.LookupOrder(req.Context(), scan); err == nil && order != nil {
				data["ScannedOrder"] = order
				data["ScannedOrderParcels"] = parcels
			}
		}
		tmpl["load"].ExecuteTemplate(w, "load.html", data)
	}))
	r.POST("/pda/load", pdaAuth(func(w http.ResponseWriter, req *http.Request) {
		orderNo := req.FormValue("order_no")
		containerNo := req.FormValue("container_no")
		opID := getOperatorID(req)
		if containerNo == "" {
			containerNo = "CONT-001"
		}
		err := ops.LoadContainer(context.Background(), opID, containerNo, orderNo)
		if err != nil {
			tmpl["result"].ExecuteTemplate(w, "result.html", map[string]any{
				"Active": "load", "Error": err.Error(), "BackUrl": "/pda/load",
			})
			return
		}
		tmpl["result"].ExecuteTemplate(w, "result.html", map[string]any{
			"Active":        "load",
			"Message":       "装柜完成",
			"OrderNo":       orderNo,
			"ContainerNo":   containerNo,
			"Status":        "loaded",
			"BackUrl":       "/pda/load",
			"NextStepUrl":   "/pda",
			"NextStepLabel": "返回首页",
		})
	}))

	// ==========================================
	// 9. EXCEPTION (标异常)
	// ==========================================
	r.GET("/pda/exception", pdaAuth(func(w http.ResponseWriter, req *http.Request) {
		data := map[string]any{"Active": "exception"}
		if scan := parseScanParam(req); scan != "" {
			if p, err := ops.GetParcelForReceive(req.Context(), scan); err == nil && p != nil {
				data["ScannedParcel"] = p
			}
		}
		tmpl["exception"].ExecuteTemplate(w, "exception.html", data)
	}))
	r.POST("/pda/exception", pdaAuth(func(w http.ResponseWriter, req *http.Request) {
		tn := req.FormValue("tracking_no")
		reasonType := req.FormValue("reason_type")
		note := req.FormValue("note")
		reason := reasonType
		if note != "" {
			reason = reasonType + ": " + note
		}
		opID := getOperatorID(req)
		err := ops.MarkException(context.Background(), opID, tn, reason)
		if err != nil {
			tmpl["result"].ExecuteTemplate(w, "result.html", map[string]any{
				"Active": "exception", "Error": err.Error(), "BackUrl": "/pda/exception",
			})
			return
		}
		tmpl["result"].ExecuteTemplate(w, "result.html", map[string]any{
			"Active":        "exception",
			"Message":       "已标记异常",
			"TrackingNo":    tn,
			"Status":        "abnormal",
			"BackUrl":       "/pda/exception",
			"NextStepUrl":   "/pda/exception",
			"NextStepLabel": "继续标记",
		})
	}))

	// ==========================================
	// 10. QUERY (查询)
	// ==========================================
	r.GET("/pda/query", pdaAuth(func(w http.ResponseWriter, req *http.Request) {
		data := map[string]any{"Active": "query"}
		if scan := parseScanParam(req); scan != "" {
			p, logs, err := svc.QueryParcel(context.Background(), "", scan)
			if err == nil && p != nil {
				data["Parcel"] = p
				data["ScanHistory"] = logs
			}
		}
		tmpl["query"].ExecuteTemplate(w, "query.html", data)
	}))
	r.POST("/pda/query", pdaAuth(func(w http.ResponseWriter, req *http.Request) {
		tn := req.FormValue("tracking_no")
		p, logs, err := svc.QueryParcel(context.Background(), "", tn)
		if err != nil || p == nil {
			tmpl["query"].ExecuteTemplate(w, "query.html", map[string]any{
				"Active": "query", "Error": fmt.Sprintf("包裹 %s 未找到", tn),
			})
			return
		}
		tmpl["query"].ExecuteTemplate(w, "query.html", map[string]any{
			"Active":      "query",
			"Parcel":      p,
			"ScanHistory": logs,
		})
	}))
}

// ---------- Template initialization with custom functions ----------
// ──────────────────────────────────────────────────────────────────────
// SSE task pool helpers
// ──────────────────────────────────────────────────────────────────────

// TaskEntry is a task pool item for SSE serialization
type TaskEntry struct {
	ID          string `json:"id"`
	TaskType    string `json:"task_type"`
	TaskLabel   string `json:"task_label"`
	Icon        string `json:"icon"`
	TargetID    string `json:"target_id"`
	TargetDesc  string `json:"target_desc"`
	ParcelCount int    `json:"parcel_count"`
	Priority    string `json:"priority"`
	CreatedAt   string `json:"created_at"`
}

func buildTaskEntries(ops *pdaSvc.PDAOperations, ctx context.Context) []TaskEntry {
	receiveList := ops.PendingReceive(ctx)
	pickList := ops.PendingPick(ctx)
	packList := ops.PendingPack(ctx)
	loadList := ops.PendingLoad(ctx)
	putawayList := ops.PendingPutAway(ctx)

	var tasks []TaskEntry
	for i, p := range receiveList {
		tasks = append(tasks, TaskEntry{
			ID: fmt.Sprintf("recv-%d", i+1), TaskType: "receive", TaskLabel: "收货入库",
			Icon: "📦", TargetID: p.TrackingNumber, TargetDesc: p.ProductName,
			Priority: "normal", CreatedAt: time.Now().Add(-time.Duration(i) * time.Minute).Format("15:04"),
		})
	}
	for i, o := range pickList {
		tasks = append(tasks, TaskEntry{
			ID: fmt.Sprintf("pick-%d", i+1), TaskType: "pick", TaskLabel: "订单拣货",
			Icon: "🛒", TargetID: o.OrderNo, TargetDesc: o.RecipientName,
			ParcelCount: o.ParcelCount, Priority: "normal", CreatedAt: time.Now().Add(-time.Duration(i) * 2 * time.Minute).Format("15:04"),
		})
	}
	for i, o := range packList {
		tasks = append(tasks, TaskEntry{
			ID: fmt.Sprintf("pack-%d", i+1), TaskType: "pack", TaskLabel: "打包复核",
			Icon: "📋", TargetID: o.OrderNo, TargetDesc: o.RecipientName,
			ParcelCount: o.ParcelCount, Priority: "normal", CreatedAt: time.Now().Add(-time.Duration(i) * 3 * time.Minute).Format("15:04"),
		})
	}
	for i, o := range loadList {
		tasks = append(tasks, TaskEntry{
			ID: fmt.Sprintf("load-%d", i+1), TaskType: "load", TaskLabel: "装柜发货",
			Icon: "🚛", TargetID: o.OrderNo, TargetDesc: o.RecipientName,
			ParcelCount: o.ParcelCount, Priority: "normal", CreatedAt: time.Now().Add(-time.Duration(i) * 5 * time.Minute).Format("15:04"),
		})
	}
	for i, p := range putawayList {
		tasks = append(tasks, TaskEntry{
			ID: fmt.Sprintf("put-%d", i+1), TaskType: "putaway", TaskLabel: "上架入库",
			Icon: "📍", TargetID: p.TrackingNumber, TargetDesc: p.ProductName,
			Priority: "normal", CreatedAt: time.Now().Add(-time.Duration(i) * time.Minute).Format("15:04"),
		})
	}
	return tasks
}

func buildTaskPoolFragmentHTML(tasks []TaskEntry) string {
	var buf bytes.Buffer
	if len(tasks) == 0 {
		buf.WriteString(`<div class="pda-empty"><div class="icon">🎉</div><div style="font-size:15px;font-weight:600;margin-bottom:6px">暂无待抢任务</div><div class="text-sm">所有任务已分配，请稍后再来</div></div>`)
		return buf.String()
	}
	buf.WriteString(fmt.Sprintf(`<div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:10px"><div class="text-sm text-muted">共 %d 个待抢任务</div></div>`, len(tasks)))
	for _, t := range tasks {
		iconBg := "rgba(59,130,246,.12)"
		switch t.TaskType {
		case "receive":
			iconBg = "rgba(14,165,233,.12)"
		case "putaway":
			iconBg = "rgba(16,185,129,.12)"
		case "pick":
			iconBg = "rgba(139,92,246,.12)"
		case "pack":
			iconBg = "rgba(249,115,22,.12)"
		}
		desc := t.TargetID
		if t.TargetDesc != "" {
			desc += " · " + t.TargetDesc
		}
		if t.ParcelCount > 0 {
			desc += fmt.Sprintf(" · %d件", t.ParcelCount)
		}
		buf.WriteString(fmt.Sprintf(
			`<div class="pda-task-item w-full" style="border:none;background:var(--pda-surface);width:100%%;text-align:left;cursor:pointer;font-family:var(--pda-font);margin-bottom:8px" onclick="claimTask('%s','%s',this)"><div class="task-icon" style="background:%s">%s</div><div class="task-info"><div class="task-type">%s</div><div class="task-desc">%s</div><div style="font-size:10px;color:var(--pda-text-muted);margin-top:4px">⏱ %s</div></div><div class="task-arrow">抢单 ›</div></div>`,
			t.TaskType, t.TargetID, iconBg, t.Icon, t.TaskLabel, desc, t.CreatedAt,
		))
	}
	return buf.String()
}

func sendTaskPoolSSE(w http.ResponseWriter, flusher http.Flusher, ops *pdaSvc.PDAOperations, ctx context.Context) {
	tasks := buildTaskEntries(ops, ctx)
	fragment := buildTaskPoolFragmentHTML(tasks)
	fmt.Fprintf(w, "event: taskPoolUpdate\ndata: %s\n\n", fragment)
	flusher.Flush()
}

func broadcastTaskPoolUpdate(hub *sse.Hub, ops *pdaSvc.PDAOperations, ctx context.Context) {
	tasks := buildTaskEntries(ops, ctx)
	fragment := buildTaskPoolFragmentHTML(tasks)
	hub.Publish("pda-task-pool", sse.Event{
		Type: "taskPoolUpdate",
		Data: fragment,
	})
}

func initPDATemplates() map[string]*template.Template {
	funcs := template.FuncMap{
		"statusDisplay": parcelDomain.StatusDisplay,
		"or": func(a, b string) string {
			if a != "" {
				return a
			}
			return b
		},
		"volWeight": func(p *parcelDomain.Parcel) float64 {
			if p == nil {
				return 0
			}
			return p.VolumetricWeight()
		},
		"chgWeight": func(p *parcelDomain.Parcel) float64 {
			if p == nil {
				return 0
			}
			return p.ChargeableWeight()
		},
		"totalWeight": func(parcels []parcelDomain.Parcel) string {
			var total float64
			for _, p := range parcels {
				total += p.ActualWeight
			}
			return fmt.Sprintf("%.2f", total)
		},
		"string": func(s parcelDomain.ParcelStatus) string { return string(s) },
		"parcelAfterStep": func(status parcelDomain.ParcelStatus, step string) bool {
			order := []parcelDomain.ParcelStatus{
				parcelDomain.StatusPreDeclared,
				parcelDomain.StatusReceived,
				parcelDomain.StatusWeighed,
				parcelDomain.StatusStored,
				parcelDomain.StatusPicked,
				parcelDomain.StatusPacked,
				parcelDomain.StatusShipped,
			}
			stepIdx := -1
			statusIdx := -1
			for i, s := range order {
				if string(s) == step {
					stepIdx = i
				}
				if s == status || string(s) == string(status) {
					statusIdx = i
				}
			}
			return statusIdx > stepIdx
		},
	}
	tmpl := map[string]*template.Template{}
	templateKeys := []string{
		"login", "dashboard", "task_pool", "my_tasks",
		"receive", "weigh", "putaway", "pick", "pack",
		"load", "exception", "query", "result",
	}
	for _, p := range templateKeys {
		t := template.New(p).Funcs(funcs)
		t = template.Must(t.ParseFiles("templates/pda/base.html", "templates/pda/camera_scanner.html", "templates/pda/"+p+".html"))
		tmpl[p] = t
	}
	return tmpl
}
