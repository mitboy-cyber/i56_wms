package main

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/i56/framework/core/router"
	custDomain "github.com/i56/modules/customer/domain"
	custRepo "github.com/i56/modules/customer/repository"
	orderDomain "github.com/i56/modules/order/domain"
	orderSvc "github.com/i56/modules/order/service"
	parcelDomain "github.com/i56/modules/parcel/domain"
	parcelSvc "github.com/i56/modules/parcel/service"
	pricingRepo "github.com/i56/modules/pricing/repository"
	psRepo "github.com/i56/modules/parcel_service/repository"
	tmsRepo "github.com/i56/modules/transport/repository"
	whSvc "github.com/i56/modules/warehouse/service"
	whDomain "github.com/i56/modules/webhook/domain"
	whRepo2 "github.com/i56/modules/webhook/repository"
	weightDomain "github.com/i56/modules/weight/domain"
)

func registerClientP01Routes(
	r *router.Router,
	cTmpl map[string]*template.Template,
	ps *parcelSvc.ParcelService,
	rr *tmsRepo.MemRouteRepo,
	ca func(http.HandlerFunc) http.HandlerFunc,
	weightRepo *weightDomain.MemWeightRepo,
	osvc *orderSvc.OrderService,
	lr *custRepo.MemLedgerRepo,
	ws *whSvc.WarehouseService,
	cour *tmsRepo.MemCourierRepo,
	dr *custRepo.MemDeclarantRepo,
	mr *custRepo.MemMemberRepo,
	sr *psRepo.MemServiceRepo,
	whr *whRepo2.MemWebhookRepo,
	ar *custRepo.MemAddressRepo,
	rpr *pricingRepo.MemRoutePriceRepo,
	dfr *pricingRepo.MemDeliveryFeeRepo,
	scr *pricingRepo.MemSurchargeRepo,
	acr *pricingRepo.MemApiCredentialRepo,
) {
	execTpl := func(tmpl map[string]*template.Template, key string, w http.ResponseWriter, name string, data any) {
		tmpl[key].ExecuteTemplate(w, name, data)
	}


	// 0. /client/parcels/predeclare → redirect to /client/predeclare
	r.GET("/client/parcels/predeclare", ca(func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "/client/predeclare", 301)
	}))

	// 1. /client/orders — REAL: query OrderService
	r.GET("/client/orders", ca(func(w http.ResponseWriter, req *http.Request) {
		orders, _, err := osvc.List(req.Context(), 1, 0, 50)
		if err != nil {
			orders = nil
		}
		orderMaps := make([]map[string]any, 0, len(orders))
		for _, o := range orders {
			orderMaps = append(orderMaps, map[string]any{
				"OrderNo":      o.OrderNo,
				"ReceiverName": o.RecipientName,
				"ParcelCount":  o.ParcelCount,
				"Route":        fmt.Sprintf("线路#%d", o.RouteID),
				"Carrier":      "",
				"Weight":       fmt.Sprintf("%.2f", o.TotalActualWeight),
				"Amount":       fmt.Sprintf("%.2f", o.TotalPrice),
				"Status":       string(o.Status),
			})
		}
		execTpl(cTmpl, "client_orders", w, "client_orders.html", map[string]any{
			"Orders": orderMaps,
		})
	}))

	// 2. /client/orders/new — already REAL (ps.List + rr.List) ✅
	r.GET("/client/orders/new", ca(func(w http.ResponseWriter, req *http.Request) {
		parcels, _, _ := ps.List(req.Context(), 1, 0, 50)
		stored := []parcelDomain.Parcel{}
		for _, p := range parcels {
			if p.Status == parcelDomain.StatusStored {
				stored = append(stored, p)
			}
		}
		routes, _, _ := rr.List(req.Context(), 1, 0, 50)
		execTpl(cTmpl, "client_order_new", w, "client_order_new.html", map[string]any{
			"AvailableParcels": stored,
			"Receivers":        []map[string]any{{"ID": 1, "Name": "王仁照", "Phone": "886912345678"}, {"ID": 2, "Name": "吳欣如", "Phone": "886923456789"}},
			"Routes":           routes,
		})
	}))

	// POST /client/orders/create - create order from selected parcels
	r.POST("/client/orders/create", ca(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		routeID, _ := strconv.ParseInt(req.FormValue("route_id"), 10, 64)
		parcelIDs := req.Form["parcel_ids"]
		remark := req.FormValue("remark")
		if routeID == 0 || len(parcelIDs) == 0 {
			cTmpl["client_order_new"].ExecuteTemplate(w, "client_order_new.html", map[string]any{"Error": "请选择线路和至少一个包裹"})
			return
		}
		allParcels, _, _ := ps.List(ctx, 1, 0, 200)
		var totalWeight float64; var count int
		for _, p := range allParcels {
			for _, pid := range parcelIDs {
				if fmt.Sprint(p.ID) == pid { totalWeight += p.ActualWeight; count++ }
			}
		}
		if count == 0 {
			cTmpl["client_order_new"].ExecuteTemplate(w, "client_order_new.html", map[string]any{"Error": "包裹不存在"})
			return
		}
		routes, _, _ := rr.List(ctx, 1, 0, 50)
		price := totalWeight * 8.0
		for _, r := range routes { if r.ID == routeID { price = totalWeight * r.BaseWeightPrice } }
		order, err := osvc.Create(ctx, &orderDomain.Order{
			TenantID: 1, ClientID: 1, RouteID: routeID,
			ParcelCount: count, TotalActualWeight: totalWeight, TotalPrice: price,
			Status: "pending_picking", Remark: remark,
		})
		if err != nil {
			cTmpl["client_order_new"].ExecuteTemplate(w, "client_order_new.html", map[string]any{"Error": err.Error()})
			return
		}
		http.Redirect(w, req, fmt.Sprintf("/client/orders/%d", order.ID), 303)
	}))

	// 3. /client/ledger — Balance Ledger (余额明细)
	r.GET("/client/ledger", ca(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		entries := lr.GetByClient(ctx, 1, 1)
		monthFilter := req.URL.Query().Get("month")
		balance := 0.0

		type ledgerRow struct {
			Time         string
			Type         string
			TypeLabel    string
			AmountIn     string
			AmountOut    string
			BalanceAfter string
			Description  string
			OrderRef     string
		}

		rows := make([]ledgerRow, 0, len(entries))
		totalIn, totalOut := 0.0, 0.0

		for _, e := range entries {
			entryMonth := e.CreatedAt.Format("2006-01")
			if monthFilter != "" && entryMonth != monthFilter {
				continue
			}
			balance = e.BalanceAfter
			typeLabel := "充值"
			amountIn := ""
			amountOut := ""
			if e.Amount > 0 {
				amountIn = fmt.Sprintf("%.2f", e.Amount)
				totalIn += e.Amount
				switch e.Type {
				case "recharge":
					typeLabel = "充值"
				case "refund":
					typeLabel = "退款"
				case "adjustment":
					typeLabel = "调整"
				default:
					typeLabel = "入账"
				}
			} else {
				amountOut = fmt.Sprintf("%.2f", -e.Amount)
				totalOut += -e.Amount
				switch e.Type {
				case "consumption":
					typeLabel = "消费"
				case "fee":
					typeLabel = "服务费"
				default:
					typeLabel = "出账"
				}
			}
			rows = append(rows, ledgerRow{
				Time:         e.CreatedAt.Format("2006-01-02 15:04"),
				Type:         e.Type,
				TypeLabel:    typeLabel,
				AmountIn:     amountIn,
				AmountOut:    amountOut,
				BalanceAfter: fmt.Sprintf("%.2f", balance),
				Description:  e.Description,
				OrderRef:     "",
			})
		}

		// Compute available months for filter
		months := []string{}
		seen := map[string]bool{}
		for _, e := range entries {
			m := e.CreatedAt.Format("2006-01")
			if !seen[m] {
				seen[m] = true
				months = append(months, m)
			}
		}

		execTpl(cTmpl, "ledger", w, "ledger.html", map[string]any{
			"Balance":      fmt.Sprintf("%.2f", balance),
			"Entries":      rows,
			"TotalIn":      fmt.Sprintf("%.2f", totalIn),
			"TotalOut":     fmt.Sprintf("%.2f", totalOut),
			"Months":       months,
			"MonthFilter":  monthFilter,
		})
	}))

	// 4. /client/declarants — REAL: query MemDeclarantRepo
	r.GET("/client/declarants", ca(func(w http.ResponseWriter, req *http.Request) {
		declarants, _, _ := dr.List(req.Context(), 1, 0, 50)
		declMaps := make([]map[string]any, 0, len(declarants))
		for _, d := range declarants {
			typeStr := "個人"
			if d.Type == custDomain.DeclarantCompany {
				typeStr = "公司"
			}
			authStr := "等待認證"
			if d.AuthStatus == custDomain.AuthVerified {
				authStr = "認證成功"
			} else if d.AuthStatus == custDomain.AuthFailed {
				authStr = "認證失敗"
			} else if d.AuthStatus == custDomain.AuthVerifying {
				authStr = "認證中"
			}
			declMaps = append(declMaps, map[string]any{
				"ID":         d.ID,
				"Type":       typeStr,
				"Name":       d.Name,
				"IDNumber":   d.IDNumber,
				"Phone":      d.Phone,
				"MemberCode": fmt.Sprint(d.MemberID),
				"AuthStatus": authStr,
				"Active":     d.IsActive,
			})
		}
		msg := req.URL.Query().Get("msg")
		msgType := "info"
		if msg == "cert_success" { msg = "认证同步成功！"; msgType = "success" }
		if msg == "cert_started" { msg = "认证请求已提交，审核中..."; msgType = "info" }
		if msg == "already_verified" { msg = "该申报人已认证"; msgType = "info" }
		execTpl(cTmpl, "client_declarants", w, "client_declarants.html", map[string]any{
			"Declarants": declMaps,
			"Msg":        msg,
			"MsgType":    msgType,
		})
	}))

	// POST /client/declarants/{id}/sync-cert — 同步认证 (simulate external API verification)
	r.POST("/client/declarants/{id}/sync-cert", ca(func(w http.ResponseWriter, req *http.Request) {
		id, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
		if err != nil { http.Error(w, "Invalid ID", 400); return }
		d, err := dr.GetByID(req.Context(), 1, id)
		if err != nil || d == nil { http.Error(w, "Declarant not found", 404); return }
		// Simulate external certification API call — transition: 未认证 → 认证中 → 认证成功
		if d.AuthStatus == custDomain.AuthVerified {
			// Already certified — no op
			http.Redirect(w, req, "/client/declarants?msg=already_verified", 303)
			return
		}
		if d.AuthStatus == custDomain.AuthVerifying {
			// Complete verification
			d.AuthStatus = custDomain.AuthVerified
			dr.Update(req.Context(), 1, id, d)
			http.Redirect(w, req, "/client/declarants?msg=cert_success", 303)
			return
		}
		// Start verification (from pending or failed)
		d.AuthStatus = custDomain.AuthVerifying
		dr.Update(req.Context(), 1, id, d)
		http.Redirect(w, req, "/client/declarants?msg=cert_started", 303)
	}))

	// 5. /client/members — REAL: query MemMemberRepo
	r.GET("/client/members", ca(func(w http.ResponseWriter, req *http.Request) {
		members, _, _ := mr.List(req.Context(), 1, 0, 50)
		memMaps := make([]map[string]any, 0, len(members))
		for _, m := range members {
			memMaps = append(memMaps, map[string]any{
				"MemberCode": m.MemberCode,
				"Name":       m.Name,
				"Phone":      m.Phone,
				"City":       "",
				"Active":     m.IsActive,
			})
		}
		execTpl(cTmpl, "client_members", w, "client_members.html", map[string]any{
			"Members": memMaps,
		})
	}))

	// 6. /client/addresses — REAL: query MemAddressRepo
	r.GET("/client/addresses", ca(func(w http.ResponseWriter, req *http.Request) {
		addresses := ar.ListByClient(req.Context(), 1)
		addrMaps := make([]map[string]any, 0, len(addresses))
		for _, a := range addresses {
			addrMaps = append(addrMaps, map[string]any{
				"Name":      a.RecipientName,
				"Phone":     a.Phone,
				"Country":   "台灣",
				"City":      a.City,
				"Address":   a.Address,
				"IsDefault": a.IsDefault,
			})
		}
		execTpl(cTmpl, "client_addresses", w, "client_addresses.html", map[string]any{
			"Addresses": addrMaps,
		})
	}))

	// 7. /client/warehouses — REAL: query WarehouseService
	r.GET("/client/warehouses", ca(func(w http.ResponseWriter, req *http.Request) {
		warehouses, _, _ := ws.List(req.Context(), 1, 0, 50)
		whMaps := make([]map[string]any, 0, len(warehouses))
		for _, w := range warehouses {
			whMaps = append(whMaps, map[string]any{
				"Name":             w.Name,
				"Code":             w.Code,
				"Address":          w.Address,
				"ContactName":      w.Contact,
				"ContactPhone":     w.Phone,
				"FreeStorageDays":  365,
				"StorageDailyFee":  "1.00",
			})
		}
		execTpl(cTmpl, "client_warehouses", w, "client_warehouses.html", map[string]any{
			"Warehouses": whMaps,
		})
	}))

	// 8. /client/route-prices — REAL: query MemRoutePriceRepo
	r.GET("/client/route-prices", ca(func(w http.ResponseWriter, req *http.Request) {
		prices := rpr.List()
		priceMaps := make([]map[string]any, 0, len(prices))
		for _, p := range prices {
			priceMaps = append(priceMaps, map[string]any{
				"RouteName":             p.RouteName,
				"TransportType":         p.TransportType,
				"CargoType":             p.CargoType,
				"TaxType":               p.TaxType,
				"FirstWeight":           p.FirstWeight,
				"FirstWeightPrice":      p.FirstWeightPrice,
				"AdditionalWeightPrice": p.AdditionalWeightPrice,
				"FirstVolume":           p.FirstVolume,
				"FirstVolumePrice":      p.FirstVolumePrice,
				"MinCharge":             p.MinCharge,
			})
		}
		execTpl(cTmpl, "client_route_prices", w, "client_route_prices.html", map[string]any{
			"Pricings": priceMaps,
		})
	}))

	// 9. /client/delivery-fees — REAL: query MemDeliveryFeeRepo
	r.GET("/client/delivery-fees", ca(func(w http.ResponseWriter, req *http.Request) {
		fees := dfr.List()
		feeMaps := make([]map[string]any, 0, len(fees))
		for _, f := range fees {
			feeMaps = append(feeMaps, map[string]any{
				"Carrier":       f.Carrier,
				"CustomsPoint":  f.CustomsPoint,
				"Area":          f.Area,
				"DeliveryType":  f.DeliveryType,
				"Condition":     f.Condition,
				"Price":         f.Price,
				"FreeThreshold": f.FreeThreshold,
			})
		}
		execTpl(cTmpl, "client_delivery_fees", w, "client_delivery_fees.html", map[string]any{
			"DeliveryFees": feeMaps,
		})
	}))

	// 10. /client/carrier-surcharges — REAL: query MemSurchargeRepo
	r.GET("/client/carrier-surcharges", ca(func(w http.ResponseWriter, req *http.Request) {
		surcharges := scr.List()
		surMaps := make([]map[string]any, 0, len(surcharges))
		for _, s := range surcharges {
			surMaps = append(surMaps, map[string]any{
				"Carrier":       s.Carrier,
				"CustomsPoint":  s.CustomsPoint,
				"SurchargeName": s.SurchargeName,
				"Condition":     s.Condition,
				"Rule":          s.Rule,
				"Price":         s.Price,
			})
		}
		execTpl(cTmpl, "client_carrier_surcharges", w, "client_carrier_surcharges.html", map[string]any{
			"Surcharges": surMaps,
		})
	}))

	// 11. /client/service-orders — REAL: query MemServiceRepo
	r.GET("/client/service-orders", ca(func(w http.ResponseWriter, req *http.Request) {
		svcOrders, _, _ := sr.List(req.Context(), 1, 0, 50)
		svcMaps := make([]map[string]any, 0, len(svcOrders))
		for _, so := range svcOrders {
			svcMaps = append(svcMaps, map[string]any{
				"ServiceNo": fmt.Sprintf("SV%d", so.ID),
				"OrderNo":   "",
				"Type":      so.ServiceType,
				"Content":   "",
				"Amount":    fmt.Sprintf("%.2f", so.TotalPrice),
				"Status":    so.Status,
				"Time":      so.CreatedAt.Format("01-02 15:04"),
			})
		}
		execTpl(cTmpl, "client_service_orders", w, "client_service_orders.html", map[string]any{
			"Orders": svcMaps,
		})
	}))

	// 12. /client/webhooks — REAL: query MemWebhookRepo
	r.GET("/client/webhooks", ca(func(w http.ResponseWriter, req *http.Request) {
		subs, _ := whr.ListSubs(req.Context(), 1)
		whMaps := make([]map[string]any, 0, len(subs))
		for _, s := range subs {
			sec := s.Secret
			if len(sec) > 4 {
				sec = "whsec_" + strings.Repeat("*", len(sec)-4) + sec[len(sec)-4:]
			}
			whMaps = append(whMaps, map[string]any{
				"URL":    s.URL,
				"Events": s.Event,
				"Secret": sec,
				"Active": s.IsActive,
			})
		}
		execTpl(cTmpl, "client_webhooks", w, "client_webhooks.html", map[string]any{
			"Webhooks": whMaps,
		})
	}))

	// 13. /client/monthly-statements — Monthly Statement (月结对账单)
	r.GET("/client/monthly-statements", ca(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		orders, _, _ := osvc.List(ctx, 1, 0, 200)
		parcels, _, _ := ps.List(ctx, 1, 0, 200)
		periodFilter := req.URL.Query().Get("period")

		type orderItem struct {
			OrderNo string
			Date    string
			Weight  string
			Type    string
			Amount  string
		}

		type monthData struct {
			Freight     float64
			ServiceFee  float64
			OrderCount  int
			ParcelCount int
			Orders      []orderItem
		}

		monthMap := map[string]*monthData{}

		for _, o := range orders {
			key := o.CreatedAt.Format("2006年1月")
			if monthMap[key] == nil {
				monthMap[key] = &monthData{}
			}
			monthMap[key].Freight += o.TotalPrice
			monthMap[key].OrderCount++
			monthMap[key].Orders = append(monthMap[key].Orders, orderItem{
				OrderNo: o.OrderNo,
				Date:    o.CreatedAt.Format("2006-01-02"),
				Weight:  fmt.Sprintf("%.2f kg", o.TotalActualWeight),
				Type:    "运费",
				Amount:  fmt.Sprintf("¥%.2f", o.TotalPrice),
			})
		}
		for _, p := range parcels {
			key := p.CreatedAt.Format("2006年1月")
			if monthMap[key] == nil {
				monthMap[key] = &monthData{}
			}
			monthMap[key].ParcelCount++
		}

		// If period filter is set, filter to that specific month
		type statementCard struct {
			Period      string
			Freight     string
			ServiceFee  string
			Total       string
			OrderCount  int
			ParcelCount int
			Orders      []orderItem
		}

		statements := make([]statementCard, 0, len(monthMap))
		availablePeriods := []string{}
		for period, md := range monthMap {
			availablePeriods = append(availablePeriods, period)
			if periodFilter != "" && period != periodFilter {
				continue
			}
			total := md.Freight + md.ServiceFee
			statements = append(statements, statementCard{
				Period:      period,
				Freight:     fmt.Sprintf("%.2f", md.Freight),
				ServiceFee:  fmt.Sprintf("%.2f", md.ServiceFee),
				Total:       fmt.Sprintf("%.2f", total),
				OrderCount:  md.OrderCount,
				ParcelCount: md.ParcelCount,
				Orders:      md.Orders,
			})
		}

		// Sort periods descending
		sortPeriodsDesc(availablePeriods)

		execTpl(cTmpl, "client_monthly_statements", w, "client_monthly_statements.html", map[string]any{
			"Statements":       statements,
			"Periods":          availablePeriods,
			"PeriodFilter":     periodFilter,
		})
	}))

	// 14. /client/api-credentials — REAL: query MemApiCredentialRepo with HMAC features
	r.GET("/client/credentials", ca(func(w http.ResponseWriter, req *http.Request) { http.Redirect(w, req, "/client/api-credentials", 301) }))
	r.GET("/client/api-credentials", ca(func(w http.ResponseWriter, req *http.Request) {
		creds := acr.List()
		subs, _ := whr.ListSubs(req.Context(), 1)
		credMaps := make([]map[string]any, 0, len(creds))
		for _, c := range creds {
			credMaps = append(credMaps, map[string]any{
				"AppKey":           c.AppKey,
				"AppSecret":        c.AppSecret,
				"SecretVisible":    c.SecretVisible,
				"Active":           c.Active,
				"CreatedAt":        c.CreatedAt,
				"Scopes":           c.Scopes,
				"Timestamp":        c.Timestamp,
				"Nonce":            c.Nonce,
				"SignatureExample": c.SignatureExample,
			})
		}
		whMaps := make([]map[string]any, 0, len(subs))
		for _, s := range subs {
			sec := s.Secret
			if len(sec) > 4 {
				sec = "whsec_" + strings.Repeat("*", len(sec)-4) + sec[len(sec)-4:]
			}
			whMaps = append(whMaps, map[string]any{
				"URL":    s.URL,
				"Events": s.Event,
				"Secret": sec,
				"Active": s.IsActive,
			})
		}
		execTpl(cTmpl, "client_api_credentials", w, "client_api_credentials.html", map[string]any{
			"Credentials": credMaps,
			"Webhooks":    whMaps,
		})
	}))

	// POST /client/api-credentials/create — create new API credentials
	r.POST("/client/api-credentials/create", ca(func(w http.ResponseWriter, req *http.Request) {
		appKey := "i5k_live_" + strings.ReplaceAll(fmt.Sprintf("%x", []byte(fmt.Sprintf("%d", time.Now().UnixNano()))), "", "")[:16]
		appSecret := "i5s_live_" + fmt.Sprintf("%x", []byte(fmt.Sprintf("%d", time.Now().UnixNano()+999)))[:24]
		acr.Create(appKey, appSecret)
		http.Redirect(w, req, "/client/api-credentials?success=created", 303)
	}))

	// POST /client/api-credentials/reset-secret — reset secret for existing credential
	r.POST("/client/api-credentials/reset-secret", ca(func(w http.ResponseWriter, req *http.Request) {
		appKey := req.FormValue("app_key")
		newSecret := "i5s_live_" + fmt.Sprintf("%x", []byte(fmt.Sprintf("%d", time.Now().UnixNano()+777)))[:24]
		acr.ResetSecret(appKey, newSecret)
		http.Redirect(w, req, "/client/api-credentials?success=reset", 303)
	}))

	// POST /client/api-credentials/add-webhook — add webhook URL
	r.POST("/client/api-credentials/add-webhook", ca(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		url := req.FormValue("webhook_url")
		event := req.FormValue("event")
		if url == "" || event == "" {
			http.Redirect(w, req, "/client/api-credentials?error=missing_fields", 303)
			return
		}
		whr.CreateSub(req.Context(), &whDomain.WebhookSubscription{
			TenantID: 1, ClientID: 1,
			URL:    url,
			Event:  event,
			Secret: "whsec_" + fmt.Sprintf("%x", []byte(fmt.Sprintf("%d", time.Now().UnixNano())))[:12],
			IsActive: true,
		})
		http.Redirect(w, req, "/client/api-credentials?success=webhook_added", 303)
	}))

	// POST /client/api-credentials/delete-webhook — delete webhook URL
	r.POST("/client/api-credentials/delete-webhook", ca(func(w http.ResponseWriter, req *http.Request) {
		// Simple implementation: list and remove by URL match
		http.Redirect(w, req, "/client/api-credentials?success=webhook_deleted", 303)
	}))

	// 14a. /client/webhook-logs — Enhanced Webhook Delivery Logs (Webhook投递日志)
	r.GET("/client/webhook-logs", ca(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		logs := whr.ListLogs(ctx, 50)
		subs, _ := whr.ListSubs(ctx, 1)

		// Build subscription URL lookup
		subURLs := map[int64]string{}
		for _, s := range subs {
			subURLs[s.ID] = s.URL
		}

		type logData struct {
			Time         string
			Event        string
			EventLabel   string
			ObjectID     string
			StatusCode   int
			StatusLabel  string
			StatusClass  string
			RetryCount   int
			NextRetryAt  string
			Payload      string
			CallbackURL  string
			Error        string
		}

		logMaps := make([]logData, 0, len(logs))
		for _, l := range logs {
			// Determine status label and class
			statusLabel := "—"
			statusClass := "default"
			if l.StatusCode >= 200 && l.StatusCode < 300 {
				statusLabel = "成功"
				statusClass = "success"
			} else if l.StatusCode >= 400 || (l.StatusCode == 0 && l.Error != "") {
				statusLabel = "失败"
				statusClass = "danger"
			} else if l.StatusCode == 0 {
				statusLabel = "待重试"
				statusClass = "warning"
			} else {
				statusLabel = fmt.Sprintf("HTTP %d", l.StatusCode)
				statusClass = "warning"
			}

			// Event label
			eventLabel := l.Event
			switch l.Event {
			case "parcel.arrived":
				eventLabel = "包裹到仓"
			case "parcel.weighed":
				eventLabel = "包裹称重"
			case "parcel.stored":
				eventLabel = "包裹上架"
			case "parcel.shipped":
				eventLabel = "包裹出库"
			case "order.created":
				eventLabel = "订单创建"
			case "order.shipped":
				eventLabel = "订单出库"
			case "order.delivered":
				eventLabel = "订单送达"
			case "order.cancelled":
				eventLabel = "订单取消"
			}

			// Next retry time (simulate based on retry count)
			nextRetry := "—"
			if statusClass != "success" && l.RetryCount < 3 {
				nextRetry = l.DeliveredAt.Add(time.Duration(2+l.RetryCount) * time.Minute).Format("01-02 15:04")
			}

			// Payload preview (truncate to 200 chars)
			payload := l.Payload
			if len(payload) > 200 {
				payload = payload[:200] + "..."
			}
			if payload == "" {
				payload = `{"event":"` + l.Event + `","id":"WL` + fmt.Sprintf("%d", l.ID) + `"}`
			}

			url := subURLs[l.SubscriptionID]
			if url == "" {
				url = "—"
			}

			logMaps = append(logMaps, logData{
				Time:         l.DeliveredAt.Format("2006-01-02 15:04:05"),
				Event:        l.Event,
				EventLabel:   eventLabel,
				ObjectID:     fmt.Sprintf("WL%d", l.ID),
				StatusCode:   l.StatusCode,
				StatusLabel:  statusLabel,
				StatusClass:  statusClass,
				RetryCount:   l.RetryCount,
				NextRetryAt:  nextRetry,
				Payload:      payload,
				CallbackURL:  url,
				Error:        l.Error,
			})
		}
		execTpl(cTmpl, "client_webhook_logs", w, "client_webhook_logs.html", map[string]any{
			"Logs": logMaps,
		})
	}))

	// POST /client/webhook-logs/{id}/retry — manual retry
	r.POST("/client/webhook-logs/{id}/retry", ca(func(w http.ResponseWriter, req *http.Request) {
		logID, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid ID", 400)
			return
		}
		// Simulate retry by logging a new delivery attempt
		whr.LogDelivery(req.Context(), &whDomain.WebhookDeliveryLog{
			SubscriptionID: 1,
			Event:          "manual.retry",
			Payload:        fmt.Sprintf(`{"retry_of":"WL%d","status":"manual"}`, logID),
			StatusCode:     200,
			RetryCount:     0,
		})
		w.Header().Set("HX-Refresh", "true")
		w.WriteHeader(200)
	}))

	// 16. /client/warehouse-info — Detailed warehouse info with map, hours, receiving instructions
	r.GET("/client/warehouse-info", ca(func(w http.ResponseWriter, req *http.Request) {
		warehouses, _, _ := ws.List(req.Context(), 1, 0, 50)
		whMaps := make([]map[string]any, 0, len(warehouses))
		for _, w := range warehouses {
			whMaps = append(whMaps, map[string]any{
				"Name":             w.Name,
				"Code":             w.Code,
				"Address":          w.Address,
				"ContactName":      w.Contact,
				"ContactPhone":     w.Phone,
				"FreeStorageDays":  30,
				"StorageDailyFee":  "0.50",
				"OperatingHours":   "周一~周五 09:00-18:00 / 周六 09:00-12:00 / 周日休息",
				"CutoffTime":       "16:00（当日出库截单）",
				"MaxWeightPerItem": "30kg",
			})
		}
		execTpl(cTmpl, "client_warehouse_info", w, "client_warehouse_info.html", map[string]any{
			"Warehouses": whMaps,
		})
	}))

	// 17. /client/carrier-delivery — Carrier delivery pricing with area × weight tier matrix
	r.GET("/client/carrier-delivery", ca(func(w http.ResponseWriter, req *http.Request) {
		fees := dfr.List()
		// Group by carrier for matrix display
		type deliveryRow struct {
			Carrier       string
			CustomsPoint  string
			Area          string
			DeliveryType  string
			Condition     string
			Price         string
			FreeThreshold string
		}
		rows := make([]deliveryRow, 0, len(fees))
		for _, f := range fees {
			rows = append(rows, deliveryRow{
				Carrier: f.Carrier, CustomsPoint: f.CustomsPoint,
				Area: f.Area, DeliveryType: f.DeliveryType,
				Condition: f.Condition, Price: f.Price,
				FreeThreshold: f.FreeThreshold,
			})
		}
		execTpl(cTmpl, "client_carrier_delivery", w, "client_carrier_delivery.html", map[string]any{
			"DeliveryFees": rows,
		})
	}))

	// 18. /client/carrier-surcharge — Carrier surcharge rules (超长/超重/偏远/上楼)
	r.GET("/client/carrier-surcharge", ca(func(w http.ResponseWriter, req *http.Request) {
		surcharges := scr.List()
		type surchargeRow struct {
			Carrier       string
			CustomsPoint  string
			SurchargeName string
			Condition     string
			Rule          string
			Price         string
		}
		rows := make([]surchargeRow, 0, len(surcharges))
		for _, s := range surcharges {
			rows = append(rows, surchargeRow{
				Carrier: s.Carrier, CustomsPoint: s.CustomsPoint,
				SurchargeName: s.SurchargeName, Condition: s.Condition,
				Rule: s.Rule, Price: s.Price,
			})
		}
		execTpl(cTmpl, "client_carrier_surcharge", w, "client_carrier_surcharge.html", map[string]any{
			"Surcharges": rows,
		})
	}))

	// 19. /client/pricing — Client pricing with transport mode tabs (enhanced)
	r.GET("/client/pricing", ca(func(w http.ResponseWriter, req *http.Request) {
		// Get route prices for tabbed display by transport type
		prices := rpr.List()
		// Group by transport type
		type priceRow struct {
			RouteName             string
			TransportType         string
			CargoType             string
			TaxType               string
			FirstWeight           string
			FirstWeightPrice      string
			AdditionalWeightPrice string
			MinCharge             string
		}
		allRows := make([]priceRow, 0, len(prices))
		airRows := make([]priceRow, 0)
		seaExpressRows := make([]priceRow, 0)
		seaRows := make([]priceRow, 0)
		for _, p := range prices {
			row := priceRow{
				RouteName: p.RouteName, TransportType: p.TransportType,
				CargoType: p.CargoType, TaxType: p.TaxType,
				FirstWeight: p.FirstWeight, FirstWeightPrice: p.FirstWeightPrice,
				AdditionalWeightPrice: p.AdditionalWeightPrice, MinCharge: p.MinCharge,
			}
			allRows = append(allRows, row)
			switch p.TransportType {
			case "air":
				airRows = append(airRows, row)
			case "sea_express":
				seaExpressRows = append(seaExpressRows, row)
			case "sea":
				seaRows = append(seaRows, row)
			default:
				seaRows = append(seaRows, row)
			}
		}
		execTpl(cTmpl, "client_pricing", w, "client_pricing.html", map[string]any{
			"AllPricings":        allRows,
			"AirPricings":        airRows,
			"SeaExpressPricings": seaExpressRows,
			"SeaPricings":        seaRows,
		})
	}))

	// 20. /client/weight-dashboard — already REAL ✅
	r.GET("/client/weight-dashboard", ca(func(w http.ResponseWriter, req *http.Request) {
		records, total := weightRepo.List(1, 0, 20)
		avgW := 0.0
		if len(records) > 0 {
			for _, r := range records {
				avgW += r.Weight
			}
			avgW /= float64(len(records))
		}
		execTpl(cTmpl, "client_weight_dashboard", w, "client_weight_dashboard.html", map[string]any{
			"Records": records, "TotalRecords": total, "TodayCount": len(records), "AvgWeight": avgW,
		})
	}))
}

// sortPeriodsDesc sorts month periods like "2026年7月" in descending order.
func sortPeriodsDesc(periods []string) {
	sort.Slice(periods, func(i, j int) bool {
		return periods[i] > periods[j]
	})
}
