package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"strikepad-backend/internal/service/mocks"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type APIHandlerTestSuite struct {
	suite.Suite
	handler    *APIHandler
	apiService *mocks.MockAPIServiceInterface
	echo       *echo.Echo
}

func (suite *APIHandlerTestSuite) SetupTest() {
	suite.apiService = &mocks.MockAPIServiceInterface{}
	suite.handler = NewAPIHandler(suite.apiService)
	suite.echo = echo.New()
}

func (suite *APIHandlerTestSuite) TestNewAPIHandler() {
	// Test handler creation
	assert.NotNil(suite.T(), suite.handler)
	assert.NotNil(suite.T(), suite.handler.apiService)
}

func (suite *APIHandlerTestSuite) TestTest() {
	// Setup mock
	expectedMessage := map[string]string{
		"message": "API is working",
	}
	suite.apiService.On("GetTestMessage").Return(expectedMessage)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := suite.echo.NewContext(req, rec)

	// Execute
	err := suite.handler.Test(c)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, rec.Code)
	assert.Contains(suite.T(), rec.Body.String(), "API is working")

	// Verify mock was called
	suite.apiService.AssertExpectations(suite.T())
}

func TestAPIHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(APIHandlerTestSuite))
}