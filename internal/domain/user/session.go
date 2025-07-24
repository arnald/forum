package user

import "time"

type Session struct {
	Expiry time.Time
	UserID string
	Token  string
}
