package unit

import (
	"testing"

	"github.com/i56/framework/core/validator"
)

func TestValidatorRequired(t *testing.T) {
	v := validator.New()
	v.Required("name", "Alice")
	if !v.Valid() {
		t.Error("expected valid when field is not empty")
	}

	v2 := validator.New()
	v2.Required("name", "")
	if v2.Valid() {
		t.Error("expected invalid when field is empty")
	}
}

func TestValidatorMaxLength(t *testing.T) {
	v := validator.New()
	v.MaxLength("name", "Alice", 10)
	if !v.Valid() {
		t.Error("expected valid when under max length")
	}

	v2 := validator.New()
	v2.MaxLength("name", "AliceAliceAlice", 10)
	if v2.Valid() {
		t.Error("expected invalid when over max length")
	}
}

func TestValidatorIn(t *testing.T) {
	v := validator.New()
	v.In("type", "platform", []string{"platform", "shopee", "major"})
	if !v.Valid() {
		t.Error("expected valid when value is in allowed list")
	}

	v2 := validator.New()
	v2.In("type", "invalid", []string{"platform", "shopee"})
	if v2.Valid() {
		t.Error("expected invalid when value is not in allowed list")
	}
}

func TestValidatorRange(t *testing.T) {
	v := validator.New()
	v.Range("age", 25, 0, 120)
	if !v.Valid() {
		t.Error("expected valid when in range")
	}

	v2 := validator.New()
	v2.Range("age", 150, 0, 120)
	if v2.Valid() {
		t.Error("expected invalid when out of range")
	}
}

func TestValidatorChain(t *testing.T) {
	input := ""
	v := validator.New()
	v.Required("name", input).
		MinLength("name", input, 2).
		MaxLength("name", input, 50)
	if v.Valid() {
		t.Error("expected invalid for empty input")
	}
	if len(v.Errors()) != 2 { // required + min_length
		t.Errorf("expected 2 errors, got %d: %v", len(v.Errors()), v.Errors())
	}
}

func TestValidatorToAppError(t *testing.T) {
	v := validator.New()
	v.Required("email", "")
	err := v.ToAppError()
	if err == nil {
		t.Error("expected app error")
	}
	if err.HTTPStatus != 422 {
		t.Errorf("expected 422, got %d", err.HTTPStatus)
	}
}
