package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/i56/framework/core/audit"
	"github.com/i56/framework/core/report"
	"github.com/i56/framework/core/router"
	"github.com/i56/framework/core/scheduler"
)

// ─── Scheduler Routes ───

func registerSchedulerRoutes(r *router.Router, sch *scheduler.Scheduler, a func(http.HandlerFunc) http.HandlerFunc, tmpl map[string]*template.Template) {
	// Admin page: scheduler status
	r.GET("/admin/system/scheduler", a(func(w http.ResponseWriter, req *http.Request) {
		jobs := sch.ListJobs()
		type jobView struct {
			Name, CronExpr, LastRun, NextRun, Status, StatusClass string
			RunCount                                                int64
		}
		views := make([]jobView, 0, len(jobs))
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
			views = append(views, jobView{
				Name: j.Name, CronExpr: j.CronExpr,
				LastRun: lastRun, NextRun: nextRun,
				RunCount: j.RunCount, Status: status, StatusClass: statusClass,
			})
		}
		tmpl["scheduler"].ExecuteTemplate(w, "scheduler.html", map[string]any{
			"Jobs":       views,
			"Active":     "scheduler",
			"Breadcrumb": "系统 / 定时任务调度器",
		})
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

func registerAuditRoutes(r *router.Router, auditLogger *audit.AuditLogger, a func(http.HandlerFunc) http.HandlerFunc, tmpl map[string]*template.Template) {
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

		type auditView struct {
			CreatedAt, UserID, Action, BadgeClass, Resource, ResourceID, Detail, DetailPreview, IP string
		}
		views := make([]auditView, 0, len(entries))
		for _, e := range entries {
			detailPreview := e.Detail
			if len(detailPreview) > 50 {
				detailPreview = detailPreview[:50] + "..."
			}
			views = append(views, auditView{
				CreatedAt:    e.CreatedAt.Format("2006-01-02 15:04:05"),
				UserID:       fmt.Sprintf("%d", e.UserID),
				Action:       e.Action,
				BadgeClass:   actionBadgeClass(e.Action),
				Resource:     e.Resource,
				ResourceID:   e.ResourceID,
				Detail:       e.Detail,
				DetailPreview: detailPreview,
				IP:           e.IP,
			})
		}

		tmpl["audit_logs"].ExecuteTemplate(w, "audit_logs.html", map[string]any{
			"Entries":         views,
			"Total":           total,
			"ActionFilter":    actionFilter,
			"ResourceFilter":  resourceFilter,
			"FromStr":         fromStr,
			"ToStr":           toStr,
			"ActionOptions":   []string{"CREATE", "UPDATE", "DELETE", "LOGIN", "LOGOUT", "EXPORT", "IMPORT", "VIEW"},
			"ResourceOptions": []string{"order", "parcel", "client", "warehouse", "user", "role"},
			"Breadcrumb":      "系统 / 操作审计日志",
		})
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

func registerReportRoutes(r *router.Router, engine *report.BuiltinEngine, a func(http.HandlerFunc) http.HandlerFunc, tmpl map[string]*template.Template) {
	// Admin page: report list
	r.GET("/admin/system/reports", a(func(w http.ResponseWriter, req *http.Request) {
		reports := engine.ListReports()
		type reportView struct {
			Name, Title, Query, DisplayIcon, DisplayLabel string
			ParamCount                                      int
		}
		views := make([]reportView, 0, len(reports))
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
			views = append(views, reportView{
				Name: rpt.Name, Title: rpt.Title, Query: rpt.Query,
				DisplayIcon: displayIcon, DisplayLabel: displayLabel,
				ParamCount: paramCount,
			})
		}
		tmpl["reports"].ExecuteTemplate(w, "reports.html", map[string]any{
			"Reports":    views,
			"Breadcrumb": "系统 / 内置报表",
		})
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
			tmpl["report_view"].ExecuteTemplate(w, "report_view.html", map[string]any{
				"Error":      err.Error(),
				"Breadcrumb": "系统 / 内置报表 / 错误",
			})
			return
		}
		tmpl["report_view"].ExecuteTemplate(w, "report_view.html", map[string]any{
			"Title":      result.Title,
			"Total":      result.Total,
			"Columns":    result.Columns,
			"Rows":       result.Rows,
			"Breadcrumb": "系统 / 内置报表 / " + result.Title,
		})
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
