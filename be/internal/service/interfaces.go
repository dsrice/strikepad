package service

import "strikepad-backend/internal/dto"

// AuthServiceInterface defines the interface for authentication service
type AuthServiceInterface interface {
	Signup(req *dto.SignupRequest) (*dto.SignupResponse, error)
	Login(req *dto.LoginRequest) (*dto.UserInfo, error)
}

// HealthServiceInterface defines the interface for health service
type HealthServiceInterface interface {
	GetHealth() *dto.HealthResponse
}
