package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var ErrUserNotFound = errors.New("user not found")

type Repository interface {
	FindByEmployeeNo(ctx context.Context, employeeNo string) (*User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	Create(ctx context.Context, user *User) error
	UpdatePassword(ctx context.Context, id uuid.UUID, password string) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindByEmployeeNo(ctx context.Context, employeeNo string) (*User, error) {
	const query = `
		SELECT
			id,
			employee_no,
			password,
			role_id,
			created_at,
			updated_at
		FROM users
		WHERE employee_no = $1
		LIMIT 1;
	`

	var user User

	err := r.db.GetContext(ctx, &user, query, employeeNo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	const query = `
		SELECT
			id,
			employee_no,
			password,
			role_id,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
		LIMIT 1;
	`

	var user User

	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *repository) Create(ctx context.Context, user *User) error {
	const query = `
		INSERT INTO users (
			employee_no,
			password,
			role_id
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
		user.EmployeeNo,
		user.Password,
		user.RoleID,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
}

func (r *repository) UpdatePassword(ctx context.Context, id uuid.UUID, password string) error {
	const query = `
		UPDATE users
		SET
			password = $1,
			updated_at = NOW()
		WHERE id = $2;
	`

	_, err := r.db.ExecContext(ctx, query, password, id)
	return err
}
