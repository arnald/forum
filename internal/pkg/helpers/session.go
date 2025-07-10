package helpers

import (
	"time"

	"github.com/arnald/forum/internal/domain/user"
	"github.com/google/uuid"
)

func NewSession(tokenUUID uuid.UUID, userUUID uuid.UUID, expiry time.Time, IPAddress string) user.Session {
	session := user.Session{
		Token:     tokenUUID[:],
		UserID:    userUUID.String(),
		Expiry:    expiry,
		IPAddress: IPAddress,
	}
	return session
}
