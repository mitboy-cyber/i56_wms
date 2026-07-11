package repository
import (
	"context"; "sync"; "sync/atomic"; "time"
	"github.com/i56/modules/pricing/domain"
)
type MemPricingRepo struct {
	mu         sync.RWMutex
	routePrices map[int64]*domain.ClientRoutePrice
	deliveryFees map[int64]*domain.ClientDeliveryFee
	surcharges  map[int64]*domain.ClientSurcharge
	storagePrices map[int64]*domain.ClientStoragePrice
	serviceOverrides map[string]*domain.ClientServiceOverride
	statements  []domain.MonthlyStatement
	nextID      int64
}
func NewMemPricingRepo() *MemPricingRepo {
	return &MemPricingRepo{
		routePrices: make(map[int64]*domain.ClientRoutePrice),
		deliveryFees: make(map[int64]*domain.ClientDeliveryFee),
		surcharges: make(map[int64]*domain.ClientSurcharge),
		storagePrices: make(map[int64]*domain.ClientStoragePrice),
		serviceOverrides: make(map[string]*domain.ClientServiceOverride),
	}
}
func (r *MemPricingRepo) next() int64 { return atomic.AddInt64(&r.nextID, 1)-1 }

// Route prices
func (r *MemPricingRepo) CreateRoutePrice(ctx context.Context, p *domain.ClientRoutePrice) error {
	r.mu.Lock(); defer r.mu.Unlock()
	p.ID=r.next(); p.CreatedAt=time.Now(); p.UpdatedAt=time.Now()
	r.routePrices[p.ID]=p; return nil
}
func (r *MemPricingRepo) ListRoutePrices(ctx context.Context, tenantID,clientID int64) []domain.ClientRoutePrice {
	r.mu.RLock(); defer r.mu.RUnlock()
	var result []domain.ClientRoutePrice
	for _,p := range r.routePrices {
		if p.TenantID==tenantID&&p.ClientID==clientID { result=append(result,*p) }
	}
	return result
}
// Delivery fees
func (r *MemPricingRepo) CreateDeliveryFee(ctx context.Context, f *domain.ClientDeliveryFee) error {
	r.mu.Lock(); defer r.mu.Unlock(); f.ID=r.next(); r.deliveryFees[f.ID]=f; return nil
}
func (r *MemPricingRepo) ListDeliveryFees(ctx context.Context, tenantID,clientID int64) []domain.ClientDeliveryFee {
	r.mu.RLock(); defer r.mu.RUnlock()
	var result []domain.ClientDeliveryFee
	for _,f := range r.deliveryFees { if f.TenantID==tenantID&&f.ClientID==clientID { result=append(result,*f) } }
	return result
}
// Surcharges
func (r *MemPricingRepo) CreateSurcharge(ctx context.Context, s *domain.ClientSurcharge) error {
	r.mu.Lock(); defer r.mu.Unlock(); s.ID=r.next(); r.surcharges[s.ID]=s; return nil
}
func (r *MemPricingRepo) ListSurcharges(ctx context.Context, tenantID,clientID int64) []domain.ClientSurcharge {
	r.mu.RLock(); defer r.mu.RUnlock()
	var result []domain.ClientSurcharge
	for _,s := range r.surcharges { if s.TenantID==tenantID&&s.ClientID==clientID { result=append(result,*s) } }
	return result
}
// Storage prices
func (r *MemPricingRepo) CreateStoragePrice(ctx context.Context, s *domain.ClientStoragePrice) error {
	r.mu.Lock(); defer r.mu.Unlock(); s.ID=r.next(); r.storagePrices[s.ID]=s; return nil
}
func (r *MemPricingRepo) ListStoragePrices(ctx context.Context, tenantID,clientID int64) []domain.ClientStoragePrice {
	r.mu.RLock(); defer r.mu.RUnlock()
	var result []domain.ClientStoragePrice
	for _,s := range r.storagePrices { if s.TenantID==tenantID&&s.ClientID==clientID { result=append(result,*s) } }
	return result
}
// Service overrides
func (r *MemPricingRepo) CreateServiceOverride(ctx context.Context, o *domain.ClientServiceOverride) error {
	r.mu.Lock(); defer r.mu.Unlock(); o.ID=r.next()
	r.serviceOverrides[o.ServiceCode]=o; return nil
}
func (r *MemPricingRepo) ListServiceOverrides(ctx context.Context, tenantID,clientID int64) []domain.ClientServiceOverride {
	r.mu.RLock(); defer r.mu.RUnlock()
	var result []domain.ClientServiceOverride
	for _,o := range r.serviceOverrides { if o.TenantID==tenantID&&o.ClientID==clientID { result=append(result,*o) } }
	return result
}
// Statements
func (r *MemPricingRepo) CreateStatement(ctx context.Context, s *domain.MonthlyStatement) error {
	r.mu.Lock(); defer r.mu.Unlock(); r.statements=append(r.statements,*s); return nil
}
func (r *MemPricingRepo) ListStatements(ctx context.Context, tenantID,clientID int64) []domain.MonthlyStatement {
	r.mu.RLock(); defer r.mu.RUnlock()
	var result []domain.MonthlyStatement
	for _,s := range r.statements { if s.TenantID==tenantID { result=append(result,s) } }
	return result
}
