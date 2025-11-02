package notification

import "time"

type Type string

const (
	NotificationTypeReply   Type = "reply"
	NotificationTypeMention Type = "mention"
	NotificationTypeLike    Type = "like"
)

type Notification struct {
	CreatedAt   time.Time
	RelatedID   string
	UserID      string
	Type        Type
	Title       string
	Message     string
	RelatedType string
	ID          int
	IsRead      bool
}
