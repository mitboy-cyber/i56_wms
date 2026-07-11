package main
import ("html/template";"net/http";"github.com/i56/framework/core/router";sysRepo "github.com/i56/modules/system/repository")

func adminSystemPages(tmpl map[string]*template.Template, r *router.Router, a func(http.HandlerFunc)http.HandlerFunc, sysCfg *sysRepo.MemSystemConfigRepo) {
	sysTmpl := template.Must(template.Must(tmpl["login"].Clone()).ParseFiles("templates/admin/sysconfig.html"))

	execSys := func(w http.ResponseWriter, title, cfgType string, items interface{}) {
		sysTmpl.ExecuteTemplate(w, "sysconfig.html", map[string]any{"Title":title,"ConfigType":cfgType,"Items":items})
	}

	r.GET("/admin/system/logistics-api", a(func(w http.ResponseWriter, req *http.Request) {
		cfgs,_:=sysCfg.ListLogisticsAPIs(1);execSys(w,"物流API对接配置","logistics_api",cfgs)
	}))
	r.GET("/admin/system/customs-broker-api", a(func(w http.ResponseWriter, req *http.Request) {
		cfgs:=sysCfg.ListBrokers(req.Context(),1);execSys(w,"清关公司API配置","customs_broker",cfgs)
	}))
	r.GET("/admin/system/printers", a(func(w http.ResponseWriter, req *http.Request) {
		ps:=sysCfg.ListPrinters(1);execSys(w,"打印机对接配置","printer",ps)
	}))
	r.GET("/admin/system/notification-channels", a(func(w http.ResponseWriter, req *http.Request) {
		chs:=sysCfg.ListChannels(req.Context(),1);execSys(w,"通知渠道配置","notification",chs)
	}))
	r.GET("/admin/system/settings", a(func(w http.ResponseWriter, req *http.Request) {
		ss:=sysCfg.ListSettings(1);execSys(w,"系统参数配置","settings",ss)
	}))
}
