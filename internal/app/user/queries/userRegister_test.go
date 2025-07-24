package queries

import (
	"context"
	"errors"
	"testing"

	"github.com/arnald/forum/internal/domain/user"
)

type mockRepository struct {
	userRegisterFunc   func(ctx context.Context, user *user.User) error
	createSessionFunc  func(session *user.Session) error
	getUserByEmailFunc func(ctx context.Context, email string) (*user.User, error)
}

func (m *mockRepository) UserRegister(ctx context.Context, user *user.User) error {
	return m.userRegisterFunc(ctx, user)
}

func (m *mockRepository) CreateSession(session *user.Session) error {
	if m.createSessionFunc != nil {
		return m.createSessionFunc(session)
	}
	return nil
}

func (m *mockRepository) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	if m.getUserByEmailFunc != nil {
		return m.getUserByEmailFunc(ctx, email)
	}
	return nil, nil
}

func (m *mockRepository) GetAll(ctx context.Context) ([]user.User, error) {
	return nil, nil
}

type mockUUIDProvider struct {
	newUUIDFunc func() string
}

func (m *mockUUIDProvider) NewUUID() string {
	return m.newUUIDFunc()
}

type mockEncryptionProvider struct {
	generateFunc func(plaintextPassword string) (string, error)
	matchesFunc  func(hashedPassword string, plaintextPassword string) error
}

func (m *mockEncryptionProvider) Generate(plaintextPassword string) (string, error) {
	return m.generateFunc(plaintextPassword)
}

func (m *mockEncryptionProvider) Matches(hashedPassword string, plaintextPassword string) error {
	if m.matchesFunc != nil {
		return m.matchesFunc(hashedPassword, plaintextPassword)
	}
	return nil
}

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
	setupMocks func(*mockRepository, *mockUUIDProvider, *mockEncryptionProvider)
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
			setupMocks: func(repo *mockRepository, uuid *mockUUIDProvider, enc *mockEncryptionProvider) {
				uuid.newUUIDFunc = func() string { return "test-uuid" }
				enc.generateFunc = func(pass string) (string, error) { return "hashed_password", nil }
				repo.userRegisterFunc = func(ctx context.Context, u *user.User) error { return nil }
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
			setupMocks: func(repo *mockRepository, uuid *mockUUIDProvider, enc *mockEncryptionProvider) {
				uuid.newUUIDFunc = func() string { return "test-uuid" }
				enc.generateFunc = func(pass string) (string, error) { return "", testErr }
				repo.userRegisterFunc = func(ctx context.Context, u *user.User) error { return nil }
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
			setupMocks: func(repo *mockRepository, uuid *mockUUIDProvider, enc *mockEncryptionProvider) {
				uuid.newUUIDFunc = func() string { return "test-uuid" }
				enc.generateFunc = func(pass string) (string, error) { return "hashed_password", nil }
				repo.userRegisterFunc = func(ctx context.Context, u *user.User) error { return testErr }
			},
			wantErr:  testErr,
			wantUser: nil,
		},
	}
}

func runUserRegisterTest(tt userRegisterTestCase) func(*testing.T) {
	return func(t *testing.T) {
		repo := &mockRepository{}
		uuid := &mockUUIDProvider{}
		enc := &mockEncryptionProvider{}
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
	repo := &mockRepository{}
	uuid := &mockUUIDProvider{}
	enc := &mockEncryptionProvider{}

	got := NewUserRegisterHandler(repo, uuid, enc)
	if got == nil {
		t.Fatal("NewUserRegisterHandler() returned nil")
	}
}
