package user

import "github.com/arnald/forum/internal/domain/comments"

type Topic struct {
	UserID        string
	OwnerUsername string
	Title         string
	Content       string
	ImagePath     string
	CreatedAt     string
	UpdatedAt     string
	Comments      []comments.Comment
	ID            int
	CategoryID    int
}
