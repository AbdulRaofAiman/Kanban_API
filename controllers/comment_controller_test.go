package controllers

import (
	"context"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"kanban-backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

type mockCommentService struct {
	createFunc       func(ctx context.Context, taskID, userID, content string) (*models.Comment, error)
	findByIDFunc     func(ctx context.Context, id, userID string) (*models.Comment, error)
	findByTaskIDFunc func(ctx context.Context, taskID, userID string) ([]*models.Comment, error)
	updateFunc       func(ctx context.Context, id, userID, content string) (*models.Comment, error)
	deleteFunc       func(ctx context.Context, id, userID string) error
}

func (m *mockCommentService) Create(ctx context.Context, taskID, userID, content string) (*models.Comment, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, taskID, userID, content)
	}
	return &models.Comment{
		ID:      "comment-1",
		TaskID:  taskID,
		UserID:  userID,
		Content: content,
	}, nil
}

func (m *mockCommentService) FindByID(ctx context.Context, id, userID string) (*models.Comment, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id, userID)
	}
	return &models.Comment{
		ID:      id,
		Content: "Test Comment",
	}, nil
}

func (m *mockCommentService) FindByTaskID(ctx context.Context, taskID, userID string) ([]*models.Comment, error) {
	if m.findByTaskIDFunc != nil {
		return m.findByTaskIDFunc(ctx, taskID, userID)
	}
	return []*models.Comment{
		{ID: "comment-1", TaskID: taskID, Content: "Comment 1"},
		{ID: "comment-2", TaskID: taskID, Content: "Comment 2"},
	}, nil
}

func (m *mockCommentService) Update(ctx context.Context, id, userID, content string) (*models.Comment, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, id, userID, content)
	}
	return &models.Comment{
		ID:      id,
		Content: content,
	}, nil
}

func (m *mockCommentService) Delete(ctx context.Context, id, userID string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id, userID)
	}
	return nil
}

func TestNewCommentController(t *testing.T) {
	mockService := &mockCommentService{}
	ctrl := NewCommentController(mockService)

	assert.NotNil(t, ctrl)
	assert.Equal(t, mockService, ctrl.commentService)
}

func TestCommentController_Create_Success(t *testing.T) {
	app := fiber.New()

	mockService := &mockCommentService{}
	ctrl := NewCommentController(mockService)
	app.Post("/comments", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.Create(c)
	})

	reqBody := `{"task_id":"task-123","content":"Test comment"}`
	req := httptest.NewRequest("POST", "/comments", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":true`)
	assert.Contains(t, respBody, `"content":"Test comment"`)
}

func TestCommentController_Create_ValidationError(t *testing.T) {
	app := fiber.New()

	mockService := &mockCommentService{}
	ctrl := NewCommentController(mockService)
	app.Post("/comments", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.Create(c)
	})

	reqBody := `{"task_id":"","content":""}`
	req := httptest.NewRequest("POST", "/comments", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestCommentController_FindByID_Success(t *testing.T) {
	app := fiber.New()

	mockService := &mockCommentService{}
	ctrl := NewCommentController(mockService)
	app.Get("/comments/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.FindByID(c)
	})

	req := httptest.NewRequest("GET", "/comments/comment-1", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":true`)
	assert.Contains(t, respBody, `"id":"comment-1"`)
}

func TestCommentController_FindByTaskID_Success(t *testing.T) {
	app := fiber.New()

	mockService := &mockCommentService{}
	ctrl := NewCommentController(mockService)
	app.Get("/comments/task/:task_id", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.FindByTaskID(c)
	})

	req := httptest.NewRequest("GET", "/comments/task/task-1", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":true`)
	assert.Contains(t, respBody, `"comment-1"`)
}

func TestCommentController_Update_Success(t *testing.T) {
	app := fiber.New()

	mockService := &mockCommentService{}
	ctrl := NewCommentController(mockService)
	app.Put("/comments/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.Update(c)
	})

	reqBody := `{"content":"Updated comment"}`
	req := httptest.NewRequest("PUT", "/comments/comment-1", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":true`)
	assert.Contains(t, respBody, `"content":"Updated comment"`)
}

func TestCommentController_Delete_Success(t *testing.T) {
	app := fiber.New()

	mockService := &mockCommentService{}
	ctrl := NewCommentController(mockService)
	app.Delete("/comments/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.Delete(c)
	})

	req := httptest.NewRequest("DELETE", "/comments/comment-1", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":true`)
	assert.Contains(t, respBody, `"message":"Comment deleted successfully"`)
}
