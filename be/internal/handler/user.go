package handler

import (
	"log/slog"
	"net/http"

	"strikepad-backend/internal/service"

	"strikepad-backend/internal/auth"
	"strikepad-backend/internal/dto"
	"strikepad-backend/internal/errors"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService service.UserServiceInterface
}

func NewUserHandler(userService service.UserServiceInterface) UserHandlerInterface {
	return &UserHandler{
		userService: userService,
	}
}

// Me returns current user information
// @Summary Get current user information
// @Description Get information about the currently authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserInfo "User information"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/user/me [get]
func (h *UserHandler) Me(c echo.Context) error {
	// Get user ID from JWT claims (set by JWT middleware)
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		slog.Error("Failed to get user ID from JWT token")
		errorInfo := errors.GetErrorInfo(errors.ErrCodeUnauthorized)
		return c.JSON(errorInfo.HTTPStatus, dto.ErrorResponse{
			Code:        string(errorInfo.Code),
			Message:     errorInfo.Message,
			Description: "Invalid token: user ID not found",
		})
	}

	// Get user information from service
	userInfo, err := h.userService.GetCurrentUser(userID)
	if err != nil {
		// Handle specific errors
		switch err {
		case auth.ErrInvalidCredentials:
			errorInfo := errors.GetErrorInfo(errors.ErrCodeUnauthorized)
			return c.JSON(errorInfo.HTTPStatus, dto.ErrorResponse{
				Code:        string(errorInfo.Code),
				Message:     errorInfo.Message,
				Description: "User not found or inactive",
			})
		default:
			slog.Error("Internal error during GetCurrentUser", "error", err, "user_id", userID)
			errorInfo := errors.GetErrorInfo(errors.ErrCodeInternalError)
			return c.JSON(errorInfo.HTTPStatus, dto.ErrorResponse{
				Code:        string(errorInfo.Code),
				Message:     errorInfo.Message,
				Description: errorInfo.Description,
			})
		}
	}

	slog.Info("GetCurrentUser successful", "user_id", userID)
	return c.JSON(http.StatusOK, userInfo)
}