package server

import "time"

const (
	notFoundMessage = "Oops! The page you're looking for has vanished into the digital void."
	backendAPIBase  = "http://localhost:8080/api/v1"
	requestTimeout  = 15 * time.Second
)

const (
	backendRegisterURL         = backendAPIBase + "/register"
	backendLoginEmailURL       = backendAPIBase + "/login/email"
	backendLoginUsernameURL    = backendAPIBase + "/login/username"
	backendLogoutURL           = backendAPIBase + "/logout"
	backendGithubRegister      = backendAPIBase + "/auth/github/login"
	backendGooglebRegister     = backendAPIBase + "/auth/google/login"
	backendGetCategoriesDomain = backendAPIBase + "/categories/all"
	backendGetTopicsDomain     = backendAPIBase + "/topics/all"
	backendGetTopicByID        = backendAPIBase + "/topic"
)
