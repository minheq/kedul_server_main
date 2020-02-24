package app

import (
	"time"

	"github.com/google/uuid"
)

// Employee ...
type Employee struct {
	ID             string    `db:"id"`
	LocationID     string    `db:"location_id"`
	Name           string    `db:"name"`
	ProfileImageID string    `db:"profile_image_id"`
	EmployeeRoleID string    `db:"employee_role_id"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

// NewEmployee constructor for Employee
func NewEmployee(locationID string, name string, employeeRoleID string) *Employee {
	now := time.Now()
	id := uuid.Must(uuid.New(), nil).String()

	newEmployee := Employee{
		ID:             id,
		LocationID:     locationID,
		EmployeeRoleID: employeeRoleID,
		Name:           name,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	return &newEmployee
}
