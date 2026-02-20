package middleware

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// Test Logger middleware function
func TestLogger(t *testing.T) {
	// Create test app
	app := fiber.New()

	// Register logger middleware
	app.Use(Logger())

	// Register test route
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{
			name:       "GET request",
			method:     "GET",
			path:       "/test",
			wantStatus: 200,
		},
		{
			name:       "POST request",
			method:     "POST",
			path:       "/test",
			wantStatus: 405, // Method not allowed
		},
		{
			name:       "Non-existent path",
			method:     "GET",
			path:       "/nonexistent",
			wantStatus: 404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test request
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			req.Header.Set("Content-Type", "application/json")

			// Execute request
			resp, err := app.Test(req)

			// Verify no error
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Verify status code
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("got status %d, want %d", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}

// Test request ID generation
func TestRequestID(t *testing.T) {
	// Create test app
	app := fiber.New()

	// Register logger middleware
	app.Use(Logger())

	// Register test route yang returns request ID
	app.Get("/test", func(c *fiber.Ctx) error {
		requestID := c.Locals("requestID")
		return c.JSON(fiber.Map{
			"request_id": requestID,
		})
	})

	// Create test request
	req, _ := http.NewRequest("GET", "/test", nil)

	// Execute request
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("got status %d, want 200", resp.StatusCode)
	}

	// Parse response
	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Verify request ID exists dan not empty
	requestID := result["request_id"]
	if requestID == "" {
		t.Error("request_id is empty")
	}

	if len(requestID) != 36 {
		t.Errorf("request_id length is %d, want 36", len(requestID))
	}
}
