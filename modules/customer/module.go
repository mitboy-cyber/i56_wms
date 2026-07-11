// Package customer implements the Customer business module.
package customer

import (
	"github.com/i56/framework/core/router"
	"github.com/i56/modules/customer/handler"
	"github.com/i56/modules/customer/routes"
)

// Module implements the core.Module interface for the customer domain.
type Module struct {
	handler *handler.ClientHandler
}

// New creates a new customer Module.
// Takes a ClientRepository implementation from the infrastructure layer.
func New() *Module {
	// In production, the repository is injected from infrastructure.
	// For now, we create the handler without a backing repo (tests will mock).
	return &Module{}
}

// Name returns the module identifier.
func (m *Module) Name() string { return "customer" }

// Version returns the module version.
func (m *Module) Version() string { return "1.1.0" }

// RegisterRoutes registers HTTP routes.
func (m *Module) RegisterRoutes(r *router.Router) {
	if m.handler != nil {
		routes.RegisterRoutes(r, m.handler)
	}
}
