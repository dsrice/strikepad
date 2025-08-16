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
	authService    service.AuthServiceInterface
	sessionService service.SessionServiceInterface
	validator      *validator.Validator
}

func NewAuthHandler(
	authService service.AuthServiceInterface,
	sessionService service.SessionServiceInterface,
) AuthHandlerInterface {
	return &AuthHandler{
		authService:    authService,
		sessionService: sessionService,
		validator:      validator.New(),
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
// @Summary User signup
// @Description Create a new user account with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.SignupRequest true "Signup request"
// @Success 201 {object} dto.AuthResponse "User created successfully"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 409 {object} dto.ErrorResponse "User already exists"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/auth/signup [post]
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

	// Create session and generate tokens
	tokenPair, err := h.sessionService.CreateSession(response.ID)
	if err != nil {
		slog.Error("Failed to create session after signup", "error", err, "user_id", response.ID)
		errorInfo := errors.GetErrorInfo(errors.ErrCodeInternalError)
		return c.JSON(errorInfo.HTTPStatus, dto.ErrorResponse{
			Code:        string(errorInfo.Code),
			Message:     errorInfo.Message,
			Description: "Failed to create session",
		})
	}

	// Create response with tokens
	signupResponse := dto.AuthResponse{
		SignupResponse: *response,
		AccessToken:    tokenPair.AccessToken,
		RefreshToken:   tokenPair.RefreshToken,
		ExpiresAt:      tokenPair.AccessTokenExpiresAt,
	}

	slog.Info("User signup successful", "user_id", response.ID, "email", response.Email)
	return c.JSON(http.StatusCreated, signupResponse)
}

// Login handles user authentication
// @Summary User login
// @Description Authenticate user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} dto.LoginResponse "Login successful"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Invalid credentials"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/auth/login [post]
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

	// Create session and generate tokens
	tokenPair, err := h.sessionService.CreateSession(userInfo.ID)
	if err != nil {
		slog.Error("Failed to create session after login", "error", err, "user_id", userInfo.ID)
		errorInfo := errors.GetErrorInfo(errors.ErrCodeInternalError)
		return c.JSON(errorInfo.HTTPStatus, dto.ErrorResponse{
			Code:        string(errorInfo.Code),
			Message:     errorInfo.Message,
			Description: "Failed to create session",
		})
	}

	// Create response with tokens
	loginResponse := dto.LoginResponse{
		UserInfo:     *userInfo,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.AccessTokenExpiresAt,
	}

	slog.Info("User login successful", "user_id", userInfo.ID, "email", userInfo.Email)
	return c.JSON(http.StatusOK, loginResponse)
}

// GoogleSignup handles user registration using Google OAuth
// @Summary Google OAuth signup
// @Description Create a new user account using Google OAuth
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.GoogleSignupRequest true "Google signup request"
// @Success 201 {object} dto.SignupResponse "User created successfully"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 409 {object} dto.ErrorResponse "User already exists"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/auth/google/signup [post]
func (h *AuthHandler) GoogleSignup(c echo.Context) error {
	var req dto.GoogleSignupRequest

	// Bind request body
	if err := c.Bind(&req); err != nil {
		slog.Warn("Invalid request body for Google signup", "error", err)
		errorInfo := errors.GetErrorInfo(errors.ErrCodeInvalidRequest)
		return c.JSON(errorInfo.HTTPStatus, dto.ErrorResponse{
			Code:        string(errorInfo.Code),
			Message:     errorInfo.Message,
			Description: errorInfo.Description,
		})
	}

	// Validate request using validator
	if err := h.validator.Validate(&req); err != nil {
		return h.handleValidationError(c, err, "Google signup")
	}

	// Call service
	response, err := h.authService.GoogleSignup(&req)
	if err != nil {
		// Handle specific errors
		switch err.Error() {
		case "invalid access token":
			errorInfo := errors.GetErrorInfo(errors.ErrCodeInvalidRequest)
			return c.JSON(errorInfo.HTTPStatus, dto.ErrorResponse{
				Code:        string(errorInfo.Code),
				Message:     errorInfo.Message,
				Description: "Invalid Google access token",
			})
		case auth.ErrUserAlreadyExists.Error():
			errorInfo := errors.GetErrorInfo(errors.ErrCodeUserExists)
			return c.JSON(errorInfo.HTTPStatus, dto.ErrorResponse{
				Code:        string(errorInfo.Code),
				Message:     errorInfo.Message,
				Description: errorInfo.Description,
			})
		default:
			slog.Error("Internal error during Google signup", "error", err)
			errorInfo := errors.GetErrorInfo(errors.ErrCodeInternalError)
			return c.JSON(errorInfo.HTTPStatus, dto.ErrorResponse{
				Code:        string(errorInfo.Code),
				Message:     errorInfo.Message,
				Description: errorInfo.Description,
			})
		}
	}

	slog.Info("Google user signup successful", "user_id", response.ID, "email", response.Email)
	return c.JSON(http.StatusCreated, response)
}

