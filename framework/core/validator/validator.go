// Package validator provides chainable request validation.
package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/i56/framework/core/errors"
)

// V is a chainable validator for request data.
type V struct {
	errs []errors.ErrorDetail
}

// New creates a new validator.
func New() *V {
	return &V{}
}

// Required checks that a value is not empty.
func (v *V) Required(field, value string) *V {
	if strings.TrimSpace(value) == "" {
		v.errs = append(v.errs, errors.ErrorDetail{
			Field:   field,
			Message: fmt.Sprintf("%s is required", field),
			Code:    "required",
		})
	}
	return v
}

// MaxLength checks string max length.
func (v *V) MaxLength(field, value string, max int) *V {
	if len(value) > max {
		v.errs = append(v.errs, errors.ErrorDetail{
			Field:   field,
			Message: fmt.Sprintf("%s must be at most %d characters", field, max),
			Code:    "max_length",
		})
	}
	return v
}

// MinLength checks string min length.
func (v *V) MinLength(field, value string, min int) *V {
	if len(value) < min {
		v.errs = append(v.errs, errors.ErrorDetail{
			Field:   field,
			Message: fmt.Sprintf("%s must be at least %d characters", field, min),
			Code:    "min_length",
		})
	}
	return v
}

// Email checks a valid email format.
func (v *V) Email(field, value string) *V {
	if value == "" {
		return v
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		v.errs = append(v.errs, errors.ErrorDetail{
			Field:   field,
			Message: fmt.Sprintf("%s must be a valid email address", field),
			Code:    "invalid_email",
		})
	}
	return v
}

// In checks that a value is in a list of allowed values.
func (v *V) In(field, value string, allowed []string) *V {
	if value == "" {
		return v
	}
	for _, a := range allowed {
		if value == a {
			return v
		}
	}
	v.errs = append(v.errs, errors.ErrorDetail{
		Field:   field,
		Message: fmt.Sprintf("%s must be one of: %s", field, strings.Join(allowed, ", ")),
		Code:    "invalid_value",
	})
	return v
}

// Custom adds a custom validation check.
func (v *V) Custom(field string, fn func() error) *V {
	if err := fn(); err != nil {
		v.errs = append(v.errs, errors.ErrorDetail{
			Field:   field,
			Message: err.Error(),
			Code:    "custom",
		})
	}
	return v
}

// Range checks an integer range.
func (v *V) Range(field string, value, min, max int) *V {
	if value < min || value > max {
		v.errs = append(v.errs, errors.ErrorDetail{
			Field:   field,
			Message: fmt.Sprintf("%s must be between %d and %d", field, min, max),
			Code:    "out_of_range",
		})
	}
	return v
}

// Valid returns true if no validation errors.
func (v *V) Valid() bool {
	return len(v.errs) == 0
}

// Errors returns collected validation errors.
func (v *V) Errors() []errors.ErrorDetail {
	return v.errs
}

// ToAppError converts validation errors to an AppError.
func (v *V) ToAppError() *errors.AppError {
	err := errors.NewValidation("validation failed")
	err.Details = v.errs
	return err
}
