package comment

type Comment struct {
	CreatedAt     string
	UpdatedAt     string
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
