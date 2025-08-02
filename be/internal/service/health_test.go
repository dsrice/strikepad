package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type HealthServiceTestSuite struct {
	suite.Suite
	healthService HealthServiceInterface
}

func (suite *HealthServiceTestSuite) SetupTest() {
	suite.healthService = NewHealthService()
}

func (suite *HealthServiceTestSuite) TestGetHealth() {
	result := suite.healthService.GetHealth()

	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "ok", result.Status)
	assert.Equal(suite.T(), "Server is healthy", result.Message)
}

func TestHealthServiceTestSuite(t *testing.T) {
	suite.Run(t, new(HealthServiceTestSuite))
}

func TestHealthService_GetHealth_Simple(t *testing.T) {
	service := NewHealthService()
	result := service.GetHealth()

	assert.Equal(t, "ok", result.Status)
	assert.Equal(t, "Server is healthy", result.Message)
}
