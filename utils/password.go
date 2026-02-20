package utils

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordTooShort = errors.New("password must be at least 8 characters long")
)

// HashPassword hashes a password using bcrypt with DefaultCost (10 rounds).
// Validates that the password is at least 8 characters long before hashing.
// Returns the hashed password as a string or an error if validation fails.
func HashPassword(password string) (string, error) {
	if len(password) < 8 {
		return "", ErrPasswordTooShort
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// CheckPassword verifies that a plain text password matches a hashed password.
// Returns nil if the password is correct, or an error if it doesn't match.
func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
