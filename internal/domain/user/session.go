package user

import "time"

type Session struct {
	Expiry       time.Time
	UserID       string `json:"userId"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken,omitzero"`
}
