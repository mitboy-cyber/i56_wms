// Package workflow provides a high-level approval workflow engine with
// built-in warehouse process definitions (inbound, outbound, qc).
package workflow

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

// ActionType defines what can happen at a step.
type ActionType string

const (
	ActionApprove ActionType = "approve"
	ActionReject  ActionType = "reject"
	ActionSkip    ActionType = "skip"
)

// Step is a named action node inside a workflow definition.
type Step struct {
	Name     string     `json:"name"`
	Assignee string     `json:"assignee"`
	SLAHours int        `json:"sla_hours"`
	Action   ActionType `json:"action"`
}

// WorkflowDef describes a reusable workflow blueprint.
type WorkflowDef struct {
	Name     string   `json:"name"`
	Steps    []Step   `json:"steps"`
	Triggers []string `json:"triggers"`
}

// WorkflowInstanceStatus tracks the lifecycle of a running instance.
type WorkflowInstanceStatus string

const (
	StatusPending   WorkflowInstanceStatus = "pending"
	StatusRunning   WorkflowInstanceStatus = "running"
	StatusApproved  WorkflowInstanceStatus = "approved"
	StatusRejected  WorkflowInstanceStatus = "rejected"
	StatusCompleted WorkflowInstanceStatus = "completed"
)

// WorkflowInstance is a single execution of a workflow definition.
type WorkflowInstance struct {
	ID           int                    `json:"id"`
	DefinitionID string                 `json:"definition_id"`
	CurrentStep  string                 `json:"current_step"`
	Status       WorkflowInstanceStatus `json:"status"`
	Data         map[string]interface{} `json:"data"`
	History      []StepRecord           `json:"history"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// StepRecord records what happened at each step.
type StepRecord struct {
	Step   string    `json:"step"`
	Actor  string    `json:"actor"`
	Action string    `json:"action"`
	At     time.Time `json:"at"`
}

// ---------------------------------------------------------------------------
// Engine
// ---------------------------------------------------------------------------

// WorkflowEngine manages workflow definitions and running instances.
type WorkflowEngine struct {
	mu          sync.RWMutex
	definitions map[string]*WorkflowDef
	instances   map[int]*WorkflowInstance
	nextID      int
}

// NewEngine creates a ready-to-use WorkflowEngine with built-in workflows
// already registered.
func NewEngine() *WorkflowEngine {
	e := &WorkflowEngine{
		definitions: make(map[string]*WorkflowDef),
		instances:   make(map[int]*WorkflowInstance),
		nextID:      1,
	}
	e.registerBuiltins()
	return e
}

// Register adds a workflow definition under the given name.
func (e *WorkflowEngine) Register(name string, def *WorkflowDef) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.definitions[name] = def
}

// Start creates a new WorkflowInstance for the named definition and returns it.
func (e *WorkflowEngine) Start(tenantID, entityType, entityID int, data map[string]interface{}) (*WorkflowInstance, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	typeKey := fmt.Sprintf("%d", entityType)
	def, ok := e.definitions[typeKey]
	if !ok {
		return nil, fmt.Errorf("workflow: no definition registered for %q", typeKey)
	}
	if len(def.Steps) == 0 {
		return nil, fmt.Errorf("workflow: definition %q has no steps", typeKey)
	}

	now := time.Now()
	inst := &WorkflowInstance{
		ID:           e.nextID,
		DefinitionID: typeKey,
		CurrentStep:  def.Steps[0].Name,
		Status:       StatusRunning,
		Data:         data,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if inst.Data == nil {
		inst.Data = make(map[string]interface{})
	}
	inst.Data["tenant_id"] = tenantID
	inst.Data["entity_type"] = entityType
	inst.Data["entity_id"] = entityID

	e.nextID++
	e.instances[inst.ID] = inst
	return inst, nil
}

// Advance moves the workflow instance forward.  stepName must match the
// current step; action and actor are recorded in the instance history.
func (e *WorkflowEngine) Advance(instanceID int, stepName, action, actor string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	inst, ok := e.instances[instanceID]
	if !ok {
		return fmt.Errorf("workflow: instance %d not found", instanceID)
	}
	if inst.Status != StatusRunning && inst.Status != StatusPending {
		return fmt.Errorf("workflow: instance %d is %s", instanceID, inst.Status)
	}
	if inst.CurrentStep != stepName {
		return fmt.Errorf("workflow: expected current step %q, got %q", inst.CurrentStep, stepName)
	}

	def := e.definitions[inst.DefinitionID]
	// Record history
	inst.History = append(inst.History, StepRecord{
		Step:   stepName,
		Actor:  actor,
		Action: action,
		At:     time.Now(),
	})

	switch action {
	case "approve":
		next := e.findNextStep(def, stepName)
		if next == "" {
			inst.Status = StatusCompleted
			inst.CurrentStep = ""
		} else {
			inst.CurrentStep = next
		}
	case "reject":
		inst.Status = StatusRejected
		inst.CurrentStep = ""
	case "skip":
		next := e.findNextStep(def, stepName)
		if next == "" {
			inst.Status = StatusCompleted
			inst.CurrentStep = ""
		} else {
			inst.CurrentStep = next
		}
	default:
		return fmt.Errorf("workflow: unknown action %q", action)
	}

	inst.UpdatedAt = time.Now()
	return nil
}

// Status returns the current WorkflowInstance by id.
func (e *WorkflowEngine) Status(instanceID int) (*WorkflowInstance, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	inst, ok := e.instances[instanceID]
	if !ok {
		return nil, fmt.Errorf("workflow: instance %d not found", instanceID)
	}
	return inst, nil
}

// GetDefinition returns a registered definition by name.
func (e *WorkflowEngine) GetDefinition(name string) (*WorkflowDef, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	def, ok := e.definitions[name]
	if !ok {
		return nil, fmt.Errorf("workflow: definition %q not found", name)
	}
	return def, nil
}

// ListDefinitions returns all registered definition names.
func (e *WorkflowEngine) ListDefinitions() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	names := make([]string, 0, len(e.definitions))
	for k := range e.definitions {
		names = append(names, k)
	}
	return names
}

// helper: find the step after `current` in the definition's step list.
func (e *WorkflowEngine) findNextStep(def *WorkflowDef, current string) string {
	for i, s := range def.Steps {
		if s.Name == current && i+1 < len(def.Steps) {
			return def.Steps[i+1].Name
		}
	}
	return ""
}

// ---------------------------------------------------------------------------
// Built-in workflows
// ---------------------------------------------------------------------------

func (e *WorkflowEngine) registerBuiltins() {
	inbound := &WorkflowDef{
		Name: "inbound",
		Steps: []Step{
			{Name: "收货确认", Assignee: "receiver", SLAHours: 2, Action: ActionApprove},
			{Name: "称重", Assignee: "weigher", SLAHours: 1, Action: ActionApprove},
			{Name: "上架", Assignee: "putaway", SLAHours: 4, Action: ActionApprove},
		},
		Triggers: []string{"parcel_received"},
	}
	e.definitions["inbound"] = inbound
	e.definitions["0"] = inbound

	outbound := &WorkflowDef{
		Name: "outbound",
		Steps: []Step{
			{Name: "拣货", Assignee: "picker", SLAHours: 2, Action: ActionApprove},
			{Name: "复核", Assignee: "reviewer", SLAHours: 1, Action: ActionApprove},
			{Name: "打包", Assignee: "packer", SLAHours: 1, Action: ActionApprove},
			{Name: "发货", Assignee: "shipper", SLAHours: 2, Action: ActionApprove},
		},
		Triggers: []string{"order_created"},
	}
	e.definitions["outbound"] = outbound
	e.definitions["1"] = outbound

	qc := &WorkflowDef{
		Name: "qc",
		Steps: []Step{
			{Name: "质检申请", Assignee: "requester", SLAHours: 1, Action: ActionApprove},
			{Name: "质检确认", Assignee: "inspector", SLAHours: 4, Action: ActionApprove},
			{Name: "处置", Assignee: "disposer", SLAHours: 2, Action: ActionApprove},
		},
		Triggers: []string{"qc_requested"},
	}
	e.definitions["qc"] = qc
	e.definitions["2"] = qc
}

// ToJSON is a convenience helper that marshals an instance to JSON.
func (i *WorkflowInstance) ToJSON() string {
	b, _ := json.MarshalIndent(i, "", "  ")
	return string(b)
}
