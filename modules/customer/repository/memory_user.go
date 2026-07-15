package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/modules/customer/domain"
)

// MemClientUserRepo is an in-memory implementation of ClientUserRepository.
type MemClientUserRepo struct {
	mu     sync.RWMutex
	users  map[int64]*domain.ClientUser
	nextID int64
}

func NewMemClientUserRepo() *MemClientUserRepo {
	r := &MemClientUserRepo{users: make(map[int64]*domain.ClientUser), nextID: 1}
	r.seed()
	return r
}

func (r *MemClientUserRepo) seed() {
	seeds := []struct {
		clientID   int64
		username   string
		passHash   string
		display    string
		email      string
		phone      string
		role       domain.ClientUserRole
		isActive   bool
	}{
		{1, "admin", "$2a$10$hash_admin", "陈管理员", "admin@shopee-tw.com", "886911111111", domain.ClientUserRoleAdmin, true},
		{1, "operator1", "$2a$10$hash_op1", "林操作员", "lin@shopee-tw.com", "886922222222", domain.ClientUserRoleOperator, true},
		{2, "admin", "$2a$10$hash_admin2", "王经理", "wang@major-tw.com", "886933333333", domain.ClientUserRoleAdmin, true},
		{2, "viewer1", "$2a$10$hash_view1", "黄查看员", "huang@major-tw.com", "886944444444", domain.ClientUserRoleViewer, true},
		{3, "admin", "$2a$10$hash_admin3", "刘负责人", "liu@normal-corp.com", "886955555555", domain.ClientUserRoleAdmin, true},
	}
	now := time.Now()
	for _, s := range seeds {
		id := atomic.AddInt64(&r.nextID, 1) - 1
		r.users[id] = &domain.ClientUser{
			ID:           id,
			ClientID:     s.clientID,
			Username:     s.username,
			PasswordHash: s.passHash,
			DisplayName:  s.display,
			Email:        s.email,
			Phone:        s.phone,
			Role:         s.role,
			IsActive:     s.isActive,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
	}
}

func (r *MemClientUserRepo) Create(ctx context.Context, clientID int64, u *domain.ClientUser) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	u.ID = atomic.AddInt64(&r.nextID, 1) - 1
	u.ClientID = clientID
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	r.users[u.ID] = u
	return nil
}

func (r *MemClientUserRepo) GetByID(ctx context.Context, id int64) (*domain.ClientUser, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[id]
	if !ok {
		return nil, nil
	}
	return u, nil
}

func (r *MemClientUserRepo) GetByUsername(ctx context.Context, clientID int64, username string) (*domain.ClientUser, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.users {
		if u.ClientID == clientID && u.Username == username {
			return u, nil
		}
	}
	return nil, nil
}

func (r *MemClientUserRepo) ListByClient(ctx context.Context, clientID int64) ([]domain.ClientUser, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.ClientUser
	for _, u := range r.users {
		if u.ClientID == clientID {
			result = append(result, *u)
		}
	}
	return result, nil
}

func (r *MemClientUserRepo) Update(ctx context.Context, id int64, u *domain.ClientUser) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[id]; !ok {
		return nil
	}
	u.ID = id
	u.UpdatedAt = time.Now()
	r.users[id] = u
	return nil
}

func (r *MemClientUserRepo) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.users, id)
	return nil
}

func (r *MemClientUserRepo) RecordLogin(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.users[id]
	if !ok {
		return nil
	}
	now := time.Now()
	u.LastLoginAt = &now
	return nil
}

var _ ClientUserRepository = (*MemClientUserRepo)(nil)
