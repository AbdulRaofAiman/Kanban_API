package repositories

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"kanban-backend/models"
)

func TestAttachmentRepository_Create(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &attachmentRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)
	column := createTestColumn(db, board.ID)
	task := &models.Task{ColumnID: column.ID, Title: "Test Task"}
	db.Create(task)

	testCases := []struct {
		name       string
		attachment *models.Attachment
		wantError  bool
	}{
		{
			name: "Valid attachment",
			attachment: &models.Attachment{
				TaskID:   task.ID,
				FileName: "test.jpg",
				FileURL:  "https://example.com/test.jpg",
				FileSize: 1024,
			},
			wantError: false,
		},
		{
			name: "Valid attachment without file size",
			attachment: &models.Attachment{
				TaskID:   task.ID,
				FileName: "document.pdf",
				FileURL:  "https://example.com/document.pdf",
			},
			wantError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.Create(ctx, tc.attachment)

			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tc.attachment.ID)
				assert.Equal(t, tc.attachment.FileName, tc.attachment.FileName)
				assert.Equal(t, task.ID, tc.attachment.TaskID)
			}
		})
	}
}

func TestAttachmentRepository_FindByID(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &attachmentRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)
	column := createTestColumn(db, board.ID)
	task := &models.Task{ColumnID: column.ID, Title: "Test Task"}
	db.Create(task)

	testAttachment := &models.Attachment{
		TaskID:   task.ID,
		FileName: "test.jpg",
		FileURL:  "https://example.com/test.jpg",
		FileSize: 1024,
	}
	err := repo.Create(ctx, testAttachment)
	require.NoError(t, err)

	testCases := []struct {
		name      string
		id        string
		wantError bool
	}{
		{
			name:      "Found attachment",
			id:        testAttachment.ID,
			wantError: false,
		},
		{
			name:      "Attachment not found",
			id:        uuid.New().String(),
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			attachment, err := repo.FindByID(ctx, tc.id)

			if tc.wantError {
				assert.Error(t, err)
				assert.Nil(t, attachment)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, attachment)
				assert.Equal(t, testAttachment.ID, attachment.ID)
				assert.Equal(t, "test.jpg", attachment.FileName)
				assert.Equal(t, "https://example.com/test.jpg", attachment.FileURL)
				assert.Equal(t, int64(1024), attachment.FileSize)
				assert.NotNil(t, attachment.Task)
				assert.Equal(t, task.ID, attachment.Task.ID)
			}
		})
	}
}

func TestAttachmentRepository_FindByTaskID(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &attachmentRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)
	column := createTestColumn(db, board.ID)
	task1 := &models.Task{ColumnID: column.ID, Title: "Task 1"}
	task2 := &models.Task{ColumnID: column.ID, Title: "Task 2"}
	db.Create(task1)
	db.Create(task2)

	attachment1 := &models.Attachment{
		TaskID:   task1.ID,
		FileName: "image1.jpg",
		FileURL:  "https://example.com/image1.jpg",
		FileSize: 1024,
	}
	attachment2 := &models.Attachment{
		TaskID:   task1.ID,
		FileName: "image2.jpg",
		FileURL:  "https://example.com/image2.jpg",
		FileSize: 2048,
	}
	attachment3 := &models.Attachment{
		TaskID:   task2.ID,
		FileName: "image3.jpg",
		FileURL:  "https://example.com/image3.jpg",
		FileSize: 3072,
	}

	err := repo.Create(ctx, attachment1)
	require.NoError(t, err)
	err = repo.Create(ctx, attachment2)
	require.NoError(t, err)
	err = repo.Create(ctx, attachment3)
	require.NoError(t, err)

	testCases := []struct {
		name          string
		taskID        string
		expectedCount int
		expectedNames []string
		wantError     bool
	}{
		{
			name:          "Found task1 attachments",
			taskID:        task1.ID,
			expectedCount: 2,
			expectedNames: []string{"image1.jpg", "image2.jpg"},
			wantError:     false,
		},
		{
			name:          "Found task2 attachments",
			taskID:        task2.ID,
			expectedCount: 1,
			expectedNames: []string{"image3.jpg"},
			wantError:     false,
		},
		{
			name:          "Task has no attachments",
			taskID:        uuid.New().String(),
			expectedCount: 0,
			expectedNames: []string{},
			wantError:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			attachments, err := repo.FindByTaskID(ctx, tc.taskID)

			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, attachments, tc.expectedCount)

				names := make([]string, len(attachments))
				for i, attachment := range attachments {
					names[i] = attachment.FileName
					assert.NotNil(t, attachment.Task)
					assert.Equal(t, tc.taskID, attachment.TaskID)
				}
				assert.ElementsMatch(t, tc.expectedNames, names)
			}
		})
	}
}

