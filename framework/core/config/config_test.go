package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad_Defaults(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if cfg.App.Name != "I56 Framework" {
		t.Errorf("expected 'I56 Framework', got %q", cfg.App.Name)
	}
	if cfg.App.Env != "development" {
		t.Errorf("expected 'development', got %q", cfg.App.Env)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Server.Port)
	}
	if cfg.Database.Host != "localhost" {
		t.Errorf("expected host 'localhost', got %q", cfg.Database.Host)
	}
	if cfg.Auth.Issuer != "i56-framework" {
		t.Errorf("expected issuer 'i56-framework', got %q", cfg.Auth.Issuer)
	}
	if cfg.Tenant.Mode != "shared" {
		t.Errorf("expected tenant mode 'shared', got %q", cfg.Tenant.Mode)
	}
	if cfg.Log.Level != "info" {
		t.Errorf("expected log level 'info', got %q", cfg.Log.Level)
	}
}

func TestLoad_WithDefaults(t *testing.T) {
	cfg, err := Load(WithDefaults(func(c *Config) {
		c.App.Name = "Custom App"
		c.Server.Port = 3000
	}))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if cfg.App.Name != "Custom App" {
		t.Errorf("expected 'Custom App', got %q", cfg.App.Name)
	}
	if cfg.Server.Port != 3000 {
		t.Errorf("expected port 3000, got %d", cfg.Server.Port)
	}
	// Other defaults should still apply
	if cfg.Database.Host != "localhost" {
		t.Errorf("expected host 'localhost', got %q", cfg.Database.Host)
	}
}

func TestLoadFromEnv(t *testing.T) {
	os.Setenv("APP_NAME", "EnvApp")
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("DATABASE_HOST", "db.example.com")
	defer func() {
		os.Unsetenv("APP_NAME")
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("DATABASE_HOST")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	// Re-apply env
	cfg = LoadFromEnv(cfg)

	if cfg.App.Name != "EnvApp" {
		t.Errorf("expected 'EnvApp', got %q", cfg.App.Name)
	}
	if cfg.Server.Port != 9090 {
		t.Errorf("expected port 9090, got %d", cfg.Server.Port)
	}
	if cfg.Database.Host != "db.example.com" {
		t.Errorf("expected 'db.example.com', got %q", cfg.Database.Host)
	}
}

func TestConfig_TenantSettings(t *testing.T) {
	cfg, _ := Load()

	if cfg.Tenant.HeaderName != "X-Tenant-ID" {
		t.Errorf("expected 'X-Tenant-ID', got %q", cfg.Tenant.HeaderName)
	}
}

func TestConfig_StorageSettings(t *testing.T) {
	cfg, _ := Load()

	if cfg.Storage.Driver != "local" {
		t.Errorf("expected 'local', got %q", cfg.Storage.Driver)
	}
	if cfg.Storage.UseSSL != true {
		t.Error("expected UseSSL=true by default")
	}
}

func TestConfig_Timeouts(t *testing.T) {
	cfg, _ := Load()

	if cfg.Server.ReadTimeout != 30*time.Second {
		t.Errorf("expected 30s read timeout, got %v", cfg.Server.ReadTimeout)
	}
	if cfg.Server.WriteTimeout != 30*time.Second {
		t.Errorf("expected 30s write timeout, got %v", cfg.Server.WriteTimeout)
	}
}

func TestConfig_RedisDefaults(t *testing.T) {
	cfg, _ := Load()

	if cfg.Redis.Host != "localhost" {
		t.Errorf("expected redis host 'localhost', got %q", cfg.Redis.Host)
	}
	if cfg.Redis.Port != 6379 {
		t.Errorf("expected redis port 6379, got %d", cfg.Redis.Port)
	}
	if cfg.Redis.PoolSize != 10 {
		t.Errorf("expected pool size 10, got %d", cfg.Redis.PoolSize)
	}
}

func TestConfig_DatabaseDefaults(t *testing.T) {
	cfg, _ := Load()

	if cfg.Database.Driver != "postgres" {
		t.Errorf("expected driver 'postgres', got %q", cfg.Database.Driver)
	}
	if cfg.Database.MaxOpenConns != 25 {
		t.Errorf("expected 25 max open conns, got %d", cfg.Database.MaxOpenConns)
	}
	if cfg.Database.ConnMaxLifetime != 5*time.Minute {
		t.Errorf("expected 5m conn max lifetime, got %v", cfg.Database.ConnMaxLifetime)
	}
}
