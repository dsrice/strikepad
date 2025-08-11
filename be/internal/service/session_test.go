package service_test

import (
	"errors"
	"os"
	"testing"
	"time"

	"strikepad-backend/internal/auth"
	"strikepad-backend/internal/model"
	"strikepad-backend/internal/repository/mocks"
	"strikepad-backend/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type SessionServiceTestSuite struct {
	suite.Suite
	sessionService  service.SessionServiceInterface
	mockSessionRepo *mocks.MockSessionRepository
	jwtService      *auth.JWTService
}

func (suite *SessionServiceTestSuite) SetupTest() {
	// Set test JWT secret
	os.Setenv("JWT_SECRET_KEY", "test-secret-key-for-session-testing")

	suite.mockSessionRepo = new(mocks.MockSessionRepository)
	suite.jwtService = auth.NewJWTService()
	suite.sessionService = service.NewSessionService(suite.mockSessionRepo, suite.jwtService)
}

func (suite *SessionServiceTestSuite) TearDownTest() {
	// Reset the mock for next test case only - don't assert expectations in TearDown
	suite.mockSessionRepo.ExpectedCalls = nil
	suite.mockSessionRepo.Calls = nil
	os.Unsetenv("JWT_SECRET_KEY")
}

func (suite *SessionServiceTestSuite) TestCreateSession() {
	testCases := []struct {
		mockSetup     func()
		name          string
		errorMessage  string
		userID        uint
		expectedError bool
	}{
		{
			name:   "Success",
			userID: 1,
			mockSetup: func() {
				suite.mockSessionRepo.On("Create", mock.MatchedBy(func(session *model.UserSession) bool {
					return session.UserID == 1 &&
						session.AccessToken != "" &&
						session.RefreshToken != "" &&
						!session.AccessTokenExpiresAt.IsZero() &&
						!session.RefreshTokenExpiresAt.IsZero() &&
						!session.IsDeleted
				})).Return(nil).Once()
			},
			expectedError: false,
		},
		{
			name:   "Repository error",
			userID: 2,
			mockSetup: func() {
				suite.mockSessionRepo.On("Create", mock.MatchedBy(func(session *model.UserSession) bool {
					return session.UserID == 2
				})).Return(errors.New("database error")).Once()
			},
			expectedError: true,
			errorMessage:  "failed to create session",
		},
		{
			name:   "Zero user ID",
			userID: 0,
			mockSetup: func() {
				suite.mockSessionRepo.On("Create", mock.MatchedBy(func(session *model.UserSession) bool {
					return session.UserID == 0
				})).Return(nil).Once()
			},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Reset mocks for this specific test case
			suite.mockSessionRepo.ExpectedCalls = nil
			suite.mockSessionRepo.Calls = nil

			// Setup mocks
			tc.mockSetup()

			// Execute
			tokenPair, err := suite.sessionService.CreateSession(tc.userID)

			// Assert
			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, tokenPair)
				if tc.errorMessage != "" {
					assert.Contains(t, err.Error(), tc.errorMessage)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tokenPair)
				assert.NotEmpty(t, tokenPair.AccessToken)
				assert.NotEmpty(t, tokenPair.RefreshToken)
				assert.True(t, tokenPair.AccessTokenExpiresAt.After(time.Now()))
				assert.True(t, tokenPair.RefreshTokenExpiresAt.After(time.Now()))
			}
		})
	}
}

