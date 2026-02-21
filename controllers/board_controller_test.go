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

type mockBoardService struct {
	createFunc                  func(ctx context.Context, userID, title, color string) (*models.Board, error)
	findByIDFunc                func(ctx context.Context, boardID, userID string) (*models.Board, error)
	findByUserIDFunc            func(ctx context.Context, userID string) ([]*models.Board, error)
	findByUserIDWithFiltersFunc func(ctx context.Context, userID string, title string, page, limit int) ([]*models.Board, int, error)
	searchFunc                  func(ctx context.Context, userID string, keyword string, page, limit int) ([]*models.Board, int, error)
	updateFunc                  func(ctx context.Context, boardID, userID, title, color string) (*models.Board, error)
	deleteFunc                  func(ctx context.Context, boardID, userID string) error
}

func (m *mockBoardService) Create(ctx context.Context, userID, title, color string) (*models.Board, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, userID, title, color)
	}
	board := &models.Board{
		ID:        "board-123",
		Title:     title,
		UserID:    userID,
		Color:     color,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Columns: []models.Column{
			{ID: "col-1", Title: "To Do", Order: 1, BoardID: "board-123"},
			{ID: "col-2", Title: "In Progress", Order: 2, BoardID: "board-123"},
			{ID: "col-3", Title: "Done", Order: 3, BoardID: "board-123"},
		},
	}
	return board, nil
}

func (m *mockBoardService) FindByID(ctx context.Context, boardID, userID string) (*models.Board, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, boardID, userID)
	}
	board := &models.Board{
		ID:        boardID,
		Title:     "Test Board",
		UserID:    userID,
		Color:     "#000000",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Columns: []models.Column{
			{ID: "col-1", Title: "To Do", Order: 1, BoardID: boardID},
		},
	}
	return board, nil
}

func (m *mockBoardService) FindByUserID(ctx context.Context, userID string) ([]*models.Board, error) {
	if m.findByUserIDFunc != nil {
		return m.findByUserIDFunc(ctx, userID)
	}
	boards := []*models.Board{
		{
			ID:        "board-1",
			Title:     "Board 1",
			UserID:    userID,
			Color:     "#000000",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Columns:   []models.Column{},
		},
		{
			ID:        "board-2",
			Title:     "Board 2",
			UserID:    userID,
			Color:     "#ffffff",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Columns:   []models.Column{},
		},
	}
	return boards, nil
}

func (m *mockBoardService) Update(ctx context.Context, boardID, userID, title, color string) (*models.Board, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, boardID, userID, title, color)
	}
	board := &models.Board{
		ID:        boardID,
		Title:     title,
		UserID:    userID,
		Color:     color,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Columns:   []models.Column{},
	}
	return board, nil
}

func (m *mockBoardService) Delete(ctx context.Context, boardID, userID string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, boardID, userID)
	}
	return nil
}

func (m *mockBoardService) FindByUserIDWithFilters(ctx context.Context, userID string, title string, page, limit int) ([]*models.Board, int, error) {
	if m.findByUserIDWithFiltersFunc != nil {
		return m.findByUserIDWithFiltersFunc(ctx, userID, title, page, limit)
	}
	boards := []*models.Board{
		{
			ID:        "board-1",
			Title:     "Board 1",
			UserID:    userID,
			Color:     "#000000",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Columns:   []models.Column{},
		},
	}
	return boards, 1, nil
}

func (m *mockBoardService) Search(ctx context.Context, userID string, keyword string, page, limit int) ([]*models.Board, int, error) {
	if m.searchFunc != nil {
		return m.searchFunc(ctx, userID, keyword, page, limit)
	}
	boards := []*models.Board{
		{
			ID:        "board-1",
			Title:     "Test Board",
			UserID:    userID,
			Color:     "#000000",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Columns:   []models.Column{},
		},
	}
	return boards, 1, nil
}

func TestNewBoardController(t *testing.T) {
	mockService := &mockBoardService{}
	ctrl := NewBoardController(mockService)

	assert.NotNil(t, ctrl)
	assert.Equal(t, mockService, ctrl.boardService)
}

