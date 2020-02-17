package auth

import (
	"time"

	"github.com/google/uuid"
)

// User ...
type User struct {
	ID                    string    `db:"id"`
	FullName              string    `db:"full_name"`
	PhoneNumber           string    `db:"phone_number"`
	CountryCode           string    `db:"country_code"`
	IsPhoneNumberVerified bool      `db:"is_phone_number_verified"`
	CreatedAt             time.Time `db:"created_at"`
	UpdatedAt             time.Time `db:"updated_at"`
}

// NewUser constructor for User
func NewUser(phoneNumber string, countryCode string) *User {
	now := time.Now()
	id := uuid.Must(uuid.New(), nil).String()

	newUser := User{
		ID:                    id,
		PhoneNumber:           phoneNumber,
		CountryCode:           countryCode,
		IsPhoneNumberVerified: false,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	return &newUser
}
