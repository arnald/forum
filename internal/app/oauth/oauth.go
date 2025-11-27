package oauthservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/arnald/forum/internal/domain/oauth"
	"github.com/arnald/forum/internal/domain/user"
	oauthrepo "github.com/arnald/forum/internal/infra/storage/sqlite/oauth"
	oauthpkg "github.com/arnald/forum/internal/pkg/oAuth"
	"github.com/arnald/forum/internal/pkg/uuid"
)

type OAuthService struct {
	oauthRepo    oauth.Repository
	uuidProvider uuid.Provider
}

func NewOAuthService(oauthRepo oauth.Repository, uuidProvider uuid.Provider) *OAuthService {
	return &OAuthService{
		oauthRepo:    oauthRepo,
		uuidProvider: uuidProvider,
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

	providerID := providerUserInfo.ProviderID

	// TODO: PROVIDER NAME VALIDATION
	providerName := oauth.Provider(provider.Name())
	existingUser, err := s.oauthRepo.GetUserByProviderID(ctx, providerName, providerID)
	if err != nil && !errors.Is(err, oauthrepo.ErrUserNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	if existingUser != nil {
		return existingUser, nil
	}

	oauthUser := &oauth.User{
		UserID:     s.uuidProvider.NewUUID(),
		ProviderID: providerID,
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
