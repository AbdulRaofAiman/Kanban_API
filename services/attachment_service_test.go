package services

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"kanban-backend/models"
	"kanban-backend/utils"
)

type mockAttachmentRepository struct {
	attachments map[string]*models.Attachment
	taskRepo    *mockTaskRepositoryForAttachment
}

func newMockAttachmentRepository(taskRepo *mockTaskRepositoryForAttachment) *mockAttachmentRepository {
	return &mockAttachmentRepository{
		attachments: make(map[string]*models.Attachment),
		taskRepo:    taskRepo,
	}
}

func (m *mockAttachmentRepository) Create(ctx context.Context, attachment *models.Attachment) error {
	if attachment.ID == "" {
		attachment.ID = fmt.Sprintf("attachment-%d-%d", time.Now().UnixNano(), rand.Intn(1000))
	}
	attachment.CreatedAt = time.Now()
	attachment.UpdatedAt = time.Now()
	m.attachments[attachment.ID] = attachment
	return nil
}

func (m *mockAttachmentRepository) FindByID(ctx context.Context, id string) (*models.Attachment, error) {
	attachment, exists := m.attachments[id]
	if !exists {
		return nil, utils.NewNotFound("attachment not found")
	}

	task, _ := m.taskRepo.FindByID(ctx, attachment.TaskID)
	attachmentCopy := &models.Attachment{
		ID:       attachment.ID,
		TaskID:   attachment.TaskID,
		FileName: attachment.FileName,
		FileURL:  attachment.FileURL,
		FileSize: attachment.FileSize,
		Task:     task,
	}

	return attachmentCopy, nil
}

func (m *mockAttachmentRepository) FindByTaskID(ctx context.Context, taskID string) ([]*models.Attachment, error) {
	var attachments []*models.Attachment
	for _, attachment := range m.attachments {
		if attachment.TaskID == taskID {
			task, _ := m.taskRepo.FindByID(ctx, attachment.TaskID)
			attachments = append(attachments, &models.Attachment{
				ID:       attachment.ID,
				TaskID:   attachment.TaskID,
				FileName: attachment.FileName,
				FileURL:  attachment.FileURL,
				FileSize: attachment.FileSize,
				Task:     task,
			})
		}
	}
	return attachments, nil
}

func (m *mockAttachmentRepository) Update(ctx context.Context, attachment *models.Attachment) error {
	if _, exists := m.attachments[attachment.ID]; !exists {
		return utils.NewNotFound("attachment not found")
	}
	attachment.UpdatedAt = time.Now()
	m.attachments[attachment.ID] = attachment
	return nil
}

func (m *mockAttachmentRepository) Delete(ctx context.Context, id string) error {
	if _, exists := m.attachments[id]; !exists {
		return utils.NewNotFound("attachment not found")
	}
	delete(m.attachments, id)
	return nil
}

func (m *mockAttachmentRepository) SoftDelete(ctx context.Context, id string) error {
	return m.Delete(ctx, id)
}

func (m *mockAttachmentRepository) FindByTaskIDWithPagination(ctx context.Context, taskID string, offset, limit int) ([]*models.Attachment, int, error) {
	var attachments []*models.Attachment
	for _, attachment := range m.attachments {
		if attachment.TaskID == taskID {
			task, _ := m.taskRepo.FindByID(ctx, attachment.TaskID)
			attachments = append(attachments, &models.Attachment{
				ID:       attachment.ID,
				TaskID:   attachment.TaskID,
				FileName: attachment.FileName,
				FileURL:  attachment.FileURL,
				FileSize: attachment.FileSize,
				Task:     task,
			})
		}
	}
	return attachments, len(attachments), nil
}

type mockTaskRepositoryForAttachment struct {
	tasks map[string]*models.Task
}

func newMockTaskRepositoryForAttachment() *mockTaskRepositoryForAttachment {
	return &mockTaskRepositoryForAttachment{
		tasks: make(map[string]*models.Task),
	}
}

func (m *mockTaskRepositoryForAttachment) FindByID(ctx context.Context, id string) (*models.Task, error) {
	task, exists := m.tasks[id]
	if !exists {
		return nil, utils.NewNotFound("task not found")
	}

	taskCopy := &models.Task{
		ID:       task.ID,
		ColumnID: task.ColumnID,
		Title:    task.Title,
		Column:   task.Column,
	}

	return taskCopy, nil
}

func (m *mockTaskRepositoryForAttachment) Create(ctx context.Context, task *models.Task) error {
	if task.ID == "" {
		task.ID = generateTestID()
	}
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	m.tasks[task.ID] = task
	return nil
}

func (m *mockTaskRepositoryForAttachment) FindByColumnID(ctx context.Context, columnID string) ([]*models.Task, error) {
	return nil, nil
}

func (m *mockTaskRepositoryForAttachment) FindByColumnIDWithFilters(ctx context.Context, columnID string, title string, offset, limit int) ([]*models.Task, int, error) {
	return nil, 0, nil
}

