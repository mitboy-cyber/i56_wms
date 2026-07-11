package repository

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/modules/pricing/domain"
)

// MemPricingModelsRepo holds all 5-tab pricing models in memory with real BFT56 seed data.
type MemPricingModelsRepo struct {
	mu             sync.RWMutex
	routePrices    map[int64]*domain.RoutePriceModel
	storagePrices  map[int64]*domain.StoragePriceModel
	deliveryFees   map[int64]*domain.DeliveryFeeModel
	surcharges     map[int64]*domain.SurchargeModel
	servicePrices  map[int64]*domain.ServicePriceModel
	nextID         int64
}

func NewMemPricingModelsRepo() *MemPricingModelsRepo {
	r := &MemPricingModelsRepo{
		routePrices:   make(map[int64]*domain.RoutePriceModel),
		storagePrices: make(map[int64]*domain.StoragePriceModel),
		deliveryFees:  make(map[int64]*domain.DeliveryFeeModel),
		surcharges:    make(map[int64]*domain.SurchargeModel),
		servicePrices: make(map[int64]*domain.ServicePriceModel),
	}
	r.seed()
	return r
}

func (r *MemPricingModelsRepo) next() int64 { return atomic.AddInt64(&r.nextID, 1) - 1 }

// ─── Seed data from real BFT56 analysis ──────────────────────────────

