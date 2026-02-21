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

type mockAttachmentService struct {
	createFunc       func(ctx context.Context, taskID, userID, fileName, fileURL string, fileSize int64) (*models.Attachment, error)
	findByIDFunc     func(ctx context.Context, id, userID string) (*models.Attachment, error)
	findByTaskIDFunc func(ctx context.Context, taskID, userID string) ([]*models.Attachment, error)
	updateFunc       func(ctx context.Context, id, userID, fileName, fileURL string, fileSize int64) (*models.Attachment, error)
	deleteFunc       func(ctx context.Context, id, userID string) error
}

func (m *mockAttachmentService) Create(ctx context.Context, taskID, userID, fileName, fileURL string, fileSize int64) (*models.Attachment, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, taskID, userID, fileName, fileURL, fileSize)
	}
	return &models.Attachment{
		ID:       "attachment-1",
		TaskID:   taskID,
		FileName: fileName,
		FileURL:  fileURL,
		FileSize: fileSize,
	}, nil
}

func (m *mockAttachmentService) FindByID(ctx context.Context, id, userID string) (*models.Attachment, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id, userID)
	}
	return &models.Attachment{
		ID:       id,
		FileName: "file.pdf",
		FileURL:  "https://example.com/file.pdf",
		FileSize: 1024,
	}, nil
}

func (m *mockAttachmentService) FindByTaskID(ctx context.Context, taskID, userID string) ([]*models.Attachment, error) {
	if m.findByTaskIDFunc != nil {
		return m.findByTaskIDFunc(ctx, taskID, userID)
	}
	return []*models.Attachment{
		{ID: "attachment-1", TaskID: taskID, FileName: "file1.pdf", FileURL: "https://example.com/file1.pdf"},
		{ID: "attachment-2", TaskID: taskID, FileName: "file2.pdf", FileURL: "https://example.com/file2.pdf"},
	}, nil
}

func (m *mockAttachmentService) Update(ctx context.Context, id, userID, fileName, fileURL string, fileSize int64) (*models.Attachment, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, id, userID, fileName, fileURL, fileSize)
	}
	return &models.Attachment{
		ID:       id,
		FileName: fileName,
		FileURL:  fileURL,
		FileSize: fileSize,
	}, nil
}

func (m *mockAttachmentService) Delete(ctx context.Context, id, userID string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id, userID)
	}
	return nil
}

func TestNewAttachmentController(t *testing.T) {
	mockService := &mockAttachmentService{}
	ctrl := NewAttachmentController(mockService)

	assert.NotNil(t, ctrl)
	assert.Equal(t, mockService, ctrl.attachmentService)
}

func TestAttachmentController_Create_Success(t *testing.T) {
	app := fiber.New()

	mockService := &mockAttachmentService{}
	ctrl := NewAttachmentController(mockService)
	app.Post("/attachments", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.Create(c)
	})

	reqBody := `{"task_id":"task-123","file_name":"file.pdf","file_url":"https://example.com/file.pdf","file_size":1024}`
	req := httptest.NewRequest("POST", "/attachments", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":true`)
	assert.Contains(t, respBody, `"file_name":"file.pdf"`)
}

func TestAttachmentController_Create_ValidationError(t *testing.T) {
	app := fiber.New()

	mockService := &mockAttachmentService{}
	ctrl := NewAttachmentController(mockService)
	app.Post("/attachments", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.Create(c)
	})

	reqBody := `{"task_id":"","file_name":"","file_url":"","file_size":0}`
	req := httptest.NewRequest("POST", "/attachments", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestAttachmentController_FindByID_Success(t *testing.T) {
	app := fiber.New()

	mockService := &mockAttachmentService{}
	ctrl := NewAttachmentController(mockService)
	app.Get("/attachments/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.FindByID(c)
	})

	req := httptest.NewRequest("GET", "/attachments/attachment-1", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":true`)
	assert.Contains(t, respBody, `"id":"attachment-1"`)
}

func TestAttachmentController_FindByTaskID_Success(t *testing.T) {
	app := fiber.New()

	mockService := &mockAttachmentService{}
	ctrl := NewAttachmentController(mockService)
	app.Get("/attachments/task/:task_id", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.FindByTaskID(c)
	})

	req := httptest.NewRequest("GET", "/attachments/task/task-1", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":true`)
	assert.Contains(t, respBody, `"attachment-1"`)
}

func TestAttachmentController_Update_Success(t *testing.T) {
	app := fiber.New()

	mockService := &mockAttachmentService{}
	ctrl := NewAttachmentController(mockService)
	app.Put("/attachments/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.Update(c)
	})

	reqBody := `{"file_name":"updated.pdf","file_url":"https://example.com/updated.pdf","file_size":2048}`
	req := httptest.NewRequest("PUT", "/attachments/attachment-1", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":true`)
	assert.Contains(t, respBody, `"file_name":"updated.pdf"`)
}

func TestAttachmentController_Delete_Success(t *testing.T) {
	app := fiber.New()

	mockService := &mockAttachmentService{}
	ctrl := NewAttachmentController(mockService)
	app.Delete("/attachments/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.Delete(c)
	})

	req := httptest.NewRequest("DELETE", "/attachments/attachment-1", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":true`)
	assert.Contains(t, respBody, `"message":"Attachment deleted successfully"`)
}
