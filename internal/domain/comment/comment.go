package comment

import "time"

type Comment struct {
	CreatedAt     time.Time
	UpdatedAt     time.Time
	UserVote      *int
	UserID        string
	Content       string
	OwnerUsername string
	TopicID       int
	ID            int
	UpvoteCount   int
	DownvoteCount int
	VoteScore     int
}
