package controllers

import (
	"errors"

	"kanban-backend/models"
	"kanban-backend/services"
	"kanban-backend/utils"

	"github.com/gofiber/fiber/v2"
)

type LabelController struct {
	labelService services.LabelService
}

func NewLabelController(labelService services.LabelService) *LabelController {
	return &LabelController{
		labelService: labelService,
	}
}

type CreateLabelRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type UpdateLabelRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type LabelResponse struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Color     string         `json:"color"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
	Tasks     []*models.Task `json:"tasks,omitempty"`
}

func toLabelResponse(label *models.Label) LabelResponse {
	tasks := make([]*models.Task, len(label.Tasks))
	for i, t := range label.Tasks {
		tasks[i] = t
	}
	return LabelResponse{
		ID:        label.ID,
		Name:      label.Name,
		Color:     label.Color,
		CreatedAt: label.CreatedAt.String(),
		UpdatedAt: label.UpdatedAt.String(),
		Tasks:     tasks,
	}
}

func toLabelResponseList(labels []*models.Label) []LabelResponse {
	responses := make([]LabelResponse, len(labels))
	for i, label := range labels {
		responses[i] = toLabelResponse(label)
	}
	return responses
}

func (ctrl *LabelController) Create(c *fiber.Ctx) error {
	var req CreateLabelRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, "Invalid request body", fiber.StatusBadRequest)
	}

	if req.Name == "" {
		return utils.ValidationError(c, "name", "name is required")
	}

	label, err := ctrl.labelService.Create(c.Context(), req.Name, req.Color)
	if err != nil {
		var validationErr utils.ErrValidation
		if errors.As(err, &validationErr) {
			return utils.Error(c, err.Error(), fiber.StatusBadRequest)
		}
		return utils.Error(c, "Failed to create label", fiber.StatusInternalServerError)
	}

	return utils.Success(c, toLabelResponse(label))
}

func (ctrl *LabelController) FindByID(c *fiber.Ctx) error {
	labelID := c.Params("id")

	if labelID == "" {
		return utils.ValidationError(c, "id", "label id is required")
	}

	label, err := ctrl.labelService.FindByID(c.Context(), labelID)
	if err != nil {
		var notFoundErr utils.ErrNotFound
		if errors.As(err, &notFoundErr) {
			return utils.Error(c, err.Error(), fiber.StatusNotFound)
		}
		return utils.Error(c, "Failed to find label", fiber.StatusInternalServerError)
	}

	return utils.Success(c, toLabelResponse(label))
}

func (ctrl *LabelController) FindAll(c *fiber.Ctx) error {
	var req utils.PaginationRequest
	c.QueryParser(&req)
	utils.ValidatePagination(&req)

	labels, total, err := ctrl.labelService.FindAllWithPagination(c.Context(), req.Page, req.Limit)
	if err != nil {
		return utils.Error(c, "Failed to find labels", fiber.StatusInternalServerError)
	}

	return utils.Success(c, utils.NewPaginatedResponse(toLabelResponseList(labels), req.Page, req.Limit, total))
}

func (ctrl *LabelController) Update(c *fiber.Ctx) error {
	labelID := c.Params("id")

	if labelID == "" {
		return utils.ValidationError(c, "id", "label id is required")
	}

	var req UpdateLabelRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, "Invalid request body", fiber.StatusBadRequest)
	}

	label, err := ctrl.labelService.Update(c.Context(), labelID, req.Name, req.Color)
	if err != nil {
		var notFoundErr utils.ErrNotFound
		if errors.As(err, &notFoundErr) {
			return utils.Error(c, err.Error(), fiber.StatusNotFound)
		}
		return utils.Error(c, "Failed to update label", fiber.StatusInternalServerError)
	}

	return utils.Success(c, toLabelResponse(label))
}

func (ctrl *LabelController) Delete(c *fiber.Ctx) error {
	labelID := c.Params("id")

	if labelID == "" {
		return utils.ValidationError(c, "id", "label id is required")
	}

	err := ctrl.labelService.Delete(c.Context(), labelID)
	if err != nil {
		var notFoundErr utils.ErrNotFound
		if errors.As(err, &notFoundErr) {
			return utils.Error(c, err.Error(), fiber.StatusNotFound)
		}
		return utils.Error(c, "Failed to delete label", fiber.StatusInternalServerError)
	}

	return utils.Success(c, fiber.Map{
		"message": "Label deleted successfully",
	})
}

func (ctrl *LabelController) AddToTask(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	taskID := c.Params("id")
	labelID := c.Params("label_id")

	if taskID == "" {
		return utils.ValidationError(c, "id", "task id is required")
	}

	if labelID == "" {
		return utils.ValidationError(c, "label_id", "label id is required")
	}

	err := ctrl.labelService.AddToTask(c.Context(), taskID, labelID, userID)
	if err != nil {
		var notFoundErr utils.ErrNotFound
		if errors.As(err, &notFoundErr) {
			return utils.Error(c, err.Error(), fiber.StatusNotFound)
		}
		var unauthorizedErr utils.ErrUnauthorized
		if errors.As(err, &unauthorizedErr) {
			return utils.Error(c, err.Error(), fiber.StatusUnauthorized)
		}
		return utils.Error(c, "Failed to add label to task", fiber.StatusInternalServerError)
	}

	return utils.Success(c, fiber.Map{
		"message": "Label added to task successfully",
	})
}

func (ctrl *LabelController) RemoveFromTask(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	taskID := c.Params("id")
	labelID := c.Params("label_id")

	if taskID == "" {
		return utils.ValidationError(c, "id", "task id is required")
	}

	if labelID == "" {
		return utils.ValidationError(c, "label_id", "label id is required")
	}

	err := ctrl.labelService.RemoveFromTask(c.Context(), taskID, labelID, userID)
	if err != nil {
		var notFoundErr utils.ErrNotFound
		if errors.As(err, &notFoundErr) {
			return utils.Error(c, err.Error(), fiber.StatusNotFound)
		}
		var unauthorizedErr utils.ErrUnauthorized
		if errors.As(err, &unauthorizedErr) {
			return utils.Error(c, err.Error(), fiber.StatusUnauthorized)
		}
		return utils.Error(c, "Failed to remove label from task", fiber.StatusInternalServerError)
	}

	return utils.Success(c, fiber.Map{
		"message": "Label removed from task successfully",
	})
}

func (ctrl *LabelController) Search(c *fiber.Ctx) error {
	var req utils.PaginationRequest
	c.QueryParser(&req)
	utils.ValidatePagination(&req)

	keyword := c.Query("keyword")
	if keyword == "" {
		return utils.ValidationError(c, "keyword", "keyword is required")
	}

	labels, total, err := ctrl.labelService.Search(c.Context(), keyword, req.Page, req.Limit)
	if err != nil {
		return utils.Error(c, "Failed to search labels", fiber.StatusInternalServerError)
	}

	return utils.Success(c, utils.NewPaginatedResponse(toLabelResponseList(labels), req.Page, req.Limit, total))
}
