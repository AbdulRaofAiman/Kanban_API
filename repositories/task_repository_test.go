package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"kanban-backend/models"
)

func TestTaskRepository_Create(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &taskRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)
	column := createTestColumn(db, board.ID)

	deadline := time.Now().Add(24 * time.Hour)
	testCases := []struct {
		name      string
		task      *models.Task
		wantError bool
	}{
		{
			name: "Valid task",
			task: &models.Task{
				ColumnID:    column.ID,
				Title:       "Test Task",
				Description: "Test Description",
				Deadline:    &deadline,
			},
			wantError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.Create(ctx, tc.task)

			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tc.task.ID)
				assert.Equal(t, "Test Task", tc.task.Title)
				assert.Equal(t, column.ID, tc.task.ColumnID)
				assert.Equal(t, "Test Description", tc.task.Description)
				assert.NotNil(t, tc.task.Deadline)
			}
		})
	}
}

func TestTaskRepository_FindByID(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &taskRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)
	column := createTestColumn(db, board.ID)

	testTask := &models.Task{
		ColumnID:    column.ID,
		Title:       "Test Task",
		Description: "Test Description",
	}
	err := repo.Create(ctx, testTask)
	require.NoError(t, err)

	testCases := []struct {
		name      string
		id        string
		wantError bool
	}{
		{
			name:      "Found task",
			id:        testTask.ID,
			wantError: false,
		},
		{
			name:      "Task not found",
			id:        uuid.New().String(),
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			task, err := repo.FindByID(ctx, tc.id)

			if tc.wantError {
				assert.Error(t, err)
				assert.Nil(t, task)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, task)
				assert.Equal(t, testTask.ID, task.ID)
				assert.Equal(t, "Test Task", task.Title)
				assert.NotNil(t, task.Column)
				assert.Equal(t, column.ID, task.Column.ID)
			}
		})
	}
}

func TestTaskRepository_FindByColumnID(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &taskRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)
	column1 := createTestColumn(db, board.ID)
	column2 := createTestColumn(db, board.ID)

	task1 := &models.Task{
		ColumnID:    column1.ID,
		Title:       "Task 1",
		Description: "Description 1",
	}
	task2 := &models.Task{
		ColumnID:    column1.ID,
		Title:       "Task 2",
		Description: "Description 2",
	}
	task3 := &models.Task{
		ColumnID:    column2.ID,
		Title:       "Task 3",
		Description: "Description 3",
	}

	err := repo.Create(ctx, task1)
	require.NoError(t, err)
	err = repo.Create(ctx, task2)
	require.NoError(t, err)
	err = repo.Create(ctx, task3)
	require.NoError(t, err)

	testCases := []struct {
		name           string
		columnID       string
		expectedCount  int
		expectedTitles []string
		wantError      bool
	}{
		{
			name:           "Found column1 tasks",
			columnID:       column1.ID,
			expectedCount:  2,
			expectedTitles: []string{"Task 1", "Task 2"},
			wantError:      false,
		},
		{
			name:           "Found column2 tasks",
			columnID:       column2.ID,
			expectedCount:  1,
			expectedTitles: []string{"Task 3"},
			wantError:      false,
		},
		{
			name:           "Column has no tasks",
			columnID:       uuid.New().String(),
			expectedCount:  0,
			expectedTitles: []string{},
			wantError:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tasks, err := repo.FindByColumnID(ctx, tc.columnID)

			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, tasks, tc.expectedCount)

				titles := make([]string, len(tasks))
				for i, task := range tasks {
					titles[i] = task.Title
					assert.NotNil(t, task.Column)
					assert.Equal(t, tc.columnID, task.ColumnID)
				}
				assert.ElementsMatch(t, tc.expectedTitles, titles)
			}
		})
	}
}

