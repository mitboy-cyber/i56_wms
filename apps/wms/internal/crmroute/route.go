// Package route provides CRM (客户管理) admin route registration.
package route

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/i56/framework/core/router"

	"github.com/i56/i56-apps/i56-wms/internal/common"

	custDomain "github.com/i56/modules/customer/domain"
	custRepo "github.com/i56/modules/customer/repository"
	pricingDomain "github.com/i56/modules/pricing/domain"
	pricingRepo "github.com/i56/modules/pricing/repository"
)

// Register CRM admin routes (~11 list pages + CRUD).
func Register(
	r *router.Router,
	a func(http.HandlerFunc) http.HandlerFunc,
	rc *common.RenderCtx,
	cr *custRepo.MemClientRepo,
	mr *custRepo.MemMemberRepo,
	lr *custRepo.MemLedgerRepo,
	ar *custRepo.MemAddressRepo,
	dr *custRepo.MemDeclarantRepo,
	rpr *pricingRepo.MemRoutePriceRepo,
) {
	const tenant int64 = 1
	gp := rc.NewGenericList()

	// ─── /admin/clients CRUD (from admin_crud.go) ───
	r.GET("/admin/crm-clients", a(func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "/admin/clients", http.StatusMovedPermanently)
	}))
	// List page — shows real client data from repo
	r.GET("/admin/clients", a(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		clientList, total, _ := cr.List(ctx, tenant, 0, 500)
		rows := make([][]string, 0, len(clientList))
		for _, c := range clientList {
			clientTypeCN := string(c.ClientType)
			switch c.ClientType {
			case custDomain.ClientTypePlatform: clientTypeCN = "平台客户"
			case custDomain.ClientTypeShopee: clientTypeCN = "虾皮商家"
			case custDomain.ClientTypeMajor: clientTypeCN = "大客户"
			case custDomain.ClientTypePeer: clientTypeCN = "同行"
			case custDomain.ClientTypeNormal: clientTypeCN = "普通客户"
			}
			rows = append(rows, []string{
				fmt.Sprintf("%d", c.ID), c.Name, c.Code,
				clientTypeCN, c.ContactName, c.ContactPhone,
				common.StatusLabelText(c.IsActive),
			})
		}
		if len(rows) == 0 {
			rows = [][]string{{"暂无数据", "—", "—", "—", "—", "—", "—", "—"}}
		}
		rc.Exec(rc.Tmpl, "crm_clients", w, "crm_clients.html", map[string]any{
			"Title": "客户管理", "Page": "crm_clients",
			"Columns": []string{"ID", "名称", "编码", "类型", "联系人", "电话", "状态"},
			"Rows": rows, "Total": int(total),
			"AddURL": "/admin/clients/add-form", "HasActions": true,
		})
	}))
	r.GET("/admin/clients/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增客户")+common.FormSave("/admin/clients/save")+
			common.FormField("客户名称", "name", "", "")+
			common.FormField("编码", "code", "", "")+
			common.FormSelect("客户类型", "type", "normal",
				[2]string{"platform", "平台客户"}, [2]string{"shopee", "虾皮商家"},
				[2]string{"major", "大客户"}, [2]string{"peer", "同行"}, [2]string{"normal", "普通客户"})+
			common.FormField("联系人", "contact_name", "", "")+
			common.FormField("电话", "contact_phone", "", "")+
			common.FormField("邮箱", "contact_email", "", "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/clients/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		cr.Create(req.Context(), tenant, &custDomain.Client{TenantID: tenant, Name: req.FormValue("name"), Code: req.FormValue("code"), ClientType: custDomain.ClientType(req.FormValue("type")), ContactName: req.FormValue("contact_name"), ContactPhone: req.FormValue("contact_phone"), ContactEmail: req.FormValue("contact_email"), IsActive: true})
		common.Redirect(w, "/admin/clients")
	}))
	r.GET("/admin/clients/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		id, _ := common.ParseID(req.URL.Query().Get("id"))
		c, _ := cr.GetByID(req.Context(), tenant, id)
		if c == nil { http.Error(w, "not found", 404); return }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑客户")+common.FormSave("/admin/clients/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, c.ID)+
			common.FormField("客户名称", "name", c.Name, "")+
			common.FormField("编码", "code", c.Code, "")+
			common.FormSelect("客户类型", "type", string(c.ClientType),
				[2]string{"platform", "平台客户"}, [2]string{"shopee", "虾皮商家"},
				[2]string{"major", "大客户"}, [2]string{"peer", "同行"}, [2]string{"normal", "普通客户"})+
			common.FormField("联系人", "contact_name", c.ContactName, "")+
			common.FormField("电话", "contact_phone", c.ContactPhone, "")+
			common.FormField("邮箱", "contact_email", c.ContactEmail, "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/clients/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := common.ParseID(req.FormValue("id"))
		cr.Update(req.Context(), tenant, id, &custDomain.Client{ID: id, TenantID: tenant, Name: req.FormValue("name"), Code: req.FormValue("code"), ClientType: custDomain.ClientType(req.FormValue("type")), ContactName: req.FormValue("contact_name"), ContactPhone: req.FormValue("contact_phone"), ContactEmail: req.FormValue("contact_email"), IsActive: true})
		common.Redirect(w, "/admin/clients")
	}))
	r.POST("/admin/clients/delete", a(func(w http.ResponseWriter, req *http.Request) {
		id, _ := common.ParseID(req.URL.Query().Get("id"))
		cr.Delete(req.Context(), tenant, id)
		common.Redirect(w, "/admin/clients")
	}))

	// ─── /admin/customer-addresses (from admin_modules.go CRM) ───
	r.GET("/admin/customer-addresses", a(func(w http.ResponseWriter, req *http.Request) {
		addrs, _ := ar.List(req.Context(), 0)
		memberNames := map[int64]string{}
		if members, _, _ := mr.List(req.Context(), 0, 0, 500); len(members) > 0 {
			for _, m := range members { memberNames[m.ID] = m.Name }
		}
		rows := make([][]string, len(addrs))
		for i, a := range addrs {
			mn := memberNames[a.MemberID]
			if mn == "" { mn = fmt.Sprintf("会员-%d", a.MemberID) }
			rows[i] = []string{mn, a.RecipientName, a.Phone, a.City, a.District, a.Address, fmt.Sprintf("%v", a.IsDefault)}
		}
		if len(rows) == 0 {
			rows = [][]string{{"暂无数据", "—", "—", "—", "—", "—", "—", "—"}}
		}
		rc.Exec(rc.Tmpl, "crm_addresses", w, "crm_addresses.html", map[string]any{
			"Title": "客户收件地址", "Page": "crm_addresses",
			"Columns": []string{"会员", "收件人", "电话", "城市", "区", "详细地址", "默认"},
			"Rows": rows, "Total": len(rows),
			"AddURL": "/admin/customer-addresses/add-form", "HasActions": true,
		})
	}))
	r.GET("/admin/customer-addresses/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增收件地址")+common.FormSave("/admin/customer-addresses/save")+
			common.FormField("会员ID", "member_id", "1", "")+
			common.FormField("收件人", "recipient_name", "", "")+
			common.FormField("电话", "phone", "", "")+
			common.FormField("邮编", "postal_code", "", "")+
			common.FormField("城市", "city", "", "")+
			common.FormField("区域", "district", "", "")+
			common.FormField("详细地址", "address", "", "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/customer-addresses/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		mid, _ := common.ParseID(req.FormValue("member_id"))
		ar.Create(req.Context(), mid, &custDomain.MemberAddress{RecipientName: req.FormValue("recipient_name"), Phone: req.FormValue("phone"), PostalCode: req.FormValue("postal_code"), City: req.FormValue("city"), District: req.FormValue("district"), Address: req.FormValue("address")})
		common.Redirect(w, "/admin/customer-addresses")
	}))
	r.GET("/admin/customer-addresses/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		id, _ := common.ParseID(req.URL.Query().Get("id"))
		a, _ := ar.GetByID(req.Context(), id)
		if a == nil { http.Error(w, "not found", 404); return }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑收件地址")+common.FormSave("/admin/customer-addresses/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, a.ID)+
			common.FormField("会员ID", "member_id", fmt.Sprintf("%d", a.MemberID), "")+
			common.FormField("收件人", "recipient_name", a.RecipientName, "")+
			common.FormField("电话", "phone", a.Phone, "")+
			common.FormField("邮编", "postal_code", a.PostalCode, "")+
			common.FormField("城市", "city", a.City, "")+
			common.FormField("区域", "district", a.District, "")+
			common.FormField("详细地址", "address", a.Address, "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/customer-addresses/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := common.ParseID(req.FormValue("id"))
		mid, _ := common.ParseID(req.FormValue("member_id"))
		ar.Update(req.Context(), id, &custDomain.MemberAddress{ID: id, MemberID: mid, RecipientName: req.FormValue("recipient_name"), Phone: req.FormValue("phone"), PostalCode: req.FormValue("postal_code"), City: req.FormValue("city"), District: req.FormValue("district"), Address: req.FormValue("address")})
		common.Redirect(w, "/admin/customer-addresses")
	}))

	// ─── /admin/customer-declarants (from admin_modules.go CRM) ───
	r.GET("/admin/customer-declarants", a(func(w http.ResponseWriter, req *http.Request) {
		decls, _, _ := dr.List(req.Context(), 0, 0, 50)
		rows := make([][]string, len(decls))
		for i, d := range decls {
			rows[i] = []string{d.Name, d.IDNumber, string(d.Type), string(d.AuthStatus), common.StatusLabelText(d.IsActive)}
		}
		if len(rows) == 0 {
			rows = [][]string{{"暂无数据", "—", "—", "—", "—", "—", "—", "—"}}
		}
		rc.Exec(rc.Tmpl, "crm_declarants", w, "crm_declarants.html", map[string]any{
			"Title": "客户申报人", "Page": "crm_declarants",
			"Columns": []string{"姓名", "证件号", "类型", "认证状态", "状态"},
			"Rows": rows, "Total": len(decls),
			"AddURL": "/admin/customer-declarants/add-form", "HasActions": true,
		})
	}))
	r.GET("/admin/customer-declarants/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增申报人")+common.FormSave("/admin/customer-declarants/save")+
			common.FormField("客户ID", "client_id", "1", "")+
			common.FormField("姓名", "name", "", "")+
			common.FormField("证件号", "id_number", "", "")+
			common.FormSelect("类型", "type", "individual", [2]string{"individual", "个人"}, [2]string{"company", "公司"})+
			common.FormField("电话", "phone", "", "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/customer-declarants/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		cid, _ := common.ParseID(req.FormValue("client_id"))
		dr.Create(req.Context(), cid, &custDomain.Declarant{ClientID: cid, Type: custDomain.DeclarantType(req.FormValue("type")), Name: req.FormValue("name"), IDNumber: req.FormValue("id_number"), Phone: req.FormValue("phone"), IsActive: true})
		common.Redirect(w, "/admin/customer-declarants")
	}))
	r.GET("/admin/customer-declarants/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		id, _ := common.ParseID(req.URL.Query().Get("id"))
		d, _ := dr.GetByID(req.Context(), 0, id)
		if d == nil { http.Error(w, "not found", 404); return }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑申报人")+common.FormSave("/admin/customer-declarants/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%d"><input type="hidden" name="client_id" value="%d">`, d.ID, d.ClientID)+
			common.FormField("姓名", "name", d.Name, "")+
			common.FormField("证件号", "id_number", d.IDNumber, "")+
			common.FormSelect("类型", "type", string(d.Type), [2]string{"individual", "个人"}, [2]string{"company", "公司"})+
			common.FormField("电话", "phone", d.Phone, "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/customer-declarants/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := common.ParseID(req.FormValue("id"))
		cid, _ := common.ParseID(req.FormValue("client_id"))
		dr.Update(req.Context(), cid, id, &custDomain.Declarant{ID: id, ClientID: cid, Type: custDomain.DeclarantType(req.FormValue("type")), Name: req.FormValue("name"), IDNumber: req.FormValue("id_number"), Phone: req.FormValue("phone"), IsActive: true})
		common.Redirect(w, "/admin/customer-declarants")
	}))

	// ─── /admin/client-accounts (from admin_modules.go CRM) ───
	r.GET("/admin/client-accounts", a(func(w http.ResponseWriter, req *http.Request) {
		clients, _, _ := cr.List(req.Context(), tenant, 0, 50)
		rows := make([][]string, 0)
		for _, c := range clients {
			rows = append(rows, []string{c.Name, c.Code, "运营", c.ContactEmail, common.StatusLabelText(c.IsActive)})
		}
		if len(rows) == 0 {
			rows = [][]string{{"暂无数据", "—", "—", "—", "—", "—", "—", "—"}}
		}
		rc.Exec(rc.Tmpl, "crm_accounts", w, "crm_accounts.html", map[string]any{
			"Title": "客户账号", "Page": "crm_accounts",
			"Columns": []string{"客户", "账号", "角色", "邮箱", "状态"},
			"Rows": rows, "Total": len(rows),
			"AddURL": "/admin/client-accounts/add-form", "HasActions": true,
		})
	}))
	r.GET("/admin/client-accounts/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增客户账号")+common.FormSave("/admin/client-accounts/save")+
			common.FormField("客户ID", "client_id", "1", "")+
			common.FormField("账号", "username", "", "")+
			common.FormField("邮箱", "email", "", "")+
			common.FormSelect("角色", "role", "operator", [2]string{"operator", "运营"}, [2]string{"admin", "管理"})+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/client-accounts/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		cid, _ := common.ParseID(req.FormValue("client_id"))
		cr.Create(req.Context(), tenant, &custDomain.Client{TenantID: tenant, Name: req.FormValue("username"), Code: req.FormValue("username"), ContactEmail: req.FormValue("email"), ClientType: custDomain.ClientTypeNormal, IsActive: true})
		_ = cid
		common.Redirect(w, "/admin/client-accounts")
	}))
	r.GET("/admin/client-accounts/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		id, _ := common.ParseID(req.URL.Query().Get("id"))
		c, _ := cr.GetByID(req.Context(), tenant, id)
		if c == nil { http.Error(w, "not found", 404); return }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑客户账号")+common.FormSave("/admin/client-accounts/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, c.ID)+
			common.FormField("客户名称", "name", c.Name, "")+
			common.FormField("邮箱", "email", c.ContactEmail, "")+
			common.FormField("电话", "phone", c.ContactPhone, "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/client-accounts/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := common.ParseID(req.FormValue("id"))
		cr.Update(req.Context(), tenant, id, &custDomain.Client{ID: id, TenantID: tenant, Name: req.FormValue("name"), ContactEmail: req.FormValue("email"), ContactPhone: req.FormValue("phone"), IsActive: true})
		common.Redirect(w, "/admin/client-accounts")
	}))

	// ─── /admin/client-members (from admin_modules.go CRM) ───
	r.GET("/admin/client-members", a(func(w http.ResponseWriter, req *http.Request) {
		members, _, _ := mr.List(req.Context(), 0, 0, 50)
		clientNames := map[int64]string{}
		if clients, _, _ := cr.List(req.Context(), tenant, 0, 200); len(clients) > 0 {
			for _, c := range clients { clientNames[c.ID] = c.Name }
		}
		rows := make([][]string, len(members))
		for i, m := range members {
			cn := clientNames[m.ClientID]
			if cn == "" { cn = fmt.Sprintf("客户-%d", m.ClientID) }
			rows[i] = []string{m.Name, m.MemberCode, m.Phone, m.Email, cn, common.StatusLabelText(m.IsActive)}
		}
		if len(rows) == 0 {
			rows = [][]string{{"暂无数据", "—", "—", "—", "—", "—", "—", "—"}}
		}
		rc.Exec(rc.Tmpl, "crm_members", w, "crm_members.html", map[string]any{
			"Title": "客户会员", "Page": "crm_members",
			"Columns": []string{"姓名", "会员编号", "电话", "邮箱", "客户", "状态"},
			"Rows": rows, "Total": len(members),
			"AddURL": "/admin/client-members/add-form", "HasActions": true,
		})
	}))
	r.GET("/admin/client-members/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增会员")+common.FormSave("/admin/client-members/save")+
			common.FormField("客户ID", "client_id", "1", "")+
			common.FormField("姓名", "name", "", "")+
			common.FormField("手机", "phone", "", "")+
			common.FormField("邮箱", "email", "", "")+
			common.FormField("会员编号", "member_code", "", "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/client-members/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		cid, _ := common.ParseID(req.FormValue("client_id"))
		mr.Create(req.Context(), cid, &custDomain.ClientMember{ClientID: cid, MemberCode: req.FormValue("member_code"), Name: req.FormValue("name"), Phone: req.FormValue("phone"), Email: req.FormValue("email"), IsActive: true})
		common.Redirect(w, "/admin/client-members")
	}))
	r.GET("/admin/client-members/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		id, _ := common.ParseID(req.URL.Query().Get("id"))
		cid, _ := common.ParseID(req.URL.Query().Get("client_id"))
		if cid == 0 { cid = 1 }
		m, _ := mr.GetByID(req.Context(), cid, id)
		if m == nil { http.Error(w, "not found", 404); return }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑会员")+common.FormSave("/admin/client-members/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%d"><input type="hidden" name="client_id" value="%d">`, m.ID, m.ClientID)+
			common.FormField("姓名", "name", m.Name, "")+
			common.FormField("手机", "phone", m.Phone, "")+
			common.FormField("邮箱", "email", m.Email, "")+
			common.FormField("会员编号", "member_code", m.MemberCode, "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/client-members/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := common.ParseID(req.FormValue("id"))
		cid, _ := common.ParseID(req.FormValue("client_id"))
		mr.Update(req.Context(), cid, id, &custDomain.ClientMember{ID: id, ClientID: cid, MemberCode: req.FormValue("member_code"), Name: req.FormValue("name"), Phone: req.FormValue("phone"), Email: req.FormValue("email"), IsActive: true})
		common.Redirect(w, "/admin/client-members")
	}))

	// ─── /admin/client-recharge (from admin_modules.go CRM) ───
	r.GET("/admin/client-recharge", a(func(w http.ResponseWriter, req *http.Request) {
		clients, _, _ := cr.List(req.Context(), tenant, 0, 50)
		clientOpts := ""
		for _, c := range clients {
			clientOpts += fmt.Sprintf(`<option value="%d">%s (余额: ¥%.2f)</option>`, c.ID, c.Name, c.Balance)
		}
		common.HtmlOK(w)
		fmt.Fprint(w, `<!DOCTYPE html><html lang="zh-CN" data-theme="light"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>客户充值 - I56</title><link rel="stylesheet" href="/static/css/i56-bdl.css"><script src="/static/js/i56-theme.js"></script><script src="https://unpkg.com/htmx.org@1.9.10"></script><style>
*{margin:0;padding:0;box-sizing:border-box}body{font-family:system-ui,sans-serif;background:var(--i56-bg-base);color:var(--i56-text-primary);padding:16px}
.i56-card{background:var(--i56-bg-surface);border:1px solid var(--i56-border);border-radius:8px;padding:24px;max-width:480px;margin:40px auto}
.i56-card h3{font-size:16px;margin-bottom:16px;color:var(--i56-brand)}
.form-group{margin-bottom:12px}.form-label{display:block;font-size:12px;color:var(--i56-text-secondary);margin-bottom:4px}
.form-input{width:100%;padding:8px 10px;font-size:13px;background:var(--i56-bg-base);color:var(--i56-text-primary);border:1px solid var(--i56-border);border-radius:6px}
.i56-btn.i56-btn-primary{background:var(--i56-brand);color:#fff;border:none;padding:10px 24px;border-radius:6px;font-size:13px;cursor:pointer;width:100%}
.i56-btn.i56-btn-primary:hover{opacity:.9}
</style></head><body><div class="i56-card"><h3>💰 客户充值</h3>
<form hx-post="/admin/client-recharge" hx-swap="none">
<div class="form-group"><label class="form-label">客户</label><select name="client_id" class="form-input">`+clientOpts+`</select></div>
<div class="form-group"><label class="form-label">充值金额 (元)</label><input name="amount" class="form-input" type="number" step="0.01" placeholder="充值金额"></div>
<div class="form-group"><label class="form-label">支付方式</label><select name="method" class="form-input"><option value="bank_transfer">银行转账</option><option value="wechat">微信支付</option><option value="alipay">支付宝</option><option value="cash">现金</option></select></div>
<div class="form-group"><label class="form-label">备注</label><input name="description" class="form-input" placeholder="备注信息"></div>
<button type="submit" class="i56-btn i56-btn-primary">确认充值</button>
</form></div></body></html>`)
	}))
	r.POST("/admin/client-recharge", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		cid, _ := common.ParseID(req.FormValue("client_id"))
		amt, _ := common.ParseFloat(req.FormValue("amount"))
		c, _ := cr.GetByID(req.Context(), tenant, cid)
		balanceAfter := amt
		if c != nil { balanceAfter = c.Balance + amt; c.Balance = balanceAfter; cr.Update(req.Context(), tenant, cid, c) }
		lr.Add(req.Context(), &custRepo.LedgerEntry{TenantID: tenant, ClientID: cid, Amount: amt, BalanceAfter: balanceAfter, Type: req.FormValue("method"), Description: req.FormValue("description")})
		common.Redirect(w, "/admin/balance-logs")
	}))

	// ─── /admin/balance-logs (from admin_modules.go CRM) ───
	r.GET("/admin/balance-logs", a(func(w http.ResponseWriter, req *http.Request) {
		entries, _, _ := lr.List(req.Context(), tenant, 0, 0, 50)
		clientNames := map[int64]string{}
		if clients, _, _ := cr.List(req.Context(), tenant, 0, 200); len(clients) > 0 {
			for _, c := range clients { clientNames[c.ID] = c.Name }
		}
		rows := make([][]string, len(entries))
		for i, e := range entries {
			typ := "扣款"; if e.Amount > 0 { typ = "充值" }
			cn := clientNames[e.ClientID]
			if cn == "" { cn = fmt.Sprintf("客户-%d", e.ClientID) }
			rows[i] = []string{cn, typ, fmt.Sprintf("¥%.2f", e.Amount), fmt.Sprintf("¥%.2f", e.BalanceAfter), e.Description, e.CreatedAt.Format("01-02 15:04")}
		}
		if len(rows) == 0 {
			rows = [][]string{{"暂无数据", "—", "—", "—", "—", "—", "—", "—"}}
		}
		gp(w, "crm_ledgers", "余额日志", len(rows), []string{"客户", "类型", "金额", "余额", "描述", "时间"}, rows, "/admin/balance-logs/add-form")
	}))
	r.GET("/admin/balance-logs/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增余额记录")+common.FormSave("/admin/balance-logs/save")+
			common.FormField("客户ID", "client_id", "1", "")+
			common.FormField("金额", "amount", "", "正数=充值, 负数=扣款")+
			common.FormSelect("类型", "type", "manual", [2]string{"recharge", "充值"}, [2]string{"charge", "扣款"}, [2]string{"refund", "退款"}, [2]string{"manual", "手动调整"})+
			common.FormField("描述", "description", "", "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/balance-logs/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		cid, _ := common.ParseID(req.FormValue("client_id"))
		amt, _ := common.ParseFloat(req.FormValue("amount"))
		c, _ := cr.GetByID(req.Context(), tenant, cid)
		bal := amt
		if c != nil { bal = c.Balance + amt; c.Balance = bal; cr.Update(req.Context(), tenant, cid, c) }
		lr.Add(req.Context(), &custRepo.LedgerEntry{TenantID: tenant, ClientID: cid, Amount: amt, BalanceAfter: bal, Type: req.FormValue("type"), Description: req.FormValue("description")})
		common.Redirect(w, "/admin/balance-logs")
	}))
	r.GET("/admin/balance-logs/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		id, _ := common.ParseID(req.URL.Query().Get("id"))
		e, _ := lr.GetByEntryID(req.Context(), id)
		if e == nil { http.Error(w, "not found", 404); return }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑余额记录")+common.FormSave("/admin/balance-logs/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, e.ID)+
			common.FormField("客户ID", "client_id", fmt.Sprintf("%d", e.ClientID), "")+
			common.FormField("金额", "amount", fmt.Sprintf("%.2f", e.Amount), "正数=充值, 负数=扣款")+
			common.FormSelect("类型", "type", e.Type, [2]string{"recharge", "充值"}, [2]string{"charge", "扣款"}, [2]string{"refund", "退款"}, [2]string{"manual", "手动调整"})+
			common.FormField("描述", "description", e.Description, "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/balance-logs/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := common.ParseID(req.FormValue("id"))
		cid, _ := common.ParseID(req.FormValue("client_id"))
		amt, _ := common.ParseFloat(req.FormValue("amount"))
		lr.Update(req.Context(), id, &custRepo.LedgerEntry{ID: id, TenantID: tenant, ClientID: cid, Amount: amt, BalanceAfter: amt, Type: req.FormValue("type"), Description: req.FormValue("description")})
		common.Redirect(w, "/admin/balance-logs")
	}))

	// ─── /admin/recharge-records (from admin_modules.go CRM) ───
	r.GET("/admin/recharge-records", a(func(w http.ResponseWriter, req *http.Request) {
		entries, _, _ := lr.List(req.Context(), tenant, 0, 0, 50)
		clientNames := map[int64]string{}
		if clients, _, _ := cr.List(req.Context(), tenant, 0, 200); len(clients) > 0 {
			for _, c := range clients { clientNames[c.ID] = c.Name }
		}
		var rows [][]string
		for _, e := range entries {
			if e.Amount > 0 {
				cn := clientNames[e.ClientID]; if cn == "" { cn = fmt.Sprintf("客户-%d", e.ClientID) }
				rows = append(rows, []string{cn, fmt.Sprintf("¥%.2f", e.Amount), e.Type, e.CreatedAt.Format("01-02 15:04"), "已完成"})
			}
		}
		if len(rows) == 0 {
			rows = [][]string{{"暂无数据", "—", "—", "—", "—", "—", "—", "—"}}
		}
		gp(w, "crm_recharge_records", "充值记录", len(rows), []string{"客户", "金额", "方式", "时间", "状态"}, rows, "/admin/recharge-records/add-form")
	}))
	r.GET("/admin/recharge-records/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增充值记录")+common.FormSave("/admin/recharge-records/save")+
			common.FormField("客户ID", "client_id", "1", "")+
			common.FormField("金额", "amount", "", "充值金额")+
			common.FormSelect("方式", "method", "bank_transfer", [2]string{"bank_transfer", "银行转账"}, [2]string{"wechat", "微信支付"}, [2]string{"alipay", "支付宝"}, [2]string{"cash", "现金"})+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/recharge-records/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		cid, _ := common.ParseID(req.FormValue("client_id"))
		amt, _ := common.ParseFloat(req.FormValue("amount"))
		c, _ := cr.GetByID(req.Context(), tenant, cid)
		balanceAfter := amt
		if c != nil { balanceAfter = c.Balance + amt; c.Balance = balanceAfter; cr.Update(req.Context(), tenant, cid, c) }
		lr.Add(req.Context(), &custRepo.LedgerEntry{TenantID: tenant, ClientID: cid, Amount: amt, BalanceAfter: balanceAfter, Type: req.FormValue("method"), Description: ""})
		common.Redirect(w, "/admin/recharge-records")
	}))

	// ─── /admin/client-pricing (from admin_modules.go CRM) ───
	r.GET("/admin/client-pricing", a(func(w http.ResponseWriter, req *http.Request) {
		prices := rpr.List()
		rows := make([][]string, len(prices))
		for i, p := range prices {
			rows[i] = []string{p.RouteName, p.TransportType, p.CargoType, p.TaxType, p.FirstWeightPrice, p.AdditionalWeightPrice, p.MinCharge}
		}
		if len(rows) == 0 {
			rows = [][]string{{"暂无数据", "—", "—", "—", "—", "—", "—", "—"}}
		}
		gp(w, "crm_pricing", "客户价格", len(rows), []string{"线路", "运输方式", "货类", "税档", "首重价", "续重价", "最低收费"}, rows, "/admin/client-pricing/add-form")
	}))
	r.GET("/admin/client-pricing/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增客户价格")+common.FormSave("/admin/client-pricing/save")+
			common.FormField("线路名", "route_name", "", "如: 厦门→台湾(空运)")+
			common.FormSelect("运输方式", "transport_type", "air", [2]string{"air", "空运"}, [2]string{"sea_express", "海快"}, [2]string{"sea", "海运"})+
			common.FormField("货类", "cargo_type", "", "如: general, class1")+
			common.FormSelect("税档", "tax_type", "full_inclusive", [2]string{"full_inclusive", "全包税"}, [2]string{"tax_excluded", "不含税"})+
			common.FormField("首重(kg)", "first_weight", "", "")+
			common.FormField("首重价格(元)", "first_weight_price", "", "")+
			common.FormField("续重价格(元/kg)", "additional_weight_price", "", "")+
			common.FormField("最低收费(元)", "min_charge", "", "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/client-pricing/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		rpr.Add(pricingRepo.ClientRoutePriceDisplay{RouteName: req.FormValue("route_name"), TransportType: req.FormValue("transport_type"), CargoType: req.FormValue("cargo_type"), TaxType: req.FormValue("tax_type"), FirstWeight: req.FormValue("first_weight"), FirstWeightPrice: req.FormValue("first_weight_price"), AdditionalWeightPrice: req.FormValue("additional_weight_price"), MinCharge: req.FormValue("min_charge")})
		common.Redirect(w, "/admin/client-pricing")
	}))
	r.GET("/admin/client-pricing/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("id")
		prices := rpr.List()
		var p *pricingRepo.ClientRoutePriceDisplay
		for i := range prices {
			if prices[i].RouteName == name { p = &prices[i]; break }
		}
		if p == nil { http.Error(w, "not found", 404); return }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑客户价格")+common.FormSave("/admin/client-pricing/update")+
			fmt.Sprintf(`<input type="hidden" name="old_name" value="%s">`, p.RouteName)+
			common.FormField("线路名", "route_name", p.RouteName, "如: 厦门→台湾(空运)")+
			common.FormSelect("运输方式", "transport_type", p.TransportType, [2]string{"air", "空运"}, [2]string{"sea_express", "海快"}, [2]string{"sea", "海运"})+
			common.FormField("货类", "cargo_type", p.CargoType, "")+
			common.FormSelect("税档", "tax_type", p.TaxType, [2]string{"full_inclusive", "全包税"}, [2]string{"tax_excluded", "不含税"})+
			common.FormField("首重(kg)", "first_weight", p.FirstWeight, "")+
			common.FormField("首重价格(元)", "first_weight_price", p.FirstWeightPrice, "")+
			common.FormField("续重价格(元/kg)", "additional_weight_price", p.AdditionalWeightPrice, "")+
			common.FormField("最低收费(元)", "min_charge", p.MinCharge, "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/client-pricing/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		oldName := req.FormValue("old_name")
		prices := rpr.List()
		for i := range prices {
			if prices[i].RouteName == oldName {
				rpr.Remove(i)
				rpr.Add(pricingRepo.ClientRoutePriceDisplay{RouteName: req.FormValue("route_name"), TransportType: req.FormValue("transport_type"), CargoType: req.FormValue("cargo_type"), TaxType: req.FormValue("tax_type"), FirstWeight: req.FormValue("first_weight"), FirstWeightPrice: req.FormValue("first_weight_price"), AdditionalWeightPrice: req.FormValue("additional_weight_price"), MinCharge: req.FormValue("min_charge")})
				break
			}
		}
		common.Redirect(w, "/admin/client-pricing")
	}))

	// ─── /admin/monthly-statements (from admin_modules.go CRM) ───
	r.GET("/admin/monthly-statements", a(func(w http.ResponseWriter, req *http.Request) {
		clients, _, _ := cr.List(req.Context(), tenant, 0, 50)
		now := time.Now()
		rows := make([][]string, len(clients))
		for i, c := range clients {
			rows[i] = []string{c.Name, now.Format("2006-01"), fmt.Sprintf("¥%.2f", c.Balance), fmt.Sprintf("¥%.2f", 0.0), "未结算"}
		}
		if len(rows) == 0 {
			rows = [][]string{
				{"EZ集运通", now.Format("2006-01"), "¥10,000.00", "¥2,500.00", "未结算"},
				{"拼拼侠", now.Format("2006-01"), "¥5,000.00", "¥1,200.00", "未结算"},
			}
		}
		gp(w, "crm_statements", "月结对账单", len(rows), []string{"客户", "账期", "期末余额", "已结算", "状态"}, rows, "/admin/monthly-statements/add-form")
	}))
	r.GET("/admin/monthly-statements/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增月结账单")+common.FormSave("/admin/monthly-statements/save")+
			common.FormField("客户ID", "client_id", "1", "")+
			common.FormField("账期", "period", "", "如: 2026-07")+
			common.FormField("期末余额", "ending_balance", "", "")+
			common.FormField("已结算金额", "settled_amount", "", "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/monthly-statements/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		cid, _ := common.ParseID(req.FormValue("client_id"))
		amt, _ := common.ParseFloat(req.FormValue("ending_balance"))
		lr.Add(req.Context(), &custRepo.LedgerEntry{TenantID: tenant, ClientID: cid, Amount: amt, BalanceAfter: amt, Type: "statement", Description: "月结账单 " + req.FormValue("period")})
		common.Redirect(w, "/admin/monthly-statements")
	}))

	// ─── /admin/client-ledgers CRUD (from admin_crud.go) ───
	// List page — shows real ledger data from repo
	r.GET("/admin/client-ledgers", a(func(w http.ResponseWriter, req *http.Request) {
		entries, _, _ := lr.List(req.Context(), tenant, 0, 0, 50)
		clientNames := map[int64]string{}
		if clients, _, _ := cr.List(req.Context(), tenant, 0, 200); len(clients) > 0 {
			for _, c := range clients { clientNames[c.ID] = c.Name }
		}
		rows := make([][]string, len(entries))
		for i, e := range entries {
			typ := "扣款"
			if e.Amount > 0 { typ = "充值" }
			cn := clientNames[e.ClientID]
			if cn == "" { cn = fmt.Sprintf("客户-%d", e.ClientID) }
			rows[i] = []string{
				fmt.Sprintf("%d", e.ID), cn, typ,
				fmt.Sprintf("¥%.2f", e.Amount),
				fmt.Sprintf("¥%.2f", e.BalanceAfter),
				e.Description,
				e.CreatedAt.Format("01-02 15:04"),
			}
		}
		if len(rows) == 0 {
			rows = [][]string{{"暂无数据", "—", "—", "—", "—", "—", "—", "—"}}
		}
		gp(w, "crm_client_ledgers", "余额流水", len(rows),
			[]string{"ID", "客户", "类型", "金额", "余额", "描述", "时间"},
			rows, "/admin/client-ledgers/add-form",
		)
	}))
	r.GET("/admin/client-ledgers/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增余额记录")+common.FormSave("/admin/client-ledgers/save")+
			common.FormField("客户ID", "client_id", "1", "")+
			common.FormField("金额", "amount", "", "正数=充值, 负数=扣款")+
			common.FormSelect("类型", "type", "manual", [2]string{"recharge", "充值"}, [2]string{"charge", "扣款"}, [2]string{"refund", "退款"}, [2]string{"manual", "手动调整"})+
			common.FormField("描述", "description", "", "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/client-ledgers/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		cid, _ := common.ParseID(req.FormValue("client_id"))
		amt, _ := common.ParseFloat(req.FormValue("amount"))
		c, _ := cr.GetByID(req.Context(), tenant, cid)
		balanceAfter := amt
		if c != nil { balanceAfter = c.Balance + amt; c.Balance = balanceAfter; cr.Update(req.Context(), tenant, cid, c) }
		lr.Add(req.Context(), &custRepo.LedgerEntry{TenantID: tenant, ClientID: cid, Amount: amt, BalanceAfter: balanceAfter, Type: req.FormValue("type"), Description: req.FormValue("description")})
		common.Redirect(w, "/admin/client-ledgers")
	}))

	// ─── /admin/client-recharges CRUD (from admin_crud.go) ───
	// List page — shows recharge records with real data from ledger
	r.GET("/admin/client-recharges", a(func(w http.ResponseWriter, req *http.Request) {
		entries, _, _ := lr.List(req.Context(), tenant, 0, 0, 50)
		clientNames := map[int64]string{}
		if clients, _, _ := cr.List(req.Context(), tenant, 0, 200); len(clients) > 0 {
			for _, c := range clients { clientNames[c.ID] = c.Name }
		}
		var rows [][]string
		for _, e := range entries {
			if e.Amount > 0 {
				cn := clientNames[e.ClientID]; if cn == "" { cn = fmt.Sprintf("客户-%d", e.ClientID) }
				methodCN := "银行转账"
				switch e.Type {
				case "bank_transfer": methodCN = "银行转账"
				case "wechat": methodCN = "微信支付"
				case "alipay": methodCN = "支付宝"
				case "cash": methodCN = "现金"
				default: methodCN = e.Type
				}
				rows = append(rows, []string{cn, fmt.Sprintf("¥%.2f", e.Amount), methodCN, e.Description, e.CreatedAt.Format("2006-01-02 15:04"), "已完成"})
			}
		}
		if len(rows) == 0 {
			rows = [][]string{{"暂无数据", "—", "—", "—", "—", "—", "—", "—"}}
		}
		gp(w, "crm_client_recharges", "客户充值记录", len(rows), []string{"客户", "金额", "方式", "备注", "时间", "状态"}, rows, "/admin/client-recharges/add-form")
	}))
	r.GET("/admin/client-recharges/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("客户充值")+common.FormSave("/admin/client-recharges/save")+
			common.FormField("客户ID", "client_id", "1", "")+
			common.FormField("充值金额", "amount", "", "")+
			common.FormSelect("方式", "method", "bank_transfer", [2]string{"bank_transfer", "银行转账"}, [2]string{"wechat", "微信支付"}, [2]string{"alipay", "支付宝"}, [2]string{"cash", "现金"})+
			common.FormField("备注", "description", "", "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/client-recharges/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		cid, _ := common.ParseID(req.FormValue("client_id"))
		amt, _ := common.ParseFloat(req.FormValue("amount"))
		c, _ := cr.GetByID(req.Context(), tenant, cid)
		balanceAfter := amt
		if c != nil { balanceAfter = c.Balance + amt; c.Balance = balanceAfter; cr.Update(req.Context(), tenant, cid, c) }
		lr.Add(req.Context(), &custRepo.LedgerEntry{TenantID: tenant, ClientID: cid, Amount: amt, BalanceAfter: balanceAfter, Type: req.FormValue("method"), Description: req.FormValue("description")})
		common.Redirect(w, "/admin/client-recharges")
	}))
	r.GET("/admin/client-recharges/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		id, _ := common.ParseID(req.URL.Query().Get("id"))
		e, _ := lr.GetByEntryID(req.Context(), id)
		if e == nil { http.Error(w, "not found", 404); return }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑充值")+common.FormSave("/admin/client-recharges/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, e.ID)+
			common.FormField("客户ID", "client_id", fmt.Sprintf("%d", e.ClientID), "")+
			common.FormField("充值金额", "amount", fmt.Sprintf("%.2f", e.Amount), "")+
			common.FormSelect("方式", "method", e.Type, [2]string{"bank_transfer", "银行转账"}, [2]string{"wechat", "微信支付"}, [2]string{"alipay", "支付宝"}, [2]string{"cash", "现金"})+
			common.FormField("备注", "description", e.Description, "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/client-recharges/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := common.ParseID(req.FormValue("id"))
		cid, _ := common.ParseID(req.FormValue("client_id"))
		amt, _ := common.ParseFloat(req.FormValue("amount"))
		lr.Update(req.Context(), id, &custRepo.LedgerEntry{ID: id, TenantID: tenant, ClientID: cid, Amount: amt, BalanceAfter: amt, Type: req.FormValue("method"), Description: req.FormValue("description")})
		common.Redirect(w, "/admin/client-recharges")
	}))

	// ─── /admin/declarants CRUD (from admin_crud.go) ───
	// List page — shows real declarant data from repo
	r.GET("/admin/declarants", a(func(w http.ResponseWriter, req *http.Request) {
		decls, total, _ := dr.List(req.Context(), 0, 0, 50)
		rows := make([][]string, len(decls))
		for i, d := range decls {
			rows[i] = []string{
				fmt.Sprintf("%d", d.ID), d.Name, d.IDNumber,
				string(d.Type), string(d.AuthStatus),
				common.StatusLabelText(d.IsActive),
			}
		}
		if len(rows) == 0 {
			rows = [][]string{{"暂无数据", "—", "—", "—", "—", "—", "—", "—"}}
		}
		gp(w, "crm_declarants2", "申报人管理", int(total),
			[]string{"ID", "姓名", "证件号", "类型", "认证状态", "状态"},
			rows, "/admin/declarants/add-form",
		)
	}))
	r.GET("/admin/declarants/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增申报人")+common.FormSave("/admin/declarants/save")+
			common.FormField("客户ID", "client_id", "1", "")+
			common.FormField("姓名", "name", "", "")+
			common.FormField("证件号", "id_number", "", "")+
			common.FormSelect("类型", "type", "individual", [2]string{"individual", "个人"}, [2]string{"company", "公司"})+
			common.FormField("电话", "phone", "", "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/declarants/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		cid, _ := common.ParseID(req.FormValue("client_id"))
		dr.Create(req.Context(), cid, &custDomain.Declarant{ClientID: cid, Type: custDomain.DeclarantType(req.FormValue("type")), Name: req.FormValue("name"), IDNumber: req.FormValue("id_number"), Phone: req.FormValue("phone"), IsActive: true})
		common.Redirect(w, "/admin/declarants")
	}))
	r.GET("/admin/declarants/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		id, _ := common.ParseID(req.URL.Query().Get("id"))
		d, _ := dr.GetByID(req.Context(), 0, id)
		if d == nil { http.Error(w, "not found", 404); return }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑申报人")+common.FormSave("/admin/declarants/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%d"><input type="hidden" name="client_id" value="%d">`, d.ID, d.ClientID)+
			common.FormField("姓名", "name", d.Name, "")+
			common.FormField("证件号", "id_number", d.IDNumber, "")+
			common.FormSelect("类型", "type", string(d.Type), [2]string{"individual", "个人"}, [2]string{"company", "公司"})+
			common.FormField("电话", "phone", d.Phone, "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/declarants/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := common.ParseID(req.FormValue("id"))
		cid, _ := common.ParseID(req.FormValue("client_id"))
		dr.Update(req.Context(), cid, id, &custDomain.Declarant{ID: id, ClientID: cid, Type: custDomain.DeclarantType(req.FormValue("type")), Name: req.FormValue("name"), IDNumber: req.FormValue("id_number"), Phone: req.FormValue("phone"), IsActive: true})
		common.Redirect(w, "/admin/declarants")
	}))

	// ─── /admin/client-permissions — 客户端权限 ───
	// Enhanced BFT56 model: permission types, approval workflow, levels, expiry
	type ClientPermission struct {
		ID          int64
		ClientID    int64
		ClientName  string
		Permission  string // type of permission
		Level       string // "basic" | "pro" | "enterprise"
		Status      string // "pending" | "approved" | "active" | "expired" | "revoked"
		RequestedAt string
		ApprovedAt  string
		ExpiresAt   string // expiry date
		ApprovedBy  string
		Remarks     string
	}
	permMu := &sync.RWMutex{}
	permSeq := int64(0)
	var permissions []ClientPermission
	// Seed enhanced permission data for EZ集运通
	clients, _, _ := cr.List(context.TODO(), tenant, 0, 50)
	for _, c := range clients {
		if c.Name == "EZ集运通" {
			permissions = []ClientPermission{
				{ID: 1, ClientID: c.ID, ClientName: c.Name, Permission: "报关申报", Level: "enterprise", Status: "active", RequestedAt: "2026-01-15", ApprovedAt: "2026-01-16", ExpiresAt: "2027-01-15", ApprovedBy: "系统管理员", Remarks: ""},
				{ID: 2, ClientID: c.ID, ClientName: c.Name, Permission: "API数据对接", Level: "pro", Status: "active", RequestedAt: "2026-03-01", ApprovedAt: "2026-03-02", ExpiresAt: "2027-03-01", ApprovedBy: "系统管理员", Remarks: ""},
				{ID: 3, ClientID: c.ID, ClientName: c.Name, Permission: "仓储服务", Level: "pro", Status: "active", RequestedAt: "2026-06-01", ApprovedAt: "2026-06-02", ExpiresAt: "2027-06-01", ApprovedBy: "运营主管", Remarks: ""},
				{ID: 4, ClientID: c.ID, ClientName: c.Name, Permission: "物流追踪", Level: "basic", Status: "active", RequestedAt: "2026-06-01", ApprovedAt: "2026-06-01", ExpiresAt: "2027-06-01", ApprovedBy: "系统管理员", Remarks: ""},
				{ID: 5, ClientID: c.ID, ClientName: c.Name, Permission: "财务结算", Level: "pro", Status: "active", RequestedAt: "2026-06-01", ApprovedAt: "2026-06-02", ExpiresAt: "2027-06-01", ApprovedBy: "财务主管", Remarks: ""},
				{ID: 6, ClientID: c.ID, ClientName: c.Name, Permission: "海快专线", Level: "pro", Status: "pending", RequestedAt: "2026-07-10", ApprovedAt: "—", ExpiresAt: "—", ApprovedBy: "—", Remarks: "待审核"},
				{ID: 7, ClientID: c.ID, ClientName: c.Name, Permission: "会员管理", Level: "basic", Status: "active", RequestedAt: "2026-06-01", ApprovedAt: "2026-06-01", ExpiresAt: "2027-06-01", ApprovedBy: "系统管理员", Remarks: ""},
				{ID: 8, ClientID: c.ID, ClientName: c.Name, Permission: "电子面单", Level: "pro", Status: "active", RequestedAt: "2026-03-01", ApprovedAt: "2026-03-02", ExpiresAt: "2027-03-01", ApprovedBy: "系统管理员", Remarks: ""},
				{ID: 9, ClientID: c.ID, ClientName: c.Name, Permission: "通知推送", Level: "basic", Status: "revoked", RequestedAt: "2026-01-15", ApprovedAt: "2026-01-16", ExpiresAt: "2026-12-31", ApprovedBy: "系统管理员", Remarks: "客户主动关闭"},
			}
			permSeq = 9
			break
		}
	}
	permLevelLabel := func(l string) string {
		switch l {
		case "basic": return "基础版"
		case "pro": return "专业版"
		case "enterprise": return "企业版"
		default: return "—"
		}
	}
	permStatusLabel := func(s string) string {
		switch s {
		case "pending": return "审核中"
		case "approved": return "已批准"
		case "active": return "已开通"
		case "expired": return "已过期"
		case "revoked": return "已停用"
		default: return s
		}
	}
	r.GET("/admin/client-permissions", a(func(w http.ResponseWriter, req *http.Request) {
		permMu.RLock()
		rows := make([][]string, len(permissions))
		for i, p := range permissions {
			rows[i] = []string{
				p.ClientName, p.Permission, permLevelLabel(p.Level),
				permStatusLabel(p.Status), p.RequestedAt, p.ApprovedAt, p.ExpiresAt,
			}
		}
		permMu.RUnlock()
		if len(rows) == 0 {
			rows = [][]string{
				{"EZ集运通", "报关申报", "企业版", "已开通", "2026-01-15", "2026-01-16", "2027-01-15"},
			}
		}
		gp(w, "crm_client_permissions", "客户端权限", len(rows), []string{"客户", "权限类型", "等级", "状态", "申请时间", "批准时间", "到期时间"}, rows, "/admin/client-permissions/add-form")
	}))
	r.GET("/admin/client-permissions/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		cl, _, _ := cr.List(req.Context(), tenant, 0, 50)
		clientOpts := ""
		for _, c := range cl {
			clientOpts += fmt.Sprintf(`<option value="%d">%s</option>`, c.ID, c.Name)
		}
		fmt.Fprint(w, common.ModalStart("新增客户端权限")+common.FormSave("/admin/client-permissions/save")+
			fmt.Sprintf(`<div class="form-group"><label class="form-label">客户</label><select name="client_id" class="form-input">%s</select></div>`, clientOpts)+
			common.FormSelect("权限类型", "permission", "报关申报",
				[2]string{"报关申报", "报关申报"}, [2]string{"API数据对接", "API数据对接"},
				[2]string{"仓储服务", "仓储服务"}, [2]string{"物流追踪", "物流追踪"},
				[2]string{"财务结算", "财务结算"}, [2]string{"会员管理", "会员管理"},
				[2]string{"电子面单", "电子面单"}, [2]string{"通知推送", "通知推送"})+
			common.FormSelect("权限等级", "level", "basic",
				[2]string{"basic", "基础版"}, [2]string{"pro", "专业版"}, [2]string{"enterprise", "企业版"})+
			common.FormSelect("状态", "status", "pending",
				[2]string{"pending", "审核中"}, [2]string{"approved", "已批准"},
				[2]string{"active", "已开通"}, [2]string{"revoked", "已停用"})+
			common.FormField("到期时间", "expires_at", "", "如: 2027-06-01")+
			common.FormField("备注", "remarks", "", "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/client-permissions/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		cid, _ := common.ParseID(req.FormValue("client_id"))
		clientName := ""
		if c, _ := cr.GetByID(req.Context(), tenant, cid); c != nil { clientName = c.Name }
		permMu.Lock()
		permSeq++
		permissions = append(permissions, ClientPermission{
			ID: permSeq, ClientID: cid, ClientName: clientName,
			Permission: req.FormValue("permission"), Level: req.FormValue("level"),
			Status: req.FormValue("status"), RequestedAt: time.Now().Format("2006-01-02"),
			ApprovedAt: "—", ExpiresAt: req.FormValue("expires_at"),
			ApprovedBy: "—", Remarks: req.FormValue("remarks"),
		})
		permMu.Unlock()
		common.Redirect(w, "/admin/client-permissions")
	}))
	r.GET("/admin/client-permissions/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		id, _ := common.ParseID(req.URL.Query().Get("id"))
		permMu.RLock()
		var p *ClientPermission
		for i := range permissions { if permissions[i].ID == id { p = &permissions[i]; break } }
		permMu.RUnlock()
		if p == nil { http.Error(w, "not found", 404); return }
		common.HtmlOK(w)
		cl, _, _ := cr.List(req.Context(), tenant, 0, 50)
		clientOpts := ""
		for _, c := range cl {
			sel := ""
			if c.ID == p.ClientID { sel = " selected" }
			clientOpts += fmt.Sprintf(`<option value="%d"%s>%s</option>`, c.ID, sel, c.Name)
		}
		fmt.Fprint(w, common.ModalStart("编辑客户端权限")+common.FormSave("/admin/client-permissions/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, p.ID)+
			fmt.Sprintf(`<div class="form-group"><label class="form-label">客户</label><select name="client_id" class="form-input">%s</select></div>`, clientOpts)+
			common.FormSelect("权限类型", "permission", p.Permission,
				[2]string{"报关申报", "报关申报"}, [2]string{"API数据对接", "API数据对接"},
				[2]string{"仓储服务", "仓储服务"}, [2]string{"物流追踪", "物流追踪"},
				[2]string{"财务结算", "财务结算"}, [2]string{"会员管理", "会员管理"},
				[2]string{"电子面单", "电子面单"}, [2]string{"通知推送", "通知推送"})+
			common.FormSelect("权限等级", "level", p.Level,
				[2]string{"basic", "基础版"}, [2]string{"pro", "专业版"}, [2]string{"enterprise", "企业版"})+
			common.FormSelect("状态", "status", p.Status,
				[2]string{"pending", "审核中"}, [2]string{"approved", "已批准"},
				[2]string{"active", "已开通"}, [2]string{"expired", "已过期"},
				[2]string{"revoked", "已停用"})+
			common.FormField("到期时间", "expires_at", p.ExpiresAt, "如: 2027-06-01")+
			common.FormField("审批人", "approved_by", p.ApprovedBy, "")+
			common.FormField("备注", "remarks", p.Remarks, "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/client-permissions/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := common.ParseID(req.FormValue("id"))
		cid, _ := common.ParseID(req.FormValue("client_id"))
		clientName := ""
		if c, _ := cr.GetByID(req.Context(), tenant, cid); c != nil { clientName = c.Name }
		permMu.Lock()
		for i := range permissions {
			if permissions[i].ID == id {
				permissions[i].ClientID = cid
				permissions[i].ClientName = clientName
				permissions[i].Permission = req.FormValue("permission")
				permissions[i].Level = req.FormValue("level")
				newStatus := req.FormValue("status")
				permissions[i].Status = newStatus
				if (newStatus == "approved" || newStatus == "active") && permissions[i].ApprovedAt == "—" {
					permissions[i].ApprovedAt = time.Now().Format("2006-01-02")
				}
				permissions[i].ExpiresAt = req.FormValue("expires_at")
				permissions[i].ApprovedBy = req.FormValue("approved_by")
				permissions[i].Remarks = req.FormValue("remarks")
				break
			}
		}
		permMu.Unlock()
		common.Redirect(w, "/admin/client-permissions")
	}))

	// ─── Pricing pages (from admin_modules.go) ───
	pmr := pricingRepo.NewMemPricingModelsRepo()
	registerPricingRoutes(r, a, gp, pmr)

	_ = pricingDomain.RoutePriceModel{}
}

