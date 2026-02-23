package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Member represents a board member relationship (many-to-many between User and Board)
// This is a placeholder - will be fully implemented in Task 12
type Member struct {
	ID        string         `gorm:"type:uuid;primaryKey" json:"id"`
	BoardID   string         `gorm:"type:uuid;not null;index" json:"board_id"`
	UserID    string         `gorm:"type:uuid;not null;index" json:"user_id"`
	Role      string         `gorm:"size:50;default:member" json:"role"` // owner, admin, member
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// BeforeCreate is a GORM hook called before creating a member
func (m *Member) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}

// Board represents a Kanban board
type Board struct {
	ID        string         `gorm:"type:uuid;primaryKey" json:"id"`
	Title     string         `gorm:"size:255;not null" json:"title"`
	UserID    string         `gorm:"type:uuid;not null;index" json:"user_id"`
	Color     string         `gorm:"size:255" json:"color"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	Columns []Column `gorm:"foreignKey:BoardID" json:"columns,omitempty"`
	Members []Member `gorm:"foreignKey:BoardID" json:"members,omitempty"`
	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// BeforeCreate is a GORM hook called before creating a board
func (b *Board) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}

// AfterCreate is a GORM hook called after creating a board
func (b *Board) AfterCreate(tx *gorm.DB) error {
	return b.logAudit(tx, "create", "Board created")
}

// logAudit logs board actions to the audit_log table
func (b *Board) logAudit(tx *gorm.DB, action, message string) error {
	type AuditLog struct {
		ID        string    `gorm:"primaryKey;type:varchar(36)"`
		UserID    string    `gorm:"not null;type:varchar(36);index"`
		Action    string    `gorm:"not null;type:varchar(50)"`
		Message   string    `gorm:"type:text"`
		CreatedAt time.Time `gorm:"autoCreateTime"`
	}

	auditLog := AuditLog{
		ID:        uuid.New().String(),
		UserID:    b.UserID,
		Action:    action,
		Message:   message,
		CreatedAt: time.Now(),
	}

	return tx.Table("audit_logs").Create(&auditLog).Error
}

// TableName specifies the table name for Board model
func (Board) TableName() string {
	return "boards"
}
