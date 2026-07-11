package wms

import (
	"github.com/gin-gonic/gin"
)

// Module represents the WMS business module.
// Implements the i56-framework module lifecycle contract.
type Module struct {
	name        string
	version     string
	permissions []string
	routes      []RouteRegistration
	subscribers []EventSubscription
}

// RouteRegistration holds a route pattern and its handler.
type RouteRegistration struct {
	Method  string
	Pattern string
	Handler gin.HandlerFunc
}

// EventSubscription maps an event name to a handler.
type EventSubscription struct {
	Event   string
	Handler func(event interface{})
}

// New creates a new WMS module instance.
func New() *Module {
	m := &Module{
		name:    "wms",
		version: "1.0.0",
	}

	// Register permissions
	m.permissions = []string{
		"wms.inventory.view",
		"wms.inventory.adjust",
		"wms.receiving.create",
		"wms.picking.execute",
		"wms.packing.execute",
		"wms.shipping.create",
	}

	// Register event subscriptions
	m.subscribers = []EventSubscription{
		{Event: "parcel.received", Handler: m.onParcelReceived},
		{Event: "order.created", Handler: m.onOrderCreated},
	}

	return m
}

// Name returns the module identifier.
func (m *Module) Name() string { return m.name }

// Version returns the module version.
func (m *Module) Version() string { return m.version }

// Permissions returns the permission slugs this module owns.
func (m *Module) Permissions() []string { return m.permissions }

// Routes returns the HTTP route registrations.
func (m *Module) Routes() []RouteRegistration { return m.routes }

// Subscriptions returns event subscriptions for decoupled communication.
func (m *Module) Subscriptions() []EventSubscription { return m.subscribers }

// Register adds an HTTP route to this module.
func (m *Module) Register(method, pattern string, handler gin.HandlerFunc) {
	m.routes = append(m.routes, RouteRegistration{Method: method, Pattern: pattern, Handler: handler})
}

// On registers an event subscription.
func (m *Module) On(event string, handler func(event interface{})) {
	m.subscribers = append(m.subscribers, EventSubscription{Event: event, Handler: handler})
}

// Event handlers (domain events from other modules)
func (m *Module) onParcelReceived(event interface{}) {
	// Update inventory when a parcel is received
}

func (m *Module) onOrderCreated(event interface{}) {
	// Reserve warehouse space when an order is created
}

// Ensure Module satisfies the framework.Module interface
var _ interface {
	Name() string
	Version() string
	Permissions() []string
	Routes() []RouteRegistration
	Subscriptions() []EventSubscription
} = (*Module)(nil)
