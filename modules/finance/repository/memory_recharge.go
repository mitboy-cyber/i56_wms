package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/modules/finance/domain"
)

// MemRechargeRepo is an in-memory implementation for recharges.
type MemRechargeRepo struct {
	mu        sync.RWMutex
	recharges map[int64]*domain.Recharge
	nextID    int64
}

func NewMemRechargeRepo() *MemRechargeRepo {
	r := &MemRechargeRepo{recharges: make(map[int64]*domain.Recharge), nextID: 1}
	r.seed()
	return r
}

func (r *MemRechargeRepo) seed() {
	now := time.Now()
	seeds := []struct {
		clientID  int64
		rechNo    string
		amount    float64
		method    domain.RechargeMethod
		status    domain.RechargeStatus
		remark    string
	}{
		{1, "RCH202401001", 5000.00, domain.RechargeMethodBankTransfer, domain.RechargeStatusConfirmed, "银行转账充值"},
		{1, "RCH202401002", 3000.00, domain.RechargeMethodWechat, domain.RechargeStatusConfirmed, "微信充值"},
		{2, "RCH202401003", 8000.00, domain.RechargeMethodAlipay, domain.RechargeStatusConfirmed, "支付宝充值"},
		{3, "RCH202401004", 2000.00, domain.RechargeMethodOffline, domain.RechargeStatusPending, "线下充值-待审核"},
		{2, "RCH202401005", 1000.00, domain.RechargeMethodWechat, domain.RechargeStatusRejected, "微信-审核驳回"},
	}
	for _, s := range seeds {
		id := atomic.AddInt64(&r.nextID, 1) - 1
		r.recharges[id] = &domain.Recharge{
			ID:         id,
			TenantID:   1,
			ClientID:   s.clientID,
			RechargeNo: s.rechNo,
			Amount:     s.amount,
			Method:     s.method,
			Status:     s.status,
			Remark:     s.remark,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
	}
}

func (r *MemRechargeRepo) Create(ctx context.Context, recharge *domain.Recharge) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	recharge.ID = atomic.AddInt64(&r.nextID, 1) - 1
	now := time.Now()
	recharge.CreatedAt = now
	recharge.UpdatedAt = now
	r.recharges[recharge.ID] = recharge
	return nil
}

func (r *MemRechargeRepo) GetByID(ctx context.Context, id int64) (*domain.Recharge, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rc, ok := r.recharges[id]
	if !ok {
		return nil, nil
	}
	return rc, nil
}

func (r *MemRechargeRepo) GetByRechargeNo(ctx context.Context, tenantID int64, rechNo string) (*domain.Recharge, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, rc := range r.recharges {
		if rc.TenantID == tenantID && rc.RechargeNo == rechNo {
			return rc, nil
		}
	}
	return nil, nil
}

func (r *MemRechargeRepo) ListByClient(ctx context.Context, tenantID, clientID int64) ([]domain.Recharge, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.Recharge
	for _, rc := range r.recharges {
		if rc.TenantID == tenantID && rc.ClientID == clientID {
			result = append(result, *rc)
		}
	}
	return result, nil
}

func (r *MemRechargeRepo) ListByStatus(ctx context.Context, tenantID int64, status domain.RechargeStatus) ([]domain.Recharge, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.Recharge
	for _, rc := range r.recharges {
		if rc.TenantID == tenantID && rc.Status == status {
			result = append(result, *rc)
		}
	}
	return result, nil
}

func (r *MemRechargeRepo) Update(ctx context.Context, recharge *domain.Recharge) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.recharges[recharge.ID]; !ok {
		return nil
	}
	recharge.UpdatedAt = time.Now()
	r.recharges[recharge.ID] = recharge
	return nil
}

var _ RechargeRepository = (*MemRechargeRepo)(nil)
