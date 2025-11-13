package domain

type LoginFormErrors struct {
	Identifier      string        `json:"-"`
	Password        string        `json:"-"`
	IdentifierError string        `json:"-"`
	PasswordError   string        `json:"-"`
	User            *LoggedInUser `json:"-"`
}

// BackendLoginRequest - sent to backend
type BackendLoginRequest struct {
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password"`
}

// BackendLoginResponse - response from backend
type BackendLoginResponse struct {
	UserID       string `json:"userId"`
	Username     string `json:"username"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
