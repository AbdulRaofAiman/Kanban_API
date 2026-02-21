package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Comment represents a comment on a task
type Comment struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	TaskID    string         `gorm:"type:varchar(36);not null;index" json:"task_id"`
	UserID    string         `gorm:"type:varchar(36);not null;index" json:"user_id"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	Task *Task `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for Comment model
func (Comment) TableName() string {
	return "comments"
}

// BeforeCreate is a GORM hook called before creating a comment
func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.NewString()
	}
	return nil
}
