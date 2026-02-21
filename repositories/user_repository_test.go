package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"kanban-backend/models"
)

func setupUserRepositoryTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.User{}, &models.RefreshToken{})
	require.NoError(t, err)

	type AuditLog struct {
		ID        string    `gorm:"primaryKey;type:varchar(36)"`
		UserID    string    `gorm:"not null;type:varchar(36);index"`
		Action    string    `gorm:"not null;type:varchar(50)"`
		Message   string    `gorm:"type:text"`
		CreatedAt time.Time `gorm:"autoCreateTime"`
	}
	err = db.Table("audit_logs").AutoMigrate(&AuditLog{})
	require.NoError(t, err)

	return db
}

func TestUserRepository_Create(t *testing.T) {
	db := setupUserRepositoryTestDB(t)
	repo := &userRepository{db: db}

	testCases := []struct {
		name      string
		user      *models.User
		wantError bool
	}{
		{
			name: "Valid user",
			user: &models.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantError: false,
		},
		{
			name: "Duplicate email",
			user: &models.User{
				Username: "testuser2",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			err := repo.Create(ctx, tc.user)

			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tc.user.ID)
				assert.Equal(t, "testuser", tc.user.Username)
				assert.Equal(t, "test@example.com", tc.user.Email)
			}
		})
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db := setupUserRepositoryTestDB(t)
	repo := &userRepository{db: db}

	ctx := context.Background()

	testUser := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	err := repo.Create(ctx, testUser)
	require.NoError(t, err)

	testCases := []struct {
		name      string
		email     string
		wantError bool
	}{
		{
			name:      "Found user",
			email:     "test@example.com",
			wantError: false,
		},
		{
			name:      "User not found",
			email:     "notfound@example.com",
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user, err := repo.FindByEmail(ctx, tc.email)

			if tc.wantError {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, "testuser", user.Username)
				assert.Equal(t, "test@example.com", user.Email)
			}
		})
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	db := setupUserRepositoryTestDB(t)
	repo := &userRepository{db: db}

	ctx := context.Background()

	testUser := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	err := repo.Create(ctx, testUser)
	require.NoError(t, err)

	testCases := []struct {
		name      string
		id        string
		wantError bool
	}{
		{
			name:      "Found user",
			id:        testUser.ID,
			wantError: false,
		},
		{
			name:      "User not found",
			id:        uuid.New().String(),
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user, err := repo.FindByID(ctx, tc.id)

			if tc.wantError {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, testUser.ID, user.ID)
				assert.Equal(t, "testuser", user.Username)
			}
		})
	}
}

func TestUserRepository_Update(t *testing.T) {
	db := setupUserRepositoryTestDB(t)
	repo := &userRepository{db: db}

	ctx := context.Background()

	testUser := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	err := repo.Create(ctx, testUser)
	require.NoError(t, err)

	time.Sleep(10 * time.Millisecond)
	testUser.Username = "updateduser"
	err = repo.Update(ctx, testUser)
	assert.NoError(t, err)

	updatedUser, err := repo.FindByID(ctx, testUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, "updateduser", updatedUser.Username)
	assert.True(t, updatedUser.UpdatedAt.After(testUser.CreatedAt))
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupUserRepositoryTestDB(t)
	repo := &userRepository{db: db}

	ctx := context.Background()

	testUser := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	err := repo.Create(ctx, testUser)
	require.NoError(t, err)

	err = repo.Delete(ctx, testUser.ID)
	assert.NoError(t, err)

	_, err = repo.FindByID(ctx, testUser.ID)
	assert.Error(t, err)

	testCases := []struct {
		name      string
		id        string
		wantError bool
	}{
		{
			name:      "Delete non-existent user",
			id:        uuid.New().String(),
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.Delete(ctx, tc.id)
			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserRepository_SoftDelete(t *testing.T) {
	db := setupUserRepositoryTestDB(t)
	repo := &userRepository{db: db}

	ctx := context.Background()

	testUser := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	err := repo.Create(ctx, testUser)
	require.NoError(t, err)

	err = repo.SoftDelete(ctx, testUser.ID)
	assert.NoError(t, err)

	_, err = repo.FindByID(ctx, testUser.ID)
	assert.Error(t, err)

	var deletedUser models.User
	err = db.Unscoped().Where("id = ?", testUser.ID).First(&deletedUser).Error
	assert.NoError(t, err)
	assert.NotNil(t, deletedUser.DeletedAt)

	testCases := []struct {
		name      string
		id        string
		wantError bool
	}{
		{
			name:      "Soft delete non-existent user",
			id:        uuid.New().String(),
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.SoftDelete(ctx, tc.id)
			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserRepository_Context(t *testing.T) {
	db := setupUserRepositoryTestDB(t)
	repo := &userRepository{db: db}

	ctx, cancel := context.WithCancel(context.Background())

	testUser := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	err := repo.Create(ctx, testUser)
	assert.NoError(t, err)

	user, err := repo.FindByEmail(ctx, "test@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, user)

	cancel()

	_, err = repo.FindByEmail(ctx, "test@example.com")
	assert.Error(t, err)
}

func TestUserRepository_PasswordHashing(t *testing.T) {
	db := setupUserRepositoryTestDB(t)
	repo := &userRepository{db: db}

	ctx := context.Background()

	testUser := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "plainpassword",
	}

	err := repo.Create(ctx, testUser)
	assert.NoError(t, err)

	assert.NotEqual(t, "plainpassword", testUser.Password)
	assert.Equal(t, 60, len(testUser.Password))
}
