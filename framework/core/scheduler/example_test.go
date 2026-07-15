package scheduler_test

import (
	"fmt"
	"sync"
	"time"

	"github.com/i56/framework/core/scheduler"
)

// ExampleScheduler demonstrates registering and running scheduled jobs.
func ExampleScheduler() {
	s := scheduler.New()

	var mu sync.Mutex
	var runTimes []string

	// Register a job that runs every 500ms
	s.AddJob("health-check", "@every 500ms", func() error {
		mu.Lock()
		runTimes = append(runTimes, time.Now().Format("15:04:05"))
		mu.Unlock()
		return nil
	})

	// Register a daily job
	s.AddJob("daily-report", "@daily", func() error {
		fmt.Println("Generating daily report...")
		return nil
	})

	// Start the scheduler
	s.Start()

	// Let it run briefly
	time.Sleep(600 * time.Millisecond)
	s.Stop()

	// Print status
	jobs := s.ListJobs()
	for _, j := range jobs {
		fmt.Printf("Job: %s, Runs: %d\n", j.Name, j.RunCount)
	}
	// Example output (at least one run):
	// Job: daily-report, Runs: 0
	// Job: health-check, Runs: 1
}
