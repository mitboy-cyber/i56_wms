package repository

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/i56/modules/customer/domain"
)

// MemDeclarantRepo is an in-memory implementation of DeclarantRepository.
type MemDeclarantRepo struct {
	mu         sync.RWMutex
	declarants map[int64]*domain.Declarant
	nextID     int64
}

func NewMemDeclarantRepo() *MemDeclarantRepo {
	return &MemDeclarantRepo{declarants: make(map[int64]*domain.Declarant), nextID: 1}
}

func (r *MemDeclarantRepo) Create(ctx context.Context, clientID int64, d *domain.Declarant) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	d.ID = atomic.AddInt64(&r.nextID, 1) - 1
	r.declarants[d.ID] = d
	return nil
}

func (r *MemDeclarantRepo) GetByID(ctx context.Context, clientID, id int64) (*domain.Declarant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.declarants[id]
	if !ok || d.ClientID != clientID {
		return nil, nil
	}
	return d, nil
}

func (r *MemDeclarantRepo) List(ctx context.Context, clientID int64, offset, limit int) ([]domain.Declarant, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.Declarant
	for _, d := range r.declarants {
		if d.ClientID == clientID {
			result = append(result, *d)
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

func (r *MemDeclarantRepo) Update(ctx context.Context, clientID, id int64, d *domain.Declarant) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.declarants[id]
	if !ok || existing.ClientID != clientID {
		return nil
	}
	r.declarants[id] = d
	return nil
}

var _ DeclarantRepository = (*MemDeclarantRepo)(nil)
