package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type HealthServiceTestSuite struct {
	suite.Suite
	healthService HealthService
}

func (suite *HealthServiceTestSuite) SetupTest() {
	suite.healthService = NewHealthService()
}

func (suite *HealthServiceTestSuite) TestCheck() {
	result := suite.healthService.Check()

	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "ok", result["status"])
}

func (suite *HealthServiceTestSuite) TestCheckReturnValue() {
	result := suite.healthService.Check()

	assert.Contains(suite.T(), result, "status")
	assert.Len(suite.T(), result, 1)
}

func TestHealthServiceTestSuite(t *testing.T) {
	suite.Run(t, new(HealthServiceTestSuite))
}

func TestHealthService_Check_Simple(t *testing.T) {
	service := NewHealthService()
	result := service.Check()

	assert.Equal(t, "ok", result["status"])
}
