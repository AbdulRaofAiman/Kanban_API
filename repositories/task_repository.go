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
	FindByColumnIDWithFilters(ctx context.Context, columnID string, title string, page, limit int) ([]*models.Task, int, error)
	Search(ctx context.Context, boardID string, keyword string, page, limit int) ([]*models.Task, int, error)
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

func (r *taskRepository) FindByColumnIDWithFilters(ctx context.Context, columnID string, title string, page, limit int) ([]*models.Task, int, error) {
	var tasks []*models.Task
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Task{}).Where("column_id = ?", columnID)

	if title != "" {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}

	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Preload("Comments").
		Preload("Labels").
		Preload("Attachments").
		Preload("Column").
		Offset(offset).
		Limit(limit).
		Find(&tasks).Error

	return tasks, int(total), err
}

func (r *taskRepository) Search(ctx context.Context, boardID string, keyword string, page, limit int) ([]*models.Task, int, error) {
	var tasks []*models.Task
	var total int64

	query := r.db.WithContext(ctx).
		Joins("JOIN columns ON columns.id = tasks.column_id").
		Joins("JOIN boards ON boards.id = columns.board_id").
		Where("boards.id = ?", boardID).
		Where("tasks.title ILIKE ? OR tasks.description ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")

	query.Model(&models.Task{}).Count(&total)

	offset := (page - 1) * limit
	err := query.Preload("Comments").
		Preload("Labels").
		Preload("Attachments").
		Preload("Column").
		Offset(offset).
		Limit(limit).
		Find(&tasks).Error

	return tasks, int(total), err
}
