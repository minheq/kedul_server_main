package handlers

import (
	"testing"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/google/uuid"
	"github.com/minheq/kedul_server_main/errors"
	"github.com/minheq/kedul_server_main/models"
	"github.com/minheq/kedul_server_main/testutils"
)

func TestLoginHappyPath(t *testing.T) {
	var code string
	var verificationID string
	var err error
	db, cleanup := testutils.SetupDB()
	store := models.NewStore(db)
	smsSender := &testutils.SmsSenderMock{}
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	defer cleanup()

	t.Run("should send code and return state when login start", func(t *testing.T) {
		verificationID, err = LoginVerify("999111333", "VN", store, smsSender)
		code = smsSender.Text

		if err != nil {
			t.Error(err)
		}

		if code == "" {
			t.Error("missing code")
		}
	})

	t.Run("should return access token when login verified", func(t *testing.T) {
		accessToken, err := LoginVerifyCheck(verificationID, code, store, tokenAuth)

		if err != nil {
			t.Error(err)
		}

		if accessToken == "" {
			t.Error("missing access token")
		}
	})
}

func TestLoginWithExpiredVerificationCode(t *testing.T) {
	db, cleanup := testutils.SetupDB()
	store := models.NewStore(db)
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	defer cleanup()

	now := time.Now()

	account := models.NewAccount("999999999", "VN")

	err := store.StoreAccount(account)

	if err != nil {
		t.Error(err)
		return
	}

	expiredVerificationCode := &models.VerificationCode{
		ID:             uuid.Must(uuid.New(), nil).String(),
		AccountID:      account.ID,
		Code:           "111111",
		CodeType:       "LOGIN",
		VerificationID: "ABC",
		PhoneNumber:    account.PhoneNumber,
		CountryCode:    account.CountryCode,
		ExpiredAt:      now.Add(time.Duration(-1) * time.Minute),
		CreatedAt:      now,
	}

	err = store.StoreVerificationCode(expiredVerificationCode)

	if err != nil {
		t.Error(err)
		return
	}

	t.Run("should return error when log in verify with expired verification code", func(t *testing.T) {
		_, err := LoginVerifyCheck(expiredVerificationCode.VerificationID, expiredVerificationCode.Code, store, tokenAuth)

		if !errors.Is(errors.KindInvalid, err) {
			t.Error("should forbid")
		}
	})
}
