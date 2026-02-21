package controllers

import (
	"context"
	"errors"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"kanban-backend/models"
	"kanban-backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

type mockAuthService struct {
	registerFunc func(username, email, password string) (*models.User, error)
	loginFunc    func(email, password string) (string, error)
}

func (m *mockAuthService) Register(ctx context.Context, username, email, password string) (*models.User, error) {
	if m.registerFunc != nil {
		return m.registerFunc(username, email, password)
	}
	user := &models.User{
		ID:        "user-123",
		Username:  username,
		Email:     email,
		Password:  "hashed-password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return user, nil
}

func (m *mockAuthService) Login(ctx context.Context, email, password string) (string, error) {
	if m.loginFunc != nil {
		return m.loginFunc(email, password)
	}
	return "test-token", nil
}

func (m *mockAuthService) GenerateToken(userID string, expiry time.Duration) (string, error) {
	return "test-token-" + userID, nil
}

func (m *mockAuthService) ValidateToken(tokenString string) (string, error) {
	return "user-123", nil
}

func (m *mockAuthService) HashPassword(password string) (string, error) {
	return "hashed-password", nil
}

func (m *mockAuthService) VerifyPassword(hashedPassword, password string) error {
	return nil
}

func TestNewAuthController(t *testing.T) {
	mockService := &mockAuthService{}
	ctrl := NewAuthController(mockService)

	assert.NotNil(t, ctrl)
	assert.Equal(t, mockService, ctrl.authService)
}

func TestAuthController_Register_Success(t *testing.T) {
	app := fiber.New()

	mockService := &mockAuthService{
		registerFunc: func(username, email, password string) (*models.User, error) {
			user := &models.User{
				ID:        "user-123",
				Username:  username,
				Email:     email,
				Password:  "hashed-password",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			return user, nil
		},
	}

	ctrl := NewAuthController(mockService)
	app.Post("/auth/register", ctrl.Register)

	reqBody := `{"username":"testuser","email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest("POST", "/auth/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":true`)
	assert.Contains(t, respBody, `"token"`)
	assert.Contains(t, respBody, `"user"`)
	assert.Contains(t, respBody, `"id":"user-123"`)
	assert.Contains(t, respBody, `"username":"testuser"`)
	assert.Contains(t, respBody, `"email":"test@example.com"`)
}

func TestAuthController_Register_ValidationErrors(t *testing.T) {
	tests := []struct {
		name       string
		reqBody    string
		wantStatus int
		wantError  string
	}{
		{
			name:       "Missing username",
			reqBody:    `{"email":"test@example.com","password":"password123"}`,
			wantStatus: fiber.StatusBadRequest,
			wantError:  "username is required",
		},
		{
			name:       "Missing email",
			reqBody:    `{"username":"testuser","password":"password123"}`,
			wantStatus: fiber.StatusBadRequest,
			wantError:  "email is required",
		},
		{
			name:       "Invalid email format",
			reqBody:    `{"username":"testuser","email":"invalidemail","password":"password123"}`,
			wantStatus: fiber.StatusBadRequest,
			wantError:  "invalid email format",
		},
		{
			name:       "Missing password",
			reqBody:    `{"username":"testuser","email":"test@example.com"}`,
			wantStatus: fiber.StatusBadRequest,
			wantError:  "password is required",
		},
		{
			name:       "Short password",
			reqBody:    `{"username":"testuser","email":"test@example.com","password":"short"}`,
			wantStatus: fiber.StatusBadRequest,
			wantError:  "password must be at least 8 characters long",
		},
		{
			name:       "Invalid JSON",
			reqBody:    `invalid json`,
			wantStatus: fiber.StatusBadRequest,
			wantError:  "Invalid request body",
		},
		{
			name:       "Empty body",
			reqBody:    `{}`,
			wantStatus: fiber.StatusBadRequest,
			wantError:  "username is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			mockService := &mockAuthService{}
			ctrl := NewAuthController(mockService)
			app.Post("/auth/register", ctrl.Register)

			req := httptest.NewRequest("POST", "/auth/register", strings.NewReader(tt.reqBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			body, _ := io.ReadAll(resp.Body)
			respBody := string(body)

			assert.Contains(t, respBody, `"success":false`)
			assert.Contains(t, respBody, tt.wantError)
		})
	}
}

func TestAuthController_Register_ConflictError(t *testing.T) {
	app := fiber.New()

	mockService := &mockAuthService{
		registerFunc: func(username, email, password string) (*models.User, error) {
			return nil, utils.NewConflict("user with this email already exists")
		},
	}

	ctrl := NewAuthController(mockService)
	app.Post("/auth/register", ctrl.Register)

	reqBody := `{"username":"testuser","email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest("POST", "/auth/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusConflict, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":false`)
	assert.Contains(t, respBody, "user with this email already exists")
}

func TestAuthController_Register_ValidationError(t *testing.T) {
	app := fiber.New()

	mockService := &mockAuthService{
		registerFunc: func(username, email, password string) (*models.User, error) {
			return nil, utils.NewValidation("password must be at least 8 characters long")
		},
	}

	ctrl := NewAuthController(mockService)
	app.Post("/auth/register", ctrl.Register)

	reqBody := `{"username":"testuser","email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest("POST", "/auth/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":false`)
	assert.Contains(t, respBody, "password must be at least 8 characters long")
}

func TestAuthController_Register_InternalServerError(t *testing.T) {
	app := fiber.New()

	mockService := &mockAuthService{
		registerFunc: func(username, email, password string) (*models.User, error) {
			return nil, errors.New("database connection failed")
		},
	}

	ctrl := NewAuthController(mockService)
	app.Post("/auth/register", ctrl.Register)

	reqBody := `{"username":"testuser","email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest("POST", "/auth/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":false`)
	assert.Contains(t, respBody, "Failed to register user")
}

