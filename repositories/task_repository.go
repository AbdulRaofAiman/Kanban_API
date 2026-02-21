package repositories

import (
	"context"
	"fmt"

	"kanban-backend/config"
	"kanban-backend/models"

	"gorm.io/gorm"
)

type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	FindByID(ctx context.Context, id string) (*models.Task, error)
	FindByColumnID(ctx context.Context, columnID string) ([]*models.Task, error)
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id string) error
	SoftDelete(ctx context.Context, id string) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository() TaskRepository {
	return &taskRepository{
		db: config.DB,
	}
}

func (r *taskRepository) Create(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *taskRepository) FindByID(ctx context.Context, id string) (*models.Task, error) {
	var task models.Task
	err := r.db.WithContext(ctx).
		Preload("Comments").
		Preload("Labels").
		Preload("Attachments").
		Preload("Column").
		Where("id = ?", id).
		First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) FindByColumnID(ctx context.Context, columnID string) ([]*models.Task, error) {
	var tasks []*models.Task
	err := r.db.WithContext(ctx).
		Preload("Comments").
		Preload("Labels").
		Preload("Attachments").
		Preload("Column").
		Where("column_id = ?", columnID).
		Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepository) Update(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).Save(task).Error
}

func (r *taskRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&models.Task{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("task with id %s not found", id)
	}
	return nil
}

func (r *taskRepository) SoftDelete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Task{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("task with id %s not found", id)
	}
	return nil
}
