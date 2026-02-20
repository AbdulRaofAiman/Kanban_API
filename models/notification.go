package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Notification represents a user notification in the system
type Notification struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID    string         `gorm:"not null;type:varchar(36);index" json:"user_id"`
	Message   string         `gorm:"not null;type:text" json:"message"`
	ReadAt    *time.Time     `gorm:"index" json:"read_at,omitempty"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

// TableName specifies the table name for Notification model
func (Notification) TableName() string {
	return "notifications"
}

// BeforeCreate is a GORM hook called before creating a notification
func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	// Set ID if not already set
	if n.ID == "" {
		n.ID = uuid.New().String()
	}
	return nil
}

// MarkAsRead marks the notification as read
func (n *Notification) MarkAsRead() {
	now := time.Now()
	n.ReadAt = &now
}

// IsRead returns true if the notification has been read
func (n *Notification) IsRead() bool {
	return n.ReadAt != nil
}
