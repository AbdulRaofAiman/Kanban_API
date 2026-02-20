package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORSConfig returns configured CORS middleware
func CORSConfig() cors.Config {
	// Get allowed origins from environment or default to "*" for development
	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		// Default to allow all origins in development
		allowedOrigins = "*"
	}

	// When AllowCredentials is true, AllowOrigins cannot be "*" (CORS security requirement)
	allowCredentials := allowedOrigins != "*"

	return cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     "GET, POST, PUT, DELETE, PATCH, OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Authorization, Accept, X-Requested-With",
		AllowCredentials: allowCredentials,
		ExposeHeaders:    "Content-Length, Content-Type",
		MaxAge:           86400, // 24 hours
	}
}
