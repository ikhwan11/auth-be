package application

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, app *Application) error
	FindAll(ctx context.Context) ([]Application, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Application, error)
	Update(ctx context.Context, id uuid.UUID, req UpdateApplicationRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
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
	app *Application,
) error {
	const query = `
		INSERT INTO applications (
			name,
			code,
			url,
			is_default
		)
		VALUES (
			$1,
			$2,
			$3,
			$4
		)
		RETURNING
			id,
			created_at,
			updated_at;
	`

	return r.db.QueryRowxContext(
		ctx,
		query,
		app.Name,
		app.Code,
		app.URL,
		app.IsDefault,
	).Scan(
		&app.ID,
		&app.CreatedAt,
		&app.UpdatedAt,
	)
}

func (r *repository) FindAll(
	ctx context.Context,
) ([]Application, error) {
	const query = `
		SELECT
			id,
			name,
			code,
			url,
			is_default,
			is_active,
			created_at,
			updated_at
		FROM applications
		ORDER BY name ASC;
	`

	var applications []Application

	err := r.db.SelectContext(
		ctx,
		&applications,
		query,
	)
	if err != nil {
		return nil, err
	}

	return applications, nil
}

func (r *repository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*Application, error) {
	const query = `
		SELECT
			id,
			name,
			code,
			url,
			is_default,
			is_active,
			created_at,
			updated_at
		FROM applications
		WHERE id = $1
		LIMIT 1;
	`

	var app Application

	err := r.db.GetContext(
		ctx,
		&app,
		query,
		id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrApplicationNotFound
		}

		return nil, err
	}

	return &app, nil
}

func (r *repository) Update(
	ctx context.Context,
	id uuid.UUID,
	req UpdateApplicationRequest,
) error {
	const query = `
		UPDATE applications
		SET
			name = $1,
			code = $2,
			url = $3,
			is_default = $4,
			is_active = $5,
			updated_at = NOW()
		WHERE id = $6;
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		req.Name,
		req.Code,
		req.URL,
		req.IsDefault,
		req.IsActive,
		id,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrApplicationNotFound
	}

	return nil
}

func (r *repository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	const query = `
		DELETE FROM applications
		WHERE id = $1;
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		id,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrApplicationNotFound
	}

	return nil
}
