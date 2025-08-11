package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"strikepad-backend/internal/middleware"
	"strikepad-backend/internal/model"
	servicemocks "strikepad-backend/internal/service/mocks"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthMiddlewareTestSuite struct {
	suite.Suite
	echo           *echo.Echo
	mockSessionSvc *servicemocks.MockSessionServiceInterface
}

func (suite *AuthMiddlewareTestSuite) SetupTest() {
	suite.echo = echo.New()
	suite.mockSessionSvc = new(servicemocks.MockSessionServiceInterface)
}

func (suite *AuthMiddlewareTestSuite) TearDownTest() {
	// Reset mocks
	suite.mockSessionSvc.ExpectedCalls = nil
	suite.mockSessionSvc.Calls = nil
}

func (suite *AuthMiddlewareTestSuite) TestJWTMiddleware() {
	testCases := []struct {
		setupRequest    func(req *http.Request)
		setupMocks      func()
		expectedError   map[string]string
		validateContext func(t *testing.T, c echo.Context)
		name            string
		description     string
		expectedStatus  int
		expectNext      bool
	}{
		{
			name:        "Valid Bearer token",
			description: "Should authenticate successfully with valid Bearer token",
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer valid-access-token")
			},
			setupMocks: func() {
				session := &model.UserSession{
					ID:     1,
					UserID: 123,
				}
				suite.mockSessionSvc.On("ValidateAccessToken", "valid-access-token").
					Return(session, nil)
			},
			expectedStatus: http.StatusOK,
			expectNext:     true,
			validateContext: func(t *testing.T, c echo.Context) {
				session, exists := c.Get("session").(*model.UserSession)
				assert.True(t, exists, "Session should be set in context")
				assert.NotNil(t, session, "Session should not be nil")
				assert.Equal(t, uint(123), session.UserID, "User ID should match")

				userID, exists := c.Get("user_id").(uint)
				assert.True(t, exists, "User ID should be set in context")
				assert.Equal(t, uint(123), userID, "User ID should match")

				token, exists := c.Get("access_token").(string)
				assert.True(t, exists, "Access token should be set in context")
				assert.Equal(t, "valid-access-token", token, "Access token should match")
			},
		},
		{
			name:        "Missing Authorization header",
			description: "Should return 401 when Authorization header is missing",
			setupRequest: func(_ *http.Request) {
				// No Authorization header
			},
			setupMocks:     func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedError: map[string]string{
				"code":    "E005",
				"message": "Unauthorized",
			},
			expectNext: false,
		},
		{
			name:        "Empty Authorization header",
			description: "Should return 401 when Authorization header is empty",
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "")
			},
			setupMocks:     func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedError: map[string]string{
				"code":    "E005",
				"message": "Unauthorized",
			},
			expectNext: false,
		},
		{
			name:        "Invalid Authorization header format - no Bearer",
			description: "Should return 401 when Authorization header doesn't start with Bearer",
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "Token some-token")
			},
			setupMocks:     func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedError: map[string]string{
				"code":    "E005",
				"message": "Invalid authorization header format",
			},
			expectNext: false,
		},
		{
			name:        "Invalid Authorization header format - missing token",
			description: "Should return 401 when Authorization header has Bearer but no token",
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer")
			},
			setupMocks:     func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedError: map[string]string{
				"code":    "E005",
				"message": "Invalid authorization header format",
			},
			expectNext: false,
		},
		{
			name:        "Invalid Authorization header format - extra parts",
			description: "Should return 401 when Authorization header has more than 2 parts",
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer token extra")
			},
			setupMocks:     func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedError: map[string]string{
				"code":    "E005",
				"message": "Invalid authorization header format",
			},
			expectNext: false,
		},
		{
			name:        "Invalid access token",
			description: "Should return 401 when access token is invalid",
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer invalid-token")
			},
			setupMocks: func() {
				suite.mockSessionSvc.On("ValidateAccessToken", "invalid-token").
					Return(nil, errors.New("invalid token"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError: map[string]string{
				"code":    "E005",
				"message": "Invalid or expired token",
			},
			expectNext: false,
		},
		{
			name:        "Expired access token",
			description: "Should return 401 when access token is expired",
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer expired-token")
			},
			setupMocks: func() {
				suite.mockSessionSvc.On("ValidateAccessToken", "expired-token").
					Return(nil, errors.New("token expired"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError: map[string]string{
				"code":    "E005",
				"message": "Invalid or expired token",
			},
			expectNext: false,
		},
		{
			name:        "Session not found",
			description: "Should return 401 when session is not found",
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer token-no-session")
			},
			setupMocks: func() {
				suite.mockSessionSvc.On("ValidateAccessToken", "token-no-session").
					Return(nil, errors.New("session not found"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError: map[string]string{
				"code":    "E005",
				"message": "Invalid or expired token",
			},
			expectNext: false,
		},
		{
			name:        "Bearer with lowercase",
			description: "Should handle case sensitivity correctly",
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "bearer valid-token")
			},
			setupMocks:     func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedError: map[string]string{
				"code":    "E005",
				"message": "Invalid authorization header format",
			},
			expectNext: false,
		},
		{
			name:        "Valid token with special characters",
			description: "Should handle tokens with special characters",
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer valid-token-123.abc_def")
			},
			setupMocks: func() {
				session := &model.UserSession{
					ID:     2,
					UserID: 456,
				}
				suite.mockSessionSvc.On("ValidateAccessToken", "valid-token-123.abc_def").
					Return(session, nil)
			},
			expectedStatus: http.StatusOK,
			expectNext:     true,
			validateContext: func(t *testing.T, c echo.Context) {
				userID, exists := c.Get("user_id").(uint)
				assert.True(t, exists, "User ID should be set in context")
				assert.Equal(t, uint(456), userID, "User ID should match")
			},
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Reset mocks for this test case
			suite.mockSessionSvc.ExpectedCalls = nil
			suite.mockSessionSvc.Calls = nil

			// Setup mocks
			tc.setupMocks()

			// Create middleware
			middleware := middleware.JWTMiddleware(suite.mockSessionSvc)

			// Create test handler
			nextCalled := false
			testHandler := func(c echo.Context) error {
				nextCalled = true
				if tc.validateContext != nil {
					tc.validateContext(t, c)
				}
				return c.JSON(http.StatusOK, map[string]string{"message": "success"})
			}

			// Create HTTP request
			req := httptest.NewRequest(http.MethodGet, "/protected", http.NoBody)
			tc.setupRequest(req)
			rec := httptest.NewRecorder()
			c := suite.echo.NewContext(req, rec)

			// Execute middleware
			handler := middleware(testHandler)
			err := handler(c)

			// Validate results
			assert.NoError(t, err, tc.description)
			assert.Equal(t, tc.expectedStatus, rec.Code, tc.description)
			assert.Equal(t, tc.expectNext, nextCalled, "Next handler call expectation: %s", tc.description)

			// Check error response
			if tc.expectedError != nil {
				// Parse JSON response
				assert.Contains(t, rec.Body.String(), tc.expectedError["code"], "Error code should match")
				assert.Contains(t, rec.Body.String(), tc.expectedError["message"], "Error message should match")
			}
		})
	}
}

