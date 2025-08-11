package repository_test

import (
	"regexp"
	"testing"
	"time"

	"strikepad-backend/internal/model"
	"strikepad-backend/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type SessionRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	mock sqlmock.Sqlmock
	repo repository.SessionRepositoryInterface
}

func (suite *SessionRepositoryTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	assert.NoError(suite.T(), err)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(suite.T(), err)

	suite.db = gormDB
	suite.mock = mock
	suite.repo = repository.NewSessionRepository(gormDB)
}

func (suite *SessionRepositoryTestSuite) TearDownTest() {
	err := suite.mock.ExpectationsWereMet()
	assert.NoError(suite.T(), err)
}

func (suite *SessionRepositoryTestSuite) TestCreate() {
	testCases := []struct {
		name        string
		session     *model.UserSession
		mockSetup   func()
		expectError bool
		errorMsg    string
	}{
		{
			name: "Success",
			session: &model.UserSession{
				UserID:                1,
				AccessToken:           "test-access-token",
				RefreshToken:          "test-refresh-token",
				AccessTokenExpiresAt:  time.Now().Add(time.Hour),
				RefreshTokenExpiresAt: time.Now().Add(24 * time.Hour),
				IsDeleted:             false,
			},
			mockSetup: func() {
				suite.mock.ExpectBegin()
				suite.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `user_sessions`")).
					WithArgs(
						sqlmock.AnyArg(), // user_id
						sqlmock.AnyArg(), // access_token
						sqlmock.AnyArg(), // refresh_token
						sqlmock.AnyArg(), // access_token_expires_at
						sqlmock.AnyArg(), // refresh_token_expires_at
						sqlmock.AnyArg(), // is_deleted
						sqlmock.AnyArg(), // deleted_at
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				suite.mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "Database error",
			session: &model.UserSession{
				UserID:       2,
				AccessToken:  "test-access-token-2",
				RefreshToken: "test-refresh-token-2",
				IsDeleted:    false,
			},
			mockSetup: func() {
				suite.mock.ExpectBegin()
				suite.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `user_sessions`")).
					WillReturnError(assert.AnError)
				suite.mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "failed to create session",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Setup mock expectations
			tc.mockSetup()

			// Execute
			err := suite.repo.Create(tc.session)

			// Assert
			if tc.expectError {
				assert.Error(t, err)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func (suite *SessionRepositoryTestSuite) TestFindByAccessToken() {
	testCases := []struct {
		name        string
		accessToken string
		mockSetup   func()
		expectError bool
		errorMsg    string
		expectedUID uint
	}{
		{
			name:        "Success",
			accessToken: "valid-access-token",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "access_token", "refresh_token",
					"access_token_expires_at", "refresh_token_expires_at",
					"created_at", "updated_at", "is_deleted", "deleted_at",
				}).AddRow(
					1, 123, "valid-access-token", "refresh-token",
					time.Now().Add(time.Hour), time.Now().Add(24*time.Hour),
					time.Now(), time.Now(), false, nil,
				)

				suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_sessions`")).
					WithArgs("valid-access-token", false).
					WillReturnRows(rows)

				// Mock for User preload
				userRows := sqlmock.NewRows([]string{
					"id", "provider_type", "provider_user_id", "email",
					"display_name", "password_hash", "email_verified",
					"created_at", "updated_at", "is_deleted", "deleted_at",
				}).AddRow(
					123, "email", nil, "test@example.com",
					"Test User", nil, false,
					time.Now(), time.Now(), false, nil,
				)

				suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
					WithArgs(123).
					WillReturnRows(userRows)
			},
			expectError: false,
			expectedUID: 123,
		},
		{
			name:        "Session not found",
			accessToken: "nonexistent-token",
			mockSetup: func() {
				suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_sessions`")).
					WithArgs("nonexistent-token", false).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectError: true,
			errorMsg:    "session not found",
		},
		{
			name:        "Database error",
			accessToken: "error-token",
			mockSetup: func() {
				suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_sessions`")).
					WithArgs("error-token", false).
					WillReturnError(assert.AnError)
			},
			expectError: true,
			errorMsg:    "failed to find session by access token",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Setup mock expectations
			tc.mockSetup()

			// Execute
			session, err := suite.repo.FindByAccessToken(tc.accessToken)

			// Assert
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, session)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, session)
				assert.Equal(t, tc.expectedUID, session.UserID)
				assert.Equal(t, tc.accessToken, session.AccessToken)
			}
		})
	}
}

func (suite *SessionRepositoryTestSuite) TestFindByRefreshToken() {
	testCases := []struct {
		name         string
		refreshToken string
		mockSetup    func()
		expectError  bool
		errorMsg     string
		expectedUID  uint
	}{
		{
			name:         "Success",
			refreshToken: "valid-refresh-token",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "access_token", "refresh_token",
					"access_token_expires_at", "refresh_token_expires_at",
					"created_at", "updated_at", "is_deleted", "deleted_at",
				}).AddRow(
					1, 456, "access-token", "valid-refresh-token",
					time.Now().Add(time.Hour), time.Now().Add(24*time.Hour),
					time.Now(), time.Now(), false, nil,
				)

				suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_sessions`")).
					WithArgs("valid-refresh-token", false).
					WillReturnRows(rows)

				// Mock for User preload
				userRows := sqlmock.NewRows([]string{
					"id", "provider_type", "provider_user_id", "email",
					"display_name", "password_hash", "email_verified",
					"created_at", "updated_at", "is_deleted", "deleted_at",
				}).AddRow(
					456, "email", nil, "test2@example.com",
					"Test User 2", nil, false,
					time.Now(), time.Now(), false, nil,
				)

				suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
					WithArgs(456).
					WillReturnRows(userRows)
			},
			expectError: false,
			expectedUID: 456,
		},
		{
			name:         "Session not found",
			refreshToken: "nonexistent-refresh-token",
			mockSetup: func() {
				suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_sessions`")).
					WithArgs("nonexistent-refresh-token", false).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectError: true,
			errorMsg:    "session not found",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Setup mock expectations
			tc.mockSetup()

			// Execute
			session, err := suite.repo.FindByRefreshToken(tc.refreshToken)

			// Assert
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, session)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, session)
				assert.Equal(t, tc.expectedUID, session.UserID)
				assert.Equal(t, tc.refreshToken, session.RefreshToken)
			}
		})
	}
}

func (suite *SessionRepositoryTestSuite) TestFindActiveByUserID() {
	testCases := []struct {
		name          string
		userID        uint
		mockSetup     func()
		expectError   bool
		errorMsg      string
		expectedCount int
	}{
		{
			name:   "Success with multiple sessions",
			userID: 789,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "access_token", "refresh_token",
					"access_token_expires_at", "refresh_token_expires_at",
					"created_at", "updated_at", "is_deleted", "deleted_at",
				}).
					AddRow(1, 789, "token1", "refresh1", time.Now().Add(time.Hour), time.Now().Add(24*time.Hour), time.Now(), time.Now(), false, nil).
					AddRow(2, 789, "token2", "refresh2", time.Now().Add(time.Hour), time.Now().Add(24*time.Hour), time.Now(), time.Now(), false, nil)

				suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_sessions`")).
					WithArgs(uint(789), sqlmock.AnyArg()).
					WillReturnRows(rows)
			},
			expectError:   false,
			expectedCount: 2,
		},
		{
			name:   "No active sessions",
			userID: 999,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "access_token", "refresh_token",
					"access_token_expires_at", "refresh_token_expires_at",
					"created_at", "updated_at", "is_deleted", "deleted_at",
				})

				suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_sessions`")).
					WithArgs(uint(999), sqlmock.AnyArg()).
					WillReturnRows(rows)
			},
			expectError:   false,
			expectedCount: 0,
		},
		{
			name:   "Database error",
			userID: 888,
			mockSetup: func() {
				suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_sessions`")).
					WithArgs(uint(888), false, sqlmock.AnyArg()).
					WillReturnError(assert.AnError)
			},
			expectError: true,
			errorMsg:    "failed to find active sessions",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Setup mock expectations
			tc.mockSetup()

			// Execute
			sessions, err := suite.repo.FindActiveByUserID(tc.userID)

			// Assert
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, sessions)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, sessions)
				assert.Len(t, sessions, tc.expectedCount)

				for _, session := range sessions {
					assert.Equal(t, tc.userID, session.UserID)
					assert.False(t, session.IsDeleted)
				}
			}
		})
	}
}

