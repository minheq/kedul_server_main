package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/minheq/kedul_server_main/errors"
)

// EmployeeRoleStore ...
type EmployeeRoleStore interface {
	GetEmployeeRoleByID(ctx context.Context, id string) (*EmployeeRole, error)
	StoreEmployeeRole(ctx context.Context, employeeRole *EmployeeRole) error
	UpdateEmployeeRole(ctx context.Context, employeeRole *EmployeeRole) error
	DeleteEmployeeRole(ctx context.Context, employeeRole *EmployeeRole) error
}

type employeeRoleStore struct {
	db *sql.DB
}

// NewEmployeeRoleStore ...
func NewEmployeeRoleStore(db *sql.DB) EmployeeRoleStore {
	return &employeeRoleStore{db: db}
}

// GetEmployeeRoleByID gets EmployeeRole by ID
func (s *employeeRoleStore) GetEmployeeRoleByID(ctx context.Context, id string) (*EmployeeRole, error) {
	const op = "app/employeeRoleStore.GetEmployeeRoleByPhoneNumber"

	query := `
		SELECT id, location_id, name, permission_ids, created_at, updated_at
		FROM employee_role
		WHERE id=$1;
	`

	employeeRole := &EmployeeRole{}

	row := s.db.QueryRow(query, id)

	if row == nil {
		return nil, nil
	}

	err := row.Scan(&employeeRole.ID, &employeeRole.LocationID, &employeeRole.Name, &employeeRole.PermissionIDs, &employeeRole.CreatedAt, &employeeRole.UpdatedAt)

	if err != nil {
		return nil, errors.Wrap(op, err, "database error")
	}

	permissions := []Permission{}

	for _, permissionID := range employeeRole.PermissionIDs {
		found := false

		for _, permission := range permissionList {
			if permission.ID == id {
				permissions = append(permissions, permission)
				found = true
			}
		}

		if found == false {
			return nil, errors.Unexpected(op, fmt.Errorf("permission for permissionID=%s not found", permissionID), "failed to retrieve permission")
		}
	}

	employeeRole.Permissions = permissions

	return employeeRole, nil
}

// StoreEmployeeRole persists EmployeeRole
func (s *employeeRoleStore) StoreEmployeeRole(ctx context.Context, employeeRole *EmployeeRole) error {
	const op = "app/employeeRoleStore.StoreEmployeeRole"

	query := `
		INSERT INTO employeeRole (id, location_id, name, permission_ids, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := s.db.Exec(query, employeeRole.ID, employeeRole.LocationID, employeeRole.Name, employeeRole.PermissionIDs, employeeRole.CreatedAt, employeeRole.UpdatedAt)

	if err != nil {
		return errors.Wrap(op, err, "database error")
	}

	return nil
}

// UpdateEmployeeRole updates EmployeeRole including all fields
func (s *employeeRoleStore) UpdateEmployeeRole(ctx context.Context, employeeRole *EmployeeRole) error {
	const op = "app/employeeRoleStore.UpdateEmployeeRole"

	query := `
		UPDATE employeeRole
		SET name=$2, permission_ids=$3, updated_at=$4
		WHERE id=$1;
	`

	_, err := s.db.Exec(query, employeeRole.ID, employeeRole.Name, employeeRole.PermissionIDs, employeeRole.UpdatedAt)

	if err != nil {
		return errors.Wrap(op, err, "database error")
	}

	return nil
}

// DeleteEmployeeRole deletes EmployeeRole
func (s *employeeRoleStore) DeleteEmployeeRole(ctx context.Context, employeeRole *EmployeeRole) error {
	const op = "app/employeeRoleStore.DeleteEmployeeRole"

	query := `
		DELETE FROM employee_role
		WHERE id=$1;
	`

	_, err := s.db.Exec(query, employeeRole.ID)

	if err != nil {
		return errors.Wrap(op, err, "database error")
	}

	return nil
}
