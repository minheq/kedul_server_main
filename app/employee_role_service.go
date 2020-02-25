package app

import (
	"context"
	"fmt"
	"time"

	"github.com/minheq/kedul_server_main/errors"
)

// EmployeeRoleService ...
type EmployeeRoleService struct {
	employeeRoleStore EmployeeRoleStore
	employeeStore     EmployeeStore
}

// NewEmployeeRoleService constructor for AuthService
func NewEmployeeRoleService(employeeStore EmployeeStore, employeeRoleStore EmployeeRoleStore) EmployeeRoleService {
	return EmployeeRoleService{employeeStore: employeeStore, employeeRoleStore: employeeRoleStore}
}

// GetEmployeeRoleByID ...
func (s *EmployeeRoleService) GetEmployeeRoleByID(ctx context.Context, id string) (*EmployeeRole, error) {
	const op = "app/employeeRoleService.GetEmployeeRoleByID"

	employeeRole, err := s.employeeRoleStore.GetEmployeeRoleByID(ctx, id)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to get employeeRole by id")
	}

	return employeeRole, nil
}

// CreateEmployeeRole creates employeeRole
func (s *EmployeeRoleService) CreateEmployeeRole(ctx context.Context, locationID string, name string, permissions []Permission, actor Actor) (*EmployeeRole, error) {
	const op = "app/employeeRoleService.CreateEmployeeRole"

	err := actor.can(ctx, opCreateEmployeeRole)

	if err != nil {
		return nil, errors.Unauthorized(op, err)
	}

	employeeRole := NewEmployeeRole(locationID, name, permissions)

	err = s.employeeRoleStore.StoreEmployeeRole(ctx, employeeRole)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to employeeRoleStore employeeRole")
	}

	return employeeRole, nil
}

// UpdateEmployeeRole updates employeeRole
func (s *EmployeeRoleService) UpdateEmployeeRole(ctx context.Context, id string, name string, permissions []Permission, actor Actor) (*EmployeeRole, error) {
	const op = "app/employeeRoleService.UpdateEmployeeRole"

	err := actor.can(ctx, opUpdateEmployeeRole)

	if err != nil {
		return nil, errors.Unauthorized(op, err)
	}

	if name == "owner" {
		return nil, errors.Invalid(op, "cannot update owner role")
	}

	employeeRole, err := s.employeeRoleStore.GetEmployeeRoleByID(ctx, id)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to get employeeRole by id")
	}

	if employeeRole == nil {
		return nil, errors.NotFound(op)
	}

	employeeRole.UpdatedAt = time.Now()
	employeeRole.Name = name
	employeeRole.Permissions = permissions

	err = s.employeeRoleStore.UpdateEmployeeRole(ctx, employeeRole)

	if err != nil {
		return nil, errors.Unexpected(op, err, "failed to update employeeRole")
	}

	return employeeRole, nil
}

// DeleteEmployeeRole updates employeeRole
func (s *EmployeeRoleService) DeleteEmployeeRole(ctx context.Context, id string, actor Actor) (*EmployeeRole, error) {
	const op = "app/employeeRoleService.DeleteEmployeeRole"

	err := actor.can(ctx, opDeleteEmployeeRole)

	if err != nil {
		return nil, errors.Unauthorized(op, err)
	}

	employeeRole, err := s.employeeRoleStore.GetEmployeeRoleByID(ctx, id)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to get employee role by id")
	}

	if employeeRole == nil {
		return nil, errors.NotFound(op)
	}

	if employeeRole.Name == "owner" {
		return nil, errors.Invalid(op, "cannot delete owner role")
	}

	employeesWithTheRole, err := s.employeeStore.GetEmployeesByEmployeeRoleID(ctx, employeeRole.ID)

	if err != nil {
		return nil, errors.Unexpected(op, err, "failed to get employees by employee role id")
	}

	if len(employeesWithTheRole) > 0 {
		return nil, errors.Invalid(op, fmt.Sprintf("employees with role=%s still exist. remove them and restart operation", employeeRole.Name))
	}

	err = s.employeeRoleStore.DeleteEmployeeRole(ctx, employeeRole)

	if err != nil {
		return nil, errors.Unexpected(op, err, "failed to update employeeRole")
	}

	return employeeRole, nil
}
