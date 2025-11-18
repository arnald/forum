package oauth

import (
	"context"
	"fmt"

	"github.com/arnald/forum/internal/domain/oauth"
	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/pkg/oAuth/githubclient"
)

type GitHubLoginService struct {
	oauthRepo oauth.Repository
}

func NewGitHubLoginService(oauthRepo oauth.Repository) *GitHubLoginService {
	return &GitHubLoginService{
		oauthRepo: oauthRepo,
	}
}

func (s *GitHubLoginService) Login(ctx context.Context, code, clientID, clientSecret, redirectURL string) (*user.User, error) {
	accessToken, err := githubclient.ExchangeCode(code, clientID, clientSecret, redirectURL)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	githubUser, err := githubclient.GetUser(accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	gitHubUserID := fmt.Sprintf("%d", githubUser.ID)
	existingUser, err := s.oauthRepo.GetUserByProviderID(ctx, oauth.ProviderGitHub, gitHubUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	if existingUser != nil {
		return existingUser, nil
	}

	oauthUser := &oauth.User{
		ProviderID: gitHubUserID,
		Provider:   oauth.ProviderGitHub,
		Email:      githubUser.Email,
		Username:   githubUser.Login,
		AvatarURL:  githubUser.AvatarURL,
		Name:       githubUser.Name,
	}

	newUser, err := s.oauthRepo.CreateOAuthUser(ctx, oauthUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, nil
}
