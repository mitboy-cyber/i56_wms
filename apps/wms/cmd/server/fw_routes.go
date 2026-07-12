package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/i56/framework/core/audit"
	"github.com/i56/framework/core/report"
	"github.com/i56/framework/core/router"
	"github.com/i56/framework/core/scheduler"

	"github.com/i56/i56-apps/i56-wms/internal/common"
)

// ─── Scheduler Routes ───

func registerSchedulerRoutes(r *router.Router, sch *scheduler.Scheduler, a func(http.HandlerFunc) http.HandlerFunc) {
	// Admin page: scheduler status
	r.GET("/admin/system/scheduler", a(func(w http.ResponseWriter, req *http.Request) {
		jobs := sch.ListJobs()
		var sb strings.Builder
		sb.WriteString(`<h1 style="margin-bottom:16px">⏰ 定时任务调度器</h1>
<table class="data-table" style="width:100%"><thead><tr><th>任务名称</th><th>Cron表达式</th><th>上次运行</th><th>下次运行</th><th>运行次数</th><th>状态</th><th>操作</th></tr></thead><tbody>`)
		for _, j := range jobs {
			lastRun := "-"
			if !j.LastRun.IsZero() {
				lastRun = j.LastRun.Format("2006-01-02 15:04:05")
			}
			nextRun := "-"
			if !j.NextRun.IsZero() {
				nextRun = j.NextRun.Format("2006-01-02 15:04:05")
			}
			status := "空闲"
			statusClass := "badge badge-success"
			if j.Running {
				status = "运行中"
				statusClass = "badge badge-primary"
			}
			if j.LastError != "" {
				status = "错误"
				statusClass = "badge badge-danger"
			}
			fmt.Fprintf(&sb, `<tr>
<td>%s</td><td><code>%s</code></td><td>%s</td><td>%s</td><td>%d</td>
<td><span class="%s">%s</span></td>
<td><button class="btn btn-sm btn-primary" onclick="triggerJob('%s')">立即执行</button></td>
</tr>`, j.Name, j.CronExpr, lastRun, nextRun, j.RunCount, statusClass, status, j.Name)
		}
		sb.WriteString(`</tbody></table>
<script>async function triggerJob(name){const r=await fetch('/api/system/scheduler/trigger?name='+encodeURIComponent(name),{method:'POST'});location.reload()}</script>
<div style="margin-top:24px;padding:16px;background:var(--i56-bg-surface);border:1px solid var(--i56-border);border-radius:var(--i56-radius-md)">
<h2 style="font-size:14px;margin-bottom:12px">内置任务</h2>
<ul style="margin:8px 0;padding-left:20px;font-size:13px;color:var(--i56-text-secondary);line-height:1.8">
<li><b>bill-generation</b> — @daily — 每日账单生成</li>
<li><b>weight-cleanup</b> — @every 6h — 清理旧重量记录</li>
<li><b>statistics-report</b> — @every 1h — 生成统计报表</li>
<li><b>backup-database</b> — 0 2 * * * — 数据库备份 (凌晨2点)</li>
<li><b>health-check</b> — @every 5m — 设备健康检查</li>
</ul></div>`)
		common.RenderAdminPage(w, "定时任务调度器", "系统 / 定时任务调度器", sb.String())
	}))

	// API: trigger a job
	r.POST("/api/system/scheduler/trigger", a(func(w http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "missing job name", 400)
			return
		}
		if err := sch.TriggerNow(name); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "triggered", "job": name})
	}))

	// API: list jobs
	r.GET("/api/system/scheduler/jobs", a(func(w http.ResponseWriter, req *http.Request) {
		jobs := sch.ListJobs()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": jobs, "count": len(jobs)})
	}))
}

// ─── Audit Routes ───

