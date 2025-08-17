package service

import (
	"errors"
	"log/slog"

	"strikepad-backend/internal/auth"
	"strikepad-backend/internal/dto"
	"strikepad-backend/internal/repository"

	"gorm.io/gorm"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserServiceInterface {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetCurrentUser retrieves user information by user ID
func (s *UserService) GetCurrentUser(userID uint) (*dto.UserInfo, error) {
	// Find user by ID
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Warn("GetCurrentUser attempt with non-existent user ID", "user_id", userID)
			return nil, auth.ErrInvalidCredentials
		}
		slog.Error("Failed to find user by ID", "user_id", userID, "error", err)
		return nil, errors.New("internal server error")
	}

	// Check if user is deleted
	if user.IsDeleted {
		slog.Warn("GetCurrentUser attempt with deleted user", "user_id", userID)
		return nil, auth.ErrInvalidCredentials
	}

	// Return user info
	userInfo := &dto.UserInfo{
		ID:            user.ID,
		Email:         *user.Email,
		DisplayName:   user.DisplayName,
		EmailVerified: user.EmailVerified,
	}

	return userInfo, nil
}
