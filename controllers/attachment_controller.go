package controllers

import (
	"errors"
	"time"

	"kanban-backend/models"
	"kanban-backend/services"
	"kanban-backend/utils"

	"github.com/gofiber/fiber/v2"
)

type AttachmentController struct {
	attachmentService services.AttachmentService
}

func NewAttachmentController(attachmentService services.AttachmentService) *AttachmentController {
	return &AttachmentController{
		attachmentService: attachmentService,
	}
}

type CreateAttachmentRequest struct {
	TaskID   string `json:"task_id"`
	FileName string `json:"file_name"`
	FileURL  string `json:"file_url"`
	FileSize int64  `json:"file_size"`
}

type UpdateAttachmentRequest struct {
	FileName string `json:"file_name"`
	FileURL  string `json:"file_url"`
	FileSize int64  `json:"file_size"`
}

type AttachmentResponse struct {
	ID        string    `json:"id"`
	TaskID    string    `json:"task_id"`
	FileName  string    `json:"file_name"`
	FileURL   string    `json:"file_url"`
	FileSize  int64     `json:"file_size"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toAttachmentResponse(attachment *models.Attachment) AttachmentResponse {
	return AttachmentResponse{
		ID:        attachment.ID,
		TaskID:    attachment.TaskID,
		FileName:  attachment.FileName,
		FileURL:   attachment.FileURL,
		FileSize:  attachment.FileSize,
		CreatedAt: attachment.CreatedAt,
		UpdatedAt: attachment.UpdatedAt,
	}
}

func toAttachmentResponseList(attachments []*models.Attachment) []AttachmentResponse {
	responses := make([]AttachmentResponse, len(attachments))
	for i, attachment := range attachments {
		responses[i] = toAttachmentResponse(attachment)
	}
	return responses
}

func (ctrl *AttachmentController) Create(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req CreateAttachmentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, "Invalid request body", fiber.StatusBadRequest)
	}

	if req.TaskID == "" {
		return utils.ValidationError(c, "task_id", "task_id is required")
	}

	if req.FileName == "" {
		return utils.ValidationError(c, "file_name", "file_name is required")
	}

	if req.FileURL == "" {
		return utils.ValidationError(c, "file_url", "file_url is required")
	}

	attachment, err := ctrl.attachmentService.Create(c.Context(), req.TaskID, userID, req.FileName, req.FileURL, req.FileSize)
	if err != nil {
		var validationErr utils.ErrValidation
		if errors.As(err, &validationErr) {
			return utils.Error(c, err.Error(), fiber.StatusBadRequest)
		}
		return utils.Error(c, "Failed to create attachment", fiber.StatusInternalServerError)
	}

	return utils.Success(c, toAttachmentResponse(attachment))
}

func (ctrl *AttachmentController) FindByID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	attachmentID := c.Params("id")

	if attachmentID == "" {
		return utils.ValidationError(c, "id", "attachment id is required")
	}

	attachment, err := ctrl.attachmentService.FindByID(c.Context(), attachmentID, userID)
	if err != nil {
		var notFoundErr utils.ErrNotFound
		if errors.As(err, &notFoundErr) {
			return utils.Error(c, err.Error(), fiber.StatusNotFound)
		}
		var unauthorizedErr utils.ErrUnauthorized
		if errors.As(err, &unauthorizedErr) {
			return utils.Error(c, err.Error(), fiber.StatusUnauthorized)
		}
		return utils.Error(c, "Failed to find attachment", fiber.StatusInternalServerError)
	}

	return utils.Success(c, toAttachmentResponse(attachment))
}

func (ctrl *AttachmentController) FindByTaskID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	taskID := c.Params("task_id")

	if taskID == "" {
		return utils.ValidationError(c, "task_id", "task id is required")
	}

	attachments, err := ctrl.attachmentService.FindByTaskID(c.Context(), taskID, userID)
	if err != nil {
		var notFoundErr utils.ErrNotFound
		if errors.As(err, &notFoundErr) {
			return utils.Error(c, err.Error(), fiber.StatusNotFound)
		}
		var unauthorizedErr utils.ErrUnauthorized
		if errors.As(err, &unauthorizedErr) {
			return utils.Error(c, err.Error(), fiber.StatusUnauthorized)
		}
		return utils.Error(c, "Failed to find attachments", fiber.StatusInternalServerError)
	}

	return utils.Success(c, toAttachmentResponseList(attachments))
}

func (ctrl *AttachmentController) Update(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	attachmentID := c.Params("id")

	if attachmentID == "" {
		return utils.ValidationError(c, "id", "attachment id is required")
	}

	var req UpdateAttachmentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, "Invalid request body", fiber.StatusBadRequest)
	}

	attachment, err := ctrl.attachmentService.Update(c.Context(), attachmentID, userID, req.FileName, req.FileURL, req.FileSize)
	if err != nil {
		var notFoundErr utils.ErrNotFound
		if errors.As(err, &notFoundErr) {
			return utils.Error(c, err.Error(), fiber.StatusNotFound)
		}
		var unauthorizedErr utils.ErrUnauthorized
		if errors.As(err, &unauthorizedErr) {
			return utils.Error(c, err.Error(), fiber.StatusUnauthorized)
		}
		return utils.Error(c, "Failed to update attachment", fiber.StatusInternalServerError)
	}

	return utils.Success(c, toAttachmentResponse(attachment))
}

func (ctrl *AttachmentController) Delete(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	attachmentID := c.Params("id")

	if attachmentID == "" {
		return utils.ValidationError(c, "id", "attachment id is required")
	}

	err := ctrl.attachmentService.Delete(c.Context(), attachmentID, userID)
	if err != nil {
		var notFoundErr utils.ErrNotFound
		if errors.As(err, &notFoundErr) {
			return utils.Error(c, err.Error(), fiber.StatusNotFound)
		}
		var unauthorizedErr utils.ErrUnauthorized
		if errors.As(err, &unauthorizedErr) {
			return utils.Error(c, err.Error(), fiber.StatusUnauthorized)
		}
		return utils.Error(c, "Failed to delete attachment", fiber.StatusInternalServerError)
	}

	return utils.Success(c, fiber.Map{
		"message": "Attachment deleted successfully",
	})
}
