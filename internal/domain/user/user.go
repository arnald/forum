package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	CreatedAt time.Time
	Password  string
	AvatarURL *string
	Username  string
	Email     string
	Role      string
	ID        uuid.UUID
}
