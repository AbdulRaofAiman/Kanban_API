package services

import (
	"context"
	"errors"

	"kanban-backend/config"
	"kanban-backend/models"
	"kanban-backend/repositories"
	"kanban-backend/utils"

	"gorm.io/gorm"
)

type LabelService interface {
	Create(ctx context.Context, name, color string) (*models.Label, error)
	FindByID(ctx context.Context, id string) (*models.Label, error)
	FindAll(ctx context.Context) ([]*models.Label, error)
	Update(ctx context.Context, id, name, color string) (*models.Label, error)
	Delete(ctx context.Context, id string) error
	AddToTask(ctx context.Context, taskID, labelID, userID string) error
	RemoveFromTask(ctx context.Context, taskID, labelID, userID string) error
}

type labelService struct {
	labelRepo repositories.LabelRepository
	taskRepo  repositories.TaskRepository
	db        *gorm.DB
}

func NewLabelService(labelRepo repositories.LabelRepository, taskRepo repositories.TaskRepository) LabelService {
	return &labelService{
		labelRepo: labelRepo,
		taskRepo:  taskRepo,
		db:        config.DB,
	}
}

func (s *labelService) Create(ctx context.Context, name, color string) (*models.Label, error) {
	if name == "" {
		return nil, utils.NewValidation("name is required")
	}

	label := &models.Label{
		Name:  name,
		Color: color,
	}

	err := s.labelRepo.Create(ctx, label)
	if err != nil {
		return nil, err
	}

	return label, nil
}

func (s *labelService) FindByID(ctx context.Context, id string) (*models.Label, error) {
	label, err := s.labelRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFound("label not found")
	}

	return label, nil
}

func (s *labelService) FindAll(ctx context.Context) ([]*models.Label, error) {
	labels, err := s.labelRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return labels, nil
}

func (s *labelService) Update(ctx context.Context, id, name, color string) (*models.Label, error) {
	label, err := s.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		label.Name = name
	}

	if color != "" {
		label.Color = color
	}

	err = s.labelRepo.Update(ctx, label)
	if err != nil {
		return nil, err
	}

	return label, nil
}

func (s *labelService) Delete(ctx context.Context, id string) error {
	_, err := s.FindByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.labelRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *labelService) AddToTask(ctx context.Context, taskID, labelID, userID string) error {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return utils.NewNotFound("task not found")
	}

	if task.Column == nil {
		return utils.NewNotFound("column not found for task")
	}

	if task.Column.Board.UserID != userID {
		return utils.NewUnauthorized("you do not have access to this task")
	}

	label, err := s.FindByID(ctx, labelID)
	if err != nil {
		return err
	}

	err = s.db.WithContext(ctx).Model(task).Association("Labels").Append(label)
	if err != nil {
		return errors.New("failed to add label to task")
	}

	return nil
}

func (s *labelService) RemoveFromTask(ctx context.Context, taskID, labelID, userID string) error {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return utils.NewNotFound("task not found")
	}

	if task.Column == nil {
		return utils.NewNotFound("column not found for task")
	}

	if task.Column.Board.UserID != userID {
		return utils.NewUnauthorized("you do not have access to this task")
	}

	label, err := s.FindByID(ctx, labelID)
	if err != nil {
		return err
	}

	_ = label
	if err != nil {
		return errors.New("failed to remove label from task")
	}

	return nil
}
