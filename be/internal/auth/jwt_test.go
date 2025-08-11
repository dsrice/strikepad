package auth_test

import (
	"os"
	"testing"
	"time"

	"strikepad-backend/internal/auth"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type JWTServiceTestSuite struct {
	suite.Suite
	jwtService *auth.JWTService
}

func (suite *JWTServiceTestSuite) SetupTest() {
	// Set a test secret key for consistency
	os.Setenv("JWT_SECRET_KEY", "test-secret-key-for-testing")
	suite.jwtService = auth.NewJWTService()
}

func (suite *JWTServiceTestSuite) TearDownTest() {
	// Clean up environment variable
	os.Unsetenv("JWT_SECRET_KEY")
}

func (suite *JWTServiceTestSuite) TestNewJWTService() {
	testCases := []struct {
		name      string
		secretKey string
		expectKey string
	}{
		{
			name:      "With environment variable",
			secretKey: "custom-secret-key",
			expectKey: "custom-secret-key",
		},
		{
			name:      "Without environment variable (default)",
			secretKey: "",
			expectKey: "your-secret-key-change-this-in-production",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Setup environment
			if tc.secretKey != "" {
				os.Setenv("JWT_SECRET_KEY", tc.secretKey)
			} else {
				os.Unsetenv("JWT_SECRET_KEY")
			}

			// Create service
			service := auth.NewJWTService()

			// Assert
			assert.NotNil(t, service)

			// Clean up
			os.Unsetenv("JWT_SECRET_KEY")
		})
	}
}

func (suite *JWTServiceTestSuite) TestGenerateTokenPair() {
	testCases := []struct {
		name           string
		userID         uint
		expectError    bool
		validateTokens bool
	}{
		{
			name:           "Valid user ID",
			userID:         1,
			expectError:    false,
			validateTokens: true,
		},
		{
			name:           "Another valid user ID",
			userID:         999,
			expectError:    false,
			validateTokens: true,
		},
		{
			name:           "Zero user ID",
			userID:         0,
			expectError:    false,
			validateTokens: true,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Execute
			tokenPair, err := suite.jwtService.GenerateTokenPair(tc.userID)

			// Assert
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, tokenPair)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tokenPair)
				assert.NotEmpty(t, tokenPair.AccessToken)
				assert.NotEmpty(t, tokenPair.RefreshToken)
				assert.True(t, tokenPair.AccessTokenExpiresAt.After(time.Now()))
				assert.True(t, tokenPair.RefreshTokenExpiresAt.After(time.Now()))

				// Access token should expire before refresh token
				assert.True(t, tokenPair.AccessTokenExpiresAt.Before(tokenPair.RefreshTokenExpiresAt))

				if tc.validateTokens {
					// Validate access token
					accessClaims, err := suite.jwtService.ValidateAccessToken(tokenPair.AccessToken)
					assert.NoError(t, err)
					assert.Equal(t, tc.userID, accessClaims.UserID)
					assert.Equal(t, "access", accessClaims.Type)

					// Validate refresh token
					refreshClaims, err := suite.jwtService.ValidateRefreshToken(tokenPair.RefreshToken)
					assert.NoError(t, err)
					assert.Equal(t, tc.userID, refreshClaims.UserID)
					assert.Equal(t, "refresh", refreshClaims.Type)
				}
			}
		})
	}
}

func (suite *JWTServiceTestSuite) TestValidateToken() {
	userID := uint(123)
	tokenPair, err := suite.jwtService.GenerateTokenPair(userID)
	assert.NoError(suite.T(), err)

	testCases := []struct {
		name         string
		token        string
		expectedType string
		expectedUID  uint
		expectError  bool
	}{
		{
			name:         "Valid access token",
			token:        tokenPair.AccessToken,
			expectError:  false,
			expectedUID:  userID,
			expectedType: "access",
		},
		{
			name:         "Valid refresh token",
			token:        tokenPair.RefreshToken,
			expectError:  false,
			expectedUID:  userID,
			expectedType: "refresh",
		},
		{
			name:        "Invalid token format",
			token:       "invalid.token.format",
			expectError: true,
		},
		{
			name:        "Empty token",
			token:       "",
			expectError: true,
		},
		{
			name:        "Malformed token",
			token:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid.signature",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Execute
			claims, err := suite.jwtService.ValidateToken(tc.token)

			// Assert
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, tc.expectedUID, claims.UserID)
				assert.Equal(t, tc.expectedType, claims.Type)
			}
		})
	}
}

