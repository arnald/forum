package user

import (
	"context"
	"net/http"
)

type SessionManager interface {
	CreateSession(ctx context.Context, userID string) (*Session, error)
	GetSession(sessionID string) (*Session, error)
	DeleteSession(sessionID string) error
	GetUserFromSession(sessionID string) (*User, error)
	GetSessionFromSessionTokens(sessionToken, refreshToken string) (*Session, error)
	ValidateSession(sessionID string) error
	NewSessionCookie(token string) *http.Cookie
}
