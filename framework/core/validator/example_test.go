package validator_test

import (
	"fmt"

	"github.com/i56/framework/core/validator"
)

// ExampleV demonstrates chainable request validation.
func ExampleV() {
	v := validator.New().
		Required("name", "Alice").
		MinLength("name", "Alice", 2).
		MaxLength("name", "Alice", 100).
		Email("email", "alice@example.com").
		In("status", "active", []string{"active", "inactive", "pending"})

	if !v.Valid() {
		for _, e := range v.Errors() {
			fmt.Printf("  %s: %s\n", e.Field, e.Message)
		}
		return
	}

	fmt.Println("All validations passed!")

	// Example with errors
	v2 := validator.New().
		Required("name", "").
		Email("email", "not-an-email")

	if !v2.Valid() {
		fmt.Printf("Found %d validation errors\n", len(v2.Errors()))
	}
	// Output:
	// All validations passed!
	// Found 2 validation errors
}

// ExampleV_ToAppError demonstrates converting to API error.
func ExampleV_ToAppError() {
	v := validator.New().
		Required("username", "").
		Required("password", "")

	if !v.Valid() {
		appErr := v.ToAppError()
		fmt.Println("Code:", appErr.Code)
		fmt.Println("Fields:", len(appErr.Details))
	}
	// Output:
	// Code: VALIDATION_ERROR
	// Fields: 2
}
