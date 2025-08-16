package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"strikepad-backend/internal/service/mocks"

	"strikepad-backend/internal/auth"
	"strikepad-backend/internal/dto"
	"strikepad-backend/internal/handler"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AuthHandlerTestSuite struct {
	suite.Suite
	authHandler        handler.AuthHandlerInterface
	mockService        *mocks.MockAuthServiceInterface
	mockSessionService *mocks.MockSessionServiceInterface
	echo               *echo.Echo
}

func (suite *AuthHandlerTestSuite) SetupTest() {
	suite.mockService = new(mocks.MockAuthServiceInterface)
	suite.mockSessionService = new(mocks.MockSessionServiceInterface)
	suite.authHandler = handler.NewAuthHandler(suite.mockService, suite.mockSessionService)
	suite.echo = echo.New()
}

func (suite *AuthHandlerTestSuite) TearDownTest() {
	suite.mockService.AssertExpectations(suite.T())
	suite.mockSessionService.AssertExpectations(suite.T())
}

func (suite *AuthHandlerTestSuite) TestSignup() {
	// Comprehensive table-driven test for signup endpoint
	tests := []struct {
		requestBody    interface{}
		mockSetup      func()
		expectedError  *dto.ErrorResponse
		expectedData   *dto.SignupResponse
		name           string
		description    string
		expectedStatus int
	}{
		{
			name: "successful signup",
			requestBody: dto.SignupRequest{
				Email:       "test@example.com",
				Password:    "Password123!",
				DisplayName: "Test User",
			},
			mockSetup: func() {
				expectedResponse := &dto.SignupResponse{
					ID:            1,
					Email:         "test@example.com",
					DisplayName:   "Test User",
					EmailVerified: false,
					CreatedAt:     time.Now(),
				}
				suite.mockService.On("Signup", mock.MatchedBy(func(req *dto.SignupRequest) bool {
					return req.Email == "test@example.com" &&
						req.Password == "Password123!" &&
						req.DisplayName == "Test User"
				})).Return(expectedResponse, nil)

				// Mock session creation for successful signup
				expectedTokenPair := &auth.TokenPair{
					AccessToken:           "test-access-token",
					RefreshToken:          "test-refresh-token",
					AccessTokenExpiresAt:  time.Now().Add(time.Hour),
					RefreshTokenExpiresAt: time.Now().Add(24 * time.Hour),
				}
				suite.mockSessionService.On("CreateSession", uint(1)).Return(expectedTokenPair, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedData: &dto.SignupResponse{
				ID:            1,
				Email:         "test@example.com",
				DisplayName:   "Test User",
				EmailVerified: false,
			},
			description: "should successfully create new user",
		},
		{
			name:           "invalid JSON",
			requestBody:    "invalid json",
			mockSetup:      func() {}, // No mock setup needed
			expectedStatus: http.StatusBadRequest,
			expectedError: &dto.ErrorResponse{
				Code:    "E002",
				Message: "Invalid request",
			},
			description: "should return error for invalid JSON",
		},
		{
			name: "validation failure - missing fields",
			requestBody: dto.SignupRequest{
				Email:    "", // Invalid - required
				Password: "", // Invalid - required
			},
			mockSetup:      func() {}, // No mock setup needed
			expectedStatus: http.StatusBadRequest,
			expectedError: &dto.ErrorResponse{
				Code:    "E003",
				Message: "Validation failed",
			},
			description: "should return validation error for missing required fields",
		},
		{
			name: "validation failure - invalid email",
			requestBody: dto.SignupRequest{
				Email:       "invalid-email",
				Password:    "Password123!",
				DisplayName: "Test User",
			},
			mockSetup:      func() {}, // No mock setup needed
			expectedStatus: http.StatusBadRequest,
			expectedError: &dto.ErrorResponse{
				Code:    "E003",
				Message: "Validation failed",
			},
			description: "should return validation error for invalid email format",
		},
		{
			name: "validation failure - password too short",
			requestBody: dto.SignupRequest{
				Email:       "test@example.com",
				Password:    "short",
				DisplayName: "Test User",
			},
			mockSetup:      func() {}, // No mock setup needed
			expectedStatus: http.StatusBadRequest,
			expectedError: &dto.ErrorResponse{
				Code:    "E003",
				Message: "Validation failed",
			},
			description: "should return validation error for short password",
		},
		{
			name: "validation failure - password too long",
			requestBody: dto.SignupRequest{
				Email:       "test@example.com",
				Password:    "Password123!" + string(make([]byte, 120)),
				DisplayName: "Test User",
			},
			mockSetup:      func() {}, // No mock setup needed
			expectedStatus: http.StatusBadRequest,
			expectedError: &dto.ErrorResponse{
				Code:    "E003",
				Message: "Validation failed",
			},
			description: "should return validation error for long password",
		},
		{
			name: "user already exists",
			requestBody: dto.SignupRequest{
				Email:       "existing@example.com",
				Password:    "Password123!",
				DisplayName: "Test User",
			},
			mockSetup: func() {
				suite.mockService.On("Signup", mock.AnythingOfType("*dto.SignupRequest")).Return(nil, auth.ErrUserAlreadyExists)
			},
			expectedStatus: http.StatusConflict,
			expectedError: &dto.ErrorResponse{
				Code:    "E102",
				Message: "User already exists",
			},
			description: "should return conflict error when user already exists",
		},
		{
			name: "internal server error",
			requestBody: dto.SignupRequest{
				Email:       "test@example.com",
				Password:    "Password123!",
				DisplayName: "Test User",
			},
			mockSetup: func() {
				suite.mockService.On("Signup", mock.AnythingOfType("*dto.SignupRequest")).Return(nil, assert.AnError)
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
			var req *http.Request
			if str, ok := tt.requestBody.(string); ok {
				req = httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBufferString(str))
			} else {
				jsonBody, _ := json.Marshal(tt.requestBody)
				req = httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonBody))
			}
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := suite.echo.NewContext(req, rec)

			// Execute
			err := suite.authHandler.Signup(c)

			// Assert
			assert.NoError(suite.T(), err, tt.description)
			assert.Equal(suite.T(), tt.expectedStatus, rec.Code, tt.description)

			if tt.expectedError != nil {
				var errorResponse dto.ErrorResponse
				err = json.Unmarshal(rec.Body.Bytes(), &errorResponse)
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tt.expectedError.Code, errorResponse.Code, tt.description)
				assert.Equal(suite.T(), tt.expectedError.Message, errorResponse.Message, tt.description)
				if tt.expectedError.Code == "E003" { // Validation failed
					assert.NotEmpty(suite.T(), errorResponse.Details, "Validation errors should have details")
				}
			}

			if tt.expectedData != nil {
				var response dto.SignupResponse
				err = json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tt.expectedData.ID, response.ID, tt.description)
				assert.Equal(suite.T(), tt.expectedData.Email, response.Email, tt.description)
				assert.Equal(suite.T(), tt.expectedData.DisplayName, response.DisplayName, tt.description)
				assert.Equal(suite.T(), tt.expectedData.EmailVerified, response.EmailVerified, tt.description)
				assert.NotZero(suite.T(), response.CreatedAt, "CreatedAt should be set")
			}
		})
	}
}

