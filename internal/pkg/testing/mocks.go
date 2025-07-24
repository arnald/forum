package mocks

import (
	"context"

	"github.com/arnald/forum/internal/domain/user"
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
	return nil, nil
}

func (m *MockRepository) GetAll(ctx context.Context) ([]user.User, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}
	return nil, nil
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
