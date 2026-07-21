// Package seed provides comprehensive demo data for all in-memory stores.
package seed

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	custDomain "github.com/i56/modules/customer/domain"
	custRepo "github.com/i56/modules/customer/repository"
	orderDomain "github.com/i56/modules/order/domain"
	orderRepo "github.com/i56/modules/order/repository"
	parcelDomain "github.com/i56/modules/parcel/domain"
	parcelRepo "github.com/i56/modules/parcel/repository"
	wmsDomain "github.com/i56/modules/wms/domain"
	wmsRepo "github.com/i56/modules/wms/repository"
)

// SeedAll populates all stores with realistic demo data.
func SeedAll(ctx context.Context,
	osvc *orderRepo.Service, psvc *parcelRepo.Service,
	cs *custRepo.Store, ms *custRepo.MemberStore,
	ds *custRepo.DeclarantStore, as *custRepo.AddressStore,
	ls *custRepo.LedgerStore, crs *custRepo.CredentialStore,
	rs *orderRepo.RouteStore, cos *orderRepo.CourierStore,
	scs *custRepo.SurchargeStore, pss *custRepo.PricingServiceStore,
	prs *custRepo.PricingRouteStore, ncs *custRepo.NotificationChannelStore,
	wss *wmsRepo.WorkflowStore, tss *wmsRepo.ServiceTemplateStore,
	sts *wmsRepo.ServiceTypeStore,
) {
	rng := rand.New(rand.NewSource(42))
	now := time.Now()

	// ── Courier companies ──
	couriers := []struct{ code, name string }{
		{"SF", "顺丰速运"}, {"YT", "圆通快递"}, {"ZTO", "中通快递"},
		{"YTO", "韵达快递"}, {"STO", "申通快递"}, {"JD", "京东物流"},
		{"EMS", "邮政EMS"}, {"DHL", "DHL国际快递"},
	}
	for i, c := range couriers {
		cos.Upsert(ctx, fmt.Sprintf("CP%03d", i+1), c.name, c.code)
	}

	// ── Routes (line templates) ──
	routes := []struct{ name, transportType, from, to string; basePrice float64 }{
		{"厦门→台湾(空运)", "air", "厦门", "桃园", 20},
		{"厦门→台湾(海快)", "sea_fast", "厦门", "基隆", 15},
		{"厦门→台湾(海运)", "sea", "厦门", "台北", 12},
		{"厦门→台湾(经济)", "sea", "厦门", "台中", 8},
		{"深圳→台湾(空运)", "air", "深圳", "桃园", 22},
	}
	for i, r := range routes {
		rs.Upsert(ctx, fmt.Sprintf("LINE-%03d", i+1), r.name, r.transportType, r.from, r.to, r.basePrice)
	}

	// ── Clients ──
	clients := []struct{ name, code string }{
		{"EZ集運通", "EZ001"}, {"琦立工作室", "QL001"}, {"跨境优选", "KJ001"},
		{"海淘达人", "HT001"}, {"台湾通商", "TW001"},
	}
	for i, c := range clients {
		cs.Create(ctx, 1, &custDomain.Client{Name: c.name, Code: c.code, ClientType: "company", ContactName: randName(rng), ContactPhone: randPhone(rng), IsActive: true, CreatedAt: now.AddDate(0, -i, 0)})
	}

	// ── Members ──
	memberNames := []string{"王仁照", "张致廷", "吴欣如", "陈小美", "林小明"}
	for i, name := range memberNames {
		ms.Create(ctx, 1, name, fmt.Sprintf("MB%05d", i+1), "1380000"+fmt.Sprintf("%04d", i+1))
	}

	// ── Declarants ──
	declarants := []struct{ typ, name, idCard, phone string }{
		{"personal", "魏立璿", "F131656261", "0981888909"},
		{"personal", "李采縈", "F224206468", "0958207309"},
		{"personal", "林昀叡", "A224010449", "0919141341"},
		{"personal", "邱筱玫", "H225539959", "0979587173"},
		{"personal", "许葳", "A227920577", "0968633815"},
	}
	for _, d := range declarants {
		ds.Create(ctx, 1, d.typ, d.name, d.idCard, d.phone)
	}

	// ── Addresses ──
	addresses := []struct{ name, phone, city, district, detail string }{
		{"王仁照", "13800000001", "台北市", "大安区", "忠孝东路四段100号"},
		{"张致廷", "13800000002", "新北市", "板桥区", "文化路一段200号"},
		{"吴欣如", "13800000003", "台中市", "西屯区", "台湾大道三段300号"},
	}
	for _, a := range addresses {
		as.Create(ctx, 1, a.name, a.phone, a.city, a.district, a.detail)
	}

	// ── Pricing Routes ──
	for i, r := range routes {
		prs.Upsert(ctx, fmt.Sprintf("PRICE-%03d", i+1), r.name, "首重", r.basePrice, r.basePrice*1.5, now.Format("2006-01-02"), now.AddDate(1, 0, 0).Format("2006-01-02"))
	}

	// ── Carriers / Shipping providers ──
	carriers := []struct{ name, code string }{ {"新竹物流", "HCT"}, {"黑猫宅急便", "YAMATO"}, {"全家超商", "FAMILY"}, {"7-11超商", "SEVEN"} }
	for i, c := range carriers {
		crs.Upsert(ctx, fmt.Sprintf("SC%03d", i+1), c.name, c.code)
	}

	// ── Surcharges ──
	surcharges := []struct{ name, typ string; amount float64 }{
		{"偏远地区附加费", "delivery", 50}, {"超大型件处理费", "size", 100},
		{"危险品处理费", "dangerous", 200}, {"保价费(3%)", "insurance", 0},
	}
	for i, s := range surcharges {
		scs.Upsert(ctx, fmt.Sprintf("SUR%03d", i+1), s.name, s.typ, s.amount)
	}

	// ── Pricing Services ──
	services := []struct{ name, typ string; price float64 }{
		{"标准打包", "packing", 3}, {"加固包装", "reinforced", 8},
		{"拍照验货", "inspection", 5}, {"合并包裹", "merge", 10},
	}
	for i, s := range services {
		pss.Upsert(ctx, fmt.Sprintf("SVC%03d", i+1), s.name, s.typ, s.price)
	}

	// ── Service Types ──
	serviceTypes := []string{"打包", "加固", "拍照", "合并", "退货", "销毁", "转寄", "拆分"}
	for i, st := range serviceTypes {
		sts.Upsert(ctx, fmt.Sprintf("STYPE%03d", i+1), st, fmt.Sprintf("附加服务-%s", st), 10+float64(i)*5)
	}

	// ── Service Templates ──
	svcTemplates := []struct{ name, desc string; price float64 }{
		{"标准打包服务", "提供纸箱+气泡膜专业打包", 0},
		{"易碎品加固", "双层包装+防震+易碎标签", 8},
		{"开箱验货拍照", "开箱检查并拍照存档", 5},
	}
	for i, t := range svcTemplates {
		tss.Upsert(ctx, fmt.Sprintf("TMPL%03d", i+1), t.name, t.desc, t.price)
	}

	// ── Workflows ──
	workflows := []struct{ name string; steps []string }{
		{"入库流程", []string{"收货", "核重", "上架", "入库完成"}},
		{"出库流程(大仓)", []string{"拣货", "打包", "核重", "送出库", "装柜"}},
		{"退件处理流程", []string{"签收", "核验", "入库", "上架"}},
	}
	for i, w := range workflows {
		wss.Upsert(ctx, fmt.Sprintf("WF%03d", i+1), w.name, w.steps)
	}

	// ── Notifications ──
	notifications := []struct{ title, typ, priority, scope, content, channel string }{
		{"系统维护通知", "系统通知", "普通", "全员(跨所有公司)", "系统将于7月20日凌晨2:00-4:00进行维护升级", "站内信"},
		{"新路线上线", "公告", "紧急", "本公司全员", "厦门→基隆 海运路线已开通，首重15元/kg", "邮件"},
		{"账户余额不足", "任务通知", "普通", "指定用户", "您的账户余额已低于100元，请及时充值", "短信"},
		{"包裹已签收", "任务通知", "普通", "指定用户", "包裹YT7625763166053已被签收", "微信"},
		{"假期运营安排", "公告", "普通", "全员(跨所有公司)", "端午假期正常运营，部分路线时效延长", "站内信"},
	}
	for i, n := range notifications {
		ncs.Upsert(ctx, fmt.Sprintf("NTF%03d", i+1), n.title, n.typ, n.priority, n.scope, n.content, n.channel, "系统管理员", true, now.AddDate(0, 0, -i))
	}

	// ── Orders + Parcels ──
	statuses := []string{"pending_picking", "picking", "pending_packing", "pending_loading", "loaded", "in_transit", "customs_clearance", "delivered", "completed"}
	productNames := []string{"手机壳", "数据线", "陶瓷杯", "T恤", "洗发水", "蓝牙耳机", "运动鞋", "充电宝", "化妆品套装", "笔记本电脑"}
	cargoTypes := []string{"普货", "易碎品", "液体", "电子产品", "食品"}

	for i := 0; i < 12; i++ {
		// Parcel
		trackingNo := fmt.Sprintf("%s%s", []string{"TN", "SF", "ZTO", "YTO", "STO", "JD", "EMS", "HTKY"}[rng.Intn(8)], fmt.Sprintf("%010d", rng.Int63n(9999999999)))
		prodIdx := rng.Intn(len(productNames))
		p, _ := psvc.PreDeclare(ctx, &parcelDomain.Parcel{
			TrackingNumber: trackingNo, ProductName: productNames[prodIdx],
			ActualWeight: float64(rng.Intn(300)+5) / 100.0,
			CourierCode: couriers[rng.Intn(len(couriers))].code,
			CargoType: cargoTypes[rng.Intn(len(cargoTypes))],
			TenantID: 1, WarehouseID: 1,
		})
		// Transition parcel status
		psvc.Transition(ctx, 1, p.ID, parcelDomain.ParcelStatus(statuses[rng.Intn(5)]))

		// Order (some parcels linked to orders)
		if i < 9 {
			ord, _ := osvc.Create(ctx, &orderDomain.Order{
				RecipientName: memberNames[rng.Intn(len(memberNames))],
				ParcelCount: rng.Intn(4) + 1,
				TotalPrice: float64(rng.Intn(5000)+500) / 100.0,
				RouteID: int64(rng.Intn(len(routes)) + 1),
				TenantID: 1,
				TrackingNumbers: trackingNo,
				Remark: []string{"空运急件", "大件运输", "已签收", "PDA演示订单-待拣货", "", ""}[rng.Intn(6)],
			})
			osvc.Transition(ctx, 1, ord.ID, orderDomain.OrderStatus(statuses[rng.Intn(len(statuses))]))
		}
	}

	fmt.Printf("[seed] Created %d couriers, %d routes, %d clients, %d members, %d declarants, %d addresses\n",
		len(couriers), len(routes), len(clients), len(memberNames), len(declarants), len(addresses))
}

func randName(rng *rand.Rand) string {
	first := []string{"张", "李", "王", "陈", "林", "黄", "赵", "周", "吴", "郑"}
	last := []string{"明", "华", "强", "丽", "芳", "伟", "静", "涛", "娜", "平"}
	return first[rng.Intn(len(first))] + last[rng.Intn(len(last))]
}

func randPhone(rng *rand.Rand) string {
	return fmt.Sprintf("138%08d", rng.Intn(99999999))
}
