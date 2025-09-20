package category

import (
	"time"
)

type Category struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
	PostCount   int // Number of posts in this category
}

type CategoryWithPosts struct {
	Category *Category
	Posts    []CategoryPost
}

type CategoryPost struct {
	ID        string
	Title     string
	UserID    string
	Username  string
	CreatedAt time.Time
}