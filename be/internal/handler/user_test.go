package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"strikepad-backend/internal/service/mocks"

	"strikepad-backend/internal/auth"
	"strikepad-backend/internal/dto"
	"strikepad-backend/internal/handler"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserHandlerTestSuite struct {
	suite.Suite
	userHandler handler.UserHandlerInterface
	mockService *mocks.MockUserServiceInterface
	echo        *echo.Echo
}

func (suite *UserHandlerTestSuite) SetupTest() {
	suite.mockService = new(mocks.MockUserServiceInterface)
	suite.userHandler = handler.NewUserHandler(suite.mockService)
	suite.echo = echo.New()
}

func (suite *UserHandlerTestSuite) TearDownTest() {
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *UserHandlerTestSuite) TestMe() {
	// Comprehensive table-driven test for me endpoint
	tests := []struct {
		setupContext   func(c echo.Context)
		mockSetup      func()
		expectedError  *dto.ErrorResponse
		expectedData   *dto.UserInfo
		name           string
		description    string
		expectedStatus int
	}{
		{
			name: "successful me request",
			setupContext: func(c echo.Context) {
				c.Set("user_id", uint(1))
			},
			mockSetup: func() {
				expectedResponse := &dto.UserInfo{
					ID:            1,
					Email:         "test@example.com",
					DisplayName:   "Test User",
					EmailVerified: true,
				}
				suite.mockService.On("GetCurrentUser", uint(1)).Return(expectedResponse, nil)
			},
			expectedStatus: http.StatusOK,
			expectedData: &dto.UserInfo{
				ID:            1,
				Email:         "test@example.com",
				DisplayName:   "Test User",
				EmailVerified: true,
			},
			description: "should successfully return current user information",
		},
		{
			name: "missing user ID in context",
			setupContext: func(_ echo.Context) {
				// Don't set user_id
			},
			mockSetup:      func() {}, // No mock setup needed
			expectedStatus: http.StatusUnauthorized,
			expectedError: &dto.ErrorResponse{
				Code:    "E005",
				Message: "Unauthorized",
			},
			description: "should return unauthorized when user ID is missing from context",
		},
		{
			name: "user not found",
			setupContext: func(c echo.Context) {
				c.Set("user_id", uint(999))
			},
			mockSetup: func() {
				suite.mockService.On("GetCurrentUser", uint(999)).Return(nil, auth.ErrInvalidCredentials)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError: &dto.ErrorResponse{
				Code:    "E005",
				Message: "Unauthorized",
			},
			description: "should return unauthorized when user is not found",
		},
		{
			name: "internal server error",
			setupContext: func(c echo.Context) {
				c.Set("user_id", uint(1))
			},
			mockSetup: func() {
				suite.mockService.On("GetCurrentUser", uint(1)).Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError: &dto.ErrorResponse{
				Code:    "E001",
				Message: "Internal server error",
			},
			description: "should return internal server error for unexpected errors",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Setup
			suite.SetupTest() // Reset mocks
			tt.mockSetup()

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/me", http.NoBody)
			rec := httptest.NewRecorder()
			c := suite.echo.NewContext(req, rec)

			// Setup context
			tt.setupContext(c)

			// Execute
			err := suite.userHandler.Me(c)

			// Assert
			assert.NoError(suite.T(), err, tt.description)
			assert.Equal(suite.T(), tt.expectedStatus, rec.Code, tt.description)

			if tt.expectedError != nil {
				var errorResponse dto.ErrorResponse
				err = json.Unmarshal(rec.Body.Bytes(), &errorResponse)
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tt.expectedError.Code, errorResponse.Code, tt.description)
				assert.Equal(suite.T(), tt.expectedError.Message, errorResponse.Message, tt.description)
			}

			if tt.expectedData != nil {
				var response dto.UserInfo
				err = json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tt.expectedData.ID, response.ID, tt.description)
				assert.Equal(suite.T(), tt.expectedData.Email, response.Email, tt.description)
				assert.Equal(suite.T(), tt.expectedData.DisplayName, response.DisplayName, tt.description)
				assert.Equal(suite.T(), tt.expectedData.EmailVerified, response.EmailVerified, tt.description)
			}
		})
	}
}

func (suite *UserHandlerTestSuite) TestNewUserHandler() {
	// Test that NewUserHandler creates a valid handler
	h := handler.NewUserHandler(suite.mockService)
	assert.NotNil(suite.T(), h)
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}
