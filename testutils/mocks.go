package testutils

type SmsSenderMock struct {
	Text string
}

func (s *SmsSenderMock) SendSMS(phoneNumber string, countryCode string, text string) error {
	s.Text = text
	return nil
}
