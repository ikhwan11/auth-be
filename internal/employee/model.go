package employee

type Employee struct {
	ID          int64  `db:"id"`
	EmployeeNo  string `db:"employee_no"`
	Name        string `db:"name"`
	Position    string `db:"position"`
	Departement string `db:"departement"`
}
