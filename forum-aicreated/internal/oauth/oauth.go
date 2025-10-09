// Package oauth handles OAuth authentication with external providers
// Supports Google, Facebook, and GitHub authentication
package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

// Config holds OAuth configuration for all providers
// Each provider config includes client credentials, scopes, and redirect URLs
type Config struct {
	Google   *oauth2.Config // Google OAuth configuration
	GitHub   *oauth2.Config // GitHub OAuth configuration
	Facebook *oauth2.Config // Facebook OAuth configuration
}

// UserInfo represents user information from OAuth providers
// This standardizes data from different providers into a common format
type UserInfo struct {
	ID        string // Provider-specific user ID
	Email     string // User's email address
	Username  string // User's display name
	Provider  string // OAuth provider name (google, github, facebook)
	AvatarURL string // Profile picture URL from provider
}

// NewConfig creates OAuth configurations for all providers
// Reads OAuth credentials from environment variables and sets up redirect URLs
// Returns a Config with all three providers configured
func NewConfig() *Config {
	// Get base URL from environment for constructing callback URLs
	baseURL := getEnv("BASE_URL", "http://localhost:8080")

	return &Config{
		// Google OAuth configuration
		Google: &oauth2.Config{
			ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			RedirectURL:  baseURL + "/auth/google/callback", // Callback URL for Google
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",   // Access user email
				"https://www.googleapis.com/auth/userinfo.profile", // Access user profile
			},
			Endpoint: google.Endpoint, // Google OAuth endpoints
		},
		// GitHub OAuth configuration
		GitHub: &oauth2.Config{
			ClientID:     getEnv("GITHUB_CLIENT_ID", ""),
			ClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
			RedirectURL:  baseURL + "/auth/github/callback", // Callback URL for GitHub
			Scopes:       []string{"user:email"},            // Access user email
			Endpoint:     github.Endpoint,                   // GitHub OAuth endpoints
		},
		// Facebook OAuth configuration
		Facebook: &oauth2.Config{
			ClientID:     getEnv("FACEBOOK_CLIENT_ID", ""),
			ClientSecret: getEnv("FACEBOOK_CLIENT_SECRET", ""),
			RedirectURL:  baseURL + "/auth/facebook/callback", // Callback URL for Facebook
			Scopes:       []string{"email", "public_profile"}, // Access email and profile
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://www.facebook.com/v12.0/dialog/oauth",            // Facebook auth URL
				TokenURL: "https://graph.facebook.com/v12.0/oauth/access_token",    // Facebook token URL
			},
		},
	}
}

// GetGoogleUserInfo retrieves user information from Google
// Makes an authenticated request to Google's userinfo API using the OAuth token
// Returns standardized UserInfo with ID, email, name, and provider identifier
func GetGoogleUserInfo(ctx context.Context, token *oauth2.Token, config *oauth2.Config) (*UserInfo, error) {
	// Create HTTP client with OAuth token authentication
	client := config.Client(ctx, token)

	// Request user information from Google's userinfo endpoint
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse JSON response into struct
	var result struct {
		ID      string `json:"id"`      // Google user ID
		Email   string `json:"email"`   // User email
		Name    string `json:"name"`    // User display name
		Picture string `json:"picture"` // Profile picture URL
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	// Return standardized user info
	return &UserInfo{
		ID:        result.ID,
		Email:     result.Email,
		Username:  result.Name,
		Provider:  "google",
		AvatarURL: result.Picture,
	}, nil
}

// GetGitHubUserInfo retrieves user information from GitHub
// GitHub requires two API calls: one for profile data and potentially one for email
// This is because GitHub users can make their email private
func GetGitHubUserInfo(ctx context.Context, token *oauth2.Token, config *oauth2.Config) (*UserInfo, error) {
	// Create HTTP client with OAuth token authentication
	client := config.Client(ctx, token)

	// Get user profile from GitHub API
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse JSON response
	var result struct {
		ID        int    `json:"id"`         // GitHub user ID (numeric)
		Login     string `json:"login"`      // GitHub username
		Email     string `json:"email"`      // User email (may be empty if private)
		AvatarURL string `json:"avatar_url"` // Profile picture URL
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	// If email is not public, fetch it from the emails endpoint
	// GitHub users can hide their email from the public profile
	email := result.Email
	if email == "" {
		emailResp, err := client.Get("https://api.github.com/user/emails")
		if err == nil {
			defer emailResp.Body.Close()
			emailData, err := io.ReadAll(emailResp.Body)
			if err == nil {
				var emails []struct {
					Email   string `json:"email"`   // Email address
					Primary bool   `json:"primary"` // Whether this is the primary email
				}
				if json.Unmarshal(emailData, &emails) == nil {
					// Look for primary email first
					for _, e := range emails {
						if e.Primary {
							email = e.Email
							break
						}
					}
					// If no primary email found, use first available email
					if email == "" && len(emails) > 0 {
						email = emails[0].Email
					}
				}
			}
		}
	}

	// Return standardized user info (convert numeric ID to string)
	return &UserInfo{
		ID:        fmt.Sprintf("%d", result.ID),
		Email:     email,
		Username:  result.Login,
		Provider:  "github",
		AvatarURL: result.AvatarURL,
	}, nil
}

// GetFacebookUserInfo retrieves user information from Facebook
// Uses Facebook Graph API to fetch user profile with id, name, and email
// Returns standardized UserInfo with Facebook-specific provider identifier
func GetFacebookUserInfo(ctx context.Context, token *oauth2.Token, config *oauth2.Config) (*UserInfo, error) {
	// Create HTTP client with OAuth token authentication
	client := config.Client(ctx, token)

	// Request user information from Facebook Graph API with specific fields
	// Include picture field to get profile picture URL
	resp, err := client.Get("https://graph.facebook.com/me?fields=id,name,email,picture")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse JSON response
	var result struct {
		ID    string `json:"id"`    // Facebook user ID
		Email string `json:"email"` // User email
		Name  string `json:"name"`  // User display name
		Picture struct {
			Data struct {
				URL string `json:"url"` // Profile picture URL
			} `json:"data"`
		} `json:"picture"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	// Return standardized user info
	return &UserInfo{
		ID:        result.ID,
		Email:     result.Email,
		Username:  result.Name,
		Provider:  "facebook",
		AvatarURL: result.Picture.Data.URL,
	}, nil
}

// getEnv retrieves environment variable or returns default value
// Used to read OAuth credentials from environment with fallback defaults
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
