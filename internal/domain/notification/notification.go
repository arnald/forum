package notification

import "time"

type Type string

const (
	NotificationTypeReply   Type = "reply"
	NotificationTypeMention Type = "mention"
	NotificationTypeLike    Type = "like"
)

type Notification struct {
	CreatedAt   time.Time `json:"createdAt"`
	UserID      string    `json:"userId"`
	ActorID     string    `json:"actorId"`
	Type        Type      `json:"type"`
	Title       string    `json:"title"`
	Message     string    `json:"message"`
	RelatedType string    `json:"relatedType,omitempty"`
	RelatedID   string    `json:"relatedId,omitempty"`
	ID          int       `json:"id"`
	IsRead      bool      `json:"isRead"`
}
