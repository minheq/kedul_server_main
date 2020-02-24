package app

import (
	"context"
)

var (
	opUpdateLocation     = Operation{Name: "update_location"}
	opCreateEmployeeRole = Operation{Name: "create_employee_role"}
	opReadEmployeeRole   = Operation{Name: "read_employee_role"}
	opUpdateEmployeeRole = Operation{Name: "update_employee_role"}
	opDeleteEmployeeRole = Operation{Name: "delete_employee_role"}
)

var (
	permManageLocation     = Permission{Name: "manage_location", Operations: []Operation{opUpdateLocation}}
	permManageEmployeeRole = Permission{Name: "manage_employee_role", Operations: []Operation{opCreateEmployeeRole, opReadEmployeeRole, opUpdateEmployeeRole, opDeleteEmployeeRole}}
)

// Actor is the current caller. It can be a user or API key or anything that is allowed to interact with our API
type Actor interface {
	can(ctx context.Context, operation Operation) error
}

type actor struct{}

func (a *actor) can(ctx context.Context, operation Operation) error {
	const op = "app/actor.can"
	// get current employee based on locationID + userID
	// get his employee role's permissions
	// find a match in the list of operations the permissions allow

	return nil
}
