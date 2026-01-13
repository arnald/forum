package topic

import "github.com/arnald/forum/internal/domain/comment"

type Topic struct {
	UserVote       *int
	UpdatedAt      string
	Title          string
	Content        string
	ImagePath      string
	CreatedAt      string
	UserID         string
	OwnerUsername  string
	CategoryName   string
	CategoryColor  string
	CategoryNames  []string
	CategoryColors []string
	Comments       []comment.Comment
	ID             int
	CategoryID     int
	CategoryIDs    []int
	UpvoteCount    int
	DownvoteCount  int
	VoteScore      int
}
