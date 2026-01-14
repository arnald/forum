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
	CategoryNames  []string
	CategoryColors []string
	Comments       []comment.Comment
	CategoryIDs    []int
	ID             int
	UpvoteCount    int
	DownvoteCount  int
	VoteScore      int
}
