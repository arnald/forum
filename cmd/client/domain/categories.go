package domain

type CategoryData struct {
	Data struct {
		Categories []Category `json:"categories"`
	} `json:"data"`
}

type Category struct {
	Name        string  `json:"name"`
	Color       string  `json:"color"`
	Slug        string  `json:"slug,omitzero"`
	Description string  `json:"description,omitzero"`
	ImagePath   string  `json:"imagePath"`
	Topics      []Topic `json:"topics,omitzero"`
	ID          int     `json:"id"`
	TopicCount  int     `json:"topicsCount,omitzero"`
}

type Pagination struct {
	Page       int `json:"page"`
	TotalPages int `json:"totalPages"`
	TotalItems int `json:"totalItems"`
	NextPage   int `json:"nextPage"`
	PrevPage   int `json:"prevPage"`
}

type Logo struct {
	URL    string `json:"url"`
	ID     int    `json:"id"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Topic struct {
	UserVote      *int      `json:"userVote,omitempty"`
	UserID        string    `json:"userId"`
	Content       string    `json:"content"`
	ImagePath     string    `json:"imagePath"`
	Title         string    `json:"title"`
	CategoryColor string    `json:"categoryColor"`
	CategoryName  string    `json:"categoryName"`
	CreatedAt     string    `json:"createdAt"`
	UpdatedAt     string    `json:"updatedAt"`
	OwnerUsername string    `json:"ownerUsername"`
	Comments      []Comment `json:"comments"`
	CategoryID    int       `json:"categoryId"`
	VoteScore     int       `json:"voteScore"`
	DownvoteCount int       `json:"downvoteCount"`
	UpvoteCount   int       `json:"upvoteCount"`
	ID            int       `json:"id"`
}

type Comment struct {
	ID            int    `json:"id"`
	TopicID       int    `json:"topicId"`
	UserID        string `json:"userId"`
	Content       string `json:"content"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
	OwnerUsername string `json:"ownerUsername"`
	UpvoteCount   int    `json:"upvoteCount"`
	DownvoteCount int    `json:"downvoteCount"`
	VoteScore     int    `json:"voteScore"`
	UserVote      *int   `json:"userVote,omitempty"`
}
