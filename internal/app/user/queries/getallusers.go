package userqueries

import (
	"context"
	"time"

	"github.com/arnald/forum/internal/domain/user"
)

type GetAllUsersResult struct {
	CreatedAt time.Time
	Name      string
	Role      string
	ID        string
}

type GetAllUsersRequestHandler interface {
	Handle(ctx context.Context) ([]GetAllUsersResult, error)
}

type getAllUsersRequestHandler struct {
	repo user.Repository
}

func NewGetAllUsersRequestHandler(repo user.Repository) GetAllUsersRequestHandler {
	return getAllUsersRequestHandler{repo: repo}
}

func (r getAllUsersRequestHandler) Handle(ctx context.Context) ([]GetAllUsersResult, error) {
	users, err := r.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	results := []GetAllUsersResult{}
	for _, u := range users {
		results = append(results, GetAllUsersResult{
			ID:        u.ID,
			Name:      u.Username,
			Role:      u.Role,
			CreatedAt: u.CreatedAt,
		})
	}

	return results, nil
}
