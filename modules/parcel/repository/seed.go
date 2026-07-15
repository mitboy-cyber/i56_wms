package repository

import (
	"sync/atomic"
	"time"

	"github.com/i56/modules/parcel/domain"
)

// seedParcels populates the repo with initial parcel data.
func (r *MemParcelRepo) seedParcels() {
	seeds := []struct {
		tenantID    int64
		warehouseID int64
		clientID    int64
		trackingNo string
		courierCode string
		cargoType   string
		productName string
		parcelName  string
		weight      float64
		length      float64
		width       float64
		height      float64
		status      domain.ParcelStatus
		location    string
	}{
		{1, 1, 1, "TN20240101001", "SF", "general", "手机壳", "手机壳-黑色", 0.15, 15, 10, 2, domain.StatusStored, "A-01-01"},
		{1, 1, 1, "TN20240101002", "SF", "general", "数据线", "数据线-1m", 0.08, 12, 8, 2, domain.StatusPacked, "A-01-02"},
		{1, 1, 2, "TN20240101003", "YT", "fragile", "陶瓷杯", "陶瓷杯-白色", 0.35, 12, 12, 10, domain.StatusReceived, "A-02-01"},
		{1, 1, 2, "TN20240101004", "YT", "general", "T恤", "T恤-L码", 0.25, 25, 20, 5, domain.StatusShipped, "B-01-03"},
		{1, 1, 3, "TN20240101005", "ZTO", "liquid", "洗发水", "洗发水-500ml", 0.55, 18, 8, 8, domain.StatusDelivered, "B-02-01"},
	}
	now := time.Now()
	for _, s := range seeds {
		id := atomic.AddInt64(&r.nextID, 1) - 1
		r.parcels[id] = &domain.Parcel{
			ID:             id,
			TenantID:       s.tenantID,
			WarehouseID:    s.warehouseID,
			ClientID:       s.clientID,
			TrackingNumber: s.trackingNo,
			CourierCode:    s.courierCode,
			CargoType:      s.cargoType,
			ProductName:    s.productName,
			ParcelName:     s.parcelName,
			ActualWeight:   s.weight,
			Length:         s.length,
			Width:          s.width,
			Height:         s.height,
			Status:         s.status,
			LocationCode:   s.location,
			CreatedAt:      now,
			UpdatedAt:      now,
		}
	}
}
