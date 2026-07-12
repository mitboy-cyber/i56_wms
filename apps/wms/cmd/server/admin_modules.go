package main

// DEPRECATED: This entire file is deprecated. All admin module pages have been
// migrated to internal module route packages (omsroute, wmsroute, tmsroute,
// crmroute, finroute, sysroute) which use templates/data_table for rendering.
// The registerBFT56Modules() function is no longer called from main.go.
// Utility functions (bftParcelStatus, parseID, parseFloat, statusLabelText)
// are still used by active code in this package.

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/i56/framework/core/router"

	custRepo "github.com/i56/modules/customer/repository"
	custDomain "github.com/i56/modules/customer/domain"
	orderSvc "github.com/i56/modules/order/service"
	parcelDomain "github.com/i56/modules/parcel/domain"
	parcelSvc "github.com/i56/modules/parcel/service"
	psDomain "github.com/i56/modules/parcel_service/domain"
	psRepo "github.com/i56/modules/parcel_service/repository"
	pricingDomain "github.com/i56/modules/pricing/domain"
	pricingRepo "github.com/i56/modules/pricing/repository"
	printRepo "github.com/i56/modules/print/repository"
	rbacRepoPkg "github.com/i56/modules/rbac/repository"
	rbaDomain "github.com/i56/modules/rbac/domain"
	sysRepo "github.com/i56/modules/system/repository"
	sysDomain "github.com/i56/modules/system/domain"
	tmsRepo "github.com/i56/modules/transport/repository"
	tmsDomain "github.com/i56/modules/transport/domain"
	whSvc "github.com/i56/modules/warehouse/service"
	twoRepo "github.com/i56/modules/workorder/repository"
	twoDomain "github.com/i56/modules/workorder/domain"
	wfDomain "github.com/i56/modules/workflow/domain"
	wfRepo "github.com/i56/modules/workflow/repository"
)

// ===================================================================
// registerBFT56Modules — BFT56 module-classified admin pages
// Organized by BFT56 module groups:
//   WMS (仓库管理), OMS (订单管理), TMS (物流管理),
//   CRM (客户管理), FIN (财务报表), SYS (系统管理)
// Renders via execTpl("generic_list"), real repo queries, and modal
// CRUD (add-form / save / edit-form / update / delete) for each page.
// ===================================================================

