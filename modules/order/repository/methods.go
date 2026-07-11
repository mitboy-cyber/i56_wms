package repository
import (
	"context"
	orderDomain "github.com/i56/modules/order/domain"
)

func (r *MemOrderRepo) ListAll(ctx context.Context, tenantID int64) []orderDomain.Order {
	r.mu.RLock(); defer r.mu.RUnlock()
	var result []orderDomain.Order
	for _, o := range r.orders { if o.TenantID == tenantID { result = append(result, *o) } }
	return result
}

func (r *MemOrderRepo) ListByStatus(ctx context.Context, tenantID int64, status orderDomain.OrderStatus) []orderDomain.Order {
	r.mu.RLock(); defer r.mu.RUnlock()
	var result []orderDomain.Order
	for _, o := range r.orders { if o.TenantID == tenantID && o.Status == status { result = append(result, *o) } }
	return result
}

