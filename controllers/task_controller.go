package controllers

import (
	"errors"
	"time"

	"kanban-backend/models"
	"kanban-backend/services"
	"kanban-backend/utils"

	"github.com/gofiber/fiber/v2"
)

type TaskController struct {
	taskService services.TaskService
}

func NewTaskController(taskService services.TaskService) *TaskController {
	return &TaskController{
		taskService: taskService,
	}
}

type CreateTaskRequest struct {
	ColumnID    string     `json:"column_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Deadline    *time.Time `json:"deadline,omitempty"`
}

type UpdateTaskRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Deadline    *time.Time `json:"deadline,omitempty"`
}

type MoveTaskRequest struct {
	ColumnID string `json:"column_id"`
}

type TaskResponse struct {
	ID          string              `json:"id"`
	ColumnID    string              `json:"column_id"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Deadline    *time.Time          `json:"deadline,omitempty"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	Comments    []models.Comment    `json:"comments,omitempty"`
	Labels      []models.Label      `json:"labels,omitempty"`
	Attachments []models.Attachment `json:"attachments,omitempty"`
}

func toTaskResponse(task *models.Task) TaskResponse {
	return TaskResponse{
		ID:          task.ID,
		ColumnID:    task.ColumnID,
		Title:       task.Title,
		Description: task.Description,
		Deadline:    task.Deadline,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		Comments:    task.Comments,
		Labels:      task.Labels,
		Attachments: task.Attachments,
	}
}

func toTaskResponseList(tasks []*models.Task) []TaskResponse {
	responses := make([]TaskResponse, len(tasks))
	for i, task := range tasks {
		responses[i] = toTaskResponse(task)
	}
	return responses
}

func (ctrl *TaskController) Create(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req CreateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, "Invalid request body", fiber.StatusBadRequest)
	}

	if req.Title == "" {
		return utils.ValidationError(c, "title", "title is required")
	}

	if req.ColumnID == "" {
		return utils.ValidationError(c, "column_id", "column_id is required")
	}

	task, err := ctrl.taskService.Create(c.Context(), userID, req.ColumnID, req.Title, req.Description, req.Deadline)
	if err != nil {
		var validationErr utils.ErrValidation
		if errors.As(err, &validationErr) {
			return utils.Error(c, err.Error(), fiber.StatusBadRequest)
		}
		return utils.Error(c, "Failed to create task", fiber.StatusInternalServerError)
	}

	return utils.Success(c, toTaskResponse(task))
}

func (ctrl *TaskController) FindByID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	taskID := c.Params("id")

	if taskID == "" {
		return utils.ValidationError(c, "id", "task id is required")
	}

	task, err := ctrl.taskService.FindByID(c.Context(), taskID, userID)
	if err != nil {
		var notFoundErr utils.ErrNotFound
		if errors.As(err, &notFoundErr) {
			return utils.Error(c, err.Error(), fiber.StatusNotFound)
		}
		var unauthorizedErr utils.ErrUnauthorized
		if errors.As(err, &unauthorizedErr) {
			return utils.Error(c, err.Error(), fiber.StatusUnauthorized)
		}
		return utils.Error(c, "Failed to find task", fiber.StatusInternalServerError)
	}

	return utils.Success(c, toTaskResponse(task))
}

func (ctrl *TaskController) FindByColumnID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	columnID := c.Params("columnId")

	if columnID == "" {
		return utils.ValidationError(c, "columnId", "column id is required")
	}

	var req utils.PaginationRequest
	c.QueryParser(&req)
	utils.ValidatePagination(&req)

	title := c.Query("title")

	tasks, total, err := ctrl.taskService.FindByColumnIDWithFilters(c.Context(), columnID, userID, title, req.Page, req.Limit)
	if err != nil {
		var notFoundErr utils.ErrNotFound
		if errors.As(err, &notFoundErr) {
			return utils.Error(c, err.Error(), fiber.StatusNotFound)
		}
		var unauthorizedErr utils.ErrUnauthorized
		if errors.As(err, &unauthorizedErr) {
			return utils.Error(c, err.Error(), fiber.StatusUnauthorized)
		}
		return utils.Error(c, "Failed to find tasks", fiber.StatusInternalServerError)
	}

	return utils.Success(c, utils.NewPaginatedResponse(toTaskResponseList(tasks), req.Page, req.Limit, total))
}

func (ctrl *TaskController) Update(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	taskID := c.Params("id")

	if taskID == "" {
		return utils.ValidationError(c, "id", "task id is required")
	}

	var req UpdateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, "Invalid request body", fiber.StatusBadRequest)
	}

	task, err := ctrl.taskService.Update(c.Context(), taskID, userID, req.Title, req.Description, req.Deadline)
	if err != nil {
		var notFoundErr utils.ErrNotFound
		if errors.As(err, &notFoundErr) {
			return utils.Error(c, err.Error(), fiber.StatusNotFound)
		}
		var unauthorizedErr utils.ErrUnauthorized
		if errors.As(err, &unauthorizedErr) {
			return utils.Error(c, err.Error(), fiber.StatusUnauthorized)
		}
		return utils.Error(c, "Failed to update task", fiber.StatusInternalServerError)
	}

	return utils.Success(c, toTaskResponse(task))
}

func (ctrl *TaskController) Delete(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	taskID := c.Params("id")

	if taskID == "" {
		return utils.ValidationError(c, "id", "task id is required")
	}

	err := ctrl.taskService.Delete(c.Context(), taskID, userID)
	if err != nil {
		var notFoundErr utils.ErrNotFound
		if errors.As(err, &notFoundErr) {
			return utils.Error(c, err.Error(), fiber.StatusNotFound)
		}
		var unauthorizedErr utils.ErrUnauthorized
		if errors.As(err, &unauthorizedErr) {
			return utils.Error(c, err.Error(), fiber.StatusUnauthorized)
		}
		return utils.Error(c, "Failed to delete task", fiber.StatusInternalServerError)
	}

	return utils.Success(c, fiber.Map{
		"message": "Task deleted successfully",
	})
}

func (ctrl *TaskController) Move(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	taskID := c.Params("id")

	if taskID == "" {
		return utils.ValidationError(c, "id", "task id is required")
	}

	var req MoveTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, "Invalid request body", fiber.StatusBadRequest)
	}

	if req.ColumnID == "" {
		return utils.ValidationError(c, "column_id", "column_id is required")
	}

	err := ctrl.taskService.Move(c.Context(), taskID, req.ColumnID, userID)
	if err != nil {
		var notFoundErr utils.ErrNotFound
		if errors.As(err, &notFoundErr) {
			return utils.Error(c, err.Error(), fiber.StatusNotFound)
		}
		var unauthorizedErr utils.ErrUnauthorized
		if errors.As(err, &unauthorizedErr) {
			return utils.Error(c, err.Error(), fiber.StatusUnauthorized)
		}
		return utils.Error(c, "Failed to move task", fiber.StatusInternalServerError)
	}

	return utils.Success(c, fiber.Map{
		"message": "Task moved successfully",
	})
}

func (ctrl *TaskController) Search(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	boardID := c.Query("board_id")

	if boardID == "" {
		return utils.ValidationError(c, "board_id", "board_id is required")
	}

	var req utils.PaginationRequest
	c.QueryParser(&req)
	utils.ValidatePagination(&req)

	keyword := c.Query("keyword")
	if keyword == "" {
		return utils.ValidationError(c, "keyword", "keyword is required")
	}

	tasks, total, err := ctrl.taskService.Search(c.Context(), boardID, userID, keyword, req.Page, req.Limit)
	if err != nil {
		var notFoundErr utils.ErrNotFound
		if errors.As(err, &notFoundErr) {
			return utils.Error(c, err.Error(), fiber.StatusNotFound)
		}
		var unauthorizedErr utils.ErrUnauthorized
		if errors.As(err, &unauthorizedErr) {
			return utils.Error(c, err.Error(), fiber.StatusUnauthorized)
		}
		return utils.Error(c, "Failed to search tasks", fiber.StatusInternalServerError)
	}

	return utils.Success(c, utils.NewPaginatedResponse(toTaskResponseList(tasks), req.Page, req.Limit, total))
}
