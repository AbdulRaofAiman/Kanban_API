package repositories

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"kanban-backend/models"
)

func TestColumnRepository_Create(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &columnRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)

	testCases := []struct {
		name      string
		column    *models.Column
		wantError bool
	}{
		{
			name: "Valid column",
			column: &models.Column{
				BoardID: board.ID,
				Title:   "To Do",
				Order:   1,
			},
			wantError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.Create(ctx, tc.column)

			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tc.column.ID)
				assert.Equal(t, "To Do", tc.column.Title)
				assert.Equal(t, board.ID, tc.column.BoardID)
				assert.Equal(t, 1, tc.column.Order)
			}
		})
	}
}

func TestColumnRepository_FindByID(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &columnRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)

	testColumn := &models.Column{
		BoardID: board.ID,
		Title:   "To Do",
		Order:   1,
	}
	err := repo.Create(ctx, testColumn)
	require.NoError(t, err)

	testCases := []struct {
		name      string
		id        string
		wantError bool
	}{
		{
			name:      "Found column",
			id:        testColumn.ID,
			wantError: false,
		},
		{
			name:      "Column not found",
			id:        uuid.New().String(),
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			column, err := repo.FindByID(ctx, tc.id)

			if tc.wantError {
				assert.Error(t, err)
				assert.Nil(t, column)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, column)
				assert.Equal(t, testColumn.ID, column.ID)
				assert.Equal(t, "To Do", column.Title)
				assert.NotNil(t, column.Board)
				assert.Equal(t, board.ID, column.Board.ID)
			}
		})
	}
}

func TestColumnRepository_FindByBoardID(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &columnRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")
	board1 := createTestBoard(db, user.ID)
	board2 := createTestBoard(db, user.ID)

	column1 := &models.Column{
		BoardID: board1.ID,
		Title:   "To Do",
		Order:   1,
	}
	column2 := &models.Column{
		BoardID: board1.ID,
		Title:   "In Progress",
		Order:   2,
	}
	column3 := &models.Column{
		BoardID: board2.ID,
		Title:   "Done",
		Order:   3,
	}

	err := repo.Create(ctx, column1)
	require.NoError(t, err)
	err = repo.Create(ctx, column2)
	require.NoError(t, err)
	err = repo.Create(ctx, column3)
	require.NoError(t, err)

	testCases := []struct {
		name           string
		boardID        string
		expectedCount  int
		expectedTitles []string
		wantError      bool
	}{
		{
			name:           "Found board1 columns",
			boardID:        board1.ID,
			expectedCount:  2,
			expectedTitles: []string{"To Do", "In Progress"},
			wantError:      false,
		},
		{
			name:           "Found board2 columns",
			boardID:        board2.ID,
			expectedCount:  1,
			expectedTitles: []string{"Done"},
			wantError:      false,
		},
		{
			name:           "Board has no columns",
			boardID:        uuid.New().String(),
			expectedCount:  0,
			expectedTitles: []string{},
			wantError:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			columns, err := repo.FindByBoardID(ctx, tc.boardID)

			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, columns, tc.expectedCount)

				titles := make([]string, len(columns))
				for i, column := range columns {
					titles[i] = column.Title
					assert.NotNil(t, column.Board)
					assert.Equal(t, tc.boardID, column.BoardID)
				}
				assert.ElementsMatch(t, tc.expectedTitles, titles)
			}
		})
	}
}

func TestColumnRepository_Update(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &columnRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)

	testColumn := &models.Column{
		BoardID: board.ID,
		Title:   "To Do",
		Order:   1,
	}
	err := repo.Create(ctx, testColumn)
	require.NoError(t, err)

	testColumn.Title = "In Progress"
	testColumn.Order = 2
	err = repo.Update(ctx, testColumn)
	assert.NoError(t, err)

	updatedColumn, err := repo.FindByID(ctx, testColumn.ID)
	assert.NoError(t, err)
	assert.Equal(t, "In Progress", updatedColumn.Title)
	assert.Equal(t, 2, updatedColumn.Order)
}

func TestColumnRepository_Delete(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &columnRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)

	testColumn := &models.Column{
		BoardID: board.ID,
		Title:   "To Do",
		Order:   1,
	}
	err := repo.Create(ctx, testColumn)
	require.NoError(t, err)

	err = repo.Delete(ctx, testColumn.ID)
	assert.NoError(t, err)

	_, err = repo.FindByID(ctx, testColumn.ID)
	assert.Error(t, err)

	testCases := []struct {
		name      string
		id        string
		wantError bool
	}{
		{
			name:      "Delete non-existent column",
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

func TestColumnRepository_SoftDelete(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &columnRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)

	testColumn := &models.Column{
		BoardID: board.ID,
		Title:   "To Do",
		Order:   1,
	}
	err := repo.Create(ctx, testColumn)
	require.NoError(t, err)

	err = repo.SoftDelete(ctx, testColumn.ID)
	assert.NoError(t, err)

	_, err = repo.FindByID(ctx, testColumn.ID)
	assert.Error(t, err)

	var deletedColumn models.Column
	err = db.Unscoped().Where("id = ?", testColumn.ID).First(&deletedColumn).Error
	assert.NoError(t, err)
	assert.NotNil(t, deletedColumn.DeletedAt)

	testCases := []struct {
		name      string
		id        string
		wantError bool
	}{
		{
			name:      "Soft delete non-existent column",
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

func TestColumnRepository_Preloading(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &columnRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)

	column := &models.Column{
		BoardID: board.ID,
		Title:   "To Do",
		Order:   1,
	}
	err := repo.Create(ctx, column)
	require.NoError(t, err)

	columnWithRelations, err := repo.FindByID(ctx, column.ID)
	assert.NoError(t, err)

	assert.NotNil(t, columnWithRelations.Board)
	assert.Equal(t, board.ID, columnWithRelations.Board.ID)
	assert.NotNil(t, columnWithRelations.Tasks)

	columns, err := repo.FindByBoardID(ctx, board.ID)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(columns), 1)

	for _, c := range columns {
		assert.NotNil(t, c.Board)
		assert.NotNil(t, c.Tasks)
	}
}

func TestColumnRepository_Context(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &columnRepository{db: db}

	ctx, cancel := context.WithCancel(context.Background())

	user := createTestUser(db, "testuser", "test@example.com")
	board := createTestBoard(db, user.ID)

	testColumn := &models.Column{
		BoardID: board.ID,
		Title:   "To Do",
		Order:   1,
	}

	err := repo.Create(ctx, testColumn)
	assert.NoError(t, err)

	column, err := repo.FindByID(ctx, testColumn.ID)
	assert.NoError(t, err)
	assert.NotNil(t, column)

	cancel()

	_, err = repo.FindByID(ctx, testColumn.ID)
	assert.Error(t, err)
}
