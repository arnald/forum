package githubclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	oauth "github.com/arnald/forum/internal/pkg/oAuth"
	"github.com/arnald/forum/internal/pkg/oAuth/httpclient"
)

const (
	tokerURL     = "https://github.com/login/oauth/access_token"
	userURL      = "https://api.github.com/user"
	userEmailURL = "https://api.github.com/user/emails"
)

type GitHubProvider struct {
	clientID     string
	clientSecret string
	redirectURL  string
	scopes       []string
}

func NewProvider(clientID, clientSecret, redirectURL string, scopes []string) *GitHubProvider {
	return &GitHubProvider{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
		scopes:       scopes,
	}
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type GitHubUser struct {
	Login     string `json:"login"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	ID        int    `json:"id"`
}

type GitHubEmail struct {
	Email      string `json:"email"`
	Visibility string `json:"visibility"`
	Primary    bool   `json:"primary"`
	Verified   bool   `json:"verified"`
}

func (p *GitHubProvider) Name() string {
	return "github"
}

func (p *GitHubProvider) GetAuthURL(state string) string {
	params := url.Values{}
	params.Add("client_id", p.clientID)
	params.Add("redirect_uri", p.redirectURL)
	params.Add("scope", strings.Join(p.scopes, " "))
	params.Add("state", state)

	return "https://github.com/login/oauth/authorize?" + params.Encode()
}

func (p *GitHubProvider) ExchangeCode(ctx context.Context, code string) (string, error) {
	body := map[string]string{
		"client_id":     p.clientID,
		"client_secret": p.clientSecret,
		"code":          code,
		"redirect_uri":  p.redirectURL,
	}

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	respBody, err := httpclient.Post(ctx, tokerURL, headers, body)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToExchangeCode, err)
	}

	var tokenResp TokenResponse
	err = json.Unmarshal(respBody, &tokenResp)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToParseToken, err)
	}

	if tokenResp.AccessToken == "" {
		return "", ErrTokenNotFound
	}

	return tokenResp.AccessToken, nil
}

func (p *GitHubProvider) GetUserInfo(ctx context.Context, accessToken string) (*oauth.ProviderUserInfo, error) {
	user, err := p.fetchUser(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	return &oauth.ProviderUserInfo{
		ProviderID: strconv.Itoa(user.ID),
		Email:      user.Email,
		Username:   user.Login,
		Name:       user.Name,
		AvatarURL:  user.AvatarURL,
	}, nil
}

func (p *GitHubProvider) fetchUser(ctx context.Context, accessToken string) (*GitHubUser, error) {
	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
		"Accept":        "application/json",
	}

	respoBody, err := httpclient.Get(ctx, userURL, headers)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetUser, err)
	}

	var user GitHubUser
	err = json.Unmarshal(respoBody, &user)
	if err != nil {
		return nil, ErrFailedToParseUser
	}

	if user.Email == "" {
		email, err := p.getPrimaryEmail(ctx, accessToken)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrFailedToGetPrimaryEmail, err)
		}
		user.Email = email
	}

	return &user, nil
}

func (p *GitHubProvider) getPrimaryEmail(ctx context.Context, accessToken string) (string, error) {
	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
		"Accept":        "application/json",
	}

	respBody, err := httpclient.Get(ctx, userEmailURL, headers)
	if err != nil {
		return "", fmt.Errorf("failed to get emails: %w", err)
	}

	var emails []GitHubEmail
	err = json.Unmarshal(respBody, &emails)
	if err != nil {
		return "", fmt.Errorf("failed to parse emails response: %w", err)
	}

	for _, email := range emails {
		if email.Primary && email.Verified {
			return email.Email, nil
		}
	}

	for _, email := range emails {
		if email.Verified {
			return email.Email, nil
		}
	}

	return "", ErrNoVerifiedEmailFound
}
