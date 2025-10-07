package session

import "time"

type Session struct {
	Expiry             time.Time `json:"expiry"`
	RefreshTokenExpiry time.Time `json:"refreshTokenExpiry"`
	UserID             string    `json:"userId"`
	AccessToken        string    `json:"accessToken"`
	RefreshToken       string    `json:"refreshToken,omitzero"`
}
