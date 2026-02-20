package utils

import (
	"github.com/gofiber/fiber/v2"
)

// Response represents the standardized JSON response format
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// Success sends a successful response with data
func Success(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Data:    data,
	})
}

// Error sends an error response with message and status code
func Error(c *fiber.Ctx, message string, statusCode int) error {
	return c.Status(statusCode).JSON(Response{
		Success: false,
		Error: fiber.Map{
			"message": message,
		},
	})
}

// ValidationError sends a validation error response (400)
func ValidationError(c *fiber.Ctx, field, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(Response{
		Success: false,
		Error: fiber.Map{
			"field":   field,
			"message": message,
		},
	})
}

// AuthError sends an authentication error response (401)
func AuthError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(Response{
		Success: false,
		Error: fiber.Map{
			"message": message,
		},
	})
}
