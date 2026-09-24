package teachingassignments

import (
	"context"
	"errors"
)

var ErrTeachingAssignmentNotFound = errors.New("mapping not found")
var ErrInvalidReference = errors.New("class, subject, or teacher does not exist")
var ErrDuplicate = errors.New("class or subject already has teacher assigned")

type TeachingAssignmentsRepository interface {
	List(ctx context.Context, teacherId int) ([]TeachingAssignment, error)
	GetById(ctx context.Context, id int) (TeachingAssignment, error)
	Create(ctx context.Context, tm TeachingAssignment) (TeachingAssignment, error)
	Delete(ctx context.Context, id int) (TeachingAssignment, error)
}