// GoogleLogin handles user authentication using Google OAuth
// @Summary Google OAuth login
// @Description Authenticate user using Google OAuth
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.GoogleLoginRequest true "Google login request"
// @Success 200 {object} dto.UserInfo "Login successful"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Invalid credentials"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/auth/google/login [post]
func (h *AuthHandler) GoogleLogin(c echo.Context) error {
	var req dto.GoogleLoginRequest

	// Bind request body
	if err := c.Bind(&req); err != nil {
		slog.Warn("Invalid request body for Google login", "error", err)
		errorInfo := errors.GetErrorInfo(errors.ErrCodeInvalidRequest)
		return c.JSON(errorInfo.HTTPStatus, dto.ErrorResponse{
			Code:        string(errorInfo.Code),
			Message:     errorInfo.Message,
			Description: errorInfo.Description,
		})
	}

	// Validate request using validator
	if err := h.validator.Validate(&req); err != nil {
		return h.handleValidationError(c, err, "Google login")
	}

	// Call service
	userInfo, err := h.authService.GoogleLogin(&req)
	if err != nil {
		// Handle specific errors
		switch err {
		case auth.ErrInvalidCredentials:
			errorInfo := errors.GetErrorInfo(errors.ErrCodeInvalidCredentials)
			return c.JSON(errorInfo.HTTPStatus, dto.ErrorResponse{
				Code:        string(errorInfo.Code),
				Message:     errorInfo.Message,
				Description: "Invalid Google credentials",
			})
		default:
			slog.Error("Internal error during Google login", "error", err)
			errorInfo := errors.GetErrorInfo(errors.ErrCodeInternalError)
			return c.JSON(errorInfo.HTTPStatus, dto.ErrorResponse{
				Code:        string(errorInfo.Code),
				Message:     errorInfo.Message,
				Description: errorInfo.Description,
			})
		}
	}

	slog.Info("Google user login successful", "user_id", userInfo.ID, "email", userInfo.Email)
	return c.JSON(http.StatusOK, userInfo)
}

// Logout handles user logout
// @Summary User logout
// @Description Logout current user and invalidate session
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "Logout successful"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(c echo.Context) error {
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

	accessToken, ok := c.Get("access_token").(string)
	if !ok {
		slog.Error("Failed to get access token from context")
		errorInfo := errors.GetErrorInfo(errors.ErrCodeInternalError)
		return c.JSON(errorInfo.HTTPStatus, dto.ErrorResponse{
			Code:        string(errorInfo.Code),
			Message:     errorInfo.Message,
			Description: "Failed to get token information",
		})
	}

	// Call session service to logout using JWT user_id
	err := h.sessionService.Logout(userID, accessToken)
	if err != nil {
		slog.Error("Failed to logout user", "error", err, "user_id", userID)
		errorInfo := errors.GetErrorInfo(errors.ErrCodeInternalError)
		return c.JSON(errorInfo.HTTPStatus, dto.ErrorResponse{
			Code:        string(errorInfo.Code),
			Message:     errorInfo.Message,
			Description: "Logout failed",
		})
	}

	slog.Info("User logout successful", "user_id", userID)
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Logout successful",
	})
}

