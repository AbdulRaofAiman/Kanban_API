package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// CustomClaims extends jwt.RegisteredClaims to include user_id
type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GetJWTSecret returns the JWT signing key
func GetJWTSecret() []byte {
	if len(jwtSecret) == 0 {
		jwtSecret = []byte("default-secret-key-change-in-production")
	}
	return jwtSecret
}

// GenerateToken creates a new JWT access token with specified expiry
func GenerateToken(userID string, expiry time.Duration) (string, error) {
	if userID == "" {
		return "", errors.New("user_id cannot be empty")
	}

	secret := GetJWTSecret()

	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates the access token and returns the claims
func ValidateToken(tokenString string) (*CustomClaims, error) {
	if tokenString == "" {
		return nil, errors.New("token cannot be empty")
	}

	secret := GetJWTSecret()

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}

// GenerateRefreshToken creates a new refresh token (7 days expiry)
func GenerateRefreshToken(userID string) (string, error) {
	if userID == "" {
		return "", errors.New("user_id cannot be empty")
	}

	// Refresh token expires in 7 days
	expiry := 7 * 24 * time.Hour

	secret := GetJWTSecret()

	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, nil
}

// ValidateRefreshToken validates the refresh token and returns the user_id
func ValidateRefreshToken(tokenString string) (string, error) {
	if tokenString == "" {
		return "", errors.New("refresh token cannot be empty")
	}

	secret := GetJWTSecret()

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse refresh token: %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid refresh token")
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		return "", errors.New("refresh token has expired")
	}

	return claims.UserID, nil
}
