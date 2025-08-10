package oauth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGoogleOAuthService(t *testing.T) {
	service := NewGoogleOAuthService()
	assert.NotNil(t, service)
}

func TestValidateAccessToken(t *testing.T) {
	service := NewGoogleOAuthService()

	tests := []struct {
		name        string
		accessToken string
		expectError bool
	}{
		{
			name:        "empty access token",
			accessToken: "",
			expectError: true,
		},
		{
			name:        "whitespace only access token",
			accessToken: "   ",
			expectError: true,
		},
		{
			name:        "invalid access token",
			accessToken: "invalid_token",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateAccessToken(tt.accessToken)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
