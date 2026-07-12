// Package core provides AI cost tracking and business context injection
// for the I56 WMS AI subsystem.
package core

import (
	"sync"
	"time"
)

// CostTracker tracks AI model usage and estimated costs.
// It is safe for concurrent use.
type CostTracker struct {
	mu   sync.RWMutex
	logs []CostLog
}

// CostLog records a single AI inference call and its estimated cost.
type CostLog struct {
	Timestamp        time.Time `json:"timestamp"`
	Model            string    `json:"model"`
	Module           string    `json:"module"`
	TenantID         int64     `json:"tenant_id"`
	PromptTokens     int       `json:"prompt_tokens"`
	CompletionTokens int       `json:"completion_tokens"`
	EstimatedCost    float64   `json:"estimated_cost"`
}

// NewCostTracker creates a new CostTracker.
func NewCostTracker() *CostTracker {
	return &CostTracker{}
}

// Track records a cost entry for an AI inference call.
func (ct *CostTracker) Track(model, module string, tenantID int64, promptTokens, completionTokens int, cost float64) {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	ct.logs = append(ct.logs, CostLog{
		Timestamp:        time.Now(),
		Model:            model,
		Module:           module,
		TenantID:         tenantID,
		PromptTokens:     promptTokens,
		CompletionTokens: completionTokens,
		EstimatedCost:    cost,
	})
}

// Query returns cost log entries for the given tenant since the specified time.
// If since is zero, returns all entries.
func (ct *CostTracker) Query(tenantID int64, since time.Time) []CostLog {
	ct.mu.RLock()
	defer ct.mu.RUnlock()
	var result []CostLog
	for _, l := range ct.logs {
		if tenantID != 0 && l.TenantID != tenantID {
			continue
		}
		if !since.IsZero() && l.Timestamp.Before(since) {
			continue
		}
		result = append(result, l)
	}
	return result
}

// TotalCost returns the total estimated cost for the given tenant since the specified time.
func (ct *CostTracker) TotalCost(tenantID int64, since time.Time) float64 {
	ct.mu.RLock()
	defer ct.mu.RUnlock()
	var total float64
	for _, l := range ct.logs {
		if tenantID != 0 && l.TenantID != tenantID {
			continue
		}
		if !since.IsZero() && l.Timestamp.Before(since) {
			continue
		}
		total += l.EstimatedCost
	}
	return total
}

// Count returns the number of tracked entries.
func (ct *CostTracker) Count() int {
	ct.mu.RLock()
	defer ct.mu.RUnlock()
	return len(ct.logs)
}
