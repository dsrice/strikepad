package handler

import (
	"log/slog"
	"net/http"

	"strikepad-backend/internal/service"

	"strikepad-backend/internal/auth"
	"strikepad-backend/internal/dto"
	"strikepad-backend/internal/errors"
	"strikepad-backend/internal/validator"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService service.AuthServiceInterface
	validator   *validator.Validator
}

func NewAuthHandler(authService service.AuthServiceInterface) AuthHandlerInterface {
	return &AuthHandler{
		authService: authService,
		validator:   validator.New(),
	}
}

// handleValidationError handles validation errors and returns appropriate JSON response
func (h *AuthHandler) handleValidationError(c echo.Context, err error, operation string) error {
	slog.Warn("Validation failed for "+operation, "error", err)
	if ve, ok := err.(validator.ValidationErrors); ok {
		errorInfo := errors.GetErrorInfo(errors.ErrCodeValidationFailed)

		// Convert validator errors to our format
		var validationErrors []dto.ValidationError
		for _, validatorErr := range ve.Errors {
			validationErrors = append(validationErrors, dto.ValidationError{
				Field:   validatorErr.Field,
				Tag:     validatorErr.Tag,
				Value:   validatorErr.Value,
				Message: validatorErr.Message,
			})
		}

		return c.JSON(errorInfo.HTTPStatus, dto.ErrorResponse{
			Code:        string(errorInfo.Code),
			Message:     errorInfo.Message,
			Description: errorInfo.Description,
			Details:     validationErrors,
		})
	}
	errorInfo := errors.GetErrorInfo(errors.ErrCodeValidationFailed)
	return c.JSON(errorInfo.HTTPStatus, dto.ErrorResponse{
		Code:        string(errorInfo.Code),
		Message:     errorInfo.Message,
		Description: err.Error(),
	})
}

// Signup handles user registration
func (h *AuthHandler) Signup(c echo.Context) error {
	var req dto.SignupRequest

	// Bind request body
	if err := c.Bind(&req); err != nil {
		slog.Warn("Invalid request body for signup", "error", err)
		errorInfo := errors.GetErrorInfo(errors.ErrCodeInvalidRequest)
		return c.JSON(errorInfo.HTTPStatus, dto.ErrorResponse{
			Code:        string(errorInfo.Code),
			Message:     errorInfo.Message,
			Description: errorInfo.Description,
		})
	}

	// Validate request using validator
	if err := h.validator.Validate(&req); err != nil {
		return h.handleValidationError(c, err, "signup")
	}

	// Call service
	response, err := h.authService.Signup(&req)
	if err != nil {
		// Handle specific errors
		switch err {
		case auth.ErrInvalidEmail:
			errorInfo := errors.GetErrorInfo(errors.ErrCodeEmailInvalid)
			return c.JSON(errorInfo.HTTPStatus, dto.ErrorResponse{
				Code:        string(errorInfo.Code),
				Message:     errorInfo.Message,
				Description: errorInfo.Description,
			})
		case auth.ErrPasswordTooShort:
			errorInfo := errors.GetErrorInfo(errors.ErrCodePasswordTooShort)
			return c.JSON(errorInfo.HTTPStatus, dto.ErrorResponse{
				Code:        string(errorInfo.Code),
				Message:     errorInfo.Message,
				Description: errorInfo.Description,
			})
		case auth.ErrPasswordTooLong:
			errorInfo := errors.GetErrorInfo(errors.ErrCodePasswordTooLong)
			return c.JSON(errorInfo.HTTPStatus, dto.ErrorResponse{
				Code:        string(errorInfo.Code),
				Message:     errorInfo.Message,
				Description: errorInfo.Description,
			})
		case auth.ErrUserAlreadyExists:
			errorInfo := errors.GetErrorInfo(errors.ErrCodeUserExists)
			return c.JSON(errorInfo.HTTPStatus, dto.ErrorResponse{
				Code:        string(errorInfo.Code),
				Message:     errorInfo.Message,
				Description: errorInfo.Description,
			})
		default:
			slog.Error("Internal error during signup", "error", err)
			errorInfo := errors.GetErrorInfo(errors.ErrCodeInternalError)
			return c.JSON(errorInfo.HTTPStatus, dto.ErrorResponse{
				Code:        string(errorInfo.Code),
				Message:     errorInfo.Message,
				Description: errorInfo.Description,
			})
		}
	}

	slog.Info("User signup successful", "user_id", response.ID, "email", response.Email)
	return c.JSON(http.StatusCreated, response)
}

// Login handles user authentication
func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest

	// Bind request body
	if err := c.Bind(&req); err != nil {
		slog.Warn("Invalid request body for login", "error", err)
		errorInfo := errors.GetErrorInfo(errors.ErrCodeInvalidRequest)
		return c.JSON(errorInfo.HTTPStatus, dto.ErrorResponse{
			Code:        string(errorInfo.Code),
			Message:     errorInfo.Message,
			Description: errorInfo.Description,
		})
	}

	// Validate request using validator
	if err := h.validator.Validate(&req); err != nil {
		return h.handleValidationError(c, err, "login")
	}

	// Call service
	userInfo, err := h.authService.Login(&req)
	if err != nil {
		// Handle specific errors
		switch err {
		case auth.ErrInvalidCredentials:
			errorInfo := errors.GetErrorInfo(errors.ErrCodeInvalidCredentials)
			return c.JSON(errorInfo.HTTPStatus, dto.ErrorResponse{
				Code:        string(errorInfo.Code),
				Message:     errorInfo.Message,
				Description: errorInfo.Description,
			})
		default:
			slog.Error("Internal error during login", "error", err)
			errorInfo := errors.GetErrorInfo(errors.ErrCodeInternalError)
			return c.JSON(errorInfo.HTTPStatus, dto.ErrorResponse{
				Code:        string(errorInfo.Code),
				Message:     errorInfo.Message,
				Description: errorInfo.Description,
			})
		}
	}

	slog.Info("User login successful", "user_id", userInfo.ID, "email", userInfo.Email)
	return c.JSON(http.StatusOK, userInfo)
}
