package students

import (
	"context"
	"errors"
)

var ErrStudentNotFound = errors.New("student not found")

type StudentsRepository interface {
	List(ctx context.Context) ([]Student, error)
	GetById(ctx context.Context, id int) (Student, error)
	Create(ctx context.Context, s Student) (Student, error)
	Delete(ctx context.Context, id int) (Student, error)
	Update(ctx context.Context, s Student) (Student, error)
}
