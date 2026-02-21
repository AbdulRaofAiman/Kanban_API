package services

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"kanban-backend/models"
	"kanban-backend/utils"
)

var taskTestIDCounter atomic.Int64

func generateTaskTestID() string {
	counter := taskTestIDCounter.Add(1)
	return "task-" + string(rune(counter))
}

type mockTaskRepository struct {
	tasks map[string]*models.Task
}

func newMockTaskRepository() *mockTaskRepository {
	return &mockTaskRepository{
		tasks: make(map[string]*models.Task),
	}
}

func (m *mockTaskRepository) Create(ctx context.Context, task *models.Task) error {
	if task.ID == "" {
		task.ID = generateTaskTestID()
	}
	m.tasks[task.ID] = task
	return nil
}

func (m *mockTaskRepository) FindByID(ctx context.Context, id string) (*models.Task, error) {
	task, exists := m.tasks[id]
	if !exists {
		return nil, errors.New("task not found")
	}
	taskCopy := *task
	if task.Column != nil {
		columnCopy := *task.Column
		if task.Column.Board != nil {
			boardCopy := *task.Column.Board
			columnCopy.Board = &boardCopy
		}
		taskCopy.Column = &columnCopy
	}
	return &taskCopy, nil
}

func (m *mockTaskRepository) FindByColumnID(ctx context.Context, columnID string) ([]*models.Task, error) {
	var tasks []*models.Task
	for _, task := range m.tasks {
		if task.ColumnID == columnID {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

func (m *mockTaskRepository) Update(ctx context.Context, task *models.Task) error {
	if _, exists := m.tasks[task.ID]; !exists {
		return errors.New("task not found")
	}
	m.tasks[task.ID] = task
	return nil
}

func (m *mockTaskRepository) Delete(ctx context.Context, id string) error {
	if _, exists := m.tasks[id]; !exists {
		return errors.New("task not found")
	}
	delete(m.tasks, id)
	return nil
}

func (m *mockTaskRepository) SoftDelete(ctx context.Context, id string) error {
	return m.Delete(ctx, id)
}

func setupTestColumn(boardID string) *models.Column {
	return &models.Column{
		ID:      generateColumnTestID(),
		BoardID: boardID,
		Title:   "Test Column",
		Order:   1,
		Board: &models.Board{
			ID:     boardID,
			UserID: "user123",
			Title:  "Test Board",
		},
	}
}

func setupTestTask(columnID string) *models.Task {
	return &models.Task{
		ID:          generateTaskTestID(),
		ColumnID:    columnID,
		Title:       "Test Task",
		Description: "Test Description",
		Column:      setupTestColumn("board123"),
	}
}

func TestNewTaskService(t *testing.T) {
	mockTaskRepo := newMockTaskRepository()
	mockColumnRepo := newMockColumnRepository()
	service := NewTaskService(mockTaskRepo, mockColumnRepo)

	if service == nil {
		t.Error("NewTaskService() should return non-nil service")
	}
}

func TestTaskService_Create(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		columnID    string
		title       string
		description string
		deadline    *time.Time
		expectError bool
		errorType   string
	}{
		{
			name:        "Valid task creation",
			userID:      "user123",
			columnID:    "col1",
			title:       "New Task",
			description: "Task description",
			expectError: false,
		},
		{
			name:        "Task creation without description",
			userID:      "user123",
			columnID:    "col1",
			title:       "New Task",
			description: "",
			expectError: false,
		},
		{
			name:        "Task creation in column from different user",
			userID:      "user456",
			columnID:    "col1",
			title:       "New Task",
			description: "",
			expectError: true,
			errorType:   "unauthorized",
		},
		{
			name:        "Task creation in non-existent column",
			userID:      "user123",
			columnID:    "nonexistent",
			title:       "New Task",
			description: "",
			expectError: true,
			errorType:   "not_found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTaskRepo := newMockTaskRepository()
			mockColumnRepo := newMockColumnRepository()
			service := NewTaskService(mockTaskRepo, mockColumnRepo)
			ctx := context.Background()

			column := setupTestColumn("board123")
			column.ID = "col1"
			column.Board.UserID = "user123"
			err := mockColumnRepo.Create(ctx, column)
			if err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			task, err := service.Create(ctx, tt.userID, tt.columnID, tt.title, tt.description, tt.deadline)

			if tt.expectError {
				if err == nil {
					t.Error("Create() expected error but got nil")
					return
				}

				if tt.errorType == "unauthorized" {
					var unauthorizedErr utils.ErrUnauthorized
					if !errors.As(err, &unauthorizedErr) {
						t.Errorf("Create() should return ErrUnauthorized, got %v", err)
					}
				} else if tt.errorType == "not_found" {
					var notFoundErr utils.ErrNotFound
					if !errors.As(err, &notFoundErr) {
						t.Errorf("Create() should return ErrNotFound, got %v", err)
					}
				}
				return
			}

			if err != nil {
				t.Errorf("Create() unexpected error = %v", err)
				return
			}

			if task == nil {
				t.Error("Create() should return non-nil task")
				return
			}

			if task.Title != tt.title {
				t.Errorf("Create() title = %v, want %v", task.Title, tt.title)
			}

			if task.Description != tt.description {
				t.Errorf("Create() description = %v, want %v", task.Description, tt.description)
			}

			if task.ID == "" {
				t.Error("Create() should generate task ID")
			}
		})
	}
}

func TestTaskService_FindByID(t *testing.T) {
	tests := []struct {
		name          string
		setupTask     bool
		taskID        string
		taskUserID    string
		requestUserID string
		expectError   bool
		errorType     string
	}{
		{
			name:          "Valid find by ID",
			setupTask:     true,
			taskID:        "task1",
			taskUserID:    "user123",
			requestUserID: "user123",
			expectError:   false,
		},
		{
			name:          "Find by ID with wrong user",
			setupTask:     true,
			taskID:        "task1",
			taskUserID:    "user123",
			requestUserID: "user456",
			expectError:   true,
			errorType:     "unauthorized",
		},
		{
			name:          "Find non-existent task",
			setupTask:     false,
			taskID:        "nonexistent",
			taskUserID:    "user123",
			requestUserID: "user123",
			expectError:   true,
			errorType:     "not_found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTaskRepo := newMockTaskRepository()
			mockColumnRepo := newMockColumnRepository()
			service := NewTaskService(mockTaskRepo, mockColumnRepo)
			ctx := context.Background()

			if tt.setupTask {
				column := setupTestColumn("board123")
				column.Board.UserID = tt.taskUserID
				task := setupTestTask(column.ID)
				task.ID = tt.taskID
				task.Column = column
				err := mockTaskRepo.Create(ctx, task)
				if err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			task, err := service.FindByID(ctx, tt.taskID, tt.requestUserID)

			if tt.expectError {
				if err == nil {
					t.Error("FindByID() expected error but got nil")
					return
				}

				if tt.errorType == "unauthorized" {
					var unauthorizedErr utils.ErrUnauthorized
					if !errors.As(err, &unauthorizedErr) {
						t.Errorf("FindByID() should return ErrUnauthorized, got %v", err)
					}
				} else if tt.errorType == "not_found" {
					var notFoundErr utils.ErrNotFound
					if !errors.As(err, &notFoundErr) {
						t.Errorf("FindByID() should return ErrNotFound, got %v", err)
					}
				}
				return
			}

			if err != nil {
				t.Errorf("FindByID() unexpected error = %v", err)
				return
			}

			if task == nil {
				t.Error("FindByID() should return non-nil task")
				return
			}

			if task.ID != tt.taskID {
				t.Errorf("FindByID() task ID = %v, want %v", task.ID, tt.taskID)
			}
		})
	}
}

func TestTaskService_FindByColumnID(t *testing.T) {
	tests := []struct {
		name          string
		setupTasks    int
		columnID      string
		columnUserID  string
		requestUserID string
		expectedCount int
		expectError   bool
		errorType     string
	}{
		{
			name:          "Find column's tasks",
			setupTasks:    3,
			columnID:      "col1",
			columnUserID:  "user123",
			requestUserID: "user123",
			expectedCount: 3,
			expectError:   false,
		},
		{
			name:          "Find tasks from column with wrong user",
			setupTasks:    3,
			columnID:      "col1",
			columnUserID:  "user123",
			requestUserID: "user456",
			expectedCount: 0,
			expectError:   true,
			errorType:     "unauthorized",
		},
		{
			name:          "Find tasks from non-existent column",
			setupTasks:    0,
			columnID:      "nonexistent",
			columnUserID:  "user123",
			requestUserID: "user123",
			expectedCount: 0,
			expectError:   true,
			errorType:     "not_found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTaskRepo := newMockTaskRepository()
			mockColumnRepo := newMockColumnRepository()
			service := NewTaskService(mockTaskRepo, mockColumnRepo)
			ctx := context.Background()

			if tt.setupTasks > 0 {
				column := setupTestColumn("board123")
				column.ID = tt.columnID
				column.Board.UserID = tt.columnUserID
				err := mockColumnRepo.Create(ctx, column)
				if err != nil {
					t.Fatalf("Setup failed: %v", err)
				}

				for i := 0; i < tt.setupTasks; i++ {
					task := setupTestTask(tt.columnID)
					task.Column = column
					err := mockTaskRepo.Create(ctx, task)
					if err != nil {
						t.Fatalf("Setup failed: %v", err)
					}
				}
			}

			tasks, err := service.FindByColumnID(ctx, tt.columnID, tt.requestUserID)

			if tt.expectError {
				if err == nil {
					t.Error("FindByColumnID() expected error but got nil")
					return
				}

				if tt.errorType == "unauthorized" {
					var unauthorizedErr utils.ErrUnauthorized
					if !errors.As(err, &unauthorizedErr) {
						t.Errorf("FindByColumnID() should return ErrUnauthorized, got %v", err)
					}
				} else if tt.errorType == "not_found" {
					var notFoundErr utils.ErrNotFound
					if !errors.As(err, &notFoundErr) {
						t.Errorf("FindByColumnID() should return ErrNotFound, got %v", err)
					}
				}
				return
			}

			if err != nil {
				t.Errorf("FindByColumnID() unexpected error = %v", err)
				return
			}

			if len(tasks) != tt.expectedCount {
				t.Errorf("FindByColumnID() returned %d tasks, want %d", len(tasks), tt.expectedCount)
			}
		})
	}
}

