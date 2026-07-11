// Package route provides FIN (财务报表) admin route registration.
package route

import (
	"fmt"
	"net/http"

	"github.com/i56/framework/core/router"

	"github.com/i56/i56-apps/i56-wms/internal/common"

	reportDomain "github.com/i56/modules/report/domain"
)

// Register FIN admin routes (~4 report pages).
func Register(
	r *router.Router,
	a func(http.HandlerFunc) http.HandlerFunc,
	rc *common.RenderCtx,
	rpt *reportDomain.ReportService,
) {
	_ = rpt
	gp := rc.NewGenericList()

	// /admin/report/order-profit
	r.GET("/admin/report/order-profit", a(func(w http.ResponseWriter, req *http.Request) {
		rows := [][]string{
			{"厦门→台湾(空运)", "ORD-001", "EZ集运通", "王仁照", "空运", "¥150.00", "¥85.00", "¥65.00", "2026-07-10 10:30"},
			{"厦门→台湾(海快)", "ORD-002", "EZ集运通", "吴欣如", "海快", "¥42.00", "¥25.00", "¥17.00", "2026-07-09 14:20"},
		}
		gp(w, "fin_order_profit", "订单利润报表", len(rows), []string{"线路", "订单号", "客户", "收件人", "运输方式", "收入", "成本", "利润", "时间"}, rows, "/admin/report/order-profit/add-form")
	}))
	r.GET("/admin/report/order-profit/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增订单利润记录")+common.FormSave("/admin/report/order-profit/save")+
			common.FormField("线路", "route", "", "如: 厦门→台湾(空运)")+
			common.FormField("订单号", "order_no", "", "如: ORD-003")+
			common.FormField("客户", "client", "", "客户名称")+
			common.FormField("收件人", "recipient", "", "收件人名称")+
			common.FormField("运输方式", "transport", "", "空运/海快/海运")+
			common.FormField("收入", "revenue", "", "如: 150.00")+
			common.FormField("成本", "cost", "", "如: 85.00")+
			common.FormField("利润", "profit", "", "如: 65.00")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/report/order-profit/save", a(func(w http.ResponseWriter, req *http.Request) {
		common.Redirect(w, "/admin/report/order-profit")
	}))
	r.GET("/admin/report/order-profit/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		route := req.URL.Query().Get("id")
		if route == "" { http.Error(w, "invalid id", 400); return }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑订单利润记录")+common.FormSave("/admin/report/order-profit/update")+
			common.FormField("线路", "route", route, "如: 厦门→台湾(空运)")+
			common.FormField("订单号", "order_no", "", "如: ORD-003")+
			common.FormField("客户", "client", "", "客户名称")+
			common.FormField("收件人", "recipient", "", "收件人名称")+
			common.FormField("运输方式", "transport", "", "空运/海快/海运")+
			common.FormField("收入", "revenue", "", "如: 150.00")+
			common.FormField("成本", "cost", "", "如: 85.00")+
			common.FormField("利润", "profit", "", "如: 65.00")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/report/order-profit/update", a(func(w http.ResponseWriter, req *http.Request) {
		common.Redirect(w, "/admin/report/order-profit")
	}))

	// /admin/report/service-profit
	r.GET("/admin/report/service-profit", a(func(w http.ResponseWriter, req *http.Request) {
		rows := [][]string{
			{"包装服务", "15", "¥150.00", "¥45.00", "¥105.00", "2026-07"},
			{"拍照服务", "8", "¥40.00", "¥10.00", "¥30.00", "2026-07"},
		}
		gp(w, "fin_service_profit", "服务利润报表", len(rows), []string{"服务类型", "数量", "收入", "成本", "利润", "月份"}, rows, "/admin/report/service-profit/add-form")
	}))
	r.GET("/admin/report/service-profit/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增服务利润记录")+common.FormSave("/admin/report/service-profit/save")+
			common.FormField("服务类型", "type", "", "如: 包装服务")+
			common.FormField("数量", "count", "", "如: 15")+
			common.FormField("收入", "revenue", "", "如: 150.00")+
			common.FormField("成本", "cost", "", "如: 45.00")+
			common.FormField("利润", "profit", "", "如: 105.00")+
			common.FormField("月份", "month", "", "如: 2026-07")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/report/service-profit/save", a(func(w http.ResponseWriter, req *http.Request) {
		common.Redirect(w, "/admin/report/service-profit")
	}))
	r.GET("/admin/report/service-profit/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		stype := req.URL.Query().Get("id")
		if stype == "" { http.Error(w, "invalid id", 400); return }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑服务利润记录")+common.FormSave("/admin/report/service-profit/update")+
			common.FormField("服务类型", "type", stype, "如: 包装服务")+
			common.FormField("数量", "count", "", "如: 15")+
			common.FormField("收入", "revenue", "", "如: 150.00")+
			common.FormField("成本", "cost", "", "如: 45.00")+
			common.FormField("利润", "profit", "", "如: 105.00")+
			common.FormField("月份", "month", "", "如: 2026-07")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/report/service-profit/update", a(func(w http.ResponseWriter, req *http.Request) {
		common.Redirect(w, "/admin/report/service-profit")
	}))

	// /admin/report/client-profit
	r.GET("/admin/report/client-profit", a(func(w http.ResponseWriter, req *http.Request) {
		rows := [][]string{
			{"EZ集运通", "25", "¥2,500.00", "¥1,500.00", "¥1,000.00", "¥5,000.00", "2026-07"},
		}
		gp(w, "fin_client_profit", "客户利润报表", len(rows), []string{"客户", "订单数", "收入", "成本", "利润", "余额", "月份"}, rows, "/admin/report/client-profit/add-form")
	}))
	r.GET("/admin/report/client-profit/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增客户利润记录")+common.FormSave("/admin/report/client-profit/save")+
			common.FormField("客户", "client", "", "客户名称")+
			common.FormField("订单数", "order_count", "", "如: 25")+
			common.FormField("收入", "revenue", "", "如: 2500.00")+
			common.FormField("成本", "cost", "", "如: 1500.00")+
			common.FormField("利润", "profit", "", "如: 1000.00")+
			common.FormField("余额", "balance", "", "如: 5000.00")+
			common.FormField("月份", "month", "", "如: 2026-07")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/report/client-profit/save", a(func(w http.ResponseWriter, req *http.Request) {
		common.Redirect(w, "/admin/report/client-profit")
	}))
	r.GET("/admin/report/client-profit/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		client := req.URL.Query().Get("id")
		if client == "" { http.Error(w, "invalid id", 400); return }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑客户利润记录")+common.FormSave("/admin/report/client-profit/update")+
			common.FormField("客户", "client", client, "客户名称")+
			common.FormField("订单数", "order_count", "", "如: 25")+
			common.FormField("收入", "revenue", "", "如: 2500.00")+
			common.FormField("成本", "cost", "", "如: 1500.00")+
			common.FormField("利润", "profit", "", "如: 1000.00")+
			common.FormField("余额", "balance", "", "如: 5000.00")+
			common.FormField("月份", "month", "", "如: 2026-07")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/report/client-profit/update", a(func(w http.ResponseWriter, req *http.Request) {
		common.Redirect(w, "/admin/report/client-profit")
	}))

	// /admin/report/route-profit
	r.GET("/admin/report/route-profit", a(func(w http.ResponseWriter, req *http.Request) {
		rows := [][]string{
			{"厦门→台湾(空运)", "12", "¥1,800.00", "¥1,200.00", "¥600.00", "33.3%", "2026-07"},
			{"厦门→台湾(海快)", "18", "¥756.00", "¥450.00", "¥306.00", "40.5%", "2026-07"},
			{"厦门→台湾(海运)", "5", "¥160.00", "¥100.00", "¥60.00", "37.5%", "2026-07"},
		}
		gp(w, "fin_route_profit", "线路利润报表", len(rows), []string{"线路", "订单数", "收入", "成本", "利润", "利润率", "月份"}, rows, "/admin/report/route-profit/add-form")
	}))
	r.GET("/admin/report/route-profit/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增线路利润记录")+common.FormSave("/admin/report/route-profit/save")+
			common.FormField("线路", "route", "", "如: 厦门→台湾(空运)")+
			common.FormField("订单数", "order_count", "", "如: 12")+
			common.FormField("收入", "revenue", "", "如: 1800.00")+
			common.FormField("成本", "cost", "", "如: 1200.00")+
			common.FormField("利润", "profit", "", "如: 600.00")+
			common.FormField("利润率", "margin", "", "如: 33.3%")+
			common.FormField("月份", "month", "", "如: 2026-07")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/report/route-profit/save", a(func(w http.ResponseWriter, req *http.Request) {
		common.Redirect(w, "/admin/report/route-profit")
	}))
	r.GET("/admin/report/route-profit/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		route := req.URL.Query().Get("id")
		if route == "" { http.Error(w, "invalid id", 400); return }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑线路利润记录")+common.FormSave("/admin/report/route-profit/update")+
			common.FormField("线路", "route", route, "如: 厦门→台湾(空运)")+
			common.FormField("订单数", "order_count", "", "如: 12")+
			common.FormField("收入", "revenue", "", "如: 1800.00")+
			common.FormField("成本", "cost", "", "如: 1200.00")+
			common.FormField("利润", "profit", "", "如: 600.00")+
			common.FormField("利润率", "margin", "", "如: 33.3%")+
			common.FormField("月份", "month", "", "如: 2026-07")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/report/route-profit/update", a(func(w http.ResponseWriter, req *http.Request) {
		common.Redirect(w, "/admin/report/route-profit")
	}))

	_ = fmt.Sprintf
}
