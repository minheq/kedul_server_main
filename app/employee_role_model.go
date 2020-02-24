package app

import (
	"time"

	"github.com/google/uuid"
)

var (
	defaultOwnerRolePermissions        = []Permission{}
	defaultAdminRolePermissions        = []Permission{}
	defaultManagerRolePermissions      = []Permission{}
	defaultReceptionistRolePermissions = []Permission{}
	defaultSpecialistRolePermissions   = []Permission{}
)

// EmployeeRole ...
type EmployeeRole struct {
	ID            string    `db:"id"`
	LocationID    string    `db:"location_id"`
	Name          string    `db:"name"`
	PermissionIDs []string  `db:"permission_ids"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
	// These permissions are retrieved in the application code based on PermissionIDs
	Permissions []Permission
}

// NewEmployeeRole constructor for EmployeeRole
func NewEmployeeRole(locationID string, name string, permissions []Permission) *EmployeeRole {
	now := time.Now()
	id := uuid.Must(uuid.New(), nil).String()

	permissionIDs := []string{}

	for _, permission := range permissions {
		permissionIDs = append(permissionIDs, permission.ID)
	}

	newEmployeeRole := EmployeeRole{
		ID:            id,
		LocationID:    locationID,
		Name:          name,
		PermissionIDs: permissionIDs,
		Permissions:   permissions,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	return &newEmployeeRole
}