func registerAuditRoutes(r *router.Router, auditLogger *audit.AuditLogger, a func(http.HandlerFunc) http.HandlerFunc) {
	// Admin page: audit logs
	r.GET("/admin/system/audit-logs", a(func(w http.ResponseWriter, req *http.Request) {
		actionFilter := req.URL.Query().Get("action")
		resourceFilter := req.URL.Query().Get("resource")
		fromStr := req.URL.Query().Get("from")
		toStr := req.URL.Query().Get("to")
		pageStr := req.URL.Query().Get("page")
		page, _ := strconv.Atoi(pageStr)
		if page < 1 { page = 1 }

		filter := audit.AuditFilter{
			Action:   actionFilter,
			Resource: resourceFilter,
			Limit:    50,
			Offset:   (page - 1) * 50,
		}
		if fromStr != "" {
			if t, err := time.Parse("2006-01-02", fromStr); err == nil {
				filter.From = t
			}
		}
		if toStr != "" {
			if t, err := time.Parse("2006-01-02", toStr); err == nil {
				filter.To = t.Add(24 * time.Hour)
			}
		}

		entries, total, _ := auditLogger.Query(req.Context(), filter)
		var sb strings.Builder
		sb.WriteString(`<h1 style="margin-bottom:16px">📋 操作审计日志</h1>
<form method="GET" style="display:flex;gap:12px;margin-bottom:16px;flex-wrap:wrap;align-items:flex-end">
<div><label style="font-size:13px;color:var(--i56-text-secondary)">操作类型</label>
<select name="action" class="form-select" style="min-width:120px">
<option value="">全部</option>`)
		for _, act := range []string{"CREATE", "UPDATE", "DELETE", "LOGIN", "LOGOUT", "EXPORT", "IMPORT", "VIEW"} {
			sel := ""
			if act == actionFilter { sel = " selected" }
			fmt.Fprintf(&sb, `<option value="%s"%s>%s</option>`, act, sel, act)
		}
		sb.WriteString(`</select></div>
<div><label style="font-size:13px;color:var(--i56-text-secondary)">资源类型</label>
<select name="resource" class="form-select" style="min-width:120px">
<option value="">全部</option>`)
		for _, res := range []string{"order", "parcel", "client", "warehouse", "user", "role"} {
			sel := ""
			if res == resourceFilter { sel = " selected" }
			fmt.Fprintf(&sb, `<option value="%s"%s>%s</option>`, res, sel, res)
		}
		fmt.Fprintf(&sb, `</select></div>
<div><label style="font-size:13px;color:var(--i56-text-secondary)">开始日期</label><input type="date" name="from" value="%s" class="form-input"></div>
<div><label style="font-size:13px;color:var(--i56-text-secondary)">结束日期</label><input type="date" name="to" value="%s" class="form-input"></div>
<button type="submit" class="btn btn-primary">筛选</button>
</form>`, fromStr, toStr)

		fmt.Fprintf(&sb, `<p style="color:var(--i56-text-secondary);margin-bottom:12px">共 %d 条记录</p>`, total)
		sb.WriteString(`<table class="data-table" style="width:100%"><thead><tr><th>时间</th><th>用户</th><th>操作</th><th>资源</th><th>资源ID</th><th>详情</th><th>IP</th></tr></thead><tbody>`)
		for _, e := range entries {
			detailPreview := e.Detail
			if len(detailPreview) > 50 {
				detailPreview = detailPreview[:50] + "..."
			}
			fmt.Fprintf(&sb, `<tr>
<td>%s</td><td>%d</td><td><span class="badge badge-%s">%s</span></td>
<td>%s</td><td>%s</td><td style="max-width:200px;overflow:hidden;text-overflow:ellipsis" title="%s">%s</td><td>%s</td>
</tr>`,
				e.CreatedAt.Format("2006-01-02 15:04:05"), e.UserID,
				actionBadgeClass(e.Action), e.Action,
				e.Resource, e.ResourceID, e.Detail, detailPreview, e.IP)
		}
		if len(entries) == 0 {
			sb.WriteString(`<tr><td colspan="7" style="text-align:center;color:var(--i56-text-muted)">暂无审计日志</td></tr>`)
		}
		sb.WriteString(`</tbody></table>`)
		common.RenderAdminPage(w, "操作审计日志", "系统 / 操作审计日志", sb.String())
	}))

	// API: list audit entries
	r.GET("/api/system/audit-logs", a(func(w http.ResponseWriter, req *http.Request) {
		entries, total, _ := auditLogger.Query(req.Context(), audit.AuditFilter{Limit: 50})
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": entries, "total": total})
	}))
}

