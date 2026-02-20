package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "valid email",
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "valid email with subdomain",
			email:   "user@mail.example.com",
			wantErr: false,
		},
		{
			name:    "valid email with plus sign",
			email:   "user+tag@example.com",
			wantErr: false,
		},
		{
			name:    "empty email",
			email:   "",
			wantErr: true,
		},
		{
			name:    "invalid email - no @",
			email:   "invalidemail.com",
			wantErr: true,
		},
		{
			name:    "invalid email - no domain",
			email:   "test@",
			wantErr: true,
		},
		{
			name:    "invalid email - no local part",
			email:   "@example.com",
			wantErr: true,
		},
		{
			name:    "invalid email - spaces",
			email:   "test @example.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateUUID(t *testing.T) {
	tests := []struct {
		name    string
		uuid    string
		wantErr bool
	}{
		{
			name:    "valid UUID v4",
			uuid:    "550e8400-e29b-41d4-a716-446655440000",
			wantErr: false,
		},
		{
			name:    "valid UUID v4 lowercase",
			uuid:    "550e8400-e29b-41d4-a716-446655440000",
			wantErr: false,
		},
		{
			name:    "valid UUID v4 uppercase",
			uuid:    "550E8400-E29B-41D4-A716-446655440000",
			wantErr: false,
		},
		{
			name:    "valid UUID v4 without dashes",
			uuid:    "550e8400e29b41d4a716446655440000",
			wantErr: false,
		},
		{
			name:    "empty UUID",
			uuid:    "",
			wantErr: true,
		},
		{
			name:    "invalid UUID - wrong length",
			uuid:    "550e8400-e29b-41d4-a716",
			wantErr: true,
		},
		{
			name:    "invalid UUID - invalid characters",
			uuid:    "550e8400-e29b-41d4-a716-44665544000g",
			wantErr: true,
		},
		{
			name:    "invalid UUID - random string",
			uuid:    "not-a-uuid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUUID(tt.uuid)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateStruct(t *testing.T) {
	// Define test struct with validation tags
	type TestStruct struct {
		Email    string `validate:"required,email"`
		Username string `validate:"required,min=3,max=20"`
		Age      int    `validate:"required,gte=0,lte=130"`
	}

	t.Run("valid struct", func(t *testing.T) {
		ts := TestStruct{
			Email:    "test@example.com",
			Username: "testuser",
			Age:      25,
		}

		err := ValidateStruct(ts)
		assert.NoError(t, err)
	})

	t.Run("invalid email", func(t *testing.T) {
		ts := TestStruct{
			Email:    "invalidemail",
			Username: "testuser",
			Age:      25,
		}

		err := ValidateStruct(ts)
		require.Error(t, err)

		fieldErrors, ok := err.(*FieldValidationErrors)
		require.True(t, ok)
		require.Len(t, fieldErrors.Errors, 1)
		assert.Equal(t, "Email", fieldErrors.Errors[0].Field)
		assert.Contains(t, fieldErrors.Errors[0].Message, "email")
	})

	t.Run("multiple validation errors", func(t *testing.T) {
		ts := TestStruct{
			Email:    "invalid",
			Username: "ab",
			Age:      150,
		}

		err := ValidateStruct(ts)
		require.Error(t, err)

		fieldErrors, ok := err.(*FieldValidationErrors)
		require.True(t, ok)
		require.Len(t, fieldErrors.Errors, 3)

		// Check all errors are present
		fields := make(map[string]bool)
		for _, fe := range fieldErrors.Errors {
			fields[fe.Field] = true
		}
		assert.True(t, fields["Email"])
		assert.True(t, fields["Username"])
		assert.True(t, fields["Age"])
	})

	t.Run("missing required field", func(t *testing.T) {
		ts := TestStruct{
			Email:    "",
			Username: "testuser",
			Age:      25,
		}

		err := ValidateStruct(ts)
		require.Error(t, err)

		fieldErrors, ok := err.(*FieldValidationErrors)
		require.True(t, ok)
		require.Len(t, fieldErrors.Errors, 1)
		assert.Equal(t, "Email", fieldErrors.Errors[0].Field)
		assert.Contains(t, fieldErrors.Errors[0].Message, "required")
	})

	t.Run("nil struct", func(t *testing.T) {
		err := ValidateStruct(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be nil")
	})
}

func TestFieldValidationErrors_Error(t *testing.T) {
	t.Run("single error", func(t *testing.T) {
		errors := FieldValidationErrors{
			Errors: []FieldValidationError{
				{Field: "Email", Message: "is required"},
			},
		}

		errStr := errors.Error()
		assert.Contains(t, errStr, "validation errors")
		assert.Contains(t, errStr, "Email")
		assert.Contains(t, errStr, "is required")
	})

	t.Run("multiple errors", func(t *testing.T) {
		errors := FieldValidationErrors{
			Errors: []FieldValidationError{
				{Field: "Email", Message: "is required"},
				{Field: "Username", Message: "must be at least 3 characters"},
			},
		}

		errStr := errors.Error()
		assert.Contains(t, errStr, "validation errors")
		assert.Contains(t, errStr, "Email")
		assert.Contains(t, errStr, "Username")
	})

	t.Run("no errors", func(t *testing.T) {
		errors := FieldValidationErrors{
			Errors: []FieldValidationError{},
		}

		errStr := errors.Error()
		assert.Equal(t, "validation failed", errStr)
	})
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{"valid email", "test@example.com", true},
		{"invalid email", "invalidemail", false},
		{"empty email", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidEmail(tt.email)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestIsValidUUID(t *testing.T) {
	tests := []struct {
		name string
		uuid string
		want bool
	}{
		{"valid UUID", "550e8400-e29b-41d4-a716-446655440000", true},
		{"invalid UUID", "not-a-uuid", false},
		{"empty UUID", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidUUID(tt.uuid)
			assert.Equal(t, tt.want, result)
		})
	}
}
