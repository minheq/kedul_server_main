package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/minheq/kedul_server_main/errors"
	"github.com/minheq/kedul_server_main/phone"
	"github.com/minheq/kedul_server_main/random"
)

// Service handles authentication
type Service struct {
	store     Store
	tokenAuth *jwtauth.JWTAuth
	smsSender phone.SMSSender
}

// NewService constructor for AuthService
func NewService(store Store, tokenAuth *jwtauth.JWTAuth, smsSender phone.SMSSender) Service {
	return Service{store: store, tokenAuth: tokenAuth, smsSender: smsSender}
}

func (as *Service) createNewVerificationCode(ctx context.Context, user *User, phoneNumber string, countryCode string, verificationCodeType string) (*VerificationCode, error) {
	const op = "auth/service.createNewVerificationCode"

	err := as.store.DeleteVerificationCodeByPhoneNumber(ctx, phoneNumber, countryCode)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to remove verification code")
	}

	verificationID := random.String(50)
	code := random.Number(6)

	verificationCode := NewVerificationCode(verificationID, code, user.ID, phoneNumber, countryCode, verificationCodeType)

	err = as.store.StoreVerificationCode(ctx, verificationCode)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to store verification code")
	}

	return verificationCode, nil
}

func (as *Service) consumeVerificationCode(ctx context.Context, verificationID string, code string) (*VerificationCode, error) {
	const op = "auth/service.consumeVerificationCode"

	verificationCode, err := as.store.GetVerificationCodeByIDAndCode(ctx, verificationID, code)

	if err != nil {
		return nil, errors.Unexpected(op, err, "failed to get verification code")
	}

	if verificationCode == nil {
		return nil, errors.Invalid(op, "verification code invalid")
	}

	if verificationCode.ExpiredAt.Before(time.Now()) {
		return nil, errors.Invalid(op, "verification code expired")
	}

	err = as.store.DeleteVerificationCodeByID(ctx, verificationCode.ID)

	if err != nil {
		return nil, errors.Unexpected(op, err, "failed to remove verification code")
	}

	return verificationCode, nil
}

// LoginVerify login verification initialization core logic
func (as *Service) LoginVerify(ctx context.Context, phoneNumber string, countryCode string) (string, error) {
	const op = "auth/service.LoginVerify"

	formattedPhoneNumber, err := phone.FormatPhoneNumber(phoneNumber, countryCode)

	if err != nil {
		return "", errors.Invalid(op, "invalid phone number")
	}

	user, err := as.store.GetUserByPhoneNumber(ctx, formattedPhoneNumber, countryCode)

	if err != nil {
		return "", errors.Unexpected(op, err, "failed to get user by phone number")
	}

	// Create and persist new user if it didn't exist yet
	if user == nil {
		user = NewUser(formattedPhoneNumber, countryCode)

		err = as.store.StoreUser(ctx, user)

		if err != nil {
			return "", errors.Unexpected(op, err, "failed to store user")
		}
	}

	verificationCode, err := as.createNewVerificationCode(ctx, user, formattedPhoneNumber, countryCode, "LOGIN")

	if err != nil {
		return "", errors.Unexpected(op, err, "failed to create new verification code")
	}

	err = as.smsSender.SendSMS(formattedPhoneNumber, verificationCode.CountryCode, verificationCode.Code)

	if err != nil {
		return "", errors.Unexpected(op, err, "failed to send sms")
	}

	return verificationCode.VerificationID, nil
}

// LoginCheck returns accessToken given the verificationID and code match the persisted verification code
func (as *Service) LoginCheck(ctx context.Context, verificationID string, code string) (string, error) {
	const op = "auth/service.LoginCheck"

	verificationCode, err := as.consumeVerificationCode(ctx, verificationID, code)

	if err != nil {
		return "", errors.Wrap(op, err, "failed to consume verification code")
	}

	user, err := as.store.GetUserByID(ctx, verificationCode.UserID)

	if err != nil {
		return "", errors.Unexpected(op, err, "failed to get user by id")
	}

	if user == nil {
		return "", errors.NotFound(op)
	}

	user.IsPhoneNumberVerified = true
	user.UpdatedAt = time.Now()

	err = as.store.UpdateUser(ctx, user)

	if err != nil {
		return "", errors.Unexpected(op, err, "failed to update user")
	}

	_, accessToken, err := as.tokenAuth.Encode(jwt.MapClaims{"user_id": user.ID})

	if err != nil {
		return "", errors.Unexpected(op, err, "failed to encode access token")
	}

	return accessToken, nil
}

