package errors

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ErrorCodesTestSuite struct {
	suite.Suite
}

func (suite *ErrorCodesTestSuite) TestGetErrorInfoValidCodes() {
	testCases := []struct {
		code           ErrorCode
		expectedCode   ErrorCode
		expectedMsg    string
		expectedStatus int
	}{
		// General errors
		{ErrCodeInternalError, ErrCodeInternalError, http.StatusInternalServerError, "Internal server error"},
		{ErrCodeInvalidRequest, ErrCodeInvalidRequest, http.StatusBadRequest, "Invalid request"},
		{ErrCodeValidationFailed, ErrCodeValidationFailed, http.StatusBadRequest, "Validation failed"},
		{ErrCodeNotFound, ErrCodeNotFound, http.StatusNotFound, "Resource not found"},
		{ErrCodeUnauthorized, ErrCodeUnauthorized, http.StatusUnauthorized, "Unauthorized"},
		{ErrCodeForbidden, ErrCodeForbidden, http.StatusForbidden, "Forbidden"},
		{ErrCodeConflict, ErrCodeConflict, http.StatusConflict, "Conflict"},

		// Authentication errors
		{ErrCodeInvalidCredentials, ErrCodeInvalidCredentials, http.StatusUnauthorized, "Invalid credentials"},
		{ErrCodeUserNotFound, ErrCodeUserNotFound, http.StatusNotFound, "User not found"},
		{ErrCodeUserExists, ErrCodeUserExists, http.StatusConflict, "User already exists"},
		{ErrCodeTokenExpired, ErrCodeTokenExpired, http.StatusUnauthorized, "Token expired"},
		{ErrCodeTokenInvalid, ErrCodeTokenInvalid, http.StatusUnauthorized, "Invalid token"},

		// Validation errors
		{ErrCodeEmailRequired, ErrCodeEmailRequired, http.StatusBadRequest, "Email is required"},
		{ErrCodeEmailInvalid, ErrCodeEmailInvalid, http.StatusBadRequest, "Invalid email format"},
		{ErrCodePasswordRequired, ErrCodePasswordRequired, http.StatusBadRequest, "Password is required"},
		{ErrCodePasswordTooShort, ErrCodePasswordTooShort, http.StatusBadRequest, "Password too short"},
		{ErrCodePasswordTooLong, ErrCodePasswordTooLong, http.StatusBadRequest, "Password too long"},
		{ErrCodePasswordComplexity, ErrCodePasswordComplexity, http.StatusBadRequest, "Password complexity requirements not met"},
		{ErrCodeDisplayNameRequired, ErrCodeDisplayNameRequired, http.StatusBadRequest, "Display name is required"},
		{ErrCodeDisplayNameTooLong, ErrCodeDisplayNameTooLong, http.StatusBadRequest, "Display name too long"},

		// Business logic errors
		{ErrCodeEmailNotVerified, ErrCodeEmailNotVerified, http.StatusForbidden, "Email not verified"},
		{ErrCodeAccountDisabled, ErrCodeAccountDisabled, http.StatusForbidden, "Account disabled"},
		{ErrCodeAccountDeleted, ErrCodeAccountDeleted, http.StatusForbidden, "Account deleted"},
	}

	for _, tc := range testCases {
		suite.T().Run(string(tc.code), func(t *testing.T) {
			errorInfo := GetErrorInfo(tc.code)

			assert.Equal(t, tc.expectedCode, errorInfo.Code)
			assert.Equal(t, tc.expectedStatus, errorInfo.HTTPStatus)
			assert.Equal(t, tc.expectedMsg, errorInfo.Message)
			assert.NotEmpty(t, errorInfo.Description)
		})
	}
}

func (suite *ErrorCodesTestSuite) TestGetErrorInfoUnknownCode() {
	unknownCode := ErrorCode("E999")
	errorInfo := GetErrorInfo(unknownCode)

	assert.Equal(suite.T(), ErrCodeInternalError, errorInfo.Code)
	assert.Equal(suite.T(), http.StatusInternalServerError, errorInfo.HTTPStatus)
	assert.Equal(suite.T(), "Unknown error", errorInfo.Message)
	assert.Equal(suite.T(), "An unknown error occurred", errorInfo.Description)
}

