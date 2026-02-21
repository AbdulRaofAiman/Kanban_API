package repositories

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"kanban-backend/models"
)

func TestCommentRepository_Create(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &commentRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)
	column := createTestColumn(db, board.ID)
	task := &models.Task{ColumnID: column.ID, Title: "Test Task"}
	db.Create(task)

	testCases := []struct {
		name      string
		comment   *models.Comment
		wantError bool
	}{
		{
			name: "Valid comment",
			comment: &models.Comment{
				TaskID:  task.ID,
				UserID:  user.ID,
				Content: "This is a test comment",
			},
			wantError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.Create(ctx, tc.comment)

			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tc.comment.ID)
				assert.Equal(t, "This is a test comment", tc.comment.Content)
				assert.Equal(t, task.ID, tc.comment.TaskID)
				assert.Equal(t, user.ID, tc.comment.UserID)
			}
		})
	}
}

func TestCommentRepository_FindByID(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &commentRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)
	column := createTestColumn(db, board.ID)
	task := &models.Task{ColumnID: column.ID, Title: "Test Task"}
	db.Create(task)

	testComment := &models.Comment{
		TaskID:  task.ID,
		UserID:  user.ID,
		Content: "Test comment content",
	}
	err := repo.Create(ctx, testComment)
	require.NoError(t, err)

	testCases := []struct {
		name      string
		id        string
		wantError bool
	}{
		{
			name:      "Found comment",
			id:        testComment.ID,
			wantError: false,
		},
		{
			name:      "Comment not found",
			id:        uuid.New().String(),
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			comment, err := repo.FindByID(ctx, tc.id)

			if tc.wantError {
				assert.Error(t, err)
				assert.Nil(t, comment)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, comment)
				assert.Equal(t, testComment.ID, comment.ID)
				assert.Equal(t, "Test comment content", comment.Content)
				assert.NotNil(t, comment.User)
				assert.NotNil(t, comment.Task)
				assert.Equal(t, user.ID, comment.User.ID)
				assert.Equal(t, task.ID, comment.Task.ID)
			}
		})
	}
}

func TestCommentRepository_FindByTaskID(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &commentRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)
	column := createTestColumn(db, board.ID)
	task1 := &models.Task{ColumnID: column.ID, Title: "Task 1"}
	task2 := &models.Task{ColumnID: column.ID, Title: "Task 2"}
	db.Create(task1)
	db.Create(task2)

	comment1 := &models.Comment{
		TaskID:  task1.ID,
		UserID:  user.ID,
		Content: "Comment 1",
	}
	comment2 := &models.Comment{
		TaskID:  task1.ID,
		UserID:  user.ID,
		Content: "Comment 2",
	}
	comment3 := &models.Comment{
		TaskID:  task2.ID,
		UserID:  user.ID,
		Content: "Comment 3",
	}

	err := repo.Create(ctx, comment1)
	require.NoError(t, err)
	err = repo.Create(ctx, comment2)
	require.NoError(t, err)
	err = repo.Create(ctx, comment3)
	require.NoError(t, err)

	testCases := []struct {
		name           string
		taskID         string
		expectedCount  int
		expectedTitles []string
		wantError      bool
	}{
		{
			name:           "Found task1 comments",
			taskID:         task1.ID,
			expectedCount:  2,
			expectedTitles: []string{"Comment 1", "Comment 2"},
			wantError:      false,
		},
		{
			name:           "Found task2 comments",
			taskID:         task2.ID,
			expectedCount:  1,
			expectedTitles: []string{"Comment 3"},
			wantError:      false,
		},
		{
			name:           "Task has no comments",
			taskID:         uuid.New().String(),
			expectedCount:  0,
			expectedTitles: []string{},
			wantError:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			comments, err := repo.FindByTaskID(ctx, tc.taskID)

			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, comments, tc.expectedCount)

				titles := make([]string, len(comments))
				for i, comment := range comments {
					titles[i] = comment.Content
					assert.NotNil(t, comment.User)
					assert.NotNil(t, comment.Task)
					assert.Equal(t, tc.taskID, comment.TaskID)
				}
				assert.ElementsMatch(t, tc.expectedTitles, titles)
			}
		})
	}
}

