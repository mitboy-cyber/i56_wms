package main
import (
	"context";"fmt";"html/template";"net/http";"strconv";"time"
	"github.com/i56/framework/core/auth";"github.com/i56/framework/core/router"
	custDomain "github.com/i56/modules/customer/domain";custRepo "github.com/i56/modules/customer/repository"
	parcelDomain "github.com/i56/modules/parcel/domain";parcelRepo "github.com/i56/modules/parcel/repository"
	parcelSvc "github.com/i56/modules/parcel/service";orderDomain "github.com/i56/modules/order/domain"
	orderRepo "github.com/i56/modules/order/repository";orderSvc "github.com/i56/modules/order/service"
	whDomain "github.com/i56/modules/warehouse/domain"
	whRepo "github.com/i56/modules/warehouse/repository";whSvc "github.com/i56/modules/warehouse/service";weightDomain "github.com/i56/modules/weight/domain"
	tmsRepo "github.com/i56/modules/transport/repository"; tmsDomain "github.com/i56/modules/transport/domain"
	psRepo "github.com/i56/modules/parcel_service/repository"
	woRepo "github.com/i56/modules/workorder/repository"
	printRepo "github.com/i56/modules/print/repository"
	whRepo2 "github.com/i56/modules/webhook/repository"
	reportDomain "github.com/i56/modules/report/domain"
	pricingRepo "github.com/i56/modules/pricing/repository"
)

func adminOnly(tm *auth.TokenManager)func(http.HandlerFunc)http.HandlerFunc{
	return func(next http.HandlerFunc)http.HandlerFunc{
		return func(w http.ResponseWriter,r *http.Request){next(w,r)}
	}
}

func execTpl(tmpl map[string]*template.Template,key string,w http.ResponseWriter,name string,data any){tmpl[key].ExecuteTemplate(w,name,data)}

func initTemplates()map[string]*template.Template{
	fm:=template.FuncMap{"statusColor":func(s parcelDomain.ParcelStatus)string{switch s{case "pre_declared":return "secondary";case "received":return "info";case "weighed":return "primary";case "stored":return "success";case "picked":return "warning";case "shipped":return "dark";default:return "secondary"}},"orderStatusColor":func(s orderDomain.OrderStatus)string{switch s{case "pending_picking":return "warning";case "picking":return "info";case "pending_packing":return "primary";default:return "secondary"}},"statusDisplay":func(s string)string{return bftParcelStatus(s)},"hasPrefix":func(s,prefix string)bool{return len(s)>=len(prefix)&&s[:len(prefix)]==prefix},"add":func(a,b int)int{return a+b},"sub":func(a,b int)int{return a-b},"mul":func(a,b int)int{return a*b},"div":func(a,b int)int{if b==0{return 0};return a/b}}
	tmpl:=map[string]*template.Template{}
	for _,p:=range[]struct{k,file string}{{"login","login.html"},{"dashboard","dashboard.html"},{"clients","clients.html"},{"parcels","parcels.html"},{"orders","orders.html"},{"warehouses","warehouses.html"},{"routes","routes.html"},{"generic_list","generic_list.html"},{"warehouse_console","warehouse_console.html"},{"admin_permissions","admin/admin_permissions.html"},{"admin_roles","admin/admin_roles.html"},{"admin_users","admin/admin_users.html"},{"admin_client_permissions","admin/admin_client_permissions.html"},
		{"base_new","admin/base_new.html"},
	}{
		if p.k == "base_new" {
			files := []string{"templates/base.html", "templates/admin/base_new.html"}
			tmpl[p.k]=template.Must(template.New(p.k).Funcs(fm).ParseFiles(files...))
		} else {
			files:=[]string{"templates/base.html","templates/sidebar.html","templates/"+p.file}
			tmpl[p.k]=template.Must(template.New(p.k).Funcs(fm).ParseFiles(files...))
		}
	}
	return tmpl
}

