package domain

// BackendErrorResponse - handles backend validation errors (shared across all features)
type BackendErrorResponse struct {
	Message string `json:"message"`
}

// LoggedInUser - user data to pass to templates and store in session
type LoggedInUser struct {
	ID        string
	Username  string
	Email     string
	AvatarURL string // For future navbar avatar display
}