func (r *MemPricingModelsRepo) seed() {
	now := time.Now()
	tenantID := int64(1)
	clientID := int64(1)
	clientName := "EZ集运通"

	// ─── Tab 1: Route Prices ─────────────────────────────────────
	// Real data from BFT56: 空運普货=¥20/kg, 海快普货=¥8.30/kg,
	// 海运家具类=¥2.50/kg+¥15/才, through 六类=¥6.80/kg+¥38/才
	routeSeeds := []struct {
		route, trans, cargo, tax      string
		wp, vp, mc, fw, fwp, cwp     float64
		fv, fvp, cvp                  float64
		volCoeff                      int
	}{
		{"深圳→台湾(空运)", "air", "普货", "全包税", 20.00, 20.00, 50, 0.5, 25.00, 20.00, 0.5, 25.00, 20.00, 6000},
		{"深圳→台湾(空运)", "air", "特货", "全包税", 28.00, 28.00, 80, 0.5, 35.00, 28.00, 0.5, 35.00, 28.00, 6000},
		{"厦门→台湾(海快)", "sea_express", "普货", "频税", 8.30, 15.00, 50, 1.0, 15.00, 8.30, 1.0, 15.00, 15.00, 6000},
		{"厦门→台湾(海快)", "sea_express", "一类", "频税", 10.50, 18.00, 60, 1.0, 18.00, 10.50, 1.0, 18.00, 18.00, 6000},
		{"厦门→台湾(海快)", "sea_express", "二类", "频税", 12.00, 20.00, 80, 1.0, 20.00, 12.00, 1.0, 20.00, 20.00, 6000},
		{"厦门→台湾(海运)", "sea", "家具类", "全包税", 2.50, 15.00, 50, 10.0, 32.00, 2.50, 1.0, 20.00, 15.00, 6000},
		{"厦门→台湾(海运)", "sea", "一类", "全包税", 3.20, 20.00, 50, 10.0, 32.00, 3.20, 1.0, 20.00, 20.00, 6000},
		{"厦门→台湾(海运)", "sea", "二类", "全包税", 4.50, 25.00, 60, 10.0, 45.00, 4.50, 1.0, 25.00, 25.00, 6000},
		{"厦门→台湾(海运)", "sea", "三类", "全包税", 5.00, 28.00, 80, 10.0, 50.00, 5.00, 1.0, 28.00, 28.00, 6000},
		{"厦门→台湾(海运)", "sea", "四类", "全包税", 5.80, 32.00, 80, 10.0, 58.00, 5.80, 1.0, 32.00, 32.00, 6000},
		{"厦门→台湾(海运)", "sea", "五类", "全包税", 6.30, 35.00, 100, 10.0, 63.00, 6.30, 1.0, 35.00, 35.00, 6000},
		{"厦门→台湾(海运)", "sea", "六类", "全包税", 6.80, 38.00, 100, 10.0, 68.00, 6.80, 1.0, 38.00, 38.00, 6000},
		{"厦门→台湾(海运)", "sea", "易碎品", "全包税", 4.00, 22.00, 60, 10.0, 40.00, 4.00, 1.0, 22.00, 22.00, 6000},
	}

	for _, s := range routeSeeds {
		id := r.next()
		r.routePrices[id] = &domain.RoutePriceModel{
			ID: id, TenantID: tenantID, ClientID: clientID, ClientName: clientName,
			RouteName: s.route, TransportType: s.trans, CargoType: s.cargo, TaxMode: s.tax,
			WeightPrice: s.wp, VolumePrice: s.vp, MinCharge: s.mc,
			FirstWeight: s.fw, FirstWeightPrice: s.fwp, ContWeightPrice: s.cwp,
			FirstVolume: s.fv, FirstVolumePrice: s.fvp, ContVolumePrice: s.cvp,
			VolumeCoeff: s.volCoeff, IsActive: true, CreatedAt: now, UpdatedAt: now,
		}
	}

	// ─── Tab 2: Storage Prices ───────────────────────────────────
	storageSeeds := []struct {
		whID   int64
		whName string
		free   int
		daily  float64
		maxD   int
	}{
		{1, "厦门仓", 30, 0.50, 90},
		{2, "深圳仓", 15, 0.80, 60},
	}
	for _, s := range storageSeeds {
		id := r.next()
		r.storagePrices[id] = &domain.StoragePriceModel{
			ID: id, TenantID: tenantID, ClientID: clientID, ClientName: clientName,
			WarehouseID: s.whID, WarehouseName: s.whName,
			FreeDays: s.free, DailyRate: s.daily, MaxStorageDays: s.maxD,
			IsActive: true, CreatedAt: now, UpdatedAt: now,
		}
	}

	// ─── Tab 3: Delivery Fees ────────────────────────────────────
	// Real BFT56 data: 新竹物流(台北/宅配)=¥20 (free ≥10kg), 新竹物流(专车)=¥3500 fixed
	deliverySeeds := []struct {
		carrier string
		cp, area, method, cond string
		fee     float64
		freeKg  float64
		freeTxt string
	}{
		{"新竹物流", "台北", "預設", "宅配", "重量>39.8kg", 20, 10, "≥10kg免运"},
		{"新竹物流", "台北", "預設", "專車", "單邊長>=600cm或重量>=500kg", 3500, 0, ""},
		{"新竹物流", "台中", "預設", "宅配", "重量>39.8kg", 20, 10, "≥10kg免运"},
		{"新竹物流", "高雄", "預設", "宅配", "重量>39.8kg", 25, 10, "≥10kg免运"},
		{"新竹物流", "台北", "东部", "宅配", "宜蘭/花蓮/台東", 50, 0, ""},
		{"黑猫宅急便", "台北", "預設", "宅配", "—", 15, 0, ""},
		{"黑猫宅急便", "高雄", "預設", "宅配", "—", 18, 0, ""},
		{"顺丰速运", "台北", "預設", "宅配", "—", 12, 0, ""},
		{"顺丰速运", "台中", "預設", "宅配", "—", 12, 0, ""},
	}
	for _, s := range deliverySeeds {
		id := r.next()
		r.deliveryFees[id] = &domain.DeliveryFeeModel{
			ID: id, TenantID: tenantID, ClientID: clientID, ClientName: clientName,
			CarrierName: s.carrier, CustomsPoint: s.cp, Area: s.area,
			DeliveryMethod: s.method, Condition: s.cond,
			Fee: s.fee, FreeThreshold: s.freeKg, FreeThresholdTxt: s.freeTxt,
			IsActive: true, CreatedAt: now, UpdatedAt: now,
		}
	}

	// ─── Tab 4: Surcharges ───────────────────────────────────────
	// Real BFT56: 超長費, 超材費, 小板, 大板 × 清關點 × 區域
	surchargeSeeds := []struct {
		carrier, chgType, tier, cp, area, cond, pd string
		price float64
	}{
		{"新竹物流", "超長費", "—", "台北", "預設", "單邊>150cm", "每件加收", 100},
		{"新竹物流", "超長費", "—", "台中", "預設", "單邊>150cm", "每件加收", 100},
		{"新竹物流", "超長費", "—", "高雄", "預設", "單邊>150cm", "每件加收", 120},
		{"新竹物流", "超材費", "—", "台北", "預設", "體積重>實重2倍", "按體積重計費差額", 0},
		{"新竹物流", "超材費", "—", "台中", "預設", "體積重>實重2倍", "按體積重計費差額", 0},
		{"新竹物流", "超材費", "—", "高雄", "預設", "體積重>實重2倍", "按體積重計費差額", 0},
		{"新竹物流", "棧板費", "小板", "台北", "預設", "棧板尺寸≤120cm", "每板加收", 150},
		{"新竹物流", "棧板費", "大板", "台北", "預設", "棧板尺寸>120cm", "每板加收", 300},
		{"新竹物流", "棧板費", "小板", "台中", "預設", "棧板尺寸≤120cm", "每板加收", 150},
		{"新竹物流", "棧板費", "大板", "台中", "預設", "棧板尺寸>120cm", "每板加收", 300},
		{"新竹物流", "偏遠費", "—", "台北", "預設", "偏遠/離島地區", "每票加收", 50},
		{"新竹物流", "偏遠費", "—", "高雄", "預設", "偏遠/離島地區", "每票加收", 60},
		{"新竹物流", "上樓費", "—", "台北", "預設", "無電梯4樓以上", "每層加收", 20},
		{"新竹物流", "上樓費", "—", "台中", "預設", "無電梯4樓以上", "每層加收", 20},
	}
	for _, s := range surchargeSeeds {
		id := r.next()
		r.surcharges[id] = &domain.SurchargeModel{
			ID: id, TenantID: tenantID, ClientID: clientID, ClientName: clientName,
			CarrierName: s.carrier, ChargeType: s.chgType, Tier: s.tier,
			CustomsPoint: s.cp, Area: s.area, Condition: s.cond,
			Price: s.price, PriceDesc: s.pd,
			IsActive: true, CreatedAt: now, UpdatedAt: now,
		}
	}

	// ─── Tab 5: Service Prices ───────────────────────────────────
	serviceSeeds := []struct {
		stype, scode, pm string
		price            float64
	}{
		{"木箱包装", "WOODEN_CRATE", "per_item", 80.00},
		{"开箱验货", "OPEN_INSPECT", "fixed", 0.00},
		{"拍照存证", "PHOTO", "per_item", 5.00},
		{"换标服务", "RELABEL", "per_item", 3.00},
		{"合箱服务", "MERGE", "per_order", 10.00},
		{"拆箱服务", "UNBOX", "per_item", 5.00},
		{"加固包装", "REINFORCE", "per_item", 15.00},
		{"代付运费", "PREPAY_FREIGHT", "per_order", 0.00},
		{"退货签收", "RETURN_SIGN", "per_item", 8.00},
		{"暂存服务", "TEMP_STORAGE", "per_kg", 0.50},
	}
	for _, s := range serviceSeeds {
		id := r.next()
		r.servicePrices[id] = &domain.ServicePriceModel{
			ID: id, TenantID: tenantID, ClientID: clientID, ClientName: clientName,
			ServiceType: s.stype, ServiceCode: s.scode,
			UnitPrice: s.price, PriceMode: s.pm,
			IsActive: true, CreatedAt: now, UpdatedAt: now,
		}
	}
}

