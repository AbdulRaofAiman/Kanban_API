package repositories

import (
	"context"
	"fmt"

	"kanban-backend/config"
	"kanban-backend/models"

	"gorm.io/gorm"
)

type AttachmentRepository interface {
	Create(ctx context.Context, attachment *models.Attachment) error
	FindByID(ctx context.Context, id string) (*models.Attachment, error)
	FindByTaskID(ctx context.Context, taskID string) ([]*models.Attachment, error)
	Update(ctx context.Context, attachment *models.Attachment) error
	Delete(ctx context.Context, id string) error
	SoftDelete(ctx context.Context, id string) error
}

type attachmentRepository struct {
	db *gorm.DB
}

func NewAttachmentRepository() AttachmentRepository {
	return &attachmentRepository{
		db: config.DB,
	}
}

func (r *attachmentRepository) Create(ctx context.Context, attachment *models.Attachment) error {
	return r.db.WithContext(ctx).Create(attachment).Error
}

func (r *attachmentRepository) FindByID(ctx context.Context, id string) (*models.Attachment, error) {
	var attachment models.Attachment
	err := r.db.WithContext(ctx).
		Preload("Task").
		Where("id = ?", id).
		First(&attachment).Error
	if err != nil {
		return nil, err
	}
	return &attachment, nil
}

func (r *attachmentRepository) FindByTaskID(ctx context.Context, taskID string) ([]*models.Attachment, error) {
	var attachments []*models.Attachment
	err := r.db.WithContext(ctx).
		Preload("Task").
		Where("task_id = ?", taskID).
		Order("created_at DESC").
		Find(&attachments).Error
	if err != nil {
		return nil, err
	}
	return attachments, nil
}

func (r *attachmentRepository) Update(ctx context.Context, attachment *models.Attachment) error {
	return r.db.WithContext(ctx).Save(attachment).Error
}

func (r *attachmentRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&models.Attachment{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("attachment with id %s not found", id)
	}
	return nil
}

func (r *attachmentRepository) SoftDelete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Attachment{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("attachment with id %s not found", id)
	}
	return nil
}
