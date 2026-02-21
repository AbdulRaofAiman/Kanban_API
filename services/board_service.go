package services

import (
	"context"

	"kanban-backend/models"
	"kanban-backend/repositories"
	"kanban-backend/utils"
)

type BoardService interface {
	Create(ctx context.Context, userID, title, color string) (*models.Board, error)
	FindByID(ctx context.Context, boardID, userID string) (*models.Board, error)
	FindByUserID(ctx context.Context, userID string) ([]*models.Board, error)
	Update(ctx context.Context, boardID, userID, title, color string) (*models.Board, error)
	Delete(ctx context.Context, boardID, userID string) error
}

type boardService struct {
	boardRepo  repositories.BoardRepository
	columnRepo repositories.ColumnRepository
}

func NewBoardService(boardRepo repositories.BoardRepository, columnRepo repositories.ColumnRepository) BoardService {
	return &boardService{
		boardRepo:  boardRepo,
		columnRepo: columnRepo,
	}
}

func (s *boardService) Create(ctx context.Context, userID, title, color string) (*models.Board, error) {
	board := &models.Board{
		Title:  title,
		UserID: userID,
		Color:  color,
	}

	err := s.boardRepo.Create(ctx, board)
	if err != nil {
		return nil, err
	}

	defaultColumns := []models.Column{
		{Title: "To Do", Order: 1, BoardID: board.ID},
		{Title: "In Progress", Order: 2, BoardID: board.ID},
		{Title: "Done", Order: 3, BoardID: board.ID},
	}

	for _, col := range defaultColumns {
		err = s.columnRepo.Create(ctx, &col)
		if err != nil {
			return nil, err
		}
	}

	board.Columns = defaultColumns

	return board, nil
}

func (s *boardService) FindByID(ctx context.Context, boardID, userID string) (*models.Board, error) {
	board, err := s.boardRepo.FindByID(ctx, boardID)
	if err != nil {
		return nil, utils.NewNotFound("board not found")
	}

	if board.UserID != userID {
		return nil, utils.NewUnauthorized("you do not have access to this board")
	}

	return board, nil
}

func (s *boardService) FindByUserID(ctx context.Context, userID string) ([]*models.Board, error) {
	boards, err := s.boardRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return boards, nil
}

func (s *boardService) Update(ctx context.Context, boardID, userID, title, color string) (*models.Board, error) {
	board, err := s.FindByID(ctx, boardID, userID)
	if err != nil {
		return nil, err
	}

	if title != "" {
		board.Title = title
	}

	if color != "" {
		board.Color = color
	}

	err = s.boardRepo.Update(ctx, board)
	if err != nil {
		return nil, err
	}

	return board, nil
}

func (s *boardService) Delete(ctx context.Context, boardID, userID string) error {
	_, err := s.FindByID(ctx, boardID, userID)
	if err != nil {
		return err
	}

	err = s.boardRepo.Delete(ctx, boardID)
	if err != nil {
		return err
	}

	return nil
}
