package utils

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// Custom error types
type ErrNotFound struct {
	Message string
}

func (e ErrNotFound) Error() string {
	if e.Message == "" {
		return "resource not found"
	}
	return e.Message
}

type ErrUnauthorized struct {
	Message string
}

func (e ErrUnauthorized) Error() string {
	if e.Message == "" {
		return "unauthorized access"
	}
	return e.Message
}

type ErrValidation struct {
	Message string
}

func (e ErrValidation) Error() string {
	if e.Message == "" {
		return "validation failed"
	}
	return e.Message
}

type ErrConflict struct {
	Message string
}

func (e ErrConflict) Error() string {
	if e.Message == "" {
		return "resource conflict"
	}
	return e.Message
}

// Helper functions to create errors with custom messages
func NewNotFound(msg string) error {
	return ErrNotFound{Message: msg}
}

func NewUnauthorized(msg string) error {
	return ErrUnauthorized{Message: msg}
}

func NewValidation(msg string) error {
	return ErrValidation{Message: msg}
}

func NewConflict(msg string) error {
	return ErrConflict{Message: msg}
}

// ErrorHandler is a Fiber error handler that maps custom errors to appropriate HTTP status codes
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Check if error is one of our custom error types
	switch e := err.(type) {
	case ErrNotFound:
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   e.Error(),
		})
	case ErrUnauthorized:
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   e.Error(),
		})
	case ErrValidation:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   e.Error(),
		})
	case ErrConflict:
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"error":   e.Error(),
		})
	default:
		// For all other errors, return 500 Internal Server Error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   fmt.Sprintf("internal server error: %s", err.Error()),
		})
	}
}
