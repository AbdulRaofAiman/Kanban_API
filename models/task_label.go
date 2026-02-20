package models

import (
	"time"
)

// TaskLabel is a join model for many-to-many relationship between Task and Label
type TaskLabel struct {
	TaskID    string    `gorm:"primaryKey;type:varchar(36);not null" json:"task_id"`
	LabelID   string    `gorm:"primaryKey;type:varchar(36);not null" json:"label_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relations
	Task  *Task  `gorm:"foreignKey:TaskID;constraint:OnDelete:CASCADE" json:"task,omitempty"`
	Label *Label `gorm:"foreignKey:LabelID;constraint:OnDelete:CASCADE" json:"label,omitempty"`
}

// TableName specifies the table name for TaskLabel model
func (TaskLabel) TableName() string {
	return "task_labels"
}
