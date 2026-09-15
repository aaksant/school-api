package students

import (
	"context"
	"database/sql"

	"github.com/aaksant/school-api/internal/db"
)

type SqliteRepository struct {
	db      *sql.DB
	queries *db.Queries
}

func NewSqliteRepository(conn *sql.DB) *SqliteRepository {
	return &SqliteRepository{
		db:      conn,
		queries: db.New(conn),
	}
}

func (repo *SqliteRepository) List(ctx context.Context) ([]Student, error) {
	raw, err := repo.queries.ListStudents(ctx)
	if err != nil {
		return []Student{}, err
	}

	students := make([]Student, 0, len(raw))
	for _, s := range raw {
		students = append(students, toStudent(s))
	}
	return students, nil
}

func (repo *SqliteRepository) GetById(ctx context.Context, id int) (Student, error) {
	student, err := repo.queries.GetStudent(ctx, int64(id))
	if err != nil {
		return Student{}, err
	}
	return toStudent(student), nil
}

func (repo *SqliteRepository) Create(ctx context.Context, s Student) (Student, error) {
	student, err := repo.queries.CreateStudent(ctx, db.CreateStudentParams{
		FirstName:   s.FirstName,
		LastName:    s.LastName,
		Email:       s.Email,
		DateOfBirth: s.DateOfBirth,
	})
	if err != nil {
		return Student{}, err
	}
	return toStudent(student), nil
}

func (repo *SqliteRepository) Delete(ctx context.Context, id int) (Student, error) {
	removed, err := repo.queries.DeleteStudent(ctx, int64(id))
	if err != nil {
		return Student{}, err
	}
	return toStudent(removed), nil
}

func (repo *SqliteRepository) Update(ctx context.Context, s Student) (Student, error) {
	updated, err := repo.queries.UpdateStudent(ctx, db.UpdateStudentParams{
		ID:          int64(s.Id),
		FirstName:   s.FirstName,
		LastName:    s.LastName,
		Email:       s.Email,
		DateOfBirth: s.DateOfBirth,
	})
	if err != nil {
		return Student{}, err
	}
	return toStudent(updated), nil
}

/* mappers to transform sqlc models to validated models */

func toStudent(s db.Student) Student {
	return Student{
		Id:        int(s.ID),
		FirstName: s.FirstName,
		LastName:  s.LastName,
		Email:     s.Email,
	}
}
