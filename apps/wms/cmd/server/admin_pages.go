package main

// DEPRECATED: This entire file is deprecated. All admin page routes have been
// migrated to internal module route packages (omsroute, wmsroute, tmsroute,
// crmroute, finroute, sysroute) which use templates/data_table for rendering.
// The adminPages() function is no longer called from main.go.

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/i56/framework/core/auth"
	"github.com/i56/framework/core/router"

	rbacRepo "github.com/i56/modules/rbac/repository"

	custRepo "github.com/i56/modules/customer/repository"
	custSvc "github.com/i56/modules/customer/service"
	orderDomain "github.com/i56/modules/order/domain"
	orderSvc "github.com/i56/modules/order/service"
	parcelDomain "github.com/i56/modules/parcel/domain"
	parcelSvc "github.com/i56/modules/parcel/service"
	psRepo "github.com/i56/modules/parcel_service/repository"
	pricingRepo "github.com/i56/modules/pricing/repository"
	printRepo "github.com/i56/modules/print/repository"
	sysRepo "github.com/i56/modules/system/repository"
	tmsRepo "github.com/i56/modules/transport/repository"
	whSvc "github.com/i56/modules/warehouse/service"
	twoRepo "github.com/i56/modules/workorder/repository"

	"github.com/i56/i56-apps/i56-wms/internal/common"
)

