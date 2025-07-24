package queries

import (
	"context"
	"errors"
	"testing"

	"github.com/arnald/forum/internal/domain/user"
	mocks "github.com/arnald/forum/internal/pkg/testing"
)

func TestUserRegisterHandler_Handle(t *testing.T) {
	t.Run("group: user registration", func(t *testing.T) {
		testCases := newUserRegisterTestCases()
		for _, tt := range testCases {
			t.Run(tt.name, runUserRegisterTest(tt))
		}
	})
}

type userRegisterTestCase struct {
	name       string
	request    UserRegisterRequest
	setupMocks func(*mocks.MockRepository, *mocks.MockUUIDProvider, *mocks.MockEncryptionProvider)
	wantErr    error
	wantUser   *user.User
}

func newUserRegisterTestCases() []userRegisterTestCase {
	testErr := errors.New("test error")

	return []userRegisterTestCase{
		{
			name: "successful registration",
			request: UserRegisterRequest{
				Name:     "testuser",
				Password: "password123",
				Email:    "test@example.com",
			},
			setupMocks: func(repo *mocks.MockRepository, uuid *mocks.MockUUIDProvider, enc *mocks.MockEncryptionProvider) {
				uuid.NewUUIDFunc = func() string { return "test-uuid" }
				enc.GenerateFunc = func(pass string) (string, error) { return "hashed_password", nil }
				repo.UserRegisterFunc = func(ctx context.Context, u *user.User) error { return nil }
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
			name: "encryption fails",
			request: UserRegisterRequest{
				Name:     "testuser",
				Password: "password123",
				Email:    "test@example.com",
			},
			setupMocks: func(repo *mocks.MockRepository, uuid *mocks.MockUUIDProvider, enc *mocks.MockEncryptionProvider) {
				uuid.NewUUIDFunc = func() string { return "test-uuid" }
				enc.GenerateFunc = func(pass string) (string, error) { return "", testErr }
				repo.UserRegisterFunc = func(ctx context.Context, u *user.User) error { return nil }
			},
			wantErr:  testErr,
			wantUser: nil,
		},
		{
			name: "repository fails",
			request: UserRegisterRequest{
				Name:     "testuser",
				Password: "password123",
				Email:    "test@example.com",
			},
			setupMocks: func(repo *mocks.MockRepository, uuid *mocks.MockUUIDProvider, enc *mocks.MockEncryptionProvider) {
				uuid.NewUUIDFunc = func() string { return "test-uuid" }
				enc.GenerateFunc = func(pass string) (string, error) { return "hashed_password", nil }
				repo.UserRegisterFunc = func(ctx context.Context, u *user.User) error { return testErr }
			},
			wantErr:  testErr,
			wantUser: nil,
		},
	}
}

func runUserRegisterTest(tt userRegisterTestCase) func(*testing.T) {
	return func(t *testing.T) {
		repo := &mocks.MockRepository{}
		uuid := &mocks.MockUUIDProvider{}
		enc := &mocks.MockEncryptionProvider{}
		tt.setupMocks(repo, uuid, enc)

		handler := NewUserRegisterHandler(repo, uuid, enc)
		got, err := handler.Handle(context.Background(), tt.request)

		if !errors.Is(err, tt.wantErr) {
			t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
			return
		}

		assertUserMatch(t, got, tt.wantUser)
	}
}

func assertUserMatch(t *testing.T, got, want *user.User) {
	t.Helper()

	if want == nil {
		if got != nil {
			t.Error("Handle() expected nil user, got non-nil")
		}
		return
	}

	if got == nil {
		t.Fatalf("Handle() got nil user, want user with ID %s", want.ID)
		return
	}

	compareUserFields(t, got, want)
}

func compareUserFields(t *testing.T, got, want *user.User) {
	t.Helper()

	if got.ID != want.ID {
		t.Errorf("Handle() got ID = %v, want %v", got.ID, want.ID)
	}
	if got.Username != want.Username {
		t.Errorf("Handle() got Username = %v, want %v", got.Username, want.Username)
	}
	if got.Email != want.Email {
		t.Errorf("Handle() got Email = %v, want %v", got.Email, want.Email)
	}
	if got.Password != want.Password {
		t.Errorf("Handle() got Password = %v, want %v", got.Password, want.Password)
	}
}

func TestNewUserRegisterHandler(t *testing.T) {
	repo := &mocks.MockRepository{}
	uuid := &mocks.MockUUIDProvider{}
	enc := &mocks.MockEncryptionProvider{}

	got := NewUserRegisterHandler(repo, uuid, enc)
	if got == nil {
		t.Fatal("NewUserRegisterHandler() returned nil")
	}
}
