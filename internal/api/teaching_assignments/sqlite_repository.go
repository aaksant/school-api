package teachingassignments

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

func (repo *SqliteRepository) List(
	ctx context.Context,
	teacherId int,
) ([]TeachingAssignment, error) {
	raw, err := repo.queries.ListTeachingAssignments(ctx, int64(teacherId))
	if err != nil {
		return []TeachingAssignment{}, err
	}

	assignments := make([]TeachingAssignment, 0, len(raw))
	for _, m := range raw {
		assignments = append(assignments, toTeachingAssignment(m))
	}
	return assignments, nil
}

func (repo *SqliteRepository) GetById(
	ctx context.Context,
	id int,
) (TeachingAssignment, error) {
	assignment, err := repo.queries.GetTeachingAssignmentById(ctx, int64(id))
	if err != nil {
		return TeachingAssignment{}, err
	}
	return toTeachingAssignment(assignment), nil
}

func (repo *SqliteRepository) Create(
	ctx context.Context,
	ta TeachingAssignment,
) (TeachingAssignment, error) {
	assignment, err := repo.queries.CreateTeacherAssignment(
		ctx,
		db.CreateTeacherAssignmentParams{
			ClassID:   int64(ta.ClassId),
			SubjectID: int64(ta.SubjectId),
			TeacherID: int64(ta.TeacherId),
		},
	)
	if err != nil {
		return TeachingAssignment{}, err
	}
	return toTeachingAssignment(assignment), nil
}

func (repo *SqliteRepository) Delete(
	ctx context.Context,
	id int,
) (TeachingAssignment, error) {
	removed, err := repo.queries.DeleteTeacherAssignment(ctx, int64(id))
	if err != nil {
		return TeachingAssignment{}, err
	}
	return toTeachingAssignment(removed), nil
}

func toTeachingAssignment(ta db.TeachingAssignment) TeachingAssignment {
	return TeachingAssignment{
		Id:        int(ta.ID),
		ClassId:   int(ta.ClassID),
		SubjectId: int(ta.SubjectID),
		TeacherId: int(ta.TeacherID),
	}
}
