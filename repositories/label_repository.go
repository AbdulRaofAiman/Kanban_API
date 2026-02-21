package repositories

import (
	"context"
	"fmt"

	"kanban-backend/config"
	"kanban-backend/models"

	"gorm.io/gorm"
)

type LabelRepository interface {
	Create(ctx context.Context, label *models.Label) error
	FindByID(ctx context.Context, id string) (*models.Label, error)
	FindAll(ctx context.Context) ([]*models.Label, error)
	FindAllWithPagination(ctx context.Context, page, limit int) ([]*models.Label, int, error)
	Search(ctx context.Context, keyword string, page, limit int) ([]*models.Label, int, error)
	Update(ctx context.Context, label *models.Label) error
	Delete(ctx context.Context, id string) error
	SoftDelete(ctx context.Context, id string) error
}

type labelRepository struct {
	db *gorm.DB
}

func NewLabelRepository() LabelRepository {
	return &labelRepository{
		db: config.DB,
	}
}

func (r *labelRepository) Create(ctx context.Context, label *models.Label) error {
	return r.db.WithContext(ctx).Create(label).Error
}

func (r *labelRepository) FindByID(ctx context.Context, id string) (*models.Label, error) {
	var label models.Label
	err := r.db.WithContext(ctx).
		Preload("Tasks").
		Where("id = ?", id).
		First(&label).Error
	if err != nil {
		return nil, err
	}
	return &label, nil
}

func (r *labelRepository) FindAll(ctx context.Context) ([]*models.Label, error) {
	var labels []*models.Label
	err := r.db.WithContext(ctx).
		Preload("Tasks").
		Order("name ASC").
		Find(&labels).Error
	if err != nil {
		return nil, err
	}
	return labels, nil
}

func (r *labelRepository) Update(ctx context.Context, label *models.Label) error {
	return r.db.WithContext(ctx).Save(label).Error
}

func (r *labelRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&models.Label{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("label with id %s not found", id)
	}
	return nil
}

func (r *labelRepository) SoftDelete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Label{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("label with id %s not found", id)
	}
	return nil
}

func (r *labelRepository) FindAllWithPagination(ctx context.Context, page, limit int) ([]*models.Label, int, error) {
	var labels []*models.Label
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Label{})

	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Preload("Tasks").
		Order("name ASC").
		Offset(offset).
		Limit(limit).
		Find(&labels).Error

	return labels, int(total), err
}

func (r *labelRepository) Search(ctx context.Context, keyword string, page, limit int) ([]*models.Label, int, error) {
	var labels []*models.Label
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Label{}).
		Where("name ILIKE ? OR color ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")

	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Preload("Tasks").
		Order("name ASC").
		Offset(offset).
		Limit(limit).
		Find(&labels).Error

	return labels, int(total), err
}
