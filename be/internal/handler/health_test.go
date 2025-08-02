package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"strikepad-backend/test/testutil"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler_Health(t *testing.T) {
	mockService := testutil.NewMockHealthService().(*testutil.MockHealthService)
	handler := NewHealthHandler(mockService)

	mockService.On("Check").Return(map[string]string{"status": "ok"})

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Health(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"status":"ok"`)
	mockService.AssertExpectations(t)
}

func TestHealthHandler_Health_MockVerification(t *testing.T) {
	mockService := &testutil.MockHealthService{}
	handler := NewHealthHandler(mockService)

	expectedResult := map[string]string{"status": "healthy"}
	mockService.On("Check").Return(expectedResult)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Health(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"status":"healthy"`)

	mockService.AssertCalled(t, "Check")
	mockService.AssertNumberOfCalls(t, "Check", 1)
}
