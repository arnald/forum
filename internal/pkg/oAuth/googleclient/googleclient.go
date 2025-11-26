package googleclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	oauth "github.com/arnald/forum/internal/pkg/oAuth"
	"github.com/arnald/forum/internal/pkg/oAuth/httpclient"
)

const (
	tokenURL    = "https://oauth2.googleapis.com/token"
	userInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
	authURL     = "https://accounts.google.com/o/oauth2/v2/auth"
)

type GoogleProvider struct {
	clientID     string
	clientSecret string
	redirectURL  string
	scopes       []string
}

func NewProvider(clientID, clientSecret, redirectURL string, scopes []string) *GoogleProvider {
	return &GoogleProvider{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
		scopes:       scopes,
	}
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope"`
	IDToken      string `json:"id_token,omitempty"`
}

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func (p *GoogleProvider) Name() string {
	return "google"
}

func (p *GoogleProvider) GetAuthURL(state string) string {
	params := url.Values{}
	params.Add("client_id", p.clientID)
	params.Add("redirect_uri", p.redirectURL)
	params.Add("response_type", "code")
	params.Add("scope", strings.Join(p.scopes, " "))
	params.Add("state", state)
	params.Add("access_type", "offline")
	params.Add("prompt", "consent")

	return authURL + "?" + params.Encode()
}

func (p *GoogleProvider) ExchangeCode(ctx context.Context, code string) (string, error) {
	formData := url.Values{}
	formData.Set("client_id", p.clientID)
	formData.Set("client_secret", p.clientSecret)
	formData.Set("code", code)
	formData.Set("redirect_uri", p.redirectURL)
	formData.Set("grant_type", "authorization_code")

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/x-www-form-urlencoded",
	}

	respBody, err := httpclient.PostWithURLEncodedParams(ctx, tokenURL, headers, strings.NewReader(formData.Encode()))
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

func (p *GoogleProvider) GetUserInfo(ctx context.Context, accessToken string) (*oauth.ProviderUserInfo, error) {
	user, err := p.fetchUser(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	return &oauth.ProviderUserInfo{
		ProviderID: user.ID,
		Email:      user.Email,
		Username:   user.Name,
		Name:       user.Name,
		AvatarURL:  user.Picture,
	}, nil
}

func (p *GoogleProvider) fetchUser(ctx context.Context, accessToken string) (*GoogleUser, error) {
	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
		"Accept":        "application/json",
	}

	respoBody, err := httpclient.Get(ctx, userInfoURL, headers)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetUser, err)
	}

	var user GoogleUser
	err = json.Unmarshal(respoBody, &user)
	if err != nil {
		return nil, ErrFailedToParseUser
	}

	if !user.VerifiedEmail {
		return nil, ErrEmailNotVerified
	}

	return &user, nil
}
