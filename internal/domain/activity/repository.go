package activity

import "context"

type Repository interface {
	GetUserActivity(ctx context.Context, userID string) (*Activity, error)
}