// GetCurrentUser ...
func (as *Service) GetCurrentUser(ctx context.Context) (*User, error) {
	const op = "auth/service.GetCurrentUser"
	_, claims, err := jwtauth.FromContext(ctx)

	if err != nil {
		return nil, errors.Wrap(op, err, "missing or invalid access token")
	}

	userID := fmt.Sprintf("%v", claims["user_id"])
	user, err := as.store.GetUserByID(ctx, userID)

	if err != nil {
		return nil, errors.Unexpected(op, err, "failed to get user by id")
	}

	return user, nil
}

// UpdatePhoneNumberVerify ...
func (as *Service) UpdatePhoneNumberVerify(ctx context.Context, phoneNumber string, countryCode string, currentUser *User) (string, error) {
	const op = "auth/service.UpdatePhoneNumberVerify"

	formattedPhoneNumber, err := phone.FormatPhoneNumber(phoneNumber, countryCode)

	if err != nil {
		return "", errors.Wrap(op, err, "invalid phone number")
	}

	user, err := as.store.GetUserByPhoneNumber(ctx, formattedPhoneNumber, countryCode)

	if err != nil {
		return "", errors.Unexpected(op, err, "failed to get user by phone number")
	}

	if user != nil {
		return "", errors.Invalid(op, "user with given phone number already exists")
	}

	verificationCode, err := as.createNewVerificationCode(ctx, currentUser, formattedPhoneNumber, countryCode, "UPDATE")

	if err != nil {
		return "", errors.Unexpected(op, err, "failed to create new verification code")
	}

	err = as.smsSender.SendSMS(formattedPhoneNumber, verificationCode.CountryCode, verificationCode.Code)

	if err != nil {
		return "", errors.Unexpected(op, err, "failed to send sms")
	}

	return verificationCode.VerificationID, nil
}

// UpdatePhoneNumberCheck ...
func (as *Service) UpdatePhoneNumberCheck(ctx context.Context, verificationID string, code string, currentUser *User) (*User, error) {
	const op = "auth/service.UpdatePhoneNumberCheck"

	verificationCode, err := as.consumeVerificationCode(ctx, verificationID, code)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to consume verification code")
	}

	user, err := as.store.GetUserByID(ctx, verificationCode.UserID)

	if err != nil {
		return nil, errors.Unexpected(op, err, "failed to get user by id")
	}

	if user == nil {
		return nil, errors.NotFound(op)
	}

	if currentUser.ID != user.ID {
		return nil, errors.Unauthorized(op, fmt.Errorf("current user not owner"))
	}

	user.UpdatedAt = time.Now()
	user.PhoneNumber = verificationCode.PhoneNumber
	user.CountryCode = verificationCode.CountryCode

	err = as.store.UpdateUser(ctx, user)

	if err != nil {
		return nil, errors.Unexpected(op, err, "failed to update user")
	}

	return user, nil
}

// UpdateUserProfileInput ...
type UpdateUserProfileInput struct {
	FullName       string `json:"full_name"`
	ProfileImageID string `json:"image_id"`
}

// UpdateUserProfile ...
func (as *Service) UpdateUserProfile(ctx context.Context, input *UpdateUserProfileInput, currentUser *User) (*User, error) {
	const op = "auth/service.UpdatePhoneNumberCheck"

	user, err := as.store.GetUserByID(ctx, currentUser.ID)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to get user by id")
	}

	user.UpdatedAt = time.Now()

	if input.FullName != "" {
		user.FullName = strings.TrimSpace(input.FullName)
	}
	if input.ProfileImageID != "" {
		user.ProfileImageID = input.ProfileImageID
	}

	err = as.store.UpdateUser(ctx, user)

	if err != nil {
		return nil, errors.Unexpected(op, err, "failed to update user")
	}

	return user, nil
}
