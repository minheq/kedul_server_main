package account

import (
	"database/sql"
	"github.com/minheq/kedulv2/service_salon/errors"
)

// Store ...
type Store struct {
	db *sql.DB
}

// NewAccountStore ...
func NewAccountStore(db *sql.DB) Store {
	return Store{db: db}
}

// GetAccountByID gets Account by ID
func (s *Store) GetAccountByID(id string) (*Account, error) {
	const op errors.Op = "account/store.GetAccountByPhoneNumber"

	query := `
		SELECT id, full_name, phone_number, country_code, is_phone_number_verified, created_at, updated_at
		FROM account
		WHERE id=$1;
	`

	var account Account

	row := s.db.QueryRow(query, id)

	err := row.Scan(&account.ID, &account.FullName, &account.PhoneNumber, &account.CountryCode, &account.IsPhoneNumberVerified, &account.CreatedAt, &account.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.E(op, err, errors.NotFound)
	}

	if err != nil {
		return nil, errors.E(op, err)
	}

	return &account, nil
}

// GetAccountByPhoneNumber gets Account by Phone Number
func (s *Store) GetAccountByPhoneNumber(phoneNumber string, countryCode string) (*Account, error) {
	const op errors.Op = "account/store.GetAccountByPhoneNumber"

	query := `
		SELECT id, full_name, phone_number, country_code, is_phone_number_verified, created_at, updated_at
		FROM account
		WHERE phone_number=$1
			AND country_code=$2;
	`

	var account Account

	row := s.db.QueryRow(query, phoneNumber, countryCode)

	err := row.Scan(&account.ID, &account.FullName, &account.PhoneNumber, &account.CountryCode, &account.IsPhoneNumberVerified, &account.CreatedAt, &account.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.E(op, err, errors.NotFound)
	}

	if err != nil {
		return nil, errors.E(op, err)
	}

	return &account, nil
}

// StoreAccount persists Account
func (s *Store) StoreAccount(account *Account) error {
	const op errors.Op = "account/store.StoreAccount"

	query := `
		INSERT INTO account (id, full_name, phone_number, country_code, is_phone_number_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := s.db.Exec(query, account.ID, account.FullName, account.PhoneNumber, account.CountryCode, account.IsPhoneNumberVerified, account.CreatedAt, account.UpdatedAt)

	if err != nil {
		return errors.E(op, err)
	}

	return nil
}

// UpdateAccount updates Account including all fields
func (s *Store) UpdateAccount(account *Account) error {
	const op errors.Op = "account/store.UpdateAccount"

	query := `
		UPDATE account
		SET full_name=$2, phone_number=$3, country_code=$4, is_phone_number_verified=$5, created_at=$6, updated_at=$7
		WHERE id=$1;
	`

	_, err := s.db.Exec(query, account.ID, account.FullName, account.PhoneNumber, account.CountryCode, account.IsPhoneNumberVerified, account.CreatedAt, account.UpdatedAt)

	if err != nil {
		return errors.E(op, err)
	}

	return nil
}

// GetVerificationCodeByClientStateAndCode gets VerificationCode by code
func (s *Store) GetVerificationCodeByClientStateAndCode(clientState string, code string) (*VerificationCode, error) {
	const op errors.Op = "account/store.GetVerificationCodeByClientStateAndCode"

	query := `
		SELECT id, account_id, code, client_state, code_type, phone_number, country_code, expired_at, created_at
		FROM verification_code
		WHERE client_state=$1
			AND code=$2;
	`

	var vc VerificationCode

	row := s.db.QueryRow(query, clientState, code)

	err := row.Scan(&vc.ID, &vc.AccountID, &vc.Code, &vc.ClientState, &vc.CodeType, &vc.PhoneNumber, &vc.CountryCode, &vc.ExpiredAt, &vc.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.E(op, err, errors.NotFound)
	}

	if err != nil {
		return nil, errors.E(op, err)
	}

	return &vc, nil
}

// StoreVerificationCode persists VerificationCode
func (s *Store) StoreVerificationCode(vc *VerificationCode) error {
	const op errors.Op = "account/store.StoreVerificationCode"

	query := `
		INSERT INTO verification_code (id, account_id, code, client_state, code_type, phone_number, country_code, expired_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := s.db.Exec(query, vc.ID, vc.AccountID, vc.Code, vc.ClientState, vc.CodeType, vc.PhoneNumber, vc.CountryCode, vc.ExpiredAt, vc.CreatedAt)

	if err != nil {
		return errors.E(op, err)
	}

	return nil
}

// RemoveVerificationCodeByPhoneNumber removes VerificationCode
func (s *Store) RemoveVerificationCodeByPhoneNumber(phoneNumber string, countryCode string) error {
	const op errors.Op = "account/store.RemoveVerificationCodeByPhoneNumber"

	query := `
		DELETE FROM verification_code
		WHERE phone_number=$1 AND country_code=$2;
	`

	_, err := s.db.Exec(query, phoneNumber, countryCode)

	if err != nil {
		return errors.E(op, err)
	}

	return nil
}

// RemoveVerificationCodeByID removes VerificationCode by Id
func (s *Store) RemoveVerificationCodeByID(id string) error {
	const op errors.Op = "account/store.RemoveVerificationCodeByID"

	query := `
		DELETE FROM verification_code
		WHERE id=$1;
	`

	_, err := s.db.Exec(query, id)

	if err != nil {
		return errors.E(op, err)
	}

	return nil
}
