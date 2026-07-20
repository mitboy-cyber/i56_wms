// Package workflow provides a lightweight state machine and process engine.
package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/i56/framework/core/logger"
)

// StateType classifies states in a process.
type StateType string

const (
	StateStart   StateType = "start"
	StateTask    StateType = "task"
	StateGateway StateType = "gateway"
	StateEnd     StateType = "end"
)

// State represents a node in a process definition.
type State struct {
	ID   string    `json:"id"`
	Name string    `json:"name"`
	Type StateType `json:"type"`
}

// Transition defines a possible move between states.
type Transition struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Condition string `json:"condition,omitempty"` // Expression to evaluate
	Action    string `json:"action,omitempty"`    // Auto-execute action name
}

// ProcessDefinition defines a workflow blueprint.
type ProcessDefinition struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	States      []State      `json:"states"`
	Transitions []Transition `json:"transitions"`
}

// ProcessInstance is a running instance of a process.
type ProcessInstance struct {
	ID             string         `json:"id"`
	DefinitionID   string         `json:"definition_id"`
	CurrentState   string         `json:"current_state"`
	Variables      map[string]any `json:"variables"`
	Status         string         `json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	CompletedAt    *time.Time     `json:"completed_at,omitempty"`
}

// Instance status constants.
const (
	StatusRunning   = "running"
	StatusCompleted = "completed"
	StatusCancelled = "cancelled"
)

// ProcessStore persists process definitions and instances.
type ProcessStore interface {
	GetDefinition(ctx context.Context, id string) (*ProcessDefinition, error)
	SaveInstance(ctx context.Context, instance *ProcessInstance) error
	GetInstance(ctx context.Context, id string) (*ProcessInstance, error)
}

// Engine executes workflow processes.
type Engine struct {
	mu          sync.RWMutex
	definitions map[string]*ProcessDefinition
	store       ProcessStore
	log         logger.Logger
}

// NewEngine creates a workflow engine.
func NewEngine(store ProcessStore, log logger.Logger) *Engine {
	return &Engine{
		definitions: make(map[string]*ProcessDefinition),
		store:       store,
		log:         log,
	}
}

// RegisterDefinition adds a process definition to the engine.
func (e *Engine) RegisterDefinition(def *ProcessDefinition) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.definitions[def.ID] = def
}

// StartProcess creates and starts a new process instance.
func (e *Engine) StartProcess(ctx context.Context, defID string, vars map[string]any) (*ProcessInstance, error) {
	e.mu.RLock()
	def, ok := e.definitions[defID]
	e.mu.RUnlock()

	if !ok {
		return nil, errors.New("process definition not found: " + defID)
	}

	// Find start state
	var startState string
	for _, s := range def.States {
		if s.Type == StateStart {
			startState = s.ID
			break
		}
	}
	if startState == "" {
		return nil, errors.New("no start state in process: " + defID)
	}

	now := time.Now()
	inst := &ProcessInstance{
		ID:           generateID(),
		DefinitionID: defID,
		CurrentState: startState,
		Variables:    vars,
		Status:       StatusRunning,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := e.store.SaveInstance(ctx, inst); err != nil {
		return nil, err
	}

	e.log.Info("process started",
		"definition", def.Name,
		"instance", inst.ID,
		"state", startState,
	)

	return inst, nil
}

// Transition attempts to move a process instance to a new state.
func (e *Engine) Transition(ctx context.Context, instanceID, toState string, vars map[string]any) error {
	inst, err := e.store.GetInstance(ctx, instanceID)
	if err != nil {
		return err
	}

	e.mu.RLock()
	def, ok := e.definitions[inst.DefinitionID]
	e.mu.RUnlock()

	if !ok {
		return errors.New("process definition not found")
	}

	// Check if transition is valid
	valid := false
	for _, t := range def.Transitions {
		if t.From == inst.CurrentState && t.To == toState {
			// Evaluate condition if present
			if t.Condition != "" && !evaluateCondition(t.Condition, inst.Variables) {
				continue // condition not met, try next transition
			}
			valid = true
			break
		}
	}
	if !valid {
		return errors.New("invalid transition: " + inst.CurrentState + " -> " + toState)
	}

	inst.CurrentState = toState
	inst.UpdatedAt = time.Now()
	if vars != nil {
		for k, v := range vars {
			inst.Variables[k] = v
		}
	}

	// Check if reached end state
	for _, s := range def.States {
		if s.ID == toState && s.Type == StateEnd {
			inst.Status = StatusCompleted
			now := time.Now()
			inst.CompletedAt = &now
		}
	}

	if err := e.store.SaveInstance(ctx, inst); err != nil {
		return err
	}

	e.log.Info("process transition",
		"instance", instanceID,
		"from", inst.CurrentState,
		"to", toState,
	)

	return nil
}

// CancelProcess cancels a running process instance.
func (e *Engine) CancelProcess(ctx context.Context, instanceID string) error {
	inst, err := e.store.GetInstance(ctx, instanceID)
	if err != nil {
		return err
	}

	inst.Status = StatusCancelled
	inst.UpdatedAt = time.Now()
	now := time.Now()
	inst.CompletedAt = &now

	return e.store.SaveInstance(ctx, inst)
}

// GetValidTransitions returns all valid transitions from the current state.
func (e *Engine) GetValidTransitions(ctx context.Context, instanceID string) ([]Transition, error) {
	inst, err := e.store.GetInstance(ctx, instanceID)
	if err != nil {
		return nil, err
	}

	e.mu.RLock()
	def, ok := e.definitions[inst.DefinitionID]
	e.mu.RUnlock()

	if !ok {
		return nil, errors.New("process definition not found")
	}

	var valid []Transition
	for _, t := range def.Transitions {
		if t.From == inst.CurrentState {
			valid = append(valid, t)
		}
	}
	return valid, nil
}

func generateID() string {
	return time.Now().Format("20060102150405") + "-" + randomStr(8)
}

func randomStr(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}

// evaluateCondition parses simple expressions like "amount>5000" and evaluates against vars.
// Supports: >, <, >=, <=, ==
func evaluateCondition(expr string, vars map[string]any) bool {
	// Split by operator
	type op struct{ sym string; fn func(a, b float64) bool }
	ops := []op{
		{">=", func(a, b float64) bool { return a >= b }},
		{"<=", func(a, b float64) bool { return a <= b }},
		{">", func(a, b float64) bool { return a > b }},
		{"<", func(a, b float64) bool { return a < b }},
		{"==", func(a, b float64) bool { return a == b }},
	}
	for _, o := range ops {
		parts := splitTrim(expr, o.sym)
		if len(parts) == 2 {
			val := toFloat(vars[parts[0]])
			target := parseFloat(parts[1])
			return o.fn(val, target)
		}
	}
	return false // unparseable condition
}

func splitTrim(s, sep string) []string {
	parts := make([]string, 0, 2)
	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			parts = append(parts, s[:i], s[i+len(sep):])
			return parts
		}
	}
	return parts
}

func toFloat(v any) float64 {
	switch n := v.(type) {
	case float64: return n
	case int: return float64(n)
	case int64: return float64(n)
	case json.Number: f, _ := n.Float64(); return f
	}
	return 0
}

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
