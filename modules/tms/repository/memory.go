package repository
import (
	"context"; "sync"; "sync/atomic"; "time"
	"github.com/i56/modules/tms/domain"
)
type MemTMSRepo struct {
	mu              sync.RWMutex
	areaGroups      map[int64]*domain.AreaGroup
	carrierNumbers  map[int64]*domain.CarrierNumber
	customsBrokers  map[int64]*domain.CustomsBroker
	customsPoints   map[int64]*domain.CustomsPoint
	customsNumbers  map[int64]*domain.CustomsNumber
	loadings        []domain.ContainerLoading
	shippingProviders map[int64]*domain.ShippingProvider
	transportTypes  []domain.TransportType
	trackings       []domain.Tracking
	nextID          int64
}
func NewMemTMSRepo() *MemTMSRepo {
	return &MemTMSRepo{
		areaGroups:make(map[int64]*domain.AreaGroup),carrierNumbers:make(map[int64]*domain.CarrierNumber),
		customsBrokers:make(map[int64]*domain.CustomsBroker),customsPoints:make(map[int64]*domain.CustomsPoint),
		customsNumbers:make(map[int64]*domain.CustomsNumber),
		shippingProviders:make(map[int64]*domain.ShippingProvider),transportTypes:domain.DefaultTransportTypes(),
	}
}
func (r *MemTMSRepo) next() int64 { return atomic.AddInt64(&r.nextID,1)-1 }

// AreaGroups
func (r *MemTMSRepo) CreateAreaGroup(ctx context.Context, a *domain.AreaGroup) error {
	r.mu.Lock(); defer r.mu.Unlock(); a.ID=r.next(); r.areaGroups[a.ID]=a; return nil
}
func (r *MemTMSRepo) ListAreaGroups(ctx context.Context, tenantID int64) []domain.AreaGroup {
	r.mu.RLock(); defer r.mu.RUnlock()
	var result []domain.AreaGroup
	for _,a := range r.areaGroups { if a.TenantID==tenantID { result=append(result,*a) } }
	return result
}
// CarrierNumbers
func (r *MemTMSRepo) CreateCarrierNumber(ctx context.Context, c *domain.CarrierNumber) error {
	r.mu.Lock(); defer r.mu.Unlock(); c.ID=r.next(); c.CurrentNo=c.StartNo; r.carrierNumbers[c.ID]=c; return nil
}
func (r *MemTMSRepo) AllocateCarrierNumber(ctx context.Context, carrierID int64) (string, error) {
	r.mu.Lock(); defer r.mu.Unlock()
	for _,c := range r.carrierNumbers {
		if c.CarrierID==carrierID&&c.IsActive&&c.CurrentNo<=c.EndNo {
			no:=c.CurrentNo; c.CurrentNo++; return c.Prefix+string(rune(no)),nil
		}
	}
	return "",nil
}
// CustomsBrokers
func (r *MemTMSRepo) CreateCustomsBroker(ctx context.Context, c *domain.CustomsBroker) error {
	r.mu.Lock(); defer r.mu.Unlock(); c.ID=r.next(); r.customsBrokers[c.ID]=c; return nil
}
func (r *MemTMSRepo) ListCustomsBrokers(ctx context.Context, tenantID int64) []domain.CustomsBroker {
	r.mu.RLock(); defer r.mu.RUnlock()
	var result []domain.CustomsBroker
	for _,c := range r.customsBrokers { if c.TenantID==tenantID { result=append(result,*c) } }
	return result
}
// CustomsPoints
func (r *MemTMSRepo) CreateCustomsPoint(ctx context.Context, c *domain.CustomsPoint) error {
	r.mu.Lock(); defer r.mu.Unlock(); c.ID=r.next(); r.customsPoints[c.ID]=c; return nil
}
func (r *MemTMSRepo) ListCustomsPoints(ctx context.Context, tenantID int64) []domain.CustomsPoint {
	r.mu.RLock(); defer r.mu.RUnlock()
	var result []domain.CustomsPoint
	for _,c := range r.customsPoints { if c.TenantID==tenantID { result=append(result,*c) } }
	return result
}
// CustomsNumbers
func (r *MemTMSRepo) CreateCustomsNumber(ctx context.Context, c *domain.CustomsNumber) error {
	r.mu.Lock(); defer r.mu.Unlock(); c.ID=r.next(); c.CurrentNo=c.StartNo; r.customsNumbers[c.ID]=c; return nil
}
func (r *MemTMSRepo) AllocateCustomsNumber(ctx context.Context, pointID int64) (string, error) {
	r.mu.Lock(); defer r.mu.Unlock()
	for _,c := range r.customsNumbers {
		if c.CustomsPointID==pointID&&c.IsActive&&c.CurrentNo<=c.EndNo {
			no:=c.CurrentNo; c.CurrentNo++; return c.Prefix+string(rune(no)),nil
		}
	}
	return "",nil
}
// ContainerLoadings
func (r *MemTMSRepo) RecordLoading(ctx context.Context, l *domain.ContainerLoading) error {
	r.mu.Lock(); defer r.mu.Unlock(); r.loadings=append(r.loadings,*l); return nil
}
func (r *MemTMSRepo) ListLoadings(ctx context.Context, containerID int64) []domain.ContainerLoading {
	r.mu.RLock(); defer r.mu.RUnlock()
	var result []domain.ContainerLoading
	for _,l := range r.loadings { if l.ContainerID==containerID { result=append(result,l) } }
	return result
}
// Shipping Providers
func (r *MemTMSRepo) CreateShippingProvider(ctx context.Context, s *domain.ShippingProvider) error {
	r.mu.Lock(); defer r.mu.Unlock(); s.ID=r.next(); r.shippingProviders[s.ID]=s; return nil
}
func (r *MemTMSRepo) ListShippingProviders(ctx context.Context) []domain.ShippingProvider {
	r.mu.RLock(); defer r.mu.RUnlock()
	var result []domain.ShippingProvider
	for _,s := range r.shippingProviders { result=append(result,*s) }
	return result
}
func (r *MemTMSRepo) ListTransportTypes() []domain.TransportType { return r.transportTypes }
// Tracking
func (r *MemTMSRepo) AddTracking(ctx context.Context, t *domain.Tracking) error {
	r.mu.Lock(); defer r.mu.Unlock(); t.CreatedAt=time.Now(); r.trackings=append(r.trackings,*t); return nil
}
func (r *MemTMSRepo) GetTrackingByOrderID(ctx context.Context, orderID int64) []domain.Tracking {
	r.mu.RLock(); defer r.mu.RUnlock()
	var result []domain.Tracking
	for _,t := range r.trackings { if t.OrderID==orderID { result=append(result,t) } }
	return result
}
