package middleware

import (
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func TestCORSConfig_Default(t *testing.T) {
	// Clear CORS_ALLOWED_ORIGINS environment variable to test default
	os.Unsetenv("CORS_ALLOWED_ORIGINS")

	config := CORSConfig()

	if config.AllowOrigins != "*" {
		t.Errorf("Expected AllowOrigins to be '*', got %s", config.AllowOrigins)
	}
	if config.AllowMethods != "GET, POST, PUT, DELETE, PATCH, OPTIONS" {
		t.Errorf("Expected AllowMethods to be 'GET, POST, PUT, DELETE, PATCH, OPTIONS', got %s", config.AllowMethods)
	}
	if config.AllowHeaders != "Origin, Content-Type, Authorization, Accept, X-Requested-With" {
		t.Errorf("Expected AllowHeaders to be 'Origin, Content-Type, Authorization, Accept, X-Requested-With', got %s", config.AllowHeaders)
	}
	// AllowCredentials must be false when using "*" (CORS security requirement)
	if config.AllowCredentials {
		t.Error("Expected AllowCredentials to be false when using wildcard origin")
	}
	if config.MaxAge != 86400 {
		t.Errorf("Expected MaxAge to be 86400, got %d", config.MaxAge)
	}
}

func TestCORSConfig_WithCustomOrigins(t *testing.T) {
	// Set custom origins
	os.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://example.com")
	defer os.Unsetenv("CORS_ALLOWED_ORIGINS")

	config := CORSConfig()

	if config.AllowOrigins != "http://localhost:3000,http://example.com" {
		t.Errorf("Expected AllowOrigins to be 'http://localhost:3000,http://example.com', got %s", config.AllowOrigins)
	}
	if config.AllowMethods != "GET, POST, PUT, DELETE, PATCH, OPTIONS" {
		t.Errorf("Expected AllowMethods to be 'GET, POST, PUT, DELETE, PATCH, OPTIONS', got %s", config.AllowMethods)
	}
	if config.AllowHeaders != "Origin, Content-Type, Authorization, Accept, X-Requested-With" {
		t.Errorf("Expected AllowHeaders to be 'Origin, Content-Type, Authorization, Accept, X-Requested-With', got %s", config.AllowHeaders)
	}
	if !config.AllowCredentials {
		t.Error("Expected AllowCredentials to be true")
	}
}

func TestCORSHeaders(t *testing.T) {
	config := CORSConfig()

	if !strings.Contains(config.AllowMethods, "GET") {
		t.Error("Expected GET in AllowMethods")
	}
	if !strings.Contains(config.AllowMethods, "POST") {
		t.Error("Expected POST in AllowMethods")
	}
	if !strings.Contains(config.AllowMethods, "PUT") {
		t.Error("Expected PUT in AllowMethods")
	}
	if !strings.Contains(config.AllowMethods, "DELETE") {
		t.Error("Expected DELETE in AllowMethods")
	}
	if !strings.Contains(config.AllowMethods, "PATCH") {
		t.Error("Expected PATCH in AllowMethods")
	}
	if !strings.Contains(config.AllowMethods, "OPTIONS") {
		t.Error("Expected OPTIONS in AllowMethods")
	}

	if !strings.Contains(config.AllowHeaders, "Authorization") {
		t.Error("Expected Authorization in AllowHeaders")
	}
	if !strings.Contains(config.AllowHeaders, "Content-Type") {
		t.Error("Expected Content-Type in AllowHeaders")
	}
	if !strings.Contains(config.AllowHeaders, "Origin") {
		t.Error("Expected Origin in AllowHeaders")
	}
}

func TestCORSActualRequest(t *testing.T) {
	app := fiber.New()
	app.Use(cors.New(CORSConfig()))

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("test")
	})

	// Test actual GET request
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	// Check CORS headers are present
	if resp.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("Expected Access-Control-Allow-Origin to be '*', got %s", resp.Header.Get("Access-Control-Allow-Origin"))
	}
}
