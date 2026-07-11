// Package agent implements a long-running EnterpriseAgent that operates
// on an event-driven model, reacting to triggers and executing AI workflows
// asynchronously.
package agent

import (
	"context"
	"sync"
	"time"
)

// State represents the current lifecycle state of an agent.
type State string

const (
	StateIdle      State = "idle"
	StateRunning   State = "running"
	StatePaused    State = "paused"
	StateStopped   State = "stopped"
	StateError     State = "error"
)

// Event represents a trigger that the agent reacts to.
type Event struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Source    string    `json:"source"`
	Payload   any       `json:"payload"`
	Timestamp time.Time `json:"timestamp"`
}

// Task is a unit of work the agent executes.
type Task struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
	Status      string    `json:"status"`
	Result      string    `json:"result,omitempty"`
}

// AgentHandle is a function the agent can use to handle specific event types.
type AgentHandle func(ctx context.Context, event Event) error

// EnterpriseAgent is a long-running goroutine that processes events and
// executes AI-powered tasks.
type EnterpriseAgent struct {
	mu       sync.RWMutex
	Name     string
	State    State
	Handlers map[string]AgentHandle // event type → handler
	Tasks    chan Task
	Events   chan Event
	stopCh   chan struct{}
	doneCh   chan struct{}
}

// New creates a new EnterpriseAgent.
func New(name string, bufferSize int) *EnterpriseAgent {
	return &EnterpriseAgent{
		Name:     name,
		State:    StateIdle,
		Handlers: make(map[string]AgentHandle),
		Tasks:    make(chan Task, bufferSize),
		Events:   make(chan Event, bufferSize),
		stopCh:   make(chan struct{}),
		doneCh:   make(chan struct{}),
	}
}

// RegisterHandler maps an event type to a handler function.
func (a *EnterpriseAgent) RegisterHandler(eventType string, handler AgentHandle) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Handlers[eventType] = handler
}

// Start launches the agent's event loop in a background goroutine.
func (a *EnterpriseAgent) Start(ctx context.Context) error {
	a.mu.Lock()
	if a.State == StateRunning {
		a.mu.Unlock()
		return nil
	}
	a.State = StateRunning
	a.mu.Unlock()

	go a.loop(ctx)
	return nil
}

// Stop gracefully shuts down the agent.
func (a *EnterpriseAgent) Stop() {
	a.mu.Lock()
	a.State = StateStopped
	a.mu.Unlock()
	close(a.stopCh)
}

// SendEvent pushes an event into the agent's event queue.
func (a *EnterpriseAgent) SendEvent(event Event) {
	a.Events <- event
}

// SubmitTask queues a task for the agent.
func (a *EnterpriseAgent) SubmitTask(task Task) {
	a.Tasks <- task
}

// GetState returns the current agent state.
func (a *EnterpriseAgent) GetState() State {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.State
}

// Wait blocks until the agent has stopped.
func (a *EnterpriseAgent) Wait() {
	<-a.doneCh
}

func (a *EnterpriseAgent) loop(ctx context.Context) {
	defer close(a.doneCh)

	for {
		select {
		case <-a.stopCh:
			return
		case <-ctx.Done():
			return
		case event := <-a.Events:
			a.handleEvent(ctx, event)
		case task := <-a.Tasks:
			a.handleTask(ctx, task)
		}
	}
}

func (a *EnterpriseAgent) handleEvent(ctx context.Context, event Event) {
	a.mu.RLock()
	handler, ok := a.Handlers[event.Type]
	a.mu.RUnlock()
	if ok && handler != nil {
		_ = handler(ctx, event) // TODO: error handling and retry
	}
}

func (a *EnterpriseAgent) handleTask(ctx context.Context, task Task) {
	// TODO: Execute task using AI pipeline
	_ = ctx
	_ = task
}
