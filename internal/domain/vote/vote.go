package vote

import "time"

type Vote struct {
	CreatedAt    time.Time
	CommentID    *int
	UserID       string
	ID           int
	TopicID      int
	ReactionType int
}

type Target struct {
	TopicID   *int
	CommentID *int
}

type Counts struct {
	Upvotes   int
	DownVotes int
	Score     int
}
