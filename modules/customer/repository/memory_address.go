package repository

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/i56/modules/customer/domain"
)

// MemAddressRepo is an in-memory implementation of AddressRepository.
type MemAddressRepo struct {
	mu        sync.RWMutex
	addresses map[int64]*domain.MemberAddress
	nextID    int64
}

func NewMemAddressRepo() *MemAddressRepo {
	r := &MemAddressRepo{addresses: make(map[int64]*domain.MemberAddress), nextID: 1}
	// Seed BFT56-style addresses for client 1, members 1 & 2
	seeds := []struct {
		memberID      int64
		recipientName string
		phone         string
		city          string
		district      string
		address       string
		isDefault     bool
	}{
		{1, "王仁照", "886912345678", "台北市", "信義區", "信義路五段7號101樓", true},
		{1, "王仁照", "886912345678", "新北市", "板橋區", "文化路一段200號", false},
		{2, "吳欣如", "886923456789", "台北市", "大安區", "忠孝東路四段300號5樓", true},
		{2, "吳欣如", "886923456789", "台中市", "西屯區", "台灣大道三段99號", false},
		{2, "張致廷", "886934567890", "高雄市", "前鎮區", "中山四路100號", false},
	}
	for _, s := range seeds {
		id := atomic.AddInt64(&r.nextID, 1) - 1
		r.addresses[id] = &domain.MemberAddress{
			ID:            id,
			MemberID:      s.memberID,
			RecipientName: s.recipientName,
			Phone:         s.phone,
			City:          s.city,
			District:      s.district,
			Address:       s.address,
			IsDefault:     s.isDefault,
		}
	}
	return r
}

func (r *MemAddressRepo) Create(ctx context.Context, memberID int64, a *domain.MemberAddress) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	a.ID = atomic.AddInt64(&r.nextID, 1) - 1
	a.MemberID = memberID
	r.addresses[a.ID] = a
	return nil
}

func (r *MemAddressRepo) List(ctx context.Context, memberID int64) ([]domain.MemberAddress, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.MemberAddress
	for _, a := range r.addresses {
		if a.MemberID == memberID {
			result = append(result, *a)
		}
	}
	return result, nil
}

// ListByClient returns all addresses across all members for a given client.
// This is a convenience method for the client portal.
func (r *MemAddressRepo) ListByClient(ctx context.Context, clientID int64) []domain.MemberAddress {
	_ = clientID // clientID is implicit through members; return all.
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.MemberAddress
	for _, a := range r.addresses {
		result = append(result, *a)
	}
	return result
}

func (r *MemAddressRepo) SetDefault(ctx context.Context, memberID, addressID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, a := range r.addresses {
		if a.MemberID == memberID {
			a.IsDefault = (a.ID == addressID)
		}
	}
	return nil
}

func (r *MemAddressRepo) GetByID(ctx context.Context, id int64) (*domain.MemberAddress, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.addresses[id]
	if !ok { return nil, nil }
	return a, nil
}

func (r *MemAddressRepo) Update(ctx context.Context, id int64, a *domain.MemberAddress) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.addresses[id]; !ok { return nil }
	a.ID = id
	r.addresses[id] = a
	return nil
}

var _ AddressRepository = (*MemAddressRepo)(nil)
