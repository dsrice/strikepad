package mocks

import (
	"strikepad-backend/internal/model"

	"github.com/stretchr/testify/mock"
)

// MockSessionRepository is a mock implementation of SessionRepositoryInterface
type MockSessionRepository struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockSessionRepository) Create(session *model.UserSession) error {
	args := m.Called(session)
	return args.Error(0)
}

// FindByAccessToken mocks the FindByAccessToken method
func (m *MockSessionRepository) FindByAccessToken(accessToken string) (*model.UserSession, error) {
	args := m.Called(accessToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserSession), args.Error(1)
}

// FindByRefreshToken mocks the FindByRefreshToken method
func (m *MockSessionRepository) FindByRefreshToken(refreshToken string) (*model.UserSession, error) {
	args := m.Called(refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserSession), args.Error(1)
}

// FindActiveByUserID mocks the FindActiveByUserID method
func (m *MockSessionRepository) FindActiveByUserID(userID uint) ([]*model.UserSession, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.UserSession), args.Error(1)
}

// Update mocks the Update method
func (m *MockSessionRepository) Update(session *model.UserSession) error {
	args := m.Called(session)
	return args.Error(0)
}

// InvalidateByUserID mocks the InvalidateByUserID method
func (m *MockSessionRepository) InvalidateByUserID(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

// InvalidateExpiredSessions mocks the InvalidateExpiredSessions method
func (m *MockSessionRepository) InvalidateExpiredSessions() error {
	args := m.Called()
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockSessionRepository) Delete(sessionID uint) error {
	args := m.Called(sessionID)
	return args.Error(0)
}
