package eventbus_test

import (
	"context"
	"fmt"
	"sync"

	"github.com/i56/framework/core/eventbus"
	"github.com/i56/framework/core/logger"
)

type noopLogger struct{}

func (l noopLogger) Debug(msg string, args ...any)      {}
func (l noopLogger) Info(msg string, args ...any)       {}
func (l noopLogger) Warn(msg string, args ...any)       {}
func (l noopLogger) Error(msg string, args ...any)      {}
func (l noopLogger) With(args ...any) logger.Logger     { return l }
func (l noopLogger) WithGroup(name string) logger.Logger { return l }

var _ logger.Logger = noopLogger{}

// ExampleEventBus demonstrates publish/subscribe with the event bus.
func ExampleEventBus() {
	eb := eventbus.New(noopLogger{})
	ctx := context.Background()

	var mu sync.Mutex
	var received []string

	// Subscribe to an event
	eb.Subscribe("order.created", func(ctx context.Context, e eventbus.Event) error {
		mu.Lock()
		defer mu.Unlock()
		received = append(received, e.EventName())
		return nil
	}, false) // synchronous

	// Publish
	ev := eventbus.NewEvent("order.created")
	eb.Publish(ctx, ev)

	mu.Lock()
	fmt.Println("Events received:", len(received))
	mu.Unlock()
	// Output:
	// Events received: 1
}

// ExampleEventBus_WithTransport demonstrates external transport integration.
func ExampleEventBus_WithTransport() {
	transport := eventbus.NewMemTransport()
	eb := eventbus.New(noopLogger{}).WithTransport(transport)

	var externalReceived bool
	transport.Subscribe("order.created", func(ctx context.Context, e eventbus.Event) error {
		externalReceived = true
		return nil
	})

	eb.Publish(context.Background(), eventbus.NewEvent("order.created"))
	fmt.Println("Transport received:", externalReceived)
	// Output:
	// Transport received: true
}
