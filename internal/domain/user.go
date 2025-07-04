package domain

import "time"

// User Model
type User struct {
	ID         []byte
	Username   string
	Email      string
	Password   *string
	Role       string
	Created_at time.Time
	Avatar_url *string
}
