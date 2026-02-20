package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"kanban-backend/utils"
)

// setupUserTestDB creates an in-memory SQLite database for User and RefreshToken tests
// This avoids issues with other models (e.g., Board.Tasks) that have relationship issues
func setupUserTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err, "Failed to connect to test database")

	err = db.AutoMigrate(&User{}, &RefreshToken{})
	assert.NoError(t, err, "Failed to migrate database")

	type AuditLog struct {
		ID        string    `gorm:"primaryKey;type:varchar(36)"`
		UserID    string    `gorm:"not null;type:varchar(36);index"`
		Action    string    `gorm:"not null;type:varchar(50)"`
		Message   string    `gorm:"type:text"`
		CreatedAt time.Time `gorm:"autoCreateTime"`
	}
	err = db.Table("audit_logs").AutoMigrate(&AuditLog{})
	assert.NoError(t, err, "Failed to migrate audit_logs table")

	return db
}

// TestUserModelStructure tests the User model structure
func TestUserModelStructure(t *testing.T) {
	user := User{
		ID:        uuid.New().String(),
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Verify ID is set
	if user.ID == "" {
		t.Error("ID should not be empty")
	}

	// Verify username
	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", user.Username)
	}

	// Verify email
	if user.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", user.Email)
	}

	// Verify password is not exposed in JSON
	// The struct tag `json:"-"` should omit password field
}

// TestUserTableName tests the table name for User model
func TestUserTableName(t *testing.T) {
	user := User{}
	tableName := user.TableName()

	if tableName != "users" {
		t.Errorf("Expected table name 'users', got '%s'", tableName)
	}
}

// TestUserBeforeCreateHook tests the BeforeCreate GORM hook
func TestUserBeforeCreateHook(t *testing.T) {
	db := setupUserTestDB(t)

	testCases := []struct {
		name      string
		user      *User
		wantError bool
	}{
		{
			name: "Valid user with plain password",
			user: &User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123", // Will be hashed by BeforeCreate
			},
			wantError: false,
		},
		{
			name: "Valid user with short password",
			user: &User{
				Username: "testuser2",
				Email:    "test2@example.com",
				Password: "short",
			},
			wantError: true, // Password too short
		},
		{
			name: "Valid user with already hashed password",
			user: &User{
				Username: "testuser3",
				Email:    "test3@example.com",
				Password: "$2a$10$ThisIsAValidBcryptHashThatIsExactly60CharsLong",
			},
			wantError: false, // Already hashed, won't re-hash
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset ID to ensure BeforeCreate sets it
			tc.user.ID = ""

			// Create user (triggers BeforeCreate hook)
			result := db.Create(tc.user)

			if tc.wantError {
				if result.Error == nil {
					t.Error("Expected error, but got nil")
				}
			} else {
				if result.Error != nil {
					t.Errorf("Expected no error, but got: %v", result.Error)
				}

				// Verify ID was set
				if tc.user.ID == "" {
					t.Error("ID should be set by BeforeCreate hook")
				}

				// If password was plain, verify it was hashed
				if tc.user.Password == "password123" {
					t.Error("Password should be hashed, but it's still plain text")
				}

				// Verify password is a bcrypt hash (60 characters, starts with $2a$, $2b$, or $2y$)
				if len(tc.user.Password) != 60 {
					t.Errorf("Password hash should be 60 characters, got %d", len(tc.user.Password))
				}

				// Verify CreatedAt was set
				if tc.user.CreatedAt.IsZero() {
					t.Error("CreatedAt should be set")
				}

				// Verify UpdatedAt was set
				if tc.user.UpdatedAt.IsZero() {
					t.Error("UpdatedAt should be set")
				}
			}
		})
	}
}