func TestTaskRepository_Update(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &taskRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)
	column := createTestColumn(db, board.ID)

	testTask := &models.Task{
		ColumnID:    column.ID,
		Title:       "Test Task",
		Description: "Original Description",
	}
	err := repo.Create(ctx, testTask)
	require.NoError(t, err)

	newDeadline := time.Now().Add(48 * time.Hour)
	testTask.Title = "Updated Task"
	testTask.Description = "Updated Description"
	testTask.Deadline = &newDeadline
	err = repo.Update(ctx, testTask)
	assert.NoError(t, err)

	updatedTask, err := repo.FindByID(ctx, testTask.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Task", updatedTask.Title)
	assert.Equal(t, "Updated Description", updatedTask.Description)
	assert.NotNil(t, updatedTask.Deadline)
}

func TestTaskRepository_Delete(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &taskRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)
	column := createTestColumn(db, board.ID)

	testTask := &models.Task{
		ColumnID:    column.ID,
		Title:       "Test Task",
		Description: "Test Description",
	}
	err := repo.Create(ctx, testTask)
	require.NoError(t, err)

	err = repo.Delete(ctx, testTask.ID)
	assert.NoError(t, err)

	_, err = repo.FindByID(ctx, testTask.ID)
	assert.Error(t, err)

	testCases := []struct {
		name      string
		id        string
		wantError bool
	}{
		{
			name:      "Delete non-existent task",
			id:        uuid.New().String(),
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.Delete(ctx, tc.id)
			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTaskRepository_SoftDelete(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &taskRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)
	column := createTestColumn(db, board.ID)

	testTask := &models.Task{
		ColumnID:    column.ID,
		Title:       "Test Task",
		Description: "Test Description",
	}
	err := repo.Create(ctx, testTask)
	require.NoError(t, err)

	err = repo.SoftDelete(ctx, testTask.ID)
	assert.NoError(t, err)

	_, err = repo.FindByID(ctx, testTask.ID)
	assert.Error(t, err)

	var deletedTask models.Task
	err = db.Unscoped().Where("id = ?", testTask.ID).First(&deletedTask).Error
	assert.NoError(t, err)
	assert.NotNil(t, deletedTask.DeletedAt)

	testCases := []struct {
		name      string
		id        string
		wantError bool
	}{
		{
			name:      "Soft delete non-existent task",
			id:        uuid.New().String(),
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.SoftDelete(ctx, tc.id)
			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTaskRepository_Preloading(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &taskRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)
	column := createTestColumn(db, board.ID)

	task := &models.Task{
		ColumnID:    column.ID,
		Title:       "Test Task",
		Description: "Test Description",
	}
	err := repo.Create(ctx, task)
	require.NoError(t, err)

	taskWithRelations, err := repo.FindByID(ctx, task.ID)
	assert.NoError(t, err)

	assert.NotNil(t, taskWithRelations.Column)
	assert.Equal(t, column.ID, taskWithRelations.Column.ID)
	assert.NotNil(t, taskWithRelations.Comments)
	assert.NotNil(t, taskWithRelations.Labels)
	assert.NotNil(t, taskWithRelations.Attachments)

	tasks, err := repo.FindByColumnID(ctx, column.ID)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(tasks), 1)

	for _, task := range tasks {
		assert.NotNil(t, task.Column)
		assert.NotNil(t, task.Comments)
		assert.NotNil(t, task.Labels)
		assert.NotNil(t, task.Attachments)
	}
}

func TestTaskRepository_Context(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &taskRepository{db: db}

	ctx, cancel := context.WithCancel(context.Background())

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)
	column := createTestColumn(db, board.ID)

	testTask := &models.Task{
		ColumnID:    column.ID,
		Title:       "Test Task",
		Description: "Test Description",
	}

	err := repo.Create(ctx, testTask)
	assert.NoError(t, err)

	task, err := repo.FindByID(ctx, testTask.ID)
	assert.NoError(t, err)
	assert.NotNil(t, task)

	cancel()

	_, err = repo.FindByID(ctx, testTask.ID)
	assert.Error(t, err)
}
