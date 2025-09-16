package topicqueries

import (
	"context"

	"github.com/arnald/forum/internal/domain/user"
)

type GetAllTopicsRequest struct {
	OrderBy string `json:"orderBy"`
	Filter  string `json:"filter"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}

type GetAllTopicsRequestHandler interface {
	Handle(ctx context.Context, req GetAllTopicsRequest) ([]user.Topic, error)
}

type getAllTopicsRequestHandler struct {
	repo user.Repository
}

func NewGetAllTopicsHandler(repo user.Repository) GetAllTopicsRequestHandler {
	return getAllTopicsRequestHandler{
		repo: repo,
	}
}

func (h getAllTopicsRequestHandler) Handle(ctx context.Context, req GetAllTopicsRequest) ([]user.Topic, error) {
	topics, err := h.repo.GetAllTopics(ctx, req.Page, req.Size, req.OrderBy, req.Filter)
	if err != nil {
		return nil, err
	}
	return topics, nil
}
