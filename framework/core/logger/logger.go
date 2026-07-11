// Package logger provides structured logging based on Go 1.21+ log/slog.
package logger

import (
	"io"
	"log/slog"
	"os"

	"github.com/i56/framework/core/config"
)

// Logger is the application-wide logger interface.
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	With(args ...any) Logger
	WithGroup(name string) Logger
}

// SLogLogger wraps slog.Logger to implement the Logger interface.
type SLogLogger struct {
	inner *slog.Logger
}

// Init initializes the global logger based on config.
func Init(cfg *config.Config, w io.Writer) Logger {
	if w == nil {
		w = os.Stdout
	}

	level := parseLevel(cfg.Log.Level)
	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if cfg.Log.Format == "text" {
		handler = slog.NewTextHandler(w, opts)
	} else {
		handler = slog.NewJSONHandler(w, opts)
	}

	return &SLogLogger{inner: slog.New(handler)}
}

func parseLevel(s string) slog.Level {
	switch s {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func (l *SLogLogger) Debug(msg string, args ...any) { l.inner.Debug(msg, args...) }
func (l *SLogLogger) Info(msg string, args ...any)  { l.inner.Info(msg, args...) }
func (l *SLogLogger) Warn(msg string, args ...any)  { l.inner.Warn(msg, args...) }
func (l *SLogLogger) Error(msg string, args ...any) { l.inner.Error(msg, args...) }
func (l *SLogLogger) With(args ...any) Logger {
	return &SLogLogger{inner: l.inner.With(args...)}
}
func (l *SLogLogger) WithGroup(name string) Logger {
	return &SLogLogger{inner: l.inner.WithGroup(name)}
}
