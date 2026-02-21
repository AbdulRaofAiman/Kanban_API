package handlers

import (
	"errors"
	"time"

	"kanban-backend/models"
	"kanban-backend/services"
	"kanban-backend/utils"

	"github.com/gofiber/fiber/v2"
)

type BoardController struct {
	boardService services.BoardService
}

func NewBoardController(boardService services.BoardService) *BoardController {
	return &BoardController{
		boardService: boardService,
	}
}

type CreateBoardRequest struct {
	Title string `json:"title"`
	Color string `json:"color"`
}

type UpdateBoardRequest struct {
	Title string `json:"title"`
	Color string `json:"color"`
}

type BoardResponse struct {
	ID        string          `json:"id"`
	Title     string          `json:"title"`
	Color     string          `json:"color"`
	UserID    string          `json:"user_id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Columns   []models.Column `json:"columns"`
	Members   []models.Member `json:"members"`
}

func toBoardResponse(board *models.Board) BoardResponse {
	return BoardResponse{
		ID:        board.ID,
		Title:     board.Title,
		Color:     board.Color,
		UserID:    board.UserID,
		CreatedAt: board.CreatedAt,
		UpdatedAt: board.UpdatedAt,
		Columns:   board.Columns,
		Members:   board.Members,
	}
}

func toBoardResponseList(boards []*models.Board) []BoardResponse {
	responses := make([]BoardResponse, len(boards))
	for i, board := range boards {
		responses[i] = toBoardResponse(board)
	}
	return responses
}

func (ctrl *BoardController) Create(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req CreateBoardRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, "Invalid request body", fiber.StatusBadRequest)
	}

	if req.Title == "" {
		return utils.ValidationError(c, "title", "title is required")
	}

	if req.Color == "" {
		return utils.ValidationError(c, "color", "color is required")
	}

	board, err := ctrl.boardService.Create(c.Context(), userID, req.Title, req.Color)
	if err != nil {
		var validationErr utils.ErrValidation
		if errors.As(err, &validationErr) {
			return utils.Error(c, err.Error(), fiber.StatusBadRequest)
		}
		return utils.Error(c, "Failed to create board", fiber.StatusInternalServerError)
	}

	return utils.Success(c, toBoardResponse(board))
}

func (ctrl *BoardController) FindByID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	boardID := c.Params("id")

	if boardID == "" {
		return utils.ValidationError(c, "id", "board id is required")
	}

	board, err := ctrl.boardService.FindByID(c.Context(), boardID, userID)
	if err != nil {
		var notFoundErr utils.ErrNotFound
		if errors.As(err, &notFoundErr) {
			return utils.Error(c, err.Error(), fiber.StatusNotFound)
		}
		var unauthorizedErr utils.ErrUnauthorized
		if errors.As(err, &unauthorizedErr) {
			return utils.Error(c, err.Error(), fiber.StatusUnauthorized)
		}
		return utils.Error(c, "Failed to find board", fiber.StatusInternalServerError)
	}

	return utils.Success(c, toBoardResponse(board))
}

func (ctrl *BoardController) FindAll(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	boards, err := ctrl.boardService.FindByUserID(c.Context(), userID)
	if err != nil {
		return utils.Error(c, "Failed to find boards", fiber.StatusInternalServerError)
	}

	return utils.Success(c, toBoardResponseList(boards))
}

func (ctrl *BoardController) Update(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	boardID := c.Params("id")

	if boardID == "" {
		return utils.ValidationError(c, "id", "board id is required")
	}

	var req UpdateBoardRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, "Invalid request body", fiber.StatusBadRequest)
	}

	board, err := ctrl.boardService.Update(c.Context(), boardID, userID, req.Title, req.Color)
	if err != nil {
		var notFoundErr utils.ErrNotFound
		if errors.As(err, &notFoundErr) {
			return utils.Error(c, err.Error(), fiber.StatusNotFound)
		}
		var unauthorizedErr utils.ErrUnauthorized
		if errors.As(err, &unauthorizedErr) {
			return utils.Error(c, err.Error(), fiber.StatusUnauthorized)
		}
		return utils.Error(c, "Failed to update board", fiber.StatusInternalServerError)
	}

	return utils.Success(c, toBoardResponse(board))
}

func (ctrl *BoardController) Delete(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	boardID := c.Params("id")

	if boardID == "" {
		return utils.ValidationError(c, "id", "board id is required")
	}

	err := ctrl.boardService.Delete(c.Context(), boardID, userID)
	if err != nil {
		var notFoundErr utils.ErrNotFound
		if errors.As(err, &notFoundErr) {
			return utils.Error(c, err.Error(), fiber.StatusNotFound)
		}
		var unauthorizedErr utils.ErrUnauthorized
		if errors.As(err, &unauthorizedErr) {
			return utils.Error(c, err.Error(), fiber.StatusUnauthorized)
		}
		return utils.Error(c, "Failed to delete board", fiber.StatusInternalServerError)
	}

	return utils.Success(c, fiber.Map{
		"message": "Board deleted successfully",
	})
}
