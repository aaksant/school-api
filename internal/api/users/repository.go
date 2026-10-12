package users

import (
	"context"
	"errors"
)

var ErrDuplicate = errors.New("user already has an account")
var ErrInvalidReference = errors.New("foregin key violation")
var ErrUserNotFound = errors.New("user not found")

type UsersRepository interface {
	GetById(ctx context.Context, id int64) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	Create(ctx context.Context, u User) (User, error)
}
