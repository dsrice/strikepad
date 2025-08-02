package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"strikepad-backen
	"strikepad-backend/internal/auth"
	"strikepad-backend/internal/dto"
	"gorm.io/gorm"
)

// MockUserRepository implements the UserRepository interface for testing
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *model.User) (*model.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(id uint) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) List() ([]model.User, error) {
	args := m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

type AuthServiceTestSuite struct {
	authService  *AuthService
	mockUserRepo *MockUserRepository
	mockUserRepo   *MockUserRepository
}

func (suite *AuthServiceTestSuite) SetupTest() {
	suite.mockUserRepo = new(MockUserRepository)
	suite.authService = NewAuthService(suite.mockUserRepo)
}

func (suite *AuthServiceTestSuite) TearDownTest() {
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *AuthServiceTestSuite) TestSignupSuccess() {
	email := "test@example.com"
	request := &dto.SignupRequest{
		Email:       email,
		Password:    "Password123!",
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
	suite.mockUserRepo.On("FindByEmail", "test@example.com").Return(nil, gorm.ErrRecordNotFound)

	// Mock: Create returns the new user
	suite.mockUserRepo.On("Create", mock.MatchedBy(func(user *model.User) bool {
		return user.ProviderType == "email" &&
			*user.Email == "test@example.com" &&
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
		Password:    "Password123!",
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
		Password:    "Password123!",
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
		Email:       "test@example.com",
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
	longPassword := "Password123!" + string(make([]byte, 120)) // 132 chars total
	request := &dto.SignupRequest{
		Email:       "test@example.com",
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
	email := "test@example.com"
	request := &dto.SignupRequest{
		Email:       email,
		Password:    "Password123!",
		DisplayName: "Test User",
	}

	// Mock: FindByEmail returns not found (user doesn't exist)
	suite.mockUserRepo.On("FindByEmail", "test@example.com").Return(nil, gorm.ErrRecordNotFound)

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
	email := "test@example.com"
	password := "Password123!"
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
	suite.mockUserRepo.On("FindByEmail", "test@example.com").Return(existingUser, nil)

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
		Password: "Password123!",
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
	email := "test@example.com"
	correctPassword := "Password123!"
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
	suite.mockUserRepo.On("FindByEmail", "test@example.com").Return(existingUser, nil)

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
		Password: "Password123!",
	}

	// Execute
	result, err := suite.authService.Login(request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), auth.ErrInvalidCredentials, err) // Service converts to invalid credentials
	assert.Nil(suite.T(), result)
}

func (suite *AuthServiceTestSuite) TestLoginUserWithoutPassword() {

	
	request := &dto.LoginRequest{
		Email:    email,
		Password: "Password123!",
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
	suite.mockUserRepo.On("FindByEmail", "test@example.com").Return(existingUser, nil)

	// Execute
	result, err := suite.authService.Login(request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), auth.ErrInvalidCredentials, err)
	assert.Nil(suite.T(), result)
}

func (suite *AuthServiceTestSuite) TestLoginRepositoryError() {
	request := &dto.LoginRequest{
		Email:    "test@example.com",
		Password: "Password123!",
	}

	// Mock: FindByEmail returns a repository error
	suite.mockUserRepo.On("FindByEmail", "test@example.com").Return(nil, assert.AnError)

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
	assert.Equal(suite.T(), suite.mockUserRepo, service.userRepo)
}

func (suite *AuthServiceTestSuite) TestEmailNormalization() {
	email := "  Test.User@EXAMPLE.COM  "
	normalizedEmail := "test.user@example.com"
	password := "Password123!"
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
t
}