package queue

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/i56/framework/core/logger"
)

type testLogger struct{}

func (l testLogger) Debug(msg string, args ...any)      {}
func (l testLogger) Info(msg string, args ...any)       {}
func (l testLogger) Warn(msg string, args ...any)       {}
func (l testLogger) Error(msg string, args ...any)      {}
func (l testLogger) With(args ...any) logger.Logger     { return l }
func (l testLogger) WithGroup(name string) logger.Logger { return l }

var _ logger.Logger = testLogger{}

func TestMemQueue_EnqueueDequeue(t *testing.T) {
	q := NewMemQueue()

	msg := &Message{ID: "1", Payload: []byte("hello")}
	if err := q.Enqueue(context.Background(), msg); err != nil {
		t.Fatalf("Enqueue: %v", err)
	}

	len, err := q.Len(context.Background())
	if err != nil {
		t.Fatalf("Len: %v", err)
	}
	if len != 1 {
		t.Errorf("expected length 1, got %d", len)
	}

	out, err := q.Dequeue(context.Background(), time.Second)
	if err != nil {
		t.Fatalf("Dequeue: %v", err)
	}
	if out.ID != "1" {
		t.Errorf("expected id '1', got %q", out.ID)
	}
	if string(out.Payload) != "hello" {
		t.Errorf("expected 'hello', got %q", out.Payload)
	}

	len, _ = q.Len(context.Background())
	if len != 0 {
		t.Errorf("expected length 0 after dequeue, got %d", len)
	}
}

func TestMemQueue_DequeueTimeout(t *testing.T) {
	q := NewMemQueue()

	_, err := q.Dequeue(context.Background(), 50*time.Millisecond)
	if err != context.DeadlineExceeded {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
}

func TestMemQueue_AckNack(t *testing.T) {
	q := NewMemQueue()
	ctx := context.Background()

	if err := q.Ack(ctx, "msg-1"); err != nil {
		t.Errorf("Ack: %v", err)
	}
	if err := q.Nack(ctx, "msg-1"); err != nil {
		t.Errorf("Nack: %v", err)
	}
}

func TestMemQueue_Close(t *testing.T) {
	q := NewMemQueue()

	if err := q.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	_, err := q.Dequeue(context.Background(), time.Second)
	if err != ErrQueueClosed {
		t.Errorf("expected ErrQueueClosed, got %v", err)
	}

	err = q.Enqueue(context.Background(), &Message{ID: "1"})
	if err != ErrQueueClosed {
		t.Errorf("expected ErrQueueClosed on enqueue, got %v", err)
	}
}

func TestMemQueue_MultipleMessages(t *testing.T) {
	q := NewMemQueue()
	ctx := context.Background()

	for i := 0; i < 100; i++ {
		if err := q.Enqueue(ctx, &Message{ID: string(rune('a' + i%26)), Payload: []byte{byte(i)}}); err != nil {
			t.Fatalf("Enqueue %d: %v", i, err)
		}
	}

	for i := 0; i < 100; i++ {
		if _, err := q.Dequeue(ctx, time.Second); err != nil {
			t.Fatalf("Dequeue %d: %v", i, err)
		}
	}

	len, _ := q.Len(ctx)
	if len != 0 {
		t.Errorf("expected empty queue, got %d", len)
	}
}

func TestConsumer_StartStop(t *testing.T) {
	q := NewMemQueue()
	var processed atomic.Int32

	consumer := NewConsumer(q, func(ctx context.Context, msg *Message) error {
		processed.Add(1)
		return nil
	}, 2, testLogger{})

	ctx, cancel := context.WithCancel(context.Background())
	consumer.Start(ctx)

	// Enqueue some messages
	q.Enqueue(context.Background(), &Message{ID: "1"})
	q.Enqueue(context.Background(), &Message{ID: "2"})

	// Wait for processing
	time.Sleep(200 * time.Millisecond)
	consumer.Stop()
	cancel()

	if processed.Load() < 1 {
		t.Errorf("expected at least 1 processed message, got %d", processed.Load())
	}
}

func TestMessage_TimestampDefault(t *testing.T) {
	q := NewMemQueue()
	msg := &Message{ID: "test", Payload: []byte("data")}
	q.Enqueue(context.Background(), msg)

	out, _ := q.Dequeue(context.Background(), time.Second)
	if out.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}
