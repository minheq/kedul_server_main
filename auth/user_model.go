package auth

import (
	"time"

	"github.com/google/uuid"
)

// User ...
type User struct {
	ID                    string
	FullName              string
	PhoneNumber           string
	CountryCode           string
	ProfileImageID        string
	IsPhoneNumberVerified bool
	CreatedAt             time.Time
	UpdatedAt             time.Time
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
