package mocks

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

// MockHealthHandler is a mock implementation of HealthHandlerInterface
type MockHealthHandler struct {
	mock.Mock
}

func (m *MockHealthHandler) Check(c echo.Context) error {
	args := m.Called(c)
	return args.Error(0)
}
