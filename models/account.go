package models

import (
	"time"

	"github.com/google/uuid"
)

// Account ...
type Account struct {
	ID                    string    `db:"id"`
	FullName              string    `db:"full_name"`
	PhoneNumber           string    `db:"phone_number"`
	CountryCode           string    `db:"country_code"`
	IsPhoneNumberVerified bool      `db:"is_phone_number_verified"`
	CreatedAt             time.Time `db:"created_at"`
	UpdatedAt             time.Time `db:"updated_at"`
}

// NewAccount constructor for Account
func NewAccount(phoneNumber string, countryCode string) *Account {
	now := time.Now()
	id := uuid.Must(uuid.New(), nil).String()

	newAccount := Account{
		ID:                    id,
		PhoneNumber:           phoneNumber,
		CountryCode:           countryCode,
		IsPhoneNumberVerified: false,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	return &newAccount
}
