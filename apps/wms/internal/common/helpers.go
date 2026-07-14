// Package common provides shared helpers for admin route modules.
package common

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// ===================================================================
// Parsing helpers
// ===================================================================

// ParseID parses an int64 from a form value string.
func ParseID(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}
	var id int64
	_, err := fmt.Sscanf(s, "%d", &id)
	return id, err
}

// ParseFloat parses a float64 from a form value string.
func ParseFloat(s string) (float64, error) {
	if s == "" {
		return 0, nil
	}
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

// StatusLabelText returns a plain text status label.
func StatusLabelText(active bool) string {
	if active {
		return "启用"
	}
	return "停用"
}

// ===================================================================
// Response helpers
// ===================================================================

// HtmlOK sets the Content-Type header for HTML responses.
func HtmlOK(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

// Redirect sends an HTMX HX-Redirect response.
func Redirect(w http.ResponseWriter, url string) {
	w.Header().Set("HX-Redirect", url)
	w.WriteHeader(200)
}

// ===================================================================
// Modal form helpers (BDL 1.0 modal overlay)
// ===================================================================

// ModalStart opens a modal overlay div.
func ModalStart(title string) string {
	return `<div class="modal-overlay" onclick="if(event.target===this) closeI56Modal()"><div class="modal-content"><div class="modal-header"><span class="modal-title">` + title + `</span><button class="modal-close" onclick="closeI56Modal()">&times;</button></div><div class="modal-body">`
}

// ModalEnd closes a modal overlay div.
func ModalEnd() string { return `</div></div></div>` }

// FormField renders a form field with label and input.
func FormField(label, name, value, placeholder string) string {
	return fmt.Sprintf(`<div class="form-group"><label class="form-label">%s</label><input name="%s" value="%s" class="form-input" placeholder="%s"></div>`, label, name, value, placeholder)
}

// FormSelect renders a form field with a select dropdown.
func FormSelect(label, name, value string, opts ...[2]string) string {
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

// FormSave opens a form element with HTMX post action.
func FormSave(action string) string {
	return fmt.Sprintf(`<form action="%s" method="POST" hx-post="%s" hx-target="#main-content" hx-swap="innerHTML">`, action, action)
}

// FormFooter closes a modal form with cancel and submit buttons.
func FormFooter() string {
	return `<div class="modal-footer"><button type="button" class="i56-btn" onclick="closeI56Modal()">取消</button><button type="submit" class="i56-btn i56-btn-primary">保存</button></div></form>`
}

// ===================================================================
// Rendering context
// ===================================================================

// ExecTpl is a function type for executing templates.
type ExecTpl func(tmpl map[string]*template.Template, key string, w http.ResponseWriter, name string, data any)

// GenericListFunc renders a page through the BDL generic_list template.
type GenericListFunc func(w http.ResponseWriter, page, title string, total int, cols []string, rows [][]string, addURL ...string)

// RenderCtx bundles template rendering helpers.
type RenderCtx struct {
	Tmpl  map[string]*template.Template
	Exec  ExecTpl
}

// DefaultExecTpl is the standard template execution helper.
func DefaultExecTpl(tmpl map[string]*template.Template, key string, w http.ResponseWriter, name string, data any) {
	tmpl[key].ExecuteTemplate(w, name, data)
}

// formatCell fixes common Sprintf format errors in cell values.
func formatCell(v string) string {
	v = strings.ReplaceAll(v, "%!(EXTRA int64=", "")
	v = strings.ReplaceAll(v, "%!(EXTRA float64=", "")
	// Clean up Sprintf artifacts like ") " that become ")"
	for strings.Contains(v, "))") {
		v = strings.ReplaceAll(v, "))", ")")
	}
	// Only strip the trailing ")" if it's a Sprintf artifact (i.e., there's no matching opening paren)
	// Don't strip closing parens that are part of legitimate parenthetical text (e.g., "厦门→台湾(海快)")
	v = strings.ReplaceAll(v, "sea_express", "海快")
	v = strings.ReplaceAll(v, "sea", "海运")
	v = strings.ReplaceAll(v, "air", "空运")
	// Cargo type Chinese labels
	v = strings.ReplaceAll(v, "general", "普货")
	v = strings.ReplaceAll(v, "sensitive", "特货")
	v = strings.ReplaceAll(v, "dangerous", "危险品")
	// Device type Chinese labels
	v = strings.ReplaceAll(v, "scale", "地磅")
	v = strings.ReplaceAll(v, "conveyor", "入库机")
	v = strings.ReplaceAll(v, "scanner", "扫码枪")
	return v
}

// ===================================================================
// RenderAdminPage — wraps content in the BDL light-theme admin layout
// with sidebar (matching templates/base.html) so all admin pages are
// visually consistent.  Use this for any standalone admin page instead
// of writing raw <html> inline.
// ===================================================================
func RenderAdminPage(w http.ResponseWriter, title, breadcrumb, content string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<!DOCTYPE html>
<html lang="zh-CN" data-theme="light">
<head>
<meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>`+title+` - I56</title>
<link rel="stylesheet" href="/static/css/i56-bdl.css">
<script src="/static/js/i56-theme.js"></script>
<script src="https://unpkg.com/htmx.org@1.9.10"></script>
</head>
<body>
<div class="app-layout">
<aside class="app-sidebar">
  <div class="sidebar-header">
    <div class="sidebar-logo">I</div>
    <a href="/admin" class="sidebar-title">I56 Framework</a>
  </div>
  <nav class="sidebar-nav">
    <a href="/admin" class="nav-item">
      <span class="nav-icon">🏠</span> 首页
    </a>
    <a href="/admin/warehouse-board" class="nav-item">
      <span class="nav-icon">📊</span> 仓库看板
    </a>
    <div class="nav-group" data-page="oms">
      <div class="nav-group-header" onclick="toggleNavGroup(this)">
        <span class="nav-group-header-label">📋 订单管理</span>
        <span class="nav-chevron">▾</span>
      </div>
      <div class="nav-group-body">
        <a href="/admin/orders" class="nav-item nav-sub-item">集运订单</a>
        <a href="/admin/service-orders" class="nav-item nav-sub-item">附加服务订单</a>
      </div>
    </div>
    <div class="nav-group" data-page="wms">
      <div class="nav-group-header" onclick="toggleNavGroup(this)">
        <span class="nav-group-header-label">🏗️ 仓库管理</span>
        <span class="nav-chevron">▾</span>
      </div>
      <div class="nav-group-body">
        <a href="/admin/parcels" class="nav-item nav-sub-item">包裹列表</a>
        <a href="/admin/service-workorders" class="nav-item nav-sub-item">附加服务工单</a>
        <a href="/admin/warehouses" class="nav-item nav-sub-item">仓库列表</a>
        <a href="/admin/inbound-board" class="nav-item nav-sub-item">入库看板</a>
        <a href="/admin/warehouse-console" class="nav-item nav-sub-item">仓库作业台</a>
        <a href="/admin/task-monitor" class="nav-item nav-sub-item">员工任务监控</a>
        <a href="/admin/exception-reports" class="nav-item nav-sub-item">异常记录</a>
      </div>
    </div>
    <div class="nav-group" data-page="fin">
      <div class="nav-group-header" onclick="toggleNavGroup(this)">
        <span class="nav-group-header-label">💰 财务报表</span>
        <span class="nav-chevron">▾</span>
      </div>
      <div class="nav-group-body">
        <a href="/admin/report/order-profit" class="nav-item nav-sub-item">集运订单盈利</a>
        <a href="/admin/report/service-profit" class="nav-item nav-sub-item">附加服务盈利</a>
        <a href="/admin/report/client-profit" class="nav-item nav-sub-item">客户盈利</a>
        <a href="/admin/report/route-profit" class="nav-item nav-sub-item">路线盈利</a>
      </div>
    </div>
    <div class="nav-group" data-page="tms">
      <div class="nav-group-header" onclick="toggleNavGroup(this)">
        <span class="nav-group-header-label">🚛 物流管理</span>
        <span class="nav-chevron">▾</span>
      </div>
      <div class="nav-group-body">
        <a href="/admin/area-groups" class="nav-item nav-sub-item">区域组管理</a>
        <a href="/admin/carriers" class="nav-item nav-sub-item">承运商列表</a>
        <a href="/admin/couriers" class="nav-item nav-sub-item">快递公司</a>
        <a href="/admin/customs-brokers" class="nav-item nav-sub-item">清关公司</a>
        <a href="/admin/route-templates" class="nav-item nav-sub-item">线路模板</a>
        <a href="/admin/shipping-providers" class="nav-item nav-sub-item">运输公司</a>
      </div>
    </div>
    <div class="nav-group" data-page="crm">
      <div class="nav-group-header" onclick="toggleNavGroup(this)">
        <span class="nav-group-header-label">👤 客户管理</span>
        <span class="nav-chevron">▾</span>
      </div>
      <div class="nav-group-body">
        <a href="/admin/clients" class="nav-item nav-sub-item">客户管理</a>
        <a href="/admin/client-accounts" class="nav-item nav-sub-item">客户账号</a>
        <a href="/admin/client-members" class="nav-item nav-sub-item">客户会员</a>
        <a href="/admin/client-recharge" class="nav-item nav-sub-item">客户充值</a>
        <a href="/admin/balance-logs" class="nav-item nav-sub-item">余额日志</a>
        <a href="/admin/monthly-statements" class="nav-item nav-sub-item">月结对账单</a>
      </div>
    </div>
    <div class="nav-group" data-page="sys">
      <div class="nav-group-header" onclick="toggleNavGroup(this)">
        <span class="nav-group-header-label">⚙️ 系统</span>
        <span class="nav-chevron">▾</span>
      </div>
      <div class="nav-group-body">
        <a href="/admin/roles" class="nav-item nav-sub-item">角色管理</a>
        <a href="/admin/employees" class="nav-item nav-sub-item">员工管理</a>
        <a href="/admin/system/api-couriers" class="nav-item nav-sub-item">物流API对接</a>
        <a href="/admin/system/api-customs" class="nav-item nav-sub-item">清关API对接</a>
        <a href="/admin/system/api-notifications" class="nav-item nav-sub-item">通知渠道配置</a>
        <a href="/admin/system/api-storage" class="nav-item nav-sub-item">存储配置</a>
        <a href="/admin/system/api-printers" class="nav-item nav-sub-item">打印机设置</a>
        <a href="/admin/system/api-devices" class="nav-item nav-sub-item">🔌 设备管理</a>
        <div class="nav-sub-divider"></div>
        <a href="/admin/system/api-ezway" class="nav-item nav-sub-item">🏛️ EZ Way实名认证</a>
        <div class="nav-sub-divider"></div>
        <a href="/admin/system/scheduler" class="nav-item nav-sub-item">⏰ 定时任务</a>
        <a href="/admin/system/audit-logs" class="nav-item nav-sub-item">📋 审计日志</a>
        <a href="/admin/system/reports" class="nav-item nav-sub-item">📊 内置报表</a>
      </div>
    </div>
    <div class="sidebar-footer">
      <span class="version">v2.0 LTS</span>
    </div>
  </nav>
</aside>
<div class="app-main">
  <header class="app-header">
    <div class="header-breadcrumb">
      I56 Admin <span style="color:var(--i56-text-muted);margin:0 4px">/</span> <span>`+breadcrumb+`</span>
    </div>
    <div class="header-actions">
      <button class="i56-btn i56-btn-sm i56-btn-ghost" onclick="I56Theme.toggle()" title="切换主题">🌓</button>
    </div>
  </header>
  <main class="app-content">
`+content+`
  </main>
</div>
</div>
<script>
function toggleNavGroup(header) {
  var group = header.parentElement;
  var wasOpen = group.classList.contains('open');
  document.querySelectorAll('.nav-group.open').forEach(function(g) {
    if (g !== group) g.classList.remove('open');
  });
  if (wasOpen) { group.classList.remove('open'); }
  else { group.classList.add('open'); }
}
(function() {
  var currentUrl = window.location.pathname;
  var subItems = document.querySelectorAll('.nav-sub-item');
  var matched = false;
  subItems.forEach(function(item) {
    var href = item.getAttribute('href');
    if (href && currentUrl.startsWith(href.split('?')[0]) && href.split('?')[0] !== '/admin') {
      item.classList.add('active');
      var group = item.closest('.nav-group');
      if (group) group.classList.add('open');
      matched = true;
    }
  });
  if (!matched) {
    var groups = document.querySelectorAll('.nav-group[data-page]');
    groups.forEach(function(group) {
      var pages = {
        'oms': ['orders', 'service-orders'],
        'wms': ['parcels', 'service-workorders', 'warehouses', 'inbound-board', 'warehouse-console', 'task-monitor', 'exception-reports'],
        'fin': ['report/order-profit', 'report/service-profit', 'report/client-profit', 'report/route-profit'],
        'tms': ['area-groups', 'carriers', 'couriers', 'customs-brokers', 'route-templates', 'shipping-providers'],
        'crm': ['clients', 'client-accounts', 'client-members', 'client-recharge', 'balance-logs', 'monthly-statements'],
        'sys': ['roles', 'employees', 'system/api-couriers', 'system/api-customs', 'system/api-ezway', 'system/api-notifications', 'system/api-storage', 'system/api-printers', 'system/api-devices', 'system/scheduler', 'system/audit-logs', 'system/reports']
      }[group.getAttribute('data-page')] || [];
      for (var i = 0; i < pages.length; i++) {
        if (currentUrl.indexOf('/admin/' + pages[i]) === 0) {
          group.classList.add('open'); matched = true; break;
        }
      }
    });
  }
  var topItems = document.querySelectorAll('a.nav-item:not(.nav-sub-item)');
  topItems.forEach(function(item) {
    var href = item.getAttribute('href');
    if (href && currentUrl === href.split('?')[0]) {
      item.classList.add('active');
    } else if (href === '/admin' && currentUrl !== '/admin') {
      item.classList.remove('active');
    }
  });
})();
</script>
</body>
</html>`)
}

// NewGenericList creates a genericList closure using the render context.
func (rc *RenderCtx) NewGenericList() GenericListFunc {
	return func(w http.ResponseWriter, page, title string, total int, cols []string, rows [][]string, addURL ...string) {
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
		if len(addURL) > 0 && addURL[0] != "" {
			data["AddURL"] = addURL[0]
		}
		rc.Exec(rc.Tmpl, "generic_list", w, "generic_list.html", data)
	}
}

// FormImageUpload renders a file input for image upload.
// Files are auto-uploaded via JS fetch to /admin/upload/parcel-image.
// URLs are stored in a hidden field "uploaded_urls" for form submission.
func FormImageUpload(label string) string {
	return fmt.Sprintf(`<div class="form-group">
		<label class="form-label">%s</label>
		<input type="file" accept="image/*" multiple
			style="display:block;font-size:12px;margin-bottom:4px"
			onchange="i56UploadImages(this)" />
		<div class="i56-preview-row" style="display:flex;gap:4px;flex-wrap:wrap;margin-top:4px"></div>
		<input type="hidden" name="uploaded_urls" value="" />
	</div>
	<script>
	if(!window._i56UploadInit){
		window._i56UploadInit=true;
		async function i56UploadImages(el){
			var files=el.files;if(!files.length)return;
			var fd=new FormData();
			for(var f of files)fd.append("images",f);
			var r=await fetch("/admin/upload/parcel-image",{method:"POST",body:fd});
			var d=await r.json();
			if(d.ok&&d.urls){
				var h=el.parentElement.querySelector("[name=uploaded_urls]");
				if(h.value)h.value+=",";
				h.value+=d.urls.join(",");
				var pv=el.parentElement.querySelector(".i56-preview-row");
				for(var u of d.urls){
					var img=document.createElement("img");
					img.src=u;img.style="width:64px;height:64px;object-fit:cover;border-radius:4px;border:1px solid #ddd";
					pv.appendChild(img);
				}
			}
		}
	}
	</script>`, label)
}
