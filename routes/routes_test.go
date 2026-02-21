package routes

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"kanban-backend/handlers"
	"kanban-backend/models"
	"kanban-backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

type MockAuthService struct{}

func (m *MockAuthService) Register(ctx context.Context, username, email, password string) (*models.User, error) {
	return &models.User{ID: "user-1", Username: username, Email: email}, nil
}

func (m *MockAuthService) Login(ctx context.Context, email, password string) (string, error) {
	return "mock-jwt-token", nil
}

func (m *MockAuthService) GenerateToken(userID string, expiry time.Duration) (string, error) {
	return "mock-jwt-token", nil
}

func (m *MockAuthService) ValidateToken(tokenString string) (string, error) {
	if tokenString == "mock-jwt-token" {
		return "user-1", nil
	}
	return "", utils.NewUnauthorized("invalid or expired token")
}

func (m *MockAuthService) HashPassword(password string) (string, error) {
	return "hashed-password", nil
}

func (m *MockAuthService) VerifyPassword(hashedPassword, password string) error {
	return nil
}

type MockBoardService struct{}

func (m *MockBoardService) Create(ctx context.Context, userID, title, color string) (*models.Board, error) {
	return &models.Board{ID: "board-1", Title: title, Color: color, UserID: userID}, nil
}

func (m *MockBoardService) FindByID(ctx context.Context, boardID, userID string) (*models.Board, error) {
	if boardID == "board-1" {
		return &models.Board{ID: boardID, Title: "Test Board", UserID: userID}, nil
	}
	return nil, utils.NewNotFound("board not found")
}

func (m *MockBoardService) FindByUserID(ctx context.Context, userID string) ([]*models.Board, error) {
	return []*models.Board{
		{ID: "board-1", Title: "Test Board", UserID: userID},
	}, nil
}

func (m *MockBoardService) Update(ctx context.Context, boardID, userID, title, color string) (*models.Board, error) {
	if boardID == "board-1" {
		return &models.Board{ID: boardID, Title: title, Color: color, UserID: userID}, nil
	}
	return nil, utils.NewNotFound("board not found")
}

func (m *MockBoardService) Delete(ctx context.Context, boardID, userID string) error {
	if boardID == "board-1" {
		return nil
	}
	return utils.NewNotFound("board not found")
}

type MockTaskService struct{}

func (m *MockTaskService) Create(ctx context.Context, userID, columnID, title, description string, deadline *time.Time) (*models.Task, error) {
	return &models.Task{ID: "task-1", ColumnID: columnID, Title: title, Description: description}, nil
}

func (m *MockTaskService) FindByID(ctx context.Context, taskID, userID string) (*models.Task, error) {
	if taskID == "task-1" {
		return &models.Task{ID: taskID, Title: "Test Task"}, nil
	}
	return nil, utils.NewNotFound("task not found")
}

func (m *MockTaskService) FindByColumnID(ctx context.Context, columnID, userID string) ([]*models.Task, error) {
	if columnID == "column-1" {
		return []*models.Task{
			{ID: "task-1", ColumnID: columnID, Title: "Test Task"},
		}, nil
	}
	return nil, utils.NewNotFound("column not found")
}

func (m *MockTaskService) Update(ctx context.Context, taskID, userID, title, description string, deadline *time.Time) (*models.Task, error) {
	if taskID == "task-1" {
		return &models.Task{ID: taskID, Title: title, Description: description}, nil
	}
	return nil, utils.NewNotFound("task not found")
}

func (m *MockTaskService) Delete(ctx context.Context, taskID, userID string) error {
	if taskID == "task-1" {
		return nil
	}
	return utils.NewNotFound("task not found")
}

func (m *MockTaskService) Move(ctx context.Context, taskID, columnID, userID string) error {
	if taskID == "task-1" {
		return nil
	}
	return utils.NewNotFound("task not found")
}

func setupApp() *fiber.App {
	app := fiber.New()

	mockAuthService := &MockAuthService{}
	mockBoardService := &MockBoardService{}
	mockTaskService := &MockTaskService{}

	authController := handlers.NewAuthController(mockAuthService)
	boardController := handlers.NewBoardController(mockBoardService)
	taskController := handlers.NewTaskController(mockTaskService)

	Setup(app, mockAuthService, authController, boardController, taskController)

	return app
}

