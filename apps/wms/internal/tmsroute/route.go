// Package route provides TMS (物流管理) admin route registration.
package route

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/i56/framework/core/router"

	"github.com/i56/i56-apps/i56-wms/internal/common"

	tmsDomain "github.com/i56/modules/transport/domain"
	tmsRepo "github.com/i56/modules/transport/repository"
)

// BftCarrier represents an in-memory carrier for TMS pages.
type BftCarrier struct {
	Name, Code, CustomsPoint, DeliveryMethod, DeliveryPrice, Surcharge, Status string
}

var bftCarriersMu sync.Mutex
var bftCarriers = []BftCarrier{
	{"新竹物流", "HCT", "台北/台中/高雄", "宅配", "¥20固定/≥10kg免运", "超长+超材", "启用"},
	{"黑猫宅急便", "YAMATO", "台北/高雄", "宅配", "¥15固定", "超材", "启用"},
	{"顺丰速运", "SF-TW", "台北/台中", "宅配", "¥12/kg", "超长", "启用"},
}

func addCarrier(name, code, customsPoint, deliveryMethod, deliveryPrice, surcharge string) {
	bftCarriersMu.Lock(); defer bftCarriersMu.Unlock()
	bftCarriers = append(bftCarriers, BftCarrier{name, code, customsPoint, deliveryMethod, deliveryPrice, surcharge, "启用"})
}

// In-memory shipping providers
var bftShippingProvidersMu sync.Mutex
var bftShippingProviders = []struct{ Name, Code, Type, Contact, Phone, Status string }{
	{"远洋航运", "OCEANLINK", "海运", "王总", "13900001111", "启用"},
	{"空港快运", "AIRPORTEX", "空运", "陈经理", "13700002222", "启用"},
	{"海陆通", "SEALAND", "海陆", "林总", "0928123456", "启用"},
}

// In-memory customs brokers
var bftCustomsBrokersMu sync.Mutex
var bftCustomsBrokers = []struct{ Name, Code, BrokerNum, Prefix, Country, Points, Contact string }{
	{"厦门电子口岸", "XM-CUS", "BR001", "XM", "中国", "厦门港/福州港", "张三"},
	{"深圳海关", "SZ-CUS", "BR002", "SZ", "中国", "深圳湾/蛇口港", "李四"},
}

