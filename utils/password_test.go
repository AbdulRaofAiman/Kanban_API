package utils

import (
	"errors"
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		expectError bool
		errorIs     error
	}{
		{
			name:        "Valid password with 8 characters",
			password:    "password",
			expectError: false,
		},
		{
			name:        "Valid password with more than 8 characters",
			password:    "mySecurePassword123!",
			expectError: false,
		},
		{
			name:        "Password with exactly 8 characters",
			password:    "12345678",
			expectError: false,
		},
		{
			name:        "Password with less than 8 characters should fail",
			password:    "short",
			expectError: true,
			errorIs:     ErrPasswordTooShort,
		},
		{
			name:        "Empty password should fail",
			password:    "",
			expectError: true,
			errorIs:     ErrPasswordTooShort,
		},
		{
			name:        "Password with exactly 7 characters should fail",
			password:    "1234567",
			expectError: true,
			errorIs:     ErrPasswordTooShort,
		},
		{
			name:        "Password with special characters",
			password:    "P@ssw0rd!",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashed, err := HashPassword(tt.password)

			if tt.expectError {
				if err == nil {
					t.Errorf("HashPassword() expected error but got nil")
					return
				}
				if tt.errorIs != nil && !errors.Is(err, tt.errorIs) {
					t.Errorf("HashPassword() error = %v, want %v", err, tt.errorIs)
				}
				return
			}

			if err != nil {
				t.Errorf("HashPassword() unexpected error = %v", err)
				return
			}

			// Verify that the hashed password is different from the original
			if hashed == tt.password {
				t.Errorf("HashPassword() hashed password should differ from original password")
			}

			// Verify that the hashed password is not empty
			if hashed == "" {
				t.Errorf("HashPassword() returned empty hash")
			}

			// Verify that hashing the same password twice produces different hashes
			// (due to bcrypt's salt)
			hashed2, err := HashPassword(tt.password)
			if err != nil {
				t.Errorf("HashPassword() second call unexpected error = %v", err)
				return
			}

			if hashed == hashed2 {
				t.Errorf("HashPassword() same password should produce different hashes due to salt")
			}

			// Verify both hashes can be checked against the original password
			if err := CheckPassword(hashed, tt.password); err != nil {
				t.Errorf("CheckPassword() first hash failed = %v", err)
			}
			if err := CheckPassword(hashed2, tt.password); err != nil {
				t.Errorf("CheckPassword() second hash failed = %v", err)
			}
		})
	}
}

func TestCheckPassword(t *testing.T) {
	// Create a test hash for a known password
	testPassword := "testPassword123"
	hashed, err := HashPassword(testPassword)
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
			name:           "Correct password should pass",
			hashedPassword: hashed,
			password:       testPassword,
			expectError:    false,
		},
		{
			name:           "Incorrect password should fail",
			hashedPassword: hashed,
			password:       "wrongPassword",
			expectError:    true,
		},
		{
			name:           "Empty password should fail",
			hashedPassword: hashed,
			password:       "",
			expectError:    true,
		},
		{
			name:           "Password with different case should fail",
			hashedPassword: hashed,
			password:       strings.ToUpper(testPassword),
			expectError:    true,
		},
		{
			name:           "Password with trailing space should fail",
			hashedPassword: hashed,
			password:       testPassword + " ",
			expectError:    true,
		},
		{
			name:           "Password with leading space should fail",
			hashedPassword: hashed,
			password:       " " + testPassword,
			expectError:    true,
		},
		{
			name:           "Invalid hash format should fail",
			hashedPassword: "invalidHash",
			password:       testPassword,
			expectError:    true,
		},
		{
			name:           "Empty hash should fail",
			hashedPassword: "",
			password:       testPassword,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPassword(tt.hashedPassword, tt.password)

			if tt.expectError {
				if err == nil {
					t.Errorf("CheckPassword() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("CheckPassword() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestPasswordIntegration(t *testing.T) {
	// Test the full workflow: hash, store, verify
	originalPassword := "MySecurePassword123!"

	// Hash the password
	hashed, err := HashPassword(originalPassword)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Verify the hashed password is different from original
	if hashed == originalPassword {
		t.Errorf("Hashed password should differ from original")
	}

	// Verify correct password works
	if err := CheckPassword(hashed, originalPassword); err != nil {
		t.Errorf("Failed to verify correct password: %v", err)
	}

	// Verify incorrect password fails
	wrongPassword := "WrongPassword123!"
	if err := CheckPassword(hashed, wrongPassword); err == nil {
		t.Errorf("CheckPassword should fail for incorrect password")
	}

	// Verify hash length is consistent (bcrypt hashes are 60 characters)
	if len(hashed) != 60 {
		t.Errorf("Bcrypt hash should be 60 characters, got %d", len(hashed))
	}

	// Verify hash starts with bcrypt prefix
	if !strings.HasPrefix(hashed, "$2a$") && !strings.HasPrefix(hashed, "$2b$") && !strings.HasPrefix(hashed, "$2y$") {
		t.Errorf("Hash should have bcrypt prefix, got: %s", hashed[:5])
	}
}

func TestHashPasswordErrorCases(t *testing.T) {
	t.Run("Nil or empty password returns ErrPasswordTooShort", func(t *testing.T) {
		_, err := HashPassword("")
		if err == nil {
			t.Error("Expected ErrPasswordTooShort for empty password")
		}
		if !errors.Is(err, ErrPasswordTooShort) {
			t.Errorf("Expected ErrPasswordTooShort, got: %v", err)
		}
	})
}

func TestCheckPasswordErrorCases(t *testing.T) {
	t.Run("Invalid bcrypt hash format", func(t *testing.T) {
		err := CheckPassword("not-a-valid-hash", "password")
		if err == nil {
			t.Error("Expected error for invalid hash format")
		}
	})

	t.Run("Hash with wrong version", func(t *testing.T) {
		err := CheckPassword("$1$salt$hash", "password")
		if err == nil {
			t.Error("Expected error for wrong bcrypt version")
		}
	})

	t.Run("Corrupted hash", func(t *testing.T) {
		err := CheckPassword("$2a$10$invalidhashlength", "password")
		if err == nil {
			t.Error("Expected error for corrupted hash")
		}
	})
}
