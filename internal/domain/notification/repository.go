package notification

import "context"

type Repository interface {
	Create(ctx context.Context, notification *Notification) error
	GetByUserID(ctx context.Context, userID string, limit int) ([]*Notification, error)
	GetUnreadCount(ctx context.Context, userID string) (int, error)
	MarkAsRead(ctx context.Context, notificationID int, userID string) error
	MarkAllAsRead(ctx context.Context, userID string) error
}
