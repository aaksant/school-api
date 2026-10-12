package refresh_tokens

import (
	"context"
	"errors"
)

var ErrTokenAlreadyRevoked = errors.New("token is no longer valid")
var ErrTokenNotFound = errors.New("token not found")

type RefreshTokensRepository interface {
	Create(ctx context.Context, rt RefreshToken) (RefreshToken, error)
	GetByHash(ctx context.Context, hash string) (RefreshToken, error)
	Revoke(ctx context.Context, id int64) error
	RevokeAllForUser(ctx context.Context, userId int64) error
	Rotate(ctx context.Context, oldId int64, next RefreshToken) (RefreshToken, error)
}
