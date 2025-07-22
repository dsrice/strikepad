package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type APIServiceTestSuite struct {
	suite.Suite
	apiService APIService
}

func (suite *APIServiceTestSuite) SetupTest() {
	suite.apiService = NewAPIService()
}

func (suite *APIServiceTestSuite) TestGetTestMessage() {
	result := suite.apiService.GetTestMessage()

	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "API endpoint working", result["message"])
}

func (suite *APIServiceTestSuite) TestGetTestMessageStructure() {
	result := suite.apiService.GetTestMessage()

	assert.Contains(suite.T(), result, "message")
	assert.Len(suite.T(), result, 1)
}

func TestAPIServiceTestSuite(t *testing.T) {
	suite.Run(t, new(APIServiceTestSuite))
}

func TestAPIService_GetTestMessage_Simple(t *testing.T) {
	service := NewAPIService()
	result := service.GetTestMessage()

	assert.Equal(t, "API endpoint working", result["message"])
}
