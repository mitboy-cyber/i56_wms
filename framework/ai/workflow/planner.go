// Package workflow implements ReAct and Plan-and-Solve DAG executors
// for orchestrating multi-step AI reasoning tasks.
package workflow

import (
	"context"
	"fmt"
	"sync"
)

// Step represents a single node in a workflow DAG.
type Step struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	DependsOn   []string `json:"depends_on"` // IDs of prerequisite steps
	Action      string   `json:"action"`     // The action/prompt to execute
}

// Result holds the output of a completed step.
type Result struct {
	StepID string `json:"step_id"`
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
	Status string `json:"status"` // "success", "failed", "skipped"
}

// Plan is a directed acyclic graph of steps.
type Plan struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Steps       []Step `json:"steps"`
}

// Executor handles the execution of a workflow plan.
type Executor interface {
	Execute(ctx context.Context, plan *Plan) ([]Result, error)
}

// ReActExecutor implements the ReAct (Reasoning + Acting) pattern.
type ReActExecutor struct {
	mu           sync.Mutex
	thoughts     []string
	observations []string
	maxIterations int
}

// NewReActExecutor creates a ReAct executor.
func NewReActExecutor(maxIterations int) *ReActExecutor {
	if maxIterations <= 0 {
		maxIterations = 10
	}
	return &ReActExecutor{
		maxIterations: maxIterations,
	}
}

// Execute runs the ReAct loop: Thought → Action → Observation → repeat.
func (e *ReActExecutor) Execute(ctx context.Context, plan *Plan) ([]Result, error) {
	var results []Result

	for i := 0; i < e.maxIterations; i++ {
		select {
		case <-ctx.Done():
			return results, ctx.Err()
		default:
		}

		// TODO: Implement ReAct loop with LLM integration
		// Thought: LLM reasons about current state
		// Action: LLM selects and invokes a tool
		// Observation: Parse tool output and decide next step
		results = append(results, Result{
			StepID: fmt.Sprintf("react-step-%d", i),
			Output: "stub: ReAct iteration",
			Status: "success",
		})

		// TODO: Check termination condition
		if i >= 2 {
			break
		}
	}

	return results, nil
}

// PlanAndSolveExecutor implements the Plan-and-Solve pattern.
type PlanAndSolveExecutor struct {
	mu sync.Mutex
}

// NewPlanAndSolveExecutor creates a Plan-and-Solve executor.
func NewPlanAndSolveExecutor() *PlanAndSolveExecutor {
	return &PlanAndSolveExecutor{}
}

// Execute runs Plan-and-Solve: Plan → Execute steps → Solve.
func (e *PlanAndSolveExecutor) Execute(ctx context.Context, plan *Plan) ([]Result, error) {
	var results []Result

	// Phase 1: Plan — resolve step dependencies and order
	sorted, err := topologicalSort(plan.Steps)
	if err != nil {
		return nil, fmt.Errorf("workflow: cannot resolve dependencies: %w", err)
	}

	// Phase 2: Execute each step in dependency order
	completed := make(map[string]string) // stepID → output
	for _, step := range sorted {
		select {
		case <-ctx.Done():
			return results, ctx.Err()
		default:
		}

		// TODO: Execute step via AI pipeline
		result := Result{
			StepID: step.ID,
			Output: fmt.Sprintf("stub: executed step %q", step.Name),
			Status: "success",
		}
		completed[step.ID] = result.Output
		results = append(results, result)
	}

	// Phase 3: Solve — synthesize final answer
	results = append(results, Result{
		StepID: "solve",
		Output: "stub: final synthesis",
		Status: "success",
	})

	return results, nil
}

// topologicalSort returns steps in dependency-safe execution order using Kahn's algorithm.
func topologicalSort(steps []Step) ([]Step, error) {
	// Build adjacency and in-degree
	graph := make(map[string][]string)
	inDegree := make(map[string]int)
	stepMap := make(map[string]Step)

	for _, s := range steps {
		stepMap[s.ID] = s
		if _, ok := inDegree[s.ID]; !ok {
			inDegree[s.ID] = 0
		}
		for _, dep := range s.DependsOn {
			graph[dep] = append(graph[dep], s.ID)
			inDegree[s.ID]++
		}
	}

	// Find all nodes with no dependencies
	var queue []string
	for id, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, id)
		}
	}

	var sorted []Step
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		sorted = append(sorted, stepMap[id])
		for _, next := range graph[id] {
			inDegree[next]--
			if inDegree[next] == 0 {
				queue = append(queue, next)
			}
		}
	}

	if len(sorted) != len(steps) {
		return nil, fmt.Errorf("cycle detected in workflow DAG")
	}

	return sorted, nil
}
