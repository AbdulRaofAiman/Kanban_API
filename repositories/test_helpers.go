package repositories

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"kanban-backend/models"
)

func setupRepositoryTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Board{}, &models.Member{}, &models.Column{}, &models.Task{}, &models.RefreshToken{}, &models.Comment{}, &models.Label{}, &models.Attachment{})
	if err != nil {
		t.Fatal(err)
	}

	type AuditLog struct {
		ID        string `gorm:"primaryKey;type:varchar(36)"`
		UserID    string `gorm:"not null;type:varchar(36);index"`
		Action    string `gorm:"not null;type:varchar(50)"`
		Message   string `gorm:"type:text"`
		CreatedAt int64  `gorm:"autoCreateTime"`
	}
	err = db.Table("audit_logs").AutoMigrate(&AuditLog{})
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func createTestUser(db *gorm.DB, username, email string) *models.User {
	user := &models.User{
		Username: username,
		Email:    email,
		Password: "password123",
	}
	db.Create(user)
	return user
}

func createTestBoard(db *gorm.DB, userID string) *models.Board {
	board := &models.Board{
		Title:  "Test Board",
		UserID: userID,
		Color:  "#FF5733",
	}
	db.Create(board)
	return board
}

func createTestColumn(db *gorm.DB, boardID string) *models.Column {
	column := &models.Column{
		BoardID: boardID,
		Title:   "Test Column",
		Order:   1,
	}
	db.Create(column)
	return column
}
