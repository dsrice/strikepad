package service

import (
	"testing"

	"strikepad-backend/internal/auth"
	"strikepad-backend/internal/dto"
	"strikepad-backend/internal/model"
	"strikepad-backend/internal/oauth"
	"strikepad-backend/internal/repository/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockGoogleOAuthService is a mock for GoogleOAuthService
type MockGoogleOAuthService struct {
	mock.Mock
}

func (m *MockGoogleOAuthService) GetUserInfo(accessToken string) (*oauth.GoogleUserInfo, error) {
	args := m.Called(accessToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*oauth.GoogleUserInfo), args.Error(1)
}

func (m *MockGoogleOAuthService) ValidateAccessToken(accessToken string) error {
	args := m.Called(accessToken)
	return args.Error(0)
}

func TestAuthService_GoogleSignup(t *testing.T) {
	mockUserRepo := &mocks.MockUserRepository{}
	authService := &AuthService{
		userRepo: mockUserRepo,
	}

	tests := []struct {
		request        *dto.GoogleSignupRequest
		setupMocks     func()
		expectedResult *dto.SignupResponse
		name           string
		expectedError  bool
	}{
		{
			name: "successful Google signup",
			request: &dto.GoogleSignupRequest{
				AccessToken: "valid_token",
			},
			setupMocks: func() {
				// Mock user repository calls
				mockUserRepo.On("FindByEmail", "test@example.com").Return(nil, gorm.ErrRecordNotFound)
				mockUserRepo.On("Create", mock.AnythingOfType("*model.User")).Return(&model.User{
					ID:            1,
					Email:         &[]string{"test@example.com"}[0],
					DisplayName:   "Test User",
					ProviderType:  "google",
					EmailVerified: true,
				}, nil)
			},
			expectedError: false,
		},
		{
			name: "user already exists",
			request: &dto.GoogleSignupRequest{
				AccessToken: "valid_token",
			},
			setupMocks: func() {
				existingUser := &model.User{
					ID:           1,
					Email:        &[]string{"test@example.com"}[0],
					DisplayName:  "Existing User",
					ProviderType: "email",
				}
				mockUserRepo.On("FindByEmail", "test@example.com").Return(existingUser, nil)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockUserRepo.ExpectedCalls = nil
			mockUserRepo.Calls = nil

			if tt.setupMocks != nil {
				tt.setupMocks()
			}

			result, err := authService.GoogleSignup(tt.request)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_GoogleLogin(t *testing.T) {
	mockUserRepo := &mocks.MockUserRepository{}
	authService := &AuthService{
		userRepo: mockUserRepo,
	}

	tests := []struct {
		expectedError error
		request       *dto.GoogleLoginRequest
		setupMocks    func()
		name          string
	}{
		{
			name: "successful Google login",
			request: &dto.GoogleLoginRequest{
				AccessToken: "valid_token",
			},
			setupMocks: func() {
				googleUserID := "google_id_123"
				user := &model.User{
					ID:             1,
					Email:          &[]string{"test@example.com"}[0],
					DisplayName:    "Test User",
					ProviderType:   "google",
					ProviderUserID: &googleUserID,
					EmailVerified:  true,
					IsDeleted:      false,
				}
				mockUserRepo.On("FindByEmail", "test@example.com").Return(user, nil)
			},
			expectedError: nil,
		},
		{
			name: "user not found",
			request: &dto.GoogleLoginRequest{
				AccessToken: "valid_token",
			},
			setupMocks: func() {
				mockUserRepo.On("FindByEmail", "test@example.com").Return(nil, gorm.ErrRecordNotFound)
			},
			expectedError: auth.ErrInvalidCredentials,
		},
		{
			name: "wrong provider type",
			request: &dto.GoogleLoginRequest{
				AccessToken: "valid_token",
			},
			setupMocks: func() {
				user := &model.User{
					ID:           1,
					Email:        &[]string{"test@example.com"}[0],
					DisplayName:  "Test User",
					ProviderType: "email",
					IsDeleted:    false,
				}
				mockUserRepo.On("FindByEmail", "test@example.com").Return(user, nil)
			},
			expectedError: auth.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockUserRepo.ExpectedCalls = nil
			mockUserRepo.Calls = nil

			if tt.setupMocks != nil {
				tt.setupMocks()
			}

			result, err := authService.GoogleLogin(tt.request)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			mockUserRepo.AssertExpectations(t)
		})
	}
}
