package app

import (
	"time"

	"github.com/google/uuid"
)

// EmployeeRole ...
type EmployeeRole struct {
	ID          string    `db:"id"`
	LocationID  string    `db:"location_id"`
	Name        string    `db:"name"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	Permissions []Permission
}

// NewEmployeeRole constructor for EmployeeRole
func NewEmployeeRole(locationID string, name string, permissions []Permission) *EmployeeRole {
	now := time.Now()
	id := uuid.Must(uuid.New(), nil).String()

	newEmployeeRole := EmployeeRole{
		ID:          id,
		LocationID:  locationID,
		Name:        name,
		Permissions: permissions,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return &newEmployeeRole
}
