// Package anomaly provides parcel anomaly detection for the I56 WMS AI subsystem.
// It implements 5 detection rules that scan parcels for weight mismatches,
// suspicious addresses, duplicate tracking, missing customs info, and
// high-value packages without insurance.
package anomaly

import (
	"fmt"
	"strings"
	"sync"
	"time"

	parcelDomain "github.com/i56/modules/parcel/domain"
)

// AnomalyDetector scans parcels against a set of detection rules
// and returns any discovered anomalies.
type AnomalyDetector struct {
	mu    sync.RWMutex
	rules []DetectionRule
}

// DetectionRule defines a named check that inspects a ParcelInfo
// and optionally returns an Anomaly if the check fails.
type DetectionRule struct {
	Name     string
	Check    func(info ParcelInfo) *Anomaly
	Severity string // low / medium / high / critical
}

// Anomaly represents a single detected issue on a parcel.
type Anomaly struct {
	ID         string    `json:"id"`
	ParcelID   int64     `json:"parcel_id"`
	Type       string    `json:"type"` // weight_mismatch, suspicious_address, duplicate_tracking, missing_customs, high_value_no_insurance
	Message    string    `json:"message"`
	Severity   string    `json:"severity"`
	DetectedAt time.Time `json:"detected_at"`
}

// ParcelInfo extends the domain Parcel with optional fields needed for
// anomaly detection (declared weight, HS code, declared value, insurance,
// recipient address). Fields that aren't directly on the domain model
// can be populated from external sources.
type ParcelInfo struct {
	parcelDomain.Parcel

	DeclaredWeight  float64 `json:"declared_weight,omitempty"`
	DeclaredValue   float64 `json:"declared_value,omitempty"`
	HSCode          string  `json:"hs_code,omitempty"`
	HasInsurance    bool    `json:"has_insurance,omitempty"`
	RecipientAddr   string  `json:"recipient_addr,omitempty"`
	CustomsInfoComplete bool `json:"customs_info_complete,omitempty"`
}

// New creates a new AnomalyDetector with the 5 built-in rules.
func New() *AnomalyDetector {
	d := &AnomalyDetector{}
	d.rules = []DetectionRule{
		{
			Name:     "weight_mismatch",
			Severity: "high",
			Check:    d.checkWeightMismatch,
		},
		{
			Name:     "suspicious_address",
			Severity: "medium",
			Check:    d.checkSuspiciousAddress,
		},
		{
			Name:     "duplicate_tracking",
			Severity: "high",
			Check:    d.checkDuplicateTracking,
		},
		{
			Name:     "missing_customs_info",
			Severity: "medium",
			Check:    d.checkMissingCustomsInfo,
		},
		{
			Name:     "high_value_no_insurance",
			Severity: "high",
			Check:    d.checkHighValueNoInsurance,
		},
	}
	return d
}

// AddRule appends a custom detection rule to the detector.
func (d *AnomalyDetector) AddRule(rule DetectionRule) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.rules = append(d.rules, rule)
}

// Scan runs all detection rules against a single ParcelInfo and returns
// any anomalies found.
func (d *AnomalyDetector) Scan(info ParcelInfo) []Anomaly {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var results []Anomaly
	for _, rule := range d.rules {
		if a := rule.Check(info); a != nil {
			results = append(results, *a)
		}
	}
	return results
}

// ScanAll runs all detection rules against every parcel and returns the
// full list of anomalies. It also provides a `visited` callback that
// can be used to build a duplicate-tracking map (the duplicate rule
// requires cross-parcel knowledge).
func (d *AnomalyDetector) ScanAll(parcels []ParcelInfo) []Anomaly {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// Build tracking number → clientID map for duplicate detection
	trackingMap := make(map[string][]int64) // trackingNumber → []clientID

	var all []Anomaly
	for _, info := range parcels {
		// Accumulate tracking numbers for cross-parcel duplicate check
		trackingMap[info.TrackingNumber] = append(trackingMap[info.TrackingNumber], info.ClientID)

		for _, rule := range d.rules {
			// Skip duplicate_tracking in per-parcel scan — handled after loop
			if rule.Name == "duplicate_tracking" {
				continue
			}
			if a := rule.Check(info); a != nil {
				all = append(all, *a)
			}
		}
	}

	// Duplicate tracking: same tracking number, different clients
	anomalyID := 0
	for trackingNum, clientIDs := range trackingMap {
		if len(clientIDs) < 2 {
			continue
		}
		// Check if they belong to different clients
		seen := make(map[int64]bool)
		for _, cid := range clientIDs {
			seen[cid] = true
		}
		if len(seen) <= 1 {
			continue
		}
		for _, info := range parcels {
			if info.TrackingNumber == trackingNum {
				anomalyID++
				all = append(all, Anomaly{
					ID:         fmt.Sprintf("ANM-%06d", anomalyID),
					ParcelID:   info.ID,
					Type:       "duplicate_tracking",
					Message:    fmt.Sprintf("快递单号 %s 被多个客户使用 (%d 个不同客户)", trackingNum, len(seen)),
					Severity:   "high",
					DetectedAt: time.Now(),
				})
				break // One anomaly per duplicate cluster
			}
		}
	}

	// Assign IDs
	for i := range all {
		if all[i].ID == "" {
			anomalyID++
			all[i].ID = fmt.Sprintf("ANM-%06d", anomalyID)
		}
	}

	return all
}

