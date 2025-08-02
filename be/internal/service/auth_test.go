package service

import (
	"testing"

	"strikepad-backend/internal/auth"
	"strikepad-backend/internal/dto"
	"strikepad-backend/internal/model"
	"strikepad-backend/internal/repository/mocks"

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
	authService  AuthServiceInterface
	mockUserRepo *mocks.MockUserRepository
}

func (suite *AuthServiceTestSuite) SetupTest() {
	suite.mockUserRepo = new(mocks.MockUserRepository)
	suite.authService = NewAuthService(suite.mockUserRepo)
}

func (suite *AuthServiceTestSuite) TearDownTest() {
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *AuthServiceTestSuite) TestSignupSuccess() {
	email := testServiceEmailConst
	request := &dto.SignupRequest{
		Email:       email,
		Password:    testServicePasswordConst,
		DisplayName: "Test User",
	}

	// Expected user to be created
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

	// Execute
	result, err := suite.authService.Signup(request)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), expectedUser.ID, result.ID)
	assert.Equal(suite.T(), *expectedUser.Email, result.Email)
	assert.Equal(suite.T(), expectedUser.DisplayName, result.DisplayName)
	assert.Equal(suite.T(), expectedUser.EmailVerified, result.EmailVerified)
}

func (suite *AuthServiceTestSuite) TestSignupUserAlreadyExists() {
	email := "existing@example.com"
	request := &dto.SignupRequest{
		Email:       email,
		Password:    testServicePasswordConst,
		DisplayName: "Test User",
	}

	existingUser := &model.User{
		ID:          1,
		Email:       &email,
		DisplayName: "Existing User",
	}

	// Mock: FindByEmail returns existing user
	suite.mockUserRepo.On("FindByEmail", "existing@example.com").Return(existingUser, nil)

	// Execute
	result, err := suite.authService.Signup(request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), auth.ErrUserAlreadyExists, err)
	assert.Nil(suite.T(), result)
}

func (suite *AuthServiceTestSuite) TestSignupInvalidEmail() {
	request := &dto.SignupRequest{
		Email:       "invalid-email",
		Password:    testServicePasswordConst,
		DisplayName: "Test User",
	}

	// Execute
	result, err := suite.authService.Signup(request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), auth.ErrInvalidEmail, err)
	assert.Nil(suite.T(), result)
}

func (suite *AuthServiceTestSuite) TestSignupPasswordTooShort() {
	request := &dto.SignupRequest{
		Email:       testServiceEmailConst,
		Password:    "short",
		DisplayName: "Test User",
	}

	// Execute
	result, err := suite.authService.Signup(request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), auth.ErrPasswordTooShort, err)
	assert.Nil(suite.T(), result)
}

func (suite *AuthServiceTestSuite) TestSignupPasswordTooLong() {
	longPassword := testServicePasswordConst + string(make([]byte, 120)) // 132 chars total
	request := &dto.SignupRequest{
		Email:       testServiceEmailConst,
		Password:    longPassword,
		DisplayName: "Test User",
	}

	// Execute
	result, err := suite.authService.Signup(request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), auth.ErrPasswordTooLong, err)
	assert.Nil(suite.T(), result)
}

func (suite *AuthServiceTestSuite) TestSignupRepositoryCreateError() {
	email := testServiceEmailConst
	request := &dto.SignupRequest{
		Email:       email,
		Password:    testServicePasswordConst,
		DisplayName: "Test User",
	}

	// Mock: FindByEmail returns not found (user doesn't exist)
	suite.mockUserRepo.On("FindByEmail", testServiceEmailConst).Return(nil, gorm.ErrRecordNotFound)

	// Mock: Create returns an error
	suite.mockUserRepo.On("Create", mock.AnythingOfType("*model.User")).Return(nil, assert.AnError)

	// Execute
	result, err := suite.authService.Signup(request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "internal server error", err.Error()) // Service converts to generic error
	assert.Nil(suite.T(), result)
}

func (suite *AuthServiceTestSuite) TestLoginSuccess() {
	email := testServiceEmailConst
	password := testServicePasswordConst
	hashedPassword, _ := auth.HashPassword(password)

	request := &dto.LoginRequest{
		Email:    email,
		Password: password,
	}

	existingUser := &model.User{
		ID:           1,
		ProviderType: "email",
		Email:        &email,
		DisplayName:  "Test User",
		PasswordHash: &hashedPassword,
	}

	// Mock: FindByEmail returns the user
	suite.mockUserRepo.On("FindByEmail", testServiceEmailConst).Return(existingUser, nil)

	// Execute
	result, err := suite.authService.Login(request)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), existingUser.ID, result.ID)
	assert.Equal(suite.T(), *existingUser.Email, result.Email)
	assert.Equal(suite.T(), existingUser.DisplayName, result.DisplayName)
}

