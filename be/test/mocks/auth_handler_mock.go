package mocks

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

// MockAuthHandler is a mock implementation of AuthHandlerInterface
type MockAuthHandler struct {
	mock.Mock
}

func (m *MockAuthHandler) Signup(c echo.Context) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockAuthHandler) Login(c echo.Context) error {
	args := m.Called(c)
	return args.Error(0)
}
