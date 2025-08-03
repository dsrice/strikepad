package repository_test

import (
	"strikepad-backend/internal/repository"
	"testing"
	"time"

	"strikepad-backend/internal/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const testEmail = "test@example.com"

type UserRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	mock sqlmock.Sqlmock
	repo repository.UserRepository
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	assert.NoError(suite.T(), err)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(suite.T(), err)

	suite.db = gormDB
	suite.mock = mock
	suite.repo = repository.NewUserRepository(gormDB)
}

func (suite *UserRepositoryTestSuite) TearDownTest() {
	err := suite.mock.ExpectationsWereMet()
	assert.NoError(suite.T(), err)
}

func (suite *UserRepositoryTestSuite) TestCreate() {
	// Table-driven test for user creation
	tests := []struct {
		user         *model.User
		mockSetup    func()
		validateUser func(*model.User)
		name         string
		description  string
		expectError  bool
	}{
		{
			name: "successful email user creation",
			user: &model.User{
				ProviderType: "email",
				DisplayName:  "Test User",
				Email:        func() *string { s := testEmail; return &s }(),
			},
			mockSetup: func() {
				suite.mock.ExpectBegin()
				suite.mock.ExpectExec("INSERT INTO `users`").
					WithArgs(nil, nil, "test@example.com", nil, "email", "Test User", false, false).
					WillReturnResult(sqlmock.NewResult(1, 1))
				suite.mock.ExpectCommit()
			},
			expectError: false,
			validateUser: func(user *model.User) {
				assert.Equal(suite.T(), "email", user.ProviderType)
				assert.Equal(suite.T(), "Test User", user.DisplayName)
				assert.Equal(suite.T(), testEmail, *user.Email)
			},
			description: "should create email user successfully",
		},
		{
			name: "oauth user creation",
			user: &model.User{
				ProviderType:   "oauth",
				ProviderUserID: func() *string { s := "oauth123"; return &s }(),
				DisplayName:    "OAuth User",
				Email:          func() *string { s := "oauth@example.com"; return &s }(),
			},
			mockSetup: func() {
				suite.mock.ExpectBegin()
				suite.mock.ExpectExec("INSERT INTO `users`").
					WithArgs(nil, "oauth123", "oauth@example.com", nil, "oauth", "OAuth User", false, false).
					WillReturnResult(sqlmock.NewResult(2, 1))
				suite.mock.ExpectCommit()
			},
			expectError: false,
			validateUser: func(user *model.User) {
				assert.Equal(suite.T(), "oauth", user.ProviderType)
				assert.Equal(suite.T(), "OAuth User", user.DisplayName)
				assert.Equal(suite.T(), "oauth123", *user.ProviderUserID)
			},
			description: "should create OAuth user successfully",
		},
		{
			name: "user with password hash",
			user: &model.User{
				ProviderType: "email",
				DisplayName:  "Password User",
				Email:        func() *string { s := "password@example.com"; return &s }(),
				PasswordHash: func() *string { s := "hashedpassword"; return &s }(),
			},
			mockSetup: func() {
				suite.mock.ExpectBegin()
				suite.mock.ExpectExec("INSERT INTO `users`").
					WithArgs(nil, nil, "password@example.com", "hashedpassword", "email", "Password User", false, false).
					WillReturnResult(sqlmock.NewResult(3, 1))
				suite.mock.ExpectCommit()
			},
			expectError: false,
			validateUser: func(user *model.User) {
				assert.Equal(suite.T(), "email", user.ProviderType)
				assert.Equal(suite.T(), "Password User", user.DisplayName)
				assert.Equal(suite.T(), "hashedpassword", *user.PasswordHash)
			},
			description: "should create user with password hash successfully",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.mockSetup()

			createdUser, err := suite.repo.Create(tt.user)

			if tt.expectError {
				assert.Error(suite.T(), err, tt.description)
			} else {
				assert.NoError(suite.T(), err, tt.description)
				assert.NotNil(suite.T(), createdUser, "Created user should not be nil")
				if tt.validateUser != nil {
					tt.validateUser(createdUser)
				}
			}
		})
	}
}

