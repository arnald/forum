package middleware

import (
	"net/http"
	"time"

	"github.com/arnald/forum/internal/domain/session"
	"github.com/arnald/forum/internal/domain/user"
)

type Key string

const (
	userIDKey Key = "user"
)

func CheckTokenExpiration(session *session.Session) (sessionExpired, refreshTokenExpired bool) {
	if session.Expiry.Before(time.Now()) {
		sessionExpired = true
	}

	if session.RefreshTokenExpiry.Before(time.Now()) {
		refreshTokenExpired = true
	}

	return
}

func GetTokensFromRequest(r *http.Request) (sessionToken, refreshToken string) {
	cookie, err := r.Cookie("session_token")
	if err == nil {
		sessionToken = cookie.Value
	}

	cookie, err = r.Cookie("refresh_token")
	if err == nil {
		refreshToken = cookie.Value
	}

	return
}

func GetUserFromContext(r *http.Request) *user.User {
	value := r.Context().Value(userIDKey)
	if value == nil {
		return nil
	}

	user, ok := value.(*user.User)
	if !ok {
		return nil
	}

	return user
}
