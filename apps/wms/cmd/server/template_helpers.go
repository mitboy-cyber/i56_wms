package main

import (
	"html/template"
	"net/http"

	parcelDomain "github.com/i56/modules/parcel/domain"
	orderDomain "github.com/i56/modules/order/domain"
)

// bftParcelStatus maps parcel status strings to Chinese display labels.
// Moved from admin_modules.go (deleted as part of HTMX cleanup).
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

// execTpl executes a pre-loaded template by key and name.
func execTpl(tmpl map[string]*template.Template, key string, w http.ResponseWriter, name string, data any) {
	tmpl[key].ExecuteTemplate(w, name, data)
}

// initTemplates initializes all admin-side Go templates with common FuncMap helpers.
func initTemplates() map[string]*template.Template {
	fm := template.FuncMap{
		"statusColor": func(s parcelDomain.ParcelStatus) string {
			switch s {
			case "pre_declared": return "secondary"
			case "received": return "info"
			case "weighed": return "primary"
			case "stored": return "success"
			case "picked": return "warning"
			case "shipped": return "dark"
			default: return "secondary"
			}
		},
		"orderStatusColor": func(s orderDomain.OrderStatus) string {
			switch s {
			case "pending_picking": return "warning"
			case "picking": return "info"
			case "pending_packing": return "primary"
			default: return "secondary"
			}
		},
		"statusDisplay": func(s string) string { return bftParcelStatus(s) },
		"hasPrefix":     func(s, prefix string) bool { return len(s) >= len(prefix) && s[:len(prefix)] == prefix },
		"add":           func(a, b int) int { return a + b },
		"sub":           func(a, b int) int { return a - b },
		"mul":           func(a, b int) int { return a * b },
		"div":           func(a, b int) int { if b == 0 { return 0 }; return a / b },
	}
	tmpl := map[string]*template.Template{}
	for _, p := range []struct{ k, file string }{
		{"login", "login.html"},
		{"dashboard", "dashboard.html"},
		{"clients", "clients.html"},
		{"parcels", "parcels.html"},
		{"orders", "orders.html"},
		{"warehouses", "warehouses.html"},
		{"routes", "routes.html"},
		{"generic_list", "generic_list.html"},
		{"warehouse_console", "warehouse_console.html"},
		{"admin_permissions", "admin/admin_permissions.html"},
		{"admin_roles", "admin/admin_roles.html"},
		{"admin_users", "admin/admin_users.html"},
		{"admin_client_permissions", "admin/admin_client_permissions.html"},
		{"base_new", "admin/base_new.html"},
		// System page templates (Phase 1: replace inline HTML)
		{"scheduler", "admin/system/scheduler.html"},
		{"audit_logs", "admin/system/audit_logs.html"},
		{"reports", "admin/system/reports.html"},
		{"report_view", "admin/system/report_view.html"},
		{"api_ezway", "admin/system/api_ezway.html"},
		{"ai_chat", "admin/system/ai_chat_page.html"},
		{"ai_settings", "admin/system/ai_settings.html"},
		{"brand_settings", "admin/system/brand_settings.html"},
		// System params new template
		{"system_params", "admin/system/system_params.html"},
		// TMS module — data_table-based templates (P3)
		{"carriers", "admin/tms/carriers.html"},
		{"couriers", "admin/tms/couriers.html"},
		{"area_groups", "admin/tms/area_groups.html"},
		{"route_templates", "admin/tms/route_templates.html"},
		// OMS module — data_table-based templates (P3b)
		{"oms_orders", "admin/oms/orders.html"},
		{"oms_service_orders", "admin/oms/service_orders.html"},
		// WMS module — data_table-based templates (P3c)
		{"wms_parcels", "admin/wms/parcels.html"},
		{"wms_warehouses", "admin/wms/warehouses.html"},
		{"wms_exceptions", "admin/wms/exceptions.html"},
		{"wms_service_workorders", "admin/wms/service_workorders.html"},
		{"wms_service_templates", "admin/wms/service_templates.html"},
		// CRM module — data_table-based templates (P3d)
		{"crm_clients", "admin/crm/crm_clients.html"},
		{"crm_accounts", "admin/crm/crm_accounts.html"},
		{"crm_members", "admin/crm/crm_members.html"},
		{"crm_addresses", "admin/crm/crm_addresses.html"},
		{"crm_declarants", "admin/crm/crm_declarants.html"},
		// SYS module — data_table-based templates (P3d)
		{"sys_roles", "admin/sys/sys_roles.html"},
		{"sys_employees", "admin/sys/sys_employees.html"},
		{"sys_print_templates", "admin/sys/sys_print_templates.html"},
		// FIN module — data_table-based templates (P3d)
		{"fin_order_profit", "admin/fin/fin_order_profit.html"},
		{"fin_route_profit", "admin/fin/fin_route_profit.html"},
	} {
		if p.k == "carriers" || p.k == "couriers" || p.k == "area_groups" || p.k == "route_templates" ||
			p.k == "oms_orders" || p.k == "oms_service_orders" ||
			p.k == "wms_parcels" || p.k == "wms_warehouses" || p.k == "wms_exceptions" || p.k == "wms_service_workorders" || p.k == "wms_service_templates" ||
			p.k == "crm_clients" || p.k == "crm_accounts" || p.k == "crm_members" || p.k == "crm_addresses" || p.k == "crm_declarants" ||
			p.k == "sys_roles" || p.k == "sys_employees" || p.k == "sys_print_templates" ||
			p.k == "ai_chat" || p.k == "ai_settings" || p.k == "brand_settings" || p.k == "system_params" ||
			p.k == "fin_order_profit" || p.k == "fin_route_profit" {
			files := []string{"templates/sidebar.html",
		"templates/base.html", "templates/sidebar.html", "templates/admin/admin_layout.html", "templates/admin/partials/data_table.html", "templates/" + p.file}
			tmpl[p.k] = template.Must(template.New(p.k).Funcs(fm).ParseFiles(files...))
		} else if p.k == "base_new" {
			files := []string{"templates/sidebar.html",
		"templates/base.html", "templates/admin/base_new.html"}
			tmpl[p.k] = template.Must(template.New(p.k).Funcs(fm).ParseFiles(files...))
		} else {
			files := []string{"templates/sidebar.html",
		"templates/base.html", "templates/sidebar.html", "templates/" + p.file}
			tmpl[p.k] = template.Must(template.New(p.k).Funcs(fm).ParseFiles(files...))
		}
	}
	return tmpl
}