func TestBoardController_Create_Success(t *testing.T) {
	app := fiber.New()

	mockService := &mockBoardService{
		createFunc: func(ctx context.Context, userID, title, color string) (*models.Board, error) {
			board := &models.Board{
				ID:        "board-123",
				Title:     title,
				UserID:    userID,
				Color:     color,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Columns: []models.Column{
					{ID: "col-1", Title: "To Do", Order: 1, BoardID: "board-123"},
					{ID: "col-2", Title: "In Progress", Order: 2, BoardID: "board-123"},
					{ID: "col-3", Title: "Done", Order: 3, BoardID: "board-123"},
				},
			}
			return board, nil
		},
	}

	ctrl := NewBoardController(mockService)
	app.Post("/boards", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.Create(c)
	})

	reqBody := `{"title":"My Board","color":"#ff0000"}`
	req := httptest.NewRequest("POST", "/boards", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":true`)
	assert.Contains(t, respBody, `"id":"board-123"`)
	assert.Contains(t, respBody, `"title":"My Board"`)
	assert.Contains(t, respBody, `"color":"#ff0000"`)
	assert.Contains(t, respBody, `"columns"`)
}

func TestBoardController_Create_ValidationErrors(t *testing.T) {
	tests := []struct {
		name       string
		reqBody    string
		wantStatus int
		wantError  string
	}{
		{
			name:       "Missing title",
			reqBody:    `{"color":"#ff0000"}`,
			wantStatus: fiber.StatusBadRequest,
			wantError:  "title is required",
		},
		{
			name:       "Missing color",
			reqBody:    `{"title":"My Board"}`,
			wantStatus: fiber.StatusBadRequest,
			wantError:  "color is required",
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
			mockService := &mockBoardService{}
			ctrl := NewBoardController(mockService)
			app.Post("/boards", func(c *fiber.Ctx) error {
				c.Locals("user_id", "user-123")
				return ctrl.Create(c)
			})

			req := httptest.NewRequest("POST", "/boards", strings.NewReader(tt.reqBody))
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

func TestBoardController_Create_ServiceError(t *testing.T) {
	app := fiber.New()

	mockService := &mockBoardService{
		createFunc: func(ctx context.Context, userID, title, color string) (*models.Board, error) {
			return nil, utils.NewValidation("board title must be at least 3 characters")
		},
	}

	ctrl := NewBoardController(mockService)
	app.Post("/boards", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.Create(c)
	})

	reqBody := `{"title":"My Board","color":"#ff0000"}`
	req := httptest.NewRequest("POST", "/boards", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":false`)
	assert.Contains(t, respBody, "board title must be at least 3 characters")
}

func TestBoardController_FindByID_Success(t *testing.T) {
	app := fiber.New()

	mockService := &mockBoardService{
		findByIDFunc: func(ctx context.Context, boardID, userID string) (*models.Board, error) {
			board := &models.Board{
				ID:        boardID,
				Title:     "Test Board",
				UserID:    userID,
				Color:     "#000000",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Columns: []models.Column{
					{ID: "col-1", Title: "To Do", Order: 1, BoardID: boardID},
				},
			}
			return board, nil
		},
	}

	ctrl := NewBoardController(mockService)
	app.Get("/boards/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.FindByID(c)
	})

	req := httptest.NewRequest("GET", "/boards/board-123", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":true`)
	assert.Contains(t, respBody, `"id":"board-123"`)
	assert.Contains(t, respBody, `"title":"Test Board"`)
}

func TestBoardController_FindByID_NotFound(t *testing.T) {
	app := fiber.New()

	mockService := &mockBoardService{
		findByIDFunc: func(ctx context.Context, boardID, userID string) (*models.Board, error) {
			return nil, utils.NewNotFound("board not found")
		},
	}

	ctrl := NewBoardController(mockService)
	app.Get("/boards/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.FindByID(c)
	})

	req := httptest.NewRequest("GET", "/boards/board-999", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":false`)
	assert.Contains(t, respBody, "board not found")
}

func TestBoardController_FindByID_Unauthorized(t *testing.T) {
	app := fiber.New()

	mockService := &mockBoardService{
		findByIDFunc: func(ctx context.Context, boardID, userID string) (*models.Board, error) {
			return nil, utils.NewUnauthorized("you do not have access to this board")
		},
	}

	ctrl := NewBoardController(mockService)
	app.Get("/boards/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-999")
		return ctrl.FindByID(c)
	})

	req := httptest.NewRequest("GET", "/boards/board-123", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":false`)
	assert.Contains(t, respBody, "you do not have access to this board")
}

func TestBoardController_FindAll_Success(t *testing.T) {
	app := fiber.New()

	mockService := &mockBoardService{
		findByUserIDWithFiltersFunc: func(ctx context.Context, userID string, title string, page, limit int) ([]*models.Board, int, error) {
			boards := []*models.Board{
				{
					ID:        "board-1",
					Title:     "Board 1",
					UserID:    userID,
					Color:     "#000000",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Columns:   []models.Column{},
				},
				{
					ID:        "board-2",
					Title:     "Board 2",
					UserID:    userID,
					Color:     "#ffffff",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Columns:   []models.Column{},
				},
			}
			return boards, 2, nil
		},
	}

	ctrl := NewBoardController(mockService)
	app.Get("/boards", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.FindAll(c)
	})

	req := httptest.NewRequest("GET", "/boards", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":true`)
	assert.Contains(t, respBody, `"id":"board-1"`)
	assert.Contains(t, respBody, `"id":"board-2"`)
	assert.Contains(t, respBody, `"title":"Board 1"`)
	assert.Contains(t, respBody, `"title":"Board 2"`)
	assert.Contains(t, respBody, `"pagination"`)
}

func TestBoardController_Update_Success(t *testing.T) {
	app := fiber.New()

	mockService := &mockBoardService{
		updateFunc: func(ctx context.Context, boardID, userID, title, color string) (*models.Board, error) {
			board := &models.Board{
				ID:        boardID,
				Title:     title,
				UserID:    userID,
				Color:     color,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Columns:   []models.Column{},
			}
			return board, nil
		},
	}

	ctrl := NewBoardController(mockService)
	app.Put("/boards/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.Update(c)
	})

	reqBody := `{"title":"Updated Board","color":"#00ff00"}`
	req := httptest.NewRequest("PUT", "/boards/board-123", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":true`)
	assert.Contains(t, respBody, `"id":"board-123"`)
	assert.Contains(t, respBody, `"title":"Updated Board"`)
	assert.Contains(t, respBody, `"color":"#00ff00"`)
}

func TestBoardController_Update_NotFound(t *testing.T) {
	app := fiber.New()

	mockService := &mockBoardService{
		updateFunc: func(ctx context.Context, boardID, userID, title, color string) (*models.Board, error) {
			return nil, utils.NewNotFound("board not found")
		},
	}

	ctrl := NewBoardController(mockService)
	app.Put("/boards/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.Update(c)
	})

	reqBody := `{"title":"Updated Board","color":"#00ff00"}`
	req := httptest.NewRequest("PUT", "/boards/board-999", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":false`)
	assert.Contains(t, respBody, "board not found")
}

func TestBoardController_Delete_Success(t *testing.T) {
	app := fiber.New()

	mockService := &mockBoardService{
		deleteFunc: func(ctx context.Context, boardID, userID string) error {
			return nil
		},
	}

	ctrl := NewBoardController(mockService)
	app.Delete("/boards/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.Delete(c)
	})

	req := httptest.NewRequest("DELETE", "/boards/board-123", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":true`)
	assert.Contains(t, respBody, `"message":"Board deleted successfully"`)
}

func TestBoardController_Delete_NotFound(t *testing.T) {
	app := fiber.New()

	mockService := &mockBoardService{
		deleteFunc: func(ctx context.Context, boardID, userID string) error {
			return utils.NewNotFound("board not found")
		},
	}

	ctrl := NewBoardController(mockService)
	app.Delete("/boards/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.Delete(c)
	})

	req := httptest.NewRequest("DELETE", "/boards/board-999", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":false`)
	assert.Contains(t, respBody, "board not found")
}

func TestBoardController_ServiceError(t *testing.T) {
	app := fiber.New()

	mockService := &mockBoardService{
		findByIDFunc: func(ctx context.Context, boardID, userID string) (*models.Board, error) {
			return nil, errors.New("database connection failed")
		},
	}

	ctrl := NewBoardController(mockService)
	app.Get("/boards/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		return ctrl.FindByID(c)
	})

	req := httptest.NewRequest("GET", "/boards/board-123", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	respBody := string(body)

	assert.Contains(t, respBody, `"success":false`)
	assert.Contains(t, respBody, "Failed to find board")
}
