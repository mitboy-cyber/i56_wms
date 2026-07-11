// Package response provides unified JSON response formatting.
package response

import (
	"encoding/json"
	"net/http"

	"github.com/i56/framework/core/errors"
)

// Envelope is the standard JSON response wrapper.
type Envelope struct {
	Data  any       `json:"data,omitempty"`
	Meta  *Meta     `json:"meta,omitempty"`
	Error *APIError `json:"error,omitempty"`
}

// Meta contains pagination and request metadata.
type Meta struct {
	Total      int64  `json:"total,omitempty"`
	Page       int    `json:"page,omitempty"`
	PageSize   int    `json:"page_size,omitempty"`
	TotalPages int    `json:"total_pages,omitempty"`
	NextCursor string `json:"next_cursor,omitempty"`
	RequestID  string `json:"request_id"`
}

// APIError is the error portion of the response envelope.
type APIError struct {
	Code    string              `json:"code"`
	Message string              `json:"message"`
	Details []errors.ErrorDetail `json:"details,omitempty"`
}

// JSON writes a successful JSON response.
func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Envelope{Data: data})
}

// JSONWithMeta writes a successful JSON response with metadata.
func JSONWithMeta(w http.ResponseWriter, status int, data any, meta *Meta) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Envelope{Data: data, Meta: meta})
}

// Error writes an error JSON response.
func Error(w http.ResponseWriter, err error) {
	if err == nil {
		err = errors.NewInternal("internal server error")
	}
	appErr, ok := err.(*errors.AppError)
	if !ok {
		appErr = errors.NewInternal(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.HTTPStatus)
	json.NewEncoder(w).Encode(Envelope{
		Error: &APIError{
			Code:    appErr.Code,
			Message: appErr.Message,
			Details: appErr.Details,
		},
	})
}

// PaginatedResponse wraps data with pagination info.
type PaginatedResponse struct {
	Data       []any `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// PaginatedJSON writes a paginated response.
func PaginatedJSON(w http.ResponseWriter, data []any, total int64, page, pageSize int) {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	JSONWithMeta(w, http.StatusOK, data, &Meta{
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// NoContent writes a 204 No Content response.
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// Created writes a 201 Created response.
func Created(w http.ResponseWriter, data any) {
	JSON(w, http.StatusCreated, data)
}
