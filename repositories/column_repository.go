package repositories

import (
	"context"
	"fmt"

	"kanban-backend/config"
	"kanban-backend/models"

	"gorm.io/gorm"
)

type ColumnRepository interface {
	Create(ctx context.Context, column *models.Column) error
	FindByID(ctx context.Context, id string) (*models.Column, error)
	FindByBoardID(ctx context.Context, boardID string) ([]*models.Column, error)
	Update(ctx context.Context, column *models.Column) error
	Delete(ctx context.Context, id string) error
	SoftDelete(ctx context.Context, id string) error
}

type columnRepository struct {
	db *gorm.DB
}

func NewColumnRepository() ColumnRepository {
	return &columnRepository{
		db: config.DB,
	}
}

func (r *columnRepository) Create(ctx context.Context, column *models.Column) error {
	return r.db.WithContext(ctx).Create(column).Error
}

func (r *columnRepository) FindByID(ctx context.Context, id string) (*models.Column, error) {
	var column models.Column
	err := r.db.WithContext(ctx).
		Preload("Tasks").
		Preload("Board").
		Where("id = ?", id).
		First(&column).Error
	if err != nil {
		return nil, err
	}
	return &column, nil
}

func (r *columnRepository) FindByBoardID(ctx context.Context, boardID string) ([]*models.Column, error) {
	var columns []*models.Column
	err := r.db.WithContext(ctx).
		Preload("Tasks").
		Preload("Board").
		Where("board_id = ?", boardID).
		Find(&columns).Error
	if err != nil {
		return nil, err
	}
	return columns, nil
}

func (r *columnRepository) Update(ctx context.Context, column *models.Column) error {
	return r.db.WithContext(ctx).Save(column).Error
}

func (r *columnRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&models.Column{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("column with id %s not found", id)
	}
	return nil
}

func (r *columnRepository) SoftDelete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Column{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("column with id %s not found", id)
	}
	return nil
}
