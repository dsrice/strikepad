package middleware

import (
	"log/slog"
	"strings"

	"strikepad-backend/internal/errors"
	"strikepad-backend/internal/service"

	"github.com/labstack/echo/v4"
)

// JWTMiddleware handles JWT token authentication
func JWTMiddleware(sessionService service.SessionServiceInterface) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				slog.Warn("Missing authorization header")
				errorInfo := errors.GetErrorInfo(errors.ErrCodeUnauthorized)
				return c.JSON(errorInfo.HTTPStatus, map[string]string{
					"code":    string(errorInfo.Code),
					"message": errorInfo.Message,
				})
			}

			// Check Bearer token format
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				slog.Warn("Invalid authorization header format")
				errorInfo := errors.GetErrorInfo(errors.ErrCodeUnauthorized)
				return c.JSON(errorInfo.HTTPStatus, map[string]string{
					"code":    string(errorInfo.Code),
					"message": "Invalid authorization header format",
				})
			}

			accessToken := tokenParts[1]

			// Validate access token
			session, err := sessionService.ValidateAccessToken(accessToken)
			if err != nil {
				slog.Warn("Invalid access token", "error", err)
				errorInfo := errors.GetErrorInfo(errors.ErrCodeUnauthorized)
				return c.JSON(errorInfo.HTTPStatus, map[string]string{
					"code":    string(errorInfo.Code),
					"message": "Invalid or expired token",
				})
			}

			// Store session and user info in context
			c.Set("session", session)
			c.Set("user_id", session.UserID)
			c.Set("access_token", accessToken)

			return next(c)
		}
	}
}

// GetUserIDFromContext extracts user ID from echo context
func GetUserIDFromContext(c echo.Context) (uint, bool) {
	userID, ok := c.Get("user_id").(uint)
	return userID, ok
}

// GetAccessTokenFromContext extracts access token from echo context
func GetAccessTokenFromContext(c echo.Context) (string, bool) {
	token, ok := c.Get("access_token").(string)
	return token, ok
}
