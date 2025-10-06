package session

import (
	"context"
	"net/http"

	"github.com/arnald/forum/internal/domain/user"
)

type Manager interface {
	CreateSession(ctx context.Context, userID string) (*Session, error)
	GetSession(sessionID string) (*Session, error)
	DeleteSession(sessionID string) error
	GetUserFromSession(sessionID string) (*user.User, error)
	GetSessionFromSessionTokens(sessionToken, refreshToken string) (*Session, error)
	ValidateSession(sessionID string) error
	NewSessionCookie(token string) *http.Cookie
	DeleteSessionWhenNewCreated(ctx context.Context, sessionID string, userID string) error
}
