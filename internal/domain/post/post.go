package post

import (
	"time"
)

type Post struct {
	ID           string
	Title        string
	Content      string
	UserID       string
	Username     string
	Categories   []string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	LikeCount    int
	DislikeCount int
	VoteScore    int // LikeCount - DislikeCount
}

type Category struct {
	ID   string
	Name string
}