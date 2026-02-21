package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"kanban-backend/models"
	"kanban-backend/utils"
)

type mockUserRepository struct {
	users map[string]*models.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[string]*models.User),
	}
}

func (m *mockUserRepository) Create(ctx context.Context, user *models.User) error {
	if _, exists := m.users[user.Email]; exists {
		return errors.New("user already exists")
	}

	if user.ID == "" {
		user.ID = generateTestID()
	}

	if user.Password != "" && len(user.Password) != 60 {
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}

	m.users[user.Email] = user
	return nil
}

func generateTestID() string {
	return "test-" + time.Now().Format("20060102150405")
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	user, exists := m.users[email]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *mockUserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *mockUserRepository) Update(ctx context.Context, user *models.User) error {
	if _, exists := m.users[user.Email]; !exists {
		return errors.New("user not found")
	}
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepository) Delete(ctx context.Context, id string) error {
	for email, user := range m.users {
		if user.ID == id {
			delete(m.users, email)
			return nil
		}
	}
	return errors.New("user not found")
}

func (m *mockUserRepository) SoftDelete(ctx context.Context, id string) error {
	return m.Delete(ctx, id)
}

func TestNewAuthService(t *testing.T) {
	mockRepo := newMockUserRepository()
	service := NewAuthService(mockRepo)

	if service == nil {
		t.Error("NewAuthService() should return non-nil service")
	}
}

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name        string
		username    string
		email       string
		password    string
		expectError bool
	}{
		{
			name:        "Valid registration",
			username:    "testuser",
			email:       "test@example.com",
			password:    "password123",
			expectError: false,
		},
		{
			name:        "Valid registration with special chars",
			username:    "user123",
			email:       "user123@test.com",
			password:    "P@ssw0rd!",
			expectError: false,
		},
		{
			name:        "Valid registration with long password",
			username:    "longuser",
			email:       "long@example.com",
			password:    "veryLongPasswordWithMixedChars123!@#",
			expectError: false,
		},
		{
			name:        "Registration with short password",
			username:    "shortuser",
			email:       "short@example.com",
			password:    "short",
			expectError: true,
		},
		{
			name:        "Registration with empty password",
			username:    "emptypass",
			email:       "empty@example.com",
			password:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockUserRepository()
			service := NewAuthService(mockRepo)

			user, err := service.Register(context.Background(), tt.username, tt.email, tt.password)

			if tt.expectError {
				if err == nil {
					t.Errorf("Register() expected error but got nil")
					return
				}
				var validationErr utils.ErrValidation
				if !errors.As(err, &validationErr) {
					t.Errorf("Register() should return ErrValidation, got %v", err)
				}
				return
			}

			if err != nil {
				t.Errorf("Register() unexpected error = %v", err)
				return
			}

			if user == nil {
				t.Error("Register() should return non-nil user")
				return
			}

			if user.Username != tt.username {
				t.Errorf("Register() username = %v, want %v", user.Username, tt.username)
			}

			if user.Email != tt.email {
				t.Errorf("Register() email = %v, want %v", user.Email, tt.email)
			}

			if user.ID == "" {
				t.Error("Register() should generate user ID")
			}

			if user.Password == tt.password {
				t.Error("Register() should hash password")
			}

			if len(user.Password) != 60 {
				t.Errorf("Register() password hash should be 60 chars, got %d", len(user.Password))
			}
		})
	}
}

func TestAuthService_RegisterDuplicateEmail(t *testing.T) {
	mockRepo := newMockUserRepository()
	service := NewAuthService(mockRepo)

	ctx := context.Background()
	username := "testuser"
	email := "test@example.com"
	password := "password123"

	_, err := service.Register(ctx, username, email, password)
	if err != nil {
		t.Fatalf("First registration failed: %v", err)
	}

	_, err = service.Register(ctx, username+"2", email, password)
	if err == nil {
		t.Error("Register() should error for duplicate email")
	}

	var conflictErr utils.ErrConflict
	if !errors.As(err, &conflictErr) {
		t.Errorf("Register() should return ErrConflict, got %v", err)
	}
}

