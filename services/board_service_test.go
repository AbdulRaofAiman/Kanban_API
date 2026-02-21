package services

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	"kanban-backend/models"
	"kanban-backend/utils"
)

var boardTestIDCounter atomic.Int64

func generateBoardTestID() string {
	counter := boardTestIDCounter.Add(1)
	return "board-" + string(rune(counter))
}

var columnTestIDCounter atomic.Int64

func generateColumnTestID() string {
	counter := columnTestIDCounter.Add(1)
	return "col-" + string(rune(counter))
}

type mockBoardRepository struct {
	boards map[string]*models.Board
}

func newMockBoardRepository() *mockBoardRepository {
	return &mockBoardRepository{
		boards: make(map[string]*models.Board),
	}
}

func (m *mockBoardRepository) Create(ctx context.Context, board *models.Board) error {
	if board.ID == "" {
		board.ID = generateBoardTestID()
	}
	m.boards[board.ID] = board
	return nil
}

func (m *mockBoardRepository) FindByID(ctx context.Context, id string) (*models.Board, error) {
	board, exists := m.boards[id]
	if !exists {
		return nil, errors.New("board not found")
	}
	return board, nil
}

func (m *mockBoardRepository) FindByUserID(ctx context.Context, userID string) ([]*models.Board, error) {
	var boards []*models.Board
	for _, board := range m.boards {
		if board.UserID == userID {
			boards = append(boards, board)
		}
	}
	return boards, nil
}

func (m *mockBoardRepository) Update(ctx context.Context, board *models.Board) error {
	if _, exists := m.boards[board.ID]; !exists {
		return errors.New("board not found")
	}
	m.boards[board.ID] = board
	return nil
}

func (m *mockBoardRepository) Delete(ctx context.Context, id string) error {
	if _, exists := m.boards[id]; !exists {
		return errors.New("board not found")
	}
	delete(m.boards, id)
	return nil
}

func (m *mockBoardRepository) SoftDelete(ctx context.Context, id string) error {
	return m.Delete(ctx, id)
}

type mockColumnRepository struct {
	columns map[string]*models.Column
}

func newMockColumnRepository() *mockColumnRepository {
	return &mockColumnRepository{
		columns: make(map[string]*models.Column),
	}
}

func (m *mockColumnRepository) Create(ctx context.Context, column *models.Column) error {
	if column.ID == "" {
		column.ID = generateColumnTestID()
	}
	m.columns[column.ID] = column
	return nil
}

func (m *mockColumnRepository) FindByID(ctx context.Context, id string) (*models.Column, error) {
	column, exists := m.columns[id]
	if !exists {
		return nil, errors.New("column not found")
	}
	columnCopy := *column
	if column.Board != nil {
		boardCopy := *column.Board
		columnCopy.Board = &boardCopy
	}
	return &columnCopy, nil
}

func (m *mockColumnRepository) FindByBoardID(ctx context.Context, boardID string) ([]*models.Column, error) {
	var columns []*models.Column
	for _, column := range m.columns {
		if column.BoardID == boardID {
			columns = append(columns, column)
		}
	}
	return columns, nil
}

func (m *mockColumnRepository) Update(ctx context.Context, column *models.Column) error {
	if _, exists := m.columns[column.ID]; !exists {
		return errors.New("column not found")
	}
	m.columns[column.ID] = column
	return nil
}

func (m *mockColumnRepository) Delete(ctx context.Context, id string) error {
	if _, exists := m.columns[id]; !exists {
		return errors.New("column not found")
	}
	delete(m.columns, id)
	return nil
}

func (m *mockColumnRepository) SoftDelete(ctx context.Context, id string) error {
	return m.Delete(ctx, id)
}

func TestNewBoardService(t *testing.T) {
	mockBoardRepo := newMockBoardRepository()
	mockColumnRepo := newMockColumnRepository()
	service := NewBoardService(mockBoardRepo, mockColumnRepo)

	if service == nil {
		t.Error("NewBoardService() should return non-nil service")
	}
}

