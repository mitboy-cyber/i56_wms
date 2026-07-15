package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/i56/modules/order/domain"
)

// MemServiceOrderRepo is an in-memory implementation for service orders.
type MemServiceOrderRepo struct {
	mu      sync.RWMutex
	orders  map[int64]*domain.ServiceOrder
	nextID  int64
}

func NewMemServiceOrderRepo() *MemServiceOrderRepo {
	r := &MemServiceOrderRepo{orders: make(map[int64]*domain.ServiceOrder), nextID: 1}
	r.seed()
	return r
}

func (r *MemServiceOrderRepo) seed() {
	seeds := []struct {
		tenantID    int64
		clientID    int64
		parcelID    int64
		orderID     int64
		serviceType domain.ServiceOrderType
		serviceName string
		status      domain.ServiceOrderStatus
		price       float64
		operatorID  int64
		resultNote  string
		resultImages []string
	}{
		{1, 1, 1, 1, domain.ServiceTypePhotos, "拍照服务", domain.ServStatusCompleted, 5.00, 1, "包裹外观完好", []string{"photo_001.jpg", "photo_002.jpg"}},
		{1, 1, 2, 1, domain.ServiceTypeRemoveBox, "拆箱服务", domain.ServStatusProcessing, 8.00, 2, "正在处理", nil},
		{1, 2, 3, 2, domain.ServiceTypeReinforce, "加固包装", domain.ServStatusPending, 15.00, 0, "", nil},
		{1, 2, 4, 2, domain.ServiceTypeInspection, "验货服务", domain.ServStatusCompleted, 10.00, 1, "货物与订单一致", []string{"inspect_001.jpg"}},
		{1, 3, 5, 3, domain.ServiceTypeInsurance, "保价服务-200元", domain.ServStatusCompleted, 20.00, 3, "已保价", nil},
	}
	now := time.Now()
	for _, s := range seeds {
		id := atomic.AddInt64(&r.nextID, 1) - 1
		var completedAt *time.Time
		if s.status == domain.ServStatusCompleted {
			completedAt = &now
		}
		r.orders[id] = &domain.ServiceOrder{
			ID:           id,
			TenantID:     s.tenantID,
			ClientID:     s.clientID,
			ParcelID:     s.parcelID,
			OrderID:      s.orderID,
			ServiceType:  s.serviceType,
			ServiceName:  s.serviceName,
			Status:       s.status,
			Price:        s.price,
			OperatorID:   s.operatorID,
			ResultNote:   s.resultNote,
			ResultImages: s.resultImages,
			CreatedAt:    now,
			UpdatedAt:    now,
			CompletedAt:  completedAt,
		}
	}
}

func (r *MemServiceOrderRepo) Create(ctx context.Context, o *domain.ServiceOrder) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	o.ID = atomic.AddInt64(&r.nextID, 1) - 1
	now := time.Now()
	o.CreatedAt = now
	o.UpdatedAt = now
	r.orders[o.ID] = o
	return nil
}

func (r *MemServiceOrderRepo) GetByID(ctx context.Context, tenantID, id int64) (*domain.ServiceOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	o, ok := r.orders[id]
	if !ok || o.TenantID != tenantID {
		return nil, nil
	}
	return o, nil
}

func (r *MemServiceOrderRepo) ListByParcel(ctx context.Context, tenantID, parcelID int64) ([]domain.ServiceOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.ServiceOrder
	for _, o := range r.orders {
		if o.TenantID == tenantID && o.ParcelID == parcelID {
			result = append(result, *o)
		}
	}
	return result, nil
}

func (r *MemServiceOrderRepo) ListByOrder(ctx context.Context, tenantID, orderID int64) ([]domain.ServiceOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.ServiceOrder
	for _, o := range r.orders {
		if o.TenantID == tenantID && o.OrderID == orderID {
			result = append(result, *o)
		}
	}
	return result, nil
}

func (r *MemServiceOrderRepo) ListByClient(ctx context.Context, tenantID, clientID int64) ([]domain.ServiceOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.ServiceOrder
	for _, o := range r.orders {
		if o.TenantID == tenantID && o.ClientID == clientID {
			result = append(result, *o)
		}
	}
	return result, nil
}

func (r *MemServiceOrderRepo) Update(ctx context.Context, o *domain.ServiceOrder) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.orders[o.ID]; !ok {
		return nil
	}
	o.UpdatedAt = time.Now()
	r.orders[o.ID] = o
	return nil
}

var _ ServiceOrderRepository = (*MemServiceOrderRepo)(nil)
