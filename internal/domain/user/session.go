package user

import "time"

type Session struct {
	Token     []byte
	UserID    string
	Expiry    time.Time
	IPAddress string
}
