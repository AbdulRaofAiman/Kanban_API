package services

import (
	"context"
	"kanban-backend/models"
	"kanban-backend/repositories"
	"kanban-backend/utils"
)

type CommentService interface {
	Create(ctx context.Context, taskID, userID, content string) (*models.Comment, error)
	FindByID(ctx context.Context, id, userID string) (*models.Comment, error)
	FindByTaskID(ctx context.Context, taskID, userID string) ([]*models.Comment, error)
	FindByTaskIDWithPagination(ctx context.Context, taskID, userID string, page, limit int) ([]*models.Comment, int, error)
	Update(ctx context.Context, id, userID, content string) (*models.Comment, error)
	Delete(ctx context.Context, id, userID string) error
}

type commentService struct {
	commentRepo repositories.CommentRepository
	taskRepo    repositories.TaskRepository
}

func NewCommentService(commentRepo repositories.CommentRepository, taskRepo repositories.TaskRepository) CommentService {
	return &commentService{
		commentRepo: commentRepo,
		taskRepo:    taskRepo,
	}
}

func (s *commentService) Create(ctx context.Context, taskID, userID, content string) (*models.Comment, error) {
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

	comment := &models.Comment{
		TaskID:  taskID,
		UserID:  userID,
		Content: content,
	}

	err = s.commentRepo.Create(ctx, comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *commentService) FindByID(ctx context.Context, id, userID string) (*models.Comment, error) {
	comment, err := s.commentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFound("comment not found")
	}

	if comment.UserID != userID {
		return nil, utils.NewUnauthorized("you do not have access to this comment")
	}

	return comment, nil
}

func (s *commentService) FindByTaskID(ctx context.Context, taskID, userID string) ([]*models.Comment, error) {
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

	comments, err := s.commentRepo.FindByTaskID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *commentService) Update(ctx context.Context, id, userID, content string) (*models.Comment, error) {
	comment, err := s.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	comment.Content = content

	err = s.commentRepo.Update(ctx, comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *commentService) Delete(ctx context.Context, id, userID string) error {
	_, err := s.FindByID(ctx, id, userID)
	if err != nil {
		return err
	}

	err = s.commentRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *commentService) FindByTaskIDWithPagination(ctx context.Context, taskID, userID string, page, limit int) ([]*models.Comment, int, error) {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, 0, utils.NewNotFound("task not found")
	}

	if task.Column == nil {
		return nil, 0, utils.NewNotFound("column not found for task")
	}

	if task.Column.Board.UserID != userID {
		return nil, 0, utils.NewUnauthorized("you do not have access to this task")
	}

	comments, total, err := s.commentRepo.FindByTaskIDWithPagination(ctx, taskID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}
