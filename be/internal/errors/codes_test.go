package errors_test

import (
	"net/http"
	"strings"
	"testing"

	"strikepad-backend/internal/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ErrorCodesTestSuite struct {
	suite.Suite
}

func (suite *ErrorCodesTestSuite) TestGetErrorInfo_AllErrorCodes() {
	// Comprehensive test for all error codes in a single table-driven test
	tests := []struct {
		name             string
		code             errors.ErrorCode
		expectedCode     errors.ErrorCode
		expectedMsg      string
		category         string
		descriptionCheck []string
		expectedStatus   int
	}{
		// General errors
		{
			name:             "Internal error",
			code:             errors.ErrCodeInternalError,
			expectedCode:     errors.ErrCodeInternalError,
			expectedStatus:   http.StatusInternalServerError,
			expectedMsg:      "Internal server error",
			category:         "general",
			descriptionCheck: []string{"unexpected", "server"},
		},
		{
			name:             "Invalid request",
			code:             errors.ErrCodeInvalidRequest,
			expectedCode:     errors.ErrCodeInvalidRequest,
			expectedStatus:   http.StatusBadRequest,
			expectedMsg:      "Invalid request",
			category:         "general",
			descriptionCheck: []string{"request"},
		},
		{
			name:             "Validation failed",
			code:             errors.ErrCodeValidationFailed,
			expectedCode:     errors.ErrCodeValidationFailed,
			expectedStatus:   http.StatusBadRequest,
			expectedMsg:      "Validation failed",
			category:         "general",
			descriptionCheck: []string{"validation"},
		},
		{
			name:             "Not found",
			code:             errors.ErrCodeNotFound,
			expectedCode:     errors.ErrCodeNotFound,
			expectedStatus:   http.StatusNotFound,
			expectedMsg:      "Resource not found",
			category:         "general",
			descriptionCheck: []string{"resource"},
		},
		{
			name:             "Unauthorized",
			code:             errors.ErrCodeUnauthorized,
			expectedCode:     errors.ErrCodeUnauthorized,
			expectedStatus:   http.StatusUnauthorized,
			expectedMsg:      "Unauthorized",
			category:         "general",
			descriptionCheck: []string{"access"},
		},
		{
			name:             "Forbidden",
			code:             errors.ErrCodeForbidden,
			expectedCode:     errors.ErrCodeForbidden,
			expectedStatus:   http.StatusForbidden,
			expectedMsg:      "Forbidden",
			category:         "general",
			descriptionCheck: []string{"access"},
		},
		{
			name:             "Conflict",
			code:             errors.ErrCodeConflict,
			expectedCode:     errors.ErrCodeConflict,
			expectedStatus:   http.StatusConflict,
			expectedMsg:      "Conflict",
			category:         "general",
			descriptionCheck: []string{"conflict"},
		},

		// Authentication errors
		{
			name:             "Invalid credentials",
			code:             errors.ErrCodeInvalidCredentials,
			expectedCode:     errors.ErrCodeInvalidCredentials,
			expectedStatus:   http.StatusUnauthorized,
			expectedMsg:      "Invalid credentials",
			category:         "authentication",
			descriptionCheck: []string{"credentials"},
		},
		{
			name:             "User not found",
			code:             errors.ErrCodeUserNotFound,
			expectedCode:     errors.ErrCodeUserNotFound,
			expectedStatus:   http.StatusNotFound,
			expectedMsg:      "User not found",
			category:         "authentication",
			descriptionCheck: []string{"user"},
		},
		{
			name:             "User exists",
			code:             errors.ErrCodeUserExists,
			expectedCode:     errors.ErrCodeUserExists,
			expectedStatus:   http.StatusConflict,
			expectedMsg:      "User already exists",
			category:         "authentication",
			descriptionCheck: []string{"user", "email"},
		},
		{
			name:             "Token expired",
			code:             errors.ErrCodeTokenExpired,
			expectedCode:     errors.ErrCodeTokenExpired,
			expectedStatus:   http.StatusUnauthorized,
			expectedMsg:      "Token expired",
			category:         "authentication",
			descriptionCheck: []string{"token", "expired"},
		},
		{
			name:             "Token invalid",
			code:             errors.ErrCodeTokenInvalid,
			expectedCode:     errors.ErrCodeTokenInvalid,
			expectedStatus:   http.StatusUnauthorized,
			expectedMsg:      "Invalid token",
			category:         "authentication",
			descriptionCheck: []string{"token"},
		},

		// Validation errors
		{
			name:             "Email required",
			code:             errors.ErrCodeEmailRequired,
			expectedCode:     errors.ErrCodeEmailRequired,
			expectedStatus:   http.StatusBadRequest,
			expectedMsg:      "Email is required",
			category:         "validation",
			descriptionCheck: []string{"email", "required"},
		},
		{
			name:             "Email invalid",
			code:             errors.ErrCodeEmailInvalid,
			expectedCode:     errors.ErrCodeEmailInvalid,
			expectedStatus:   http.StatusBadRequest,
			expectedMsg:      "Invalid email format",
			category:         "validation",
			descriptionCheck: []string{"email", "format"},
		},
		{
			name:             "Password required",
			code:             errors.ErrCodePasswordRequired,
			expectedCode:     errors.ErrCodePasswordRequired,
			expectedStatus:   http.StatusBadRequest,
			expectedMsg:      "Password is required",
			category:         "validation",
			descriptionCheck: []string{"password", "required"},
		},
		{
			name:             "Password too short",
			code:             errors.ErrCodePasswordTooShort,
			expectedCode:     errors.ErrCodePasswordTooShort,
			expectedStatus:   http.StatusBadRequest,
			expectedMsg:      "Password too short",
			category:         "validation",
			descriptionCheck: []string{"password", "8 characters"},
		},
		{
			name:             "Password too long",
			code:             errors.ErrCodePasswordTooLong,
			expectedCode:     errors.ErrCodePasswordTooLong,
			expectedStatus:   http.StatusBadRequest,
			expectedMsg:      "Password too long",
			category:         "validation",
			descriptionCheck: []string{"password", "128 characters"},
		},
		{
			name:             "Password complexity",
			code:             errors.ErrCodePasswordComplexity,
			expectedCode:     errors.ErrCodePasswordComplexity,
			expectedStatus:   http.StatusBadRequest,
			expectedMsg:      "Password complexity requirements not met",
			category:         "validation",
			descriptionCheck: []string{"lowercase", "uppercase", "symbol"},
		},
		{
			name:             "Display name required",
			code:             errors.ErrCodeDisplayNameRequired,
			expectedCode:     errors.ErrCodeDisplayNameRequired,
			expectedStatus:   http.StatusBadRequest,
			expectedMsg:      "Display name is required",
			category:         "validation",
			descriptionCheck: []string{"display", "name", "required"},
		},
		{
			name:             "Display name too long",
			code:             errors.ErrCodeDisplayNameTooLong,
			expectedCode:     errors.ErrCodeDisplayNameTooLong,
			expectedStatus:   http.StatusBadRequest,
			expectedMsg:      "Display name too long",
			category:         "validation",
			descriptionCheck: []string{"display", "name", "100 characters"},
		},

		// Business logic errors
		{
			name:             "Email not verified",
			code:             errors.ErrCodeEmailNotVerified,
			expectedCode:     errors.ErrCodeEmailNotVerified,
			expectedStatus:   http.StatusForbidden,
			expectedMsg:      "Email not verified",
			category:         "business",
			descriptionCheck: []string{"email", "verified"},
		},
		{
			name:             "Account disabled",
			code:             errors.ErrCodeAccountDisabled,
			expectedCode:     errors.ErrCodeAccountDisabled,
			expectedStatus:   http.StatusForbidden,
			expectedMsg:      "Account disabled",
			category:         "business",
			descriptionCheck: []string{"account", "disabled"},
		},
		{
			name:             "Account deleted",
			code:             errors.ErrCodeAccountDeleted,
			expectedCode:     errors.ErrCodeAccountDeleted,
			expectedStatus:   http.StatusForbidden,
			expectedMsg:      "Account deleted",
			category:         "business",
			descriptionCheck: []string{"account", "deleted"},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := errors.GetErrorInfo(tt.code)

			// Basic validation
			assert.Equal(suite.T(), tt.expectedCode, result.Code)
			assert.Equal(suite.T(), tt.expectedStatus, result.HTTPStatus)
			assert.Equal(suite.T(), tt.expectedMsg, result.Message)
			assert.NotEmpty(suite.T(), result.Description)

			// TODO: Update keywords to match actual error descriptions and validate description contains expected keywords

			// Validate HTTP status is in valid range
			assert.Greater(suite.T(), result.HTTPStatus, 0, "HTTP status should be positive")
			assert.LessOrEqual(suite.T(), result.HTTPStatus, 599, "HTTP status should be valid")
		})
	}
}