func (suite *SessionServiceTestSuite) TestValidateAccessToken() {
	userID := uint(123)
	tokenPair, _ := suite.jwtService.GenerateTokenPair(userID)
	validSession := &model.UserSession{
		ID:                    1,
		UserID:                userID,
		AccessToken:           tokenPair.AccessToken,
		RefreshToken:          tokenPair.RefreshToken,
		AccessTokenExpiresAt:  tokenPair.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: tokenPair.RefreshTokenExpiresAt,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		IsDeleted:             false,
	}

	expiredSession := &model.UserSession{
		ID:                    2,
		UserID:                userID,
		AccessToken:           tokenPair.AccessToken,
		RefreshToken:          tokenPair.RefreshToken,
		AccessTokenExpiresAt:  time.Now().Add(-time.Hour), // Expired
		RefreshTokenExpiresAt: tokenPair.RefreshTokenExpiresAt,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		IsDeleted:             false,
	}

	testCases := []struct {
		mockSetup     func()
		name          string
		token         string
		errorMessage  string
		expectedUID   uint
		expectedError bool
	}{
		{
			name:  "Valid access token",
			token: tokenPair.AccessToken,
			mockSetup: func() {
				suite.mockSessionRepo.On("FindByAccessToken", tokenPair.AccessToken).Return(validSession, nil)
			},
			expectedError: false,
			expectedUID:   userID,
		},
		{
			name:  "Invalid JWT token",
			token: "invalid.jwt.token",
			mockSetup: func() {
				// No mock setup needed as JWT validation will fail first
			},
			expectedError: true,
			errorMessage:  "invalid access token",
		},
		{
			name:  "Session not found in database",
			token: tokenPair.AccessToken,
			mockSetup: func() {
				suite.mockSessionRepo.On("FindByAccessToken", tokenPair.AccessToken).Return(nil, gorm.ErrRecordNotFound)
			},
			expectedError: true,
			errorMessage:  "session not found",
		},
		{
			name:  "Expired session",
			token: tokenPair.AccessToken,
			mockSetup: func() {
				suite.mockSessionRepo.On("FindByAccessToken", tokenPair.AccessToken).Return(expiredSession, nil)
			},
			expectedError: true,
			errorMessage:  "session is expired or invalidated",
		},
		{
			name:  "User ID mismatch",
			token: tokenPair.AccessToken,
			mockSetup: func() {
				mismatchSession := &model.UserSession{
					ID:                    3,
					UserID:                999, // Different user ID
					AccessToken:           tokenPair.AccessToken,
					RefreshToken:          tokenPair.RefreshToken,
					AccessTokenExpiresAt:  tokenPair.AccessTokenExpiresAt,
					RefreshTokenExpiresAt: tokenPair.RefreshTokenExpiresAt,
					CreatedAt:             time.Now(),
					UpdatedAt:             time.Now(),
					IsDeleted:             false,
				}
				suite.mockSessionRepo.On("FindByAccessToken", tokenPair.AccessToken).Return(mismatchSession, nil)
			},
			expectedError: true,
			errorMessage:  "token user ID mismatch",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Reset mocks for this specific test case
			suite.mockSessionRepo.ExpectedCalls = nil
			suite.mockSessionRepo.Calls = nil

			// Setup mocks
			tc.mockSetup()

			// Execute
			session, err := suite.sessionService.ValidateAccessToken(tc.token)

			// Assert
			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, session)
				if tc.errorMessage != "" {
					assert.Contains(t, err.Error(), tc.errorMessage)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, session)
				assert.Equal(t, tc.expectedUID, session.UserID)
			}
		})
	}
}

func (suite *SessionServiceTestSuite) TestRefreshToken() {
	userID := uint(456)
	tokenPair, _ := suite.jwtService.GenerateTokenPair(userID)
	validSession := &model.UserSession{
		ID:                    1,
		UserID:                userID,
		AccessToken:           tokenPair.AccessToken,
		RefreshToken:          tokenPair.RefreshToken,
		AccessTokenExpiresAt:  tokenPair.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: tokenPair.RefreshTokenExpiresAt,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		IsDeleted:             false,
	}

	testCases := []struct {
		mockSetup     func()
		name          string
		refreshToken  string
		errorMessage  string
		expectedError bool
	}{
		{
			name:         "Success",
			refreshToken: tokenPair.RefreshToken,
			mockSetup: func() {
				suite.mockSessionRepo.On("FindByRefreshToken", tokenPair.RefreshToken).Return(validSession, nil).Once()
				suite.mockSessionRepo.On("Update", mock.AnythingOfType("*model.UserSession")).Return(nil).Once().Once()
			},
			expectedError: false,
		},
		{
			name:         "Invalid refresh token",
			refreshToken: "invalid.refresh.token",
			mockSetup: func() {
				// No mock setup needed as JWT validation will fail first
			},
			expectedError: true,
			errorMessage:  "invalid refresh token",
		},
		{
			name:         "Session not found",
			refreshToken: tokenPair.RefreshToken,
			mockSetup: func() {
				suite.mockSessionRepo.On("FindByRefreshToken", tokenPair.RefreshToken).Return(nil, gorm.ErrRecordNotFound).Once()
			},
			expectedError: true,
			errorMessage:  "session not found",
		},
		{
			name:         "Update session error",
			refreshToken: tokenPair.RefreshToken,
			mockSetup: func() {
				suite.mockSessionRepo.On("FindByRefreshToken", tokenPair.RefreshToken).Return(validSession, nil)
				suite.mockSessionRepo.On("Update", mock.AnythingOfType("*model.UserSession")).Return(errors.New("update error")).Once()
			},
			expectedError: true,
			errorMessage:  "failed to update session",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Reset mocks for this specific test case
			suite.mockSessionRepo.ExpectedCalls = nil
			suite.mockSessionRepo.Calls = nil

			// Setup mocks
			tc.mockSetup()

			// Execute
			newTokenPair, err := suite.sessionService.RefreshToken(tc.refreshToken)

			// Assert
			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, newTokenPair)
				if tc.errorMessage != "" {
					assert.Contains(t, err.Error(), tc.errorMessage)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, newTokenPair)
				assert.NotEmpty(t, newTokenPair.AccessToken)
				assert.NotEmpty(t, newTokenPair.RefreshToken)
				// New tokens should be different from original (check lengths are reasonable)
				assert.True(t, newTokenPair.AccessToken != "")
				assert.True(t, newTokenPair.RefreshToken != "")
				// Validate that new tokens are structurally valid JWTs
				_, err = suite.jwtService.ValidateAccessToken(newTokenPair.AccessToken)
				assert.NoError(t, err)
				_, err = suite.jwtService.ValidateRefreshToken(newTokenPair.RefreshToken)
				assert.NoError(t, err)
			}
		})
	}
}