func (suite *JWTServiceTestSuite) TestValidateAccessToken() {
	userID := uint(456)
	tokenPair, err := suite.jwtService.GenerateTokenPair(userID)
	assert.NoError(suite.T(), err)

	testCases := []struct {
		name        string
		token       string
		expectError bool
		expectedUID uint
	}{
		{
			name:        "Valid access token",
			token:       tokenPair.AccessToken,
			expectError: false,
			expectedUID: userID,
		},
		{
			name:        "Refresh token (should fail)",
			token:       tokenPair.RefreshToken,
			expectError: true,
		},
		{
			name:        "Invalid token",
			token:       "invalid.token",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Execute
			claims, err := suite.jwtService.ValidateAccessToken(tc.token)

			// Assert
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, tc.expectedUID, claims.UserID)
				assert.Equal(t, "access", claims.Type)
			}
		})
	}
}

func (suite *JWTServiceTestSuite) TestValidateRefreshToken() {
	userID := uint(789)
	tokenPair, err := suite.jwtService.GenerateTokenPair(userID)
	assert.NoError(suite.T(), err)

	testCases := []struct {
		name        string
		token       string
		expectError bool
		expectedUID uint
	}{
		{
			name:        "Valid refresh token",
			token:       tokenPair.RefreshToken,
			expectError: false,
			expectedUID: userID,
		},
		{
			name:        "Access token (should fail)",
			token:       tokenPair.AccessToken,
			expectError: true,
		},
		{
			name:        "Invalid token",
			token:       "invalid.token",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Execute
			claims, err := suite.jwtService.ValidateRefreshToken(tc.token)

			// Assert
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, tc.expectedUID, claims.UserID)
				assert.Equal(t, "refresh", claims.Type)
			}
		})
	}
}

func (suite *JWTServiceTestSuite) TestTokenExpiration() {
	// Test with a very short duration to test expiration
	testCases := []struct {
		name         string
		userID       uint
		waitDuration time.Duration
		expectValid  bool
	}{
		{
			name:         "Token not yet expired",
			userID:       1,
			waitDuration: 0,
			expectValid:  true,
		},
		// Note: We can't easily test actual expiration without waiting or mocking time
		// In a real test environment, you might use libraries like testify/mock to mock time
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Generate token
			tokenPair, err := suite.jwtService.GenerateTokenPair(tc.userID)
			assert.NoError(t, err)

			// Wait if specified
			if tc.waitDuration > 0 {
				time.Sleep(tc.waitDuration)
			}

			// Validate token
			claims, err := suite.jwtService.ValidateAccessToken(tokenPair.AccessToken)

			if tc.expectValid {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, tc.userID, claims.UserID)
			} else {
				assert.Error(t, err)
				assert.Nil(t, claims)
			}
		})
	}
}

func (suite *JWTServiceTestSuite) TestTokenWithDifferentSigningKey() {
	// Create another JWT service with different secret
	os.Setenv("JWT_SECRET_KEY", "different-secret-key")
	differentJWTService := auth.NewJWTService()

	// Generate token with original service
	userID := uint(123)
	tokenPair, err := suite.jwtService.GenerateTokenPair(userID)
	assert.NoError(suite.T(), err)

	// Try to validate with different service (should fail)
	_, err = differentJWTService.ValidateAccessToken(tokenPair.AccessToken)
	assert.Error(suite.T(), err)

	// Clean up
	os.Unsetenv("JWT_SECRET_KEY")
}

func TestJWTServiceTestSuite(t *testing.T) {
	suite.Run(t, new(JWTServiceTestSuite))
}
