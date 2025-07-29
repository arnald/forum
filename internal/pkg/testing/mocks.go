package testhelpers

import (
	"context"
	"errors"
	"net/http"

	"github.com/arnald/forum/internal/domain/user"
)

var (
	ErrTest = errors.New("test error")
)

type MockRepository struct {
	UserRegisterFunc   func(ctx context.Context, user *user.User) error
	GetUserByEmailFunc func(ctx context.Context, email string) (*user.User, error)
	GetAllFunc         func(ctx context.Context) ([]user.User, error)
}

func (m *MockRepository) UserRegister(ctx context.Context, user *user.User) error {
	return m.UserRegisterFunc(ctx, user)
}

func (m *MockRepository) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	if m.GetUserByEmailFunc != nil {
		return m.GetUserByEmailFunc(ctx, email)
	}
	return nil, ErrTest
}

func (m *MockRepository) GetAll(ctx context.Context) ([]user.User, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}
	return nil, ErrTest
}

type MockUUIDProvider struct {
	NewUUIDFunc func() string
}

func (m *MockUUIDProvider) NewUUID() string {
	return m.NewUUIDFunc()
}

type MockEncryptionProvider struct {
	GenerateFunc func(plaintextPassword string) (string, error)
	MatchesFunc  func(hashedPassword string, plaintextPassword string) error
}

func (m *MockEncryptionProvider) Generate(plaintextPassword string) (string, error) {
	return m.GenerateFunc(plaintextPassword)
}

func (m *MockEncryptionProvider) Matches(hashedPassword string, plaintextPassword string) error {
	if m.MatchesFunc != nil {
		return m.MatchesFunc(hashedPassword, plaintextPassword)
	}
	return nil
}

type MockSessionManager struct {
	GetSessionFunc       func(sessionID string) (*user.Session, error)
	CreateSessionFunc    func(userID string) (*user.Session, error)
	DeleteSessionFunc    func(sessionID string) error
	NewSessionCookieFunc func(token string) *http.Cookie
}

func (m *MockSessionManager) GetSession(sessionID string) (*user.Session, error) {
	if m.GetSessionFunc != nil {
		return m.GetSessionFunc(sessionID)
	}
	return nil, ErrTest
}

func (m *MockSessionManager) CreateSession(ctx context.Context, userID string) (*user.Session, error) {
	if m.CreateSessionFunc != nil {
		return m.CreateSessionFunc(userID)
	}
	return nil, ErrTest
}

func (m *MockSessionManager) DeleteSession(sessionID string) error {
	if m.DeleteSessionFunc != nil {
		return m.DeleteSessionFunc(sessionID)
	}
	return ErrTest
}

func (m *MockSessionManager) NewSessionCookie(token string) *http.Cookie {
	if m.NewSessionCookieFunc != nil {
		return m.NewSessionCookieFunc(token)
	}
	return &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	}
}