func TestCommentRepository_Update(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &commentRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)
	column := createTestColumn(db, board.ID)
	task := &models.Task{ColumnID: column.ID, Title: "Test Task"}
	db.Create(task)

	testComment := &models.Comment{
		TaskID:  task.ID,
		UserID:  user.ID,
		Content: "Original comment",
	}
	err := repo.Create(ctx, testComment)
	require.NoError(t, err)

	testComment.Content = "Updated comment"
	err = repo.Update(ctx, testComment)
	assert.NoError(t, err)

	updatedComment, err := repo.FindByID(ctx, testComment.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated comment", updatedComment.Content)
}

func TestCommentRepository_Delete(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &commentRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)
	column := createTestColumn(db, board.ID)
	task := &models.Task{ColumnID: column.ID, Title: "Test Task"}
	db.Create(task)

	testComment := &models.Comment{
		TaskID:  task.ID,
		UserID:  user.ID,
		Content: "Test comment",
	}
	err := repo.Create(ctx, testComment)
	require.NoError(t, err)

	err = repo.Delete(ctx, testComment.ID)
	assert.NoError(t, err)

	_, err = repo.FindByID(ctx, testComment.ID)
	assert.Error(t, err)

	testCases := []struct {
		name      string
		id        string
		wantError bool
	}{
		{
			name:      "Delete non-existent comment",
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

func TestCommentRepository_SoftDelete(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &commentRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)
	column := createTestColumn(db, board.ID)
	task := &models.Task{ColumnID: column.ID, Title: "Test Task"}
	db.Create(task)

	testComment := &models.Comment{
		TaskID:  task.ID,
		UserID:  user.ID,
		Content: "Test comment",
	}
	err := repo.Create(ctx, testComment)
	require.NoError(t, err)

	err = repo.SoftDelete(ctx, testComment.ID)
	assert.NoError(t, err)

	_, err = repo.FindByID(ctx, testComment.ID)
	assert.Error(t, err)

	var deletedComment models.Comment
	err = db.Unscoped().Where("id = ?", testComment.ID).First(&deletedComment).Error
	assert.NoError(t, err)
	assert.NotNil(t, deletedComment.DeletedAt)

	testCases := []struct {
		name      string
		id        string
		wantError bool
	}{
		{
			name:      "Soft delete non-existent comment",
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

func TestCommentRepository_Preloading(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &commentRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)
	column := createTestColumn(db, board.ID)
	task := &models.Task{ColumnID: column.ID, Title: "Test Task"}
	db.Create(task)

	comment := &models.Comment{
		TaskID:  task.ID,
		UserID:  user.ID,
		Content: "Test comment",
	}
	err := repo.Create(ctx, comment)
	require.NoError(t, err)

	commentWithRelations, err := repo.FindByID(ctx, comment.ID)
	assert.NoError(t, err)

	assert.NotNil(t, commentWithRelations.User)
	assert.Equal(t, user.ID, commentWithRelations.User.ID)
	assert.NotNil(t, commentWithRelations.Task)
	assert.Equal(t, task.ID, commentWithRelations.Task.ID)

	comments, err := repo.FindByTaskID(ctx, task.ID)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(comments), 1)

	for _, comment := range comments {
		assert.NotNil(t, comment.User)
		assert.NotNil(t, comment.Task)
	}
}

func TestCommentRepository_Context(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &commentRepository{db: db}

	ctx, cancel := context.WithCancel(context.Background())

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)
	column := createTestColumn(db, board.ID)
	task := &models.Task{ColumnID: column.ID, Title: "Test Task"}
	db.Create(task)

	testComment := &models.Comment{
		TaskID:  task.ID,
		UserID:  user.ID,
		Content: "Test comment",
	}

	err := repo.Create(ctx, testComment)
	assert.NoError(t, err)

	comment, err := repo.FindByID(ctx, testComment.ID)
	assert.NoError(t, err)
	assert.NotNil(t, comment)

	cancel()

	_, err = repo.FindByID(ctx, testComment.ID)
	assert.Error(t, err)
}
