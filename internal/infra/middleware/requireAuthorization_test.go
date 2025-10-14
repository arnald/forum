//nolint:gocognit // TODO: refactor this test
package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/arnald/forum/internal/domain/session"
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
	accessToken      *http.Cookie
	refreshToken     *http.Cookie
	setupMockSession func(*testhelpers.MockSessionManager)
	wantUserID       string
	wantNextCalled   bool
	wantStatusCode   int
}

func newRequireAuthorizationTestCases() []requireAuthorizationTestCase {
	return []requireAuthorizationTestCase{
		{
			name:         "no cookie present",
			accessToken:  nil,
			refreshToken: nil,
			setupMockSession: func(sm *testhelpers.MockSessionManager) {
			},
			wantUserID:     "",
			wantNextCalled: false,
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "valid session",
			accessToken: &http.Cookie{
				Name:  "session_token",
				Value: "valid-session",
			},
			refreshToken: &http.Cookie{
				Name:  "refresh_token",
				Value: "valid-refresh",
			},
			setupMockSession: func(sm *testhelpers.MockSessionManager) {
				sm.GetSessionFromSessionTokensFunc = func(sessionToken, refreshToken string) (*session.Session, error) {
					if sessionToken == "valid-session" && refreshToken == "valid-refresh" {
						return &session.Session{
							UserID:             "test-user-id",
							AccessToken:        "valid-session",
							RefreshToken:       "valid-refresh",
							Expiry:             time.Now().Add(1 * time.Hour),
							RefreshTokenExpiry: time.Now().Add(1 * time.Hour),
						}, nil
					}
					return nil, testhelpers.ErrTest
				}
				sm.GetUserFromSessionFunc = func(sessionID string) (*user.User, error) {
					if sessionID == "valid-session" {
						return &user.User{
							ID:       "test-user-id",
							Username: "Test User",
						}, nil
					}
					return nil, testhelpers.ErrTest
				}
			},
			wantUserID:     "test-user-id",
			wantNextCalled: true,
			wantStatusCode: http.StatusOK,
		},
		{
			name: "invalid session",
			accessToken: &http.Cookie{
				Name:  "session_token",
				Value: "invalid-session",
			},
			refreshToken: &http.Cookie{
				Name:  "refresh_token",
				Value: "invalid-refresh",
			},
			setupMockSession: func(sm *testhelpers.MockSessionManager) {
				sm.GetSessionFromSessionTokensFunc = func(sessionToken, refreshToken string) (*session.Session, error) {
					if sessionToken == "invalid-session" && refreshToken == "invalid-refresh" {
						return nil, testhelpers.ErrTest
					}
					return nil, testhelpers.ErrTest
				}
			},
			wantUserID:     "",
			wantNextCalled: false,
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "session expired, refresh token valid",
			accessToken: &http.Cookie{
				Name:  "session_token",
				Value: "invalid-session",
			},
			refreshToken: &http.Cookie{
				Name:  "refresh_token",
				Value: "valid-refresh",
			},
			setupMockSession: func(sm *testhelpers.MockSessionManager) {
				sm.GetSessionFromSessionTokensFunc = func(sessionToken, refreshToken string) (*session.Session, error) {
					if sessionToken == "invalid-session" && refreshToken == "valid-refresh" {
						return &session.Session{
							UserID:             "test-user-id",
							AccessToken:        "invalid-session",
							RefreshToken:       "valid-refresh",
							Expiry:             time.Now().Add(-1 * time.Hour),
							RefreshTokenExpiry: time.Now().Add(1 * time.Hour),
						}, nil
					}
					return nil, testhelpers.ErrTest
				}
				sm.DeleteSessionFunc = func(sessionID string) error {
					return nil
				}
				sm.CreateSessionFunc = func(userID string) (*session.Session, error) {
					return &session.Session{
						UserID:             userID,
						AccessToken:        "new-access-token",
						RefreshToken:       "new-refresh-token",
						Expiry:             time.Now().Add(1 * time.Hour),
						RefreshTokenExpiry: time.Now().Add(1 * time.Hour),
					}, nil
				}
				sm.GetUserFromSessionFunc = func(sessionID string) (*user.User, error) {
					if sessionID == "new-access-token" {
						return &user.User{
							ID:       "test-user-id",
							Username: "Test User",
						}, nil
					}
					return nil, testhelpers.ErrTest
				}
			},
			wantUserID:     "test-user-id",
			wantNextCalled: true,
			wantStatusCode: http.StatusOK,
		},
		{
			name: "valid session, expired refresh token",
			accessToken: &http.Cookie{
				Name:  "session_token",
				Value: "valid-session",
			},
			refreshToken: &http.Cookie{
				Name:  "refresh_token",
				Value: "expired-refresh",
			},
			setupMockSession: func(sm *testhelpers.MockSessionManager) {
				sm.GetSessionFromSessionTokensFunc = func(sessionToken, refreshToken string) (*session.Session, error) {
					if sessionToken == "valid-session" && refreshToken == "expired-refresh" {
						return &session.Session{
							UserID:             "test-user-id",
							AccessToken:        "valid-session",
							RefreshToken:       "expired-refresh",
							Expiry:             time.Now().Add(1 * time.Hour),
							RefreshTokenExpiry: time.Now().Add(-1 * time.Hour),
						}, nil
					}
					return nil, testhelpers.ErrTest
				}
			},
			wantUserID:     "",
			wantNextCalled: false,
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "valid session, expired refresh token",
			accessToken: &http.Cookie{
				Name:  "session_token",
				Value: "valid-session",
			},
			refreshToken: &http.Cookie{
				Name:  "refresh_token",
				Value: "expired-refresh",
			},
			setupMockSession: func(sm *testhelpers.MockSessionManager) {
				sm.GetSessionFromSessionTokensFunc = func(sessionToken, refreshToken string) (*session.Session, error) {
					if sessionToken == "valid-session" && refreshToken == "expired-refresh" {
						return &session.Session{
							UserID:             "test-user-id",
							AccessToken:        "valid-session",
							RefreshToken:       "expired-refresh",
							Expiry:             time.Now().Add(1 * time.Hour),
							RefreshTokenExpiry: time.Now().Add(-1 * time.Hour),
						}, nil
					}
					return nil, testhelpers.ErrTest
				}
			},
			wantUserID:     "",
			wantNextCalled: false,
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "session and refresh token expired",
			accessToken: &http.Cookie{
				Name:  "session_token",
				Value: "expired-session",
			},
			refreshToken: &http.Cookie{
				Name:  "refresh_token",
				Value: "expired-refresh",
			},
			setupMockSession: func(sm *testhelpers.MockSessionManager) {
				sm.GetSessionFromSessionTokensFunc = func(sessionToken, refreshToken string) (*session.Session, error) {
					if sessionToken == "expired-session" && refreshToken == "expired-refresh" {
						return &session.Session{
							UserID:             "test-user-id",
							AccessToken:        "expired-session",
							RefreshToken:       "expired-refresh",
							Expiry:             time.Now().Add(-1 * time.Hour),
							RefreshTokenExpiry: time.Now().Add(-1 * time.Hour),
						}, nil
					}
					return nil, testhelpers.ErrTest
				}
			},
			wantUserID:     "",
			wantNextCalled: false,
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "user not found",
			accessToken: &http.Cookie{
				Name:  "session_token",
				Value: "valid-session",
			},
			refreshToken: &http.Cookie{
				Name:  "refresh_token",
				Value: "valid-refresh",
			},
			setupMockSession: func(sm *testhelpers.MockSessionManager) {
				sm.GetSessionFromSessionTokensFunc = func(sessionToken, refreshToken string) (*session.Session, error) {
					if sessionToken == "valid-session" && refreshToken == "valid-refresh" {
						return &session.Session{
							UserID:             "non-existent-user",
							AccessToken:        "valid-session",
							RefreshToken:       "valid-refresh",
							Expiry:             time.Now().Add(1 * time.Hour),
							RefreshTokenExpiry: time.Now().Add(1 * time.Hour),
						}, nil
					}
					return nil, testhelpers.ErrTest
				}
				sm.GetUserFromSessionFunc = func(sessionID string) (*user.User, error) {
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

		middleware := NewAuthorizationMiddleware(mockSessionManager)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		if tt.accessToken != nil {
			req.AddCookie(tt.accessToken)
		}
		if tt.refreshToken != nil {
			req.AddCookie(tt.refreshToken)
		}

		rr := httptest.NewRecorder()

		nextCalled := false
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nextCalled = true
			if user, ok := r.Context().Value(userIDKey).(*user.User); ok {
				if user.ID != tt.wantUserID {
					t.Errorf("expected user ID %s, got %s", tt.wantUserID, user.ID)
				}
			} else if tt.wantUserID != "" {
				t.Error("expected user ID to be set in context")
			}
		})

		handler := middleware.Required(next)
		handler.ServeHTTP(rr, req)

		if nextCalled != tt.wantNextCalled {
			t.Errorf("next handler called = %v, want %v", nextCalled, tt.wantNextCalled)
		}

		if rr.Code != tt.wantStatusCode {
			t.Errorf("status code = %v, want %v", rr.Code, tt.wantStatusCode)
		}
	}
}

func TestNewAuthorizationMiddleware(t *testing.T) {
	mockSessionManager := &testhelpers.MockSessionManager{}
	middleware := NewAuthorizationMiddleware(mockSessionManager)

	if middleware == nil {
		t.Fatal("NewRequireAuthMiddleware returned nil")
	}
}
