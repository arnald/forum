package middleware

import (
	"testing"

	testhelpers "github.com/arnald/forum/internal/pkg/testing"
)

func TestServices(t *testing.T) {
	mockSessionManager := &testhelpers.MockSessionManager{}

	middleware := NewMiddleware(mockSessionManager)

	var (
		auth     = middleware.Authorization
		optional = middleware.OptionalAuth
	)

	if auth == nil {
		t.Error("Authorization middleware is nil")
	}

	if optional == nil {
		t.Error("OptionalAuth middleware is nil")
	}
}
