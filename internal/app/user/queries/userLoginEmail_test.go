package userqueries

import (
	"context"
	"errors"
	"testing"

	"github.com/arnald/forum/internal/domain/user"
	testhelpers "github.com/arnald/forum/internal/pkg/testing"
)

func TestUserLoginEmailHandler_Handle(t *testing.T) {
	t.Run("group: user login", func(t *testing.T) {
		testCases := newUserLoginEmailTestCases()
		for _, tt := range testCases {
			t.Run(tt.name, runUserLoginEmailTest(tt))
		}
	})
}

type userLoginEmailTestCase struct {
	name       string
	request    UserLoginEmailRequest
	setupMocks func(*testhelpers.MockRepository, *testhelpers.MockEncryptionProvider)
	wantErr    error
	wantUser   *user.User
}

func newUserLoginEmailTestCases() []userLoginEmailTestCase {
	return []userLoginEmailTestCase{
		{
			name: "successful login",
			request: UserLoginEmailRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMocks: func(repo *testhelpers.MockRepository, enc *testhelpers.MockEncryptionProvider) {
				storedUser := &user.User{
					ID:       "test-uuid",
					Username: "testuser",
					Email:    "test@example.com",
					Password: "hashed_password",
				}
				repo.GetUserByEmailFunc = func(ctx context.Context, email string) (*user.User, error) {
					if email != "test@example.com" {
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
			request: UserLoginEmailRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
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
			name: "wrong password",
			request: UserLoginEmailRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			setupMocks: func(repo *testhelpers.MockRepository, enc *testhelpers.MockEncryptionProvider) {
				storedUser := &user.User{
					ID:       "test-uuid",
					Username: "testuser",
					Email:    "test@example.com",
					Password: "hashed_password",
				}
				repo.GetUserByEmailFunc = func(ctx context.Context, email string) (*user.User, error) {
					return storedUser, nil
				}
				enc.MatchesFunc = func(hashedPassword, plaintextPassword string) error {
					return ErrPasswordMismatch
				}
			},
			wantErr:  ErrPasswordMismatch,
			wantUser: nil,
		},
		{
			name: "invalid email format",
			request: UserLoginEmailRequest{
				Email:    "invalid-email",
				Password: "password123",
			},
			setupMocks: func(repo *testhelpers.MockRepository, enc *testhelpers.MockEncryptionProvider) {
			},
			wantErr:  testhelpers.ErrTest,
			wantUser: nil,
		},
	}
}

func runUserLoginEmailTest(tt userLoginEmailTestCase) func(*testing.T) {
	return func(t *testing.T) {
		repo := &testhelpers.MockRepository{}
		enc := &testhelpers.MockEncryptionProvider{}
		tt.setupMocks(repo, enc)

		handler := NewUserLoginEmailHandler(repo, enc)
		got, err := handler.Handle(context.Background(), tt.request)

		if !errors.Is(err, tt.wantErr) {
			t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
			return
		}

		testhelpers.AssertUserMatch(t, got, tt.wantUser)
	}
}

func TestNewUserLoginEmailHandler(t *testing.T) {
	repo := &testhelpers.MockRepository{}
	enc := &testhelpers.MockEncryptionProvider{}

	got := NewUserLoginEmailHandler(repo, enc)
	if got == nil {
		t.Fatal("NewUserLoginEmailHandler() returned nil")
	}
}
