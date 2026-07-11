package scheduler

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/i56/framework/core/logger"
)

// testLogger implements logger.Logger.
type testLogger struct{}

func (l testLogger) Debug(msg string, args ...any)  {}
func (l testLogger) Info(msg string, args ...any)   {}
func (l testLogger) Warn(msg string, args ...any)   {}
func (l testLogger) Error(msg string, args ...any)  {}
func (l testLogger) With(args ...any) logger.Logger { return l }
func (l testLogger) WithGroup(name string) logger.Logger { return l }

var _ logger.Logger = testLogger{}

func TestScheduler_AddAndRemoveJob(t *testing.T) {
	s := New(testLogger{})

	err := s.AddJob(&Job{Name: "cleanup", Schedule: "@every 1h", Handler: func(ctx context.Context) error {
		return nil
	}})
	if err != nil {
		t.Fatalf("AddJob: %v", err)
	}

	if len(s.jobs) != 1 {
		t.Errorf("expected 1 job, got %d", len(s.jobs))
	}

	s.RemoveJob("cleanup")
	if len(s.jobs) != 0 {
		t.Errorf("expected 0 jobs after remove, got %d", len(s.jobs))
	}
}

func TestScheduler_StartAndStop(t *testing.T) {
	s := New(testLogger{})

	var count atomic.Int32
	s.AddJob(&Job{Name: "tick", Schedule: "@every 1s", Handler: func(ctx context.Context) error {
		count.Add(1)
		return nil
	}})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	go func() {
		s.Start(ctx)
	}()

	// Wait briefly then stop
	time.Sleep(500 * time.Millisecond)
	s.Stop()

	// At least one run should have happened
	if count.Load() < 1 {
		t.Errorf("expected at least 1 run, got %d", count.Load())
	}
}

func TestScheduler_ContextCancellation(t *testing.T) {
	s := New(testLogger{})
	s.AddJob(&Job{Name: "cancelled", Schedule: "@every 1s", Handler: func(ctx context.Context) error {
		return nil
	}})

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() {
		done <- s.Start(ctx)
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Error("timed out waiting for scheduler to stop")
	}
}
