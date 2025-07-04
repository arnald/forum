package domain

import "time"

// OAuth Model
type OAuth struct {
	Provider       string
	ProviderUserID string
	UserID         []byte
	Access_token   string
	Refresh_token  *string
	Expiry         *time.Time
}
