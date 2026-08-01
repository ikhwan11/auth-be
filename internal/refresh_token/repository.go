package refreshtoken

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var ErrRefreshTokenNotFound = errors.New("refresh token not found")

type Repository interface {
	Create(ctx context.Context, token *RefreshToken) error
	FindByToken(ctx context.Context, token string) (*RefreshToken, error)

	Revoke(ctx context.Context, id uuid.UUID) error

	RevokeByUserID(ctx context.Context, userID uuid.UUID) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(
	ctx context.Context,
	token *RefreshToken,
) error {
	const query = `
		INSERT INTO refresh_tokens (
			user_id,
			token,
			expires_at,
			revoked
		)
		VALUES (
			$1,
			$2,
			$3,
			$4
		)
		RETURNING
			id,
			created_at;
	`

	return r.db.QueryRowxContext(
		ctx,
		query,
		token.UserID,
		token.Token,
		token.ExpiresAt,
		token.Revoked,
	).Scan(
		&token.ID,
		&token.CreatedAt,
	)
}

func (r *repository) FindByToken(
	ctx context.Context,
	token string,
) (*RefreshToken, error) {
	const query = `
		SELECT
			id,
			user_id,
			token,
			expires_at,
			revoked,
			created_at
		FROM refresh_tokens
		WHERE token = $1
		LIMIT 1;
	`

	var rt RefreshToken

	err := r.db.GetContext(
		ctx,
		&rt,
		query,
		token,
	)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRefreshTokenNotFound
		}

		return nil, err
	}

	return &rt, nil
}

func (r *repository) Revoke(
	ctx context.Context,
	id uuid.UUID,
) error {
	const query = `
		UPDATE refresh_tokens
		SET revoked = TRUE
		WHERE id = $1;
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		id,
	)

	return err
}

func (r *repository) RevokeByUserID(
	ctx context.Context,
	userID uuid.UUID,
) error {
	const query = `
		UPDATE refresh_tokens
		SET revoked = TRUE
		WHERE user_id = $1
		AND revoked = FALSE;
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		userID,
	)

	return err
}
