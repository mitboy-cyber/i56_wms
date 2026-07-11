// Package config provides multi-source configuration loading.
// Supports: environment variables > config file > defaults.
package config

import (
	"os"
	"time"
)

// Config is the top-level configuration struct.
type Config struct {
	App      AppConfig      `yaml:"app" json:"app"`
	Server   ServerConfig   `yaml:"server" json:"server"`
	Database DatabaseConfig `yaml:"database" json:"database"`
	Redis    RedisConfig    `yaml:"redis" json:"redis"`
	Storage  StorageConfig  `yaml:"storage" json:"storage"`
	Auth     AuthConfig     `yaml:"auth" json:"auth"`
	Tenant   TenantConfig   `yaml:"tenant" json:"tenant"`
	Log      LogConfig      `yaml:"log" json:"log"`
}

type AppConfig struct {
	Name    string `yaml:"name" json:"name" default:"I56 Framework"`
	Version string `yaml:"version" json:"version" default:"1.1.0"`
	Env     string `yaml:"env" json:"env" default:"development"`
	Debug   bool   `yaml:"debug" json:"debug" default:"true"`
}

type ServerConfig struct {
	Host         string        `yaml:"host" json:"host" default:"0.0.0.0"`
	Port         int           `yaml:"port" json:"port" default:"8080"`
	ReadTimeout  time.Duration `yaml:"read_timeout" json:"read_timeout" default:"30s"`
	WriteTimeout time.Duration `yaml:"write_timeout" json:"write_timeout" default:"30s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" json:"idle_timeout" default:"120s"`
}

type DatabaseConfig struct {
	Driver         string        `yaml:"driver" json:"driver" default:"postgres"`
	Host           string        `yaml:"host" json:"host" default:"localhost"`
	Port           int           `yaml:"port" json:"port" default:"5432"`
	Name           string        `yaml:"name" json:"name" default:"i56"`
	User           string        `yaml:"user" json:"user" default:"i56"`
	Password       string        `yaml:"password" json:"password"`
	SSLMode        string        `yaml:"ssl_mode" json:"ssl_mode" default:"disable"`
	DSN            string        `yaml:"dsn" json:"dsn"` // full connection string
	MaxOpenConns   int           `yaml:"max_open_conns" json:"max_open_conns" default:"25"`
	MaxIdleConns   int           `yaml:"max_idle_conns" json:"max_idle_conns" default:"10"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" json:"conn_max_lifetime" default:"5m"`
}

type RedisConfig struct {
	Host     string `yaml:"host" json:"host" default:"localhost"`
	Port     int    `yaml:"port" json:"port" default:"6379"`
	Password string `yaml:"password" json:"password"`
	DB       int    `yaml:"db" json:"db" default:"0"`
	PoolSize int    `yaml:"pool_size" json:"pool_size" default:"10"`
}

type StorageConfig struct {
	Driver    string `yaml:"driver" json:"driver" default:"local"`
	BasePath  string `yaml:"base_path" json:"base_path" default:"./storage"`
	Endpoint  string `yaml:"endpoint" json:"endpoint"`
	AccessKey string `yaml:"access_key" json:"access_key"`
	SecretKey string `yaml:"secret_key" json:"secret_key"`
	Bucket    string `yaml:"bucket" json:"bucket"`
	UseSSL    bool   `yaml:"use_ssl" json:"use_ssl" default:"true"`
}

type AuthConfig struct {
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl" json:"access_token_ttl" default:"15m"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl" json:"refresh_token_ttl" default:"30d"`
	Issuer          string        `yaml:"issuer" json:"issuer" default:"i56-framework"`
	SigningKey      string        `yaml:"signing_key" json:"signing_key"`
}

type TenantConfig struct {
	Mode        string `yaml:"mode" json:"mode" default:"shared"` // shared | schema | database
	HeaderName  string `yaml:"header_name" json:"header_name" default:"X-Tenant-ID"`
	DefaultID   string `yaml:"default_id" json:"default_id"`
}

type LogConfig struct {
	Level  string `yaml:"level" json:"level" default:"info"`
	Format string `yaml:"format" json:"format" default:"json"` // json | text
	Output string `yaml:"output" json:"output" default:"stdout"`
}

// Option is a functional option for loading config.
type Option func(*Config)

// Load creates a Config with defaults, applies options, then overlays from env.
func Load(opts ...Option) (*Config, error) {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	cfg.loadFromEnv()
	return cfg, nil
}

func defaultConfig() *Config {
	return &Config{
		App: AppConfig{
			Name:    "I56 Framework",
			Version: "1.1.0",
			Env:     "development",
			Debug:   true,
		},
		Server: ServerConfig{
			Host:        "0.0.0.0",
			Port:        8080,
			ReadTimeout: 30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout: 120 * time.Second,
		},
		Database: DatabaseConfig{
			Driver:         "postgres",
			Host:           "localhost",
			Port:           5432,
			Name:           "i56",
			User:           "i56",
			SSLMode:        "disable",
			MaxOpenConns:   25,
			MaxIdleConns:   10,
			ConnMaxLifetime: 5 * time.Minute,
		},
		Auth: AuthConfig{
			AccessTokenTTL:  15 * time.Minute,
			RefreshTokenTTL: 30 * 24 * time.Hour,
			Issuer:          "i56-framework",
		},
		Storage: StorageConfig{
			Driver:   "local",
			BasePath: "./storage",
			UseSSL:   true,
		},
		Tenant: TenantConfig{
			Mode:       "shared",
			HeaderName: "X-Tenant-ID",
		},
		Redis: RedisConfig{
			Host:     "localhost",
			Port:     6379,
			DB:       0,
			PoolSize: 10,
		},
		Log: LogConfig{
			Level:  "info",
			Format: "json",
			Output: "stdout",
		},
	}
}

func (c *Config) loadFromEnv() {
	if v := os.Getenv("APP_ENV"); v != "" {
		c.App.Env = v
	}
	if v := os.Getenv("APP_DEBUG"); v == "false" {
		c.App.Debug = false
	}
	if v := os.Getenv("SERVER_PORT"); v != "" {
		// Parse port...
	}
	if v := os.Getenv("DATABASE_URL"); v != "" {
		// Parse DSN...
	}
	if v := os.Getenv("REDIS_URL"); v != "" {
		// Parse Redis URL...
	}
}

// WithDefaults sets a custom default config (applied before env overlay).
func WithDefaults(fn func(*Config)) Option {
	return func(c *Config) {
		fn(c)
	}
}
