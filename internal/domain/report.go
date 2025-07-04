package domain

import "time"

// Report Model
type Report struct {
	ID         int64
	ReporterID []byte
	PostID     int64
	CommentID  int64
	Reason     string
	Status     string
	ResolvedBy []byte
	ResolvedAt time.Time
	CreatedAt  time.Time
}
