package mocks

import (
	"strikepad-backend/internal/auth"
	"strikepad-backend/internal/model"

	"github.com/stretchr/testify/mock"
)

// MockSessionServiceInterface is a mock implementation of SessionServiceInterface
type MockSessionServiceInterface struct {
	mock.Mock
}

// CreateSession mocks the CreateSession method
func (m *MockSessionServiceInterface) CreateSession(userID uint) (*auth.TokenPair, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.TokenPair), args.Error(1)
}

// ValidateAccessToken mocks the ValidateAccessToken method
func (m *MockSessionServiceInterface) ValidateAccessToken(token string) (*model.UserSession, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserSession), args.Error(1)
}

// RefreshToken mocks the RefreshToken method
func (m *MockSessionServiceInterface) RefreshToken(refreshToken string) (*auth.TokenPair, error) {
	args := m.Called(refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.TokenPair), args.Error(1)
}

// InvalidateSession mocks the InvalidateSession method
func (m *MockSessionServiceInterface) InvalidateSession(accessToken string) error {
	args := m.Called(accessToken)
	return args.Error(0)
}

// InvalidateAllUserSessions mocks the InvalidateAllUserSessions method
func (m *MockSessionServiceInterface) InvalidateAllUserSessions(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

// Logout mocks the Logout method
func (m *MockSessionServiceInterface) Logout(userID uint, accessToken string) error {
	args := m.Called(userID, accessToken)
	return args.Error(0)
}

// CleanupExpiredSessions mocks the CleanupExpiredSessions method
func (m *MockSessionServiceInterface) CleanupExpiredSessions() error {
	args := m.Called()
	return args.Error(0)
}