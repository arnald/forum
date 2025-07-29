package middleware

import (
	"context"
	"net/http"

	"github.com/arnald/forum/internal/domain/user"
)

type optionalAuthMiddleware struct {
	sessionManager user.SessionManager
}

type OptionalAuthMiddleware interface {
	OptionalAuth(next http.Handler) http.Handler
}

func NewOptionalAuthMiddleware(sessionManager user.SessionManager) OptionalAuthMiddleware {
	return optionalAuthMiddleware{
		sessionManager: sessionManager,
	}
}

func (a optionalAuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := r.Cookie("session_token")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		session, err := a.sessionManager.GetSession(sessionID.Value)
		if err != nil || session == nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