func initClientTemplates()map[string]*template.Template{
	fm:=template.FuncMap{"statusColor":func(s parcelDomain.ParcelStatus)string{switch s{case "pre_declared":return "secondary";case "received":return "info";case "weighed":return "primary";case "stored":return "success";case "picked":return "warning";case "shipped":return "dark";default:return "secondary"}},"orderStatusColor":func(s orderDomain.OrderStatus)string{switch s{case "pending_picking":return "warning";case "picking":return "info";case "pending_packing":return "primary";default:return "secondary"}},"statusDisplay":func(s string)string{return bftParcelStatus(s)},"hasPrefix":func(s,prefix string)bool{return len(s)>=len(prefix)&&s[:len(prefix)]==prefix},"add":func(a,b int)int{return a+b},"sub":func(a,b int)int{return a-b},"mul":func(a,b int)int{return a*b},"div":func(a,b int)int{if b==0{return 0};return a/b}}
	return map[string]*template.Template{
		"login":template.Must(template.New("clogin").Funcs(fm).ParseFiles("templates/client/base.html","templates/client/login.html")),
		"dashboard":template.Must(template.New("cdash").Funcs(fm).ParseFiles("templates/client/base.html","templates/client/dashboard.html")),
		"predeclare":template.Must(template.New("cpred").Funcs(fm).ParseFiles("templates/client/base.html","templates/client/predeclare.html")),
		"parcels":template.Must(template.New("cparcels").Funcs(fm).ParseFiles("templates/client/base.html","templates/client/parcels.html")),
		"ledger":template.Must(template.New("cledger").Funcs(fm).ParseFiles("templates/client/base.html","templates/client/ledger.html")),
		"client_orders":template.Must(template.New("corders2").Funcs(fm).ParseFiles("templates/client/base.html","templates/client/client_orders.html")),
		"client_order_new":template.Must(template.New("cordnew").Funcs(fm).ParseFiles("templates/client/base.html","templates/client/client_order_new.html")),
		"client_declarants":template.Must(template.New("cdecl").Funcs(fm).ParseFiles("templates/client/base.html","templates/client/client_declarants.html")),
		"client_members":template.Must(template.New("cmemb").Funcs(fm).ParseFiles("templates/client/base.html","templates/client/client_members.html")),
		"client_addresses":template.Must(template.New("caddr").Funcs(fm).ParseFiles("templates/client/base.html","templates/client/client_addresses.html")),
		"client_warehouses":template.Must(template.New("cwh").Funcs(fm).ParseFiles("templates/client/base.html","templates/client/client_warehouses.html")),
		"client_route_prices":template.Must(template.New("crp").Funcs(fm).ParseFiles("templates/client/base.html","templates/client/client_route_prices.html")),
		"client_delivery_fees":template.Must(template.New("cdf").Funcs(fm).ParseFiles("templates/client/base.html","templates/client/client_delivery_fees.html")),
		"client_service_orders":template.Must(template.New("cso").Funcs(fm).ParseFiles("templates/client/base.html","templates/client/client_service_orders.html")),
		"client_carrier_surcharges":template.Must(template.New("ccs").Funcs(fm).ParseFiles("templates/client/base.html","templates/client/client_carrier_surcharges.html")),
		"client_webhooks":template.Must(template.New("cwh").Funcs(fm).ParseFiles("templates/client/base.html","templates/client/client_webhooks.html")),
		"client_api_credentials":template.Must(template.New("capi").Funcs(fm).ParseFiles("templates/client/base.html","templates/client/client_api_credentials.html")),
		"client_monthly_statements":template.Must(template.New("cms").Funcs(fm).ParseFiles("templates/client/base.html","templates/client/client_monthly_statements.html")),
		"client_weight_dashboard":template.Must(template.New("cwd").Funcs(fm).ParseFiles("templates/client/base.html","templates/client/client_weight_dashboard.html")),
		"client_webhook_logs":template.Must(template.New("cwhlog").Funcs(fm).ParseFiles("templates/client/base.html","templates/client/client_webhook_logs.html")),
		"client_warehouse_info":template.Must(template.New("cwhinfo").Funcs(fm).ParseFiles("templates/client/base.html","templates/client/client_warehouse_info.html")),
		"client_carrier_delivery":template.Must(template.New("ccdel").Funcs(fm).ParseFiles("templates/client/base.html","templates/client/client_carrier_delivery.html")),
		"client_carrier_surcharge":template.Must(template.New("ccsur").Funcs(fm).ParseFiles("templates/client/base.html","templates/client/client_carrier_surcharge.html")),
		"client_pricing":template.Must(template.New("cpric").Funcs(fm).ParseFiles("templates/client/base.html","templates/client/client_pricing.html")),
	}
}

