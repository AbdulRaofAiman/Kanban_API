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

type mockLabelService struct {
	createFunc                func(ctx context.Context, name, color string) (*models.Label, error)
	findByIDFunc              func(ctx context.Context, id string) (*models.Label, error)
	findAllFunc               func(ctx context.Context) ([]*models.Label, error)
	findAllWithPaginationFunc func(ctx context.Context, page, limit int) ([]*models.Label, int, error)
	searchFunc                func(ctx context.Context, keyword string, page, limit int) ([]*models.Label, int, error)
	updateFunc                func(ctx context.Context, id, name, color string) (*models.Label, error)
	deleteFunc                func(ctx context.Context, id string) error
	addToTaskFunc             func(ctx context.Context, taskID, labelID, userID string) error
	removeFromTaskFunc        func(ctx context.Context, taskID, labelID, userID string) error
}

func (m *mockLabelService) Create(ctx context.Context, name, color string) (*models.Label, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, name, color)
	}
	return &models.Label{
		ID:    "label-1",
		Name:  name,
		Color: color,
	}, nil
}

func (m *mockLabelService) FindByID(ctx context.Context, id string) (*models.Label, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return &models.Label{
		ID:    id,
		Name:  "Bug",
		Color: "#FF0000",
	}, nil
}

func (m *mockLabelService) FindAll(ctx context.Context) ([]*models.Label, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc(ctx)
	}
	return []*models.Label{
		{ID: "label-1", Name: "Bug", Color: "#FF0000"},
		{ID: "label-2", Name: "Feature", Color: "#00FF00"},
	}, nil
}

func (m *mockLabelService) Update(ctx context.Context, id, name, color string) (*models.Label, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, id, name, color)
	}
	return &models.Label{
		ID:    id,
		Name:  name,
		Color: color,
	}, nil
}

func (m *mockLabelService) Delete(ctx context.Context, id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func (m *mockLabelService) AddToTask(ctx context.Context, taskID, labelID, userID string) error {
	if m.addToTaskFunc != nil {
		return m.addToTaskFunc(ctx, taskID, labelID, userID)
	}
	return nil
}

func (m *mockLabelService) RemoveFromTask(ctx context.Context, taskID, labelID, userID string) error {
	if m.removeFromTaskFunc != nil {
		return m.removeFromTaskFunc(ctx, taskID, labelID, userID)
	}
	return nil
}

func (m *mockLabelService) FindAllWithPagination(ctx context.Context, page, limit int) ([]*models.Label, int, error) {
	if m.findAllWithPaginationFunc != nil {
		return m.findAllWithPaginationFunc(ctx, page, limit)
	}
	return []*models.Label{
		{ID: "label-1", Name: "Bug", Color: "#FF0000"},
	}, 1, nil
}

func (m *mockLabelService) Search(ctx context.Context, keyword string, page, limit int) ([]*models.Label, int, error) {
	if m.searchFunc != nil {
		return m.searchFunc(ctx, keyword, page, limit)
	}
	return []*models.Label{
		{ID: "label-1", Name: "Bug", Color: "#FF0000"},
	}, 1, nil
}

func TestNewLabelController(t *testing.T) {
	mockService := &mockLabelService{}
	ctrl := NewLabelController(mockService)

	assert.NotNil(t, ctrl)
	assert.Equal(t, mockService, ctrl.labelService)
}

func TestLabelController_Create_Success(t *testing.T) {
	app := fiber.New()

	mockService := &mockLabelService{}
	ctrl := NewLabelController(mockService)
	app.Post("/labels", ctrl.Create)

	reqBody := `{"name":"Bug","color":"#FF0000"}`
	req := httptest.NewRequest("POST", "/labels", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":true`)
	assert.Contains(t, respBody, `"name":"Bug"`)
}

func TestLabelController_Create_ValidationError(t *testing.T) {
	app := fiber.New()

	mockService := &mockLabelService{}
	ctrl := NewLabelController(mockService)
	app.Post("/labels", ctrl.Create)

	reqBody := `{"name":"","color":""}`
	req := httptest.NewRequest("POST", "/labels", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestLabelController_FindAll_Success(t *testing.T) {
	app := fiber.New()

	mockService := &mockLabelService{
		findAllWithPaginationFunc: func(ctx context.Context, page, limit int) ([]*models.Label, int, error) {
			labels := []*models.Label{
				{ID: "label-1", Name: "Bug", Color: "#FF0000"},
				{ID: "label-2", Name: "Feature", Color: "#00FF00"},
			}
			return labels, 2, nil
		},
	}
	ctrl := NewLabelController(mockService)
	app.Get("/labels", ctrl.FindAll)

	req := httptest.NewRequest("GET", "/labels", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":true`)
	assert.Contains(t, respBody, `"label-1"`)
	assert.Contains(t, respBody, `"label-2"`)
	assert.Contains(t, respBody, `"pagination"`)
}

func TestLabelController_Update_Success(t *testing.T) {
	app := fiber.New()

	mockService := &mockLabelService{}
	ctrl := NewLabelController(mockService)
	app.Put("/labels/:id", ctrl.Update)

	reqBody := `{"name":"Critical","color":"#FF0000"}`
	req := httptest.NewRequest("PUT", "/labels/label-1", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":true`)
	assert.Contains(t, respBody, `"name":"Critical"`)
}

func TestLabelController_Delete_Success(t *testing.T) {
	app := fiber.New()

	mockService := &mockLabelService{}
	ctrl := NewLabelController(mockService)
	app.Delete("/labels/:id", ctrl.Delete)

	req := httptest.NewRequest("DELETE", "/labels/label-1", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":true`)
	assert.Contains(t, respBody, `"message":"Label deleted successfully"`)
}

func TestLabelController_AddToTask_Success(t *testing.T) {
	app := fiber.New()

	mockService := &mockLabelService{}
	ctrl := NewLabelController(mockService)
	app.Post("/tasks/:id/labels/:label_id", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.AddToTask(c)
	})

	req := httptest.NewRequest("POST", "/tasks/task-1/labels/label-1", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":true`)
	assert.Contains(t, respBody, `"message":"Label added to task successfully"`)
}

func TestLabelController_RemoveFromTask_Success(t *testing.T) {
	app := fiber.New()

	mockService := &mockLabelService{}
	ctrl := NewLabelController(mockService)
	app.Delete("/tasks/:id/labels/:label_id", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.RemoveFromTask(c)
	})

	req := httptest.NewRequest("DELETE", "/tasks/task-1/labels/label-1", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":true`)
	assert.Contains(t, respBody, `"message":"Label removed from task successfully"`)
}
