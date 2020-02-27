package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minheq/kedul_server_main/errors"
)

// EmployeeRole ...
type EmployeeRole struct {
	ID            string    `json:"id"`
	LocationID    string    `json:"location_id"`
	Name          string    `json:"name"`
	PermissionIDs []string  `json:"permission_ids"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	// These permissions are retrieved in the application code based on PermissionIDs
	Permissions []Permission
}

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

// CreateEmployeeRoleInput ...
type CreateEmployeeRoleInput struct {
	LocationID    string   `db:"location_id"`
	Name          string   `db:"name"`
	PermissionIDs []string `db:"permission_ids"`
}

// CreateEmployeeRole creates employeeRole
func (s *EmployeeRoleService) CreateEmployeeRole(ctx context.Context, input *CreateEmployeeRoleInput, actor Actor) (*EmployeeRole, error) {
	const op = "app/employeeRoleService.CreateEmployeeRole"

	err := actor.can(ctx, opCreateEmployeeRole)

	if err != nil {
		return nil, errors.Unauthorized(op, err)
	}

	now := time.Now()

	permissions, err := getPermissionsByPermissionIDs(input.PermissionIDs)

	if err != nil {
		return nil, errors.Unexpected(op, err, "failed to get permissions")
	}

	if input.Name == "" {
		return nil, errors.Invalid(op, "name field required")
	}

	if input.PermissionIDs == nil {
		return nil, errors.Invalid(op, "permissions field required")
	}

	employeeRole := &EmployeeRole{
		ID:            uuid.Must(uuid.New(), nil).String(),
		LocationID:    input.LocationID,
		Name:          strings.TrimSpace(input.Name),
		PermissionIDs: input.PermissionIDs,
		Permissions:   permissions,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	err = s.employeeRoleStore.StoreEmployeeRole(ctx, employeeRole)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to employeeRoleStore employeeRole")
	}

	return employeeRole, nil
}

// UpdateEmployeeRoleInput ...
type UpdateEmployeeRoleInput struct {
	Name          string   `db:"name"`
	PermissionIDs []string `db:"permission_ids"`
}

// UpdateEmployeeRole updates employeeRole
func (s *EmployeeRoleService) UpdateEmployeeRole(ctx context.Context, id string, input *UpdateEmployeeRoleInput, actor Actor) (*EmployeeRole, error) {
	const op = "app/employeeRoleService.UpdateEmployeeRole"

	err := actor.can(ctx, opUpdateEmployeeRole)

	if err != nil {
		return nil, errors.Unauthorized(op, err)
	}

	employeeRole, err := s.employeeRoleStore.GetEmployeeRoleByID(ctx, id)

	if employeeRole.Name == "owner" {
		return nil, errors.Invalid(op, "cannot update owner role")
	}

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to get employeeRole by id")
	}

	if employeeRole == nil {
		return nil, errors.NotFound(op)
	}

	employeeRole.UpdatedAt = time.Now()

	if input.Name != "" {
		employeeRole.Name = strings.TrimSpace(input.Name)
	}
	if input.PermissionIDs != nil {
		employeeRole.PermissionIDs = input.PermissionIDs
		permissions, err := getPermissionsByPermissionIDs(input.PermissionIDs)

		if err != nil {
			return nil, errors.Unexpected(op, err, "failed to get permissions")
		}

		employeeRole.Permissions = permissions
	}

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

	if employeesWithTheRole == nil {
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
