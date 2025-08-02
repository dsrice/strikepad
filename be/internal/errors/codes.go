package errors

import "net/http"

// ErrorCode represents standard error codes for the API
type ErrorCode string

const (
	// General error codes (E001-E099)
	ErrCodeInternalError    ErrorCode = "E001"
	ErrCodeInvalidRequest   ErrorCode = "E002"
	ErrCodeValidationFailed ErrorCode = "E003"
	ErrCodeNotFound         ErrorCode = "E004"
	ErrCodeUnauthorized     ErrorCode = "E005"
	ErrCodeForbidden        ErrorCode = "E006"
	ErrCodeConflict         ErrorCode = "E007"

	// Authentication error codes (E100-E199)
	ErrCodeInvalidCredentials ErrorCode = "E100"
	ErrCodeUserNotFound       ErrorCode = "E101"
	ErrCodeUserExists         ErrorCode = "E102"
	ErrCodeTokenExpired       ErrorCode = "E103"
	ErrCodeTokenInvalid       ErrorCode = "E104"

	// Validation error codes (E200-E299)
	ErrCodeEmailRequired       ErrorCode = "E200"
	ErrCodeEmailInvalid        ErrorCode = "E201"
	ErrCodePasswordRequired    ErrorCode = "E202"
	ErrCodePasswordTooShort    ErrorCode = "E203"
	ErrCodePasswordTooLong     ErrorCode = "E204"
	ErrCodePasswordComplexity  ErrorCode = "E205"
	ErrCodeDisplayNameRequired ErrorCode = "E206"
	ErrCodeDisplayNameTooLong  ErrorCode = "E207"

	// Business logic error codes (E300-E399)
	ErrCodeEmailNotVerified ErrorCode = "E300"
	ErrCodeAccountDisabled  ErrorCode = "E301"
	ErrCodeAccountDeleted   ErrorCode = "E302"
)

// ErrorInfo contains error information including code, message, description, and HTTP status
type ErrorInfo struct {
	Code        ErrorCode `json:"code"`
	Message     string    `json:"message"`
	Description string    `json:"description,omitempty"`
	HTTPStatus  int       `json:"-"` // Not included in JSON response
}