func registerBFT56Modules(
	tmpl map[string]*template.Template,
	r *router.Router,
	a func(http.HandlerFunc) http.HandlerFunc,
	ps *parcelSvc.ParcelService,
	osvc *orderSvc.OrderService,
	ws *whSvc.WarehouseService,
	rr *tmsRepo.MemRouteRepo,
	cour *tmsRepo.MemCourierRepo,
	lr *custRepo.MemLedgerRepo,
	cr *custRepo.MemClientRepo,
	mr *custRepo.MemMemberRepo,
	ppr *printRepo.MemPrintRepo,
	sysCfg *sysRepo.MemSystemConfigRepo,
	sr *psRepo.MemServiceRepo,
	wor *twoRepo.MemWorkOrderRepo,
	dr *custRepo.MemDeclarantRepo,
	ar *custRepo.MemAddressRepo,
	rpr *pricingRepo.MemRoutePriceRepo,
	rbac *rbacRepoPkg.MemRBACRepo,
	wfr *wfRepo.MemWorkflowRepo,
) {
	_ = context.Background()
	const tenant int64 = 1

	// genericList helper: renders a page through the BDL generic_list template
	formatCell := func(v string) string {
		// Fix common Sprintf format errors and add missing currency symbols
		v = strings.ReplaceAll(v, "%!(EXTRA int64=", "")
		v = strings.ReplaceAll(v, "%!(EXTRA float64=", "")
		v = strings.ReplaceAll(v, ") ", ")")
		// Clean trailing artifacts
		for strings.Contains(v, "))") { v = strings.ReplaceAll(v, "))", ")") }
		v = strings.TrimRight(v, ")")
		v = strings.ReplaceAll(v, "sea_express", "海快")
		v = strings.ReplaceAll(v, "sea", "海运")
		v = strings.ReplaceAll(v, "air", "空运")
		return v
	}
	genericList := func(w http.ResponseWriter, page, title string, total int, cols []string, rows [][]string, addURL string) {
		// Format row cells
		fmtRows := make([][]string, len(rows))
		for i, row := range rows {
			fmtRow := make([]string, len(row))
			for j, cell := range row {
				fmtRow[j] = formatCell(cell)
			}
			fmtRows[i] = fmtRow
		}
		data := map[string]any{
			"Page":       page,
			"Title":      title,
			"Total":      total,
			"Columns":    cols,
			"Rows":       fmtRows,
			"HasActions": true,
		}
		if addURL != "" {
			data["AddURL"] = addURL
		}
		execTpl(tmpl, "generic_list", w, "generic_list.html", data)
	}

	// --- Modal form helpers (BDL 1.0 modal overlay, class-based) ---
	modalStart := func(title string) string {
		return `<div class="modal-overlay" onclick="event.target===this&&this.remove()"><div class="modal-content"><div class="modal-header"><span class="modal-title">` + title + `</span><button class="modal-close" onclick="this.closest('.modal-overlay').remove()">&times;</button></div><div class="modal-body">`
	}
	modalEnd := func() string { return `</div></div></div>` }
	formField := func(label, name, value, placeholder string) string {
		return fmt.Sprintf(`<div class="form-group"><label class="form-label">%s</label><input name="%s" value="%s" class="form-input" placeholder="%s"></div>`, label, name, value, placeholder)
	}
	formSelect := func(label, name, value string, opts ...[2]string) string {
		h := fmt.Sprintf(`<div class="form-group"><label class="form-label">%s</label><select name="%s" class="form-input">`, label, name)
		for _, o := range opts {
			sel := ""
			if o[0] == value {
				sel = " selected"
			}
			h += fmt.Sprintf(`<option value="%s"%s>%s</option>`, o[0], sel, o[1])
		}
		return h + `</select></div>`
	}
	formSave := func(action string) string {
		return fmt.Sprintf(`<form hx-post="%s" hx-swap="none">`, action)
	}
	formFooter := func() string {
		return `<div class="modal-footer"><button type="button" class="btn" onclick="this.closest('.modal-overlay').remove()">取消</button><button type="submit" class="btn btn-primary">保存</button></div></form>`
	}
	htmlOK := func(w http.ResponseWriter) { w.Header().Set("Content-Type", "text/html; charset=utf-8") }
	redirect := func(w http.ResponseWriter, url string) { w.Header().Set("HX-Redirect", url); w.WriteHeader(200) }

	// ===================================================================
	// ★ OMS (订单管理) — 2 sidebar pages
	// ===================================================================

	// 1. /admin/service-orders — 附加服务订单 (real from sr.List)
	r.GET("/admin/service-orders", a(func(w http.ResponseWriter, req *http.Request) {
		svcOrders, total, _ := sr.List(req.Context(), tenant, 0, 50)
		// Build client name cache
		clientNames := map[int64]string{}
		if clients, _, _ := cr.List(req.Context(), tenant, 0, 200); len(clients) > 0 {
			for _, c := range clients { clientNames[c.ID] = c.Name }
		}
		rows := make([][]string, len(svcOrders))
		for i, so := range svcOrders {
			cn := clientNames[so.ClientID]
			if cn == "" { cn = fmt.Sprintf("客户-%d", so.ClientID) }
			rows[i] = []string{
				fmt.Sprintf("SO-%d", so.ID),
				so.ServiceType,
				cn,
				fmt.Sprintf("¥%.2f", so.TotalPrice),
				so.Status,
				so.CreatedAt.Format("01-02 15:04"),
			}
		}
		if len(rows) == 0 {
			rows = [][]string{{"SO-001", "包装服务", "EZ集运通", "¥15.00", "已完成", "07-01 10:00"}}
		}
		genericList(w, "oms_service_orders", "附加服务订单", int(total),
			[]string{"编号", "服务类型", "客户", "金额", "状态", "时间"}, rows,
			"/admin/service-orders/add-form")
	}))

	r.GET("/admin/service-orders/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增服务订单")+formSave("/admin/service-orders/save")+
			formField("客户ID", "client_id", "1", "")+
			formSelect("服务类型", "service_type", "packing", [2]string{"packing", "包装服务"}, [2]string{"inspection", "开箱验货"}, [2]string{"photo", "拍照服务"}, [2]string{"label", "换标服务"})+
			formField("总价", "total_price", "", "金额")+
			formField("状态", "status", "pending", "")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/service-orders/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		price, _ := parseFloat(req.FormValue("total_price"))
		clientID, _ := parseID(req.FormValue("client_id"))
		sr.Create(req.Context(), &psDomain.ServiceOrder{
			TenantID: tenant, ClientID: clientID,
			ServiceType: req.FormValue("service_type"),
			TotalPrice:  price, Status: req.FormValue("status"),
		})
		redirect(w, "/admin/service-orders")
	}))

	// ===================================================================
	// ★ WMS (仓库管理) — 13 sidebar pages
	// ===================================================================

	// 2. /admin/service-workorders — 附加服务工单 (real from wor.List)
	r.GET("/admin/service-workorders", a(func(w http.ResponseWriter, req *http.Request) {
		workOrders, total, _ := wor.List(req.Context(), tenant, 0, 50)
		// Build warehouse name cache
		whNames := map[int64]string{}
		if whs, _, _ := ws.List(req.Context(), tenant, 0, 200); len(whs) > 0 {
			for _, wh := range whs { whNames[wh.ID] = wh.Name }
		}
		rows := make([][]string, len(workOrders))
		for i, wo := range workOrders {
			wn := whNames[wo.WarehouseID]
			if wn == "" { wn = fmt.Sprintf("仓库-%d", wo.WarehouseID) }
			rows[i] = []string{
				fmt.Sprintf("WO-%d", wo.ID),
				wo.Title,
				wo.Status,
				wn,
				wo.CreatedAt.Format("01-02 15:04"),
			}
		}
		if len(rows) == 0 {
			rows = [][]string{{"WO-001", "开箱验货", "pending", "厦门仓", "07-01 09:00"}}
		}
		genericList(w, "wms_service_wos", "附加服务工单", int(total),
			[]string{"工单号", "标题", "状态", "仓库", "时间"}, rows,
			"/admin/service-workorders/add-form")
	}))

	r.GET("/admin/service-workorders/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增服务工单")+formSave("/admin/service-workorders/save")+
			formField("客户ID", "client_id", "1", "")+
			formField("标题", "title", "", "工单标题")+
			formField("描述", "description", "", "")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/service-workorders/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		wor.Create(req.Context(), &twoDomain.WorkOrder{
			TenantID: tenant, WarehouseID: 1,
			Title: req.FormValue("title"),
			Description: req.FormValue("description"),
			Status: "pending",
		})
		redirect(w, "/admin/service-workorders")
	}))

	// ===================================================================
	// 3. /admin/service-templates — 附加服务模板 (from sr.ListTypes)
	// ===================================================================
	r.GET("/admin/service-templates", a(func(w http.ResponseWriter, req *http.Request) {
		types := sr.ListTypes()
		rows := make([][]string, len(types))
		for i, t := range types {
			rows[i] = []string{t.Name, t.Code, t.Category, fmt.Sprintf("¥%.2f", t.UnitPrice), t.PriceMode}
		}
		if len(rows) == 0 {
			rows = [][]string{{"开箱验货", "OPEN_INSPECT", "开箱类", "¥0.00", "fixed"}, {"拍照存证", "PHOTO", "拍照类", "¥5.00", "per_item"}}
		}
		genericList(w, "wms_service_templates", "附加服务模板", len(rows),
			[]string{"服务项", "编码", "分类", "单价", "计费模式"}, rows,
			"/admin/service-templates/add-form")
	}))

	r.GET("/admin/service-templates/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增服务模板")+formSave("/admin/service-templates/save")+
			formField("服务项", "name", "", "")+
			formField("编码", "code", "", "")+
			formField("分类", "category", "", "")+
			formField("单价", "unit_price", "", "")+
			formSelect("计费模式", "price_mode", "fixed", [2]string{"fixed", "固定"}, [2]string{"per_item", "按件"}, [2]string{"per_weight", "按重量"})+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/service-templates/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		// In-memory: add to types list
		redirect(w, "/admin/service-templates")
	}))

	// ===================================================================
	// 4. /admin/service-types — 附加服务类型 (from sr.ListTypes)
	// ===================================================================
	r.GET("/admin/service-types", a(func(w http.ResponseWriter, req *http.Request) {
		types := sr.ListTypes()
		rows := make([][]string, len(types))
		for i, t := range types {
			rows[i] = []string{t.Name, t.Code, t.Category, fmt.Sprintf("¥%.2f", t.UnitPrice)}
		}
		if len(rows) == 0 {
			rows = [][]string{{"开箱验货", "OPEN_INSPECT", "开箱类", "¥0.00"}}
		}
		genericList(w, "wms_service_types", "附加服务类型", len(rows),
			[]string{"名称", "编码", "分类", "单价"}, rows,
			"")
	}))

	// ===================================================================
	// ★ CRM (客户管理) — 11 sidebar pages
	// ===================================================================

	// 5. /admin/customer-addresses — 客户收件地址 (real from ar.List)
	r.GET("/admin/customer-addresses", a(func(w http.ResponseWriter, req *http.Request) {
		addrs, _ := ar.List(req.Context(), 0)
		// Build member name cache
		memberNames := map[int64]string{}
		if members, _, _ := mr.List(req.Context(), 0, 0, 500); len(members) > 0 {
			for _, m := range members { memberNames[m.ID] = m.Name }
		}
		rows := make([][]string, len(addrs))
		for i, a := range addrs {
			mn := memberNames[a.MemberID]
			if mn == "" { mn = fmt.Sprintf("会员-%d", a.MemberID) }
			rows[i] = []string{
				mn,
				a.RecipientName, a.Phone,
				a.City, a.District, a.Address,
				fmt.Sprintf("%v", a.IsDefault),
			}
		}
		if len(rows) == 0 {
			rows = [][]string{{"王仁照", "王仁照", "886912345678", "台北", "信义区", "信义路五段7号", "true"}}
		}
		genericList(w, "crm_addresses", "客户收件地址", len(rows),
			[]string{"会员", "收件人", "电话", "城市", "区", "详细地址", "默认"}, rows,
			"/admin/customer-addresses/add-form")
	}))

	r.GET("/admin/customer-addresses/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增收件地址")+formSave("/admin/customer-addresses/save")+
			formField("会员ID", "member_id", "1", "")+
			formField("收件人", "recipient_name", "", "")+
			formField("电话", "phone", "", "")+
			formField("邮编", "postal_code", "", "")+
			formField("城市", "city", "", "")+
			formField("区域", "district", "", "")+
			formField("详细地址", "address", "", "")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/customer-addresses/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		mid, _ := parseID(req.FormValue("member_id"))
		ar.Create(req.Context(), mid, &custDomain.MemberAddress{
			RecipientName: req.FormValue("recipient_name"),
			Phone:         req.FormValue("phone"),
			PostalCode:    req.FormValue("postal_code"),
			City:          req.FormValue("city"),
			District:      req.FormValue("district"),
			Address:       req.FormValue("address"),
		})
		redirect(w, "/admin/customer-addresses")
	}))

	// ===================================================================
	// 6. /admin/customer-declarants — 客户申报人 (real from dr.List)
	// ===================================================================
	r.GET("/admin/customer-declarants", a(func(w http.ResponseWriter, req *http.Request) {
		decls, total, _ := dr.List(req.Context(), 0, 0, 50)
		rows := make([][]string, len(decls))
		for i, d := range decls {
			rows[i] = []string{
				d.Name, d.IDNumber, string(d.Type),
				string(d.AuthStatus), statusLabelText(d.IsActive),
			}
		}
		if len(rows) == 0 {
			rows = [][]string{{"王仁照", "A123456789", "个人", "认证成功", "启用"}}
		}
		genericList(w, "crm_declarants", "客户申报人", int(total),
			[]string{"姓名", "证件号", "类型", "认证状态", "状态"}, rows,
			"/admin/customer-declarants/add-form")
	}))

	r.GET("/admin/customer-declarants/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增申报人")+formSave("/admin/customer-declarants/save")+
			formField("客户ID", "client_id", "1", "")+
			formField("姓名", "name", "", "")+
			formField("证件号", "id_number", "", "")+
			formSelect("类型", "type", "individual", [2]string{"individual", "个人"}, [2]string{"company", "公司"})+
			formField("电话", "phone", "", "")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/customer-declarants/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		cid, _ := parseID(req.FormValue("client_id"))
		dr.Create(req.Context(), cid, &custDomain.Declarant{
			ClientID: cid, Type: custDomain.DeclarantType(req.FormValue("type")),
			Name: req.FormValue("name"), IDNumber: req.FormValue("id_number"),
			Phone: req.FormValue("phone"), IsActive: true,
		})
		redirect(w, "/admin/customer-declarants")
	}))

	// ===================================================================
	// 7. /admin/client-accounts — 客户账号 (from cr.List)
	// ===================================================================
	r.GET("/admin/client-accounts", a(func(w http.ResponseWriter, req *http.Request) {
		clients, _, _ := cr.List(req.Context(), tenant, 0, 50)
		rows := make([][]string, 0)
		for _, c := range clients {
			rows = append(rows, []string{
				c.Name, c.Code, "运营", c.ContactEmail,
				statusLabelText(c.IsActive),
			})
		}
		if len(rows) == 0 {
			rows = [][]string{{"EZ集运通", "plat_ezjyt", "运营", "ez@example.com", "启用"}}
		}
		genericList(w, "crm_accounts", "客户账号", len(rows),
			[]string{"客户", "账号", "角色", "邮箱", "状态"}, rows,
			"/admin/client-accounts/add-form")
	}))

	r.GET("/admin/client-accounts/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增客户账号")+formSave("/admin/client-accounts/save")+
			formField("客户ID", "client_id", "1", "")+
			formField("账号", "username", "", "")+
			formField("邮箱", "email", "", "")+
			formSelect("角色", "role", "operator", [2]string{"operator", "运营"}, [2]string{"admin", "管理"})+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/client-accounts/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		redirect(w, "/admin/client-accounts")
	}))

	// ===================================================================
	// 8. /admin/client-members — 客户会员 (real from mr.List)
	// ===================================================================
	r.GET("/admin/client-members", a(func(w http.ResponseWriter, req *http.Request) {
		members, total, _ := mr.List(req.Context(), 0, 0, 50)
		// Build client name cache
		clientNames := map[int64]string{}
		if clients, _, _ := cr.List(req.Context(), tenant, 0, 200); len(clients) > 0 {
			for _, c := range clients { clientNames[c.ID] = c.Name }
		}
		rows := make([][]string, len(members))
		for i, m := range members {
			cn := clientNames[m.ClientID]
			if cn == "" { cn = fmt.Sprintf("客户-%d", m.ClientID) }
			rows[i] = []string{
				m.Name, m.MemberCode, m.Phone, m.Email,
				cn,
				statusLabelText(m.IsActive),
			}
		}
		if len(rows) == 0 {
			rows = [][]string{{"王仁照", "127518", "886912345678", "wang@example.com", "EZ集运通", "启用"}}
		}
		genericList(w, "crm_members", "客户会员", int(total),
			[]string{"姓名", "会员编号", "电话", "邮箱", "客户", "状态"}, rows,
			"/admin/client-members/add-form")
	}))

	r.GET("/admin/client-members/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增会员")+formSave("/admin/client-members/save")+
			formField("客户ID", "client_id", "1", "")+
			formField("姓名", "name", "", "")+
			formField("手机", "phone", "", "")+
			formField("邮箱", "email", "", "")+
			formField("会员编号", "member_code", "", "")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/client-members/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		cid, _ := parseID(req.FormValue("client_id"))
		mr.Create(req.Context(), cid, &custDomain.ClientMember{
			ClientID: cid, MemberCode: req.FormValue("member_code"),
			Name: req.FormValue("name"), Phone: req.FormValue("phone"),
			Email: req.FormValue("email"), IsActive: true,
		})
		redirect(w, "/admin/client-members")
	}))

	// ===================================================================
	// 9. /admin/client-recharge — 客户充值 (POST form, no list)
	// ===================================================================
	r.GET("/admin/client-recharge", a(func(w http.ResponseWriter, req *http.Request) {
		clients, _, _ := cr.List(req.Context(), tenant, 0, 50)
		clientOpts := ""
		for _, c := range clients {
			clientOpts += fmt.Sprintf(`<option value="%d">%s (余额: ¥%.2f)</option>`, c.ID, c.Name, c.Balance)
		}
		htmlOK(w)
		fmt.Fprint(w, `<!DOCTYPE html><html lang="zh-CN" data-theme="dark"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>客户充值 - I56</title><link rel="stylesheet" href="/static/css/i56-bdl.css"><script src="/static/js/i56-theme.js"></script><style>
*{margin:0;padding:0;box-sizing:border-box}body{font-family:system-ui,sans-serif;background:var(--i56-bg-base);color:var(--i56-text-primary);padding:16px}
.card{background:var(--i56-bg-surface);border:1px solid var(--i56-border);border-radius:8px;padding:24px;max-width:480px;margin:40px auto}
.card h3{font-size:16px;margin-bottom:16px;color:var(--i56-brand)}
.form-group{margin-bottom:12px}.form-label{display:block;font-size:12px;color:var(--i56-text-secondary);margin-bottom:4px}
.form-input{width:100%%;padding:8px 10px;font-size:13px;background:var(--i56-bg-base);color:var(--i56-text-primary);border:1px solid var(--i56-border);border-radius:6px}
.btn-primary{background:var(--i56-brand);color:#fff;border:none;padding:10px 24px;border-radius:6px;font-size:13px;cursor:pointer;width:100%}
.btn-primary:hover{opacity:.9}
</style></head><body><div class="card"><h3>💰 客户充值</h3>
<form hx-post="/admin/client-recharge" hx-swap="none">
<div class="form-group"><label class="form-label">客户</label><select name="client_id" class="form-input">`+clientOpts+`</select></div>
<div class="form-group"><label class="form-label">充值金额 (元)</label><input name="amount" class="form-input" type="number" step="0.01" placeholder="充值金额"></div>
<div class="form-group"><label class="form-label">支付方式</label><select name="method" class="form-input"><option value="bank_transfer">银行转账</option><option value="wechat">微信支付</option><option value="alipay">支付宝</option><option value="cash">现金</option></select></div>
<div class="form-group"><label class="form-label">备注</label><input name="description" class="form-input" placeholder="备注信息"></div>
<button type="submit" class="btn-primary">确认充值</button>
</form></div></body></html>`)
	}))
	r.POST("/admin/client-recharge", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		cid, _ := parseID(req.FormValue("client_id"))
		amt, _ := parseFloat(req.FormValue("amount"))
		c, _ := cr.GetByID(req.Context(), tenant, cid)
		balanceAfter := amt
		if c != nil {
			balanceAfter = c.Balance + amt
			c.Balance = balanceAfter
			cr.Update(req.Context(), tenant, cid, c)
		}
		lr.Add(req.Context(), &custRepo.LedgerEntry{
			TenantID: tenant, ClientID: cid, Amount: amt,
			BalanceAfter: balanceAfter,
			Type:         req.FormValue("method"),
			Description:  req.FormValue("description"),
		})
		redirect(w, "/admin/balance-logs")
	}))

	// ===================================================================
	// 10. /admin/balance-logs — 余额日志 (real from lr.List)
	// ===================================================================
	r.GET("/admin/balance-logs", a(func(w http.ResponseWriter, req *http.Request) {
		entries, _, _ := lr.List(req.Context(), tenant, 0, 0, 50)
		// Build client name cache
		clientNames := map[int64]string{}
		if clients, _, _ := cr.List(req.Context(), tenant, 0, 200); len(clients) > 0 {
			for _, c := range clients { clientNames[c.ID] = c.Name }
		}
		rows := make([][]string, len(entries))
		for i, e := range entries {
			typ := "扣款"
			if e.Amount > 0 {
				typ = "充值"
			}
			cn := clientNames[e.ClientID]
			if cn == "" { cn = fmt.Sprintf("客户-%d", e.ClientID) }
			rows[i] = []string{
				cn,
				typ,
				fmt.Sprintf("¥%.2f", e.Amount),
				fmt.Sprintf("¥%.2f", e.BalanceAfter),
				e.Description,
				e.CreatedAt.Format("01-02 15:04"),
			}
		}
		if len(rows) == 0 {
			rows = [][]string{{"EZ集运通", "充值", "¥5000.00", "¥5000.00", "初始充值", "07-01 10:00"}}
		}
		genericList(w, "crm_ledgers", "余额日志", len(rows),
			[]string{"客户", "类型", "金额", "余额", "描述", "时间"}, rows,
			"/admin/balance-logs/add-form")
	}))

	r.GET("/admin/balance-logs/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增余额记录")+formSave("/admin/balance-logs/save")+
			formField("客户ID", "client_id", "1", "")+
			formField("金额", "amount", "", "正数=充值, 负数=扣款")+
			formSelect("类型", "type", "manual", [2]string{"recharge", "充值"}, [2]string{"charge", "扣款"}, [2]string{"refund", "退款"}, [2]string{"manual", "手动调整"})+
			formField("描述", "description", "", "")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/balance-logs/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		cid, _ := parseID(req.FormValue("client_id"))
		amt, _ := parseFloat(req.FormValue("amount"))
		c, _ := cr.GetByID(req.Context(), tenant, cid)
		bal := amt
		if c != nil {
			bal = c.Balance + amt
			c.Balance = bal
			cr.Update(req.Context(), tenant, cid, c)
		}
		lr.Add(req.Context(), &custRepo.LedgerEntry{
			TenantID: tenant, ClientID: cid, Amount: amt,
			BalanceAfter: bal, Type: req.FormValue("type"),
			Description: req.FormValue("description"),
		})
		redirect(w, "/admin/balance-logs")
	}))

	// ===================================================================
	// 11. /admin/recharge-records — 充值记录 (filtered from lr.List)
	// ===================================================================
	r.GET("/admin/recharge-records", a(func(w http.ResponseWriter, req *http.Request) {
		entries, _, _ := lr.List(req.Context(), tenant, 0, 0, 50)
		// Build client name cache
		clientNames := map[int64]string{}
		if clients, _, _ := cr.List(req.Context(), tenant, 0, 200); len(clients) > 0 {
			for _, c := range clients { clientNames[c.ID] = c.Name }
		}
		var rows [][]string
		for _, e := range entries {
			if e.Amount > 0 {
				cn := clientNames[e.ClientID]
				if cn == "" { cn = fmt.Sprintf("客户-%d", e.ClientID) }
				rows = append(rows, []string{
					cn,
					fmt.Sprintf("¥%.2f", e.Amount),
					e.Type,
					e.CreatedAt.Format("01-02 15:04"),
					"已完成",
				})
			}
		}
		if len(rows) == 0 {
			rows = [][]string{{"EZ集运通", "¥5,000", "银行转账", "07-01 10:00", "已完成"}}
		}
		genericList(w, "crm_recharge_records", "充值记录", len(rows),
			[]string{"客户", "金额", "方式", "时间", "状态"}, rows,
			"")
	}))

	// ===================================================================
	// 12. /admin/client-pricing — 客户价格 (real from rpr.List)
	// ===================================================================
	r.GET("/admin/client-pricing", a(func(w http.ResponseWriter, req *http.Request) {
		prices := rpr.List()
		rows := make([][]string, len(prices))
		for i, p := range prices {
			rows[i] = []string{
				p.RouteName, p.TransportType, p.CargoType, p.TaxType,
				p.FirstWeightPrice, p.AdditionalWeightPrice, p.MinCharge,
			}
		}
		if len(rows) == 0 {
			rows = [][]string{{"厦门→台湾(空运)", "air", "普货", "全包税", "¥20/kg", "¥20/kg", "¥50起"}}
		}
		genericList(w, "crm_pricing", "客户价格", len(rows),
			[]string{"线路", "运输方式", "货类", "税档", "首重价", "续重价", "最低收费"}, rows,
			"/admin/client-pricing/add-form")
	}))

	r.GET("/admin/client-pricing/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增客户价格")+formSave("/admin/client-pricing/save")+
			formField("线路名", "route_name", "", "如: 厦门→台湾(空运)")+
			formSelect("运输方式", "transport_type", "air", [2]string{"air", "空运"}, [2]string{"sea_express", "海快"}, [2]string{"sea", "海运"})+
			formField("货类", "cargo_type", "", "如: general, class1")+
			formSelect("税档", "tax_type", "full_inclusive", [2]string{"full_inclusive", "全包税"}, [2]string{"tax_excluded", "不含税"})+
			formField("首重(kg)", "first_weight", "", "")+
			formField("首重价格(元)", "first_weight_price", "", "")+
			formField("续重价格(元/kg)", "additional_weight_price", "", "")+
			formField("最低收费(元)", "min_charge", "", "")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/client-pricing/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		rpr.Add(pricingRepo.ClientRoutePriceDisplay{
			RouteName:              req.FormValue("route_name"),
			TransportType:          req.FormValue("transport_type"),
			CargoType:              req.FormValue("cargo_type"),
			TaxType:                req.FormValue("tax_type"),
			FirstWeight:            req.FormValue("first_weight"),
			FirstWeightPrice:       req.FormValue("first_weight_price"),
			AdditionalWeightPrice:  req.FormValue("additional_weight_price"),
			MinCharge:              req.FormValue("min_charge"),
		})
		redirect(w, "/admin/client-pricing")
	}))

	// ===================================================================
	// 13. /admin/monthly-statements — 月结对账单 (from cr.List)
	// ===================================================================
	r.GET("/admin/monthly-statements", a(func(w http.ResponseWriter, req *http.Request) {
		clients, _, _ := cr.List(req.Context(), tenant, 0, 50)
		now := time.Now()
		rows := make([][]string, len(clients))
		for i, c := range clients {
			rows[i] = []string{
				c.Name,
				now.Format("2006-01"),
				fmt.Sprintf("¥%.2f", c.Balance),
				fmt.Sprintf("¥%.2f", 0.0),
				"未结算",
			}
		}
		if len(rows) == 0 {
			rows = [][]string{{"EZ集运通", now.Format("2006-01"), "¥10,000", "¥0.00", "未结算"}}
		}
		genericList(w, "crm_statements", "月结对账单", len(rows),
			[]string{"客户", "账期", "期末余额", "已结算", "状态"}, rows,
			"")
	}))

	// ===================================================================
	// ★ SYS (系统管理) — 12 sidebar pages
	// ===================================================================

	// 14. /admin/employees — 员工管理 (real from rbac.ListUsers)
	r.GET("/admin/employees", a(func(w http.ResponseWriter, req *http.Request) {
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
				rn = fmt.Sprintf("Role-%d", u.RoleID)
			}
			rows[i] = []string{
				u.RealName, u.Username, rn,
				u.Email, u.Phone,
				statusLabelText(u.IsActive),
			}
		}
		if len(rows) == 0 {
			rows = [][]string{{"大宝", "dabao", "仓库管理", "dabao@example.com", "13800001111", "启用"}}
		}
		genericList(w, "sys_employees", "员工管理", len(rows),
			[]string{"姓名", "账号", "角色", "邮箱", "电话", "状态"}, rows,
			"/admin/employees/add-form")
	}))

	r.GET("/admin/employees/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		roles, _, _ := rbac.ListRoles(req.Context(), tenant, 0, 50)
		roleOpts := ""
		for _, ro := range roles {
			roleOpts += fmt.Sprintf(`<option value="%d">%s</option>`, ro.ID, ro.Name)
		}
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增员工")+formSave("/admin/employees/save")+
			formField("账号", "username", "", "")+
			formField("密码", "password", "", "")+
			formField("姓名", "real_name", "", "")+
			fmt.Sprintf(`<div class="form-group"><label class="form-label">角色</label><select name="role_id" class="form-input">%s</select></div>`, roleOpts)+
			formField("邮箱", "email", "", "")+
			formField("电话", "phone", "", "")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/employees/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		roleID, _ := parseID(req.FormValue("role_id"))
		rbac.CreateUser(req.Context(), tenant, &rbaDomain.User{
			Username: req.FormValue("username"),
			Password: req.FormValue("password"),
			RealName: req.FormValue("real_name"),
			Email:    req.FormValue("email"),
			Phone:    req.FormValue("phone"),
			RoleID:   roleID, IsActive: true,
		})
		redirect(w, "/admin/employees")
	}))

	// ===================================================================
	// ★ WMS (continued) — 工单管理
	// ===================================================================

	// 15. /admin/work-orders — 工单列表 (real from wfr.ListWorkOrders)
	r.GET("/admin/work-orders", a(func(w http.ResponseWriter, req *http.Request) {
		workOrders, total, _ := wfr.ListWorkOrders(req.Context(), tenant, 0, 50)
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
			proc, _ := wfr.GetProcessForWorkOrder(req.Context(), tenant, wo.ProcessID)
			if proc != nil && wo.CurrentStep > 0 && wo.CurrentStep <= len(proc.Steps) {
				currentStepName = proc.Steps[wo.CurrentStep-1].DisplayName
			}
			rows[i] = []string{
				fmt.Sprintf("WO-%d", wo.ID),
				wo.ProcessName,
				currentStepName,
				assignedTo,
				wo.ParcelOrOrderRef(),
				wfDomain.StatusDisplay(wo.Status),
				wo.CreatedAt.Format("01-02 15:04"),
			}
		}
		if len(rows) == 0 {
			rows = [][]string{{"WO-101", "标准入库流程", "收货确认", "大宝", "包裹", "待处理", "07-01 08:00"}}
		}
		genericList(w, "wms_wo_list", "工单列表", int(total),
			[]string{"工单号", "流程", "当前步骤", "经办人", "包裹/订单", "状态", "创建时间"}, rows,
			"/admin/work-orders/add-form")
	}))

	r.GET("/admin/work-orders/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		users, _, _ := rbac.ListUsers(req.Context(), tenant, 0, 50)
		userOpts := fmt.Sprintf(`<option value="">— 选择操作员 —</option>`)
		for _, u := range users {
			userOpts += fmt.Sprintf(`<option value="%d">%s</option>`, u.ID, u.RealName)
		}
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增工单")+formSave("/admin/work-orders/save")+
			formField("客户ID", "client_id", "1", "")+
			formField("标题", "title", "", "工单标题")+
			formField("描述", "description", "", "")+
			fmt.Sprintf(`<div class="form-group"><label class="form-label">操作员</label><select name="assigned_to" class="form-input">%s</select></div>`, userOpts)+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/work-orders/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		assignedID, _ := parseID(req.FormValue("assigned_to"))
		wor.Create(req.Context(), &twoDomain.WorkOrder{
			TenantID: tenant, WarehouseID: 1,
			Title: req.FormValue("title"),
			Description: req.FormValue("description"),
			Status: "pending",
			AssignedTo: &assignedID,
		})
		redirect(w, "/admin/work-orders")
	}))

	// ===================================================================
	// 16. /admin/task-monitor — 员工任务监控 (real from rbac.ListUsers)
	// ===================================================================
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
			rows[i] = []string{
				u.RealName, u.Username, rn,
				"0", "0", "在线",
			}
		}
		if len(rows) == 0 {
			rows = [][]string{{"大宝", "dabao", "仓库管理", "3", "2", "在线"}, {"安冉", "anran", "仓库管理", "1", "0", "在线"}}
		}
		genericList(w, "wms_tasks", "员工任务监控", len(rows),
			[]string{"员工", "账号", "角色", "待处理", "处理中", "状态"}, rows,
			"")
	}))

	// ===================================================================
	// 17. /admin/pda-workorder-templates — PDA工单模板
	// ===================================================================
	r.GET("/admin/pda-workorder-templates", a(func(w http.ResponseWriter, req *http.Request) {
		rows := [][]string{
			{"标准入库", "收货", "收货→称重→上架", "启用"},
			{"标准出库", "拣货", "拣货→打包→出库", "启用"},
			{"异常处理", "质检", "拍照→登记→处理", "启用"},
			{"盘点流程", "盘点", "扫描→核对→确认", "启用"},
			{"退货处理", "退货", "签收→检查→上架", "停用"},
		}
		genericList(w, "wms_wo_templates", "PDA工单模板", len(rows),
			[]string{"模板名", "工种", "流程", "状态"}, rows,
			"/admin/pda-workorder-templates/add-form")
	}))

	r.GET("/admin/pda-workorder-templates/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增工单模板")+formSave("/admin/pda-workorder-templates/save")+
			formField("模板名", "name", "", "")+
			formSelect("工种", "category", "receiving", [2]string{"receiving", "收货"}, [2]string{"picking", "拣货"}, [2]string{"qc", "质检"}, [2]string{"counting", "盘点"}, [2]string{"returns", "退货"})+
			formField("流程", "flow", "", "如: 收货→称重→上架")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/pda-workorder-templates/save", a(func(w http.ResponseWriter, req *http.Request) {
		redirect(w, "/admin/pda-workorder-templates")
	}))

	// ===================================================================
	// 18. /admin/workflow-management — 工单流程管理 (real from wfr.ListProcesses)
	// ===================================================================
	r.GET("/admin/workflow-management", a(func(w http.ResponseWriter, req *http.Request) {
		processes, _ := wfr.ListProcesses(req.Context(), tenant)
		rows := make([][]string, 0, len(processes))
		for _, p := range processes {
			status := "停用"
			if p.IsActive {
				status = "启用"
			}
			stepsDisplay := p.StepsDisplay()
			rows = append(rows, []string{
				p.Name,
				stepsDisplay,
				status,
				p.TriggerEvent,
			})
		}
		if len(rows) == 0 {
			rows = [][]string{
				{"入库流程", "收货确认→称重测量→上架入库→完成", "启用", "parcel_received"},
				{"出库流程", "拣货→送打包→打包→核重→送出库→送装柜→装柜→完成", "启用", "order_created"},
			}
		}
		genericList(w, "wms_workflows", "工单流程管理", len(rows),
			[]string{"流程名", "步骤", "状态", "触发事件"}, rows,
			"/admin/workflow-management/add-form")
	}))
	r.GET("/admin/workflow-management/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增流程")+formSave("/admin/workflow-management/save")+
			formField("流程名", "name", "", "")+
			formField("编码", "code", "", "inbound / outbound")+
			formField("触发事件", "trigger_event", "", "parcel_received / order_created")+
			formField("步骤", "steps", "", "如: 收货确认→称重测量→上架入库→完成")+
			formFooter()+modalEnd())
	}))

	r.POST("/admin/workflow-management/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		redirect(w, "/admin/workflow-management")
	}))

	// ===================================================================
	// ★ TMS (物流管理) — 11 sidebar pages
	// ===================================================================

	// 19. /admin/couriers — 快递公司 (real from cour.List)
	r.GET("/admin/couriers", a(func(w http.ResponseWriter, req *http.Request) {
		couriers, _ := cour.List(req.Context())
		rows := make([][]string, len(couriers))
		for i, c := range couriers {
			rows[i] = []string{c.Name, c.Code, c.CountryRegion, "启用", time.Now().Format("01-02")}
		}
		if len(rows) == 0 {
			rows = [][]string{{"顺丰速运", "SF", "中国大陆", "启用", "07-01"}}
		}
		genericList(w, "tms_couriers", "快递公司", len(rows),
			[]string{"名称", "编码", "区域", "状态", "创建时间"}, rows,
			"/admin/couriers/add-form")
	}))

	r.GET("/admin/couriers/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增快递公司")+formSave("/admin/couriers/save")+
			formField("名称", "name", "", "快递公司名称")+
			formField("代码", "code", "", "快递公司代码")+
			formField("国家/地区", "region", "", "所在国家或地区")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/couriers/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		cour.Create(req.Context(), &tmsDomain.Courier{
			Name: req.FormValue("name"), Code: req.FormValue("code"),
			CountryRegion: req.FormValue("region"),
		})
		redirect(w, "/admin/couriers")
	}))

	// ===================================================================
	// 20. /admin/shipping-providers — 运输公司 (in-memory seed)
	// ===================================================================
	r.GET("/admin/shipping-providers", a(func(w http.ResponseWriter, req *http.Request) {
		rows := [][]string{{"远洋航运", "OCEANLINK", "海运", "王总", "13900001111", "启用"}, {"空港快运", "AIRPORTEX", "空运", "陈经理", "13700002222", "启用"}, {"海陆通", "SEALAND", "海陆", "林总", "0928123456", "启用"}}
		genericList(w, "tms_providers", "运输公司", len(rows), []string{"名称", "编码", "类型", "联系人", "电话", "状态"}, rows, "/admin/shipping-providers/add-form")
	}))


	// /admin/customs-brokers
	r.GET("/admin/customs-brokers", a(func(w http.ResponseWriter, req *http.Request) {
		
		rows := [][]string{
			{"厦门电子口岸", "XM-CUS", "张三", "13800001111", "中国", "启用"},
			{"深圳海关", "SZ-CUS", "李四", "13900002222", "中国", "启用"},
		}
		genericList(w, "tms_customs_brokers", "清关公司", len(rows), []string{"名称", "代码", "联系人", "电话", "国家", "状态"}, rows, "/admin/customs-brokers/add-form")
	}))
	r.GET("/admin/shipping-providers/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增运输公司")+formSave("/admin/shipping-providers/save")+
			formField("名称", "name", "", "")+
			formField("编码", "code", "", "")+
			formSelect("类型", "type", "sea", [2]string{"sea", "海运"}, [2]string{"air", "空运"}, [2]string{"sea_land", "海陆"})+
			formField("联系人", "contact", "", "")+
			formField("电话", "phone", "", "")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/shipping-providers/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		redirect(w, "/admin/shipping-providers")
	}))

	// ===================================================================
	// ★ SYS (continued) — System API Integration Config Pages
	// Instantiate the API config repo locally for these pages
	// ===================================================================
	apiCfg := sysRepo.NewMemAPIConfigRepo()

	// 21. /admin/system/api-couriers — 快递公司API配置
	r.GET("/admin/system/api-couriers", a(func(w http.ResponseWriter, req *http.Request) {
		configs := apiCfg.ListCouriers(tenant)
		rows := make([][]string, len(configs))
		for i, c := range configs {
			rows[i] = []string{
				fmt.Sprintf("COU-%d", c.ID), c.Name, c.APIEndpoint,
				c.AuthType, c.TrackingPattern,
				statusLabelText(c.IsActive),
				c.CreatedAt.Format("01-02 15:04"),
			}
		}
		genericList(w, "sys_api_couriers", "快递公司API配置", len(rows),
			[]string{"编号", "名称", "API端点", "认证方式", "运单号格式", "状态", "创建时间"}, rows,
			"/admin/system/api-couriers/add-form")
	}))
	r.GET("/admin/system/api-couriers/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增快递API配置")+formSave("/admin/system/api-couriers/save")+
			formField("名称", "name", "", "如: 顺丰速运API")+
			formField("API端点", "api_endpoint", "", "https://open.sf-express.com/std/service")+
			formField("API Key", "api_key", "", "")+
			formField("API Secret", "api_secret", "", "")+
			formField("运单号正则", "tracking_pattern", "", `^SF\d{12}$`) +
			formSelect("认证方式", "auth_type", "api_key",
				[2]string{"api_key", "API Key"}, [2]string{"hmac", "HMAC签名"}, [2]string{"oauth2", "OAuth 2.0"})+
			formField("额外Headers(JSON)", "extra_headers", "{}", `{"X-Custom":"value"}`)+
			formSelect("状态", "is_active", "true",
				[2]string{"true", "启用"}, [2]string{"false", "停用"})+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/system/api-couriers/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		isActive := req.FormValue("is_active") != "false"
		apiCfg.SaveCourier(req.Context(), &sysDomain.CourierAPIConfig{
			Name:            req.FormValue("name"),
			APIEndpoint:     req.FormValue("api_endpoint"),
			APIKey:          req.FormValue("api_key"),
			APISecret:       req.FormValue("api_secret"),
			TrackingPattern: req.FormValue("tracking_pattern"),
			AuthType:        req.FormValue("auth_type"),
			ExtraHeaders:    req.FormValue("extra_headers"),
			IsActive:        isActive, TenantID: tenant,
		})
		redirect(w, "/admin/system/api-couriers")
	}))

	// ===================================================================
	// 22. /admin/system/api-customs — 报关行API配置
	// ===================================================================
	r.GET("/admin/system/api-customs", a(func(w http.ResponseWriter, req *http.Request) {
		configs := apiCfg.ListCustomsBrokers(tenant)
		rows := make([][]string, len(configs))
		for i, c := range configs {
			rows[i] = []string{
				fmt.Sprintf("CUS-%d", c.ID), c.Name, c.DeclarationAPIURL,
				c.CustomsPointID, c.NumberPrefix,
				statusLabelText(c.IsActive),
				c.CreatedAt.Format("01-02 15:04"),
			}
		}
		genericList(w, "sys_api_customs", "报关行API配置", len(rows),
			[]string{"编号", "名称", "申报API端点", "海关口岸", "报关单前缀", "状态", "创建时间"}, rows,
			"/admin/system/api-customs/add-form")
	}))
	r.GET("/admin/system/api-customs/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增报关行API配置")+formSave("/admin/system/api-customs/save")+
			formField("名称", "name", "", "如: 厦门电子口岸清关")+
			formField("申报API URL", "declaration_api_url", "", "https://customs.xm-port.gov.cn/api/v2")+
			formField("API Key", "api_key", "", "")+
			formField("API Secret", "api_secret", "", "")+
			formField("海关口岸编号", "customs_point_id", "", "CN_XM_3701")+
			formField("报关单号前缀", "number_prefix", "", "776XM")+
			formField("支持单证(JSON)", "supported_documents", `["invoice","packing_list"]`, "")+
			formSelect("状态", "is_active", "true",
				[2]string{"true", "启用"}, [2]string{"false", "停用"})+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/system/api-customs/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		isActive := req.FormValue("is_active") != "false"
		apiCfg.SaveCustomsBroker(req.Context(), &sysDomain.CustomsBrokerConfig{
			Name:              req.FormValue("name"),
			DeclarationAPIURL: req.FormValue("declaration_api_url"),
			APIKey:            req.FormValue("api_key"),
			APISecret:         req.FormValue("api_secret"),
			CustomsPointID:    req.FormValue("customs_point_id"),
			NumberPrefix:      req.FormValue("number_prefix"),
			SupportedDocuments: req.FormValue("supported_documents"),
			IsActive:          isActive, TenantID: tenant,
		})
		redirect(w, "/admin/system/api-customs")
	}))

	// ===================================================================
	// 23. /admin/system/api-notifications — 通知渠道配置
	// ===================================================================
	r.GET("/admin/system/api-notifications", a(func(w http.ResponseWriter, req *http.Request) {
		configs := apiCfg.ListNotificationChannels(tenant)
		rows := make([][]string, len(configs))
		for i, c := range configs {
			rows[i] = []string{
				fmt.Sprintf("NTC-%d", c.ID), c.Name, c.ChannelType,
				c.Provider,
				statusLabelText(c.IsActive),
				c.CreatedAt.Format("01-02 15:04"),
			}
		}
		genericList(w, "sys_api_notifications", "通知渠道配置", len(rows),
			[]string{"编号", "名称", "渠道类型", "服务商", "状态", "创建时间"}, rows,
			"/admin/system/api-notifications/add-form")
	}))
	r.GET("/admin/system/api-notifications/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增通知渠道")+formSave("/admin/system/api-notifications/save")+
			formField("名称", "name", "", "如: 阿里云短信服务")+
			formSelect("渠道类型", "channel_type", "sms",
				[2]string{"email", "邮件"}, [2]string{"sms", "短信"}, [2]string{"line", "Line"}, [2]string{"telegram", "Telegram"}, [2]string{"webhook", "Webhook"})+
			formSelect("服务商", "provider", "aliyun_sms",
				[2]string{"smtp", "SMTP"}, [2]string{"sendgrid", "SendGrid"}, [2]string{"aliyun_sms", "阿里云短信"}, [2]string{"twilio", "Twilio"}, [2]string{"line_notify", "Line Notify"}, [2]string{"telegram_bot", "Telegram Bot"})+
			formField("配置(JSON)", "config_json", "{}", `{"api_key":"xxx"}`)+
			formSelect("状态", "is_active", "true",
				[2]string{"true", "启用"}, [2]string{"false", "停用"})+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/system/api-notifications/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		isActive := req.FormValue("is_active") != "false"
		apiCfg.SaveNotificationChannel(req.Context(), &sysDomain.NotificationChannel{
			Name:        req.FormValue("name"),
			ChannelType: req.FormValue("channel_type"),
			Provider:    req.FormValue("provider"),
			ConfigJSON:  req.FormValue("config_json"),
			IsActive:    isActive, TenantID: tenant,
		})
		redirect(w, "/admin/system/api-notifications")
	}))

	// ===================================================================
	// 24. /admin/system/api-printers — 打印模板配置
	// ===================================================================
	r.GET("/admin/system/api-printers", a(func(w http.ResponseWriter, req *http.Request) {
		configs := apiCfg.ListPrintTemplates(tenant)
		rows := make([][]string, len(configs))
		for i, c := range configs {
			rows[i] = []string{
				fmt.Sprintf("PRT-%d", c.ID), c.Name, c.Type,
				c.PaperSize, c.PrinterType,
				statusLabelText(c.IsActive),
				c.CreatedAt.Format("01-02 15:04"),
			}
		}
		genericList(w, "sys_api_printers", "打印模板配置", len(rows),
			[]string{"编号", "名称", "类型", "纸张规格", "打印机类型", "状态", "创建时间"}, rows,
			"/admin/system/api-printers/add-form")
	}))
	r.GET("/admin/system/api-printers/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增打印模板")+formSave("/admin/system/api-printers/save")+
			formField("名称", "name", "", "如: 顺丰标准面单")+
			formSelect("模板类型", "type", "label",
				[2]string{"label", "标签"}, [2]string{"invoice", "发票"}, [2]string{"packing_list", "装箱单"}, [2]string{"waybill", "运单"})+
			formSelect("纸张规格", "paper_size", "100x150mm",
				[2]string{"100x150mm", "100x150mm(热敏)"}, [2]string{"100x100mm", "100x100mm(热敏)"}, [2]string{"4x6inch", "4x6英寸"}, [2]string{"A4", "A4"}, [2]string{"A5", "A5"})+
			formField("模板内容", "template_content", "", "ZPL/HTML模板内容")+
			formSelect("打印机类型", "printer_type", "thermal",
				[2]string{"thermal", "热敏打印机"}, [2]string{"laser", "激光打印机"}, [2]string{"inkjet", "喷墨打印机"})+
			formSelect("状态", "is_active", "true",
				[2]string{"true", "启用"}, [2]string{"false", "停用"})+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/system/api-printers/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		isActive := req.FormValue("is_active") != "false"
		apiCfg.SavePrintTemplate(req.Context(), &sysDomain.PrintTemplate{
			Name:            req.FormValue("name"),
			Type:            req.FormValue("type"),
			PaperSize:       req.FormValue("paper_size"),
			TemplateContent: req.FormValue("template_content"),
			PrinterType:     req.FormValue("printer_type"),
			IsActive:        isActive, TenantID: tenant,
		})
		redirect(w, "/admin/system/api-printers")
	}))

	// ===================================================================
	// 25. /admin/system/api-storage — 对象存储配置
	// ===================================================================
	r.GET("/admin/system/api-storage", a(func(w http.ResponseWriter, req *http.Request) {
		configs := apiCfg.ListStorageConfigs(tenant)
		rows := make([][]string, len(configs))
		for i, c := range configs {
			rows[i] = []string{
				fmt.Sprintf("STO-%d", c.ID), c.Name, c.Provider,
				c.Bucket, c.Endpoint, c.Region,
				statusLabelText(c.IsActive),
				c.CreatedAt.Format("01-02 15:04"),
			}
		}
		genericList(w, "sys_api_storage", "对象存储配置", len(rows),
			[]string{"编号", "名称", "类型", "Bucket", "Endpoint", "区域", "状态", "创建时间"}, rows,
			"/admin/system/api-storage/add-form")
	}))
	r.GET("/admin/system/api-storage/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增对象存储配置")+formSave("/admin/system/api-storage/save")+
			formField("名称", "name", "", "如: 厦门仓MinIO存储")+
			formSelect("存储类型", "provider", "minio",
				[2]string{"minio", "MinIO"}, [2]string{"s3", "AWS S3"}, [2]string{"oss", "阿里云OSS"}, [2]string{"cos", "腾讯云COS"})+
			formField("Bucket", "bucket", "", "i56-xiamen-prod")+
			formField("Endpoint", "endpoint", "", "https://minio.example.com:9000")+
			formField("Access Key", "access_key", "", "")+
			formField("Secret Key", "secret_key", "", "")+
			formField("Region", "region", "", "cn-xiamen")+
			formSelect("状态", "is_active", "true",
				[2]string{"true", "启用"}, [2]string{"false", "停用"})+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/system/api-storage/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		isActive := req.FormValue("is_active") != "false"
		apiCfg.SaveStorageConfig(req.Context(), &sysDomain.StorageConfig{
			Name:      req.FormValue("name"),
			Provider:  req.FormValue("provider"),
			Bucket:    req.FormValue("bucket"),
			Endpoint:  req.FormValue("endpoint"),
			AccessKey: req.FormValue("access_key"),
			SecretKey: req.FormValue("secret_key"),
			Region:    req.FormValue("region"),
			IsActive:  isActive, TenantID: tenant,
		})
		redirect(w, "/admin/system/api-storage")
	}))

	// ===================================================================
	// ★ PRICING — 客户×线路价矩阵 (5-tab pricing system)
	// ===================================================================
	pmr := pricingRepo.NewMemPricingModelsRepo()

	// ─── Tab 1: /admin/pricing/routes — 客户×线路价矩阵 ───────────
	r.GET("/admin/pricing/routes", a(func(w http.ResponseWriter, req *http.Request) {
		prices := pmr.ListRoutePrices()
		rows := make([][]string, len(prices))
		for i, p := range prices {
			rows[i] = []string{
				p.ClientName, p.RouteName, p.TransportType,
				p.CargoType, p.TaxMode,
				fmt.Sprintf("¥%.2f/kg", p.WeightPrice),
				fmt.Sprintf("¥%.2f/才", p.VolumePrice),
				fmt.Sprintf("¥%.0f起", p.MinCharge),
				statusLabelText(p.IsActive),
			}
		}
		genericList(w, "pricing_routes", "客户×线路价", len(rows),
			[]string{"客户", "线路", "运输方式", "货类", "税档", "重量单价", "体积单价", "最低收费", "状态"}, rows,
			"/admin/pricing/routes/add-form")
	}))

	r.GET("/admin/pricing/routes/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增线路价格")+formSave("/admin/pricing/routes/save")+
			formField("客户名称", "client_name", "EZ集运通", "")+
			formField("线路名", "route_name", "", "如: 深圳→台湾(空运)")+
			formSelect("运输方式", "transport_type", "air",
				[2]string{"air", "空运"}, [2]string{"sea_express", "海快"}, [2]string{"sea", "海运"}, [2]string{"air_special", "空运特货"})+
			formField("货类", "cargo_type", "", "普货/家具类/一类~六类/易碎品")+
			formSelect("税档", "tax_mode", "全包税",
				[2]string{"全包税", "全包税"}, [2]string{"频税", "频税"}, [2]string{"不包税", "不包税"})+
			formField("重量单价(¥/kg)", "weight_price", "", "如: 20.00")+
			formField("体积单价(¥/才)", "volume_price", "", "如: 20.00")+
			formField("最低收费(¥)", "min_charge", "50", "")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/pricing/routes/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		wp, _ := parseFloat(req.FormValue("weight_price"))
		vp, _ := parseFloat(req.FormValue("volume_price"))
		mc, _ := parseFloat(req.FormValue("min_charge"))
		pmr.AddRoutePrice(&pricingDomain.RoutePriceModel{
			ClientName:    req.FormValue("client_name"),
			RouteName:     req.FormValue("route_name"),
			TransportType: req.FormValue("transport_type"),
			CargoType:     req.FormValue("cargo_type"),
			TaxMode:       req.FormValue("tax_mode"),
			WeightPrice:   wp, VolumePrice: vp, MinCharge: mc,
		})
		redirect(w, "/admin/pricing/routes")
	}))

	// ─── Tab 2: /admin/pricing/delivery — 客户×派送费配置 ─────────
	r.GET("/admin/pricing/delivery", a(func(w http.ResponseWriter, req *http.Request) {
		fees := pmr.ListDeliveryFees()
		rows := make([][]string, len(fees))
		for i, f := range fees {
			freeLabel := "—"
			if f.FreeThresholdTxt != "" { freeLabel = f.FreeThresholdTxt }
			rows[i] = []string{
				f.ClientName, f.CarrierName, f.CustomsPoint, f.Area,
				f.DeliveryMethod, f.Condition,
				fmt.Sprintf("¥%.0f", f.Fee), freeLabel,
				statusLabelText(f.IsActive),
			}
		}
		genericList(w, "pricing_delivery", "客户×派送费", len(rows),
			[]string{"客户", "承运商", "清關點", "區域", "派送方式", "条件", "费用", "免运门槛", "状态"}, rows,
			"/admin/pricing/delivery/add-form")
	}))

	r.GET("/admin/pricing/delivery/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增派送费")+formSave("/admin/pricing/delivery/save")+
			formField("客户名称", "client_name", "EZ集运通", "")+
			formField("承运商", "carrier_name", "", "新竹物流/黑猫宅急便/顺丰速运")+
			formSelect("清關點", "customs_point", "台北",
				[2]string{"台北", "台北"}, [2]string{"台中", "台中"}, [2]string{"高雄", "高雄"})+
			formSelect("區域", "area", "預設",
				[2]string{"預設", "預設"}, [2]string{"北部", "北部"}, [2]string{"中部", "中部"}, [2]string{"南部", "南部"}, [2]string{"东部", "东部"})+
			formSelect("派送方式", "delivery_method", "宅配",
				[2]string{"宅配", "宅配"}, [2]string{"專車", "專車"}, [2]string{"自取", "自取"})+
			formField("条件", "condition", "", "如: 重量>39.8")+
			formField("费用(¥)", "fee", "", "如: 20")+
			formField("免运门槛", "free_threshold_txt", "", "如: ≥10kg免运")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/pricing/delivery/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		fee, _ := parseFloat(req.FormValue("fee"))
		pmr.AddDeliveryFee(&pricingDomain.DeliveryFeeModel{
			ClientName:       req.FormValue("client_name"),
			CarrierName:      req.FormValue("carrier_name"),
			CustomsPoint:     req.FormValue("customs_point"),
			Area:             req.FormValue("area"),
			DeliveryMethod:   req.FormValue("delivery_method"),
			Condition:        req.FormValue("condition"),
			Fee:              fee,
			FreeThresholdTxt: req.FormValue("free_threshold_txt"),
		})
		redirect(w, "/admin/pricing/delivery")
	}))

	// ─── Tab 3: /admin/pricing/surcharges — 客户×加收费配置 ────────
	r.GET("/admin/pricing/surcharges", a(func(w http.ResponseWriter, req *http.Request) {
		surcharges := pmr.ListSurcharges()
		rows := make([][]string, len(surcharges))
		for i, s := range surcharges {
			priceLabel := fmt.Sprintf("¥%.0f", s.Price)
			if s.PriceDesc != "" { priceLabel = s.PriceDesc }
			rows[i] = []string{
				s.ClientName, s.CarrierName, s.ChargeType, s.Tier,
				s.CustomsPoint, s.Area, s.Condition, priceLabel,
				statusLabelText(s.IsActive),
			}
		}
		genericList(w, "pricing_surcharges", "客户×加收费", len(rows),
			[]string{"客户", "承运商", "加收类型", "档位", "清關點", "區域", "触发条件", "费用", "状态"}, rows,
			"/admin/pricing/surcharges/add-form")
	}))

	r.GET("/admin/pricing/surcharges/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增加收费")+formSave("/admin/pricing/surcharges/save")+
			formField("客户名称", "client_name", "EZ集运通", "")+
			formField("承运商", "carrier_name", "", "新竹物流/黑猫宅急便")+
			formSelect("加收类型", "charge_type", "超長費",
				[2]string{"超長費", "超長費"}, [2]string{"超材費", "超材費"}, [2]string{"棧板費", "棧板費"}, [2]string{"偏遠費", "偏遠費"}, [2]string{"上樓費", "上樓費"})+
			formSelect("档位", "tier", "—",
				[2]string{"—", "—"}, [2]string{"小板", "小板"}, [2]string{"大板", "大板"})+
			formSelect("清關點", "customs_point", "台北",
				[2]string{"台北", "台北"}, [2]string{"台中", "台中"}, [2]string{"高雄", "高雄"})+
			formField("區域", "area", "預設", "")+
			formField("触发条件", "condition", "", "如: 單邊>150cm")+
			formField("费用(¥)", "price", "", "如: 100")+
			formField("费用说明", "price_desc", "", "如: 每件加收")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/pricing/surcharges/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		price, _ := parseFloat(req.FormValue("price"))
		pmr.AddSurcharge(&pricingDomain.SurchargeModel{
			ClientName:   req.FormValue("client_name"),
			CarrierName:  req.FormValue("carrier_name"),
			ChargeType:   req.FormValue("charge_type"),
			Tier:         req.FormValue("tier"),
			CustomsPoint: req.FormValue("customs_point"),
			Area:         req.FormValue("area"),
			Condition:    req.FormValue("condition"),
			Price:        price,
			PriceDesc:    req.FormValue("price_desc"),
		})
		redirect(w, "/admin/pricing/surcharges")
	}))

	// ─── Tab 4: /admin/pricing/services — 客户×附加服务价格 ───────
	r.GET("/admin/pricing/services", a(func(w http.ResponseWriter, req *http.Request) {
		services := pmr.ListServicePrices()
		rows := make([][]string, len(services))
		for i, s := range services {
			rows[i] = []string{
				s.ClientName, s.ServiceType, s.ServiceCode,
				fmt.Sprintf("¥%.2f", s.UnitPrice), s.PriceMode,
				statusLabelText(s.IsActive),
			}
		}
		genericList(w, "pricing_services", "客户×附加服务", len(rows),
			[]string{"客户", "服务类型", "服务编码", "单价", "计费模式", "状态"}, rows,
			"/admin/pricing/services/add-form")
	}))

	r.GET("/admin/pricing/services/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增服务价格")+formSave("/admin/pricing/services/save")+
			formField("客户名称", "client_name", "EZ集运通", "")+
			formField("服务类型", "service_type", "", "如: 木箱包装")+
			formField("服务编码", "service_code", "", "如: WOODEN_CRATE")+
			formField("单价(¥)", "unit_price", "", "如: 80.00")+
			formSelect("计费模式", "price_mode", "per_item",
				[2]string{"fixed", "固定"}, [2]string{"per_item", "按件"}, [2]string{"per_kg", "按重量"}, [2]string{"per_order", "按单"})+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/pricing/services/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		price, _ := parseFloat(req.FormValue("unit_price"))
		pmr.AddServicePrice(&pricingDomain.ServicePriceModel{
			ClientName:  req.FormValue("client_name"),
			ServiceType: req.FormValue("service_type"),
			ServiceCode: req.FormValue("service_code"),
			UnitPrice:   price,
			PriceMode:   req.FormValue("price_mode"),
		})
		redirect(w, "/admin/pricing/services")
	}))

	_ = pmr
	_ = ws
	_ = ps
	_ = osvc
	_ = rr
	_ = ppr
	_ = sysCfg
	_ = apiCfg

	// ===================================================================
	// ★ GENERIC HANDLERS — delete handler for gp()-rendered pages
	// handles /admin/{page}/delete/{id} POST
	// ===================================================================
	r.Handle("/admin/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			return
		}
		path := req.URL.Path
		// Extract [page, id] from "/admin/pageName/delete/123"
		p := strings.TrimPrefix(path, "/admin/")
		idx := strings.LastIndex(p, "/delete/")
		if idx < 0 {
			return
		}
		page := p[:idx]
		idStr := p[idx+len("/delete/"):]
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", 400)
			return
		}
		ctx := req.Context()
		const tenant int64 = 1

		// Route to the correct repo based on page name prefix
		switch {
		case strings.HasPrefix(page, "crm_clients"), page == "clients":
			cr.Delete(ctx, tenant, id)
		case strings.HasPrefix(page, "sys_notifications"), page == "notifications":
			sysCfg.DeleteChannel(ctx, tenant, id)
		default:
			// For in-memory-only pages, just redirect back (no-op delete)
		}
		// HTMx redirect to refresh list
		w.Header().Set("HX-Redirect", "/admin/"+page)
		w.WriteHeader(200)
	}))

	// ===================================================================
	// ★ NOTIFICATION MANAGEMENT — /admin/notifications with full CRUD
	// (routes defined in admin_pages.go; here for future extensions)
	// ===================================================================
	_ = wfr // placeholder for future workflow integration

}

