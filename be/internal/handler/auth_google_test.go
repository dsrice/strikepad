package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"strikepad-backend/internal/auth"
	"strikepad-backend/internal/dto"
	"strikepad-backend/internal/service/mocks"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthHandler_GoogleSignup(t *testing.T) {
	tests := []struct {
		requestBody    map[string]interface{}
		setupMocks     func(*mocks.MockAuthServiceInterface)
		name           string
		expectedStatus int
		expectError    bool
	}{
		{
			name: "successful Google signup",
			requestBody: map[string]interface{}{
				"access_token": "valid_google_token",
			},
			setupMocks: func(mockService *mocks.MockAuthServiceInterface) {
				mockService.On("GoogleSignup", mock.AnythingOfType("*dto.GoogleSignupRequest")).Return(
					&dto.SignupResponse{
						ID:            1,
						Email:         "test@example.com",
						DisplayName:   "Test User",
						EmailVerified: true,
					}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name: "invalid request body",
			requestBody: map[string]interface{}{
				"invalid_field": "value",
			},
			setupMocks: func(_ *mocks.MockAuthServiceInterface) {
				// No service call expected
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "user already exists",
			requestBody: map[string]interface{}{
				"access_token": "valid_google_token",
			},
			setupMocks: func(mockService *mocks.MockAuthServiceInterface) {
				mockService.On("GoogleSignup", mock.AnythingOfType("*dto.GoogleSignupRequest")).Return(
					nil, auth.ErrUserAlreadyExists)
			},
			expectedStatus: http.StatusConflict,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := &mocks.MockAuthServiceInterface{}
			handler := NewAuthHandler(mockService)

			if tt.setupMocks != nil {
				tt.setupMocks(mockService)
			}

			// Create request
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/auth/google/signup", bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(req, rec)

			// Execute
			err := handler.GoogleSignup(c)

			// Assert
			assert.NoError(t, err) // Echo handlers don't return errors for HTTP errors
			assert.Equal(t, tt.expectedStatus, rec.Code)

			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_GoogleLogin(t *testing.T) {
	tests := []struct {
		requestBody    map[string]interface{}
		setupMocks     func(*mocks.MockAuthServiceInterface)
		name           string
		expectedStatus int
		expectError    bool
	}{
		{
			name: "successful Google login",
			requestBody: map[string]interface{}{
				"access_token": "valid_google_token",
			},
			setupMocks: func(mockService *mocks.MockAuthServiceInterface) {
				mockService.On("GoogleLogin", mock.AnythingOfType("*dto.GoogleLoginRequest")).Return(
					&dto.UserInfo{
						ID:            1,
						Email:         "test@example.com",
						DisplayName:   "Test User",
						EmailVerified: true,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "invalid credentials",
			requestBody: map[string]interface{}{
				"access_token": "invalid_token",
			},
			setupMocks: func(mockService *mocks.MockAuthServiceInterface) {
				mockService.On("GoogleLogin", mock.AnythingOfType("*dto.GoogleLoginRequest")).Return(
					nil, auth.ErrInvalidCredentials)
			},
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name: "missing access token",
			requestBody: map[string]interface{}{
				"access_token": "",
			},
			setupMocks: func(_ *mocks.MockAuthServiceInterface) {
				// Validation should fail before service call
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := &mocks.MockAuthServiceInterface{}
			handler := NewAuthHandler(mockService)

			if tt.setupMocks != nil {
				tt.setupMocks(mockService)
			}

			// Create request
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/auth/google/login", bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(req, rec)

			// Execute
			err := handler.GoogleLogin(c)

			// Assert
			assert.NoError(t, err) // Echo handlers don't return errors for HTTP errors
			assert.Equal(t, tt.expectedStatus, rec.Code)

			mockService.AssertExpectations(t)
		})
	}
}
