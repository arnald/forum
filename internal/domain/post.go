package domain

import "time"

// Post Model
type Post struct {
	ID           int64
	UserID       []byte
	Title        string
	Content      string
	Image_path   *string
	Created_at   time.Time
	Published    bool
	Published_at *time.Time
}
