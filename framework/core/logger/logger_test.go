package logger

import (
	"bytes"
	"testing"

	"github.com/i56/framework/core/config"
)

func TestInit_JSON(t *testing.T) {
	cfg := &config.Config{
		Log: config.LogConfig{Level: "info", Format: "json", Output: "stdout"},
	}
	var buf bytes.Buffer
	log := Init(cfg, &buf)

	log.Info("test message", "key", "value")
	output := buf.String()
	if output == "" {
		t.Error("expected non-empty log output")
	}
	if !bytes.Contains(buf.Bytes(), []byte("test message")) {
		t.Errorf("expected log to contain 'test message', got %q", output)
	}
}

func TestInit_Text(t *testing.T) {
	cfg := &config.Config{
		Log: config.LogConfig{Level: "info", Format: "text", Output: "stdout"},
	}
	var buf bytes.Buffer
	log := Init(cfg, &buf)

	log.Info("hello")
	output := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte("hello")) {
		t.Errorf("expected log to contain 'hello', got %q", output)
	}
}

func TestLogger_With(t *testing.T) {
	cfg := &config.Config{
		Log: config.LogConfig{Level: "info", Format: "json"},
	}
	var buf bytes.Buffer
	log := Init(cfg, &buf)

	child := log.With("component", "auth")
	child.Info("auth started")

	output := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte("auth")) {
		t.Errorf("expected 'auth' in output, got %q", output)
	}
}

func TestLogger_WithGroup(t *testing.T) {
	cfg := &config.Config{
		Log: config.LogConfig{Level: "info", Format: "json"},
	}
	var buf bytes.Buffer
	log := Init(cfg, &buf)

	child := log.WithGroup("http")
	child.Info("request handled")

	output := buf.String()
	if output == "" {
		t.Error("expected non-empty output")
	}
}

func TestLogger_DebugLevel(t *testing.T) {
	cfg := &config.Config{
		Log: config.LogConfig{Level: "debug", Format: "json"},
	}
	var buf bytes.Buffer
	log := Init(cfg, &buf)

	log.Debug("debug message")
	if !bytes.Contains(buf.Bytes(), []byte("debug message")) {
		t.Errorf("expected debug message to appear, got %q", buf.String())
	}
}

func TestLogger_WarnError(t *testing.T) {
	cfg := &config.Config{
		Log: config.LogConfig{Level: "warn", Format: "json"},
	}
	var buf bytes.Buffer
	log := Init(cfg, &buf)

	log.Warn("warning")
	log.Error("error")

	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("warning")) {
		t.Error("expected warning in output")
	}
	if !bytes.Contains([]byte(output), []byte("error")) {
		t.Error("expected error in output")
	}
}

func TestParseLevel(t *testing.T) {
	if parseLevel("debug") != -4 { // slog.LevelDebug
		t.Error("expected debug level")
	}
	if parseLevel("warn") != 4 { // slog.LevelWarn
		t.Error("expected warn level")
	}
	if parseLevel("error") != 8 { // slog.LevelError
		t.Error("expected error level")
	}
	if parseLevel("") != 0 { // info
		t.Error("expected info level for empty")
	}
}

func TestInit_WithNilWriter(t *testing.T) {
	cfg := &config.Config{
		Log: config.LogConfig{Level: "error", Format: "json"},
	}
	// nil writer → defaults to os.Stdout
	log := Init(cfg, nil)
	if log == nil {
		t.Error("expected non-nil logger")
	}
	// Don't actually write to os.Stdout in tests, just check non-nil
}
