// Package plugin provides a lightweight HTTP-oriented plugin system for the
// I56 framework.  Plugins contribute routes and navigation menu items.
package plugin

import (
	"fmt"
	"net/http"
	"sync"
)

// ---------------------------------------------------------------------------
// Core types
// ---------------------------------------------------------------------------

// Plugin is the interface every plugin must implement.
type Plugin interface {
	Name() string
	Version() string
	Init(config map[string]interface{}) error
	Routes() []Route
	MenuItems() []MenuItem
	Shutdown() error
}

// Route describes a single HTTP endpoint contributed by a plugin.
type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

// MenuItem describes a navigation entry contributed by a plugin.
type MenuItem struct {
	Group string
	Label string
	URL   string
	Icon  string
}

// ---------------------------------------------------------------------------
// PluginManager
// ---------------------------------------------------------------------------

// PluginManager manages the lifecycle of registered plugins.
type PluginManager struct {
	mu      sync.RWMutex
	plugins map[string]Plugin
	order   []string // registration order
}

// NewManager creates an empty PluginManager.
func NewManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]Plugin),
	}
}

// Register adds a plugin and immediately calls its Init method.
func (pm *PluginManager) Register(p Plugin, config map[string]interface{}) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if _, ok := pm.plugins[p.Name()]; ok {
		return fmt.Errorf("plugin: %q already registered", p.Name())
	}

	if err := p.Init(config); err != nil {
		return fmt.Errorf("plugin %s init: %w", p.Name(), err)
	}

	pm.plugins[p.Name()] = p
	pm.order = append(pm.order, p.Name())
	return nil
}

// Unregister calls Shutdown on the named plugin and removes it.
func (pm *PluginManager) Unregister(name string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	p, ok := pm.plugins[name]
	if !ok {
		return fmt.Errorf("plugin: %q not found", name)
	}
	if err := p.Shutdown(); err != nil {
		return fmt.Errorf("plugin %s shutdown: %w", name, err)
	}
	delete(pm.plugins, name)
	return nil
}

// ShutdownAll calls Shutdown on every registered plugin (reverse order).
func (pm *PluginManager) ShutdownAll() []error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	var errs []error
	for i := len(pm.order) - 1; i >= 0; i-- {
		name := pm.order[i]
		if p, ok := pm.plugins[name]; ok {
			if err := p.Shutdown(); err != nil {
				errs = append(errs, fmt.Errorf("plugin %s shutdown: %w", name, err))
			}
		}
	}
	pm.plugins = make(map[string]Plugin)
	pm.order = nil
	return errs
}

// GetAllRoutes aggregates routes from every registered plugin.
func (pm *PluginManager) GetAllRoutes() []Route {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	var all []Route
	for _, name := range pm.order {
		p := pm.plugins[name]
		all = append(all, p.Routes()...)
	}
	return all
}

// GetAllMenuItems aggregates menu items from every registered plugin.
func (pm *PluginManager) GetAllMenuItems() []MenuItem {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	var all []MenuItem
	for _, name := range pm.order {
		p := pm.plugins[name]
		all = append(all, p.MenuItems()...)
	}
	return all
}

// Get returns a plugin by name (nil if not found).
func (pm *PluginManager) Get(name string) Plugin {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.plugins[name]
}

// List returns registered plugin names in registration order.
func (pm *PluginManager) List() []string {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	out := make([]string, len(pm.order))
	copy(out, pm.order)
	return out
}