// ─── List methods ───────────────────────────────────────────────────

func (r *MemPricingModelsRepo) ListRoutePrices() []domain.RoutePriceModel {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]domain.RoutePriceModel, 0, len(r.routePrices))
	for _, p := range r.routePrices {
		result = append(result, *p)
	}
	return result
}

func (r *MemPricingModelsRepo) ListStoragePrices() []domain.StoragePriceModel {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]domain.StoragePriceModel, 0, len(r.storagePrices))
	for _, p := range r.storagePrices {
		result = append(result, *p)
	}
	return result
}

func (r *MemPricingModelsRepo) ListDeliveryFees() []domain.DeliveryFeeModel {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]domain.DeliveryFeeModel, 0, len(r.deliveryFees))
	for _, f := range r.deliveryFees {
		result = append(result, *f)
	}
	return result
}

func (r *MemPricingModelsRepo) ListSurcharges() []domain.SurchargeModel {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]domain.SurchargeModel, 0, len(r.surcharges))
	for _, s := range r.surcharges {
		result = append(result, *s)
	}
	return result
}

func (r *MemPricingModelsRepo) ListServicePrices() []domain.ServicePriceModel {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]domain.ServicePriceModel, 0, len(r.servicePrices))
	for _, p := range r.servicePrices {
		result = append(result, *p)
	}
	return result
}

// ─── Add methods ────────────────────────────────────────────────────

