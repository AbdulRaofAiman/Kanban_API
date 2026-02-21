package handlers

import (
	"errors"
	"regexp"
	"time"

	"kanban-backend/services"
	"kanban-backend/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

func (ctrl *AuthController) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, "Invalid request body", fiber.StatusBadRequest)
	}

	if req.Username == "" {
		return utils.ValidationError(c, "username", "username is required")
	}

	if req.Email == "" {
		return utils.ValidationError(c, "email", "email is required")
	}

	if !isValidEmail(req.Email) {
		return utils.ValidationError(c, "email", "invalid email format")
	}

	if req.Password == "" {
		return utils.ValidationError(c, "password", "password is required")
	}

	if len(req.Password) < 8 {
		return utils.ValidationError(c, "password", "password must be at least 8 characters long")
	}

	user, err := ctrl.authService.Register(c.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		var conflictErr utils.ErrConflict
		if errors.As(err, &conflictErr) {
			return utils.Error(c, err.Error(), fiber.StatusConflict)
		}
		var validationErr utils.ErrValidation
		if errors.As(err, &validationErr) {
			return utils.Error(c, err.Error(), fiber.StatusBadRequest)
		}
		return utils.Error(c, "Failed to register user", fiber.StatusInternalServerError)
	}

	token, err := ctrl.authService.GenerateToken(user.ID, 24*time.Hour)
	if err != nil {
		return utils.Error(c, "Failed to generate token", fiber.StatusInternalServerError)
	}

	return utils.Success(c, AuthResponse{
		User: UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Token: token,
	})
}

func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, "Invalid request body", fiber.StatusBadRequest)
	}

	if req.Email == "" {
		return utils.ValidationError(c, "email", "email is required")
	}

	if !isValidEmail(req.Email) {
		return utils.ValidationError(c, "email", "invalid email format")
	}

	if req.Password == "" {
		return utils.ValidationError(c, "password", "password is required")
	}

	token, err := ctrl.authService.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		var unauthorizedErr utils.ErrUnauthorized
		if errors.As(err, &unauthorizedErr) {
			return utils.AuthError(c, err.Error())
		}
		return utils.Error(c, "Failed to login", fiber.StatusInternalServerError)
	}

	return utils.Success(c, fiber.Map{
		"token": token,
	})
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
