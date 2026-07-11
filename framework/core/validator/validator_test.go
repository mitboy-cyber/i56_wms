package validator

import (
	"errors"
	"testing"
)

func TestV_Required(t *testing.T) {
	v := New().Required("name", "")
	if v.Valid() {
		t.Error("expected invalid for empty required field")
	}
	if !v.Valid() {
		errs := v.Errors()
		if len(errs) != 1 || errs[0].Field != "name" {
			t.Errorf("unexpected errors: %+v", errs)
		}
	}

	v2 := New().Required("name", "John")
	if !v2.Valid() {
		t.Error("expected valid for non-empty required field")
	}
}

func TestV_MaxLength(t *testing.T) {
	v := New().MaxLength("code", "abcdefghij", 5)
	if v.Valid() {
		t.Error("expected invalid for exceeding max length")
	}

	v2 := New().MaxLength("code", "abc", 5)
	if !v2.Valid() {
		t.Error("expected valid for within max length")
	}

	v3 := New().MaxLength("code", "abcde", 5)
	if !v3.Valid() {
		t.Error("expected valid for exact max length")
	}
}

func TestV_MinLength(t *testing.T) {
	v := New().MinLength("code", "ab", 3)
	if v.Valid() {
		t.Error("expected invalid for below min length")
	}

	v2 := New().MinLength("code", "abcd", 3)
	if !v2.Valid() {
		t.Error("expected valid for meeting min length")
	}
}

func TestV_Email(t *testing.T) {
	valid := []string{"test@example.com", "user.name+tag@domain.co.jp"}
	for _, email := range valid {
		v := New().Email("email", email)
		if !v.Valid() {
			t.Errorf("expected valid email: %q", email)
		}
	}

	invalid := []string{"not-email", "@missing", "no-domain@"}
	for _, email := range invalid {
		v := New().Email("email", email)
		if !v.Valid() {
			// Empty is allowed (not required)
			continue
		}
		t.Errorf("expected invalid email: %q", email)
	}
}

func TestV_In(t *testing.T) {
	v := New().In("status", "active", []string{"active", "inactive", "pending"})
	if !v.Valid() {
		t.Error("expected valid for allowed value")
	}

	v2 := New().In("status", "deleted", []string{"active", "inactive"})
	if v2.Valid() {
		t.Error("expected invalid for disallowed value")
	}

	v3 := New().In("status", "", []string{"active"})
	if !v3.Valid() {
		t.Error("empty value should be valid (not required)")
	}
}

func TestV_Range(t *testing.T) {
	v := New().Range("age", 25, 0, 120)
	if !v.Valid() {
		t.Error("expected valid within range")
	}

	v2 := New().Range("age", -5, 0, 120)
	if v2.Valid() {
		t.Error("expected invalid below range")
	}

	v3 := New().Range("age", 150, 0, 120)
	if v3.Valid() {
		t.Error("expected invalid above range")
	}

	v4 := New().Range("age", 0, 0, 120)
	if !v4.Valid() {
		t.Error("expected valid at boundary")
	}
}

func TestV_Custom(t *testing.T) {
	v := New().Custom("custom", func() error {
		return errors.New("custom error")
	})
	if v.Valid() {
		t.Error("expected invalid for custom error")
	}

	v2 := New().Custom("custom", func() error {
		return nil
	})
	if !v2.Valid() {
		t.Error("expected valid for custom OK")
	}
}

func TestV_Chaining(t *testing.T) {
	v := New().
		Required("name", "John").
		MinLength("name", "John", 3).
		MaxLength("name", "John", 100).
		Email("email", "john@example.com").
		In("status", "active", []string{"active", "inactive"})

	if !v.Valid() {
		t.Error("expected all validations to pass")
	}
}

func TestV_ToAppError(t *testing.T) {
	v := New().
		Required("name", "").
		Required("email", "")

	err := v.ToAppError()
	if err == nil {
		t.Error("expected AppError")
	}
	if len(err.Details) != 2 {
		t.Errorf("expected 2 details, got %d", len(err.Details))
	}
	if err.Code != "VALIDATION_ERROR" {
		t.Errorf("expected VALIDATION_ERROR, got %q", err.Code)
	}
}
