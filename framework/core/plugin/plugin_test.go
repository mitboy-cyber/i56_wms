package plugin

import (
	"context"
	"testing"
)

type testLogger struct{}

func (l testLogger) Debug(msg string, args ...any) {}
func (l testLogger) Info(msg string, args ...any)  {}
func (l testLogger) Warn(msg string, args ...any)  {}
func (l testLogger) Error(msg string, args ...any) {}
func (l testLogger) With(args ...any) logger       { return l }
func (l testLogger) WithGroup(name string) logger  { return l }

// duplicated logger interface (minimal, avoids circular deps)
type logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	With(args ...any) logger
	WithGroup(name string) logger
}

func TestRegistry_Start(t *testing.T) {
	reg := NewRegistry(testLogger{})

	initCalled := false
	reg.Register(NewBasePlugin("test", "1.0", func(ctx context.Context, reg *Registry) error {
		initCalled = true
		reg.Provide("test_svc", "hello")
		return nil
	}), PriorityCore)

	if reg.Status() != StatusUninitialized {
		t.Error("expected uninitialized status")
	}

	err := reg.Start(context.Background())
	if err != nil {
		t.Fatalf("Start: %v", err)
	}

	if !initCalled {
		t.Error("expected init to be called")
	}
	if reg.Status() != StatusRunning {
		t.Errorf("expected running status, got %d", reg.Status())
	}
	if svc := reg.Resolve("test_svc"); svc != "hello" {
		t.Errorf("expected 'hello', got %v", svc)
	}
}

func TestRegistry_StartOrder(t *testing.T) {
	reg := NewRegistry(testLogger{})

	var order []string
	reg.Register(NewBasePlugin("last", "1.0", func(ctx context.Context, reg *Registry) error {
		order = append(order, "last")
		return nil
	}), PriorityAPI)
	reg.Register(NewBasePlugin("first", "1.0", func(ctx context.Context, reg *Registry) error {
		order = append(order, "first")
		return nil
	}), PrioritySystem)

	_ = reg.Start(context.Background())

	if len(order) != 2 {
		t.Fatalf("expected 2 plugin inits, got %d", len(order))
	}
	if order[0] != "first" {
		t.Errorf("expected 'first' first, got %q", order[0])
	}
	if order[1] != "last" {
		t.Errorf("expected 'last' second, got %q", order[1])
	}
}

func TestRegistry_AlreadyStarted(t *testing.T) {
	reg := NewRegistry(testLogger{})
	_ = reg.Start(context.Background())

	err := reg.Start(context.Background())
	if err == nil {
		t.Error("expected error for already started registry")
	}
}

func TestRegistry_ResolveAndMustResolve(t *testing.T) {
	reg := NewRegistry(testLogger{})
	reg.Provide("svc", 42)

	if reg.Resolve("nonexistent") != nil {
		t.Error("expected nil for nonexistent service")
	}
	if reg.Resolve("svc") != 42 {
		t.Error("expected 42")
	}

	// MustResolve panics on missing
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	reg.MustResolve("nonexistent")
}

func TestRegistry_ListAndCount(t *testing.T) {
	reg := NewRegistry(testLogger{})
	reg.Register(NewBasePlugin("a", "1.0", nil), PriorityCore)
	reg.Register(NewBasePlugin("b", "1.0", nil), PriorityCore)

	if reg.PluginCount() != 2 {
		t.Errorf("expected 2 plugins, got %d", reg.PluginCount())
	}

	names := reg.List()
	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}
}

func TestBasePlugin(t *testing.T) {
	p := NewBasePlugin("my-plugin", "2.0.1", nil)
	if p.Name() != "my-plugin" {
		t.Errorf("expected 'my-plugin', got %q", p.Name())
	}
	if p.Version() != "2.0.1" {
		t.Errorf("expected '2.0.1', got %q", p.Version())
	}
	err := p.Init(context.Background(), nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
