package utils

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func setupTestJWTSecret() {
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")
	jwtSecret = []byte("test-secret-key-for-testing")
}

func teardownTestJWTSecret() {
	os.Unsetenv("JWT_SECRET")
	jwtSecret = []byte{}
}

func TestGenerateToken(t *testing.T) {
	setupTestJWTSecret()
	defer teardownTestJWTSecret()

	tests := []struct {
		name    string
		userID  string
		expiry  time.Duration
		wantErr bool
	}{
		{
			name:    "Valid token generation",
			userID:  "user123",
			expiry:  1 * time.Hour,
			wantErr: false,
		},
		{
			name:    "Empty user ID",
			userID:  "",
			expiry:  1 * time.Hour,
			wantErr: true,
		},
		{
			name:    "Short expiry",
			userID:  "user456",
			expiry:  5 * time.Minute,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.userID, tt.expiry)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && token == "" {
				t.Error("GenerateToken() returned empty token")
			}

			if !tt.wantErr {
				// Verify token can be parsed
				parsedToken, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte(os.Getenv("JWT_SECRET")), nil
				})
				if err != nil {
					t.Errorf("Failed to parse generated token: %v", err)
				}

				claims, ok := parsedToken.Claims.(*CustomClaims)
				if !ok {
					t.Error("Failed to extract claims from token")
				}

				if claims.UserID != tt.userID {
					t.Errorf("Expected user_id %s, got %s", tt.userID, claims.UserID)
				}
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	setupTestJWTSecret()
	defer teardownTestJWTSecret()

	// Generate a valid token
	validToken, err := GenerateToken("user123", 1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to generate test token: %v", err)
	}

	tests := []struct {
		name       string
		token      string
		wantErr    bool
		errorMsg   string
		wantUserID string
	}{
		{
			name:       "Valid token",
			token:      validToken,
			wantErr:    false,
			wantUserID: "user123",
		},
		{
			name:     "Empty token",
			token:    "",
			wantErr:  true,
			errorMsg: "token cannot be empty",
		},
		{
			name:     "Invalid token format",
			token:    "invalid.token.format",
			wantErr:  true,
			errorMsg: "failed to parse token",
		},
		{
			name:     "Malformed token",
			token:    "notavalidtoken",
			wantErr:  true,
			errorMsg: "failed to parse token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorMsg != "" && err != nil {
				// Check error message contains expected string
				if !contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to contain '%s', got '%s'", tt.errorMsg, err.Error())
				}
			}

			if !tt.wantErr && claims.UserID != tt.wantUserID {
				t.Errorf("Expected user_id %s, got %s", tt.wantUserID, claims.UserID)
			}
		})
	}
}

func TestValidateExpiredToken(t *testing.T) {
	setupTestJWTSecret()
	defer teardownTestJWTSecret()

	// Create an expired token by using negative duration
	expiredToken, err := GenerateToken("user123", -1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to generate expired token: %v", err)
	}

	claims, err := ValidateToken(expiredToken)
	if err == nil {
		t.Error("Expected error for expired token, got nil")
	}

	if claims != nil {
		t.Error("Expected nil claims for expired token")
	}

	if err != nil && !contains(err.Error(), "expired") {
		t.Errorf("Expected expired error, got: %v", err)
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	setupTestJWTSecret()
	defer teardownTestJWTSecret()

	tests := []struct {
		name    string
		userID  string
		wantErr bool
	}{
		{
			name:    "Valid refresh token generation",
			userID:  "user123",
			wantErr: false,
		},
		{
			name:    "Empty user ID",
			userID:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateRefreshToken(tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && token == "" {
				t.Error("GenerateRefreshToken() returned empty token")
			}

			if !tt.wantErr {
				// Verify token can be parsed
				parsedToken, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte(os.Getenv("JWT_SECRET")), nil
				})
				if err != nil {
					t.Errorf("Failed to parse generated refresh token: %v", err)
				}

				claims, ok := parsedToken.Claims.(*CustomClaims)
				if !ok {
					t.Error("Failed to extract claims from refresh token")
				}

				if claims.UserID != tt.userID {
					t.Errorf("Expected user_id %s, got %s", tt.userID, claims.UserID)
				}

				// Verify expiry is approximately 7 days
				expectedExpiry := time.Now().Add(7 * 24 * time.Hour)
				timeDiff := expectedExpiry.Sub(claims.ExpiresAt.Time)
				if timeDiff > 1*time.Second || timeDiff < -1*time.Second {
					t.Errorf("Refresh token expiry mismatch. Expected ~7 days, got %v", claims.ExpiresAt.Time)
				}
			}
		})
	}
}

