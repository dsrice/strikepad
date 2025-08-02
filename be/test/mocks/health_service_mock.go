package mocks

import (
	"strikepad-backend/internal/dto"

	"github.com/stretchr/testify/mock"
)

// MockHealthService is a mock implementation of HealthServiceInterface
type MockHealthService struct {
	mock.Mock
}

func (m *MockHealthService) GetHealth() *dto.HealthResponse {
	args := m.Called()
	return args.Get(0).(*dto.HealthResponse)
}
