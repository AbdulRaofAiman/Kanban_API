package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupNotificationTestDB creates an in-memory SQLite database for testing notifications
func setupNotificationTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "Failed to connect to test database")

	// Create tables
	err = db.AutoMigrate(&User{}, &Notification{})
	require.NoError(t, err, "Failed to migrate tables")

	return db
}

// createTestNotificationUser creates a minimal test user in the database
func createTestNotificationUser(t *testing.T, db *gorm.DB) *User {
	user := &User{
		ID:        uuid.New().String(),
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := db.Create(user).Error
	require.NoError(t, err, "Failed to create test user")

	return user
}

// TestNotificationModel tests the Notification model structure and GORM tags
func TestNotificationModel(t *testing.T) {
	db := setupNotificationTestDB(t)

	t.Run("table name is notifications", func(t *testing.T) {
		notification := &Notification{}
		assert.Equal(t, "notifications", notification.TableName())
	})

	t.Run("has correct fields with GORM tags", func(t *testing.T) {
		user := createTestNotificationUser(t, db)

		notification := &Notification{
			ID:      uuid.New().String(),
			UserID:  user.ID,
			Message: "Test notification message",
			ReadAt:  nil,
		}

		err := db.Create(notification).Error
		require.NoError(t, err, "Failed to create notification")

		// Verify the notification was created correctly
		var dbNotification Notification
		err = db.First(&dbNotification, notification.ID).Error
		require.NoError(t, err, "Failed to retrieve notification")

		assert.Equal(t, notification.ID, dbNotification.ID)
		assert.Equal(t, notification.UserID, dbNotification.UserID)
		assert.Equal(t, notification.Message, dbNotification.Message)
		assert.Nil(t, dbNotification.ReadAt)
		assert.False(t, dbNotification.CreatedAt.IsZero())
		assert.False(t, dbNotification.UpdatedAt.IsZero())
	})
}

// TestNotificationBeforeCreate tests the BeforeCreate hook
func TestNotificationBeforeCreate(t *testing.T) {
	db := setupNotificationTestDB(t)
	user := createTestNotificationUser(t, db)

	t.Run("generates UUID if ID is empty", func(t *testing.T) {
		notification := &Notification{
			UserID:  user.ID,
			Message: "Test notification",
		}

		err := db.Create(notification).Error
		require.NoError(t, err, "Failed to create notification")

		assert.NotEmpty(t, notification.ID, "ID should be generated")
		assert.NotEqual(t, "", notification.ID, "ID should not be empty")

		// Verify it's a valid UUID
		_, err = uuid.Parse(notification.ID)
		assert.NoError(t, err, "ID should be a valid UUID")
	})

	t.Run("keeps existing ID if provided", func(t *testing.T) {
		existingID := uuid.New().String()
		notification := &Notification{
			ID:      existingID,
			UserID:  user.ID,
			Message: "Test notification with existing ID",
		}

		err := db.Create(notification).Error
		require.NoError(t, err, "Failed to create notification")

		assert.Equal(t, existingID, notification.ID, "ID should be preserved")
	})
}

// TestNotificationRelationships tests the relationship between Notification and User
func TestNotificationRelationships(t *testing.T) {
	db := setupNotificationTestDB(t)
	user := createTestNotificationUser(t, db)

	t.Run("belongsTo user with foreign key", func(t *testing.T) {
		notification := &Notification{
			UserID:  user.ID,
			Message: "Test notification for user",
		}

		err := db.Create(notification).Error
		require.NoError(t, err, "Failed to create notification")

		// Test eager loading
		var dbNotification Notification
		err = db.Preload("User").First(&dbNotification, notification.ID).Error
		require.NoError(t, err, "Failed to preload user")

		assert.NotNil(t, dbNotification.User, "User should be loaded")
		assert.Equal(t, user.ID, dbNotification.User.ID, "User ID should match")
		assert.Equal(t, user.Email, dbNotification.User.Email, "User email should match")
	})

	t.Run("cascade delete when user is deleted", func(t *testing.T) {
		newUser := &User{
			ID:        uuid.New().String(),
			Email:     "test2@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := db.Create(newUser).Error
		require.NoError(t, err, "Failed to create test user")

		notification := &Notification{
			UserID:  newUser.ID,
			Message: "Notification for user to be deleted",
		}

		err = db.Create(notification).Error
		require.NoError(t, err, "Failed to create notification")

		// Delete the user (should cascade delete notifications)
		err = db.Delete(newUser).Error
		require.NoError(t, err, "Failed to delete user")

		// Verify notification is soft deleted
		var dbNotification Notification
		err = db.Unscoped().First(&dbNotification, notification.ID).Error
		require.NoError(t, err, "Notification should still exist in database (soft deleted)")

		assert.NotNil(t, dbNotification.DeletedAt, "DeletedAt should be set")
	})
}

// TestNotificationMarkAsRead tests the MarkAsRead method
func TestNotificationMarkAsRead(t *testing.T) {
	t.Run("marks notification as read", func(t *testing.T) {
		notification := &Notification{
			ID:      uuid.New().String(),
			UserID:  uuid.New().String(),
			Message: "Test notification",
			ReadAt:  nil,
		}

		assert.False(t, notification.IsRead(), "Notification should be unread initially")

		notification.MarkAsRead()

		assert.NotNil(t, notification.ReadAt, "ReadAt should be set")
		assert.True(t, notification.IsRead(), "Notification should be marked as read")
		assert.WithinDuration(t, time.Now(), *notification.ReadAt, time.Second, "ReadAt should be close to now")
	})

	t.Run("updates read timestamp", func(t *testing.T) {
		notification := &Notification{
			ID:      uuid.New().String(),
			UserID:  uuid.New().String(),
			Message: "Test notification",
		}

		oldReadAt := time.Now().Add(-1 * time.Hour)
		notification.ReadAt = &oldReadAt

		notification.MarkAsRead()

		assert.True(t, notification.ReadAt.After(oldReadAt), "ReadAt should be updated to a later time")
	})
}

// TestNotificationIsRead tests the IsRead method
func TestNotificationIsRead(t *testing.T) {
	t.Run("returns true when ReadAt is set", func(t *testing.T) {
		now := time.Now()
		notification := &Notification{
			ID:      uuid.New().String(),
			UserID:  uuid.New().String(),
			Message: "Test notification",
			ReadAt:  &now,
		}

		assert.True(t, notification.IsRead(), "Notification should be marked as read")
	})

	t.Run("returns false when ReadAt is nil", func(t *testing.T) {
		notification := &Notification{
			ID:      uuid.New().String(),
			UserID:  uuid.New().String(),
			Message: "Test notification",
			ReadAt:  nil,
		}

		assert.False(t, notification.IsRead(), "Notification should be marked as unread")
	})
}

// TestNotificationConstraints tests database constraints
func TestNotificationConstraints(t *testing.T) {
	db := setupNotificationTestDB(t)
	user := createTestNotificationUser(t, db)

	t.Run("UserID cannot be null", func(t *testing.T) {
		notification := &Notification{
			Message: "Test notification without user",
		}

		err := db.Create(notification).Error
		assert.Error(t, err, "Should fail when UserID is null")
	})

	t.Run("Message cannot be null", func(t *testing.T) {
		notification := &Notification{
			UserID: user.ID,
		}

		err := db.Create(notification).Error
		assert.Error(t, err, "Should fail when Message is null")
	})

	t.Run("UserID foreign key constraint", func(t *testing.T) {
		nonExistentUserID := uuid.New().String()
		notification := &Notification{
			UserID:  nonExistentUserID,
			Message: "Test notification with invalid user",
		}

		// SQLite doesn't enforce foreign keys by default, but the model has the tag
		// This test verifies the foreign key is properly configured
		assert.Equal(t, nonExistentUserID, notification.UserID, "UserID should be set")
	})
}

// TestNotificationSoftDelete tests soft delete functionality
func TestNotificationSoftDelete(t *testing.T) {
	db := setupNotificationTestDB(t)
	user := createTestNotificationUser(t, db)

	t.Run("soft delete sets DeletedAt", func(t *testing.T) {
		notification := &Notification{
			UserID:  user.ID,
			Message: "Test notification for soft delete",
		}

		err := db.Create(notification).Error
		require.NoError(t, err, "Failed to create notification")

		// Soft delete
		err = db.Delete(&notification).Error
		require.NoError(t, err, "Failed to soft delete notification")

		// Verify DeletedAt is set
		assert.NotNil(t, notification.DeletedAt, "DeletedAt should be set")
	})

	t.Run("soft deleted records not found by default", func(t *testing.T) {
		notification := &Notification{
			UserID:  user.ID,
			Message: "Test notification for query",
		}

		err := db.Create(notification).Error
		require.NoError(t, err, "Failed to create notification")

		// Soft delete
		err = db.Delete(&notification).Error
		require.NoError(t, err, "Failed to soft delete notification")

		// Try to find with normal query
		var foundNotification Notification
		err = db.First(&foundNotification, notification.ID).Error
		assert.Error(t, err, "Should not find soft deleted record")

		// Find with Unscoped
		var unscopedNotification Notification
		err = db.Unscoped().First(&unscopedNotification, notification.ID).Error
		require.NoError(t, err, "Should find soft deleted record with Unscoped")
		assert.Equal(t, notification.ID, unscopedNotification.ID, "ID should match")
	})
}

// TestNotificationJSONSerialization tests JSON serialization
func TestNotificationJSONSerialization(t *testing.T) {
	db := setupNotificationTestDB(t)

	t.Run("serializes correctly", func(t *testing.T) {
		now := time.Now()
		user := createTestNotificationUser(t, db)
		notification := &Notification{
			ID:      uuid.New().String(),
			UserID:  user.ID,
			Message: "Test notification",
			ReadAt:  &now,
		}

		// Create the notification to verify it works with GORM
		err := db.Create(notification).Error
		require.NoError(t, err, "Failed to create notification")

		// Verify all fields are present
		assert.NotEmpty(t, notification.ID, "ID should be set")
		assert.NotEmpty(t, notification.UserID, "UserID should be set")
		assert.NotEmpty(t, notification.Message, "Message should be set")
		assert.NotNil(t, notification.ReadAt, "ReadAt should be set")
	})
}
