package repository

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/i56/modules/customer/domain"
)

// MemMemberRepo is an in-memory implementation of MemberRepository.
type MemMemberRepo struct {
	mu      sync.RWMutex
	members map[int64]*domain.ClientMember
	nextID  int64
}

func NewMemMemberRepo() *MemMemberRepo {
	return &MemMemberRepo{members: make(map[int64]*domain.ClientMember), nextID: 1}
}

func (r *MemMemberRepo) Create(ctx context.Context, clientID int64, m *domain.ClientMember) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	m.ID = atomic.AddInt64(&r.nextID, 1) - 1
	r.members[m.ID] = m
	return nil
}

func (r *MemMemberRepo) GetByID(ctx context.Context, clientID, id int64) (*domain.ClientMember, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.members[id]
	if !ok || m.ClientID != clientID {
		return nil, nil
	}
	return m, nil
}

func (r *MemMemberRepo) List(ctx context.Context, clientID int64, offset, limit int) ([]domain.ClientMember, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.ClientMember
	for _, m := range r.members {
		if m.ClientID == clientID {
			result = append(result, *m)
		}
	}
	total := int64(len(result))
	if offset >= int(total) {
		return nil, total, nil
	}
	end := offset + limit
	if end > int(total) {
		end = int(total)
	}
	return result[offset:end], total, nil
}

func (r *MemMemberRepo) Update(ctx context.Context, clientID, id int64, m *domain.ClientMember) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.members[id]; !ok { return nil }
	m.ID = id
	m.ClientID = clientID
	r.members[id] = m
	return nil
}

var _ MemberRepository = (*MemMemberRepo)(nil)
