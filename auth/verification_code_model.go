package auth

import (
	"time"

	"github.com/google/uuid"
)

// VerificationCode ...
type VerificationCode struct {
	ID             string    `db:"id"`
	UserID         string    `db:"user_id"`
	Code           string    `db:"code"`
	VerificationID string    `db:"verification_id"`
	CodeType       string    `db:"code_type"`
	PhoneNumber    string    `db:"phone_number"`
	CountryCode    string    `db:"country_code"`
	CreatedAt      time.Time `db:"created_at"`
	ExpiredAt      time.Time `db:"expired_at"`
}

// NewVerificationCode constructor for VerificationCode
func NewVerificationCode(verificationID string, code string, userID string, phoneNumber string, countryCode string, codeType string) *VerificationCode {
	now := time.Now()
	id := uuid.Must(uuid.New(), nil).String()

	verificationCode := VerificationCode{
		ID:             id,
		UserID:         userID,
		Code:           code,
		CodeType:       codeType,
		VerificationID: verificationID,
		PhoneNumber:    phoneNumber,
		CountryCode:    countryCode,
		ExpiredAt:      now.Add(time.Duration(10) * time.Minute),
		CreatedAt:      now,
	}

	return &verificationCode
}