// Register TMS admin routes (~11 list pages + CRUD).
func Register(
	r *router.Router,
	a func(http.HandlerFunc) http.HandlerFunc,
	rc *common.RenderCtx,
	rr *tmsRepo.MemRouteRepo,
	cour *tmsRepo.MemCourierRepo,
) {
	const tenant int64 = 1
	gp := rc.NewGenericList()

	// ─── /admin/couriers — 快递公司 (template-based P3) ───
	r.GET("/admin/couriers", a(func(w http.ResponseWriter, req *http.Request) {
		couriers, _ := cour.List(req.Context())
		rows := make([][]string, len(couriers))
		for i, c := range couriers {
			rows[i] = []string{c.Name, c.Code, c.CountryRegion, "启用", time.Now().Format("01-02")}
		}
		if len(rows) == 0 {
			rows = [][]string{{"顺丰速运", "SF", "中国大陆", "启用", "07-01"}}
		}
		rc.Tmpl["couriers"].ExecuteTemplate(w, "couriers.html", map[string]any{
			"Title": "快递公司", "Page": "tms_couriers",
			"Columns": []string{"名称", "编码", "区域", "状态", "创建时间"},
			"Rows": rows, "Total": len(rows),
			"AddURL": "/admin/couriers/add-form", "HasActions": true,
		})
	}))
	r.GET("/admin/couriers/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增快递公司")+common.FormSave("/admin/couriers/save")+
			common.FormField("名称", "name", "", "快递公司名称")+
			common.FormField("代码", "code", "", "快递公司代码")+
			common.FormField("国家/地区", "region", "", "所在国家或地区")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/couriers/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		cour.Create(req.Context(), &tmsDomain.Courier{Name: req.FormValue("name"), Code: req.FormValue("code"), CountryRegion: req.FormValue("region")})
		common.Redirect(w, "/admin/couriers")
	}))
	r.GET("/admin/couriers/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("id")
		couriers, _ := cour.List(req.Context())
		var c *tmsDomain.Courier
		for i := range couriers {
			if couriers[i].Name == name {
				c = &couriers[i]
				break
			}
		}
		if c == nil { http.Error(w, "not found", 404); return }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑快递公司")+common.FormSave("/admin/couriers/update")+
			fmt.Sprintf(`<input type="hidden" name="old_name" value="%s">`, c.Name)+
			common.FormField("名称", "name", c.Name, "")+
			common.FormField("代码", "code", c.Code, "")+
			common.FormField("国家/地区", "region", c.CountryRegion, "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/couriers/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		oldName := req.FormValue("old_name")
		couriers, _ := cour.List(req.Context())
		for i := range couriers {
			if couriers[i].Name == oldName {
				cour.Update(req.Context(), couriers[i].Code, &tmsDomain.Courier{Name: req.FormValue("name"), Code: req.FormValue("code"), CountryRegion: req.FormValue("region")})
				break
			}
		}
		common.Redirect(w, "/admin/couriers")
	}))
	r.POST("/admin/couriers", a(func(w http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("delete")
		if name == "" { http.Error(w, "bad request", 400); return }
		couriers, _ := cour.List(req.Context())
		for _, c := range couriers {
			if c.Name == name || c.Code == name {
				cour.Delete(req.Context(), c.Code)
				break
			}
		}
		common.Redirect(w, "/admin/couriers")
	}))

	// ─── /admin/shipping-providers — 运输公司 (from admin_modules.go TMS) ───
	r.GET("/admin/shipping-providers", a(func(w http.ResponseWriter, req *http.Request) {
		bftShippingProvidersMu.Lock()
		providers := make([]struct{ Name, Code, Type, Contact, Phone, Status string }, len(bftShippingProviders))
		copy(providers, bftShippingProviders)
		bftShippingProvidersMu.Unlock()
		rows := make([][]string, len(providers))
		for i, p := range providers {
			rows[i] = []string{p.Name, p.Code, p.Type, p.Contact, p.Phone, p.Status}
		}
		gp(w, "tms_providers", "运输公司", len(rows), []string{"名称", "编码", "类型", "联系人", "电话", "状态"}, rows, "/admin/shipping-providers/add-form")
	}))
	r.GET("/admin/shipping-providers/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增运输公司")+common.FormSave("/admin/shipping-providers/save")+
			common.FormField("名称", "name", "", "")+
			common.FormField("编码", "code", "", "")+
			common.FormSelect("类型", "type", "sea", [2]string{"sea", "海运"}, [2]string{"air", "空运"}, [2]string{"sea_land", "海陆"})+
			common.FormField("联系人", "contact", "", "")+
			common.FormField("电话", "phone", "", "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/shipping-providers/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		bftShippingProvidersMu.Lock()
		bftShippingProviders = append(bftShippingProviders, struct{ Name, Code, Type, Contact, Phone, Status string }{
			req.FormValue("name"), req.FormValue("code"), req.FormValue("type"), req.FormValue("contact"), req.FormValue("phone"), "启用",
		})
		bftShippingProvidersMu.Unlock()
		common.Redirect(w, "/admin/shipping-providers")
	}))
	r.GET("/admin/shipping-providers/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("id")
		bftShippingProvidersMu.Lock()
		var sp *struct{ Name, Code, Type, Contact, Phone, Status string }
		for i := range bftShippingProviders {
			if bftShippingProviders[i].Name == name {
				sp = &bftShippingProviders[i]
				break
			}
		}
		bftShippingProvidersMu.Unlock()
		if sp == nil { http.Error(w, "not found", 404); return }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑运输公司")+common.FormSave("/admin/shipping-providers/update")+
			fmt.Sprintf(`<input type="hidden" name="old_name" value="%s">`, sp.Name)+
			common.FormField("名称", "name", sp.Name, "")+
			common.FormField("编码", "code", sp.Code, "")+
			common.FormSelect("类型", "type", sp.Type, [2]string{"sea", "海运"}, [2]string{"air", "空运"}, [2]string{"sea_land", "海陆"})+
			common.FormField("联系人", "contact", sp.Contact, "")+
			common.FormField("电话", "phone", sp.Phone, "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/shipping-providers/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		oldName := req.FormValue("old_name")
		bftShippingProvidersMu.Lock()
		for i := range bftShippingProviders {
			if bftShippingProviders[i].Name == oldName {
				bftShippingProviders[i] = struct{ Name, Code, Type, Contact, Phone, Status string }{
					req.FormValue("name"), req.FormValue("code"), req.FormValue("type"), req.FormValue("contact"), req.FormValue("phone"), bftShippingProviders[i].Status,
				}
				break
			}
		}
		bftShippingProvidersMu.Unlock()
		common.Redirect(w, "/admin/shipping-providers")
	}))
	r.POST("/admin/shipping-providers", a(func(w http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("delete")
		if name == "" { http.Error(w, "bad request", 400); return }
		bftShippingProvidersMu.Lock()
		filtered := make([]struct{ Name, Code, Type, Contact, Phone, Status string }, 0, len(bftShippingProviders))
		for _, sp := range bftShippingProviders {
			if sp.Name != name && sp.Code != name {
				filtered = append(filtered, sp)
			}
		}
		bftShippingProviders = filtered
		bftShippingProvidersMu.Unlock()
		common.Redirect(w, "/admin/shipping-providers")
	}))

	// ─── /admin/customs-brokers — 清关公司 (from admin_modules.go TMS) ───
	r.GET("/admin/customs-brokers", a(func(w http.ResponseWriter, req *http.Request) {
		bftCustomsBrokersMu.Lock()
		brokers := make([]struct{ Name, Code, BrokerNum, Prefix, Country, Points, Contact string }, len(bftCustomsBrokers))
		copy(brokers, bftCustomsBrokers)
		bftCustomsBrokersMu.Unlock()
		rows := make([][]string, len(brokers))
		for i, b := range brokers {
			rows[i] = []string{b.Name, b.Code, b.BrokerNum, b.Prefix, b.Contact, b.Country, "启用"}
		}
		gp(w, "tms_customs_brokers", "清关公司", len(rows), []string{"名称", "代码", "报关号", "前缀", "联系人", "国家", "状态"}, rows, "/admin/customs-brokers/add-form")
	}))
	r.GET("/admin/customs-brokers/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增清关公司")+common.FormSave("/admin/customs-brokers/save")+
			common.FormField("名称", "name", "", "")+
			common.FormField("代码", "code", "", "")+
			common.FormField("编号", "broker_num", "", "")+
			common.FormField("前缀", "prefix", "", "")+
			common.FormField("国家", "country", "", "")+
			common.FormField("清关点", "points", "", "")+
			common.FormField("联系人", "contact", "", "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/customs-brokers/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		bftCustomsBrokersMu.Lock()
		bftCustomsBrokers = append(bftCustomsBrokers, struct{ Name, Code, BrokerNum, Prefix, Country, Points, Contact string }{
			req.FormValue("name"), req.FormValue("code"), req.FormValue("broker_num"), req.FormValue("prefix"), req.FormValue("country"), req.FormValue("points"), req.FormValue("contact"),
		})
		bftCustomsBrokersMu.Unlock()
		common.Redirect(w, "/admin/customs-brokers")
	}))
	r.GET("/admin/customs-brokers/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("id")
		bftCustomsBrokersMu.Lock()
		var cb *struct{ Name, Code, BrokerNum, Prefix, Country, Points, Contact string }
		for i := range bftCustomsBrokers {
			if bftCustomsBrokers[i].Name == name {
				cb = &bftCustomsBrokers[i]
				break
			}
		}
		bftCustomsBrokersMu.Unlock()
		if cb == nil { http.Error(w, "not found", 404); return }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑清关公司")+common.FormSave("/admin/customs-brokers/update")+
			fmt.Sprintf(`<input type="hidden" name="old_name" value="%s">`, cb.Name)+
			common.FormField("名称", "name", cb.Name, "")+
			common.FormField("代码", "code", cb.Code, "")+
			common.FormField("编号", "broker_num", cb.BrokerNum, "")+
			common.FormField("前缀", "prefix", cb.Prefix, "")+
			common.FormField("国家", "country", cb.Country, "")+
			common.FormField("清关点", "points", cb.Points, "")+
			common.FormField("联系人", "contact", cb.Contact, "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/customs-brokers/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		oldName := req.FormValue("old_name")
		bftCustomsBrokersMu.Lock()
		for i := range bftCustomsBrokers {
			if bftCustomsBrokers[i].Name == oldName {
				bftCustomsBrokers[i] = struct{ Name, Code, BrokerNum, Prefix, Country, Points, Contact string }{
					req.FormValue("name"), req.FormValue("code"), req.FormValue("broker_num"), req.FormValue("prefix"), req.FormValue("country"), req.FormValue("points"), req.FormValue("contact"),
				}
				break
			}
		}
		bftCustomsBrokersMu.Unlock()
		common.Redirect(w, "/admin/customs-brokers")
	}))
	r.POST("/admin/customs-brokers", a(func(w http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("delete")
		if name == "" { http.Error(w, "bad request", 400); return }
		bftCustomsBrokersMu.Lock()
		filtered := make([]struct{ Name, Code, BrokerNum, Prefix, Country, Points, Contact string }, 0, len(bftCustomsBrokers))
		for _, cb := range bftCustomsBrokers {
			if cb.Name != name && cb.Code != name {
				filtered = append(filtered, cb)
			}
		}
		bftCustomsBrokers = filtered
		bftCustomsBrokersMu.Unlock()
		common.Redirect(w, "/admin/customs-brokers")
	}))

	// ─── /admin/route-templates — 线路模板 CRUD (template-based P3) ───
	r.GET("/admin/route-templates", a(func(w http.ResponseWriter, req *http.Request) {
		routes, _, _ := rr.List(req.Context(), tenant, 0, 100)
		rows := make([][]string, len(routes))
		for i, rt := range routes {
			rows[i] = []string{rt.Name, fmt.Sprintf("%d", rt.WarehouseID), rt.TransportType, fmt.Sprintf("%d-%d天", rt.MinDays, rt.MaxDays), "台湾", func()string{if rt.IsActive{return "启用"};return "停用"}()}
		}
		if len(rows) == 0 {
			rows = [][]string{{"厦门→台湾空运", "厦门仓", "空运", "3-5天", "台湾", "启用"}, {"深圳→台湾海快", "深圳仓", "海快", "5-7天", "台湾", "启用"}, {"厦门→台湾海运", "厦门仓", "海运", "7-10天", "台湾", "启用"}}
		}
		rc.Tmpl["route_templates"].ExecuteTemplate(w, "route_templates.html", map[string]any{
			"Title": "线路模板", "Page": "route_templates",
			"Columns": []string{"线路名称", "始发仓", "运输方式", "时效", "目的地", "状态"},
			"Rows": rows, "Total": len(rows),
			"AddURL": "/admin/route-templates/add-form", "HasActions": true,
		})
	}))
	r.GET("/admin/route-templates/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增线路模板")+common.FormSave("/admin/route-templates/save")+
			common.FormField("线路名", "name", "", "线路名称")+
			common.FormSelect("运输方式", "transport_type", "air", [2]string{"air", "空运"}, [2]string{"sea_express", "海快"}, [2]string{"sea", "海运"}, [2]string{"land", "陆运"})+
			common.FormField("体积系数", "volume_coeff", "6000", "计费体积系数")+
			common.FormField("最低重量(kg)", "min_weight", "", "最低计费重量")+
			common.FormField("重量单价", "base_weight_price", "", "元/kg")+
			common.FormField("体积单价", "base_volume_price", "", "元/才")+
			common.FormField("最低收费", "min_amount", "", "最低收费金额")+
			common.FormField("最少天数", "min_days", "", "运输最少天数")+
			common.FormField("最多天数", "max_days", "", "运输最多天数")+
			common.FormFooter()+common.ModalEnd())
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
		rr.Create(req.Context(), &tmsDomain.Route{TenantID: tenant, WarehouseID: 1, Name: req.FormValue("name"), TransportType: req.FormValue("transport_type"), VolumeCoeff: vc, MinWeight: mw, BaseWeightPrice: bwp, BaseVolumePrice: bvp, MinAmount: ma, MinDays: mind, MaxDays: maxd, IsActive: true})
		common.Redirect(w, "/admin/route-templates")
	}))
	r.GET("/admin/route-templates/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		id, _ := common.ParseID(req.URL.Query().Get("id"))
		rt, _ := rr.GetByID(req.Context(), tenant, id)
		if rt == nil { http.Error(w, "not found", 404); return }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑线路模板")+common.FormSave("/admin/route-templates/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%d">`, rt.ID)+
			common.FormField("线路名", "name", rt.Name, "")+
			common.FormSelect("运输方式", "transport_type", rt.TransportType, [2]string{"air", "空运"}, [2]string{"sea_express", "海快"}, [2]string{"sea", "海运"}, [2]string{"land", "陆运"})+
			common.FormField("体积系数", "volume_coeff", fmt.Sprintf("%d", rt.VolumeCoeff), "")+
			common.FormField("最低重量(kg)", "min_weight", fmt.Sprintf("%.2f", rt.MinWeight), "")+
			common.FormField("重量单价", "base_weight_price", fmt.Sprintf("%.2f", rt.BaseWeightPrice), "")+
			common.FormField("体积单价", "base_volume_price", fmt.Sprintf("%.2f", rt.BaseVolumePrice), "")+
			common.FormField("最低收费", "min_amount", fmt.Sprintf("%.2f", rt.MinAmount), "")+
			common.FormField("最少天数", "min_days", fmt.Sprintf("%d", rt.MinDays), "")+
			common.FormField("最多天数", "max_days", fmt.Sprintf("%d", rt.MaxDays), "")+
			common.FormFooter()+common.ModalEnd())
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
		rr.Update(req.Context(), &tmsDomain.Route{ID: id, TenantID: tenant, WarehouseID: 1, Name: req.FormValue("name"), TransportType: req.FormValue("transport_type"), VolumeCoeff: vc, MinWeight: mw, BaseWeightPrice: bwp, BaseVolumePrice: bvp, MinAmount: ma, MinDays: mind, MaxDays: maxd, IsActive: true})
		common.Redirect(w, "/admin/route-templates")
	}))
	r.POST("/admin/route-templates", a(func(w http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("delete")
		if name == "" { http.Error(w, "bad request", 400); return }
		all, _, _ := rr.List(req.Context(), tenant, 0, 1000)
		for _, rt := range all {
			if rt.Name == name {
				rr.Delete(req.Context(), tenant, rt.ID)
				break
			}
		}
		common.Redirect(w, "/admin/route-templates")
	}))

	// ─── /admin/carriers — 承运商列表 (template-based P3) ───
	r.GET("/admin/carriers", a(func(w http.ResponseWriter, req *http.Request) {
		bftCarriersMu.Lock()
		carriers := make([]BftCarrier, len(bftCarriers))
		copy(carriers, bftCarriers)
		bftCarriersMu.Unlock()
		rows := make([][]string, len(carriers))
		for i, c := range carriers {
			rows[i] = []string{c.Name, c.Code, c.CustomsPoint, c.DeliveryMethod, c.DeliveryPrice, c.Surcharge, c.Status}
		}
		rc.Tmpl["carriers"].ExecuteTemplate(w, "carriers.html", map[string]any{
			"Title": "承运商列表", "Page": "tms_carriers",
			"Columns": []string{"名称", "编码", "清关点", "派送方式", "派送价", "加收费", "状态"},
			"Rows": rows, "Total": len(rows),
			"AddURL": "/admin/carriers/add-form", "HasActions": true,
		})
	}))
	r.GET("/admin/carriers/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增承运商")+common.FormSave("/admin/carriers/save")+
			common.FormField("名称", "name", "", "承运商名称")+
			common.FormField("编码", "code", "", "承运商编码")+
			common.FormField("清关点", "customs_point", "", "")+
			common.FormField("派送方式", "delivery_method", "", "宅配/专车/自取")+
			common.FormField("派送价", "delivery_price", "", "")+
			common.FormField("加收费", "surcharge", "", "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/carriers/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		addCarrier(req.FormValue("name"), req.FormValue("code"), req.FormValue("customs_point"), req.FormValue("delivery_method"), req.FormValue("delivery_price"), req.FormValue("surcharge"))
		common.Redirect(w, "/admin/carriers")
	}))
	r.GET("/admin/carriers/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("id")
		bftCarriersMu.Lock()
		var c *BftCarrier
		for i := range bftCarriers {
			if bftCarriers[i].Name == name {
				c = &bftCarriers[i]
				break
			}
		}
		bftCarriersMu.Unlock()
		if c == nil { http.Error(w, "not found", 404); return }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑承运商")+common.FormSave("/admin/carriers/update")+
			fmt.Sprintf(`<input type="hidden" name="old_name" value="%s">`, c.Name)+
			common.FormField("名称", "name", c.Name, "")+
			common.FormField("编码", "code", c.Code, "")+
			common.FormField("清关点", "customs_point", c.CustomsPoint, "")+
			common.FormField("派送方式", "delivery_method", c.DeliveryMethod, "")+
			common.FormField("派送价", "delivery_price", c.DeliveryPrice, "")+
			common.FormField("加收费", "surcharge", c.Surcharge, "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/carriers/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		oldName := req.FormValue("old_name")
		bftCarriersMu.Lock()
		for i := range bftCarriers {
			if bftCarriers[i].Name == oldName {
				bftCarriers[i] = BftCarrier{req.FormValue("name"), req.FormValue("code"), req.FormValue("customs_point"), req.FormValue("delivery_method"), req.FormValue("delivery_price"), req.FormValue("surcharge"), bftCarriers[i].Status}
				break
			}
		}
		bftCarriersMu.Unlock()
		common.Redirect(w, "/admin/carriers")
	}))
	r.POST("/admin/carriers", a(func(w http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("delete")
		if name == "" { http.Error(w, "bad request", 400); return }
		bftCarriersMu.Lock()
		filtered := make([]BftCarrier, 0, len(bftCarriers))
		for _, c := range bftCarriers {
			if c.Name != name && c.Code != name {
				filtered = append(filtered, c)
			}
		}
		bftCarriers = filtered
		bftCarriersMu.Unlock()
		common.Redirect(w, "/admin/carriers")
	}))

	// ─── /admin/carriers/{id}/delivery-price — 承运商派送价明细 ───
	r.GET("/admin/carriers/{id}/delivery-price", a(func(w http.ResponseWriter, req *http.Request) {
		id := req.PathValue("id")
		carrierName := id
		bftCarriersMu.Lock()
		for _, c := range bftCarriers {
			if fmt.Sprint(len(bftCarriers)) == id || c.Code == id || c.Name == id {
				carrierName = c.Name; break
			}
		}
		bftCarriersMu.Unlock()
		common.HtmlOK(w)
		// Sample delivery price tiers for the carrier
		rows := [][]string{
			{"台北", "宅配", "≤10kg", "¥20", "固定"},
			{"台中", "宅配", "≤10kg", "¥25", "固定"},
			{"高雄", "宅配", "≤10kg", "¥30", "固定"},
			{"台北", "专车", "不限", "¥150", "按趟"},
			{"台中", "专车", "不限", "¥180", "按趟"},
			{"全台", "自取", "不限", "¥0", "免费"},
		}
		coloredRows := make([][]string, len(rows)); copy(coloredRows, rows)
		gp(w, "tms_carrier_delivery", carrierName+" — 派送价", len(rows),
			[]string{"区域", "派送方式", "条件", "价格", "计费方式"}, coloredRows)
	}))

	// ─── /admin/carriers/{id}/surcharge — 承运商加收价明细 ───
	r.GET("/admin/carriers/{id}/surcharge", a(func(w http.ResponseWriter, req *http.Request) {
		id := req.PathValue("id")
		carrierName := id
		bftCarriersMu.Lock()
		for _, c := range bftCarriers {
			if c.Code == id || c.Name == id {
				carrierName = c.Name; break
			}
		}
		bftCarriersMu.Unlock()
		common.HtmlOK(w)
		// Sample surcharge rules for the carrier
		rows := [][]string{
			{"超长费", "单边 > 120cm", "¥30/件", "叠加"},
			{"超材费", "材积 > 200cm", "¥50/件", "叠加"},
			{"偏远费", "偏远地区", "¥80/票", "叠加"},
			{"上楼费", "无电梯", "¥50/件", "可选"},
			{"保价费", "申报价值 > ¥5000", "1%", "可选"},
		}
		coloredRows := make([][]string, len(rows)); copy(coloredRows, rows)
		gp(w, "tms_carrier_surcharge", carrierName+" — 加收价", len(rows),
			[]string{"加收项目", "触发条件", "加收标准", "计费模式"}, coloredRows)
	}))

	// ─── /admin/cargo-types — 货物类型 (from admin_crud.go) ───
	var cargoTypesMu sync.Mutex
	var cargoTypes = []struct{ Name, Code, Description string }{{"普货", "general", "普通货物"}, {"一类", "class1", "一类危险品"}, {"家具类", "furniture", "家具"}, {"易碎品", "fragile", "易碎品"}}

	r.GET("/admin/cargo-types", a(func(w http.ResponseWriter, req *http.Request) {
		cargoTypesMu.Lock()
		cts := append([]struct{ Name, Code, Description string }{}, cargoTypes...)
		cargoTypesMu.Unlock()
		rows := make([][]string, len(cts))
		for i, ct := range cts { rows[i] = []string{ct.Name, ct.Code, ct.Description, "启用"} }
		gp(w, "tms_cargo_types", "货物类型", len(rows), []string{"名称", "编码", "描述", "状态"}, rows, "/admin/cargo-types/add-form")
	}))
	r.GET("/admin/cargo-types/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增货物类型")+common.FormSave("/admin/cargo-types/save")+
			common.FormField("名称", "name", "", "如: 普货")+
			common.FormField("编码", "code", "", "如: general")+
			common.FormField("描述", "description", "", "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/cargo-types/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		cargoTypesMu.Lock()
		cargoTypes = append(cargoTypes, struct{ Name, Code, Description string }{req.FormValue("name"), req.FormValue("code"), req.FormValue("description")})
		cargoTypesMu.Unlock()
		common.Redirect(w, "/admin/cargo-types")
	}))
	r.GET("/admin/cargo-types/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("id")
		cargoTypesMu.Lock()
		var ct *struct{ Name, Code, Description string }
		for i := range cargoTypes {
			if cargoTypes[i].Name == name {
				ct = &cargoTypes[i]
				break
			}
		}
		cargoTypesMu.Unlock()
		if ct == nil { http.Error(w, "not found", 404); return }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑货物类型")+common.FormSave("/admin/cargo-types/update")+
			fmt.Sprintf(`<input type="hidden" name="old_name" value="%s">`, ct.Name)+
			common.FormField("名称", "name", ct.Name, "")+
			common.FormField("编码", "code", ct.Code, "")+
			common.FormField("描述", "description", ct.Description, "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/cargo-types/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		oldName := req.FormValue("old_name")
		cargoTypesMu.Lock()
		for i := range cargoTypes {
			if cargoTypes[i].Name == oldName {
				cargoTypes[i] = struct{ Name, Code, Description string }{req.FormValue("name"), req.FormValue("code"), req.FormValue("description")}
				break
			}
		}
		cargoTypesMu.Unlock()
		common.Redirect(w, "/admin/cargo-types")
	}))
	r.POST("/admin/cargo-types", a(func(w http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("delete")
		if name == "" { http.Error(w, "bad request", 400); return }
		cargoTypesMu.Lock()
		filtered := make([]struct{ Name, Code, Description string }, 0, len(cargoTypes))
		for _, ct := range cargoTypes {
			if ct.Name != name && ct.Code != name {
				filtered = append(filtered, ct)
			}
		}
		cargoTypes = filtered
		cargoTypesMu.Unlock()
		common.Redirect(w, "/admin/cargo-types")
	}))

	// ─── /admin/area-groups — 区域组 (template-based P3) ───
	var areaGroupsMu sync.Mutex
	var areaGroups = []struct{ Name, Code, Coverage string }{{"华南区", "CN-SOUTH", "广东/福建/海南"}, {"华东区", "CN-EAST", "上海/浙江/江苏"}, {"台湾区", "TW", "全台湾"}}

	r.GET("/admin/area-groups", a(func(w http.ResponseWriter, req *http.Request) {
		areaGroupsMu.Lock()
		ags := append([]struct{ Name, Code, Coverage string }{}, areaGroups...)
		areaGroupsMu.Unlock()
		rows := make([][]string, len(ags))
		for i, ag := range ags { rows[i] = []string{ag.Name, ag.Code, ag.Coverage, "启用"} }
		rc.Tmpl["area_groups"].ExecuteTemplate(w, "area_groups.html", map[string]any{
			"Title": "区域组管理", "Page": "tms_area_groups",
			"Columns": []string{"区域名", "编码", "覆盖范围", "状态"},
			"Rows": rows, "Total": len(rows),
			"AddURL": "/admin/area-groups/add-form", "HasActions": true,
		})
	}))
	r.GET("/admin/area-groups/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增区域组")+common.FormSave("/admin/area-groups/save")+
			common.FormField("区域名", "name", "", "如: 华南区")+
			common.FormField("编码", "code", "", "如: CN-SOUTH")+
			common.FormField("覆盖范围", "coverage", "", "如: 广东/福建/海南")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/area-groups/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		areaGroupsMu.Lock()
		areaGroups = append(areaGroups, struct{ Name, Code, Coverage string }{req.FormValue("name"), req.FormValue("code"), req.FormValue("coverage")})
		areaGroupsMu.Unlock()
		common.Redirect(w, "/admin/area-groups")
	}))
	r.GET("/admin/area-groups/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("id")
		areaGroupsMu.Lock()
		var ag *struct{ Name, Code, Coverage string }
		for i := range areaGroups {
			if areaGroups[i].Name == name {
				ag = &areaGroups[i]
				break
			}
		}
		areaGroupsMu.Unlock()
		if ag == nil { http.Error(w, "not found", 404); return }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑区域组")+common.FormSave("/admin/area-groups/update")+
			fmt.Sprintf(`<input type="hidden" name="old_name" value="%s">`, ag.Name)+
			common.FormField("区域名", "name", ag.Name, "")+
			common.FormField("编码", "code", ag.Code, "")+
			common.FormField("覆盖范围", "coverage", ag.Coverage, "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/area-groups/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		oldName := req.FormValue("old_name")
		areaGroupsMu.Lock()
		for i := range areaGroups {
			if areaGroups[i].Name == oldName {
				areaGroups[i] = struct{ Name, Code, Coverage string }{req.FormValue("name"), req.FormValue("code"), req.FormValue("coverage")}
				break
			}
		}
		areaGroupsMu.Unlock()
		common.Redirect(w, "/admin/area-groups")
	}))
	r.POST("/admin/area-groups", a(func(w http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("delete")
		if name == "" { http.Error(w, "bad request", 400); return }
		areaGroupsMu.Lock()
		filtered := make([]struct{ Name, Code, Coverage string }, 0, len(areaGroups))
		for _, ag := range areaGroups {
			if ag.Name != name && ag.Code != name {
				filtered = append(filtered, ag)
			}
		}
		areaGroups = filtered
		areaGroupsMu.Unlock()
		common.Redirect(w, "/admin/area-groups")
	}))

	// ─── /admin/transport-modes — 运输方式 (from admin_crud.go) ───
	var transportModesMu sync.Mutex
	var transportModes = []struct{ Name, Code, Description string }{{"空运", "air", "航空运输"}, {"海快", "sea_express", "海运快件"}, {"海运", "sea", "海上运输"}, {"陆运", "land", "陆路运输"}}

	r.GET("/admin/transport-modes", a(func(w http.ResponseWriter, req *http.Request) {
		transportModesMu.Lock()
		tms := append([]struct{ Name, Code, Description string }{}, transportModes...)
		transportModesMu.Unlock()
		rows := make([][]string, len(tms))
		for i, tm := range tms { rows[i] = []string{tm.Name, tm.Code, tm.Description, "启用"} }
		gp(w, "tms_modes", "运输方式", len(rows), []string{"方式", "编码", "描述", "状态"}, rows, "/admin/transport-modes/add-form")
	}))
	r.GET("/admin/transport-modes/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增运输方式")+common.FormSave("/admin/transport-modes/save")+
			common.FormField("方式", "name", "", "如: 空运")+
			common.FormField("编码", "code", "", "如: air")+
			common.FormField("描述", "description", "", "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/transport-modes/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		transportModesMu.Lock()
		transportModes = append(transportModes, struct{ Name, Code, Description string }{req.FormValue("name"), req.FormValue("code"), req.FormValue("description")})
		transportModesMu.Unlock()
		common.Redirect(w, "/admin/transport-modes")
	}))
	r.GET("/admin/transport-modes/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("id")
		transportModesMu.Lock()
		var tm *struct{ Name, Code, Description string }
		for i := range transportModes {
			if transportModes[i].Name == name {
				tm = &transportModes[i]
				break
			}
		}
		transportModesMu.Unlock()
		if tm == nil { http.Error(w, "not found", 404); return }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑运输方式")+common.FormSave("/admin/transport-modes/update")+
			fmt.Sprintf(`<input type="hidden" name="old_name" value="%s">`, tm.Name)+
			common.FormField("方式", "name", tm.Name, "")+
			common.FormField("编码", "code", tm.Code, "")+
			common.FormField("描述", "description", tm.Description, "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/transport-modes/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		oldName := req.FormValue("old_name")
		transportModesMu.Lock()
		for i := range transportModes {
			if transportModes[i].Name == oldName {
				transportModes[i] = struct{ Name, Code, Description string }{req.FormValue("name"), req.FormValue("code"), req.FormValue("description")}
				break
			}
		}
		transportModesMu.Unlock()
		common.Redirect(w, "/admin/transport-modes")
	}))
	r.POST("/admin/transport-modes", a(func(w http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("delete")
		if name == "" { http.Error(w, "bad request", 400); return }
		transportModesMu.Lock()
		filtered := make([]struct{ Name, Code, Description string }, 0, len(transportModes))
		for _, tm := range transportModes {
			if tm.Name != name && tm.Code != name {
				filtered = append(filtered, tm)
			}
		}
		transportModes = filtered
		transportModesMu.Unlock()
		common.Redirect(w, "/admin/transport-modes")
	}))

	// ─── /admin/logistics-tracking — 物流追踪 ───
	r.GET("/admin/logistics-tracking", a(func(w http.ResponseWriter, req *http.Request) {
		rows := [][]string{
			{"CT-8837291", "运输中", "台北配送站", "新竹物流", "07-11 09:30"},
			{"CT-8837292", "已签收", "台中", "黑猫宅急便", "07-10 15:00"},
			{"CT-8837293", "清关中", "台北港", "新竹物流", "07-10 11:00"},
			{"CT-8837294", "运输中", "基隆港", "顺丰速运", "07-09 18:20"},
			{"CT-8837295", "已装柜", "厦门仓", "新竹物流", "07-11 08:00"},
		}
		gp(w, "tms_tracking", "物流追踪", len(rows), []string{"运单号", "状态", "位置", "承运商", "更新时间"}, rows, "/admin/logistics-tracking/add-form")
	}))
	r.GET("/admin/logistics-tracking/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增物流追踪")+common.FormSave("/admin/logistics-tracking/save")+
			common.FormField("运单号", "tracking_no", "", "运单号")+
			common.FormField("状态", "status", "", "运输中/已签收/清关中")+
			common.FormField("位置", "location", "", "当前位置")+
			common.FormField("承运商", "carrier", "", "承运商名称")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.GET("/admin/logistics-tracking/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		id := req.URL.Query().Get("id")
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑物流追踪")+common.FormSave("/admin/logistics-tracking/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%s">`, id)+
			common.FormField("运单号", "tracking_no", id, "")+
			common.FormField("状态", "status", "运输中", "")+
			common.FormField("位置", "location", "", "")+
			common.FormField("承运商", "carrier", "", "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/logistics-tracking/update", a(func(w http.ResponseWriter, req *http.Request) { common.Redirect(w, "/admin/logistics-tracking") }))

	// ─── /admin/container-loadings — 装柜记录 ───
	r.GET("/admin/container-loadings", a(func(w http.ResponseWriter, req *http.Request) {
		rows := [][]string{
			{"CTNR-001", "20GP", "MSC-1234", "厦门仓", "7", "2026-07-11 08:00", "已装柜"},
			{"CTNR-002", "40HQ", "MSK-5678", "厦门仓", "12", "2026-07-10 14:30", "已装柜"},
			{"CTNR-003", "20GP", "COSCO-9012", "厦门仓", "5", "2026-07-11 10:00", "待装柜"},
			{"CTNR-004", "40GP", "OOCL-3456", "厦门仓", "8", "2026-07-09 16:00", "运输中"},
		}
		gp(w, "tms_containers", "装柜记录", len(rows), []string{"柜号", "柜型", "船名航次", "装柜仓库", "包裹数", "装柜时间", "状态"}, rows, "/admin/container-loadings/add-form")
	}))
	r.GET("/admin/container-loadings/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增装柜记录")+common.FormSave("/admin/container-loadings/save")+
			common.FormField("柜号", "container_no", "", "如: CTNR-001")+
			common.FormSelect("柜型", "container_type", "20GP", [2]string{"20GP", "20GP"}, [2]string{"40GP", "40GP"}, [2]string{"40HQ", "40HQ"})+
			common.FormField("船名航次", "vessel", "", "")+
			common.FormField("装柜仓库", "warehouse", "", "厦门仓")+
			common.FormField("包裹数", "parcel_count", "", "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/container-loadings/save", a(func(w http.ResponseWriter, req *http.Request) { common.Redirect(w, "/admin/container-loadings") }))
	r.GET("/admin/container-loadings/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		id := req.URL.Query().Get("id")
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑装柜记录")+common.FormSave("/admin/container-loadings/update")+
			fmt.Sprintf(`<input type="hidden" name="id" value="%s">`, id)+
			common.FormField("柜号", "container_no", id, "")+
			common.FormSelect("柜型", "container_type", "20GP", [2]string{"20GP", "20GP"}, [2]string{"40GP", "40GP"}, [2]string{"40HQ", "40HQ"})+
			common.FormField("船名航次", "vessel", "", "")+
			common.FormField("装柜仓库", "warehouse", "", "")+
			common.FormField("包裹数", "parcel_count", "", "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/container-loadings/update", a(func(w http.ResponseWriter, req *http.Request) { common.Redirect(w, "/admin/container-loadings") }))

	// ─── /admin/customs-points — 清关点管理 ───
	var customsPointsData = [][]string{
		{"台北港", "TPE-PORT", "台湾", "新竹物流", "启用"},
		{"基隆港", "KEE-PORT", "台湾", "黑猫宅急便", "启用"},
		{"台中港", "TXG-PORT", "台湾", "顺丰速运", "启用"},
		{"高雄港", "KHH-PORT", "台湾", "新竹物流", "启用"},
	}
	r.GET("/admin/customs-points", a(func(w http.ResponseWriter, req *http.Request) {
		rows := make([][]string, len(customsPointsData))
		copy(rows, customsPointsData)
		gp(w, "tms_customs_points", "清关点管理", len(rows), []string{"名称", "编码", "国家", "承运商", "状态"}, rows, "/admin/customs-points/add-form")
	}))
	r.GET("/admin/customs-points/add-form", a(func(w http.ResponseWriter, req *http.Request) {
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("新增清关点")+common.FormSave("/admin/customs-points/save")+
			common.FormField("名称", "name", "", "如: 台北港")+
			common.FormField("编码", "code", "", "如: TPE-PORT")+
			common.FormField("国家/地区", "country", "", "如: 台湾")+
			common.FormField("承运商", "carrier", "", "如: 新竹物流")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/customs-points/save", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		customsPointsData = append(customsPointsData, []string{req.FormValue("name"), req.FormValue("code"), req.FormValue("country"), req.FormValue("carrier"), "启用"})
		common.Redirect(w, "/admin/customs-points")
	}))
	r.GET("/admin/customs-points/edit-form", a(func(w http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("id")
		var found []string
		for _, row := range customsPointsData {
			if row[0] == name { found = row; break }
		}
		if found == nil { http.Error(w, "not found", 404); return }
		common.HtmlOK(w)
		fmt.Fprint(w, common.ModalStart("编辑清关点")+common.FormSave("/admin/customs-points/update")+
			fmt.Sprintf(`<input type="hidden" name="old_name" value="%s">`, found[0])+
			common.FormField("名称", "name", found[0], "")+
			common.FormField("编码", "code", found[1], "")+
			common.FormField("国家/地区", "country", found[2], "")+
			common.FormField("承运商", "carrier", found[3], "")+
			common.FormFooter()+common.ModalEnd())
	}))
	r.POST("/admin/customs-points/update", a(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		oldName := req.FormValue("old_name")
		for i, row := range customsPointsData {
			if row[0] == oldName {
				customsPointsData[i] = []string{req.FormValue("name"), req.FormValue("code"), req.FormValue("country"), req.FormValue("carrier"), row[4]}
				break
			}
		}
		common.Redirect(w, "/admin/customs-points")
	}))
	r.POST("/admin/customs-points", a(func(w http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("delete")
		if name == "" { http.Error(w, "bad request", 400); return }
		filtered := make([][]string, 0, len(customsPointsData))
		for _, row := range customsPointsData {
			if row[0] != name && row[1] != name {
				filtered = append(filtered, row)
			}
		}
		customsPointsData = filtered
		common.Redirect(w, "/admin/customs-points")
	}))
}
