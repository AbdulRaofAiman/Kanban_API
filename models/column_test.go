package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestColumnModel verifies Column struct fields and GORM tags
func TestColumnModel(t *testing.T) {
	column := Column{
		ID:        uuid.NewString(),
		BoardID:   uuid.NewString(),
		Title:     "To Do",
		Order:     1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Verify ID is a valid UUID
	_, err := uuid.Parse(column.ID)
	assert.NoError(t, err, "Column ID should be a valid UUID")

	// Verify BoardID is a valid UUID
	_, err = uuid.Parse(column.BoardID)
	assert.NoError(t, err, "BoardID should be a valid UUID")

	// Verify required fields are set
	assert.NotEmpty(t, column.ID, "ID should not be empty")
	assert.NotEmpty(t, column.BoardID, "BoardID should not be empty")
	assert.NotEmpty(t, column.Title, "Title should not be empty")
	assert.GreaterOrEqual(t, column.Order, 0, "Order should be non-negative")
	assert.False(t, column.CreatedAt.IsZero(), "CreatedAt should be set")
	assert.False(t, column.UpdatedAt.IsZero(), "UpdatedAt should be set")

	// Verify TableName
	assert.Equal(t, "columns", column.TableName(), "TableName should be 'columns'")
}

// TestColumnTableName verifies TableName method
func TestColumnTableName(t *testing.T) {
	column := Column{}
	assert.Equal(t, "columns", column.TableName(), "TableName should return 'columns'")
}

// TestColumnDatabase verifies Column can be saved to database
func TestColumnDatabase(t *testing.T) {
	db := setupTestDB(t)

	// Create a test user first
	userID := uuid.NewString()
	user := User{
		ID:        userID,
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := db.Create(&user).Error
	assert.NoError(t, err, "Failed to create test user")

	// Create a test board
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

	// Create a test column
	columnID := uuid.NewString()
	column := Column{
		ID:        columnID,
		BoardID:   boardID,
		Title:     "To Do",
		Order:     1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = db.Create(&column).Error
	assert.NoError(t, err, "Failed to create column")

	// Verify column was created with correct ID
	var retrievedColumn Column
	err = db.First(&retrievedColumn, "id = ?", columnID).Error
	assert.NoError(t, err, "Failed to retrieve column")
	assert.Equal(t, column.ID, retrievedColumn.ID, "Column ID should match")
	assert.Equal(t, column.Title, retrievedColumn.Title, "Column Title should match")
	assert.Equal(t, column.BoardID, retrievedColumn.BoardID, "Column BoardID should match")
	assert.Equal(t, column.Order, retrievedColumn.Order, "Column Order should match")
}

// TestColumnRelationships verifies Column relationships to Task and Board
func TestColumnRelationships(t *testing.T) {
	db := setupTestDB(t)

	// Create test user
	userID := uuid.NewString()
	user := User{
		ID:        userID,
		Email:     "test@example.com",
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

	// Create test column
	columnID := uuid.NewString()
	column := Column{
		ID:        columnID,
		BoardID:   boardID,
		Title:     "In Progress",
		Order:     2,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = db.Create(&column).Error
	assert.NoError(t, err, "Failed to create column")

	// Test HasMany Tasks relationship
	task1 := Task{
		ID:        uuid.NewString(),
		ColumnID:  columnID,
		Title:     "Test Task 1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	task2 := Task{
		ID:        uuid.NewString(),
		ColumnID:  columnID,
		Title:     "Test Task 2",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = db.Create(&task1).Error
	assert.NoError(t, err, "Failed to create task1")
	err = db.Create(&task2).Error
	assert.NoError(t, err, "Failed to create task2")

	// Verify tasks can be loaded via relationship
	var columnWithTasks Column
	err = db.Preload("Tasks").First(&columnWithTasks, "id = ?", columnID).Error
	assert.NoError(t, err, "Failed to load column with tasks")
	assert.Len(t, columnWithTasks.Tasks, 2, "Column should have 2 tasks")
	assert.Equal(t, "Test Task 1", columnWithTasks.Tasks[0].Title, "First task should be 'Test Task 1'")
	assert.Equal(t, "Test Task 2", columnWithTasks.Tasks[1].Title, "Second task should be 'Test Task 2'")

	// Test BelongsTo Board relationship
	var columnWithBoard Column
	err = db.Preload("Board").First(&columnWithBoard, "id = ?", columnID).Error
	assert.NoError(t, err, "Failed to load column with board")
	assert.NotNil(t, columnWithBoard.Board, "Column should have a board")
	assert.Equal(t, boardID, columnWithBoard.Board.ID, "Board ID should match")
	assert.Equal(t, "Test Board", columnWithBoard.Board.Title, "Board title should match")
}

// TestColumnSoftDelete verifies Column soft delete functionality
func TestColumnSoftDelete(t *testing.T) {
	db := setupTestDB(t)

	// Create test user
	userID := uuid.NewString()
	user := User{
		ID:        userID,
		Email:     "test@example.com",
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

	// Create test column
	column := Column{
		ID:        uuid.NewString(),
		BoardID:   boardID,
		Title:     "To Do",
		Order:     1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = db.Create(&column).Error
	assert.NoError(t, err, "Failed to create column")

	// Verify column exists
	var retrievedColumn Column
	err = db.First(&retrievedColumn, "id = ?", column.ID).Error
	assert.NoError(t, err, "Column should exist before delete")

	// Soft delete column
	err = db.Delete(&column).Error
	assert.NoError(t, err, "Failed to soft delete column")

	// Verify column is soft deleted (not found with regular query)
	err = db.First(&retrievedColumn, "id = ?", column.ID).Error
	assert.Error(t, err, "Column should not be found with regular query after soft delete")

	// Verify column can be found with Unscoped
	err = db.Unscoped().First(&retrievedColumn, "id = ?", column.ID).Error
	assert.NoError(t, err, "Column should be found with Unscoped query")
	assert.Equal(t, column.ID, retrievedColumn.ID, "Column ID should match after soft delete")
}

// TestColumnGORMTags verifies GORM tags are correct
func TestColumnGORMTags(t *testing.T) {
	db := setupTestDB(t)

	// Get schema for Column model
	stmt := &gorm.Statement{DB: db}
	err := stmt.Parse(&Column{})
	assert.NoError(t, err, "Failed to parse Column model")

	// Verify primary key
	assert.Equal(t, "ID", stmt.Schema.PrimaryFields[0].Name, "Primary key should be ID")

	// Verify required fields have not null tag
	for _, field := range stmt.Schema.Fields {
		switch field.Name {
		case "BoardID":
			assert.True(t, field.NotNull, "BoardID should have not null tag")
		case "Title":
			assert.True(t, field.NotNull, "Title should have not null tag")
		case "Order":
			assert.True(t, field.NotNull, "Order should have not null tag")
		}
	}

	// Verify relationships (at minimum 2 relationships should exist)
	assert.GreaterOrEqual(t, len(stmt.Schema.Relationships.Relations), 2, "Column should have at least 2 relationships")

	// Verify foreign keys for expected relationships
	foundTasks := false
	foundBoard := false
	for _, rel := range stmt.Schema.Relationships.Relations {
		switch rel.Name {
		case "Tasks":
			assert.Equal(t, "ColumnID", rel.References[0].ForeignKey.Name, "Tasks foreign key should be ColumnID")
			foundTasks = true
		case "Board":
			assert.Equal(t, "BoardID", rel.References[0].ForeignKey.Name, "Board foreign key should be BoardID")
			foundBoard = true
		}
	}
	assert.True(t, foundTasks, "Tasks relationship should exist")
	assert.True(t, foundBoard, "Board relationship should exist")
}