func (suite *SessionRepositoryTestSuite) TestUpdate() {
	testCases := []struct {
		name        string
		session     *model.UserSession
		mockSetup   func()
		expectError bool
		errorMsg    string
	}{
		{
			name: "Success",
			session: &model.UserSession{
				ID:                    1,
				UserID:                123,
				AccessToken:           "updated-access-token",
				RefreshToken:          "updated-refresh-token",
				AccessTokenExpiresAt:  time.Now().Add(2 * time.Hour),
				RefreshTokenExpiresAt: time.Now().Add(48 * time.Hour),
				IsDeleted:             false,
			},
			mockSetup: func() {
				suite.mock.ExpectBegin()
				suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE `user_sessions`")).
					WithArgs(
						sqlmock.AnyArg(), // user_id
						sqlmock.AnyArg(), // access_token
						sqlmock.AnyArg(), // refresh_token
						sqlmock.AnyArg(), // access_token_expires_at
						sqlmock.AnyArg(), // refresh_token_expires_at
						sqlmock.AnyArg(), // created_at
						sqlmock.AnyArg(), // updated_at
						sqlmock.AnyArg(), // is_deleted
						sqlmock.AnyArg(), // deleted_at
						sqlmock.AnyArg(), // id
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				suite.mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "Database error",
			session: &model.UserSession{
				ID:     2,
				UserID: 456,
			},
			mockSetup: func() {
				suite.mock.ExpectBegin()
				suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE `user_sessions`")).
					WillReturnError(assert.AnError)
				suite.mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "failed to update session",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Setup mock expectations
			tc.mockSetup()

			// Execute
			err := suite.repo.Update(tc.session)

			// Assert
			if tc.expectError {
				assert.Error(t, err)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func (suite *SessionRepositoryTestSuite) TestInvalidateByUserID() {
	testCases := []struct {
		name        string
		userID      uint
		mockSetup   func()
		expectError bool
		errorMsg    string
	}{
		{
			name:   "Success",
			userID: 123,
			mockSetup: func() {
				suite.mock.ExpectBegin()
				suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE `user_sessions`")).
					WithArgs(true, sqlmock.AnyArg(), sqlmock.AnyArg(), uint(123), false).
					WillReturnResult(sqlmock.NewResult(0, 2)) // 2 rows affected
				suite.mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name:   "Database error",
			userID: 456,
			mockSetup: func() {
				suite.mock.ExpectBegin()
				suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE `user_sessions`")).
					WithArgs(true, sqlmock.AnyArg(), sqlmock.AnyArg(), uint(456), false).
					WillReturnError(assert.AnError)
				suite.mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "failed to invalidate sessions for user",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Setup mock expectations
			tc.mockSetup()

			// Execute
			err := suite.repo.InvalidateByUserID(tc.userID)

			// Assert
			if tc.expectError {
				assert.Error(t, err)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func (suite *SessionRepositoryTestSuite) TestInvalidateExpiredSessions() {
	testCases := []struct {
		name        string
		mockSetup   func()
		expectError bool
		errorMsg    string
	}{
		{
			name: "Success",
			mockSetup: func() {
				suite.mock.ExpectBegin()
				suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE `user_sessions`")).
					WithArgs(true, sqlmock.AnyArg(), sqlmock.AnyArg(), false, sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(0, 3)) // 3 expired sessions
				suite.mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "Database error",
			mockSetup: func() {
				suite.mock.ExpectBegin()
				suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE `user_sessions`")).
					WithArgs(true, sqlmock.AnyArg(), sqlmock.AnyArg(), false, sqlmock.AnyArg()).
					WillReturnError(assert.AnError)
				suite.mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "failed to invalidate expired sessions",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Setup mock expectations
			tc.mockSetup()

			// Execute
			err := suite.repo.InvalidateExpiredSessions()

			// Assert
			if tc.expectError {
				assert.Error(t, err)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func (suite *SessionRepositoryTestSuite) TestDelete() {
	testCases := []struct {
		name        string
		sessionID   uint
		mockSetup   func()
		expectError bool
		errorMsg    string
	}{
		{
			name:      "Success",
			sessionID: 1,
			mockSetup: func() {
				suite.mock.ExpectBegin()
				suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE `user_sessions` SET `deleted_at`")).
					WithArgs(sqlmock.AnyArg(), uint(1)).
					WillReturnResult(sqlmock.NewResult(0, 1))
				suite.mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name:      "Database error",
			sessionID: 2,
			mockSetup: func() {
				suite.mock.ExpectBegin()
				suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE `user_sessions` SET `deleted_at`")).
					WithArgs(sqlmock.AnyArg(), uint(2)).
					WillReturnError(assert.AnError)
				suite.mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "failed to delete session",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Setup mock expectations
			tc.mockSetup()

			// Execute
			err := suite.repo.Delete(tc.sessionID)

			// Assert
			if tc.expectError {
				assert.Error(t, err)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSessionRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(SessionRepositoryTestSuite))
}