func (m *mockTaskRepositoryForAttachment) Search(ctx context.Context, boardID string, keyword string, offset, limit int) ([]*models.Task, int, error) {
	return nil, 0, nil
}

func (m *mockTaskRepositoryForAttachment) Update(ctx context.Context, task *models.Task) error {
	return nil
}

func (m *mockTaskRepositoryForAttachment) Delete(ctx context.Context, id string) error {
	return nil
}

func (m *mockTaskRepositoryForAttachment) SoftDelete(ctx context.Context, id string) error {
	return nil
}

func TestNewAttachmentService(t *testing.T) {
	taskRepo := newMockTaskRepositoryForAttachment()
	attachmentRepo := newMockAttachmentRepository(taskRepo)
	service := NewAttachmentService(attachmentRepo, taskRepo)

	if service == nil {
		t.Error("NewAttachmentService() should return non-nil service")
	}
}

func TestAttachmentService_Create(t *testing.T) {
	taskRepo := newMockTaskRepositoryForAttachment()
	attachmentRepo := newMockAttachmentRepository(taskRepo)
	service := NewAttachmentService(attachmentRepo, taskRepo)

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

	attachment, err := service.Create(context.Background(), "task-1", userID, "file.pdf", "https://example.com/file.pdf", 1024)
	if err != nil {
		t.Errorf("Create() error = %v", err)
		return
	}

	if attachment.FileName != "file.pdf" {
		t.Errorf("Create() file_name = %v, want %v", attachment.FileName, "file.pdf")
	}

	if attachment.FileURL != "https://example.com/file.pdf" {
		t.Errorf("Create() file_url = %v, want %v", attachment.FileURL, "https://example.com/file.pdf")
	}

	if attachment.FileSize != 1024 {
		t.Errorf("Create() file_size = %v, want %v", attachment.FileSize, 1024)
	}
}

func TestAttachmentService_Create_Unauthorized(t *testing.T) {
	taskRepo := newMockTaskRepositoryForAttachment()
	attachmentRepo := newMockAttachmentRepository(taskRepo)
	service := NewAttachmentService(attachmentRepo, taskRepo)

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

	_, err := service.Create(context.Background(), "task-1", userID, "file.pdf", "https://example.com/file.pdf", 1024)
	if err == nil {
		t.Error("Create() should return error for unauthorized user")
	}
}

func TestAttachmentService_FindByTaskID(t *testing.T) {
	taskRepo := newMockTaskRepositoryForAttachment()
	attachmentRepo := newMockAttachmentRepository(taskRepo)
	service := NewAttachmentService(attachmentRepo, taskRepo)

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

	service.Create(context.Background(), "task-1", userID, "file1.pdf", "https://example.com/file1.pdf", 1024)
	service.Create(context.Background(), "task-1", userID, "file2.pdf", "https://example.com/file2.pdf", 2048)

	attachments, err := service.FindByTaskID(context.Background(), "task-1", userID)
	if err != nil {
		t.Errorf("FindByTaskID() error = %v", err)
		return
	}

	if len(attachments) != 2 {
		t.Errorf("FindByTaskID() returned %d attachments, want 2", len(attachments))
	}
}

func TestAttachmentService_Update(t *testing.T) {
	taskRepo := newMockTaskRepositoryForAttachment()
	attachmentRepo := newMockAttachmentRepository(taskRepo)
	service := NewAttachmentService(attachmentRepo, taskRepo)

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

	attachment, _ := service.Create(context.Background(), "task-1", userID, "file.pdf", "https://example.com/file.pdf", 1024)

	updatedAttachment, err := service.Update(context.Background(), attachment.ID, userID, "updated.pdf", "https://example.com/updated.pdf", 2048)
	if err != nil {
		t.Errorf("Update() error = %v", err)
		return
	}

	if updatedAttachment.FileName != "updated.pdf" {
		t.Errorf("Update() file_name = %v, want %v", updatedAttachment.FileName, "updated.pdf")
	}

	if updatedAttachment.FileURL != "https://example.com/updated.pdf" {
		t.Errorf("Update() file_url = %v, want %v", updatedAttachment.FileURL, "https://example.com/updated.pdf")
	}

	if updatedAttachment.FileSize != 2048 {
		t.Errorf("Update() file_size = %v, want %v", updatedAttachment.FileSize, 2048)
	}
}

func TestAttachmentService_Delete(t *testing.T) {
	taskRepo := newMockTaskRepositoryForAttachment()
	attachmentRepo := newMockAttachmentRepository(taskRepo)
	service := NewAttachmentService(attachmentRepo, taskRepo)

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

	attachment, _ := service.Create(context.Background(), "task-1", userID, "file.pdf", "https://example.com/file.pdf", 1024)

	err := service.Delete(context.Background(), attachment.ID, userID)
	if err != nil {
		t.Errorf("Delete() error = %v", err)
		return
	}

	_, err = service.FindByID(context.Background(), attachment.ID, userID)
	if err == nil {
		t.Error("Delete() should have removed the attachment")
	}
}
