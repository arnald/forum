package middleware

import (
	"github.com/arnald/forum/internal/domain/user"
)

type Middleware struct {
	Authorization RequireAuthMiddleware
}

func NewMiddleware(sessionManager user.SessionManager) *Middleware {
	return &Middleware{
		Authorization: NewRequireAuthMiddleware(sessionManager),
	}
}
