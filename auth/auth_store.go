package auth

import (
	"context"
	"database/sql"

	"github.com/minheq/kedul_server_main/errors"
)

// Store ...
type Store struct {
	db *sql.DB
}

// NewStore ...
func NewStore(db *sql.DB) Store {
	return Store{db: db}
}

// GetVerificationCodeByIDAndCode gets VerificationCode by code
func (s *Store) GetVerificationCodeByIDAndCode(ctx context.Context, verificationID string, code string) (*VerificationCode, error) {
	const op = "auth/auth_store.GetVerificationCodeByIDAndCode"

	query := `
		SELECT id, user_id, code, verification_id, code_type, phone_number, country_code, expired_at, created_at
		FROM verification_code
		WHERE verification_id=$1
			AND code=$2;
	`

	var vc VerificationCode

	row := s.db.QueryRow(query, verificationID, code)

	err := row.Scan(&vc.ID, &vc.UserID, &vc.Code, &vc.VerificationID, &vc.CodeType, &vc.PhoneNumber, &vc.CountryCode, &vc.ExpiredAt, &vc.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.NotFound(op)
	}

	if err != nil {
		return nil, errors.Wrap(op, err, "database error")
	}

	return &vc, nil
}

// StoreVerificationCode persists VerificationCode
func (s *Store) StoreVerificationCode(ctx context.Context, vc *VerificationCode) error {
	const op = "auth/auth_store.StoreVerificationCode"

	query := `
		INSERT INTO verification_code (id, user_id, code, verification_id, code_type, phone_number, country_code, expired_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := s.db.Exec(query, vc.ID, vc.UserID, vc.Code, vc.VerificationID, vc.CodeType, vc.PhoneNumber, vc.CountryCode, vc.ExpiredAt, vc.CreatedAt)

	if err != nil {
		return errors.Wrap(op, err, "database error")
	}

	return nil
}

// RemoveVerificationCodeByPhoneNumber removes VerificationCode
func (s *Store) RemoveVerificationCodeByPhoneNumber(ctx context.Context, phoneNumber string, countryCode string) error {
	const op = "auth/auth_store.RemoveVerificationCodeByPhoneNumber"

	query := `
		DELETE FROM verification_code
		WHERE phone_number=$1 AND country_code=$2;
	`

	_, err := s.db.Exec(query, phoneNumber, countryCode)

	if err != nil {
		return errors.Wrap(op, err, "database error")
	}

	return nil
}

// RemoveVerificationCodeByID removes VerificationCode by Id
func (s *Store) RemoveVerificationCodeByID(ctx context.Context, id string) error {
	const op = "auth/auth_store.RemoveVerificationCodeByID"

	query := `
		DELETE FROM verification_code
		WHERE id=$1;
	`

	_, err := s.db.Exec(query, id)

	if err != nil {
		return errors.Wrap(op, err, "database error")
	}

	return nil
}

// GetUserByID gets User by ID
func (s *Store) GetUserByID(ctx context.Context, id string) (*User, error) {
	const op = "auth/auth_store.GetUserByPhoneNumber"

	query := `
		SELECT id, full_name, phone_number, country_code, is_phone_number_verified, created_at, updated_at
		FROM kedul_user
		WHERE id=$1;
	`

	var user User

	row := s.db.QueryRow(query, id)

	err := row.Scan(&user.ID, &user.FullName, &user.PhoneNumber, &user.CountryCode, &user.IsPhoneNumberVerified, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.NotFound(op)
	}

	if err != nil {
		return nil, errors.Wrap(op, err, "database error")
	}

	return &user, nil
}

// GetUserByPhoneNumber gets User by Phone Number
func (s *Store) GetUserByPhoneNumber(ctx context.Context, phoneNumber string, countryCode string) (*User, error) {
	const op = "auth/auth_store.GetUserByPhoneNumber"

	query := `
		SELECT id, full_name, phone_number, country_code, is_phone_number_verified, created_at, updated_at
		FROM kedul_user
		WHERE phone_number=$1
			AND country_code=$2;
	`

	var user User

	row := s.db.QueryRow(query, phoneNumber, countryCode)

	err := row.Scan(&user.ID, &user.FullName, &user.PhoneNumber, &user.CountryCode, &user.IsPhoneNumberVerified, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.NotFound(op)
	}

	if err != nil {
		return nil, errors.Wrap(op, err, "database error")
	}

	return &user, nil
}

// StoreUser persists User
func (s *Store) StoreUser(ctx context.Context, user *User) error {
	const op = "auth/auth_store.StoreUser"

	query := `
		INSERT INTO kedul_user (id, full_name, phone_number, country_code, is_phone_number_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := s.db.Exec(query, user.ID, user.FullName, user.PhoneNumber, user.CountryCode, user.IsPhoneNumberVerified, user.CreatedAt, user.UpdatedAt)

	if err != nil {
		return errors.Wrap(op, err, "database error")
	}

	return nil
}

// UpdateUser updates User including all fields
func (s *Store) UpdateUser(ctx context.Context, user *User) error {
	const op = "auth/auth_store.UpdateUser"

	query := `
		UPDATE kedul_user
		SET full_name=$2, phone_number=$3, country_code=$4, is_phone_number_verified=$5, created_at=$6, updated_at=$7
		WHERE id=$1;
	`

	_, err := s.db.Exec(query, user.ID, user.FullName, user.PhoneNumber, user.CountryCode, user.IsPhoneNumberVerified, user.CreatedAt, user.UpdatedAt)

	if err != nil {
		return errors.Wrap(op, err, "database error")
	}

	return nil
}
