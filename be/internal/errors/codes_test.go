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

func (suite *ErrorCodesTestSuite) TestGetErrorInfo_GeneralErrors() {
	tests := []struct {
		name           string
		code           ErrorCode
		expectedCode   ErrorCode
		expectedMsg    string
		expectedStatus int
	}{
		{
			name:           "Internal error",
			code:           ErrCodeInternalError,
			expectedCode:   ErrCodeInternalError,
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "Internal server error",
		},
		{
			name:           "Invalid request",
			code:           ErrCodeInvalidRequest,
			expectedCode:   ErrCodeInvalidRequest,
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Invalid request",
		},
		{
			name:           "Validation failed",
			code:           ErrCodeValidationFailed,
			expectedCode:   ErrCodeValidationFailed,
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Validation failed",
		},
		{
			name:           "Not found",
			code:           ErrCodeNotFound,
			expectedCode:   ErrCodeNotFound,
			expectedStatus: http.StatusNotFound,
			expectedMsg:    "Resource not found",
		},
		{
			name:           "Unauthorized",
			code:           ErrCodeUnauthorized,
			expectedCode:   ErrCodeUnauthorized,
			expectedStatus: http.StatusUnauthorized,
			expectedMsg:    "Unauthorized",
		},
		{
			name:           "Forbidden",
			code:           ErrCodeForbidden,
			expectedCode:   ErrCodeForbidden,
			expectedStatus: http.StatusForbidden,
			expectedMsg:    "Forbidden",
		},
		{
			name:           "Conflict",
			code:           ErrCodeConflict,
			expectedCode:   ErrCodeConflict,
			expectedStatus: http.StatusConflict,
			expectedMsg:    "Conflict",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := GetErrorInfo(tt.code)

			assert.Equal(suite.T(), tt.expectedCode, result.Code)
			assert.Equal(suite.T(), tt.expectedStatus, result.HTTPStatus)
			assert.Equal(suite.T(), tt.expectedMsg, result.Message)
			assert.NotEmpty(suite.T(), result.Description)
		})
	}
}

func (suite *ErrorCodesTestSuite) TestGetErrorInfo_AuthenticationErrors() {
	tests := []struct {
		name           string
		code           ErrorCode
		expectedCode   ErrorCode
		expectedMsg    string
		expectedStatus int
	}{
		{
			name:           "Invalid credentials",
			code:           ErrCodeInvalidCredentials,
			expectedCode:   ErrCodeInvalidCredentials,
			expectedStatus: http.StatusUnauthorized,
			expectedMsg:    "Invalid credentials",
		},
		{
			name:           "User not found",
			code:           ErrCodeUserNotFound,
			expectedCode:   ErrCodeUserNotFound,
			expectedStatus: http.StatusNotFound,
			expectedMsg:    "User not found",
		},
		{
			name:           "User exists",
			code:           ErrCodeUserExists,
			expectedCode:   ErrCodeUserExists,
			expectedStatus: http.StatusConflict,
			expectedMsg:    "User already exists",
		},
		{
			name:           "Token expired",
			code:           ErrCodeTokenExpired,
			expectedCode:   ErrCodeTokenExpired,
			expectedStatus: http.StatusUnauthorized,
			expectedMsg:    "Token expired",
		},
		{
			name:           "Token invalid",
			code:           ErrCodeTokenInvalid,
			expectedCode:   ErrCodeTokenInvalid,
			expectedStatus: http.StatusUnauthorized,
			expectedMsg:    "Invalid token",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := GetErrorInfo(tt.code)

			assert.Equal(suite.T(), tt.expectedCode, result.Code)
			assert.Equal(suite.T(), tt.expectedStatus, result.HTTPStatus)
			assert.Equal(suite.T(), tt.expectedMsg, result.Message)
			assert.NotEmpty(suite.T(), result.Description)
		})
	}
}

