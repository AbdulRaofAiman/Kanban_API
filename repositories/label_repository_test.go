package repositories

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"kanban-backend/models"
)

func TestLabelRepository_Create(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &labelRepository{db: db}
	ctx := context.Background()

	testCases := []struct {
		name      string
		label     *models.Label
		wantError bool
	}{
		{
			name: "Valid label",
			label: &models.Label{
				Name:  "Bug",
				Color: "#FF0000",
			},
			wantError: false,
		},
		{
			name: "Valid label without color",
			label: &models.Label{
				Name: "Feature",
			},
			wantError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.Create(ctx, tc.label)

			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tc.label.ID)
				assert.Equal(t, tc.label.Name, tc.label.Name)
			}
		})
	}
}

func TestLabelRepository_FindByID(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &labelRepository{db: db}
	ctx := context.Background()

	testLabel := &models.Label{
		Name:  "Bug",
		Color: "#FF0000",
	}
	err := repo.Create(ctx, testLabel)
	require.NoError(t, err)

	testCases := []struct {
		name      string
		id        string
		wantError bool
	}{
		{
			name:      "Found label",
			id:        testLabel.ID,
			wantError: false,
		},
		{
			name:      "Label not found",
			id:        uuid.New().String(),
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			label, err := repo.FindByID(ctx, tc.id)

			if tc.wantError {
				assert.Error(t, err)
				assert.Nil(t, label)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, label)
				assert.Equal(t, testLabel.ID, label.ID)
				assert.Equal(t, "Bug", label.Name)
				assert.Equal(t, "#FF0000", label.Color)
			}
		})
	}
}

func TestLabelRepository_FindAll(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &labelRepository{db: db}
	ctx := context.Background()

	label1 := &models.Label{
		Name:  "Bug",
		Color: "#FF0000",
	}
	label2 := &models.Label{
		Name:  "Feature",
		Color: "#00FF00",
	}
	label3 := &models.Label{
		Name:  "High Priority",
		Color: "#FFFF00",
	}

	err := repo.Create(ctx, label1)
	require.NoError(t, err)
	err = repo.Create(ctx, label2)
	require.NoError(t, err)
	err = repo.Create(ctx, label3)
	require.NoError(t, err)

	labels, err := repo.FindAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, labels, 3)

	names := make([]string, len(labels))
	for i, label := range labels {
		names[i] = label.Name
	}
	assert.ElementsMatch(t, []string{"Bug", "Feature", "High Priority"}, names)
}

func TestLabelRepository_Update(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &labelRepository{db: db}
	ctx := context.Background()

	testLabel := &models.Label{
		Name:  "Bug",
		Color: "#FF0000",
	}
	err := repo.Create(ctx, testLabel)
	require.NoError(t, err)

	testLabel.Name = "Critical Bug"
	testLabel.Color = "#FF5500"
	err = repo.Update(ctx, testLabel)
	assert.NoError(t, err)

	updatedLabel, err := repo.FindByID(ctx, testLabel.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Critical Bug", updatedLabel.Name)
	assert.Equal(t, "#FF5500", updatedLabel.Color)
}

func TestLabelRepository_Delete(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &labelRepository{db: db}
	ctx := context.Background()

	testLabel := &models.Label{
		Name:  "Bug",
		Color: "#FF0000",
	}
	err := repo.Create(ctx, testLabel)
	require.NoError(t, err)

	err = repo.Delete(ctx, testLabel.ID)
	assert.NoError(t, err)

	_, err = repo.FindByID(ctx, testLabel.ID)
	assert.Error(t, err)

	testCases := []struct {
		name      string
		id        string
		wantError bool
	}{
		{
			name:      "Delete non-existent label",
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

func TestLabelRepository_SoftDelete(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &labelRepository{db: db}
	ctx := context.Background()

	testLabel := &models.Label{
		Name:  "Bug",
		Color: "#FF0000",
	}
	err := repo.Create(ctx, testLabel)
	require.NoError(t, err)

	err = repo.SoftDelete(ctx, testLabel.ID)
	assert.NoError(t, err)

	_, err = repo.FindByID(ctx, testLabel.ID)
	assert.Error(t, err)

	var deletedLabel models.Label
	err = db.Unscoped().Where("id = ?", testLabel.ID).First(&deletedLabel).Error
	assert.NoError(t, err)
	assert.NotNil(t, deletedLabel.DeletedAt)

	testCases := []struct {
		name      string
		id        string
		wantError bool
	}{
		{
			name:      "Soft delete non-existent label",
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

func TestLabelRepository_TaskAssociations(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &labelRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)
	column := createTestColumn(db, board.ID)
	task1 := &models.Task{ColumnID: column.ID, Title: "Task 1"}
	task2 := &models.Task{ColumnID: column.ID, Title: "Task 2"}
	db.Create(task1)
	db.Create(task2)

	label1 := &models.Label{
		Name: "Bug",
	}
	label2 := &models.Label{
		Name: "Feature",
	}

	err := repo.Create(ctx, label1)
	require.NoError(t, err)
	err = repo.Create(ctx, label2)
	require.NoError(t, err)

	err = db.Model(label1).Association("Tasks").Append([]*models.Task{task1, task2})
	require.NoError(t, err)
	err = db.Model(label2).Association("Tasks").Append([]*models.Task{task1})
	require.NoError(t, err)

	labelWithTasks1, err := repo.FindByID(ctx, label1.ID)
	assert.NoError(t, err)
	assert.Len(t, labelWithTasks1.Tasks, 2)

	labelWithTasks2, err := repo.FindByID(ctx, label2.ID)
	assert.NoError(t, err)
	assert.Len(t, labelWithTasks2.Tasks, 1)

	labels, err := repo.FindAll(ctx)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(labels), 2)

	taskCounts := make(map[string]int)
	for _, label := range labels {
		taskCounts[label.ID] = len(label.Tasks)
	}
	assert.Equal(t, 2, taskCounts[label1.ID])
	assert.Equal(t, 1, taskCounts[label2.ID])
}

func TestLabelRepository_Preloading(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &labelRepository{db: db}
	ctx := context.Background()

	label := &models.Label{
		Name: "Bug",
	}
	err := repo.Create(ctx, label)
	require.NoError(t, err)

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)
	column := createTestColumn(db, board.ID)
	task := &models.Task{ColumnID: column.ID, Title: "Test Task"}
	db.Create(task)

	err = db.Model(label).Association("Tasks").Append([]*models.Task{task})
	require.NoError(t, err)

	labelWithTasks, err := repo.FindByID(ctx, label.ID)
	assert.NoError(t, err)

	assert.NotNil(t, labelWithTasks.Tasks)
	assert.Len(t, labelWithTasks.Tasks, 1)
	assert.Equal(t, task.ID, labelWithTasks.Tasks[0].ID)
}

func TestLabelRepository_Context(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &labelRepository{db: db}

	ctx, cancel := context.WithCancel(context.Background())

	testLabel := &models.Label{
		Name: "Bug",
	}

	err := repo.Create(ctx, testLabel)
	assert.NoError(t, err)

	label, err := repo.FindByID(ctx, testLabel.ID)
	assert.NoError(t, err)
	assert.NotNil(t, label)

	cancel()

	_, err = repo.FindByID(ctx, testLabel.ID)
	assert.Error(t, err)
}
