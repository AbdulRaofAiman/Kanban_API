package utils

import (
	"errors"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	fasthttp "github.com/valyala/fasthttp"
)

// Test Error Types implement error interface
func TestErrorTypes(t *testing.T) {
	tests := []struct {
		name       string
		error      error
		wantMsg    string
		wantCustom bool
	}{
		{
			name:       "ErrNotFound with default message",
			error:      ErrNotFound{},
			wantMsg:    "resource not found",
			wantCustom: true,
		},
		{
			name:       "ErrNotFound with custom message",
			error:      ErrNotFound{Message: "user not found"},
			wantMsg:    "user not found",
			wantCustom: true,
		},
		{
			name:       "ErrUnauthorized with default message",
			error:      ErrUnauthorized{},
			wantMsg:    "unauthorized access",
			wantCustom: true,
		},
		{
			name:       "ErrUnauthorized with custom message",
			error:      ErrUnauthorized{Message: "invalid credentials"},
			wantMsg:    "invalid credentials",
			wantCustom: true,
		},
		{
			name:       "ErrValidation with default message",
			error:      ErrValidation{},
			wantMsg:    "validation failed",
			wantCustom: true,
		},
		{
			name:       "ErrValidation with custom message",
			error:      ErrValidation{Message: "email is required"},
			wantMsg:    "email is required",
			wantCustom: true,
		},
		{
			name:       "ErrConflict with default message",
			error:      ErrConflict{},
			wantMsg:    "resource conflict",
			wantCustom: true,
		},
		{
			name:       "ErrConflict with custom message",
			error:      ErrConflict{Message: "email already exists"},
			wantMsg:    "email already exists",
			wantCustom: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that error implements error interface
			var _ error = tt.error

			// Test Error() method returns expected message
			if got := tt.error.Error(); got != tt.wantMsg {
				t.Errorf("Error() = %v, want %v", got, tt.wantMsg)
			}

			// Test type assertion works
			switch tt.error.(type) {
			case ErrNotFound, ErrUnauthorized, ErrValidation, ErrConflict:
				// Success - error is one of our custom types
			default:
				t.Errorf("error is not a custom error type")
			}
		})
	}
}

// Test helper functions
func TestNewErrorFunctions(t *testing.T) {
	tests := []struct {
		name     string
		fn       func(string) error
		msg      string
		wantType interface{}
	}{
		{
			name:     "NewNotFound",
			fn:       NewNotFound,
			msg:      "board not found",
			wantType: ErrNotFound{},
		},
		{
			name:     "NewUnauthorized",
			fn:       NewUnauthorized,
			msg:      "token expired",
			wantType: ErrUnauthorized{},
		},
		{
			name:     "NewValidation",
			fn:       NewValidation,
			msg:      "invalid password",
			wantType: ErrValidation{},
		},
		{
			name:     "NewConflict",
			fn:       NewConflict,
			msg:      "username taken",
			wantType: ErrConflict{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn(tt.msg)
			var _ error = err
			if err == nil {
				t.Errorf("expected error, got nil")
			}
		})
	}
}

// Test ErrorHandler maps errors to correct HTTP status codes
func TestErrorHandler(t *testing.T) {
	app := fiber.New()

	tests := []struct {
		name         string
		err          error
		wantStatus   int
		wantSuccess  string
		wantContains string
	}{
		{
			name:         "ErrNotFound returns 404",
			err:          ErrNotFound{Message: "resource not found"},
			wantStatus:   fiber.StatusNotFound,
			wantSuccess:  "false",
			wantContains: "resource not found",
		},
		{
			name:         "ErrUnauthorized returns 401",
			err:          ErrUnauthorized{Message: "unauthorized"},
			wantStatus:   fiber.StatusUnauthorized,
			wantSuccess:  "false",
			wantContains: "unauthorized",
		},
		{
			name:         "ErrValidation returns 400",
			err:          ErrValidation{Message: "validation error"},
			wantStatus:   fiber.StatusBadRequest,
			wantSuccess:  "false",
			wantContains: "validation error",
		},
		{
			name:         "ErrConflict returns 409",
			err:          ErrConflict{Message: "conflict"},
			wantStatus:   fiber.StatusConflict,
			wantSuccess:  "false",
			wantContains: "conflict",
		},
		{
			name:         "Generic error returns 500",
			err:          errors.New("generic error"),
			wantStatus:   fiber.StatusInternalServerError,
			wantSuccess:  "false",
			wantContains: "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test context
			ctx := &fasthttp.RequestCtx{}

			c := app.AcquireCtx(ctx)
			defer app.ReleaseCtx(c)

			// Call error handler
			err := ErrorHandler(c, tt.err)

			// Should return nil (Fiber handles the response)
			if err != nil {
				t.Errorf("ErrorHandler() returned error: %v", err)
			}

			// Check status code
			if got := c.Response().StatusCode(); got != tt.wantStatus {
				t.Errorf("StatusCode() = %v, want %v", got, tt.wantStatus)
			}

			// Check response body contains expected fields
			body := string(c.Response().Body())
			if !strings.Contains(body, `"success":false`) {
				t.Errorf("Response body missing success:false, got: %s", body)
			}
			if !strings.Contains(body, tt.wantContains) {
				t.Errorf("Response body missing expected text %q, got: %s", tt.wantContains, body)
			}
		})
	}
}
