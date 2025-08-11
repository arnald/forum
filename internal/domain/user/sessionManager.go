package user

import (
	"context"
	"net/http"
)

type SessionManager interface {
	CreateSession(ctx context.Context, userID string) (*Session, error)
	GetSession(sessionID string) (*Session, error)
	DeleteSession(sessionID string) error
	NewSessionCookie(token string) *http.Cookie
	DeleteSessionWhenNewCreated(ctx context.Context, sessionID string, userID string) error
}
