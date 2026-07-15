package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	apperr "github.com/i56/framework/core/errors"
)

func TestJSON_WritesEnvelope(t *testing.T) {
	w := httptest.NewRecorder()
	JSON(w, http.StatusOK, map[string]string{"key": "value"})

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}

	var envelope Envelope
	if err := json.NewDecoder(w.Body).Decode(&envelope); err != nil {
		t.Fatalf("decode: %v", err)
	}
	data, ok := envelope.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("expected data to be map, got %T", envelope.Data)
	}
	if data["key"] != "value" {
		t.Errorf("expected key=value, got %v", data["key"])
	}
	if envelope.Error != nil {
		t.Error("expected no error in success response")
	}
}

func TestJSONWithMeta(t *testing.T) {
	w := httptest.NewRecorder()
	meta := &Meta{
		Total:      100,
		Page:       1,
		PageSize:   10,
		TotalPages: 10,
		RequestID:  "abc-123",
	}
	JSONWithMeta(w, http.StatusOK, []string{"item1", "item2"}, meta)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var envelope Envelope
	json.NewDecoder(w.Body).Decode(&envelope)

	if envelope.Meta == nil {
		t.Fatal("expected meta to be set")
	}
	if envelope.Meta.Total != 100 {
		t.Errorf("expected total 100, got %d", envelope.Meta.Total)
	}
	if envelope.Meta.RequestID != "abc-123" {
		t.Errorf("expected request_id abc-123, got %q", envelope.Meta.RequestID)
	}
}

func TestError_AppError(t *testing.T) {
	w := httptest.NewRecorder()
	appErr := apperr.NewValidation("name is required")
	Error(w, appErr)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", w.Code)
	}

	var envelope Envelope
	json.NewDecoder(w.Body).Decode(&envelope)

	if envelope.Error == nil {
		t.Fatal("expected error in response")
	}
	if envelope.Error.Code != "VALIDATION_ERROR" {
		t.Errorf("expected VALIDATION_ERROR, got %q", envelope.Error.Code)
	}
}

func TestError_NilError(t *testing.T) {
	w := httptest.NewRecorder()
	Error(w, nil)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 for nil error, got %d", w.Code)
	}
}

func TestError_GenericError(t *testing.T) {
	w := httptest.NewRecorder()
	Error(w, errors.New("something broke"))

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 for generic error, got %d", w.Code)
	}

	var envelope Envelope
	json.NewDecoder(w.Body).Decode(&envelope)

	if envelope.Error == nil {
		t.Fatal("expected error in response")
	}
	if envelope.Error.Message == "" {
		t.Error("expected non-empty error message")
	}
}

func TestError_NotFound(t *testing.T) {
	w := httptest.NewRecorder()
	Error(w, apperr.NewNotFound("Order"))

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestError_Unauthorized(t *testing.T) {
	w := httptest.NewRecorder()
	Error(w, apperr.NewUnauthorized("invalid token"))

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestError_Forbidden(t *testing.T) {
	w := httptest.NewRecorder()
	Error(w, apperr.NewForbidden(""))

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestPaginatedJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := []any{"a", "b", "c"}
	PaginatedJSON(w, data, 3, 1, 3)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var envelope Envelope
	json.NewDecoder(w.Body).Decode(&envelope)

	if envelope.Meta == nil {
		t.Fatal("expected meta")
	}
	if envelope.Meta.Total != 3 {
		t.Errorf("expected total 3, got %d", envelope.Meta.Total)
	}
	if envelope.Meta.TotalPages != 1 {
		t.Errorf("expected 1 total page, got %d", envelope.Meta.TotalPages)
	}
}

func TestPaginatedJSON_MultiplePages(t *testing.T) {
	w := httptest.NewRecorder()
	data := []any{"a", "b", "c", "d", "e"}
	PaginatedJSON(w, data, 5, 1, 2)

	var envelope Envelope
	json.NewDecoder(w.Body).Decode(&envelope)

	if envelope.Meta.TotalPages != 3 {
		t.Errorf("expected 3 total pages, got %d", envelope.Meta.TotalPages)
	}
}

func TestCreated(t *testing.T) {
	w := httptest.NewRecorder()
	Created(w, map[string]int{"id": 42})

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
}

func TestNoContent(t *testing.T) {
	w := httptest.NewRecorder()
	NoContent(w)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
}

func TestEnvelope_OmitEmpty(t *testing.T) {
	w := httptest.NewRecorder()
	JSON(w, http.StatusOK, nil)

	var envelope Envelope
	json.NewDecoder(w.Body).Decode(&envelope)

	// Data should be null/omitted when nil
	if envelope.Data != nil {
		t.Errorf("expected nil data, got %v", envelope.Data)
	}
	if envelope.Meta != nil {
		t.Errorf("expected nil meta, got %v", envelope.Meta)
	}
	if envelope.Error != nil {
		t.Errorf("expected nil error, got %v", envelope.Error)
	}
}
