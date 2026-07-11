package repository
import (
	"context"; "sync"; "sync/atomic"; "time"
	"github.com/i56/modules/parcel_service/domain"
)
type MemServiceRepo struct {
	mu      sync.RWMutex
	orders  map[int64]*domain.ServiceOrder
	types   []domain.ServiceType
	nextID  int64
}
func NewMemServiceRepo() *MemServiceRepo {
	return &MemServiceRepo{orders:make(map[int64]*domain.ServiceOrder),types:domain.DefaultServiceTypes(),nextID:1}
}
func (r *MemServiceRepo) ListTypes() []domain.ServiceType { return r.types }
func (r *MemServiceRepo) GetTypeByCode(code string) *domain.ServiceType {
	for _,t := range r.types { if t.Code==code { return &t } }; return nil
}
func (r *MemServiceRepo) Create(ctx context.Context, o *domain.ServiceOrder) error {
	r.mu.Lock(); defer r.mu.Unlock()
	o.ID = atomic.AddInt64(&r.nextID,1)-1
	o.CreatedAt=time.Now(); o.UpdatedAt=time.Now()
	r.orders[o.ID]=o; return nil
}
func (r *MemServiceRepo) List(ctx context.Context, tenantID int64, offset,limit int) ([]domain.ServiceOrder,int64,error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	var result []domain.ServiceOrder
	for _,o := range r.orders { if o.TenantID==tenantID { result=append(result,*o) } }
	total:=int64(len(result))
	if offset>=int(total){return nil,total,nil}
	end:=offset+limit; if end>int(total){end=int(total)}
	return result[offset:end],total,nil
}

// AddType appends a new service type.
func (r *MemServiceRepo) AddType(name, code, category string, unitPrice float64, priceMode string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := int64(len(r.types) + 1)
	r.types = append(r.types, domain.ServiceType{ID: id, Name: name, Code: code, Category: category, UnitPrice: unitPrice, PriceMode: priceMode})
}

// DeleteType removes a service type by code.
// GetByID returns a service order by ID.
func (r *MemServiceRepo) GetByID(ctx context.Context, id int64) (*domain.ServiceOrder, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	if o, ok := r.orders[id]; ok { return o, nil }; return nil, nil
}

// Update modifies an existing service order.
func (r *MemServiceRepo) Update(ctx context.Context, o *domain.ServiceOrder) error {
	r.mu.Lock(); defer r.mu.Unlock()
	if _, ok := r.orders[o.ID]; ok {
		o.UpdatedAt = time.Now()
		r.orders[o.ID] = o
	}
	return nil
}

func (r *MemServiceRepo) DeleteType(code string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, t := range r.types {
		if t.Code == code {
			r.types = append(r.types[:i], r.types[i+1:]...)
			return
		}
	}
}
