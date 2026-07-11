// Package handler provides HTTP handlers for the customer module.
package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/i56/framework/core/response"
	"github.com/i56/modules/customer/dto"
	"github.com/i56/modules/customer/service"
)

// ClientHandler handles HTTP requests for clients.
type ClientHandler struct {
	svc *service.ClientService
}

// NewClientHandler creates a ClientHandler.
func NewClientHandler(svc *service.ClientService) *ClientHandler {
	return &ClientHandler{svc: svc}
}

// List handles GET /clients
func (h *ClientHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantID := int64(1) // TODO: from context
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	clients, total, err := h.svc.List(r.Context(), tenantID, offset, pageSize)
	if err != nil {
		response.Error(w, err)
		return
	}

	// Convert to []any for paginated response
	data := make([]any, len(clients))
	for i, c := range clients {
		data[i] = c
	}
	response.PaginatedJSON(w, data, total, page, pageSize)
}

// GetByID handles GET /clients/{id}
func (h *ClientHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path, "/api/v1/clients/")
	if err != nil {
		response.Error(w, err)
		return
	}

	tenantID := int64(1) // TODO: from context
	client, err := h.svc.GetByID(r.Context(), tenantID, id)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, client)
}

// Create handles POST /clients
func (h *ClientHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateClientRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, err)
		return
	}

	tenantID := int64(1) // TODO: from context
	client, err := h.svc.Create(r.Context(), tenantID, input)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, client)
}

// Update handles PATCH /clients/{id}
func (h *ClientHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path, "/api/v1/clients/")
	if err != nil {
		response.Error(w, err)
		return
	}

	var input dto.UpdateClientRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, err)
		return
	}

	tenantID := int64(1)
	client, err := h.svc.Update(r.Context(), tenantID, id, input)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, client)
}

// Delete handles DELETE /clients/{id}
func (h *ClientHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path, "/api/v1/clients/")
	if err != nil {
		response.Error(w, err)
		return
	}

	tenantID := int64(1)
	if err := h.svc.Delete(r.Context(), tenantID, id); err != nil {
		response.Error(w, err)
		return
	}
	response.NoContent(w)
}

func extractID(path, prefix string) (int64, error) {
	idStr := strings.TrimPrefix(path, prefix)
	idStr = strings.TrimSuffix(idStr, "/")
	// Remove any sub-paths
	if idx := strings.Index(idStr, "/"); idx >= 0 {
		idStr = idStr[:idx]
	}
	return strconv.ParseInt(idStr, 10, 64)
}
