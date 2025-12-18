package server

import "time"

const (
	notFoundMessage = "Oops! The page you're looking for has vanished into the digital void."
	backendAPIBase  = "http://localhost:8080/api/v1"
	requestTimeout  = 15 * time.Second
)

const (
	// User endpoints
	backendRegisterURL      = backendAPIBase + "/register"
	backendLoginEmailURL    = backendAPIBase + "/login/email"
	backendLoginUsernameURL = backendAPIBase + "/login/username"
	backendLogoutURL        = backendAPIBase + "/logout"
	// OAuth endpoints
	backendGithubRegister  = backendAPIBase + "/auth/github/login"
	backendGooglebRegister = backendAPIBase + "/auth/google/login"
	// Category endpoints
	backendGetCategoriesDomain = backendAPIBase + "/categories/all"
	// Topic endpoints
	backendGetTopicsDomain = backendAPIBase + "/topics/all"
	backendGetTopicByID    = backendAPIBase + "/topic"
	backendCreateTopic     = backendAPIBase + "/topics/create"
	backendUpdateTopic     = backendAPIBase + "/topics/update"
	backendDeleteTopic     = backendAPIBase + "/topics/delete"
	// Comment endpoints
	backendCreateComment = backendAPIBase + "/comments/create"z
	backendUpdateComment = backendAPIBase + "/comments/update"
	backendDeleteComment = backendAPIBase + "/comments/delete"
	// Vote endpoints
	backendCastVote      = backendAPIBase + "/vote/cast"
	backendDeleteVote    = backendAPIBase + "/vote/delete"
	backendGetVoteCounts = backendAPIBase + "/vote/counts"
)
