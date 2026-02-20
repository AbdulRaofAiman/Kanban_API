package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Label represents a tag or category for tasks (e.g., Bug, Feature, High Priority)
type Label struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name      string         `gorm:"not null;size:255;uniqueIndex" json:"name"`
	Color     string         `gorm:"size:255" json:"color"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Many-to-many relationship with Task via TaskLabel join table
	Tasks []*Task `gorm:"many2many:task_labels;foreignKey:ID;joinForeignKey:LabelID;references:ID;joinReferences:TaskID" json:"tasks,omitempty"`
}

// TableName specifies the table name for Label model
func (Label) TableName() string {
	return "labels"
}

// BeforeCreate is a GORM hook called before creating a label
func (l *Label) BeforeCreate(tx *gorm.DB) error {
	// Set ID if not already set
	if l.ID == "" {
		l.ID = uuid.NewString()
	}
	return nil
}
