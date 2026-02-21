package main

import (
	"log"
	"os"

	"kanban-backend/config"
	"kanban-backend/handlers"
	"kanban-backend/repositories"
	"kanban-backend/routes"
	"kanban-backend/services"
	"kanban-backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()
	config.ConnectS3()

	userRepo := repositories.NewUserRepository()
	boardRepo := repositories.NewBoardRepository()
	columnRepo := repositories.NewColumnRepository()
	taskRepo := repositories.NewTaskRepository()

	authService := services.NewAuthService(userRepo)
	boardService := services.NewBoardService(boardRepo, columnRepo)
	taskService := services.NewTaskService(taskRepo, columnRepo)

	authController := handlers.NewAuthController(authService)
	boardController := handlers.NewBoardController(boardService)
	taskController := handlers.NewTaskController(taskService)

	app := fiber.New(fiber.Config{
		AppName:      "Kanban API v1.0",
		ErrorHandler: utils.ErrorHandler,
	})

	routes.Setup(app, authService, authController, boardController, taskController)

	port := os.Getenv("PORT")
	log.Printf("🚀 Server running on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
