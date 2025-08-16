package dto

import "time"

// SignupRequest represents the request payload for user signup
type SignupRequest struct {
	Email       string `json:"email" validate:"required,email,max=255" example:"user@example.com" swaggertype:"string" format:"email"`
	Password    string `json:"password" validate:"required,min=8,max=128,password_complex" example:"Password123!" swaggertype:"string" minLength:"8" maxLength:"128"`
	DisplayName string `json:"display_name" validate:"required,min=1,max=100" example:"John Doe" swaggertype:"string" minLength:"1" maxLength:"100"`
}

// GoogleSignupRequest represents the request payload for Google OAuth signup
type GoogleSignupRequest struct {
	AccessToken string `json:"access_token" validate:"required" example:"ya29.a0ARrdaM..." swaggertype:"string"`
}

// SignupResponse represents the response payload for user signup
type SignupResponse struct {
	CreatedAt     time.Time `json:"created_at" example:"2025-01-27T10:15:30Z" swaggertype:"string" format:"date-time"`
	Email         string    `json:"email" example:"user@example.com" swaggertype:"string" format:"email"`
	DisplayName   string    `json:"display_name" example:"John Doe" swaggertype:"string"`
	ID            uint      `json:"id" example:"1" swaggertype:"integer"`
	EmailVerified bool      `json:"email_verified" example:"false" swaggertype:"boolean"`
}

// LoginRequest represents the request payload for user login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=255" example:"user@example.com" swaggertype:"string" format:"email"`
	Password string `json:"password" validate:"required,min=1,max=128" example:"password123" swaggertype:"string"`
}

// GoogleLoginRequest represents the request payload for Google OAuth login
type GoogleLoginRequest struct {
	AccessToken string `json:"access_token" validate:"required" example:"ya29.a0ARrdaM..." swaggertype:"string"`
}

// LoginResponse represents the response payload for user login
type LoginResponse struct {
	ExpiresAt    time.Time `json:"expires_at" swaggertype:"string" format:"date-time"`
	AccessToken  string    `json:"access_token" swaggertype:"string"`
	RefreshToken string    `json:"refresh_token" swaggertype:"string"`
}

// UserInfo represents basic user information
type UserInfo struct {
	Email         string `json:"email" swaggertype:"string" format:"email"`
	DisplayName   string `json:"display_name" swaggertype:"string"`
	ID            uint   `json:"id" swaggertype:"integer"`
	EmailVerified bool   `json:"email_verified" swaggertype:"boolean"`
}

// ErrorResponse represents a unified error response structure
type ErrorResponse struct {
	Code        string            `json:"code" swaggertype:"string"`
	Message     string            `json:"message" swaggertype:"string"`
	Description string            `json:"description,omitempty" swaggertype:"string"`
	Details     []ValidationError `json:"details,omitempty"`
}

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string `json:"field" swaggertype:"string"`
	Tag     string `json:"tag" swaggertype:"string"`
	Value   string `json:"value" swaggertype:"string"`
	Message string `json:"message" swaggertype:"string"`
}

// AuthResponse represents the response payload for signup with tokens
type AuthResponse struct {
	ExpiresAt      time.Time `json:"expires_at" swaggertype:"string" format:"date-time"`
	AccessToken    string    `json:"access_token" swaggertype:"string"`
	RefreshToken   string    `json:"refresh_token" swaggertype:"string"`
	SignupResponse `json:",inline"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status  string `json:"status" swaggertype:"string"`
	Message string `json:"message" swaggertype:"string"`
}