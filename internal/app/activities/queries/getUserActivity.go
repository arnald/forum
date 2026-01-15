package activityqueries

import (
	"context"

	"github.com/arnald/forum/internal/domain/activity"
)

type GetUserActivityRequest struct {
	UserID string
}

type GetUserActivityHandler interface {
	Handle(ctx context.Context, req GetUserActivityRequest) (*activity.Activity, error)
}

type getUserActivityHandler struct {
	repo activity.Repository
}

func NewGetUserActivityHandler(repo activity.Repository) GetUserActivityHandler {
	return &getUserActivityHandler{repo: repo}
}

func (h *getUserActivityHandler) Handle(ctx context.Context, req GetUserActivityRequest) (*activity.Activity, error) {
	return h.repo.GetUserActivity(ctx, req.UserID)
}
