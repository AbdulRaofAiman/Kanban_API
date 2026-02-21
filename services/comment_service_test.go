package services

import (
	"context"
	"testing"
	"time"

	"kanban-backend/models"
	"kanban-backend/utils"
)

type mockCommentRepository struct {
	comments map[string]*models.Comment
}

func newMockCommentRepository() *mockCommentRepository {
	return &mockCommentRepository{
		comments: make(map[string]*models.Comment),
	}
}

func (m *mockCommentRepository) Create(ctx context.Context, comment *models.Comment) error {
	if comment.ID == "" {
		comment.ID = generateTestID()
	}
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	m.comments[comment.ID] = comment
	return nil
}

func (m *mockCommentRepository) FindByID(ctx context.Context, id string) (*models.Comment, error) {
	comment, exists := m.comments[id]
	if !exists {
		return nil, utils.NewNotFound("comment not found")
	}
	return comment, nil
}

func (m *mockCommentRepository) FindByTaskID(ctx context.Context, taskID string) ([]*models.Comment, error) {
	var comments []*models.Comment
	for _, comment := range m.comments {
		if comment.TaskID == taskID {
			comments = append(comments, comment)
		}
	}
	return comments, nil
}

func (m *mockCommentRepository) Update(ctx context.Context, comment *models.Comment) error {
	if _, exists := m.comments[comment.ID]; !exists {
		return utils.NewNotFound("comment not found")
	}
	comment.UpdatedAt = time.Now()
	m.comments[comment.ID] = comment
	return nil
}

func (m *mockCommentRepository) Delete(ctx context.Context, id string) error {
	if _, exists := m.comments[id]; !exists {
		return utils.NewNotFound("comment not found")
	}
	delete(m.comments, id)
	return nil
}

func (m *mockCommentRepository) SoftDelete(ctx context.Context, id string) error {
	return m.Delete(ctx, id)
}

type mockTaskRepositoryForComment struct {
	tasks map[string]*models.Task
}

func newMockTaskRepositoryForComment() *mockTaskRepositoryForComment {
	return &mockTaskRepositoryForComment{
		tasks: make(map[string]*models.Task),
	}
}

func (m *mockTaskRepositoryForComment) FindByID(ctx context.Context, id string) (*models.Task, error) {
	task, exists := m.tasks[id]
	if !exists {
		return nil, utils.NewNotFound("task not found")
	}
	return task, nil
}

func (m *mockTaskRepositoryForComment) Create(ctx context.Context, task *models.Task) error {
	if task.ID == "" {
		task.ID = generateTestID()
	}
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	m.tasks[task.ID] = task
	return nil
}

func (m *mockTaskRepositoryForComment) FindByColumnID(ctx context.Context, columnID string) ([]*models.Task, error) {
	return nil, nil
}

func (m *mockTaskRepositoryForComment) Update(ctx context.Context, task *models.Task) error {
	return nil
}

func (m *mockTaskRepositoryForComment) Delete(ctx context.Context, id string) error {
	return nil
}

func (m *mockTaskRepositoryForComment) SoftDelete(ctx context.Context, id string) error {
	return nil
}

func TestNewCommentService(t *testing.T) {
	commentRepo := newMockCommentRepository()
	taskRepo := newMockTaskRepositoryForComment()
	service := NewCommentService(commentRepo, taskRepo)

	if service == nil {
		t.Error("NewCommentService() should return non-nil service")
	}
}

func TestCommentService_Create(t *testing.T) {
	commentRepo := newMockCommentRepository()
	taskRepo := newMockTaskRepositoryForComment()
	service := NewCommentService(commentRepo, taskRepo)

	userID := "user-1"
	task := &models.Task{
		ID:       "task-1",
		ColumnID: "col-1",
		Title:    "Test Task",
		Column: &models.Column{
			ID:      "col-1",
			BoardID: "board-1",
			Board: &models.Board{
				ID:     "board-1",
				UserID: userID,
			},
		},
	}
	taskRepo.Create(context.Background(), task)

	comment, err := service.Create(context.Background(), "task-1", userID, "Test comment")
	if err != nil {
		t.Errorf("Create() error = %v", err)
		return
	}

	if comment.Content != "Test comment" {
		t.Errorf("Create() content = %v, want %v", comment.Content, "Test comment")
	}
}

func TestCommentService_Create_Unauthorized(t *testing.T) {
	commentRepo := newMockCommentRepository()
	taskRepo := newMockTaskRepositoryForComment()
	service := NewCommentService(commentRepo, taskRepo)

	userID := "user-1"
	task := &models.Task{
		ID:       "task-1",
		ColumnID: "col-1",
		Title:    "Test Task",
		Column: &models.Column{
			ID:      "col-1",
			BoardID: "board-1",
			Board: &models.Board{
				ID:     "board-1",
				UserID: "user-2",
			},
		},
	}
	taskRepo.Create(context.Background(), task)

	_, err := service.Create(context.Background(), "task-1", userID, "Test comment")
	if err == nil {
		t.Error("Create() should return error for unauthorized user")
	}
}

func TestCommentService_Update(t *testing.T) {
	commentRepo := newMockCommentRepository()
	taskRepo := newMockTaskRepositoryForComment()
	service := NewCommentService(commentRepo, taskRepo)

	userID := "user-1"
	task := &models.Task{
		ID:       "task-1",
		ColumnID: "col-1",
		Title:    "Test Task",
		Column: &models.Column{
			ID:      "col-1",
			BoardID: "board-1",
			Board: &models.Board{
				ID:     "board-1",
				UserID: userID,
			},
		},
	}
	taskRepo.Create(context.Background(), task)

	comment, _ := service.Create(context.Background(), "task-1", userID, "Original comment")

	updatedComment, err := service.Update(context.Background(), comment.ID, userID, "Updated comment")
	if err != nil {
		t.Errorf("Update() error = %v", err)
		return
	}

	if updatedComment.Content != "Updated comment" {
		t.Errorf("Update() content = %v, want %v", updatedComment.Content, "Updated comment")
	}
}