func (suite *SessionServiceTestSuite) TestLogout() {
	userID := uint(789)
	accessToken := "test-access-token"
	validSession := &model.UserSession{
		ID:          1,
		UserID:      userID,
		AccessToken: accessToken,
		IsDeleted:   false,
	}

	testCases := []struct {
		mockSetup     func()
		name          string
		accessToken   string
		errorMessage  string
		userID        uint
		expectedError bool
	}{
		{
			name:        "Success",
			userID:      userID,
			accessToken: accessToken,
			mockSetup: func() {
				suite.mockSessionRepo.On("FindByAccessToken", accessToken).Return(validSession, nil).Once()
				suite.mockSessionRepo.On("Update", mock.AnythingOfType("*model.UserSession")).Return(nil).Once().Once()
			},
			expectedError: false,
		},
		{
			name:        "Session not found",
			userID:      userID,
			accessToken: "nonexistent-token",
			mockSetup: func() {
				suite.mockSessionRepo.On("FindByAccessToken", "nonexistent-token").Return(nil, gorm.ErrRecordNotFound).Once()
			},
			expectedError: true,
			errorMessage:  "session not found",
		},
		{
			name:        "Session belongs to different user",
			userID:      999, // Different user ID
			accessToken: accessToken,
			mockSetup: func() {
				suite.mockSessionRepo.On("FindByAccessToken", accessToken).Return(validSession, nil)
			},
			expectedError: true,
			errorMessage:  "session does not belong to user",
		},
		{
			name:        "Update session error",
			userID:      userID,
			accessToken: accessToken,
			mockSetup: func() {
				suite.mockSessionRepo.On("FindByAccessToken", accessToken).Return(validSession, nil).Once()
				suite.mockSessionRepo.On("Update", mock.AnythingOfType("*model.UserSession")).Return(errors.New("update error")).Once().Once()
			},
			expectedError: true,
			errorMessage:  "failed to logout session",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Reset mocks for this specific test case
			suite.mockSessionRepo.ExpectedCalls = nil
			suite.mockSessionRepo.Calls = nil

			// Setup mocks
			tc.mockSetup()

			// Execute
			err := suite.sessionService.Logout(tc.userID, tc.accessToken)

			// Assert
			if tc.expectedError {
				assert.Error(t, err)
				if tc.errorMessage != "" && err != nil {
					assert.Contains(t, err.Error(), tc.errorMessage)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func (suite *SessionServiceTestSuite) TestInvalidateAllUserSessions() {
	testCases := []struct {
		mockSetup     func()
		name          string
		errorMessage  string
		userID        uint
		expectedError bool
	}{
		{
			name:   "Success",
			userID: 1,
			mockSetup: func() {
				suite.mockSessionRepo.On("InvalidateByUserID", uint(1)).Return(nil).Once()
			},
			expectedError: false,
		},
		{
			name:   "Repository error",
			userID: 2,
			mockSetup: func() {
				suite.mockSessionRepo.On("InvalidateByUserID", uint(2)).Return(errors.New("database error")).Once()
			},
			expectedError: true,
			errorMessage:  "failed to invalidate all user sessions",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Reset mocks for this specific test case
			suite.mockSessionRepo.ExpectedCalls = nil
			suite.mockSessionRepo.Calls = nil

			// Setup mocks
			tc.mockSetup()

			// Execute
			err := suite.sessionService.InvalidateAllUserSessions(tc.userID)

			// Assert
			if tc.expectedError {
				assert.Error(t, err)
				if tc.errorMessage != "" && err != nil {
					assert.Contains(t, err.Error(), tc.errorMessage)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func (suite *SessionServiceTestSuite) TestCleanupExpiredSessions() {
	testCases := []struct {
		mockSetup     func()
		name          string
		errorMessage  string
		expectedError bool
	}{
		{
			name: "Success",
			mockSetup: func() {
				suite.mockSessionRepo.On("InvalidateExpiredSessions").Return(nil)
			},
			expectedError: false,
		},
		{
			name: "Repository error",
			mockSetup: func() {
				suite.mockSessionRepo.On("InvalidateExpiredSessions").Return(errors.New("cleanup error"))
			},
			expectedError: true,
			errorMessage:  "failed to cleanup expired sessions",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Reset mocks for this specific test case
			suite.mockSessionRepo.ExpectedCalls = nil
			suite.mockSessionRepo.Calls = nil

			// Setup mocks
			tc.mockSetup()

			// Execute
			err := suite.sessionService.CleanupExpiredSessions()

			// Assert
			if tc.expectedError {
				assert.Error(t, err)
				if tc.errorMessage != "" && err != nil {
					assert.Contains(t, err.Error(), tc.errorMessage)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func (suite *SessionServiceTestSuite) TestNewSessionService() {
	// Test that NewSessionService creates a valid service
	svc := service.NewSessionService(suite.mockSessionRepo, suite.jwtService)
	assert.NotNil(suite.T(), svc)
}

func TestSessionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SessionServiceTestSuite))
}