// initClientTemplates initializes all client-side Go templates with common FuncMap helpers.
func initClientTemplates() map[string]*template.Template {
	fm := template.FuncMap{
		"statusColor": func(s parcelDomain.ParcelStatus) string {
			switch s {
			case "pre_declared": return "secondary"
			case "received": return "info"
			case "weighed": return "primary"
			case "stored": return "success"
			case "picked": return "warning"
			case "shipped": return "dark"
			default: return "secondary"
			}
		},
		"orderStatusColor": func(s orderDomain.OrderStatus) string {
			switch s {
			case "pending_picking": return "warning"
			case "picking": return "info"
			case "pending_packing": return "primary"
			default: return "secondary"
			}
		},
		"statusDisplay": func(s string) string { return bftParcelStatus(s) },
		"hasPrefix":     func(s, prefix string) bool { return len(s) >= len(prefix) && s[:len(prefix)] == prefix },
		"add":           func(a, b int) int { return a + b },
		"sub":           func(a, b int) int { return a - b },
		"mul":           func(a, b int) int { return a * b },
		"div":           func(a, b int) int { if b == 0 { return 0 }; return a / b },
	}
	return map[string]*template.Template{
		"login":                    template.Must(template.New("clogin").Funcs(fm).ParseFiles("templates/client/base.html", "templates/client/login.html")),
		"dashboard":                template.Must(template.New("cdash").Funcs(fm).ParseFiles("templates/client/base.html", "templates/client/dashboard.html")),
		"predeclare":               template.Must(template.New("cpred").Funcs(fm).ParseFiles("templates/client/base.html", "templates/client/predeclare.html")),
		"parcels":                  template.Must(template.New("cparcels").Funcs(fm).ParseFiles("templates/client/base.html", "templates/client/parcels.html")),
		"ledger":                   template.Must(template.New("cledger").Funcs(fm).ParseFiles("templates/client/base.html", "templates/client/ledger.html")),
		"client_orders":            template.Must(template.New("corders2").Funcs(fm).ParseFiles("templates/client/base.html", "templates/client/client_orders.html")),
		"client_order_new":         template.Must(template.New("cordnew").Funcs(fm).ParseFiles("templates/client/base.html", "templates/client/client_order_new.html")),
		"client_declarants":        template.Must(template.New("cdecl").Funcs(fm).ParseFiles("templates/client/base.html", "templates/client/client_declarants.html")),
		"client_members":           template.Must(template.New("cmemb").Funcs(fm).ParseFiles("templates/client/base.html", "templates/client/client_members.html")),
		"client_addresses":         template.Must(template.New("caddr").Funcs(fm).ParseFiles("templates/client/base.html", "templates/client/client_addresses.html")),
		"client_warehouses":        template.Must(template.New("cwh").Funcs(fm).ParseFiles("templates/client/base.html", "templates/client/client_warehouses.html")),
		"client_route_prices":      template.Must(template.New("crp").Funcs(fm).ParseFiles("templates/client/base.html", "templates/client/client_route_prices.html")),
		"client_delivery_fees":     template.Must(template.New("cdf").Funcs(fm).ParseFiles("templates/client/base.html", "templates/client/client_delivery_fees.html")),
		"client_service_orders":    template.Must(template.New("cso").Funcs(fm).ParseFiles("templates/client/base.html", "templates/client/client_service_orders.html")),
		"client_carrier_surcharges": template.Must(template.New("ccs").Funcs(fm).ParseFiles("templates/client/base.html", "templates/client/client_carrier_surcharges.html")),
		"client_webhooks":          template.Must(template.New("cwh").Funcs(fm).ParseFiles("templates/client/base.html", "templates/client/client_webhooks.html")),
		"client_api_credentials":   template.Must(template.New("capi").Funcs(fm).ParseFiles("templates/client/base.html", "templates/client/client_api_credentials.html")),
		"client_monthly_statements": template.Must(template.New("cms").Funcs(fm).ParseFiles("templates/client/base.html", "templates/client/client_monthly_statements.html")),
		"client_weight_dashboard":  template.Must(template.New("cwd").Funcs(fm).ParseFiles("templates/client/base.html", "templates/client/client_weight_dashboard.html")),
		"client_webhook_logs":      template.Must(template.New("cwhlog").Funcs(fm).ParseFiles("templates/client/base.html", "templates/client/client_webhook_logs.html")),
		"client_warehouse_info":    template.Must(template.New("cwhinfo").Funcs(fm).ParseFiles("templates/client/base.html", "templates/client/client_warehouse_info.html")),
		"client_carrier_delivery":  template.Must(template.New("ccdel").Funcs(fm).ParseFiles("templates/client/base.html", "templates/client/client_carrier_delivery.html")),
		"client_carrier_surcharge": template.Must(template.New("ccsur").Funcs(fm).ParseFiles("templates/client/base.html", "templates/client/client_carrier_surcharge.html")),
		"client_pricing":           template.Must(template.New("cpric").Funcs(fm).ParseFiles("templates/client/base.html", "templates/client/client_pricing.html")),
		"client_order_detail":      template.Must(template.New("codet").Funcs(fm).ParseFiles("templates/client/base.html", "templates/client/client_order_detail.html")),
	}
}
