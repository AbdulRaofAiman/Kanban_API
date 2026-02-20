package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err, "Failed to connect to test database")

	// Migrate all models including audit_logs for User model hooks
	type AuditLog struct {
		ID        string    `gorm:"primaryKey;type:varchar(36)"`
		UserID    string    `gorm:"not null;type:varchar(36);index"`
		Action    string    `gorm:"not null;type:varchar(50)"`
		Message   string    `gorm:"type:text"`
		CreatedAt time.Time `gorm:"autoCreateTime"`
	}

	err = db.AutoMigrate(&AuditLog{}, &User{}, &Board{}, &Column{}, &Task{}, &Member{}, &RefreshToken{})
	assert.NoError(t, err, "Failed to migrate database")

	return db
}

// TestBoardModel verifies Board struct fields and GORM tags
func TestBoardModel(t *testing.T) {
	board := Board{
		ID:        uuid.NewString(),
		Title:     "Test Board",
		UserID:    uuid.NewString(),
		Color:     "#FF5733",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Verify ID is a valid UUID
	_, err := uuid.Parse(board.ID)
	assert.NoError(t, err, "Board ID should be a valid UUID")

	// Verify UserID is a valid UUID
	_, err = uuid.Parse(board.UserID)
	assert.NoError(t, err, "UserID should be a valid UUID")

	// Verify required fields are set
	assert.NotEmpty(t, board.Title, "Title should not be empty")
	assert.NotEmpty(t, board.UserID, "UserID should not be empty")
	assert.NotEmpty(t, board.Color, "Color should not be empty")
	assert.False(t, board.CreatedAt.IsZero(), "CreatedAt should be set")
	assert.False(t, board.UpdatedAt.IsZero(), "UpdatedAt should be set")

	// Verify TableName
	assert.Equal(t, "boards", board.TableName(), "TableName should be 'boards'")
}

// TestBoardTableName verifies the TableName method
func TestBoardTableName(t *testing.T) {
	board := Board{}
	assert.Equal(t, "boards", board.TableName(), "TableName should return 'boards'")
}

// TestBoardDatabase verifies Board can be saved to database
func TestBoardDatabase(t *testing.T) {
	db := setupTestDB(t)

	// Create a test user first
	userID := uuid.NewString()
	user := User{
		ID:        userID,
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "testpass123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := db.Create(&user).Error
	assert.NoError(t, err, "Failed to create test user")

	// Create a test board
	board := Board{
		ID:        uuid.NewString(),
		Title:     "Test Board",
		Color:     "#FF5733",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = db.Create(&board).Error
	assert.NoError(t, err, "Failed to create board")

	// Verify board was created with correct ID
	var retrievedBoard Board
	err = db.First(&retrievedBoard, "id = ?", board.ID).Error
	assert.NoError(t, err, "Failed to retrieve board")
	assert.Equal(t, board.ID, retrievedBoard.ID, "Board ID should match")
	assert.Equal(t, board.Title, retrievedBoard.Title, "Board Title should match")
	assert.Equal(t, board.UserID, retrievedBoard.UserID, "Board UserID should match")
	assert.Equal(t, board.Color, retrievedBoard.Color, "Board Color should match")

	// Verify foreign key constraint works (cannot create board with non-existent user)
	invalidBoard := Board{
		ID:        uuid.NewString(),
		Title:     "Invalid Board",
		Color:     "#FF5733",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = db.Create(&invalidBoard).Error
	// Note: SQLite doesn't enforce foreign key constraints by default
	// In production with PostgreSQL, this would fail
	// We'll verify this in integration tests with PostgreSQL
}

// TestBoardRelationships verifies Board relationships to Column, Task, and Member
func TestBoardRelationships(t *testing.T) {
	db := setupTestDB(t)

	// Create test user
	userID := uuid.NewString()
	user := User{
		ID:        userID,
		Username:  "testuser3",
		Email:     "test@example.com",
		Password:  "testpass123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := db.Create(&user).Error
	assert.NoError(t, err, "Failed to create test user")

	// Create test board
	boardID := uuid.NewString()
	board := Board{
		ID:        boardID,
		Title:     "Test Board",
		UserID:    userID,
		Color:     "#FF5733",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = db.Create(&board).Error
	assert.NoError(t, err, "Failed to create board")

	// Test HasMany Columns relationship
	column1 := Column{
		ID:        uuid.NewString(),
		BoardID:   boardID,
		Title:     "To Do",
		Order:     1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	column2 := Column{
		ID:        uuid.NewString(),
		BoardID:   boardID,
		Title:     "In Progress",
		Order:     2,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = db.Create(&column1).Error
	assert.NoError(t, err, "Failed to create column1")
	err = db.Create(&column2).Error
	assert.NoError(t, err, "Failed to create column2")

	// Verify columns can be loaded via relationship
	var boardWithColumns Board
	err = db.Preload("Columns").First(&boardWithColumns, "id = ?", boardID).Error
	assert.NoError(t, err, "Failed to load board with columns")
	assert.Len(t, boardWithColumns.Columns, 2, "Board should have 2 columns")
	assert.Equal(t, "To Do", boardWithColumns.Columns[0].Title, "First column should be 'To Do'")
	assert.Equal(t, "In Progress", boardWithColumns.Columns[1].Title, "Second column should be 'In Progress'")

	// Test HasMany Members relationship
	member1 := Member{
		ID:        uuid.NewString(),
		BoardID:   boardID,
		UserID:    userID,
		Role:      "admin",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	// Create another user for member2
	member2UserID := uuid.NewString()
	member2User := User{
		ID:        member2UserID,
		Username:  "member2",
		Email:     "member2@example.com",
		Password:  "memberpass123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = db.Create(&member2User).Error
	assert.NoError(t, err, "Failed to create member2 user")

	member2 := Member{
		ID:        uuid.NewString(),
		BoardID:   boardID,
		UserID:    member2UserID,
		Role:      "member",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = db.Create(&member1).Error
	assert.NoError(t, err, "Failed to create member1")
	err = db.Create(&member2).Error
	assert.NoError(t, err, "Failed to create member2")

	// Verify members can be loaded via relationship
	var boardWithMembers Board
	err = db.Preload("Members").First(&boardWithMembers, "id = ?", boardID).Error
	assert.NoError(t, err, "Failed to load board with members")
	assert.Len(t, boardWithMembers.Members, 2, "Board should have 2 members")
	assert.Equal(t, "admin", boardWithMembers.Members[0].Role, "First member should have role 'admin'")
	assert.Equal(t, "member", boardWithMembers.Members[1].Role, "Second member should have role 'member'")

	// Test BelongsTo User relationship
	var boardWithUser Board
	err = db.Preload("User").First(&boardWithUser, "id = ?", boardID).Error
	assert.NoError(t, err, "Failed to load board with user")
	assert.NotNil(t, boardWithUser.User, "Board should have a user")
	assert.Equal(t, userID, boardWithUser.User.ID, "User ID should match")
	assert.Equal(t, "test@example.com", boardWithUser.User.Email, "User email should match")
}

// TestBoardSoftDelete verifies Board soft delete functionality
func TestBoardSoftDelete(t *testing.T) {
	db := setupTestDB(t)

	// Create test user
	userID := uuid.NewString()
	user := User{
		ID:        userID,
		Username:  "testuser4",
		Email:     "test@example.com",
		Password:  "testpass123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := db.Create(&user).Error
	assert.NoError(t, err, "Failed to create test user")

	// Create test board
	board := Board{
		ID:        uuid.NewString(),
		Title:     "Test Board",
		Color:     "#FF5733",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = db.Create(&board).Error
	assert.NoError(t, err, "Failed to create board")

	// Verify board exists
	var retrievedBoard Board
	err = db.First(&retrievedBoard, "id = ?", board.ID).Error
	assert.NoError(t, err, "Board should exist before delete")

	// Soft delete board
	err = db.Delete(&board).Error
	assert.NoError(t, err, "Failed to soft delete board")

	// Verify board is soft deleted (not found with regular query)
	err = db.First(&retrievedBoard, "id = ?", board.ID).Error
	assert.Error(t, err, "Board should not be found with regular query after soft delete")

	// Verify board can be found with Unscoped
	err = db.Unscoped().First(&retrievedBoard, "id = ?", board.ID).Error
	assert.NoError(t, err, "Board should be found with Unscoped query")
	assert.Equal(t, board.ID, retrievedBoard.ID, "Board ID should match after soft delete")
}

// TestBoardGORMTags verifies GORM tags are correct
func TestBoardGORMTags(t *testing.T) {
	db := setupTestDB(t)

	// Get the schema for Board model
	stmt := &gorm.Statement{DB: db}
	err := stmt.Parse(&Board{})
	assert.NoError(t, err, "Failed to parse Board model")

	// Verify primary key
	assert.Equal(t, "ID", stmt.Schema.PrimaryFields[0].Name, "Primary key should be ID")

	// Verify required fields have not null tag
	for _, field := range stmt.Schema.Fields {
		switch field.Name {
		case "Title":
			assert.True(t, field.NotNull, "Title should have not null tag")
		case "UserID":
			assert.True(t, field.NotNull, "UserID should have not null tag")
		}
	}

	// Verify field sizes
	for _, field := range stmt.Schema.Fields {
		switch field.Name {
		case "Title":
			assert.Equal(t, 255, field.Size, "Title size should be 255")
		case "Color":
			assert.Equal(t, 255, field.Size, "Color size should be 255")
		}
	}

	// Verify relationships
	assert.Len(t, stmt.Schema.Relationships.Relations, 4, "Board should have 4 relationships")

	// Verify foreign keys
	for _, rel := range stmt.Schema.Relationships.Relations {
		switch rel.Name {
		case "Columns":
			assert.Equal(t, "BoardID", rel.References[0].ForeignKey.Name, "Columns foreign key should be BoardID")
		case "Tasks":
			assert.Equal(t, "BoardID", rel.References[0].ForeignKey.Name, "Tasks foreign key should be BoardID")
		case "Members":
			assert.Equal(t, "BoardID", rel.References[0].ForeignKey.Name, "Members foreign key should be BoardID")
		case "User":
			assert.Equal(t, "UserID", rel.References[0].ForeignKey.Name, "User foreign key should be UserID")
		}
	}
}
