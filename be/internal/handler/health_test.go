package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"strikepad-backend/internal/service/mocks"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"strikepad-backend/internal/dto"
)

func TestHealthHandler_Check(t *testing.T) {
	mockService := &mocks.MockHealthServiceInterface{}
	handler := NewHealthHandler(mockService)

	expectedResponse := &dto.HealthResponse{
		Status:  "ok",
		Message: "Server is healthy",
	}
	mockService.On("GetHealth").Return(expectedResponse)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Check(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"status":"ok"`)
	assert.Contains(t, rec.Body.String(), `"message":"Server is healthy"`)
	mockService.AssertExpectations(t)
}

func TestHealthHandler_Check_MockVerification(t *testing.T) {
	mockService := &mocks.MockHealthServiceInterface{}
	handler := NewHealthHandler(mockService)

	expectedResponse := &dto.HealthResponse{
		Status:  "healthy",
		Message: "All systems operational",
	}
	mockService.On("GetHealth").Return(expectedResponse)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Check(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"status":"healthy"`)
	assert.Contains(t, rec.Body.String(), `"message":"All systems operational"`)

	mockService.AssertCalled(t, "GetHealth")
	mockService.AssertNumberOfCalls(t, "GetHealth", 1)
}