func TestAuthService_Login(t *testing.T) {
	tests := []struct {
		name        string
		setupUser   bool
		setupEmail  string
		setupPass   string
		loginEmail  string
		loginPass   string
		expectError bool
	}{
		{
			name:        "Valid login",
			setupUser:   true,
			setupEmail:  "test@example.com",
			setupPass:   "password123",
			loginEmail:  "test@example.com",
			loginPass:   "password123",
			expectError: false,
		},
		{
			name:        "Login with wrong password",
			setupUser:   true,
			setupEmail:  "test@example.com",
			setupPass:   "password123",
			loginEmail:  "test@example.com",
			loginPass:   "wrongpassword",
			expectError: true,
		},
		{
			name:        "Login with non-existent user",
			setupUser:   false,
			loginEmail:  "nonexistent@example.com",
			loginPass:   "password123",
			expectError: true,
		},
		{
			name:        "Login with empty password",
			setupUser:   true,
			setupEmail:  "test@example.com",
			setupPass:   "password123",
			loginEmail:  "test@example.com",
			loginPass:   "",
			expectError: true,
		},
		{
			name:        "Login with empty email",
			setupUser:   true,
			setupEmail:  "test@example.com",
			setupPass:   "password123",
			loginEmail:  "",
			loginPass:   "password123",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockUserRepository()
			service := NewAuthService(mockRepo)
			ctx := context.Background()

			if tt.setupUser {
				_, err := service.Register(ctx, "testuser", tt.setupEmail, tt.setupPass)
				if err != nil {
					t.Fatalf("Setup registration failed: %v", err)
				}
			}

			token, err := service.Login(ctx, tt.loginEmail, tt.loginPass)

			if tt.expectError {
				if err == nil {
					t.Error("Login() expected error but got nil")
					return
				}
				var unauthorizedErr utils.ErrUnauthorized
				if !errors.As(err, &unauthorizedErr) {
					t.Errorf("Login() should return ErrUnauthorized, got %v", err)
				}
				return
			}

			if err != nil {
				t.Errorf("Login() unexpected error = %v", err)
				return
			}

			if token == "" {
				t.Error("Login() should return non-empty token")
			}
		})
	}
}

func TestAuthService_GenerateToken(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		expiry      time.Duration
		expectError bool
	}{
		{
			name:        "Valid token generation",
			userID:      "user123",
			expiry:      24 * time.Hour,
			expectError: false,
		},
		{
			name:        "Token with short expiry",
			userID:      "user456",
			expiry:      time.Minute,
			expectError: false,
		},
		{
			name:        "Token with long expiry",
			userID:      "user789",
			expiry:      30 * 24 * time.Hour,
			expectError: false,
		},
		{
			name:        "Empty userID should error",
			userID:      "",
			expiry:      24 * time.Hour,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockUserRepository()
			service := NewAuthService(mockRepo)

			token, err := service.GenerateToken(tt.userID, tt.expiry)

			if tt.expectError {
				if err == nil {
					t.Error("GenerateToken() expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("GenerateToken() unexpected error = %v", err)
				return
			}

			if token == "" {
				t.Error("GenerateToken() should return non-empty token")
			}

			userID, err := service.ValidateToken(token)
			if err != nil {
				t.Errorf("Generated token should be valid: %v", err)
			}

			if userID != tt.userID {
				t.Errorf("ValidateToken() userID = %v, want %v", userID, tt.userID)
			}
		})
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		expectError bool
		expectUser  string
	}{
		{
			name:        "Valid token",
			token:       generateTestToken("user123"),
			expectError: false,
			expectUser:  "user123",
		},
		{
			name:        "Empty token",
			token:       "",
			expectError: true,
		},
		{
			name:        "Invalid token format",
			token:       "invalid.token.string",
			expectError: true,
		},
		{
			name:        "Malformed token",
			token:       "notavalidtoken",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockUserRepository()
			service := NewAuthService(mockRepo)

			userID, err := service.ValidateToken(tt.token)

			if tt.expectError {
				if err == nil {
					t.Error("ValidateToken() expected error but got nil")
				}
				var unauthorizedErr utils.ErrUnauthorized
				if !errors.As(err, &unauthorizedErr) {
					t.Errorf("ValidateToken() should return ErrUnauthorized, got %v", err)
				}
				return
			}

			if err != nil {
				t.Errorf("ValidateToken() unexpected error = %v", err)
				return
			}

			if userID != tt.expectUser {
				t.Errorf("ValidateToken() userID = %v, want %v", userID, tt.expectUser)
			}
		})
	}
}

