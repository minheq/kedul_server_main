package account

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/google/uuid"
	"github.com/minheq/kedulv2/service_salon/errors"
	"github.com/minheq/kedulv2/service_salon/random"
	"github.com/minheq/kedulv2/service_salon/sms"
	"github.com/nyaruka/phonenumbers"
)

// Service for AccountService
type Service struct {
	store     Store
	smsSender sms.Sender
	tokenAuth *jwtauth.JWTAuth
}

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

// VerificationCode ...
type VerificationCode struct {
	ID          string    `db:"id"`
	AccountID   string    `db:"account_id"`
	Code        string    `db:"code"`
	ClientState string    `db:"client_state"`
	CodeType    string    `db:"code_type"`
	PhoneNumber string    `db:"phone_number"`
	CountryCode string    `db:"country_code"`
	CreatedAt   time.Time `db:"created_at"`
	ExpiredAt   time.Time `db:"expired_at"`
}

// NewVerificationCode constructor for VerificationCode
func NewVerificationCode(clientState string, code string, accountID string, phoneNumber string, countryCode string, codeType string) *VerificationCode {
	now := time.Now()
	id := uuid.Must(uuid.New(), nil).String()

	verificationCode := VerificationCode{
		ID:          id,
		AccountID:   accountID,
		Code:        code,
		CodeType:    codeType,
		ClientState: clientState,
		PhoneNumber: phoneNumber,
		CountryCode: countryCode,
		ExpiredAt:   now.Add(time.Duration(10) * time.Minute),
		CreatedAt:   now,
	}

	return &verificationCode
}

// NewAccountService ...
func NewAccountService(store Store, smsSender sms.Sender, tokenAuth *jwtauth.JWTAuth) Service {
	return Service{store: store, smsSender: smsSender, tokenAuth: tokenAuth}
}

// LogInStart returns clientState which needs to be verified during LogInVerify step
func (s *Service) LogInStart(phoneNumber string, countryCode string) (string, error) {
	const op errors.Op = "account/service.LogInStart"

	parsedPhoneNumber, err := phonenumbers.Parse(phoneNumber, countryCode)

	if err != nil {
		return "", errors.E(op, err, "failed to parse phone number")
	}

	formattedPhoneNumber := phonenumbers.Format(parsedPhoneNumber, phonenumbers.NATIONAL)

	clientState := random.String(50)
	code := random.Number(6)

	account, err := s.store.GetAccountByPhoneNumber(formattedPhoneNumber, countryCode)

	if err != nil && errors.Is(errors.NotFound, err) == false {
		return "", errors.E(op, err, "could not get account")
	}

	// Create and persist new account if it didn't exist yet
	if account == nil {
		account = NewAccount(formattedPhoneNumber, countryCode)

		err = s.store.StoreAccount(account)

		if err != nil {
			return "", errors.E(op, err, "could not store account")
		}
	}

	err = s.store.RemoveVerificationCodeByPhoneNumber(formattedPhoneNumber, countryCode)

	if err != nil {
		return "", errors.E(op, err, "could not remove verification code")
	}

	newVerificationCode := NewVerificationCode(clientState, code, account.ID, formattedPhoneNumber, countryCode, "LOGIN")

	err = s.store.StoreVerificationCode(newVerificationCode)

	if err != nil {
		return "", errors.E(op, err, "could not store verification code")
	}

	err = s.smsSender.SendSMS(formattedPhoneNumber, countryCode, code)

	if err != nil {
		return "", errors.E(op, err, "send sms failed")
	}

	return clientState, nil
}

// LogInVerify returns accessToken given the clientState and code match the persisted verification code
func (s *Service) LogInVerify(clientState string, code string) (string, error) {
	const op errors.Op = "account/service.LogInVerify"

	verificationCode, err := s.store.GetVerificationCodeByClientStateAndCode(clientState, code)

	if err != nil {
		return "", errors.E(op, err, "could not get verification code")
	}

	if verificationCode.ExpiredAt.Before(time.Now()) {
		return "", errors.E(op, "verification code expired", errors.Forbidden)
	}

	// Verify

	account, err := s.store.GetAccountByID(verificationCode.AccountID)

	if err != nil {
		return "", errors.E(op, err, "could not get account")
	}

	account.IsPhoneNumberVerified = true
	account.UpdatedAt = time.Now()

	err = s.store.UpdateAccount(account)

	if err != nil {
		return "", errors.E(op, err, "could not update account")
	}

	err = s.store.RemoveVerificationCodeByID(verificationCode.ID)

	if err != nil {
		return "", errors.E(op, err, "could not remove verification code")
	}

	_, accessToken, err := s.tokenAuth.Encode(jwt.MapClaims{"account_id": account.ID})

	if err != nil {
		return "", errors.E(op, err, "failed to generate access token")
	}

	return accessToken, nil
}
