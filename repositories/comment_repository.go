package repositories

import (
	"context"
	"fmt"

	"kanban-backend/config"
	"kanban-backend/models"

	"gorm.io/gorm"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *models.Comment) error
	FindByID(ctx context.Context, id string) (*models.Comment, error)
	FindByTaskID(ctx context.Context, taskID string) ([]*models.Comment, error)
	FindByTaskIDWithPagination(ctx context.Context, taskID string, page, limit int) ([]*models.Comment, int, error)
	Update(ctx context.Context, comment *models.Comment) error
	Delete(ctx context.Context, id string) error
	SoftDelete(ctx context.Context, id string) error
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository() CommentRepository {
	return &commentRepository{
		db: config.DB,
	}
}

func (r *commentRepository) Create(ctx context.Context, comment *models.Comment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

func (r *commentRepository) FindByID(ctx context.Context, id string) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Task").
		Where("id = ?", id).
		First(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) FindByTaskID(ctx context.Context, taskID string) ([]*models.Comment, error) {
	var comments []*models.Comment
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Task").
		Where("task_id = ?", taskID).
		Order("created_at ASC").
		Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *commentRepository) Update(ctx context.Context, comment *models.Comment) error {
	return r.db.WithContext(ctx).Save(comment).Error
}

func (r *commentRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&models.Comment{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("comment with id %s not found", id)
	}
	return nil
}

func (r *commentRepository) SoftDelete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Comment{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("comment with id %s not found", id)
	}
	return nil
}

func (r *commentRepository) FindByTaskIDWithPagination(ctx context.Context, taskID string, page, limit int) ([]*models.Comment, int, error) {
	var comments []*models.Comment
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Comment{}).Where("task_id = ?", taskID)

	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Preload("User").
		Preload("Task").
		Order("created_at ASC").
		Offset(offset).
		Limit(limit).
		Find(&comments).Error

	return comments, int(total), err
}
