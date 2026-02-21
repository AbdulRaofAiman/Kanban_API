package services

import (
	"context"
	"time"

	"kanban-backend/models"
	"kanban-backend/repositories"
	"kanban-backend/utils"
)

type AuthService interface {
	Register(ctx context.Context, username, email, password string) (*models.User, error)
	Login(ctx context.Context, email, password string) (string, error)
	GenerateToken(userID string, expiry time.Duration) (string, error)
	ValidateToken(tokenString string) (string, error)
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) error
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) Register(ctx context.Context, username, email, password string) (*models.User, error) {
	if len(password) < 8 {
		return nil, utils.NewValidation("password must be at least 8 characters long")
	}

	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, utils.NewConflict("user with this email already exists")
	}

	user := &models.User{
		Username: username,
		Email:    email,
		Password: password,
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", utils.NewUnauthorized("invalid email or password")
	}

	err = utils.CheckPassword(user.Password, password)
	if err != nil {
		return "", utils.NewUnauthorized("invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID, 24*time.Hour)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) GenerateToken(userID string, expiry time.Duration) (string, error) {
	return utils.GenerateToken(userID, expiry)
}

func (s *authService) ValidateToken(tokenString string) (string, error) {
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		return "", utils.NewUnauthorized("invalid or expired token")
	}
	return claims.UserID, nil
}

func (s *authService) HashPassword(password string) (string, error) {
	return utils.HashPassword(password)
}

func (s *authService) VerifyPassword(hashedPassword, password string) error {
	return utils.CheckPassword(hashedPassword, password)
}
