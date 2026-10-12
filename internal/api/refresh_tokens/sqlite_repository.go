package refresh_tokens

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aaksant/school-api/internal/db"
)

// created_at in database is datetime('now')
// datetime('now') stored as YYYY-MM-DD HH:MM:S layout.
// revoked_at is set by CURRENT_TIMESTAMP
// CURRENT_TIMESTAMP has the same layout as datetime('now')
//
// only expires_at is manual
const sqliteTimeLayout = "2006-01-02 15:04:05"

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

func (repo *SqliteRepository) Create(
	ctx context.Context,
	rt RefreshToken,
) (RefreshToken, error) {
	return create(ctx, repo.queries, rt)
}

func (repo *SqliteRepository) GetByHash(
	ctx context.Context,
	hash string,
) (RefreshToken, error) {
	raw, err := repo.queries.GetRefreshTokenByHash(ctx, hash)
	if err != nil {
		return RefreshToken{}, ErrTokenNotFound
	}

	refreshToken, err := toRefreshToken(raw)
	if err != nil {
		return RefreshToken{}, err
	}
	return refreshToken, nil
}

func (repo *SqliteRepository) Revoke(ctx context.Context, id int64) error {
	return revoke(ctx, repo.queries, id)
}

func (repo *SqliteRepository) RevokeAllForUser(
	ctx context.Context,
	userId int64,
) error {
	err := repo.queries.RevokeAllRefreshTokensForUser(ctx, userId)
	if err != nil {
		return fmt.Errorf("revoke all refresh tokens for user %d: %w", userId, err)
	}
	return nil
}

func (repo *SqliteRepository) Rotate(
	ctx context.Context,
	oldId int64,
	next RefreshToken,
) (RefreshToken, error) {
	// use transaction because a query might fail
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return RefreshToken{}, fmt.Errorf("begin rotate tx: %w", err)
	}
	defer tx.Rollback()

	// use only q and never repo.queries.
	// call through repo.queries runs on a different pooled connection,
	// outside the transaction
	q := repo.queries.WithTx(tx)

	if err := revoke(ctx, q, oldId); err != nil {
		return RefreshToken{}, err
	}

	new, err := create(ctx, q, next)
	if err != nil {
		return RefreshToken{}, err
	}

	if err := tx.Commit(); err != nil {
		return RefreshToken{}, fmt.Errorf("commit rotate tx: %w", err)
	}
	return new, nil
}

func create(
	ctx context.Context,
	q *db.Queries,
	rt RefreshToken,
) (RefreshToken, error) {
	created, err := q.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		UserID:    rt.UserId,
		TokenHash: rt.TokenHash,
		ExpiresAt: rt.ExpiresAt.UTC().Format(sqliteTimeLayout),
	})
	if err != nil {
		return RefreshToken{}, err
	}

	refreshToken, err := toRefreshToken(created)
	if err != nil {
		return RefreshToken{}, err
	}
	return refreshToken, nil
}

func revoke(
	ctx context.Context,
	q *db.Queries,
	id int64,
) error {
	rowsAffected, err := q.RevokeRefreshToken(ctx, id)
	if err != nil {
		return fmt.Errorf("revoke refresh token %d: %w", id, err)
	}
	if rowsAffected == 0 {
		return ErrTokenAlreadyRevoked
	}
	return nil
}

func toRefreshToken(rt db.RefreshToken) (RefreshToken, error) {
	parsedExpiresAt, err := time.Parse(sqliteTimeLayout, rt.ExpiresAt)
	if err != nil {
		return RefreshToken{}, fmt.Errorf("error parsing exp: %w", err)
	}

	parsedCreatedAt, err := time.Parse(sqliteTimeLayout, rt.CreatedAt)
	if err != nil {
		return RefreshToken{}, fmt.Errorf("error parsing iat: %w", err)
	}

	// RevokedAt is nullable, more guard
	var parsedRevokedAt *time.Time
	if rt.RevokedAt != nil {
		t, err := time.Parse(sqliteTimeLayout, *rt.RevokedAt)
		if err != nil {
			return RefreshToken{}, fmt.Errorf(
				"parse revoked_at %v: %w",
				*rt.RevokedAt, err,
			)
		}
		parsedRevokedAt = &t
	}

	return RefreshToken{
		Id:        rt.ID,
		UserId:    rt.UserID,
		TokenHash: rt.TokenHash,
		ExpiresAt: parsedExpiresAt,
		CreatedAt: parsedCreatedAt,
		RevokedAt: parsedRevokedAt,
	}, nil
}