// DEPRECATED: use internal module route packages (omsroute, wmsroute, etc.)
// and templates/data_table.html instead.
func adminPages(
	tmpl map[string]*template.Template,
	r *router.Router,
	tm *auth.TokenManager,
	ps *parcelSvc.ParcelService,
	osvc *orderSvc.OrderService,
	cs *custSvc.ClientService,
	ws *whSvc.WarehouseService,
	rr *tmsRepo.MemRouteRepo,
	cour *tmsRepo.MemCourierRepo,
	rbac *rbacRepo.MemRBACRepo,
	cr *custRepo.MemClientRepo,
	mr *custRepo.MemMemberRepo,
	dr *custRepo.MemDeclarantRepo,
	ar *custRepo.MemAddressRepo,
	lr *custRepo.MemLedgerRepo,
	rpr *pricingRepo.MemRoutePriceRepo,
	ppr *printRepo.MemPrintRepo,
	sysCfg *sysRepo.MemSystemConfigRepo,
	sr *psRepo.MemServiceRepo,
	wor *twoRepo.MemWorkOrderRepo,
) {
	a := adminOnly(tm)

	// DEPRECATED: gp() was the old pattern before data_table templates.
// Use data_table template directly via RenderCtx.NewGenericList() instead.
// This function is in a deprecated file and no longer called.
gp := func(w http.ResponseWriter, page, title, icon string, total int, cols []string, rows [][]string, addURL ...string) {
		if cols == nil && len(rows) > 0 {
			cols = make([]string, len(rows[0]))
			for i := range cols {
				cols[i] = fmt.Sprintf("列%d", i+1)
			}
		}
		data := map[string]any{
			"Page": page, "Title": title, "Icon": icon, "Total": total,
			"Columns": cols, "Rows": rows, "HasActions": true,
		}
		if len(addURL) > 0 && addURL[0] != "" {
			data["AddURL"] = addURL[0]
		}
		execTpl(tmpl, "generic_list", w, "generic_list.html", data)
	}

	// ===================================================================
	// Dashboard — real KPI stats from repos (mirrored from wmsroute/route.go)
	// ===================================================================
	r.GET("/admin", a(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		now := time.Now()
		todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		parcels, totalParcels, _ := ps.List(ctx, 1, 0, 200)
		orders, _, _ := osvc.List(ctx, 1, 0, 100)
		clients, _, _ := cr.List(ctx, 1, 0, 50)
		routes, _, _ := rr.List(ctx, 1, 0, 50)

		activeOrders := 0; todayRevenue := 0.0; abnormalParcels := 0
		parcelStatusCounts := map[string]int{}
		for _, o := range orders {
			if o.Status != orderDomain.StatusCancelled && o.Status != orderDomain.StatusCompleted { activeOrders++ }
			if o.CreatedAt.After(todayStart) || o.CreatedAt.Equal(todayStart) { todayRevenue += o.TotalPrice }
		}
		for _, p := range parcels {
			parcelStatusCounts[string(p.Status)]++
			if p.Status == parcelDomain.StatusAbnormal || p.Status == parcelDomain.StatusReturned { abnormalParcels++ }
		}
		// Status distribution
		type sd struct { Label string; Count int; Pct float64; Color string }
		totalSI := 0; for _, c := range parcelStatusCounts { totalSI += c }
		var sdList []sd
		colors := []string{"#6366f1","#22c55e","#f59e0b","#3b82f6","#ec4899","#14b8a6","#8b5cf6","#f97316","#06b6d4","#ef4444"}
		cnLabels := map[string]string{"pre_declared":"预报","received":"已入仓","weighed":"已称重","stored":"已上架","picked":"已拣货","packed":"已打包","loaded":"已装柜","shipped":"运输中","delivered":"已签收","abnormal":"异常","returned":"已退货","outbound":"已出货","delivering":"配送中"}
		ci := 0
		for k, c := range parcelStatusCounts { l := cnLabels[k]; if l == "" { l = k }; pct := 0.0; if totalSI > 0 { pct = float64(c)/float64(totalSI)*100 }; sdList = append(sdList, sd{l,c,pct,colors[ci%len(colors)]}); ci++ }
		// 7-day bar chart
		type bar struct { Date string; Count int; Max int }
		var barData []bar; maxC := 0
		for d := 6; d >= 0; d-- {
			day := now.Add(-time.Duration(d)*24*time.Hour); ds := time.Date(day.Year(),day.Month(),day.Day(),0,0,0,0,day.Location()); de := ds.Add(24*time.Hour); cnt := 0
			for _, o := range orders { if (o.CreatedAt.After(ds)||o.CreatedAt.Equal(ds)) && o.CreatedAt.Before(de) { cnt++ } }
			barData = append(barData, bar{day.Format("01/02"),cnt,0}); if cnt>maxC{maxC=cnt}
		}
		for i := range barData { if maxC > 0 { barData[i].Max = maxC } }
		// Activity feed
		type act struct { Icon, Time, Text string }
		var acts []act; ei := 0
		for _, o := range orders { if ei>=8{break}; a:=act{Time:o.CreatedAt.Format("01-02 15:04")}
			switch o.Status{
			case orderDomain.StatusPendingPicking: a.Icon="📋"; a.Text=fmt.Sprintf("新订单 %s 待拣货 (%s)",o.OrderNo,o.RecipientName)
			case orderDomain.StatusInTransit: a.Icon="🚛"; a.Text=fmt.Sprintf("订单 %s 已发运 (%s)",o.OrderNo,o.RecipientName)
			case orderDomain.StatusCompleted: a.Icon="✅"; a.Text=fmt.Sprintf("订单 %s 已签收 (%s)",o.OrderNo,o.RecipientName)
			case orderDomain.StatusPendingLoading: a.Icon="📦"; a.Text=fmt.Sprintf("订单 %s 待装柜 (%s)",o.OrderNo,o.RecipientName)
			case orderDomain.StatusCustomsClearance: a.Icon="🛃"; a.Text=fmt.Sprintf("订单 %s 清关中 (%s)",o.OrderNo,o.RecipientName)
			case orderDomain.StatusLoaded: a.Icon="🏗️"; a.Text=fmt.Sprintf("订单 %s 已装柜 (%s)",o.OrderNo,o.RecipientName)
			default: a.Icon="📋"; a.Text=fmt.Sprintf("订单 %s %s (%s)",o.OrderNo,string(o.Status),o.RecipientName)
			}; acts=append(acts,a); ei++ }
		for _, p := range parcels { if ei>=8{break}; a:=act{Time:p.CreatedAt.Format("01-02 15:04")}
			switch p.Status{case parcelDomain.StatusReceived:a.Icon="📥";a.Text=fmt.Sprintf("包裹 %s 已入仓",p.TrackingNumber);case parcelDomain.StatusStored:a.Icon="🏗️";a.Text=fmt.Sprintf("包裹 %s 已上架",p.TrackingNumber);default:a.Icon="📦";a.Text=fmt.Sprintf("包裹 %s %s",p.TrackingNumber,string(p.Status))}; acts=append(acts,a); ei++ }
		// Recent parcels & orders
		type rd struct { TrackingNumber, ProductName, Status string; ActualWeight float64 }
		pds := func(s parcelDomain.ParcelStatus)string{switch s{case parcelDomain.StatusPreDeclared:return"预报";case parcelDomain.StatusReceived:return"已入仓";case parcelDomain.StatusWeighed:return"已称重";case parcelDomain.StatusStored:return"已上架";case parcelDomain.StatusPicked:return"已拣货";case parcelDomain.StatusPacked:return"已打包";case parcelDomain.StatusLoaded:return"已装柜";case parcelDomain.StatusOutbound:return"已出货";case parcelDomain.StatusShipped,parcelDomain.StatusDelivering:return"运输中";case parcelDomain.StatusDelivered:return"已签收";case parcelDomain.StatusAbnormal:return"异常";case parcelDomain.StatusReturned:return"已退货";default:return string(s)}}
		var recent []rd; for i, p := range parcels { if i>=8{break}; recent=append(recent,rd{p.TrackingNumber,p.ProductName,pds(p.Status),p.ActualWeight}) }
		type rod struct { OrderNo, ClientName string; ParcelCount int; Status string; TotalPrice float64 }
		cm:=map[int64]string{}; for _,c:=range clients{cm[c.ID]=c.Name}
		var recentOrds []rod
		for i,o:=range orders{if i>=8{break};os:=string(o.Status)
			switch o.Status{case orderDomain.StatusPendingPicking:os="待拣货";case orderDomain.StatusPicking:os="拣货中";case orderDomain.StatusPendingPacking:os="待打包";case orderDomain.StatusPendingLoading:os="待装柜";case orderDomain.StatusLoaded:os="已装柜";case orderDomain.StatusInTransit:os="运输中";case orderDomain.StatusCustomsClearance:os="清关中";case orderDomain.StatusOutForDelivery:os="派送中";case orderDomain.StatusCompleted:os="已完成";case orderDomain.StatusCancelled:os="已取消";case orderDomain.StatusShipped:os="已发货"}
			cn:=cm[o.ClientID];if cn==""{cn="EZ集运通"};recentOrds=append(recentOrds,rod{o.OrderNo,cn,o.ParcelCount,os,o.TotalPrice})}
		execTpl(tmpl, "dashboard", w, "dashboard.html", map[string]any{
			"Page":"dashboard","ParcelCount":totalParcels,"OrderCount":activeOrders,
			"ClientCount":len(clients),"RouteCount":len(routes),
			"TodayRevenue":fmt.Sprintf("%.2f",todayRevenue),"AbnormalParcels":abnormalParcels,
			"RecentParcels":recent,"RecentOrders":recentOrds,
			"StatusDistribution":sdList,"BarData":barData,"Activities":acts,
		})
	}))
	r.GET("/admin/warehouse-board", a(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		parcels, totalParcels, _ := ps.List(ctx, 1, 0, 500)
		warehouses, _, _ := ws.List(ctx, 1, 0, 50)

		type whStats struct {
			Name    string
			Total   int
			Stored  int
			Picked  int
			Shipped int
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
		execTpl(tmpl, "dashboard", w, "dashboard.html", map[string]any{
			"Page":           "warehouse_board",
			"Title":          "仓库看板",
			"TotalParcels":   totalParcels,
			"WarehouseStats": whStatsList,
			"ShowBoard":      true,
		})
	}))

	// ===================================================================
	// 订单管理 — real data via BFT56 modules
	// ===================================================================
	r.GET("/admin/orders", a(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		orders, total, _ := osvc.List(ctx, 1, 0, 50)

		// Build name caches for ID resolution
		whName := map[int64]string{}
		if whs, _, _ := ws.List(ctx, 1, 0, 200); len(whs) > 0 {
			for _, wh := range whs { whName[wh.ID] = wh.Name }
		}
		clientName := map[int64]string{}
		if clients, _, _ := cr.List(ctx, 1, 0, 200); len(clients) > 0 {
			for _, c := range clients { clientName[c.ID] = c.Name }
		}
		routeName := map[int64]string{}
		if routes, _, _ := rr.List(ctx, 1, 0, 200); len(routes) > 0 {
			for _, rt := range routes { routeName[rt.ID] = rt.Name }
		}
		memberCode := map[int64]string{}
		if members, _, _ := mr.List(ctx, 0, 0, 500); len(members) > 0 {
			for _, m := range members { memberCode[m.ID] = m.MemberCode }
		}

		rows := make([][]string, len(orders))
		for i, o := range orders {
			wh := whName[o.WarehouseID]; if wh == "" { wh = fmt.Sprintf("仓库-%d", o.WarehouseID) }
			cn := clientName[o.ClientID]; if cn == "" { cn = fmt.Sprintf("客户-%d", o.ClientID) }
			mc := memberCode[o.MemberID]; if mc == "" { mc = fmt.Sprintf("会员-%d", o.MemberID) }
			rn := routeName[o.RouteID]; if rn == "" { rn = fmt.Sprintf("路线-%d", o.RouteID) }
			rows[i] = []string{
				o.OrderNo,
				wh,
				o.RecipientName,
				cn,
				mc,
				rn,
				fmt.Sprintf("%d", o.ParcelCount),
				common.OrderStatusCN(string(o.Status)),
				fmt.Sprintf("%.2f", o.TotalActualWeight),
				fmt.Sprintf("¥%.2f", o.TotalPrice),
				o.CreatedAt.Format("01-02 15:04"),
			}
		}
		if len(rows) == 0 {
			rows = [][]string{{"—", "—", "暂无订单", "—", "—", "—", "—", "—", "—", "—", "—"}}
		}
		gp(w, "oms_orders", "集运订单", "cart-check", int(total), []string{
			"订单号", "仓库", "收件人", "客户", "会员", "线路", "件数", "状态", "实重(kg)", "金额", "时间",
		}, rows)
	}))
	// Order Detail View
	r.GET("/admin/orders/{id}", a(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		id, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid ID", 400)
			return
		}
		order, err := osvc.GetByID(ctx, 1, id)
		if err != nil || order == nil {
			http.Error(w, "Order not found", 404)
			return
		}
		client, _ := cr.GetByID(ctx, 1, order.ClientID)
		route, _ := rr.GetByID(ctx, 1, order.RouteID)
		warehouse, _ := ws.GetByID(ctx, 1, order.WarehouseID)
		member, _ := mr.GetByID(ctx, order.ClientID, order.MemberID)
		allParcels, _, _ := ps.List(ctx, 1, 0, 500)
		entries := lr.GetByClient(ctx, 1, order.ClientID)

		clientName := "—"
		if client != nil {
			clientName = client.Name
		}
		routeName := "—"
		routeType := "—"
		if route != nil {
			routeName = route.Name
			routeType = route.TransportType
		}
		warehouseName := "—"
		if warehouse != nil {
			warehouseName = warehouse.Name
		}
		memberDisplay := fmt.Sprintf("会员-%d", order.MemberID)
		if member != nil && member.MemberCode != "" {
			memberDisplay = fmt.Sprintf("%s (%s)", member.Name, member.MemberCode)
		} else if member != nil {
			memberDisplay = member.Name
		}
		// Compute cost breakdown
		baseFreight := 0.0
		if route != nil && route.BaseWeightPrice > 0 && order.TotalChargeableWeight > 0 {
			baseFreight = route.BaseWeightPrice * order.TotalChargeableWeight
		}
		carrierCost := order.TotalPrice*0.3
		serviceCost := order.TotalPrice*0.15
		totalCost := baseFreight + carrierCost + serviceCost

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<!DOCTYPE html><html lang="zh-CN" data-theme="dark"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>订单详情 - I56</title><link rel="stylesheet" href="/static/css/i56-bdl.css"><script src="/static/js/i56-theme.js"></script><style>
*{margin:0;padding:0;box-sizing:border-box}body{font-family:system-ui,sans-serif;background:var(--i56-bg-base);color:var(--i56-text-primary);padding:16px}
.card{background:var(--i56-bg-surface);border:1px solid var(--i56-border);border-radius:8px;padding:16px;margin-bottom:12px}
.card-header{font-size:14px;font-weight:600;color:var(--i56-text-primary);margin-bottom:12px;padding-bottom:8px;border-bottom:1px solid var(--i56-border)}
.tabs{display:flex;gap:0;margin-bottom:12px;border-bottom:1px solid var(--i56-border)}
.tab{padding:8px 16px;font-size:12px;color:var(--i56-text-secondary);cursor:pointer;border-bottom:2px solid transparent;transition:all .2s}
.tab:hover,.tab.active{color:var(--i56-brand);border-bottom-color:var(--i56-brand)}
.tab-content{display:none}.tab-content.active{display:block}
.info-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(240px,1fr));gap:8px}
.info-item{display:flex;padding:6px 0;font-size:12px}
.info-label{color:var(--i56-text-secondary);min-width:80px;flex-shrink:0}
.info-value{color:var(--i56-text-primary);font-weight:500}
table.data-table{width:100%;border-collapse:collapse;font-size:12px}
table.data-table th{padding:8px 10px;text-align:left;font-weight:600;color:var(--i56-text-secondary);border-bottom:1px solid var(--i56-border);background:var(--i56-bg-base);font-size:11px}
table.data-table td{padding:8px 10px;border-bottom:1px solid var(--i56-border);color:var(--i56-text-primary)}
table.data-table tr:hover td{background:var(--i56-bg-surface-hover)}
.cost-row{display:flex;justify-content:space-between;padding:6px 0;font-size:12px;border-bottom:1px solid var(--i56-border)}
.cost-row:last-child{border-bottom:none;font-weight:600;color:var(--i56-brand)}
.btn-back{display:inline-block;padding:6px 12px;font-size:12px;background:var(--i56-bg-surface);color:var(--i56-text-primary);border:1px solid var(--i56-border);border-radius:4px;text-decoration:none;margin-bottom:12px;cursor:pointer}
.btn-back:hover{background:var(--i56-bg-surface-hover)}
</style></head><body>`)
		fmt.Fprint(w, `<a href="/admin/orders" class="btn-back">&larr; 返回订单列表</a>`)
		fmt.Fprintf(w, `<div class="card"><div class="card-header">📦 订单详情 — %s</div>`, order.OrderNo)
		// Tabs
		fmt.Fprint(w, `<div class="tabs">