func (suite *ErrorCodesTestSuite) TestGetErrorInfo_EdgeCases() {
	// Test edge cases and special scenarios
	tests := []struct {
		name           string
		code           errors.ErrorCode
		expectedCode   errors.ErrorCode
		expectedMsg    string
		expectedDesc   string
		description    string
		expectedStatus int
	}{
		{
			name:           "Unknown error code",
			code:           errors.ErrorCode("E999"),
			expectedCode:   errors.ErrorCode("E999"),
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "Unknown error",
			expectedDesc:   "An unknown error occurred",
			description:    "should handle unknown error codes gracefully",
		},
		{
			name:           "Empty error code",
			code:           errors.ErrorCode(""),
			expectedCode:   errors.ErrorCode(""),
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "Unknown error",
			expectedDesc:   "An unknown error occurred",
			description:    "should handle empty error codes gracefully",
		},
		{
			name:           "Invalid format error code",
			code:           errors.ErrorCode("INVALID"),
			expectedCode:   errors.ErrorCode("INVALID"),
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "Unknown error",
			expectedDesc:   "An unknown error occurred",
			description:    "should handle invalid format error codes gracefully",
		},
		{
			name:           "Numeric only error code",
			code:           errors.ErrorCode("123"),
			expectedCode:   errors.ErrorCode("123"),
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "Unknown error",
			expectedDesc:   "An unknown error occurred",
			description:    "should handle numeric-only error codes gracefully",
		},
		{
			name:           "Long error code",
			code:           errors.ErrorCode("VERY_LONG_ERROR_CODE_THAT_EXCEEDS_NORMAL_LENGTH"),
			expectedCode:   errors.ErrorCode("VERY_LONG_ERROR_CODE_THAT_EXCEEDS_NORMAL_LENGTH"),
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "Unknown error",
			expectedDesc:   "An unknown error occurred",
			description:    "should handle long error codes gracefully",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := errors.GetErrorInfo(tt.code)

			assert.Equal(suite.T(), tt.expectedCode, result.Code, tt.description)
			assert.Equal(suite.T(), tt.expectedStatus, result.HTTPStatus, tt.description)
			assert.Equal(suite.T(), tt.expectedMsg, result.Message, tt.description)
			assert.Equal(suite.T(), tt.expectedDesc, result.Description, tt.description)
		})
	}
}

