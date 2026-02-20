package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Attachment represents a file attachment to a task
// This is a placeholder - will be fully implemented in Task 14
type Attachment struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	TaskID    string         `gorm:"type:varchar(36);not null;index" json:"task_id"`
	FileName  string         `gorm:"size:255;not null" json:"file_name"`
	FileURL   string         `gorm:"size:500;not null" json:"file_url"`
	FileSize  int64          `gorm:"default:0" json:"file_size"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	Task *Task `gorm:"foreignKey:TaskID;constraint:OnDelete:CASCADE" json:"task,omitempty"`
}

// TableName specifies the table name for Attachment model
func (Attachment) TableName() string {
	return "attachments"
}

// BeforeCreate is a GORM hook called before creating an attachment
func (a *Attachment) BeforeCreate(tx *gorm.DB) error {
	// Set ID if not already set
	if a.ID == "" {
		a.ID = uuid.NewString()
	}
	return nil
}
