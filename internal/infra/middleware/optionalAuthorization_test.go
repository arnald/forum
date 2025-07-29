package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arnald/forum/internal/domain/user"
	testhelpers "github.com/arnald/forum/internal/pkg/testing"
)

func TestOptionalAuthMiddleware(t *testing.T) {
	t.Run("group: optional authorization", func(t *testing.T) {
		testCases := newOptionalAuthorizationTestCases()
		for _, tt := range testCases {
			t.Run(tt.name, runOptionalAuthorizationTest(tt))
		}
	})
}

type optionalAuthorizationTestCase struct {
	name             string
	cookie           *http.Cookie
	setupMockSession func(*testhelpers.MockSessionManager)
	wantUserID       string
	wantNextCalled   bool
}

func newOptionalAuthorizationTestCases() []optionalAuthorizationTestCase {
	return []optionalAuthorizationTestCase{
		{
			name:   "no cookie present",
			cookie: nil,
			setupMockSession: func(sm *testhelpers.MockSessionManager) {
			},
			wantUserID:     "",
			wantNextCalled: true,
		},
		{
			name: "valid session",
			cookie: &http.Cookie{
				Name:  "session_token",
				Value: "valid-session",
			},
			setupMockSession: func(sm *testhelpers.MockSessionManager) {
				sm.GetSessionFunc = func(sessionID string) (*user.Session, error) {
					return &user.Session{
						UserID: "test-user-id",
					}, nil
				}
			},
			wantUserID:     "test-user-id",
			wantNextCalled: true,
		},
		{
			name: "invalid session",
			cookie: &http.Cookie{
				Name:  "session_token",
				Value: "invalid-session",
			},
			setupMockSession: func(sm *testhelpers.MockSessionManager) {
				sm.GetSessionFunc = func(sessionID string) (*user.Session, error) {
					return nil, testhelpers.ErrTest
				}
			},
			wantUserID:     "",
			wantNextCalled: true,
		},
	}
}

func runOptionalAuthorizationTest(tt optionalAuthorizationTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		mockSessionManager := &testhelpers.MockSessionManager{}
		tt.setupMockSession(mockSessionManager)

		middleware := NewOptionalAuthMiddleware(mockSessionManager)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		if tt.cookie != nil {
			req.AddCookie(tt.cookie)
		}

		rr := httptest.NewRecorder()

		nextCalled := false
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nextCalled = true
			if userID, ok := r.Context().Value(userIDKey).(string); ok {
				if userID != tt.wantUserID {
					t.Errorf("expected user ID %s, got %s", tt.wantUserID, userID)
				}
			} else if tt.wantUserID != "" {
				t.Error("expected user ID to be set in context")
			}
		})

		handler := middleware.OptionalAuth(next)
		handler.ServeHTTP(rr, req)

		if nextCalled != tt.wantNextCalled {
			t.Errorf("next handler called = %v, want %v", nextCalled, tt.wantNextCalled)
		}
	}
}

func TestNewOptionalAuthMiddleware(t *testing.T) {
	mockSessionManager := &testhelpers.MockSessionManager{}
	middleware := NewOptionalAuthMiddleware(mockSessionManager)

	if middleware == nil {
		t.Fatal("NewOptionalAuthMiddleware returned nil")
	}
}