func (suite *ErrorCodesTestSuite) TestGetErrorInfoEmptyCode() {
	emptyCode := ErrorCode("")
	errorInfo := GetErrorInfo(emptyCode)

	assert.Equal(suite.T(), ErrCodeInternalError, errorInfo.Code)
	assert.Equal(suite.T(), http.StatusInternalServerError, errorInfo.HTTPStatus)
	assert.Equal(suite.T(), "Unknown error", errorInfo.Message)
	assert.Equal(suite.T(), "An unknown error occurred", errorInfo.Description)
}

func (suite *ErrorCodesTestSuite) TestErrorCodeConstants() {
	// Test that error codes follow the expected format
	testCases := []struct {
		code     ErrorCode
		expected string
	}{
		{ErrCodeInternalError, "E001"},
		{ErrCodeInvalidRequest, "E002"},
		{ErrCodeValidationFailed, "E003"},
		{ErrCodeInvalidCredentials, "E100"},
		{ErrCodeUserExists, "E102"},
		{ErrCodePasswordComplexity, "E205"},
		{ErrCodeEmailNotVerified, "E300"},
	}

	for _, tc := range testCases {
		assert.Equal(suite.T(), tc.expected, string(tc.code))
	}
}

func (suite *ErrorCodesTestSuite) TestErrorInfoStructure() {
	errorInfo := GetErrorInfo(ErrCodeInvalidRequest)

	// Test that all required fields are present
	assert.NotEmpty(suite.T(), errorInfo.Code)
	assert.NotEmpty(suite.T(), errorInfo.Message)
	assert.NotEmpty(suite.T(), errorInfo.Description)
	assert.NotZero(suite.T(), errorInfo.HTTPStatus)

	// Test that HTTPStatus is a valid HTTP status code
	validStatuses := []int{
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusConflict,
		http.StatusInternalServerError,
	}

	found := false
	for _, status := range validStatuses {
		if errorInfo.HTTPStatus == status {
			found = true
			break
		}
	}
	assert.True(suite.T(), found, "HTTPStatus should be a valid HTTP status code")
}

func (suite *ErrorCodesTestSuite) TestHTTPStatusMapping() {
	// Test that HTTP status codes are mapped correctly for error categories
	testCases := []struct {
		code           ErrorCode
		category       string
		expectedStatus int
	}{
		// 400 errors
		{ErrCodeInvalidRequest, http.StatusBadRequest, "client error"},
		{ErrCodeValidationFailed, http.StatusBadRequest, "client error"},
		{ErrCodeEmailInvalid, http.StatusBadRequest, "client error"},
		{ErrCodePasswordTooShort, http.StatusBadRequest, "client error"},

		// 401 errors
		{ErrCodeUnauthorized, http.StatusUnauthorized, "authentication error"},
		{ErrCodeInvalidCredentials, http.StatusUnauthorized, "authentication error"},
		{ErrCodeTokenExpired, http.StatusUnauthorized, "authentication error"},

		// 403 errors
		{ErrCodeForbidden, http.StatusForbidden, "authorization error"},
		{ErrCodeEmailNotVerified, http.StatusForbidden, "authorization error"},
		{ErrCodeAccountDisabled, http.StatusForbidden, "authorization error"},

		// 404 errors
		{ErrCodeNotFound, http.StatusNotFound, "not found error"},
		{ErrCodeUserNotFound, http.StatusNotFound, "not found error"},

		// 409 errors
		{ErrCodeConflict, http.StatusConflict, "conflict error"},
		{ErrCodeUserExists, http.StatusConflict, "conflict error"},

		// 500 errors
		{ErrCodeInternalError, http.StatusInternalServerError, "server error"},
	}

	for _, tc := range testCases {
		suite.T().Run(string(tc.code)+"_"+tc.category, func(t *testing.T) {
			errorInfo := GetErrorInfo(tc.code)
			assert.Equal(t, tc.expectedStatus, errorInfo.HTTPStatus,
				"Error code %s should have HTTP status %d", tc.code, tc.expectedStatus)
		})
	}
}

