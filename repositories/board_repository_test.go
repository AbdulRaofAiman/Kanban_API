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

func TestBoardRepository_Create(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &boardRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")

	testCases := []struct {
		name      string
		board     *models.Board
		wantError bool
	}{
		{
			name: "Valid board",
			board: &models.Board{
				Title:  "Test Board",
				UserID: user.ID,
				Color:  "#FF5733",
			},
			wantError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.Create(ctx, tc.board)

			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tc.board.ID)
				assert.Equal(t, "Test Board", tc.board.Title)
				assert.Equal(t, user.ID, tc.board.UserID)
			}
		})
	}
}

func TestBoardRepository_FindByID(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &boardRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")

	testBoard := &models.Board{
		Title:  "Test Board",
		UserID: user.ID,
		Color:  "#FF5733",
	}
	err := repo.Create(ctx, testBoard)
	require.NoError(t, err)

	testCases := []struct {
		name      string
		id        string
		wantError bool
	}{
		{
			name:      "Found board",
			id:        testBoard.ID,
			wantError: false,
		},
		{
			name:      "Board not found",
			id:        uuid.New().String(),
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			board, err := repo.FindByID(ctx, tc.id)

			if tc.wantError {
				assert.Error(t, err)
				assert.Nil(t, board)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, board)
				assert.Equal(t, testBoard.ID, board.ID)
				assert.Equal(t, "Test Board", board.Title)
				assert.NotNil(t, board.User)
				assert.Equal(t, user.ID, board.User.ID)
			}
		})
	}
}

func TestBoardRepository_FindByUserID(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &boardRepository{db: db}
	ctx := context.Background()

	user1 := createTestUser(db, "user1", "user1@example.com")
	user2 := createTestUser(db, "user2", "user2@example.com")

	board1 := &models.Board{
		Title:  "User1 Board 1",
		UserID: user1.ID,
		Color:  "#FF5733",
	}
	board2 := &models.Board{
		Title:  "User1 Board 2",
		UserID: user1.ID,
		Color:  "#33FF57",
	}
	board3 := &models.Board{
		Title:  "User2 Board 1",
		UserID: user2.ID,
		Color:  "#3357FF",
	}

	err := repo.Create(ctx, board1)
	require.NoError(t, err)
	err = repo.Create(ctx, board2)
	require.NoError(t, err)
	err = repo.Create(ctx, board3)
	require.NoError(t, err)

	testCases := []struct {
		name           string
		userID         string
		expectedCount  int
		expectedTitles []string
		wantError      bool
	}{
		{
			name:           "Found user1 boards",
			userID:         user1.ID,
			expectedCount:  2,
			expectedTitles: []string{"User1 Board 1", "User1 Board 2"},
			wantError:      false,
		},
		{
			name:           "Found user2 boards",
			userID:         user2.ID,
			expectedCount:  1,
			expectedTitles: []string{"User2 Board 1"},
			wantError:      false,
		},
		{
			name:           "User has no boards",
			userID:         uuid.New().String(),
			expectedCount:  0,
			expectedTitles: []string{},
			wantError:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			boards, err := repo.FindByUserID(ctx, tc.userID)

			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, boards, tc.expectedCount)

				titles := make([]string, len(boards))
				for i, board := range boards {
					titles[i] = board.Title
					assert.NotNil(t, board.User)
					assert.Equal(t, tc.userID, board.UserID)
				}
				assert.ElementsMatch(t, tc.expectedTitles, titles)
			}
		})
	}
}

func TestBoardRepository_Update(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &boardRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")

	testBoard := &models.Board{
		Title:  "Test Board",
		UserID: user.ID,
		Color:  "#FF5733",
	}
	err := repo.Create(ctx, testBoard)
	require.NoError(t, err)

	time.Sleep(10 * time.Millisecond)
	testBoard.Title = "Updated Board"
	testBoard.Color = "#33FF57"
	err = repo.Update(ctx, testBoard)
	assert.NoError(t, err)

	updatedBoard, err := repo.FindByID(ctx, testBoard.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Board", updatedBoard.Title)
	assert.Equal(t, "#33FF57", updatedBoard.Color)
	assert.True(t, updatedBoard.UpdatedAt.After(testBoard.CreatedAt))
}

func TestBoardRepository_Delete(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &boardRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")

	testBoard := &models.Board{
		Title:  "Test Board",
		UserID: user.ID,
		Color:  "#FF5733",
	}
	err := repo.Create(ctx, testBoard)
	require.NoError(t, err)

	err = repo.Delete(ctx, testBoard.ID)
	assert.NoError(t, err)

	_, err = repo.FindByID(ctx, testBoard.ID)
	assert.Error(t, err)

	testCases := []struct {
		name      string
		id        string
		wantError bool
	}{
		{
			name:      "Delete non-existent board",
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

func TestBoardRepository_SoftDelete(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &boardRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")

	testBoard := &models.Board{
		Title:  "Test Board",
		UserID: user.ID,
		Color:  "#FF5733",
	}
	err := repo.Create(ctx, testBoard)
	require.NoError(t, err)

	err = repo.SoftDelete(ctx, testBoard.ID)
	assert.NoError(t, err)

	_, err = repo.FindByID(ctx, testBoard.ID)
	assert.Error(t, err)

	var deletedBoard models.Board
	err = db.Unscoped().Where("id = ?", testBoard.ID).First(&deletedBoard).Error
	assert.NoError(t, err)
	assert.NotNil(t, deletedBoard.DeletedAt)

	testCases := []struct {
		name      string
		id        string
		wantError bool
	}{
		{
			name:      "Soft delete non-existent board",
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

func TestBoardRepository_Preloading(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &boardRepository{db: db}
	ctx := context.Background()

	user := createTestUser(db, "testuser", "test@example.com")

	board := &models.Board{
		Title:  "Test Board",
		UserID: user.ID,
		Color:  "#FF5733",
	}
	err := repo.Create(ctx, board)
	require.NoError(t, err)

	boardWithRelations, err := repo.FindByID(ctx, board.ID)
	assert.NoError(t, err)

	assert.NotNil(t, boardWithRelations.User)
	assert.Equal(t, user.ID, boardWithRelations.User.ID)
	assert.NotNil(t, boardWithRelations.Columns)
	assert.NotNil(t, boardWithRelations.Members)

	boards, err := repo.FindByUserID(ctx, user.ID)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(boards), 1)

	for _, b := range boards {
		assert.NotNil(t, b.User)
		assert.NotNil(t, b.Columns)
		assert.NotNil(t, b.Members)
	}
}

func TestBoardRepository_Context(t *testing.T) {
	db := setupRepositoryTestDB(t)
	repo := &boardRepository{db: db}

	ctx, cancel := context.WithCancel(context.Background())

	user := createTestUser(db, "testuser", "test@example.com")

	testBoard := &models.Board{
		Title:  "Test Board",
		UserID: user.ID,
		Color:  "#FF5733",
	}

	err := repo.Create(ctx, testBoard)
	assert.NoError(t, err)

	board, err := repo.FindByID(ctx, testBoard.ID)
	assert.NoError(t, err)
	assert.NotNil(t, board)

	cancel()

	_, err = repo.FindByID(ctx, testBoard.ID)
	assert.Error(t, err)
}
