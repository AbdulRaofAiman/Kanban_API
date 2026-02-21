package routes

import (
	"kanban-backend/controllers"
	"kanban-backend/middleware"
	"kanban-backend/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Setup(app *fiber.App, authService services.AuthService, authController *controllers.AuthController, boardController *controllers.BoardController, taskController *controllers.TaskController, commentController *controllers.CommentController, labelController *controllers.LabelController, attachmentController *controllers.AttachmentController) {
	app.Use(middleware.Logger())
	app.Use(cors.New(middleware.CORSConfig()))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	auth := app.Group("/api/v1/auth")
	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)

	boards := app.Group("/api/v1/boards")
	boards.Use(middleware.AuthMiddleware(authService))
	boards.Post("/", boardController.Create)
	boards.Get("/:id", boardController.FindByID)
	boards.Get("/", boardController.FindAll)
	boards.Put("/:id", boardController.Update)
	boards.Delete("/:id", boardController.Delete)

	tasks := app.Group("/api/v1/tasks")
	tasks.Use(middleware.AuthMiddleware(authService))
	tasks.Post("/", taskController.Create)
	tasks.Get("/:id", taskController.FindByID)
	tasks.Get("/column/:columnId", taskController.FindByColumnID)
	tasks.Put("/:id", taskController.Update)
	tasks.Delete("/:id", taskController.Delete)
	tasks.Put("/:id/move", taskController.Move)
	tasks.Post("/:id/labels/:label_id", labelController.AddToTask)
	tasks.Delete("/:id/labels/:label_id", labelController.RemoveFromTask)

	comments := app.Group("/api/v1/comments")
	comments.Use(middleware.AuthMiddleware(authService))
	comments.Post("/", commentController.Create)
	comments.Get("/:id", commentController.FindByID)
	comments.Get("/task/:task_id", commentController.FindByTaskID)
	comments.Put("/:id", commentController.Update)
	comments.Delete("/:id", commentController.Delete)

	labels := app.Group("/api/v1/labels")
	labels.Use(middleware.AuthMiddleware(authService))
	labels.Post("/", labelController.Create)
	labels.Get("/", labelController.FindAll)
	labels.Get("/:id", labelController.FindByID)
	labels.Put("/:id", labelController.Update)
	labels.Delete("/:id", labelController.Delete)

	attachments := app.Group("/api/v1/attachments")
	attachments.Use(middleware.AuthMiddleware(authService))
	attachments.Post("/", attachmentController.Create)
	attachments.Get("/:id", attachmentController.FindByID)
	attachments.Get("/task/:task_id", attachmentController.FindByTaskID)
	attachments.Put("/:id", attachmentController.Update)
	attachments.Delete("/:id", attachmentController.Delete)
}
