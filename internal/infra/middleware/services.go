package middleware

import (
	"github.com/arnald/forum/internal/domain/session"
)

type Middleware struct {
	Authorization RequireAuthMiddleware
	OptionalAuth  OptionalAuthMiddleware
}

func NewMiddleware(sessionManager session.Manager) *Middleware {
	return &Middleware{
		Authorization: NewRequireAuthMiddleware(sessionManager),
		OptionalAuth:  NewOptionalAuthMiddleware(sessionManager),
	}
}
