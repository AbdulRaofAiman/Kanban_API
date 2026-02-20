package models

import (
	"time"

	"gorm.io/gorm"
)

// RefreshToken represents a refresh token for user authentication
type RefreshToken struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID    string         `gorm:"not null;type:varchar(36);index" json:"user_id"`
	Token     string         `gorm:"not null;type:text;uniqueIndex" json:"token"`
	ExpiresAt time.Time      `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

// TableName specifies the table name for RefreshToken model
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
