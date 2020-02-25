package app

import (
	"context"
	"fmt"
)

// Actor is the current caller. It can be a user or API key or anything that is allowed to interact with our API
type Actor interface {
	can(ctx context.Context, operation Operation) error
}

type actor struct {
	permissions []Permission
}

// NewActor ...
func NewActor(permissions []Permission) Actor {
	return &actor{permissions: permissions}
}

func (a *actor) can(ctx context.Context, operation Operation) error {
	const op = "app/actor.can"

	operations := []Operation{}

	for _, permission := range a.permissions {
		operations = append(operations, permission.Operations...)
	}

	for _, op := range operations {
		if op == operation {
			return nil
		}
	}

	return fmt.Errorf("permission for operation=%s not found", operation)
}
