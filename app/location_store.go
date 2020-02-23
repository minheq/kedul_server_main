package app

import (
	"context"
	"database/sql"

	"github.com/minheq/kedul_server_main/errors"
)

// LocationStore ...
type LocationStore interface {
	GetLocationByID(ctx context.Context, id string) (*Location, error)
	StoreLocation(ctx context.Context, location *Location) error
	UpdateLocation(ctx context.Context, location *Location) error
	DeleteLocation(ctx context.Context, location *Location) error
}

type locationStore struct {
	db *sql.DB
}

// NewLocationStore ...
func NewLocationStore(db *sql.DB) LocationStore {
	return &locationStore{db: db}
}

// GetLocationByID gets Location by ID
func (s *locationStore) GetLocationByID(ctx context.Context, id string) (*Location, error) {
	const op = "location/locationStore.GetLocationByPhoneNumber"

	query := `
		SELECT id, business_id, name, profile_image_id, created_at, updated_at
		FROM location
		WHERE id=$1;
	`

	var location Location

	row := s.db.QueryRow(query, id)

	err := row.Scan(&location.ID, &location.BusinessID, &location.Name, &location.ProfileImageID, &location.CreatedAt, &location.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(op, err, "database error")
	}

	return &location, nil
}

// StoreLocation persists Location
func (s *locationStore) StoreLocation(ctx context.Context, location *Location) error {
	const op = "location/locationStore.StoreLocation"

	query := `
		INSERT INTO location (id, business_id, name, profile_image_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := s.db.Exec(query, location.ID, location.BusinessID, location.Name, location.ProfileImageID, location.CreatedAt, location.UpdatedAt)

	if err != nil {
		return errors.Wrap(op, err, "database error")
	}

	return nil
}

// UpdateLocation updates Location including all fields
func (s *locationStore) UpdateLocation(ctx context.Context, location *Location) error {
	const op = "location/locationStore.UpdateLocation"

	query := `
		UPDATE location
		SET name=$2, profile_image_id=$3, updated_at=$4
		WHERE id=$1;
	`

	_, err := s.db.Exec(query, location.ID, location.Name, location.ProfileImageID, location.UpdatedAt)

	if err != nil {
		return errors.Wrap(op, err, "database error")
	}

	return nil
}

// DeleteLocation deletes Location
func (s *locationStore) DeleteLocation(ctx context.Context, location *Location) error {
	const op = "location/locationStore.DeleteLocation"

	query := `
		DELETE FROM location
		WHERE id=$1;
	`

	_, err := s.db.Exec(query, location.ID)

	if err != nil {
		return errors.Wrap(op, err, "database error")
	}

	return nil
}