func (suite *AuthMiddlewareTestSuite) TestGetUserIDFromContext() {
	testCases := []struct {
		setupContext func(c echo.Context)
		name         string
		description  string
		expectedID   uint
		expectedOK   bool
	}{
		{
			name:        "Valid user ID in context",
			description: "Should extract valid user ID from context",
			setupContext: func(c echo.Context) {
				c.Set("user_id", uint(123))
			},
			expectedID: 123,
			expectedOK: true,
		},
		{
			name:        "Missing user ID in context",
			description: "Should return false when user_id is not in context",
			setupContext: func(_ echo.Context) {
				// Don't set user_id
			},
			expectedID: 0,
			expectedOK: false,
		},
		{
			name:        "Wrong type in context",
			description: "Should return false when user_id is wrong type",
			setupContext: func(c echo.Context) {
				c.Set("user_id", "123") // string instead of uint
			},
			expectedID: 0,
			expectedOK: false,
		},
		{
			name:        "Zero user ID",
			description: "Should handle zero user ID correctly",
			setupContext: func(c echo.Context) {
				c.Set("user_id", uint(0))
			},
			expectedID: 0,
			expectedOK: true,
		},
		{
			name:        "Large user ID",
			description: "Should handle large user ID correctly",
			setupContext: func(c echo.Context) {
				c.Set("user_id", uint(999999))
			},
			expectedID: 999999,
			expectedOK: true,
		},
		{
			name:        "Nil value in context",
			description: "Should return false when user_id is nil",
			setupContext: func(c echo.Context) {
				c.Set("user_id", nil)
			},
			expectedID: 0,
			expectedOK: false,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Create context
			req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			rec := httptest.NewRecorder()
			c := suite.echo.NewContext(req, rec)

			// Setup context
			tc.setupContext(c)

			// Execute function
			userID, ok := middleware.GetUserIDFromContext(c)

			// Validate results
			assert.Equal(t, tc.expectedID, userID, tc.description)
			assert.Equal(t, tc.expectedOK, ok, tc.description)
		})
	}
}

