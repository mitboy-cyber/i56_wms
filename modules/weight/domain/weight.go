package weight

import (
	"sync"
	"time"
)

// TransportType represents shipping method
type TransportType string

const (
	TransportAirNormal   TransportType = "air_normal"   // 空普
	TransportAirSpecial  TransportType = "air_special"  // 空特
	TransportSeaExpress  TransportType = "sea_express"  // 海快
	TransportSeaFreight  TransportType = "sea_freight"  // 海运
)

// CargoType represents cargo classification
type CargoType string

const (
	CargoTypeNormal   CargoType = "normal"   // 普货
	CargoTypeSpecial  CargoType = "special"  // 特货
)

// WeightRecord represents a weighing record
type WeightRecord struct {
	ID              int64         `json:"id"`
	TenantID        int64         `json:"tenant_id"`
	TrackingNumber  string        `json:"tracking_number"`
	CourierCompany  string        `json:"courier_company"`
	MemberID        string        `json:"member_id"`
	Platform        string        `json:"platform"`
	Weight          float64       `json:"weight_kg"`
	Length          float64       `json:"length_cm"`
	Width           float64       `json:"width_cm"`
	Height          float64       `json:"height_cm"`
	Volume          float64       `json:"volume_cm3"`
	ParcelCount     int           `json:"parcel_count"`
	ParcelType      CargoType     `json:"parcel_type"`
	ProductName     string        `json:"product_name"`
	ParcelName      string        `json:"parcel_name"`
	TransportTypes  []TransportType `json:"transport_types"`
	ImageURL        string        `json:"image_url,omitempty"`
	Remark          string        `json:"remark,omitempty"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

// CreateWeightRecordRequest is the input for creating a record
type CreateWeightRecordRequest struct {
	TrackingNumber string          `json:"tracking_number"`
	CourierCompany string          `json:"courier_company"`
	MemberID       string          `json:"member_id"`
	Platform       string          `json:"platform"`
	Weight         float64         `json:"weight_kg"`
	Length         float64         `json:"length_cm"`
	Width          float64         `json:"width_cm"`
	Height         float64         `json:"height_cm"`
	ParcelCount    int             `json:"parcel_count"`
	ParcelType     CargoType       `json:"parcel_type"`
	ProductName    string          `json:"product_name"`
	ParcelName     string          `json:"parcel_name"`
	TransportTypes []TransportType `json:"transport_types"`
	Remark         string          `json:"remark,omitempty"`
}

// MemWeightRepo is an in-memory repository for weight records
type MemWeightRepo struct {
	mu      sync.RWMutex
	records []WeightRecord
	nextID  int64
}

// NewMemWeightRepo creates a new in-memory weight record repository
func NewMemWeightRepo() *MemWeightRepo {
	return &MemWeightRepo{records: make([]WeightRecord, 0), nextID: 1}
}

// Create adds a new weight record
func (r *MemWeightRepo) Create(tenantID int64, req CreateWeightRecordRequest) *WeightRecord {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	vol := req.Length * req.Width * req.Height

	rec := WeightRecord{
		ID:             r.nextID,
		TenantID:       tenantID,
		TrackingNumber: req.TrackingNumber,
		CourierCompany: req.CourierCompany,
		MemberID:       req.MemberID,
		Platform:       req.Platform,
		Weight:         req.Weight,
		Length:         req.Length,
		Width:          req.Width,
		Height:         req.Height,
		Volume:         vol,
		ParcelCount:    req.ParcelCount,
		ParcelType:     req.ParcelType,
		ProductName:    req.ProductName,
		ParcelName:     req.ParcelName,
		TransportTypes: req.TransportTypes,
		Remark:         req.Remark,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	r.nextID++
	r.records = append(r.records, rec)
	return &rec
}

// List returns weight records with pagination
func (r *MemWeightRepo) List(tenantID int64, offset, limit int) ([]WeightRecord, int64) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []WeightRecord
	for _, rec := range r.records {
		if rec.TenantID == tenantID {
			filtered = append(filtered, rec)
		}
	}

	total := int64(len(filtered))
	if offset >= len(filtered) {
		return []WeightRecord{}, total
	}

	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}

	// Return in reverse chronological order
	result := make([]WeightRecord, end-offset)
	for i := offset; i < end; i++ {
		result[i-offset] = filtered[len(filtered)-1-i]
	}

	return result, total
}

// GetByID returns a single record
func (r *MemWeightRepo) GetByID(tenantID, id int64) *WeightRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, rec := range r.records {
		if rec.ID == id && rec.TenantID == tenantID {
			return &rec
		}
	}
	return nil
}

// Search searches by tracking number or member ID
func (r *MemWeightRepo) Search(tenantID int64, query string, limit int) []WeightRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []WeightRecord
	for _, rec := range r.records {
		if rec.TenantID == tenantID && (contains(rec.TrackingNumber, query) || contains(rec.MemberID, query)) {
			results = append(results, rec)
		}
	}

	// Reverse order
	for i, j := 0, len(results)-1; i < j; i, j = i+1, j-1 {
		results[i], results[j] = results[j], results[i]
	}

	if len(results) > limit {
		results = results[:limit]
	}

	return results
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || searchSubstring(s, sub))
}

func searchSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
