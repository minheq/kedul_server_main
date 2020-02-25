package app

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/minheq/kedul_server_main/errors"
)

// Employee ...
type Employee struct {
	ID             string    `json:"id"`
	LocationID     string    `json:"location_id"`
	Name           string    `json:"name"`
	UserID         string    `json:"user_id"`
	ProfileImageID string    `json:"profile_image_id"`
	EmployeeRoleID string    `json:"employee_role_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

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

// CreateEmployeeInput ...
type CreateEmployeeInput struct {
	LocationID     string `json:"location_id"`
	Name           string `json:"name"`
	ProfileImageID string `json:"profile_image_id"`
}

// CreateEmployee creates employee
func (s *EmployeeService) CreateEmployee(ctx context.Context, input *CreateEmployeeInput, actor Actor) (*Employee, error) {
	const op = "app/employeeService.CreateEmployee"

	err := actor.can(ctx, opCreateEmployee)

	if err != nil {
		return nil, errors.Unauthorized(op, err)
	}

	now := time.Now()
	id := uuid.Must(uuid.New(), nil).String()

	if input.Name == "" {
		return nil, errors.Invalid(op, "name field required")
	}

	employee := &Employee{
		ID:             id,
		LocationID:     input.LocationID,
		ProfileImageID: input.ProfileImageID,
		Name:           input.Name,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	err = s.employeeStore.StoreEmployee(ctx, employee)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to employeeStore employee")
	}

	return employee, nil
}

// UpdateEmployeeInput ...
type UpdateEmployeeInput struct {
	Name           string `json:"name"`
	ProfileImageID string `json:"profile_image_id"`
}

// UpdateEmployee updates employee
func (s *EmployeeService) UpdateEmployee(ctx context.Context, id string, input *UpdateEmployeeInput, actor Actor) (*Employee, error) {
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

	if input.Name != "" {
		employee.Name = input.Name
	}
	if input.ProfileImageID != "" {
		employee.ProfileImageID = input.ProfileImageID
	}

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
