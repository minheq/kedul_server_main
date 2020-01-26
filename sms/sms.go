package sms

import (
	"fmt"
)

// Sender is a SMTP mailer.
type Sender interface {
	SendSMS(phoneNumber string, countryCode string, text string) error
}

type sender struct{}

// NewSender constructor for SMSSender
func NewSender() Sender {
	return &sender{}
}

func (s *sender) SendSMS(phoneNumber string, countryCode string, text string) error {
	fmt.Println(phoneNumber, countryCode, text)
	return nil
}