func TestAttachmentRepository_Update(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &attachmentRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)
	column := createTestColumn(db, board.ID)
	task := &models.Task{ColumnID: column.ID, Title: "Test Task"}
	db.Create(task)

	testAttachment := &models.Attachment{
		TaskID:   task.ID,
		FileName: "old.jpg",
		FileURL:  "https://example.com/old.jpg",
		FileSize: 1024,
	}
	err := repo.Create(ctx, testAttachment)
	require.NoError(t, err)

	testAttachment.FileName = "new.jpg"
	testAttachment.FileURL = "https://example.com/new.jpg"
	testAttachment.FileSize = 2048
	err = repo.Update(ctx, testAttachment)
	assert.NoError(t, err)

	updatedAttachment, err := repo.FindByID(ctx, testAttachment.ID)
	assert.NoError(t, err)
	assert.Equal(t, "new.jpg", updatedAttachment.FileName)
	assert.Equal(t, "https://example.com/new.jpg", updatedAttachment.FileURL)
	assert.Equal(t, int64(2048), updatedAttachment.FileSize)
}

func TestAttachmentRepository_Delete(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &attachmentRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)
	column := createTestColumn(db, board.ID)
	task := &models.Task{ColumnID: column.ID, Title: "Test Task"}
	db.Create(task)

	testAttachment := &models.Attachment{
		TaskID:   task.ID,
		FileName: "test.jpg",
		FileURL:  "https://example.com/test.jpg",
		FileSize: 1024,
	}
	err := repo.Create(ctx, testAttachment)
	require.NoError(t, err)

	err = repo.Delete(ctx, testAttachment.ID)
	assert.NoError(t, err)

	_, err = repo.FindByID(ctx, testAttachment.ID)
	assert.Error(t, err)

	testCases := []struct {
		name      string
		id        string
		wantError bool
	}{
		{
			name:      "Delete non-existent attachment",
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

func TestAttachmentRepository_SoftDelete(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &attachmentRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)
	column := createTestColumn(db, board.ID)
	task := &models.Task{ColumnID: column.ID, Title: "Test Task"}
	db.Create(task)

	testAttachment := &models.Attachment{
		TaskID:   task.ID,
		FileName: "test.jpg",
		FileURL:  "https://example.com/test.jpg",
		FileSize: 1024,
	}
	err := repo.Create(ctx, testAttachment)
	require.NoError(t, err)

	err = repo.SoftDelete(ctx, testAttachment.ID)
	assert.NoError(t, err)

	_, err = repo.FindByID(ctx, testAttachment.ID)
	assert.Error(t, err)

	var deletedAttachment models.Attachment
	err = db.Unscoped().Where("id = ?", testAttachment.ID).First(&deletedAttachment).Error
	assert.NoError(t, err)
	assert.NotNil(t, deletedAttachment.DeletedAt)

	testCases := []struct {
		name      string
		id        string
		wantError bool
	}{
		{
			name:      "Soft delete non-existent attachment",
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

func TestAttachmentRepository_Preloading(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &attachmentRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)
	column := createTestColumn(db, board.ID)
	task := &models.Task{ColumnID: column.ID, Title: "Test Task"}
	db.Create(task)

	attachment := &models.Attachment{
		TaskID:   task.ID,
		FileName: "test.jpg",
		FileURL:  "https://example.com/test.jpg",
		FileSize: 1024,
	}
	err := repo.Create(ctx, attachment)
	require.NoError(t, err)

	attachmentWithTask, err := repo.FindByID(ctx, attachment.ID)
	assert.NoError(t, err)

	assert.NotNil(t, attachmentWithTask.Task)
	assert.Equal(t, task.ID, attachmentWithTask.Task.ID)

	attachments, err := repo.FindByTaskID(ctx, task.ID)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(attachments), 1)

	for _, attachment := range attachments {
		assert.NotNil(t, attachment.Task)
	}
}

func TestAttachmentRepository_Context(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &attachmentRepository{db: db}

	ctx, cancel := context.WithCancel(context.Background())

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)
	column := createTestColumn(db, board.ID)
	task := &models.Task{ColumnID: column.ID, Title: "Test Task"}
	db.Create(task)

	testAttachment := &models.Attachment{
		TaskID:   task.ID,
		FileName: "test.jpg",
		FileURL:  "https://example.com/test.jpg",
		FileSize: 1024,
	}

	err := repo.Create(ctx, testAttachment)
	assert.NoError(t, err)

	attachment, err := repo.FindByID(ctx, testAttachment.ID)
	assert.NoError(t, err)
	assert.NotNil(t, attachment)

	cancel()

	_, err = repo.FindByID(ctx, testAttachment.ID)
	assert.Error(t, err)
}
