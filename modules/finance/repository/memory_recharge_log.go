package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/modules/finance/domain"
)

// MemRechargeLogRepo is an in-memory implementation for recharge audit logs.
type MemRechargeLogRepo struct {
	mu    sync.RWMutex
	logs  map[int64]*domain.RechargeLog
	nextID int64
}

func NewMemRechargeLogRepo() *MemRechargeLogRepo {
	r := &MemRechargeLogRepo{logs: make(map[int64]*domain.RechargeLog), nextID: 1}
	r.seed()
	return r
}

func (r *MemRechargeLogRepo) seed() {
	now := time.Now()
	seeds := []struct {
		rechargeID int64
		action     string
		operatorID int64
		note       string
	}{
		{1, "submitted", 1, "客户提交充值申请"},
		{1, "confirmed", 2, "财务审核通过"},
		{2, "submitted", 1, "客户提交微信充值"},
		{2, "confirmed", 2, "自动确认"},
		{5, "rejected", 2, "金额不匹配，驳回"},
	}
	for _, s := range seeds {
		id := atomic.AddInt64(&r.nextID, 1) - 1
		r.logs[id] = &domain.RechargeLog{
			ID:         id,
			RechargeID: s.rechargeID,
			Action:     s.action,
			OperatorID: s.operatorID,
			Note:       s.note,
			CreatedAt:  now,
		}
	}
}

func (r *MemRechargeLogRepo) Create(ctx context.Context, log *domain.RechargeLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	log.ID = atomic.AddInt64(&r.nextID, 1) - 1
	log.CreatedAt = time.Now()
	r.logs[log.ID] = log
	return nil
}

func (r *MemRechargeLogRepo) ListByRecharge(ctx context.Context, rechargeID int64) ([]domain.RechargeLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.RechargeLog
	for _, l := range r.logs {
		if l.RechargeID == rechargeID {
			result = append(result, *l)
		}
	}
	return result, nil
}

var _ RechargeLogRepository = (*MemRechargeLogRepo)(nil)
