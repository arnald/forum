package domain

import "time"

// Votes Model
type Vote struct {
	ID         int64
	UserID     []byte
	PostID     *int64
	CommentID  *int64
	Value      int8
	Created_at time.Time
}
