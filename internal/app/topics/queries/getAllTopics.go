package topicqueries

import (
	"context"

	"github.com/arnald/forum/internal/domain/topic"
)

type GetAllTopicsRequest struct {
	OrderBy string `json:"orderBy"`
	Filter  string `json:"filter"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}

type GetAllTopicsRequestHandler interface {
	Handle(ctx context.Context, req GetAllTopicsRequest) ([]topic.Topic, int, error)
}

type getAllTopicsRequestHandler struct {
	repo topic.Repository
}

func NewGetAllTopicsHandler(repo topic.Repository) GetAllTopicsRequestHandler {
	return getAllTopicsRequestHandler{
		repo: repo,
	}
}

func (h getAllTopicsRequestHandler) Handle(ctx context.Context, req GetAllTopicsRequest) ([]topic.Topic, int, error) {
	topics, err := h.repo.GetAllTopics(ctx, req.Page, req.Size, req.OrderBy, req.Filter)
	if err != nil {
		return nil, 0, err
	}

	count, err := h.repo.GetTotalTopicsCount(ctx, req.Filter)
	if err != nil {
		return nil, 0, err
	}
	return topics, count, nil
}
