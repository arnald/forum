package userqueries

import (
	"context"
	"errors"
	"testing"

	"github.com/arnald/forum/internal/domain/user"
	testhelpers "github.com/arnald/forum/internal/pkg/testing"
)

func TestUserLoginUsernameHandler_Handle(t *testing.T) {
	t.Run("group: user login", func(t *testing.T) {
		testCases := newUserLoginUsernameTestCases()
		for _, tt := range testCases {
			t.Run(tt.name, runUserLoginUsernameTest(tt))
		}
	})
}

type userLoginUsernameTestCase struct {
	name       string
	request    UserLoginUsernameRequest
	setupMocks func(*testhelpers.MockRepository, *testhelpers.MockEncryptionProvider)
	wantErr    error
	wantUser   *user.User
}

func newUserLoginUsernameTestCases() []userLoginUsernameTestCase {
	return []userLoginUsernameTestCase{
		{
			name: "successful login",
			request: UserLoginUsernameRequest{
				Username: "testuser",
				Password: "password123",
			},
			setupMocks: func(repo *testhelpers.MockRepository, enc *testhelpers.MockEncryptionProvider) {
				storedUser := &user.User{
					ID:       "test-uuid",
					Username: "testuser",
					Email:    "test@example.com",
					Password: "hashed_password",
				}
				repo.GetUserByUsernameFunc = func(ctx context.Context, username string) (*user.User, error) {
					if username != "testuser" {
						return nil, testhelpers.ErrTest
					}
					return storedUser, nil
				}
				enc.MatchesFunc = func(hashedPassword, plaintextPassword string) error {
					if hashedPassword != storedUser.Password {
						return testhelpers.ErrTest
					}
					if plaintextPassword != "password123" {
						return testhelpers.ErrTest
					}
					return nil
				}
			},
			wantErr: nil,
			wantUser: &user.User{
				ID:       "test-uuid",
				Username: "testuser",
				Email:    "test@example.com",
				Password: "hashed_password",
			},
		},
		{
			name: "user not found",
			request: UserLoginUsernameRequest{
				Username: "nonexistent",
				Password: "password123",
			},
			setupMocks: func(repo *testhelpers.MockRepository, enc *testhelpers.MockEncryptionProvider) {
				repo.GetUserByUsernameFunc = func(ctx context.Context, username string) (*user.User, error) {
					return nil, testhelpers.ErrTest
				}
			},
			wantErr:  testhelpers.ErrTest,
			wantUser: nil,
		},
		{
			name: "wrong password",
			request: UserLoginUsernameRequest{
				Username: "testuser",
				Password: "wrongpassword",
			},
			setupMocks: func(repo *testhelpers.MockRepository, enc *testhelpers.MockEncryptionProvider) {
				storedUser := &user.User{
					ID:       "test-uuid",
					Username: "testuser",
					Email:    "test@example.com",
					Password: "hashed_password",
				}
				repo.GetUserByUsernameFunc = func(ctx context.Context, username string) (*user.User, error) {
					return storedUser, nil
				}
				enc.MatchesFunc = func(hashedPassword, plaintextPassword string) error {
					return ErrPasswordMismatch
				}
			},
			wantErr:  ErrPasswordMismatch,
			wantUser: nil,
		},
	}
}

func runUserLoginUsernameTest(tt userLoginUsernameTestCase) func(*testing.T) {
	return func(t *testing.T) {
		repo := &testhelpers.MockRepository{}
		enc := &testhelpers.MockEncryptionProvider{}
		tt.setupMocks(repo, enc)

		handler := NewUserLoginUsernameHandler(repo, enc)
		got, err := handler.Handle(context.Background(), tt.request)

		if !errors.Is(err, tt.wantErr) {
			t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
			return
		}

		testhelpers.AssertUserMatch(t, got, tt.wantUser)
	}
}

func TestNewUserLoginUsernameHandler(t *testing.T) {
	repo := &testhelpers.MockRepository{}
	enc := &testhelpers.MockEncryptionProvider{}

	got := NewUserLoginUsernameHandler(repo, enc)
	if got == nil {
		t.Fatal("NewUserLoginUsernameHandler() returned nil")
	}
}
