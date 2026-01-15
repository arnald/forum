package activity

type Activity struct {
	CreatedTopics    []TopicActivity
	LikedTopics      []TopicActivity
	DislikedTopics   []TopicActivity
	LikedComments    []CommentVoteActivity
	DislikedComments []CommentVoteActivity
	UserComments     []CommentActivity
}

type TopicActivity struct {
	CreatedAt string
	Title     string
	ID        int
}

type CommentActivity struct {
	CreatedAt  string
	Content    string
	TopicTitle string
	ID         int
	TopicID    int
}

type CommentVoteActivity struct {
	CreatedAt  string
	TopicTitle string
	CommentID  int
	TopicID    int
}
