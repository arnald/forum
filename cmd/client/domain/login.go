package domain

type LoginFormErrors struct {
	Username      string        `json:"-"`
	Email         string        `json:"-"`
	Password      string        `json:"-"`
	UsernameError string        `json:"username,omitempty"`
	EmailError    string        `json:"email,omitempty"`
	PasswordError string        `json:"password,omitempty"`
	User          *LoggedInUser `json:"-"`
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
