package handlers

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

type mockTaskService struct {
	createFunc         func(ctx context.Context, userID, columnID, title, description string, deadline *time.Time) (*models.Task, error)
	findByIDFunc       func(ctx context.Context, taskID, userID string) (*models.Task, error)
	findByColumnIDFunc func(ctx context.Context, columnID, userID string) ([]*models.Task, error)
	updateFunc         func(ctx context.Context, taskID, userID, title, description string, deadline *time.Time) (*models.Task, error)
	deleteFunc         func(ctx context.Context, taskID, userID string) error
	moveFunc           func(ctx context.Context, taskID, columnID, userID string) error
}

func (m *mockTaskService) Create(ctx context.Context, userID, columnID, title, description string, deadline *time.Time) (*models.Task, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, userID, columnID, title, description, deadline)
	}
	task := &models.Task{
		ID:          "task-123",
		ColumnID:    columnID,
		Title:       title,
		Description: description,
		Deadline:    deadline,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return task, nil
}

func (m *mockTaskService) FindByID(ctx context.Context, taskID, userID string) (*models.Task, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, taskID, userID)
	}
	task := &models.Task{
		ID:          taskID,
		ColumnID:    "col-123",
		Title:       "Test Task",
		Description: "Test Description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return task, nil
}

func (m *mockTaskService) FindByColumnID(ctx context.Context, columnID, userID string) ([]*models.Task, error) {
	if m.findByColumnIDFunc != nil {
		return m.findByColumnIDFunc(ctx, columnID, userID)
	}
	tasks := []*models.Task{
		{
			ID:          "task-1",
			ColumnID:    columnID,
			Title:       "Task 1",
			Description: "Description 1",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "task-2",
			ColumnID:    columnID,
			Title:       "Task 2",
			Description: "Description 2",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
	return tasks, nil
}

func (m *mockTaskService) Update(ctx context.Context, taskID, userID, title, description string, deadline *time.Time) (*models.Task, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, taskID, userID, title, description, deadline)
	}
	task := &models.Task{
		ID:          taskID,
		ColumnID:    "col-123",
		Title:       title,
		Description: description,
		Deadline:    deadline,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return task, nil
}

func (m *mockTaskService) Delete(ctx context.Context, taskID, userID string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, taskID, userID)
	}
	return nil
}

func (m *mockTaskService) Move(ctx context.Context, taskID, columnID, userID string) error {
	if m.moveFunc != nil {
		return m.moveFunc(ctx, taskID, columnID, userID)
	}
	return nil
}

func TestNewTaskController(t *testing.T) {
	mockService := &mockTaskService{}
	ctrl := NewTaskController(mockService)

	assert.NotNil(t, ctrl)
	assert.Equal(t, mockService, ctrl.taskService)
}