func TestAuthService_HashPassword(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		expectError bool
	}{
		{
			name:        "Valid password",
			password:    "password123",
			expectError: false,
		},
		{
			name:        "Password with special chars",
			password:    "P@ssw0rd!",
			expectError: false,
		},
		{
			name:        "Long password",
			password:    "veryLongPasswordWithMixedChars123!@#",
			expectError: false,
		},
		{
			name:        "Short password",
			password:    "short",
			expectError: true,
		},
		{
			name:        "Empty password",
			password:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockUserRepository()
			service := NewAuthService(mockRepo)

			hashed, err := service.HashPassword(tt.password)

			if tt.expectError {
				if err == nil {
					t.Error("HashPassword() expected error but got nil")
				}
				if !errors.Is(err, utils.ErrPasswordTooShort) {
					t.Errorf("HashPassword() should return ErrPasswordTooShort, got %v", err)
				}
				return
			}

			if err != nil {
				t.Errorf("HashPassword() unexpected error = %v", err)
				return
			}

			if hashed == tt.password {
				t.Error("HashPassword() should hash password")
			}

			if len(hashed) != 60 {
				t.Errorf("HashPassword() hash should be 60 chars, got %d", len(hashed))
			}
		})
	}
}

func TestAuthService_VerifyPassword(t *testing.T) {
	password := "password123"
	hashed, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to create test hash: %v", err)
	}

	tests := []struct {
		name           string
		hashedPassword string
		password       string
		expectError    bool
	}{
		{
			name:           "Correct password",
			hashedPassword: hashed,
			password:       password,
			expectError:    false,
		},
		{
			name:           "Wrong password",
			hashedPassword: hashed,
			password:       "wrongpassword",
			expectError:    true,
		},
		{
			name:           "Empty password",
			hashedPassword: hashed,
			password:       "",
			expectError:    true,
		},
		{
			name:           "Invalid hash",
			hashedPassword: "invalidhash",
			password:       password,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockUserRepository()
			service := NewAuthService(mockRepo)

			err := service.VerifyPassword(tt.hashedPassword, tt.password)

			if tt.expectError {
				if err == nil {
					t.Error("VerifyPassword() expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("VerifyPassword() unexpected error = %v", err)
			}
		})
	}
}

func TestAuthService_Integration(t *testing.T) {
	mockRepo := newMockUserRepository()
	service := NewAuthService(mockRepo)
	ctx := context.Background()

	username := "testuser"
	email := "test@example.com"
	password := "password123"

	user, err := service.Register(ctx, username, email, password)
	if err != nil {
		t.Fatalf("Register() failed: %v", err)
	}

	if user == nil || user.ID == "" {
		t.Fatal("Register() should return user with ID")
	}

	token, err := service.Login(ctx, email, password)
	if err != nil {
		t.Fatalf("Login() failed: %v", err)
	}

	if token == "" {
		t.Fatal("Login() should return token")
	}

	userID, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() failed: %v", err)
	}

	if userID != user.ID {
		t.Errorf("ValidateToken() userID = %v, want %v", userID, user.ID)
	}

	err = service.VerifyPassword(user.Password, password)
	if err != nil {
		t.Errorf("VerifyPassword() should verify correct password: %v", err)
	}

	err = service.VerifyPassword(user.Password, "wrongpassword")
	if err == nil {
		t.Error("VerifyPassword() should reject wrong password")
	}
}

func generateTestToken(userID string) string {
	token, _ := utils.GenerateToken(userID, 24*time.Hour)
	return token
}
