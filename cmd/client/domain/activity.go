package domain

// ActivityData represents the data structure for the activity page.
type ActivityData struct {
	User             *LoggedInUser
	CreatedTopics    []ActivityTopic       `json:"createdTopics"`
	LikedTopics      []ActivityTopic       `json:"likedTopics"`
	DislikedTopics   []ActivityTopic       `json:"dislikedTopics"`
	LikedComments    []ActivityCommentVote `json:"likedComments"`
	DislikedComments []ActivityCommentVote `json:"dislikedComments"`
	UserComments     []ActivityComment     `json:"userComments"`
}

// ActivityTopic represents a topic in the activity feed.
type ActivityTopic struct {
	Title     string
	CreatedAt string
	ID        int
}

// ActivityComment represents a comment the user made.
type ActivityComment struct {
	Content    string
	TopicTitle string
	CreatedAt  string
	ID         int
	TopicID    int
}

// ActivityCommentVote represents a comment the user liked/disliked.
type ActivityCommentVote struct {
	TopicTitle string
	CreatedAt  string
	CommentID  int
	TopicID    int
}
