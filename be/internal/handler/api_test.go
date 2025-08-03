package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"strikepad-backend/internal/service/mocks"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type APIHandlerTestSuite struct {
	suite.Suite
	handler    *APIHandler
	apiService *mocks.MockAPIServiceInterface
	echo       *echo.Echo
}

func (suite *APIHandlerTestSuite) SetupTest() {
	suite.apiService = &mocks.MockAPIServiceInterface{}
	suite.handler = NewAPIHandler(suite.apiService)
	suite.echo = echo.New()
}

func (suite *APIHandlerTestSuite) TestNewAPIHandler() {
	// Test handler creation with various scenarios
	tests := []struct {
		service     *mocks.MockAPIServiceInterface
		name        string
		description string
		expectNil   bool
	}{
		{
			name:        "valid service",
			service:     &mocks.MockAPIServiceInterface{},
			expectNil:   false,
			description: "should create handler with valid service",
		},
		{
			name:        "different service instance",
			service:     &mocks.MockAPIServiceInterface{},
			expectNil:   false,
			description: "should create handler with different service instance",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			hd := NewAPIHandler(tt.service)

			if tt.expectNil {
				assert.Nil(suite.T(), hd, tt.description)
			} else {
				assert.NotNil(suite.T(), hd, tt.description)
				// Since we're in the handler package, we can't directly access unexported fields
				// Instead, we'll verify the handler works by calling a method
				result := make(map[string]string)
				tt.service.On("GetTestMessage").Return(result)

				// Create a proper echo.Context for testing
				req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
				rec := httptest.NewRecorder()
				c := echo.New().NewContext(req, rec)

				// Test that the handler works without panicking
				assert.NotPanics(suite.T(), func() { hd.Test(c) }, "Handler should not panic when used")
				tt.service.AssertExpectations(suite.T())
			}
		})
	}
}

func (suite *APIHandlerTestSuite) TestTest() {
	// Table-driven test for API test endpoint
	tests := []struct {
		mockResponse      map[string]string
		name              string
		description       string
		expectedInBody    []string
		expectedStatus    int
		expectedCallCount int
		checkMockCalls    bool
	}{
		{
			name: "standard API test",
			mockResponse: map[string]string{
				"message": "API is working",
			},
			expectedStatus:    http.StatusOK,
			expectedInBody:    []string{"API is working"},
			checkMockCalls:    true,
			expectedCallCount: 1,
			description:       "should return standard API working message",
		},
		{
			name: "custom API message",
			mockResponse: map[string]string{
				"message": "Service is operational",
				"status":  "active",
			},
			expectedStatus:    http.StatusOK,
			expectedInBody:    []string{"Service is operational", "active"},
			checkMockCalls:    true,
			expectedCallCount: 1,
			description:       "should return custom message with additional fields",
		},
		{
			name: "minimal response",
			mockResponse: map[string]string{
				"status": "OK",
			},
			expectedStatus:    http.StatusOK,
			expectedInBody:    []string{"OK"},
			checkMockCalls:    false,
			expectedCallCount: 0,
			description:       "should handle minimal response",
		},
		{
			name: "detailed API info",
			mockResponse: map[string]string{
				"message": "API endpoint functioning normally",
				"version": "1.0.0",
				"env":     "test",
			},
			expectedStatus:    http.StatusOK,
			expectedInBody:    []string{"API endpoint functioning normally", "1.0.0", "test"},
			checkMockCalls:    true,
			expectedCallCount: 1,
			description:       "should return detailed API information",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Setup fresh mock for each test
			apiService := &mocks.MockAPIServiceInterface{}
			hd := NewAPIHandler(apiService)
			apiService.On("GetTestMessage").Return(tt.mockResponse)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
			rec := httptest.NewRecorder()
			c := suite.echo.NewContext(req, rec)

			// Execute
			err := hd.Test(c)

			// Assert
			assert.NoError(suite.T(), err, tt.description)
			assert.Equal(suite.T(), tt.expectedStatus, rec.Code, tt.description)

			// Check response body contains expected content
			for _, expectedContent := range tt.expectedInBody {
				assert.Contains(suite.T(), rec.Body.String(), expectedContent,
					"Response should contain: %s", expectedContent)
			}

			// Verify mock expectations
			apiService.AssertExpectations(suite.T())

			// Additional mock verification if requested
			if tt.checkMockCalls {
				apiService.AssertCalled(suite.T(), "GetTestMessage")
				apiService.AssertNumberOfCalls(suite.T(), "GetTestMessage", tt.expectedCallCount)
			}
		})
	}
}

func TestAPIHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(APIHandlerTestSuite))
}