// GetErrorInfo returns error information for a given error code
func GetErrorInfo(code ErrorCode) ErrorInfo {
	errorMap := map[ErrorCode]ErrorInfo{
		// General errors (E001-E099)
		ErrCodeInternalError: {
			Code:        ErrCodeInternalError,
			Message:     "Internal server error",
			Description: "An unexpected error occurred on the server",
			HTTPStatus:  http.StatusInternalServerError,
		},
		ErrCodeInvalidRequest: {
			Code:        ErrCodeInvalidRequest,
			Message:     "Invalid request",
			Description: "The request format is invalid or malformed",
			HTTPStatus:  http.StatusBadRequest,
		},
		ErrCodeValidationFailed: {
			Code:        ErrCodeValidationFailed,
			Message:     "Validation failed",
			Description: "One or more fields failed validation",
			HTTPStatus:  http.StatusBadRequest,
		},
		ErrCodeNotFound: {
			Code:        ErrCodeNotFound,
			Message:     "Resource not found",
			Description: "The requested resource was not found",
			HTTPStatus:  http.StatusNotFound,
		},
		ErrCodeUnauthorized: {
			Code:        ErrCodeUnauthorized,
			Message:     "Unauthorized",
			Description: "Authentication is required to access this resource",
			HTTPStatus:  http.StatusUnauthorized,
		},
		ErrCodeForbidden: {
			Code:        ErrCodeForbidden,
			Message:     "Forbidden",
			Description: "You do not have permission to access this resource",
			HTTPStatus:  http.StatusForbidden,
		},
		ErrCodeConflict: {
			Code:        ErrCodeConflict,
			Message:     "Conflict",
			Description: "The request conflicts with the current state of the resource",
			HTTPStatus:  http.StatusConflict,
		},

		// Authentication errors (E100-E199)
		ErrCodeInvalidCredentials: {
			Code:        ErrCodeInvalidCredentials,
			Message:     "Invalid credentials",
			Description: "The provided email or password is incorrect",
			HTTPStatus:  http.StatusUnauthorized,
		},
		ErrCodeUserNotFound: {
			Code:        ErrCodeUserNotFound,
			Message:     "User not found",
			Description: "No user found with the provided email address",
			HTTPStatus:  http.StatusNotFound,
		},
		ErrCodeUserExists: {
			Code:        ErrCodeUserExists,
			Message:     "User already exists",
			Description: "A user with this email address already exists",
			HTTPStatus:  http.StatusConflict,
		},
		ErrCodeTokenExpired: {
			Code:        ErrCodeTokenExpired,
			Message:     "Token expired",
			Description: "The authentication token has expired",
			HTTPStatus:  http.StatusUnauthorized,
		},
		ErrCodeTokenInvalid: {
			Code:        ErrCodeTokenInvalid,
			Message:     "Invalid token",
			Description: "The authentication token is invalid or malformed",
			HTTPStatus:  http.StatusUnauthorized,
		},

		// Validation errors (E200-E299)
		ErrCodeEmailRequired: {
			Code:        ErrCodeEmailRequired,
			Message:     "Email is required",
			Description: "Email address must be provided",
			HTTPStatus:  http.StatusBadRequest,
		},
		ErrCodeEmailInvalid: {
			Code:        ErrCodeEmailInvalid,
			Message:     "Invalid email format",
			Description: "The email address format is invalid",
			HTTPStatus:  http.StatusBadRequest,
		},
		ErrCodePasswordRequired: {
			Code:        ErrCodePasswordRequired,
			Message:     "Password is required",
			Description: "Password must be provided",
			HTTPStatus:  http.StatusBadRequest,
		},
		ErrCodePasswordTooShort: {
			Code:        ErrCodePasswordTooShort,
			Message:     "Password too short",
			Description: "Password must be at least 8 characters long",
			HTTPStatus:  http.StatusBadRequest,
		},
		ErrCodePasswordTooLong: {
			Code:        ErrCodePasswordTooLong,
			Message:     "Password too long",
			Description: "Password must be at most 128 characters long",
			HTTPStatus:  http.StatusBadRequest,
		},
		ErrCodePasswordComplexity: {
			Code:        ErrCodePasswordComplexity,
			Message:     "Password complexity requirements not met",
			Description: "Password must contain at least one lowercase letter, one uppercase letter, and one symbol",
			HTTPStatus:  http.StatusBadRequest,
		},
		ErrCodeDisplayNameRequired: {
			Code:        ErrCodeDisplayNameRequired,
			Message:     "Display name is required",
			Description: "Display name must be provided",
			HTTPStatus:  http.StatusBadRequest,
		},
		ErrCodeDisplayNameTooLong: {
			Code:        ErrCodeDisplayNameTooLong,
			Message:     "Display name too long",
			Description: "Display name must be at most 100 characters long",
			HTTPStatus:  http.StatusBadRequest,
		},

		// Business logic errors (E300-E399)
		ErrCodeEmailNotVerified: {
			Code:        ErrCodeEmailNotVerified,
			Message:     "Email not verified",
			Description: "Email address must be verified before performing this action",
			HTTPStatus:  http.StatusForbidden,
		},
		ErrCodeAccountDisabled: {
			Code:        ErrCodeAccountDisabled,
			Message:     "Account disabled",
			Description: "This account has been disabled",
			HTTPStatus:  http.StatusForbidden,
		},
		ErrCodeAccountDeleted: {
			Code:        ErrCodeAccountDeleted,
			Message:     "Account deleted",
			Description: "This account has been deleted",
			HTTPStatus:  http.StatusForbidden,
		},
	}

	if info, exists := errorMap[code]; exists {
		return info
	}

	// Return default error if code not found
	return ErrorInfo{
		Code:        code,
		Message:     "Unknown error",
		Description: "An unknown error occurred",
		HTTPStatus:  http.StatusInternalServerError,
	}
}