func (suite *UserRepositoryTestSuite) TestGetByID() {
	// Table-driven test for getting user by ID
	tests := []struct {
		mockSetup    func()
		validateUser func(*model.User)
		name         string
		description  string
		userID       uint
		expectError  bool
	}{
		{
			name:   "successful get by ID",
			userID: 1,
			mockSetup: func() {
				email := testEmail
				now := time.Now()
				suite.mock.ExpectQuery("SELECT \\* FROM `users` WHERE `users`.`id` = \\? ORDER BY `users`.`id` LIMIT \\?").
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "provider_type", "provider_user_id", "email", "display_name", "password_hash", "email_verified", "created_at", "updated_at", "is_deleted", "deleted_at"}).
						AddRow(1, "email", nil, email, "Test User", nil, false, now, now, false, nil))
			},
			expectError: false,
			validateUser: func(user *model.User) {
				assert.Equal(suite.T(), uint(1), user.ID)
				assert.Equal(suite.T(), "Test User", user.DisplayName)
				assert.Equal(suite.T(), testEmail, *user.Email)
				assert.Equal(suite.T(), "email", user.ProviderType)
			},
			description: "should get user by ID successfully",
		},
		{
			name:   "get oauth user by ID",
			userID: 2,
			mockSetup: func() {
				email := "oauth@example.com"
				providerUserID := "oauth123"
				now := time.Now()
				suite.mock.ExpectQuery("SELECT \\* FROM `users` WHERE `users`.`id` = \\? ORDER BY `users`.`id` LIMIT \\?").
					WithArgs(2, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "provider_type", "provider_user_id", "email", "display_name", "password_hash", "email_verified", "created_at", "updated_at", "is_deleted", "deleted_at"}).
						AddRow(2, "oauth", providerUserID, email, "OAuth User", nil, true, now, now, false, nil))
			},
			expectError: false,
			validateUser: func(user *model.User) {
				assert.Equal(suite.T(), uint(2), user.ID)
				assert.Equal(suite.T(), "OAuth User", user.DisplayName)
				assert.Equal(suite.T(), "oauth@example.com", *user.Email)
				assert.Equal(suite.T(), "oauth", user.ProviderType)
				assert.Equal(suite.T(), "oauth123", *user.ProviderUserID)
				assert.Equal(suite.T(), true, user.EmailVerified)
			},
			description: "should get OAuth user by ID successfully",
		},
		{
			name:   "get verified user by ID",
			userID: 3,
			mockSetup: func() {
				email := "verified@example.com"
				now := time.Now()
				suite.mock.ExpectQuery("SELECT \\* FROM `users` WHERE `users`.`id` = \\? ORDER BY `users`.`id` LIMIT \\?").
					WithArgs(3, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "provider_type", "provider_user_id", "email", "display_name", "password_hash", "email_verified", "created_at", "updated_at", "is_deleted", "deleted_at"}).
						AddRow(3, "email", nil, email, "Verified User", "passwordhash", true, now, now, false, nil))
			},
			expectError: false,
			validateUser: func(user *model.User) {
				assert.Equal(suite.T(), uint(3), user.ID)
				assert.Equal(suite.T(), "Verified User", user.DisplayName)
				assert.Equal(suite.T(), true, user.EmailVerified)
				assert.Equal(suite.T(), "passwordhash", *user.PasswordHash)
			},
			description: "should get verified user with password hash successfully",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.mockSetup()

			found, err := suite.repo.GetByID(tt.userID)

			if tt.expectError {
				assert.Error(suite.T(), err, tt.description)
			} else {
				assert.NoError(suite.T(), err, tt.description)
				assert.NotNil(suite.T(), found, "Found user should not be nil")
				if tt.validateUser != nil {
					tt.validateUser(found)
				}
			}
		})
	}
}

