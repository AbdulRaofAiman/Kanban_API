package repositories

import (
	"context"
	"fmt"

	"kanban-backend/config"
	"kanban-backend/models"

	"gorm.io/gorm"
)

type BoardRepository interface {
	Create(ctx context.Context, board *models.Board) error
	FindByID(ctx context.Context, id string) (*models.Board, error)
	FindByUserID(ctx context.Context, userID string) ([]*models.Board, error)
	Update(ctx context.Context, board *models.Board) error
	Delete(ctx context.Context, id string) error
	SoftDelete(ctx context.Context, id string) error
}

type boardRepository struct {
	db *gorm.DB
}

func NewBoardRepository() BoardRepository {
	return &boardRepository{
		db: config.DB,
	}
}

func (r *boardRepository) Create(ctx context.Context, board *models.Board) error {
	return r.db.WithContext(ctx).Create(board).Error
}

func (r *boardRepository) FindByID(ctx context.Context, id string) (*models.Board, error) {
	var board models.Board
	err := r.db.WithContext(ctx).
		Preload("Columns").
		Preload("Members").
		Preload("User").
		Where("id = ?", id).
		First(&board).Error
	if err != nil {
		return nil, err
	}
	return &board, nil
}

func (r *boardRepository) FindByUserID(ctx context.Context, userID string) ([]*models.Board, error) {
	var boards []*models.Board
	err := r.db.WithContext(ctx).
		Preload("Columns").
		Preload("Members").
		Preload("User").
		Where("user_id = ?", userID).
		Find(&boards).Error
	if err != nil {
		return nil, err
	}
	return boards, nil
}

func (r *boardRepository) Update(ctx context.Context, board *models.Board) error {
	return r.db.WithContext(ctx).Save(board).Error
}

func (r *boardRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&models.Board{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("board with id %s not found", id)
	}
	return nil
}

func (r *boardRepository) SoftDelete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Board{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("board with id %s not found", id)
	}
	return nil
}
