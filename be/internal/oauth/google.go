package oauth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

type GoogleOAuthService struct{}

func NewGoogleOAuthService() *GoogleOAuthService {
	return &GoogleOAuthService{}
}

func (g *GoogleOAuthService) GetUserInfo(accessToken string) (*GoogleUserInfo, error) {
	ctx := context.Background()

	service, err := oauth2.NewService(ctx, option.WithHTTPClient(&http.Client{}))
	if err != nil {
		return nil, fmt.Errorf("failed to create OAuth2 service: %w", err)
	}

	userInfoCall := service.Userinfo.Get()
	userInfoCall.Header().Set("Authorization", "Bearer "+accessToken)

	userInfo, err := userInfoCall.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	verifiedEmail := false
	if userInfo.VerifiedEmail != nil {
		verifiedEmail = *userInfo.VerifiedEmail
	}

	return &GoogleUserInfo{
		ID:            userInfo.Id,
		Email:         userInfo.Email,
		VerifiedEmail: verifiedEmail,
		Name:          userInfo.Name,
		Picture:       userInfo.Picture,
	}, nil
}

func (g *GoogleOAuthService) ValidateAccessToken(accessToken string) error {
	if strings.TrimSpace(accessToken) == "" {
		return fmt.Errorf("access token is empty")
	}

	_, err := g.GetUserInfo(accessToken)
	return err
}