package services

import (
	"context"
	"kanban-backend/models"
	"kanban-backend/repositories"
	"kanban-backend/utils"
)

type AttachmentService interface {
	Create(ctx context.Context, taskID, userID, fileName, fileURL string, fileSize int64) (*models.Attachment, error)
	FindByID(ctx context.Context, id, userID string) (*models.Attachment, error)
	FindByTaskID(ctx context.Context, taskID, userID string) ([]*models.Attachment, error)
	FindByTaskIDWithPagination(ctx context.Context, taskID, userID string, page, limit int) ([]*models.Attachment, int, error)
	Update(ctx context.Context, id, userID, fileName, fileURL string, fileSize int64) (*models.Attachment, error)
	Delete(ctx context.Context, id, userID string) error
}

type attachmentService struct {
	attachmentRepo repositories.AttachmentRepository
	taskRepo       repositories.TaskRepository
}

func NewAttachmentService(attachmentRepo repositories.AttachmentRepository, taskRepo repositories.TaskRepository) AttachmentService {
	return &attachmentService{
		attachmentRepo: attachmentRepo,
		taskRepo:       taskRepo,
	}
}

func (s *attachmentService) Create(ctx context.Context, taskID, userID, fileName, fileURL string, fileSize int64) (*models.Attachment, error) {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, utils.NewNotFound("task not found")
	}

	if task.Column == nil {
		return nil, utils.NewNotFound("column not found for task")
	}

	if task.Column.Board.UserID != userID {
		return nil, utils.NewUnauthorized("you do not have access to this task")
	}

	attachment := &models.Attachment{
		TaskID:   taskID,
		FileName: fileName,
		FileURL:  fileURL,
		FileSize: fileSize,
	}

	err = s.attachmentRepo.Create(ctx, attachment)
	if err != nil {
		return nil, err
	}

	return attachment, nil
}

func (s *attachmentService) FindByID(ctx context.Context, id, userID string) (*models.Attachment, error) {
	attachment, err := s.attachmentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFound("attachment not found")
	}

	if attachment.Task == nil {
		return nil, utils.NewNotFound("task not found for attachment")
	}

	if attachment.Task.Column == nil {
		return nil, utils.NewNotFound("column not found for task")
	}

	if attachment.Task.Column.Board.UserID != userID {
		return nil, utils.NewUnauthorized("you do not have access to this attachment")
	}

	return attachment, nil
}

func (s *attachmentService) FindByTaskID(ctx context.Context, taskID, userID string) ([]*models.Attachment, error) {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, utils.NewNotFound("task not found")
	}

	if task.Column == nil {
		return nil, utils.NewNotFound("column not found for task")
	}

	if task.Column.Board.UserID != userID {
		return nil, utils.NewUnauthorized("you do not have access to this task")
	}

	attachments, err := s.attachmentRepo.FindByTaskID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	return attachments, nil
}

func (s *attachmentService) Update(ctx context.Context, id, userID, fileName, fileURL string, fileSize int64) (*models.Attachment, error) {
	attachment, err := s.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	if fileName != "" {
		attachment.FileName = fileName
	}

	if fileURL != "" {
		attachment.FileURL = fileURL
	}

	if fileSize > 0 {
		attachment.FileSize = fileSize
	}

	err = s.attachmentRepo.Update(ctx, attachment)
	if err != nil {
		return nil, err
	}

	return attachment, nil
}

func (s *attachmentService) Delete(ctx context.Context, id, userID string) error {
	_, err := s.FindByID(ctx, id, userID)
	if err != nil {
		return err
	}

	err = s.attachmentRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *attachmentService) FindByTaskIDWithPagination(ctx context.Context, taskID, userID string, page, limit int) ([]*models.Attachment, int, error) {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, 0, utils.NewNotFound("task not found")
	}

	if task.Column == nil {
		return nil, 0, utils.NewNotFound("column not found for task")
	}

	if task.Column.Board.UserID != userID {
		return nil, 0, utils.NewUnauthorized("you do not have access to this task")
	}

	attachments, total, err := s.attachmentRepo.FindByTaskIDWithPagination(ctx, taskID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return attachments, total, nil
}
