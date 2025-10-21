package middleware

import (
	"context"
	"net/http"
)

// type optionalAuthMiddleware struct {
// 	sessionManager session.Manager
// }

// type OptionalAuthMiddleware interface {
// 	Optional(next http.HandlerFunc) http.HandlerFunc
// }

// func NewOptionalAuthMiddleware(sessionManager session.Manager) AuthorizationInterface {
// 	return Authorization{
// 		sessionManager: sessionManager,
// 	}
// }

func (a authorization) Optional(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionToken, refreshToken := GetTokensFromRequest(r)
		if sessionToken == "" && refreshToken == "" {
			next.ServeHTTP(w, r)
			return
		}

		session, err := a.sessionManager.GetSessionFromSessionTokens(sessionToken, refreshToken)
		if err != nil || session == nil {
			next.ServeHTTP(w, r)
			return
		}

		sessionExpired, refreshTokenExpired := CheckTokenExpiration(session)
		if sessionExpired && !refreshTokenExpired {
			_ = a.sessionManager.DeleteSession(session.AccessToken)
			session, _ = a.sessionManager.CreateSession(r.Context(), session.UserID)
		} else if sessionExpired && refreshTokenExpired {
			_ = a.sessionManager.DeleteSession(session.AccessToken)
			next.ServeHTTP(w, r)
			return
		}

		user, err := a.sessionManager.GetUserFromSession(session.AccessToken)
		if err != nil || user == nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
