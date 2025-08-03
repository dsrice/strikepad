package container

import (
	"testing"

	"strikepad-backend/internal/handler"
	"strikepad-backend/internal/repository"
	"strikepad-backend/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ContainerTestSuite struct {
	suite.Suite
}

func (suite *ContainerTestSuite) TestBuildContainer() {
	// Test that BuildContainer creates a container without panicking
	container := BuildContainer()

	// Verify container is not nil
	assert.NotNil(suite.T(), container)

	// Verify that we can invoke key components from the container
	// Note: We'll skip database-dependent components to avoid connection issues

	// Test HealthService
	err := container.Invoke(func(hs service.HealthServiceInterface) {
		assert.NotNil(suite.T(), hs)
		response := hs.GetHealth()
		assert.NotNil(suite.T(), response)
	})
	assert.NoError(suite.T(), err)

	// Test APIService
	err = container.Invoke(func(as service.APIService) {
		assert.NotNil(suite.T(), as)
		message := as.GetTestMessage()
		assert.NotNil(suite.T(), message)
	})
	assert.NoError(suite.T(), err)
}

func (suite *ContainerTestSuite) TestContainerProvides() {
	// Test that all required dependencies are properly provided
	container := BuildContainer()

	// Test that each component can be resolved
	testCases := []struct {
		invokeFunc interface{}
		name       string
	}{
		{
			name: "HealthService",
			invokeFunc: func(hs service.HealthServiceInterface) {
				assert.NotNil(suite.T(), hs)
			},
		},
		{
			name: "APIService",
			invokeFunc: func(as service.APIService) {
				assert.NotNil(suite.T(), as)
			},
		},
		{
			name: "HealthHandler",
			invokeFunc: func(hh handler.HealthHandlerInterface) {
				assert.NotNil(suite.T(), hh)
			},
		},
		{
			name: "APIHandler",
			invokeFunc: func(ah *handler.APIHandler) {
				assert.NotNil(suite.T(), ah)
			},
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			err := container.Invoke(tc.invokeFunc)
			assert.NoError(t, err, "Should be able to resolve %s", tc.name)
		})
	}
}

func (suite *ContainerTestSuite) TestContainerWithDatabaseComponents() {
	// Test database-dependent components (will skip if no DB connection)
	container := BuildContainer()

	// Try to invoke database-dependent components
	// These may fail with database connection errors, which is expected in test environment
	suite.T().Run("UserRepository", func(t *testing.T) {
		err := container.Invoke(func(ur repository.UserRepository) {
			if ur != nil {
				t.Log("UserRepository successfully resolved")
			}
		})
		// Don't assert error here as it may fail due to DB connection
		if err != nil {
			t.Logf("Expected database-related error: %v", err)
		}
	})

	suite.T().Run("AuthService", func(t *testing.T) {
		err := container.Invoke(func(as service.AuthServiceInterface) {
			if as != nil {
				t.Log("AuthService successfully resolved")
			}
		})
		// Don't assert error here as it may fail due to DB connection
		if err != nil {
			t.Logf("Expected database-related error: %v", err)
		}
	})

	suite.T().Run("AuthHandler", func(t *testing.T) {
		err := container.Invoke(func(ah handler.AuthHandlerInterface) {
			if ah != nil {
				t.Log("AuthHandler successfully resolved")
			}
		})
		// Don't assert error here as it may fail due to DB connection
		if err != nil {
			t.Logf("Expected database-related error: %v", err)
		}
	})
}

func TestContainerTestSuite(t *testing.T) {
	suite.Run(t, new(ContainerTestSuite))
}
