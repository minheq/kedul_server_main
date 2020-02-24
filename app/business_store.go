package app

import (
	"context"
	"database/sql"

	"github.com/minheq/kedul_server_main/errors"
)

// BusinessStore ...
type BusinessStore interface {
	GetBusinessByID(ctx context.Context, id string) (*Business, error)
	GetBusinessByName(ctx context.Context, name string) (*Business, error)
	StoreBusiness(ctx context.Context, b *Business) error
	UpdateBusiness(ctx context.Context, b *Business) error
	DeleteBusiness(ctx context.Context, b *Business) error
}

// businessStore ...
type businessStore struct {
	db *sql.DB
}

// NewStore ...
func NewStore(db *sql.DB) BusinessStore {
	return &businessStore{db: db}
}

// GetBusinessByID gets Business by ID
func (s *businessStore) GetBusinessByID(ctx context.Context, id string) (*Business, error) {
	const op = "business/businessStore.GetBusinessByPhoneNumber"

	query := `
		SELECT id, user_id, name, profile_image_id, created_at, updated_at
		FROM business
		WHERE id=$1;
	`

	var b Business

	row := s.db.QueryRow(query, id)

	err := row.Scan(&b.ID, &b.UserID, &b.Name, &b.ProfileImageID, &b.CreatedAt, &b.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(op, err, "database error")
	}

	return &b, nil
}

// GetBusinessByName gets Business by name
func (s *businessStore) GetBusinessByName(ctx context.Context, name string) (*Business, error) {
	const op = "business/businessStore.GetBusinessByName"

	query := `
		SELECT id, user_id, name, profile_image_id, created_at, updated_at
		FROM business
		WHERE name=$1;
	`

	var b Business

	row := s.db.QueryRow(query, name)

	err := row.Scan(&b.ID, &b.UserID, &b.Name, &b.ProfileImageID, &b.CreatedAt, &b.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(op, err, "database error")
	}

	return &b, nil
}

// StoreBusiness persists Business
func (s *businessStore) StoreBusiness(ctx context.Context, b *Business) error {
	const op = "business/businessStore.StoreBusiness"

	query := `
		INSERT INTO business (id, user_id, name, profile_image_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := s.db.Exec(query, b.ID, b.UserID, b.Name, b.ProfileImageID, b.CreatedAt, b.UpdatedAt)

	if err != nil {
		return errors.Wrap(op, err, "database error")
	}

	return nil
}

// UpdateBusiness updates Business including all fields
func (s *businessStore) UpdateBusiness(ctx context.Context, b *Business) error {
	const op = "business/businessStore.UpdateBusiness"

	query := `
		UPDATE business
		SET name=$2, profile_image_id=$3, updated_at=$4
		WHERE id=$1;
	`

	_, err := s.db.Exec(query, b.ID, b.Name, b.ProfileImageID, b.UpdatedAt)

	if err != nil {
		return errors.Wrap(op, err, "database error")
	}

	return nil
}

// DeleteBusiness deletes Business
func (s *businessStore) DeleteBusiness(ctx context.Context, b *Business) error {
	const op = "business/businessStore.DeleteBusiness"

	query := `
		DELETE FROM business
		WHERE id=$1;
	`

	_, err := s.db.Exec(query, b.ID)

	if err != nil {
		return errors.Wrap(op, err, "database error")
	}

	return nil
}