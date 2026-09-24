package classes

import (
	"context"
	"errors"
)

var ErrClassNotFound = errors.New("class not found")
var ErrInvalidReference = errors.New("homeroom teacher does not exist")
var ErrDuplicate = errors.New("this teacher is already the homeroom teacher of another class")

type ClassesRepository interface {
	List(ctx context.Context) ([]Class, error)
	GetById(ctx context.Context, id int) (Class, error)
	Create(ctx context.Context, c Class) (Class, error)
	Delete(ctx context.Context, id int) (Class, error)
	Update(ctx context.Context, c Class) (Class, error)
}
