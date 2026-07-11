package repository
import (
	"context"
	parcelDomain "github.com/i56/modules/parcel/domain"
)

func (r *MemParcelRepo) ListByStatus(ctx context.Context, tenantID int64, status parcelDomain.ParcelStatus) []parcelDomain.Parcel {
	r.mu.RLock(); defer r.mu.RUnlock()
	var result []parcelDomain.Parcel
	for _, p := range r.parcels { if p.TenantID == tenantID && string(p.Status) == string(status) { result = append(result, *p) } }
	return result
}