func (suite *UserRepositoryTestSuite) TestGetByEmail() {
	// Table-driven test for getting user by email
	tests := []struct {
		mockSetup    func()
		validateUser func(*model.User)
		name         string
		email        string
		description  string
		expectError  bool
	}{
		{
			name:  "successful get by email",
			email: testEmail,
			mockSetup: func() {
				now := time.Now()
				suite.mock.ExpectQuery("SELECT \\* FROM `users` WHERE email = \\? ORDER BY `users`.`id` LIMIT \\?").
					WithArgs(testEmail, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "provider_type", "provider_user_id", "email", "display_name", "password_hash", "email_verified", "created_at", "updated_at", "is_deleted", "deleted_at"}).
						AddRow(1, "email", nil, testEmail, "Test User", nil, false, now, now, false, nil))
			},
			expectError: false,
			validateUser: func(user *model.User) {
				assert.Equal(suite.T(), uint(1), user.ID)
				assert.Equal(suite.T(), "Test User", user.DisplayName)
				assert.Equal(suite.T(), testEmail, *user.Email)
				assert.Equal(suite.T(), "email", user.ProviderType)
			},
			description: "should get user by email successfully",
		},
		{
			name:  "get oauth user by email",
			email: "oauth@example.com",
			mockSetup: func() {
				now := time.Now()
				suite.mock.ExpectQuery("SELECT \\* FROM `users` WHERE email = \\? ORDER BY `users`.`id` LIMIT \\?").
					WithArgs("oauth@example.com", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "provider_type", "provider_user_id", "email", "display_name", "password_hash", "email_verified", "created_at", "updated_at", "is_deleted", "deleted_at"}).
						AddRow(2, "oauth", "oauth123", "oauth@example.com", "OAuth User", nil, true, now, now, false, nil))
			},
			expectError: false,
			validateUser: func(user *model.User) {
				assert.Equal(suite.T(), uint(2), user.ID)
				assert.Equal(suite.T(), "OAuth User", user.DisplayName)
				assert.Equal(suite.T(), "oauth@example.com", *user.Email)
				assert.Equal(suite.T(), "oauth", user.ProviderType)
				assert.Equal(suite.T(), true, user.EmailVerified)
			},
			description: "should get OAuth user by email successfully",
		},
		{
			name:  "get user with different case email",
			email: "Mixed@Example.Com",
			mockSetup: func() {
				now := time.Now()
				suite.mock.ExpectQuery("SELECT \\* FROM `users` WHERE email = \\? ORDER BY `users`.`id` LIMIT \\?").
					WithArgs("Mixed@Example.Com", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "provider_type", "provider_user_id", "email", "display_name", "password_hash", "email_verified", "created_at", "updated_at", "is_deleted", "deleted_at"}).
						AddRow(3, "email", nil, "Mixed@Example.Com", "Mixed Case User", "hash123", false, now, now, false, nil))
			},
			expectError: false,
			validateUser: func(user *model.User) {
				assert.Equal(suite.T(), uint(3), user.ID)
				assert.Equal(suite.T(), "Mixed Case User", user.DisplayName)
				assert.Equal(suite.T(), "Mixed@Example.Com", *user.Email)
				assert.Equal(suite.T(), "hash123", *user.PasswordHash)
			},
			description: "should handle email with mixed case",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.mockSetup()

			found, err := suite.repo.GetByEmail(tt.email)

			if tt.expectError {
				assert.Error(suite.T(), err, tt.description)
			} else {
				assert.NoError(suite.T(), err, tt.description)
				assert.NotNil(suite.T(), found, "Found user should not be nil")
				if tt.validateUser != nil {
					tt.validateUser(found)
				}
			}
		})
	}
}

