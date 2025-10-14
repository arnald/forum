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
	accessToken      *http.Cookie
	refreshToken     *http.Cookie
	setupMockSession func(*testhelpers.MockSessionManager)
	wantUserID       string
	wantNextCalled   bool
}

func newOptionalAuthorizationTestCases() []optionalAuthorizationTestCase {
	return []optionalAuthorizationTestCase{
		{
			name:         "no cookie present",
			accessToken:  nil,
			refreshToken: nil,
			setupMockSession: func(sm *testhelpers.MockSessionManager) {
			},
			wantUserID:     "",
			wantNextCalled: true,
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
				sm.GetUserFromSessionFunc = func(sessionID string) (*user.User, error) {
					return &user.User{
						ID:       "test-user-id",
						Username: "Test User",
					}, nil
				}
				sm.GetSessionFromSessionTokensFunc = func(sessionToken string, refreshToken string) (*session.Session, error) {
					return &session.Session{
						UserID:             "test-user-id",
						AccessToken:        "valid-session",
						RefreshToken:       "valid-refresh",
						Expiry:             time.Now().Add(1 * time.Hour),
						RefreshTokenExpiry: time.Now().Add(1 * time.Hour),
					}, nil
				}
			},
			wantUserID:     "test-user-id",
			wantNextCalled: true,
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
				sm.GetSessionFromSessionTokensFunc = func(sessionToken string, refreshToken string) (*session.Session, error) {
					if sessionToken == "invalid-session" && refreshToken == "valid-refresh" {
						return &session.Session{
							UserID:             "test-user-id",
							AccessToken:        "valid-session",
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
		},
		{
			name: "invalid session and refresh token",
			accessToken: &http.Cookie{
				Name:  "session_token",
				Value: "invalid-session",
			},
			refreshToken: &http.Cookie{
				Name:  "refresh_token",
				Value: "invalid-refresh",
			},
			setupMockSession: func(sm *testhelpers.MockSessionManager) {
				sm.GetSessionFromSessionTokensFunc = func(sessionToken string, refreshToken string) (*session.Session, error) {
					if sessionToken == "invalid-session" && refreshToken == "invalid-refresh" {
						return nil, testhelpers.ErrTest
					}
					return nil, testhelpers.ErrTest
				}
			},
			wantUserID:     "",
			wantNextCalled: true,
		},
		{
			name: "session and refresh token expired",
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
							Expiry:             time.Now().Add(-1 * time.Hour),
							RefreshTokenExpiry: time.Now().Add(-1 * time.Hour),
						}, nil
					}
					return nil, testhelpers.ErrTest
				}
				sm.GetUserFromSessionFunc = func(sessionID string) (*user.User, error) {
					return nil, testhelpers.ErrTest
				}
				sm.DeleteSessionFunc = func(sessionID string) error {
					return nil
				}
			},
			wantUserID:     "",
			wantNextCalled: true,
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
					return nil, testhelpers.ErrTest
				}
				sm.DeleteSessionFunc = func(sessionID string) error {
					return nil
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

		handler := middleware.Optional(next)
		handler.ServeHTTP(rr, req)

		if nextCalled != tt.wantNextCalled {
			t.Errorf("next handler called = %v, want %v", nextCalled, tt.wantNextCalled)
		}
	}
}

func TestNewOptionalAuthMiddleware(t *testing.T) {
	mockSessionManager := &testhelpers.MockSessionManager{}
	middleware := NewAuthorizationMiddleware(mockSessionManager)

	if middleware == nil {
		t.Fatal("NewOptionalAuthMiddleware returned nil")
	}
}