// ─── Rule 1: Weight mismatch — declared_weight vs actual_weight > 30% ───

func (d *AnomalyDetector) checkWeightMismatch(info ParcelInfo) *Anomaly {
	if info.DeclaredWeight <= 0 || info.ActualWeight <= 0 {
		return nil
	}
	diff := info.DeclaredWeight - info.ActualWeight
	if diff < 0 {
		diff = -diff
	}
	pct := diff / info.DeclaredWeight * 100
	if pct > 30 {
		return &Anomaly{
			ParcelID: info.ID,
			Type:     "weight_mismatch",
			Message: fmt.Sprintf("申报重量 %.2fkg 与实际重量 %.2fkg 差异 %.1f%% (超过30%%阈值)",
				info.DeclaredWeight, info.ActualWeight, pct),
			Severity:   "high",
			DetectedAt: time.Now(),
		}
	}
	return nil
}

// ─── Rule 2: Suspicious address — keywords like PO Box, 代收, 转寄 ───

var suspiciousKeywords = []string{"PO Box", "P.O.Box", "代收", "转寄", "转运仓", "中转"}

func (d *AnomalyDetector) checkSuspiciousAddress(info ParcelInfo) *Anomaly {
	addr := info.RecipientAddr
	if addr == "" {
		return nil
	}
	lower := strings.ToLower(addr)
	for _, kw := range suspiciousKeywords {
		if strings.Contains(lower, strings.ToLower(kw)) {
			return &Anomaly{
				ParcelID: info.ID,
				Type:     "suspicious_address",
				Message:  fmt.Sprintf("收件地址包含可疑关键词: %q (地址: %s)", kw, addr),
				Severity:   "medium",
				DetectedAt: time.Now(),
			}
		}
	}
	return nil
}

// ─── Rule 3: Duplicate tracking — handled in ScanAll ───

func (d *AnomalyDetector) checkDuplicateTracking(info ParcelInfo) *Anomaly {
	// This is handled at the ScanAll level with cross-parcel knowledge.
	return nil
}

// ─── Rule 4: Missing customs info — no HS code, no declared value ───

func (d *AnomalyDetector) checkMissingCustomsInfo(info ParcelInfo) *Anomaly {
	missing := []string{}
	if info.HSCode == "" {
		missing = append(missing, "HS编码")
	}
	if info.DeclaredValue <= 0 {
		missing = append(missing, "申报价值")
	}
	if !info.CustomsInfoComplete {
		missing = append(missing, "报关信息不完整")
	}
	if len(missing) > 0 {
		return &Anomaly{
			ParcelID: info.ID,
			Type:     "missing_customs",
			Message:  fmt.Sprintf("缺少报关信息: %s", strings.Join(missing, "、")),
			Severity:   "medium",
			DetectedAt: time.Now(),
		}
	}
	return nil
}

// ─── Rule 5: High value no insurance — declared value > ¥1000 but no insurance ──

func (d *AnomalyDetector) checkHighValueNoInsurance(info ParcelInfo) *Anomaly {
	if info.DeclaredValue > 1000 && !info.HasInsurance {
		return &Anomaly{
			ParcelID: info.ID,
			Type:     "high_value_no_insurance",
			Message:  fmt.Sprintf("申报价值 ¥%.2f 超过 ¥1000 但未购买保险", info.DeclaredValue),
			Severity:   "high",
			DetectedAt: time.Now(),
		}
	}
	return nil
}
