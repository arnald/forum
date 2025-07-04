package domain

import "time"

// Session Model
type Session struct {
	Token     []byte
	UserID    []byte
	Expiry    time.Time
	UserAgent *string
	IPAdress  *string
}
