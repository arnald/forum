package topic

import "github.com/arnald/forum/internal/domain/comment"

type Topic struct {
	UserID        string
	OwnerUsername string
	Title         string
	Content       string
	ImagePath     string
	CreatedAt     string
	UpdatedAt     string
	Comments      []comment.Comment
	ID            int
	CategoryID    int
}
