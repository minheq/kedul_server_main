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
)

var permissionsTable = map[string]Permission{
	"1": permManageLocation,
	"2": permManageEmployeeRole,
}

type permissionService struct {
	employeeRoleStore EmployeeRoleStore
	employeeStore     EmployeeStore
}

func (p *permissionService) GetEmployeePermissions(ctx context.Context, userID string, locationID string) ([]Permission, error) {
	const op = "app/permissionService.GetEmployeePermissions"

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

	return employeeRole.Permissions, nil
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
