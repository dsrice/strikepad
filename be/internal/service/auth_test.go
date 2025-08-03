package service_test

import (
	"testing"

	"strikepad-backend/internal/auth"
	"strikepad-backend/internal/dto"
	"strikepad-backend/internal/model"
	"strikepad-backend/internal/repository/mocks"
	"strikepad-backend/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

const (
	testServiceEmailConst    = "test@example.com"
	testServicePasswordConst = "Password123!"
)

type AuthServiceTestSuite struct {
	suite.Suite
	authService  service.AuthServiceInterface
	mockUserRepo *mocks.MockUserRepository
}

func (suite *AuthServiceTestSuite) SetupTest() {
	suite.mockUserRepo = new(mocks.MockUserRepository)
	suite.authService = service.NewAuthService(suite.mockUserRepo)
}

func (suite *AuthServiceTestSuite) TearDownTest() {
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *AuthServiceTestSuite) TestSignup() {
	testCases := []struct {
		name           string
		request        *dto.SignupRequest
		mockSetup      func()
		expectedError  error
		expectedErrMsg string // For cases where we check the error message instead of the error itself
		checkResult    bool
	}{
		{
			name: "Success",
			request: &dto.SignupRequest{
				Email:       testServiceEmailConst,
				Password:    testServicePasswordConst,
				DisplayName: "Test User",
			},
			mockSetup: func() {
				email := testServiceEmailConst
				expectedUser := &model.User{
					ID:            1,
					ProviderType:  "email",
					Email:         &email,
					DisplayName:   "Test User",
					EmailVerified: false,
				}

				// Mock: FindByEmail returns not found error (user doesn't exist)
				suite.mockUserRepo.On("FindByEmail", testServiceEmailConst).Return(nil, gorm.ErrRecordNotFound)

				// Mock: Create returns the new user
				suite.mockUserRepo.On("Create", mock.MatchedBy(func(user *model.User) bool {
					return user.ProviderType == "email" &&
						*user.Email == testServiceEmailConst &&
						user.DisplayName == "Test User" &&
						user.PasswordHash != nil &&
						!user.EmailVerified
				})).Return(expectedUser, nil)
			},
			expectedError:  nil,
			expectedErrMsg: "",
			checkResult:    true,
		},
		{
			name: "User already exists",
			request: &dto.SignupRequest{
				Email:       "existing@example.com",
				Password:    testServicePasswordConst,
				DisplayName: "Test User",
			},
			mockSetup: func() {
				email := "existing@example.com"
				existingUser := &model.User{
					ID:          1,
					Email:       &email,
					DisplayName: "Existing User",
				}

				// Mock: FindByEmail returns existing user
				suite.mockUserRepo.On("FindByEmail", "existing@example.com").Return(existingUser, nil)
			},
			expectedError:  auth.ErrUserAlreadyExists,
			expectedErrMsg: "",
			checkResult:    false,
		},
		{
			name: "Invalid email",
			request: &dto.SignupRequest{
				Email:       "invalid-email",
				Password:    testServicePasswordConst,
				DisplayName: "Test User",
			},
			mockSetup:      func() {},
			expectedError:  auth.ErrInvalidEmail,
			expectedErrMsg: "",
			checkResult:    false,
		},
		{
			name: "Password too short",
			request: &dto.SignupRequest{
				Email:       testServiceEmailConst,
				Password:    "short",
				DisplayName: "Test User",
			},
			mockSetup:      func() {},
			expectedError:  auth.ErrPasswordTooShort,
			expectedErrMsg: "",
			checkResult:    false,
		},
		{
			name: "Password too long",
			request: &dto.SignupRequest{
				Email:       testServiceEmailConst,
				Password:    testServicePasswordConst + string(make([]byte, 120)), // 132 chars total
				DisplayName: "Test User",
			},
			mockSetup:      func() {},
			expectedError:  auth.ErrPasswordTooLong,
			expectedErrMsg: "",
			checkResult:    false,
		},
		{
			name: "Database error when checking existing user",
			request: &dto.SignupRequest{
				Email:       "dberror-signup@example.com",
				Password:    testServicePasswordConst,
				DisplayName: "Test User",
			},
			mockSetup: func() {
				// Mock: FindByEmail returns a database error
				suite.mockUserRepo.On("FindByEmail", "dberror-signup@example.com").Return(nil, assert.AnError)
			},
			expectedError:  nil,
			expectedErrMsg: "internal server error",
			checkResult:    false,
		},
		{
			name: "Repository create error",
			request: &dto.SignupRequest{
				Email:       "create-error@example.com",
				Password:    testServicePasswordConst,
				DisplayName: "Test User",
			},
			mockSetup: func() {
				// Mock: FindByEmail returns not found (user doesn't exist)
				suite.mockUserRepo.On("FindByEmail", "create-error@example.com").Return(nil, gorm.ErrRecordNotFound)

				// Mock: Create returns an error
				suite.mockUserRepo.On("Create", mock.MatchedBy(func(user *model.User) bool {
					return user.ProviderType == "email" &&
						*user.Email == "create-error@example.com" &&
						user.DisplayName == "Test User" &&
						user.PasswordHash != nil &&
						!user.EmailVerified
				})).Return(nil, assert.AnError)
			},
			expectedError:  nil,
			expectedErrMsg: "internal server error", // Service converts to generic error
			checkResult:    false,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Setup mocks for this test case
			tc.mockSetup()

			// Execute
			result, err := suite.authService.Signup(tc.request)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, result)
			} else if tc.expectedErrMsg != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErrMsg, err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)

				if tc.checkResult {
					email := tc.request.Email
					expectedUser := &model.User{
						ID:            1,
						ProviderType:  "email",
						Email:         &email,
						DisplayName:   "Test User",
						EmailVerified: false,
					}
					assert.Equal(t, expectedUser.ID, result.ID)
					assert.Equal(t, *expectedUser.Email, result.Email)
					assert.Equal(t, expectedUser.DisplayName, result.DisplayName)
					assert.Equal(t, expectedUser.EmailVerified, result.EmailVerified)
				}
			}
		})
	}
}

