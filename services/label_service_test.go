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

type mockLabelRepository struct {
	labels map[string]*models.Label
}

func newMockLabelRepository() *mockLabelRepository {
	return &mockLabelRepository{
		labels: make(map[string]*models.Label),
	}
}

func (m *mockLabelRepository) Create(ctx context.Context, label *models.Label) error {
	if label.ID == "" {
		label.ID = fmt.Sprintf("label-%d-%d", time.Now().UnixNano(), rand.Intn(1000))
	}
	label.CreatedAt = time.Now()
	label.UpdatedAt = time.Now()
	m.labels[label.ID] = label
	return nil
}

func (m *mockLabelRepository) FindByID(ctx context.Context, id string) (*models.Label, error) {
	label, exists := m.labels[id]
	if !exists {
		return nil, utils.NewNotFound("label not found")
	}
	return label, nil
}

func (m *mockLabelRepository) FindAll(ctx context.Context) ([]*models.Label, error) {
	var labels []*models.Label
	for _, label := range m.labels {
		labels = append(labels, label)
	}
	return labels, nil
}

func (m *mockLabelRepository) Update(ctx context.Context, label *models.Label) error {
	if _, exists := m.labels[label.ID]; !exists {
		return utils.NewNotFound("label not found")
	}
	label.UpdatedAt = time.Now()
	m.labels[label.ID] = label
	return nil
}

func (m *mockLabelRepository) Delete(ctx context.Context, id string) error {
	if _, exists := m.labels[id]; !exists {
		return utils.NewNotFound("label not found")
	}
	delete(m.labels, id)
	return nil
}

func (m *mockLabelRepository) SoftDelete(ctx context.Context, id string) error {
	return m.Delete(ctx, id)
}

type mockTaskRepositoryForLabel struct {
	tasks map[string]*models.Task
}

func newMockTaskRepositoryForLabel() *mockTaskRepositoryForLabel {
	return &mockTaskRepositoryForLabel{
		tasks: make(map[string]*models.Task),
	}
}

func (m *mockTaskRepositoryForLabel) FindByID(ctx context.Context, id string) (*models.Task, error) {
	task, exists := m.tasks[id]
	if !exists {
		return nil, utils.NewNotFound("task not found")
	}
	return task, nil
}

func (m *mockTaskRepositoryForLabel) Create(ctx context.Context, task *models.Task) error {
	if task.ID == "" {
		task.ID = generateTestID()
	}
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	m.tasks[task.ID] = task
	return nil
}

func (m *mockTaskRepositoryForLabel) FindByColumnID(ctx context.Context, columnID string) ([]*models.Task, error) {
	return nil, nil
}

func (m *mockTaskRepositoryForLabel) Update(ctx context.Context, task *models.Task) error {
	return nil
}

func (m *mockTaskRepositoryForLabel) Delete(ctx context.Context, id string) error {
	return nil
}

func (m *mockTaskRepositoryForLabel) SoftDelete(ctx context.Context, id string) error {
	return nil
}

func TestNewLabelService(t *testing.T) {
	labelRepo := newMockLabelRepository()
	taskRepo := newMockTaskRepositoryForLabel()
	service := NewLabelService(labelRepo, taskRepo)

	if service == nil {
		t.Error("NewLabelService() should return non-nil service")
	}
}

func TestLabelService_Create(t *testing.T) {
	labelRepo := newMockLabelRepository()
	taskRepo := newMockTaskRepositoryForLabel()
	service := NewLabelService(labelRepo, taskRepo)

	label, err := service.Create(context.Background(), "Bug", "#FF0000")
	if err != nil {
		t.Errorf("Create() error = %v", err)
		return
	}

	if label.Name != "Bug" {
		t.Errorf("Create() name = %v, want %v", label.Name, "Bug")
	}

	if label.Color != "#FF0000" {
		t.Errorf("Create() color = %v, want %v", label.Color, "#FF0000")
	}
}

func TestLabelService_Create_ValidationError(t *testing.T) {
	labelRepo := newMockLabelRepository()
	taskRepo := newMockTaskRepositoryForLabel()
	service := NewLabelService(labelRepo, taskRepo)

	_, err := service.Create(context.Background(), "", "#FF0000")
	if err == nil {
		t.Error("Create() should return error for empty name")
	}
}

func TestLabelService_FindAll(t *testing.T) {
	labelRepo := newMockLabelRepository()
	taskRepo := newMockTaskRepositoryForLabel()
	service := NewLabelService(labelRepo, taskRepo)

	service.Create(context.Background(), "Bug", "#FF0000")
	service.Create(context.Background(), "Feature", "#00FF00")

	labels, err := service.FindAll(context.Background())
	if err != nil {
		t.Errorf("FindAll() error = %v", err)
		return
	}

	if len(labels) != 2 {
		t.Errorf("FindAll() returned %d labels, want 2", len(labels))
	}
}

func TestLabelService_Update(t *testing.T) {
	labelRepo := newMockLabelRepository()
	taskRepo := newMockTaskRepositoryForLabel()
	service := NewLabelService(labelRepo, taskRepo)

	label, _ := service.Create(context.Background(), "Bug", "#FF0000")

	updatedLabel, err := service.Update(context.Background(), label.ID, "Critical", "#FF0000")
	if err != nil {
		t.Errorf("Update() error = %v", err)
		return
	}

	if updatedLabel.Name != "Critical" {
		t.Errorf("Update() name = %v, want %v", updatedLabel.Name, "Critical")
	}
}

func TestLabelService_Delete(t *testing.T) {
	labelRepo := newMockLabelRepository()
	taskRepo := newMockTaskRepositoryForLabel()
	service := NewLabelService(labelRepo, taskRepo)

	label, _ := service.Create(context.Background(), "Bug", "#FF0000")

	err := service.Delete(context.Background(), label.ID)
	if err != nil {
		t.Errorf("Delete() error = %v", err)
		return
	}

	_, err = service.FindByID(context.Background(), label.ID)
	if err == nil {
		t.Error("Delete() should have removed the label")
	}
}

func TestLabelService_AddToTask_Unauthorized(t *testing.T) {
	labelRepo := newMockLabelRepository()
	taskRepo := newMockTaskRepositoryForLabel()
	service := NewLabelService(labelRepo, taskRepo)

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

	label, _ := service.Create(context.Background(), "Bug", "#FF0000")

	err := service.AddToTask(context.Background(), "task-1", label.ID, userID)
	if err == nil {
		t.Error("AddToTask() should return error for unauthorized user")
	}
}

func TestLabelService_RemoveFromTask_Unauthorized(t *testing.T) {
	labelRepo := newMockLabelRepository()
	taskRepo := newMockTaskRepositoryForLabel()
	service := NewLabelService(labelRepo, taskRepo)

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

	label, _ := service.Create(context.Background(), "Bug", "#FF0000")

	err := service.RemoveFromTask(context.Background(), "task-1", label.ID, userID)
	if err == nil {
		t.Error("RemoveFromTask() should return error for unauthorized user")
	}
}