func (suite *AuthServiceTestSuite) TestLoginUserNotFound() {
	request := &dto.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: testServicePasswordConst,
	}

	// Mock: FindByEmail returns not found error
	suite.mockUserRepo.On("FindByEmail", "nonexistent@example.com").Return(nil, gorm.ErrRecordNotFound)

	// Execute
	result, err := suite.authService.Login(request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), auth.ErrInvalidCredentials, err)
	assert.Nil(suite.T(), result)
}

func (suite *AuthServiceTestSuite) TestLoginInvalidPassword() {
	email := testServiceEmailConst
	correctPassword := testServicePasswordConst
	wrongPassword := "WrongPassword456!"
	hashedPassword, _ := auth.HashPassword(correctPassword)

	request := &dto.LoginRequest{
		Email:    email,
		Password: wrongPassword,
	}

	existingUser := &model.User{
		ID:           1,
		ProviderType: "email",
		Email:        &email,
		DisplayName:  "Test User",
		PasswordHash: &hashedPassword,
	}

	// Mock: FindByEmail returns the user
	suite.mockUserRepo.On("FindByEmail", testServiceEmailConst).Return(existingUser, nil)

	// Execute
	result, err := suite.authService.Login(request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), auth.ErrInvalidCredentials, err)
	assert.Nil(suite.T(), result)
}

func (suite *AuthServiceTestSuite) TestLoginInvalidEmail() {
	request := &dto.LoginRequest{
		Email:    "invalid-email",
		Password: testServicePasswordConst,
	}

	// Execute
	result, err := suite.authService.Login(request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), auth.ErrInvalidCredentials, err) // Service converts to invalid credentials
	assert.Nil(suite.T(), result)
}

func (suite *AuthServiceTestSuite) TestLoginUserWithoutPassword() {
	email := testServiceEmailConst

	request := &dto.LoginRequest{
		Email:    email,
		Password: testServicePasswordConst,
	}

	// User exists but has no password (OAuth user, for example)
	existingUser := &model.User{
		ID:           1,
		ProviderType: "oauth",
		Email:        &email,
		DisplayName:  "Test User",
		PasswordHash: nil, // No password hash
	}

	// Mock: FindByEmail returns the user
	suite.mockUserRepo.On("FindByEmail", testServiceEmailConst).Return(existingUser, nil)

	// Execute
	result, err := suite.authService.Login(request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), auth.ErrInvalidCredentials, err)
	assert.Nil(suite.T(), result)
}

func (suite *AuthServiceTestSuite) TestLoginRepositoryError() {
	request := &dto.LoginRequest{
		Email:    testServiceEmailConst,
		Password: testServicePasswordConst,
	}

	// Mock: FindByEmail returns a repository error
	suite.mockUserRepo.On("FindByEmail", testServiceEmailConst).Return(nil, assert.AnError)

	// Execute
	result, err := suite.authService.Login(request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "internal server error", err.Error()) // Service converts to generic error
	assert.Nil(suite.T(), result)
}

func (suite *AuthServiceTestSuite) TestNewAuthService() {
	// Test that NewAuthService creates a valid service
	service := NewAuthService(suite.mockUserRepo)
	assert.NotNil(suite.T(), service)
}

func (suite *AuthServiceTestSuite) TestEmailNormalization() {
	email := "  Test.User@EXAMPLE.COM  "
	normalizedEmail := "test.user@example.com"
	password := testServicePasswordConst
	hashedPassword, _ := auth.HashPassword(password)

	request := &dto.LoginRequest{
		Email:    email, // Email with spaces and mixed case
		Password: password,
	}

	existingUser := &model.User{
		ID:           1,
		ProviderType: "email",
		Email:        &normalizedEmail, // Stored email is normalized
		DisplayName:  "Test User",
		PasswordHash: &hashedPassword,
	}

	// Mock: FindByEmail should be called with normalized email
	suite.mockUserRepo.On("FindByEmail", normalizedEmail).Return(existingUser, nil)

	// Execute
	result, err := suite.authService.Login(request)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), normalizedEmail, result.Email)
}

func TestAuthServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}
