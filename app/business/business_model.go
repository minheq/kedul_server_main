package business

import (
	"time"

	"github.com/google/uuid"
)

// Business ...
type Business struct {
	ID             string    `db:"id"`
	UserID         string    `db:"user_id"`
	Name           string    `db:"name"`
	ProfileImageID string    `db:"profile_image_id"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

// NewBusiness constructor for Business
func NewBusiness(userID string, name string) *Business {
	now := time.Now()
	id := uuid.Must(uuid.New(), nil).String()

	newBusiness := Business{
		ID:        id,
		UserID:    userID,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return &newBusiness
}
