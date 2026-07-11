package repository
import (
	"context"; "sync"; "time"
	"github.com/i56/modules/print/domain"
)
type MemPrintRepo struct {
	mu        sync.RWMutex
	templates map[int64]*domain.PrintTemplate
	nextID    int64
}
func NewMemPrintRepo() *MemPrintRepo {
	r:=&MemPrintRepo{templates:make(map[int64]*domain.PrintTemplate),nextID:1}
	for _,t := range domain.DefaultTemplates() {
		t.TenantID=1; t.CreatedAt=time.Now(); t.UpdatedAt=time.Now()
		r.templates[r.nextID]=&t; r.nextID++
	}
	return r
}
func (r *MemPrintRepo) List(ctx context.Context, tenantID int64) ([]domain.PrintTemplate,error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	var result []domain.PrintTemplate
	for _,t := range r.templates { if t.TenantID==tenantID { result=append(result,*t) } }
	return result,nil
}
func (r *MemPrintRepo) SetDefault(ctx context.Context, tenantID,id int64) error {
	r.mu.Lock(); defer r.mu.Unlock()
	for _,t := range r.templates { t.IsDefault=false }
	if t,ok:=r.templates[id]; ok&&t.TenantID==tenantID { t.IsDefault=true }
	return nil
}

// Add creates a new print template.
func (r *MemPrintRepo) Add(ctx context.Context, tenantID int64, name, ttype, content string) error {
	r.mu.Lock(); defer r.mu.Unlock()
	id := r.nextID; r.nextID++
	now := time.Now()
	r.templates[id] = &domain.PrintTemplate{
		ID: id, TenantID: tenantID, Name: name, Type: ttype,
		Content: content, CreatedAt: now, UpdatedAt: now,
	}
	return nil
}

// Delete removes a print template.
func (r *MemPrintRepo) Delete(ctx context.Context, tenantID, id int64) error {
	r.mu.Lock(); defer r.mu.Unlock()
	if t, ok := r.templates[id]; ok && t.TenantID == tenantID {
		delete(r.templates, id)
	}
	return nil
}
