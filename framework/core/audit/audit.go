// Package audit provides operation audit logging.
package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

// AuditEntry represents a single audit log record.
type AuditEntry struct {
	ID         int64     `json:"id"`
	TenantID   int64     `json:"tenant_id"`
	UserID     int64     `json:"user_id"`
	Action     string    `json:"action"`
	Resource   string    `json:"resource"`
	ResourceID string    `json:"resource_id"`
	Detail     string    `json:"detail,omitempty"`
	IP         string    `json:"ip"`
	CreatedAt  time.Time `json:"created_at"`
}

// AuditFilter defines query parameters for audit log retrieval.
type AuditFilter struct {
	TenantID   int64     `json:"tenant_id,omitempty"`
	UserID     int64     `json:"user_id,omitempty"`
	Action     string    `json:"action,omitempty"`
	Resource   string    `json:"resource,omitempty"`
	ResourceID string    `json:"resource_id,omitempty"`
	From       time.Time `json:"from,omitempty"`
	To         time.Time `json:"to,omitempty"`
	Limit      int       `json:"limit,omitempty"`
	Offset     int       `json:"offset,omitempty"`
}

// AuditRepo defines the storage interface for audit entries.
type AuditRepo interface {
	Save(ctx context.Context, entry *AuditEntry) error
	Query(ctx context.Context, filter AuditFilter) ([]AuditEntry, int64, error)
}

// AuditLogger provides audit logging capabilities.
type AuditLogger struct {
	repo AuditRepo
}

// New creates a new AuditLogger with the given repo.
func New(repo AuditRepo) *AuditLogger {
	return &AuditLogger{repo: repo}
}

// Log records an audit entry.
func (a *AuditLogger) Log(ctx context.Context, action, resource, resourceID string, detail interface{}) error {
	var detailStr string
	if detail != nil {
		if b, err := json.Marshal(detail); err == nil {
			detailStr = string(b)
		} else {
			detailStr = fmt.Sprintf("%v", detail)
		}
	}

	entry := &AuditEntry{
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Detail:     detailStr,
		CreatedAt:  time.Now(),
	}

	// Extract user info from context if available
	if v := ctx.Value("tenant_id"); v != nil {
		if tid, ok := v.(int64); ok {
			entry.TenantID = tid
		}
	}
	if v := ctx.Value("user_id"); v != nil {
		if uid, ok := v.(int64); ok {
			entry.UserID = uid
		}
	}

	return a.repo.Save(ctx, entry)
}

// Query retrieves audit entries matching the filter.
func (a *AuditLogger) Query(ctx context.Context, filter AuditFilter) ([]AuditEntry, int64, error) {
	return a.repo.Query(ctx, filter)
}

// MemAuditRepo is an in-memory implementation of AuditRepo.
type MemAuditRepo struct {
	mu      sync.RWMutex
	entries []AuditEntry
	nextID  int64
}

// NewMemAuditRepo creates a new in-memory audit repository.
func NewMemAuditRepo() *MemAuditRepo {
	return &MemAuditRepo{
		entries: make([]AuditEntry, 0),
	}
}

// Save stores an audit entry in memory.
func (r *MemAuditRepo) Save(ctx context.Context, entry *AuditEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	entry.ID = atomic.AddInt64(&r.nextID, 1)
	r.entries = append(r.entries, *entry)

	// Keep only last 10000 entries
	if len(r.entries) > 10000 {
		r.entries = r.entries[len(r.entries)-10000:]
	}
	return nil
}

// Query retrieves audit entries matching the filter.
func (r *MemAuditRepo) Query(ctx context.Context, filter AuditFilter) ([]AuditEntry, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []AuditEntry
	for _, e := range r.entries {
		if filter.Action != "" && e.Action != filter.Action {
			continue
		}
		if filter.Resource != "" && e.Resource != filter.Resource {
			continue
		}
		if filter.ResourceID != "" && e.ResourceID != filter.ResourceID {
			continue
		}
		if filter.UserID != 0 && e.UserID != filter.UserID {
			continue
		}
		if filter.TenantID != 0 && e.TenantID != filter.TenantID {
			continue
		}
		if !filter.From.IsZero() && e.CreatedAt.Before(filter.From) {
			continue
		}
		if !filter.To.IsZero() && e.CreatedAt.After(filter.To) {
			continue
		}
		filtered = append(filtered, e)
	}

	// Sort by CreatedAt descending
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
	})

	total := int64(len(filtered))

	// Apply pagination
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	if filter.Offset >= len(filtered) {
		return []AuditEntry{}, total, nil
	}

	end := filter.Offset + filter.Limit
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[filter.Offset:end], total, nil
}

// Common action names for consistency.
const (
	ActionCreate  = "CREATE"
	ActionUpdate  = "UPDATE"
	ActionDelete  = "DELETE"
	ActionView    = "VIEW"
	ActionLogin   = "LOGIN"
	ActionLogout  = "LOGOUT"
	ActionExport  = "EXPORT"
	ActionImport  = "IMPORT"
	ActionApprove = "APPROVE"
	ActionReject  = "REJECT"
	ActionCancel  = "CANCEL"
	ActionPrint   = "PRINT"
)