func TestTaskService_Update(t *testing.T) {
	tests := []struct {
		name          string
		setupTask     bool
		taskID        string
		taskUserID    string
		requestUserID string
		title         string
		description   string
		deadline      *time.Time
		expectError   bool
		errorType     string
	}{
		{
			name:          "Valid update title",
			setupTask:     true,
			taskID:        "task1",
			taskUserID:    "user123",
			requestUserID: "user123",
			title:         "Updated Title",
			expectError:   false,
		},
		{
			name:          "Valid update description",
			setupTask:     true,
			taskID:        "task1",
			taskUserID:    "user123",
			requestUserID: "user123",
			title:         "",
			description:   "Updated Description",
			expectError:   false,
		},
		{
			name:          "Valid update deadline",
			setupTask:     true,
			taskID:        "task1",
			taskUserID:    "user123",
			requestUserID: "user123",
			deadline:      func() *time.Time { t := time.Now().Add(24 * time.Hour); return &t }(),
			expectError:   false,
		},
		{
			name:          "Update by wrong user",
			setupTask:     true,
			taskID:        "task1",
			taskUserID:    "user123",
			requestUserID: "user456",
			title:         "Updated Title",
			expectError:   true,
			errorType:     "unauthorized",
		},
		{
			name:          "Update non-existent task",
			setupTask:     false,
			taskID:        "nonexistent",
			taskUserID:    "user123",
			requestUserID: "user123",
			title:         "Updated Title",
			expectError:   true,
			errorType:     "not_found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTaskRepo := newMockTaskRepository()
			mockColumnRepo := newMockColumnRepository()
			service := NewTaskService(mockTaskRepo, mockColumnRepo)
			ctx := context.Background()

			if tt.setupTask {
				column := setupTestColumn("board123")
				column.Board.UserID = tt.taskUserID
				task := setupTestTask(column.ID)
				task.ID = tt.taskID
				task.Column = column
				err := mockTaskRepo.Create(ctx, task)
				if err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			task, err := service.Update(ctx, tt.taskID, tt.requestUserID, tt.title, tt.description, tt.deadline)

			if tt.expectError {
				if err == nil {
					t.Error("Update() expected error but got nil")
					return
				}

				if tt.errorType == "unauthorized" {
					var unauthorizedErr utils.ErrUnauthorized
					if !errors.As(err, &unauthorizedErr) {
						t.Errorf("Update() should return ErrUnauthorized, got %v", err)
					}
				} else if tt.errorType == "not_found" {
					var notFoundErr utils.ErrNotFound
					if !errors.As(err, &notFoundErr) {
						t.Errorf("Update() should return ErrNotFound, got %v", err)
					}
				}
				return
			}

			if err != nil {
				t.Errorf("Update() unexpected error = %v", err)
				return
			}

			if task == nil {
				t.Error("Update() should return non-nil task")
				return
			}

			if tt.title != "" && task.Title != tt.title {
				t.Errorf("Update() title = %v, want %v", task.Title, tt.title)
			}

			if tt.description != "" && task.Description != tt.description {
				t.Errorf("Update() description = %v, want %v", task.Description, tt.description)
			}

			if tt.deadline != nil && (task.Deadline == nil || !task.Deadline.Equal(*tt.deadline)) {
				t.Errorf("Update() deadline = %v, want %v", task.Deadline, tt.deadline)
			}
		})
	}
}