func (r *MemPricingModelsRepo) AddRoutePrice(p *domain.RoutePriceModel) {
	r.mu.Lock()
	defer r.mu.Unlock()
	p.ID = r.next()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	if p.TenantID == 0 { p.TenantID = 1 }
	if p.ClientID == 0 { p.ClientID = 1 }
	if p.ClientName == "" { p.ClientName = "EZ集运通" }
	p.IsActive = true
	r.routePrices[p.ID] = p
}

func (r *MemPricingModelsRepo) AddStoragePrice(p *domain.StoragePriceModel) {
	r.mu.Lock()
	defer r.mu.Unlock()
	p.ID = r.next()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	if p.TenantID == 0 { p.TenantID = 1 }
	if p.ClientID == 0 { p.ClientID = 1 }
	if p.ClientName == "" { p.ClientName = "EZ集运通" }
	p.IsActive = true
	r.storagePrices[p.ID] = p
}

func (r *MemPricingModelsRepo) AddDeliveryFee(f *domain.DeliveryFeeModel) {
	r.mu.Lock()
	defer r.mu.Unlock()
	f.ID = r.next()
	f.CreatedAt = time.Now()
	f.UpdatedAt = time.Now()
	if f.TenantID == 0 { f.TenantID = 1 }
	if f.ClientID == 0 { f.ClientID = 1 }
	if f.ClientName == "" { f.ClientName = "EZ集运通" }
	f.IsActive = true
	r.deliveryFees[f.ID] = f
}

func (r *MemPricingModelsRepo) AddSurcharge(s *domain.SurchargeModel) {
	r.mu.Lock()
	defer r.mu.Unlock()
	s.ID = r.next()
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	if s.TenantID == 0 { s.TenantID = 1 }
	if s.ClientID == 0 { s.ClientID = 1 }
	if s.ClientName == "" { s.ClientName = "EZ集运通" }
	s.IsActive = true
	r.surcharges[s.ID] = s
}

// ─── GetByID methods ────────────────────────────────────────────────

func (r *MemPricingModelsRepo) GetRoutePriceByID(id int64) *domain.RoutePriceModel {
	r.mu.RLock(); defer r.mu.RUnlock()
	if p, ok := r.routePrices[id]; ok { return p }; return nil
}
func (r *MemPricingModelsRepo) GetDeliveryFeeByID(id int64) *domain.DeliveryFeeModel {
	r.mu.RLock(); defer r.mu.RUnlock()
	if f, ok := r.deliveryFees[id]; ok { return f }; return nil
}
func (r *MemPricingModelsRepo) GetSurchargeByID(id int64) *domain.SurchargeModel {
	r.mu.RLock(); defer r.mu.RUnlock()
	if s, ok := r.surcharges[id]; ok { return s }; return nil
}

// ─── Update methods ─────────────────────────────────────────────────

func (r *MemPricingModelsRepo) UpdateRoutePrice(id int64, p *domain.RoutePriceModel) {
	r.mu.Lock(); defer r.mu.Unlock()
	if existing, ok := r.routePrices[id]; ok {
		p.ID = id; p.CreatedAt = existing.CreatedAt; p.UpdatedAt = time.Now()
		r.routePrices[id] = p
	}
}
func (r *MemPricingModelsRepo) UpdateDeliveryFee(id int64, f *domain.DeliveryFeeModel) {
	r.mu.Lock(); defer r.mu.Unlock()
	if existing, ok := r.deliveryFees[id]; ok {
		f.ID = id; f.CreatedAt = existing.CreatedAt; f.UpdatedAt = time.Now()
		r.deliveryFees[id] = f
	}
}
func (r *MemPricingModelsRepo) UpdateSurcharge(id int64, s *domain.SurchargeModel) {
	r.mu.Lock(); defer r.mu.Unlock()
	if existing, ok := r.surcharges[id]; ok {
		s.ID = id; s.CreatedAt = existing.CreatedAt; s.UpdatedAt = time.Now()
		r.surcharges[id] = s
	}
}

func (r *MemPricingModelsRepo) AddServicePrice(p *domain.ServicePriceModel) {
	r.mu.Lock()
	defer r.mu.Unlock()
	p.ID = r.next()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	if p.TenantID == 0 { p.TenantID = 1 }
	if p.ClientID == 0 { p.ClientID = 1 }
	if p.ClientName == "" { p.ClientName = "EZ集运通" }
	p.IsActive = true
	r.servicePrices[p.ID] = p
}
