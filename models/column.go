package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Column represents a column in a kanban board
type Column struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	BoardID   string         `gorm:"not null;type:varchar(36);index:column_board" json:"board_id"`
	Title     string         `gorm:"not null;type:varchar(255)" json:"title"`
	Order     int            `gorm:"not null" json:"order"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Tasks []Task `gorm:"foreignKey:ColumnID;constraint:OnDelete:CASCADE" json:"tasks,omitempty"`
	Board *Board `gorm:"foreignKey:BoardID;constraint:OnDelete:CASCADE" json:"board,omitempty"`
}

// TableName specifies the table name for Column model
func (Column) TableName() string {
	return "columns"
}

// BeforeCreate hook to generate UUID if ID is empty
func (c *Column) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}