func actionBadgeClass(action string) string {
	switch action {
	case "CREATE":
		return "success"
	case "UPDATE":
		return "brand"
	case "DELETE":
		return "danger"
	case "LOGIN":
		return "info"
	case "EXPORT":
		return "warning"
	default:
		return "secondary"
	}
}

// ─── Report Routes ───

func registerReportRoutes(r *router.Router, engine *report.BuiltinEngine, a func(http.HandlerFunc) http.HandlerFunc) {
	// Admin page: report list
	r.GET("/admin/system/reports", a(func(w http.ResponseWriter, req *http.Request) {
		reports := engine.ListReports()
		var sb strings.Builder
		sb.WriteString(`<h1 style="margin-bottom:16px">📊 内置报表</h1>
<table class="data-table" style="width:100%"><thead><tr><th>报表名称</th><th>描述</th><th>显示方式</th><th>参数</th><th>操作</th></tr></thead><tbody>`)
		for _, rpt := range reports {
			paramCount := len(rpt.Params)
			displayIcon := "📋"
			if rpt.Display == "chart" {
				displayIcon = "📈"
			}
			displayLabel := "表格"
			if rpt.Display == "chart" {
				displayLabel = "图表"
			}
			fmt.Fprintf(&sb, `<tr>
<td><b>%s %s</b></td><td style="font-size:12px;color:var(--i56-text-muted)"><code>%s</code></td>
<td>%s</td><td>%d 个参数</td>
<td><a href="/admin/system/reports/view?name=%s" class="btn btn-sm btn-primary">查看</a></td>
</tr>`, displayIcon, rpt.Title, rpt.Query, displayLabel, paramCount, rpt.Name)
		}
		sb.WriteString(`</tbody></table>`)
		common.RenderAdminPage(w, "内置报表", "系统 / 内置报表", sb.String())
	}))

	// View report
	r.GET("/admin/system/reports/view", a(func(w http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("name")
		if name == "" {
			http.Redirect(w, req, "/admin/system/reports", 303)
			return
		}
		result, err := engine.Execute(req.Context(), name, nil)
		if err != nil {
			common.RenderAdminPage(w, "报表错误", "系统 / 内置报表 / 错误", fmt.Sprintf(`<div class="data-table-wrapper" style="padding:32px;text-align:center"><p style="color:var(--i56-error);font-size:14px">报表执行失败: %s</p><a href="/admin/system/reports" class="btn btn-sm btn-primary mt-4">← 返回报表列表</a></div>`, err.Error()))
			return
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, `<h1 style="margin-bottom:16px">%s</h1>
<p style="color:var(--i56-text-secondary);margin-bottom:16px">共 %d 行数据</p>
<table class="data-table" style="width:100%"><thead><tr>`, result.Title, result.Total)
		for _, col := range result.Columns {
			fmt.Fprintf(&sb, "<th>%s</th>", col)
		}
		sb.WriteString("</tr></thead><tbody>")
		for _, row := range result.Rows {
			sb.WriteString("<tr>")
			for _, cell := range row {
				fmt.Fprintf(&sb, "<td>%v</td>", cell)
			}
			sb.WriteString("</tr>")
		}
		sb.WriteString(`</tbody></table>`)
		common.RenderAdminPage(w, result.Title, "系统 / 内置报表 / "+result.Title, sb.String())
	}))

	// API: execute report
	r.GET("/api/system/reports/execute", a(func(w http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "missing report name", 400)
			return
		}
		result, err := engine.Execute(req.Context(), name, nil)
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": result,
		})
	}))
}

// ─── OpenAPI Demo Routes ───

