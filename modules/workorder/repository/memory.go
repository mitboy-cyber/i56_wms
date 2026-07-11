package repository
import (
	"context"; "sync"; "sync/atomic"; "time"
	"github.com/i56/modules/workorder/domain"
)
type MemWorkOrderRepo struct {
	mu     sync.RWMutex
	orders map[int64]*domain.WorkOrder
	nextID int64
}
func NewMemWorkOrderRepo() *MemWorkOrderRepo {
	return &MemWorkOrderRepo{orders:make(map[int64]*domain.WorkOrder),nextID:1}
}
func (r *MemWorkOrderRepo) Create(ctx context.Context, wo *domain.WorkOrder) error {
	r.mu.Lock(); defer r.mu.Unlock()
	wo.ID=atomic.AddInt64(&r.nextID,1)-1
	wo.CreatedAt=time.Now(); wo.UpdatedAt=time.Now()
	r.orders[wo.ID]=wo; return nil
}
func (r *MemWorkOrderRepo) List(ctx context.Context, tenantID int64, offset,limit int) ([]domain.WorkOrder,int64,error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	var result []domain.WorkOrder
	for _,o := range r.orders { if o.TenantID==tenantID { result=append(result,*o) } }
	total:=int64(len(result))
	if offset>=int(total){return nil,total,nil}
	end:=offset+limit; if end>int(total){end=int(total)}
	return result[offset:end],total,nil
}
func (r *MemWorkOrderRepo) GetByID(ctx context.Context, tenantID,id int64) (*domain.WorkOrder,error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	wo,ok:=r.orders[id]
	if !ok||wo.TenantID!=tenantID{return nil,nil}
	return wo,nil
}

// Update modifies an existing work order.
func (r *MemWorkOrderRepo) Update(ctx context.Context, wo *domain.WorkOrder) error {
	r.mu.Lock(); defer r.mu.Unlock()
	wo.UpdatedAt=time.Now()
	r.orders[wo.ID]=wo; return nil
}

// Delete removes a work order.
func (r *MemWorkOrderRepo) Delete(ctx context.Context, tenantID, id int64) error {
	r.mu.Lock(); defer r.mu.Unlock()
	if wo, ok := r.orders[id]; ok && wo.TenantID == tenantID {
		delete(r.orders, id)
	}
	return nil
}
