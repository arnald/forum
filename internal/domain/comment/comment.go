package comment

import "time"

type Comment struct {
	CreatedAt     time.Time
	UpdatedAt     time.Time
	UserID        string
	Content       string
	OwnerUsername string
	TopicID       int
	ID            int
}
