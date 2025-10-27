package topicqueries

import (
	"context"

	"github.com/arnald/forum/internal/domain/topic"
)

type GetAllTopicsRequest struct {
	OrderBy    string  `json:"orderBy"`
	Order      string  `json:"order"`
	Filter     string  `json:"filter"`
	Page       int     `json:"page"`
	Size       int     `json:"size"`
	Offset     int     `json:"offset"`
	CategoryID int     `json:"categoryId"`
	UserID     *string `json:"userID,omitempty"`
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
	count, err := h.repo.GetTotalTopicsCount(ctx, req.Filter, req.CategoryID)
	if err != nil {
		return nil, 0, err
	}

	topics, err := h.repo.GetAllTopics(ctx, req.Page, req.Size, req.CategoryID, req.OrderBy, req.Order, req.Filter, req.UserID)
	if err != nil {
		return nil, 0, err
	}

	return topics, count, nil
}