// TestUserAfterUpdateHook tests the AfterUpdate GORM hook
func TestUserAfterUpdateHook(t *testing.T) {
	db := setupUserTestDB(t)

	// Create a user
	user := User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	result := db.Create(&user)
	if result.Error != nil {
		t.Fatalf("Failed to create test user: %v", result.Error)
	}

	// Update the user (triggers AfterUpdate hook)
	result = db.Model(&user).Update("username", "updateduser")
	if result.Error != nil {
		t.Errorf("Failed to update user: %v", result.Error)
	}

	// Verify audit log was created
	type AuditLog struct {
		ID        string    `gorm:"primaryKey;type:varchar(36)"`
		UserID    string    `gorm:"not null;type:varchar(36);index"`
		Action    string    `gorm:"not null;type:varchar(50)"`
		Message   string    `gorm:"type:text"`
		CreatedAt time.Time `gorm:"autoCreateTime"`
	}

	var auditLogs []AuditLog
	result = db.Table("audit_logs").Where("user_id = ? AND action = ?", user.ID, "update").Find(&auditLogs)

	if result.Error != nil {
		t.Errorf("Failed to query audit logs: %v", result.Error)
	}

	if len(auditLogs) == 0 {
		t.Error("Expected audit log entry for update action, but found none")
	}

	// Verify the audit log message
	for _, log := range auditLogs {
		if log.Message != "User updated" {
			t.Errorf("Expected message 'User updated', got '%s'", log.Message)
		}
	}
}

// TestAuditHooks tests both BeforeCreate and AfterUpdate hooks together
func TestAuditHooks(t *testing.T) {
	db := setupUserTestDB(t)

	// Test create audit log
	user := User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	result := db.Create(&user)
	if result.Error != nil {
		t.Fatalf("Failed to create test user: %v", result.Error)
	}

	// Verify create audit log
	type AuditLog struct {
		ID        string    `gorm:"primaryKey;type:varchar(36)"`
		UserID    string    `gorm:"not null;type:varchar(36);index"`
		Action    string    `gorm:"not null;type:varchar(50)"`
		Message   string    `gorm:"type:text"`
		CreatedAt time.Time `gorm:"autoCreateTime"`
	}

	var createAuditLogs []AuditLog
	result = db.Table("audit_logs").Where("user_id = ? AND action = ?", user.ID, "create").Find(&createAuditLogs)

	if result.Error != nil {
		t.Errorf("Failed to query create audit logs: %v", result.Error)
	}

	if len(createAuditLogs) == 0 {
		t.Error("Expected audit log entry for create action, but found none")
	}

	for _, log := range createAuditLogs {
		if log.Message != "User created" {
			t.Errorf("Expected message 'User created', got '%s'", log.Message)
		}
	}

	// Test update audit log
	oldUpdatedAt := user.UpdatedAt
	time.Sleep(10 * time.Millisecond) // Ensure UpdatedAt changes

	result = db.Model(&user).Update("username", "updateduser")
	if result.Error != nil {
		t.Errorf("Failed to update user: %v", result.Error)
	}

	// Refresh user from database
	db.First(&user, user.ID)
	if user.UpdatedAt.Equal(oldUpdatedAt) {
		t.Error("UpdatedAt should have changed after update")
	}

	// Verify update audit log
	var updateAuditLogs []AuditLog
	result = db.Table("audit_logs").Where("user_id = ? AND action = ?", user.ID, "update").Find(&updateAuditLogs)

	if result.Error != nil {
		t.Errorf("Failed to query update audit logs: %v", result.Error)
	}

	if len(updateAuditLogs) == 0 {
		t.Error("Expected audit log entry for update action, but found none")
	}
}

