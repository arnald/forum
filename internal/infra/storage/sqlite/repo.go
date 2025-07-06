package sqlite

import (
	"context"

	"github.com/arnald/forum/internal/domain/user"
)

type Repo struct {
	// users map[string]user.User
}

func NewRepo() Repo {
	return Repo{
		// users: make(map[string]user.User),
	}
}

// TODO: retrieves all users from the repository.
func (r Repo) GetAll(_ context.Context) ([]user.User, error) {
	return nil, nil
}