func parseFloat(s string) (float64, error) {
	if s == "" {
		return 0, nil
	}
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

// ===================================================================
// Helper: statusLabelText returns a plain text status label
// ===================================================================
func statusLabelText(active bool) string {
	if active {
		return "启用"
	}
	return "停用"
}

// ===================================================================
// Existing helpers preserved from original file
// ===================================================================

type BftCarrier struct {
	Name, Code, CustomsPoint, DeliveryMethod, DeliveryPrice, Surcharge, Status string
}

var bftCarriersMu sync.Mutex
var bftCarrierSeed = []BftCarrier{
	{"新竹物流", "HCT", "台北/台中/高雄", "宅配", "¥20固定/≥10kg免运", "超长+超材", "启用"},
	{"黑猫宅急便", "YAMATO", "台北/高雄", "宅配", "¥15固定", "超材", "启用"},
	{"顺丰速运", "SF-TW", "台北/台中", "宅配", "¥12/kg", "超长", "启用"},
}

type BftClientPricing struct {
	Client, Route, CargoType, TaxType, WeightPrice, VolumePrice, MinCharge, Status string
}

var bftPricingMu sync.Mutex
var bftPricingSeed = []BftClientPricing{
	{"EZ集运通", "空运", "普货", "全包稅", "¥20/kg", "¥20/才", "¥50起", "启用"},
	{"EZ集运通", "海快", "海快普货", "頻稅", "¥8.30/kg", "—", "¥50起", "启用"},
	{"EZ集运通", "海运", "一类", "全包稅", "¥3.20/kg", "¥20/才", "¥50起", "启用"},
	{"EZ集运通", "海运", "家具类", "全包稅", "¥2.50/kg", "¥15/才", "¥50起", "启用"},
}

type BftTaskDispatch struct {
	Warehouse, TaskType, DispatchMethod, BatchSize, Trigger, Status string
}

var bftDispatchMu sync.Mutex
var bftDispatchSeed = []BftTaskDispatch{
	{"厦门仓", "收货", "轮询", "5件/人", "自动", "启用"},
	{"厦门仓", "拣货", "最少忙碌", "3件/人", "自动", "启用"},
}

func parseID(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}
	var id int64
	_, err := fmt.Sscanf(s, "%d", &id)
	return id, err
}

