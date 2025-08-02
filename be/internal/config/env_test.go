package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const defaultVal = "default_value"

type EnvConfigTestSuite struct {
	suite.Suite
	originalEnvVars map[string]string
}

func (suite *EnvConfigTestSuite) SetupTest() {
	// Save original environment variables
	suite.originalEnvVars = make(map[string]string)
	envVars := []string{"TEST_KEY", "EMPTY_KEY", "WHITESPACE_KEY"}

	for _, envVar := range envVars {
		if value, exists := os.LookupEnv(envVar); exists {
			suite.originalEnvVars[envVar] = value
		}
	}

	// Clear environment variables for clean test state
	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
}

func (suite *EnvConfigTestSuite) TearDownTest() {
	// Restore original environment variables
	envVars := []string{"TEST_KEY", "EMPTY_KEY", "WHITESPACE_KEY"}

	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}

	for envVar, value := range suite.originalEnvVars {
		os.Setenv(envVar, value)
	}
}

func (suite *EnvConfigTestSuite) TestGetEnvWithValue() {
	key := "TEST_KEY"
	expectedValue := "test_value"
	defaultVal := defaultVal

	os.Setenv(key, expectedValue)

	result := getEnv(key, defaultVal)
	assert.Equal(suite.T(), expectedValue, result)
}

func (suite *EnvConfigTestSuite) TestGetEnvWithDefault() {
	key := "NONEXISTENT_KEY"
	defaultVal := defaultVal

	// Ensure the key doesn't exist
	os.Unsetenv(key)

	result := getEnv(key, defaultVal)
	assert.Equal(suite.T(), defaultVal, result)
}

func (suite *EnvConfigTestSuite) TestGetEnvWithEmptyValue() {
	key := "EMPTY_KEY"
	defaultVal := defaultVal

	os.Setenv(key, "")

	result := getEnv(key, defaultVal)
	assert.Equal(suite.T(), defaultVal, result)
}

func (suite *EnvConfigTestSuite) TestGetEnvWithWhitespaceValue() {
	key := "WHITESPACE_KEY"
	value := "  value_with_spaces  "
	defaultVal := defaultVal

	os.Setenv(key, value)

	result := getEnv(key, defaultVal)
	assert.Equal(suite.T(), value, result) // getEnv doesn't trim whitespace
}

func (suite *EnvConfigTestSuite) TestGetEnvVariousValues() {
	testCases := []struct {
		name       string
		envValue   string
		defaultVal string
		expected   string
	}{
		{"Normal value", "hello", "default", "hello"},
		{"Numeric value", "123", "default", "123"},
		{"Special chars", "user@domain.com", "default", "user@domain.com"},
		{"Path value", "/path/to/file", "default", "/path/to/file"},
		{"URL value", "https://example.com", "default", "https://example.com"},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			key := "TEST_VAR_" + tc.name
			os.Setenv(key, tc.envValue)

			result := getEnv(key, tc.defaultVal)
			assert.Equal(t, tc.expected, result)

			os.Unsetenv(key)
		})
	}
}

func (suite *EnvConfigTestSuite) TestGetEnvDefaultValues() {
	// Test common default patterns
	testCases := []struct {
		key        string
		defaultVal string
	}{
		{"DB_HOST", "localhost"},
		{"DB_PORT", "5432"},
		{"DB_USER", "postgres"},
		{"ENV", "development"},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.key, func(t *testing.T) {
			// Ensure env var is not set
			os.Unsetenv(tc.key)

			result := getEnv(tc.key, tc.defaultVal)
			assert.Equal(t, tc.defaultVal, result)
		})
	}
}

func TestEnvConfigTestSuite(t *testing.T) {
	suite.Run(t, new(EnvConfigTestSuite))
}
