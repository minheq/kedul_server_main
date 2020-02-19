package phone

import (
	"fmt"
)

// SMSSender sends sms.
type SMSSender interface {
	SendSMS(phoneNumber string, countryCode string, text string) error
}

type sender struct{}

// NewSMSSender constructor for SMSSender
func NewSMSSender() SMSSender {
	return &sender{}
}

func (s *sender) SendSMS(phoneNumber string, countryCode string, text string) error {
	fmt.Println(phoneNumber, countryCode, text)
	return nil
}
