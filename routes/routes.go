package routes

import (
	"kanban-backend/handlers"
	"kanban-backend/middleware"
	"kanban-backend/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Setup(app *fiber.App, authService services.AuthService, authController *handlers.AuthController, boardController *handlers.BoardController, taskController *handlers.TaskController) {
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
}
