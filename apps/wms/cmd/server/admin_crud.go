package main

// DEPRECATED: This entire file is deprecated. All admin CRUD routes have been
// migrated to internal module route packages which use common modal helpers
// (common.ModalStart, common.FormField, etc.) instead of local closures.
// The registerAdminCRUD() function is no longer called from main.go.
// Utility types/functions at the bottom (CargoTypeSeed, cargoTypeSeed, etc.)
// are still used by active code in this package.

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/i56/framework/core/router"

	custDomain "github.com/i56/modules/customer/domain"
	custRepo "github.com/i56/modules/customer/repository"
	orderDomain "github.com/i56/modules/order/domain"
	orderRepo "github.com/i56/modules/order/repository"
	orderSvc "github.com/i56/modules/order/service"
	parcelDomain "github.com/i56/modules/parcel/domain"
	parcelRepo "github.com/i56/modules/parcel/repository"
	parcelSvc "github.com/i56/modules/parcel/service"
	psRepo "github.com/i56/modules/parcel_service/repository"
	pricingRepo "github.com/i56/modules/pricing/repository"
	printRepo "github.com/i56/modules/print/repository"
	rbaDomain "github.com/i56/modules/rbac/domain"
	rbacRepo "github.com/i56/modules/rbac/repository"
	sysRepo "github.com/i56/modules/system/repository"
	tmsDomain "github.com/i56/modules/transport/domain"
	tmsRepo "github.com/i56/modules/transport/repository"
	whDomain "github.com/i56/modules/warehouse/domain"
	whRepo "github.com/i56/modules/warehouse/repository"
	whSvc "github.com/i56/modules/warehouse/service"
	woDomain "github.com/i56/modules/workorder/domain"
	woRepo "github.com/i56/modules/workorder/repository"
)