<div class="tab active" onclick="switchTab(event,'tab-basic')">基础信息</div>
<div class="tab" onclick="switchTab(event,'tab-cost')">费用明细</div>
<div class="tab" onclick="switchTab(event,'tab-parcels')">包裹列表</div>
<div class="tab" onclick="switchTab(event,'tab-tracking')">单号/状态</div>
</div>`)
		// Tab: Basic Info
		fmt.Fprint(w, `<div id="tab-basic" class="tab-content active"><div class="info-grid">`)
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">订单号</span><span class="info-value">%s</span></div>`, order.OrderNo)
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">客户</span><span class="info-value">%s</span></div>`, clientName)
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">会员</span><span class="info-value">%s</span></div>`, memberDisplay)
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">收件人</span><span class="info-value">%s</span></div>`, order.RecipientName)
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">仓库</span><span class="info-value">%s</span></div>`, warehouseName)
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">线路</span><span class="info-value">%s</span></div>`, routeName)
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">运输方式</span><span class="info-value">%s</span></div>`, routeType)
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">总实重</span><span class="info-value">%.2f kg</span></div>`, order.TotalActualWeight)
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">总计费重</span><span class="info-value">%.2f kg</span></div>`, order.TotalChargeableWeight)
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">总价</span><span class="info-value" style="color:var(--i56-brand);font-weight:700">¥%.2f</span></div>`, order.TotalPrice)
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">状态</span><span class="info-value">%s</span></div>`, common.OrderStatusCN(string(order.Status)))
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">件数</span><span class="info-value">%d</span></div>`, order.ParcelCount)
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">创建时间</span><span class="info-value">%s</span></div>`, order.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">更新时间</span><span class="info-value">%s</span></div>`, order.UpdatedAt.Format("2006-01-02 15:04:05"))
		fmt.Fprint(w, `</div></div>`)

		// Tab: Cost
		fmt.Fprint(w, `<div id="tab-cost" class="tab-content"><div class="card" style="margin:0">`)
		fmt.Fprintf(w, `<div class="cost-row"><span>基础运费 (%.2fkg × ¥%.2f/kg)</span><span>¥%.2f</span></div>`, order.TotalChargeableWeight, func()float64{if route!=nil{return route.BaseWeightPrice};return 0}(), baseFreight)
		fmt.Fprintf(w, `<div class="cost-row"><span>承运商运输费</span><span>¥%.2f</span></div>`, carrierCost)
		fmt.Fprintf(w, `<div class="cost-row"><span>附加服务费</span><span>¥%.2f</span></div>`, serviceCost)
		fmt.Fprintf(w, `<div class="cost-row"><span>合计</span><span>¥%.2f</span></div>`, totalCost)
		fmt.Fprintln(w, `</div></div>`)
		_ = entries

		// Tab: Parcels
		fmt.Fprint(w, `<div id="tab-parcels" class="tab-content"><table class="data-table"><thead><tr><th>快递单号</th><th>品名</th><th>货类</th><th>实重(kg)</th><th>尺寸(cm)</th><th>数量</th><th>到仓时间</th></tr></thead><tbody>`)
		for _, p := range allParcels {
			if p.ClientID == order.ClientID {
				dims := "—"
				if p.Length > 0 {
					dims = fmt.Sprintf("%.0f×%.0f×%.0f", p.Length, p.Width, p.Height)
				}
				fmt.Fprintf(w, `<tr><td><a href="/admin/parcels/%d" style="color:var(--i56-brand);text-decoration:none">%s</a></td><td>%s</td><td>%s</td><td>%.2f</td><td>%s</td><td>1</td><td>%s</td></tr>`,
					p.ID, p.TrackingNumber, p.ProductName, p.CargoType, p.ActualWeight, dims, p.CreatedAt.Format("2006-01-02"))
			}
		}
		fmt.Fprint(w, `</tbody></table></div>`)

		// Tab: Tracking/Status
		fmt.Fprint(w, `<div id="tab-tracking" class="tab-content"><div class="info-grid">`)
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">承运商单号</span><span class="info-value">%s</span></div>`, func()string{if order.CarrierTrackingNo!=""{return order.CarrierTrackingNo};return "—"}())
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">清关单号</span><span class="info-value">%s</span></div>`, func()string{if order.CustomsNumber!=""{return order.CustomsNumber};return "—"}())
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">柜号</span><span class="info-value">—</span></div>`)
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">快递单号</span><span class="info-value">%s</span></div>`, order.TrackingNumbers)
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">备注</span><span class="info-value">%s</span></div>`, func()string{if order.Remark!=""{return order.Remark};return "—"}())
		fmt.Fprint(w, `</div></div>`)

		fmt.Fprint(w, `</div><script>
function switchTab(e,id){var tabs=e.target.parentElement.children;for(var i=0;i<tabs.length;i++)tabs[i].classList.remove('active');e.target.classList.add('active');var contents=document.querySelectorAll('.tab-content');for(var i=0;i<contents.length;i++)contents[i].classList.remove('active');document.getElementById(id).classList.add('active')}
</script></body></html>`)
	}))


	// ===================================================================
	// 仓库管理 — real data from repos
	// ===================================================================
	r.GET("/admin/warehouses", a(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		whs, total, _ := ws.List(ctx, 1, 0, 50)
		parcels, _, _ := ps.List(ctx, 1, 0, 500)
		whParcelCount := map[int64]int{}
		for _, p := range parcels { whParcelCount[p.WarehouseID]++ }
		rows := make([][]string, len(whs))
		for i := range whs {
			rows[i] = []string{
				whs[i].Name, whs[i].Code, whs[i].Address,
				whs[i].Contact, whs[i].Phone,
				fmt.Sprintf("%d件", whParcelCount[whs[i].ID]),
				common.StatusLabelText(whs[i].IsActive),
			}
		}
		gp(w, "wms_warehouses", "仓库列表", "building", int(total), []string{"仓库", "编码", "地址", "联系人", "电话", "包裹数", "状态"}, rows, "/admin/warehouses/add-form")
	}))

	// Warehouse Detail View
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

		// Count parcels in this warehouse
		parcels, totalParcels, _ := ps.List(ctx, 1, 0, 500)
		var whParcels int
		var stored, picked, packed, shipped int
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
		fmt.Fprint(w, `<!DOCTYPE html><html lang="zh-CN" data-theme="dark"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>仓库详情 - I56</title><link rel="stylesheet" href="/static/css/i56-bdl.css"><script src="/static/js/i56-theme.js"></script><style>
*{margin:0;padding:0;box-sizing:border-box}body{font-family:system-ui,sans-serif;background:var(--i56-bg-base);color:var(--i56-text-primary);padding:16px}
.card{background:var(--i56-bg-surface);border:1px solid var(--i56-border);border-radius:8px;padding:16px;margin-bottom:12px}
.card-header{font-size:14px;font-weight:600;color:var(--i56-text-primary);margin-bottom:12px;padding-bottom:8px;border-bottom:1px solid var(--i56-border)}
.section{margin-bottom:16px}
.section-title{font-size:13px;font-weight:600;color:var(--i56-brand);margin-bottom:8px;padding-bottom:6px;border-bottom:1px solid var(--i56-border)}
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
		// Section: Basic Info
		fmt.Fprint(w, `<div class="section"><div class="section-title">📋 基本信息</div><div class="info-grid">`)
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">仓库名</span><span class="info-value">%s</span></div>`, wh.Name)
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">编码</span><span class="info-value">%s</span></div>`, wh.Code)
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">地址</span><span class="info-value">%s</span></div>`, wh.Address)
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">联系人</span><span class="info-value">%s</span></div>`, wh.Contact)
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">电话</span><span class="info-value">%s</span></div>`, wh.Phone)
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">状态</span><span class="info-value">%s</span></div>`, common.StatusLabelText(wh.IsActive))
		fmt.Fprint(w, `</div></div>`)
		// Section: Parcel Stats
		fmt.Fprint(w, `<div class="section"><div class="section-title">📦 包裹统计</div><div class="stat-grid">`)
		fmt.Fprintf(w, `<div class="stat-card"><div class="stat-value">%d</div><div class="stat-label">库存总量</div></div>`, whParcels)
		fmt.Fprintf(w, `<div class="stat-card"><div class="stat-value">%d</div><div class="stat-label">已上架</div></div>`, stored)
		fmt.Fprintf(w, `<div class="stat-card"><div class="stat-value">%d</div><div class="stat-label">已拣货</div></div>`, picked)
		fmt.Fprintf(w, `<div class="stat-card"><div class="stat-value">%d</div><div class="stat-label">已打包</div></div>`, packed)
		fmt.Fprintf(w, `<div class="stat-card"><div class="stat-value">%d</div><div class="stat-label">已出库</div></div>`, shipped)
		fmt.Fprint(w, `</div></div>`)
		fmt.Fprint(w, `</div></body></html>`)
		_ = totalParcels
	}))

	r.GET("/admin/parcels", a(func(w http.ResponseWriter, req *http.Request) {		ctx := req.Context()
		allParcels, _, _ := ps.List(ctx, 1, 0, 500)
		rows := make([][]string, 0, len(allParcels))
		for _, p := range allParcels {
			statusCN := common.ParcelStatusCN(string(p.Status))
			dims := "—"
			if p.Length > 0 { dims = fmt.Sprintf("%.0f×%.0f×%.0f", p.Length, p.Width, p.Height) }
			rows = append(rows, []string{
				p.TrackingNumber, p.ProductName, p.CargoType, statusCN,
				fmt.Sprintf("%.2f", p.ActualWeight), dims,
				"1", p.CreatedAt.Format("2006-01-02"),
			})
		}
		if len(rows) == 0 {
			rows = [][]string{{"SF1234567890", "手机壳", "普货", "预报", "0.35", "15×10×5", "1", "2026-07-11"}}
		}
		execTpl(tmpl, "generic_list", w, "generic_list.html", map[string]any{
			"Page":"parcels","Title":"包裹列表","Total":len(rows),
			"Columns":[]string{"快递单号","品名","货类","状态","实重(kg)","尺寸(cm)","数量","到仓时间"},
			"Rows":rows,"HasActions":true,"AddURL":"/admin/parcels/add-form",
		})
	}))

}
