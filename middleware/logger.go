package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Logger middleware untuk log HTTP request details
// Log: HTTP method, path, status code, duration, request ID
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generate unique request ID
		requestID := uuid.New().String()

		// Set request ID dalam context untuk use dalam handler lain
		c.Locals("requestID", requestID)

		// Capture start time
		start := time.Now()

		// Continue ke next handler
		err := c.Next()

		// Capture response details
		duration := time.Since(start)
		method := c.Method()
		path := c.Path()
		status := c.Response().StatusCode()

		// Log request details
		c.Locals("log", fiber.Map{
			"request_id": requestID,
			"method":     method,
			"path":       path,
			"status":     status,
			"duration":   duration.String(),
		})

		return err
	}
}
