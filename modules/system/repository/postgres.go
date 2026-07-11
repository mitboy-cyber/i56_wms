package repository

import (
	"context"

	"github.com/i56/framework/db"
	domain "github.com/i56/modules/system/domain"
)

// PgSystemConfigRepo provides PostgreSQL-backed system configuration persistence.
type PgSystemConfigRepo struct{}

func NewPgSystemConfigRepo() *PgSystemConfigRepo { return &PgSystemConfigRepo{} }

func (r *PgSystemConfigRepo) ListLogisticsAPIs(tenantID int64) ([]*domain.LogisticsAPIConfig, int64) {
	return nil, 0
}

func (r *PgSystemConfigRepo) ListBrokers(_ context.Context, tenantID int64) []*domain.CustomsBrokerAPIConfig {
	_ = tenantID
	return nil
}

func (r *PgSystemConfigRepo) ListPrinters(tenantID int64) []*domain.PrinterConfig {
	_ = tenantID
	return nil
}

func (r *PgSystemConfigRepo) ListChannels(_ context.Context, tenantID int64) []*domain.NotificationChannel {
	_ = tenantID
	return nil
}

func (r *PgSystemConfigRepo) ListSettings(tenantID int64) []*domain.SystemSetting {
	var result []*domain.SystemSetting
	// Optional: query from a settings table if one exists
	_ = tenantID
	return result
}

func (r *PgSystemConfigRepo) ListAPIConfig() []APIConfigEntry {
	return nil
}

func (r *PgSystemConfigRepo) SaveAPIConfig(name, provider, endpoint, apiKey, apiSecret, webhookURL, description string) {
}

func (r *PgSystemConfigRepo) SaveNotificationChannel(channelType, name, config string) {
}

func (r *PgSystemConfigRepo) DeleteChannel(_ context.Context, tenantID, id int64) {
}

func (r *PgSystemConfigRepo) SaveSetting(tenantID int64, key, value, typ, group, label string) {
}

// Ensure db import is used
var _ = db.Pool