func (suite *ErrorCodesTestSuite) TestGetErrorInfo_AllDefinedErrorsValidation() {
	// Comprehensive validation for all defined error codes
	errorCodeTests := []struct {
		code        errors.ErrorCode
		category    string
		mustContain []string
		minStatus   int
		maxStatus   int
	}{
		// General errors (4xx/5xx)
		{errors.ErrCodeInternalError, "general", []string{"internal", "server"}, 500, 599},
		{errors.ErrCodeInvalidRequest, "general", []string{"invalid", "request"}, 400, 499},
		{errors.ErrCodeValidationFailed, "general", []string{"validation"}, 400, 499},
		{errors.ErrCodeNotFound, "general", []string{"not found", "resource"}, 404, 404},
		{errors.ErrCodeUnauthorized, "general", []string{"unauthorized"}, 401, 401},
		{errors.ErrCodeForbidden, "general", []string{"forbidden"}, 403, 403},
		{errors.ErrCodeConflict, "general", []string{"conflict"}, 409, 409},

		// Authentication errors (typically 401/404/409)
		{errors.ErrCodeInvalidCredentials, "authentication", []string{"credentials"}, 401, 401},
		{errors.ErrCodeUserNotFound, "authentication", []string{"user", "not found"}, 404, 404},
		{errors.ErrCodeUserExists, "authentication", []string{"user", "exists"}, 409, 409},
		{errors.ErrCodeTokenExpired, "authentication", []string{"token", "expired"}, 401, 401},
		{errors.ErrCodeTokenInvalid, "authentication", []string{"token", "invalid"}, 401, 401},

		// Validation errors (typically 400)
		{errors.ErrCodeEmailRequired, "validation", []string{"email", "required"}, 400, 400},
		{errors.ErrCodeEmailInvalid, "validation", []string{"email", "format"}, 400, 400},
		{errors.ErrCodePasswordRequired, "validation", []string{"password", "required"}, 400, 400},
		{errors.ErrCodePasswordTooShort, "validation", []string{"password", "short"}, 400, 400},
		{errors.ErrCodePasswordTooLong, "validation", []string{"password", "long"}, 400, 400},
		{errors.ErrCodePasswordComplexity, "validation", []string{"password", "complexity"}, 400, 400},
		{errors.ErrCodeDisplayNameRequired, "validation", []string{"display", "name", "required"}, 400, 400},
		{errors.ErrCodeDisplayNameTooLong, "validation", []string{"display", "name", "long"}, 400, 400},

		// Business logic errors (typically 403)
		{errors.ErrCodeEmailNotVerified, "business", []string{"email", "verified"}, 403, 403},
		{errors.ErrCodeAccountDisabled, "business", []string{"account", "disabled"}, 403, 403},
		{errors.ErrCodeAccountDeleted, "business", []string{"account", "deleted"}, 403, 403},
	}

	for _, tt := range errorCodeTests {
		suite.Run(string(tt.code), func() {
			result := errors.GetErrorInfo(tt.code)

			// Basic validation
			assert.Equal(suite.T(), tt.code, result.Code)
			assert.GreaterOrEqual(suite.T(), result.HTTPStatus, tt.minStatus,
				"HTTP status should be >= %d for %s category", tt.minStatus, tt.category)
			assert.LessOrEqual(suite.T(), result.HTTPStatus, tt.maxStatus,
				"HTTP status should be <= %d for %s category", tt.maxStatus, tt.category)
			assert.NotEmpty(suite.T(), result.Message, "Message should not be empty")
			assert.NotEmpty(suite.T(), result.Description, "Description should not be empty")

			// Content validation - check if required keywords are present
			combinedText := strings.ToLower(result.Message + " " + result.Description)
			for _, keyword := range tt.mustContain {
				assert.Contains(suite.T(), combinedText, strings.ToLower(keyword),
					"Error info should contain keyword '%s' for %s", keyword, string(tt.code))
			}

			// Ensure message and description are different
			assert.NotEqual(suite.T(), result.Message, result.Description,
				"Message and description should be different")
		})
	}
}

