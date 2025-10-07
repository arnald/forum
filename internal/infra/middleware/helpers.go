package middleware

import (
	"net/http"
	"time"

	"github.com/arnald/forum/internal/domain/session"
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
