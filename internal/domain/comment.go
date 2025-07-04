package domain

import "time"

// Comment Model
type Comment struct {
	ID         int64
	UserID     []byte
	PostID     int64
	Content    string
	Created_at time.Time
}