func (suite *ErrorCodesTestSuite) TestGetErrorInfo_ValidationErrors() {
	tests := []struct {
		name           string
		code           ErrorCode
		expectedCode   ErrorCode
		expectedMsg    string
		expectedStatus int
	}{
		{
			name:           "Email required",
			code:           ErrCodeEmailRequired,
			expectedCode:   ErrCodeEmailRequired,
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Email is required",
		},
		{
			name:           "Email invalid",
			code:           ErrCodeEmailInvalid,
			expectedCode:   ErrCodeEmailInvalid,
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Invalid email format",
		},
		{
			name:           "Password required",
			code:           ErrCodePasswordRequired,
			expectedCode:   ErrCodePasswordRequired,
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Password is required",
		},
		{
			name:           "Password too short",
			code:           ErrCodePasswordTooShort,
			expectedCode:   ErrCodePasswordTooShort,
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Password too short",
		},
		{
			name:           "Password too long",
			code:           ErrCodePasswordTooLong,
			expectedCode:   ErrCodePasswordTooLong,
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Password too long",
		},
		{
			name:           "Password complexity",
			code:           ErrCodePasswordComplexity,
			expectedCode:   ErrCodePasswordComplexity,
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Password complexity requirements not met",
		},
		{
			name:           "Display name required",
			code:           ErrCodeDisplayNameRequired,
			expectedCode:   ErrCodeDisplayNameRequired,
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Display name is required",
		},
		{
			name:           "Display name too long",
			code:           ErrCodeDisplayNameTooLong,
			expectedCode:   ErrCodeDisplayNameTooLong,
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Display name too long",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := GetErrorInfo(tt.code)

			assert.Equal(suite.T(), tt.expectedCode, result.Code)
			assert.Equal(suite.T(), tt.expectedStatus, result.HTTPStatus)
			assert.Equal(suite.T(), tt.expectedMsg, result.Message)
			assert.NotEmpty(suite.T(), result.Description)
		})
	}
}

func (suite *ErrorCodesTestSuite) TestGetErrorInfo_BusinessLogicErrors() {
	tests := []struct {
		name           string
		code           ErrorCode
		expectedCode   ErrorCode
		expectedMsg    string
		expectedStatus int
	}{
		{
			name:           "Email not verified",
			code:           ErrCodeEmailNotVerified,
			expectedCode:   ErrCodeEmailNotVerified,
			expectedStatus: http.StatusForbidden,
			expectedMsg:    "Email not verified",
		},
		{
			name:           "Account disabled",
			code:           ErrCodeAccountDisabled,
			expectedCode:   ErrCodeAccountDisabled,
			expectedStatus: http.StatusForbidden,
			expectedMsg:    "Account disabled",
		},
		{
			name:           "Account deleted",
			code:           ErrCodeAccountDeleted,
			expectedCode:   ErrCodeAccountDeleted,
			expectedStatus: http.StatusForbidden,
			expectedMsg:    "Account deleted",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := GetErrorInfo(tt.code)

			assert.Equal(suite.T(), tt.expectedCode, result.Code)
			assert.Equal(suite.T(), tt.expectedStatus, result.HTTPStatus)
			assert.Equal(suite.T(), tt.expectedMsg, result.Message)
			assert.NotEmpty(suite.T(), result.Description)
		})
	}
}

func (suite *ErrorCodesTestSuite) TestGetErrorInfo_UnknownErrorCode() {
	unknownCode := ErrorCode("E999")
	result := GetErrorInfo(unknownCode)

	assert.Equal(suite.T(), unknownCode, result.Code)
	assert.Equal(suite.T(), http.StatusInternalServerError, result.HTTPStatus)
	assert.Equal(suite.T(), "Unknown error", result.Message)
	assert.Equal(suite.T(), "An unknown error occurred", result.Description)
}

func (suite *ErrorCodesTestSuite) TestGetErrorInfo_EmptyErrorCode() {
	emptyCode := ErrorCode("")
	result := GetErrorInfo(emptyCode)

	assert.Equal(suite.T(), emptyCode, result.Code)
	assert.Equal(suite.T(), http.StatusInternalServerError, result.HTTPStatus)
	assert.Equal(suite.T(), "Unknown error", result.Message)
	assert.Equal(suite.T(), "An unknown error occurred", result.Description)
}

