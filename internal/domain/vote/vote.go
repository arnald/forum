package vote

import "time"

type Vote struct {
	ID           int
	UserID       string
	TopicID      int
	CommentID    *int
	ReactionType int
	CreatedAt    time.Time
}

type VoteTarget struct {
	TopicID   int
	CommentID *int
}

type VoteCounts struct {
	Upvotes   int
	DownVotes int
	Score     int
}
