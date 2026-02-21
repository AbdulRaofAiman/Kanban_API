package controllers

import (
	"errors"
	"time"

	"kanban-backend/models"
	"kanban-backend/services"
	"kanban-backend/utils"

	"github.com/gofiber/fiber/v2"
)

type CommentController struct {
	commentService services.CommentService
}

func NewCommentController(commentService services.CommentService) *CommentController {
	return &CommentController{
		commentService: commentService,
	}
}

type CreateCommentRequest struct {
	TaskID  string `json:"task_id"`
	Content string `json:"content"`
}

type UpdateCommentRequest struct {
	Content string `json:"content"`
}

type CommentResponse struct {
	ID        string       `json:"id"`
	TaskID    string       `json:"task_id"`
	UserID    string       `json:"user_id"`
	Content   string       `json:"content"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	User      *models.User `json:"user,omitempty"`
}

func toCommentResponse(comment *models.Comment) CommentResponse {
	return CommentResponse{
		ID:        comment.ID,
		TaskID:    comment.TaskID,
		UserID:    comment.UserID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
		User:      comment.User,
	}
}

func toCommentResponseList(comments []*models.Comment) []CommentResponse {
	responses := make([]CommentResponse, len(comments))
	for i, comment := range comments {
		responses[i] = toCommentResponse(comment)
	}
	return responses
}

func (ctrl *CommentController) Create(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req CreateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, "Invalid request body", fiber.StatusBadRequest)
	}

	if req.TaskID == "" {
		return utils.ValidationError(c, "task_id", "task_id is required")
	}

	if req.Content == "" {
		return utils.ValidationError(c, "content", "content is required")
	}

	comment, err := ctrl.commentService.Create(c.Context(), req.TaskID, userID, req.Content)
	if err != nil {
		var validationErr utils.ErrValidation
		if errors.As(err, &validationErr) {
			return utils.Error(c, err.Error(), fiber.StatusBadRequest)
		}
		return utils.Error(c, "Failed to create comment", fiber.StatusInternalServerError)
	}

	return utils.Success(c, toCommentResponse(comment))
}

func (ctrl *CommentController) FindByID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	commentID := c.Params("id")

	if commentID == "" {
		return utils.ValidationError(c, "id", "comment id is required")
	}

	comment, err := ctrl.commentService.FindByID(c.Context(), commentID, userID)
	if err != nil {
		var notFoundErr utils.ErrNotFound
		if errors.As(err, &notFoundErr) {
			return utils.Error(c, err.Error(), fiber.StatusNotFound)
		}
		var unauthorizedErr utils.ErrUnauthorized
		if errors.As(err, &unauthorizedErr) {
			return utils.Error(c, err.Error(), fiber.StatusUnauthorized)
		}
		return utils.Error(c, "Failed to find comment", fiber.StatusInternalServerError)
	}

	return utils.Success(c, toCommentResponse(comment))
}

func (ctrl *CommentController) FindByTaskID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	taskID := c.Params("task_id")

	if taskID == "" {
		return utils.ValidationError(c, "task_id", "task id is required")
	}

	comments, err := ctrl.commentService.FindByTaskID(c.Context(), taskID, userID)
	if err != nil {
		var notFoundErr utils.ErrNotFound
		if errors.As(err, &notFoundErr) {
			return utils.Error(c, err.Error(), fiber.StatusNotFound)
		}
		var unauthorizedErr utils.ErrUnauthorized
		if errors.As(err, &unauthorizedErr) {
			return utils.Error(c, err.Error(), fiber.StatusUnauthorized)
		}
		return utils.Error(c, "Failed to find comments", fiber.StatusInternalServerError)
	}

	return utils.Success(c, toCommentResponseList(comments))
}

func (ctrl *CommentController) Update(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	commentID := c.Params("id")

	if commentID == "" {
		return utils.ValidationError(c, "id", "comment id is required")
	}

	var req UpdateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, "Invalid request body", fiber.StatusBadRequest)
	}

	if req.Content == "" {
		return utils.ValidationError(c, "content", "content is required")
	}

	comment, err := ctrl.commentService.Update(c.Context(), commentID, userID, req.Content)
	if err != nil {
		var notFoundErr utils.ErrNotFound
		if errors.As(err, &notFoundErr) {
			return utils.Error(c, err.Error(), fiber.StatusNotFound)
		}
		var unauthorizedErr utils.ErrUnauthorized
		if errors.As(err, &unauthorizedErr) {
			return utils.Error(c, err.Error(), fiber.StatusUnauthorized)
		}
		return utils.Error(c, "Failed to update comment", fiber.StatusInternalServerError)
	}

	return utils.Success(c, toCommentResponse(comment))
}

func (ctrl *CommentController) Delete(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	commentID := c.Params("id")

	if commentID == "" {
		return utils.ValidationError(c, "id", "comment id is required")
	}

	err := ctrl.commentService.Delete(c.Context(), commentID, userID)
	if err != nil {
		var notFoundErr utils.ErrNotFound
		if errors.As(err, &notFoundErr) {
			return utils.Error(c, err.Error(), fiber.StatusNotFound)
		}
		var unauthorizedErr utils.ErrUnauthorized
		if errors.As(err, &unauthorizedErr) {
			return utils.Error(c, err.Error(), fiber.StatusUnauthorized)
		}
		return utils.Error(c, "Failed to delete comment", fiber.StatusInternalServerError)
	}

	return utils.Success(c, fiber.Map{
		"message": "Comment deleted successfully",
	})
}
