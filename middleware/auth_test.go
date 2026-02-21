package middleware

import (
	"context"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"kanban-backend/models"
	"kanban-backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

type mockAuthServiceForAuth struct {
	validateTokenFunc func(token string) (string, error)
}

func (m *mockAuthServiceForAuth) Register(ctx context.Context, username, email, password string) (*models.User, error) {
	return nil, nil
}

func (m *mockAuthServiceForAuth) Login(ctx context.Context, email, password string) (string, error) {
	return "valid-token", nil
}

func (m *mockAuthServiceForAuth) GenerateToken(userID string, expiry time.Duration) (string, error) {
	return "valid-token-" + userID, nil
}

func (m *mockAuthServiceForAuth) ValidateToken(token string) (string, error) {
	if m.validateTokenFunc != nil {
		return m.validateTokenFunc(token)
	}
	return "user-123", nil
}

func (m *mockAuthServiceForAuth) HashPassword(password string) (string, error) {
	return "hashed-password", nil
}

func (m *mockAuthServiceForAuth) VerifyPassword(hashedPassword, password string) error {
	return nil
}

func TestAuthMiddleware_Success(t *testing.T) {
	app := fiber.New()

	mockAuthService := &mockAuthServiceForAuth{
		validateTokenFunc: func(token string) (string, error) {
			return "user-123", nil
		},
	}

	app.Use(AuthMiddleware(mockAuthService))
	app.Get("/protected", func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(string)
		return c.JSON(fiber.Map{"user_id": userID})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	respBody := string(body)

	assert.Contains(t, respBody, `"user_id":"user-123"`)
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	app := fiber.New()

	mockAuthService := &mockAuthServiceForAuth{}
	app.Use(AuthMiddleware(mockAuthService))
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/protected", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestAuthMiddleware_InvalidHeaderFormat(t *testing.T) {
	tests := []struct {
		name       string
		authHeader string
	}{
		{
			name:       "Missing Bearer prefix",
			authHeader: "invalid-token",
		},
		{
			name:       "Wrong prefix",
			authHeader: "Basic valid-token",
		},
		{
			name:       "Empty token",
			authHeader: "Bearer ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			mockAuthService := &mockAuthServiceForAuth{}
			app.Use(AuthMiddleware(mockAuthService))
			app.Get("/protected", func(c *fiber.Ctx) error {
				return c.JSON(fiber.Map{"message": "success"})
			})

			req := httptest.NewRequest("GET", "/protected", nil)
			req.Header.Set("Authorization", tt.authHeader)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
		})
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	app := fiber.New()

	mockAuthService := &mockAuthServiceForAuth{
		validateTokenFunc: func(token string) (string, error) {
			return "", utils.NewUnauthorized("invalid token")
		},
	}

	app.Use(AuthMiddleware(mockAuthService))
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}