func TestHealthCheck(t *testing.T) {
	app := setupApp()

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestAuthRegister(t *testing.T) {
	app := setupApp()

	payload := `{"username":"testuser","email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestAuthLogin(t *testing.T) {
	app := setupApp()

	payload := `{"email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestBoardCreate_WithoutToken(t *testing.T) {
	app := setupApp()

	payload := `{"title":"Test Board","color":"#FF5733"}`
	req := httptest.NewRequest("POST", "/api/v1/boards", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestBoardCreate_WithValidToken(t *testing.T) {
	app := setupApp()

	payload := `{"title":"Test Board","color":"#FF5733"}`
	req := httptest.NewRequest("POST", "/api/v1/boards", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer mock-jwt-token")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestBoardFindAll_WithoutToken(t *testing.T) {
	app := setupApp()

	req := httptest.NewRequest("GET", "/api/v1/boards", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestBoardFindAll_WithValidToken(t *testing.T) {
	app := setupApp()

	req := httptest.NewRequest("GET", "/api/v1/boards", nil)
	req.Header.Set("Authorization", "Bearer mock-jwt-token")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestBoardFindByID_WithValidToken(t *testing.T) {
	app := setupApp()

	req := httptest.NewRequest("GET", "/api/v1/boards/board-1", nil)
	req.Header.Set("Authorization", "Bearer mock-jwt-token")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestBoardUpdate_WithValidToken(t *testing.T) {
	app := setupApp()

	payload := `{"title":"Updated Board","color":"#00FF00"}`
	req := httptest.NewRequest("PUT", "/api/v1/boards/board-1", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer mock-jwt-token")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestBoardDelete_WithValidToken(t *testing.T) {
	app := setupApp()

	req := httptest.NewRequest("DELETE", "/api/v1/boards/board-1", nil)
	req.Header.Set("Authorization", "Bearer mock-jwt-token")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestTaskCreate_WithoutToken(t *testing.T) {
	app := setupApp()

	payload := `{"title":"Test Task","column_id":"column-1"}`
	req := httptest.NewRequest("POST", "/api/v1/tasks", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestTaskCreate_WithValidToken(t *testing.T) {
	app := setupApp()

	payload := `{"title":"Test Task","column_id":"column-1"}`
	req := httptest.NewRequest("POST", "/api/v1/tasks", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer mock-jwt-token")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestTaskFindByID_WithValidToken(t *testing.T) {
	app := setupApp()

	req := httptest.NewRequest("GET", "/api/v1/tasks/task-1", nil)
	req.Header.Set("Authorization", "Bearer mock-jwt-token")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestTaskFindByColumnID_WithValidToken(t *testing.T) {
	app := setupApp()

	req := httptest.NewRequest("GET", "/api/v1/tasks/column/column-1", nil)
	req.Header.Set("Authorization", "Bearer mock-jwt-token")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestTaskUpdate_WithValidToken(t *testing.T) {
	app := setupApp()

	payload := `{"title":"Updated Task"}`
	req := httptest.NewRequest("PUT", "/api/v1/tasks/task-1", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer mock-jwt-token")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestTaskDelete_WithValidToken(t *testing.T) {
	app := setupApp()

	req := httptest.NewRequest("DELETE", "/api/v1/tasks/task-1", nil)
	req.Header.Set("Authorization", "Bearer mock-jwt-token")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestTaskMove_WithValidToken(t *testing.T) {
	app := setupApp()

	payload := `{"column_id":"column-2"}`
	req := httptest.NewRequest("PUT", "/api/v1/tasks/task-1/move", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer mock-jwt-token")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestInvalidRoute(t *testing.T) {
	app := setupApp()

	req := httptest.NewRequest("GET", "/api/v1/invalid", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestInvalidMethod(t *testing.T) {
	app := setupApp()

	req := httptest.NewRequest("POST", "/api/v1/boards", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}