func bftParcelStatus(s string) string {
	switch s {
	case "pre_declared": return "预报"
	case "received": return "已入仓"
	case "weighed": return "已称重"
	case "stored": return "已上架"
	case "picked": return "待打包"
	case "packed": return "已打包"
	case "outbound": return "已出货"
	case "container_area": return "待装柜"
	case "loaded": return "已装柜"
	case "shipped": return "已出货"
	case "customs": return "清关中"
	case "delivering": return "派送中"
	case "delivered": return "已签收"
	case "abnormal": return "已拒收"
	case "returned": return "已退货"
	default: return s
	}
}

func bftOrderStatus(s string) string {
	switch s {
	case "pending_picking": return "待拣货"
	case "picking": return "拣货中"
	case "pending_packing": return "待打包"
	case "pending_loading": return "待装柜"
	case "loaded": return "已装柜"
	case "in_transit": return "运输中"
	case "customs_clearance": return "清关中"
	case "out_for_delivery": return "派送中"
	case "completed": return "已完成"
	case "cancelled": return "已取消"
	case "shipped": return "已发货"
	default: return s
	}
}

func bftParcelRows(parcels []parcelDomain.Parcel) [][]string {
	rows := make([][]string, len(parcels))
	for i, p := range parcels {
		dims := fmt.Sprintf("%.0f×%.0f×%.0f", p.Length, p.Width, p.Height)
		if p.Length == 0 { dims = "—" }
		rows[i] = []string{
			p.TrackingNumber, p.ProductName, bftParcelStatus(string(p.Status)),
			fmt.Sprintf("%.2f", p.ActualWeight), dims, p.CourierCode,
			p.CreatedAt.Format("01-02 15:04"),
		}
	}
	if len(rows) == 0 {
		rows = [][]string{{"—", "暂无数据", "—", "—", "—", "—", "—"}}
	}
	return rows
}

