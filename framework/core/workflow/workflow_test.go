package workflow

import (
	"context"
	"testing"

	"github.com/i56/framework/core/logger"
)

type testStore struct {
	defs      map[string]*ProcessDefinition
	instances map[string]*ProcessInstance
}

func newTestStore() *testStore {
	return &testStore{
		defs:      make(map[string]*ProcessDefinition),
		instances: make(map[string]*ProcessInstance),
	}
}

func (s *testStore) GetDefinition(ctx context.Context, id string) (*ProcessDefinition, error) {
	d, ok := s.defs[id]
	if !ok {
		return nil, nil
	}
	return d, nil
}

func (s *testStore) SaveInstance(ctx context.Context, instance *ProcessInstance) error {
	s.instances[instance.ID] = instance
	return nil
}

func (s *testStore) GetInstance(ctx context.Context, id string) (*ProcessInstance, error) {
	i, ok := s.instances[id]
	if !ok {
		return nil, nil
	}
	return i, nil
}

// testLogger implements logger.Logger.
type testLogger struct{}

func (l testLogger) Debug(msg string, args ...any)  {}
func (l testLogger) Info(msg string, args ...any)   {}
func (l testLogger) Warn(msg string, args ...any)   {}
func (l testLogger) Error(msg string, args ...any)  {}
func (l testLogger) With(args ...any) logger.Logger { return l }
func (l testLogger) WithGroup(name string) logger.Logger { return l }

var _ logger.Logger = testLogger{}

func TestEngine_StartProcess(t *testing.T) {
	store := newTestStore()
	engine := NewEngine(store, testLogger{})

	def := &ProcessDefinition{
		ID:   "order-fulfillment",
		Name: "Order Fulfillment",
		States: []State{
			{ID: "start", Name: "Order Received", Type: StateStart},
			{ID: "packing", Name: "Packing", Type: StateTask},
			{ID: "shipped", Name: "Shipped", Type: StateEnd},
		},
		Transitions: []Transition{
			{From: "start", To: "packing"},
			{From: "packing", To: "shipped"},
		},
	}
	engine.RegisterDefinition(def)

	inst, err := engine.StartProcess(context.Background(), "order-fulfillment", map[string]any{"order_id": "123"})
	if err != nil {
		t.Fatalf("StartProcess: %v", err)
	}
	if inst.CurrentState != "start" {
		t.Errorf("expected start state, got %q", inst.CurrentState)
	}
	if inst.Status != StatusRunning {
		t.Errorf("expected running status, got %q", inst.Status)
	}
	if inst.Variables["order_id"] != "123" {
		t.Errorf("expected order_id=123, got %v", inst.Variables["order_id"])
	}
}

func TestEngine_Transition(t *testing.T) {
	store := newTestStore()
	engine := NewEngine(store, testLogger{})

	def := &ProcessDefinition{
		ID:   "simple",
		Name: "Simple",
		States: []State{
			{ID: "start", Name: "Start", Type: StateStart},
			{ID: "end", Name: "End", Type: StateEnd},
		},
		Transitions: []Transition{
			{From: "start", To: "end"},
		},
	}
	engine.RegisterDefinition(def)

	inst, _ := engine.StartProcess(context.Background(), "simple", nil)

	err := engine.Transition(context.Background(), inst.ID, "end", nil)
	if err != nil {
		t.Fatalf("Transition: %v", err)
	}

	updated, _ := store.GetInstance(context.Background(), inst.ID)
	if updated.Status != StatusCompleted {
		t.Errorf("expected completed status, got %q", updated.Status)
	}
	if updated.CompletedAt == nil {
		t.Error("expected CompletedAt to be set")
	}
}

func TestEngine_InvalidTransition(t *testing.T) {
	store := newTestStore()
	engine := NewEngine(store, testLogger{})

	def := &ProcessDefinition{
		ID:   "locked",
		Name: "Locked",
		States: []State{
			{ID: "a", Name: "A", Type: StateStart},
			{ID: "b", Name: "B", Type: StateTask},
			{ID: "c", Name: "C", Type: StateEnd},
		},
		Transitions: []Transition{
			{From: "a", To: "b"},
			{From: "b", To: "c"},
		},
	}
	engine.RegisterDefinition(def)

	inst, _ := engine.StartProcess(context.Background(), "locked", nil)

	// Try to jump from a → c (invalid)
	err := engine.Transition(context.Background(), inst.ID, "c", nil)
	if err == nil {
		t.Error("expected error for invalid transition")
	}
}

func TestEngine_CancelProcess(t *testing.T) {
	store := newTestStore()
	engine := NewEngine(store, testLogger{})

	def := &ProcessDefinition{
		ID:   "cancelable",
		Name: "Cancelable",
		States: []State{
			{ID: "start", Name: "Start", Type: StateStart},
			{ID: "end", Name: "End", Type: StateEnd},
		},
		Transitions: []Transition{
			{From: "start", To: "end"},
		},
	}
	engine.RegisterDefinition(def)

	inst, _ := engine.StartProcess(context.Background(), "cancelable", nil)

	err := engine.CancelProcess(context.Background(), inst.ID)
	if err != nil {
		t.Fatalf("CancelProcess: %v", err)
	}

	updated, _ := store.GetInstance(context.Background(), inst.ID)
	if updated.Status != StatusCancelled {
		t.Errorf("expected cancelled status, got %q", updated.Status)
	}
}

func TestEngine_GetValidTransitions(t *testing.T) {
	store := newTestStore()
	engine := NewEngine(store, testLogger{})

	def := &ProcessDefinition{
		ID:   "multi",
		Name: "Multi",
		States: []State{
			{ID: "start", Name: "Start", Type: StateStart},
			{ID: "step1", Name: "Step 1", Type: StateTask},
			{ID: "step2", Name: "Step 2", Type: StateTask},
			{ID: "end", Name: "End", Type: StateEnd},
		},
		Transitions: []Transition{
			{From: "start", To: "step1"},
			{From: "start", To: "step2"},
			{From: "step1", To: "end"},
			{From: "step2", To: "end"},
		},
	}
	engine.RegisterDefinition(def)

	inst, _ := engine.StartProcess(context.Background(), "multi", nil)

	transitions, err := engine.GetValidTransitions(context.Background(), inst.ID)
	if err != nil {
		t.Fatalf("GetValidTransitions: %v", err)
	}
	if len(transitions) != 2 {
		t.Errorf("expected 2 transitions from start, got %d", len(transitions))
	}
}

func TestEngine_StartProcessUnknownDefinition(t *testing.T) {
	store := newTestStore()
	engine := NewEngine(store, testLogger{})

	_, err := engine.StartProcess(context.Background(), "nonexistent", nil)
	if err == nil {
		t.Error("expected error for unknown definition")
	}
}
