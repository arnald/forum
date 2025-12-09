package topicqueries

import (
	"context"

	"github.com/arnald/forum/internal/domain/category"
	"github.com/arnald/forum/internal/domain/topic"
)

type GetAllTopicsRequest struct {
	UserID     *string `json:"userId,omitempty"`
	OrderBy    string  `json:"orderBy"`
	Order      string  `json:"order"`
	Filter     string  `json:"filter"`
	Page       int     `json:"page"`
	Size       int     `json:"size"`
	Offset     int     `json:"offset"`
	CategoryID int     `json:"categoryId"`
}

type GetAllTopicsResponse struct {
	Topics     []topic.Topic
	Categories []category.Category
	Count      int
}
type GetAllTopicsRequestHandler interface {
	Handle(ctx context.Context, req GetAllTopicsRequest) (*GetAllTopicsResponse, error)
}

type getAllTopicsRequestHandler struct {
	topicRepo    topic.Repository
	categoryRepo category.Repository
}

func NewGetAllTopicsHandler(topicRepo topic.Repository, categoryRepo category.Repository) GetAllTopicsRequestHandler {
	return getAllTopicsRequestHandler{
		topicRepo:    topicRepo,
		categoryRepo: categoryRepo,
	}
}

func (h getAllTopicsRequestHandler) Handle(ctx context.Context, req GetAllTopicsRequest) (*GetAllTopicsResponse, error) {
	count, err := h.topicRepo.GetTotalTopicsCount(ctx, req.Filter, req.CategoryID)
	if err != nil {
		return nil, err
	}

	topics, err := h.topicRepo.GetAllTopics(ctx, req.Page, req.Size, req.CategoryID, req.OrderBy, req.Order, req.Filter, req.UserID)
	if err != nil {
		return nil, err
	}

	categories, err := h.categoryRepo.GetAllCategorieNamesAndIDs(ctx)
	if err != nil {
		return nil, err
	}

	response := &GetAllTopicsResponse{
		Topics:     topics,
		Count:      count,
		Categories: categories,
	}

	return response, nil
}