func TestTaskController_Create_Success(t *testing.T) {
	app := fiber.New()

	mockService := &mockTaskService{
		createFunc: func(ctx context.Context, userID, columnID, title, description string, deadline *time.Time) (*models.Task, error) {
			task := &models.Task{
				ID:          "task-123",
				ColumnID:    columnID,
				Title:       title,
				Description: description,
				Deadline:    deadline,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			return task, nil
		},
	}

	ctrl := NewTaskController(mockService)
	app.Post("/tasks", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.Create(c)
	})

	reqBody := `{"column_id":"col-123","title":"My Task","description":"Task description"}`
	req := httptest.NewRequest("POST", "/tasks", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":true`)
	assert.Contains(t, respBody, `"id":"task-123"`)
	assert.Contains(t, respBody, `"title":"My Task"`)
	assert.Contains(t, respBody, `"description":"Task description"`)
}

func TestTaskController_Create_ValidationErrors(t *testing.T) {
	tests := []struct {
		name       string
		reqBody    string
		wantStatus int
		wantError  string
	}{
		{
			name:       "Missing title",
			reqBody:    `{"column_id":"col-123"}`,
			wantStatus: fiber.StatusBadRequest,
			wantError:  "title is required",
		},
		{
			name:       "Missing column_id",
			reqBody:    `{"title":"My Task"}`,
			wantStatus: fiber.StatusBadRequest,
			wantError:  "column_id is required",
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
			wantError:  "title is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			mockService := &mockTaskService{}
			ctrl := NewTaskController(mockService)
			app.Post("/tasks", func(c *fiber.Ctx) error {
				c.Locals("user_id", "user-123")
				return ctrl.Create(c)
			})

			req := httptest.NewRequest("POST", "/tasks", strings.NewReader(tt.reqBody))
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

func TestTaskController_Create_ServiceError(t *testing.T) {
	app := fiber.New()

	mockService := &mockTaskService{
		createFunc: func(ctx context.Context, userID, columnID, title, description string, deadline *time.Time) (*models.Task, error) {
			return nil, utils.NewValidation("task title must be at least 3 characters")
		},
	}

	ctrl := NewTaskController(mockService)
	app.Post("/tasks", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.Create(c)
	})

	reqBody := `{"column_id":"col-123","title":"My Task"}`
	req := httptest.NewRequest("POST", "/tasks", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":false`)
	assert.Contains(t, respBody, "task title must be at least 3 characters")
}

func TestTaskController_FindByID_Success(t *testing.T) {
	app := fiber.New()

	mockService := &mockTaskService{
		findByIDFunc: func(ctx context.Context, taskID, userID string) (*models.Task, error) {
			task := &models.Task{
				ID:          taskID,
				ColumnID:    "col-123",
				Title:       "Test Task",
				Description: "Test Description",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			return task, nil
		},
	}

	ctrl := NewTaskController(mockService)
	app.Get("/tasks/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.FindByID(c)
	})

	req := httptest.NewRequest("GET", "/tasks/task-123", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":true`)
	assert.Contains(t, respBody, `"id":"task-123"`)
	assert.Contains(t, respBody, `"title":"Test Task"`)
}

func TestTaskController_FindByID_NotFound(t *testing.T) {
	app := fiber.New()

	mockService := &mockTaskService{
		findByIDFunc: func(ctx context.Context, taskID, userID string) (*models.Task, error) {
			return nil, utils.NewNotFound("task not found")
		},
	}

	ctrl := NewTaskController(mockService)
	app.Get("/tasks/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.FindByID(c)
	})

	req := httptest.NewRequest("GET", "/tasks/task-999", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":false`)
	assert.Contains(t, respBody, "task not found")
}

func TestTaskController_FindByID_Unauthorized(t *testing.T) {
	app := fiber.New()

	mockService := &mockTaskService{
		findByIDFunc: func(ctx context.Context, taskID, userID string) (*models.Task, error) {
			return nil, utils.NewUnauthorized("you do not have access to this task")
		},
	}

	ctrl := NewTaskController(mockService)
	app.Get("/tasks/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-999")
		return ctrl.FindByID(c)
	})

	req := httptest.NewRequest("GET", "/tasks/task-123", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":false`)
	assert.Contains(t, respBody, "you do not have access to this task")
}

func TestTaskController_FindByColumnID_Success(t *testing.T) {
	app := fiber.New()

	mockService := &mockTaskService{
		findByColumnIDFunc: func(ctx context.Context, columnID, userID string) ([]*models.Task, error) {
			tasks := []*models.Task{
				{
					ID:          "task-1",
					ColumnID:    columnID,
					Title:       "Task 1",
					Description: "Description 1",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
				{
					ID:          "task-2",
					ColumnID:    columnID,
					Title:       "Task 2",
					Description: "Description 2",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			}
			return tasks, nil
		},
	}

	ctrl := NewTaskController(mockService)
	app.Get("/columns/:columnId/tasks", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.FindByColumnID(c)
	})

	req := httptest.NewRequest("GET", "/columns/col-123/tasks", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":true`)
	assert.Contains(t, respBody, `"id":"task-1"`)
	assert.Contains(t, respBody, `"id":"task-2"`)
	assert.Contains(t, respBody, `"title":"Task 1"`)
	assert.Contains(t, respBody, `"title":"Task 2"`)
}

func TestTaskController_FindByColumnID_NotFound(t *testing.T) {
	app := fiber.New()

	mockService := &mockTaskService{
		findByColumnIDFunc: func(ctx context.Context, columnID, userID string) ([]*models.Task, error) {
			return nil, utils.NewNotFound("column not found")
		},
	}

	ctrl := NewTaskController(mockService)
	app.Get("/columns/:columnId/tasks", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.FindByColumnID(c)
	})

	req := httptest.NewRequest("GET", "/columns/col-999/tasks", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":false`)
	assert.Contains(t, respBody, "column not found")
}

func TestTaskController_Update_Success(t *testing.T) {
	app := fiber.New()

	mockService := &mockTaskService{
		updateFunc: func(ctx context.Context, taskID, userID, title, description string, deadline *time.Time) (*models.Task, error) {
			task := &models.Task{
				ID:          taskID,
				ColumnID:    "col-123",
				Title:       title,
				Description: description,
				Deadline:    deadline,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			return task, nil
		},
	}

	ctrl := NewTaskController(mockService)
	app.Put("/tasks/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.Update(c)
	})

	reqBody := `{"title":"Updated Task","description":"Updated description"}`
	req := httptest.NewRequest("PUT", "/tasks/task-123", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":true`)
	assert.Contains(t, respBody, `"id":"task-123"`)
	assert.Contains(t, respBody, `"title":"Updated Task"`)
	assert.Contains(t, respBody, `"description":"Updated description"`)
}

func TestTaskController_Update_NotFound(t *testing.T) {
	app := fiber.New()

	mockService := &mockTaskService{
		updateFunc: func(ctx context.Context, taskID, userID, title, description string, deadline *time.Time) (*models.Task, error) {
			return nil, utils.NewNotFound("task not found")
		},
	}

	ctrl := NewTaskController(mockService)
	app.Put("/tasks/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.Update(c)
	})

	reqBody := `{"title":"Updated Task"}`
	req := httptest.NewRequest("PUT", "/tasks/task-999", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":false`)
	assert.Contains(t, respBody, "task not found")
}

func TestTaskController_Delete_Success(t *testing.T) {
	app := fiber.New()

	mockService := &mockTaskService{
		deleteFunc: func(ctx context.Context, taskID, userID string) error {
			return nil
		},
	}

	ctrl := NewTaskController(mockService)
	app.Delete("/tasks/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.Delete(c)
	})

	req := httptest.NewRequest("DELETE", "/tasks/task-123", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":true`)
	assert.Contains(t, respBody, `"message":"Task deleted successfully"`)
}

func TestTaskController_Delete_NotFound(t *testing.T) {
	app := fiber.New()

	mockService := &mockTaskService{
		deleteFunc: func(ctx context.Context, taskID, userID string) error {
			return utils.NewNotFound("task not found")
		},
	}

	ctrl := NewTaskController(mockService)
	app.Delete("/tasks/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.Delete(c)
	})

	req := httptest.NewRequest("DELETE", "/tasks/task-999", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":false`)
	assert.Contains(t, respBody, "task not found")
}

func TestTaskController_Move_Success(t *testing.T) {
	app := fiber.New()

	mockService := &mockTaskService{
		moveFunc: func(ctx context.Context, taskID, columnID, userID string) error {
			return nil
		},
	}

	ctrl := NewTaskController(mockService)
	app.Put("/tasks/:id/move", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.Move(c)
	})

	reqBody := `{"column_id":"col-456"}`
	req := httptest.NewRequest("PUT", "/tasks/task-123/move", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":true`)
	assert.Contains(t, respBody, `"message":"Task moved successfully"`)
}

func TestTaskController_Move_ValidationError(t *testing.T) {
	app := fiber.New()

	mockService := &mockTaskService{}

	ctrl := NewTaskController(mockService)
	app.Put("/tasks/:id/move", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.Move(c)
	})

	reqBody := `{}`
	req := httptest.NewRequest("PUT", "/tasks/task-123/move", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":false`)
	assert.Contains(t, respBody, "column_id is required")
}

func TestTaskController_Move_NotFound(t *testing.T) {
	app := fiber.New()

	mockService := &mockTaskService{
		moveFunc: func(ctx context.Context, taskID, columnID, userID string) error {
			return utils.NewNotFound("task not found")
		},
	}

	ctrl := NewTaskController(mockService)
	app.Put("/tasks/:id/move", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.Move(c)
	})

	reqBody := `{"column_id":"col-456"}`
	req := httptest.NewRequest("PUT", "/tasks/task-999/move", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":false`)
	assert.Contains(t, respBody, "task not found")
}

func TestTaskController_ServiceError(t *testing.T) {
	app := fiber.New()

	mockService := &mockTaskService{
		findByIDFunc: func(ctx context.Context, taskID, userID string) (*models.Task, error) {
			return nil, errors.New("database connection failed")
		},
	}

	ctrl := NewTaskController(mockService)
	app.Get("/tasks/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.FindByID(c)
	})

	req := httptest.NewRequest("GET", "/tasks/task-123", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":false`)
	assert.Contains(t, respBody, "Failed to find task")
}
