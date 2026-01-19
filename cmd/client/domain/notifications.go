package domain

import "time"

type NotificationType string

const (
	NotificationTypeReply   NotificationType = "reply"
	NotificationTypeMention NotificationType = "mention"
	NotificationTypeLike    NotificationType = "like"
)

type Notification struct {
	CreatedAt   time.Time        `json:"createdAt"`
	UserID      string           `json:"userId"`
	ActorID     string           `json:"actorId"`
	Type        NotificationType `json:"type"`
	Title       string           `json:"title"`
	Message     string           `json:"message"`
	RelatedType string           `json:"relatedType,omitempty"`
	RelatedID   string           `json:"relatedId,omitempty"`
	ID          int              `json:"id"`
	IsRead      bool             `json:"isRead"`
}

type UnreadCountResponse struct {
	Count int `json:"count"`
}
