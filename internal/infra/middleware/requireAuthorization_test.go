package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arnald/forum/internal/domain/user"
	testhelpers "github.com/arnald/forum/internal/pkg/testing"
)

func TestRequireAuthMiddleware(t *testing.T) {
	t.Run("group: require authorization", func(t *testing.T) {
		testCases := newRequireAuthorizationTestCases()
		for _, tt := range testCases {
			t.Run(tt.name, runRequireAuthorizationTest(tt))
		}
	})
}

type requireAuthorizationTestCase struct {
	name             string
	cookie           *http.Cookie
	setupMockSession func(*testhelpers.MockSessionManager)
	wantUserID       string
	wantNextCalled   bool
	wantStatusCode   int
}

func newRequireAuthorizationTestCases() []requireAuthorizationTestCase {
	return []requireAuthorizationTestCase{
		{
			name:   "no cookie present",
			cookie: nil,
			setupMockSession: func(sm *testhelpers.MockSessionManager) {
			},
			wantUserID:     "",
			wantNextCalled: false,
			wantStatusCode: http.StatusUnauthorized,
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
			wantStatusCode: http.StatusOK,
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
			wantNextCalled: false,
			wantStatusCode: http.StatusUnauthorized,
		},
	}
}

func runRequireAuthorizationTest(tt requireAuthorizationTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		mockSessionManager := &testhelpers.MockSessionManager{}
		tt.setupMockSession(mockSessionManager)

		middleware := NewRequireAuthMiddleware(mockSessionManager)

		req := httptest.NewRequest("GET", "/", nil)
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

		handler := middleware.RequireAuth(next)
		handler.ServeHTTP(rr, req)

		if nextCalled != tt.wantNextCalled {
			t.Errorf("next handler called = %v, want %v", nextCalled, tt.wantNextCalled)
		}

		if rr.Code != tt.wantStatusCode {
			t.Errorf("status code = %v, want %v", rr.Code, tt.wantStatusCode)
		}
	}
}

func TestNewRequireAuthMiddleware(t *testing.T) {
	mockSessionManager := &testhelpers.MockSessionManager{}
	middleware := NewRequireAuthMiddleware(mockSessionManager)

	if middleware == nil {
		t.Fatal("NewRequireAuthMiddleware returned nil")
	}
}
