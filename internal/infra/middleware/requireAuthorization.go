package middleware

import (
	"context"
	"net/http"

	"github.com/arnald/forum/internal/domain/session"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type requireAuthMiddleware struct {
	sessionManager session.Manager
}

type RequireAuthMiddleware interface {
	RequireAuth(next http.HandlerFunc) http.HandlerFunc
}

func NewRequireAuthMiddleware(sessionManager session.Manager) RequireAuthMiddleware {
	return requireAuthMiddleware{
		sessionManager: sessionManager,
	}
}

func (a requireAuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionToken, refreshToken := GetTokensFromRequest(r)

		session, err := a.sessionManager.GetSessionFromSessionTokens(sessionToken, refreshToken)
		if err != nil {
			helpers.RespondWithJSON(
				w,
				http.StatusUnauthorized,
				nil,
				"Unauthorized: Invalid session")
			return
		}

		sessionExpired, refreshTokenExpired := CheckTokenExpiration(session)
		switch {
		case sessionExpired && refreshTokenExpired:
			helpers.RespondWithError(w,
				http.StatusUnauthorized,
				"Unauthorized: Session and refresh token expired")
			return
		case sessionExpired && !refreshTokenExpired:
			_ = a.sessionManager.DeleteSession(session.AccessToken)
			session, _ = a.sessionManager.CreateSession(r.Context(), session.UserID)
		case !sessionExpired && refreshTokenExpired:
			helpers.RespondWithError(w,
				http.StatusUnauthorized,
				"Unauthorized: Refresh token expired")
			return
		}

		user, err := a.sessionManager.GetUserFromSession(session.AccessToken)
		if err != nil {
			helpers.RespondWithError(
				w,
				http.StatusUnauthorized,
				"Unauthorized: User not found")
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
