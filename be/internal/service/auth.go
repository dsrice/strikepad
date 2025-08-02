package service

import (
	"errors"
	"log/slog"

	"strikepad-backend/internal/auth"
	"strikepad-backend/internal/dto"
	"strikepad-backend/internal/model"
	"strikepad-backend/internal/repository"

	"gorm.io/gorm"
)

type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// Signup creates a new user account
func (s *AuthService) Signup(req *dto.SignupRequest) (*dto.SignupResponse, error) {
	// Validate email format
	if err := auth.ValidateEmail(req.Email); err != nil {
		slog.Warn("Invalid email format during signup", "email", req.Email, "error", err)
		return nil, err
	}

	// Validate password
	if err := auth.ValidatePassword(req.Password); err != nil {
		slog.Warn("Invalid password during signup", "error", err)
		return nil, err
	}

	// Normalize email
	normalizedEmail := auth.NormalizeEmail(req.Email)

	// Check if user already exists
	existingUser, err := s.userRepo.FindByEmail(normalizedEmail)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Error("Failed to check existing user", "email", normalizedEmail, "error", err)
		return nil, errors.New("internal server error")
	}
	if existingUser != nil {
		slog.Warn("User already exists", "email", normalizedEmail)
		return nil, auth.ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		slog.Error("Failed to hash password", "error", err)
		return nil, errors.New("internal server error")
	}

	// Create user
	user := &model.User{
		ProviderType:   "email",
		ProviderUserID: nil,
		Email:          &normalizedEmail,
		DisplayName:    req.DisplayName,
		PasswordHash:   &hashedPassword,
		EmailVerified:  false,
		IsDeleted:      false,
	}

	createdUser, err := s.userRepo.Create(user)
	if err != nil {
		slog.Error("Failed to create user", "email", normalizedEmail, "error", err)
		return nil, errors.New("internal server error")
	}

	slog.Info("User created successfully", "user_id", createdUser.ID, "email", normalizedEmail)

	// Return response
	response := &dto.SignupResponse{
		ID:            createdUser.ID,
		Email:         normalizedEmail,
		DisplayName:   createdUser.DisplayName,
		EmailVerified: createdUser.EmailVerified,
		CreatedAt:     createdUser.CreatedAt,
	}

	return response, nil
}

// Login authenticates a user and returns user information
func (s *AuthService) Login(req *dto.LoginRequest) (*dto.UserInfo, error) {
	// Validate email format
	if err := auth.ValidateEmail(req.Email); err != nil {
		slog.Warn("Invalid email format during login", "email", req.Email, "error", err)
		return nil, auth.ErrInvalidCredentials
	}

	// Normalize email
	normalizedEmail := auth.NormalizeEmail(req.Email)

	// Find user by email
	user, err := s.userRepo.FindByEmail(normalizedEmail)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Warn("Login attempt with non-existent email", "email", normalizedEmail)
			return nil, auth.ErrInvalidCredentials
		}
		slog.Error("Failed to find user during login", "email", normalizedEmail, "error", err)
		return nil, errors.New("internal server error")
	}

	// Check if user is deleted
	if user.IsDeleted {
		slog.Warn("Login attempt with deleted user", "user_id", user.ID, "email", normalizedEmail)
		return nil, auth.ErrInvalidCredentials
	}

	// Check if password hash exists (for email provider)
	if user.PasswordHash == nil {
		slog.Warn("Login attempt for user without password", "user_id", user.ID, "email", normalizedEmail)
		return nil, auth.ErrInvalidCredentials
	}

	// Verify password
	if !auth.CheckPasswordHash(req.Password, *user.PasswordHash) {
		slog.Warn("Invalid password during login", "user_id", user.ID, "email", normalizedEmail)
		return nil, auth.ErrInvalidCredentials
	}

	slog.Info("User logged in successfully", "user_id", user.ID, "email", normalizedEmail)

	// Return user info
	userInfo := &dto.UserInfo{
		ID:            user.ID,
		Email:         normalizedEmail,
		DisplayName:   user.DisplayName,
		EmailVerified: user.EmailVerified,
	}

	return userInfo, nil
}