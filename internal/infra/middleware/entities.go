package middleware

import (
	"net/http"

	"github.com/arnald/forum/internal/domain/session"
)

type authorization struct {
	sessionManager session.Manager
}
type Authorization interface {
	Required(next http.HandlerFunc) http.HandlerFunc
	Optional(next http.HandlerFunc) http.HandlerFunc
}

func NewAuthorizationMiddleware(sessionManager session.Manager) Authorization {
	return authorization{
		sessionManager: sessionManager,
	}
}
