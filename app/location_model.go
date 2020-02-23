package app

import (
	"time"

	"github.com/google/uuid"
)

// Location ...
type Location struct {
	ID             string    `db:"id"`
	BusinessID     string    `db:"business_id"`
	Name           string    `db:"name"`
	ProfileImageID string    `db:"profile_image_id"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

// NewLocation constructor for Location
func NewLocation(businessID string, name string) *Location {
	now := time.Now()
	id := uuid.Must(uuid.New(), nil).String()

	newLocation := Location{
		ID:         id,
		BusinessID: businessID,
		Name:       name,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	return &newLocation
}
