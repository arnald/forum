package middleware

import (
	"context"
	"net/http"

	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type Key string

const (
	userIDKey Key = "user_id"
)

type requireAuthMiddleware struct {
	sessionManager user.SessionManager
}

type RequireAuthMiddleware interface {
	RequireAuth(next http.Handler) http.Handler
}

func NewRequireAuthMiddleware(sessionManager user.SessionManager) RequireAuthMiddleware {
	return requireAuthMiddleware{
		sessionManager: sessionManager,
	}
}

func (a requireAuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := r.Cookie("session_token")
		if err != nil {
			helpers.RespondWithError(w,
				http.StatusUnauthorized,
				"Unauthorized: No session cookie found")
			return
		}
		session, err := a.sessionManager.GetSession(sessionID.Value)
		if err != nil || session == nil {
			helpers.RespondWithError(w,
				http.StatusUnauthorized,
				"Unauthorized: Invalid session")
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
