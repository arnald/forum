package domain

// Post-Category Model
type PostCategory struct {
	PostID     int64 `db:"post_id"`     // Foreign key to posts(id)
	CategoryID int64 `db:"category_id"` // Foreign key to categories(id)
}
