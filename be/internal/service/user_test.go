package service_test

import (
	"testing"

	"strikepad-backend/internal/auth"
	"strikepad-backend/internal/model"
	"strikepad-backend/internal/repository/mocks"
	"strikepad-backend/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

const (
	testUserServiceEmailConst = "test@example.com"
)

type UserServiceTestSuite struct {
	suite.Suite
	userService  service.UserServiceInterface
	mockUserRepo *mocks.MockUserRepository
}

func (suite *UserServiceTestSuite) SetupTest() {
	suite.mockUserRepo = new(mocks.MockUserRepository)
	suite.userService = service.NewUserService(suite.mockUserRepo)
}

func (suite *UserServiceTestSuite) TearDownTest() {
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestGetCurrentUser() {
	testCases := []struct {
		expectedError  error
		mockSetup      func()
		name           string
		expectedErrMsg string
		userID         uint
		checkResult    bool
	}{
		{
			name:   "Success",
			userID: 1,
			mockSetup: func() {
				email := testUserServiceEmailConst
				existingUser := &model.User{
					ID:            1,
					ProviderType:  "email",
					Email:         &email,
					DisplayName:   "Test User",
					EmailVerified: true,
					IsDeleted:     false,
				}
				suite.mockUserRepo.On("GetByID", uint(1)).Return(existingUser, nil)
			},
			expectedError:  nil,
			expectedErrMsg: "",
			checkResult:    true,
		},
		{
			name:   "User not found",
			userID: 999,
			mockSetup: func() {
				suite.mockUserRepo.On("GetByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)
			},
			expectedError:  auth.ErrInvalidCredentials,
			expectedErrMsg: "",
			checkResult:    false,
		},
		{
			name:   "User is deleted",
			userID: 2,
			mockSetup: func() {
				email := "deleted@example.com"
				deletedUser := &model.User{
					ID:            2,
					ProviderType:  "email",
					Email:         &email,
					DisplayName:   "Deleted User",
					EmailVerified: false,
					IsDeleted:     true,
				}
				suite.mockUserRepo.On("GetByID", uint(2)).Return(deletedUser, nil)
			},
			expectedError:  auth.ErrInvalidCredentials,
			expectedErrMsg: "",
			checkResult:    false,
		},
		{
			name:   "Database error",
			userID: 3,
			mockSetup: func() {
				suite.mockUserRepo.On("GetByID", uint(3)).Return(nil, assert.AnError)
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
			result, err := suite.userService.GetCurrentUser(tc.userID)

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
					assert.Equal(t, tc.userID, result.ID)
					assert.Equal(t, testUserServiceEmailConst, result.Email)
					assert.Equal(t, "Test User", result.DisplayName)
					assert.True(t, result.EmailVerified)
				}
			}
		})
	}
}

func (suite *UserServiceTestSuite) TestNewUserService() {
	// Test that NewUserService creates a valid service
	svc := service.NewUserService(suite.mockUserRepo)
	assert.NotNil(suite.T(), svc)
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
