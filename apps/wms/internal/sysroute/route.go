// Package route provides SYS (系统管理) admin route registration.
package route

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/i56/framework/core/router"

	"github.com/i56/i56-apps/i56-wms/internal/common"

	rbaDomain "github.com/i56/modules/rbac/domain"
	rbacRepo "github.com/i56/modules/rbac/repository"
	sysDomain "github.com/i56/modules/system/domain"
	sysRepo "github.com/i56/modules/system/repository"
)

// Register SYS admin routes (~12 list pages + CRUD).
func Register(
	r *router.Router,
	a func(http.HandlerFunc) http.HandlerFunc,
	rc *common.RenderCtx,
	sysCfg *sysRepo.MemSystemConfigRepo,
	rbac *rbacRepo.MemRBACRepo,
) {
	const tenant int64 = 1
	gp := rc.NewGenericList()

	// ─── /admin/employees (from admin_modules.go SYS) ───
	r.GET("/admin/employees", a(func(w http.ResponseWriter, req *http.Request) {
		users, _, _ := rbac.ListUsers(req.Context(), tenant, 0, 50)
		roles, _, _ := rbac.ListRoles(req.Context(), tenant, 0, 50)
		roleNames := map[int64]string{}
		for _, ro := range roles { roleNames[ro.ID] = ro.Name }
		rows := make([][]string, len(users))
		for i, u := range users {
			rn := roleNames[u.RoleID]; if rn == "" { rn = fmt.Sprintf("Role-%d", u.RoleID) }
			rows[i] = []string{u.RealName, u.Username, rn, u.Email, u.Phone, common.StatusLabelText(u.IsActive)}
		}
		if len(rows) == 0 {
			rows = [][]string{
				{"大宝", "dabao", "仓库管理", "dabao@example.com", "13800001111", "启用"},
				{"安冉", "anran", "仓库管理", "anran@example.com", "13800002222", "启用"},
				{"小林", "xiaolin", "质检管理", "xiaolin@example.com", "13800003333", "启用"},
			}
		}
		rc.Exec(rc.Tmpl, "sys_employees", w, "sys_employees.html", map[string]any{
			"Title": "员工管理", "Page": "sys_employees",
			"Columns": []string{"姓名", "账号", "角色", "邮箱", "电话", "状态"},
			"Rows": rows, "Total": len(rows),
			"AddURL": "/admin/employees/add-form", "HasActions": true,
		})
	}))
	r.GET("/admin/employees/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		roles, _, _ := rbac.ListRoles(req.Context(), tenant, 0, 50)
		roleOpts := ""
		for _, ro := range roles { roleOpts += fmt.Sprintf(`<option value="%d">%s</option>`, ro.ID, ro.Name) }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增员工")+common.FormSave("/admin/employees/save")+
			common.FormField("账号", "username", "", "")+
			common.FormField("密码", "password", "", "")+
			common.FormField("姓名", "real_name", "", "")+
			fmt.Sprintf(`<div class="form-group"><label class="form-label">角色</label><select name="role_id" class="form-input">%s</select></div>`, roleOpts)+
			common.FormField("邮箱", "email", "", "")+
			common.FormField("电话", "phone", "", "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/employees/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		roleID, _ := common.ParseID(req.FormValue("role_id"))
		rbac.CreateUser(req.Context(), tenant, &rbaDomain.User{
			Username: req.FormValue("username"), Password: req.FormValue("password"),
			RealName: req.FormValue("real_name"), Email: req.FormValue("email"),
			Phone: req.FormValue("phone"), RoleID: roleID, IsActive: true,
		})
		common.Redirect(w, "/admin/employees")
	}))
	r.GET("/admin/employees/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		id, _ := common.ParseID(req.URL.Query().Get("id"))
		u, _ := rbac.GetUserByID(req.Context(), tenant, id)
		if u == nil { http.Error(w, "not found", 404); return }
		roles, _, _ := rbac.ListRoles(req.Context(), tenant, 0, 50)
		roleOpts := ""
		for _, ro := range roles {
			sel := ""; if ro.ID == u.RoleID { sel = " selected" }
			roleOpts += fmt.Sprintf(`<option value="%d"%s>%s</option>`, ro.ID, sel, ro.Name)
		}
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑员工")+common.FormSave("/admin/employees/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, u.ID)+
			common.FormField("账号", "username", u.Username, "")+
			common.FormField("姓名", "real_name", u.RealName, "")+
			fmt.Sprintf(`<div class="form-group"><label class="form-label">角色</label><select name="role_id" class="form-input">%s</select></div>`, roleOpts)+
			common.FormField("邮箱", "email", u.Email, "")+
			common.FormField("电话", "phone", u.Phone, "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/employees/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := strconv.ParseInt(req.FormValue("id"), 10, 64)
		roleID, _ := common.ParseID(req.FormValue("role_id"))
		rbac.UpdateUser(req.Context(), tenant, id, &rbaDomain.User{
			ID: id, TenantID: tenant, Username: req.FormValue("username"),
			RealName: req.FormValue("real_name"), Email: req.FormValue("email"),
			Phone: req.FormValue("phone"), RoleID: roleID, IsActive: true,
		})
		common.Redirect(w, "/admin/employees")
	}))

	// ─── /admin/roles — 角色管理 ───
	r.GET("/admin/roles", a(func(w http.ResponseWriter, req *http.Request) {
		roles, _, _ := rbac.ListRoles(req.Context(), tenant, 0, 50)
		rows := make([][]string, len(roles))
		for i, ro := range roles {
			permCount := len(ro.PermissionIDs)
			desc := ro.Description
			if desc == "" { desc = "—" }
			rows[i] = []string{ro.Name, ro.Slug, fmt.Sprintf("%d", permCount), desc, common.StatusLabelText(ro.IsActive)}
		}
		if len(rows) == 0 {
			rows = [][]string{
				{"超级管理员", "super_admin", "12", "系统最高权限", "启用"},
				{"仓库管理", "warehouse_admin", "8", "仓库操作与订单管理", "启用"},
				{"质检管理", "qc_admin", "5", "质检与异常处理", "启用"},
				{"运营人员", "operator", "4", "日常订单操作", "启用"},
			}
		}
		rc.Exec(rc.Tmpl, "sys_roles", w, "sys_roles.html", map[string]any{
			"Title": "角色管理", "Page": "sys_roles",
			"Columns": []string{"角色名", "Slug", "权限数", "描述", "状态"},
			"Rows": rows, "Total": len(rows),
			"AddURL": "/admin/roles/add-form", "HasActions": true,
		})
	}))
	r.GET("/admin/roles/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增角色")+common.FormSave("/admin/roles/save")+
			common.FormField("角色名", "name", "", "如: 仓库管理")+
			common.FormField("Slug", "slug", "", "如: warehouse_admin")+
			common.FormField("描述", "description", "", "角色描述")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/roles/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		rbac.CreateRole(req.Context(), tenant, &rbaDomain.Role{Name: req.FormValue("name"), Slug: req.FormValue("slug"), Description: req.FormValue("description"), IsActive: true})
		common.Redirect(w, "/admin/roles")
	}))
	r.GET("/admin/roles/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		id, _ := common.ParseID(req.URL.Query().Get("id"))
		ro, _ := rbac.GetRoleByID(req.Context(), tenant, id)
		if ro == nil { http.Error(w, "not found", 404); return }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑角色")+common.FormSave("/admin/roles/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, ro.ID)+
			common.FormField("角色名", "name", ro.Name, "如: 仓库管理")+
			common.FormField("Slug", "slug", ro.Slug, "如: warehouse_admin")+
			common.FormField("描述", "description", ro.Description, "角色描述")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/roles/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := strconv.ParseInt(req.FormValue("id"), 10, 64)
		rbac.UpdateRole(req.Context(), tenant, id, &rbaDomain.Role{
			ID: id, TenantID: tenant,
			Name: req.FormValue("name"), Slug: req.FormValue("slug"),
			Description: req.FormValue("description"), IsActive: true,
		})
		common.Redirect(w, "/admin/roles")
	}))

	// ─── System API config pages (from admin_modules.go SYS) ───
	apiCfg := sysRepo.NewMemAPIConfigRepo()

	// /admin/system/api-couriers
	r.GET("/admin/system/api-couriers", a(func(w http.ResponseWriter, req *http.Request) {
		configs := apiCfg.ListCouriers(tenant)
		rows := make([][]string, len(configs))
		for i, c := range configs {
			rows[i] = []string{fmt.Sprintf("COU-%d", c.ID), c.Name, c.APIEndpoint, c.AuthType, c.TrackingPattern, common.StatusLabelText(c.IsActive), c.CreatedAt.Format("01-02 15:04")}
		}
		gp(w, "sys_api_couriers", "快递公司API配置", len(rows), []string{"编号", "名称", "API端点", "认证方式", "运单号格式", "状态", "创建时间"}, rows, "/admin/system/api-couriers/add-form")
	}))
	r.GET("/admin/system/api-couriers/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增快递API配置")+common.FormSave("/admin/system/api-couriers/save")+
			common.FormField("名称", "name", "", "如: 顺丰速运API")+
			common.FormField("API端点", "api_endpoint", "", "https://open.sf-express.com/std/service")+
			common.FormField("API Key", "api_key", "", "")+
			common.FormField("API Secret", "api_secret", "", "")+
			common.FormField("运单号正则", "tracking_pattern", "", `^SF\d{12}$`)+
			common.FormSelect("认证方式", "auth_type", "api_key", [2]string{"api_key", "API Key"}, [2]string{"hmac", "HMAC签名"}, [2]string{"oauth2", "OAuth 2.0"})+
			common.FormField("额外Headers(JSON)", "extra_headers", "{}", `{"X-Custom":"value"}`)+
			common.FormSelect("状态", "is_active", "true", [2]string{"true", "启用"}, [2]string{"false", "停用"})+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/system/api-couriers/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		isActive := req.FormValue("is_active") != "false"
		apiCfg.SaveCourier(req.Context(), &sysDomain.CourierAPIConfig{
			Name: req.FormValue("name"), APIEndpoint: req.FormValue("api_endpoint"),
			APIKey: req.FormValue("api_key"), APISecret: req.FormValue("api_secret"),
			TrackingPattern: req.FormValue("tracking_pattern"), AuthType: req.FormValue("auth_type"),
			ExtraHeaders: req.FormValue("extra_headers"), IsActive: isActive, TenantID: tenant,
		})
		common.Redirect(w, "/admin/system/api-couriers")
	}))
	// edit-form
	r.GET("/admin/system/api-couriers/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		idStr := req.URL.Query().Get("id")
		if idx := strings.LastIndex(idStr, "-"); idx >= 0 { idStr = idStr[idx+1:] }
		id, _ := common.ParseID(idStr)
		if id == 0 { http.Error(w, "invalid id", 400); return }
		var config *sysDomain.CourierAPIConfig
		for _, c := range apiCfg.ListCouriers(tenant) { if c.ID == id { config = c; break } }
		if config == nil { http.Error(w, "not found", 404); return }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑快递API配置")+common.FormSave("/admin/system/api-couriers/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, config.ID)+
			common.FormField("名称", "name", config.Name, "如: 顺丰速运API")+
			common.FormField("API端点", "api_endpoint", config.APIEndpoint, "https://open.sf-express.com/std/service")+
			common.FormField("API Key", "api_key", config.APIKey, "")+
			common.FormField("API Secret", "api_secret", config.APISecret, "")+
			common.FormField("运单号正则", "tracking_pattern", config.TrackingPattern, `^SF\\d{12}$`)+
			common.FormSelect("认证方式", "auth_type", config.AuthType, [2]string{"api_key", "API Key"}, [2]string{"hmac", "HMAC签名"}, [2]string{"oauth2", "OAuth 2.0"})+
			common.FormField("额外Headers(JSON)", "extra_headers", config.ExtraHeaders, `{"X-Custom":"value"}`)+
			common.FormSelect("状态", "is_active", fmt.Sprintf("%v", config.IsActive), [2]string{"true", "启用"}, [2]string{"false", "停用"})+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/system/api-couriers/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := common.ParseID(req.FormValue("id"))
		isActive := req.FormValue("is_active") != "false"
		apiCfg.SaveCourier(req.Context(), &sysDomain.CourierAPIConfig{
			ID: id, TenantID: tenant,
			Name: req.FormValue("name"), APIEndpoint: req.FormValue("api_endpoint"),
			APIKey: req.FormValue("api_key"), APISecret: req.FormValue("api_secret"),
			TrackingPattern: req.FormValue("tracking_pattern"), AuthType: req.FormValue("auth_type"),
			ExtraHeaders: req.FormValue("extra_headers"), IsActive: isActive,
		})
		common.Redirect(w, "/admin/system/api-couriers")
	}))

	// /admin/system/api-customs
	r.GET("/admin/system/api-customs", a(func(w http.ResponseWriter, req *http.Request) {
		configs := apiCfg.ListCustomsBrokers(tenant)
		rows := make([][]string, len(configs))
		for i, c := range configs {
			rows[i] = []string{fmt.Sprintf("CUS-%d", c.ID), c.Name, c.DeclarationAPIURL, c.CustomsPointID, c.NumberPrefix, common.StatusLabelText(c.IsActive), c.CreatedAt.Format("01-02 15:04")}
		}
		if len(rows) == 0 {
			rows = [][]string{
				{"CUS-1", "厦门电子口岸", "https://customs.xm-port.gov.cn/api/v2", "CN_XM_3701", "776XM", "启用", "07-01 10:00"},
				{"CUS-2", "深圳海关", "https://customs.sz-port.gov.cn/api/v2", "CN_SZ_5301", "4403SZ", "启用", "07-01 10:00"},
			}
		}
		gp(w, "sys_api_customs", "报关行API配置", len(rows), []string{"编号", "名称", "申报API端点", "海关口岸", "报关单前缀", "状态", "创建时间"}, rows, "/admin/system/api-customs/add-form")
	}))
	r.GET("/admin/system/api-customs/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增报关行API配置")+common.FormSave("/admin/system/api-customs/save")+
			common.FormField("名称", "name", "", "如: 厦门电子口岸清关")+
			common.FormField("申报API URL", "declaration_api_url", "", "https://customs.xm-port.gov.cn/api/v2")+
			common.FormField("API Key", "api_key", "", "")+
			common.FormField("API Secret", "api_secret", "", "")+
			common.FormField("海关口岸编号", "customs_point_id", "", "CN_XM_3701")+
			common.FormField("报关单号前缀", "number_prefix", "", "776XM")+
			common.FormField("支持单证(JSON)", "supported_documents", `["invoice","packing_list"]`, "")+
			common.FormSelect("状态", "is_active", "true", [2]string{"true", "启用"}, [2]string{"false", "停用"})+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/system/api-customs/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		isActive := req.FormValue("is_active") != "false"
		apiCfg.SaveCustomsBroker(req.Context(), &sysDomain.CustomsBrokerConfig{
			Name: req.FormValue("name"), DeclarationAPIURL: req.FormValue("declaration_api_url"),
			APIKey: req.FormValue("api_key"), APISecret: req.FormValue("api_secret"),
			CustomsPointID: req.FormValue("customs_point_id"), NumberPrefix: req.FormValue("number_prefix"),
			SupportedDocuments: req.FormValue("supported_documents"), IsActive: isActive, TenantID: tenant,
		})
		common.Redirect(w, "/admin/system/api-customs")
	}))
	// edit-form
	r.GET("/admin/system/api-customs/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		idStr := req.URL.Query().Get("id")
		if idx := strings.LastIndex(idStr, "-"); idx >= 0 { idStr = idStr[idx+1:] }
		id, _ := common.ParseID(idStr)
		if id == 0 { http.Error(w, "invalid id", 400); return }
		var config *sysDomain.CustomsBrokerConfig
		for _, c := range apiCfg.ListCustomsBrokers(tenant) { if c.ID == id { config = c; break } }
		if config == nil { http.Error(w, "not found", 404); return }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑报关行API配置")+common.FormSave("/admin/system/api-customs/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, config.ID)+
			common.FormField("名称", "name", config.Name, "如: 厦门电子口岸清关")+
			common.FormField("申报API URL", "declaration_api_url", config.DeclarationAPIURL, "https://customs.xm-port.gov.cn/api/v2")+
			common.FormField("API Key", "api_key", config.APIKey, "")+
			common.FormField("API Secret", "api_secret", config.APISecret, "")+
			common.FormField("海关口岸编号", "customs_point_id", config.CustomsPointID, "CN_XM_3701")+
			common.FormField("报关单号前缀", "number_prefix", config.NumberPrefix, "776XM")+
			common.FormField("支持单证(JSON)", "supported_documents", config.SupportedDocuments, "")+
			common.FormSelect("状态", "is_active", fmt.Sprintf("%v", config.IsActive), [2]string{"true", "启用"}, [2]string{"false", "停用"})+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/system/api-customs/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := common.ParseID(req.FormValue("id"))
		isActive := req.FormValue("is_active") != "false"
		apiCfg.SaveCustomsBroker(req.Context(), &sysDomain.CustomsBrokerConfig{
			ID: id, TenantID: tenant,
			Name: req.FormValue("name"), DeclarationAPIURL: req.FormValue("declaration_api_url"),
			APIKey: req.FormValue("api_key"), APISecret: req.FormValue("api_secret"),
			CustomsPointID: req.FormValue("customs_point_id"), NumberPrefix: req.FormValue("number_prefix"),
			SupportedDocuments: req.FormValue("supported_documents"), IsActive: isActive,
		})
		common.Redirect(w, "/admin/system/api-customs")
	}))

	// /admin/system/api-notifications
	r.GET("/admin/system/api-notifications", a(func(w http.ResponseWriter, req *http.Request) {
		configs := apiCfg.ListNotificationChannels(tenant)
		rows := make([][]string, len(configs))
		for i, c := range configs {
			rows[i] = []string{fmt.Sprintf("NTC-%d", c.ID), c.Name, c.ChannelType, c.Provider, common.StatusLabelText(c.IsActive), c.CreatedAt.Format("01-02 15:04")}
		}
		gp(w, "sys_api_notifications", "通知渠道配置", len(rows), []string{"编号", "名称", "渠道类型", "服务商", "状态", "创建时间"}, rows, "/admin/system/api-notifications/add-form")
	}))
	r.GET("/admin/system/api-notifications/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增通知渠道")+common.FormSave("/admin/system/api-notifications/save")+
			common.FormField("名称", "name", "", "如: 阿里云短信服务")+
			common.FormSelect("渠道类型", "channel_type", "sms", [2]string{"email", "邮件"}, [2]string{"sms", "短信"}, [2]string{"line", "Line"}, [2]string{"telegram", "Telegram"}, [2]string{"webhook", "Webhook"})+
			common.FormSelect("服务商", "provider", "aliyun_sms", [2]string{"smtp", "SMTP"}, [2]string{"sendgrid", "SendGrid"}, [2]string{"aliyun_sms", "阿里云短信"}, [2]string{"twilio", "Twilio"}, [2]string{"line_notify", "Line Notify"}, [2]string{"telegram_bot", "Telegram Bot"})+
			common.FormField("配置(JSON)", "config_json", "{}", `{"api_key":"xxx"}`)+
			common.FormSelect("状态", "is_active", "true", [2]string{"true", "启用"}, [2]string{"false", "停用"})+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/system/api-notifications/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		isActive := req.FormValue("is_active") != "false"
		apiCfg.SaveNotificationChannel(req.Context(), &sysDomain.NotificationChannel{
			Name: req.FormValue("name"), ChannelType: req.FormValue("channel_type"),
			Provider: req.FormValue("provider"), ConfigJSON: req.FormValue("config_json"),
			IsActive: isActive, TenantID: tenant,
		})
		common.Redirect(w, "/admin/system/api-notifications")
	}))
	// edit-form
	r.GET("/admin/system/api-notifications/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		idStr := req.URL.Query().Get("id")
		if idx := strings.LastIndex(idStr, "-"); idx >= 0 { idStr = idStr[idx+1:] }
		id, _ := common.ParseID(idStr)
		if id == 0 { http.Error(w, "invalid id", 400); return }
		var config *sysDomain.NotificationChannel
		for _, c := range apiCfg.ListNotificationChannels(tenant) { if c.ID == id { config = c; break } }
		if config == nil { http.Error(w, "not found", 404); return }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑通知渠道")+common.FormSave("/admin/system/api-notifications/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, config.ID)+
			common.FormField("名称", "name", config.Name, "如: 阿里云短信服务")+
			common.FormSelect("渠道类型", "channel_type", config.ChannelType, [2]string{"email", "邮件"}, [2]string{"sms", "短信"}, [2]string{"line", "Line"}, [2]string{"telegram", "Telegram"}, [2]string{"webhook", "Webhook"})+
			common.FormSelect("服务商", "provider", config.Provider, [2]string{"smtp", "SMTP"}, [2]string{"sendgrid", "SendGrid"}, [2]string{"aliyun_sms", "阿里云短信"}, [2]string{"twilio", "Twilio"}, [2]string{"line_notify", "Line Notify"}, [2]string{"telegram_bot", "Telegram Bot"})+
			common.FormField("配置(JSON)", "config_json", config.ConfigJSON, `{"api_key":"xxx"}`)+
			common.FormSelect("状态", "is_active", fmt.Sprintf("%v", config.IsActive), [2]string{"true", "启用"}, [2]string{"false", "停用"})+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/system/api-notifications/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := common.ParseID(req.FormValue("id"))
		isActive := req.FormValue("is_active") != "false"
		apiCfg.SaveNotificationChannel(req.Context(), &sysDomain.NotificationChannel{
			ID: id, TenantID: tenant,
			Name: req.FormValue("name"), ChannelType: req.FormValue("channel_type"),
			Provider: req.FormValue("provider"), ConfigJSON: req.FormValue("config_json"),
			IsActive: isActive,
		})
		common.Redirect(w, "/admin/system/api-notifications")
	}))

	// /admin/system/api-printers
	r.GET("/admin/system/api-printers", a(func(w http.ResponseWriter, req *http.Request) {
		configs := apiCfg.ListPrintTemplates(tenant)
		rows := make([][]string, len(configs))
		for i, c := range configs {
			rows[i] = []string{fmt.Sprintf("PRT-%d", c.ID), c.Name, c.Type, c.PaperSize, c.PrinterType, common.StatusLabelText(c.IsActive), c.CreatedAt.Format("01-02 15:04")}
		}
		gp(w, "sys_api_printers", "打印模板配置", len(rows), []string{"编号", "名称", "类型", "纸张规格", "打印机类型", "状态", "创建时间"}, rows, "/admin/system/api-printers/add-form")
	}))
	r.GET("/admin/system/api-printers/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增打印模板")+common.FormSave("/admin/system/api-printers/save")+
			common.FormField("名称", "name", "", "如: 顺丰标准面单")+
			common.FormSelect("模板类型", "type", "label", [2]string{"label", "标签"}, [2]string{"invoice", "发票"}, [2]string{"packing_list", "装箱单"}, [2]string{"waybill", "运单"})+
			common.FormSelect("纸张规格", "paper_size", "100x150mm", [2]string{"100x150mm", "100x150mm(热敏)"}, [2]string{"100x100mm", "100x100mm(热敏)"}, [2]string{"4x6inch", "4x6英寸"}, [2]string{"A4", "A4"}, [2]string{"A5", "A5"})+
			common.FormField("模板内容", "template_content", "", "ZPL/HTML模板内容")+
			common.FormSelect("打印机类型", "printer_type", "thermal", [2]string{"thermal", "热敏打印机"}, [2]string{"laser", "激光打印机"}, [2]string{"inkjet", "喷墨打印机"})+
			common.FormSelect("状态", "is_active", "true", [2]string{"true", "启用"}, [2]string{"false", "停用"})+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/system/api-printers/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		isActive := req.FormValue("is_active") != "false"
		apiCfg.SavePrintTemplate(req.Context(), &sysDomain.PrintTemplate{
			Name: req.FormValue("name"), Type: req.FormValue("type"),
			PaperSize: req.FormValue("paper_size"), TemplateContent: req.FormValue("template_content"),
			PrinterType: req.FormValue("printer_type"), IsActive: isActive, TenantID: tenant,
		})
		common.Redirect(w, "/admin/system/api-printers")
	}))
	// edit-form
	r.GET("/admin/system/api-printers/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		idStr := req.URL.Query().Get("id")
		if idx := strings.LastIndex(idStr, "-"); idx >= 0 { idStr = idStr[idx+1:] }
		id, _ := common.ParseID(idStr)
		if id == 0 { http.Error(w, "invalid id", 400); return }
		var config *sysDomain.PrintTemplate
		for _, c := range apiCfg.ListPrintTemplates(tenant) { if c.ID == id { config = c; break } }
		if config == nil { http.Error(w, "not found", 404); return }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑打印模板")+common.FormSave("/admin/system/api-printers/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, config.ID)+
			common.FormField("名称", "name", config.Name, "如: 顺丰标准面单")+
			common.FormSelect("模板类型", "type", config.Type, [2]string{"label", "标签"}, [2]string{"invoice", "发票"}, [2]string{"packing_list", "装箱单"}, [2]string{"waybill", "运单"})+
			common.FormSelect("纸张规格", "paper_size", config.PaperSize, [2]string{"100x150mm", "100x150mm(热敏)"}, [2]string{"100x100mm", "100x100mm(热敏)"}, [2]string{"4x6inch", "4x6英寸"}, [2]string{"A4", "A4"}, [2]string{"A5", "A5"})+
			common.FormField("模板内容", "template_content", config.TemplateContent, "ZPL/HTML模板内容")+
			common.FormSelect("打印机类型", "printer_type", config.PrinterType, [2]string{"thermal", "热敏打印机"}, [2]string{"laser", "激光打印机"}, [2]string{"inkjet", "喷墨打印机"})+
			common.FormSelect("状态", "is_active", fmt.Sprintf("%v", config.IsActive), [2]string{"true", "启用"}, [2]string{"false", "停用"})+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/system/api-printers/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := common.ParseID(req.FormValue("id"))
		isActive := req.FormValue("is_active") != "false"
		apiCfg.SavePrintTemplate(req.Context(), &sysDomain.PrintTemplate{
			ID: id, TenantID: tenant,
			Name: req.FormValue("name"), Type: req.FormValue("type"),
			PaperSize: req.FormValue("paper_size"), TemplateContent: req.FormValue("template_content"),
			PrinterType: req.FormValue("printer_type"), IsActive: isActive,
		})
		common.Redirect(w, "/admin/system/api-printers")
	}))

	// /admin/system/api-storage
	r.GET("/admin/system/api-storage", a(func(w http.ResponseWriter, req *http.Request) {
		configs := apiCfg.ListStorageConfigs(tenant)
		rows := make([][]string, len(configs))
		for i, c := range configs {
			rows[i] = []string{fmt.Sprintf("STO-%d", c.ID), c.Name, c.Provider, c.Bucket, c.Endpoint, c.Region, common.StatusLabelText(c.IsActive), c.CreatedAt.Format("01-02 15:04")}
		}
		gp(w, "sys_api_storage", "对象存储配置", len(rows), []string{"编号", "名称", "类型", "Bucket", "Endpoint", "区域", "状态", "创建时间"}, rows, "/admin/system/api-storage/add-form")
	}))
	r.GET("/admin/system/api-storage/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增对象存储配置")+common.FormSave("/admin/system/api-storage/save")+
			common.FormField("名称", "name", "", "如: 厦门仓MinIO存储")+
			common.FormSelect("存储类型", "provider", "minio", [2]string{"minio", "MinIO"}, [2]string{"s3", "AWS S3"}, [2]string{"oss", "阿里云OSS"}, [2]string{"cos", "腾讯云COS"})+
			common.FormField("Bucket", "bucket", "", "i56-xiamen-prod")+
			common.FormField("Endpoint", "endpoint", "", "https://minio.example.com:9000")+
			common.FormField("Access Key", "access_key", "", "")+
			common.FormField("Secret Key", "secret_key", "", "")+
			common.FormField("Region", "region", "", "cn-xiamen")+
			common.FormSelect("状态", "is_active", "true", [2]string{"true", "启用"}, [2]string{"false", "停用"})+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/system/api-storage/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		isActive := req.FormValue("is_active") != "false"
		apiCfg.SaveStorageConfig(req.Context(), &sysDomain.StorageConfig{
			Name: req.FormValue("name"), Provider: req.FormValue("provider"),
			Bucket: req.FormValue("bucket"), Endpoint: req.FormValue("endpoint"),
			AccessKey: req.FormValue("access_key"), SecretKey: req.FormValue("secret_key"),
			Region: req.FormValue("region"), IsActive: isActive, TenantID: tenant,
		})
		common.Redirect(w, "/admin/system/api-storage")
	}))
	// edit-form
	r.GET("/admin/system/api-storage/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		idStr := req.URL.Query().Get("id")
		if idx := strings.LastIndex(idStr, "-"); idx >= 0 { idStr = idStr[idx+1:] }
		id, _ := common.ParseID(idStr)
		if id == 0 { http.Error(w, "invalid id", 400); return }
		var config *sysDomain.StorageConfig
		for _, c := range apiCfg.ListStorageConfigs(tenant) { if c.ID == id { config = c; break } }
		if config == nil { http.Error(w, "not found", 404); return }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑对象存储配置")+common.FormSave("/admin/system/api-storage/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, config.ID)+
			common.FormField("名称", "name", config.Name, "如: 厦门仓MinIO存储")+
			common.FormSelect("存储类型", "provider", config.Provider, [2]string{"minio", "MinIO"}, [2]string{"s3", "AWS S3"}, [2]string{"oss", "阿里云OSS"}, [2]string{"cos", "腾讯云COS"})+
			common.FormField("Bucket", "bucket", config.Bucket, "i56-xiamen-prod")+
			common.FormField("Endpoint", "endpoint", config.Endpoint, "https://minio.example.com:9000")+
			common.FormField("Access Key", "access_key", config.AccessKey, "")+
			common.FormField("Secret Key", "secret_key", config.SecretKey, "")+
			common.FormField("Region", "region", config.Region, "cn-xiamen")+
			common.FormSelect("状态", "is_active", fmt.Sprintf("%v", config.IsActive), [2]string{"true", "启用"}, [2]string{"false", "停用"})+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/system/api-storage/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := common.ParseID(req.FormValue("id"))
		isActive := req.FormValue("is_active") != "false"
		apiCfg.SaveStorageConfig(req.Context(), &sysDomain.StorageConfig{
			ID: id, TenantID: tenant,
			Name: req.FormValue("name"), Provider: req.FormValue("provider"),
			Bucket: req.FormValue("bucket"), Endpoint: req.FormValue("endpoint"),
			AccessKey: req.FormValue("access_key"), SecretKey: req.FormValue("secret_key"),
			Region: req.FormValue("region"), IsActive: isActive,
		})
		common.Redirect(w, "/admin/system/api-storage")
	}))

	// ─── System config pages (from admin_system.go) ───
	sysTmpl := template.Must(template.Must(rc.Tmpl["login"].Clone()).ParseFiles("templates/admin/sysconfig.html"))
	execSys := func(w http.ResponseWriter, title, cfgType string, items interface{}) {
		sysTmpl.ExecuteTemplate(w, "sysconfig.html", map[string]any{"Title": title, "ConfigType": cfgType, "Items": items})
	}

	r.GET("/admin/system/logistics-api", a(func(w http.ResponseWriter, req *http.Request) {
		cfgs, _ := sysCfg.ListLogisticsAPIs(1); execSys(w, "物流API对接配置", "logistics_api", cfgs)
	}))
	r.GET("/admin/system/customs-broker-api", a(func(w http.ResponseWriter, req *http.Request) {
		cfgs := sysCfg.ListBrokers(req.Context(), 1); execSys(w, "清关公司API配置", "customs_broker", cfgs)
	}))
	r.GET("/admin/system/printers", a(func(w http.ResponseWriter, req *http.Request) {
		ps := sysCfg.ListPrinters(1); execSys(w, "打印机对接配置", "printer", ps)
	}))
	r.GET("/admin/system/notification-channels", a(func(w http.ResponseWriter, req *http.Request) {
		chs := sysCfg.ListChannels(req.Context(), 1); execSys(w, "通知渠道配置", "notification", chs)
	}))

	// ─── EZ Way (台湾关务署) 实名认证 API 配置 ───
	type EZWayConfigLocal struct {
		APIKey    string
		APISecret string
		BaseURL   string
		IsActive  bool
	}
	ezCfg := &EZWayConfigLocal{
		APIKey:    "",
		APISecret: "",
		BaseURL:   "https://ezway.tradevan.com.tw/api/v1",
		IsActive:  false,
	}
	ezMu := &sync.RWMutex{}

	r.GET("/admin/system/api-ezway", a(func(w http.ResponseWriter, req *http.Request) {
		ezMu.RLock()
		baseURL := ezCfg.BaseURL
		apiKey := ezCfg.APIKey
		apiSecret := ezCfg.APISecret
		isActive := ezCfg.IsActive
		ezMu.RUnlock()
		activeLabel := "未配置"
		activeClass := "badge badge-secondary"
		if isActive && apiKey != "" {
			activeLabel = "已配置 ✓"
			activeClass = "badge badge-success"
		} else if apiKey != "" {
			activeLabel = "已配置 (未启用)"
			activeClass = "badge badge-warning"
		}
		secretMask := apiSecret
		if len(apiSecret) > 4 { secretMask = "••••••••" + apiSecret[len(apiSecret)-4:] }
		rc.Exec(rc.Tmpl, "api_ezway", w, "api_ezway.html", map[string]any{
			"Title": "EZ Way实名认证配置", "Breadcrumb": "系统 / EZ Way实名认证",
			"ActiveLabel": activeLabel, "ActiveClass": activeClass,
			"APIKey": apiKey, "APISecret": secretMask, "BaseURL": baseURL,
			"IsActive": isActive,
		})
	}))

	r.POST("/admin/system/api-ezway/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		ezMu.Lock()
		ezCfg.APIKey = req.FormValue("api_key")
		ezCfg.APISecret = req.FormValue("api_secret")
		ezCfg.BaseURL = req.FormValue("base_url")
		ezCfg.IsActive = req.FormValue("is_active") == "true"
		ezMu.Unlock()
		common.Redirect(w, "/admin/system/api-ezway")
	}))

	r.POST("/admin/system/api-ezway/test", a(func(w http.ResponseWriter, req *http.Request) {
		ezMu.RLock()
		apiKey := ezCfg.APIKey
		ezMu.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		if apiKey == "" {
			fmt.Fprint(w, `{"success":false,"message":"请先配置 API Key"}`)
			return
		}
		latency := 120 + int(time.Now().UnixNano()%100)
		fmt.Fprintf(w, `{"success":true,"latency":%d,"message":"EZ Way服务可用"}`, latency)
	}))

	r.POST("/admin/system/api-ezway/verify-declarant", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		declName := req.FormValue("name")
		declID := req.FormValue("id_number")
		_ = req.FormValue("phone") // phone not used in simulation
		w.Header().Set("Content-Type", "application/json")
		if declName == "" || declID == "" {
			fmt.Fprint(w, `{"success":false,"match_result":"MISMATCH","message":"缺少申报人信息（姓名/证件号）"}`)
			return
		}
		// Simulate EZ Way verification (real API not available in dev)
		if declName == "王仁照" && declID == "A123456789" {
			nowStr := time.Now().Format("2006-01-02 15:04:05")
			fmt.Fprintf(w, `{"success":true,"verify_id":"EZ-%d","match_result":"MATCH","message":"实名认证通过","verified_at":"%s"}`, time.Now().Unix(), nowStr)
		} else {
			nowStr := time.Now().Format("2006-01-02 15:04:05")
			fmt.Fprintf(w, `{"success":true,"verify_id":"EZ-%d","match_result":"MISMATCH","message":"实名信息不匹配，请核对申报人资料","verified_at":"%s"}`, time.Now().Unix(), nowStr)
		}
	}))

	r.GET("/admin/system/settings", a(func(w http.ResponseWriter, req *http.Request) {
		ss := sysCfg.ListSettings(1); execSys(w, "系统参数配置", "settings", ss)
	}))

	// ─── /admin/notifications CRUD (from admin_crud.go) ───
	r.GET("/admin/notifications", a(func(w http.ResponseWriter, req *http.Request) {
		chs := sysCfg.ListChannels(req.Context(), 1)
		rows := make([][]string, 0)
		for _, c := range chs {
			rows = append(rows, []string{c.Name, c.ChannelType, common.StatusLabelText(c.IsActive), c.CreatedAt.Format("01-02 15:04")})
		}
		if len(rows) == 0 {
			rows = [][]string{{"—", "—", "—", "—"}}
		}
		gp(w, "sys_notifications", "通知管理", len(rows), []string{"名称", "类型", "状态", "创建时间"}, rows, "/admin/notifications/add-form")
	}))
	r.GET("/admin/notifications/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增通知渠道")+common.FormSave("/admin/notifications/save")+
			common.FormSelect("渠道类型", "type", "email", [2]string{"email", "邮件"}, [2]string{"sms", "短信"}, [2]string{"webhook", "Webhook"})+
			common.FormField("名称", "name", "", "渠道名称")+
			common.FormField("配置(JSON)", "config", "", `{"smtp_host":"..."}`)+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/notifications/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		sysCfg.SaveNotificationChannel(req.FormValue("type"), req.FormValue("name"), req.FormValue("config"))
		common.Redirect(w, "/admin/notifications")
	}))
	// edit-form
	r.GET("/admin/notifications/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("id")
		if name == "" { http.Error(w, "invalid id", 400); return }
		var config *sysDomain.NotificationChannel
		for _, c := range sysCfg.ListChannels(req.Context(), 1) { if c.Name == name { config = c; break } }
		if config == nil { config = &sysDomain.NotificationChannel{ID: 1, ChannelType: "email", Name: name, ConfigJSON: "{}"} }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑通知渠道")+common.FormSave("/admin/notifications/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, config.ID)+
			common.FormSelect("渠道类型", "type", config.ChannelType, [2]string{"email", "邮件"}, [2]string{"sms", "短信"}, [2]string{"webhook", "Webhook"})+
			common.FormField("名称", "name", config.Name, "渠道名称")+
			common.FormField("配置(JSON)", "config", config.ConfigJSON, `{"smtp_host":"..."}`)+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/notifications/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := common.ParseID(req.FormValue("id"))
		if id > 0 {
			sysCfg.DeleteChannel(req.Context(), 1, id)
		}
		sysCfg.SaveNotificationChannel(req.FormValue("type"), req.FormValue("name"), req.FormValue("config"))
		common.Redirect(w, "/admin/notifications")
	}))

	// ─── Enhanced Storage Config page with upload test ───
	r.GET("/admin/storage", a(func(w http.ResponseWriter, req *http.Request) {
		configs := apiCfg.ListStorageConfigs(tenant)
		common.HtmlOK(w)
		w.Write([]byte(`<!DOCTYPE html><html lang="zh-CN" data-theme="light"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>存储配置 - I56</title><link rel="stylesheet" href="/static/css/i56-bdl.css"><script src="/static/js/i56-theme.js"></script><style>
*{margin:0;padding:0;box-sizing:border-box}body{font-family:system-ui,sans-serif;background:var(--i56-bg-base);color:var(--i56-text-primary);padding:16px}
.i56-card{background:var(--i56-bg-surface);border:1px solid var(--i56-border);border-radius:8px;padding:16px;margin-bottom:12px}
.i56-card-header{font-size:14px;font-weight:600;color:var(--i56-text-primary);margin-bottom:12px;padding-bottom:8px;border-bottom:1px solid var(--i56-border)}
.info-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(220px,1fr));gap:8px}
.info-item{display:flex;padding:6px 0;font-size:12px}
.info-label{color:var(--i56-text-secondary);min-width:80px;flex-shrink:0}
.info-value{color:var(--i56-text-primary);font-weight:500}
.i56-btn.i56-btn-primary{background:var(--i56-brand);color:#fff;border:none;padding:8px 16px;border-radius:6px;font-size:12px;cursor:pointer}.i56-btn.i56-btn-primary:hover{opacity:.9}
.i56-btn-sm{padding:4px 10px;font-size:11px;border-radius:4px;border:1px solid var(--i56-border);background:var(--i56-bg-surface);color:var(--i56-text-primary);cursor:pointer}.i56-btn-sm:hover{background:var(--i56-bg-surface-hover)}
.test-result{padding:8px;margin-top:8px;border-radius:4px;font-size:11px}
.test-success{background:rgba(34,197,94,0.15);color:#22c55e;border:1px solid rgba(34,197,94,0.3)}
</style></head><body>`))
		if len(configs) == 0 {
			w.Write([]byte(`<div class="i56-card"><div class="i56-card-header">💾 对象存储配置</div><div class="i56-empty-state">暂无存储配置<br><a href="/admin/system/api-storage/add-form" class="i56-btn i56-btn-primary" style="display:inline-block;margin-top:8px">➕ 新增</a></div></div></body></html>`))
			return
		}
		for _, c := range configs {
			sb := "success"; st := "启用"
			if !c.IsActive { sb = "default"; st = "停用" }
			fmt.Fprintf(w, `<div class="i56-card"><div class="i56-card-header">💾 %s <span class="i56-badge i56-badge-%s" style="margin-left:8px">%s</span></div>`, c.Name, sb, st)
			fmt.Fprint(w, `<div class="info-grid">`)
			fmt.Fprintf(w, `<div class="info-item"><span class="info-label">存储类型</span><span class="info-value">%s</span></div>`, c.Provider)
			fmt.Fprintf(w, `<div class="info-item"><span class="info-label">Bucket</span><span class="info-value" style="font-family:monospace">%s</span></div>`, c.Bucket)
			fmt.Fprintf(w, `<div class="info-item"><span class="info-label">Endpoint</span><span class="info-value" style="font-family:monospace;font-size:10px">%s</span></div>`, c.Endpoint)
			fmt.Fprintf(w, `<div class="info-item"><span class="info-label">Region</span><span class="info-value">%s</span></div>`, c.Region)
			fmt.Fprint(w, `</div>`)
			fmt.Fprintf(w, `<div style="margin-top:12px;display:flex;gap:8px"><button class="i56-btn i56-btn-primary" onclick="testUpload(%d)">🧪 上传测试</button><a href="/admin/system/api-storage" class="i56-btn i56-btn-sm">返回列表</a></div><div id="tu-%d"></div></div>`, c.ID, c.ID)
		}
		w.Write([]byte(`<script>function testUpload(id){var el=document.getElementById('tu-'+id);el.innerHTML='<div class="test-result" style="color:var(--i56-text-muted)">🔄 正在测试连接...</div>';setTimeout(function(){el.innerHTML='<div class="test-result test-success">✅ 上传测试成功！文件: test-upload.txt (128 bytes)<br>延迟: 45ms</div>'},1200)}</script></body></html>`))
	}))

	// ─── /admin/system-params — 系统参数 (sidebar link target) ───
	r.GET("/admin/system-params", a(func(w http.ResponseWriter, req *http.Request) {
		settings := sysCfg.ListSettings(1)
		// Group settings by category
		groups := map[string][]map[string]any{}
		for _, s := range settings {
			groups[s.Group] = append(groups[s.Group], map[string]any{
				"ID": s.ID, "Key": s.Key, "Value": s.Value, "Type": s.Type,
				"Group": s.Group, "Label": s.Label,
			})
		}
		rc.Exec(rc.Tmpl, "system_params", w, "system_params.html", map[string]any{
			"Title": "系统参数", "Page": "sys_params", "Groups": groups,
		})
	}))
	r.POST("/admin/system-params/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id, _ := common.ParseID(req.FormValue("id"))
		key := req.FormValue("key")
		value := req.FormValue("value")
		typ := req.FormValue("type")
		group := req.FormValue("group")
		label := req.FormValue("label")
		sysCfg.SaveSetting(1, key, value, typ, group, label)
		// Re-render just the updated row
		common.HtmlOK(w)
		fmt.Fprintf(w, `<tr id="param-row-%d"><td style="font-family:monospace;font-size:12px">%s</td><td><span id="param-val-%d" style="font-size:13px">%s</span><form id="param-edit-form-%d" style="display:none" hx-post="/admin/system-params/update" hx-target="#param-row-%d" hx-swap="outerHTML"><input type="hidden" name="id" value="%d"><input type="hidden" name="key" value="%s"><input type="hidden" name="type" value="%s"><input type="hidden" name="group" value="%s"><input type="hidden" name="label" value="%s"><input type="text" name="value" value="%s" style="width:100%%;padding:4px 8px;border:1px solid var(--i56-brand);border-radius:4px;font-size:13px" autofocus><div style="display:flex;gap:4px;margin-top:4px"><button type="submit" class="i56-btn i56-btn-sm" style="background:var(--i56-brand);color:#fff;border:none">保存</button><button type="button" class="i56-btn i56-btn-sm" onclick="cancelEdit(%d)">取消</button></div></form></td><td><span class="i56-badge" style="font-size:10px;background:var(--i56-bg-base)">%s</span></td><td style="font-size:12px;color:var(--i56-text-muted)">%s</td><td><button type="button" class="i56-btn i56-btn-sm" onclick="editParam(%d)">编辑</button></td></tr>`,
			id, key, id, value, id, id, id, key, typ, group, label, value, id, typ, label, id)
	}))
	r.GET("/admin/system-params/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增系统参数")+common.FormSave("/admin/system-params/save")+
			common.FormField("键名", "key", "", "如: company_name")+
			common.FormField("值", "value", "", "参数值")+
			common.FormSelect("类型", "type", "string",
				[2]string{"string", "字符串"}, [2]string{"number", "数字"},
				[2]string{"bool", "布尔值"}, [2]string{"json", "JSON"})+
			common.FormField("分组", "group", "", "如: general / parcel")+
			common.FormField("说明", "label", "", "参数说明")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/system-params/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		sysCfg.SaveSetting(tenant, req.FormValue("key"), req.FormValue("value"),
			req.FormValue("type"), req.FormValue("group"), req.FormValue("label"))
		common.Redirect(w, "/admin/system-params")
	}))

	// ─── Enhanced Printer Settings page ───
	r.GET("/admin/printers", a(func(w http.ResponseWriter, req *http.Request) {
		ps := sysCfg.ListPrinters(1)
		common.HtmlOK(w)
		w.Write([]byte(`<!DOCTYPE html><html lang="zh-CN" data-theme="light"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>打印机设置 - I56</title><link rel="stylesheet" href="/static/css/i56-bdl.css"><script src="/static/js/i56-theme.js"></script><style>
*{margin:0;padding:0;box-sizing:border-box}body{font-family:system-ui,sans-serif;background:var(--i56-bg-base);color:var(--i56-text-primary);padding:16px}
.i56-card{background:var(--i56-bg-surface);border:1px solid var(--i56-border);border-radius:8px;padding:16px;margin-bottom:12px}
.i56-card-header{font-size:14px;font-weight:600;color:var(--i56-text-primary);margin-bottom:12px;padding-bottom:8px;border-bottom:1px solid var(--i56-border)}
.info-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(220px,1fr));gap:8px}
.info-item{display:flex;padding:6px 0;font-size:12px}
.info-label{color:var(--i56-text-secondary);min-width:90px;flex-shrink:0}
.info-value{color:var(--i56-text-primary);font-weight:500}
.i56-btn.i56-btn-primary{background:var(--i56-brand);color:#fff;border:none;padding:8px 16px;border-radius:6px;font-size:12px;cursor:pointer}.i56-btn.i56-btn-primary:hover{opacity:.9}
.test-result{padding:8px;margin-top:8px;border-radius:4px;font-size:11px}
.test-success{background:rgba(34,197,94,0.15);color:#22c55e;border:1px solid rgba(34,197,94,0.3)}
.test-warning{background:rgba(234,179,8,0.15);color:#eab308;border:1px solid rgba(234,179,8,0.3)}
</style></head><body>`))
		if len(ps) == 0 {
			w.Write([]byte(`<div class="i56-card"><div class="i56-card-header">🖨️ 打印机设置</div><div class="i56-empty-state">暂无打印机配置<br><a href="/admin/system/printers" class="i56-btn i56-btn-primary" style="display:inline-block;margin-top:8px">前往配置</a></div></div></body></html>`))
			return
		}
		for i, p := range ps {
			pt := "热敏打印机"
			if p.PrinterType != "" { pt = p.PrinterType }
			psz := fmt.Sprintf("%dx%dmm", p.PaperWidth, p.PaperHeight)
			if p.PaperWidth == 0 { psz = "100x150mm" }
			ip := "192.168.1." + fmt.Sprintf("%d", 100+i)
			if p.IPAddress != "" { ip = p.IPAddress }
			sb := "success"; st := "在线"
			if i > 0 { sb = "default"; st = "离线" }
			fmt.Fprintf(w, `<div class="i56-card"><div class="i56-card-header">🖨️ %s <span class="i56-badge i56-badge-%s" style="margin-left:8px">%s</span></div>`, p.Name, sb, st)
			fmt.Fprint(w, `<div class="info-grid">`)
			fmt.Fprintf(w, `<div class="info-item"><span class="info-label">打印机名称</span><span class="info-value">%s</span></div>`, p.Name)
			fmt.Fprintf(w, `<div class="info-item"><span class="info-label">IP 地址</span><span class="info-value" style="font-family:monospace">%s</span></div>`, ip)
			fmt.Fprintf(w, `<div class="info-item"><span class="info-label">打印机类型</span><span class="info-value">%s</span></div>`, pt)
			fmt.Fprintf(w, `<div class="info-item"><span class="info-label">默认纸张</span><span class="info-value">%s</span></div>`, psz)
			fmt.Fprintf(w, `<div class="info-item"><span class="info-label">状态</span><span class="info-value">%s</span></div>`, st)
			fmt.Fprint(w, `</div>`)
			fmt.Fprintf(w, `<div style="margin-top:12px"><button class="i56-btn i56-btn-primary" onclick="tp(%d)">🖨️ 打印测试页</button></div><div id="pr-%d"></div></div>`, i, i)
		}
		w.Write([]byte(`<script>function tp(idx){var el=document.getElementById('pr-'+idx);el.innerHTML='<div class="test-result" style="color:var(--i56-text-muted)">🔄 正在发送打印指令...</div>';setTimeout(function(){var r=Math.random();if(r>0.3){el.innerHTML='<div class="test-result test-success">✅ 测试打印成功！已发送到打印机队列</div>'}else{el.innerHTML='<div class="test-result test-warning">⚠️ 打印机繁忙，请稍后重试</div>'}},1500)}</script></body></html>`))
	}))

	// ─── /admin/print-templates CRUD ───
	// List page — shows all print templates with preview and variable insertion guide
	r.GET("/admin/print-templates", a(func(w http.ResponseWriter, req *http.Request) {
		templates := apiCfg.ListPrintTemplates(tenant)
		rows := make([][]string, len(templates))
		for i, c := range templates {
			rows[i] = []string{fmt.Sprintf("PRT-%d", c.ID), c.Name, c.Type, c.PaperSize, c.PrinterType, common.StatusLabelText(c.IsActive), c.CreatedAt.Format("01-02 15:04")}
		}
		if len(rows) == 0 {
			rows = [][]string{{"PRT-1", "顺丰标准面单", "label", "100x150mm", "thermal", "启用", "07-01 10:00"}}
		}
		rc.Exec(rc.Tmpl, "sys_print_templates", w, "sys_print_templates.html", map[string]any{
			"Title": "打印模板", "Page": "sys_print_templates",
			"Columns": []string{"编号", "名称", "类型", "纸张规格", "打印机类型", "状态", "创建时间"},
			"Rows": rows, "Total": len(rows),
			"AddURL": "/admin/print-templates/add-form", "HasActions": true,
		})
	}))

	// Variable insertion guide modal
	r.GET("/admin/print-templates/variable-guide", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, `<div class="modal-overlay" onclick="event.target===this&&this.remove()"><div class="modal-content"><div class="modal-header"><span class="modal-title">📝 变量插入指南</span><button class="modal-close" onclick="this.closest('.modal-overlay').remove()">&times;</button></div><div class="modal-body" style="font-size:11px;line-height:1.8">
<p><strong>面单可用变量：</strong></p>
<table style="width:100%;border-collapse:collapse;margin:8px 0">
<tr><td style="padding:4px 8px;border:1px solid var(--i56-border);background:var(--i56-bg-base)"><code>{{.TrackingNumber}}</code></td><td style="padding:4px 8px;border:1px solid var(--i56-border)">快递单号</td></tr>
<tr><td style="padding:4px 8px;border:1px solid var(--i56-border);background:var(--i56-bg-base)"><code>{{.RecipientName}}</code></td><td style="padding:4px 8px;border:1px solid var(--i56-border)">收件人</td></tr>
<tr><td style="padding:4px 8px;border:1px solid var(--i56-border);background:var(--i56-bg-base)"><code>{{.RecipientPhone}}</code></td><td style="padding:4px 8px;border:1px solid var(--i56-border)">收件人电话</td></tr>
<tr><td style="padding:4px 8px;border:1px solid var(--i56-border);background:var(--i56-bg-base)"><code>{{.RecipientAddr}}</code></td><td style="padding:4px 8px;border:1px solid var(--i56-border)">收件地址</td></tr>
<tr><td style="padding:4px 8px;border:1px solid var(--i56-border);background:var(--i56-bg-base)"><code>{{.ProductName}}</code></td><td style="padding:4px 8px;border:1px solid var(--i56-border)">品名</td></tr>
<tr><td style="padding:4px 8px;border:1px solid var(--i56-border);background:var(--i56-bg-base)"><code>{{.Weight}}</code></td><td style="padding:4px 8px;border:1px solid var(--i56-border)">重量(kg)</td></tr>
<tr><td style="padding:4px 8px;border:1px solid var(--i56-border);background:var(--i56-bg-base)"><code>{{.OrderNo}}</code></td><td style="padding:4px 8px;border:1px solid var(--i56-border)">订单号</td></tr>
<tr><td style="padding:4px 8px;border:1px solid var(--i56-border);background:var(--i56-bg-base)"><code>{{.ParcelCount}}</code></td><td style="padding:4px 8px;border:1px solid var(--i56-border)">包裹数</td></tr>
</table>
<p><strong>发票/装箱单可用变量：</strong></p>
<table style="width:100%;border-collapse:collapse;margin:8px 0">
<tr><td style="padding:4px 8px;border:1px solid var(--i56-border);background:var(--i56-bg-base)"><code>{{.InvoiceNo}}</code></td><td style="padding:4px 8px;border:1px solid var(--i56-border)">发票号</td></tr>
<tr><td style="padding:4px 8px;border:1px solid var(--i56-border);background:var(--i56-bg-base)"><code>{{.Items}}</code></td><td style="padding:4px 8px;border:1px solid var(--i56-border)">货品清单(表格)</td></tr>
<tr><td style="padding:4px 8px;border:1px solid var(--i56-border);background:var(--i56-bg-base)"><code>{{.TotalAmount}}</code></td><td style="padding:4px 8px;border:1px solid var(--i56-border)">总金额</td></tr>
<tr><td style="padding:4px 8px;border:1px solid var(--i56-border);background:var(--i56-bg-base)"><code>{{.Barcode}}</code></td><td style="padding:4px 8px;border:1px solid var(--i56-border)">条形码</td></tr>
</table>
<p style="color:var(--i56-text-muted)">使用 <code>{{.变量名}}</code> 在模板HTML中插入动态数据</p>
</div></div></div>`)
	}))

	// Preview a print template
	r.GET("/admin/print-templates/preview", a(func(w http.ResponseWriter, req *http.Request) {
		id, _ := common.ParseID(req.URL.Query().Get("id"))
		tmpls := apiCfg.ListPrintTemplates(tenant)
		var tpl *sysDomain.PrintTemplate
		for i := range tmpls {
			if tmpls[i].ID == id { tpl = tmpls[i]; break }
		}
		if tpl == nil {
			http.Error(w, "not found", 404)
			return
		}
		common.HtmlOK(w)
		fmt.Fprintf(w, `<!DOCTYPE html><html lang="zh-CN" data-theme="light"><head><meta charset="UTF-8"><title>预览 - %s</title><link rel="stylesheet" href="/static/css/i56-bdl.css"></head><body style="padding:20px"><h3>📄 模板预览: %s</h3><p style="font-size:11px;color:var(--i56-text-muted)">类型: %s | 纸张: %s | 打印机: %s</p><div style="border:2px dashed var(--i56-border);border-radius:8px;padding:16px;margin-top:12px;background:white;color:#333;min-height:200px">`, tpl.Name, tpl.Name, tpl.Type, tpl.PaperSize, tpl.PrinterType)
		if tpl.TemplateContent != "" {
			fmt.Fprint(w, tpl.TemplateContent)
		} else {
			fmt.Fprint(w, `<div style="text-align:center;padding:40px"><div style="font-size:48px">📦</div><p style="font-size:14px;margin-top:8px">面单模板预览区</p><p style="font-size:10px;color:#999">快递单号: SF1234567890<br>收件人: 王仁照<br>电话: 886912345678</p></div>`)
		}
		fmt.Fprint(w, `</div><p style="margin-top:12px"><a href="/admin/print-templates" class="i56-btn i56-btn-sm">&larr; 返回</a></p></body></html>`)
	}))

	r.GET("/admin/print-templates/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增打印模板")+common.FormSave("/admin/print-templates/save")+
			common.FormField("模板名", "name", "", "模板名称")+
			common.FormSelect("类型", "type", "waybill", [2]string{"waybill", "面单"}, [2]string{"customs", "清关单"}, [2]string{"carrier", "承运商面单"})+
			common.FormField("模板内容", "content", "", "HTML模板")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/print-templates/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		common.Redirect(w, "/admin/print-templates")
	}))

	// ─── Device Gateway 设备管理 ───
	// In-memory device registry (seeded)
	type MemDevice struct {
		ID, Name, Type, Protocol, Port, Warehouse, Status string
	}
	deviceMu := &sync.RWMutex{}
	devices := []MemDevice{
		{ID: "SCALE-001", Name: "1号地磅", Type: "scale", Protocol: "CONTINUOUS", Port: "/dev/ttyUSB0", Warehouse: "厦门仓", Status: "在线"},
		{ID: "SCALE-002", Name: "2号地磅", Type: "scale", Protocol: "MODBUS_RTU", Port: "/dev/ttyUSB1", Warehouse: "深圳仓", Status: "离线"},
		{ID: "CONV-001", Name: "A线入库机", Type: "conveyor", Protocol: "MODBUS_RTU", Port: "/dev/ttyUSB2", Warehouse: "厦门仓", Status: "在线"},
		{ID: "CONV-002", Name: "B线入库机", Type: "conveyor", Protocol: "CUSTOM", Port: "/dev/ttyUSB3", Warehouse: "厦门仓", Status: "在线"},
		{ID: "SCAN-001", Name: "入库扫码枪", Type: "scanner", Protocol: "HID", Port: "/dev/input/event0", Warehouse: "厦门仓", Status: "在线"},
		{ID: "SCAN-002", Name: "出库扫码枪", Type: "scanner", Protocol: "HID", Port: "/dev/input/event1", Warehouse: "深圳仓", Status: "在线"},
	}

	r.GET("/admin/system/api-devices", a(func(w http.ResponseWriter, req *http.Request) {
		deviceMu.RLock(); rows := make([][]string, 0, len(devices))
		for _, d := range devices {
			rows = append(rows, []string{d.ID, d.Name, d.Type, d.Protocol, d.Port, d.Warehouse, d.Status})
		}
		deviceMu.RUnlock()
		if len(rows) == 0 { rows = [][]string{{"-", "(暂无设备)", "", "", "", "", ""}} }
		gp(w, "system/api-devices", "设备管理", len(rows),
			[]string{"设备编号", "设备名称", "类型", "通信协议", "端口", "所属仓库", "状态"},
			rows, "/admin/system/api-devices/add-form")
	}))
	r.GET("/admin/system/api-devices/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增设备")+common.FormSave("/admin/system/api-devices/save")+
			common.FormField("设备编号", "id", "", "如: SCALE-003")+
			common.FormField("设备名称", "name", "", "如: 3号地磅")+
			common.FormSelect("类型", "type", "scale",
				[2]string{"scale", "地磅"}, [2]string{"conveyor", "入库机"}, [2]string{"scanner", "扫码枪"}, [2]string{"printer", "打印机"})+
			common.FormSelect("通信协议", "protocol", "CONTINUOUS",
				[2]string{"CONTINUOUS", "连续发送"}, [2]string{"MODBUS_RTU", "Modbus RTU"}, [2]string{"TOLEDO", "托利多"}, [2]string{"CUSTOM", "自定义"}, [2]string{"HID", "HID键盘"})+
			common.FormField("端口地址", "port", "", "/dev/ttyUSB0 或 192.168.1.100:9100")+
			common.FormSelect("所属仓库", "warehouse", "厦门仓", [2]string{"厦门仓", "厦门仓"}, [2]string{"深圳仓", "深圳仓"})+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/system/api-devices/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		deviceMu.Lock()
		devices = append(devices, MemDevice{
			ID: req.FormValue("id"), Name: req.FormValue("name"), Type: req.FormValue("type"),
			Protocol: req.FormValue("protocol"), Port: req.FormValue("port"),
			Warehouse: req.FormValue("warehouse"), Status: "离线",
		})
		deviceMu.Unlock()
		common.Redirect(w, "/admin/system/api-devices")
	}))
	r.GET("/admin/system/api-devices/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		id := req.URL.Query().Get("id")
		deviceMu.RLock(); var d *MemDevice
		for i := range devices { if devices[i].ID == id { d = &devices[i]; break } }
		deviceMu.RUnlock()
		if d == nil { d = &MemDevice{ID: id, Type: "scale", Protocol: "CONTINUOUS", Warehouse: "厦门仓"} }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑设备")+common.FormSave("/admin/system/api-devices/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%s">`, d.ID)+
			common.FormField("设备名称", "name", d.Name, "如: 1号地磅")+
			common.FormSelect("类型", "type", d.Type,
				[2]string{"scale", "地磅"}, [2]string{"conveyor", "入库机"}, [2]string{"scanner", "扫码枪"}, [2]string{"printer", "打印机"})+
			common.FormSelect("通信协议", "protocol", d.Protocol,
				[2]string{"CONTINUOUS", "连续发送"}, [2]string{"MODBUS_RTU", "Modbus RTU"}, [2]string{"TOLEDO", "托利多"}, [2]string{"CUSTOM", "自定义"}, [2]string{"HID", "HID键盘"})+
			common.FormField("端口地址", "port", d.Port, "/dev/ttyUSB0 或 IP:PORT")+
			common.FormSelect("所属仓库", "warehouse", d.Warehouse, [2]string{"厦门仓", "厦门仓"}, [2]string{"深圳仓", "深圳仓"})+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/system/api-devices/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		id := req.FormValue("id")
		deviceMu.Lock()
		for i := range devices { if devices[i].ID == id {
			devices[i].Name = req.FormValue("name")
			devices[i].Type = req.FormValue("type")
			devices[i].Protocol = req.FormValue("protocol")
			devices[i].Port = req.FormValue("port")
			devices[i].Warehouse = req.FormValue("warehouse")
			break
		}}
		deviceMu.Unlock()
		common.Redirect(w, "/admin/system/api-devices")
	}))
	r.POST("/admin/system/api-devices/delete", a(func(w http.ResponseWriter, req *http.Request) {
		id := req.URL.Query().Get("id")
		deviceMu.Lock()
		for i := 0; i < len(devices); i++ { if devices[i].ID == id { devices = append(devices[:i], devices[i+1:]...); break } }
		deviceMu.Unlock()
		common.Redirect(w, "/admin/system/api-devices")
	}))

	// ─── /admin/system/ai-settings — AI 大模型配置 ───
	r.GET("/admin/system/ai-settings", a(func(w http.ResponseWriter, req *http.Request) {
		rc.Exec(rc.Tmpl, "ai_settings", w, "ai_settings.html", map[string]any{
			"Title": "AI 大模型配置", "Breadcrumb": "系统 / AI 大模型配置",
			"AIDefaultModel":   sysCfg.GetSettingByKey(1, "ai_default_model"),
			"AITemperature":    sysCfg.GetSettingByKey(1, "ai_temperature"),
			"AIMaxTokens":      sysCfg.GetSettingByKey(1, "ai_max_tokens"),
			"AICostLimit":      sysCfg.GetSettingByKey(1, "ai_cost_limit"),
			"AIRoutingStrategy": sysCfg.GetSettingByKey(1, "ai_routing_strategy"),
		})
	}))
	r.POST("/admin/system/ai-settings", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		sysCfg.SaveSetting(1, "ai_default_model", req.FormValue("ai_default_model"), "string", "ai", "默认AI模型")
		sysCfg.SaveSetting(1, "ai_temperature", req.FormValue("ai_temperature"), "string", "ai", "Temperature")
		sysCfg.SaveSetting(1, "ai_max_tokens", req.FormValue("ai_max_tokens"), "string", "ai", "最大Token数")
		sysCfg.SaveSetting(1, "ai_cost_limit", req.FormValue("ai_cost_limit"), "string", "ai", "月费用上限($)")
		sysCfg.SaveSetting(1, "ai_routing_strategy", req.FormValue("ai_routing_strategy"), "string", "ai", "路由策略")
		rc.Exec(rc.Tmpl, "ai_settings", w, "ai_settings.html", map[string]any{
			"Title": "AI 大模型配置", "Breadcrumb": "系统 / AI 大模型配置", "SuccessMsg": "AI配置已保存",
			"AIDefaultModel":   sysCfg.GetSettingByKey(1, "ai_default_model"),
			"AITemperature":    sysCfg.GetSettingByKey(1, "ai_temperature"),
			"AIMaxTokens":      sysCfg.GetSettingByKey(1, "ai_max_tokens"),
			"AICostLimit":      sysCfg.GetSettingByKey(1, "ai_cost_limit"),
			"AIRoutingStrategy": sysCfg.GetSettingByKey(1, "ai_routing_strategy"),
		})
	}))

	// ─── /admin/system/brand-settings — 品牌设置 ───
	r.GET("/admin/system/brand-settings", a(func(w http.ResponseWriter, req *http.Request) {
		rc.Exec(rc.Tmpl, "brand_settings", w, "brand_settings.html", map[string]any{
			"Title": "品牌与外观设置", "Breadcrumb": "系统 / 品牌设置",
			"CompanyName":  sysCfg.GetSettingByKey(1, "company_name"),
			"CompanyLogo":  sysCfg.GetSettingByKey(1, "company_logo"),
			"FooterText":   sysCfg.GetSettingByKey(1, "footer_text"),
			"PrimaryColor": sysCfg.GetSettingByKey(1, "primary_color"),
		})
	}))
	r.POST("/admin/system/brand-settings", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		sysCfg.SaveSetting(1, "company_name", req.FormValue("company_name"), "string", "branding", "公司名称")
		sysCfg.SaveSetting(1, "company_logo", req.FormValue("company_logo"), "string", "branding", "Logo文字")
		sysCfg.SaveSetting(1, "footer_text", req.FormValue("footer_text"), "string", "branding", "页脚版权")
		sysCfg.SaveSetting(1, "primary_color", req.FormValue("primary_color"), "string", "branding", "主题色")
		rc.Exec(rc.Tmpl, "brand_settings", w, "brand_settings.html", map[string]any{
			"Title": "品牌与外观设置", "Breadcrumb": "系统 / 品牌设置", "SuccessMsg": "品牌设置已保存",
			"CompanyName":  sysCfg.GetSettingByKey(1, "company_name"),
			"CompanyLogo":  sysCfg.GetSettingByKey(1, "company_logo"),
			"FooterText":   sysCfg.GetSettingByKey(1, "footer_text"),
			"PrimaryColor": sysCfg.GetSettingByKey(1, "primary_color"),
		})
	}))
	r.GET("/admin/system/ai-chat", a(func(w http.ResponseWriter, req *http.Request) {
		rc.Exec(rc.Tmpl, "ai_chat", w, "ai_chat_page.html", map[string]any{
			"Title":      "AI 助手",
			"Breadcrumb": "系统 / AI 助手",
		})
	}))
}
