package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// LoadFromEnv overrides config defaults from environment variables.
// Supports: APP_PORT, DATABASE_HOST, REDIS_URL, etc.
func LoadFromEnv(cfg *Config) *Config {
	if v := os.Getenv("APP_ENV"); v != "" { cfg.App.Env = v }
	if v := os.Getenv("APP_NAME"); v != "" { cfg.App.Name = v }
	if v := os.Getenv("APP_DEBUG"); v != "" { cfg.App.Debug = v == "true" || v == "1" }
	if v := os.Getenv("SERVER_PORT"); v != "" { cfg.Server.Port, _ = strconv.Atoi(v) }
	if v := os.Getenv("SERVER_HOST"); v != "" { cfg.Server.Host = v }
	if v := os.Getenv("DATABASE_HOST"); v != "" { cfg.Database.Host = v }
	if v := os.Getenv("DATABASE_PORT"); v != "" { cfg.Database.Port, _ = strconv.Atoi(v) }
	if v := os.Getenv("DATABASE_NAME"); v != "" { cfg.Database.Name = v }
	if v := os.Getenv("DATABASE_USER"); v != "" { cfg.Database.User = v }
	if v := os.Getenv("DATABASE_PASSWORD"); v != "" { cfg.Database.Password = v }
	if v := os.Getenv("DATABASE_SSL_MODE"); v != "" { cfg.Database.SSLMode = v }
	if v := os.Getenv("DATABASE_URL"); v != "" { cfg.Database.DSN = v }
	if v := os.Getenv("REDIS_HOST"); v != "" { cfg.Redis.Host = v }
	if v := os.Getenv("REDIS_PORT"); v != "" { cfg.Redis.Port, _ = strconv.Atoi(v) }
	if v := os.Getenv("REDIS_DB"); v != "" { cfg.Redis.DB, _ = strconv.Atoi(v) }
	if v := os.Getenv("AUTH_ISSUER"); v != "" { cfg.Auth.Issuer = v }
	if v := os.Getenv("STORAGE_DRIVER"); v != "" { cfg.Storage.Driver = v }
	if v := os.Getenv("STORAGE_BASE_PATH"); v != "" { cfg.Storage.BasePath = v }
	if v := os.Getenv("LOG_LEVEL"); v != "" { cfg.Log.Level = v }
	if v := os.Getenv("LOG_FORMAT"); v != "" { cfg.Log.Format = v }

	// Load .env file if present
	loadEnvFile(".env", cfg)
	return cfg
}

func loadEnvFile(path string, cfg *Config) {
	f, err := os.Open(path)
	if err != nil { return }
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") { continue }
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 { continue }
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		val = strings.Trim(val, `"'`)
		os.Setenv(key, val)
	}
	LoadFromEnv(cfg)
}

// LoadWithEnv loads config with defaults, overridden by .env and OS env.