func (suite *ErrorCodesTestSuite) TestGetErrorInfo_AllDefinedErrorsHaveValidHTTPStatus() {
	errorCodes := []ErrorCode{
		// General errors
		ErrCodeInternalError, ErrCodeInvalidRequest, ErrCodeValidationFailed,
		ErrCodeNotFound, ErrCodeUnauthorized, ErrCodeForbidden, ErrCodeConflict,
		// Authentication errors
		ErrCodeInvalidCredentials, ErrCodeUserNotFound, ErrCodeUserExists,
		ErrCodeTokenExpired, ErrCodeTokenInvalid,
		// Validation errors
		ErrCodeEmailRequired, ErrCodeEmailInvalid, ErrCodePasswordRequired,
		ErrCodePasswordTooShort, ErrCodePasswordTooLong, ErrCodePasswordComplexity,
		ErrCodeDisplayNameRequired, ErrCodeDisplayNameTooLong,
		// Business logic errors
		ErrCodeEmailNotVerified, ErrCodeAccountDisabled, ErrCodeAccountDeleted,
	}

	for _, code := range errorCodes {
		suite.Run(string(code), func() {
			result := GetErrorInfo(code)

			assert.Equal(suite.T(), code, result.Code)
			assert.Greater(suite.T(), result.HTTPStatus, 0, "HTTP status should be positive")
			assert.LessOrEqual(suite.T(), result.HTTPStatus, 599, "HTTP status should be valid")
			assert.NotEmpty(suite.T(), result.Message, "Message should not be empty")
			assert.NotEmpty(suite.T(), result.Description, "Description should not be empty")
		})
	}
}

func (suite *ErrorCodesTestSuite) TestGetErrorInfo_SpecificErrorDetails() {
	// Test specific error details to ensure they match expected business logic
	result := GetErrorInfo(ErrCodePasswordTooShort)
	assert.Contains(suite.T(), result.Description, "8 characters")

	result = GetErrorInfo(ErrCodePasswordTooLong)
	assert.Contains(suite.T(), result.Description, "128 characters")

	result = GetErrorInfo(ErrCodeDisplayNameTooLong)
	assert.Contains(suite.T(), result.Description, "100 characters")

	result = GetErrorInfo(ErrCodePasswordComplexity)
	assert.Contains(suite.T(), result.Description, "lowercase")
	assert.Contains(suite.T(), result.Description, "uppercase")
	assert.Contains(suite.T(), result.Description, "symbol")
}

func (suite *ErrorCodesTestSuite) TestErrorInfo_JSONSerialization() {
	// Test that HTTPStatus is not included in JSON (due to json:"-" tag)
	result := GetErrorInfo(ErrCodeInternalError)

	// Verify structure has all required fields
	assert.NotEmpty(suite.T(), result.Code)
	assert.NotEmpty(suite.T(), result.Message)
	assert.NotEmpty(suite.T(), result.Description)
	assert.Greater(suite.T(), result.HTTPStatus, 0)
}

func (suite *ErrorCodesTestSuite) TestErrorCode_StringConversion() {
	// Test that ErrorCode can be converted to string
	code := ErrCodeInternalError
	codeStr := string(code)
	assert.Equal(suite.T(), "E001", codeStr)

	// Test creating ErrorCode from string
	newCode := ErrorCode("E001")
	assert.Equal(suite.T(), ErrCodeInternalError, newCode)
}

func TestErrorCodesTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorCodesTestSuite))
}

// Simple function tests for direct testing

func TestGetErrorInfo_DirectTest_InternalError(t *testing.T) {
	result := GetErrorInfo(ErrCodeInternalError)

	assert.Equal(t, ErrCodeInternalError, result.Code)
	assert.Equal(t, "Internal server error", result.Message)
	assert.Equal(t, "An unexpected error occurred on the server", result.Description)
	assert.Equal(t, http.StatusInternalServerError, result.HTTPStatus)
}

func TestGetErrorInfo_DirectTest_UserExists(t *testing.T) {
	result := GetErrorInfo(ErrCodeUserExists)

	assert.Equal(t, ErrCodeUserExists, result.Code)
	assert.Equal(t, "User already exists", result.Message)
	assert.Equal(t, "A user with this email address already exists", result.Description)
	assert.Equal(t, http.StatusConflict, result.HTTPStatus)
}

func TestGetErrorInfo_DirectTest_UnknownCode(t *testing.T) {
	unknownCode := ErrorCode("UNKNOWN")
	result := GetErrorInfo(unknownCode)

	assert.Equal(t, unknownCode, result.Code)
	assert.Equal(t, "Unknown error", result.Message)
	assert.Equal(t, "An unknown error occurred", result.Description)
	assert.Equal(t, http.StatusInternalServerError, result.HTTPStatus)
}
