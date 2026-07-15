// Package scheduler provides cron-like task scheduling.
package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// JobStatus represents the current state of a scheduled job.
type JobStatus struct {
	Name      string    `json:"name"`
	CronExpr  string    `json:"cron_expr"`
	LastRun   time.Time `json:"last_run"`
	NextRun   time.Time `json:"next_run"`
	LastError string    `json:"last_error,omitempty"`
	RunCount  int64     `json:"run_count"`
	Running   bool      `json:"running"`
}

// Job represents a scheduled task.
type Job struct {
	Name     string
	CronExpr string // cron expression or "@every 5m"
	Handler  func() error
	LastRun  time.Time
	NextRun  time.Time
	LastErr  string
	RunCount int64
	running  bool
	mu       sync.Mutex
}

// Scheduler manages and runs scheduled jobs.
type Scheduler struct {
	mu       sync.RWMutex
	jobs     map[string]*Job
	stopCh   chan struct{}
	running  bool
}

// New creates a new Scheduler.
func New() *Scheduler {
	return &Scheduler{
		jobs:   make(map[string]*Job),
		stopCh: make(chan struct{}),
	}
}

// AddJob registers a new job with a cron expression.
func (s *Scheduler) AddJob(name, cronExpr string, handler func() error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.jobs[name]; exists {
		return fmt.Errorf("scheduler: job %q already registered", name)
	}

	next, err := nextRun(cronExpr, time.Now())
	if err != nil {
		return fmt.Errorf("scheduler: invalid cron expression %q: %w", cronExpr, err)
	}

	s.jobs[name] = &Job{
		Name:     name,
		CronExpr: cronExpr,
		Handler:  handler,
		NextRun:  next,
	}
	return nil
}

// Start begins the scheduler loop.
func (s *Scheduler) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	go s.loop()
}

// Stop gracefully stops the scheduler.
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.running {
		close(s.stopCh)
		s.running = false
	}
}

// ListJobs returns the status of all registered jobs.
func (s *Scheduler) ListJobs() []JobStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]JobStatus, 0, len(s.jobs))
	for _, j := range s.jobs {
		j.mu.Lock()
		result = append(result, JobStatus{
			Name:      j.Name,
			CronExpr:  j.CronExpr,
			LastRun:   j.LastRun,
			NextRun:   j.NextRun,
			LastError: j.LastErr,
			RunCount:  j.RunCount,
			Running:   j.running,
		})
		j.mu.Unlock()
	}
	return result
}

// TriggerNow triggers a job immediately by name.
func (s *Scheduler) TriggerNow(name string) error {
	s.mu.RLock()
	job, ok := s.jobs[name]
	s.mu.RUnlock()
	if !ok {
		return fmt.Errorf("scheduler: job %q not found", name)
	}

	go s.runJob(job)
	return nil
}

// loop is the main scheduling loop.
func (s *Scheduler) loop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.runDueJobs()
		}
	}
}

func (s *Scheduler) runDueJobs() {
	now := time.Now()

	s.mu.RLock()
	jobs := make([]*Job, 0, len(s.jobs))
	for _, j := range s.jobs {
		jobs = append(jobs, j)
	}
	s.mu.RUnlock()

	for _, job := range jobs {
		job.mu.Lock()
		due := !job.NextRun.IsZero() && !now.Before(job.NextRun)
		job.mu.Unlock()
		if due {
			go s.runJob(job)
		}
	}
}

func (s *Scheduler) runJob(job *Job) {
	job.mu.Lock()
	job.running = true
	job.mu.Unlock()

	start := time.Now()
	err := job.Handler()

	job.mu.Lock()
	job.LastRun = start
	job.RunCount++
	job.running = false
	if err != nil {
		job.LastErr = err.Error()
	} else {
		job.LastErr = ""
	}
	next, nerr := nextRun(job.CronExpr, time.Now())
	if nerr == nil {
		job.NextRun = next
	}
	job.mu.Unlock()
}

