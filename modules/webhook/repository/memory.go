package repository
import (
	"context"; "sync"; "sync/atomic"; "time"
	"github.com/i56/modules/webhook/domain"
)
type MemWebhookRepo struct {
	mu     sync.RWMutex
	subs   map[int64]*domain.WebhookSubscription
	logs   []domain.WebhookDeliveryLog
	nextID int64
}
func NewMemWebhookRepo() *MemWebhookRepo {
	return &MemWebhookRepo{subs:make(map[int64]*domain.WebhookSubscription),nextID:1}
}
func (r *MemWebhookRepo) CreateSub(ctx context.Context, s *domain.WebhookSubscription) error {
	r.mu.Lock(); defer r.mu.Unlock()
	s.ID=atomic.AddInt64(&r.nextID,1)-1; s.CreatedAt=time.Now()
	r.subs[s.ID]=s; return nil
}
func (r *MemWebhookRepo) ListSubs(ctx context.Context, tenantID int64) ([]domain.WebhookSubscription,error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	var result []domain.WebhookSubscription
	for _,s := range r.subs { if s.TenantID==tenantID { result=append(result,*s) } }
	return result,nil
}
func (r *MemWebhookRepo) FindByEvent(ctx context.Context, event string) []domain.WebhookSubscription {
	r.mu.RLock(); defer r.mu.RUnlock()
	var result []domain.WebhookSubscription
	for _,s := range r.subs { if s.Event==event&&s.IsActive { result=append(result,*s) } }
	return result
}
func (r *MemWebhookRepo) LogDelivery(ctx context.Context, log *domain.WebhookDeliveryLog) {
	r.mu.Lock(); defer r.mu.Unlock()
	log.DeliveredAt=time.Now(); r.logs=append(r.logs,*log)
}
func (r *MemWebhookRepo) ListLogs(ctx context.Context, limit int) []domain.WebhookDeliveryLog {
	r.mu.RLock(); defer r.mu.RUnlock()
	if limit<=0||limit>len(r.logs){limit=len(r.logs)}
	return r.logs[len(r.logs)-limit:]
}
