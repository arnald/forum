package oauth

type Provider string

const (
	ProviderGitHub Provider = "github"
	ProviderGoogle Provider = "google"
)

type User struct {
	UserID     string
	ProviderID string
	Provider   Provider
	Email      string
	Username   string
	AvatarURL  string
	Name       string
}

type State struct {
	Token     string
	CreatedAt int
}
