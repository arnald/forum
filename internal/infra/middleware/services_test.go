package middleware

import (
	"testing"

	testhelpers "github.com/arnald/forum/internal/pkg/testing"
)

func TestServices(t *testing.T) {
	mockSessionManager := &testhelpers.MockSessionManager{}
	middleware := NewMiddleware(mockSessionManager)

	if middleware == nil {
		t.Fatal("NewMiddleware returned nil")
	}

	if middleware.Authorization == nil {
		t.Fatal("Authorization middleware is nil")
	}

	if middleware.OptionalAuth == nil {
		t.Fatal("OptionalAuth middleware is nil")
	}
}
