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

func (suite *EnvConfigTestSuite) TestGetEnv() {
	testCases := []struct {
		name         string
		key          string
		envValue     *string // nil means env var is not set
		defaultValue string
		expected     string
		description  string
	}{
		{
			name:         "existing value",
			key:          "TEST_KEY",
			envValue:     stringPtr("test_value"),
			defaultValue: defaultVal,
			expected:     "test_value",
			description:  "should return env value when set",
		},
		{
			name:         "nonexistent key",
			key:          "NONEXISTENT_KEY",
			envValue:     nil,
			defaultValue: defaultVal,
			expected:     defaultVal,
			description:  "should return default when env var doesn't exist",
		},
		{
			name:         "empty value",
			key:          "EMPTY_KEY",
			envValue:     stringPtr(""),
			defaultValue: defaultVal,
			expected:     defaultVal,
			description:  "should return default when env var is empty",
		},
		{
			name:         "whitespace value",
			key:          "WHITESPACE_KEY",
			envValue:     stringPtr("  value_with_spaces  "),
			defaultValue: defaultVal,
			expected:     "  value_with_spaces  ",
			description:  "should preserve whitespace in env values",
		},
		{
			name:         "numeric value",
			key:          "NUMERIC_KEY",
			envValue:     stringPtr("123456"),
			defaultValue: "0",
			expected:     "123456",
			description:  "should handle numeric strings correctly",
		},
		{
			name:         "special characters",
			key:          "SPECIAL_KEY",
			envValue:     stringPtr("user@domain.com!#$%"),
			defaultValue: "default@example.com",
			expected:     "user@domain.com!#$%",
			description:  "should handle special characters correctly",
		},
		{
			name:         "url value",
			key:          "URL_KEY",
			envValue:     stringPtr("https://api.example.com/v1"),
			defaultValue: "http://localhost",
			expected:     "https://api.example.com/v1",
			description:  "should handle URL strings correctly",
		},
		{
			name:         "path value",
			key:          "PATH_KEY",
			envValue:     stringPtr("/usr/local/bin:/usr/bin"),
			defaultValue: "/bin",
			expected:     "/usr/local/bin:/usr/bin",
			description:  "should handle file paths correctly",
		},
		{
			name:         "unicode value",
			key:          "UNICODE_KEY",
			envValue:     stringPtr("こんにちは世界"),
			defaultValue: "hello",
			expected:     "こんにちは世界",
			description:  "should handle unicode characters correctly",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Clean up any existing value
			os.Unsetenv(tc.key)

			// Set environment variable if specified
			if tc.envValue != nil {
				os.Setenv(tc.key, *tc.envValue)
			}

			result := getEnv(tc.key, tc.defaultValue)
			assert.Equal(t, tc.expected, result, tc.description)

			// Clean up
			os.Unsetenv(tc.key)
		})
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}

func (suite *EnvConfigTestSuite) TestGetEnvCommonDefaults() {
	// Test common application configuration patterns
	testCases := []struct {
		name         string
		key          string
		defaultValue string
		category     string
	}{
		// Database configuration
		{"database host", "DB_HOST", "localhost", "database"},
		{"database port", "DB_PORT", "5432", "database"},
		{"database user", "DB_USER", "postgres", "database"},
		{"database name", "DB_NAME", "strikepad", "database"},
		{"database ssl mode", "DB_SSLMODE", "disable", "database"},

		// Application configuration
		{"environment", "ENV", "development", "application"},
		{"port", "PORT", "8080", "application"},
		{"host", "HOST", "0.0.0.0", "application"},

		// Authentication configuration
		{"jwt secret", "JWT_SECRET", "default-secret", "auth"},
		{"jwt expiry", "JWT_EXPIRY", "24h", "auth"},

		// External service configuration
		{"redis url", "REDIS_URL", "redis://localhost:6379", "cache"},
		{"api base url", "API_BASE_URL", "http://localhost:8080", "external"},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Ensure env var is not set
			os.Unsetenv(tc.key)

			result := getEnv(tc.key, tc.defaultValue)
			assert.Equal(t, tc.defaultValue, result,
				"getEnv should return default value for %s configuration", tc.category)
		})
	}
}

func TestEnvConfigTestSuite(t *testing.T) {
	suite.Run(t, new(EnvConfigTestSuite))
}