package testutil

import (
	"github.com/stretchr/testify/mock"

	"strikepad-backend/internal/model"
	"strikepad-backend/internal/service"
)

type MockHealthService struct {
	mock.Mock
}

func (m *MockHealthService) Check() map[string]string {
	args := m.Called()
	result, _ := args.Get(0).(map[string]string)
	return result
}

func NewMockHealthService() service.HealthService {
	return &MockHealthService{}
}

type MockAPIService struct {
	mock.Mock
}

func (m *MockAPIService) GetTestMessage() map[string]string {
	args := m.Called()
	result, _ := args.Get(0).(map[string]string)
	return result
}

func NewMockAPIService() service.APIService {
	return &MockAPIService{}
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id uint) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	result, _ := args.Get(0).(*model.User)
	return result, args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	result, _ := args.Get(0).(*model.User)
	return result, args.Error(1)
}

func (m *MockUserRepository) Update(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) List() ([]model.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	result, _ := args.Get(0).([]model.User)
	return result, args.Error(1)
}
