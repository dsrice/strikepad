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
	testCases := []struct {
		name           string
		expectedMsg    string
		expectedLength int
		checkStructure bool
	}{
		{
			name:           "Check message content",
			expectedMsg:    "API endpoint working",
			expectedLength: 1,
			checkStructure: false,
		},
		{
			name:           "Check message structure",
			expectedMsg:    "API endpoint working",
			expectedLength: 1,
			checkStructure: true,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			result := suite.apiService.GetTestMessage()

			assert.NotNil(t, result)
			assert.Equal(t, tc.expectedMsg, result["message"])

			if tc.checkStructure {
				assert.Contains(t, result, "message")
				assert.Len(t, result, tc.expectedLength)
			}
		})
	}
}

func TestAPIServiceTestSuite(t *testing.T) {
	suite.Run(t, new(APIServiceTestSuite))
}

func TestAPIService_GetTestMessage_Simple(t *testing.T) {
	testCases := []struct {
		name        string
		expectedMsg string
	}{
		{
			name:        "Simple test for message content",
			expectedMsg: "API endpoint working",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := NewAPIService()
			result := service.GetTestMessage()

			assert.Equal(t, tc.expectedMsg, result["message"])
		})
	}
}
