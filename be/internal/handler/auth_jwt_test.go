package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"strikepad-backend/internal/dto"
	"strikepad-backend/internal/handler"
	authmocks "strikepad-backend/internal/service/mocks"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthJWTHandlerTestSuite struct {
	suite.Suite
	authHandler    handler.AuthHandlerInterface
	mockAuthSvc    *authmocks.MockAuthServiceInterface
	mockSessionSvc *authmocks.MockSessionServiceInterface
	echo           *echo.Echo
}

func (suite *AuthJWTHandlerTestSuite) SetupTest() {
	suite.mockAuthSvc = new(authmocks.MockAuthServiceInterface)
	suite.mockSessionSvc = new(authmocks.MockSessionServiceInterface)
	suite.authHandler = handler.NewAuthHandler(suite.mockAuthSvc, suite.mockSessionSvc)
	suite.echo = echo.New()
}

func (suite *AuthJWTHandlerTestSuite) TearDownTest() {
	// Reset mocks for next test case
	suite.mockAuthSvc.ExpectedCalls = nil
	suite.mockAuthSvc.Calls = nil
	suite.mockSessionSvc.ExpectedCalls = nil
	suite.mockSessionSvc.Calls = nil
}

func (suite *AuthJWTHandlerTestSuite) TestLogout() {
	testCases := []struct {
		setupContext   func(c echo.Context)
		mockSetup      func()
		expectedError  *dto.ErrorResponse
		name           string
		expectedMsg    string
		expectedStatus int
	}{
		{
			name: "Success",
			setupContext: func(c echo.Context) {
				c.Set("user_id", uint(123))
				c.Set("access_token", "valid-access-token")
			},
			mockSetup: func() {
				suite.mockSessionSvc.On("Logout", uint(123), "valid-access-token").Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedMsg:    "Logout successful",
		},
		{
			name: "Missing user ID",
			setupContext: func(c echo.Context) {
				c.Set("access_token", "valid-access-token")
				// user_id not set
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedError: &dto.ErrorResponse{
				Code:    "E005",
				Message: "Unauthorized",
			},
		},
		{
			name: "Missing access token",
			setupContext: func(c echo.Context) {
				c.Set("user_id", uint(123))
				// access_token not set
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusInternalServerError,
			expectedError: &dto.ErrorResponse{
				Code:    "E001",
				Message: "Internal server error",
			},
		},
		{
			name: "Invalid user ID type",
			setupContext: func(c echo.Context) {
				c.Set("user_id", "invalid-type") // string instead of uint
				c.Set("access_token", "valid-access-token")
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedError: &dto.ErrorResponse{
				Code:    "E005",
				Message: "Unauthorized",
			},
		},
		{
			name: "Invalid access token type",
			setupContext: func(c echo.Context) {
				c.Set("user_id", uint(123))
				c.Set("access_token", 12345) // int instead of string
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusInternalServerError,
			expectedError: &dto.ErrorResponse{
				Code:    "E001",
				Message: "Internal server error",
			},
		},
		{
			name: "Session service error",
			setupContext: func(c echo.Context) {
				c.Set("user_id", uint(456))
				c.Set("access_token", "error-token")
			},
			mockSetup: func() {
				suite.mockSessionSvc.On("Logout", uint(456), "error-token").Return(errors.New("session not found"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError: &dto.ErrorResponse{
				Code:    "E001",
				Message: "Internal server error",
			},
		},
		{
			name: "Session does not belong to user",
			setupContext: func(c echo.Context) {
				c.Set("user_id", uint(789))
				c.Set("access_token", "other-user-token")
			},
			mockSetup: func() {
				suite.mockSessionSvc.On("Logout", uint(789), "other-user-token").Return(errors.New("session does not belong to user"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError: &dto.ErrorResponse{
				Code:    "E001",
				Message: "Internal server error",
			},
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Setup mocks
			tc.mockSetup()

			// Create HTTP request and response recorder
			req := httptest.NewRequest(http.MethodPost, "/auth/logout", http.NoBody)
			rec := httptest.NewRecorder()
			c := suite.echo.NewContext(req, rec)

			// Setup context
			tc.setupContext(c)

			// Execute
			err := suite.authHandler.Logout(c)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedError != nil {
				var errorResponse dto.ErrorResponse
				err = json.Unmarshal(rec.Body.Bytes(), &errorResponse)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedError.Code, errorResponse.Code)
				assert.Equal(t, tc.expectedError.Message, errorResponse.Message)
			} else if tc.expectedMsg != "" {
				var response map[string]string
				err = json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedMsg, response["message"])
			}
		})
	}
}

func TestAuthJWTHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(AuthJWTHandlerTestSuite))
}
