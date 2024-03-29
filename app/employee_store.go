package app

import (
	"context"
	"database/sql"

	"github.com/minheq/kedul_server_main/errors"
)

// EmployeeStore ...
type EmployeeStore interface {
	GetEmployeesByUserID(ctx context.Context, userID string) ([]*Employee, error)
	GetEmployeesByEmployeeRoleID(ctx context.Context, employeeRoleID string) ([]*Employee, error)
	GetEmployeeByUserIDAndLocationID(ctx context.Context, userID string, locationID string) (*Employee, error)
	GetEmployeeByID(ctx context.Context, id string) (*Employee, error)
	StoreEmployee(ctx context.Context, employee *Employee) error
	UpdateEmployee(ctx context.Context, employee *Employee) error
	DeleteEmployee(ctx context.Context, employee *Employee) error
}

type employeeStore struct {
	db *sql.DB
}

// NewEmployeeStore ...
func NewEmployeeStore(db *sql.DB) EmployeeStore {
	return &employeeStore{db: db}
}

// GetEmployeeByUserIDAndLocationID gets Employee by UserID and LocationID
func (s *employeeStore) GetEmployeesByUserID(ctx context.Context, userID string) ([]*Employee, error) {
	const op = "app/employeeStore.GetEmployeesByUserID"

	query := `
		SELECT id, location_id, name, user_id, employee_role_id, profile_image_id, created_at, updated_at
		FROM employee
		WHERE user_id=$1;
	`
	employees := make([]*Employee, 0)

	rows, err := s.db.Query(query, userID)

	if err != nil {
		return nil, errors.Wrap(op, err, "database error")
	}

	for rows.Next() {
		employee := &Employee{}

		err := rows.Scan(&employee.ID, &employee.LocationID, &employee.Name, &employee.UserID, &employee.EmployeeRoleID, &employee.ProfileImageID, &employee.CreatedAt, &employee.UpdatedAt)

		if err != nil {
			return nil, errors.Wrap(op, err, "row scan error")
		}

		employees = append(employees, employee)
	}

	err = rows.Err()

	if err != nil {
		return nil, errors.Wrap(op, err, "database error")
	}

	return employees, nil
}

// GetEmployeeByUserIDAndLocationID gets Employee by UserID and LocationID
func (s *employeeStore) GetEmployeesByEmployeeRoleID(ctx context.Context, employeeRoleID string) ([]*Employee, error) {
	const op = "app/employeeStore.GetEmployeesByEmployeeRoleID"

	query := `
		SELECT id, location_id, name, user_id, employee_role_id, profile_image_id, created_at, updated_at
		FROM employee
		WHERE employee_role_id=$1;
	`
	employees := make([]*Employee, 0)

	rows, err := s.db.Query(query, employeeRoleID)

	if err != nil {
		return nil, errors.Wrap(op, err, "database error")
	}

	for rows.Next() {
		employee := &Employee{}

		err = rows.Scan(&employee.ID, &employee.LocationID, &employee.Name, &employee.UserID, &employee.EmployeeRoleID, &employee.ProfileImageID, &employee.CreatedAt, &employee.UpdatedAt)

		if err != nil {
			return nil, errors.Wrap(op, err, "row scan error")
		}

		employees = append(employees, employee)
	}

	err = rows.Err()

	if err != nil {
		return nil, errors.Wrap(op, err, "database error")
	}

	return employees, nil
}

// GetEmployeeByUserIDAndLocationID gets Employee by UserID and LocationID
func (s *employeeStore) GetEmployeeByUserIDAndLocationID(ctx context.Context, userID string, locationID string) (*Employee, error) {
	const op = "app/employeeStore.GetEmployeeByID"

	query := `
		SELECT id, location_id, name, user_id, employee_role_id, profile_image_id, created_at, updated_at
		FROM employee
		WHERE user_id=$1
			AND location_id=$2;
	`

	employee := &Employee{}

	row := s.db.QueryRow(query, userID, locationID)

	if row == nil {
		return nil, nil
	}

	err := row.Scan(&employee.ID, &employee.LocationID, &employee.Name, &employee.UserID, &employee.EmployeeRoleID, &employee.ProfileImageID, &employee.CreatedAt, &employee.UpdatedAt)

	if err != nil {
		return nil, errors.Wrap(op, err, "database error")
	}

	return employee, nil
}

// GetEmployeeByID gets Employee by ID
func (s *employeeStore) GetEmployeeByID(ctx context.Context, id string) (*Employee, error) {
	const op = "app/employeeStore.GetEmployeeByID"

	query := `
		SELECT id, location_id, name, user_id, employee_role_id, profile_image_id, created_at, updated_at
		FROM employee
		WHERE id=$1;
	`

	employee := &Employee{}

	row := s.db.QueryRow(query, id)

	if row == nil {
		return nil, nil
	}

	err := row.Scan(&employee.ID, &employee.LocationID, &employee.Name, &employee.UserID, &employee.EmployeeRoleID, &employee.ProfileImageID, &employee.CreatedAt, &employee.UpdatedAt)

	if err != nil {
		return nil, errors.Wrap(op, err, "database error")
	}

	return employee, nil
}

// StoreEmployee persists Employee
func (s *employeeStore) StoreEmployee(ctx context.Context, employee *Employee) error {
	const op = "app/employeeStore.StoreEmployee"

	query := `
		INSERT INTO employee (id, location_id, user_id, name, employee_role_id, profile_image_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := s.db.Exec(query, employee.ID, employee.LocationID, employee.UserID, employee.Name, employee.EmployeeRoleID, employee.ProfileImageID, employee.CreatedAt, employee.UpdatedAt)

	if err != nil {
		return errors.Wrap(op, err, "database error")
	}

	return nil
}

// UpdateEmployee updates Employee including all fields
func (s *employeeStore) UpdateEmployee(ctx context.Context, employee *Employee) error {
	const op = "app/employeeStore.UpdateEmployee"

	query := `
		UPDATE employee
		SET user_id=$2, name=$3, employee_role_id=$4, profile_image_id=$5, updated_at=$6
		WHERE id=$1;
	`

	_, err := s.db.Exec(query, employee.ID, employee.UserID, employee.Name, employee.EmployeeRoleID, employee.ProfileImageID, employee.UpdatedAt)

	if err != nil {
		return errors.Wrap(op, err, "database error")
	}

	return nil
}

// DeleteEmployee deletes Employee
func (s *employeeStore) DeleteEmployee(ctx context.Context, employee *Employee) error {
	const op = "app/employeeStore.DeleteEmployee"

	query := `
		DELETE FROM employee
		WHERE id=$1;
	`

	_, err := s.db.Exec(query, employee.ID)

	if err != nil {
		return errors.Wrap(op, err, "database error")
	}

	return nil
}
