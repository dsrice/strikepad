package dto

import "time"

// SignupRequest represents the request payload for user signup
type SignupRequest struct {
	Email       string `json:"email" validate:"required,email,max=255" example:"user@example.com"`
	Password    string `json:"password" validate:"required,min=8,max=128,password_complex" example:"Password123!"`
	DisplayName string `json:"display_name" validate:"required,min=1,max=100" example:"John Doe"`
}

// SignupResponse represents the response payload for user signup
type SignupResponse struct {
	ID            uint      `json:"id" example:"1"`
	EmailVerified bool      `json:"email_verified" example:"false"`
	Email         string    `json:"email" example:"user@example.com"`
	DisplayName   string    `json:"display_name" example:"John Doe"`
	CreatedAt     time.Time `json:"created_at" example:"2025-01-27T10:15:30Z"`
}

// LoginRequest represents the request payload for user login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=255" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=1,max=128" example:"password123"`
}

// LoginResponse represents the response payload for user login
type LoginResponse struct {
	ExpiresAt   time.Time `json:"expires_at"`
	User        UserInfo  `json:"user"`
	AccessToken string    `json:"access_token"`
}

// UserInfo represents basic user information
type UserInfo struct {
	ID            uint   `json:"id"`
	EmailVerified bool   `json:"email_verified"`
	Email         string `json:"email"`
	DisplayName   string `json:"display_name"`
}

// ErrorResponse represents a unified error response structure
type ErrorResponse struct {
	Code        string            `json:"code"`
	Message     string            `json:"message"`
	Description string            `json:"description,omitempty"`
	Details     []ValidationError `json:"details,omitempty"`
}

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}
