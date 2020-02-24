package app

import (
	"context"
	"time"

	"github.com/minheq/kedul_server_main/errors"
)

// EmployeeService ...
type EmployeeService struct {
	employeeStore     EmployeeStore
	employeeRoleStore EmployeeRoleStore
}

// NewEmployeeService constructor for AuthService
func NewEmployeeService(employeeStore EmployeeStore, employeeRoleStore EmployeeRoleStore) EmployeeService {
	return EmployeeService{employeeStore: employeeStore, employeeRoleStore: employeeRoleStore}
}

// GetEmployeeByID ...
func (s *EmployeeService) GetEmployeeByID(ctx context.Context, id string) (*Employee, error) {
	const op = "app/employeeService.CreateEmployee"

	employee, err := s.employeeStore.GetEmployeeByID(ctx, id)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to get employee by id")
	}

	return employee, nil
}

// CreateEmployee creates employee
func (s *EmployeeService) CreateEmployee(ctx context.Context, locationID string, name string, employeeRoleID string) (*Employee, error) {
	const op = "app/employeeService.CreateEmployee"

	employee := NewEmployee(locationID, name, employeeRoleID)

	err := s.employeeStore.StoreEmployee(ctx, employee)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to employeeStore employee")
	}

	ownerRole := NewEmployeeRole(employee.ID, name, defaultOwnerRolePermissions)
	adminRole := NewEmployeeRole(employee.ID, name, defaultAdminRolePermissions)
	managerRole := NewEmployeeRole(employee.ID, name, defaultManagerRolePermissions)
	receptionistRole := NewEmployeeRole(employee.ID, name, defaultReceptionistRolePermissions)
	specialistRole := NewEmployeeRole(employee.ID, name, defaultSpecialistRolePermissions)

	err = s.employeeRoleStore.StoreEmployeeRole(ctx, ownerRole)
	if err != nil {
		return nil, errors.Wrap(op, err, "failed to create default owner role")
	}

	err = s.employeeRoleStore.StoreEmployeeRole(ctx, adminRole)
	if err != nil {
		return nil, errors.Wrap(op, err, "failed to create default admin role")
	}

	err = s.employeeRoleStore.StoreEmployeeRole(ctx, managerRole)
	if err != nil {
		return nil, errors.Wrap(op, err, "failed to create default manager role")
	}

	err = s.employeeRoleStore.StoreEmployeeRole(ctx, receptionistRole)
	if err != nil {
		return nil, errors.Wrap(op, err, "failed to create default receptionist role")
	}

	err = s.employeeRoleStore.StoreEmployeeRole(ctx, specialistRole)
	if err != nil {
		return nil, errors.Wrap(op, err, "failed to create default employee role")
	}

	return employee, nil
}

// UpdateEmployee updates employee
func (s *EmployeeService) UpdateEmployee(ctx context.Context, id string, name string, profileImageID string, actor Actor) (*Employee, error) {
	const op = "app/employeeService.UpdateEmployee"

	err := actor.can(ctx, opUpdateEmployee)

	if err != nil {
		return nil, errors.Unauthorized(op, err)
	}

	employee, err := s.employeeStore.GetEmployeeByID(ctx, id)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to get employee by id")
	}

	if employee == nil {
		return nil, errors.NotFound(op)
	}

	employee.UpdatedAt = time.Now()
	employee.Name = name
	employee.ProfileImageID = profileImageID

	err = s.employeeStore.UpdateEmployee(ctx, employee)

	if err != nil {
		return nil, errors.Unexpected(op, err, "failed to update employee")
	}

	return employee, nil
}

// DeleteEmployee updates employee
func (s *EmployeeService) DeleteEmployee(ctx context.Context, id string) (*Employee, error) {
	const op = "app/employeeService.DeleteEmployee"

	employee, err := s.employeeStore.GetEmployeeByID(ctx, id)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to get employee by id")
	}

	if employee == nil {
		return nil, errors.NotFound(op)
	}

	err = s.employeeStore.DeleteEmployee(ctx, employee)

	if err != nil {
		return nil, errors.Unexpected(op, err, "failed to update employee")
	}

	return employee, nil
}