func (suite *AuthServiceTestSuite) TestLogin() {
	// Pre-compute a hashed password for the correct password
	correctPassword := testServicePasswordConst
	hashedPassword, _ := auth.HashPassword(correctPassword)

	testCases := []struct {
		name           string
		request        *dto.LoginRequest
		mockSetup      func()
		expectedError  error
		expectedErrMsg string // For cases where we check the error message instead of the error itself
		checkResult    bool
	}{
		{
			name: "Success",
			request: &dto.LoginRequest{
				Email:    testServiceEmailConst,
				Password: correctPassword,
			},
			mockSetup: func() {
				email := testServiceEmailConst
				existingUser := &model.User{
					ID:           1,
					ProviderType: "email",
					Email:        &email,
					DisplayName:  "Test User",
					PasswordHash: &hashedPassword,
				}
				// Mock: FindByEmail returns the user
				suite.mockUserRepo.On("FindByEmail", testServiceEmailConst).Return(existingUser, nil)
			},
			expectedError:  nil,
			expectedErrMsg: "",
			checkResult:    true,
		},
		{
			name: "User not found",
			request: &dto.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: testServicePasswordConst,
			},
			mockSetup: func() {
				// Mock: FindByEmail returns not found error
				suite.mockUserRepo.On("FindByEmail", "nonexistent@example.com").Return(nil, gorm.ErrRecordNotFound)
			},
			expectedError:  auth.ErrInvalidCredentials,
			expectedErrMsg: "",
			checkResult:    false,
		},
		{
			name: "Invalid password",
			request: &dto.LoginRequest{
				Email:    testServiceEmailConst,
				Password: "WrongPassword456!",
			},
			mockSetup: func() {
				email := testServiceEmailConst
				existingUser := &model.User{
					ID:           1,
					ProviderType: "email",
					Email:        &email,
					DisplayName:  "Test User",
					PasswordHash: &hashedPassword,
				}
				// Mock: FindByEmail returns the user
				suite.mockUserRepo.On("FindByEmail", testServiceEmailConst).Return(existingUser, nil)
			},
			expectedError:  auth.ErrInvalidCredentials,
			expectedErrMsg: "",
			checkResult:    false,
		},
		{
			name: "Invalid email",
			request: &dto.LoginRequest{
				Email:    "invalid-email",
				Password: testServicePasswordConst,
			},
			mockSetup: func() {
				// No mock setup needed
			},
			expectedError:  auth.ErrInvalidCredentials,
			expectedErrMsg: "",
			checkResult:    false,
		},
		{
			name: "User is deleted",
			request: &dto.LoginRequest{
				Email:    "deleted@example.com",
				Password: correctPassword,
			},
			mockSetup: func() {
				email := "deleted@example.com"
				existingUser := &model.User{
					ID:           1,
					ProviderType: "email",
					Email:        &email,
					DisplayName:  "Test User",
					PasswordHash: &hashedPassword,
					IsDeleted:    true,
				}
				// Mock: FindByEmail returns a deleted user
				suite.mockUserRepo.On("FindByEmail", "deleted@example.com").Return(existingUser, nil)
			},
			expectedError:  auth.ErrInvalidCredentials,
			expectedErrMsg: "",
			checkResult:    false,
		},
		{
			name: "User without password hash",
			request: &dto.LoginRequest{
				Email:    "oauth@example.com",
				Password: correctPassword,
			},
			mockSetup: func() {
				email := "oauth@example.com"
				existingUser := &model.User{
					ID:           1,
					ProviderType: "oauth",
					Email:        &email,
					DisplayName:  "Test User",
					PasswordHash: nil, // No password hash
				}
				// Mock: FindByEmail returns a user without password hash
				suite.mockUserRepo.On("FindByEmail", "oauth@example.com").Return(existingUser, nil)
			},
			expectedError:  auth.ErrInvalidCredentials,
			expectedErrMsg: "",
			checkResult:    false,
		},
		{
			name: "Database error when finding user",
			request: &dto.LoginRequest{
				Email:    "dberror@example.com",
				Password: correctPassword,
			},
			mockSetup: func() {
				// Mock: FindByEmail returns a database error
				suite.mockUserRepo.On("FindByEmail", "dberror@example.com").Return(nil, assert.AnError)
			},
			expectedError:  nil,
			expectedErrMsg: "internal server error",
			checkResult:    false,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Setup mocks for this test case
			tc.mockSetup()

			// Execute
			result, err := suite.authService.Login(tc.request)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, result)
			} else if tc.expectedErrMsg != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErrMsg, err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)

				if tc.checkResult {
					// For successful login, check the result fields
					email := tc.request.Email
					existingUser := &model.User{
						ID:           1,
						ProviderType: "email",
						Email:        &email,
						DisplayName:  "Test User",
					}
					assert.Equal(t, existingUser.ID, result.ID)
					assert.Equal(t, *existingUser.Email, result.Email)
					assert.Equal(t, existingUser.DisplayName, result.DisplayName)
				}
			}
		})
	}
}

