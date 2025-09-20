package post

import (
	"context"
)

type Repository interface {
	CreatePost(ctx context.Context, post *Post) error
	GetPostByID(ctx context.Context, id string) (*Post, error)
	GetAllPosts(ctx context.Context) ([]*Post, error)
	GetPostsByUserID(ctx context.Context, userID string) ([]*Post, error)
	GetPostsByCategory(ctx context.Context, category string) ([]*Post, error)
	GetPostsByCategoryID(ctx context.Context, categoryID string) ([]*Post, error)
	UpdatePost(ctx context.Context, post *Post) error
	DeletePost(ctx context.Context, id string) error
	GetCategories(ctx context.Context) ([]*Category, error)
}