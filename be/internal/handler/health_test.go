package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strikepad-backend/internal/handler"
	"testing"

	"strikepad-backend/internal/service/mocks"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"strikepad-backend/internal/dto"
)

func TestHealthHandler_Check(t *testing.T) {
	// Table-driven test for health check endpoint
	tests := []struct {
		mockResponse      *dto.HealthResponse
		name              string
		description       string
		expectedInBody    []string
		expectedStatus    int
		expectedCallCount int
		checkMockCalls    bool
	}{
		{
			name: "standard health check",
			mockResponse: &dto.HealthResponse{
				Status:  "ok",
				Message: "Server is healthy",
			},
			expectedStatus:    http.StatusOK,
			expectedInBody:    []string{`"status":"ok"`, `"message":"Server is healthy"`},
			checkMockCalls:    true,
			expectedCallCount: 1,
			description:       "should return healthy status with standard message",
		},
		{
			name: "alternate health message",
			mockResponse: &dto.HealthResponse{
				Status:  "healthy",
				Message: "All systems operational",
			},
			expectedStatus:    http.StatusOK,
			expectedInBody:    []string{`"status":"healthy"`, `"message":"All systems operational"`},
			checkMockCalls:    true,
			expectedCallCount: 1,
			description:       "should return healthy status with alternate message",
		},
		{
			name: "service ready status",
			mockResponse: &dto.HealthResponse{
				Status:  "ready",
				Message: "Service is ready to accept requests",
			},
			expectedStatus:    http.StatusOK,
			expectedInBody:    []string{`"status":"ready"`, `"message":"Service is ready to accept requests"`},
			checkMockCalls:    true,
			expectedCallCount: 1,
			description:       "should return ready status with custom message",
		},
		{
			name: "minimal response",
			mockResponse: &dto.HealthResponse{
				Status:  "up",
				Message: "OK",
			},
			expectedStatus:    http.StatusOK,
			expectedInBody:    []string{`"status":"up"`, `"message":"OK"`},
			checkMockCalls:    false, // Just check basic functionality
			expectedCallCount: 0,
			description:       "should handle minimal response correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := &mocks.MockHealthServiceInterface{}
			hd := handler.NewHealthHandler(mockService)
			mockService.On("GetHealth").Return(tt.mockResponse)

			// Create request
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/health", http.NoBody)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Execute
			err := hd.Check(c)

			// Assert
			assert.NoError(t, err, tt.description)
			assert.Equal(t, tt.expectedStatus, rec.Code, tt.description)

			// Check response body contains expected content
			for _, expectedContent := range tt.expectedInBody {
				assert.Contains(t, rec.Body.String(), expectedContent,
					"Response should contain: %s", expectedContent)
			}

			// Verify mock expectations
			mockService.AssertExpectations(t)

			// Additional mock verification if requested
			if tt.checkMockCalls {
				mockService.AssertCalled(t, "GetHealth")
				mockService.AssertNumberOfCalls(t, "GetHealth", tt.expectedCallCount)
			}
		})
	}
}

func TestHealthHandler_NewHealthHandler(t *testing.T) {
	// Test handler creation
	mockService := &mocks.MockHealthServiceInterface{}
	hd := handler.NewHealthHandler(mockService)

	assert.NotNil(t, hd, "Handler should not be nil")
	assert.NotNil(t, hd, "Handler should be properly initialized")
}