func (suite *UserRepositoryTestSuite) TestList() {
	// Table-driven test for listing users
	tests := []struct {
		mockSetup      func()
		validateResult func([]model.User)
		name           string
		description    string
		expectedCount  int
		expectError    bool
	}{
		{
			name: "list multiple users",
			mockSetup: func() {
				email1 := "user1@example.com"
				email2 := "user2@example.com"
				now := time.Now()
				suite.mock.ExpectQuery("SELECT \\* FROM `users`").
					WillReturnRows(sqlmock.NewRows([]string{"id", "provider_type", "provider_user_id", "email", "display_name", "password_hash", "email_verified", "created_at", "updated_at", "is_deleted", "deleted_at"}).
						AddRow(1, "email", nil, email1, "User 1", nil, false, now, now, false, nil).
						AddRow(2, "email", nil, email2, "User 2", nil, false, now, now, false, nil))
			},
			expectError:   false,
			expectedCount: 2,
			validateResult: func(users []model.User) {
				assert.Equal(suite.T(), "User 1", users[0].DisplayName)
				assert.Equal(suite.T(), "User 2", users[1].DisplayName)
				assert.Equal(suite.T(), "user1@example.com", *users[0].Email)
				assert.Equal(suite.T(), "user2@example.com", *users[1].Email)
			},
			description: "should list multiple users successfully",
		},
		{
			name: "list mixed provider types",
			mockSetup: func() {
				now := time.Now()
				suite.mock.ExpectQuery("SELECT \\* FROM `users`").
					WillReturnRows(sqlmock.NewRows([]string{"id", "provider_type", "provider_user_id", "email", "display_name", "password_hash", "email_verified", "created_at", "updated_at", "is_deleted", "deleted_at"}).
						AddRow(1, "email", nil, "email@example.com", "Email User", "hash123", true, now, now, false, nil).
						AddRow(2, "oauth", "oauth456", "oauth@example.com", "OAuth User", nil, true, now, now, false, nil).
						AddRow(3, "email", nil, "another@example.com", "Another User", "hash456", false, now, now, false, nil))
			},
			expectError:   false,
			expectedCount: 3,
			validateResult: func(users []model.User) {
				assert.Equal(suite.T(), "email", users[0].ProviderType)
				assert.Equal(suite.T(), "oauth", users[1].ProviderType)
				assert.Equal(suite.T(), "email", users[2].ProviderType)
				assert.Equal(suite.T(), "oauth456", *users[1].ProviderUserID)
				assert.True(suite.T(), users[0].EmailVerified)
				assert.False(suite.T(), users[2].EmailVerified)
			},
			description: "should list users with different provider types",
		},
		{
			name: "empty user list",
			mockSetup: func() {
				suite.mock.ExpectQuery("SELECT \\* FROM `users`").
					WillReturnRows(sqlmock.NewRows([]string{"id", "provider_type", "provider_user_id", "email", "display_name", "password_hash", "email_verified", "created_at", "updated_at", "is_deleted", "deleted_at"}))
			},
			expectError:   false,
			expectedCount: 0,
			validateResult: func(users []model.User) {
				assert.Empty(suite.T(), users)
			},
			description: "should handle empty user list",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.mockSetup()

			result, err := suite.repo.List()

			if tt.expectError {
				assert.Error(suite.T(), err, tt.description)
			} else {
				assert.NoError(suite.T(), err, tt.description)
				assert.Len(suite.T(), result, tt.expectedCount, tt.description)
				if tt.validateResult != nil {
					tt.validateResult(result)
				}
			}
		})
	}
}

func (suite *UserRepositoryTestSuite) TestDelete() {
	// Table-driven test for user deletion
	tests := []struct {
		mockSetup   func()
		name        string
		description string
		userID      uint
		expectError bool
	}{
		{
			name:   "successful deletion",
			userID: 1,
			mockSetup: func() {
				suite.mock.ExpectBegin()
				suite.mock.ExpectExec("DELETE FROM `users` WHERE `users`.`id` = \\?").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				suite.mock.ExpectCommit()
			},
			expectError: false,
			description: "should delete user successfully",
		},
		{
			name:   "delete different user ID",
			userID: 99,
			mockSetup: func() {
				suite.mock.ExpectBegin()
				suite.mock.ExpectExec("DELETE FROM `users` WHERE `users`.`id` = \\?").
					WithArgs(99).
					WillReturnResult(sqlmock.NewResult(1, 1))
				suite.mock.ExpectCommit()
			},
			expectError: false,
			description: "should delete user with different ID successfully",
		},
		{
			name:   "delete zero ID",
			userID: 0,
			mockSetup: func() {
				suite.mock.ExpectBegin()
				suite.mock.ExpectExec("DELETE FROM `users` WHERE `users`.`id` = \\?").
					WithArgs(0).
					WillReturnResult(sqlmock.NewResult(1, 1))
				suite.mock.ExpectCommit()
			},
			expectError: false,
			description: "should handle zero ID deletion",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.mockSetup()

			err := suite.repo.Delete(tt.userID)

			if tt.expectError {
				assert.Error(suite.T(), err, tt.description)
			} else {
				assert.NoError(suite.T(), err, tt.description)
			}
		})
	}
}