func TestTaskService_Delete(t *testing.T) {
	tests := []struct {
		name          string
		setupTask     bool
		taskID        string
		taskUserID    string
		requestUserID string
		expectError   bool
		errorType     string
	}{
		{
			name:          "Valid delete",
			setupTask:     true,
			taskID:        "task1",
			taskUserID:    "user123",
			requestUserID: "user123",
			expectError:   false,
		},
		{
			name:          "Delete by wrong user",
			setupTask:     true,
			taskID:        "task1",
			taskUserID:    "user123",
			requestUserID: "user456",
			expectError:   true,
			errorType:     "unauthorized",
		},
		{
			name:          "Delete non-existent task",
			setupTask:     false,
			taskID:        "nonexistent",
			taskUserID:    "user123",
			requestUserID: "user123",
			expectError:   true,
			errorType:     "not_found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTaskRepo := newMockTaskRepository()
			mockColumnRepo := newMockColumnRepository()
			service := NewTaskService(mockTaskRepo, mockColumnRepo)
			ctx := context.Background()

			if tt.setupTask {
				column := setupTestColumn("board123")
				column.Board.UserID = tt.taskUserID
				task := setupTestTask(column.ID)
				task.ID = tt.taskID
				task.Column = column
				err := mockTaskRepo.Create(ctx, task)
				if err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			err := service.Delete(ctx, tt.taskID, tt.requestUserID)

			if tt.expectError {
				if err == nil {
					t.Error("Delete() expected error but got nil")
					return
				}

				if tt.errorType == "unauthorized" {
					var unauthorizedErr utils.ErrUnauthorized
					if !errors.As(err, &unauthorizedErr) {
						t.Errorf("Delete() should return ErrUnauthorized, got %v", err)
					}
				} else if tt.errorType == "not_found" {
					var notFoundErr utils.ErrNotFound
					if !errors.As(err, &notFoundErr) {
						t.Errorf("Delete() should return ErrNotFound, got %v", err)
					}
				}
				return
			}

			if err != nil {
				t.Errorf("Delete() unexpected error = %v", err)
				return
			}

			_, err = service.FindByID(ctx, tt.taskID, tt.requestUserID)
			if err == nil {
				t.Error("Delete() task should no longer be accessible")
			}
		})
	}
}

