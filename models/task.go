package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Task represents a kanban task within a column
type Task struct {
	ID          string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	ColumnID    string         `gorm:"not null;type:varchar(36);index:task_column" json:"column_id"`
	Title       string         `gorm:"not null;type:varchar(255)" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Deadline    *time.Time     `gorm:"index:task_deadline" json:"deadline,omitempty"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Comments    []Comment    `gorm:"foreignKey:TaskID;constraint:OnDelete:CASCADE" json:"comments,omitempty"`
	Labels      []Label      `gorm:"many2many:task_labels;constraint:OnDelete:CASCADE" json:"labels,omitempty"`
	Attachments []Attachment `gorm:"foreignKey:TaskID;constraint:OnDelete:CASCADE" json:"attachments,omitempty"`
	Column      *Column      `gorm:"foreignKey:ColumnID;constraint:OnDelete:CASCADE" json:"column,omitempty"`
}

// TableName specifies the table name for Task model
func (Task) TableName() string {
	return "tasks"
}

// BeforeCreate hook to generate UUID before insertion
func (t *Task) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.NewString()
	}
	return nil
}
