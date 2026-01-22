package server

import "time"

const (
	notFoundMessage = "Oops! The page you're looking for has vanished into the digital void."
	requestTimeout  = 15 * time.Second
)

// Endpoint path constants (without base URL).
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

// BackendURLs holds all backend API endpoint URLs.
type BackendURLs struct {
	baseURL string
}

// NewBackendURLs creates a new BackendURLs instance.
func NewBackendURLs(baseURL string) *BackendURLs {
	return &BackendURLs{baseURL: baseURL}
}

// Methods to get full URLs.
func (b *BackendURLs) RegisterURL() string            { return b.baseURL + pathRegister }
func (b *BackendURLs) LoginEmailURL() string          { return b.baseURL + pathLoginEmail }
func (b *BackendURLs) LoginUsernameURL() string       { return b.baseURL + pathLoginUsername }
func (b *BackendURLs) LogoutURL() string              { return b.baseURL + pathLogout }
func (b *BackendURLs) MeURL() string                  { return b.baseURL + pathMe }
func (b *BackendURLs) GithubRegisterURL() string      { return b.baseURL + pathGithubAuth }
func (b *BackendURLs) GoogleRegisterURL() string      { return b.baseURL + pathGoogleAuth }
func (b *BackendURLs) CategoriesAllURL() string       { return b.baseURL + pathCategoriesAll }
func (b *BackendURLs) TopicsAllURL() string           { return b.baseURL + pathTopicsAll }
func (b *BackendURLs) TopicURL() string               { return b.baseURL + pathTopic }
func (b *BackendURLs) CreateTopicURL() string         { return b.baseURL + pathTopicsCreate }
func (b *BackendURLs) UpdateTopicURL() string         { return b.baseURL + pathTopicsUpdate }
func (b *BackendURLs) DeleteTopicURL() string         { return b.baseURL + pathTopicsDelete }
func (b *BackendURLs) CreateCommentURL() string       { return b.baseURL + pathCommentsCreate }
func (b *BackendURLs) UpdateCommentURL() string       { return b.baseURL + pathCommentsUpdate }
func (b *BackendURLs) DeleteCommentURL() string       { return b.baseURL + pathCommentsDelete }
func (b *BackendURLs) CastVoteURL() string            { return b.baseURL + pathVoteCast }
func (b *BackendURLs) DeleteVoteURL() string          { return b.baseURL + pathVoteDelete }
func (b *BackendURLs) VoteCountsURL() string          { return b.baseURL + pathVoteCounts }
func (b *BackendURLs) UserActivityURL() string        { return b.baseURL + pathUserActivity }
func (b *BackendURLs) NotificationsStreamURL() string { return b.baseURL + pathNotificationsStream }
func (b *BackendURLs) NotificationsListURL() string   { return b.baseURL + pathNotificationsList }
func (b *BackendURLs) UnreadCountURL() string         { return b.baseURL + pathNotificationsUnread }
func (b *BackendURLs) MarkAsReadURL() string          { return b.baseURL + pathNotificationsRead }
func (b *BackendURLs) MarkAllAsReadURL() string       { return b.baseURL + pathNotificationsAllRead }
