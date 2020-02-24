package app

import (
	"context"
	"time"

	"github.com/minheq/kedul_server_main/errors"
)

// EmployeeService ...
type EmployeeService struct {
	employeeStore EmployeeStore
}

// NewEmployeeService constructor for AuthService
func NewEmployeeService(employeeStore EmployeeStore) EmployeeService {
	return EmployeeService{employeeStore: employeeStore}
}

// GetEmployeeByID ...
func (s *EmployeeService) GetEmployeeByID(ctx context.Context, id string, actor Actor) (*Employee, error) {
	const op = "app/employeeService.GetEmployeeByID"

	err := actor.can(ctx, opReadEmployee)

	if err != nil {
		return nil, errors.Unauthorized(op, err)
	}

	employee, err := s.employeeStore.GetEmployeeByID(ctx, id)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to get employee by id")
	}

	return employee, nil
}

// CreateEmployee creates employee
func (s *EmployeeService) CreateEmployee(ctx context.Context, locationID string, name string, employeeRoleID string, actor Actor) (*Employee, error) {
	const op = "app/employeeService.CreateEmployee"

	err := actor.can(ctx, opCreateEmployee)

	if err != nil {
		return nil, errors.Unauthorized(op, err)
	}

	employee := NewEmployee(locationID, name, employeeRoleID)

	err = s.employeeStore.StoreEmployee(ctx, employee)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to employeeStore employee")
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
func (s *EmployeeService) DeleteEmployee(ctx context.Context, id string, actor Actor) (*Employee, error) {
	const op = "app/employeeService.DeleteEmployee"

	err := actor.can(ctx, opDeleteEmployee)

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

	err = s.employeeStore.DeleteEmployee(ctx, employee)

	if err != nil {
		return nil, errors.Unexpected(op, err, "failed to update employee")
	}

	return employee, nil
}
