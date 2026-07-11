// Package audit provides operation audit logging.
package audit

import (
	"context"
	"time"

	"github.com/i56/framework/core/logger"
)

// Entry represents a single audit log record.
type Entry struct {
	ID         string         `json:"id"`
	TenantID   string         `json:"tenant_id"`
	UserID     string         `json:"user_id"`
	Action     string         `json:"action"`
	Resource   string         `json:"resource"`
	ResourceID string         `json:"resource_id"`
	Details    map[string]any `json:"details,omitempty"`
	IP         string         `json:"ip"`
	UserAgent  string         `json:"user_agent"`
	Timestamp  time.Time      `json:"timestamp"`
}

// AuditFilter defines query parameters for audit log retrieval.
type AuditFilter struct {
	TenantID   string
	UserID     string
	Resource   string
	ResourceID string
	Action     string
	From       time.Time
	To         time.Time
	Limit      int
	Offset     int
}

// Storage defines how audit entries are persisted.
type Storage interface {
	Save(ctx context.Context, entry Entry) error
	Query(ctx context.Context, filter AuditFilter) ([]Entry, int64, error)
}

// Service provides audit logging capabilities.
type Service struct {
	log     logger.Logger
	storage Storage
}

// NewService creates an audit service.
func NewService(log logger.Logger, storage Storage) *Service {
	return &Service{log: log, storage: storage}
}

// Log records an audit entry.
func (s *Service) Log(ctx context.Context, entry Entry) error {
	entry.Timestamp = time.Now()
	s.log.Info("audit",
		"action", entry.Action,
		"resource", entry.Resource,
		"resource_id", entry.ResourceID,
		"user_id", entry.UserID,
		"tenant_id", entry.TenantID,
	)
	return s.storage.Save(ctx, entry)
}

// Query retrieves audit entries matching the filter.
func (s *Service) Query(ctx context.Context, filter AuditFilter) ([]Entry, int64, error) {
	return s.storage.Query(ctx, filter)
}

// Common action names for consistency.
const (
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
	ActionView   = "view"
	ActionList   = "list"
	ActionExport = "export"
	ActionImport = "import"
	ActionLogin  = "login"
	ActionLogout = "logout"
	ActionApprove = "approve"
	ActionReject  = "reject"
	ActionCancel  = "cancel"
	ActionPrint   = "print"
)
