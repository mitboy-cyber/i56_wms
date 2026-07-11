package eventbus

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	"github.com/i56/framework/core/logger"
)

func TestEventBus_PublishSync(t *testing.T) {
	log := testLogger{}
	eb := New(log)

	var callCount int32
	handler := func(ctx context.Context, e Event) error {
		atomic.AddInt32(&callCount, 1)
		return nil
	}

	eb.Subscribe("order.created", handler, false)

	err := eb.Publish(context.Background(), NewEvent("order.created"))
	if err != nil {
		t.Fatalf("Publish: %v", err)
	}

	if atomic.LoadInt32(&callCount) != 1 {
		t.Errorf("expected 1 handler call, got %d", callCount)
	}
}

func TestEventBus_PublishAsync(t *testing.T) {
	log := testLogger{}
	eb := New(log)

	ch := make(chan struct{})
	handler := func(ctx context.Context, e Event) error {
		ch <- struct{}{}
		return nil
	}

	eb.Subscribe("order.created", handler, true)

	_ = eb.Publish(context.Background(), NewEvent("order.created"))

	// Wait for async handler
	<-ch
}

func TestEventBus_PublishSyncWaitsForAll(t *testing.T) {
	log := testLogger{}
	eb := New(log)

	var callCount int32
	handler := func(ctx context.Context, e Event) error {
		atomic.AddInt32(&callCount, 1)
		return nil
	}

	eb.Subscribe("order.created", handler, false)
	eb.Subscribe("order.created", handler, false)

	err := eb.PublishSync(context.Background(), NewEvent("order.created"))
	if err != nil {
		t.Fatalf("PublishSync: %v", err)
	}

	if atomic.LoadInt32(&callCount) != 2 {
		t.Errorf("expected 2 handler calls, got %d", callCount)
	}
}

func TestEventBus_ErrorReturns(t *testing.T) {
	log := testLogger{}
	eb := New(log)

	expectedErr := errors.New("handler failed")
	handler := func(ctx context.Context, e Event) error {
		return expectedErr
	}

	eb.Subscribe("order.created", handler, false)

	err := eb.Publish(context.Background(), NewEvent("order.created"))
	if err == nil {
		t.Error("expected error from handler")
	}
}

func TestEventBus_HandlerCount(t *testing.T) {
	log := testLogger{}
	eb := New(log)

	handler := func(ctx context.Context, e Event) error { return nil }

	eb.Subscribe("e1", handler, false)
	eb.Subscribe("e1", handler, true)
	eb.Subscribe("e2", handler, false)

	if count := eb.HandlerCount("e1"); count != 2 {
		t.Errorf("expected 2 handlers for e1, got %d", count)
	}
	if count := eb.HandlerCount("e2"); count != 1 {
		t.Errorf("expected 1 handler for e2, got %d", count)
	}
	if count := eb.HandlerCount("nonexistent"); count != 0 {
		t.Errorf("expected 0 handlers for nonexistent, got %d", count)
	}
}

func TestEventBus_SubscribedEvents(t *testing.T) {
	log := testLogger{}
	eb := New(log)
	handler := func(ctx context.Context, e Event) error { return nil }

	eb.Subscribe("a", handler, false)
	eb.Subscribe("b", handler, false)

	events := eb.SubscribedEvents()
	if len(events) != 2 {
		t.Errorf("expected 2 subscribed events, got %d", len(events))
	}
}

func TestMemTransport_PublishSubscribe(t *testing.T) {
	transport := NewMemTransport()

	var received Event
	handler := func(ctx context.Context, e Event) error {
		received = e
		return nil
	}

	transport.Subscribe("order.created", handler)
	ev := NewEvent("order.created")

	err := transport.Publish(context.Background(), "order.created", ev)
	if err != nil {
		t.Fatalf("Publish: %v", err)
	}

	if received == nil {
		t.Error("expected event to be received")
	}
	if received.EventName() != "order.created" {
		t.Errorf("expected 'order.created', got %q", received.EventName())
	}
}

func TestEventBus_WithTransport(t *testing.T) {
	log := testLogger{}
	transport := NewMemTransport()
	eb := New(log).WithTransport(transport)

	var transportReceived atomic.Bool
	transport.Subscribe("order.created", func(ctx context.Context, e Event) error {
		transportReceived.Store(true)
		return nil
	})

	err := eb.Publish(context.Background(), NewEvent("order.created"))
	if err != nil {
		t.Fatalf("Publish: %v", err)
	}

	if !transportReceived.Load() {
		t.Error("expected transport to receive event")
	}
}

func TestEventBus_Close(t *testing.T) {
	log := testLogger{}
	transport := NewMemTransport()
	eb := New(log).WithTransport(transport)

	if err := eb.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	// Without transport
	eb2 := New(log)
	if err := eb2.Close(); err != nil {
		t.Fatalf("Close without transport: %v", err)
	}
}

func TestBaseEvent(t *testing.T) {
	ev := NewEvent("test.event")
	if ev.EventName() != "test.event" {
		t.Errorf("expected 'test.event', got %q", ev.EventName())
	}
	if ev.OccurredAt().IsZero() {
		t.Error("expected non-zero time")
	}
}

// testLogger implements logger.Logger for testing (no-op).
type testLogger struct{}

func (l testLogger) Debug(msg string, args ...any)  {}
func (l testLogger) Info(msg string, args ...any)   {}
func (l testLogger) Warn(msg string, args ...any)   {}
func (l testLogger) Error(msg string, args ...any)  {}
func (l testLogger) With(args ...any) logger.Logger { return l }
func (l testLogger) WithGroup(name string) logger.Logger { return l }

var _ logger.Logger = testLogger{}
