package userqueries

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
			name: "successful login",
			request: UserLoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMocks: func(repo *testhelpers.MockRepository, enc *testhelpers.MockEncryptionProvider) {
				repo.GetUserByEmailFunc = func(ctx context.Context, email string) (*user.User, error) {
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
			name: "user not found",
			request: UserLoginRequest{
				Email: "notfound@example.com",
			},
			setupMocks: func(repo *testhelpers.MockRepository, enc *testhelpers.MockEncryptionProvider) {
				repo.GetUserByEmailFunc = func(ctx context.Context, email string) (*user.User, error) {
					return nil, testhelpers.ErrTest
				}
			},
			wantErr:  testhelpers.ErrTest,
			wantUser: nil,
		},
		{
			name: "password mismatch",
			request: UserLoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			setupMocks: func(repo *testhelpers.MockRepository, enc *testhelpers.MockEncryptionProvider) {
				repo.GetUserByEmailFunc = func(ctx context.Context, email string) (*user.User, error) {
					return &user.User{
						ID:       "test-uuid",
						Username: "testuser",
						Email:    "test@example.com",
					}, nil
				}
				enc.MatchesFunc = func(hashedPassword string, plaintextPassword string) error {
					return testhelpers.ErrTest
				}
			},
			wantErr:  testhelpers.ErrTest,
			wantUser: nil,
		},
		{
			name: "encryption provider fails",
			request: UserLoginRequest{
				Email:    "test@example.com",
				Password: "Password123",
			},
			setupMocks: func(repo *testhelpers.MockRepository, enc *testhelpers.MockEncryptionProvider) {
				repo.GetUserByEmailFunc = func(ctx context.Context, email string) (*user.User, error) {
					return &user.User{
						ID:       "test-uuid",
						Username: "testuser",
						Email:    "test@example.com",
					}, nil
				}
				enc.MatchesFunc = func(hashedPassword string, plaintextPassword string) error {
					return testhelpers.ErrTest
				}
			},
			wantErr:  testhelpers.ErrTest,
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
