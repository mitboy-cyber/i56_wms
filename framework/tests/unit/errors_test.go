package unit

import (
	"net/http"
	"testing"

	"github.com/i56/framework/core/errors"
)

func TestNewNotFound(t *testing.T) {
	err := errors.NewNotFound("Order")
	if err.Code != errors.ErrCodeNotFound {
		t.Errorf("expected code %s, got %s", errors.ErrCodeNotFound, err.Code)
	}
	if err.HTTPStatus != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, err.HTTPStatus)
	}
	if err.Message == "" {
		t.Error("expected non-empty message")
	}
}

func TestNewValidation(t *testing.T) {
	err := errors.NewValidation("name is required")
	if err.Code != errors.ErrCodeValidation {
		t.Errorf("expected code %s, got %s", errors.ErrCodeValidation, err.Code)
	}
	if err.HTTPStatus != http.StatusUnprocessableEntity {
		t.Errorf("expected status %d, got %d", http.StatusUnprocessableEntity, err.HTTPStatus)
	}
}

func TestNewUnauthorized(t *testing.T) {
	err := errors.NewUnauthorized("")
	if err.Code != errors.ErrCodeUnauthorized {
		t.Errorf("expected code %s, got %s", errors.ErrCodeUnauthorized, err.Code)
	}
	if err.HTTPStatus != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, err.HTTPStatus)
	}
}

func TestNewForbidden(t *testing.T) {
	err := errors.NewForbidden("")
	if err.Code != errors.ErrCodeForbidden {
		t.Errorf("expected code %s, got %s", errors.ErrCodeForbidden, err.Code)
	}
	if err.HTTPStatus != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, err.HTTPStatus)
	}
}

func TestNewConflict(t *testing.T) {
	err := errors.NewConflict("duplicate key")
	if err.Code != errors.ErrCodeConflict {
		t.Errorf("expected code %s, got %s", errors.ErrCodeConflict, err.Code)
	}
	if err.HTTPStatus != http.StatusConflict {
		t.Errorf("expected status %d, got %d", http.StatusConflict, err.HTTPStatus)
	}
}

func TestNewInternal(t *testing.T) {
	err := errors.NewInternal("")
	if err.Code != errors.ErrCodeInternal {
		t.Errorf("expected code %s, got %s", errors.ErrCodeInternal, err.Code)
	}
	if err.HTTPStatus != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, err.HTTPStatus)
	}
}

func TestNewInvalidTransition(t *testing.T) {
	err := errors.NewInvalidTransition("pending", "shipped")
	if err.Code != errors.ErrCodeInvalidTransition {
		t.Errorf("expected code %s, got %s", errors.ErrCodeInvalidTransition, err.Code)
	}
	if err.HTTPStatus != http.StatusUnprocessableEntity {
		t.Errorf("expected status %d, got %d", http.StatusUnprocessableEntity, err.HTTPStatus)
	}
}

func TestAppErrorWithDetail(t *testing.T) {
	err := errors.NewValidation("validation failed").
		WithDetail("name", "name is required").
		WithDetail("email", "invalid format")

	if len(err.Details) != 2 {
		t.Errorf("expected 2 details, got %d", len(err.Details))
	}
	if err.Details[0].Field != "name" {
		t.Errorf("expected field 'name', got '%s'", err.Details[0].Field)
	}
}

func TestAppErrorWithCause(t *testing.T) {
	cause := errors.NewNotFound("User")
	err := errors.NewInternal("wrapper").WithCause(cause)
	if err.Unwrap() != cause {
		t.Error("expected cause to be preserved")
	}
}

func TestAppErrorErrorString(t *testing.T) {
	err := errors.NewNotFound("Order")
	s := err.Error()
	if s == "" {
		t.Error("expected non-empty error string")
	}
}
