// Package eventbus provides an in-process pub/sub event bus with a pluggable
// external transport interface (Kafka, NATS, Redis Streams, RabbitMQ).
package eventbus

import (
	"context"
	"sync"
	"time"

	"github.com/i56/framework/core/logger"
)

// ---------------------------------------------------------------------------
// Domain Events
// ---------------------------------------------------------------------------

// Event is the interface that all domain events must implement.
type Event interface {
	EventName() string
	OccurredAt() time.Time
}

// BaseEvent provides default Event implementation.
type BaseEvent struct {
	Name string    `json:"name"`
	Time time.Time `json:"occurred_at"`
}

func (e BaseEvent) EventName() string     { return e.Name }
func (e BaseEvent) OccurredAt() time.Time { return e.Time }

// NewEvent creates a new BaseEvent.
func NewEvent(name string) BaseEvent {
	return BaseEvent{Name: name, Time: time.Now()}
}

// ---------------------------------------------------------------------------
// Handler
// ---------------------------------------------------------------------------

// EventHandler is a function that handles events.
type EventHandler func(ctx context.Context, event Event) error

// ---------------------------------------------------------------------------
// Transport (pluggable external transport)
// ---------------------------------------------------------------------------

// Transport abstracts an external messaging system (Kafka, NATS, Redis, RabbitMQ).
// Plug in a real transport to publish events out-of-process.
type Transport interface {
	// Publish sends an event to the external transport.
	Publish(ctx context.Context, topic string, event Event) error
	// Subscribe registers an event handler for a topic.
	Subscribe(topic string, handler EventHandler) error
	// Close shuts down the transport gracefully.
	Close() error
}

// ---------------------------------------------------------------------------
// In-memory EventBus
// ---------------------------------------------------------------------------

// EventBus manages event subscriptions and publishing in-process.
type EventBus struct {
	mu            sync.RWMutex
	syncHandlers  map[string][]EventHandler
	asyncHandlers map[string][]EventHandler
	transport     Transport
	log           logger.Logger
}

// New creates a new EventBus.
func New(log logger.Logger) *EventBus {
	return &EventBus{
		syncHandlers:  make(map[string][]EventHandler),
		asyncHandlers: make(map[string][]EventHandler),
		log:           log,
	}
}

// WithTransport attaches an external transport to the event bus.
// When set, Publish/PublishSync will also forward to the transport.
func (eb *EventBus) WithTransport(t Transport) *EventBus {
	eb.transport = t
	return eb
}

// Subscribe registers a handler for an event name.
// If async is true, the handler runs in a goroutine (fire-and-forget).
func (eb *EventBus) Subscribe(eventName string, handler EventHandler, async bool) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if async {
		eb.asyncHandlers[eventName] = append(eb.asyncHandlers[eventName], handler)
	} else {
		eb.syncHandlers[eventName] = append(eb.syncHandlers[eventName], handler)
	}
}

// Publish fires an event to all registered handlers.
// Sync handlers run in order; async handlers run in goroutines.
// If an external transport is set, the event is also published there.
func (eb *EventBus) Publish(ctx context.Context, event Event) error {
	name := event.EventName()

	// Forward to external transport
	if eb.transport != nil {
		if err := eb.transport.Publish(ctx, name, event); err != nil {
			eb.log.Error("eventbus: transport publish failed", "event", name, "error", err)
		}
	}

	eb.mu.RLock()
	asyncCopy := eb.asyncHandlers[name]
	syncCopy := eb.syncHandlers[name]
	eb.mu.RUnlock()

	// Run async handlers
	for _, handler := range asyncCopy {
		go func(h EventHandler) {
			if err := h(ctx, event); err != nil {
				eb.log.Error("async event handler failed",
					"event", name,
					"error", err,
				)
			}
		}(handler)
	}

	// Run sync handlers
	for _, handler := range syncCopy {
		if err := handler(ctx, event); err != nil {
			eb.log.Error("sync event handler failed",
				"event", name,
				"error", err,
			)
			return err
		}
	}

	return nil
}

// PublishSync fires an event and waits for all handlers (sync + async).
// If an external transport is set, the event is also published there.
func (eb *EventBus) PublishSync(ctx context.Context, event Event) error {
	name := event.EventName()

	// Forward to external transport (blocking)
	if eb.transport != nil {
		if err := eb.transport.Publish(ctx, name, event); err != nil {
			eb.log.Error("eventbus: transport publish sync failed", "event", name, "error", err)
		}
	}

	eb.mu.RLock()
	asyncCopy := eb.asyncHandlers[name]
	syncCopy := eb.syncHandlers[name]
	eb.mu.RUnlock()

	var wg sync.WaitGroup
	errCh := make(chan error, len(asyncCopy)+len(syncCopy))

	// Run all handlers concurrently
	all := append(syncCopy, asyncCopy...)
	for _, handler := range all {
		wg.Add(1)
		go func(h EventHandler) {
			defer wg.Done()
			if err := h(ctx, event); err != nil {
				errCh <- err
			}
		}(handler)
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return err
		}
	}
	return nil
}

// HandlerCount returns the total number of subscribed handlers for an event.
func (eb *EventBus) HandlerCount(eventName string) int {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	return len(eb.syncHandlers[eventName]) + len(eb.asyncHandlers[eventName])
}

// SubscribedEvents returns all event names with registered handlers.
func (eb *EventBus) SubscribedEvents() []string {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	seen := make(map[string]bool)
	for k := range eb.syncHandlers {
		seen[k] = true
	}
	for k := range eb.asyncHandlers {
		seen[k] = true
	}
	names := make([]string, 0, len(seen))
	for k := range seen {
		names = append(names, k)
	}
	return names
}

// Close shuts down the external transport if one is set.
func (eb *EventBus) Close() error {
	if eb.transport != nil {
		return eb.transport.Close()
	}
	return nil
}

// ---------------------------------------------------------------------------
// In-memory Transport (for testing)
// ---------------------------------------------------------------------------

// MemTransport is an in-memory transport that implements Transport.
// Useful for testing event routing without a real Kafka/NATS.
type MemTransport struct {
	mu       sync.RWMutex
	handlers map[string][]EventHandler
}

// NewMemTransport creates an in-memory transport.
func NewMemTransport() *MemTransport {
	return &MemTransport{
		handlers: make(map[string][]EventHandler),
	}
}

func (t *MemTransport) Publish(ctx context.Context, topic string, event Event) error {
	t.mu.RLock()
	handlers := t.handlers[topic]
	t.mu.RUnlock()

	for _, h := range handlers {
		if err := h(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

func (t *MemTransport) Subscribe(topic string, handler EventHandler) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.handlers[topic] = append(t.handlers[topic], handler)
	return nil
}

func (t *MemTransport) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.handlers = nil
	return nil
}