func (suite *UserRepositoryTestSuite) TestFindByEmail() {
	// Table-driven test for finding user by email (non-deleted)
	tests := []struct {
		mockSetup    func()
		validateUser func(*model.User)
		name         string
		email        string
		description  string
		expectError  bool
	}{
		{
			name:  "find active user by email",
			email: testEmail,
			mockSetup: func() {
				now := time.Now()
				suite.mock.ExpectQuery("SELECT \\* FROM `users` WHERE email = \\? AND is_deleted = \\? ORDER BY `users`.`id` LIMIT \\?").
					WithArgs(testEmail, false, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "provider_type", "provider_user_id", "email", "display_name", "password_hash", "email_verified", "created_at", "updated_at", "is_deleted", "deleted_at"}).
						AddRow(1, "email", nil, testEmail, "Test User", nil, false, now, now, false, nil))
			},
			expectError: false,
			validateUser: func(user *model.User) {
				assert.Equal(suite.T(), uint(1), user.ID)
				assert.Equal(suite.T(), "Test User", user.DisplayName)
				assert.Equal(suite.T(), testEmail, *user.Email)
				assert.Equal(suite.T(), false, user.IsDeleted)
				assert.Equal(suite.T(), "email", user.ProviderType)
			},
			description: "should find active user by email successfully",
		},
		{
			name:  "find verified oauth user by email",
			email: "oauth@example.com",
			mockSetup: func() {
				now := time.Now()
				suite.mock.ExpectQuery("SELECT \\* FROM `users` WHERE email = \\? AND is_deleted = \\? ORDER BY `users`.`id` LIMIT \\?").
					WithArgs("oauth@example.com", false, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "provider_type", "provider_user_id", "email", "display_name", "password_hash", "email_verified", "created_at", "updated_at", "is_deleted", "deleted_at"}).
						AddRow(2, "oauth", "oauth123", "oauth@example.com", "OAuth User", nil, true, now, now, false, nil))
			},
			expectError: false,
			validateUser: func(user *model.User) {
				assert.Equal(suite.T(), uint(2), user.ID)
				assert.Equal(suite.T(), "OAuth User", user.DisplayName)
				assert.Equal(suite.T(), "oauth@example.com", *user.Email)
				assert.Equal(suite.T(), false, user.IsDeleted)
				assert.Equal(suite.T(), "oauth", user.ProviderType)
				assert.Equal(suite.T(), true, user.EmailVerified)
				assert.Equal(suite.T(), "oauth123", *user.ProviderUserID)
			},
			description: "should find verified OAuth user by email successfully",
		},
		{
			name:  "find user with password hash",
			email: "secure@example.com",
			mockSetup: func() {
				now := time.Now()
				suite.mock.ExpectQuery("SELECT \\* FROM `users` WHERE email = \\? AND is_deleted = \\? ORDER BY `users`.`id` LIMIT \\?").
					WithArgs("secure@example.com", false, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "provider_type", "provider_user_id", "email", "display_name", "password_hash", "email_verified", "created_at", "updated_at", "is_deleted", "deleted_at"}).
						AddRow(3, "email", nil, "secure@example.com", "Secure User", "$2a$10$hashedpassword", true, now, now, false, nil))
			},
			expectError: false,
			validateUser: func(user *model.User) {
				assert.Equal(suite.T(), uint(3), user.ID)
				assert.Equal(suite.T(), "Secure User", user.DisplayName)
				assert.Equal(suite.T(), "secure@example.com", *user.Email)
				assert.Equal(suite.T(), false, user.IsDeleted)
				assert.Equal(suite.T(), true, user.EmailVerified)
				assert.Equal(suite.T(), "$2a$10$hashedpassword", *user.PasswordHash)
			},
			description: "should find user with password hash successfully",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.mockSetup()

			found, err := suite.repo.FindByEmail(tt.email)

			if tt.expectError {
				assert.Error(suite.T(), err, tt.description)
			} else {
				assert.NoError(suite.T(), err, tt.description)
				assert.NotNil(suite.T(), found, "Found user should not be nil")
				if tt.validateUser != nil {
					tt.validateUser(found)
				}
			}
		})
	}
}

