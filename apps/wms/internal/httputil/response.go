// Package httputil provides HTTP response helpers with safe error handling.
package httputil

import (
	"encoding/json"
	"log"
	"net/http"
)

// APIResponse is the standard API response envelope.
type APIResponse struct {
	Success bool           `json:"success"`
	Data    interface{}    `json:"data,omitempty"`
	Error   string         `json:"error,omitempty"`
	Fields  map[string]string `json:"fields,omitempty"` // per-field validation errors
}

// OK sends a 200 success response.
func OK(w http.ResponseWriter, data interface{}) {
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: data})
}

// Created sends a 201 response.
func Created(w http.ResponseWriter, data interface{}) {
	writeJSON(w, http.StatusCreated, APIResponse{Success: true, Data: data})
}

// BadRequest sends a 400 with user-safe error message (hides DB details).
func BadRequest(w http.ResponseWriter, msg string) {
	writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: msg})
}

// ValidationError sends a 422 with per-field validation errors.
func ValidationError(w http.ResponseWriter, fields map[string]string) {
	writeJSON(w, http.StatusUnprocessableEntity, APIResponse{Success: false, Error: "数据校验失败", Fields: fields})
}

// NotFound sends a 404.
func NotFound(w http.ResponseWriter, msg string) {
	writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: msg})
}

// InternalError logs the real error and returns a sanitized message.
func InternalError(w http.ResponseWriter, realErr error) {
	log.Printf("[HTTP 500] %v", realErr)
	writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "服务器内部错误，请稍后重试"})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
