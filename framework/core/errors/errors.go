// Package errors provides unified error codes and AppError type.
package errors

import (
	"fmt"
	"net/http"
)

// Predefined error codes.
const (
	ErrCodeNotFound          = "NOT_FOUND"
	ErrCodeValidation        = "VALIDATION_ERROR"
	ErrCodeUnauthorized      = "UNAUTHORIZED"
	ErrCodeForbidden         = "FORBIDDEN"
	ErrCodeConflict          = "CONFLICT"
	ErrCodeInternal          = "INTERNAL_ERROR"
	ErrCodeTenantRequired    = "TENANT_REQUIRED"
	ErrCodeInvalidTransition = "INVALID_STATE_TRANSITION"
	ErrCodeRateLimited       = "RATE_LIMITED"
)

// ErrorDetail carries per-field validation details.
type ErrorDetail struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// AppError is the unified application error type.
type AppError struct {
	Code       string        `json:"code"`
	Message    string        `json:"message"`
	HTTPStatus int           `json:"-"`
	Details    []ErrorDetail `json:"details,omitempty"`
	Cause      error         `json:"-"`
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error { return e.Cause }

// WithDetail adds a validation detail.
func (e *AppError) WithDetail(field, msg string) *AppError {
	e.Details = append(e.Details, ErrorDetail{Field: field, Message: msg})
	return e
}

// WithCause wraps an underlying error.
func (e *AppError) WithCause(err error) *AppError {
	e.Cause = err
	return e
}

// Factory functions for common errors.

func NewNotFound(resource string) *AppError {
	return &AppError{
		Code:       ErrCodeNotFound,
		Message:    fmt.Sprintf("%s not found", resource),
		HTTPStatus: http.StatusNotFound,
	}
}

func NewValidation(msg string) *AppError {
	return &AppError{
		Code:       ErrCodeValidation,
		Message:    msg,
		HTTPStatus: http.StatusUnprocessableEntity,
	}
}

func NewUnauthorized(msg string) *AppError {
	if msg == "" {
		msg = "authentication required"
	}
	return &AppError{
		Code:       ErrCodeUnauthorized,
		Message:    msg,
		HTTPStatus: http.StatusUnauthorized,
	}
}

func NewForbidden(msg string) *AppError {
	if msg == "" {
		msg = "insufficient permissions"
	}
	return &AppError{
		Code:       ErrCodeForbidden,
		Message:    msg,
		HTTPStatus: http.StatusForbidden,
	}
}

func NewConflict(msg string) *AppError {
	return &AppError{
		Code:       ErrCodeConflict,
		Message:    msg,
		HTTPStatus: http.StatusConflict,
	}
}

func NewInternal(msg string) *AppError {
	if msg == "" {
		msg = "internal server error"
	}
	return &AppError{
		Code:       ErrCodeInternal,
		Message:    msg,
		HTTPStatus: http.StatusInternalServerError,
	}
}

func NewInvalidTransition(from, to string) *AppError {
	return &AppError{
		Code:       ErrCodeInvalidTransition,
		Message:    fmt.Sprintf("cannot transition from %s to %s", from, to),
		HTTPStatus: http.StatusUnprocessableEntity,
	}
}