func (suite *UserRepositoryTestSuite) TestUpdate() {
	// Table-driven test for user updates
	tests := []struct {
		user        *model.User
		mockSetup   func()
		name        string
		description string
		expectError bool
	}{
		{
			name: "update display name",
			user: &model.User{
				ID:           1,
				ProviderType: "email",
				DisplayName:  "Updated User",
				Email:        func() *string { s := testEmail; return &s }(),
			},
			mockSetup: func() {
				suite.mock.ExpectBegin()
				suite.mock.ExpectExec("UPDATE `users` SET").
					WillReturnResult(sqlmock.NewResult(1, 1))
				suite.mock.ExpectCommit()
			},
			expectError: false,
			description: "should update user display name successfully",
		},
		{
			name: "update email verification status",
			user: &model.User{
				ID:            2,
				ProviderType:  "email",
				DisplayName:   "Email User",
				Email:         func() *string { s := "verify@example.com"; return &s }(),
				EmailVerified: true,
			},
			mockSetup: func() {
				suite.mock.ExpectBegin()
				suite.mock.ExpectExec("UPDATE `users` SET").
					WillReturnResult(sqlmock.NewResult(1, 1))
				suite.mock.ExpectCommit()
			},
			expectError: false,
			description: "should update email verification status successfully",
		},
		{
			name: "update password hash",
			user: &model.User{
				ID:           3,
				ProviderType: "email",
				DisplayName:  "Password User",
				Email:        func() *string { s := "password@example.com"; return &s }(),
				PasswordHash: func() *string { s := "$2a$10$newhash"; return &s }(),
			},
			mockSetup: func() {
				suite.mock.ExpectBegin()
				suite.mock.ExpectExec("UPDATE `users` SET").
					WillReturnResult(sqlmock.NewResult(1, 1))
				suite.mock.ExpectCommit()
			},
			expectError: false,
			description: "should update password hash successfully",
		},
		{
			name: "update oauth user",
			user: &model.User{
				ID:             4,
				ProviderType:   "oauth",
				ProviderUserID: func() *string { s := "oauth456"; return &s }(),
				DisplayName:    "Updated OAuth User",
				Email:          func() *string { s := "oauth@example.com"; return &s }(),
				EmailVerified:  true,
			},
			mockSetup: func() {
				suite.mock.ExpectBegin()
				suite.mock.ExpectExec("UPDATE `users` SET").
					WillReturnResult(sqlmock.NewResult(1, 1))
				suite.mock.ExpectCommit()
			},
			expectError: false,
			description: "should update OAuth user successfully",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tt.mockSetup()

			err := suite.repo.Update(tt.user)

			if tt.expectError {
				assert.Error(suite.T(), err, tt.description)
			} else {
				assert.NoError(suite.T(), err, tt.description)
			}
		})
	}
}

func (suite *UserRepositoryTestSuite) TestNewUserRepository() {
	// Test that NewUserRepository creates a repository with the provided DB
	repo := repository.NewUserRepository(suite.db)
	assert.NotNil(suite.T(), repo)

	// Since we're in the repository_test package, we can't access the unexported userRepository type
	// We can only verify that the repository is not nil, which is sufficient for this test
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}