func TestBoardService_Create(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		title       string
		color       string
		expectError bool
	}{
		{
			name:        "Valid board creation",
			userID:      "user123",
			title:       "My Board",
			color:       "#FF5733",
			expectError: false,
		},
		{
			name:        "Board creation without color",
			userID:      "user456",
			title:       "Test Board",
			color:       "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBoardRepo := newMockBoardRepository()
			mockColumnRepo := newMockColumnRepository()
			service := NewBoardService(mockBoardRepo, mockColumnRepo)
			ctx := context.Background()

			board, err := service.Create(ctx, tt.userID, tt.title, tt.color)

			if tt.expectError {
				if err == nil {
					t.Errorf("Create() expected error but got nil")
					return
				}
				return
			}

			if err != nil {
				t.Errorf("Create() unexpected error = %v", err)
				return
			}

			if board == nil {
				t.Error("Create() should return non-nil board")
				return
			}

			if board.Title != tt.title {
				t.Errorf("Create() title = %v, want %v", board.Title, tt.title)
			}

			if board.UserID != tt.userID {
				t.Errorf("Create() userID = %v, want %v", board.UserID, tt.userID)
			}

			if board.ID == "" {
				t.Error("Create() should generate board ID")
			}

			if len(board.Columns) != 3 {
				t.Errorf("Create() should create 3 default columns, got %d", len(board.Columns))
			}

			columnTitles := []string{board.Columns[0].Title, board.Columns[1].Title, board.Columns[2].Title}
			expectedTitles := []string{"To Do", "In Progress", "Done"}
			for i, expected := range expectedTitles {
				if columnTitles[i] != expected {
					t.Errorf("Create() column %d title = %v, want %v", i, columnTitles[i], expected)
				}
			}
		})
	}
}

func TestBoardService_FindByID(t *testing.T) {
	tests := []struct {
		name          string
		setupBoard    bool
		boardID       string
		boardUserID   string
		requestUserID string
		expectError   bool
		errorType     string
	}{
		{
			name:          "Valid find by ID",
			setupBoard:    true,
			boardID:       "board123",
			boardUserID:   "user123",
			requestUserID: "user123",
			expectError:   false,
		},
		{
			name:          "Find by ID with wrong user",
			setupBoard:    true,
			boardID:       "board123",
			boardUserID:   "user123",
			requestUserID: "user456",
			expectError:   true,
			errorType:     "unauthorized",
		},
		{
			name:          "Find non-existent board",
			setupBoard:    false,
			boardID:       "nonexistent",
			boardUserID:   "user123",
			requestUserID: "user123",
			expectError:   true,
			errorType:     "not_found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBoardRepo := newMockBoardRepository()
			mockColumnRepo := newMockColumnRepository()
			service := NewBoardService(mockBoardRepo, mockColumnRepo)
			ctx := context.Background()

			if tt.setupBoard {
				board := &models.Board{
					ID:     tt.boardID,
					Title:  "Test Board",
					UserID: tt.boardUserID,
				}
				err := mockBoardRepo.Create(ctx, board)
				if err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			board, err := service.FindByID(ctx, tt.boardID, tt.requestUserID)

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

			if board == nil {
				t.Error("FindByID() should return non-nil board")
				return
			}

			if board.ID != tt.boardID {
				t.Errorf("FindByID() board ID = %v, want %v", board.ID, tt.boardID)
			}
		})
	}
}

func TestBoardService_FindByUserID(t *testing.T) {
	tests := []struct {
		name          string
		setupBoards   int
		userID        string
		expectedCount int
		expectError   bool
	}{
		{
			name:          "Find user's boards",
			setupBoards:   3,
			userID:        "user123",
			expectedCount: 3,
			expectError:   false,
		},
		{
			name:          "Find boards for user with none",
			setupBoards:   0,
			userID:        "user456",
			expectedCount: 0,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBoardRepo := newMockBoardRepository()
			mockColumnRepo := newMockColumnRepository()
			service := NewBoardService(mockBoardRepo, mockColumnRepo)
			ctx := context.Background()

			for i := 0; i < tt.setupBoards; i++ {
				board := &models.Board{
					Title:  "Test Board",
					UserID: tt.userID,
				}
				err := mockBoardRepo.Create(ctx, board)
				if err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			boards, err := service.FindByUserID(ctx, tt.userID)

			if tt.expectError {
				if err == nil {
					t.Error("FindByUserID() expected error but got nil")
					return
				}
				return
			}

			if err != nil {
				t.Errorf("FindByUserID() unexpected error = %v", err)
				return
			}

			if len(boards) != tt.expectedCount {
				t.Errorf("FindByUserID() returned %d boards, want %d", len(boards), tt.expectedCount)
			}
		})
	}
}

