package repository
import (
	"context"; "sync"; "sync/atomic"; "time"
	"github.com/i56/modules/wms/domain"
)
type MemWMSRepo struct {
	mu             sync.RWMutex
	zones          map[int64]*domain.Zone
	zoneTypes      []domain.ZoneType
	locations      map[int64]*domain.Location
	locationTypes  []domain.LocationType
	containers     map[int64]*domain.Container
	inboundMachines map[int64]*domain.InboundMachine
	nextID         int64
}
func NewMemWMSRepo() *MemWMSRepo {
	return &MemWMSRepo{
		zones:make(map[int64]*domain.Zone),zoneTypes:domain.DefaultZoneTypes(),
		locations:make(map[int64]*domain.Location),locationTypes:domain.DefaultLocationTypes(),
		containers:make(map[int64]*domain.Container),inboundMachines:make(map[int64]*domain.InboundMachine),
	}
}
func (r *MemWMSRepo) next() int64 { return atomic.AddInt64(&r.nextID,1)-1 }

func (r *MemWMSRepo) ListZoneTypes() []domain.ZoneType { return r.zoneTypes }
func (r *MemWMSRepo) ListLocationTypes() []domain.LocationType { return r.locationTypes }

func (r *MemWMSRepo) CreateZone(ctx context.Context, z *domain.Zone) error {
	r.mu.Lock(); defer r.mu.Unlock(); z.ID=r.next(); z.CreatedAt=time.Now(); r.zones[z.ID]=z; return nil
}
func (r *MemWMSRepo) ListZones(ctx context.Context, warehouseID int64) []domain.Zone {
	r.mu.RLock(); defer r.mu.RUnlock()
	var result []domain.Zone
	for _,z := range r.zones { if z.WarehouseID==warehouseID { result=append(result,*z) } }
	return result
}
func (r *MemWMSRepo) CreateLocation(ctx context.Context, l *domain.Location) error {
	r.mu.Lock(); defer r.mu.Unlock(); l.ID=r.next(); l.CreatedAt=time.Now(); r.locations[l.ID]=l; return nil
}
func (r *MemWMSRepo) ListLocations(ctx context.Context, zoneID int64) []domain.Location {
	r.mu.RLock(); defer r.mu.RUnlock()
	var result []domain.Location
	for _,l := range r.locations { if l.ZoneID==zoneID { result=append(result,*l) } }
	return result
}
func (r *MemWMSRepo) CreateContainer(ctx context.Context, c *domain.Container) error {
	r.mu.Lock(); defer r.mu.Unlock(); c.ID=r.next(); c.CreatedAt=time.Now(); r.containers[c.ID]=c; return nil
}
func (r *MemWMSRepo) ListContainers(ctx context.Context, warehouseID int64) []domain.Container {
	r.mu.RLock(); defer r.mu.RUnlock()
	var result []domain.Container
	for _,c := range r.containers { if c.WarehouseID==warehouseID { result=append(result,*c) } }
	return result
}
func (r *MemWMSRepo) CreateInboundMachine(ctx context.Context, m *domain.InboundMachine) error {
	r.mu.Lock(); defer r.mu.Unlock(); m.ID=r.next(); r.inboundMachines[m.ID]=m; return nil
}
func (r *MemWMSRepo) ListInboundMachines(ctx context.Context, warehouseID int64) []domain.InboundMachine {
	r.mu.RLock(); defer r.mu.RUnlock()
	var result []domain.InboundMachine
	for _,m := range r.inboundMachines { if m.WarehouseID==warehouseID { result=append(result,*m) } }
	return result
}
func (r *MemWMSRepo) GetLocationByBarcode(ctx context.Context, barcode string) *domain.Location {
	r.mu.RLock(); defer r.mu.RUnlock()
	for _,l := range r.locations { if l.Barcode==barcode { return l } }
	return nil
}