func TestTaskService_Move(t *testing.T) {
	tests := []struct {
		name            string
		setupTask       bool
		setupTarget     bool
		targetSameBoard bool
		taskID          string
		taskUserID      string
		requestUserID   string
		targetColumnID  string
		targetUserID    string
		expectError     bool
		errorType       string
	}{
		{
			name:            "Valid move within same board",
			setupTask:       true,
			setupTarget:     true,
			targetSameBoard: true,
			taskID:          "task1",
			taskUserID:      "user123",
			requestUserID:   "user123",
			targetColumnID:  "col2",
			targetUserID:    "user123",
			expectError:     false,
		},
		{
			name:            "Move task not owned by user",
			setupTask:       true,
			setupTarget:     true,
			targetSameBoard: true,
			taskID:          "task1",
			taskUserID:      "user123",
			requestUserID:   "user456",
			targetColumnID:  "col2",
			expectError:     true,
			errorType:       "unauthorized",
		},
		{
			name:            "Move to column not owned by user",
			setupTask:       true,
			setupTarget:     true,
			targetSameBoard: true,
			taskID:          "task1",
			taskUserID:      "user123",
			requestUserID:   "user123",
			targetColumnID:  "col2",
			targetUserID:    "user456",
			expectError:     true,
			errorType:       "unauthorized",
		},
		{
			name:            "Move to column in different board",
			setupTask:       true,
			setupTarget:     true,
			targetSameBoard: false,
			taskID:          "task1",
			taskUserID:      "user123",
			requestUserID:   "user123",
			targetColumnID:  "col2",
			targetUserID:    "user123",
			expectError:     true,
		},
		{
			name:            "Move non-existent task",
			setupTask:       false,
			setupTarget:     true,
			targetSameBoard: true,
			taskID:          "nonexistent",
			requestUserID:   "user123",
			targetColumnID:  "col2",
			expectError:     true,
			errorType:       "not_found",
		},
		{
			name:            "Move to non-existent column",
			setupTask:       true,
			setupTarget:     false,
			targetSameBoard: true,
			taskID:          "task1",
			requestUserID:   "user123",
			targetColumnID:  "nonexistent",
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTaskRepo := newMockTaskRepository()
			mockColumnRepo := newMockColumnRepository()
			service := NewTaskService(mockTaskRepo, mockColumnRepo)
			ctx := context.Background()

			sourceColumn := setupTestColumn("board123")
			sourceColumn.ID = "col1"
			sourceColumn.Board.UserID = tt.taskUserID
			err := mockColumnRepo.Create(ctx, sourceColumn)
			if err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			if tt.setupTarget {
				targetBoardID := "board123"
				if !tt.targetSameBoard {
					targetBoardID = "board456"
				}
				targetColumn := setupTestColumn(targetBoardID)
				targetColumn.ID = tt.targetColumnID
				targetColumn.Board.UserID = tt.targetUserID
				if tt.targetUserID == "" {
					targetColumn.Board.UserID = tt.taskUserID
				}
				err := mockColumnRepo.Create(ctx, targetColumn)
				if err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			if tt.setupTask {
				task := setupTestTask("col1")
				task.ID = tt.taskID
				task.Column = sourceColumn
				err := mockTaskRepo.Create(ctx, task)
				if err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			err = service.Move(ctx, tt.taskID, tt.targetColumnID, tt.requestUserID)

			if tt.expectError {
				if err == nil {
					t.Error("Move() expected error but got nil")
					return
				}

				if tt.errorType == "unauthorized" {
					var unauthorizedErr utils.ErrUnauthorized
					if !errors.As(err, &unauthorizedErr) {
						t.Errorf("Move() should return ErrUnauthorized, got %v", err)
					}
				} else if tt.errorType == "not_found" {
					var notFoundErr utils.ErrNotFound
					if !errors.As(err, &notFoundErr) {
						t.Errorf("Move() should return ErrNotFound, got %v", err)
					}
				}
				return
			}

			if err != nil {
				t.Errorf("Move() unexpected error = %v", err)
				return
			}

			task, err := service.FindByID(ctx, tt.taskID, tt.requestUserID)
			if err != nil {
				t.Fatalf("FindByID() after move failed: %v", err)
			}

			if task.ColumnID != tt.targetColumnID {
				t.Errorf("Move() task column ID = %v, want %v", task.ColumnID, tt.targetColumnID)
			}
		})
	}
}

func TestTaskService_Integration(t *testing.T) {
	mockTaskRepo := newMockTaskRepository()
	mockColumnRepo := newMockColumnRepository()
	service := NewTaskService(mockTaskRepo, mockColumnRepo)
	ctx := context.Background()

	userID := "user123"

	board := &models.Board{
		ID:     "board123",
		Title:  "Test Board",
		UserID: userID,
	}

	column1 := &models.Column{
		ID:      generateColumnTestID(),
		BoardID: board.ID,
		Title:   "To Do",
		Order:   1,
		Board:   board,
	}
	err := mockColumnRepo.Create(ctx, column1)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	column2 := &models.Column{
		ID:      generateColumnTestID(),
		BoardID: board.ID,
		Title:   "In Progress",
		Order:   2,
		Board:   board,
	}
	err = mockColumnRepo.Create(ctx, column2)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	title := "My Task"
	description := "Task description"
	task, err := service.Create(ctx, userID, column1.ID, title, description, nil)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	if task == nil || task.ID == "" {
		t.Fatal("Create() should return task with ID")
	}

	storedTask, _ := mockTaskRepo.FindByID(ctx, task.ID)
	storedTask.Column = column1
	storedTask.Column.Board = board
	mockTaskRepo.Update(ctx, storedTask)

	foundTask, err := service.FindByID(ctx, task.ID, userID)
	if err != nil {
		t.Fatalf("FindByID() failed: %v", err)
	}

	if foundTask.ID != task.ID {
		t.Errorf("FindByID() returned wrong task: got %v, want %v", foundTask.ID, task.ID)
	}

	updatedTitle := "Updated Task"
	updatedTask, err := service.Update(ctx, task.ID, userID, updatedTitle, "", nil)
	if err != nil {
		t.Fatalf("Update() failed: %v", err)
	}

	if updatedTask.Title != updatedTitle {
		t.Errorf("Update() title = %v, want %v", updatedTask.Title, updatedTitle)
	}

	tasks, err := service.FindByColumnID(ctx, column1.ID, userID)
	if err != nil {
		t.Fatalf("FindByColumnID() failed: %v", err)
	}

	if len(tasks) != 1 {
		t.Errorf("FindByColumnID() returned %d tasks, want 1", len(tasks))
	}

	err = service.Move(ctx, task.ID, column2.ID, userID)
	if err != nil {
		t.Fatalf("Move() failed: %v", err)
	}

	movedTask, err := service.FindByID(ctx, task.ID, userID)
	if err != nil {
		t.Fatalf("FindByID() after move failed: %v", err)
	}

	if movedTask.ColumnID != column2.ID {
		t.Errorf("Move() column ID = %v, want %v", movedTask.ColumnID, column2.ID)
	}

	err = service.Delete(ctx, task.ID, userID)
	if err != nil {
		t.Fatalf("Delete() failed: %v", err)
	}

	_, err = service.FindByID(ctx, task.ID, userID)
	if err == nil {
		t.Error("Task should not be found after deletion")
	}
}
