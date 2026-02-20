package utils

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

// FieldValidationError represents a field-level validation error
type FieldValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidateEmail validates an email address format
func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	// Use validator package's built-in email validation
	err := validate.Var(email, "email")
	if err != nil {
		return fmt.Errorf("invalid email format: %s", email)
	}

	return nil
}

// ValidateUUID validates a UUID string format
func ValidateUUID(uuidStr string) error {
	if uuidStr == "" {
		return errors.New("uuid is required")
	}

	// Parse the UUID to validate its format
	_, err := uuid.Parse(uuidStr)
	if err != nil {
		return fmt.Errorf("invalid uuid format: %s", uuidStr)
	}

	return nil
}

// ValidateStruct validates a struct using struct tags and returns field-level errors
func ValidateStruct(v interface{}) error {
	if v == nil {
		return errors.New("struct to validate cannot be nil")
	}

	err := validate.Struct(v)
	if err == nil {
		return nil
	}

	// Type assertion to get validation errors
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return fmt.Errorf("validation failed: %w", err)
	}

	var errors []FieldValidationError
	for _, fieldError := range validationErrors {
		errors = append(errors, FieldValidationError{
			Field:   fieldError.Field(),
			Message: getValidationErrorMessage(fieldError),
		})
	}

	return &FieldValidationErrors{Errors: errors}
}

// FieldValidationErrors is a collection of field-level validation errors
type FieldValidationErrors struct {
	Errors []FieldValidationError
}

// Error implements the error interface
func (ve *FieldValidationErrors) Error() string {
	if len(ve.Errors) == 0 {
		return "validation failed"
	}

	msg := "validation errors:"
	for _, err := range ve.Errors {
		msg += fmt.Sprintf(" %s: %s;", err.Field, err.Message)
	}
	return msg
}

// getValidationErrorMessage returns a human-readable error message for a validation error
func getValidationErrorMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return fmt.Sprintf("must be at least %s characters", fieldError.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters", fieldError.Param())
	case "len":
		return fmt.Sprintf("must be %s characters", fieldError.Param())
	case "uuid":
		return "must be a valid UUID"
	case "oneof":
		return fmt.Sprintf("must be one of: %s", fieldError.Param())
	default:
		return fmt.Sprintf("failed on '%s' validation", fieldError.Tag())
	}
}

// IsValidEmail checks if an email is valid without returning an error
func IsValidEmail(email string) bool {
	return ValidateEmail(email) == nil
}

// IsValidUUID checks if a UUID is valid without returning an error
func IsValidUUID(uuidStr string) bool {
	return ValidateUUID(uuidStr) == nil
}