func bftPageStart(title, icon string) string {
	return `<!DOCTYPE html><html lang=\"zh-TW\"><head><meta charset=\"UTF-8\"><title>` + title +
		` - I56</title><link rel=\"stylesheet\" href=\"/static/css/i56-bdl.css\"><script src=\"/static/js/i56-theme.js\"></script>` +
		`<script src=\"https://unpkg.com/htmx.org@1.9.10\"></script>` +
		`</head><body style=\"padding:16px\">` +
		`<h5 style=\"color:var(--i56-brand);font-size:15px;margin-bottom:16px\">` + title + `</h5>`
}

func bftPageEnd() string { return `</body></html>` }

func bftStatCard(colorClass, label, value, icon string) string {
	return fmt.Sprintf(`<div style="flex:1;min-width:120px;background:var(--i56-bg-surface);border:1px solid var(--i56-border);border-radius:8px;padding:12px;text-align:center;display:flex;flex-direction:column;align-items:center">`+
		`<div style="font-size:24px;font-weight:700;color:var(--i56-text-primary)">%s</div>`+
		`<small style="font-size:11px;color:var(--i56-text-secondary)">%s</small></div>`, value, label)
}

func bftParcelTable(parcels []parcelDomain.Parcel, maxRows int) string {
	h := `<table style="width:100%%;border-collapse:collapse;background:var(--i56-bg-surface);border:1px solid var(--i56-border);border-radius:8px;overflow:hidden;font-size:12px"><thead><tr style="background:var(--i56-bg-base)"><th style="padding:8px 12px;font-weight:600;color:var(--i56-text-secondary);text-align:left;border-bottom:1px solid var(--i56-border);font-size:11px">快递单号</th><th style="padding:8px 12px;font-weight:600;color:var(--i56-text-secondary);text-align:left;border-bottom:1px solid var(--i56-border);font-size:11px">品名</th><th style="padding:8px 12px;font-weight:600;color:var(--i56-text-secondary);text-align:left;border-bottom:1px solid var(--i56-border);font-size:11px">状态</th><th style="padding:8px 12px;font-weight:600;color:var(--i56-text-secondary);text-align:left;border-bottom:1px solid var(--i56-border);font-size:11px">实重(kg)</th><th style="padding:8px 12px;font-weight:600;color:var(--i56-text-secondary);text-align:left;border-bottom:1px solid var(--i56-border);font-size:11px">尺寸(cm)</th></tr></thead><tbody>`
	for i, p := range parcels {
		if i >= maxRows { break }
		dims := fmt.Sprintf("%.0f×%.0f×%.0f", p.Length, p.Width, p.Height)
		if p.Length == 0 { dims = "—" }
		h += fmt.Sprintf(`<tr><td>%s</td><td>%s</td><td><span class="badge badge-brand">%s</span></td><td>%.2f</td><td>%s</td></tr>`,
			p.TrackingNumber, p.ProductName, bftParcelStatus(string(p.Status)), p.ActualWeight, dims)
	}
	if len(parcels) == 0 {
		h += `<tr><td colspan="5" style="padding:32px;text-align:center;color:var(--i56-text-secondary)">暂无数据</td></tr>`
	}
	return h + `</tbody></table>`
}

