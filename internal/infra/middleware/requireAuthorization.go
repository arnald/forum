package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/infra/session"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type Key string

const (
	userIDKey Key = "user"
)

type requireAuthMiddleware struct {
	sessionManager user.SessionManager
}

type RequireAuthMiddleware interface {
	RequireAuth(next http.HandlerFunc) http.HandlerFunc
}

func NewRequireAuthMiddleware(sessionManager user.SessionManager) RequireAuthMiddleware {
	return requireAuthMiddleware{
		sessionManager: sessionManager,
	}
}

func (a requireAuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := r.Cookie("session_token")
		if err != nil {
			helpers.RespondWithError(w,
				http.StatusUnauthorized,
				"Unauthorized: No session cookie found")
			return
		}

		err = a.sessionManager.ValidateSession(sessionID.Value)
		if err != nil {
			switch {
			case errors.Is(err, session.ErrSessionNotFound):
				helpers.RespondWithError(w,
					http.StatusUnauthorized,
					"Unauthorized: Session not found")
			case errors.Is(err, session.ErrSessionExpired):
				helpers.RespondWithError(w,
					http.StatusUnauthorized,
					"Unauthorized: Session expired")
			default:
				helpers.RespondWithError(w,
					http.StatusInternalServerError,
					"Internal Server Error: Unable to validate session")
			}
			return
		}

		user, err := a.sessionManager.GetUserFromSession(sessionID.Value)
		if err != nil {
			switch {
			case errors.Is(err, session.ErrUserNotFound):
				helpers.RespondWithError(w,
					http.StatusUnauthorized,
					"Unauthorized: User not found")
			default:
				helpers.RespondWithError(w,
					http.StatusInternalServerError,
					"Internal Server Error: Unable to retrieve user from session")
			}
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
