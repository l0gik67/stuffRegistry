package postgres

import (
	"context"
	"os"
	"stuffRegistry/dto"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	createEmpSql = "INSERT INTO employees (full_name, position) VALUES ($1, $2) RETURNING id, full_name, position"
	getAllSql    = "SELECT id, full_name, position FROM employees"
	deleteSql    = "DELETE FROM employees WHERE id = $1 RETURNING id, full_name, position"
)

type EmployeeService struct {
	pool *pgxpool.Pool
}

type Employee struct {
	ID       int64  `json:"id"`
	FullName string `json:"full_name"`
	Position string `json:"position"`
}

func NewEmployeeService() *EmployeeService {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		panic("DATABASE_URL environment variable not set")
	}
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		panic(err)
	}
	return &EmployeeService{
		pool: pool,
	}
}

func (empService EmployeeService) NewEmp(emp dto.EmployeeDto) (Employee, error) {
	connection, err := empService.pool.Acquire(context.Background())
	if err != nil {
		return Employee{}, err
	}
	defer connection.Release()

	result := connection.QueryRow(
		context.Background(),
		createEmpSql,
		emp.FullName,
		emp.Position)

	employee := Employee{}
	err = result.Scan(&employee.ID, &employee.FullName, &employee.Position)
	if err != nil {
		return Employee{}, err
	}
	return employee, nil
}

func (empService EmployeeService) GetEmployees() ([]Employee, error) {
	connection, err := empService.pool.Acquire(context.Background())
	if err != nil {
		return []Employee{}, err
	}
	defer connection.Release()

	rows, err := connection.Query(context.Background(), getAllSql)
	if err != nil {
		return []Employee{}, err
	}
	defer rows.Close()
	employees := make([]Employee, 0, 100)

	for rows.Next() {
		employee := Employee{}
		err = rows.Scan(
			&employee.ID,
			&employee.FullName,
			&employee.Position)
		if err != nil {
			return []Employee{}, err
		}
		employees = append(employees, employee)
	}
	if err = rows.Err(); err != nil {
		return []Employee{}, err
	}
	return employees, nil
}

func (empService EmployeeService) DeleteEmployee(id int) (Employee, error) {
	connection, err := empService.pool.Acquire(context.Background())
	if err != nil {
		return Employee{}, err
	}
	defer connection.Release()

	res := connection.QueryRow(context.Background(), deleteSql, id)

	employee := Employee{}
	err = res.Scan(
		&employee.ID,
		&employee.FullName,
		&employee.Position)
	if err != nil {
		return Employee{}, err
	}
	return employee, nil
}
