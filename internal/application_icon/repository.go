package application_icon

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var ErrApplicationIconNotFound = errors.New("application icon not found")

type Repository interface {
	Create(ctx context.Context, icon *ApplicationIcon) error
	FindAll(ctx context.Context) ([]ApplicationIcon, error)
	FindByID(ctx context.Context, id uuid.UUID) (*ApplicationIcon, error)
	Update(ctx context.Context, id uuid.UUID, icon *ApplicationIcon) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindFileByID(ctx context.Context, id uuid.UUID) (*ApplicationIcon, error)
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
	icon *ApplicationIcon,
) error {
	const query = `
		INSERT INTO application_icons (
			name,
			file_data,
			mime_type
		)
		VALUES (
			$1,
			$2,
			$3
		)
		RETURNING
			id,
			created_at,
			updated_at;
	`

	return r.db.QueryRowxContext(
		ctx,
		query,
		icon.Name,
		icon.FileData,
		icon.MimeType,
	).Scan(
		&icon.ID,
		&icon.CreatedAt,
		&icon.UpdatedAt,
	)
}

func (r *repository) FindAll(
	ctx context.Context,
) ([]ApplicationIcon, error) {
	const query = `
		SELECT
			id,
			name,
			mime_type,
			created_at,
			updated_at
		FROM application_icons
		ORDER BY name ASC;
	`

	var icons []ApplicationIcon

	err := r.db.SelectContext(
		ctx,
		&icons,
		query,
	)
	if err != nil {
		return nil, err
	}

	return icons, nil
}

func (r *repository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*ApplicationIcon, error) {
	const query = `
		SELECT
			id,
			name,
			mime_type,
			created_at,
			updated_at
		FROM application_icons
		WHERE id = $1
		LIMIT 1;
	`

	var icon ApplicationIcon

	err := r.db.GetContext(
		ctx,
		&icon,
		query,
		id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrApplicationIconNotFound
		}

		return nil, err
	}

	return &icon, nil
}

func (r *repository) FindFileByID(
	ctx context.Context,
	id uuid.UUID,
) (*ApplicationIcon, error) {
	const query = `
		SELECT
			id,
			name,
			file_data,
			mime_type
		FROM application_icons
		WHERE id = $1
		LIMIT 1;
	`

	var icon ApplicationIcon

	err := r.db.GetContext(
		ctx,
		&icon,
		query,
		id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrApplicationIconNotFound
		}

		return nil, err
	}

	return &icon, nil
}

func (r *repository) Update(
	ctx context.Context,
	id uuid.UUID,
	icon *ApplicationIcon,
) error {
	const query = `
		UPDATE application_icons
		SET
			name = $1,
			file_data = $2,
			mime_type = $3,
			updated_at = NOW()
		WHERE id = $4;
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		icon.Name,
		icon.FileData,
		icon.MimeType,
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
		return ErrApplicationIconNotFound
	}

	return nil
}

func (r *repository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	const query = `
		DELETE FROM application_icons
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
		return ErrApplicationIconNotFound
	}

	return nil
}
