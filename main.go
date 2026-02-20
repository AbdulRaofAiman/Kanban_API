package main

import (
	"log"
	"os"

	"kanban-backend/config"
	"kanban-backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect database
	config.ConnectDB()
	config.ConnectS3()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Kanban API v1.0",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	app.Use(utils.ErrorHandler)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Kanban API is running 🚀",
		})
	})

	// Start server
	port := os.Getenv("PORT")
	log.Printf("🚀 Server running on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
