package scheduler

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestScheduler_AddAndListJobs(t *testing.T) {
	s := New()

	err := s.AddJob("cleanup", "@every 1h", func() error {
		return nil
	})
	if err != nil {
		t.Fatalf("AddJob: %v", err)
	}

	jobs := s.ListJobs()
	if len(jobs) != 1 {
		t.Errorf("expected 1 job, got %d", len(jobs))
	}
	if jobs[0].Name != "cleanup" {
		t.Errorf("expected job 'cleanup', got %q", jobs[0].Name)
	}
}

func TestScheduler_StartAndStop(t *testing.T) {
	s := New()

	var count atomic.Int32
	s.AddJob("tick", "@every 500ms", func() error {
		count.Add(1)
		return nil
	})

	s.Start()

	// Wait for at least one run
	time.Sleep(1200 * time.Millisecond)
	s.Stop()

	if count.Load() < 1 {
		t.Errorf("expected at least 1 run, got %d", count.Load())
	}
}

func TestScheduler_TriggerNow(t *testing.T) {
	s := New()

	var count atomic.Int32
	s.AddJob("manual", "@daily", func() error {
		count.Add(1)
		return nil
	})

	s.Start()
	defer s.Stop()

	err := s.TriggerNow("manual")
	if err != nil {
		t.Fatalf("TriggerNow: %v", err)
	}

	// Give it a moment to run
	time.Sleep(200 * time.Millisecond)

	if count.Load() != 1 {
		t.Errorf("expected 1 run, got %d", count.Load())
	}
}

func TestScheduler_TriggerNowNotFound(t *testing.T) {
	s := New()

	err := s.TriggerNow("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent job")
	}
}

func TestScheduler_DuplicateJob(t *testing.T) {
	s := New()

	err := s.AddJob("dup", "@daily", func() error { return nil })
	if err != nil {
		t.Fatalf("first AddJob: %v", err)
	}

	err = s.AddJob("dup", "@daily", func() error { return nil })
	if err == nil {
		t.Error("expected error for duplicate job")
	}
}
