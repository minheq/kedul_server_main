package phone

import (
	"github.com/minheq/kedul_server_main/errors"
	"github.com/nyaruka/phonenumbers"
)

// FormatPhoneNumber formats phone number
func FormatPhoneNumber(phoneNumber string, countryCode string) (string, error) {
	const op = "phone/phone_number.formatPhoneNumber"
	parsedPhoneNumber, err := phonenumbers.Parse(phoneNumber, countryCode)

	if err != nil {
		return "", errors.Wrap(op, err, "failed to parse phone number")
	}

	formattedPhoneNumber := phonenumbers.Format(parsedPhoneNumber, phonenumbers.NATIONAL)

	return formattedPhoneNumber, nil
}
