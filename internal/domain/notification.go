package domain

import "time"

// Notification Model
type Notification struct {
	ID        int64
	UserID    []byte
	Type      string
	SourceID  int64
	IsRead    bool
	CreatedAt time.Time
}