func TestValidateRefreshToken(t *testing.T) {
	setupTestJWTSecret()
	defer teardownTestJWTSecret()

	// Generate a valid refresh token
	validToken, err := GenerateRefreshToken("user123")
	if err != nil {
		t.Fatalf("Failed to generate test refresh token: %v", err)
	}

	tests := []struct {
		name       string
		token      string
		wantErr    bool
		errorMsg   string
		wantUserID string
	}{
		{
			name:       "Valid refresh token",
			token:      validToken,
			wantErr:    false,
			wantUserID: "user123",
		},
		{
			name:     "Empty refresh token",
			token:    "",
			wantErr:  true,
			errorMsg: "refresh token cannot be empty",
		},
		{
			name:     "Invalid refresh token format",
			token:    "invalid.token.format",
			wantErr:  true,
			errorMsg: "failed to parse refresh token",
		},
		{
			name:     "Malformed refresh token",
			token:    "notavalidtoken",
			wantErr:  true,
			errorMsg: "failed to parse refresh token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := ValidateRefreshToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorMsg != "" && err != nil {
				if !contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to contain '%s', got '%s'", tt.errorMsg, err.Error())
				}
			}

			if !tt.wantErr && userID != tt.wantUserID {
				t.Errorf("Expected user_id %s, got %s", tt.wantUserID, userID)
			}
		})
	}
}

func TestValidateExpiredRefreshToken(t *testing.T) {
	setupTestJWTSecret()
	defer teardownTestJWTSecret()

	// Manually create an expired refresh token
	expiryTime := time.Now().Add(-1 * time.Hour)
	claims := CustomClaims{
		UserID: "user123",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiryTime),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-8 * 24 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-8 * 24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiredToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		t.Fatalf("Failed to create expired refresh token: %v", err)
	}

	userID, err := ValidateRefreshToken(expiredToken)
	if err == nil {
		t.Error("Expected error for expired refresh token, got nil")
	}

	if userID != "" {
		t.Error("Expected empty user_id for expired refresh token")
	}

	if err != nil && !contains(err.Error(), "expired") {
		t.Errorf("Expected expired error, got: %v", err)
	}
}

func TestTokenWithDifferentSecret(t *testing.T) {
	setupTestJWTSecret()
	defer teardownTestJWTSecret()

	// Generate token with one secret
	token, err := GenerateToken("user123", 1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Try to validate with different secret (simulating environment change)
	jwtSecret = []byte("different-secret-key")

	_, err = ValidateToken(token)
	if err == nil {
		t.Error("Expected error when validating token with different secret")
	}
}

func TestAccessTokenExpiryOneHour(t *testing.T) {
	setupTestJWTSecret()
	defer teardownTestJWTSecret()

	// Generate access token with 1 hour expiry
	token, err := GenerateToken("user123", 1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Parse and verify expiry
	parsedToken, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	claims, ok := parsedToken.Claims.(*CustomClaims)
	if !ok {
		t.Fatal("Failed to extract claims")
	}

	// Verify expiry is approximately 1 hour
	expectedExpiry := time.Now().Add(1 * time.Hour)
	timeDiff := expectedExpiry.Sub(claims.ExpiresAt.Time)
	if timeDiff > 1*time.Second || timeDiff < -1*time.Second {
		t.Errorf("Access token expiry mismatch. Expected ~1 hour, got %v", claims.ExpiresAt.Time)
	}
}

func TestRefreshTokenExpirySevenDays(t *testing.T) {
	setupTestJWTSecret()
	defer teardownTestJWTSecret()

	// Generate refresh token
	token, err := GenerateRefreshToken("user123")
	if err != nil {
		t.Fatalf("Failed to generate refresh token: %v", err)
	}

	// Parse and verify expiry
	parsedToken, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		t.Fatalf("Failed to parse refresh token: %v", err)
	}

	claims, ok := parsedToken.Claims.(*CustomClaims)
	if !ok {
		t.Fatal("Failed to extract claims")
	}

	// Verify expiry is approximately 7 days
	expectedExpiry := time.Now().Add(7 * 24 * time.Hour)
	timeDiff := expectedExpiry.Sub(claims.ExpiresAt.Time)
	if timeDiff > 1*time.Second || timeDiff < -1*time.Second {
		t.Errorf("Refresh token expiry mismatch. Expected ~7 days, got %v", claims.ExpiresAt.Time)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (s[:len(substr)] == substr || contains(s[1:], substr)))
}
