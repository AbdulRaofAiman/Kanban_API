package utils

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestSuccessResponse(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		return Success(c, fiber.Map{"message": "test data"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, string(body), `"success":true`)
	assert.Contains(t, string(body), `"message":"test data"`)
}

func TestErrorResponse(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		return Error(c, "not found", fiber.StatusNotFound)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, string(body), `"success":false`)
	assert.Contains(t, string(body), `"message":"not found"`)
}

func TestValidationErrorResponse(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		return ValidationError(c, "email", "invalid email format")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, string(body), `"success":false`)
	assert.Contains(t, string(body), `"field":"email"`)
	assert.Contains(t, string(body), `"message":"invalid email format"`)
}

func TestAuthErrorResponse(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		return AuthError(c, "unauthorized access")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, string(body), `"success":false`)
	assert.Contains(t, string(body), `"message":"unauthorized access"`)
}
