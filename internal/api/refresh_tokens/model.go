package refresh_tokens

import "time"

// no CRUDs here

type RefreshToken struct {
	Id        int64
	UserId    int64
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}
