package teachers

import (
	"context"
	"errors"
)

var ErrTeacherNotFound = errors.New("teacher not found")

type TeachersRepository interface {
	List(ctx context.Context) ([]Teacher, error)
	GetById(ctx context.Context, id int) (Teacher, error)
	Create(ctx context.Context, t Teacher) (Teacher, error)
	Delete(ctx context.Context, id int) (Teacher, error)
	Update(ctx context.Context, t Teacher) (Teacher, error)
}
