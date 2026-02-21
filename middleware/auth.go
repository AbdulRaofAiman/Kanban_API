package middleware

import (
	"strings"

	"kanban-backend/services"
	"kanban-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(authService services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return utils.AuthError(c, "Authorization header is required")
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return utils.AuthError(c, "Invalid authorization header format")
		}

		token := tokenParts[1]

		userID, err := authService.ValidateToken(token)
		if err != nil {
			return utils.AuthError(c, err.Error())
		}

		c.Locals("user_id", userID)

		return c.Next()
	}
}