// nextRun calculates the next run time from a cron-like expression.
func nextRun(expr string, now time.Time) (time.Time, error) {
	// Handle "@every X" format
	if len(expr) > 7 && expr[:7] == "@every " {
		d, err := time.ParseDuration(expr[7:])
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid duration: %s", expr[7:])
		}
		return now.Add(d), nil
	}

	// Handle "@daily" and "@hourly"
	switch expr {
	case "@daily", "@midnight":
		y, m, d := now.Date()
		return time.Date(y, m, d+1, 0, 0, 0, 0, now.Location()), nil
	case "@hourly":
		return now.Truncate(time.Hour).Add(time.Hour), nil
	case "@every 5m":
		return now.Add(5 * time.Minute), nil
	case "@every 1h":
		return now.Add(time.Hour), nil
	case "@every 6h":
		return now.Add(6 * time.Hour), nil
	}

	// Handle standard 5-field cron: "min hour dom mon dow"
	fields := parseCronFields(expr)
	if len(fields) != 5 {
		return time.Time{}, fmt.Errorf("expected 5 cron fields, got %d: %q", len(fields), expr)
	}

	minute := parseCronField(fields[0], 0, 59)
	hour := parseCronField(fields[1], 0, 23)
	dayOfMonth := parseCronField(fields[2], 1, 31)
	month := parseCronField(fields[3], 1, 12)
	// dayOfWeek is field[4], we skip for simplicity

	// Find next matching time
	candidate := now.Truncate(time.Minute).Add(time.Minute)
	for i := 0; i < 525600; i++ { // search up to 1 year
		if month != nil && !month[int(candidate.Month())] {
			candidate = time.Date(candidate.Year(), candidate.Month()+1, 1, 0, 0, 0, 0, candidate.Location())
			continue
		}
		if dayOfMonth != nil && !dayOfMonth[candidate.Day()] {
			candidate = time.Date(candidate.Year(), candidate.Month(), candidate.Day()+1, 0, 0, 0, 0, candidate.Location())
			continue
		}
		if hour != nil && !hour[candidate.Hour()] {
			candidate = time.Date(candidate.Year(), candidate.Month(), candidate.Day(), candidate.Hour()+1, 0, 0, 0, candidate.Location())
			continue
		}
		if minute != nil && !minute[candidate.Minute()] {
			candidate = candidate.Truncate(time.Hour).Add(time.Hour)
			continue
		}
		return candidate, nil
	}
	return time.Time{}, fmt.Errorf("no matching time found for cron %q within 1 year", expr)
}

func parseCronFields(expr string) []string {
	fields := make([]string, 0, 5)
	current := ""
	for _, ch := range expr {
		if ch == ' ' || ch == '\t' {
			if current != "" {
				fields = append(fields, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		fields = append(fields, current)
	}
	return fields
}

func parseCronField(field string, min, max int) map[int]bool {
	if field == "*" {
		return nil // wildcard — match all
	}

	result := make(map[int]bool)
	// Handle comma-separated values: "1,3,5"
	parts := splitCron(field, ',')
	for _, part := range parts {
		result[min] = false // placeholder, will overwrite
		_ = result
		// Handle range: "1-5"
		if dashIdx := indexOf(part, '-'); dashIdx >= 0 {
			start := atoi(part[:dashIdx])
			end := atoi(part[dashIdx+1:])
			for v := start; v <= end; v++ {
				if v >= min && v <= max {
					result[v] = true
				}
			}
			continue
		}
		// Handle step: "*/5"
		if len(part) >= 2 && part[:2] == "*/" {
			step := atoi(part[2:])
			for v := min; v <= max; v += step {
				result[v] = true
			}
			continue
		}
		// Single value
		v := atoi(part)
		if v >= min && v <= max {
			result[v] = true
		}
	}
	return result
}

func splitCron(s string, sep byte) []string {
	var result []string
	current := ""
	for _, ch := range s {
		if byte(ch) == sep {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

func indexOf(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}

func atoi(s string) int {
	var n int
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			n = n*10 + int(ch-'0')
		}
	}
	return n
}

// DemoJobs pre-registers 5 demo jobs for the WMS system.
func DemoJobs(s *Scheduler) {
	_ = s.AddJob("bill-generation", "@daily", func() error {
		// Daily bill generation
		return nil
	})
	_ = s.AddJob("weight-cleanup", "@every 6h", func() error {
		// Clean old weight records
		return nil
	})
	_ = s.AddJob("statistics-report", "@every 1h", func() error {
		// Generate statistics
		return nil
	})
	_ = s.AddJob("backup-database", "0 2 * * *", func() error {
		// Database backup at 2 AM
		return nil
	})
	_ = s.AddJob("health-check", "@every 5m", func() error {
		// Device health check
		return nil
	})
}

// Ensure context import is available for handlers that need it.
var _ context.Context