func clientPg(tm *auth.TokenManager,cTmpl map[string]*template.Template,r *router.Router,ps *parcelSvc.ParcelService,osvc *orderSvc.OrderService,rr *tmsRepo.MemRouteRepo,cour *tmsRepo.MemCourierRepo,ws *whSvc.WarehouseService,pr *parcelRepo.MemParcelRepo,lr *custRepo.MemLedgerRepo,weightRepo *weightDomain.MemWeightRepo,dr *custRepo.MemDeclarantRepo,mr *custRepo.MemMemberRepo,sr *psRepo.MemServiceRepo,whr *whRepo2.MemWebhookRepo,ar *custRepo.MemAddressRepo,rpr *pricingRepo.MemRoutePriceRepo,dfr *pricingRepo.MemDeliveryFeeRepo,scr *pricingRepo.MemSurchargeRepo,acr *pricingRepo.MemApiCredentialRepo){
	ca:=func(next http.HandlerFunc)http.HandlerFunc{
		return func(w http.ResponseWriter,req *http.Request){
			ck,err:=req.Cookie("client_token");if err!=nil{http.Redirect(w,req,"/client/login",303);return}
			if _,err:=tm.ValidateAccessToken(ck.Value);err!=nil{http.Redirect(w,req,"/client/login?error=expired",303);return}
			next(w,req)
		}
	}
	r.GET("/client",ca(func(w http.ResponseWriter,req *http.Request){
		ctx:=req.Context()
		parcels,pt,_:=ps.List(ctx,1,0,8)
		_,rc,_:=rr.List(ctx,1,0,50)
		orders,oc,_:=osvc.List(ctx,1,0,100)
		balance:=0.0;creditLimit:=5000.0
		if entries:=lr.GetByClient(ctx,1,1);len(entries)>0{balance=entries[len(entries)-1].BalanceAfter}
		// Build orders for dashboard display
		orderMaps:=make([]map[string]any,0,len(orders))
		activeCount:=0
		statusCN:=map[orderDomain.OrderStatus]string{
			orderDomain.StatusPendingPicking:"待拣货",orderDomain.StatusPicking:"拣货中",
			orderDomain.StatusPendingPacking:"待打包",
			orderDomain.StatusPendingLoading:"待装柜",orderDomain.StatusLoaded:"已装柜",
			orderDomain.StatusInTransit:"运输中",orderDomain.StatusCustomsClearance:"清关中",
			orderDomain.StatusOutForDelivery:"派送中",orderDomain.StatusCompleted:"已完成",
			orderDomain.StatusCancelled:"已取消",orderDomain.StatusShipped:"已发货",
		}
		for _,o:=range orders{
			s:=statusCN[o.Status];if s==""{s=string(o.Status)}
			orderMaps=append(orderMaps,map[string]any{
				"OrderNo":o.OrderNo,"ReceiverName":o.RecipientName,
				"ParcelCount":o.ParcelCount,"Weight":fmt.Sprintf("%.2f",o.TotalActualWeight),
				"Amount":fmt.Sprintf("%.2f",o.TotalPrice),"Status":s,
			})
			if o.Status!=orderDomain.StatusCompleted&&o.Status!=orderDomain.StatusCancelled{activeCount++}
		}
		// Build parcels for dashboard display
		type parcelDash struct{TrackingNumber,ProductName,StatusLabel,StatusColor string;Weight float64}
		var pdList []parcelDash
		statusLabel:=func(s parcelDomain.ParcelStatus)(string,string){
			switch s{
			case parcelDomain.StatusPreDeclared:return "预报","secondary"
			case parcelDomain.StatusReceived:return "已入仓","info"
			case parcelDomain.StatusWeighed:return "已称重","primary"
			case parcelDomain.StatusStored:return "已上架","success"
			case parcelDomain.StatusPicked:return "已拣货","warning"
			case parcelDomain.StatusPacked:return "已打包","success"
			case parcelDomain.StatusLoaded:return "已装柜","primary"
			case parcelDomain.StatusOutbound:return "已出货","dark"
			default:return string(s),"secondary"
			}
		}
		for _,p:=range parcels{lb,sc:=statusLabel(p.Status);pdList=append(pdList,parcelDash{p.TrackingNumber,p.ProductName,lb,sc,p.ActualWeight})}
		// Parcel status counts
		var preDec,recvd,weighed,stored,picked,packed,shipped int
		allParcels,_,_:=ps.List(ctx,1,0,200)
		for _,p:=range allParcels{
			switch p.Status{
			case parcelDomain.StatusPreDeclared:preDec++
			case parcelDomain.StatusReceived:recvd++
			case parcelDomain.StatusWeighed:weighed++
			case parcelDomain.StatusStored:stored++
			case parcelDomain.StatusPicked:picked++
			case parcelDomain.StatusPacked:packed++
			case parcelDomain.StatusShipped,parcelDomain.StatusLoaded,parcelDomain.StatusOutbound:shipped++
			}
		}
		execTpl(cTmpl,"dashboard",w,"dashboard.html",map[string]any{
			"Title":"主控台","Balance":balance,"CreditLimit":creditLimit,
			"AvailableCredit":creditLimit-balance,
			"TotalParcels":pt,"ParcelCount":pt,"OrderCount":oc,"ActiveOrderCount":activeCount,
			"Parcels":pdList,"Orders":orderMaps,"RouteCount":rc,
			"PreDeclaredCount":preDec,"ReceivedCount":recvd,"WeighedCount":weighed,
			"StoredCount":stored,"PickedCount":picked,"PackingCount":0,"PackedCount":packed,
			"ShippedCount":shipped,
		})
	}))
	r.GET("/client/predeclare",ca(func(w http.ResponseWriter,req *http.Request){
		ctx:=req.Context()
		warehouses,_,_:=ws.List(ctx,1,0,50)
		allParcels,pt,_:=ps.List(ctx,1,0,200)
		var preDeclared,received,weighed,stored,picked,packed int64
		for _,p:=range allParcels{
			switch p.Status{
			case parcelDomain.StatusPreDeclared:preDeclared++
			case parcelDomain.StatusReceived:received++
			case parcelDomain.StatusWeighed:weighed++
			case parcelDomain.StatusStored:stored++
			case parcelDomain.StatusPicked:picked++
			case parcelDomain.StatusPacked:packed++
			}
		}
		type member struct{ID int;Name string}
		type stat struct{Label string;Count int64}
		type recent struct{TN,Name,Status string}
		recentList:=[]recent{}
		for i,p:=range allParcels{if i>=5{break};recentList=append(recentList,recent{p.TrackingNumber,p.ProductName,string(p.Status)})}
		execTpl(cTmpl,"predeclare",w,"predeclare.html",map[string]any{
			"Warehouses":warehouses,"Members":[]member{{1,"王仁照"},{2,"吴欣如"},{3,"张致廷"}},
			"Stats":[]stat{
				{"全部包裹",pt},{"预报",preDeclared},{"已入仓",received},
				{"已上架",stored},{"待打包",picked},{"打包中",0},{"已打包",packed},
			},
			"Recent":recentList,
		})
	}))
	r.GET("/client/parcels",ca(func(w http.ResponseWriter,req *http.Request){parcels,_,_:=ps.List(req.Context(),1,0,100);execTpl(cTmpl,"parcels",w,"parcels.html",map[string]any{"Parcels":parcels,"Total":len(parcels)})}))
	r.POST("/client/predeclare",ca(func(w http.ResponseWriter,req *http.Request){
		req.ParseForm()
		tn:=req.FormValue("tracking_number")
		if tn==""{execTpl(cTmpl,"predeclare",w,"predeclare.html",map[string]any{"Error":"快递单号必填"});return}
		whID:=int64(1)
		if v:=req.FormValue("warehouse_id");v!=""{if id,err:=strconv.ParseInt(v,10,64);err==nil{whID=id}}
		ps.PreDeclare(req.Context(),&parcelDomain.Parcel{TrackingNumber:tn,ProductName:req.FormValue("product_name"),
			TenantID:1,WarehouseID:whID,CourierCode:req.FormValue("courier_code")})
		w.Header().Set("HX-Redirect","/client/parcels")
		w.WriteHeader(200)
	}))

	// P0/P1 route registration complete
	r.GET("/client/logout",func(w http.ResponseWriter,req *http.Request){http.SetCookie(w,&http.Cookie{Name:"client_token",Value:"",Path:"/client",MaxAge:-1});http.Redirect(w,req,"/client/login",303)})
	registerClientP01Routes(r, cTmpl, ps, rr, ca, weightRepo, osvc, lr, ws, cour, dr, mr, sr, whr, ar, rpr, dfr, scr, acr)
	_=pr

	// ==========================================
	// Client template initialization
}

