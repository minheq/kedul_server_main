package auth

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/google/uuid"
	"github.com/minheq/kedul_server_main/errors"
	"github.com/minheq/kedul_server_main/testutils"
)

var (
	store       Store
	authService Service
	smsSender   *testutils.SmsSenderMock
)

func TestMain(m *testing.M) {
	db, cleanup := testutils.SetupDB()

	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	store = NewStore(db)
	smsSender = &testutils.SmsSenderMock{}

	authService = NewService(store, tokenAuth, smsSender)

	code := m.Run()

	cleanup()

	os.Exit(code)
}

func TestLoginHappyPath(t *testing.T) {
	var code string
	var verificationID string
	var err error

	t.Run("should send code and return verificationID when login start", func(t *testing.T) {
		verificationID, err = authService.LoginVerify(context.Background(), "999111333", "VN")
		code = smsSender.Text

		if err != nil {
			t.Error(err)
			return
		}

		if code == "" {
			t.Error("missing code")
			return
		}
	})

	t.Run("should return access token when login verified", func(t *testing.T) {
		accessToken, err := authService.LoginVerifyCheck(context.Background(), verificationID, code)

		if err != nil {
			t.Error(err)
			return
		}

		if accessToken == "" {
			t.Error("missing access token")
		}
	})
}

func TestLoginWithExpiredVerificationCode(t *testing.T) {
	now := time.Now()

	user := NewUser("999999999", "VN")

	err := store.StoreUser(context.Background(), user)

	if err != nil {
		t.Error(err)
		return
	}

	expiredVerificationCode := &VerificationCode{
		ID:             uuid.Must(uuid.New(), nil).String(),
		UserID:         user.ID,
		Code:           "111111",
		CodeType:       "LOGIN",
		VerificationID: "ABC",
		PhoneNumber:    user.PhoneNumber,
		CountryCode:    user.CountryCode,
		ExpiredAt:      now.Add(time.Duration(-1) * time.Minute),
		CreatedAt:      now,
	}

	err = store.StoreVerificationCode(context.Background(), expiredVerificationCode)

	if err != nil {
		t.Error(err)
		return
	}

	t.Run("should return error when log in verify with expired verification code", func(t *testing.T) {
		_, err := authService.LoginVerifyCheck(context.Background(), expiredVerificationCode.VerificationID, expiredVerificationCode.Code)

		if !errors.Is(errors.KindInvalid, err) {
			t.Error("should forbid")
		}
	})
}

func TestLoginVerifyTwice(t *testing.T) {
	var codeOne string
	var verificationIDOne string
	var codeTwo string
	var verificationIDTwo string
	var err error

	t.Run("should send code and return verificationID when login start first time", func(t *testing.T) {
		verificationIDOne, err = authService.LoginVerify(context.Background(), "999111333", "VN")
		codeOne = smsSender.Text

		if err != nil {
			t.Error(err)
			return
		}

		if codeOne == "" {
			t.Error("missing code")
		}
	})

	// This behaves like "resending"
	t.Run("should send different code and verificationID when login start second time", func(t *testing.T) {
		verificationIDTwo, err = authService.LoginVerify(context.Background(), "999111333", "VN")
		codeTwo = smsSender.Text

		if err != nil {
			t.Error(err)
			return
		}

		if codeOne == codeTwo {
			t.Error("same code")
			return
		}

		if verificationIDOne == verificationIDTwo {
			t.Error("same verification id")
		}
	})
}