func TestBoardService_Update(t *testing.T) {
	tests := []struct {
		name          string
		setupBoard    bool
		boardID       string
		boardUserID   string
		requestUserID string
		title         string
		color         string
		expectError   bool
		errorType     string
	}{
		{
			name:          "Valid update title",
			setupBoard:    true,
			boardID:       "board123",
			boardUserID:   "user123",
			requestUserID: "user123",
			title:         "Updated Title",
			color:         "",
			expectError:   false,
		},
		{
			name:          "Valid update color",
			setupBoard:    true,
			boardID:       "board123",
			boardUserID:   "user123",
			requestUserID: "user123",
			title:         "",
			color:         "#00FF00",
			expectError:   false,
		},
		{
			name:          "Update by wrong user",
			setupBoard:    true,
			boardID:       "board123",
			boardUserID:   "user123",
			requestUserID: "user456",
			title:         "Updated Title",
			color:         "",
			expectError:   true,
			errorType:     "unauthorized",
		},
		{
			name:          "Update non-existent board",
			setupBoard:    false,
			boardID:       "nonexistent",
			boardUserID:   "user123",
			requestUserID: "user123",
			title:         "Updated Title",
			color:         "",
			expectError:   true,
			errorType:     "not_found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBoardRepo := newMockBoardRepository()
			mockColumnRepo := newMockColumnRepository()
			service := NewBoardService(mockBoardRepo, mockColumnRepo)
			ctx := context.Background()

			if tt.setupBoard {
				board := &models.Board{
					ID:     tt.boardID,
					Title:  "Original Title",
					UserID: tt.boardUserID,
					Color:  "#FF5733",
				}
				err := mockBoardRepo.Create(ctx, board)
				if err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			board, err := service.Update(ctx, tt.boardID, tt.requestUserID, tt.title, tt.color)

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

			if board == nil {
				t.Error("Update() should return non-nil board")
				return
			}

			if tt.title != "" && board.Title != tt.title {
				t.Errorf("Update() title = %v, want %v", board.Title, tt.title)
			}

			if tt.color != "" && board.Color != tt.color {
				t.Errorf("Update() color = %v, want %v", board.Color, tt.color)
			}
		})
	}
}

func TestBoardService_Delete(t *testing.T) {
	tests := []struct {
		name          string
		setupBoard    bool
		boardID       string
		boardUserID   string
		requestUserID string
		expectError   bool
		errorType     string
	}{
		{
			name:          "Valid delete",
			setupBoard:    true,
			boardID:       "board123",
			boardUserID:   "user123",
			requestUserID: "user123",
			expectError:   false,
		},
		{
			name:          "Delete by wrong user",
			setupBoard:    true,
			boardID:       "board123",
			boardUserID:   "user123",
			requestUserID: "user456",
			expectError:   true,
			errorType:     "unauthorized",
		},
		{
			name:          "Delete non-existent board",
			setupBoard:    false,
			boardID:       "nonexistent",
			boardUserID:   "user123",
			requestUserID: "user123",
			expectError:   true,
			errorType:     "not_found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBoardRepo := newMockBoardRepository()
			mockColumnRepo := newMockColumnRepository()
			service := NewBoardService(mockBoardRepo, mockColumnRepo)
			ctx := context.Background()

			if tt.setupBoard {
				board := &models.Board{
					ID:     tt.boardID,
					Title:  "Test Board",
					UserID: tt.boardUserID,
				}
				err := mockBoardRepo.Create(ctx, board)
				if err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			err := service.Delete(ctx, tt.boardID, tt.requestUserID)

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

			_, err = service.FindByID(ctx, tt.boardID, tt.requestUserID)
			if err == nil {
				t.Error("Delete() board should no longer be accessible")
			}
		})
	}
}

func TestBoardService_Integration(t *testing.T) {
	mockBoardRepo := newMockBoardRepository()
	mockColumnRepo := newMockColumnRepository()
	service := NewBoardService(mockBoardRepo, mockColumnRepo)
	ctx := context.Background()

	userID := "user123"
	title := "My Kanban Board"
	color := "#FF5733"

	board, err := service.Create(ctx, userID, title, color)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	if board == nil || board.ID == "" {
		t.Fatal("Create() should return board with ID")
	}

	foundBoard, err := service.FindByID(ctx, board.ID, userID)
	if err != nil {
		t.Fatalf("FindByID() failed: %v", err)
	}

	if foundBoard.ID != board.ID {
		t.Errorf("FindByID() returned wrong board: got %v, want %v", foundBoard.ID, board.ID)
	}

	updatedTitle := "Updated Kanban Board"
	updatedBoard, err := service.Update(ctx, board.ID, userID, updatedTitle, "")
	if err != nil {
		t.Fatalf("Update() failed: %v", err)
	}

	if updatedBoard.Title != updatedTitle {
		t.Errorf("Update() title = %v, want %v", updatedBoard.Title, updatedTitle)
	}

	boards, err := service.FindByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("FindByUserID() failed: %v", err)
	}

	if len(boards) != 1 {
		t.Errorf("FindByUserID() returned %d boards, want 1", len(boards))
	}

	err = service.Delete(ctx, board.ID, userID)
	if err != nil {
		t.Fatalf("Delete() failed: %v", err)
	}

	_, err = service.FindByID(ctx, board.ID, userID)
	if err == nil {
		t.Error("Board should not be found after deletion")
	}
}
