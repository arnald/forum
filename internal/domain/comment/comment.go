package comment

import (
	"time"
)

type Comment struct {
	ID           string
	Content      string
	PostID       string
	UserID       string
	Username     string
	ParentID     *string // For nested comments/replies
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Level        int // Nesting level (0 = top-level, 1 = reply, etc.)
	LikeCount    int
	DislikeCount int
	VoteScore    int // LikeCount - DislikeCount
}

type CommentTree struct {
	Comment  *Comment
	Replies  []*CommentTree
}