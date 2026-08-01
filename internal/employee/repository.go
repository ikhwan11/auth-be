package employee

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

var ErrEmployeeNotFound = errors.New("employee not found")

type Repository interface {
	FindByEmployeeNo(ctx context.Context, employeeNo string) (*Employee, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindByEmployeeNo(ctx context.Context, employeeNo string) (*Employee, error) {
	const query = `
		SELECT
			id,
			employee_no,
			name,
			position,
			departement
		FROM employees
		WHERE employee_no = $1
		LIMIT 1;
	`

	var employee Employee

	err := r.db.GetContext(ctx, &employee, query, employeeNo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEmployeeNotFound
		}

		return nil, err
	}

	return &employee, nil
}