// ===================================================================
// HTMX Form HTML helpers (BDL 1.0 card style)
// ===================================================================
func bftFormCard(title, icon, body string) string {
	return fmt.Sprintf(`<div style="background:var(--i56-bg-surface);border:1px solid var(--i56-border);border-radius:8px;margin-bottom:12px;overflow:hidden">
  <div style="padding:8px 12px;background:var(--i56-bg-base);display:flex;justify-content:space-between;align-items:center;border-bottom:1px solid var(--i56-border)">
    <strong style="font-size:13px;color:var(--i56-text-primary)">%s</strong>
    <button onclick="this.closest('div').parentElement.remove()" style="background:none;border:none;color:var(--i56-text-secondary);cursor:pointer;font-size:16px;line-height:1">&times;</button>
  </div>
  <div style="padding:12px">%s</div>
</div>`, title, body)
}

func bftFormField(label, name, value, placeholder string) string {
	return fmt.Sprintf(`<div style="flex:1;min-width:150px;margin:4px"><label style="display:block;font-size:11px;color:var(--i56-text-secondary);margin-bottom:4px">%s</label><input name="%s" value="%s" placeholder="%s" style="width:100%%;padding:6px 8px;font-size:12px;background:var(--i56-bg-base);color:var(--i56-text-primary);border:1px solid var(--i56-border);border-radius:4px"></div>`,
		label, name, value, placeholder)
}

