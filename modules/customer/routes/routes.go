// Package routes registers customer module HTTP endpoints.
package routes

import (
	"github.com/i56/framework/core/router"
	"github.com/i56/modules/customer/handler"
)

// RegisterRoutes registers all customer endpoints on the given router.
func RegisterRoutes(r *router.Router, h *handler.ClientHandler) {
	r.GET("/clients", h.List)
	r.GET("/clients/{id}", h.GetByID)
	r.POST("/clients", h.Create)
	r.PATCH("/clients/{id}", h.Update)
	r.DELETE("/clients/{id}", h.Delete)
}
