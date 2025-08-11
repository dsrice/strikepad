package service

import "strikepad-backend/internal/dto"

// AuthServiceInterface defines the interface for authentication service
type AuthServiceInterface interface {
	Signup(req *dto.SignupRequest) (*dto.SignupResponse, error)
	Login(req *dto.LoginRequest) (*dto.UserInfo, error)
	GoogleSignup(req *dto.GoogleSignupRequest) (*dto.SignupResponse, error)
	GoogleLogin(req *dto.GoogleLoginRequest) (*dto.UserInfo, error)
}

// HealthServiceInterface defines the interface for health service
type HealthServiceInterface interface {
	GetHealth() *dto.HealthResponse
}

// APIServiceInterface defines the interface for API service
type APIServiceInterface interface {
	GetTestMessage() map[string]string
}