func registerOpenAPIDemoRoutes(r *router.Router, gen *router.OpenAPIGenerator, a func(http.HandlerFunc) http.HandlerFunc) {
	// Register some well-documented routes for OpenAPI
	gen.RegisterRoute(router.RouteDoc{
		Method:  "GET",
		Path:    "/api/v1/health",
		Summary: "Health check endpoint",
		Tags:    []string{"System"},
		Params:  []router.ParamDoc{},
	})
	gen.RegisterRoute(router.RouteDoc{
		Method:  "GET",
		Path:    "/api/v1/parcels",
		Summary: "List all parcels",
		Tags:    []string{"Parcels"},
		Params: []router.ParamDoc{
			{Name: "status", In: "query", Required: false, Description: "Filter by status", Type: "string"},
			{Name: "page", In: "query", Required: false, Description: "Page number", Type: "integer"},
			{Name: "limit", In: "query", Required: false, Description: "Items per page", Type: "integer"},
		},
	})
	gen.RegisterRoute(router.RouteDoc{
		Method:  "GET",
		Path:    "/api/v1/parcels/{tracking_number}",
		Summary: "Get parcel by tracking number",
		Tags:    []string{"Parcels"},
		Params: []router.ParamDoc{
			{Name: "tracking_number", In: "path", Required: true, Description: "Tracking number", Type: "string"},
		},
	})
	gen.RegisterRoute(router.RouteDoc{
		Method:  "GET",
		Path:    "/api/v1/orders",
		Summary: "List all orders",
		Tags:    []string{"Orders"},
		Params: []router.ParamDoc{
			{Name: "status", In: "query", Required: false, Description: "Filter by status", Type: "string"},
			{Name: "page", In: "query", Required: false, Description: "Page number", Type: "integer"},
		},
	})
	gen.RegisterRoute(router.RouteDoc{
		Method:  "POST",
		Path:    "/api/v1/parcels",
		Summary: "Create a new parcel (pre-declare)",
		Tags:    []string{"Parcels"},
		Params:  []router.ParamDoc{},
	})
	gen.RegisterRoute(router.RouteDoc{
		Method:  "GET",
		Path:    "/api/v1/warehouses",
		Summary: "List all warehouses",
		Tags:    []string{"Warehouses"},
		Params:  []router.ParamDoc{},
	})
	gen.RegisterRoute(router.RouteDoc{
		Method:  "GET",
		Path:    "/api/v1/clients",
		Summary: "List all clients",
		Tags:    []string{"Clients"},
		Params:  []router.ParamDoc{},
	})
	gen.RegisterRoute(router.RouteDoc{
		Method:  "GET",
		Path:    "/api/system/scheduler/jobs",
		Summary: "List all scheduled jobs",
		Tags:    []string{"System", "Scheduler"},
		Params:  []router.ParamDoc{},
	})
	gen.RegisterRoute(router.RouteDoc{
		Method:  "POST",
		Path:    "/api/system/scheduler/trigger",
		Summary: "Trigger a scheduled job immediately",
		Tags:    []string{"System", "Scheduler"},
		Params: []router.ParamDoc{
			{Name: "name", In: "query", Required: true, Description: "Job name to trigger", Type: "string"},
		},
	})
	gen.RegisterRoute(router.RouteDoc{
		Method:  "GET",
		Path:    "/api/system/audit-logs",
		Summary: "List audit log entries",
		Tags:    []string{"System", "Audit"},
		Params: []router.ParamDoc{
			{Name: "action", In: "query", Required: false, Description: "Filter by action", Type: "string"},
			{Name: "resource", In: "query", Required: false, Description: "Filter by resource type", Type: "string"},
		},
	})
	gen.RegisterRoute(router.RouteDoc{
		Method:  "GET",
		Path:    "/api/system/reports/execute",
		Summary: "Execute a built-in report",
		Tags:    []string{"System", "Reports"},
		Params: []router.ParamDoc{
			{Name: "name", In: "query", Required: true, Description: "Report name", Type: "string"},
		},
	})
	gen.RegisterRoute(router.RouteDoc{
		Method:  "POST",
		Path:    "/api/ai/chat",
		Summary: "Chat with AI assistant",
		Tags:    []string{"AI"},
		Params:  []router.ParamDoc{},
	})
}
