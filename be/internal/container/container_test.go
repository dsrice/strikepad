package container_test

import (
	"testing"

	"strikepad-backend/internal/container"
	"strikepad-backend/internal/handler"
	"strikepad-backend/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/dig"
)

// buildTestContainer creates a container without database dependencies for testing
func buildTestContainer() *dig.Container {
	container := dig.New()

	// Only provide non-database-dependent components
	if err := container.Provide(service.NewHealthService); err != nil {
		panic(err)
	}
	if err := container.Provide(service.NewAPIService); err != nil {
		panic(err)
	}
	if err := container.Provide(handler.NewHealthHandler); err != nil {
		panic(err)
	}
	if err := container.Provide(handler.NewAPIHandler); err != nil {
		panic(err)
	}

	return container
}

type ContainerTestSuite struct {
	suite.Suite
}

func (suite *ContainerTestSuite) TestBuildContainer() {
	// Test that BuildContainer creates a container without panicking
	c := container.BuildContainer()

	// Verify container is not nil
	assert.NotNil(suite.T(), c)

	// Verify that we can invoke key components from the container
	// Note: We'll skip database-dependent components to avoid connection issues

	// Test HealthService
	err := c.Invoke(func(hs service.HealthServiceInterface) {
		assert.NotNil(suite.T(), hs)
		response := hs.GetHealth()
		assert.NotNil(suite.T(), response)
	})
	assert.NoError(suite.T(), err)

	// Test APIService
	err = c.Invoke(func(as service.APIService) {
		assert.NotNil(suite.T(), as)
		message := as.GetTestMessage()
		assert.NotNil(suite.T(), message)
	})
	assert.NoError(suite.T(), err)
}

func (suite *ContainerTestSuite) TestContainerProvides() {
	// Test that all required non-database dependencies are properly provided
	c := buildTestContainer()

	// Test that each non-database component can be resolved
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
			err := c.Invoke(tc.invokeFunc)
			assert.NoError(t, err, "Should be able to resolve %s", tc.name)
		})
	}
}

func (suite *ContainerTestSuite) TestContainerWithDatabaseComponents() {
	t := suite.T()
	t.Skip("Skipping database-dependent tests in test environment")

	// The following tests are skipped because they require a database connection
	// In a real environment with a database, these tests would verify that:
	// 1. UserRepository can be resolved from the container
	// 2. AuthService can be resolved from the container
	// 3. AuthHandler can be resolved from the container

	// For testing purposes, we use buildTestContainer() which doesn't include
	// database-dependent components to avoid connection errors
}

func TestContainerTestSuite(t *testing.T) {
	suite.Run(t, new(ContainerTestSuite))
}