func TestAuthController_Login_Success(t *testing.T) {
	app := fiber.New()

	mockService := &mockAuthService{
		loginFunc: func(email, password string) (string, error) {
			return "test-auth-token-123", nil
		},
	}

	ctrl := NewAuthController(mockService)
	app.Post("/auth/login", ctrl.Login)

	reqBody := `{"email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest("POST", "/auth/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":true`)
	assert.Contains(t, respBody, `"token":"test-auth-token-123"`)
}

func TestAuthController_Login_ValidationErrors(t *testing.T) {
	tests := []struct {
		name       string
		reqBody    string
		wantStatus int
		wantError  string
	}{
		{
			name:       "Missing email",
			reqBody:    `{"password":"password123"}`,
			wantStatus: fiber.StatusBadRequest,
			wantError:  "email is required",
		},
		{
			name:       "Invalid email format",
			reqBody:    `{"email":"invalidemail","password":"password123"}`,
			wantStatus: fiber.StatusBadRequest,
			wantError:  "invalid email format",
		},
		{
			name:       "Missing password",
			reqBody:    `{"email":"test@example.com"}`,
			wantStatus: fiber.StatusBadRequest,
			wantError:  "password is required",
		},
		{
			name:       "Invalid JSON",
			reqBody:    `invalid json`,
			wantStatus: fiber.StatusBadRequest,
			wantError:  "Invalid request body",
		},
		{
			name:       "Empty body",
			reqBody:    `{}`,
			wantStatus: fiber.StatusBadRequest,
			wantError:  "email is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			mockService := &mockAuthService{}
			ctrl := NewAuthController(mockService)
			app.Post("/auth/login", ctrl.Login)

			req := httptest.NewRequest("POST", "/auth/login", strings.NewReader(tt.reqBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			body, _ := io.ReadAll(resp.Body)
			respBody := string(body)

			assert.Contains(t, respBody, `"success":false`)
			assert.Contains(t, respBody, tt.wantError)
		})
	}
}

func TestAuthController_Login_UnauthorizedError(t *testing.T) {
	app := fiber.New()

	mockService := &mockAuthService{
		loginFunc: func(email, password string) (string, error) {
			return "", utils.NewUnauthorized("invalid email or password")
		},
	}

	ctrl := NewAuthController(mockService)
	app.Post("/auth/login", ctrl.Login)

	reqBody := `{"email":"test@example.com","password":"wrongpassword"}`
	req := httptest.NewRequest("POST", "/auth/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":false`)
	assert.Contains(t, respBody, "invalid email or password")
}

func TestAuthController_Login_InternalServerError(t *testing.T) {
	app := fiber.New()

	mockService := &mockAuthService{
		loginFunc: func(email, password string) (string, error) {
			return "", errors.New("database connection failed")
		},
	}

	ctrl := NewAuthController(mockService)
	app.Post("/auth/login", ctrl.Login)

	reqBody := `{"email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest("POST", "/auth/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":false`)
	assert.Contains(t, respBody, "Failed to login")
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"test@example.com", true},
		{"user.name@example.com", true},
		{"user+tag@example.com", true},
		{"user123@test-domain.co.uk", true},
		{"invalidemail", false},
		{"@example.com", false},
		{"user@", false},
		{"user@.com", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			result := isValidEmail(tt.email)
			assert.Equal(t, tt.valid, result)
		})
	}
}