// ===================================================================
// registerAdminCRUD registers complete CRUD (add-form, save, edit-form,
// update, delete) for all BFT56-aligned admin modules.
// Uses modal overlay forms with BDL styling.
// ===================================================================
func registerAdminCRUD(
	r *router.Router,
	a func(http.HandlerFunc) http.HandlerFunc,
	ps *parcelSvc.ParcelService,
	osvc *orderSvc.OrderService,
	ws *whSvc.WarehouseService,
	wr *whRepo.MemWarehouseRepo,
	pr *parcelRepo.MemParcelRepo,
	or *orderRepo.MemOrderRepo,
	rr *tmsRepo.MemRouteRepo,
	cour *tmsRepo.MemCourierRepo,
	lr *custRepo.MemLedgerRepo,
	cr *custRepo.MemClientRepo,
	mr *custRepo.MemMemberRepo,
	dr *custRepo.MemDeclarantRepo,
	ar *custRepo.MemAddressRepo,
	rpr *pricingRepo.MemRoutePriceRepo,
	ppr *printRepo.MemPrintRepo,
	sysCfg *sysRepo.MemSystemConfigRepo,
	sr *psRepo.MemServiceRepo,
	wor *woRepo.MemWorkOrderRepo,
	rbac *rbacRepo.MemRBACRepo,
) {
	const t int64 = 1

	// ===================================================================
	// Modal HTML helpers (BDL 1.0 modal overlay styling)
	// ===================================================================
	modalStart := func(title string) string {
		return `<div class="modal-overlay" onclick="event.target===this&&this.remove()"><div class="modal-content"><div class="modal-header"><span class="modal-title">` + title + `</span><button class="modal-close" onclick="this.closest('.modal-overlay').remove()">&times;</button></div><div class="modal-body">`
	}
	modalEnd := func() string { return `</div></div></div>` }
	formField := func(label, name, value, placeholder string) string {
		return fmt.Sprintf(`<div class="form-group"><label class="form-label">%s</label><input name="%s" value="%s" class="form-input" placeholder="%s"></div>`, label, name, value, placeholder)
	}
	formSelect := func(label, name, value string, opts [][2]string) string {
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
		return `<div class="modal-footer"><button type="button" class="i56-btn" onclick="this.closest('.modal-overlay').remove()">取消</button><button type="submit" class="i56-btn i56-btn-primary">保存</button></div></form>`
	}
	htmlOK := func(w http.ResponseWriter) { w.Header().Set("Content-Type", "text/html; charset=utf-8") }
	redirect := func(w http.ResponseWriter, url string) { w.Header().Set("HX-Redirect", url); w.WriteHeader(200) }
	fieldArea := func(label, name, value, placeholder string, rows int) string {
		return fmt.Sprintf(`<div class="form-group"><label class="form-label">%s</label><textarea name="%s" class="form-input" rows="%d" placeholder="%s">%s</textarea></div>`, label, name, rows, placeholder, value)
	}

	// ===================================================================
	// 1. WAREHOUSE CRUD
	// ===================================================================
	r.GET("/admin/warehouses/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增仓库")+formSave("/admin/warehouses/save")+
			formField("仓库名", "name", "", "仓库名称")+
			formField("编码", "code", "", "仓库编码")+
			formField("地址", "address", "", "详细地址")+
			formField("联系人", "contact", "", "联系人")+
			formField("电话", "phone", "", "联系电话")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/warehouses/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		ws.Create(req.Context(), t, req.FormValue("name"), req.FormValue("code"), req.FormValue("address"), req.FormValue("contact"), req.FormValue("phone"))
		redirect(w, "/admin/warehouses")
	}))
	r.GET("/admin/warehouses/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		id, _ := parseID(req.URL.Query().Get("id"))
		wh, _ := wr.GetByID(req.Context(), t, id)
		if wh == nil {
			http.Error(w, "not found", 404)
			return
		}
		htmlOK(w)
		fmt.Fprint(w, modalStart("编辑仓库")+formSave("/admin/warehouses/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, wh.ID)+
			formField("仓库名", "name", wh.Name, "")+
			formField("编码", "code", wh.Code, "")+
			formField("地址", "address", wh.Address, "")+
			formField("联系人", "contact", wh.Contact, "")+
			formField("电话", "phone", wh.Phone, "")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/warehouses/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := strconv.ParseInt(req.FormValue("id"), 10, 64)
		wr.Update(req.Context(), t, id, &whDomain.Warehouse{ID: id, TenantID: t, Name: req.FormValue("name"), Code: req.FormValue("code"), Address: req.FormValue("address"), Contact: req.FormValue("contact"), Phone: req.FormValue("phone"), IsActive: true})
		redirect(w, "/admin/warehouses")
	}))

	// ===================================================================
	// 2. ORDERS CRUD
	// ===================================================================
	r.GET("/admin/orders/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增集运订单")+formSave("/admin/orders/save")+
			formField("收件人", "recipient_name", "", "收件人姓名")+
			formField("路线ID", "route_id", "", "路线编号")+
			formField("备注", "remark", "", "备注信息")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/orders/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		routeID, _ := parseID(req.FormValue("route_id"))
		osvc.Create(req.Context(), &orderDomain.Order{TenantID: t, WarehouseID: 1, ClientID: 1, RouteID: routeID, RecipientName: req.FormValue("recipient_name"), Status: orderDomain.StatusPendingPicking})
		redirect(w, "/admin/orders")
	}))
	r.GET("/admin/orders/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		id, _ := parseID(req.URL.Query().Get("id"))
		o, _ := or.GetByID(req.Context(), t, id)
		if o == nil {
			http.Error(w, "not found", 404)
			return
		}
		htmlOK(w)
		fmt.Fprint(w, modalStart("编辑集运订单")+formSave("/admin/orders/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, o.ID)+
			formField("收件人", "recipient_name", o.RecipientName, "")+
			formField("路线ID", "route_id", fmt.Sprintf("%d", o.RouteID), "")+
			formField("备注", "remark", o.Remark, "")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/orders/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := strconv.ParseInt(req.FormValue("id"), 10, 64)
		routeID, _ := parseID(req.FormValue("route_id"))
		o, _ := or.GetByID(req.Context(), t, id)
		if o != nil {
			o.RecipientName = req.FormValue("recipient_name")
			o.RouteID = routeID
			o.Remark = req.FormValue("remark")
			or.Update(req.Context(), o)
		}
		redirect(w, "/admin/orders")
	}))
	r.POST("/admin/orders/delete", a(func(w http.ResponseWriter, req *http.Request) {
		id, _ := parseID(req.URL.Query().Get("id"))
		osvc.Cancel(req.Context(), t, id)
		redirect(w, "/admin/orders")
	}))

	// ===================================================================
	// 3. PARCELS CRUD
	// ===================================================================
	r.GET("/admin/parcels/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增包裹")+formSave("/admin/parcels/save")+
			formField("快递单号", "tracking_number", "", "快递单号")+
			formField("品名", "product_name", "", "商品名称")+
			formField("重量(kg)", "actual_weight", "", "实际重量")+
			formField("长(cm)", "length", "", "长度")+
			formField("宽(cm)", "width", "", "宽度")+
			formField("高(cm)", "height", "", "高度")+
			formField("库位", "location_code", "", "库位编码")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/parcels/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		wgt, _ := strconv.ParseFloat(req.FormValue("actual_weight"), 64)
		l, _ := strconv.ParseFloat(req.FormValue("length"), 64)
		wi, _ := strconv.ParseFloat(req.FormValue("width"), 64)
		h, _ := strconv.ParseFloat(req.FormValue("height"), 64)
		pr.Create(req.Context(), &parcelDomain.Parcel{TenantID: t, WarehouseID: 1, ClientID: 1, TrackingNumber: req.FormValue("tracking_number"), ProductName: req.FormValue("product_name"), ParcelName: req.FormValue("product_name"), ActualWeight: wgt, Length: l, Width: wi, Height: h, LocationCode: req.FormValue("location_code"), Status: parcelDomain.StatusPreDeclared, CourierCode: "SF", CargoType: "general"})
		redirect(w, "/admin/parcels")
	}))
	r.GET("/admin/parcels/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		id, _ := parseID(req.URL.Query().Get("id"))
		p, _ := pr.GetByID(req.Context(), t, id)
		if p == nil {
			http.Error(w, "not found", 404)
			return
		}
		htmlOK(w)
		fmt.Fprint(w, modalStart("编辑包裹")+formSave("/admin/parcels/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, p.ID)+
			formField("快递单号", "tracking_number", p.TrackingNumber, "")+
			formField("品名", "product_name", p.ProductName, "")+
			formField("重量(kg)", "actual_weight", fmt.Sprintf("%.2f", p.ActualWeight), "")+
			formField("长(cm)", "length", fmt.Sprintf("%.0f", p.Length), "")+
			formField("宽(cm)", "width", fmt.Sprintf("%.0f", p.Width), "")+
			formField("高(cm)", "height", fmt.Sprintf("%.0f", p.Height), "")+
			formField("库位", "location_code", p.LocationCode, "")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/parcels/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := strconv.ParseInt(req.FormValue("id"), 10, 64)
		p, _ := pr.GetByID(req.Context(), t, id)
		if p != nil {
			wgt, _ := strconv.ParseFloat(req.FormValue("actual_weight"), 64)
			l, _ := strconv.ParseFloat(req.FormValue("length"), 64)
			wi, _ := strconv.ParseFloat(req.FormValue("width"), 64)
			h, _ := strconv.ParseFloat(req.FormValue("height"), 64)
			p.TrackingNumber = req.FormValue("tracking_number")
			p.ProductName = req.FormValue("product_name")
			p.ActualWeight = wgt
			p.Length = l
			p.Width = wi
			p.Height = h
			p.LocationCode = req.FormValue("location_code")
			pr.Update(req.Context(), p)
		}
		redirect(w, "/admin/parcels")
	}))

	// ===================================================================
	// 4. COURIERS CRUD
	// ===================================================================

	// ===================================================================
	// 5. ROUTE TEMPLATES CRUD
	// ===================================================================
	r.GET("/admin/route-templates/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增线路模板")+formSave("/admin/route-templates/save")+
			formField("线路名", "name", "", "线路名称")+
			formSelect("运输方式", "transport_type", "air", [][2]string{{"air", "空运"}, {"sea_express", "海快"}, {"sea", "海运"}, {"land", "陆运"}})+
			formField("体积系数", "volume_coeff", "6000", "计费体积系数")+
			formField("最低重量(kg)", "min_weight", "", "最低计费重量")+
			formField("重量单价", "base_weight_price", "", "元/kg")+
			formField("体积单价", "base_volume_price", "", "元/才")+
			formField("最低收费", "min_amount", "", "最低收费金额")+
			formField("最少天数", "min_days", "", "运输最少天数")+
			formField("最多天数", "max_days", "", "运输最多天数")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/route-templates/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		vc, _ := strconv.Atoi(req.FormValue("volume_coeff"))
		mw, _ := strconv.ParseFloat(req.FormValue("min_weight"), 64)
		bwp, _ := strconv.ParseFloat(req.FormValue("base_weight_price"), 64)
		bvp, _ := strconv.ParseFloat(req.FormValue("base_volume_price"), 64)
		ma, _ := strconv.ParseFloat(req.FormValue("min_amount"), 64)
		mind, _ := strconv.Atoi(req.FormValue("min_days"))
		maxd, _ := strconv.Atoi(req.FormValue("max_days"))
		rr.Create(req.Context(), &tmsDomain.Route{TenantID: t, WarehouseID: 1, Name: req.FormValue("name"), TransportType: req.FormValue("transport_type"), VolumeCoeff: vc, MinWeight: mw, BaseWeightPrice: bwp, BaseVolumePrice: bvp, MinAmount: ma, MinDays: mind, MaxDays: maxd, IsActive: true})
		redirect(w, "/admin/route-templates")
	}))
	r.GET("/admin/route-templates/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		id, _ := parseID(req.URL.Query().Get("id"))
		rt, _ := rr.GetByID(req.Context(), t, id)
		if rt == nil {
			http.Error(w, "not found", 404)
			return
		}
		htmlOK(w)
		fmt.Fprint(w, modalStart("编辑线路模板")+formSave("/admin/route-templates/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, rt.ID)+
			formField("线路名", "name", rt.Name, "")+
			formSelect("运输方式", "transport_type", rt.TransportType, [][2]string{{"air", "空运"}, {"sea_express", "海快"}, {"sea", "海运"}, {"land", "陆运"}})+
			formField("体积系数", "volume_coeff", fmt.Sprintf("%d", rt.VolumeCoeff), "")+
			formField("最低重量(kg)", "min_weight", fmt.Sprintf("%.2f", rt.MinWeight), "")+
			formField("重量单价", "base_weight_price", fmt.Sprintf("%.2f", rt.BaseWeightPrice), "")+
			formField("体积单价", "base_volume_price", fmt.Sprintf("%.2f", rt.BaseVolumePrice), "")+
			formField("最低收费", "min_amount", fmt.Sprintf("%.2f", rt.MinAmount), "")+
			formField("最少天数", "min_days", fmt.Sprintf("%d", rt.MinDays), "")+
			formField("最多天数", "max_days", fmt.Sprintf("%d", rt.MaxDays), "")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/route-templates/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := strconv.ParseInt(req.FormValue("id"), 10, 64)
		vc, _ := strconv.Atoi(req.FormValue("volume_coeff"))
		mw, _ := strconv.ParseFloat(req.FormValue("min_weight"), 64)
		bwp, _ := strconv.ParseFloat(req.FormValue("base_weight_price"), 64)
		bvp, _ := strconv.ParseFloat(req.FormValue("base_volume_price"), 64)
		ma, _ := strconv.ParseFloat(req.FormValue("min_amount"), 64)
		mind, _ := strconv.Atoi(req.FormValue("min_days"))
		maxd, _ := strconv.Atoi(req.FormValue("max_days"))
		rr.Update(req.Context(), &tmsDomain.Route{ID: id, TenantID: t, WarehouseID: 1, Name: req.FormValue("name"), TransportType: req.FormValue("transport_type"), VolumeCoeff: vc, MinWeight: mw, BaseWeightPrice: bwp, BaseVolumePrice: bvp, MinAmount: ma, MinDays: mind, MaxDays: maxd, IsActive: true})
		redirect(w, "/admin/route-templates")
	}))

	// ===================================================================
	// 6. CLIENTS CRUD
	// ===================================================================
	r.GET("/admin/clients/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增客户")+formSave("/admin/clients/save")+
			formField("客户名称", "name", "", "")+
			formField("编码", "code", "", "")+
			formSelect("客户类型", "type", "normal", [][2]string{{"platform", "平台客户"}, {"shopee", "虾皮商家"}, {"major", "大客户"}, {"peer", "同行"}, {"normal", "普通客户"}})+
			formField("联系人", "contact_name", "", "")+
			formField("电话", "contact_phone", "", "")+
			formField("邮箱", "contact_email", "", "")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/clients/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		cr.Create(req.Context(), t, &custDomain.Client{TenantID: t, Name: req.FormValue("name"), Code: req.FormValue("code"), ClientType: custDomain.ClientType(req.FormValue("type")), ContactName: req.FormValue("contact_name"), ContactPhone: req.FormValue("contact_phone"), ContactEmail: req.FormValue("contact_email"), IsActive: true})
		redirect(w, "/admin/clients")
	}))
	r.GET("/admin/clients/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		id, _ := parseID(req.URL.Query().Get("id"))
		c, _ := cr.GetByID(req.Context(), t, id)
		if c == nil {
			http.Error(w, "not found", 404)
			return
		}
		htmlOK(w)
		fmt.Fprint(w, modalStart("编辑客户")+formSave("/admin/clients/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, c.ID)+
			formField("客户名称", "name", c.Name, "")+
			formField("编码", "code", c.Code, "")+
			formSelect("客户类型", "type", string(c.ClientType), [][2]string{{"platform", "平台客户"}, {"shopee", "虾皮商家"}, {"major", "大客户"}, {"peer", "同行"}, {"normal", "普通客户"}})+
			formField("联系人", "contact_name", c.ContactName, "")+
			formField("电话", "contact_phone", c.ContactPhone, "")+
			formField("邮箱", "contact_email", c.ContactEmail, "")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/clients/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := strconv.ParseInt(req.FormValue("id"), 10, 64)
		cr.Update(req.Context(), t, id, &custDomain.Client{ID: id, TenantID: t, Name: req.FormValue("name"), Code: req.FormValue("code"), ClientType: custDomain.ClientType(req.FormValue("type")), ContactName: req.FormValue("contact_name"), ContactPhone: req.FormValue("contact_phone"), ContactEmail: req.FormValue("contact_email"), IsActive: true})
		redirect(w, "/admin/clients")
	}))
	r.POST("/admin/clients/delete", a(func(w http.ResponseWriter, req *http.Request) {
		id, _ := parseID(req.URL.Query().Get("id"))
		cr.Delete(req.Context(), t, id)
		redirect(w, "/admin/clients")
	}))

	// ===================================================================
	// 7. CLIENT ADDRESSES CRUD
	// ===================================================================
	r.GET("/admin/client-addresses/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增收件地址")+formSave("/admin/client-addresses/save")+
			formField("会员ID", "member_id", "", "")+
			formField("收件人", "recipient_name", "", "")+
			formField("电话", "phone", "", "")+
			formField("邮编", "postal_code", "", "")+
			formField("城市", "city", "", "")+
			formField("区域", "district", "", "")+
			formField("详细地址", "address", "", "")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/client-addresses/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		mid, _ := parseID(req.FormValue("member_id"))
		ar.Create(req.Context(), mid, &custDomain.MemberAddress{RecipientName: req.FormValue("recipient_name"), Phone: req.FormValue("phone"), PostalCode: req.FormValue("postal_code"), City: req.FormValue("city"), District: req.FormValue("district"), Address: req.FormValue("address")})
		redirect(w, "/admin/client-addresses")
	}))

	// ===================================================================
	// 8. DECLARANTS CRUD
	// ===================================================================
	r.GET("/admin/declarants/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增申报人")+formSave("/admin/declarants/save")+
			formField("客户ID", "client_id", "1", "")+
			formField("姓名", "name", "", "")+
			formField("证件号", "id_number", "", "")+
			formSelect("类型", "type", "individual", [][2]string{{"individual", "个人"}, {"company", "公司"}})+
			formField("电话", "phone", "", "")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/declarants/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		cid, _ := parseID(req.FormValue("client_id"))
		dr.Create(req.Context(), cid, &custDomain.Declarant{ClientID: cid, Type: custDomain.DeclarantType(req.FormValue("type")), Name: req.FormValue("name"), IDNumber: req.FormValue("id_number"), Phone: req.FormValue("phone"), IsActive: true})
		redirect(w, "/admin/declarants")
	}))
	r.GET("/admin/declarants/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		id, _ := parseID(req.URL.Query().Get("id"))
		d, _ := dr.GetByID(req.Context(), 0, id)
		if d == nil {
			http.Error(w, "not found", 404)
			return
		}
		htmlOK(w)
		fmt.Fprint(w, modalStart("编辑申报人")+formSave("/admin/declarants/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%d"><input type="hidden" name="client_id" value="%d">`, d.ID, d.ClientID)+
			formField("姓名", "name", d.Name, "")+
			formField("证件号", "id_number", d.IDNumber, "")+
			formSelect("类型", "type", string(d.Type), [][2]string{{"individual", "个人"}, {"company", "公司"}})+
			formField("电话", "phone", d.Phone, "")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/declarants/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := strconv.ParseInt(req.FormValue("id"), 10, 64)
		cid, _ := parseID(req.FormValue("client_id"))
		dr.Update(req.Context(), cid, id, &custDomain.Declarant{ID: id, ClientID: cid, Type: custDomain.DeclarantType(req.FormValue("type")), Name: req.FormValue("name"), IDNumber: req.FormValue("id_number"), Phone: req.FormValue("phone"), IsActive: true})
		redirect(w, "/admin/declarants")
	}))

	// ===================================================================
	// 9. CLIENT MEMBERS CRUD
	// ===================================================================

	// ===================================================================
	// 10. CLIENT PRICING CRUD (uses MemRoutePriceRepo)
	// ===================================================================

	// ===================================================================
	// 11. CLIENT LEDGERS CRUD
	// ===================================================================
	r.GET("/admin/client-ledgers/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增余额记录")+formSave("/admin/client-ledgers/save")+
			formField("客户ID", "client_id", "1", "")+
			formField("金额", "amount", "", "正数=充值, 负数=扣款")+
			formSelect("类型", "type", "manual", [][2]string{{"recharge", "充值"}, {"charge", "扣款"}, {"refund", "退款"}, {"manual", "手动调整"}})+
			formField("描述", "description", "", "")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/client-ledgers/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		cid, _ := parseID(req.FormValue("client_id"))
		amt, _ := strconv.ParseFloat(req.FormValue("amount"), 64)
		c, _ := cr.GetByID(req.Context(), t, cid)
		balanceAfter := amt
		if c != nil {
			balanceAfter = c.Balance + amt
			c.Balance = balanceAfter
			cr.Update(req.Context(), t, cid, c)
		}
		lr.Add(req.Context(), &custRepo.LedgerEntry{TenantID: t, ClientID: cid, Amount: amt, BalanceAfter: balanceAfter, Type: req.FormValue("type"), Description: req.FormValue("description")})
		redirect(w, "/admin/client-ledgers")
	}))

	// ===================================================================
	// 12. CLIENT RECHARGES (manual recharge)
	// ===================================================================
	r.GET("/admin/client-recharges/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("客户充值")+formSave("/admin/client-recharges/save")+
			formField("客户ID", "client_id", "1", "")+
			formField("充值金额", "amount", "", "")+
			formSelect("方式", "method", "bank_transfer", [][2]string{{"bank_transfer", "银行转账"}, {"wechat", "微信支付"}, {"alipay", "支付宝"}, {"cash", "现金"}})+
			formField("备注", "description", "", "")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/client-recharges/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		cid, _ := parseID(req.FormValue("client_id"))
		amt, _ := strconv.ParseFloat(req.FormValue("amount"), 64)
		c, _ := cr.GetByID(req.Context(), t, cid)
		balanceAfter := amt
		if c != nil {
			balanceAfter = c.Balance + amt
			c.Balance = balanceAfter
			cr.Update(req.Context(), t, cid, c)
		}
		lr.Add(req.Context(), &custRepo.LedgerEntry{TenantID: t, ClientID: cid, Amount: amt, BalanceAfter: balanceAfter, Type: req.FormValue("method"), Description: req.FormValue("description")})
		redirect(w, "/admin/client-recharges")
	}))

	// ===================================================================
	// 13. PRINT TEMPLATES CRUD
	// ===================================================================
	r.GET("/admin/print-templates/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增打印模板")+formSave("/admin/print-templates/save")+
			formField("模板名", "name", "", "模板名称")+
			formSelect("类型", "type", "waybill", [][2]string{{"waybill", "面单"}, {"customs", "清关单"}, {"carrier", "承运商面单"}})+
			fieldArea("模板内容", "content", "", "HTML模板", 5)+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/print-templates/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		ppr.Add(req.Context(), t, req.FormValue("name"), req.FormValue("type"), req.FormValue("content"))
		redirect(w, "/admin/print-templates")
	}))

	// ===================================================================
	// 14. NOTIFICATIONS CRUD
	// ===================================================================
	r.GET("/admin/notifications/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增通知渠道")+formSave("/admin/notifications/save")+
			formSelect("渠道类型", "type", "email", [][2]string{{"email", "邮件"}, {"sms", "短信"}, {"webhook", "Webhook"}})+
			formField("名称", "name", "", "渠道名称")+
			fieldArea("配置(JSON)", "config", "", "{\"smtp_host\":\"...\"}", 3)+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/notifications/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		sysCfg.SaveNotificationChannel(req.FormValue("type"), req.FormValue("name"), req.FormValue("config"))
		redirect(w, "/admin/notifications")
	}))

	// ===================================================================
	// 15. CUSTOMS BROKERS CRUD (in-memory seed)
	// ===================================================================
	r.GET("/admin/customs-brokers/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增清关公司")+formSave("/admin/customs-brokers/save")+
			formField("名称", "name", "", "")+
			formField("代码", "code", "", "")+
			formField("编号", "broker_num", "", "")+
			formField("前缀", "prefix", "", "")+
			formField("国家", "country", "", "")+
			formField("清关点", "points", "", "")+
			formField("联系人", "contact", "", "")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/customs-brokers/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		addBroker(req.FormValue("name"), req.FormValue("code"), req.FormValue("broker_num"), req.FormValue("prefix"), req.FormValue("country"), req.FormValue("points"), req.FormValue("contact"))
		redirect(w, "/admin/customs-brokers")
	}))

	// ===================================================================
	// 16. CARGO TYPES CRUD (in-memory seed)
	// ===================================================================
	r.GET("/admin/cargo-types/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增货物类型")+formSave("/admin/cargo-types/save")+
			formField("名称", "name", "", "如: 普货")+
			formField("编码", "code", "", "如: general")+
			formField("描述", "description", "", "")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/cargo-types/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		addCargoType(req.FormValue("name"), req.FormValue("code"), req.FormValue("description"))
		redirect(w, "/admin/cargo-types")
	}))

	// ===================================================================
	// 17. AREA GROUPS CRUD (in-memory seed)
	// ===================================================================
	r.GET("/admin/area-groups/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增区域组")+formSave("/admin/area-groups/save")+
			formField("区域名", "name", "", "如: 华南区")+
			formField("编码", "code", "", "如: CN-SOUTH")+
			formField("覆盖范围", "coverage", "", "如: 广东/福建/海南")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/area-groups/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		addAreaGroup(req.FormValue("name"), req.FormValue("code"), req.FormValue("coverage"))
		redirect(w, "/admin/area-groups")
	}))

	// ===================================================================
	// 18. TRANSPORT TYPES / TRANSPORT MODES CRUD (in-memory seed)
	// ===================================================================
	r.GET("/admin/transport-modes/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增运输方式")+formSave("/admin/transport-modes/save")+
			formField("方式", "name", "", "如: 空运")+
			formField("编码", "code", "", "如: air")+
			formField("描述", "description", "", "")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/transport-modes/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		addTransportType(req.FormValue("name"), req.FormValue("code"), req.FormValue("description"))
		redirect(w, "/admin/transport-modes")
	}))

	// ===================================================================
	// 19. TRANSPORT COMPANIES / SHIPPING PROVIDERS CRUD (in-memory)
	// ===================================================================

	// ===================================================================
	// 20. CUSTOMS POINTS CRUD (in-memory seed)
	// ===================================================================
	r.GET("/admin/customs-points/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增清关点")+formSave("/admin/customs-points/save")+
			formField("名称", "name", "", "如: 台北港")+
			formField("编码", "code", "", "如: TPE-PORT")+
			formField("国家/地区", "country", "", "如: 台湾")+
			formField("承运商", "carrier", "", "如: 新竹物流")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/customs-points/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		addCustomsPoint(req.FormValue("name"), req.FormValue("code"), req.FormValue("country"), req.FormValue("carrier"))
		redirect(w, "/admin/customs-points")
	}))

	// ===================================================================
	// 21. SERVICE TYPES CRUD (uses MemServiceRepo)
	// ===================================================================
	r.GET("/admin/service-types/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增服务类型")+formSave("/admin/service-types/save")+
			formField("名称", "name", "", "如: 打木箱")+
			formField("编码", "code", "", "如: WOODEN_CRATE")+
			formSelect("分类", "category", "开箱类", [][2]string{{"开箱类", "开箱类"}, {"加固类", "加固类"}, {"打包类", "打包类"}, {"退货类", "退货类"}})+
			formField("单价", "unit_price", "0.00", "")+
			formSelect("计费模式", "price_mode", "fixed", [][2]string{{"fixed", "固定"}, {"per_qty", "按数量"}, {"per_kg", "按重量"}})+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/service-types/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		up, _ := strconv.ParseFloat(req.FormValue("unit_price"), 64)
		sr.AddType(req.FormValue("name"), req.FormValue("code"), req.FormValue("category"), up, req.FormValue("price_mode"))
		redirect(w, "/admin/service-types")
	}))

	// ===================================================================
	// 22. SERVICE TEMPLATES CRUD (in-memory seed)
	// ===================================================================

	// ===================================================================
	// 23. SERVICE WORKORDERS CRUD
	// ===================================================================
	r.GET("/admin/service-workorders/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		id, _ := parseID(req.URL.Query().Get("id"))
		wo, _ := wor.GetByID(req.Context(), t, id)
		if wo == nil {
			http.Error(w, "not found", 404)
			return
		}
		htmlOK(w)
		fmt.Fprint(w, modalStart("编辑服务工单")+formSave("/admin/service-workorders/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, wo.ID)+
			formField("标题", "title", wo.Title, "")+
			formField("描述", "description", wo.Description, "")+
			formSelect("状态", "status", wo.Status, [][2]string{{"pending", "待处理"}, {"in_progress", "进行中"}, {"completed", "已完成"}, {"cancelled", "已取消"}})+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/service-workorders/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := strconv.ParseInt(req.FormValue("id"), 10, 64)
		wor.Update(req.Context(), &woDomain.WorkOrder{ID: id, TenantID: t, Title: req.FormValue("title"), Description: req.FormValue("description"), Status: req.FormValue("status")})
		redirect(w, "/admin/service-workorders")
	}))

	// ===================================================================
	// 24. SERVICE ORDERS CRUD
	// ===================================================================

	// ===================================================================
	// 25. RBAC ROLES CRUD — add edit form
	// ===================================================================
	r.GET("/admin/roles/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		id, _ := parseID(req.URL.Query().Get("id"))
		ro, _ := rbac.GetRoleByID(req.Context(), t, id)
		if ro == nil {
			http.Error(w, "not found", 404)
			return
		}
		perms, _, _ := rbac.ListPermissions(req.Context(), 0, 200)
		htmlOK(w)
		h := modalStart("编辑角色") + formSave("/admin/roles/update") +
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, ro.ID) +
			formField("角色名", "name", ro.Name, "角色名称") +
			formField("Slug", "slug", ro.Slug, "role_slug") +
			formSelect("状态", "is_active", fmt.Sprintf("%v", ro.IsActive), [][2]string{{"true", "启用"}, {"false", "停用"}}) +
			formField("描述", "description", ro.Description, "角色描述") +
			`<div class="form-group"><label class="form-label">权限</label><div style="max-height:200px;overflow-y:auto;border:1px solid var(--i56-border);border-radius:4px;padding:8px;background:var(--i56-bg-base)">`
		pidSet := map[int64]bool{}
		for _, pid := range ro.PermissionIDs {
			pidSet[pid] = true
		}
		for _, p := range perms {
			ck := ""
			if pidSet[p.ID] {
				ck = " checked"
			}
			h += fmt.Sprintf(`<div style="font-size:12px;padding:2px 0"><label><input type="checkbox" name="perm_ids" value="%d"%s> %s</label></div>`, p.ID, ck, p.Name)
		}
		h += `</div></div>` + formFooter() + modalEnd()
		fmt.Fprint(w, h)
	}))
	r.POST("/admin/roles/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := strconv.ParseInt(req.FormValue("id"), 10, 64)
		isActive := req.FormValue("is_active") == "true"
		var permIDs []int64
		for _, s := range req.Form["perm_ids"] {
			if pid, err := strconv.ParseInt(s, 10, 64); err == nil {
				permIDs = append(permIDs, pid)
			}
		}
		rbac.UpdateRole(req.Context(), t, id, &rbaDomain.Role{Name: req.FormValue("name"), Slug: req.FormValue("slug"), Description: req.FormValue("description"), PermissionIDs: permIDs, IsActive: isActive})
		redirect(w, "/admin/roles")
	}))

	// ===================================================================
	// 26. RBAC USERS CRUD — add edit form
	// ===================================================================
	r.GET("/admin/users/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		id, _ := parseID(req.URL.Query().Get("id"))
		u, _ := rbac.GetUserByID(req.Context(), t, id)
		if u == nil {
			http.Error(w, "not found", 404)
			return
		}
		roles, _, _ := rbac.ListRoles(req.Context(), t, 0, 50)
		htmlOK(w)
		h := modalStart("编辑员工") + formSave("/admin/users/update") +
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, u.ID) +
			formField("账号", "username", u.Username, "登录账号") +
			formField("密码", "password", "", "留空不修改") +
			formField("姓名", "real_name", u.RealName, "真实姓名") +
			`<div class="form-group"><label class="form-label">角色</label><select name="role_id" class="form-input">`
		for _, ro := range roles {
			sel := ""
			if ro.ID == u.RoleID {
				sel = " selected"
			}
			h += fmt.Sprintf(`<option value="%d"%s>%s</option>`, ro.ID, sel, ro.Name)
		}
		h += `</select></div>` +
			formField("邮箱", "email", u.Email, "") +
			formField("电话", "phone", u.Phone, "") +
			formFooter() + modalEnd()
		fmt.Fprint(w, h)
	}))
	r.POST("/admin/users/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := strconv.ParseInt(req.FormValue("id"), 10, 64)
		roleID, _ := strconv.ParseInt(req.FormValue("role_id"), 10, 64)
		existing, _ := rbac.GetUserByID(req.Context(), t, id)
		pwd := ""
		if existing != nil {
			pwd = existing.Password
		}
		if np := req.FormValue("password"); np != "" {
			pwd = np
		}
		rbac.UpdateUser(req.Context(), t, id, &rbaDomain.User{Username: req.FormValue("username"), Password: pwd, RealName: req.FormValue("real_name"), Email: req.FormValue("email"), Phone: req.FormValue("phone"), RoleID: roleID, IsActive: true})
		redirect(w, "/admin/users")
	}))

	// ===================================================================
	// 27. RBAC PERMISSIONS CRUD — add edit form
	// ===================================================================
	r.GET("/admin/permissions/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		id, _ := parseID(req.URL.Query().Get("id"))
		p, _ := rbac.GetPermissionByID(req.Context(), id)
		if p == nil {
			http.Error(w, "not found", 404)
			return
		}
		htmlOK(w)
		fmt.Fprint(w, modalStart("编辑权限")+formSave("/admin/permissions/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, p.ID)+
			formField("权限名", "name", p.Name, "权限名称")+
			formField("Slug", "slug", p.Slug, "module:action")+
			formField("模块", "module", p.Module, "模块名")+
			formSelect("状态", "is_active", fmt.Sprintf("%v", p.IsActive), [][2]string{{"true", "启用"}, {"false", "停用"}})+
			formField("描述", "description", p.Description, "描述")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/permissions/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := strconv.ParseInt(req.FormValue("id"), 10, 64)
		rbac.UpdatePermission(req.Context(), id, &rbaDomain.Permission{Name: req.FormValue("name"), Slug: req.FormValue("slug"), Module: req.FormValue("module"), Description: req.FormValue("description"), IsActive: req.FormValue("is_active") == "true"})
		redirect(w, "/admin/permissions")
	}))

	// ===================================================================
	// 28. CLIENT PERMISSIONS CRUD — add edit form
	// ===================================================================
	r.GET("/admin/client-permissions/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		id, _ := parseID(req.URL.Query().Get("id"))
		cp, _ := rbac.GetClientPermissionByID(req.Context(), t, id)
		if cp == nil {
			http.Error(w, "not found", 404)
			return
		}
		clients, _, _ := cr.List(req.Context(), t, 0, 50)
		perms, _, _ := rbac.ListPermissions(req.Context(), 0, 200)
		htmlOK(w)
		h := modalStart("编辑客户端权限") + formSave("/admin/client-permissions/update") +
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, cp.ID) +
			`<div class="form-group"><label class="form-label">客户</label><select name="client_id" class="form-input">`
		for _, c := range clients {
			sel := ""
			if c.ID == cp.ClientID {
				sel = " selected"
			}
			h += fmt.Sprintf(`<option value="%d"%s>%s</option>`, c.ID, sel, c.Name)
		}
		h += `</select></div><div class="form-group"><label class="form-label">权限</label><div style="max-height:200px;overflow-y:auto;border:1px solid var(--i56-border);border-radius:4px;padding:8px;background:var(--i56-bg-base)">`
		slugSet := map[string]bool{}
		for _, s := range cp.PermissionSlugs {
			slugSet[s] = true
		}
		for _, p := range perms {
			ck := ""
			if slugSet[p.Slug] {
				ck = " checked"
			}
			h += fmt.Sprintf(`<div style="font-size:12px;padding:2px 0"><label><input type="checkbox" name="perm_slugs" value="%s"%s> %s</label></div>`, p.Slug, ck, p.Name)
		}
		h += `</div></div>` + formFooter() + modalEnd()
		fmt.Fprint(w, h)
	}))
	r.POST("/admin/client-permissions/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := strconv.ParseInt(req.FormValue("id"), 10, 64)
		clientID, _ := strconv.ParseInt(req.FormValue("client_id"), 10, 64)
		rbac.UpdateClientPermission(req.Context(), t, id, &rbaDomain.ClientPermission{ClientID: clientID, ClientName: req.FormValue("client_name"), PermissionSlugs: req.Form["perm_slugs"], IsActive: true})
		redirect(w, "/admin/client-permissions")
	}))

	// ===================================================================
	// 29. CARRIERS CRUD (in-memory seed BftCarrier)
	// ===================================================================
	r.GET("/admin/carriers/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增承运商")+formSave("/admin/carriers/save")+
			formField("承运商名称", "name", "", "")+
			formField("编码", "code", "", "")+
			formField("清关点", "customs_point", "", "")+
			formSelect("派送方式", "delivery_method", "宅配", [][2]string{{"宅配", "宅配"}, {"专车", "专车"}})+
			formField("派送价", "delivery_price", "", "如: ¥20固定/≥10kg免运")+
			formField("加收费", "surcharge", "", "如: 超长+超材")+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/carriers/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		bftCarriersMu.Lock()
		bftCarrierSeed = append(bftCarrierSeed, BftCarrier{Name: req.FormValue("name"), Code: req.FormValue("code"), CustomsPoint: req.FormValue("customs_point"), DeliveryMethod: req.FormValue("delivery_method"), DeliveryPrice: req.FormValue("delivery_price"), Surcharge: req.FormValue("surcharge"), Status: "启用"})
		bftCarriersMu.Unlock()
		redirect(w, "/admin/carriers")
	}))

	// ===================================================================
	// 30. TASK DISPATCH PARAMS CRUD (in-memory)
	// ===================================================================
	r.GET("/admin/task-dispatch/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		htmlOK(w)
		fmt.Fprint(w, modalStart("新增任务派发参数")+formSave("/admin/task-dispatch/save")+
			formField("仓库", "warehouse", "", "厦门仓")+
			formSelect("任务类型", "task_type", "收货", [][2]string{{"收货", "收货"}, {"拣货", "拣货"}})+
			formSelect("分派方式", "dispatch_method", "轮询", [][2]string{{"轮询", "轮询"}, {"最少忙碌", "最少忙碌"}})+
			formField("每批量", "batch_size", "5件/人", "")+
			formSelect("触发方式", "trigger", "自动", [][2]string{{"自动", "自动"}, {"手动", "手动"}})+
			formFooter()+modalEnd())
	}))
	r.POST("/admin/task-dispatch/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		bftDispatchMu.Lock()
		bftDispatchSeed = append(bftDispatchSeed, BftTaskDispatch{Warehouse: req.FormValue("warehouse"), TaskType: req.FormValue("task_type"), DispatchMethod: req.FormValue("dispatch_method"), BatchSize: req.FormValue("batch_size"), Trigger: req.FormValue("trigger"), Status: "启用"})
		bftDispatchMu.Unlock()
		redirect(w, "/admin/task-dispatch")
	}))

	_ = pr
	_ = or
}

// ===================================================================
// In-memory seed data and helpers for modules without real repos
// ===================================================================

// Cargo types
type CargoTypeSeed struct{ Name, Code, Description, Status string }

var cargoTypesMu sync.Mutex
var cargoTypeSeed = []CargoTypeSeed{
	{"普货", "general", "通用普通货物", "启用"},
	{"一类", "class1", "纺织品/服装", "启用"},
	{"二类", "class2", "日用品/小家电", "启用"},
	{"三类", "class3", "食品/保健品", "启用"},
	{"四类", "class4", "电子产品/电器", "启用"},
	{"特货", "special", "液体/粉末/电池", "启用"},
}

func addCargoType(name, code, desc string) {
	cargoTypesMu.Lock()
	defer cargoTypesMu.Unlock()
	cargoTypeSeed = append(cargoTypeSeed, CargoTypeSeed{name, code, desc, "启用"})
}

// Area groups
type AreaGroupSeed struct{ Name, Code, Coverage, Status string }

var areaGroupsMu sync.Mutex
var areaGroupSeed = []AreaGroupSeed{
	{"华南区", "CN-SOUTH", "广东/福建/海南", "启用"},
	{"华东区", "CN-EAST", "上海/江苏/浙江", "启用"},
	{"台湾北部", "TW-NORTH", "台北/新北/基隆", "启用"},
	{"台湾南部", "TW-SOUTH", "台中/台南/高雄", "启用"},
}

func addAreaGroup(name, code, coverage string) {
	areaGroupsMu.Lock()
	defer areaGroupsMu.Unlock()
	areaGroupSeed = append(areaGroupSeed, AreaGroupSeed{name, code, coverage, "启用"})
}

// Transport types / modes
type TransportTypeSeed struct{ Name, Code, Description, Status string }

var transportTypesMu sync.Mutex
var transportTypeSeed = []TransportTypeSeed{
	{"空运", "air", "航空快递", "启用"},
	{"海快", "sea_express", "海运快速", "启用"},
	{"海运", "sea", "普通海运", "启用"},
	{"陆运", "land", "公路运输", "启用"},
}

func addTransportType(name, code, desc string) {
	transportTypesMu.Lock()
	defer transportTypesMu.Unlock()
	transportTypeSeed = append(transportTypeSeed, TransportTypeSeed{name, code, desc, "启用"})
}

// Customs points
type CustomsPointSeed struct{ Name, Code, Country, Carrier string }

var customsPointsMu sync.Mutex
var customsPointSeed = []CustomsPointSeed{
	{"台北港", "TPE-PORT", "台湾", "新竹物流"},
	{"基隆港", "KEE-PORT", "台湾", "黑猫宅急便"},
	{"台中港", "TXG-PORT", "台湾", "顺丰速运"},
	{"高雄港", "KHH-PORT", "台湾", "新竹物流"},
}

func addCustomsPoint(name, code, country, carrier string) {
	customsPointsMu.Lock()
	defer customsPointsMu.Unlock()
	customsPointSeed = append(customsPointSeed, CustomsPointSeed{name, code, country, carrier})
}

// Customs brokers
type BrokerSeed struct{ Name, Code, BrokerNum, Prefix, Country, Points, Contact, Status string }

var brokersMu sync.Mutex
var brokerSeed = []BrokerSeed{
	{"厦门清关", "XM-CS", "776XM", "", "中国", "厦门/福州", "陈经理", "启用"},
}

func addBroker(name, code, brokerNum, prefix, country, points, contact string) {
	brokersMu.Lock()
	defer brokersMu.Unlock()
	brokerSeed = append(brokerSeed, BrokerSeed{name, code, brokerNum, prefix, country, points, contact, "启用"})
}

// Service templates
type ServiceTemplateSeed struct{ Name, ServiceCode, UnitPrice, PriceMode string }

var serviceTemplatesMu sync.Mutex
var serviceTemplateSeed = []ServiceTemplateSeed{
	{"开箱验货", "OPEN_INSPECT", "免费", "fixed"},
}

func addServiceTemplate(name, serviceCode, unitPrice, priceMode string) {
	serviceTemplatesMu.Lock()
	defer serviceTemplatesMu.Unlock()
	serviceTemplateSeed = append(serviceTemplateSeed, ServiceTemplateSeed{name, serviceCode, unitPrice, priceMode})
}
