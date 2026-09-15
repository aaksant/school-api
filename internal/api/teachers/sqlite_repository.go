package teachers

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

func (repo *SqliteRepository) List(ctx context.Context) ([]Teacher, error) {
	dbTeachers, err := repo.queries.ListTeachers(context.Background())
	if err != nil {
		return []Teacher{}, err
	}

	teachers := make([]Teacher, 0, len(dbTeachers))
	for _, t := range dbTeachers {
		teachers = append(teachers, toTeacher(t))
	}
	return teachers, nil
}

func (repo *SqliteRepository) GetById(ctx context.Context, id int) (Teacher, error) {
	teacher, err := repo.queries.GetTeacher(ctx, int64(id))
	if err != nil {
		return Teacher{}, err
	}
	return toTeacher(teacher), nil
}

func (repo *SqliteRepository) Create(ctx context.Context, t Teacher) (Teacher, error) {
	teacher, err := repo.queries.CreateTeacher(ctx, db.CreateTeacherParams{
		FirstName:   t.FirstName,
		LastName:    t.LastName,
		Email:       t.Email,
		DateOfBirth: t.DateOfBirth,
	})
	if err != nil {
		return Teacher{}, err
	}
	return toTeacher(teacher), nil
}

func (repo *SqliteRepository) Delete(ctx context.Context, id int) (Teacher, error) {
	removed, err := repo.queries.DeleteTeacher(ctx, int64(id))
	if err != nil {
		return Teacher{}, err
	}
	return toTeacher(removed), nil
}

func (repo *SqliteRepository) Update(ctx context.Context, t Teacher) (Teacher, error) {
	updated, err := repo.queries.UpdateTeacher(ctx, db.UpdateTeacherParams{
		ID:          int64(t.Id),
		FirstName:   t.FirstName,
		LastName:    t.LastName,
		Email:       t.Email,
		DateOfBirth: t.DateOfBirth,
	})
	if err != nil {
		return Teacher{}, err
	}
	return toTeacher(updated), nil
}

func toTeacher(t db.Teacher) Teacher {
	return Teacher{
		Id:          int(t.ID),
		FirstName:   t.FirstName,
		LastName:    t.LastName,
		Email:       t.Email,
		DateOfBirth: t.DateOfBirth,
	}
}