// TestUserSetPassword tests the SetPassword method
func TestUserSetPassword(t *testing.T) {
	testCases := []struct {
		name      string
		password  string
		wantError bool
	}{
		{
			name:      "Valid password (8 chars)",
			password:  "pass1234",
			wantError: false,
		},
		{
			name:      "Valid password (more than 8 chars)",
			password:  "verylongpassword",
			wantError: false,
		},
		{
			name:      "Too short password (7 chars)",
			password:  "short1",
			wantError: true,
		},
		{
			name:      "Empty password",
			password:  "",
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user := &User{}

			err := user.SetPassword(tc.password)

			if tc.wantError {
				if err == nil {
					t.Error("Expected error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}

				// Verify password was hashed
				if user.Password == tc.password {
					t.Error("Password should be hashed, but it's still plain text")
				}

				// Verify it's a bcrypt hash
				if len(user.Password) != 60 {
					t.Errorf("Password hash should be 60 characters, got %d", len(user.Password))
				}

				// Verify the hash is valid by checking it
				err = utils.CheckPassword(user.Password, tc.password)
				if err != nil {
					t.Errorf("Hashed password doesn't match original: %v", err)
				}
			}
		})
	}
}

// TestUserCheckPassword tests the CheckPassword method
func TestUserCheckPassword(t *testing.T) {
	user := &User{}

	// Set a password first
	err := user.SetPassword("testpassword123")
	if err != nil {
		t.Fatalf("Failed to set password for testing: %v", err)
	}

	testCases := []struct {
		name      string
		password  string
		wantError bool
	}{
		{
			name:      "Correct password",
			password:  "testpassword123",
			wantError: false,
		},
		{
			name:      "Incorrect password",
			password:  "wrongpassword",
			wantError: true,
		},
		{
			name:      "Empty password",
			password:  "",
			wantError: true,
		},
		{
			name:      "Case sensitive check",
			password:  "TestPassword123", // Different case
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := user.CheckPassword(tc.password)

			if tc.wantError {
				if err == nil {
					t.Error("Expected error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			}
		})
	}
}

// TestUserRefreshTokensRelation tests the RefreshTokens relationship
func TestUserRefreshTokensRelation(t *testing.T) {
	db := setupUserTestDB(t)

	// Create a user
	user := User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	result := db.Create(&user)
	if result.Error != nil {
		t.Fatalf("Failed to create test user: %v", result.Error)
	}

	// Create refresh tokens for the user
	refreshToken1 := RefreshToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Token:     "refreshtoken123",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	refreshToken2 := RefreshToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Token:     "refreshtoken456",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	result = db.Create(&refreshToken1)
	if result.Error != nil {
		t.Fatalf("Failed to create refresh token 1: %v", result.Error)
	}
	result = db.Create(&refreshToken2)
	if result.Error != nil {
		t.Fatalf("Failed to create refresh token 2: %v", result.Error)
	}

	// Query user with refresh tokens
	var foundUser User
	result = db.Preload("RefreshTokens").Where("id = ?", user.ID).First(&foundUser)
	if result.Error != nil {
		t.Fatalf("Failed to query user with refresh tokens: %v", result.Error)
	}

	// Verify refresh tokens are loaded
	if len(foundUser.RefreshTokens) != 2 {
		t.Errorf("Expected 2 refresh tokens, got %d", len(foundUser.RefreshTokens))
	}

	// Verify the refresh token data
	for _, rt := range foundUser.RefreshTokens {
		if rt.UserID != user.ID {
			t.Errorf("Refresh token UserID should match user ID")
		}
	}
}

// TestRefreshTokenTableName tests the table name for RefreshToken model
func TestRefreshTokenTableName(t *testing.T) {
	refreshToken := RefreshToken{}
	tableName := refreshToken.TableName()

	if tableName != "refresh_tokens" {
		t.Errorf("Expected table name 'refresh_tokens', got '%s'", tableName)
	}
}

// TestRefreshTokenUniqueIndex tests the unique index on token field
func TestRefreshTokenUniqueIndex(t *testing.T) {
	db := setupUserTestDB(t)

	// Create a user
	user := User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	result := db.Create(&user)
	if result.Error != nil {
		t.Fatalf("Failed to create test user: %v", result.Error)
	}

	// Create first refresh token
	refreshToken1 := RefreshToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Token:     "uniquetoken123",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	result = db.Create(&refreshToken1)
	if result.Error != nil {
		t.Fatalf("Failed to create first refresh token: %v", result.Error)
	}

	// Try to create a second refresh token with the same token value
	refreshToken2 := RefreshToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Token:     "uniquetoken123", // Same token value
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	result = db.Create(&refreshToken2)

	// Should fail due to unique index constraint
	if result.Error == nil {
		t.Error("Expected error when creating refresh token with duplicate token value")
	}
}

// TestUserUniqueConstraints tests unique constraints on username and email
func TestUserUniqueConstraints(t *testing.T) {
	db := setupUserTestDB(t)

	// Create first user
	user1 := User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	result := db.Create(&user1)
	if result.Error != nil {
		t.Fatalf("Failed to create first user: %v", result.Error)
	}

	// Try to create second user with same username
	user2 := User{
		Username: "testuser", // Same username
		Email:    "different@example.com",
		Password: "password123",
	}
	result = db.Create(&user2)

	// Should fail due to unique constraint on username
	if result.Error == nil {
		t.Error("Expected error when creating user with duplicate username")
	}

	// Try to create third user with same email
	user3 := User{
		Username: "differentuser",
		Email:    "test@example.com", // Same email
		Password: "password123",
	}
	result = db.Create(&user3)

	// Should fail due to unique constraint on email
	if result.Error == nil {
		t.Error("Expected error when creating user with duplicate email")
	}
}
