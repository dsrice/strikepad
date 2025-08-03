package service_test

import (
	"testing"

	"strikepad-backend/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type HealthServiceTestSuite struct {
	suite.Suite
	healthService service.HealthServiceInterface
}

func (suite *HealthServiceTestSuite) SetupTest() {
	suite.healthService = service.NewHealthService()
}

func (suite *HealthServiceTestSuite) TestGetHealth() {
	testCases := []struct {
		name            string
		expectedStatus  string
		expectedMessage string
	}{
		{
			name:            "Health check returns ok status",
			expectedStatus:  "ok",
			expectedMessage: "Server is healthy",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			result := suite.healthService.GetHealth()

			assert.NotNil(t, result)
			assert.Equal(t, tc.expectedStatus, result.Status)
			assert.Equal(t, tc.expectedMessage, result.Message)
		})
	}
}

func TestHealthServiceTestSuite(t *testing.T) {
	suite.Run(t, new(HealthServiceTestSuite))
}

func TestHealthService_GetHealth_Simple(t *testing.T) {
	testCases := []struct {
		name            string
		expectedStatus  string
		expectedMessage string
	}{
		{
			name:            "Simple health check test",
			expectedStatus:  "ok",
			expectedMessage: "Server is healthy",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := service.NewHealthService()
			result := svc.GetHealth()

			assert.Equal(t, tc.expectedStatus, result.Status)
			assert.Equal(t, tc.expectedMessage, result.Message)
		})
	}
}