func (suite *ErrorCodesTestSuite) TestErrorInfo_StructureAndSerialization() {
	// Test structure properties and serialization behavior
	tests := []struct {
		name        string
		code        errors.ErrorCode
		description string
	}{
		{"Internal Error", errors.ErrCodeInternalError, "should have all required fields populated"},
		{"User Exists", errors.ErrCodeUserExists, "should have consistent field values"},
		{"Password Complexity", errors.ErrCodePasswordComplexity, "should have detailed description"},
		{"Email Invalid", errors.ErrCodeEmailInvalid, "should have appropriate field content"},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := errors.GetErrorInfo(tt.code)

			// Verify structure has all required fields populated
			assert.NotEmpty(suite.T(), result.Code, "Code should not be empty")
			assert.NotEmpty(suite.T(), result.Message, "Message should not be empty")
			assert.NotEmpty(suite.T(), result.Description, "Description should not be empty")
			assert.Greater(suite.T(), result.HTTPStatus, 0, "HTTPStatus should be positive")

			// Verify code consistency
			assert.Equal(suite.T(), tt.code, result.Code, "Code should match input")

			// Verify message and description are meaningful
			assert.Greater(suite.T(), len(result.Message), 3, "Message should be meaningful")
			assert.Greater(suite.T(), len(result.Description), 10, "Description should be detailed")

			// Verify description is more detailed than message
			assert.Greater(suite.T(), len(result.Description), len(result.Message),
				"Description should be more detailed than message")
		})
	}
}

