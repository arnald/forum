package user

import (
	"time"
)

type User struct {
	CreatedAt time.Time
	Password  string
	AvatarURL *string
	Username  string
	Email     string
	Role      string
	ID        string
}
