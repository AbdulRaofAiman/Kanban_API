package models

import (
	"fmt"
	"time"

	"kanban-backend/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Username  string         `gorm:"not null;type:varchar(50);uniqueIndex" json:"username"`
	Email     string         `gorm:"not null;type:varchar(255);uniqueIndex" json:"email"`
	Password  string         `gorm:"not null;type:text" json:"-"` // Never expose password in JSON
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	RefreshTokens []RefreshToken `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"refresh_tokens,omitempty"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

// BeforeCreate is a GORM hook called before creating a user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Set ID if not already set
	if u.ID == "" {
		u.ID = uuid.New().String()
	}

	// Hash password if it's not already hashed (bcrypt hashes start with $2a$, $2b$, or $2y$)
	if u.Password != "" && len(u.Password) != 60 {
		hashedPassword, err := utils.HashPassword(u.Password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		u.Password = hashedPassword
	}

	// Log to audit table
	return u.logAudit(tx, "create", "User created")
}

// AfterUpdate is a GORM hook called after updating a user
func (u *User) AfterUpdate(tx *gorm.DB) error {
	// Log to audit table
	return u.logAudit(tx, "update", "User updated")
}

// logAudit logs user actions to the audit_log table
func (u *User) logAudit(tx *gorm.DB, action, message string) error {
	// Note: This will create the audit_log table if it doesn't exist
	// We'll implement the AuditLog model in a separate file if needed
	type AuditLog struct {
		ID        string    `gorm:"primaryKey;type:varchar(36)"`
		UserID    string    `gorm:"not null;type:varchar(36);index"`
		Action    string    `gorm:"not null;type:varchar(50)"`
		Message   string    `gorm:"type:text"`
		CreatedAt time.Time `gorm:"autoCreateTime"`
	}

	auditLog := AuditLog{
		ID:        uuid.New().String(),
		UserID:    u.ID,
		Action:    action,
		Message:   message,
		CreatedAt: time.Now(),
	}

	return tx.Table("audit_logs").Create(&auditLog).Error
}

// SetPassword sets a new password for the user (automatically hashed)
func (u *User) SetPassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	u.Password = hashedPassword
	return nil
}

// CheckPassword verifies if the provided password matches the stored password
func (u *User) CheckPassword(password string) error {
	return utils.CheckPassword(u.Password, password)
}
