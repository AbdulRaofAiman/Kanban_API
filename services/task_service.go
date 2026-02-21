package services

import (
	"context"
	"errors"
	"time"

	"kanban-backend/models"
	"kanban-backend/repositories"
	"kanban-backend/utils"
)

type TaskService interface {
	Create(ctx context.Context, userID, columnID, title, description string, deadline *time.Time) (*models.Task, error)
	FindByID(ctx context.Context, taskID, userID string) (*models.Task, error)
	FindByColumnID(ctx context.Context, columnID, userID string) ([]*models.Task, error)
	FindByColumnIDWithFilters(ctx context.Context, columnID, userID string, title string, page, limit int) ([]*models.Task, int, error)
	Search(ctx context.Context, boardID, userID string, keyword string, page, limit int) ([]*models.Task, int, error)
	Update(ctx context.Context, taskID, userID, title, description string, deadline *time.Time) (*models.Task, error)
	Delete(ctx context.Context, taskID, userID string) error
	Move(ctx context.Context, taskID, columnID, userID string) error
}

type taskService struct {
	taskRepo   repositories.TaskRepository
	columnRepo repositories.ColumnRepository
}

func NewTaskService(taskRepo repositories.TaskRepository, columnRepo repositories.ColumnRepository) TaskService {
	return &taskService{
		taskRepo:   taskRepo,
		columnRepo: columnRepo,
	}
}

func (s *taskService) Create(ctx context.Context, userID, columnID, title, description string, deadline *time.Time) (*models.Task, error) {
	column, err := s.columnRepo.FindByID(ctx, columnID)
	if err != nil {
		return nil, utils.NewNotFound("column not found")
	}

	if column.Board.UserID != userID {
		return nil, utils.NewUnauthorized("you do not have access to this column")
	}

	task := &models.Task{
		ColumnID:    columnID,
		Title:       title,
		Description: description,
		Deadline:    deadline,
	}

	err = s.taskRepo.Create(ctx, task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *taskService) FindByID(ctx context.Context, taskID, userID string) (*models.Task, error) {
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

	return task, nil
}

func (s *taskService) FindByColumnID(ctx context.Context, columnID, userID string) ([]*models.Task, error) {
	column, err := s.columnRepo.FindByID(ctx, columnID)
	if err != nil {
		return nil, utils.NewNotFound("column not found")
	}

	if column.Board.UserID != userID {
		return nil, utils.NewUnauthorized("you do not have access to this column")
	}

	tasks, err := s.taskRepo.FindByColumnID(ctx, columnID)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *taskService) Update(ctx context.Context, taskID, userID, title, description string, deadline *time.Time) (*models.Task, error) {
	task, err := s.FindByID(ctx, taskID, userID)
	if err != nil {
		return nil, err
	}

	if title != "" {
		task.Title = title
	}

	if description != "" {
		task.Description = description
	}

	if deadline != nil {
		task.Deadline = deadline
	}

	err = s.taskRepo.Update(ctx, task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *taskService) Delete(ctx context.Context, taskID, userID string) error {
	_, err := s.FindByID(ctx, taskID, userID)
	if err != nil {
		return err
	}

	err = s.taskRepo.Delete(ctx, taskID)
	if err != nil {
		return err
	}

	return nil
}

func (s *taskService) Move(ctx context.Context, taskID, columnID, userID string) error {
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

	column, err := s.columnRepo.FindByID(ctx, columnID)
	if err != nil {
		return utils.NewNotFound("target column not found")
	}

	if column.Board.UserID != userID {
		return utils.NewUnauthorized("you do not have access to the target column")
	}

	if task.Column.BoardID != column.BoardID {
		return errors.New("cannot move task to a different board")
	}

	task.ColumnID = columnID
	return s.taskRepo.Update(ctx, task)
}

func (s *taskService) FindByColumnIDWithFilters(ctx context.Context, columnID, userID string, title string, page, limit int) ([]*models.Task, int, error) {
	column, err := s.columnRepo.FindByID(ctx, columnID)
	if err != nil {
		return nil, 0, utils.NewNotFound("column not found")
	}

	if column.Board.UserID != userID {
		return nil, 0, utils.NewUnauthorized("you do not have access to this column")
	}

	tasks, total, err := s.taskRepo.FindByColumnIDWithFilters(ctx, columnID, title, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (s *taskService) Search(ctx context.Context, boardID, userID string, keyword string, page, limit int) ([]*models.Task, int, error) {
	boardRepo := repositories.NewBoardRepository()
	board, err := boardRepo.FindByID(ctx, boardID)
	if err != nil {
		return nil, 0, utils.NewNotFound("board not found")
	}

	if board.UserID != userID {
		return nil, 0, utils.NewUnauthorized("you do not have access to this board")
	}

	tasks, total, err := s.taskRepo.Search(ctx, boardID, keyword, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}