func (suite *AuthHandlerTestSuite) TestLogin() {
	// Comprehensive table-driven test for login endpoint
	tests := []struct {
		requestBody    interface{}
		mockSetup      func()
		expectedError  *dto.ErrorResponse
		expectedData   *dto.UserInfo
		name           string
		description    string
		expectedStatus int
	}{
		{
			name: "successful login",
			requestBody: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "Password123!",
			},
			mockSetup: func() {
				expectedResponse := &dto.UserInfo{
					ID:            1,
					Email:         "test@example.com",
					DisplayName:   "Test User",
					EmailVerified: false,
				}
				suite.mockService.On("Login", mock.MatchedBy(func(req *dto.LoginRequest) bool {
					return req.Email == "test@example.com" && req.Password == "Password123!"
				})).Return(expectedResponse, nil)

				// Mock session creation for successful login
				expectedTokenPair := &auth.TokenPair{
					AccessToken:           "test-access-token",
					RefreshToken:          "test-refresh-token",
					AccessTokenExpiresAt:  time.Now().Add(time.Hour),
					RefreshTokenExpiresAt: time.Now().Add(24 * time.Hour),
				}
				suite.mockSessionService.On("CreateSession", uint(1)).Return(expectedTokenPair, nil)
			},
			expectedStatus: http.StatusOK,
			expectedData: &dto.UserInfo{
				ID:            1,
				Email:         "test@example.com",
				DisplayName:   "Test User",
				EmailVerified: false,
			},
			description: "should successfully authenticate user",
		},
		{
			name:           "invalid JSON",
			requestBody:    "invalid json",
			mockSetup:      func() {}, // No mock setup needed
			expectedStatus: http.StatusBadRequest,
			expectedError: &dto.ErrorResponse{
				Code:    "E002",
				Message: "Invalid request",
			},
			description: "should return error for invalid JSON",
		},
		{
			name: "validation failure - missing fields",
			requestBody: dto.LoginRequest{
				Email:    "", // Invalid - required
				Password: "", // Invalid - required
			},
			mockSetup:      func() {}, // No mock setup needed
			expectedStatus: http.StatusBadRequest,
			expectedError: &dto.ErrorResponse{
				Code:    "E003",
				Message: "Validation failed",
			},
			description: "should return validation error for missing required fields",
		},
		{
			name: "invalid credentials",
			requestBody: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func() {
				suite.mockService.On("Login", mock.AnythingOfType("*dto.LoginRequest")).Return(nil, auth.ErrInvalidCredentials)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError: &dto.ErrorResponse{
				Code:    "E100",
				Message: "Invalid credentials",
			},
			description: "should return unauthorized for invalid credentials",
		},
		{
			name: "internal server error",
			requestBody: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "Password123!",
			},
			mockSetup: func() {
				suite.mockService.On("Login", mock.AnythingOfType("*dto.LoginRequest")).Return(nil, assert.AnError)
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
			var req *http.Request
			if str, ok := tt.requestBody.(string); ok {
				req = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(str))
			} else {
				jsonBody, _ := json.Marshal(tt.requestBody)
				req = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
			}
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := suite.echo.NewContext(req, rec)

			// Execute
			err := suite.authHandler.Login(c)

			// Assert
			assert.NoError(suite.T(), err, tt.description)
			assert.Equal(suite.T(), tt.expectedStatus, rec.Code, tt.description)

			if tt.expectedError != nil {
				var errorResponse dto.ErrorResponse
				err = json.Unmarshal(rec.Body.Bytes(), &errorResponse)
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tt.expectedError.Code, errorResponse.Code, tt.description)
				assert.Equal(suite.T(), tt.expectedError.Message, errorResponse.Message, tt.description)
				if tt.expectedError.Code == "E003" { // Validation failed
					assert.NotEmpty(suite.T(), errorResponse.Details, "Validation errors should have details")
				}
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

func (suite *AuthHandlerTestSuite) TestNewAuthHandler() {
	// Test that NewAuthHandler creates a valid handler
	h := handler.NewAuthHandler(suite.mockService, suite.mockSessionService)
	assert.NotNil(suite.T(), h)
}

func TestAuthHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(AuthHandlerTestSuite))
}