package business

import (
	"context"
	"database/sql"

	"github.com/minheq/kedul_server_main/errors"
)

// Store ...
type Store interface {
	GetBusinessByID(ctx context.Context, id string) (*Business, error)
	GetBusinessByName(ctx context.Context, name string) (*Business, error)
	StoreBusiness(ctx context.Context, b *Business) error
	UpdateBusiness(ctx context.Context, b *Business) error
	DeleteBusiness(ctx context.Context, b *Business) error
}

// store ...
type store struct {
	db *sql.DB
}

// NewStore ...
func NewStore(db *sql.DB) Store {
	return &store{db: db}
}

// GetBusinessByID gets Business by ID
func (s *store) GetBusinessByID(ctx context.Context, id string) (*Business, error) {
	const op = "business/store.GetBusinessByPhoneNumber"

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
func (s *store) GetBusinessByName(ctx context.Context, name string) (*Business, error) {
	const op = "business/store.GetBusinessByName"

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
func (s *store) StoreBusiness(ctx context.Context, b *Business) error {
	const op = "business/store.StoreBusiness"

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
func (s *store) UpdateBusiness(ctx context.Context, b *Business) error {
	const op = "business/store.UpdateBusiness"

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
func (s *store) DeleteBusiness(ctx context.Context, b *Business) error {
	const op = "business/store.DeleteBusiness"

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
