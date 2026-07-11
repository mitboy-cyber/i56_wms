package unit

import (
	"testing"
	"time"

	"github.com/i56/framework/core/config"
)

func TestLoadDefaults(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.App.Name != "I56 Framework" {
		t.Errorf("expected 'I56 Framework', got '%s'", cfg.App.Name)
	}
	if cfg.App.Version != "1.0.0" {
		t.Errorf("expected '1.0.0', got '%s'", cfg.App.Version)
	}
	if cfg.App.Env != "development" {
		t.Errorf("expected 'development', got '%s'", cfg.App.Env)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("expected 8080, got %d", cfg.Server.Port)
	}
}

func TestLoadWithOverride(t *testing.T) {
	cfg, err := config.Load(config.WithDefaults(func(c *config.Config) {
		c.App.Name = "CustomApp"
		c.Server.Port = 9090
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.App.Name != "CustomApp" {
		t.Errorf("expected 'CustomApp', got '%s'", cfg.App.Name)
	}
	if cfg.Server.Port != 9090 {
		t.Errorf("expected 9090, got %d", cfg.Server.Port)
	}
}

func TestAuthDefaults(t *testing.T) {
	cfg, _ := config.Load()
	if cfg.Auth.AccessTokenTTL != 15*time.Minute {
		t.Errorf("expected 15m, got %v", cfg.Auth.AccessTokenTTL)
	}
	if cfg.Auth.RefreshTokenTTL != 30*24*time.Hour {
		t.Errorf("expected 30d, got %v", cfg.Auth.RefreshTokenTTL)
	}
}

func TestDatabaseDefaults(t *testing.T) {
	cfg, _ := config.Load()
	if cfg.Database.Driver != "postgres" {
		t.Errorf("expected 'postgres', got '%s'", cfg.Database.Driver)
	}
	if cfg.Database.MaxOpenConns != 25 {
		t.Errorf("expected 25, got %d", cfg.Database.MaxOpenConns)
	}
}