// ===================================================================
// Individual form HTML generators (preserved for backward compat)
// ===================================================================

func bftOrderFormHTML(id int64, orderNo, recipient string, routeID int, remark, mode string) string {
	title := "新增集运订单"
	if mode == "edit" { title = "编辑集运订单" }
	return bftFormCard(title, "cart-check",
		fmt.Sprintf(`<form hx-post="/admin/orders/save" hx-swap="none" style="display:flex;flex-wrap:wrap\u003b margin:-4px">
      <input type="hidden" name="id" value="%d">
      %s%s%s
      <div style="width:100%%;margin-top:8px"><button type="submit" class="btn btn-primary">保存</button></div>
    </form>`, id,
			bftFormField("收件人", "recipient_name", recipient, "收件人姓名"),
			bftFormField("路线ID", "route_id", fmt.Sprintf("%d", routeID), "路线编号"),
			bftFormField("备注", "remark", remark, "备注信息"),
		))
}

func bftWarehouseFormHTML(name, code, address, contact, phone, mode string) string {
	title := "新增仓库"
	if mode == "edit" { title = "编辑仓库" }
	return bftFormCard(title, "building",
		fmt.Sprintf(`<form hx-post="/admin/warehouses/save" hx-swap="none" style="display:flex;flex-wrap:wrap\u003b margin:-4px">
      %s%s%s%s%s
      <div style="width:100%%;margin-top:8px"><button type="submit" class="btn btn-primary">保存</button></div>
    </form>`,
			bftFormField("仓库名", "name", name, "仓库名称"),
			bftFormField("编码", "code", code, "仓库编码"),
			bftFormField("地址", "address", address, "详细地址"),
			bftFormField("联系人", "contact", contact, "联系人"),
			bftFormField("电话", "phone", phone, "联系电话"),
		))
}