func (suite *AuthServiceTestSuite) TestNewAuthService() {
	// Test that NewAuthService creates a valid service
	svc := service.NewAuthService(suite.mockUserRepo)
	assert.NotNil(suite.T(), svc)
}

func (suite *AuthServiceTestSuite) TestEmailNormalization() {
	password := testServicePasswordConst
	hashedPassword, _ := auth.HashPassword(password)

	testCases := []struct {
		name            string
		inputEmail      string
		normalizedEmail string
	}{
		{
			name:            "Email with spaces and mixed case",
			inputEmail:      "  Test.User@EXAMPLE.COM  ",
			normalizedEmail: "test.user@example.com",
		},
		{
			name:            "Email with only mixed case",
			inputEmail:      "Another.USER@Example.COM",
			normalizedEmail: "another.user@example.com",
		},
		{
			name:            "Email already normalized",
			inputEmail:      "simple@example.com",
			normalizedEmail: "simple@example.com",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			request := &dto.LoginRequest{
				Email:    tc.inputEmail,
				Password: password,
			}

			existingUser := &model.User{
				ID:           1,
				ProviderType: "email",
				Email:        &tc.normalizedEmail, // Stored email is normalized
				DisplayName:  "Test User",
				PasswordHash: &hashedPassword,
			}

			// Mock: FindByEmail should be called with normalized email
			suite.mockUserRepo.On("FindByEmail", tc.normalizedEmail).Return(existingUser, nil)

			// Execute
			result, err := suite.authService.Login(request)

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tc.normalizedEmail, result.Email)
		})
	}
}

func TestAuthServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}