func seed(cr *custRepo.MemClientRepo,wr *whRepo.MemWarehouseRepo,rr *tmsRepo.MemRouteRepo,pr *parcelRepo.MemParcelRepo,or *orderRepo.MemOrderRepo,_ *tmsRepo.MemCourierRepo,lr *custRepo.MemLedgerRepo,sr *psRepo.MemServiceRepo,wor *woRepo.MemWorkOrderRepo,_ *printRepo.MemPrintRepo,whr *whRepo2.MemWebhookRepo,rpt *reportDomain.ReportService){
	ctx:=context.Background();now:=time.Now()
	wr.Create(ctx,1,&whDomain.Warehouse{Name:"厦门仓",Code:"XM",Address:"福建省厦门市集美区",Contact:"仓库管理员",Phone:"0592-1234567",IsActive:true,TenantID:1})
	c:=&custDomain.Client{Name:"EZ集运通",Code:"EZ001",ClientType:custDomain.ClientTypePlatform,ContactName:"运营经理",ContactPhone:"13800001111",ContactEmail:"ez@example.com",Balance:10000,IsActive:true,TenantID:1}
	cr.Create(ctx,1,c)
	// Seed routes for real route data
	rr.Create(ctx,&tmsDomain.Route{TenantID:1,WarehouseID:1,Name:"厦门→台湾(空运)",TransportType:"air",MinWeight:0.5,VolumeCoeff:6000,BaseWeightPrice:20.0,BaseVolumePrice:20.0,MinAmount:50,MinDays:1,MaxDays:3,IsActive:true})
	rr.Create(ctx,&tmsDomain.Route{TenantID:1,WarehouseID:1,Name:"厦门→台湾(海快)",TransportType:"sea_express",MinWeight:1.0,VolumeCoeff:6000,BaseWeightPrice:8.30,BaseVolumePrice:15.0,MinAmount:50,MinDays:3,MaxDays:7,IsActive:true})
	rr.Create(ctx,&tmsDomain.Route{TenantID:1,WarehouseID:1,Name:"厦门→台湾(海运)",TransportType:"sea",MinWeight:10.0,VolumeCoeff:6000,BaseWeightPrice:3.20,BaseVolumePrice:10.0,MinAmount:50,MinDays:5,MaxDays:14,IsActive:true})
	for i,pd:=range[]struct{tn,pn string;s parcelDomain.ParcelStatus;w float64}{{"SF1234567890","手机壳","pre_declared",0.12},{"ZTO9876543210","运动鞋","received",0.80},{"YTO1111222233","T恤","weighed",0.25},{"STO4444555566","蓝牙耳机","stored",0.15},{"HTKY7777888899","数据线","stored",0.08},{"JD9999000011","充电宝","stored",0.30},{"EMS1213141516","化妆品套装","stored",1.20}}{
		tn:=now.Add(-time.Duration(10-i)*24*time.Hour)
		pr.Create(ctx,&parcelDomain.Parcel{TenantID:1,WarehouseID:1,ClientID:c.ID,TrackingNumber:pd.tn,ProductName:pd.pn,ParcelName:pd.pn,Status:parcelDomain.ParcelStatus(pd.s),CourierCode:"SF",CargoType:"general",ActualWeight:pd.w,CreatedAt:tn,UpdatedAt:tn})
	}
	// Seed 8 orders spread across last 7 days with real OrderNo, weights, prices
	type orderSeed struct {
		orderNo       string
		memberID      int64
		routeID       int64
		recipient     string
		tracking      string
		status        orderDomain.OrderStatus
		weight        float64
		chgWeight     float64
		price         float64
		daysAgo       int
		parcelCount   int
		carrierTrack  string
		customsNo     string
		remark        string
	}
	today := now
	orders := []orderSeed{
		{"ORD-20260711-001", 1, 2, "王仁照", "80020737681100020001", orderDomain.StatusInTransit, 0.56, 0.60, 8.00, 0, 1, "CT-8837291", "CN-20260711001", "空运急件"},
		{"ORD-20260711-002", 2, 2, "琦立工作室", "YT7631606603205", orderDomain.StatusPendingLoading, 1.05, 1.50, 18.00, 0, 2, "", "", ""},
		{"ORD-20260710-001", 1, 1, "张致廷", "HTKY7777888899,JD9999000011", orderDomain.StatusPendingPacking, 0.33, 0.50, 11.50, 1, 2, "", "", ""},
		{"ORD-20260709-001", 2, 3, "吴欣如", "ZTO20250601001,SF120011223344", orderDomain.StatusCompleted, 12.80, 15.00, 56.20, 2, 3, "CT-8837292", "CN-20260709001", "已签收"},
		{"ORD-20260708-001", 1, 2, "王仁照", "YTO8822110011", orderDomain.StatusCustomsClearance, 2.30, 2.50, 22.00, 3, 1, "CT-8837293", "CN-20260708001", ""},
		{"ORD-20260707-001", 2, 1, "琦立工作室", "STO5555666677", orderDomain.StatusLoaded, 4.50, 5.00, 45.00, 4, 2, "", "", "大件运输"},
		{"ORD-20260706-001", 1, 3, "张致廷", "EMS9988776655,EMS1122334455", orderDomain.StatusShipped, 28.50, 30.00, 98.00, 5, 4, "CT-8837294", "CN-20260706001", ""},
		{"ORD-20260705-001", 2, 2, "吴欣如", "SF5566778899,YTO4433221100", orderDomain.StatusPendingPicking, 0.78, 1.00, 9.50, 6, 2, "", "", "待拣货"},
	}
	for _, od := range orders {
		or.Create(ctx, &orderDomain.Order{
			TenantID: 1, WarehouseID: 1, ClientID: c.ID,
			OrderNo: od.orderNo, MemberID: od.memberID, RouteID: od.routeID,
			RecipientName: od.recipient, TrackingNumbers: od.tracking,
			Status: od.status, ParcelCount: od.parcelCount,
			TotalActualWeight: od.weight, TotalChargeableWeight: od.chgWeight,
			TotalPrice: od.price, CarrierTrackingNo: od.carrierTrack,
			CustomsNumber: od.customsNo, Remark: od.remark,
		})
		// Override the auto-set CreatedAt to match the seed date
		// (MemOrderRepo.Create sets CreatedAt=time.Now(), so we patch after)
	}
	// Patch seed order dates to match their order_no dates
	if o1, _ := or.GetByOrderNo(ctx, 1, "ORD-20260711-001"); o1 != nil { o1.CreatedAt = today; o1.UpdatedAt = today; or.Update(ctx, o1) }
	if o2, _ := or.GetByOrderNo(ctx, 1, "ORD-20260711-002"); o2 != nil { o2.CreatedAt = today; o2.UpdatedAt = today; or.Update(ctx, o2) }
	if o3, _ := or.GetByOrderNo(ctx, 1, "ORD-20260710-001"); o3 != nil { o3.CreatedAt = today.Add(-24 * time.Hour); o3.UpdatedAt = today.Add(-24 * time.Hour); or.Update(ctx, o3) }
	if o4, _ := or.GetByOrderNo(ctx, 1, "ORD-20260709-001"); o4 != nil { o4.CreatedAt = today.Add(-2 * 24 * time.Hour); o4.UpdatedAt = today.Add(-2 * 24 * time.Hour); or.Update(ctx, o4) }
	if o5, _ := or.GetByOrderNo(ctx, 1, "ORD-20260708-001"); o5 != nil { o5.CreatedAt = today.Add(-3 * 24 * time.Hour); o5.UpdatedAt = today.Add(-3 * 24 * time.Hour); or.Update(ctx, o5) }
	if o6, _ := or.GetByOrderNo(ctx, 1, "ORD-20260707-001"); o6 != nil { o6.CreatedAt = today.Add(-4 * 24 * time.Hour); o6.UpdatedAt = today.Add(-4 * 24 * time.Hour); or.Update(ctx, o6) }
	if o7, _ := or.GetByOrderNo(ctx, 1, "ORD-20260706-001"); o7 != nil { o7.CreatedAt = today.Add(-5 * 24 * time.Hour); o7.UpdatedAt = today.Add(-5 * 24 * time.Hour); or.Update(ctx, o7) }
	if o8, _ := or.GetByOrderNo(ctx, 1, "ORD-20260705-001"); o8 != nil { o8.CreatedAt = today.Add(-6 * 24 * time.Hour); o8.UpdatedAt = today.Add(-6 * 24 * time.Hour); or.Update(ctx, o8) }
	lr.Add(ctx,&custRepo.LedgerEntry{TenantID:1,ClientID:c.ID,Amount:5000,BalanceAfter:5000,Type:"recharge",Description:""})
	_=sr;_=wor;_=whr;_=rpt
}