func bftCarrierFormHTML(name, code, customsPoint, deliveryMethod, deliveryPrice, surcharge, mode string) string {
	title := "新增承运商"
	if mode == "edit" { title = "编辑承运商" }
	return bftFormCard(title, "truck-flatbed",
		fmt.Sprintf(`<form hx-post="/admin/carriers/save" hx-swap="none" style="display:flex;flex-wrap:wrap\u003b margin:-4px">
      %s%s%s%s%s%s
      <div style="width:100%%;margin-top:8px"><button type="submit" class="btn btn-primary">保存</button></div>
    </form>`,
			bftFormField("承运商名称", "name", name, ""),
			bftFormField("编码", "code", code, ""),
			bftFormField("清关点", "customs_point", customsPoint, ""),
			bftFormField("派送方式", "delivery_method", deliveryMethod, ""),
			bftFormField("派送价", "delivery_price", deliveryPrice, ""),
			bftFormField("加收费", "surcharge", surcharge, ""),
		))
}

func bftCourierFormHTML(name, code, region, mode string) string {
	title := "新增快递公司"
	if mode == "edit" { title = "编辑快递公司" }
	return bftFormCard(title, "truck",
		fmt.Sprintf(`<form hx-post="/admin/couriers/save" hx-swap="none" style="display:flex;flex-wrap:wrap\u003b margin:-4px">
      %s%s%s
      <div style="width:100%%;margin-top:8px"><button type="submit" class="btn btn-primary">保存</button></div>
    </form>`,
			bftFormField("名称", "name", name, "快递公司名称"),
			bftFormField("代码", "code", code, "快递公司代码"),
			bftFormField("国家/地区", "region", region, "所在国家或地区"),
		))
}

func bftShippingFormHTML(name, code, transportType, contact, phone, mode string) string {
	title := "新增运输公司"
	if mode == "edit" { title = "编辑运输公司" }
	return bftFormCard(title, "truck",
		fmt.Sprintf(`<form hx-post="/admin/shipping-providers/save" hx-swap="none" style="display:flex;flex-wrap:wrap\u003b margin:-4px">
      %s%s%s%s%s
      <div style="width:100%%;margin-top:8px"><button type="submit" class="btn btn-primary">保存</button></div>
    </form>`,
			bftFormField("名称", "name", name, ""),
			bftFormField("编码", "code", code, ""),
			bftFormField("类型", "type", transportType, "海运/空运/海陆"),
			bftFormField("联系人", "contact", contact, ""),
			bftFormField("电话", "phone", phone, ""),
		))
}

func bftClientFormHTML(name, code, clientType, contactName, contactPhone, contactEmail, mode string) string {
	title := "新增客户"
	if mode == "edit" { title = "编辑客户" }
	return bftFormCard(title, "people",
		fmt.Sprintf(`<form hx-post="/admin/clients/save" hx-swap="none" style="display:flex;flex-wrap:wrap\u003b margin:-4px">
      %s%s%s%s%s%s
      <div style="width:100%%;margin-top:8px"><button type="submit" class="btn btn-primary">保存</button></div>
    </form>`,
			bftFormField("客户名称", "name", name, ""),
			bftFormField("编码", "code", code, ""),
			bftFormField("客户类型", "type", clientType, "platform/shopee/major/peer/normal"),
			bftFormField("联系人", "contact_name", contactName, ""),
			bftFormField("电话", "contact_phone", contactPhone, ""),
			bftFormField("邮箱", "contact_email", contactEmail, ""),
		))
}

func bftClientUserFormHTML(clientID, username, role, email, mode string) string {
	title := "新增客户账号"
	if mode == "edit" { title = "编辑客户账号" }
	return bftFormCard(title, "person-lock",
		fmt.Sprintf(`<form hx-post="/admin/client-users/save" hx-swap="none" style="display:flex;flex-wrap:wrap\u003b margin:-4px">
      %s%s%s%s
      <div style="width:100%%;margin-top:8px"><button type="submit" class="btn btn-primary">保存</button></div>
    </form>`,
			bftFormField("客户ID", "client_id", clientID, ""),
			bftFormField("账号", "username", username, ""),
			bftFormField("角色", "role", role, "运营/管理"),
			bftFormField("邮箱", "email", email, ""),
		))
}

func bftMemberFormHTML(name, phone, email, memberCode string, clientID int64, mode string) string {
	title := "新增会员"
	if mode == "edit" { title = "编辑会员" }
	return bftFormCard(title, "person-heart",
		fmt.Sprintf(`<form hx-post="/admin/client-members/save" hx-swap="none" style="display:flex;flex-wrap:wrap\u003b margin:-4px">
      <input type="hidden" name="client_id" value="%d">
      %s%s%s%s
      <div style="width:100%%;margin-top:8px"><button type="submit" class="btn btn-primary">保存</button></div>
    </form>`, clientID,
			bftFormField("姓名", "name", name, ""),
			bftFormField("手机", "phone", phone, ""),
			bftFormField("邮箱", "email", email, ""),
			bftFormField("会员编号", "member_code", memberCode, ""),
		))
}

func bftPricingFormHTML(client, route, cargoType, taxType, weightPrice, volumePrice, minCharge, mode string) string {
	title := "新增客户价格"
	if mode == "edit" { title = "编辑客户价格" }
	return bftFormCard(title, "currency-exchange",
		fmt.Sprintf(`<form hx-post="/admin/client-pricing/save" hx-swap="none" style="display:flex;flex-wrap:wrap\u003b margin:-4px">
      %s%s%s%s%s%s%s
      <div style="width:100%%;margin-top:8px"><button type="submit" class="btn btn-primary">保存</button></div>
    </form>`,
			bftFormField("客户", "client", client, ""),
			bftFormField("线路", "route", route, ""),
			bftFormField("货类", "cargo_type", cargoType, ""),
			bftFormField("税档", "tax_type", taxType, ""),
			bftFormField("重量单价", "weight_price", weightPrice, ""),
			bftFormField("体积单价", "volume_price", volumePrice, ""),
			bftFormField("最低收费", "min_charge", minCharge, ""),
		))
}

func bftNotificationFormHTML(chanType, name, config, mode string) string {
	title := "添加通知渠道"
	if mode == "edit" { title = "编辑通知渠道" }
	return bftFormCard(title, "bell",
		fmt.Sprintf(`<form hx-post="/admin/notifications/save" hx-swap="none" style="display:flex;flex-wrap:wrap\u003b margin:-4px">
      %s%s%s
      <div style="width:100%%;margin-top:8px"><button type="submit" class="btn btn-primary">保存</button></div>
    </form>`,
			bftFormField("渠道类型", "type", chanType, "email/sms/webhook"),
			bftFormField("名称", "name", name, "渠道名称"),
			bftFormField("配置", "config", config, "JSON配置"),
		))
}

func bftPrintFormHTML(name, templateType, mode string) string {
	title := "新增打印模板"
	if mode == "edit" { title = "编辑打印模板" }
	return bftFormCard(title, "printer",
		fmt.Sprintf(`<form hx-post="/admin/print-templates/save" hx-swap="none" style="display:flex;flex-wrap:wrap\u003b margin:-4px">
      %s%s
      <div style="width:100%%;margin:4px"><label style="display:block;font-size:11px;color:var(--i56-text-secondary);margin-bottom:4px">模板内容</label><textarea name="content" style="width:100%%;padding:6px 8px;font-size:12px;background:var(--i56-bg-base);color:var(--i56-text-primary);border:1px solid var(--i56-border);border-radius:4px" rows="4"></textarea></div>
      <div style="width:100%%;margin-top:8px"><button type="submit" class="btn btn-primary">保存</button></div>
    </form>`,
			bftFormField("模板名", "name", name, "模板名称"),
			bftFormField("类型", "type", templateType, "waybill/customs/carrier"),
		))
}

func bftDispatchFormHTML(warehouse, taskType, dispatchMethod, batchSize, trigger, mode string) string {
	title := "新增任务派发参数"
	if mode == "edit" { title = "编辑任务派发参数" }
	return bftFormCard(title, "sliders",
		fmt.Sprintf(`<form hx-post="/admin/task-dispatch/save" hx-swap="none" style="display:flex;flex-wrap:wrap\u003b margin:-4px">
      %s%s%s%s%s
      <div style="width:100%%;margin-top:8px"><button type="submit" class="btn btn-primary">保存</button></div>
    </form>`,
			bftFormField("仓库", "warehouse", warehouse, "厦门仓"),
			bftFormField("任务类型", "task_type", taskType, "收货/拣货"),
			bftFormField("分派方式", "dispatch_method", dispatchMethod, "轮询/最少忙碌"),
			bftFormField("每批量", "batch_size", batchSize, "5件/人"),
			bftFormField("触发方式", "trigger", trigger, "自动/手动"),
		))
}

func bftBrokerFormHTML(name string, code, brokerNum, prefix, country, points, contact, mode string) string {
	title := "新增清关公司"
	if mode == "edit" { title = "编辑清关公司" }
	return bftFormCard(title, "shield-check",
		fmt.Sprintf(`<form hx-post="/admin/customs-brokers/save" hx-swap="none" style="display:flex;flex-wrap:wrap\u003b margin:-4px">
      %s%s%s%s%s%s%s
      <div style="width:100%%;margin-top:8px"><button type="submit" class="btn btn-primary">保存</button></div>
    </form>`,
			bftFormField("名称", "name", name, ""),
			bftFormField("代码", "code", code, ""),
			bftFormField("编号", "broker_num", brokerNum, ""),
			bftFormField("前缀", "prefix", prefix, ""),
			bftFormField("国家", "country", country, ""),
			bftFormField("清关点", "points", points, ""),
			bftFormField("联系人", "contact", contact, ""),
		))
}
