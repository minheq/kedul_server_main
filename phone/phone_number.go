package phone

import (
	"github.com/nyaruka/phonenumbers"
)

// FormatPhoneNumber formats phone number
func FormatPhoneNumber(phoneNumber string, countryCode string) (string, error) {
	const op = "auth/auth_service.formatPhoneNumber"
	parsedPhoneNumber, err := phonenumbers.Parse(phoneNumber, countryCode)

	if err != nil {
		return "", err
	}

	formattedPhoneNumber := phonenumbers.Format(parsedPhoneNumber, phonenumbers.NATIONAL)

	return formattedPhoneNumber, nil
}
