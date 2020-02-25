package app

import (
	"context"
	"fmt"

	"github.com/minheq/kedul_server_main/errors"
)

var (
	opUpdateLocation     = Operation{Name: "update_location"}
	opCreateEmployeeRole = Operation{Name: "create_employee_role"}
	opReadEmployeeRole   = Operation{Name: "read_employee_role"}
	opUpdateEmployeeRole = Operation{Name: "update_employee_role"}
	opDeleteEmployeeRole = Operation{Name: "delete_employee_role"}
	opCreateEmployee     = Operation{Name: "create_employee"}
	opReadEmployee       = Operation{Name: "read_employee"}
	opUpdateEmployee     = Operation{Name: "update_employee"}
	opDeleteEmployee     = Operation{Name: "delete_employee"}
)

var (
	permManageLocation     = Permission{ID: "1", Name: "manage_location", Operations: []Operation{opUpdateLocation}}
	permManageEmployeeRole = Permission{ID: "2", Name: "manage_employee_role", Operations: []Operation{opCreateEmployeeRole, opReadEmployeeRole, opUpdateEmployeeRole, opDeleteEmployeeRole}}
	permManageEmployee     = Permission{ID: "3", Name: "manage_employee", Operations: []Operation{opCreateEmployee, opReadEmployee, opUpdateEmployee, opDeleteEmployee}}
)

var permissionsTable = map[string]Permission{
	permManageLocation.ID:     permManageLocation,
	permManageEmployeeRole.ID: permManageEmployeeRole,
	permManageEmployee.ID:     permManageEmployee,
}

// PermissionService ...
type PermissionService struct {
	employeeRoleStore EmployeeRoleStore
	employeeStore     EmployeeStore
}

// NewPermissionService constructor for AuthService
func NewPermissionService(employeeRoleStore EmployeeRoleStore, employeeStore EmployeeStore) PermissionService {
	return PermissionService{employeeRoleStore: employeeRoleStore, employeeStore: employeeStore}
}

// GetEmployeeActor ...
func (p *PermissionService) GetEmployeeActor(ctx context.Context, userID string, locationID string) (Actor, error) {
	const op = "app/permissionService.GetEmployeeActor"

	permissions := []Permission{}

	employee, err := p.employeeStore.GetEmployeeByUserIDAndLocationID(ctx, userID, locationID)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to get employee role by id")
	}

	if employee == nil {
		return nil, errors.NotFound(op)
	}

	employeeRole, err := p.employeeRoleStore.GetEmployeeRoleByID(ctx, employee.EmployeeRoleID)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to get employee role by id")
	}

	if employee == nil {
		return nil, errors.NotFound(op)
	}

	permissions = employeeRole.Permissions

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to get actor permissions")
	}

	actor := NewActor(permissions)

	return actor, nil
}

func getPermissionsByPermissionIDs(permissionIDs []string) ([]Permission, error) {
	const op = "app/getPermissionsByPermissionIDs"

	permissions := []Permission{}

	for _, permissionID := range permissionIDs {
		permission := permissionsTable[permissionID]

		if permission.ID == "" {
			return nil, errors.Unexpected(op, fmt.Errorf("permission for permissionID=%s not found", permissionID), "failed to retrieve permission")
		}

		permissions = append(permissions, permission)
	}

	return permissions, nil
}
