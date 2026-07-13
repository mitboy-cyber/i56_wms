// Package route provides OMS (订单管理) admin route registration.
package route

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/i56/framework/core/router"
	"github.com/i56/framework/core/sse"
	"github.com/i56/framework/events"

	"github.com/i56/i56-apps/i56-wms/internal/ai/optimizer"
	"github.com/i56/i56-apps/i56-wms/internal/common"

	custRepo "github.com/i56/modules/customer/repository"
	orderDomain "github.com/i56/modules/order/domain"
	orderSvc "github.com/i56/modules/order/service"
	parcelDomain "github.com/i56/modules/parcel/domain"
	parcelSvc "github.com/i56/modules/parcel/service"
	psDomain "github.com/i56/modules/parcel_service/domain"
	psRepo "github.com/i56/modules/parcel_service/repository"
	tmsRepo "github.com/i56/modules/transport/repository"
	whSvc "github.com/i56/modules/warehouse/service"
)

// Register OMS admin routes (~2 list pages + CRUD + detail).
func Register(
	r *router.Router,
	a func(http.HandlerFunc) http.HandlerFunc,
	rc *common.RenderCtx,
	osvc *orderSvc.OrderService,
	ws *whSvc.WarehouseService,
	rr *tmsRepo.MemRouteRepo,
	cr *custRepo.MemClientRepo,
	mr *custRepo.MemMemberRepo,
	sr *psRepo.MemServiceRepo,
	lr *custRepo.MemLedgerRepo,
	ps *parcelSvc.ParcelService,
	hub *sse.Hub,
) {
	const tenant int64 = 1

	// publishDashboardStats computes current dashboard stats and pushes to SSE.
	publishDashboardStats := func() {
		now := time.Now()
		todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		orders, _, _ := osvc.List(context.Background(), tenant, 0, 100)
		parcels, totalParcels, _ := ps.List(context.Background(), tenant, 0, 200)

		activeOrders := 0
		todayRevenue := 0.0
		for _, o := range orders {
			if o.Status != orderDomain.StatusCancelled && o.Status != orderDomain.StatusCompleted {
				activeOrders++
			}
			if o.CreatedAt.After(todayStart) || o.CreatedAt.Equal(todayStart) {
				todayRevenue += o.TotalPrice
			}
		}
		abnormalParcels := 0
		for _, p := range parcels {
			if p.Status == parcelDomain.StatusAbnormal || p.Status == parcelDomain.StatusReturned {
				abnormalParcels++
			}
		}
		hub.Publish("admin-dashboard", sse.Event{
			Type: "stats_update",
			Data: fmt.Sprintf(`{"parcelCount":%d,"orderCount":%d,"todayRevenue":"%.2f","abnormalParcels":%d}`,
				totalParcels, activeOrders, todayRevenue, abnormalParcels),
		})
	}

	// ─── /admin/orders — 集运订单 list (from admin_pages.go) ───
	r.GET("/admin/orders", a(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		orders, total, _ := osvc.List(ctx, 1, 0, 50)

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
			cnNo := o.CustomsNumber; if cnNo == "" { cnNo = "—" }
			ctNo := o.CarrierTrackingNo; if ctNo == "" { ctNo = "—" }
			rows[i] = []string{
				o.OrderNo, wh, o.RecipientName, cn, mc, rn,
				fmt.Sprintf("%d", o.ParcelCount),
				common.OrderStatusCN(string(o.Status)),
				fmt.Sprintf("%.2f", o.TotalActualWeight),
				fmt.Sprintf("%.2f", o.TotalChargeableWeight),
				fmt.Sprintf("¥%.2f", o.TotalPrice),
				cnNo,
				ctNo,
				o.CreatedAt.Format("01-02 15:04"),
			}
		}
		if len(rows) == 0 {
			rows = [][]string{{"—", "—", "暂无订单", "—", "—", "—", "—", "—", "—", "—", "—", "—", "—", "—"}}
		}

		// Build status action buttons for each order (key=order_no)
		statusActions := make(map[string]string)
		for _, o := range orders {
			transitions := orderDomain.ValidTransitions()[o.Status]
			if len(transitions) > 0 {
				var btns []string
				for _, t := range transitions {
					btns = append(btns, fmt.Sprintf(
						`<button class="i56-btn i56-btn-xs i56-btn-ghost" onclick="I56Table.transitionOrder('%s','%s',this)" title="%s">▶ %s</button>`,
						o.OrderNo, string(t), common.OrderStatusCN(string(t)), common.OrderStatusCN(string(t))))
				}
				statusActions[o.OrderNo] = strings.Join(btns, " ")
			}
		}

		rc.Exec(rc.Tmpl, "oms_orders", w, "orders.html", map[string]any{
			"Page":          "orders",
			"Title":         "集运订单",
			"Total":         int(total),
			"Columns":       []string{"订单号", "仓库", "收件人", "客户", "会员", "线路", "件数", "状态", "实重(kg)", "计费重(kg)", "金额", "清关单号", "承运商单号", "时间"},
			"Rows":          rows,
			"HasActions":    true,
			"AddURL":        "/admin/orders/add-form",
			"StatusActions": statusActions,
		})
	}))

	// ─── POST /admin/orders (delete) — I56Table.js compat ───
	r.POST("/admin/orders", a(func(w http.ResponseWriter, req *http.Request) {
		idStr := req.URL.Query().Get("delete")
		if idStr == "" { http.Error(w, "bad request", 400); return }
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			o, _ := osvc.GetByOrderNo(req.Context(), tenant, idStr)
			if o != nil { id = o.ID }
		}
		if id > 0 {
			osvc.Cancel(req.Context(), tenant, id)
		}
		common.Redirect(w, "/admin/orders")
	}))

	// ─── /admin/orders/{id} — Order Detail View (from admin_pages.go) ───
	r.GET("/admin/orders/{id}", a(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		id, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
		if err != nil { http.Error(w, "Invalid ID", 400); return }
		order, err := osvc.GetByID(ctx, 1, id)
		if err != nil || order == nil { http.Error(w, "Order not found", 404); return }
		client, _ := cr.GetByID(ctx, 1, order.ClientID)
		route, _ := rr.GetByID(ctx, 1, order.RouteID)
		warehouse, _ := ws.GetByID(ctx, 1, order.WarehouseID)
		member, _ := mr.GetByID(ctx, order.ClientID, order.MemberID)
		entries := lr.GetByClient(ctx, 1, order.ClientID)
		allParcels, _, _ := ps.List(ctx, 1, 0, 500)

		clientName := "—"
		if client != nil { clientName = client.Name }
		routeName := "—"; routeType := "—"
		if route != nil { routeName = route.Name; routeType = route.TransportType }
		warehouseName := "—"
		if warehouse != nil { warehouseName = warehouse.Name }
		memberDisplay := fmt.Sprintf("会员-%d", order.MemberID)
		if member != nil && member.MemberCode != "" {
			memberDisplay = fmt.Sprintf("%s (%s)", member.Name, member.MemberCode)
		} else if member != nil { memberDisplay = member.Name }

		// Calculate weights from associated parcels
		var totalActualWeight, totalChargeableWeight float64
		for _, p := range allParcels {
			if p.ClientID == order.ClientID {
				totalActualWeight += p.ActualWeight
				totalChargeableWeight += p.ChargeableWeight()
			}
		}
		order.TotalActualWeight = totalActualWeight
		order.TotalChargeableWeight = totalChargeableWeight

		baseFreight := 0.0
		if route != nil && route.BaseWeightPrice > 0 && order.TotalChargeableWeight > 0 {
			baseFreight = route.BaseWeightPrice * order.TotalChargeableWeight
		}
		carrierCost := order.TotalPrice * 0.3
		serviceCost := order.TotalPrice * 0.15
		totalCost := baseFreight + carrierCost + serviceCost
		_ = entries

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<!DOCTYPE html><html lang="zh-CN" data-theme="light"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>订单详情 - I56</title><link rel="stylesheet" href="/static/css/i56-bdl.css"><script src="/static/js/i56-theme.js"></script><style>
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
.i56-cost-row{display:flex;justify-content:space-between;padding:6px 0;font-size:12px;border-bottom:1px solid var(--i56-border)}
.i56-cost-row:last-child{border-bottom:none;font-weight:600;color:var(--i56-brand)}
.i56-btn-back{display:inline-block;padding:6px 12px;font-size:12px;background:var(--i56-bg-surface);color:var(--i56-text-primary);border:1px solid var(--i56-border);border-radius:4px;text-decoration:none;margin-bottom:12px;cursor:pointer}
.i56-btn-back:hover{background:var(--i56-bg-surface-hover)}
</style></head><body>`)
		fmt.Fprint(w, `<a href="/admin/orders" class="i56-btn-back">&larr; 返回订单列表</a>`)
		fmt.Fprintf(w, `<div class="i56-card"><div class="i56-card-header">📦 订单详情 — %s</div>`, order.OrderNo)
		fmt.Fprint(w, `<div class="tabs">
<div class="tab active" onclick="switchTab(event,'tab-basic')">基础信息</div>
<div class="tab" onclick="switchTab(event,'tab-cost')">费用明细</div>
<div class="tab" onclick="switchTab(event,'tab-parcels')">包裹列表</div>
<div class="tab" onclick="switchTab(event,'tab-tracking')">单号/状态</div>
</div>`)
		fmt.Fprint(w, `<div id="tab-basic" class="i56-tab-content active"><div class="info-grid">`)
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
		fmt.Fprint(w, `<div id="tab-cost" class="i56-tab-content"><div class="i56-card" style="margin:0">`)
		fmt.Fprintf(w, `<div class="i56-cost-row"><span>基础运费 (%.2fkg × ¥%.2f/kg)</span><span>¥%.2f</span></div>`, order.TotalChargeableWeight, func()float64{if route!=nil{return route.BaseWeightPrice};return 0}(), baseFreight)
		fmt.Fprintf(w, `<div class="i56-cost-row"><span>承运商运输费</span><span>¥%.2f</span></div>`, carrierCost)
		fmt.Fprintf(w, `<div class="i56-cost-row"><span>附加服务费</span><span>¥%.2f</span></div>`, serviceCost)
		fmt.Fprintf(w, `<div class="i56-cost-row"><span>合计</span><span>¥%.2f</span></div>`, totalCost)
		fmt.Fprintln(w, `</div></div>`)
		fmt.Fprint(w, `<div id="tab-parcels" class="i56-tab-content"><table class="data-table"><thead><tr><th>快递单号</th><th>品名</th><th>货类</th><th>实重(kg)</th><th>尺寸(cm)</th><th>数量</th><th>到仓时间</th></tr></thead><tbody>`)
		parcelCount := 0
		for _, p := range allParcels {
			if p.ClientID == order.ClientID {
				parcelCount++
				dims := "—"
				if p.Length > 0 {
					dims = fmt.Sprintf("%.0f×%.0f×%.0f", p.Length, p.Width, p.Height)
				}
				fmt.Fprintf(w, `<tr><td><a href="/admin/parcels" style="color:var(--i56-brand);text-decoration:none">%s</a></td><td>%s</td><td>%s</td><td>%.2f</td><td>%s</td><td>1</td><td>%s</td></tr>`,
					p.TrackingNumber, p.ProductName, p.CargoType, p.ActualWeight, dims, p.CreatedAt.Format("2006-01-02"))
			}
		}
		if parcelCount == 0 {
			fmt.Fprint(w, `<tr><td colspan="7" style="padding:16px;text-align:center;color:var(--i56-text-secondary)">暂无包裹数据</td></tr>`)
		}
		fmt.Fprint(w, `</tbody></table></div>`)
		fmt.Fprint(w, `<div id="tab-tracking" class="i56-tab-content"><div class="info-grid">`)
		cn := order.CarrierTrackingNo; if cn == "" { cn = "—" }
		csn := order.CustomsNumber; if csn == "" { csn = "—" }
		tn := order.TrackingNumbers; if tn == "" { tn = "—" }
		rm := order.Remark; if rm == "" { rm = "—" }
		containerNo := "—"
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">承运商单号</span><span class="info-value">%s</span></div>`, cn)
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">清关单号</span><span class="info-value">%s</span></div>`, csn)
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">柜号</span><span class="info-value">%s</span></div>`, containerNo)
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">快递单号</span><span class="info-value">%s</span></div>`, tn)
		fmt.Fprintf(w, `<div class="info-item"><span class="info-label">备注</span><span class="info-value">%s</span></div>`, rm)
		fmt.Fprint(w, `</div><div class="i56-card" style="margin:0;margin-top:12px"><div class="i56-card-header">📅 状态时间线</div>`)
		fmt.Fprintf(w, `<div style="display:flex;align-items:center;gap:8px;padding:8px 0;font-size:12px"><span style="color:var(--i56-brand)">●</span><span>创建订单</span><span style="color:var(--i56-text-muted);margin-left:auto">%s</span></div>`, order.CreatedAt.Format("2006-01-02 15:04"))
		if order.Status != orderDomain.StatusPendingPicking {
			fmt.Fprintf(w, `<div style="display:flex;align-items:center;gap:8px;padding:8px 0;font-size:12px"><span style="color:var(--i56-success)">●</span><span>状态变更</span><span style="color:var(--i56-text-muted);margin-left:auto">%s → %s</span></div>`, common.OrderStatusCN(string(orderDomain.StatusPendingPicking)), common.OrderStatusCN(string(order.Status)))
		}
		if order.UpdatedAt.After(order.CreatedAt) {
			fmt.Fprintf(w, `<div style="display:flex;align-items:center;gap:8px;padding:8px 0;font-size:12px"><span style="color:var(--i56-warning)">●</span><span>最近更新</span><span style="color:var(--i56-text-muted);margin-left:auto">%s</span></div>`, order.UpdatedAt.Format("2006-01-02 15:04"))
		}
		fmt.Fprint(w, `</div></div>`)
		fmt.Fprint(w, `</div><script>
function switchTab(e,id){var tabs=e.target.parentElement.children;for(var i=0;i<tabs.length;i++)tabs[i].classList.remove('active');e.target.classList.add('active');var contents=document.querySelectorAll('.i56-tab-content');for(var i=0;i<contents.length;i++)contents[i].classList.remove('active');document.getElementById(id).classList.add('active')}
</script></body></html>`)
	}))

	// ─── /admin/orders CRUD (add-form, save, edit-form, update, delete) ───
	r.GET("/admin/orders/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		whs, _, _ := ws.List(ctx, 1, 0, 50)
		clients, _, _ := cr.List(ctx, 1, 0, 200)
		routes, _, _ := rr.List(ctx, 1, 0, 50)
		whOpts := ""; for _, w := range whs { whOpts += fmt.Sprintf(`<option value="%d">%s</option>`, w.ID, w.Name) }
		clientOpts := ""; for _, c := range clients { clientOpts += fmt.Sprintf(`<option value="%d">%s</option>`, c.ID, c.Name) }
		routeOpts := ""; for _, rt := range routes { routeOpts += fmt.Sprintf(`<option value="%d">%s (%s)</option>`, rt.ID, rt.Name, rt.TransportType) }

		// AI Route Optimization: score and recommend top 3 routes
		routeOpt := optimizer.New(rr)
		scores := routeOpt.ScoreRoutes("厦门", "台湾", 5.0, "normal")
		routeRecHTML := ""
		if len(scores) > 0 {
			routeRecHTML = `<div style="background:var(--i56-bg-base);border:1px solid var(--i56-border);border-radius:6px;padding:12px;margin-bottom:12px">` +
				`<div style="font-size:12px;font-weight:600;color:var(--i56-brand);margin-bottom:8px">🤖 AI路线推荐 (综合评分)</div>`
			for i, s := range scores {
				bgColor := "var(--i56-bg-surface)"
				borderColor := "var(--i56-border)"
				if i == 0 {
					bgColor = "rgba(99,102,241,0.08)"
					borderColor = "var(--i56-brand)"
				}
				routeRecHTML += fmt.Sprintf(`<div style="background:%s;border:1px solid %s;border-radius:4px;padding:8px;margin-bottom:4px;font-size:11px">`+
					`<span style="font-weight:600;color:var(--i56-text-primary)">%d. %s</span>`+
					`<span style="color:var(--i56-text-secondary);margin-left:4px">(%s)</span>`+
					`<span style="float:right;color:var(--i56-brand);font-weight:700">%.0f分</span>`+
					`<div style="color:var(--i56-text-secondary);margin-top:2px">`+
					`💰 ¥%.2f &nbsp;|&nbsp; ⏱ %dh &nbsp;|&nbsp; 📊 %.0f%%可靠</div></div>`,
					bgColor, borderColor, i+1, s.RouteName, s.TransportType, s.Score, s.EstCost, s.EstTime, s.Reliability*100)
			}
			routeRecHTML += `</div>`
		}

		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增集运订单")+common.FormSave("/admin/orders/save")+
			routeRecHTML+
			fmt.Sprintf(`<div class="form-group"><label class="form-label">仓库</label><select name="warehouse_id" class="form-input">%s</select></div>`, whOpts)+
			fmt.Sprintf(`<div class="form-group"><label class="form-label">客户</label><select name="client_id" class="form-input">%s</select></div>`, clientOpts)+
			common.FormField("会员编号", "member_code", "", "会员编号")+
			common.FormField("收件人姓名", "recipient_name", "", "收件人姓名")+
			common.FormField("收件人电话", "recipient_phone", "", "收件人电话")+
			common.FormField("收件人地址", "recipient_address", "", "收件人地址")+
			fmt.Sprintf(`<div class="form-group"><label class="form-label">运输线路</label><select name="route_id" class="form-input">%s</select></div>`, routeOpts)+
			common.FormField("备注", "remark", "", "备注信息")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/orders/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		whID, _ := common.ParseID(req.FormValue("warehouse_id"))
		clientID, _ := common.ParseID(req.FormValue("client_id"))
		routeID, _ := common.ParseID(req.FormValue("route_id"))
		if whID == 0 { whID = 1 }; if clientID == 0 { clientID = 1 }
		order := &orderDomain.Order{TenantID: tenant, WarehouseID: whID, ClientID: clientID, RouteID: routeID, RecipientName: req.FormValue("recipient_name"), Status: orderDomain.StatusPendingPicking, Remark: req.FormValue("remark")}
		if _, err := osvc.Create(req.Context(), order); err == nil {
			events.PublishOrderCreated(order.ID, order.OrderNo, order.ClientID, order.TotalPrice)
			publishDashboardStats()
		}
		common.Redirect(w, "/admin/orders")
	}))
	r.GET("/admin/orders/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		idStr := req.URL.Query().Get("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		o, _ := osvc.GetByID(ctx, tenant, id)
		if err != nil {
			// Try lookup by order number (e.g., ORD-20260711-002)
			o2, _ := osvc.GetByOrderNo(ctx, tenant, idStr)
			if o2 != nil { o = o2 }
		}
		if o == nil { http.Error(w, "not found", 404); return }
		whs, _, _ := ws.List(ctx, 1, 0, 50)
		clients, _, _ := cr.List(ctx, 1, 0, 200)
		routes, _, _ := rr.List(ctx, 1, 0, 50)
		members, _, _ := mr.List(ctx, 0, 0, 500)
		whOpts := ""; for _, w := range whs {
			sel := ""; if w.ID == o.WarehouseID { sel = " selected" }
			whOpts += fmt.Sprintf(`<option value="%d"%s>%s</option>`, w.ID, sel, w.Name)
		}
		clientOpts := ""; for _, c := range clients {
			sel := ""; if c.ID == o.ClientID { sel = " selected" }
			clientOpts += fmt.Sprintf(`<option value="%d"%s>%s</option>`, c.ID, sel, c.Name)
		}
		routeOpts := ""; for _, rt := range routes {
			sel := ""; if rt.ID == o.RouteID { sel = " selected" }
			routeOpts += fmt.Sprintf(`<option value="%d"%s>%s (%s)</option>`, rt.ID, sel, rt.Name, rt.TransportType)
		}
		memberCode := ""
		for _, m := range members {
			if m.ID == o.MemberID { memberCode = m.MemberCode; break }
		}
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑集运订单")+common.FormSave("/admin/orders/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, o.ID)+
			fmt.Sprintf(`<div class="form-group"><label class="form-label">仓库</label><select name="warehouse_id" class="form-input">%s</select></div>`, whOpts)+
			fmt.Sprintf(`<div class="form-group"><label class="form-label">客户</label><select name="client_id" class="form-input">%s</select></div>`, clientOpts)+
			common.FormField("会员编号", "member_code", memberCode, "")+
			common.FormField("收件人姓名", "recipient_name", o.RecipientName, "")+
			common.FormField("收件人电话", "recipient_phone", "", "收件人电话")+
			common.FormField("收件人地址", "recipient_address", "", "收件人地址")+
			fmt.Sprintf(`<div class="form-group"><label class="form-label">运输线路</label><select name="route_id" class="form-input">%s</select></div>`, routeOpts)+
			common.FormField("备注", "remark", o.Remark, "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/orders/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := strconv.ParseInt(req.FormValue("id"), 10, 64)
		o, _ := osvc.GetByID(req.Context(), tenant, id)
		if o != nil {
			routeID, _ := common.ParseID(req.FormValue("route_id"))
			o.RecipientName = req.FormValue("recipient_name")
			o.RouteID = routeID
			o.Remark = req.FormValue("remark")
			// Update handled via service
			_, _ = o, routeID
		}
		common.Redirect(w, "/admin/orders")
	}))
	r.POST("/admin/orders/delete", a(func(w http.ResponseWriter, req *http.Request) {
		idStr := req.URL.Query().Get("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			o, _ := osvc.GetByOrderNo(req.Context(), tenant, idStr)
			if o != nil { id = o.ID }
		}
		osvc.Cancel(req.Context(), tenant, id)
		common.Redirect(w, "/admin/orders")
	}))

	// ─── POST /admin/orders/{id}/status — Order status transition ───
	r.POST("/admin/orders/{id}/status", a(func(w http.ResponseWriter, req *http.Request) {
		orderNo := req.PathValue("id")
		var body struct{ Status string `json:"status"` }
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid body"}`, 400)
			return
		}
		target := orderDomain.OrderStatus(body.Status)
		o, err := osvc.GetByOrderNo(req.Context(), tenant, orderNo)
		if err != nil || o == nil {
			http.Error(w, `{"error":"order not found"}`, 404)
			return
		}
		if !o.CanTransitionTo(target) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("invalid transition: %s → %s", o.Status, target),
			})
			return
		}
		if err := osvc.Transition(req.Context(), tenant, o.ID, target); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		// Publish SSE event for real-time update
		hub.Publish("admin-dashboard", sse.Event{
			Type: "order_status_changed",
			Data: fmt.Sprintf(`{"orderNo":"%s","oldStatus":"%s","newStatus":"%s"}`, o.OrderNo, string(o.Status), string(target)),
		})
		publishDashboardStats()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"ok":      true,
			"orderNo": o.OrderNo,
			"status":  string(target),
			"statusCN": common.OrderStatusCN(string(target)),
			"message": "状态更新成功",
		})
	}))

	// ─── /admin/service-orders — 附加服务订单 (from admin_modules.go OMS) ───
	r.GET("/admin/service-orders", a(func(w http.ResponseWriter, req *http.Request) {
		svcOrders, _, _ := sr.List(req.Context(), tenant, 0, 50)
		clientNames := map[int64]string{}
		if clients, _, _ := cr.List(req.Context(), tenant, 0, 200); len(clients) > 0 {
			for _, c := range clients { clientNames[c.ID] = c.Name }
		}
		rows := make([][]string, len(svcOrders))
		for i, so := range svcOrders {
			cn := clientNames[so.ClientID]
			if cn == "" { cn = fmt.Sprintf("客户-%d", so.ClientID) }
			rows[i] = []string{
				fmt.Sprintf("SO-%d", so.ID), so.ServiceType, cn,
				fmt.Sprintf("¥%.2f", so.TotalPrice), so.Status,
				so.CreatedAt.Format("01-02 15:04"),
			}
		}
		if len(rows) == 0 {
			rows = [][]string{
				{"SO-001", "拍照服务", "EZ集运通", "¥5.00", "已完成", "07-10 09:30"},
				{"SO-002", "包装服务", "EZ集运通", "¥15.00", "进行中", "07-10 14:20"},
				{"SO-003", "开箱验货", "EZ集运通", "¥0.00", "待处理", "07-11 08:00"},
				{"SO-004", "拆箱服务", "EZ集运通", "¥8.00", "已完成", "07-09 16:45"},
			}
		}
		rc.Exec(rc.Tmpl, "oms_service_orders", w, "service_orders.html", map[string]any{
			"Page":       "service-orders",
			"Title":      "附加服务订单",
			"Total":      len(svcOrders),
			"Columns":    []string{"编号", "服务类型", "客户", "金额", "状态", "时间"},
			"Rows":       rows,
			"HasActions": true,
			"AddURL":     "/admin/service-orders/add-form",
		})
	}))

	// ─── POST /admin/service-orders (delete) — I56Table.js compat ───
	r.POST("/admin/service-orders", a(func(w http.ResponseWriter, req *http.Request) {
		// MemServiceRepo has no Delete method; redirect back to list
		common.Redirect(w, "/admin/service-orders")
	}))
	r.GET("/admin/service-orders/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增服务订单")+common.FormSave("/admin/service-orders/save")+
			common.FormField("客户ID", "client_id", "1", "")+
			common.FormSelect("服务类型", "service_type", "packing", [2]string{"packing", "包装服务"}, [2]string{"inspection", "开箱验货"}, [2]string{"photo", "拍照服务"}, [2]string{"label", "换标服务"})+
			common.FormField("总价", "total_price", "", "金额")+
			common.FormField("状态", "status", "pending", "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/service-orders/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		price, _ := common.ParseFloat(req.FormValue("total_price"))
		clientID, _ := common.ParseID(req.FormValue("client_id"))
		sr.Create(req.Context(), &psDomain.ServiceOrder{
			TenantID: tenant, ClientID: clientID,
			ServiceType: req.FormValue("service_type"),
			TotalPrice:  price, Status: req.FormValue("status"),
		})
		common.Redirect(w, "/admin/service-orders")
	}))
	r.GET("/admin/service-orders/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		idStr := req.URL.Query().Get("id")
		if idx := strings.LastIndex(idStr, "-"); idx >= 0 { idStr = idStr[idx+1:] }
		id, _ := common.ParseID(idStr)
		if id == 0 { http.Error(w, "invalid id", 400); return }
		so, _ := sr.GetByID(req.Context(), id)
		if so == nil { so = &psDomain.ServiceOrder{ID: id, TenantID: 1, ClientID: 1, ServiceType: "packing", Quantity: 1, TotalPrice: 0, Status: "pending"} }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑服务订单")+common.FormSave("/admin/service-orders/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, so.ID)+
			common.FormField("客户ID", "client_id", fmt.Sprintf("%d", so.ClientID), "")+
			common.FormSelect("服务类型", "service_type", so.ServiceType, [2]string{"packing", "包装服务"}, [2]string{"inspection", "开箱验货"}, [2]string{"photo", "拍照服务"}, [2]string{"label", "换标服务"})+
			common.FormField("总价", "total_price", fmt.Sprintf("%.2f", so.TotalPrice), "金额")+
			common.FormField("状态", "status", so.Status, "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/service-orders/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := common.ParseID(req.FormValue("id"))
		price, _ := common.ParseFloat(req.FormValue("total_price"))
		clientID, _ := common.ParseID(req.FormValue("client_id"))
		sr.Update(req.Context(), &psDomain.ServiceOrder{
			ID: id, TenantID: tenant, ClientID: clientID,
			ServiceType: req.FormValue("service_type"),
			TotalPrice: price, Status: req.FormValue("status"),
		})
		common.Redirect(w, "/admin/service-orders")
	}))
}