func (suite *ErrorCodesTestSuite) TestErrorCodeRanges() {
	// Test that error codes are in the expected ranges
	testCases := []struct {
		code  ErrorCode
		start string
		end   string
		desc  string
	}{
		{ErrCodeInternalError, "E001", "E099", "General errors"},
		{ErrCodeInvalidRequest, "E001", "E099", "General errors"},
		{ErrCodeInvalidCredentials, "E100", "E199", "Authentication errors"},
		{ErrCodeUserExists, "E100", "E199", "Authentication errors"},
		{ErrCodeEmailRequired, "E200", "E299", "Validation errors"},
		{ErrCodePasswordComplexity, "E200", "E299", "Validation errors"},
		{ErrCodeEmailNotVerified, "E300", "E399", "Business logic errors"},
	}

	for _, tc := range testCases {
		suite.T().Run(string(tc.code)+"_range", func(t *testing.T) {
			codeStr := string(tc.code)
			assert.GreaterOrEqual(t, codeStr, tc.start,
				"Code %s should be >= %s for %s", codeStr, tc.start, tc.desc)
			assert.LessOrEqual(t, codeStr, tc.end,
				"Code %s should be <= %s for %s", codeStr, tc.end, tc.desc)
		})
	}
}

func (suite *ErrorCodesTestSuite) TestErrorInfoConsistency() {
	// Test that all error codes have consistent structure
	allCodes := []ErrorCode{
		ErrCodeInternalError, ErrCodeInvalidRequest, ErrCodeValidationFailed,
		ErrCodeNotFound, ErrCodeUnauthorized, ErrCodeForbidden, ErrCodeConflict,
		ErrCodeInvalidCredentials, ErrCodeUserNotFound, ErrCodeUserExists,
		ErrCodeTokenExpired, ErrCodeTokenInvalid,
		ErrCodeEmailRequired, ErrCodeEmailInvalid, ErrCodePasswordRequired,
		ErrCodePasswordTooShort, ErrCodePasswordTooLong, ErrCodePasswordComplexity,
		ErrCodeDisplayNameRequired, ErrCodeDisplayNameTooLong,
		ErrCodeEmailNotVerified, ErrCodeAccountDisabled, ErrCodeAccountDeleted,
	}

	for _, code := range allCodes {
		suite.T().Run(string(code)+"_consistency", func(t *testing.T) {
			errorInfo := GetErrorInfo(code)

			// All error info should have these properties
			assert.Equal(t, code, errorInfo.Code, "Code field should match input")
			assert.NotEmpty(t, errorInfo.Message, "Message should not be empty")
			assert.NotEmpty(t, errorInfo.Description, "Description should not be empty")
			assert.Greater(t, errorInfo.HTTPStatus, 0, "HTTPStatus should be positive")
			assert.Less(t, errorInfo.HTTPStatus, 600, "HTTPStatus should be valid HTTP status")

			// Error codes should follow E### format
			codeStr := string(code)
			assert.Regexp(t, `^E\d{3}$`, codeStr, "Error code should follow E### format")
		})
	}
}

func (suite *ErrorCodesTestSuite) TestPasswordComplexityErrorDetails() {
	// Specific test for the new password complexity error
	errorInfo := GetErrorInfo(ErrCodePasswordComplexity)

	assert.Equal(suite.T(), ErrCodePasswordComplexity, errorInfo.Code)
	assert.Equal(suite.T(), "E205", string(errorInfo.Code))
	assert.Equal(suite.T(), http.StatusBadRequest, errorInfo.HTTPStatus)
	assert.Equal(suite.T(), "Password complexity requirements not met", errorInfo.Message)
	assert.Contains(suite.T(), errorInfo.Description, "lowercase")
	assert.Contains(suite.T(), errorInfo.Description, "uppercase")
	assert.Contains(suite.T(), errorInfo.Description, "symbol")
}

func TestErrorCodesTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorCodesTestSuite))
}
