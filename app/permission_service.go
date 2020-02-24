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
)

var (
	permManageLocation     = Permission{ID: "1", Name: "manage_location", Operations: []Operation{opUpdateLocation}}
	permManageEmployeeRole = Permission{ID: "2", Name: "manage_employee_role", Operations: []Operation{opCreateEmployeeRole, opReadEmployeeRole, opUpdateEmployeeRole, opDeleteEmployeeRole}}
)

var (
	// Exhaustive list of the above permissions
	permissionList = []Permission{permManageLocation, permManageEmployeeRole}
)

type permissionService struct {
	employeeRoleStore EmployeeRoleStore
}

func (p *permissionService) GetEmployeePermissions(ctx context.Context, employeeRoleID string) ([]Permission, error) {
	const op = "app/permissionService.GetEmployeePermissions"

	employeeRole, err := p.employeeRoleStore.GetEmployeeRoleByID(ctx, employeeRoleID)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to get employee role by id")
	}

	if employeeRole == nil {
		return nil, errors.NotFound(op)
	}

	return employeeRole.Permissions, nil
}

func getPermissions(employeeRole *EmployeeRole) ([]Permission, error) {
	const op = "app/getPermissions"

	permissions := []Permission{}

	for _, permissionID := range employeeRole.PermissionIDs {
		found := false

		for _, permission := range permissionList {
			if permission.ID == permissionID {
				permissions = append(permissions, permission)
				found = true
			}
		}

		if found == false {
			return nil, errors.Unexpected(op, fmt.Errorf("permission for permissionID=%s not found", permissionID), "failed to retrieve permission")
		}
	}

	return permissions, nil
}
