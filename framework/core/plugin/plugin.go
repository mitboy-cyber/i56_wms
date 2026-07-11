// Package plugin provides a lightweight plugin system with lifecycle management.
// Plugins register implementations of service interfaces; the framework wire them
// together through the PluginRegistry.
package plugin

import (
	"context"
	"fmt"
	"sync"

	"github.com/i56/framework/core/logger"
)

// Plugin is the interface that all plugins must implement.
type Plugin interface {
	Name() string
	Version() string
	Init(ctx context.Context, reg *Registry) error
}

// Priority groups for ordered startup/shutdown.
type Priority int

const (
	PrioritySystem   Priority = 0  // config, logger (lowest = starts first)
	PriorityCore     Priority = 10 // auth, rbac, tenant
	PriorityData     Priority = 20 // database, cache, storage
	PriorityService  Priority = 30 // business services
	PriorityAPI      Priority = 40 // http, grpc servers
	PriorityExternal Priority = 50 // kafka, nats, integrations
)

// Registry holds all registered plugins and service lookups.
type Registry struct {
	mu       sync.RWMutex
	plugins  []pluginEntry
	services map[string]any // name → implementation
	log      logger.Logger
	status   Status
}

type pluginEntry struct {
	plugin   Plugin
	priority Priority
}

// Status tracks the registry lifecycle.
type Status int

const (
	StatusUninitialized Status = iota
	StatusInitializing
	StatusRunning
	StatusStopped
)

// NewRegistry creates a new PluginRegistry.
func NewRegistry(log logger.Logger) *Registry {
	return &Registry{
		plugins:  make([]pluginEntry, 0),
		services: make(map[string]any),
		log:      log,
		status:   StatusUninitialized,
	}
}

// Register adds a plugin to the registry (order not guaranteed — use priority).
func (r *Registry) Register(p Plugin, priority Priority) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.plugins = append(r.plugins, pluginEntry{plugin: p, priority: priority})
}

// Provide registers a service implementation by name for later lookup.
func (r *Registry) Provide(name string, svc any) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.services[name] = svc
	r.log.Debug("plugin: service registered", "name", name)
}

// Resolve retrieves a registered service by name. Returns nil if not found.
func (r *Registry) Resolve(name string) any {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.services[name]
}

// MustResolve retrieves a service and panics if not found.
func (r *Registry) MustResolve(name string) any {
	svc := r.Resolve(name)
	if svc == nil {
		panic(fmt.Sprintf("plugin: service %q not registered", name))
	}
	return svc
}

// Start initializes all plugins in priority order and sets status to Running.
func (r *Registry) Start(ctx context.Context) error {
	r.mu.Lock()
	if r.status != StatusUninitialized {
		r.mu.Unlock()
		return fmt.Errorf("plugin: registry already started (status=%d)", r.status)
	}
	r.status = StatusInitializing

	// Sort by priority (stable for same-priority entries)
	sortByPriority(r.plugins)
	entries := make([]pluginEntry, len(r.plugins))
	copy(entries, r.plugins)
	r.mu.Unlock()

	for _, e := range entries {
		r.log.Info("plugin: initializing", "name", e.plugin.Name(), "version", e.plugin.Version(), "priority", int(e.priority))
		if err := e.plugin.Init(ctx, r); err != nil {
			r.mu.Lock()
			r.status = StatusStopped
			r.mu.Unlock()
			return fmt.Errorf("plugin %s init: %w", e.plugin.Name(), err)
		}
	}

	r.mu.Lock()
	r.status = StatusRunning
	r.mu.Unlock()
	r.log.Info("plugin: all plugins initialized", "count", len(r.plugins))
	return nil
}

// Status returns the current registry status.
func (r *Registry) Status() Status {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.status
}

// List returns all registered plugin names.
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, len(r.plugins))
	for i, e := range r.plugins {
		names[i] = e.plugin.Name()
	}
	return names
}

// PluginCount returns the number of registered plugins.
func (r *Registry) PluginCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.plugins)
}

// sortByPriority performs insertion sort on plugin entries (stable, O(n²) but n is small).
func sortByPriority(entries []pluginEntry) {
	for i := 1; i < len(entries); i++ {
		j := i
		for j > 0 && entries[j].priority < entries[j-1].priority {
			entries[j], entries[j-1] = entries[j-1], entries[j]
			j--
		}
	}
}

// BasePlugin provides a minimal Plugin implementation for embedding.
type BasePlugin struct {
	name    string
	version string
	initFn  func(ctx context.Context, reg *Registry) error
}

// NewBasePlugin creates a plugin from a name, version, and init function.
func NewBasePlugin(name, version string, initFn func(ctx context.Context, reg *Registry) error) *BasePlugin {
	return &BasePlugin{name: name, version: version, initFn: initFn}
}

func (p *BasePlugin) Name() string    { return p.name }
func (p *BasePlugin) Version() string { return p.version }
func (p *BasePlugin) Init(ctx context.Context, reg *Registry) error {
	if p.initFn != nil {
		return p.initFn(ctx, reg)
	}
	return nil
}
