package handler

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

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AuthHandlerTestSuite struct {
	suite.Suite
	authHandler AuthHandlerInterface
	mockService *mocks.MockAuthServiceInterface
	echo        *echo.Echo
}

func (suite *AuthHandlerTestSuite) SetupTest() {
	suite.mockService = new(mocks.MockAuthServiceInterface)
	suite.authHandler = NewAuthHandler(suite.mockService)
	suite.echo = echo.New()
}

func (suite *AuthHandlerTestSuite) TearDownTest() {
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *AuthHandlerTestSuite) TestSignupSuccess() {
	// Arrange
	requestBody := dto.SignupRequest{
		Email:       "test@example.com",
		Password:    "Password123!",
		DisplayName: "Test User",
	}

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

	// Create request
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	// Act
	err := suite.authHandler.Signup(c)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusCreated, rec.Code)

	var response dto.SignupResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedResponse.ID, response.ID)
	assert.Equal(suite.T(), expectedResponse.Email, response.Email)
	assert.Equal(suite.T(), expectedResponse.DisplayName, response.DisplayName)
	assert.Equal(suite.T(), expectedResponse.EmailVerified, response.EmailVerified)
}

func (suite *AuthHandlerTestSuite) TestSignupInvalidJSON() {
	// Create invalid JSON request
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	// Act
	err := suite.authHandler.Signup(c)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)

	var errorResponse dto.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &errorResponse)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "E002", errorResponse.Code)
	assert.Equal(suite.T(), "Invalid request", errorResponse.Message)
}

func (suite *AuthHandlerTestSuite) TestSignupValidationFailure() {
	// Arrange - missing required fields
	requestBody := dto.SignupRequest{
		Email:    "", // Invalid - required
		Password: "", // Invalid - required
	}

	// Create request
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	// Act
	err := suite.authHandler.Signup(c)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)

	var errorResponse dto.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &errorResponse)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "E003", errorResponse.Code)
	assert.Equal(suite.T(), "Validation failed", errorResponse.Message)
	assert.NotEmpty(suite.T(), errorResponse.Details)
}

func (suite *AuthHandlerTestSuite) TestSignupUserAlreadyExists() {
	// Arrange
	requestBody := dto.SignupRequest{
		Email:       "existing@example.com",
		Password:    "Password123!",
		DisplayName: "Test User",
	}

	suite.mockService.On("Signup", mock.AnythingOfType("*dto.SignupRequest")).Return(nil, auth.ErrUserAlreadyExists)

	// Create request
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	// Act
	err := suite.authHandler.Signup(c)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusConflict, rec.Code)

	var errorResponse dto.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &errorResponse)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "E102", errorResponse.Code)
	assert.Equal(suite.T(), "User already exists", errorResponse.Message)
}

func (suite *AuthHandlerTestSuite) TestSignupInvalidEmail() {
	// Arrange - invalid email will be caught by validator, not service
	requestBody := dto.SignupRequest{
		Email:       "invalid-email",
		Password:    "Password123!",
		DisplayName: "Test User",
	}

	// No mock needed as validator catches this before service is called

	// Create request
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	// Act
	err := suite.authHandler.Signup(c)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)

	var errorResponse dto.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &errorResponse)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "E003", errorResponse.Code) // Validation failed
	assert.Equal(suite.T(), "Validation failed", errorResponse.Message)
	assert.NotEmpty(suite.T(), errorResponse.Details) // Should have validation details
}

func (suite *AuthHandlerTestSuite) TestSignupPasswordTooShort() {
	// Arrange - password too short will be caught by validator
	requestBody := dto.SignupRequest{
		Email:       "test@example.com",
		Password:    "short",
		DisplayName: "Test User",
	}

	// No mock needed as validator catches this before service is called

	// Create request
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	// Act
	err := suite.authHandler.Signup(c)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)

	var errorResponse dto.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &errorResponse)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "E003", errorResponse.Code) // Validation failed
	assert.Equal(suite.T(), "Validation failed", errorResponse.Message)
	assert.NotEmpty(suite.T(), errorResponse.Details) // Should have validation details
}

