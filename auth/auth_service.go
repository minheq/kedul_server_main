package auth

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/minheq/kedul_server_main/errors"
	"github.com/minheq/kedul_server_main/random"
	"github.com/minheq/kedul_server_main/sms"
	"github.com/nyaruka/phonenumbers"
)

// Service handles authentication
type Service struct {
	store     Store
	tokenAuth *jwtauth.JWTAuth
	smsSender sms.Sender
}

// NewService constructor for AuthService
func NewService(store Store, tokenAuth *jwtauth.JWTAuth, smsSender sms.Sender) Service {
	return Service{store: store, tokenAuth: tokenAuth, smsSender: smsSender}
}

// LoginVerify login verification initialization core logic
func (as *Service) LoginVerify(ctx context.Context, phoneNumber string, countryCode string) (string, error) {
	const op = "handlers/login_verify.LoginVerify"

	parsedPhoneNumber, err := phonenumbers.Parse(phoneNumber, countryCode)

	if err != nil {
		return "", errors.Invalid(op, "phone number supplied is invalid")
	}

	formattedPhoneNumber := phonenumbers.Format(parsedPhoneNumber, phonenumbers.NATIONAL)

	user, err := as.store.GetUserByPhoneNumber(ctx, formattedPhoneNumber, countryCode)

	if err != nil && errors.Is(errors.KindNotFound, err) == false {
		return "", errors.Unexpected(op, err)
	}

	// Create and persist new user if it didn't exist yet
	if user == nil {
		user = NewUser(formattedPhoneNumber, countryCode)

		err = as.store.StoreUser(ctx, user)

		if err != nil {
			return "", errors.Unexpected(op, err)
		}
	}

	err = as.store.RemoveVerificationCodeByPhoneNumber(ctx, formattedPhoneNumber, countryCode)

	if err != nil {
		return "", errors.Unexpected(op, err)
	}

	verificationID := random.String(50)
	code := random.Number(6)

	newVerificationCode := NewVerificationCode(verificationID, code, user.ID, formattedPhoneNumber, countryCode, "LOGIN")

	err = as.store.StoreVerificationCode(ctx, newVerificationCode)

	if err != nil {
		return "", errors.Unexpected(op, err)
	}

	err = as.smsSender.SendSMS(formattedPhoneNumber, countryCode, code)

	if err != nil {
		return "", errors.Unexpected(op, err)
	}

	return verificationID, nil
}

// LoginVerifyCheck returns accessToken given the verificationID and code match the persisted verification code
func (as *Service) LoginVerifyCheck(ctx context.Context, verificationID string, code string) (string, error) {
	const op = "handlers/login_verify_check.LoginVerifyCheck"

	verificationCode, err := as.store.GetVerificationCodeByIDAndCode(ctx, verificationID, code)

	if err != nil {
		return "", errors.Invalid(op, "verification code not found")
	}

	if verificationCode.ExpiredAt.Before(time.Now()) {
		return "", errors.Invalid(op, "verification code expired")
	}

	user, err := as.store.GetUserByID(ctx, verificationCode.UserID)

	if err != nil {
		return "", errors.Invalid(op, "user not found")
	}

	user.IsPhoneNumberVerified = true
	user.UpdatedAt = time.Now()

	err = as.store.UpdateUser(ctx, user)

	if err != nil {
		return "", errors.Unexpected(op, err)
	}

	err = as.store.RemoveVerificationCodeByID(ctx, verificationCode.ID)

	if err != nil {
		return "", errors.Unexpected(op, err)
	}

	_, accessToken, err := as.tokenAuth.Encode(jwt.MapClaims{"user_id": user.ID})

	if err != nil {
		return "", errors.Unexpected(op, err)
	}

	return accessToken, nil
}
