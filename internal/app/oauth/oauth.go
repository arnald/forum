package oauthservice

import (
	"context"
	"fmt"

	"github.com/arnald/forum/internal/domain/oauth"
	"github.com/arnald/forum/internal/domain/user"
	oauthpkg "github.com/arnald/forum/internal/pkg/oAuth"
)

type OAuthService struct {
	oauthRepo oauth.Repository
}

func NewOAuthService(oauthRepo oauth.Repository) *OAuthService {
	return &OAuthService{
		oauthRepo: oauthRepo,
	}
}

func (s *OAuthService) Login(ctx context.Context, code string, provider oauthpkg.Provider) (*user.User, error) {
	accessToken, err := provider.ExchangeCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	providerUserInfo, err := provider.GetUserInfo(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	gitHubUserID := providerUserInfo.ProviderID

	// TODO: PROVIDER NAME VALIDATION
	providerName := oauth.Provider(provider.Name())
	existingUser, err := s.oauthRepo.GetUserByProviderID(ctx, providerName, gitHubUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	if existingUser != nil {
		return existingUser, nil
	}

	oauthUser := &oauth.User{
		ProviderID: gitHubUserID,
		Provider:   providerName,
		Email:      providerUserInfo.Email,
		Username:   providerUserInfo.Username,
		AvatarURL:  providerUserInfo.AvatarURL,
		Name:       providerUserInfo.Name,
	}

	newUser, err := s.oauthRepo.CreateOAuthUser(ctx, oauthUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, nil
}
