package queries

import (
	"context"
	"errors"
	"testing"

	"github.com/arnald/forum/internal/domain/user"
	testhelpers "github.com/arnald/forum/internal/pkg/testing"
)

func TestUserLoginHandler_Handle(t *testing.T) {
	t.Run("group: user login", func(t *testing.T) {
		testCases := newUserLoginTestCases()
		for _, tt := range testCases {
			t.Run(tt.name, runUserLoginTest(tt))
		}
	})
}

type userLoginTestCase struct {
	name       string
	request    UserLoginRequest
	setupMocks func(*testhelpers.MockRepository, *testhelpers.MockEncryptionProvider)
	wantErr    error
	wantUser   *user.User
}

func newUserLoginTestCases() []userLoginTestCase {
	return []userLoginTestCase{
		{
			name: "successful login with email",
			request: UserLoginRequest{
				Identifier: "test@example.com",
				Password:   "password123",
			},
			setupMocks: func(repo *testhelpers.MockRepository, enc *testhelpers.MockEncryptionProvider) {
				repo.GetUserByIdentifierFunc = func(ctx context.Context, identifier string) (*user.User, error) {
					return &user.User{
						ID:       "test-uuid",
						Username: "testuser",
						Email:    "test@example.com",
					}, nil
				}
				enc.MatchesFunc = func(hashedPassword string, plaintextPassword string) error {
					return nil
				}
			},
			wantErr: nil,
			wantUser: &user.User{
				ID:       "test-uuid",
				Username: "testuser",
				Email:    "test@example.com",
			},
		},
		{
			name: "successful login with username",
			request: UserLoginRequest{
				Identifier: "testuser",
				Password:   "password123",
			},
			setupMocks: func(repo *testhelpers.MockRepository, enc *testhelpers.MockEncryptionProvider) {
				repo.GetUserByIdentifierFunc = func(ctx context.Context, identifier string) (*user.User, error) {
					return &user.User{
						ID:       "test-uuid",
						Username: "testuser",
						Email:    "test@example.com",
					}, nil
				}
				enc.MatchesFunc = func(hashedPassword string, plaintextPassword string) error {
					return nil
				}
			},
			wantErr: nil,
			wantUser: &user.User{
				ID:       "test-uuid",
				Username: "testuser",
				Email:    "test@example.com",
			},
		},
		{
			name: "empty identifier",
			request: UserLoginRequest{
				Identifier: "",
				Password:   "password123",
			},
			setupMocks: func(repo *testhelpers.MockRepository, enc *testhelpers.MockEncryptionProvider) {
			},
			wantErr:  ErrEmptyLoginCreds,
			wantUser: nil,
		},
		{
			name: "empty password",
			request: UserLoginRequest{
				Identifier: "test@example.com",
				Password:   "",
			},
			setupMocks: func(repo *testhelpers.MockRepository, enc *testhelpers.MockEncryptionProvider) {
			},
			wantErr:  ErrEmptyLoginCreds,
			wantUser: nil,
		},
		{
			name: "user not found with email",
			request: UserLoginRequest{
				Identifier: "notfound@example.com",
				Password:   "password123",
			},
			setupMocks: func(repo *testhelpers.MockRepository, enc *testhelpers.MockEncryptionProvider) {
				repo.GetUserByIdentifierFunc = func(ctx context.Context, identifier string) (*user.User, error) {
					return nil, ErrUserNotFound
				}
			},
			wantErr:  ErrUserNotFound,
			wantUser: nil,
		},
		{
			name: "user not found with username",
			request: UserLoginRequest{
				Identifier: "nonexistentuser",
				Password:   "password123",
			},
			setupMocks: func(repo *testhelpers.MockRepository, enc *testhelpers.MockEncryptionProvider) {
				repo.GetUserByIdentifierFunc = func(ctx context.Context, identifier string) (*user.User, error) {
					return nil, ErrUserNotFound
				}
			},
			wantErr:  ErrUserNotFound,
			wantUser: nil,
		},
		{
			name: "password mismatch",
			request: UserLoginRequest{
				Identifier: "test@example.com",
				Password:   "wrongpassword",
			},
			setupMocks: func(repo *testhelpers.MockRepository, enc *testhelpers.MockEncryptionProvider) {
				repo.GetUserByIdentifierFunc = func(ctx context.Context, identifier string) (*user.User, error) {
					return &user.User{
						ID:       "test-uuid",
						Username: "testuser",
						Email:    "test@example.com",
					}, nil
				}
				enc.MatchesFunc = func(hashedPassword string, plaintextPassword string) error {
					return ErrPasswordMismatch
				}
			},
			wantErr:  ErrPasswordMismatch,
			wantUser: nil,
		},
		{
			name: "encryption provider fails",
			request: UserLoginRequest{
				Identifier: "test@example.com",
				Password:   "password123",
			},
			setupMocks: func(repo *testhelpers.MockRepository, enc *testhelpers.MockEncryptionProvider) {
				repo.GetUserByIdentifierFunc = func(ctx context.Context, identifier string) (*user.User, error) {
					return &user.User{
						ID:       "test-uuid",
						Username: "testuser",
						Email:    "test@example.com",
					}, nil
				}
				enc.MatchesFunc = func(hashedPassword string, plaintextPassword string) error {
					return ErrPasswordMismatch
				}
			},
			wantErr:  ErrPasswordMismatch,
			wantUser: nil,
		},
	}
}

func runUserLoginTest(tt userLoginTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		repo := &testhelpers.MockRepository{}
		enc := &testhelpers.MockEncryptionProvider{}

		tt.setupMocks(repo, enc)
		handler := NewUserLoginHandler(repo, enc)
		user, err := handler.Handle(context.Background(), tt.request)
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
			return
		}

		testhelpers.AssertUserMatch(t, user, tt.wantUser)
	}
}

func TestNewUserLoginHandler(t *testing.T) {
	repo := &testhelpers.MockRepository{}
	enc := &testhelpers.MockEncryptionProvider{}

	got := NewUserLoginHandler(repo, enc)
	if got == nil {
		t.Fatal("NewUserLoginHandler() returned nil")
	}
}
