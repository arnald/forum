package testhelpers

import (
	"context"
	"errors"
	"net/http"

	"github.com/arnald/forum/internal/domain/session"
	"github.com/arnald/forum/internal/domain/user"
)

var ErrTest = errors.New("test error")

type MockRepository struct {
	UserRegisterFunc        func(ctx context.Context, user *user.User) error
	GetUserByEmailFunc      func(ctx context.Context, email string) (*user.User, error)
	GetUserByUsernameFunc   func(ctx context.Context, username string) (*user.User, error)
	GetAllFunc              func(ctx context.Context) ([]user.User, error)
	CreateTopicFunc         func(ctx context.Context, topic *user.Topic) error
	UpdateTopicFunc         func(ctx context.Context, topic *user.Topic) error
	DeleteTopicFunc         func(ctx context.Context, userID string, topicID int) error
	GetTopicByIDFunc        func(ctx context.Context, topicID int) (*user.Topic, error)
	GetAllTopicsFunc        func(ctx context.Context, page, size int, orderBy, filter string) ([]user.Topic, error)
	GetTotalTopicsCountFunc func(ctx context.Context, filter string) (int, error)
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

func (m *MockRepository) GetUserByUsername(ctx context.Context, username string) (*user.User, error) {
	if m.GetUserByUsernameFunc != nil {
		return m.GetUserByUsernameFunc(ctx, username)
	}
	return nil, ErrTest
}

func (m *MockRepository) GetAll(ctx context.Context) ([]user.User, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}
	return nil, ErrTest
}

func (m *MockRepository) CreateTopic(ctx context.Context, topic *user.Topic) error {
	if m.CreateTopicFunc != nil {
		return m.CreateTopicFunc(ctx, topic)
	}
	return ErrTest
}

func (m *MockRepository) UpdateTopic(ctx context.Context, topic *user.Topic) error {
	if m.UpdateTopicFunc != nil {
		return m.UpdateTopicFunc(ctx, topic)
	}
	return ErrTest
}

func (m *MockRepository) DeleteTopic(ctx context.Context, userID string, topicID int) error {
	if m.DeleteTopicFunc != nil {
		return m.DeleteTopicFunc(ctx, userID, topicID)
	}
	return ErrTest
}

func (m *MockRepository) GetTopicByID(ctx context.Context, topicID int) (*user.Topic, error) {
	if m.GetTopicByIDFunc != nil {
		return m.GetTopicByIDFunc(ctx, topicID)
	}
	return nil, ErrTest
}

func (m *MockRepository) GetAllTopics(ctx context.Context, page, size int, orderBy, filter string) ([]user.Topic, error) {
	if m.GetAllTopicsFunc != nil {
		return m.GetAllTopicsFunc(ctx, page, size, orderBy, filter)
	}
	return nil, ErrTest
}

func (m *MockRepository) GetTotalTopicsCount(ctx context.Context, filter string) (int, error) {
	if m.GetTotalTopicsCountFunc != nil {
		return m.GetTotalTopicsCountFunc(ctx, filter)
	}
	return 0, ErrTest
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
	GetSessionFunc                  func(sessionID string) (*session.Session, error)
	CreateSessionFunc               func(userID string) (*session.Session, error)
	DeleteSessionFunc               func(sessionID string) error
	NewSessionCookieFunc            func(token string) *http.Cookie
	GetUserFromSessionFunc          func(sessionID string) (*user.User, error)
	GetSessionFromSessionTokensFunc func(sessionToken, refreshToken string) (*session.Session, error)
	DeleteSessionWhenNewCreatedFunc func(ctx context.Context, sessionID string, userID string) error
}

func (m *MockSessionManager) GetSession(sessionID string) (*session.Session, error) {
	if m.GetSessionFunc != nil {
		return m.GetSessionFunc(sessionID)
	}
	return nil, ErrTest
}

func (m *MockSessionManager) GetUserFromSession(sessionID string) (*user.User, error) {
	if m.GetUserFromSessionFunc != nil {
		return m.GetUserFromSessionFunc(sessionID)
	}
	return nil, ErrTest
}

func (m *MockSessionManager) CreateSession(_ context.Context, userID string) (*session.Session, error) {
	if m.CreateSessionFunc != nil {
		return m.CreateSessionFunc(userID)
	}
	return nil, ErrTest
}

func (m *MockSessionManager) ValidateSession(sessionID string) error {
	if m.GetSessionFunc != nil {
		_, err := m.GetSessionFunc(sessionID)
		if err != nil {
			return err
		}
		return nil
	}
	return ErrTest
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

func (m *MockSessionManager) GetSessionFromSessionTokens(sessionToken, refreshToken string) (*session.Session, error) {
	if m.GetSessionFromSessionTokensFunc != nil {
		return m.GetSessionFromSessionTokensFunc(sessionToken, refreshToken)
	}
	return nil, ErrTest
}

func (m *MockSessionManager) DeleteSessionWhenNewCreated(ctx context.Context, sessionID string, userID string) error {
	if m.DeleteSessionWhenNewCreatedFunc != nil {
		return m.DeleteSessionWhenNewCreatedFunc(ctx, sessionID, userID)
	}
	return ErrTest
}
