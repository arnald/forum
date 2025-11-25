package oauth

import "context"

type Provider interface {
	Name() string
	GetAuthURL(state string) string
	ExchangeCode(ctx context.Context, code string) (string, error)
	GetUserInfo(ctx context.Context, accessToken string) (*ProviderUserInfo, error)
}

type ProviderUserInfo struct {
	ProviderID string
	Email      string
	Username   string
	Name       string
	AvatarURL  string
}
