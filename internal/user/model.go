package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `db:"id"`
	EmployeeNo string    `db:"employee_no"`
	Password   string    `db:"password"`
	RoleID     *int      `db:"role_id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