func registerPricingRoutes(
	r *router.Router,
	a func(http.HandlerFunc) http.HandlerFunc,
	gp common.GenericListFunc,
	pmr interface {
		ListRoutePrices() []pricingDomain.RoutePriceModel
		AddRoutePrice(p *pricingDomain.RoutePriceModel)
		GetRoutePriceByID(id int64) *pricingDomain.RoutePriceModel
		UpdateRoutePrice(id int64, p *pricingDomain.RoutePriceModel)
		ListDeliveryFees() []pricingDomain.DeliveryFeeModel
		AddDeliveryFee(f *pricingDomain.DeliveryFeeModel)
		GetDeliveryFeeByID(id int64) *pricingDomain.DeliveryFeeModel
		UpdateDeliveryFee(id int64, f *pricingDomain.DeliveryFeeModel)
		ListSurcharges() []pricingDomain.SurchargeModel
		AddSurcharge(s *pricingDomain.SurchargeModel)
		GetSurchargeByID(id int64) *pricingDomain.SurchargeModel
		UpdateSurcharge(id int64, s *pricingDomain.SurchargeModel)
		ListServicePrices() []pricingDomain.ServicePriceModel
		AddServicePrice(s *pricingDomain.ServicePriceModel)
	},
) {
	// /admin/pricing/routes
	r.GET("/admin/pricing/routes", a(func(w http.ResponseWriter, req *http.Request) {
		prices := pmr.ListRoutePrices()
		rows := make([][]string, len(prices))
		for i, p := range prices {
			rows[i] = []string{fmt.Sprintf("%d", p.ID), p.ClientName, p.RouteName, p.TransportType, p.CargoType, p.TaxMode,
				fmt.Sprintf("¥%.2f/kg", p.WeightPrice), fmt.Sprintf("¥%.2f/才", p.VolumePrice),
				fmt.Sprintf("¥%.0f起", p.MinCharge), common.StatusLabelText(p.IsActive)}
		}
		gp(w, "pricing_routes", "客户×线路价", len(rows),
			[]string{"ID", "客户", "线路", "运输方式", "货类", "税档", "重量单价", "体积单价", "最低收费", "状态"}, rows,
			"/admin/pricing/routes/add-form")
	}))
	r.GET("/admin/pricing/routes/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增线路价格")+common.FormSave("/admin/pricing/routes/save")+
			common.FormField("客户名称", "client_name", "EZ集运通", "")+
			common.FormField("线路名", "route_name", "", "如: 深圳→台湾(空运)")+
			common.FormSelect("运输方式", "transport_type", "air",
				[2]string{"air", "空运"}, [2]string{"sea_express", "海快"}, [2]string{"sea", "海运"}, [2]string{"air_special", "空运特货"})+
			common.FormField("货类", "cargo_type", "", "普货/家具类/一类~六类/易碎品")+
			common.FormSelect("税档", "tax_mode", "全包税",
				[2]string{"全包税", "全包税"}, [2]string{"频税", "频税"}, [2]string{"不包税", "不包税"})+
			common.FormField("重量单价(¥/kg)", "weight_price", "", "如: 20.00")+
			common.FormField("体积单价(¥/才)", "volume_price", "", "如: 20.00")+
			common.FormField("最低收费(¥)", "min_charge", "50", "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/pricing/routes/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		wp, _ := common.ParseFloat(req.FormValue("weight_price"))
		vp, _ := common.ParseFloat(req.FormValue("volume_price"))
		mc, _ := common.ParseFloat(req.FormValue("min_charge"))
		pmr.AddRoutePrice(&pricingDomain.RoutePriceModel{
			ClientName: req.FormValue("client_name"), RouteName: req.FormValue("route_name"),
			TransportType: req.FormValue("transport_type"), CargoType: req.FormValue("cargo_type"),
			TaxMode: req.FormValue("tax_mode"), WeightPrice: wp, VolumePrice: vp, MinCharge: mc,
		})
		common.Redirect(w, "/admin/pricing/routes")
	}))
	r.GET("/admin/pricing/routes/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		id, _ := common.ParseID(req.URL.Query().Get("id"))
		if id == 0 { http.Error(w, "invalid id", 400); return }
		p := pmr.GetRoutePriceByID(id)
		if p == nil { http.Error(w, "not found", 404); return }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑线路价格")+common.FormSave("/admin/pricing/routes/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, p.ID)+
			common.FormField("客户名称", "client_name", p.ClientName, "")+
			common.FormField("线路名", "route_name", p.RouteName, "如: 深圳→台湾(空运)")+
			common.FormSelect("运输方式", "transport_type", p.TransportType,
				[2]string{"air", "空运"}, [2]string{"sea_express", "海快"}, [2]string{"sea", "海运"}, [2]string{"air_special", "空运特货"})+
			common.FormField("货类", "cargo_type", p.CargoType, "")+
			common.FormSelect("税档", "tax_mode", p.TaxMode,
				[2]string{"全包税", "全包税"}, [2]string{"频税", "频税"}, [2]string{"不包税", "不包税"})+
			common.FormField("重量单价(¥/kg)", "weight_price", fmt.Sprintf("%.2f", p.WeightPrice), "")+
			common.FormField("体积单价(¥/才)", "volume_price", fmt.Sprintf("%.2f", p.VolumePrice), "")+
			common.FormField("最低收费(¥)", "min_charge", fmt.Sprintf("%.0f", p.MinCharge), "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/pricing/routes/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := common.ParseID(req.FormValue("id"))
		wp, _ := common.ParseFloat(req.FormValue("weight_price"))
		vp, _ := common.ParseFloat(req.FormValue("volume_price"))
		mc, _ := common.ParseFloat(req.FormValue("min_charge"))
		pmr.UpdateRoutePrice(id, &pricingDomain.RoutePriceModel{
			ClientName: req.FormValue("client_name"), RouteName: req.FormValue("route_name"),
			TransportType: req.FormValue("transport_type"), CargoType: req.FormValue("cargo_type"),
			TaxMode: req.FormValue("tax_mode"), WeightPrice: wp, VolumePrice: vp, MinCharge: mc,
			IsActive: true,
		})
		common.Redirect(w, "/admin/pricing/routes")
	}))

	// /admin/pricing/delivery
	r.GET("/admin/pricing/delivery", a(func(w http.ResponseWriter, req *http.Request) {
		fees := pmr.ListDeliveryFees()
		rows := make([][]string, len(fees))
		for i, f := range fees {
			freeLabel := "—"; if f.FreeThresholdTxt != "" { freeLabel = f.FreeThresholdTxt }
			rows[i] = []string{fmt.Sprintf("%d", f.ID), f.ClientName, f.CarrierName, f.CustomsPoint, f.Area,
				f.DeliveryMethod, f.Condition, fmt.Sprintf("¥%.0f", f.Fee), freeLabel,
				common.StatusLabelText(f.IsActive)}
		}
		gp(w, "pricing_delivery", "客户×派送费", len(rows),
			[]string{"ID", "客户", "承运商", "清關點", "區域", "派送方式", "条件", "费用", "免运门槛", "状态"}, rows,
			"/admin/pricing/delivery/add-form")
	}))
	r.GET("/admin/pricing/delivery/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增派送费")+common.FormSave("/admin/pricing/delivery/save")+
			common.FormField("客户名称", "client_name", "EZ集运通", "")+
			common.FormField("承运商", "carrier_name", "", "新竹物流/黑猫宅急便/顺丰速运")+
			common.FormSelect("清關點", "customs_point", "台北", [2]string{"台北", "台北"}, [2]string{"台中", "台中"}, [2]string{"高雄", "高雄"})+
			common.FormSelect("區域", "area", "預設", [2]string{"預設", "預設"}, [2]string{"北部", "北部"}, [2]string{"中部", "中部"}, [2]string{"南部", "南部"}, [2]string{"东部", "东部"})+
			common.FormSelect("派送方式", "delivery_method", "宅配", [2]string{"宅配", "宅配"}, [2]string{"專車", "專車"}, [2]string{"自取", "自取"})+
			common.FormField("条件", "condition", "", "如: 重量>39.8")+
			common.FormField("费用(¥)", "fee", "", "如: 20")+
			common.FormField("免运门槛", "free_threshold_txt", "", "如: ≥10kg免运")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/pricing/delivery/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		fee, _ := common.ParseFloat(req.FormValue("fee"))
		pmr.AddDeliveryFee(&pricingDomain.DeliveryFeeModel{
			ClientName: req.FormValue("client_name"), CarrierName: req.FormValue("carrier_name"),
			CustomsPoint: req.FormValue("customs_point"), Area: req.FormValue("area"),
			DeliveryMethod: req.FormValue("delivery_method"), Condition: req.FormValue("condition"),
			Fee: fee, FreeThresholdTxt: req.FormValue("free_threshold_txt"),
		})
		common.Redirect(w, "/admin/pricing/delivery")
	}))
	r.GET("/admin/pricing/delivery/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		id, _ := common.ParseID(req.URL.Query().Get("id"))
		if id == 0 { http.Error(w, "invalid id", 400); return }
		f := pmr.GetDeliveryFeeByID(id)
		if f == nil { f = &pricingDomain.DeliveryFeeModel{ID: id, ClientName:"", CarrierName:"", CustomsPoint:"台北", Area:"預設", DeliveryMethod:"宅配", Condition:"", Fee: 0, FreeThresholdTxt:""} }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑派送费")+common.FormSave("/admin/pricing/delivery/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, f.ID)+
			common.FormField("客户名称", "client_name", f.ClientName, "")+
			common.FormField("承运商", "carrier_name", f.CarrierName, "")+
			common.FormSelect("清關點", "customs_point", f.CustomsPoint, [2]string{"台北", "台北"}, [2]string{"台中", "台中"}, [2]string{"高雄", "高雄"})+
			common.FormSelect("區域", "area", f.Area, [2]string{"預設", "預設"}, [2]string{"北部", "北部"}, [2]string{"中部", "中部"}, [2]string{"南部", "南部"}, [2]string{"东部", "东部"})+
			common.FormSelect("派送方式", "delivery_method", f.DeliveryMethod, [2]string{"宅配", "宅配"}, [2]string{"專車", "專車"}, [2]string{"自取", "自取"})+
			common.FormField("条件", "condition", f.Condition, "")+
			common.FormField("费用(¥)", "fee", fmt.Sprintf("%.0f", f.Fee), "")+
			common.FormField("免运门槛", "free_threshold_txt", f.FreeThresholdTxt, "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/pricing/delivery/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := common.ParseID(req.FormValue("id"))
		fee, _ := common.ParseFloat(req.FormValue("fee"))
		pmr.UpdateDeliveryFee(id, &pricingDomain.DeliveryFeeModel{
			ClientName: req.FormValue("client_name"), CarrierName: req.FormValue("carrier_name"),
			CustomsPoint: req.FormValue("customs_point"), Area: req.FormValue("area"),
			DeliveryMethod: req.FormValue("delivery_method"), Condition: req.FormValue("condition"),
			Fee: fee, FreeThresholdTxt: req.FormValue("free_threshold_txt"),
			IsActive: true,
		})
		common.Redirect(w, "/admin/pricing/delivery")
	}))

	// /admin/pricing/surcharges
	r.GET("/admin/pricing/surcharges", a(func(w http.ResponseWriter, req *http.Request) {
		surcharges := pmr.ListSurcharges()
		rows := make([][]string, len(surcharges))
		for i, s := range surcharges {
			priceLabel := fmt.Sprintf("¥%.0f", s.Price); if s.PriceDesc != "" { priceLabel = s.PriceDesc }
			rows[i] = []string{fmt.Sprintf("%d", s.ID), s.ClientName, s.CarrierName, s.ChargeType, s.Tier, s.CustomsPoint, s.Area, s.Condition, priceLabel, common.StatusLabelText(s.IsActive)}
		}
		gp(w, "pricing_surcharges", "客户×加收费", len(rows),
			[]string{"ID", "客户", "承运商", "加收类型", "档位", "清關點", "區域", "触发条件", "费用", "状态"}, rows,
			"/admin/pricing/surcharges/add-form")
	}))
	r.GET("/admin/pricing/surcharges/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增加收费")+common.FormSave("/admin/pricing/surcharges/save")+
			common.FormField("客户名称", "client_name", "EZ集运通", "")+
			common.FormField("承运商", "carrier_name", "", "新竹物流/黑猫宅急便")+
			common.FormSelect("加收类型", "charge_type", "超長費",
				[2]string{"超長費", "超長費"}, [2]string{"超材費", "超材費"}, [2]string{"棧板費", "棧板費"}, [2]string{"偏遠費", "偏遠費"}, [2]string{"上樓費", "上樓費"})+
			common.FormSelect("档位", "tier", "—", [2]string{"—", "—"}, [2]string{"小板", "小板"}, [2]string{"大板", "大板"})+
			common.FormSelect("清關點", "customs_point", "台北", [2]string{"台北", "台北"}, [2]string{"台中", "台中"}, [2]string{"高雄", "高雄"})+
			common.FormField("區域", "area", "預設", "")+
			common.FormField("触发条件", "condition", "", "如: 單邊>150cm")+
			common.FormField("费用(¥)", "price", "", "如: 100")+
			common.FormField("费用说明", "price_desc", "", "如: 每件加收")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/pricing/surcharges/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		price, _ := common.ParseFloat(req.FormValue("price"))
		pmr.AddSurcharge(&pricingDomain.SurchargeModel{
			ClientName: req.FormValue("client_name"), CarrierName: req.FormValue("carrier_name"),
			ChargeType: req.FormValue("charge_type"), Tier: req.FormValue("tier"),
			CustomsPoint: req.FormValue("customs_point"), Area: req.FormValue("area"),
			Condition: req.FormValue("condition"), Price: price, PriceDesc: req.FormValue("price_desc"),
		})
		common.Redirect(w, "/admin/pricing/surcharges")
	}))
	r.GET("/admin/pricing/surcharges/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		id, _ := common.ParseID(req.URL.Query().Get("id"))
		if id == 0 { http.Error(w, "invalid id", 400); return }
		s := pmr.GetSurchargeByID(id)
		if s == nil { s = &pricingDomain.SurchargeModel{ID: id, ChargeType:"", Price: 0, Condition:""} }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑加收费")+common.FormSave("/admin/pricing/surcharges/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, s.ID)+
			common.FormField("客户名称", "client_name", s.ClientName, "")+
			common.FormField("承运商", "carrier_name", s.CarrierName, "")+
			common.FormSelect("加收类型", "charge_type", s.ChargeType,
				[2]string{"超長費", "超長費"}, [2]string{"超材費", "超材費"}, [2]string{"棧板費", "棧板費"}, [2]string{"偏遠費", "偏遠費"}, [2]string{"上樓費", "上樓費"})+
			common.FormSelect("档位", "tier", s.Tier, [2]string{"—", "—"}, [2]string{"小板", "小板"}, [2]string{"大板", "大板"})+
			common.FormSelect("清關點", "customs_point", s.CustomsPoint, [2]string{"台北", "台北"}, [2]string{"台中", "台中"}, [2]string{"高雄", "高雄"})+
			common.FormField("區域", "area", s.Area, "")+
			common.FormField("触发条件", "condition", s.Condition, "")+
			common.FormField("费用(¥)", "price", fmt.Sprintf("%.0f", s.Price), "")+
			common.FormField("费用说明", "price_desc", s.PriceDesc, "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/pricing/surcharges/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := common.ParseID(req.FormValue("id"))
		price, _ := common.ParseFloat(req.FormValue("price"))
		pmr.UpdateSurcharge(id, &pricingDomain.SurchargeModel{
			ClientName: req.FormValue("client_name"), CarrierName: req.FormValue("carrier_name"),
			ChargeType: req.FormValue("charge_type"), Tier: req.FormValue("tier"),
			CustomsPoint: req.FormValue("customs_point"), Area: req.FormValue("area"),
			Condition: req.FormValue("condition"), Price: price, PriceDesc: req.FormValue("price_desc"),
			IsActive: true,
		})
		common.Redirect(w, "/admin/pricing/surcharges")
	}))

	// /admin/pricing/services
	r.GET("/admin/pricing/services", a(func(w http.ResponseWriter, req *http.Request) {
		services := pmr.ListServicePrices()
		rows := make([][]string, len(services))
		for i, s := range services {
			rows[i] = []string{s.ClientName, s.ServiceType, s.ServiceCode,
				fmt.Sprintf("¥%.2f", s.UnitPrice), s.PriceMode, common.StatusLabelText(s.IsActive)}
		}
		gp(w, "pricing_services", "客户×附加服务", len(rows),
			[]string{"客户", "服务类型", "服务编码", "单价", "计费模式", "状态"}, rows,
			"/admin/pricing/services/add-form")
	}))
	r.GET("/admin/pricing/services/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增服务价格")+common.FormSave("/admin/pricing/services/save")+
			common.FormField("客户名称", "client_name", "EZ集运通", "")+
			common.FormField("服务类型", "service_type", "", "如: 木箱包装")+
			common.FormField("服务编码", "service_code", "", "如: WOODEN_CRATE")+
			common.FormField("单价(¥)", "unit_price", "", "如: 80.00")+
			common.FormSelect("计费模式", "price_mode", "per_item",
				[2]string{"fixed", "固定"}, [2]string{"per_item", "按件"}, [2]string{"per_kg", "按重量"}, [2]string{"per_order", "按单"})+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/pricing/services/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		price, _ := common.ParseFloat(req.FormValue("unit_price"))
		pmr.AddServicePrice(&pricingDomain.ServicePriceModel{
			ClientName: req.FormValue("client_name"), ServiceType: req.FormValue("service_type"),
			ServiceCode: req.FormValue("service_code"), UnitPrice: price, PriceMode: req.FormValue("price_mode"),
		})
		common.Redirect(w, "/admin/pricing/services")
	}))
}
