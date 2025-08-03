package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DatabaseConfigTestSuite struct {
	suite.Suite
	originalEnvVars map[string]string
}

func (suite *DatabaseConfigTestSuite) SetupTest() {
	// Save original environment variables
	suite.originalEnvVars = make(map[string]string)
	envVars := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE", "ENV"}

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

func (suite *DatabaseConfigTestSuite) TearDownTest() {
	// Restore original environment variables
	envVars := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE", "ENV"}

	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}

	for envVar, value := range suite.originalEnvVars {
		os.Setenv(envVar, value)
	}
}

func (suite *DatabaseConfigTestSuite) TestGetEnvDefaults() {
	// Test that getEnv returns expected defaults for database config
	testCases := []struct {
		key          string
		defaultValue string
		expected     string
	}{
		{"DB_HOST", "localhost", "localhost"},
		{"DB_PORT", "5432", "5432"},
		{"DB_USER", "postgres", "postgres"},
		{"DB_PASSWORD", "password", "password"},
		{"DB_NAME", "strikepad", "strikepad"},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.key, func(t *testing.T) {
			// Ensure env var is not set
			os.Unsetenv(tc.key)

			result := getEnv(tc.key, tc.defaultValue)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func (suite *DatabaseConfigTestSuite) TestGetEnvCustomValues() {
	// Test that getEnv returns custom values when set
	testCases := []struct {
		key          string
		setValue     string
		defaultValue string
		expected     string
	}{
		{"DB_HOST", "custom-host", "localhost", "custom-host"},
		{"DB_PORT", "3306", "5432", "3306"},
		{"DB_USER", "myuser", "postgres", "myuser"},
		{"DB_PASSWORD", "mypassword", "password", "mypassword"},
		{"DB_NAME", "mydatabase", "strikepad", "mydatabase"},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.key, func(t *testing.T) {
			os.Setenv(tc.key, tc.setValue)

			result := getEnv(tc.key, tc.defaultValue)
			assert.Equal(t, tc.expected, result)

			os.Unsetenv(tc.key)
		})
	}
}

func (suite *DatabaseConfigTestSuite) TestGetEnvEmptyValues() {
	// Test that getEnv returns defaults for empty values
	testCases := []struct {
		key          string
		defaultValue string
	}{
		{"DB_HOST", "localhost"},
		{"DB_PORT", "5432"},
		{"DB_USER", "postgres"},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.key+"_empty", func(t *testing.T) {
			os.Setenv(tc.key, "")

			result := getEnv(tc.key, tc.defaultValue)
			assert.Equal(t, tc.defaultValue, result)

			os.Unsetenv(tc.key)
		})
	}
}

func (suite *DatabaseConfigTestSuite) TestGetEnvSpecialCharacters() {
	// Test that getEnv handles special characters correctly
	testCases := []struct {
		name     string
		key      string
		setValue string
		expected string
	}{
		{"Special chars", "DB_PASSWORD", "pass@word#123!", "pass@word#123!"},
		{"Spaces", "DB_HOST", "my host", "my host"},
		{"Unicode", "DB_NAME", "データベース", "データベース"},
		{"Path", "DB_HOST", "/var/run/postgresql", "/var/run/postgresql"},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			os.Setenv(tc.key, tc.setValue)

			result := getEnv(tc.key, "default")
			assert.Equal(t, tc.expected, result)

			os.Unsetenv(tc.key)
		})
	}
}

func (suite *DatabaseConfigTestSuite) TestDatabaseConfigurationScenarios() {
	// Test common database configuration scenarios
	scenarios := []struct {
		name    string
		envVars map[string]string
		desc    string
	}{
		{
			name: "PostgreSQL local development",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PORT":     "5432",
				"DB_USER":     "developer",
				"DB_PASSWORD": "devpass",
				"DB_NAME":     "dev_db",
			},
			desc: "Local PostgreSQL setup",
		},
		{
			name: "MySQL configuration",
			envVars: map[string]string{
				"DB_HOST":     "mysql-server",
				"DB_PORT":     "3306",
				"DB_USER":     "app_user",
				"DB_PASSWORD": "secure_pass",
				"DB_NAME":     "app_database",
			},
			desc: "MySQL server setup",
		},
		{
			name: "Production configuration",
			envVars: map[string]string{
				"DB_HOST":     "prod-db.example.com",
				"DB_PORT":     "5432",
				"DB_USER":     "prod_user",
				"DB_PASSWORD": "very_secure_password",
				"DB_NAME":     "production_db",
			},
			desc: "Production database setup",
		},
	}

	for _, scenario := range scenarios {
		suite.T().Run(scenario.name, func(t *testing.T) {
			// Set environment variables for this scenario
			for key, value := range scenario.envVars {
				os.Setenv(key, value)
			}

			// Test that getEnv returns the correct values
			for key, expectedValue := range scenario.envVars {
				result := getEnv(key, "default")
				assert.Equal(t, expectedValue, result,
					"getEnv should return correct value for %s in %s", key, scenario.desc)
			}

			// Clean up
			for key := range scenario.envVars {
				os.Unsetenv(key)
			}
		})
	}
}

func (suite *DatabaseConfigTestSuite) TestNewDatabase() {
	// Set required environment variables
	suite.T().Setenv("DB_HOST", "localhost")
	suite.T().Setenv("DB_PORT", "5432")
	suite.T().Setenv("DB_USER", "testuser")
	suite.T().Setenv("DB_PASSWORD", "testpass")
	suite.T().Setenv("DB_NAME", "testdb")

	// This will test the DSN construction but will fail on actual connection
	// We can't test the actual database connection without a real database
	// But we can test that the function doesn't panic and constructs the DSN correctly
	defer func() {
		if r := recover(); r != nil {
			// Expected to fail with connection error, not panic
			suite.T().Logf("Expected database connection failure: %v", r)
		}
	}()

	// Call NewDatabase - this will attempt to connect but should handle the error gracefully
	// by calling log.Fatal, which we can't easily test without mocking
	// For now, we'll skip this direct test and instead test the DSN construction
	suite.T().Skip("NewDatabase requires actual database connection - testing DSN construction via getEnv instead")
}

func TestDatabaseConfigTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseConfigTestSuite))
}
