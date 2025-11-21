package githubclient

import (
	"encoding/json"
	"fmt"

	"github.com/arnald/forum/internal/pkg/oAuth/httpclient"
)

const (
	tokerURL     = "https://github.com/login/oauth/access_token"
	userURL      = "https://api.github.com/user"
	userEmailURL = "https://api.github.com/user/emails"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type GitHubUser struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

type GitHubEmail struct {
	Email      string `json:"email"`
	Primary    bool   `json:"primary"`
	Verified   bool   `json:"verified"`
	Visibility string `json:"visibility"`
}

func ExchangeCode(code, clientID, clientSecret, redirectURL string) (string, error) {
	body := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"code":          code,
		"redirect_uri":  redirectURL,
	}

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	respBody, err := httpclient.Post(tokerURL, headers, body)
	if err != nil {
		return "", fmt.Errorf("failed to exchange code: %w", err)
	}

	var tokenResp TokenResponse
	err = json.Unmarshal(respBody, &tokenResp)
	if err != nil {
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("no access token in response")
	}

	return tokenResp.AccessToken, nil
}

func GetUser(accessToken string) (*GitHubUser, error) {
	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
		"Accept":        "application/json",
	}

	respoBody, err := httpclient.Get(userURL, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	var user GitHubUser
	err = json.Unmarshal(respoBody, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user response")
	}

	if user.Email == "" {
		email, err := getPrimaryEmail(accessToken)
		if err != nil {
			return nil, fmt.Errorf("failed to get primary email: %w", err)
		}
		user.Email = email
	}

	return &user, nil
}

func getPrimaryEmail(accessToken string) (string, error) {
	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
		"Accept":        "application/json",
	}

	respBody, err := httpclient.Get(userEmailURL, headers)
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

	return "", fmt.Errorf("no verified email found")
}
