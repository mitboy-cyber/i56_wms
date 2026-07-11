// Package scheduler provides cron-like task scheduling.
package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/i56/framework/core/logger"
)

// Job represents a scheduled task.
type Job struct {
	Name     string
	Schedule string // cron expression or "@every 5m"
	Handler  func(ctx context.Context) error
}

// Scheduler manages and runs scheduled jobs.
type Scheduler struct {
	mu    sync.RWMutex
	jobs  map[string]*Job
	log   logger.Logger
	stopCh chan struct{}
}

// New creates a new Scheduler.
func New(log logger.Logger) *Scheduler {
	return &Scheduler{
		jobs:   make(map[string]*Job),
		log:    log,
		stopCh: make(chan struct{}),
	}
}

// AddJob registers a new job.
func (s *Scheduler) AddJob(job *Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[job.Name] = job
	s.log.Info("job registered", "name", job.Name, "schedule", job.Schedule)
	return nil
}

// RemoveJob removes a job by name.
func (s *Scheduler) RemoveJob(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.jobs, name)
	s.log.Info("job removed", "name", name)
}

// Start begins the scheduler loop.
func (s *Scheduler) Start(ctx context.Context) error {
	s.log.Info("scheduler starting", "jobs", len(s.jobs))

	ticker := time.NewTicker(30 * time.Second) // Check every 30s
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-s.stopCh:
			return nil
		case <-ticker.C:
			s.runDueJobs(ctx)
		}
	}
}

// Stop gracefully stops the scheduler.
func (s *Scheduler) Stop() error {
	close(s.stopCh)
	s.log.Info("scheduler stopped")
	return nil
}

func (s *Scheduler) runDueJobs(ctx context.Context) {
	s.mu.RLock()
	jobs := make([]*Job, 0, len(s.jobs))
	for _, j := range s.jobs {
		jobs = append(jobs, j)
	}
	s.mu.RUnlock()

	for _, job := range jobs {
		// TODO: implement proper cron expression parsing
		// For now, run all jobs on each tick (prototype)
		go func(j *Job) {
			s.log.Debug("running job", "name", j.Name)
			if err := j.Handler(ctx); err != nil {
				s.log.Error("job failed", "name", j.Name, "error", err)
			}
		}(job)
	}
}
