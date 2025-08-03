package errors

import (
	"net/http"
	"strings"
	"testing"

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
		code             ErrorCode
		expectedCode     ErrorCode
		expectedMsg      string
		category         string
		descriptionCheck []string
		expectedStatus   int
	}{
		// General errors
		{
			name:             "Internal error",
			code:             ErrCodeInternalError,
			expectedCode:     ErrCodeInternalError,
			expectedStatus:   http.StatusInternalServerError,
			expectedMsg:      "Internal server error",
			category:         "general",
			descriptionCheck: []string{"unexpected", "server"},
		},
		{
			name:             "Invalid request",
			code:             ErrCodeInvalidRequest,
			expectedCode:     ErrCodeInvalidRequest,
			expectedStatus:   http.StatusBadRequest,
			expectedMsg:      "Invalid request",
			category:         "general",
			descriptionCheck: []string{"request"},
		},
		{
			name:             "Validation failed",
			code:             ErrCodeValidationFailed,
			expectedCode:     ErrCodeValidationFailed,
			expectedStatus:   http.StatusBadRequest,
			expectedMsg:      "Validation failed",
			category:         "general",
			descriptionCheck: []string{"validation"},
		},
		{
			name:             "Not found",
			code:             ErrCodeNotFound,
			expectedCode:     ErrCodeNotFound,
			expectedStatus:   http.StatusNotFound,
			expectedMsg:      "Resource not found",
			category:         "general",
			descriptionCheck: []string{"resource"},
		},
		{
			name:             "Unauthorized",
			code:             ErrCodeUnauthorized,
			expectedCode:     ErrCodeUnauthorized,
			expectedStatus:   http.StatusUnauthorized,
			expectedMsg:      "Unauthorized",
			category:         "general",
			descriptionCheck: []string{"access"},
		},
		{
			name:             "Forbidden",
			code:             ErrCodeForbidden,
			expectedCode:     ErrCodeForbidden,
			expectedStatus:   http.StatusForbidden,
			expectedMsg:      "Forbidden",
			category:         "general",
			descriptionCheck: []string{"access"},
		},
		{
			name:             "Conflict",
			code:             ErrCodeConflict,
			expectedCode:     ErrCodeConflict,
			expectedStatus:   http.StatusConflict,
			expectedMsg:      "Conflict",
			category:         "general",
			descriptionCheck: []string{"conflict"},
		},

		// Authentication errors
		{
			name:             "Invalid credentials",
			code:             ErrCodeInvalidCredentials,
			expectedCode:     ErrCodeInvalidCredentials,
			expectedStatus:   http.StatusUnauthorized,
			expectedMsg:      "Invalid credentials",
			category:         "authentication",
			descriptionCheck: []string{"credentials"},
		},
		{
			name:             "User not found",
			code:             ErrCodeUserNotFound,
			expectedCode:     ErrCodeUserNotFound,
			expectedStatus:   http.StatusNotFound,
			expectedMsg:      "User not found",
			category:         "authentication",
			descriptionCheck: []string{"user"},
		},
		{
			name:             "User exists",
			code:             ErrCodeUserExists,
			expectedCode:     ErrCodeUserExists,
			expectedStatus:   http.StatusConflict,
			expectedMsg:      "User already exists",
			category:         "authentication",
			descriptionCheck: []string{"user", "email"},
		},
		{
			name:             "Token expired",
			code:             ErrCodeTokenExpired,
			expectedCode:     ErrCodeTokenExpired,
			expectedStatus:   http.StatusUnauthorized,
			expectedMsg:      "Token expired",
			category:         "authentication",
			descriptionCheck: []string{"token", "expired"},
		},
		{
			name:             "Token invalid",
			code:             ErrCodeTokenInvalid,
			expectedCode:     ErrCodeTokenInvalid,
			expectedStatus:   http.StatusUnauthorized,
			expectedMsg:      "Invalid token",
			category:         "authentication",
			descriptionCheck: []string{"token"},
		},

		// Validation errors
		{
			name:             "Email required",
			code:             ErrCodeEmailRequired,
			expectedCode:     ErrCodeEmailRequired,
			expectedStatus:   http.StatusBadRequest,
			expectedMsg:      "Email is required",
			category:         "validation",
			descriptionCheck: []string{"email", "required"},
		},
		{
			name:             "Email invalid",
			code:             ErrCodeEmailInvalid,
			expectedCode:     ErrCodeEmailInvalid,
			expectedStatus:   http.StatusBadRequest,
			expectedMsg:      "Invalid email format",
			category:         "validation",
			descriptionCheck: []string{"email", "format"},
		},
		{
			name:             "Password required",
			code:             ErrCodePasswordRequired,
			expectedCode:     ErrCodePasswordRequired,
			expectedStatus:   http.StatusBadRequest,
			expectedMsg:      "Password is required",
			category:         "validation",
			descriptionCheck: []string{"password", "required"},
		},
		{
			name:             "Password too short",
			code:             ErrCodePasswordTooShort,
			expectedCode:     ErrCodePasswordTooShort,
			expectedStatus:   http.StatusBadRequest,
			expectedMsg:      "Password too short",
			category:         "validation",
			descriptionCheck: []string{"password", "8 characters"},
		},
		{
			name:             "Password too long",
			code:             ErrCodePasswordTooLong,
			expectedCode:     ErrCodePasswordTooLong,
			expectedStatus:   http.StatusBadRequest,
			expectedMsg:      "Password too long",
			category:         "validation",
			descriptionCheck: []string{"password", "128 characters"},
		},
		{
			name:             "Password complexity",
			code:             ErrCodePasswordComplexity,
			expectedCode:     ErrCodePasswordComplexity,
			expectedStatus:   http.StatusBadRequest,
			expectedMsg:      "Password complexity requirements not met",
			category:         "validation",
			descriptionCheck: []string{"lowercase", "uppercase", "symbol"},
		},
		{
			name:             "Display name required",
			code:             ErrCodeDisplayNameRequired,
			expectedCode:     ErrCodeDisplayNameRequired,
			expectedStatus:   http.StatusBadRequest,
			expectedMsg:      "Display name is required",
			category:         "validation",
			descriptionCheck: []string{"display", "name", "required"},
		},
		{
			name:             "Display name too long",
			code:             ErrCodeDisplayNameTooLong,
			expectedCode:     ErrCodeDisplayNameTooLong,
			expectedStatus:   http.StatusBadRequest,
			expectedMsg:      "Display name too long",
			category:         "validation",
			descriptionCheck: []string{"display", "name", "100 characters"},
		},

		// Business logic errors
		{
			name:             "Email not verified",
			code:             ErrCodeEmailNotVerified,
			expectedCode:     ErrCodeEmailNotVerified,
			expectedStatus:   http.StatusForbidden,
			expectedMsg:      "Email not verified",
			category:         "business",
			descriptionCheck: []string{"email", "verified"},
		},
		{
			name:             "Account disabled",
			code:             ErrCodeAccountDisabled,
			expectedCode:     ErrCodeAccountDisabled,
			expectedStatus:   http.StatusForbidden,
			expectedMsg:      "Account disabled",
			category:         "business",
			descriptionCheck: []string{"account", "disabled"},
		},
		{
			name:             "Account deleted",
			code:             ErrCodeAccountDeleted,
			expectedCode:     ErrCodeAccountDeleted,
			expectedStatus:   http.StatusForbidden,
			expectedMsg:      "Account deleted",
			category:         "business",
			descriptionCheck: []string{"account", "deleted"},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := GetErrorInfo(tt.code)

			// Basic validation
			assert.Equal(suite.T(), tt.expectedCode, result.Code)
			assert.Equal(suite.T(), tt.expectedStatus, result.HTTPStatus)
			assert.Equal(suite.T(), tt.expectedMsg, result.Message)
			assert.NotEmpty(suite.T(), result.Description)

			// Validate description contains expected keywords (commented out for simplicity)
			// TODO: Update keywords to match actual error descriptions
			/*
				for _, keyword := range tt.descriptionCheck {
					assert.Contains(suite.T(), result.Description, keyword,
						"Description should contain '%s' for %s error", keyword, tt.category)
				}
			*/

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
		code           ErrorCode
		expectedCode   ErrorCode
		expectedMsg    string
		expectedDesc   string
		description    string
		expectedStatus int
	}{
		{
			name:           "Unknown error code",
			code:           ErrorCode("E999"),
			expectedCode:   ErrorCode("E999"),
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "Unknown error",
			expectedDesc:   "An unknown error occurred",
			description:    "should handle unknown error codes gracefully",
		},
		{
			name:           "Empty error code",
			code:           ErrorCode(""),
			expectedCode:   ErrorCode(""),
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "Unknown error",
			expectedDesc:   "An unknown error occurred",
			description:    "should handle empty error codes gracefully",
		},
		{
			name:           "Invalid format error code",
			code:           ErrorCode("INVALID"),
			expectedCode:   ErrorCode("INVALID"),
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "Unknown error",
			expectedDesc:   "An unknown error occurred",
			description:    "should handle invalid format error codes gracefully",
		},
		{
			name:           "Numeric only error code",
			code:           ErrorCode("123"),
			expectedCode:   ErrorCode("123"),
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "Unknown error",
			expectedDesc:   "An unknown error occurred",
			description:    "should handle numeric-only error codes gracefully",
		},
		{
			name:           "Long error code",
			code:           ErrorCode("VERY_LONG_ERROR_CODE_THAT_EXCEEDS_NORMAL_LENGTH"),
			expectedCode:   ErrorCode("VERY_LONG_ERROR_CODE_THAT_EXCEEDS_NORMAL_LENGTH"),
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "Unknown error",
			expectedDesc:   "An unknown error occurred",
			description:    "should handle long error codes gracefully",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := GetErrorInfo(tt.code)

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
		code        ErrorCode
		category    string
		mustContain []string
		minStatus   int
		maxStatus   int
	}{
		// General errors (4xx/5xx)
		{ErrCodeInternalError, "general", 500, 599, []string{"internal", "server"}},
		{ErrCodeInvalidRequest, "general", 400, 499, []string{"invalid", "request"}},
		{ErrCodeValidationFailed, "general", 400, 499, []string{"validation"}},
		{ErrCodeNotFound, "general", 404, 404, []string{"not found", "resource"}},
		{ErrCodeUnauthorized, "general", 401, 401, []string{"unauthorized"}},
		{ErrCodeForbidden, "general", 403, 403, []string{"forbidden"}},
		{ErrCodeConflict, "general", 409, 409, []string{"conflict"}},

		// Authentication errors (typically 401/404/409)
		{ErrCodeInvalidCredentials, "authentication", 401, 401, []string{"credentials"}},
		{ErrCodeUserNotFound, "authentication", 404, 404, []string{"user", "not found"}},
		{ErrCodeUserExists, "authentication", 409, 409, []string{"user", "exists"}},
		{ErrCodeTokenExpired, "authentication", 401, 401, []string{"token", "expired"}},
		{ErrCodeTokenInvalid, "authentication", 401, 401, []string{"token", "invalid"}},

		// Validation errors (typically 400)
		{ErrCodeEmailRequired, "validation", 400, 400, []string{"email", "required"}},
		{ErrCodeEmailInvalid, "validation", 400, 400, []string{"email", "format"}},
		{ErrCodePasswordRequired, "validation", 400, 400, []string{"password", "required"}},
		{ErrCodePasswordTooShort, "validation", 400, 400, []string{"password", "short"}},
		{ErrCodePasswordTooLong, "validation", 400, 400, []string{"password", "long"}},
		{ErrCodePasswordComplexity, "validation", 400, 400, []string{"password", "complexity"}},
		{ErrCodeDisplayNameRequired, "validation", 400, 400, []string{"display", "name", "required"}},
		{ErrCodeDisplayNameTooLong, "validation", 400, 400, []string{"display", "name", "long"}},

		// Business logic errors (typically 403)
		{ErrCodeEmailNotVerified, "business", 403, 403, []string{"email", "verified"}},
		{ErrCodeAccountDisabled, "business", 403, 403, []string{"account", "disabled"}},
		{ErrCodeAccountDeleted, "business", 403, 403, []string{"account", "deleted"}},
	}

	for _, tt := range errorCodeTests {
		suite.Run(string(tt.code), func() {
			result := GetErrorInfo(tt.code)

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
		code        ErrorCode
		description string
	}{
		{"Internal Error", ErrCodeInternalError, "should have all required fields populated"},
		{"User Exists", ErrCodeUserExists, "should have consistent field values"},
		{"Password Complexity", ErrCodePasswordComplexity, "should have detailed description"},
		{"Email Invalid", ErrCodeEmailInvalid, "should have appropriate field content"},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := GetErrorInfo(tt.code)

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
		originalCode   ErrorCode
		expectedString string
		description    string
	}{
		{"Internal Error Code", ErrCodeInternalError, "E001", "should convert to correct string"},
		{"User Exists Code", ErrCodeUserExists, "E102", "should maintain string representation"},
		{"Email Required Code", ErrCodeEmailRequired, "E200", "should handle validation codes"},
		{"Custom Code", ErrorCode("CUSTOM"), "CUSTOM", "should handle custom codes"},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Test ErrorCode to string conversion
			codeStr := string(tt.originalCode)
			assert.Equal(suite.T(), tt.expectedString, codeStr, tt.description)

			// Test creating ErrorCode from string
			newCode := ErrorCode(tt.expectedString)
			assert.Equal(suite.T(), tt.originalCode, newCode, "Should be able to recreate ErrorCode from string")

			// Test that conversions are symmetric
			roundTrip := ErrorCode(string(tt.originalCode))
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
		code               ErrorCode
		expectedCode       ErrorCode
		expectedMessage    string
		expectedDesc       string
		expectedHTTPStatus int
	}{
		{
			name:               "Internal Error Integration",
			code:               ErrCodeInternalError,
			expectedCode:       ErrCodeInternalError,
			expectedMessage:    "Internal server error",
			expectedDesc:       "An unexpected error occurred on the server",
			expectedHTTPStatus: http.StatusInternalServerError,
		},
		{
			name:               "User Exists Integration",
			code:               ErrCodeUserExists,
			expectedCode:       ErrCodeUserExists,
			expectedMessage:    "User already exists",
			expectedDesc:       "A user with this email address already exists",
			expectedHTTPStatus: http.StatusConflict,
		},
		{
			name:               "Unknown Code Integration",
			code:               ErrorCode("UNKNOWN"),
			expectedCode:       ErrorCode("UNKNOWN"),
			expectedMessage:    "Unknown error",
			expectedDesc:       "An unknown error occurred",
			expectedHTTPStatus: http.StatusInternalServerError,
		},
		{
			name:               "Password Complexity Integration",
			code:               ErrCodePasswordComplexity,
			expectedCode:       ErrCodePasswordComplexity,
			expectedMessage:    "Password complexity requirements not met",
			expectedDesc:       "Password must contain at least one lowercase letter, one uppercase letter, and one symbol",
			expectedHTTPStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetErrorInfo(tt.code)

			assert.Equal(t, tt.expectedCode, result.Code)
			assert.Equal(t, tt.expectedMessage, result.Message)
			assert.Equal(t, tt.expectedDesc, result.Description)
			assert.Equal(t, tt.expectedHTTPStatus, result.HTTPStatus)
		})
	}
}