func (suite *AuthHandlerTestSuite) TestSignupPasswordTooLong() {
	// Arrange - password too long will be caught by validator
	requestBody := dto.SignupRequest{
		Email:       "test@example.com",
		Password:    "Password123!" + string(make([]byte, 120)), // Too long
		DisplayName: "Test User",
	}

	// No mock needed as validator catches this before service is called

	// Create request
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	// Act
	err := suite.authHandler.Signup(c)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)

	var errorResponse dto.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &errorResponse)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "E003", errorResponse.Code) // Validation failed
	assert.Equal(suite.T(), "Validation failed", errorResponse.Message)
	assert.NotEmpty(suite.T(), errorResponse.Details) // Should have validation details
}

func (suite *AuthHandlerTestSuite) TestSignupInternalError() {
	// Arrange
	requestBody := dto.SignupRequest{
		Email:       "test@example.com",
		Password:    "Password123!",
		DisplayName: "Test User",
	}

	suite.mockService.On("Signup", mock.AnythingOfType("*dto.SignupRequest")).Return(nil, assert.AnError)

	// Create request
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	// Act
	err := suite.authHandler.Signup(c)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusInternalServerError, rec.Code)

	var errorResponse dto.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &errorResponse)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "E001", errorResponse.Code)
	assert.Equal(suite.T(), "Internal server error", errorResponse.Message)
}

func (suite *AuthHandlerTestSuite) TestLoginSuccess() {
	// Arrange
	requestBody := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "Password123!",
	}

	expectedResponse := &dto.UserInfo{
		ID:            1,
		Email:         "test@example.com",
		DisplayName:   "Test User",
		EmailVerified: false,
	}

	suite.mockService.On("Login", mock.MatchedBy(func(req *dto.LoginRequest) bool {
		return req.Email == "test@example.com" && req.Password == "Password123!"
	})).Return(expectedResponse, nil)

	// Create request
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	// Act
	err := suite.authHandler.Login(c)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	var response dto.UserInfo
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedResponse.ID, response.ID)
	assert.Equal(suite.T(), expectedResponse.Email, response.Email)
	assert.Equal(suite.T(), expectedResponse.DisplayName, response.DisplayName)
	assert.Equal(suite.T(), expectedResponse.EmailVerified, response.EmailVerified)
}

func (suite *AuthHandlerTestSuite) TestLoginInvalidJSON() {
	// Create invalid JSON request
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	// Act
	err := suite.authHandler.Login(c)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)

	var errorResponse dto.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &errorResponse)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "E002", errorResponse.Code)
	assert.Equal(suite.T(), "Invalid request", errorResponse.Message)
}

func (suite *AuthHandlerTestSuite) TestLoginValidationFailure() {
	// Arrange - missing required fields
	requestBody := dto.LoginRequest{
		Email:    "", // Invalid - required
		Password: "", // Invalid - required
	}

	// Create request
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	// Act
	err := suite.authHandler.Login(c)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)

	var errorResponse dto.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &errorResponse)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "E003", errorResponse.Code)
	assert.Equal(suite.T(), "Validation failed", errorResponse.Message)
	assert.NotEmpty(suite.T(), errorResponse.Details)
}

func (suite *AuthHandlerTestSuite) TestLoginInvalidCredentials() {
	// Arrange
	requestBody := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	suite.mockService.On("Login", mock.AnythingOfType("*dto.LoginRequest")).Return(nil, auth.ErrInvalidCredentials)

	// Create request
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	// Act
	err := suite.authHandler.Login(c)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusUnauthorized, rec.Code)

	var errorResponse dto.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &errorResponse)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "E100", errorResponse.Code)
	assert.Equal(suite.T(), "Invalid credentials", errorResponse.Message)
}

func (suite *AuthHandlerTestSuite) TestLoginInternalError() {
	// Arrange
	requestBody := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "Password123!",
	}

	suite.mockService.On("Login", mock.AnythingOfType("*dto.LoginRequest")).Return(nil, assert.AnError)

	// Create request
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	// Act
	err := suite.authHandler.Login(c)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusInternalServerError, rec.Code)

	var errorResponse dto.ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &errorResponse)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "E001", errorResponse.Code)
	assert.Equal(suite.T(), "Internal server error", errorResponse.Message)
}

func (suite *AuthHandlerTestSuite) TestNewAuthHandler() {
	// Test that NewAuthHandler creates a valid handler
	handler := NewAuthHandler(suite.mockService)
	assert.NotNil(suite.T(), handler)
}

func TestAuthHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(AuthHandlerTestSuite))
}
