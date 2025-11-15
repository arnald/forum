package domain

type RegisterFormErrors struct {
	Username      string `json:"-"`
	Email         string `json:"-"`
	Password      string `json:"-"`
	UsernameError string `json:"username,omitempty"`
	EmailError    string `json:"email,omitempty"`
	PasswordError string `json:"password,omitempty"`
}

// BackendRegisterRequest - matches backend RegisterUserReguestModel.
type BackendRegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// BackendRegisterResponse - matches backend RegisterUserResponse.
type BackendRegisterResponse struct {
	UserID  string `json:"userId"`
	Message string `json:"Message"`
}

// BackendErrorResponse - handles backend validation errors.
type BackendErrorResponse struct {
	Message string `json:"message"`
}
