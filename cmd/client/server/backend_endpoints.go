package server

import "time"

const (
	notFoundMessage = "Oops! The page you're looking for has vanished into the digital void."
	requestTimeout  = 15 * time.Second
)

// backendAPIBase is set dynamically from client config.
var backendAPIBase = "http://localhost:8080/api/v1"

// Endpoint path constants (without base URL)
const (
	pathRegister             = "/register"
	pathLoginEmail           = "/login/email"
	pathLoginUsername        = "/login/username"
	pathLogout               = "/logout"
	pathMe                   = "/me"
	pathGithubAuth           = "/auth/github/login"
	pathGoogleAuth           = "/auth/google/login"
	pathCategoriesAll        = "/categories/all"
	pathTopicsAll            = "/topics/all"
	pathTopic                = "/topic"
	pathTopicsCreate         = "/topics/create"
	pathTopicsUpdate         = "/topics/update"
	pathTopicsDelete         = "/topics/delete"
	pathCommentsCreate       = "/comments/create"
	pathCommentsUpdate       = "/comments/update"
	pathCommentsDelete       = "/comments/delete"
	pathVoteCast             = "/vote/cast"
	pathVoteDelete           = "/vote/delete"
	pathVoteCounts           = "/vote/counts"
	pathUserActivity         = "/user/activity"
	pathNotificationsStream  = "/notifications/stream"
	pathNotificationsList    = "/notifications"
	pathNotificationsUnread  = "/notifications/unread-count"
	pathNotificationsRead    = "/notifications/mark-read"
	pathNotificationsAllRead = "/notifications/mark-all-read"
)

// Dynamic endpoint URLs
var (
	backendRegisterURL         = func() string { return backendAPIBase + pathRegister }
	backendLoginEmailURL       = func() string { return backendAPIBase + pathLoginEmail }
	backendLoginUsernameURL    = func() string { return backendAPIBase + pathLoginUsername }
	backendLogoutURL           = func() string { return backendAPIBase + pathLogout }
	backendMeURL               = func() string { return backendAPIBase + pathMe }
	backendGithubRegister      = func() string { return backendAPIBase + pathGithubAuth }
	backendGooglebRegister     = func() string { return backendAPIBase + pathGoogleAuth }
	backendGetCategoriesDomain = func() string { return backendAPIBase + pathCategoriesAll }
	backendGetTopicsDomain     = func() string { return backendAPIBase + pathTopicsAll }
	backendGetTopicByID        = func() string { return backendAPIBase + pathTopic }
	backendCreateTopic         = func() string { return backendAPIBase + pathTopicsCreate }
	backendUpdateTopic         = func() string { return backendAPIBase + pathTopicsUpdate }
	backendDeleteTopic         = func() string { return backendAPIBase + pathTopicsDelete }
	backendCreateComment       = func() string { return backendAPIBase + pathCommentsCreate }
	backendUpdateComment       = func() string { return backendAPIBase + pathCommentsUpdate }
	backendDeleteComment       = func() string { return backendAPIBase + pathCommentsDelete }
	backendCastVote            = func() string { return backendAPIBase + pathVoteCast }
	backendDeleteVote          = func() string { return backendAPIBase + pathVoteDelete }
	backendGetVoteCounts       = func() string { return backendAPIBase + pathVoteCounts }
	backedGetUserActivity      = func() string { return backendAPIBase + pathUserActivity }
	backendNotificationsStream = func() string { return backendAPIBase + pathNotificationsStream }
	backendNotificationsList   = func() string { return backendAPIBase + pathNotificationsList }
	backendUnreadCount         = func() string { return backendAPIBase + pathNotificationsUnread }
	backendMarkAsRead          = func() string { return backendAPIBase + pathNotificationsRead }
	backendMarkAllAsRead       = func() string { return backendAPIBase + pathNotificationsAllRead }
)
