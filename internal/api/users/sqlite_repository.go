package users

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aaksant/school-api/internal/auth"
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

func (repo *SqliteRepository) GetById(ctx context.Context, id int64) (User, error) {
	raw, err := repo.queries.GetById(ctx, id)
	if err != nil {
		return User{}, ErrUserNotFound
	}

	user, err := toUser(raw)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (repo *SqliteRepository) GetByEmail(
	ctx context.Context,
	email string,
) (User, error) {
	raw, err := repo.queries.GetByEmail(ctx, email)
	if err != nil {
		return User{}, ErrUserNotFound
	}

	user, err := toUser(raw)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (repo *SqliteRepository) Create(ctx context.Context, u User) (User, error) {
	params := db.CreateUserParams{
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Role:         string(u.Role),
	}
	switch u.Role {
	case auth.RoleStudent:
		params.StudentID = &u.PersonID
	case auth.RoleTeacher:
		params.TeacherID = &u.PersonID
	case auth.RoleStaff:
		params.StaffID = &u.PersonID
	case auth.RoleExecutive:
		params.ExecutiveID = &u.PersonID
	default:
		return User{}, fmt.Errorf("unsupported role %q", u.Role)
	}

	created, err := repo.queries.CreateUser(ctx, params)
	if err != nil {
		return User{}, err
	}

	user, err := toUser(created)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func toUser(u db.User) (User, error) {
	role := auth.Role(u.Role)
	var personId *int64

	switch role {
	case auth.RoleStudent:
		personId = u.StudentID
	case auth.RoleTeacher:
		personId = u.TeacherID
	case auth.RoleStaff:
		personId = u.StaffID
	case auth.RoleExecutive:
		personId = u.ExecutiveID
	}
	if personId == nil {
		return User{}, fmt.Errorf(
			"user %d: role %q has no matching person id", u.ID, u.Role,
		)
	}

	return User{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Role:         role,
		PersonID:     *personId,
		IsActive:     u.IsActive != 0,
	}, nil
}
