// Package queue provides message queue abstraction with in-memory and pluggable backends.
// Supports: in-memory (sync/async), Redis Streams, NATS, Kafka.
package queue

import (
	"context"
	"sync"
	"time"

	"github.com/i56/framework/core/logger"
)

// Message represents a queued item with metadata.
type Message struct {
	ID        string            `json:"id"`
	Payload   []byte            `json:"payload"`
	Headers   map[string]string `json:"headers,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
	Attempts  int               `json:"attempts"`
}

// Handler is the callback for processing queued messages.
type Handler func(ctx context.Context, msg *Message) error

// Queue is the abstract message queue interface.
type Queue interface {
	// Enqueue adds a message to the queue.
	Enqueue(ctx context.Context, msg *Message) error
	// Dequeue retrieves and removes the next message (blocking with timeout).
	Dequeue(ctx context.Context, timeout time.Duration) (*Message, error)
	// Ack acknowledges successful processing.
	Ack(ctx context.Context, msgID string) error
	// Nack negatively acknowledges (requeue or dead-letter).
	Nack(ctx context.Context, msgID string) error
	// Len returns the approximate queue length.
	Len(ctx context.Context) (int64, error)
	// Close shuts down the queue gracefully.
	Close() error
}

// Consumer processes messages from a queue using a handler.
type Consumer struct {
	queue       Queue
	handler     Handler
	concurrency int
	log         logger.Logger
	stopCh      chan struct{}
	wg          sync.WaitGroup
}

// NewConsumer creates a message consumer.
func NewConsumer(queue Queue, handler Handler, concurrency int, log logger.Logger) *Consumer {
	if concurrency < 1 {
		concurrency = 1
	}
	return &Consumer{
		queue:       queue,
		handler:     handler,
		concurrency: concurrency,
		log:         log,
		stopCh:      make(chan struct{}),
	}
}

// Start begins consuming messages on N goroutines.
func (c *Consumer) Start(ctx context.Context) {
	for i := 0; i < c.concurrency; i++ {
		c.wg.Add(1)
		go func(workerID int) {
			defer c.wg.Done()
			c.consume(ctx, workerID)
		}(i)
	}
}

// Stop signals the consumer to stop and waits.
func (c *Consumer) Stop() {
	close(c.stopCh)
	c.wg.Wait()
}

func (c *Consumer) consume(ctx context.Context, workerID int) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-c.stopCh:
			return
		default:
		}

		msg, err := c.queue.Dequeue(ctx, 5*time.Second)
		if err != nil {
			if err == context.Canceled || err == context.DeadlineExceeded {
				continue
			}
			c.log.Error("queue: dequeue error", "worker", workerID, "error", err)
			time.Sleep(time.Second)
			continue
		}

		if err := c.handler(ctx, msg); err != nil {
			c.log.Error("queue: handler error", "worker", workerID, "msg_id", msg.ID, "error", err)
			_ = c.queue.Nack(ctx, msg.ID)
		} else {
			_ = c.queue.Ack(ctx, msg.ID)
		}
	}
}

// MemQueue is an in-memory queue implementation (FIFO, unbounded).
type MemQueue struct {
	mu       sync.Mutex
	cond     *sync.Cond
	messages []*Message
	acked    map[string]bool
	nacked   map[string]bool
	closed   bool
}

// NewMemQueue creates an in-memory queue.
func NewMemQueue() *MemQueue {
	mq := &MemQueue{
		messages: make([]*Message, 0),
		acked:    make(map[string]bool),
		nacked:   make(map[string]bool),
	}
	mq.cond = sync.NewCond(&mq.mu)
	return mq
}

func (q *MemQueue) Enqueue(ctx context.Context, msg *Message) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return ErrQueueClosed
	}
	if msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}
	q.messages = append(q.messages, msg)
	q.cond.Signal()
	return nil
}

func (q *MemQueue) Dequeue(ctx context.Context, timeout time.Duration) (*Message, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	deadline := time.Now().Add(timeout)
	for len(q.messages) == 0 && !q.closed {
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return nil, context.DeadlineExceeded
		}
		// Use a timer-based wait to respect the deadline
		done := make(chan struct{})
		go func() {
			q.cond.Wait()
			close(done)
		}()

		q.mu.Unlock()
		select {
		case <-done:
		case <-time.After(remaining):
			q.cond.Signal() // wake up the goroutine waiting on cond
		case <-ctx.Done():
			q.mu.Lock()
			return nil, ctx.Err()
		}
		q.mu.Lock()
	}

	if q.closed && len(q.messages) == 0 {
		return nil, ErrQueueClosed
	}

	msg := q.messages[0]
	q.messages = q.messages[1:]
	if msg.Attempts == 0 {
		msg.Attempts = 1
	}
	return msg, nil
}

func (q *MemQueue) Ack(ctx context.Context, msgID string) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.acked[msgID] = true
	return nil
}

func (q *MemQueue) Nack(ctx context.Context, msgID string) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.nacked[msgID] = true
	return nil
}

func (q *MemQueue) Len(ctx context.Context) (int64, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	return int64(len(q.messages)), nil
}

func (q *MemQueue) Close() error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.closed = true
	q.cond.Broadcast()
	return nil
}

// ErrQueueClosed indicates the queue has been closed.
var ErrQueueClosed = &queueError{"queue is closed"}

type queueError struct{ msg string }

func (e *queueError) Error() string { return e.msg }
