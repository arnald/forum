package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/pkg/helpers"
)

type requireAuthMiddleware struct {
	sessionManager user.SessionManager
}

func NewRequireAuthMiddleware(sessionManager user.SessionManager) requireAuthMiddleware {
	return requireAuthMiddleware{
		sessionManager: sessionManager,
	}
}

type RequireAuthMiddleware interface {
	RequireAuth(next http.Handler) http.Handler
}

func (a requireAuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := r.Cookie("session_id")
		if err != nil {
			helpers.RespondWithError(w,
				http.StatusUnauthorized,
				"Unauthorized: No session cookie found")
			return
		}
		session, err := a.sessionManager.GetSession(sessionID.Value)
		if err != nil || session == nil {
			fmt.Println("Error retrieving session:", err)
			helpers.RespondWithError(w,
				http.StatusUnauthorized,
				"Unauthorized: Invalid session")
			return
		}

		ctx := context.WithValue(r.Context(), "userID", session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
