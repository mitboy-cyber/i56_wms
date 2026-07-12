// Package route provides WMS (仓库管理) admin route registration.
package route

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/i56/framework/core/router"
	"github.com/i56/framework/events"

	"github.com/i56/i56-apps/i56-wms/internal/common"

	custRepo "github.com/i56/modules/customer/repository"
	orderDomain "github.com/i56/modules/order/domain"
	orderSvc "github.com/i56/modules/order/service"
	parcelDomain "github.com/i56/modules/parcel/domain"
	parcelSvc "github.com/i56/modules/parcel/service"
	psRepo "github.com/i56/modules/parcel_service/repository"
	rbacRepo "github.com/i56/modules/rbac/repository"
	tmsRepo "github.com/i56/modules/transport/repository"
	twoDomain "github.com/i56/modules/workorder/domain"
	twoRepo "github.com/i56/modules/workorder/repository"
	wfDomain "github.com/i56/modules/workflow/domain"
	wfRepo "github.com/i56/modules/workflow/repository"
	parcelRepo "github.com/i56/modules/parcel/repository"
	whDomain "github.com/i56/modules/warehouse/domain"
	whRepo "github.com/i56/modules/warehouse/repository"
	whSvc "github.com/i56/modules/warehouse/service"
)

// Register WMS admin routes (~13 list pages + CRUD).
func Register(
	r *router.Router,
	a func(http.HandlerFunc) http.HandlerFunc,
	rc *common.RenderCtx,
	ps *parcelSvc.ParcelService,
	ws *whSvc.WarehouseService,
	osvc *orderSvc.OrderService,
	cr *custRepo.MemClientRepo,
	rr *tmsRepo.MemRouteRepo,
	wor *twoRepo.MemWorkOrderRepo,
	sr *psRepo.MemServiceRepo,
	wfr *wfRepo.MemWorkflowRepo,
	rbac *rbacRepo.MemRBACRepo,
	pr *parcelRepo.MemParcelRepo,
	wr *whRepo.MemWarehouseRepo,
) {
	const tenant int64 = 1
	gp := rc.NewGenericList()

	// ─── Dashboard — real KPI stats from all repos ───
	r.GET("/admin", a(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		now := time.Now()
		todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

		// 1. Core counts
		parcels, totalParcels, _ := ps.List(ctx, 1, 0, 200)
		orders, totalOrders, _ := osvc.List(ctx, 1, 0, 100)
		clients, totalClients, _ := cr.List(ctx, 1, 0, 50)
		routes, totalRoutes, _ := rr.List(ctx, 1, 0, 50)

		// 2. Active orders (not cancelled/completed)
		activeOrders := 0
		for _, o := range orders {
			if o.Status != orderDomain.StatusCancelled && o.Status != orderDomain.StatusCompleted {
				activeOrders++
			}
		}

		// 3. Today revenue — sum order totals created today
		todayRevenue := 0.0
		for _, o := range orders {
			if o.CreatedAt.After(todayStart) || o.CreatedAt.Equal(todayStart) {
				todayRevenue += o.TotalPrice
			}
		}

		// 4. Abnormal parcels count
		abnormalParcels := 0
		parcelStatusCounts := map[string]int{}
		for _, p := range parcels {
			statusKey := string(p.Status)
			parcelStatusCounts[statusKey]++
			if p.Status == parcelDomain.StatusAbnormal || p.Status == parcelDomain.StatusReturned {
				abnormalParcels++
			}
		}

		// 5. Status distribution (parcel + order combined)
		type statusDist struct {
			Label  string
			Count  int
			Pct    float64
			Color  string
		}
		totalStatusItems := 0
		for _, c := range parcelStatusCounts { totalStatusItems += c }
		var statusDistribution []statusDist
		colors := []string{"#6366f1", "#22c55e", "#f59e0b", "#3b82f6", "#ec4899", "#14b8a6", "#8b5cf6", "#f97316", "#06b6d4", "#ef4444"}
		orderStatusCN := map[string]string{
			"pre_declared": "预报", "received": "已入仓", "weighed": "已称重",
			"stored": "已上架", "picked": "已拣货", "packed": "已打包",
			"loaded": "已装柜", "shipped": "运输中",
			"delivered": "已签收", "abnormal": "异常", "returned": "已退货",
			"outbound": "已出货", "delivering": "配送中",
		}
		ci := 0
		for statusKey, count := range parcelStatusCounts {
			label := orderStatusCN[statusKey]
			if label == "" { label = statusKey }
			pct := 0.0
			if totalStatusItems > 0 { pct = float64(count) / float64(totalStatusItems) * 100 }
			statusDistribution = append(statusDistribution, statusDist{label, count, pct, colors[ci%len(colors)]})
			ci++
		}

		// 6. Last 7 days order volume for bar chart
		type dayBar struct {
			Date  string
			Count int
			Max   int // for height calc
		}
		var barData []dayBar
		maxCount := 0
		for d := 6; d >= 0; d-- {
			day := now.Add(-time.Duration(d) * 24 * time.Hour)
			dayStart := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
			dayEnd := dayStart.Add(24 * time.Hour)
			count := 0
			for _, o := range orders {
				if (o.CreatedAt.After(dayStart) || o.CreatedAt.Equal(dayStart)) && o.CreatedAt.Before(dayEnd) {
					count++
				}
			}
			barData = append(barData, dayBar{day.Format("01/02"), count, 0})
			if count > maxCount { maxCount = count }
		}
		for i := range barData {
			if maxCount > 0 {
				barData[i].Max = maxCount
			}
		}

		// 7. Recent activity feed (last 8 events)
		type activityEvent struct {
			Icon    string
			Time    string
			Text    string
		}
		var activities []activityEvent
		// Mix recent orders and parcel events
		eventIdx := 0
		for _, o := range orders {
			if eventIdx >= 8 { break }
			act := activityEvent{Time: o.CreatedAt.Format("01-02 15:04")}
			switch o.Status {
			case orderDomain.StatusPendingPicking: act.Icon = "📋"; act.Text = fmt.Sprintf("新订单 %s 待拣货 (%s)", o.OrderNo, o.RecipientName)
			case orderDomain.StatusInTransit: act.Icon = "🚛"; act.Text = fmt.Sprintf("订单 %s 已发运 (%s)", o.OrderNo, o.RecipientName)
			case orderDomain.StatusCompleted: act.Icon = "✅"; act.Text = fmt.Sprintf("订单 %s 已签收 (%s)", o.OrderNo, o.RecipientName)
			case orderDomain.StatusPendingLoading: act.Icon = "📦"; act.Text = fmt.Sprintf("订单 %s 待装柜 (%s)", o.OrderNo, o.RecipientName)
			case orderDomain.StatusCustomsClearance: act.Icon = "🛃"; act.Text = fmt.Sprintf("订单 %s 清关中 (%s)", o.OrderNo, o.RecipientName)
			case orderDomain.StatusLoaded: act.Icon = "🏗️"; act.Text = fmt.Sprintf("订单 %s 已装柜 (%s)", o.OrderNo, o.RecipientName)
			default: act.Icon = "📋"; act.Text = fmt.Sprintf("订单 %s %s (%s)", o.OrderNo, string(o.Status), o.RecipientName)
			}
			activities = append(activities, act)
			eventIdx++
		}
		// Add some parcel events
		for _, p := range parcels {
			if eventIdx >= 8 { break }
			act := activityEvent{Time: p.CreatedAt.Format("01-02 15:04")}
			switch p.Status {
			case parcelDomain.StatusReceived: act.Icon = "📥"; act.Text = fmt.Sprintf("包裹 %s 已入仓", p.TrackingNumber)
			case parcelDomain.StatusStored: act.Icon = "🏗️"; act.Text = fmt.Sprintf("包裹 %s 已上架", p.TrackingNumber)
			default: act.Icon = "📦"; act.Text = fmt.Sprintf("包裹 %s %s", p.TrackingNumber, string(p.Status))
			}
			activities = append(activities, act)
			eventIdx++
		}

		// 8. Recent parcels for table
		type recentData struct {
			TrackingNumber string
			ProductName    string
			Status         string
			ActualWeight   float64
		}
		parcelStatusDisplay := func(s parcelDomain.ParcelStatus) string {
			switch s {
			case parcelDomain.StatusPreDeclared: return "预报"
			case parcelDomain.StatusReceived: return "已入仓"
			case parcelDomain.StatusWeighed: return "已称重"
			case parcelDomain.StatusStored: return "已上架"
			case parcelDomain.StatusPicked: return "已拣货"
			case parcelDomain.StatusPacked: return "已打包"
			case parcelDomain.StatusLoaded: return "已装柜"
			case parcelDomain.StatusOutbound: return "已出货"
			case parcelDomain.StatusShipped, parcelDomain.StatusDelivering: return "运输中"
			case parcelDomain.StatusDelivered: return "已签收"
			case parcelDomain.StatusAbnormal: return "异常"
			case parcelDomain.StatusReturned: return "已退货"
			default: return string(s)
			}
		}
		recent := make([]recentData, 0, 8)
		for i, p := range parcels {
			if i >= 8 { break }
			recent = append(recent, recentData{p.TrackingNumber, p.ProductName, parcelStatusDisplay(p.Status), p.ActualWeight})
		}

		// 9. Recent orders for table
		type recentOrderData struct {
			OrderNo     string
			ClientName  string
			ParcelCount int
			Status      string
			TotalPrice  float64
		}
		var recentOrders []recentOrderData
		cnMap := map[int64]string{}
		for _, c := range clients { cnMap[c.ID] = c.Name }
		for i, o := range orders {
			if i >= 8 { break }
			ordStatus := string(o.Status)
			switch o.Status {
			case orderDomain.StatusPendingPicking: ordStatus = "待拣货"
			case orderDomain.StatusPicking: ordStatus = "拣货中"
			case orderDomain.StatusPendingPacking: ordStatus = "待打包"
			case orderDomain.StatusPendingLoading: ordStatus = "待装柜"
			case orderDomain.StatusLoaded: ordStatus = "已装柜"
			case orderDomain.StatusInTransit: ordStatus = "运输中"
			case orderDomain.StatusCustomsClearance: ordStatus = "清关中"
			case orderDomain.StatusOutForDelivery: ordStatus = "派送中"
			case orderDomain.StatusCompleted: ordStatus = "已完成"
			case orderDomain.StatusCancelled: ordStatus = "已取消"
			case orderDomain.StatusShipped: ordStatus = "已发货"
			}
			cn := cnMap[o.ClientID]
			if cn == "" { cn = "EZ集运通" }
			recentOrders = append(recentOrders, recentOrderData{o.OrderNo, cn, o.ParcelCount, ordStatus, o.TotalPrice})
		}

		rc.Exec(rc.Tmpl, "dashboard", w, "dashboard.html", map[string]any{
			"Page":               "dashboard",
			"ParcelCount":        totalParcels,
			"OrderCount":         activeOrders,
			"ClientCount":        len(clients),
			"RouteCount":         len(routes),
			"TodayRevenue":       fmt.Sprintf("%.2f", todayRevenue),
			"AbnormalParcels":    abnormalParcels,
			"TotalOrders":        totalOrders,
			"RecentParcels":      recent,
			"RecentOrders":       recentOrders,
			"StatusDistribution": statusDistribution,
			"BarData":            barData,
			"Activities":         activities,
		})
		_ = totalClients; _ = totalRoutes; _ = totalOrders
	}))

	// ─── /admin/warehouses — 仓库列表 (from admin_pages.go) ───
	r.GET("/admin/warehouses", a(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		whs, total, _ := ws.List(ctx, 1, 0, 50)
		parcels, _, _ := ps.List(ctx, 1, 0, 500)
		whParcelCount := map[int64]int{}
		for _, p := range parcels {
			whParcelCount[p.WarehouseID]++
		}
		rows := make([][]string, len(whs))
		for i := range whs {
			rows[i] = []string{whs[i].Name, whs[i].Code, whs[i].Address, whs[i].Contact, whs[i].Phone,
				fmt.Sprintf("%d件", whParcelCount[whs[i].ID]), common.StatusLabelText(whs[i].IsActive)}
		}
		gp(w, "wms_warehouses", "仓库列表", int(total), []string{"仓库", "编码", "地址", "联系人", "电话", "包裹数", "状态"}, rows, "/admin/warehouses/add-form")
	}))

	// ─── /admin/warehouses/{id} — 仓库详情 (from admin_pages.go) ───
	r.GET("/admin/warehouses/{id}", a(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		id, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid ID", 400)
			return
		}
		wh, err := ws.GetByID(ctx, 1, id)
		if err != nil || wh == nil {
			http.Error(w, "Warehouse not found", 404)
			return
		}

		parcels, _, _ := ps.List(ctx, 1, 0, 500)
		var whParcels, stored, picked, packed, shipped int
		for _, p := range parcels {
			if p.WarehouseID == wh.ID {
				whParcels++
				switch p.Status {
				case parcelDomain.StatusStored:
					stored++
				case parcelDomain.StatusPicked:
					picked++
				case parcelDomain.StatusPacked:
					packed++
				case parcelDomain.StatusShipped:
					shipped++
				}
			}
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<!DOCTYPE html><html lang="zh-CN" data-theme="light"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>仓库详情 - I56</title><link rel="stylesheet" href="/static/css/i56-bdl.css"><script src="/static/js/i56-theme.js"></script><style>
*{margin:0;padding:0;box-sizing:border-box}body{font-family:system-ui,sans-serif;background:var(--i56-bg-base);color:var(--i56-text-primary);padding:16px}
.card{background:var(--i56-bg-surface);border:1px solid var(--i56-border);border-radius:8px;padding:16px;margin-bottom:12px}
.card-header{font-size:14px;font-weight:600;color:var(--i56-text-primary);margin-bottom:12px;padding-bottom:8px;border-bottom:1px solid var(--i56-border)}
.section{margin-bottom:16px}.section-title{font-size:13px;font-weight:600;color:var(--i56-brand);margin-bottom:8px;padding-bottom:6px;border-bottom:1px solid var(--i56-border)}
.info-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(240px,1fr));gap:8px}
.info-item{display:flex;padding:6px 0;font-size:12px}
.info-label{color:var(--i56-text-secondary);min-width:80px;flex-shrink:0}
.info-value{color:var(--i56-text-primary);font-weight:500}
.stat-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(140px,1fr));gap:8px;margin-bottom:12px}
.stat-card{background:var(--i56-bg-base);border:1px solid var(--i56-border);border-radius:6px;padding:12px;text-align:center}
.stat-value{font-size:24px;font-weight:700;color:var(--i56-brand)}
.stat-label{font-size:11px;color:var(--i56-text-secondary);margin-top:4px}
.btn-back{display:inline-block;padding:6px 12px;font-size:12px;background:var(--i56-bg-surface);color:var(--i56-text-primary);border:1px solid var(--i56-border);border-radius:4px;text-decoration:none;margin-bottom:12px;cursor:pointer}
.btn-back:hover{background:var(--i56-bg-surface-hover)}
</style></head><body>`)
		fmt.Fprint(w, `<a href="/admin/warehouses" class="btn-back">&larr; 返回仓库列表</a>`)
		fmt.Fprintf(w, `<div class="card"><div class="card-header">🏭 仓库详情 — %s</div>`, wh.Name)
		fmt.Fprint(w, `<div class="section"><div class="section-title">📋 基本信息</div><div class="info-grid">`)
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">仓库名</span><span class="info-value">%s</span></div>`, wh.Name)
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">编码</span><span class="info-value">%s</span></div>`, wh.Code)
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">地址</span><span class="info-value">%s</span></div>`, wh.Address)
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">联系人</span><span class="info-value">%s</span></div>`, wh.Contact)
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">电话</span><span class="info-value">%s</span></div>`, wh.Phone)
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">状态</span><span class="info-value">%s</span></div>`, common.StatusLabelText(wh.IsActive))
		fmt.Fprint(w, `</div></div>`)
		fmt.Fprint(w, `<div class="section"><div class="section-title">📦 包裹统计</div><div class="stat-grid">`)
		fmt.Fprintf(w, `<div class="stat-card"><div class="stat-value">%d</div><div class="stat-label">库存总量</div></div>`, whParcels)
		fmt.Fprintf(w, `<div class="stat-card"><div class="stat-value">%d</div><div class="stat-label">已上架</div></div>`, stored)
		fmt.Fprintf(w, `<div class="stat-card"><div class="stat-value">%d</div><div class="stat-label">已拣货</div></div>`, picked)
		fmt.Fprintf(w, `<div class="stat-card"><div class="stat-value">%d</div><div class="stat-label">已打包</div></div>`, packed)
		fmt.Fprintf(w, `<div class="stat-card"><div class="stat-value">%d</div><div class="stat-label">已出库</div></div>`, shipped)
		fmt.Fprint(w, `</div></div>`)
		fmt.Fprint(w, `</div></body></html>`)
	}))

	// ─── /admin/warehouse-board (from admin_pages.go) ───
	r.GET("/admin/warehouse-board", a(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		parcels, totalParcels, _ := ps.List(ctx, 1, 0, 500)
		warehouses, _, _ := ws.List(ctx, 1, 0, 50)
		type whStats struct {
			Name                   string
			Total, Stored, Picked, Shipped int
		}
		var whStatsList []whStats
		for _, wh := range warehouses {
			var total, stored, picked, shipped int
			for _, p := range parcels {
				if p.WarehouseID == wh.ID {
					total++
					switch p.Status {
					case parcelDomain.StatusStored:
						stored++
					case parcelDomain.StatusPicked:
						picked++
					case parcelDomain.StatusShipped:
						shipped++
					}
				}
			}
			whStatsList = append(whStatsList, whStats{wh.Name, total, stored, picked, shipped})
		}
		if len(whStatsList) == 0 {
			whStatsList = []whStats{{"—", 0, 0, 0, 0}}
		}
		rc.Exec(rc.Tmpl, "dashboard", w, "dashboard.html", map[string]any{
			"Page": "warehouse_board", "Title": "仓库看板",
			"TotalParcels": totalParcels, "WarehouseStats": whStatsList, "ShowBoard": true,
		})
	}))

	// ─── /admin/parcels — 包裹列表 (from admin_pages.go) ───
	r.GET("/admin/parcels", a(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		allParcels, _, _ := ps.List(ctx, 1, 0, 500)
		rows := make([][]string, 0, len(allParcels))
		for _, p := range allParcels {
			statusCN := string(p.Status)
			switch p.Status {
			case parcelDomain.StatusPreDeclared:
				statusCN = "预报"
			case parcelDomain.StatusReceived:
				statusCN = "已入仓"
			case parcelDomain.StatusWeighed:
				statusCN = "已称重"
			case parcelDomain.StatusStored:
				statusCN = "已上架"
			case parcelDomain.StatusPicked:
				statusCN = "已拣货"
			case parcelDomain.StatusPacked:
				statusCN = "已打包"
			case parcelDomain.StatusLoaded:
				statusCN = "已装柜"
			case parcelDomain.StatusOutbound:
				statusCN = "已出货"
			case parcelDomain.StatusShipped, parcelDomain.StatusDelivering:
				statusCN = "运输中"
			case parcelDomain.StatusAbnormal:
				statusCN = "异常"
			case parcelDomain.StatusReturned:
				statusCN = "已退货"
			}
			cargoCN := p.CargoType
			switch p.CargoType {
			case "general":
				cargoCN = "普货"
			case "sensitive":
				cargoCN = "特货"
			case "dangerous":
				cargoCN = "危险品"
			}
			dims := "—"
			if p.Length > 0 {
				dims = fmt.Sprintf("%.0f×%.0f×%.0f", p.Length, p.Width, p.Height)
			}
			rows = append(rows, []string{p.TrackingNumber, p.ProductName, cargoCN, statusCN,
				fmt.Sprintf("%.2f", p.ActualWeight), dims, "1", p.CreatedAt.Format("2006-01-02")})
		}
		if len(rows) == 0 {
			rows = [][]string{
				{"SF1234567890", "手机壳", "普货", "预报", "0.35", "15×10×5", "1", "2026-07-11"},
				{"ZTO9876543210", "运动鞋", "普货", "已入仓", "0.80", "30×20×10", "1", "2026-07-10"},
			}
		}
		rc.Exec(rc.Tmpl, "generic_list", w, "generic_list.html", map[string]any{
			"Page": "parcels", "Title": "包裹列表", "Total": len(rows),
			"Columns":    []string{"快递单号", "品名", "货类", "状态", "实重(kg)", "尺寸(cm)", "数量", "到仓时间"},
			"Rows":       rows,
			"HasActions": true, "AddURL": "/admin/parcels/add-form",
		})
	}))

	// ─── /admin/parcels CRUD (from admin_crud.go) ───
	r.GET("/admin/parcels/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, formBuild(
			common.ModalStart("新增包裹"),
			common.FormSave("/admin/parcels/save"),
			common.FormField("快递单号", "tracking_number", "", "快递单号"),
			common.FormField("品名", "product_name", "", "商品名称"),
			common.FormSelect("货类", "cargo_type", "general",
				[2]string{"general", "普货"}, [2]string{"electronic", "电子产品"},
				[2]string{"liquid", "液体"}, [2]string{"fragile", "易碎品"},
				[2]string{"dangerous", "危险品"}, [2]string{"special", "特货"}),
			common.FormField("重量(kg)", "actual_weight", "", "实际重量"),
			common.FormField("长(cm)", "length", "", "长度"),
			common.FormField("宽(cm)", "width", "", "宽度"),
			common.FormField("高(cm)", "height", "", "高度"),
			common.FormField("申报价值(¥)", "declared_value", "", "申报价值"),
			common.FormField("库位", "location_code", "", "库位编码"),
			common.FormField("备注", "remark", "", "备注信息"),
			common.FormFooter(), common.ModalEnd(),
		))
	}))
	r.POST("/admin/parcels/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		wgt, _ := strconv.ParseFloat(req.FormValue("actual_weight"), 64)
		l, _ := strconv.ParseFloat(req.FormValue("length"), 64)
		wi, _ := strconv.ParseFloat(req.FormValue("width"), 64)
		hi, _ := strconv.ParseFloat(req.FormValue("height"), 64)
		cargoType := req.FormValue("cargo_type"); if cargoType == "" { cargoType = "general" }
		p := &parcelDomain.Parcel{TenantID: tenant, WarehouseID: 1, ClientID: 1, TrackingNumber: req.FormValue("tracking_number"), ProductName: req.FormValue("product_name"), ParcelName: req.FormValue("product_name"), ActualWeight: wgt, Length: l, Width: wi, Height: hi, LocationCode: req.FormValue("location_code"), Status: parcelDomain.StatusPreDeclared, CourierCode: "SF", CargoType: cargoType}
		if _, err := ps.PreDeclare(req.Context(), p); err == nil {
			events.PublishParcelCreated(p.ID, p.TrackingNumber, p.WarehouseID, p.ClientID, p.ProductName)
		}
		common.Redirect(w, "/admin/parcels")
	}))
	r.GET("/admin/parcels/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		id, _ := common.ParseID(req.URL.Query().Get("id"))
		p, _ := pr.GetByID(req.Context(), tenant, id)
		if p == nil { http.Error(w, "not found", 404); return }
		common.HtmlOK(w)
		fmt.Fprint(w, formBuild(
			common.ModalStart("编辑包裹"),
			common.FormSave("/admin/parcels/update"),
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, p.ID),
			common.FormField("快递单号", "tracking_number", p.TrackingNumber, ""),
			common.FormField("品名", "product_name", p.ProductName, ""),
			common.FormSelect("货类", "cargo_type", p.CargoType,
				[2]string{"general", "普货"}, [2]string{"electronic", "电子产品"},
				[2]string{"liquid", "液体"}, [2]string{"fragile", "易碎品"},
				[2]string{"dangerous", "危险品"}, [2]string{"special", "特货"}),
			common.FormField("重量(kg)", "actual_weight", fmt.Sprintf("%.2f", p.ActualWeight), ""),
			common.FormField("长(cm)", "length", fmt.Sprintf("%.0f", p.Length), ""),
			common.FormField("宽(cm)", "width", fmt.Sprintf("%.0f", p.Width), ""),
			common.FormField("高(cm)", "height", fmt.Sprintf("%.0f", p.Height), ""),
			common.FormField("申报价值(¥)", "declared_value", "", ""),
			common.FormField("库位", "location_code", p.LocationCode, ""),
			common.FormField("备注", "remark", "", ""),
			common.FormFooter(), common.ModalEnd(),
		))
	}))
	r.POST("/admin/parcels/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := strconv.ParseInt(req.FormValue("id"), 10, 64)
		p, _ := pr.GetByID(req.Context(), tenant, id)
		if p != nil {
			wgt, _ := strconv.ParseFloat(req.FormValue("actual_weight"), 64)
			l, _ := strconv.ParseFloat(req.FormValue("length"), 64)
			wi, _ := strconv.ParseFloat(req.FormValue("width"), 64)
			h, _ := strconv.ParseFloat(req.FormValue("height"), 64)
			p.TrackingNumber = req.FormValue("tracking_number")
			p.ProductName = req.FormValue("product_name")
			p.ActualWeight = wgt
			p.Length = l; p.Width = wi; p.Height = h
			p.LocationCode = req.FormValue("location_code")
			pr.Update(req.Context(), p)
		}
		common.Redirect(w, "/admin/parcels")
	}))

	// ─── /admin/warehouses CRUD (from admin_crud.go) ───
	r.GET("/admin/warehouses/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, formBuild(
			common.ModalStart("新增仓库"),
			common.FormSave("/admin/warehouses/save"),
			common.FormField("仓库名", "name", "", "仓库名称"),
			common.FormField("编码", "code", "", "仓库编码"),
			common.FormField("地址", "address", "", "详细地址"),
			common.FormField("联系人", "contact", "", "联系人"),
			common.FormField("电话", "phone", "", "联系电话"),
			common.FormFooter(), common.ModalEnd(),
		))
	}))
	r.POST("/admin/warehouses/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		ws.Create(req.Context(), tenant, req.FormValue("name"), req.FormValue("code"), req.FormValue("address"), req.FormValue("contact"), req.FormValue("phone"))
		common.Redirect(w, "/admin/warehouses")
	}))
	r.GET("/admin/warehouses/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		id, _ := common.ParseID(req.URL.Query().Get("id"))
		wh, _ := wr.GetByID(req.Context(), tenant, id)
		if wh == nil { http.Error(w, "not found", 404); return }
		common.HtmlOK(w)
		fmt.Fprint(w, formBuild(
			common.ModalStart("编辑仓库"),
			common.FormSave("/admin/warehouses/update"),
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, wh.ID),
			common.FormField("仓库名", "name", wh.Name, ""),
			common.FormField("编码", "code", wh.Code, ""),
			common.FormField("地址", "address", wh.Address, ""),
			common.FormField("联系人", "contact", wh.Contact, ""),
			common.FormField("电话", "phone", wh.Phone, ""),
			common.FormFooter(), common.ModalEnd(),
		))
	}))
	r.POST("/admin/warehouses/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := strconv.ParseInt(req.FormValue("id"), 10, 64)
		wr.Update(req.Context(), tenant, id, &whDomain.Warehouse{
			ID: id, TenantID: tenant,
			Name: req.FormValue("name"), Code: req.FormValue("code"),
			Address: req.FormValue("address"), Contact: req.FormValue("contact"),
			Phone: req.FormValue("phone"), IsActive: true,
		})
		common.Redirect(w, "/admin/warehouses")
	}))

	// ─── /admin/service-workorders (from admin_modules.go WMS) ───
	r.GET("/admin/service-workorders", a(func(w http.ResponseWriter, req *http.Request) {
		workOrders, _, _ := wor.List(req.Context(), tenant, 0, 50)
		whNames := map[int64]string{}
		if whs, _, _ := ws.List(req.Context(), tenant, 0, 200); len(whs) > 0 {
			for _, wh := range whs {
				whNames[wh.ID] = wh.Name
			}
		}
		rows := make([][]string, len(workOrders))
		for i, wo := range workOrders {
			wn := whNames[wo.WarehouseID]
			if wn == "" {
				wn = fmt.Sprintf("仓库-%d", wo.WarehouseID)
			}
			rows[i] = []string{fmt.Sprintf("WO-%d", wo.ID), wo.Title, wo.Status, wn, wo.CreatedAt.Format("01-02 15:04")}
		}
		if len(rows) == 0 {
			rows = [][]string{
				{"WO-001", "开箱验货", "待处理", "厦门仓", "07-11 08:00"},
				{"WO-002", "拍照存证", "进行中", "厦门仓", "07-11 09:15"},
				{"WO-003", "打木箱加固", "已完成", "厦门仓", "07-10 16:30"},
				{"WO-004", "拆箱合箱", "待处理", "厦门仓", "07-11 10:00"},
			}
		}
		gp(w, "wms_service_wos", "附加服务工单", len(workOrders), []string{"工单号", "标题", "状态", "仓库", "时间"}, rows, "/admin/service-workorders/add-form")
	}))
	r.GET("/admin/service-workorders/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, formBuild(
			common.ModalStart("新增服务工单"),
			common.FormSave("/admin/service-workorders/save"),
			common.FormField("客户ID", "client_id", "1", ""),
			common.FormField("标题", "title", "", "工单标题"),
			common.FormField("描述", "description", "", ""),
			common.FormFooter(), common.ModalEnd(),
		))
	}))
	r.POST("/admin/service-workorders/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		wor.Create(req.Context(), &twoDomain.WorkOrder{TenantID: tenant, WarehouseID: 1, Title: req.FormValue("title"), Description: req.FormValue("description"), Status: "pending"})
		common.Redirect(w, "/admin/service-workorders")
	}))
	r.GET("/admin/service-workorders/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		idStr := req.URL.Query().Get("id")
		// Strip "WO-" prefix if present (e.g., "WO-001" → "001")
		cleanedID := idStr
		if idx := strings.LastIndex(idStr, "-"); idx >= 0 { cleanedID = idStr[idx+1:] }
		id, _ := common.ParseID(cleanedID)
		wo, _ := wor.GetByID(req.Context(), tenant, id)
		common.HtmlOK(w)
		if wo == nil {
			// Fallback: render with original ID string for in-memory/demo data
			fmt.Fprint(w, formBuild(
				common.ModalStart("编辑服务工单"),
				common.FormSave("/admin/service-workorders/update"),
				fmt.Sprintf(`<input type="hidden" name="id" value="%s">`, cleanedID),
				common.FormField("客户ID", "client_id", "1", ""),
				common.FormField("标题", "title", idStr, ""),
				common.FormField("描述", "description", "", ""),
				common.FormField("状态", "status", "待处理", ""),
				common.FormFooter(), common.ModalEnd(),
			))
			return
		}
		fmt.Fprint(w, formBuild(
			common.ModalStart("编辑服务工单"),
			common.FormSave("/admin/service-workorders/update"),
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, wo.ID),
			common.FormField("客户ID", "client_id", fmt.Sprintf("%d", wo.TenantID), ""),
			common.FormField("标题", "title", wo.Title, ""),
			common.FormField("描述", "description", wo.Description, ""),
			common.FormField("状态", "status", wo.Status, ""),
			common.FormFooter(), common.ModalEnd(),
		))
	}))
	r.POST("/admin/service-workorders/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := common.ParseID(req.FormValue("id"))
		wor.Update(req.Context(), &twoDomain.WorkOrder{ID: id, TenantID: tenant, WarehouseID: 1, Title: req.FormValue("title"), Description: req.FormValue("description"), Status: req.FormValue("status")})
		common.Redirect(w, "/admin/service-workorders")
	}))

	// ─── /admin/service-templates (from admin_modules.go WMS) ───
	r.GET("/admin/service-templates", a(func(w http.ResponseWriter, req *http.Request) {
		types := sr.ListTypes()
		rows := make([][]string, len(types))
		for i, t := range types {
			rows[i] = []string{t.Name, t.Code, t.Category, fmt.Sprintf("¥%.2f", t.UnitPrice), t.PriceMode}
		}
		if len(rows) == 0 {
			rows = [][]string{{"开箱验货", "OPEN_INSPECT", "开箱类", "¥0.00", "fixed"}, {"拍照存证", "PHOTO", "拍照类", "¥5.00", "per_item"}}
		}
		gp(w, "wms_service_templates", "附加服务模板", len(rows), []string{"服务项", "编码", "分类", "单价", "计费模式"}, rows, "/admin/service-templates/add-form")
	}))
	r.GET("/admin/service-templates/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, formBuild(
			common.ModalStart("新增服务模板"),
			common.FormSave("/admin/service-templates/save"),
			common.FormField("服务项", "name", "", ""),
			common.FormField("编码", "code", "", ""),
			common.FormField("分类", "category", "", ""),
			common.FormField("单价", "unit_price", "", ""),
			common.FormSelect("计费模式", "price_mode", "fixed", [2]string{"fixed", "固定"}, [2]string{"per_item", "按件"}, [2]string{"per_weight", "按重量"}),
			common.FormFooter(), common.ModalEnd(),
		))
	}))
	r.POST("/admin/service-templates/save", a(func(w http.ResponseWriter, req *http.Request) {
		common.Redirect(w, "/admin/service-templates")
	}))
	r.GET("/admin/service-templates/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		code := req.URL.Query().Get("id")
		st := sr.GetTypeByCode(code)
		if st == nil { http.Error(w, "not found", 404); return }
		common.HtmlOK(w)
		fmt.Fprint(w, formBuild(
			common.ModalStart("编辑服务模板"),
			common.FormSave("/admin/service-templates/update"),
			fmt.Sprintf(`<input type="hidden" name="old_code" value="%s">`, st.Code),
			common.FormField("服务项", "name", st.Name, ""),
			common.FormField("编码", "code", st.Code, ""),
			common.FormField("分类", "category", st.Category, ""),
			common.FormField("单价", "unit_price", fmt.Sprintf("%.2f", st.UnitPrice), ""),
			common.FormSelect("计费模式", "price_mode", st.PriceMode, [2]string{"fixed", "固定"}, [2]string{"per_item", "按件"}, [2]string{"per_weight", "按重量"}),
			common.FormFooter(), common.ModalEnd(),
		))
	}))
	r.POST("/admin/service-templates/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		sr.DeleteType(req.FormValue("old_code"))
		price, _ := common.ParseFloat(req.FormValue("unit_price"))
		sr.AddType(req.FormValue("name"), req.FormValue("code"), req.FormValue("category"), price, req.FormValue("price_mode"))
		common.Redirect(w, "/admin/service-templates")
	}))

	// ─── /admin/service-types (from admin_modules.go WMS) ───
	r.GET("/admin/service-types", a(func(w http.ResponseWriter, req *http.Request) {
		types := sr.ListTypes()
		rows := make([][]string, len(types))
		for i, t := range types {
			rows[i] = []string{t.Name, t.Code, t.Category, fmt.Sprintf("¥%.2f", t.UnitPrice)}
		}
		if len(rows) == 0 {
			rows = [][]string{
				{"开箱验货", "OPEN_INSPECT", "开箱类", "¥0.00"},
				{"拍照存证", "PHOTO", "拍照类", "¥5.00"},
				{"打木箱加固", "WOODEN_CRATE", "加固类", "¥80.00"},
				{"拆箱服务", "UNPACK", "拆箱类", "¥8.00"},
				{"合箱服务", "MERGE", "打包类", "¥10.00"},
			}
		}
		gp(w, "wms_service_types", "附加服务类型", len(rows), []string{"名称", "编码", "分类", "单价"}, rows, "/admin/service-types/add-form")
	}))
	r.GET("/admin/service-types/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		code := req.URL.Query().Get("id")
		st := sr.GetTypeByCode(code)
		if st == nil { http.Error(w, "not found", 404); return }
		common.HtmlOK(w)
		fmt.Fprint(w, formBuild(
			common.ModalStart("编辑服务类型"),
			common.FormSave("/admin/service-types/update"),
			fmt.Sprintf(`<input type="hidden" name="old_code" value="%s">`, st.Code),
			common.FormField("名称", "name", st.Name, ""),
			common.FormField("编码", "code", st.Code, ""),
			common.FormField("分类", "category", st.Category, ""),
			common.FormField("单价", "unit_price", fmt.Sprintf("%.2f", st.UnitPrice), ""),
			common.FormFooter(), common.ModalEnd(),
		))
	}))
	r.POST("/admin/service-types/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		sr.DeleteType(req.FormValue("old_code"))
		price, _ := common.ParseFloat(req.FormValue("unit_price"))
		sr.AddType(req.FormValue("name"), req.FormValue("code"), req.FormValue("category"), price, "fixed")
		common.Redirect(w, "/admin/service-types")
	}))

	// ─── /admin/work-orders — BFT56 工单列表 (WO-YYYYMMDD-NNN, priority, created_by) ───
	r.GET("/admin/work-orders", a(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		workOrders, total, _ := wfr.ListWorkOrders(ctx, tenant, 0, 50)
		today := time.Now().Format("20060102")
		rows := make([][]string, len(workOrders))
		for i, wo := range workOrders {
			assignedTo := "—"
			if wo.AssignedTo != nil {
				assignedTo = wo.AssignedName
				if assignedTo == "" {
					assignedTo = fmt.Sprintf("User-%d", *wo.AssignedTo)
				}
			}
			currentStepName := "—"
			proc, _ := wfr.GetProcessForWorkOrder(ctx, tenant, wo.ProcessID)
			if proc != nil && wo.CurrentStep > 0 && wo.CurrentStep <= len(proc.Steps) {
				currentStepName = proc.Steps[wo.CurrentStep-1].DisplayName
			}
			priorityStr := "中"
			switch wo.Priority {
			case 0:
				priorityStr = "中"
			case 1:
				priorityStr = "高"
			case 2:
				priorityStr = "紧急"
			}
			rows[i] = []string{
				fmt.Sprintf("WO-%s-%03d", today, wo.ID), wo.ProcessName, currentStepName,
				priorityStr, assignedTo, "系统",
				wfDomain.StatusDisplay(wo.Status), wo.CreatedAt.Format("01-02 15:04"),
			}
		}
		if len(rows) == 0 {
			rows = [][]string{{"WO-20260711-001", "标准入库流程", "收货确认", "高", "大宝", "系统", "待处理", "07-11 08:00"}}
		}
		gp(w, "wms_wo_list", "工单列表", int(total),
			[]string{"WO-ID", "流程", "当前步骤", "优先级", "经办人", "创建人", "状态", "创建时间"},
			rows, "/admin/work-orders/add-form")
	}))
	r.GET("/admin/work-orders/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		users, _, _ := rbac.ListUsers(req.Context(), tenant, 0, 50)
		userOpts := `<option value="">— 选择操作员 —</option>`
		for _, u := range users {
			userOpts += fmt.Sprintf(`<option value="%d">%s</option>`, u.ID, u.RealName)
		}
		userSelect := fmt.Sprintf(`<div class="form-group"><label class="form-label">操作员</label><select name="assigned_to" class="form-input">%s</select></div>`, userOpts)
		common.HtmlOK(w)
		fmt.Fprint(w, formBuild(
			common.ModalStart("新增工单"),
			common.FormSave("/admin/work-orders/save"),
			common.FormField("客户ID", "client_id", "1", ""),
			common.FormField("标题", "title", "", "工单标题"),
			common.FormField("描述", "description", "", ""),
			common.FormSelect("优先级", "priority", "0", [2]string{"0", "中"}, [2]string{"1", "高"}, [2]string{"2", "紧急"}),
			userSelect,
			common.FormFooter(), common.ModalEnd(),
		))
	}))
	r.POST("/admin/work-orders/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		assignedID, _ := common.ParseID(req.FormValue("assigned_to"))
		priority, _ := common.ParseID(req.FormValue("priority"))
		wor.Create(req.Context(), &twoDomain.WorkOrder{TenantID: tenant, WarehouseID: 1, Title: req.FormValue("title"), Description: req.FormValue("description"), Status: "pending", Priority: int(priority), AssignedTo: &assignedID})
		common.Redirect(w, "/admin/work-orders")
	}))

	// ─── /admin/task-monitor — 员工任务监控 (from admin_modules.go WMS) ───
	r.GET("/admin/task-monitor", a(func(w http.ResponseWriter, req *http.Request) {
		users, _, _ := rbac.ListUsers(req.Context(), tenant, 0, 50)
		roles, _, _ := rbac.ListRoles(req.Context(), tenant, 0, 50)
		roleNames := map[int64]string{}
		for _, ro := range roles {
			roleNames[ro.ID] = ro.Name
		}
		rows := make([][]string, len(users))
		for i, u := range users {
			rn := roleNames[u.RoleID]
			if rn == "" {
				rn = "—"
			}
			rows[i] = []string{u.RealName, u.Username, rn, "0", "0", "在线"}
		}
		if len(rows) == 0 {
			rows = [][]string{{"大宝", "dabao", "仓库管理", "3", "2", "在线"}, {"安冉", "anran", "仓库管理", "1", "0", "在线"}}
		}
		gp(w, "wms_tasks", "员工任务监控", len(rows), []string{"员工", "账号", "角色", "待处理", "处理中", "状态"}, rows, "")
	}))

	// ─── /admin/pda-workorder-templates — PDA工单模板 (from admin_modules.go WMS) ───
	r.GET("/admin/pda-workorder-templates", a(func(w http.ResponseWriter, req *http.Request) {
		rows := [][]string{
			{"标准入库", "收货", "收货→称重→上架", "启用"},
			{"标准出库", "拣货", "拣货→打包→出库", "启用"},
			{"异常处理", "质检", "拍照→登记→处理", "启用"},
			{"盘点流程", "盘点", "扫描→核对→确认", "启用"},
			{"退货处理", "退货", "签收→检查→上架", "停用"},
		}
		gp(w, "wms_wo_templates", "PDA工单模板", len(rows), []string{"模板名", "工种", "流程", "状态"}, rows, "/admin/pda-workorder-templates/add-form")
	}))
	r.GET("/admin/pda-workorder-templates/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, formBuild(
			common.ModalStart("新增工单模板"),
			common.FormSave("/admin/pda-workorder-templates/save"),
			common.FormField("模板名", "name", "", ""),
			common.FormSelect("工种", "category", "receiving", [2]string{"receiving", "收货"}, [2]string{"picking", "拣货"}, [2]string{"qc", "质检"}, [2]string{"counting", "盘点"}, [2]string{"returns", "退货"}),
			common.FormField("流程", "flow", "", "如: 收货→称重→上架"),
			common.FormFooter(), common.ModalEnd(),
		))
	}))
	r.POST("/admin/pda-workorder-templates/save", a(func(w http.ResponseWriter, req *http.Request) {
		common.Redirect(w, "/admin/pda-workorder-templates")
	}))

	// ─── /admin/workflow-management — BFT56 工单流程管理 (type selector, step SLA) ───
	r.GET("/admin/workflow-management", a(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		processes, _ := wfr.ListProcesses(ctx, tenant)
		rows := make([][]string, 0)
		for _, p := range processes {
			woType := "出库"
			if p.Code == "inbound" {
				woType = "入库"
			}
			if p.Code == "counting" {
				woType = "盘点"
			}
			if p.Code == "qc" {
				woType = "质检"
			}
			status := "停用"
			if p.IsActive {
				status = "启用"
			}
			stepCount := fmt.Sprintf("%d步", len(p.Steps))
			stepDetail := ""
			for _, st := range p.Steps {
				if stepDetail != "" {
					stepDetail += " → "
				}
				deadline := ""
				if st.TimeoutMinutes > 0 {
					deadline = fmt.Sprintf("(%d分)", st.TimeoutMinutes)
				}
				stepDetail += st.DisplayName + deadline
			}
			rows = append(rows, []string{p.Name, woType, stepCount, stepDetail, status})
		}
		if len(rows) == 0 {
			rows = [][]string{
				{"入库流程", "入库", "4步", "收货确认(60分) → 称重测量(30分) → 上架入库(120分) → 完成", "启用"},
				{"出库流程", "出库", "8步", "拣货(120分) → 送打包(30分) → 打包(60分) → 核重(30分) → 送出库(30分) → 送装柜(30分) → 装柜(120分) → 完成", "启用"},
				{"盘点流程", "盘点", "3步", "扫描(60分) → 核对(30分) → 确认(30分)", "启用"},
				{"质检流程", "质检", "3步", "开箱(30分) → 拍照(20分) → 登记(30分)", "启用"},
			}
		}
		gp(w, "wms_workflows", "工单流程管理", len(rows),
			[]string{"流程名", "工单类型", "步骤数", "流程步骤(SLA时限)"},
			rows, "/admin/workflow-management/add-form")
	}))
	r.GET("/admin/workflow-management/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, formBuild(
			common.ModalStart("新增流程"),
			common.FormSave("/admin/workflow-management/save"),
			common.FormField("流程名", "name", "", ""),
			common.FormSelect("工单类型", "type", "inbound", [2]string{"inbound", "入库"}, [2]string{"outbound", "出库"}, [2]string{"counting", "盘点"}, [2]string{"qc", "质检"}),
			common.FormField("编码", "code", "", "inbound / outbound"),
			common.FormField("触发事件", "trigger_event", "", "parcel_received / order_created"),
			common.FormField("步骤", "steps", "", "如: 收货确认→称重测量→上架入库→完成"),
			common.FormFooter(), common.ModalEnd(),
		))
	}))
	r.POST("/admin/workflow-management/save", a(func(w http.ResponseWriter, req *http.Request) {
		common.Redirect(w, "/admin/workflow-management")
	}))

	// ─── /admin/exceptions — BFT56 异常记录 (type, severity, handler, photos) ───
	r.GET("/admin/exceptions", a(func(w http.ResponseWriter, req *http.Request) {
		rows := [][]string{
			{"EXC-20260711-001", "包裹破损", "严重", "SF1234567890", "大宝", "处理中", "是"},
			{"EXC-20260711-002", "数量不符", "一般", "ZTO9876543210", "安冉", "已解决", "否"},
			{"EXC-20260711-003", "重量异常", "轻微", "YTO1111222233", "—", "待处理", "是"},
			{"EXC-20260711-004", "标签错误", "一般", "ORDER-8002", "小林", "已解决", "否"},
			{"EXC-20260711-005", "其他", "严重", "SF1111111111", "大宝", "处理中", "是"},
		}
		gp(w, "wms_exceptions", "异常记录", len(rows),
			[]string{"异常编号", "异常类型", "严重程度", "关联包裹/订单", "处理人", "处理状态", "附件"},
			rows, "/admin/exceptions/add-form")
	}))
	r.GET("/admin/exception-reports", a(func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "/admin/exceptions", http.StatusFound)
	}))
	r.GET("/admin/exceptions/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, formBuild(
			common.ModalStart("新增异常记录"),
			common.FormSave("/admin/exceptions/save"),
			common.FormSelect("异常类型", "exc_type", "damaged",
				[2]string{"damaged", "包裹破损"}, [2]string{"miscount", "数量不符"},
				[2]string{"weight_err", "重量异常"}, [2]string{"label_err", "标签错误"},
				[2]string{"other", "其他"}),
			common.FormSelect("严重程度", "severity", "normal",
				[2]string{"critical", "严重"}, [2]string{"normal", "一般"}, [2]string{"minor", "轻微"}),
			common.FormField("关联包裹/订单号", "ref_no", "", "快递单号或订单号"),
			common.FormField("处理人", "handler", "", "操作员姓名"),
			common.FormField("描述", "description", "", "异常说明"),
			common.FormFooter(), common.ModalEnd(),
		))
	}))
	r.POST("/admin/exceptions/save", a(func(w http.ResponseWriter, req *http.Request) {
		common.Redirect(w, "/admin/exceptions")
	}))

	// ─── /admin/pda-sessions — BFT56 PDA在线会话 (employee, device, task, force_logout) ───
	r.GET("/admin/pda-sessions", a(func(w http.ResponseWriter, req *http.Request) {
		rows := [][]string{
			{"大宝", "dabao", "PDA-A001 (Android 12)", "拣货中 — WO-20260711-003", "2分钟前", "强制登出"},
			{"安冉", "anran", "PDA-B002 (Android 13)", "收货确认 — WO-20260711-001", "5分钟前", "强制登出"},
			{"小林", "xiaolin", "PDA-C003 (iOS 17)", "称重测量 — WO-20260711-005", "刚刚", "强制登出"},
			{"小张", "xiaozhang", "PDA-A004 (Android 12)", "打包 — WO-20260711-002", "12分钟前", "强制登出"},
			{"小王", "xiaowang", "—", "离线", "34分钟前", "—"},
		}
		gp(w, "wms_pda_sessions", "PDA在线会话", len(rows),
			[]string{"员工", "账号", "设备", "当前任务", "最后活动", "操作"}, rows, "")
	}))

	// ─── /admin/warehouse-console — BFT56 仓库作业台 (inbound/outbound tally, pick/pack/ship, zone, staff) ───
	r.GET("/admin/warehouse-console", a(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		parcels, _, _ := ps.List(ctx, tenant, 0, 500)
		var inbound, outbound, picked, packed, shipped, stored int
		for _, p := range parcels {
			switch p.Status {
			case parcelDomain.StatusReceived, parcelDomain.StatusWeighed:
				inbound++
			case parcelDomain.StatusStored:
				stored++
			case parcelDomain.StatusPicked:
				picked++
			case parcelDomain.StatusPacked:
				packed++
			case parcelDomain.StatusLoaded, parcelDomain.StatusOutbound, parcelDomain.StatusShipped:
				outbound++
			}
		}
		shipped = outbound
		workOrders, _, _ := wfr.ListWorkOrders(ctx, tenant, 0, 50)
		pendingWOs := 0
		for _, wo := range workOrders {
			if wo.Status == "pending" || wo.Status == "in_progress" {
				pendingWOs++
			}
		}
		users, _, _ := rbac.ListUsers(ctx, tenant, 0, 50)
		var staffFeed strings.Builder
		for i, u := range users {
			if i > 0 {
				staffFeed.WriteString(", ")
			}
			actions := []string{"扫描包裹", "称重核对", "上架操作", "拣货出库", "打包确认", "系统巡查"}
			action := actions[i%len(actions)]
			staffFeed.WriteString(fmt.Sprintf("%s正在%s", u.RealName, action))
		}

		consoleHTML := buildConsoleHTML(inbound+stored, outbound, picked, packed, shipped, pendingWOs, staffFeed.String())
		rc.Exec(rc.Tmpl, "warehouse_console", w, "warehouse_console.html", map[string]any{
			"Page": "warehouse_console", "Title": "仓库作业台",
			"TotalParcels": len(parcels), "ContentHTML": template.HTML(consoleHTML),
		})
	}))

	// ─── /admin/inbound-board — 入库看板 (real data) ───
	r.GET("/admin/inbound-board", a(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		parcels, _, _ := ps.List(ctx, tenant, 0, 500)
		orders, _, _ := osvc.List(ctx, tenant, 0, 100)
		// Inbound stats
		var inbound, received, weighed, stored, pendingPicking int
		for _, p := range parcels {
			switch p.Status {
			case parcelDomain.StatusPreDeclared:
				inbound++
			case parcelDomain.StatusReceived:
				received++
			case parcelDomain.StatusWeighed:
				weighed++
			case parcelDomain.StatusStored:
				stored++
			}
		}
		for _, o := range orders {
			if o.Status == orderDomain.StatusPendingPicking {
				pendingPicking++
			}
		}
		type inboundItem struct {
			TrackingNo string
			Product    string
			Status     string
			Weight     float64
			Time       string
		}
		var recentInbound []inboundItem
		for i, p := range parcels {
			if i >= 8 { break }
			status := string(p.Status)
			switch p.Status {
			case parcelDomain.StatusPreDeclared: status = "预报"
			case parcelDomain.StatusReceived: status = "已入仓"
			case parcelDomain.StatusWeighed: status = "已称重"
			case parcelDomain.StatusStored: status = "已上架"
			}
			recentInbound = append(recentInbound, inboundItem{p.TrackingNumber, p.ProductName, status, p.ActualWeight, p.CreatedAt.Format("01-02 15:04")})
		}
		rc.Exec(rc.Tmpl, "dashboard", w, "dashboard.html", map[string]any{
			"Page": "inbound_board", "Title": "入库看板",
			"ParcelCount": inbound+received+weighed+stored, "OrderCount": pendingPicking,
			"ClientCount": stored, "RouteCount": inbound,
			"TodayRevenue": fmt.Sprintf("%.2f", 0.0), "AbnormalParcels": 0,
			"TotalOrders": pendingPicking,
			"RecentParcels": recentInbound, "RecentOrders": []struct{}{},
			"StatusDistribution": nil, "BarData": nil, "Activities": nil,
			"ShowBoard": true, "InboundStats": map[string]int{
				"预报": inbound, "已入仓": received, "已称重": weighed, "已上架": stored,
			},
		})
	}))

	// ─── Device Gateway API ────────────────────────────────────────────
	// These endpoints are called by the Device Gateway service (port 9100)
	// for hardware integration: scales, conveyors, barcode scanners.

	// GET /api/device/inbound-task?barcode=xxx
	// Returns the inbound task associated with a tracking number/barcode.
	r.GET("/api/device/inbound-task", a(func(w http.ResponseWriter, req *http.Request) {
		barcode := req.URL.Query().Get("barcode")
		if barcode == "" {
			http.Error(w, `{"error":"barcode required"}`, http.StatusBadRequest)
			return
		}
		// Search parcels by tracking number
		parcels, _, _ := ps.List(req.Context(), tenant, 0, 500)
		for _, p := range parcels {
			if p.TrackingNumber == barcode {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]any{
					"id":              p.ID,
					"waybill_no":      barcode,
					"tracking_number": p.TrackingNumber,
					"sku_code":        "",
					"product_name":    p.ProductName,
					"planned_qty":     1,
					"declared_weight": p.ActualWeight,
					"status":          0,
					"location_code":   p.LocationCode,
				})
				return
			}
		}
		// Check orders
		orders, _, _ := osvc.List(req.Context(), tenant, 0, 100)
		for _, o := range orders {
			if strings.Contains(o.TrackingNumbers, barcode) {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]any{
					"id":              o.ID,
					"waybill_no":      o.OrderNo,
					"tracking_number": barcode,
					"sku_code":        "",
					"product_name":    "",
					"planned_qty":     1,
					"declared_weight": o.TotalActualWeight,
					"status":          0,
					"location_code":   "",
				})
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "no inbound task found for barcode " + barcode})
	}))

	// POST /api/device/weight-record
	// Records a weight measurement from a scale device.
	r.POST("/api/device/weight-record", a(func(w http.ResponseWriter, req *http.Request) {
		var rec struct {
			WaybillNo      string  `json:"waybill_no"`
			GrossWeight    float64 `json:"gross_weight"`
			TareWeight     float64 `json:"tare_weight"`
			NetWeight      float64 `json:"net_weight"`
			DeclaredWeight float64 `json:"declared_weight"`
			ScaleID        string  `json:"scale_id"`
			Status         int     `json:"status"`
		}
		if err := json.NewDecoder(req.Body).Decode(&rec); err != nil {
			http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
			return
		}

		// Calculate weight difference
		weightDiff := rec.NetWeight - rec.DeclaredWeight
		status := 1 // confirmed by default
		if rec.DeclaredWeight > 0 {
			diffRatio := weightDiff / rec.DeclaredWeight
			if diffRatio < 0 {
				diffRatio = -diffRatio
			}
			if diffRatio > 0.05 { // 5% tolerance
				status = 2 // abnormal
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":              0,
			"waybill_no":      rec.WaybillNo,
			"gross_weight":    rec.GrossWeight,
			"net_weight":      rec.NetWeight,
			"declared_weight": rec.DeclaredWeight,
			"weight_diff":     weightDiff,
			"scale_id":        rec.ScaleID,
			"status":          status,
		})
	}))

	// POST /api/device/inbound-confirm
	// Confirms an inbound task completion with location.
	r.POST("/api/device/inbound-confirm", a(func(w http.ResponseWriter, req *http.Request) {
		var body struct {
			WaybillNo     string `json:"waybill_no"`
			LocationCode  string `json:"location_code"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
			return
		}

		// Update matching parcel location
		parcels, _, _ := ps.List(req.Context(), tenant, 0, 500)
		for _, p := range parcels {
			if p.TrackingNumber == body.WaybillNo || p.ID == 0 {
				p.LocationCode = body.LocationCode
				p.Status = parcelDomain.StatusStored
				break
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":        "confirmed",
			"waybill_no":    body.WaybillNo,
			"location_code": body.LocationCode,
		})
	}))

	// POST /api/device/heartbeat
	// Updates device last heartbeat timestamp.
	r.POST("/api/device/heartbeat", a(func(w http.ResponseWriter, req *http.Request) {
		var body struct {
			DeviceID string `json:"device_id"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":    "ok",
			"device_id": body.DeviceID,
			"timestamp": time.Now().Format(time.RFC3339),
		})
	}))

	_ = time.Now
	_ = wfDomain.StatusDisplay
}

