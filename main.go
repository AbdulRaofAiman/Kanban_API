package main

import (
	"log"
	"os"

	"kanban-backend/config"
	"kanban-backend/controllers"
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
	commentRepo := repositories.NewCommentRepository()
	labelRepo := repositories.NewLabelRepository()
	attachmentRepo := repositories.NewAttachmentRepository()

	authService := services.NewAuthService(userRepo)
	boardService := services.NewBoardService(boardRepo, columnRepo)
	taskService := services.NewTaskService(taskRepo, columnRepo)
	commentService := services.NewCommentService(commentRepo, taskRepo)
	labelService := services.NewLabelService(labelRepo, taskRepo)
	attachmentService := services.NewAttachmentService(attachmentRepo, taskRepo)

	authController := controllers.NewAuthController(authService)
	boardController := controllers.NewBoardController(boardService)
	taskController := controllers.NewTaskController(taskService)
	commentController := controllers.NewCommentController(commentService)
	labelController := controllers.NewLabelController(labelService)
	attachmentController := controllers.NewAttachmentController(attachmentService)

	app := fiber.New(fiber.Config{
		AppName:      "Kanban API v1.0",
		ErrorHandler: utils.ErrorHandler,
	})

	routes.Setup(app, authService, authController, boardController, taskController, commentController, labelController, attachmentController)

	port := os.Getenv("PORT")
	log.Printf("🚀 Server running on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
