package mocks

import (
	"fmt"
	"strikepad-backend/internal/dto"

	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock implementation of AuthServiceInterface
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Signup(req *dto.SignupRequest) (*dto.SignupResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.SignupResponse), args.Error(1)
}

func (m *MockAuthService) Login(req *dto.LoginRequest) (*dto.UserInfo, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	userInfo, ok := args.Get(0).(*dto.UserInfo)
	if !ok {
		// 型アサーション失敗時の安全なハンドリング
		return nil, fmt.Errorf("invalid type for UserInfo in mock result")
	}

	return userInfo, args.Error(1)
}