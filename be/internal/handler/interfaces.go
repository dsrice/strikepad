package handler

import "github.com/labstack/echo/v4"

// AuthHandlerInterface defines the interface for authentication handlers
type AuthHandlerInterface interface {
	Signup(c echo.Context) error
	Login(c echo.Context) error
	GoogleSignup(c echo.Context) error
	GoogleLogin(c echo.Context) error
	Logout(c echo.Context) error
	Refresh(c echo.Context) error
}

// HealthHandlerInterface defines the interface for health handlers
type HealthHandlerInterface interface {
	Check(c echo.Context) error
}

// UserHandlerInterface defines the interface for user handlers
type UserHandlerInterface interface {
	Me(c echo.Context) error
}