func (suite *AuthMiddlewareTestSuite) TestGetAccessTokenFromContext() {
	testCases := []struct {
		name          string
		description   string
		setupContext  func(c echo.Context)
		expectedToken string
		expectedOK    bool
	}{
		{
			name:        "Valid access token in context",
			description: "Should extract valid access token from context",
			setupContext: func(c echo.Context) {
				c.Set("access_token", "test-access-token")
			},
			expectedToken: "test-access-token",
			expectedOK:    true,
		},
		{
			name:        "Missing access token in context",
			description: "Should return false when access_token is not in context",
			setupContext: func(_ echo.Context) {
				// Don't set access_token
			},
			expectedToken: "",
			expectedOK:    false,
		},
		{
			name:        "Wrong type in context",
			description: "Should return false when access_token is wrong type",
			setupContext: func(c echo.Context) {
				c.Set("access_token", 12345) // int instead of string
			},
			expectedToken: "",
			expectedOK:    false,
		},
		{
			name:        "Empty access token",
			description: "Should handle empty access token correctly",
			setupContext: func(c echo.Context) {
				c.Set("access_token", "")
			},
			expectedToken: "",
			expectedOK:    true,
		},
		{
			name:        "Long access token",
			description: "Should handle long access token correctly",
			setupContext: func(c echo.Context) {
				longToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
				c.Set("access_token", longToken)
			},
			expectedToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			expectedOK:    true,
		},
		{
			name:        "Token with special characters",
			description: "Should handle token with special characters correctly",
			setupContext: func(c echo.Context) {
				c.Set("access_token", "token-123.abc_def+ghi/jkl==")
			},
			expectedToken: "token-123.abc_def+ghi/jkl==",
			expectedOK:    true,
		},
		{
			name:        "Nil value in context",
			description: "Should return false when access_token is nil",
			setupContext: func(c echo.Context) {
				c.Set("access_token", nil)
			},
			expectedToken: "",
			expectedOK:    false,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Create context
			req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			rec := httptest.NewRecorder()
			c := suite.echo.NewContext(req, rec)

			// Setup context
			tc.setupContext(c)

			// Execute function
			token, ok := middleware.GetAccessTokenFromContext(c)

			// Validate results
			assert.Equal(t, tc.expectedToken, token, tc.description)
			assert.Equal(t, tc.expectedOK, ok, tc.description)
		})
	}
}

func (suite *AuthMiddlewareTestSuite) TestMiddlewareIntegration() {
	testCases := []struct {
		setupRequest func(req *http.Request)
		setupMocks   func()
		testFlow     func(t *testing.T, c echo.Context)
		name         string
		description  string
	}{
		{
			name:        "Complete authentication flow",
			description: "Should handle complete authentication flow correctly",
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer integration-token")
			},
			setupMocks: func() {
				session := &model.UserSession{
					ID:          10,
					UserID:      999,
					AccessToken: "integration-token",
				}
				suite.mockSessionSvc.On("ValidateAccessToken", "integration-token").
					Return(session, nil)
			},
			testFlow: func(t *testing.T, c echo.Context) {
				// Test all helper functions work together
				userID, userOK := middleware.GetUserIDFromContext(c)
				token, tokenOK := middleware.GetAccessTokenFromContext(c)
				session := c.Get("session").(*model.UserSession)

				assert.True(t, userOK, "Should extract user ID successfully")
				assert.True(t, tokenOK, "Should extract token successfully")
				assert.Equal(t, uint(999), userID, "User ID should match")
				assert.Equal(t, "integration-token", token, "Token should match")
				assert.NotNil(t, session, "Session should be available")
				assert.Equal(t, uint(10), session.ID, "Session ID should match")
			},
		},
		{
			name:        "Multiple middleware calls consistency",
			description: "Should maintain consistency across multiple calls",
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer consistent-token")
			},
			setupMocks: func() {
				session := &model.UserSession{
					ID:     5,
					UserID: 555,
				}
				suite.mockSessionSvc.On("ValidateAccessToken", "consistent-token").
					Return(session, nil)
			},
			testFlow: func(t *testing.T, c echo.Context) {
				// Call helper functions multiple times
				for i := 0; i < 3; i++ {
					userID, userOK := middleware.GetUserIDFromContext(c)
					token, tokenOK := middleware.GetAccessTokenFromContext(c)

					assert.True(t, userOK, "User ID extraction should be consistent, iteration %d", i+1)
					assert.True(t, tokenOK, "Token extraction should be consistent, iteration %d", i+1)
					assert.Equal(t, uint(555), userID, "User ID should be consistent, iteration %d", i+1)
					assert.Equal(t, "consistent-token", token, "Token should be consistent, iteration %d", i+1)
				}
			},
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Reset mocks
			suite.mockSessionSvc.ExpectedCalls = nil
			suite.mockSessionSvc.Calls = nil

			// Setup mocks
			tc.setupMocks()

			// Create middleware
			middleware := middleware.JWTMiddleware(suite.mockSessionSvc)

			// Create test handler
			testHandler := func(c echo.Context) error {
				tc.testFlow(t, c)
				return c.JSON(http.StatusOK, map[string]string{"message": "success"})
			}

			// Create HTTP request
			req := httptest.NewRequest(http.MethodGet, "/protected", http.NoBody)
			tc.setupRequest(req)
			rec := httptest.NewRecorder()
			c := suite.echo.NewContext(req, rec)

			// Execute middleware
			handler := middleware(testHandler)
			err := handler(c)

			// Validate
			assert.NoError(t, err, tc.description)
			assert.Equal(t, http.StatusOK, rec.Code, tc.description)
		})
	}
}

func TestAuthMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareTestSuite))
}
