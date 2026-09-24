package classes

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

func (repo *SqliteRepository) List(ctx context.Context) ([]Class, error) {
	raw, err := repo.queries.ListClasses(ctx)
	if err != nil {
		return []Class{}, err
	}

	classes := make([]Class, 0, len(raw))
	for _, c := range raw {
		classes = append(classes, toClass(c))
	}
	return classes, nil
}

func (repo *SqliteRepository) GetById(ctx context.Context, id int) (Class, error) {
	class, err := repo.queries.GetClass(ctx, int64(id))
	if err != nil {
		return Class{}, err
	}
	return toClass(class), nil
}

func (repo *SqliteRepository) Create(ctx context.Context, c Class) (Class, error) {
	class, err := repo.queries.CreateClass(ctx, db.CreateClassParams{
		Name:              c.Name,
		HomeroomTeacherID: toInt64Ptr(c.HomeroomTeacherId),
	})
	if err != nil {
		return Class{}, err
	}
	return toClass(class), nil
}

func (repo *SqliteRepository) Delete(ctx context.Context, id int) (Class, error) {
	removed, err := repo.queries.DeleteClass(ctx, int64(id))
	if err != nil {
		return Class{}, err
	}
	return toClass(removed), nil
}

func (repo *SqliteRepository) Update(ctx context.Context, c Class) (Class, error) {
	updated, err := repo.queries.UpdateClass(ctx, db.UpdateClassParams{
		ID:                int64(c.Id),
		Name:              c.Name,
		HomeroomTeacherID: toInt64Ptr(c.HomeroomTeacherId),
	})
	if err != nil {
		return Class{}, err
	}
	return toClass(updated), nil
}

func toClass(c db.Class) Class {
	return Class{
		Id:                int(c.ID),
		Name:              c.Name,
		HomeroomTeacherId: toIntPtr(c.HomeroomTeacherID),
	}
}

// homeroom_teacher_id is nullable
func toInt64Ptr(i *int) *int64 {
	if i == nil {
		return nil
	}
	v := int64(*i)
	return &v
}

func toIntPtr(i *int64) *int {
	if i == nil {
		return nil
	}
	v := int(*i)
	return &v
}
