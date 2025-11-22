package domain

// BackendErrorResponse - handles backend validation errors (shared across all features)
type BackendErrorResponse struct {
	Message string `json:"message"`
}

// BackendMeResponse - response from backend /me endpoint
type BackendMeResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// LoggedInUser - user data to pass to templates and store in session
type LoggedInUser struct {
	ID       string
	Username string
	Email    string
	// AvatarURL string // For future navbar avatar display
}