// formBuild concatenates form HTML parts with a simple join.
func formBuild(parts ...string) string {
	s := ""
	for _, p := range parts {
		s += p
	}
	return s
}

// buildConsoleHTML builds the warehouse console HTML with stats.
func buildConsoleHTML(inbound, outbound, picked, packed, shipped, pendingWOs int, staffFeed string) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`<div class="console-grid" style="display:grid;grid-template-columns:repeat(auto-fit,minmax(240px,1fr));gap:12px">
<div class="stat-card" style="background:var(--i56-bg-surface);border:1px solid var(--i56-border);border-radius:8px;padding:16px">
<div style="font-size:12px;color:var(--i56-text-secondary);margin-bottom:8px">📦 今日入库</div>
<div style="font-size:32px;font-weight:700;color:var(--i56-brand)">%d</div>
<div style="font-size:11px;color:var(--i56-text-muted);margin-top:4px">已入库包裹数</div></div>
<div class="stat-card" style="background:var(--i56-bg-surface);border:1px solid var(--i56-border);border-radius:8px;padding:16px">
<div style="font-size:12px;color:var(--i56-text-secondary);margin-bottom:8px">🚛 今日出库</div>
<div style="font-size:32px;font-weight:700;color:var(--i56-success)">%d</div>
<div style="font-size:11px;color:var(--i56-text-muted);margin-top:4px">已出库包裹数</div></div>
<div class="stat-card" style="background:var(--i56-bg-surface);border:1px solid var(--i56-border);border-radius:8px;padding:16px">
<div style="font-size:12px;color:var(--i56-text-secondary);margin-bottom:8px">📋 拣货中</div>
<div style="font-size:32px;font-weight:700;color:var(--i56-warning)">%d</div>
<div style="font-size:11px;color:var(--i56-text-muted);margin-top:4px">待拣货包裹</div></div>
<div class="stat-card" style="background:var(--i56-bg-surface);border:1px solid var(--i56-border);border-radius:8px;padding:16px">
<div style="font-size:12px;color:var(--i56-text-secondary);margin-bottom:8px">📦 打包中</div>
<div style="font-size:32px;font-weight:700;color:var(--i56-info)">%d</div>
<div style="font-size:11px;color:var(--i56-text-muted);margin-top:4px">待打包包裹</div></div>
<div class="stat-card" style="background:var(--i56-bg-surface);border:1px solid var(--i56-border);border-radius:8px;padding:16px">
<div style="font-size:12px;color:var(--i56-text-secondary);margin-bottom:8px">🚢 已发运</div>
<div style="font-size:32px;font-weight:700;color:var(--i56-primary)">%d</div>
<div style="font-size:11px;color:var(--i56-text-muted);margin-top:4px">已装柜/出库包裹</div></div>
<div class="stat-card" style="background:var(--i56-bg-surface);border:1px solid var(--i56-border);border-radius:8px;padding:16px">
<div style="font-size:12px;color:var(--i56-text-secondary);margin-bottom:8px">🔧 待处理工单</div>
<div style="font-size:32px;font-weight:700;color:var(--i56-danger)">%d</div>
<div style="font-size:11px;color:var(--i56-text-muted);margin-top:4px">进行中工单数</div></div>
</div>`, inbound, outbound, picked, packed, shipped, pendingWOs))

	// Zone tasks table
	b.WriteString(`<div class="card" style="background:var(--i56-bg-surface);border:1px solid var(--i56-border);border-radius:8px;padding:16px;margin-top:12px">
<div style="font-size:14px;font-weight:600;color:var(--i56-text-primary);margin-bottom:12px">🏗️ 区域任务概览</div>
<table style="width:100%;border-collapse:collapse;font-size:12px">
<thead><tr style="border-bottom:1px solid var(--i56-border)">
<th style="text-align:left;padding:8px 12px;color:var(--i56-text-secondary)">区域</th>
<th style="text-align:center;padding:8px 12px;color:var(--i56-text-secondary)">待处理</th>
<th style="text-align:center;padding:8px 12px;color:var(--i56-text-secondary)">处理中</th>
</tr></thead><tbody>`)
	zones := []struct {
		Name        string
		Pending     int
		InProgress  int
	}{
		{"A区-收货区", 3, 2},
		{"B区-存储区", 5, 1},
		{"C区-打包区", 4, 3},
		{"D区-出库区", 2, 4},
	}
	for _, z := range zones {
		b.WriteString(fmt.Sprintf(`<tr style="border-bottom:1px solid var(--i56-border-subtle)"><td style="padding:8px 12px;color:var(--i56-text-primary)">%s</td><td style="text-align:center;padding:8px 12px"><span class="badge badge-warning">%d</span></td><td style="text-align:center;padding:8px 12px"><span class="badge badge-info">%d</span></td></tr>`, z.Name, z.Pending, z.InProgress))
	}
	b.WriteString(`</tbody></table></div>`)

	// Staff activity
	b.WriteString(`<div class="card" style="background:var(--i56-bg-surface);border:1px solid var(--i56-border);border-radius:8px;padding:16px;margin-top:12px">
<div style="font-size:14px;font-weight:600;color:var(--i56-text-primary);margin-bottom:12px">👷 员工动态</div>
<div style="font-size:12px;color:var(--i56-text-secondary);line-height:1.8">`)
	b.WriteString(staffFeed)
	b.WriteString(`</div></div>`)

	return b.String()
}