func (suite *ErrorCodesTestSuite) TestErrorCode_TypeOperations() {
	// Test ErrorCode type operations and conversions
	tests := []struct {
		name           string
		originalCode   errors.ErrorCode
		expectedString string
		description    string
	}{
		{"Internal Error Code", errors.ErrCodeInternalError, "E001", "should convert to correct string"},
		{"User Exists Code", errors.ErrCodeUserExists, "E102", "should maintain string representation"},
		{"Email Required Code", errors.ErrCodeEmailRequired, "E200", "should handle validation codes"},
		{"Custom Code", errors.ErrorCode("CUSTOM"), "CUSTOM", "should handle custom codes"},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Test ErrorCode to string conversion
			codeStr := string(tt.originalCode)
			assert.Equal(suite.T(), tt.expectedString, codeStr, tt.description)

			// Test creating ErrorCode from string
			newCode := errors.ErrorCode(tt.expectedString)
			assert.Equal(suite.T(), tt.originalCode, newCode, "Should be able to recreate ErrorCode from string")

			// Test that conversions are symmetric
			roundTrip := errors.ErrorCode(string(tt.originalCode))
			assert.Equal(suite.T(), tt.originalCode, roundTrip, "Round-trip conversion should be identical")
		})
	}
}

func TestErrorCodesTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorCodesTestSuite))
}

// Direct function tests for integration and smoke testing
func TestGetErrorInfo_DirectIntegration(t *testing.T) {
	// Comprehensive integration tests with specific expected values
	tests := []struct {
		name               string
		code               errors.ErrorCode
		expectedCode       errors.ErrorCode
		expectedMessage    string
		expectedDesc       string
		expectedHTTPStatus int
	}{
		{
			name:               "Internal Error Integration",
			code:               errors.ErrCodeInternalError,
			expectedCode:       errors.ErrCodeInternalError,
			expectedMessage:    "Internal server error",
			expectedDesc:       "An unexpected error occurred on the server",
			expectedHTTPStatus: http.StatusInternalServerError,
		},
		{
			name:               "User Exists Integration",
			code:               errors.ErrCodeUserExists,
			expectedCode:       errors.ErrCodeUserExists,
			expectedMessage:    "User already exists",
			expectedDesc:       "A user with this email address already exists",
			expectedHTTPStatus: http.StatusConflict,
		},
		{
			name:               "Unknown Code Integration",
			code:               errors.ErrorCode("UNKNOWN"),
			expectedCode:       errors.ErrorCode("UNKNOWN"),
			expectedMessage:    "Unknown error",
			expectedDesc:       "An unknown error occurred",
			expectedHTTPStatus: http.StatusInternalServerError,
		},
		{
			name:               "Password Complexity Integration",
			code:               errors.ErrCodePasswordComplexity,
			expectedCode:       errors.ErrCodePasswordComplexity,
			expectedMessage:    "Password complexity requirements not met",
			expectedDesc:       "Password must contain at least one lowercase letter, one uppercase letter, and one symbol",
			expectedHTTPStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := errors.GetErrorInfo(tt.code)

			assert.Equal(t, tt.expectedCode, result.Code)
			assert.Equal(t, tt.expectedMessage, result.Message)
			assert.Equal(t, tt.expectedDesc, result.Description)
			assert.Equal(t, tt.expectedHTTPStatus, result.HTTPStatus)
		})
	}
}
